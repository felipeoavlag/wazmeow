package whatsapp

import (
	"testing"

	"wazmeow/internal/domain/entity"
)

func TestClientFactory_buildProxyURL(t *testing.T) {
	cf := &ClientFactory{}

	tests := []struct {
		name        string
		proxyConfig *entity.ProxyConfig
		expected    string
	}{
		{
			name: "HTTP proxy sem autenticação",
			proxyConfig: &entity.ProxyConfig{
				Type: "http",
				Host: "proxy.example.com",
				Port: 8080,
			},
			expected: "http://proxy.example.com:8080",
		},
		{
			name: "HTTP proxy com autenticação",
			proxyConfig: &entity.ProxyConfig{
				Type:     "http",
				Host:     "proxy.example.com",
				Port:     8080,
				Username: "user",
				Password: "pass",
			},
			expected: "http://user:pass@proxy.example.com:8080",
		},
		{
			name: "SOCKS5 proxy sem autenticação",
			proxyConfig: &entity.ProxyConfig{
				Type: "socks5",
				Host: "socks.example.com",
				Port: 1080,
			},
			expected: "socks5://socks.example.com:1080",
		},
		{
			name: "SOCKS5 proxy com autenticação",
			proxyConfig: &entity.ProxyConfig{
				Type:     "socks5",
				Host:     "socks.example.com",
				Port:     1080,
				Username: "user",
				Password: "pass",
			},
			expected: "socks5://user:pass@socks.example.com:1080",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cf.buildProxyURL(tt.proxyConfig)
			if result != tt.expected {
				t.Errorf("buildProxyURL() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestClientFactory_configureProxy_validation(t *testing.T) {
	// Mock do domain service seria necessário para teste completo
	// Por enquanto, apenas teste básico de nil
	cf := &ClientFactory{}

	err := cf.configureProxy(nil)
	if err != nil {
		t.Errorf("configureProxy(nil) should not return error, got: %v", err)
	}
}
