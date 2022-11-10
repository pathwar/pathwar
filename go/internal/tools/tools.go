//go:build tools
// +build tools

// Package tools ensures that `go mod` can detect some required dependencies.
// This package should not be imported directly.
package tools

import (
	_ "github.com/githubnemo/CompileDaemon"                               // required by Makefile
	_ "github.com/gobuffalo/packr/v2/packr2/cmd"                          // required by Makefile
	_ "github.com/gogo/protobuf/gogoproto"                                // required by protoc
	_ "github.com/gogo/protobuf/types"                                    // required by protoc
	_ "github.com/golang/protobuf/proto"                                  // required by protoc
	_ "github.com/golang/protobuf/ptypes/timestamp"                       // required by protoc
	_ "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger/options" // required by protoc
	_ "github.com/tailscale/depaware"                                     // required by Makefile
)
