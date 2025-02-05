package discover

import (
	"cf-application-discovery/pkg/models"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Health Checks tests", func() {

	When("parsing health check probe", func() {
		DescribeTable("validate the correctness of the parsing logic", func(app AppManifestProcess, expected models.ProbeSpec) {
			result := parseHealthCheck(app.HealthCheckType, app.HealthCheckHTTPEndpoint, app.HealthCheckInterval, app.HealthCheckInvocationTimeout)
			// Use Gomega's Expect function for assertions
			Expect(result).To(Equal(expected))
		},
			Entry("with default values",
				AppManifestProcess{},
				models.ProbeSpec{
					Type:     models.PortProbeType,
					Endpoint: "/",
					Timeout:  1,
					Interval: 30,
				}),
			Entry("with process type and endpoint",
				AppManifestProcess{
					HealthCheckType:              "process",
					HealthCheckHTTPEndpoint:      "/health",
					HealthCheckInterval:          10,
					HealthCheckInvocationTimeout: 5,
				},
				models.ProbeSpec{
					Type:     models.ProcessProbeType,
					Endpoint: "/health",
					Timeout:  5,
					Interval: 10,
				}),
			Entry("with custom timeout and interval",
				AppManifestProcess{
					HealthCheckInterval:          15,
					HealthCheckInvocationTimeout: 3,
				},
				models.ProbeSpec{
					Type:     models.PortProbeType,
					Endpoint: "/",
					Timeout:  3,
					Interval: 15,
				}),
			Entry("with only endpoint specified",
				AppManifestProcess{
					HealthCheckHTTPEndpoint: "/custom-endpoint",
				},
				models.ProbeSpec{
					Type:     models.PortProbeType,
					Endpoint: "/custom-endpoint",
					Timeout:  1,
					Interval: 30,
				}),
			Entry("with only type specified",
				AppManifestProcess{
					HealthCheckType: "http",
				},
				models.ProbeSpec{
					Type:     models.HTTPProbeType,
					Endpoint: "/",
					Timeout:  1,
					Interval: 30,
				}),
		)
	})

	When("parsing readiness health check probe", func() {
		DescribeTable("validate the correctness of the parsing logic", func(app AppManifestProcess, expected models.ProbeSpec) {
			result := parseReadinessHealthCheck(app.ReadinessHealthCheckType, app.ReadinessHealthCheckHttpEndpoint, app.ReadinessHealthCheckInterval, app.ReadinessHealthInvocationTimeout)
			// Use Gomega's Expect function for assertions
			Expect(result).To(Equal(expected))
		},
			Entry("with default values",
				AppManifestProcess{},
				models.ProbeSpec{
					Type:     models.ProcessProbeType,
					Endpoint: "/",
					Timeout:  1,
					Interval: 30,
				}),
			Entry("with custom type and endpoint",
				AppManifestProcess{
					ReadinessHealthCheckType:         "http",
					ReadinessHealthCheckHttpEndpoint: "/ready",
					ReadinessHealthCheckInterval:     10,
					ReadinessHealthInvocationTimeout: 5,
				},
				models.ProbeSpec{
					Type:     models.HTTPProbeType,
					Endpoint: "/ready",
					Timeout:  5,
					Interval: 10,
				}),
			Entry("with custom timeout and interval",
				AppManifestProcess{
					ReadinessHealthCheckInterval:     15,
					ReadinessHealthInvocationTimeout: 3,
				},
				models.ProbeSpec{
					Type:     models.ProcessProbeType,
					Endpoint: "/",
					Interval: 15,
					Timeout:  3,
				}),
			Entry("with only endpoint specified",
				AppManifestProcess{
					ReadinessHealthCheckHttpEndpoint: "/readiness-check",
				},
				models.ProbeSpec{
					Type:     models.ProcessProbeType,
					Endpoint: "/readiness-check",
					Timeout:  1,
					Interval: 30,
				}),
			Entry("with custom type with empty endpoint but valid timeout and interval",
				AppManifestProcess{
					ReadinessHealthCheckType:         "port",
					ReadinessHealthCheckInterval:     5,
					ReadinessHealthInvocationTimeout: 3,
				},
				models.ProbeSpec{
					Type:     models.PortProbeType,
					Endpoint: "/",
					Timeout:  3,
					Interval: 5,
				}),
			Entry("with empty type with valid endpoint and custom interval/timeout",
				AppManifestProcess{
					ReadinessHealthCheckHttpEndpoint: "/status",
					ReadinessHealthCheckInterval:     20,
					ReadinessHealthInvocationTimeout: 2,
				},
				models.ProbeSpec{
					Type:     models.ProcessProbeType,
					Endpoint: "/status",
					Timeout:  2,
					Interval: 20,
				}),
		)
	})
})

var _ = Describe("Parse Process", func() {

	When("parsing a process", func() {
		DescribeTable("validate the correctness of the parsing logic", func(app AppManifestProcess, expected models.ProcessSpec) {
			result := parseProcess(app)
			Expect(result).To(Equal(expected))
		},
			Entry("default values",
				AppManifestProcess{
					Type: Web,
				},
				models.ProcessSpec{
					Type:   models.Web,
					Memory: "1G",
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
				},
			),
			Entry("custom memory and instances",
				AppManifestProcess{
					Type:                         Worker,
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
				models.ProcessSpec{
					Type:      models.Worker,
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
			),
			Entry("custom log rate limit and lifecycle",
				AppManifestProcess{
					Type:                  Worker,
					Command:               "run_worker",
					DiskQuota:             "512M",
					LogRateLimitPerSecond: "64K",
					Lifecycle:             "cnb",
				},
				models.ProcessSpec{
					Type:      models.Worker,
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
			),
			Entry("no health check specified",
				AppManifestProcess{
					Type:      "worker",
					Command:   "run_worker",
					DiskQuota: "256M",
					Instances: uintPtr(2),
				},
				models.ProcessSpec{
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
			),
		)
	})
})

// Helper function to create a pointer to a uint value.
func uintPtr(i uint) *uint {
	return &i
}
