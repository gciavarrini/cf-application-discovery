package discover

import (
	"cf-application-discovery/pkg/models"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Health Checks", func() {
	ginkgo.Context("parseHealthCheck", func() {
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
			ginkgo.It(tt.name, func() {
				result := parseHealthCheck(tt.cfType, tt.cfEndpoint, tt.cfInterval, tt.cfTimeout)

				// Use Gomega's Expect function for assertions
				gomega.Expect(result).To(gomega.Equal(tt.expected))
			})
		}
	})

	ginkgo.Context("parseReadinessHealthCheck", func() {
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
			ginkgo.It(tt.name, func() {
				result := parseReadinessHealthCheck(tt.cfType, tt.cfEndpoint, tt.cfInterval, tt.cfTimeout)

				gomega.Expect(result).To(gomega.Equal(tt.expected))
			})
		}
	})
})

var _ = ginkgo.Describe("Parse Process", func() {
	tests := []struct {
		name      string
		cfProcess AppManifestProcess
		expected  models.ProcessSpec
	}{
		{
			name: "default values",
			cfProcess: AppManifestProcess{
				Type: AppProcessType(models.Web),
			},
			expected: models.ProcessSpec{
				Type:      models.Web,
				Command:   "",
				DiskQuota: "",
				Memory:    "1G",
				HealthCheck: models.ProbeSpec{
					Type:     models.PortProbeType,
					Endpoint: "/",
					Timeout:  1,
					Interval: 30,
				},
				ReadinessCheck: models.ProbeSpec{
					Type:     models.ProcessProbeType,
					Endpoint: "/",
					Timeout:  1,
					Interval: 30,
				},
				Instances:    1,
				LogRateLimit: "16K",
				Lifecycle:    "",
			},
		},
		{
			name: "custom memory and instances",
			cfProcess: AppManifestProcess{
				Type:                         "worker",
				Command:                      "run_worker",
				DiskQuota:                    "256M",
				Memory:                       "512M",
				Instances:                    uintPtr(3),
				LogRateLimitPerSecond:        "32K",
				HealthCheckType:              "http",
				HealthCheckHTTPEndpoint:      "/health",
				HealthCheckInterval:          10,
				HealthCheckInvocationTimeout: 5,
			},
			expected: models.ProcessSpec{
				Type:      models.ProcessType("worker"),
				Command:   "run_worker",
				DiskQuota: "256M",
				Memory:    "512M",
				HealthCheck: models.ProbeSpec{
					Type:     models.ProbeType("http"),
					Endpoint: "/health",
					Timeout:  5,
					Interval: 10,
				},
				ReadinessCheck: models.ProbeSpec{
					Type:     models.ProbeType("http"),
					Endpoint: "/health",
					Timeout:  5,
					Interval: 10,
				},
				Instances:    3,
				LogRateLimit: "32K",
				Lifecycle:    models.LifecycleType(""),
			},
		},
		{
			name: "custom log rate limit and lifecycle",
			cfProcess: AppManifestProcess{
				Type:                  "worker",
				Command:               "run_worker",
				DiskQuota:             "512M",
				LogRateLimitPerSecond: "64K",
				Lifecycle:             "cnb",
			},
			expected: models.ProcessSpec{
				Type:      models.ProcessType("worker"),
				Command:   "run_worker",
				DiskQuota: "512M",
				Memory:    "1G",
				HealthCheck: models.ProbeSpec{
					Type:     models.PortProbeType,
					Endpoint: "/",
					Timeout:  1,
					Interval: 30,
				},
				ReadinessCheck: models.ProbeSpec{
					Type:     models.ProcessProbeType,
					Endpoint: "/",
					Timeout:  1,
					Interval: 30,
				},
				Instances:    1,
				LogRateLimit: "64K",
				Lifecycle:    models.LifecycleType("cnb"),
			},
		},
		{
			name: "no health check specified",
			cfProcess: AppManifestProcess{
				Type:      "worker",
				Command:   "run_worker",
				DiskQuota: "256M",
				Instances: uintPtr(2),
			},
			expected: models.ProcessSpec{
				Type:      models.ProcessType("worker"),
				Command:   "run_worker",
				DiskQuota: "256M",
				Memory:    "1G",
				HealthCheck: models.ProbeSpec{
					Type:     models.PortProbeType,
					Endpoint: "/",
					Timeout:  1,
					Interval: 30,
				},
				ReadinessCheck: models.ProbeSpec{
					Type:     models.ProcessProbeType,
					Endpoint: "/",
					Timeout:  1,
					Interval: 30,
				},
				Instances:    2,
				LogRateLimit: "16K",
				Lifecycle:    models.LifecycleType(""),
			},
		},
	}

	for _, tt := range tests {
		ginkgo.It(tt.name, func() {
			result := parseProcess(tt.cfProcess)

			gomega.Expect(result).To(gomega.Equal(tt.expected))
		})
	}
})

// Helper function to create a pointer to a uint value.
func uintPtr(i uint) *uint {
	return &i
}
