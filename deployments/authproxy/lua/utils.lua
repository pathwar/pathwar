local conf = require('auth_conf')
local http = require('socket.http')
local json = require('cjson.safe')
local cache = require('shcache')

local M = {}
local DEBUG=false

local function api_request(endpoint, where, embedded)
   local base_url = conf.api_scheme .. conf.api_user .. ':' .. conf.api_password .. '@' .. conf.api_host
   local request = base_url .. endpoint
   if where then
      request = request .. '?where={'
      for k,v in pairs(where) do
	 request = request .. '"' .. k .. '":"' .. v .. '",'
      end
      request = string.sub(request, 1, -2)
      request = request .. '}'
   end
   if embedded then
      if where then
	 request = request .. '&embedded={'
      else
	 request = request .. '?embedded={'
      end
      for _,v in ipairs(embedded) do
	 request = request .. '"' .. v .. '":1,'
      end
      request = string.sub(request, 1, -2)
      request = request .. '}'
   end
   if DEBUG then
      ngx.say(request)
   end
   r,c,q = http.request(request)
   if c ~= 200 then
      ngx.log(ngx.ERR, string.format("Bad response from API (%s) for request : %s",c, request))
      ngx.exit(ngx.HTTP_INTERNAL_SERVER_ERROR)
   end
   return r,c,q
end

function string:split(delimiter)
   local result = { }
   local from  = 1
   local delim_from, delim_to = string.find( self, delimiter, from  )
   while delim_from do
      table.insert( result, string.sub( self, from , delim_from-1 ) )
      from  = delim_to + 1
      delim_from, delim_to = string.find( self, delimiter, from  )
   end
   table.insert( result, string.sub( self, from  ) )
   return result
end

function table:length()
   local count = 0
   for _ in pairs(self) do count = count + 1 end
   return count
end

function tprint (tbl, indent)
  if not indent then indent = 0 end
  for k, v in pairs(tbl) do
     formatting = string.rep("  ", indent) .. k .. ": "
     if type(v) == "table" then
	ngx.say(formatting)
	tprint(v, indent+1)
     elseif type(v) == 'boolean' then
	ngx.say(formatting .. tostring(v))
    else
       ngx.say(formatting .. v)
    end
  end
end
M.tprint = tprint

local function get_level_id()
   return ngx.var.level_id
end
M.get_level_id = get_level_id

local function parse_auth()
   local auth_header = ngx.req.get_headers()['Authorization']
   local infos = string.split(auth_header, ' ')
   if table.length(infos) ~= 2 then
      return nil
   end
   if infos[1] ~= 'Basic' then
      return nil
   end
   local decoded = ngx.decode_base64(infos[2])
   if decoded == nil then
      return nil
   end
   login_pass = string.split(decoded, ':')
   if table.length(login_pass) ~= 2 then
      return nil
   end
   return login_pass
end
M.parse_auth = parse_auth

local function need_auth()
   ngx.header['WWW-Authenticate'] = 'Basic Realm="Pathwar Authentication"'
   ngx.exit(ngx.HTTP_UNAUTHORIZED)
end
M.need_auth = need_auth

local function send_cookie(login)
   cookie_data = login .. '-' .. get_level_id() .. '-' .. os.time()+conf.cookie_duration
   ngx.header['Set-Cookie'] = 'pathwar: ' .. cookie_data .. '|' .. 
      ngx.encode_base64(ngx.hmac_sha1(conf.hmac_key, cookie_data))
end
M.send_cookie = send_cookie

local function __user_has_access(login, pass)
   local r,c,q = api_request('/users', {login=login})
   local data = json.decode(r)
   if data['_meta']['total'] == 0 then
      ngx.log(ngx.ERR, string.format('No matching user found for login %s', login))
      return nil
   end
   local user_id = ''
   for index,user_infos in pairs(data["_items"]) do
      if user_infos['login'] == login then
	 user_id = user_infos['_id']
      end
   end
   r,c,q = api_request('/raw-organization-users',{user=user_id})   
   local data = json.decode(r)
   local level_id = get_level_id()
   if data['_meta']['total'] == 0 then
      ngx.log(ngx.ERR, string.format('No org for user %s (%s) ??', login, user_id))
      return nil
   end
   local ok = false
   for _,org_user in pairs(data['_items']) do
      local org_id = org_user['organization']
      r,c,q = api_request('/raw-level-instance-users',{user=user_id,
				        	   organization=org_id,
						   hash=pass,
						   level=level_id})
      data = json.decode(r)
      if data['_meta']['total'] == 0 then
	 ngx.log(ngx.ERR, string.format('Could not find a level %s for user %s with hash %s in org %s', 
					level_id, 
					login, 
					pass, 
					org_id
				       )
	 )
      else
        ok = true
      end
   end
   if ok == true then
      return 1
   end
   return nil
end

local function user_has_access(login, pass)
   -- TODO: Get data from cache
   -- TODO: Escape login

   local lookup_func = function ()
      return __user_has_access(login, pass)
   end

   local cache_table, err = cache:new(ngx.shared.pathwar_auth,
				      {external_lookup=lookup_func},
				      {positive_ttl = conf.cache_duration_valid,
				       negative_ttl = conf.cache_duration_invalid})
   if cache_table == nil then
      ngx.log(ngx.ERR, string.format('Error while accessing cache: %s', err))
      return false
   end
   local cache_data, from_cache = cache_table:load(login..'||'..pass..'||'..get_level_id())
   if cache_data then
      return true
   else
      return false
   end        
end
M.user_has_access = user_has_access

return M