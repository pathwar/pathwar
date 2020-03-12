package pwapi

import (
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

var (
	// trace header to propagate.
	traceHeaders = []string{
		"x-ot-span-context",
		"x-request-id",

		// Zipkin headers
		"b3",
		"x-b3-traceid",
		"x-b3-spanid",
		"x-b3-parentspanid",
		"x-b3-sampled",
		"X-b3-flags",

		// Jaeger header (for native client)
		"uber-trace-id",
	}
)

func incomingHeaderMatcherFunc(key string) (string, bool) {
	k := strings.ToLower(key)
	for _, v := range traceHeaders {
		if v == k {
			return key, true
		}
	}
	return runtime.DefaultHeaderMatcher(key)
}
