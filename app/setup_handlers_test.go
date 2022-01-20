package app

import (
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestSetupHandlers_CreateName(t *testing.T) {
	const releaseVersion = "0.10.0"
	const releaseCandidate4 = "rc.4"
	const releaseCandidate5 = "rc.5"
	const releaseCandidate6 = "rc.6"
	require.Equal(t, "0.10.0-rc.4", strings.Join([]string{releaseVersion, releaseCandidate4}, "-"))
	require.Equal(t, "0.10.0-rc.5", strings.Join([]string{releaseVersion, releaseCandidate5}, "-"))
	require.Equal(t, "0.10.0-rc.6", strings.Join([]string{releaseVersion, releaseCandidate6}, "-"))
}
