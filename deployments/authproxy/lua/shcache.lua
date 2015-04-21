-- Copyright (C) 2013 Matthieu Tourne
-- @author Matthieu Tourne <matthieu@cloudflare.com>

-- small overlay over shdict, smart cache load mechanism

local M = {}

local resty_lock = require("lock")

local DEBUG = false

-- defaults in secs
local DEFAULT_POSITIVE_TTL = 10     -- cache for, successful lookup
local DEFAULT_NEGATIVE_TTL = 2      -- cache for, failed lookup
local DEFAULT_ACTUALIZE_TTL = 2     -- stale data, actualize data for

-- default lock options, in secs
local function _get_default_lock_options()
   return {
      exptime = 1,     -- max wait if failing to call unlock()
      timeout = 0.5,   -- max waiting time of lock()
      max_step = 0.1,  -- max sleeping interval
   }
end

local conf = require("auth_conf")
if conf then
   DEFAULT_NEGATIVE_TTL = conf.DEFAULT_NEGATIVE_TTL or DEFAULT_NEGATIVE_TTL
   DEFAULT_ACTUALIZE_TTL = conf.DEFAULT_ACTUALIZE_TTL or DEFAULT_ACTUALIZE_TTL
end

local band = bit.band
local bor = bit.bor
local st_format = string.format

-- there are only really 5 states total
                                 -- is_stale    is_neg  is_from_cache
local MISS_STATE = 0             -- 0           0       0
local HIT_POSITIVE_STATE = 1     -- 0           0       1
local HIT_NEGATIVE_STATE = 3     -- 0           1       1
local STALE_POSITIVE_STATE = 5   -- 1           0       1

-- stale negative doesn't really make sense, use HIT_NEGATIVE instead
-- local STALE_NEGATIVE_STATE = 7   -- 1           1       1

-- xor to set
local NEGATIVE_FLAG = 2
local STALE_FLAG = 4

local STATES = {
   [MISS_STATE] = 'MISS',
   [HIT_POSITIVE_STATE] = 'HIT',
   [HIT_NEGATIVE_STATE] = 'HIT_NEGATIVE',
   [STALE_POSITIVE_STATE] = 'STALE',
   -- [STALE_NEGATIVE_STATE] = 'STALE_NEGATIVE',
}

local function get_status(flags)
   return STATES[flags] or st_format('UNDEF (0x%x)', flags)
end

local EMPTY_DATA = '_EMPTY_'

-- install debug functions
if DEBUG then
   local resty_lock_lock = resty_lock.lock

   resty_lock.lock = function (...)
      local _, key = unpack({...})
      print("lock key: ", tostring(key))
      return resty_lock_lock(...)
   end

   local resty_lock_unlock = resty_lock.unlock

   resty_lock.unlock = function (...)
      print("unlock")
      return resty_lock_unlock(...)
   end
end


-- store the object in the context
-- useful for debugging and tracking cache status
local function _store_object(self, name)
   if DEBUG then
      print('storing shcache: ', name, ' into ngx.ctx')
   end

   local ngx_ctx = ngx.ctx

   if not ngx_ctx.shcache then
      ngx_ctx.shcache = {}
   end
   ngx_ctx.shcache[name] = self
end

local obj_mt = {
   __index = M,
}

-- default function for callbacks.encode / decode.
local function _identity(data)
   return data
end

-- shdict: ngx.shared.DICT, created by the lua_shared_dict directive
-- callbacks: see shcache state machine for user defined functions
--    * callbacks.external_lookup is required
--    * callbacks.encode    : optional encoding before saving to shmem
--    * callbacks.decode    : optional decoding when retreiving from shmem
-- opts:
--   * opts.positive_ttl    : save a valid external loookup for, in seconds
--   * opts.positive_ttl    : save a invalid loookup for, in seconds
--   * opts.actualize_ttl   : re-actualize a stale record for, in seconds
--   * opts.lock_options    : set option to lock see : http://github.com/agentzh/lua-resty-lock
--                            for more details.
--   * opts.locks_shdict    : specificy the name of the shdict containing the locks
--                            (useful if you might have locks key collisions)
--                            uses "locks" by default.
--   * opts.name            : if shcache object is named, it will automatically
--                            register itself in ngx.ctx.shcache (useful for logging).
local function new(self, shdict, callbacks, opts)
   if not shdict then
      return nil, "shdict does not exist"
   end

   -- check that callbacks.external_lookup is set
   if not callbacks or not callbacks.external_lookup then
      return nil, "no external_lookup function defined"
   end

   if not callbacks.encode then
      callbacks.encode = _identity
   end

   if not callbacks.decode then
      callbacks.decode = _identity
   end

   local opts = opts or {}

   -- merge default lock options with the ones passed to new()
   local lock_options = _get_default_lock_options()
   if opts.lock_options then
      for k, v in pairs(opts.lock_options) do
         lock_options[k] = v
      end
   end

   local name = opts.name

   local obj = {
      shdict = shdict,
      callbacks = callbacks,

      positive_ttl = opts.positive_ttl or DEFAULT_POSITIVE_TTL,
      negative_ttl = opts.negative_ttl or DEFAULT_NEGATIVE_TTL,

      -- ttl to actualize stale data to
      actualize_ttl = opts.actualize_ttl or DEFAULT_ACTUALIZE_TTL,

      lock_options = lock_options,

      locks_shdict = opts.lock_shdict or "locks",

      -- STATUS --

      from_cache = false,
      cache_status = 'UNDEF',
      cache_state = MISS_STATE,
      lock_status = 'NO_LOCK',

      -- shdict:set() pushed out another value
      forcible_set = false,

      -- cache hit on second attempt (post lock)
      hit2 = false,

      name = name,
   }

   local locks = ngx.shared[obj.locks_shdict]

   -- check for existence, locks is not directly used
   if not locks then
      ngx.log(ngx.CRIT, 'shared mem locks is missing.\n',
              '## add to you lua conf: lua_shared_dict locks 5M; ##')
       return nil
   end

   local self = setmetatable(obj, obj_mt)

   -- if the shcache object is named
   -- keep track of the object in the context
   -- (useful for gathering stats at log phase)
   if name then
      _store_object(self, name)
   end

   return self
end
M.new = new

-- acquire a lock
local function _get_lock(self)
   local lock = self.lock
   if not lock then
      lock = resty_lock:new(self.locks_shdict, self.lock_options)
      self.lock = lock
   end
   return lock
end

-- remove the lock if there is any
local function _unlock(self)
   local lock = self.lock
   if lock then
      local ok, err = lock:unlock()
      if not ok then
         ngx.log(ngx.ERR, "failed to unlock :" , err)
      end
      self.lock = nil
   end
end

local function _return(self, data, flags)
   -- make sure we remove the locks if any before returning data
   _unlock(self)

   -- set cache status
   local cache_status = get_status(self.cache_state)

   if cache_status == 'MISS' and not data then
      cache_status = 'NO_DATA'
   end

   self.cache_status = cache_status

   return data, self.from_cache
end

local function _set(self, ...)
   if DEBUG then
      local key, data, ttl, flags = unpack({...})
      print("saving key: ", key, ", for: ", ttl)
   end

   local ok, err, forcible = self.shdict:set(...)

   self.forcible_set = forcible

   if not ok then
      local key, data, ttl, flags = unpack({...})
      ngx.log(ngx.ERR, 'failed to set key: ', key, ', err: ', err)
   end

   return ok
end

-- check if the data returned by :get() is considered empty
local function _is_empty(data, flags)
   return flags and band(flags, NEGATIVE_FLAG) and data == EMPTY_DATA
end

-- save positive, encode the data if needed before :set()
local function _save_positive(self, key, data)
   if DEBUG then
      print("key: ", key, ". save positive, ttl: ", self.positive_ttl)
   end
   data = self.callbacks.encode(data)
   return _set(self, key, data, self.positive_ttl, HIT_POSITIVE_STATE)
end

-- save negative, no encoding required (no data actually saved)
local function _save_negative(self, key)
   if DEBUG then
      print("key: ", key, ". save negative, ttl: ", self.negative_ttl)
   end
   return _set(self, key, EMPTY_DATA, self.negative_ttl, HIT_NEGATIVE_STATE)
end

-- save actualize, will boost a stale record to a live one
local function _save_actualize(self, key, data, flags)
   local new_flags = bor(flags, STALE_FLAG)

   if DEBUG then
      print("key: ", key, ". save actualize, ttl: ", self.actualize_ttl,
            ". new state: ", get_status(new_flags))
   end

   _set(self, key, data, self.actualize_ttl, new_flags)
   return new_flags
end

local function _process_cached_data(self, data, flags)
   if DEBUG then
      print("data: ", data, st_format(", flags: %x", flags))
   end

   self.cache_state = flags
   self.from_cache = true

   if _is_empty(data, flags) then
      -- empty cached data
      return nil
   else
      return self.callbacks.decode(data)
   end
end

-- wrapper to get data from the shdict
local function _get(self, key)
   -- always call get_stale() as it does not free element
   -- like get does on each call
   local data, flags, stale = self.shdict:get_stale(key)

   if data and stale then
      if DEBUG then
         print("found stale data for key : ", key)
      end

      self.stale_data = { data, flags }

      return nil, nil
   end

   return data, flags
end

local function _get_stale(self)
   local stale_data = self.stale_data
   if stale_data then
      return unpack(stale_data)
   end

   return nil, nil
end

local function load(self, key)
   -- start: check for existing cache
   local data, flags = _get(self, key)

   -- hit: process_cache_hit
   if data then
      data = _process_cached_data(self, data, flags)
      return _return(self, data)
   end

   -- miss: set lock

   -- lock: set a lock before performing external lookup
   local lock = _get_lock(self)
   local elapsed, err = lock:lock(key)

   if not elapsed then
      -- failed to acquire lock, still proceed normally to external_lookup
      -- unlock() might fail.
      ngx.log(ngx.ERR, "failed to acquire the lock: ", err)
      self.lock_status = 'ERROR'
      -- _unlock won't try to unlock() without a valid lock
      self.lock = nil
   else
      -- lock acquired successfuly

      if elapsed > 0 then

         -- elapsed > 0 => waited lock (other thread might have :set() the data)
         -- (more likely to get a HIT on cache_load 2)
         self.lock_status = 'WAITED'

      else

         -- elapsed == 0 => immediate lock
         -- it is less likely to get a HIT on cache_load 2
         -- but still perform it (race condition cases)
         self.lock_status = 'IMMEDIATE'
      end

      -- perform cache_load 2
      data, flags = _get(self, key)
      if data then
         -- hit2 : process cache hit

         self.hit2 = true

         -- unlock before de-serializing cached data
         _unlock(self)
         data = _process_cached_data(self, data, flags)
         return _return(self, data)
      end

      -- continue to external lookup
   end

   -- perform external lookup
   data, err = self.callbacks.external_lookup()

   if data then
      -- succ: save positive and return the data

      _save_positive(self, key, data)
      return _return(self, data)
   else
      ngx.log(ngx.WARN, 'external lookup failed: ', err)
   end

   -- external lookup failed
   -- attempt to load stale data
   data, flags = _get_stale(self)
   if data and not _is_empty(data, flags) then
      -- hit_stale + valid (positive) data

      flags = _save_actualize(self, key, data, flags)
      -- unlock before de-serializing data
      _unlock(self)
      data = _process_cached_data(self, data, flags)
      return _return(self, data)
   end

   if DEBUG and data then
      -- there is data, but it failed _is_empty() => stale negative data
      print('STALE_NEGATIVE data => cache as a new HIT_NEGATIVE')
   end

   -- nothing has worked, save negative and return empty
   _save_negative(self, key)
   return _return(self, nil)
end
M.load = load

return M
