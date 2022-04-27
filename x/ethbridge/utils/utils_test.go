package utils_test

import (
	"github.com/Sifchain/sifnode/x/ethbridge/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseStringToBool(t *testing.T) {
	tt := []struct {
		name     string
		input    string
		expected bool
		isError  bool
	}{
		{
			name:     "TC_1",
			input:    "true",
			expected: true,
			isError:  false,
		},
		{
			name:     "TC_2",
			input:    "true",
			expected: true,
			isError:  false,
		},
		{
			name:     "TC_3",
			input:    "false",
			expected: false,
			isError:  false,
		},
		{
			name:     "TC_4",
			input:    "True",
			expected: true,
			isError:  false,
		},
		{
			name:     "TC_5",
			input:    "False",
			expected: false,
			isError:  false,
		},
		{
			name:     "TC_6",
			input:    "No",
			expected: false,
			isError:  true,
		},
		{
			name:     "TC_7",
			input:    "Pause",
			expected: false,
			isError:  true,
		},
	}
	for _, test := range tt {
		tc := test
		t.Run(tc.name, func(t *testing.T) {
			toBool, err := utils.ParseStringToBool(tc.input)
			assert.Equal(t, tc.expected, toBool)
			if tc.isError {
				assert.Error(t, err)
			}
		})
	}
}
