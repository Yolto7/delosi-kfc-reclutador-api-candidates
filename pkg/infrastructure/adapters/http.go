package adapters

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

// HTTPClientConfig configuración básica para el cliente HTTP
type HTTPClientConfig struct {
	Timeout         time.Duration
	MaxIdleConns    int
	MaxConnsPerHost int
	IdleConnTimeout time.Duration
	KeepAlive       time.Duration
	DisableKeepAlives bool
}

// NewHTTPClient crea un cliente HTTP optimizado para producción
func NewHTTPClient(config HTTPClientConfig) (*http.Client, error) {
	// Aplicar defaults si no se especifican
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}
	if config.MaxIdleConns == 0 {
		config.MaxIdleConns = 100
	}
	if config.MaxConnsPerHost == 0 {
		config.MaxConnsPerHost = 10
	}
	if config.IdleConnTimeout == 0 {
		config.IdleConnTimeout = 90 * time.Second
	}
	if config.KeepAlive == 0 {
		config.KeepAlive = 30 * time.Second
	}

	// Crear transport optimizado para performance
	transport := &http.Transport{
		// Connection pooling optimizado
		MaxIdleConns:        config.MaxIdleConns,
		MaxIdleConnsPerHost: config.MaxConnsPerHost,
		MaxConnsPerHost:     config.MaxConnsPerHost,
		IdleConnTimeout:     config.IdleConnTimeout,

		// Timeouts optimizados
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,  // DNS + conexión
			KeepAlive: config.KeepAlive,
			DualStack: true, // IPv4 + IPv6
		}).DialContext,

		// TLS optimizado
		TLSHandshakeTimeout: 10 * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false,
		},

		// HTTP/2 support
		ForceAttemptHTTP2: true,

		// Compression
		DisableCompression: false,

		// Keep-alive
		DisableKeepAlives: false,
	}

	// Crear cliente HTTP
	client := &http.Client{
		Transport: transport,
		Timeout:   config.Timeout,
		
		// Manejar redirects (máximo 10)
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return http.ErrUseLastResponse
			}
			return nil
		},
	}

	return client, nil
}

// NewDefaultHTTPClient crea un cliente con configuración por defecto optimizada
func NewDefaultHTTPClient() (*http.Client, error) {
	return NewHTTPClient(HTTPClientConfig{
		Timeout:         30 * time.Second,
		MaxIdleConns:    100,
		MaxConnsPerHost: 10,
		IdleConnTimeout: 90 * time.Second,
		KeepAlive:       30 * time.Second,
		DisableKeepAlives: false,
	})
}

// Para casos específicos como ElevenLabs que puede necesitar timeouts más largos
func NewLongTimeoutHTTPClient() (*http.Client, error) {
	return NewHTTPClient(HTTPClientConfig{
		Timeout:         60 * time.Second,
		MaxIdleConns:    50,
		MaxConnsPerHost: 5,
		IdleConnTimeout: 90 * time.Second,
		KeepAlive:       30 * time.Second,
	})
}