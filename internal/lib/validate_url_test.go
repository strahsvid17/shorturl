package lib

import "testing"

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name  string
		url   string
		valid bool
	}{
		{"https", "https://example.com", true},
		{"http", "http://example.com/path?key=value#section", true},
		{"localhost", "http://localhost:9090", true},
		{"IPv6", "http://[::1]:9090/path", true},
		{"escaped space", "https://example.com/a%20b", true},
		{"empty", "", false},
		{"missing scheme", "example.com", false},
		{"relative", "/path", false},
		{"unsupported scheme", "ftp://example.com", false},
		{"missing hostname", "https:///path", false},
		{"port without hostname", "http://:9090", false},
		{"invalid port", "http://example.com:abc", false},
		{"invalid escape", "https://example.com/%zz", false},
		{"space in path", "https://example.com/a b", false},
		{"leading space", " https://example.com", false},
		{"newline", "https://example.com/\n", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateURL(tt.url)
			if (err == nil) != tt.valid {
				t.Fatalf("ValidateURL(%q) = %v; want valid=%v", tt.url, err, tt.valid)
			}
		})
	}
}
