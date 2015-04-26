os = require('os')

local M = {}

M.hmac_key = 'XX_CHANGEME_XX'
M.cache_duration_valid = tonumber(os.getenv("PATHWAR_AUTH_VALID_DURATION")) or 60 * 60 * 60 -- 1 hour
M.cache_duration_invalid = 1 -- 1s (cannot be less)

M.api_scheme = os.getenv('PATHWAR_API_SCHEME') or 'http://'
M.api_user = os.getenv('PATHWAR_API_USER') or 'default'
M.api_password = os.getenv('PATHWAR_API_PASSWORD') or ''
M.api_host = os.getenv('PATHWAR_API_HOST') or 'localhost:5000'

return M