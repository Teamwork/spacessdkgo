package client

import "testing"

func TestNormalizeBaseURL(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// empty input
		{
			name:  "empty string",
			input: "",
			want:  "",
		},

		// already correct
		{
			name:  "already has suffix, no trailing slash",
			input: "https://example.com/spaces/api/v1",
			want:  "https://example.com/spaces/api/v1",
		},

		// trailing slashes on correct URL
		{
			name:  "correct suffix with one trailing slash",
			input: "https://example.com/spaces/api/v1/",
			want:  "https://example.com/spaces/api/v1",
		},
		{
			name:  "correct suffix with many trailing slashes",
			input: "https://example.com/spaces/api/v1///",
			want:  "https://example.com/spaces/api/v1",
		},

		// bare host — suffix must be appended
		{
			name:  "bare host, no path",
			input: "https://example.com",
			want:  "https://example.com/spaces/api/v1",
		},
		{
			name:  "bare host with trailing slash",
			input: "https://example.com/",
			want:  "https://example.com/spaces/api/v1",
		},
		{
			name:  "bare host with multiple trailing slashes",
			input: "https://example.com///",
			want:  "https://example.com/spaces/api/v1",
		},

		// host with an unrelated path — suffix must be appended, no double slash
		{
			name:  "host with extra path segment",
			input: "https://example.com/v2",
			want:  "https://example.com/v2/spaces/api/v1",
		},
		{
			name:  "host with extra path and trailing slash",
			input: "https://example.com/v2/",
			want:  "https://example.com/v2/spaces/api/v1",
		},

		// suffix appears in the middle of the path
		{
			name:  "suffix already present but with extra path after it",
			input: "https://example.com/spaces/api/v1/extra",
			want:  "https://example.com/spaces/api/v1",
		},
		{
			name:  "suffix in middle with trailing slash",
			input: "https://example.com/spaces/api/v1/extra/",
			want:  "https://example.com/spaces/api/v1",
		},

		// no scheme — treated as an opaque string, suffix still appended correctly
		{
			name:  "relative-style base, no suffix",
			input: "example.com",
			want:  "example.com/spaces/api/v1",
		},
		{
			name:  "relative-style base with suffix",
			input: "example.com/spaces/api/v1",
			want:  "example.com/spaces/api/v1",
		},

		// case insensitivity — suffix matching must be case-insensitive,
		// canonical lowercase suffix is always written to the output
		{
			name:  "suffix all uppercase",
			input: "https://example.com/SPACES/API/V1",
			want:  "https://example.com/spaces/api/v1",
		},
		{
			name:  "suffix mixed case",
			input: "https://example.com/Spaces/Api/V1",
			want:  "https://example.com/spaces/api/v1",
		},
		{
			name:  "suffix mixed case with trailing slash",
			input: "https://example.com/Spaces/Api/V1/",
			want:  "https://example.com/spaces/api/v1",
		},
		{
			name:  "suffix mixed case mid-path",
			input: "https://example.com/Spaces/Api/V1/extra",
			want:  "https://example.com/spaces/api/v1",
		},
		{
			name:  "uppercase scheme and host are preserved",
			input: "HTTPS://EXAMPLE.COM/SPACES/API/V1",
			want:  "HTTPS://EXAMPLE.COM/spaces/api/v1",
		},

		// port in URL
		{
			name:  "host with port, no suffix",
			input: "http://localhost:8080",
			want:  "http://localhost:8080/spaces/api/v1",
		},
		{
			name:  "host with port and trailing slash",
			input: "http://localhost:8080/",
			want:  "http://localhost:8080/spaces/api/v1",
		},
		{
			name:  "host with port and correct suffix",
			input: "http://localhost:8080/spaces/api/v1",
			want:  "http://localhost:8080/spaces/api/v1",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := normalizeBaseURL(tc.input)
			if got != tc.want {
				t.Errorf("normalizeBaseURL(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}
