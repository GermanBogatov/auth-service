package tracer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHostIPNamePort(t *testing.T) {
	tests := []struct {
		name     string
		host     string
		wantIp   string
		wantName string
		wantPort int
	}{
		{
			name:     "test-1",
			host:     "127.0.0.1:8080",
			wantIp:   "127.0.0.1",
			wantName: "",
			wantPort: 8080,
		},
		{
			name:     "test-2",
			host:     "localhost:8080",
			wantIp:   "",
			wantName: "localhost",
			wantPort: 8080,
		},
		{
			name:     "test-2",
			host:     "user-tag.com:8080",
			wantIp:   "",
			wantName: "user-tag.com",
			wantPort: 8080,
		},
	}

	for _, tt := range tests {
		ip, name, port := hostIPNamePort(tt.host)
		assert.Equal(t, tt.wantIp, ip)
		assert.Equal(t, tt.wantName, name)
		assert.Equal(t, tt.wantPort, port)
	}
}
