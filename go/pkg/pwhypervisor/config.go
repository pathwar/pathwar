package pwhypervisor

type HypervisorOpts struct {
	DomainSuffix      string             // .127.0.0.1.xip.io, .fr1.pathwar.pw, ...
	HostIP            string             // 0.0.0.0, ...
	HostPort          string             // 8000, 80 ...
	ModeratorPassword string             // s3cur3
	Salt              string             // s3cur3-t0o
	AllowedUsers      map[string][]int64 // map[INSTANCE_ID][]USER_ID, map[42][]string{4242, 4343}
	ForceRecreate     bool
	NginxDockerImage  string
}

type NginxConfigData struct {
	Upstreams []NginxUpstream
	Opts      HypervisorOpts
}

type NginxUpstream struct {
	Name   string
	Host   string
	Port   string
	Hashes []string
}

// https://github.com/jwilder/docker-gen/blob/master/templates/nginx.tmpl
// # the whole struct: {{. | toPrettyJson}}
const NginxConfigTemplate = `
	{{$root := .}}
	
	max_upload_size 80m;
	nice 502 page

	server {
		listen 80 default_server;
		server_name _; # This is just an invalid value which will never trigger on a real hostname.
		error_log /proc/self/fd/2;
		access_log /proc/self/fd/1;
		return 503;
	}

	{{range .Upstreams}}
	upstream { http://{{.Host}}:{{.Port}} {{.Name}} }
	{{end}}

	{{range .Upstreams}}
	server {
		listen {{.Port}};
		server_name {{range .Hashes}}{{.}}{{$root.Opts.DomainSuffix}}{{end}};
		location / {
			proxy_pass http://{{.Name}}
			# nice 404, 403
		}
	}
	{{end}}
`

/* const NginxConfigTemplate = `
  # the whole struct: {{. | toPrettyJson}}
	{{$root := .}}

	{{range .Upstreams}}
	upstream { http://CONTAINER:PORT NAME }
	{{end}}

	max_upload_size 80m;
	nice 502 page

	{{range .Upstreams}}
	server {
		listen PORT;
		server_name {{range .Users}}{{.Hash}}.{{$root.Opts.Suffix}}{{end}};
		location / {
			proxy_pass http://NAME
			# nice 404, 403
		}

	server {
		listen PORT;
		server_name moderator-{{.flavor}}.{{$root.Opts.Suffix}};
		location / {
			auth
			proxy_pass http://NAME
		}
	}

	#server {
		#lister PORT;
		#server_name status-{{.flavor}}.{{$root.Opts.Suffix}};
	#}
	{{end}}
`*/
