local utils = require('utils')

local auth_cookie = ngx.var.cookie_pathwar

if auth_cookie == nil or string.find(auth_cookie, "|") == nil then -- no cookie or invalid cookie
   -- check if we have a login and a password
   if ngx.req.get_headers()['Authorization'] == nil then -- no cookie and no login/password
      utils.need_auth()
   else -- parse the login/pass, make a request to the API, and set the cookie
      
      user = utils.parse_auth()
      if user == nil then
	 utils.need_auth()
      end
      if utils.user_has_access(user[1], user[2]) then
--	 utils.send_cookie(user[1])
	 return
      else
	 utils.need_auth()
      end
   end
-- elseif string.find(auth_cookie, '|') ~= nil then
--    local cookie_values = string.split(auth_cookie,'|')
--    local user_infos = cookie_values[1]
--    local hmac = cookie_values[2]
--    ngx.say(user_infos)
--    local computed_hmac = ngx.encode_base64(ngx.hmac_sha1(hmac_key, user_infos))
--    if hmac ~= computed_hmac then -- bad cookie, let's ask for login/password
--       ngx.log(ngx.ERR, string.format('Signature check failed for cookie %s (expected %s for HMAC) !', 
-- 				     cookie_values, computed_hmac))
--       utils.need_auth()
--    else -- good cookie
--       -- we need to check for timestamp validity, and if the user has still access to the current level
--       current_time = os.time()
--       cookies_infos = string.split(user_infos,'-')
--       if cookies_infos[3] + cookie_duration < current_time then --expired cookie, check if user has still access
-- 	 if utils.user_has_access(cookies_infos[1], nil) then
-- 	    utils.send_cookie(cookies_infos[1])
-- 	 else
-- 	    utils.need_auth()
-- 	 end
--       end
--       -- cookie is still valid, go ahead
--       return
--    end
end
