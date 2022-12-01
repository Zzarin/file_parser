package internal

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_stringToMap(t *testing.T) {
	tests := []struct {
		name          string
		giveNewString string
		wantMap       map[string]int
		wantErrText   string
	}{
		{
			name:          "correct_string",
			giveNewString: "2 543",
			wantMap: map[string]int{
				"key":   2,
				"value": 543,
			},
			wantErrText: "",
		},
		{
			name:          "key_not_an_integer",
			giveNewString: "sdf 543",
			wantMap:       nil,
			wantErrText:   "value string conversion to int",
		},
		{
			name:          "value_not_an_integer",
			giveNewString: "2 sdfsd",
			wantMap:       nil,
			wantErrText:   "value string conversion to int",
		},
		{
			name:          "key_empty",
			giveNewString: " 543",
			wantMap:       nil,
			wantErrText:   "value string conversion to int",
		},
		{
			name:          "value_empty",
			giveNewString: "4 ",
			wantMap:       nil,
			wantErrText:   "value string conversion to int",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			newString, err := stringToMap(test.giveNewString)

			require.Equal(t, test.wantMap, newString)
			if test.wantErrText != "" {
				require.ErrorContains(t, err, test.wantErrText)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
