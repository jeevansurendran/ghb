package ghb

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_extractURLParams(t *testing.T) {
	tests := []struct {
		name     string
		pattern  string
		path     string
		expected map[string]string
		wantErr  bool
	}{
		{
			name:    "Single parameter",
			pattern: "v1/getUsers/{id}",
			path:    "v1/getUsers/123",
			expected: map[string]string{
				"id": "123",
			},
		},
		{
			name:    "Multiple parameters",
			pattern: "v1/getUsers/{id}/workspace/{wid}",
			path:    "v1/getUsers/123/workspace/456",
			expected: map[string]string{
				"id":  "123",
				"wid": "456",
			},
		},
		{
			name:     "No parameters",
			pattern:  "v1/getUsers",
			path:     "v1/getUsers",
			expected: map[string]string{},
		},
		{
			name:    "Special characters parameter",
			pattern: "v1/getUsers/{user-id}/posts/{post-id}",
			path:    "v1/getUsers/user-123/posts/1",
			expected: map[string]string{
				"user-id": "user-123",
				"post-id": "1",
			},
		},
		{
			name:    "Trailing slashes",
			pattern: "v1/getUsers/{id}",
			path:    "v1/getUsers/123/",
			expected: map[string]string{
				"id": "123",
			},
		},
		{
			name:    "Missing patterns",
			pattern: "v1/getUsers/{id}/workspace/{wid}",
			path:    "v1/getUsers/123/",
			expected: map[string]string{
				"id": "123",
			},
		},
		{
			name:     "Missing parts",
			pattern:  "v1/getUsers/{id}/posts/{pid}",
			path:     "v1/getUsers/123//posts/22",
			expected: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, _ := extractURLParams(tt.pattern, tt.path)
			require.Equal(t, actual, tt.expected)
		})
	}
}
