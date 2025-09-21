package ghb

import (
	"testing"
	"time"

	"github.com/malayanand/ghb/test"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Test_marshalBytes(t *testing.T) {
	tests := []struct {
		name     string
		input    *test.TestUser
		options  *marshalOptions
		expected string
	}{
		{
			name: "basic fields",
			input: &test.TestUser{
				Id:   "abc",
				Name: "Alice",
				Age:  25,
			},
			expected: "{\"age\":25,\"id\":\"abc\",\"name\":\"Alice\"}",
		},
		{
			name: "with timestamp",
			input: &test.TestUser{
				Id:        "xyz",
				Name:      "Bob",
				Age:       40,
				CreatedAt: timestamppb.New(time.Unix(1758454323, 0)),
			},
			options: &marshalOptions{
				timeFormat: ISOTimeFormat,
			},
			expected: "{\"age\":40,\"created_at\":\"2025-09-21T11:32:03Z\",\"id\":\"xyz\",\"name\":\"Bob\"}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := tt.options
			if opts == nil {
				opts = &marshalOptions{}
			}
			b, err := marshalBytes(tt.input, opts)
			require.NoError(t, err)
			require.Equal(t, tt.expected, string(b))
		})
	}
}
