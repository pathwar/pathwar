server {

       listen _LISTEN_PORT_;
       server_name _SERVER_NAME_;

       set $level_id '_LEVEL_ID_';

       location / {
       		access_by_lua '/pathwar/lua/auth.lua';
		proxy_set_header Host $level_id;
		proxy_pass $scheme://_LEVEL_URL_;
       }
}