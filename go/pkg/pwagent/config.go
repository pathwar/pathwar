package pwagent

type AgentOpts struct {
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
	Opts      AgentOpts
}

type NginxUpstream struct {
	Name   string
	Host   string
	Port   string
	Hashes []string
}
