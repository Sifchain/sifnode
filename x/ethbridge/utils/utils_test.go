package utils_test

import (
	"github.com/Sifchain/sifnode/x/ethbridge/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseStringToBool(t *testing.T) {
	tc := []struct {
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
			name:     "TC_1",
			input:    "true",
			expected: true,
			isError:  false,
		},
		{
			name:     "TC_2",
			input:    "false",
			expected: false,
			isError:  false,
		},
		{
			name:     "TC_1",
			input:    "True",
			expected: true,
			isError:  false,
		},
		{
			name:     "TC_1",
			input:    "False",
			expected: false,
			isError:  false,
		},
		{
			name:     "TC_1",
			input:    "No",
			expected: false,
			isError:  true,
		},
		{
			name:     "TC_1",
			input:    "Pause",
			expected: false,
			isError:  true,
		},
	}
	for _, test := range tc {
		t.Run(test.name, func(t *testing.T) {
			toBool, err := utils.ParseStringToBool(test.input)
			assert.Equal(t, test.expected, toBool)
			if test.isError {
				assert.Error(t, err)
			}
		})
	}
}
