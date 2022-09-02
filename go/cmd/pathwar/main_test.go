package main

import (
	"testing"

	"github.com/stretchr/testify/require"
	"moul.io/u"
)

func Test(t *testing.T) {
	cleanup, err := u.CaptureStdout()
	require.NoError(t, err)

	err = runMain([]string{"version"})
	require.NoError(t, err)
	stdout := cleanup()
	require.Contains(t, stdout, "version=\"dev\"")
}
