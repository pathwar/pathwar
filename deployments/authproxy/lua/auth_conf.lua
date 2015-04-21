local M = {}

M.hmac_key = 'XX_CHANGEME_XX'
M.cache_duration_valid = 60 * 60 * 60 -- 1 hour
M.cache_duration_invalid = 1 -- 1s (cannot be less)

M.api_scheme = 'http://'
M.api_user = 'root-token'
M.api_password = ''
M.api_host = 'localhost:5000'

return M