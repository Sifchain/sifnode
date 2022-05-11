package keeper_test

import (
	"errors"
	"testing"

	clpkeeper "github.com/Sifchain/sifnode/x/clp/keeper"
	"github.com/stretchr/testify/require"
)

func TestKeeper_Int64ToUint8Safe(t *testing.T) {

	testcases := []struct {
		name      string
		x         int64
		expected  uint8
		errString error
	}{
		{
			name:     "success",
			x:        128,
			expected: 128,
		},
		{
			name:     "success 0",
			x:        0,
			expected: 0,
		},
		{
			name:     "success 255",
			x:        255,
			expected: 255,
		},
		{
			name:      "fail - below range",
			x:         -1,
			errString: errors.New("Could not perform type cast"),
		},
		{
			name:      "fail - above range",
			x:         256,
			errString: errors.New("Could not perform type cast"),
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {

			y, err := clpkeeper.Int64ToUint8Safe(tc.x)

			if tc.errString != nil {
				require.EqualError(t, err, tc.errString.Error())
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.expected, y)
		})
	}
}

func TestKeeper_Abs(t *testing.T) {

	testcases := []struct {
		name     string
		x        int16
		expected uint16
	}{
		{
			name:     "no change",
			x:        128,
			expected: 128,
		},
		{
			name:     "0 case",
			x:        0,
			expected: 0,
		},
		{
			name:     "flip sign",
			x:        -100,
			expected: 100,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {

			y := clpkeeper.Abs(tc.x)

			require.Equal(t, tc.expected, y)
		})
	}

}
