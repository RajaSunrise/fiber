package fiber

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBugRemoveEscapeChar(t *testing.T) {
	t.Parallel()
	// Fixed behavior:
	require.Equal(t, "abc", RemoveEscapeChar("\\abc"))
	require.Equal(t, ":", RemoveEscapeChar("\\:"))
	require.Equal(t, "\\", RemoveEscapeChar("\\\\")) // Should return \
}

func TestBugRouterDoubleHyphen(t *testing.T) {
	t.Parallel()
	rp := parseRoute("/api/:day-:month?-:year?")
	var params [maxParams]string
	match := rp.getMatch("/api/1-/-", "/api/1-/-", &params, false)

	require.False(t, match, "Should not match /api/1-/-")
}
