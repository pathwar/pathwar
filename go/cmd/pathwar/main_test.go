package main

import "os"

func Example() {
	os.Args = []string{"-h"}
	flagOutput = os.Stdout
	main()
	// Output:
	// USAGE
	//   pathwar [global flags] <subcommand> [flags] [args...]
	//
	// More info here: https://github.com/pathwar/pathwar/wiki/CLI
	//
	// SUBCOMMANDS
	//   rawclient  make API calls
	//   cli        CLI replacement for the web portal
	//   api        manage the Pathwar API
	//   compose    manage a challenge
	//   agent      manage an agent node (multiple challenges)
	//   misc       misc contains advanced commands
	//   admin      admin commands
	//   version    show version
	//
	// FLAGS
	//   -bearer-secretkey ...  bearer.sh secret key
	//   -debug false           debug mode
	//   -sentry-dsn ...        Sentry DSN
	//   -zipkin-endpoint ...   optional opentracing server
}
