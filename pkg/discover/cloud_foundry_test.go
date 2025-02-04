package discover

import (
	"cf-application-discovery/pkg/models"
	"testing"
)

func TestParseHealthCheck(t *testing.T) {
	tests := []struct {
		name       string
		cfType     AppHealthCheckType
		cfEndpoint string
		cfInterval uint
		cfTimeout  uint
		expected   models.ProbeSpec
	}{
		{
			name:       "default values",
			cfType:     "",
			cfEndpoint: "",
			cfInterval: 0,
			cfTimeout:  0,
			expected: models.ProbeSpec{
				Type:     models.PortProbeType,
				Endpoint: "/",
				Timeout:  1,
				Interval: 30,
			},
		},
		{
			name:       "custom type and endpoint",
			cfType:     "http",
			cfEndpoint: "/health",
			cfInterval: 10,
			cfTimeout:  5,
			expected: models.ProbeSpec{
				Type:     models.ProbeType("http"),
				Endpoint: "/health",
				Timeout:  5,
				Interval: 10,
			},
		},
		{
			name:       "custom timeout and interval",
			cfType:     "",
			cfEndpoint: "",
			cfInterval: 15,
			cfTimeout:  3,
			expected: models.ProbeSpec{
				Type:     models.PortProbeType,
				Endpoint: "/",
				Timeout:  3,
				Interval: 15,
			},
		},
		{
			name:       "only endpoint specified",
			cfType:     "",
			cfEndpoint: "/custom-endpoint",
			cfInterval: 0,
			cfTimeout:  0,
			expected: models.ProbeSpec{
				Type:     models.PortProbeType,
				Endpoint: "/custom-endpoint",
				Timeout:  1,
				Interval: 30,
			},
		},
		{
			name:       "only cf type specified",
			cfType:     "http",
			cfEndpoint: "",
			cfInterval: 0,
			cfTimeout:  0,
			expected: models.ProbeSpec{
				Type:     models.HTTPProbeType,
				Endpoint: "/",
				Timeout:  1,
				Interval: 30,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseHealthCheck(tt.cfType, tt.cfEndpoint, tt.cfInterval, tt.cfTimeout)

			if result != tt.expected {
				t.Errorf("expected %+v, got %+v", tt.expected, result)
			}
		})
	}
}

func TestParseReadinessHealthCheck(t *testing.T) {
	tests := []struct {
		name       string
		cfType     AppHealthCheckType
		cfEndpoint string
		cfInterval uint
		cfTimeout  uint
		expected   models.ProbeSpec
	}{
		{
			name:       "default values",
			cfType:     "",
			cfEndpoint: "",
			cfInterval: 0,
			cfTimeout:  0,
			expected: models.ProbeSpec{
				Type:     models.ProcessProbeType,
				Endpoint: "/",
				Timeout:  1,
				Interval: 30,
			},
		},
		{
			name:       "custom type and endpoint",
			cfType:     "http",
			cfEndpoint: "/ready",
			cfInterval: 10,
			cfTimeout:  5,
			expected: models.ProbeSpec{
				Type:     models.ProbeType("http"),
				Endpoint: "/ready",
				Timeout:  5,
				Interval: 10,
			},
		},
		{
			name:       "custom timeout and interval",
			cfType:     "",
			cfEndpoint: "",
			cfInterval: 15,
			cfTimeout:  3,
			expected: models.ProbeSpec{
				Type:     models.ProcessProbeType,
				Endpoint: "/",
				Timeout:  3,
				Interval: 15,
			},
		},
		{
			name:       "only endpoint specified",
			cfType:     "",
			cfEndpoint: "/health-check",
			cfInterval: 0,
			cfTimeout:  0,
			expected: models.ProbeSpec{
				Type:     models.ProcessProbeType,
				Endpoint: "/health-check",
				Timeout:  1,
				Interval: 30,
			},
		},
		{
			name:       "custom type with empty endpoint but valid timeout and interval",
			cfType:     "tcp",
			cfEndpoint: "",
			cfInterval: 5,
			cfTimeout:  3,
			expected: models.ProbeSpec{
				Type:     models.ProbeType("tcp"),
				Endpoint: "/",
				Timeout:  3,
				Interval: 5,
			},
		},
		{
			name:       "empty type with valid endpoint and custom interval/timeout",
			cfType:     "",
			cfEndpoint: "/status",
			cfInterval: 20,
			cfTimeout:  2,
			expected: models.ProbeSpec{
				Type:     models.ProcessProbeType,
				Endpoint: "/status",
				Timeout:  2,
				Interval: 20,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseReadinessHealthCheck(tt.cfType, tt.cfEndpoint, tt.cfInterval, tt.cfTimeout)

			if result != tt.expected {
				t.Errorf("expected %+v, got %+v", tt.expected, result)
			}
		})
	}
}
