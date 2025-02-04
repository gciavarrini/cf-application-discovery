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
