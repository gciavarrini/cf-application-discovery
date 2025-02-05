package discover

import (
	"cf-application-discovery/pkg/models"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Health Checks tests", func() {

	When("parsing health check probe", func() {
		defaultProbeSpec := models.ProbeSpec{
			Type:     models.PortProbeType,
			Endpoint: "/",
			Timeout:  1,
			Interval: 30,
		}
		overrideDefaultProbeSpec := func(overrides ...func(*models.ProbeSpec)) models.ProbeSpec {
			spec := defaultProbeSpec
			for _, override := range overrides {
				override(&spec)
			}
			return spec
		}
		DescribeTable("validate the correctness of the parsing logic", func(app AppManifestProcess, expected models.ProbeSpec) {
			result := parseHealthCheck(app.HealthCheckType, app.HealthCheckHTTPEndpoint, app.HealthCheckInterval, app.HealthCheckInvocationTimeout)
			// Use Gomega's Expect function for assertions
			Expect(result).To(Equal(expected))
		},
			Entry("with default values",
				AppManifestProcess{},
				defaultProbeSpec),
			Entry("with endpoint only",
				AppManifestProcess{
					HealthCheckHTTPEndpoint: "/example.com",
				},
				overrideDefaultProbeSpec(func(spec *models.ProbeSpec) {
					spec.Endpoint = "/example.com"
				})),
			Entry("with interval only",
				AppManifestProcess{
					HealthCheckInterval: 42,
				},
				overrideDefaultProbeSpec(func(spec *models.ProbeSpec) {
					spec.Interval = 42
				})),
			Entry("with timeout only",
				AppManifestProcess{
					HealthCheckInvocationTimeout: 42,
				},
				overrideDefaultProbeSpec(func(spec *models.ProbeSpec) {
					spec.Timeout = 42
				})),
			Entry("with type only",
				AppManifestProcess{
					HealthCheckType: "http",
				},
				overrideDefaultProbeSpec(func(spec *models.ProbeSpec) {
					spec.Type = models.HTTPProbeType
				})),
		)
	})

	When("parsing readiness health check probe", func() {
		defaultProbeSpec := models.ProbeSpec{
			Type:     models.ProcessProbeType,
			Endpoint: "/",
			Timeout:  1,
			Interval: 30,
		}
		overrideDefaultProbeSpec := func(overrides ...func(*models.ProbeSpec)) models.ProbeSpec {
			spec := defaultProbeSpec
			for _, override := range overrides {
				override(&spec)
			}
			return spec
		}
		DescribeTable("validate the correctness of the parsing logic", func(app AppManifestProcess, expected models.ProbeSpec) {
			result := parseReadinessHealthCheck(app.ReadinessHealthCheckType, app.ReadinessHealthCheckHttpEndpoint, app.ReadinessHealthCheckInterval, app.ReadinessHealthInvocationTimeout)
			// Use Gomega's Expect function for assertions
			Expect(result).To(Equal(expected))
		},
			Entry("with default values",
				AppManifestProcess{},
				defaultProbeSpec),
			Entry("with type only",
				AppManifestProcess{
					ReadinessHealthCheckType: Http,
				},
				overrideDefaultProbeSpec(func(spec *models.ProbeSpec) {
					spec.Type = models.HTTPProbeType
				})),
			Entry("with endpoint only",
				AppManifestProcess{
					ReadinessHealthCheckHttpEndpoint: "/example.com",
				},
				overrideDefaultProbeSpec(func(spec *models.ProbeSpec) {
					spec.Endpoint = "/example.com"
				})),
			Entry("with interval only",
				AppManifestProcess{
					ReadinessHealthCheckInterval: 42,
				},
				overrideDefaultProbeSpec(func(spec *models.ProbeSpec) {
					spec.Interval = 42
				})),
			Entry("with timeout only",
				AppManifestProcess{
					ReadinessHealthInvocationTimeout: 42,
				},
				overrideDefaultProbeSpec(func(spec *models.ProbeSpec) {
					spec.Timeout = 42
				})),
		)
	})
})
var _ = Describe("Parse Process", func() {

	When("parsing a process", func() {
		defaultProcessSpec := models.ProcessSpec{
			Type:   "",
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
		}
		overrideDefaultProcessSpec := func(overrides ...func(*models.ProcessSpec)) models.ProcessSpec {
			spec := defaultProcessSpec
			for _, override := range overrides {
				override(&spec)
			}
			return spec
		}

		DescribeTable("validate the correctness of the parsing logic", func(app AppManifestProcess, expected models.ProcessSpec) {
			result := parseProcess(app)
			Expect(result).To(Equal(expected))
		},
			Entry("default values",
				AppManifestProcess{},
				defaultProcessSpec,
			),
			Entry("with memory only",
				AppManifestProcess{
					Memory: "512M",
				},
				overrideDefaultProcessSpec(func(spec *models.ProcessSpec) {
					spec.Memory = "512M"
				}),
			),
			Entry("with instance only",
				AppManifestProcess{
					Instances: uintPtr(42),
				},
				overrideDefaultProcessSpec(func(spec *models.ProcessSpec) {
					spec.Instances = 42
				}),
			),
			Entry("with only lograte",
				AppManifestProcess{
					LogRateLimitPerSecond: "42K",
				},
				overrideDefaultProcessSpec(func(spec *models.ProcessSpec) {
					spec.LogRateLimit = "42K"
				}),
			),
		)
	})
	When("parsing a process type", func() {
		DescribeTable("validate the correctness of the parsing logic", func(cfProcessTypes []string, expected []models.ProcessType) {
			result := parseProcessTypes(cfProcessTypes)
			Expect(result).To(Equal(expected))
		},
			Entry("default values with nil input",
				nil,
				[]models.ProcessType{},
			),
			Entry("default values with empty input",
				[]string{},
				[]models.ProcessType{},
			),
			Entry("with web type",
				[]string{"web"},
				[]models.ProcessType{models.Web},
			),
			Entry("with worker type",
				[]string{"worker"},
				[]models.ProcessType{models.Worker},
			),
			Entry("with unknown type",
				[]string{"unknown"},
				[]models.ProcessType{"unknown"},
			),
			Entry("with multiple type",
				[]string{"web", "worker", "unknown"},
				[]models.ProcessType{models.Web, models.Worker, "unknown"},
			),
		)
	})
})

var _ = Describe("Parse Sidecars", func() {

	When("parsing sidecars", func() {
		DescribeTable("validate the correctness of the parsing logic", func(cfSidecars *AppManifestSideCars, expected models.Sidecars) {
			result := parseSidecars(cfSidecars)
			Expect(result).To(Equal(expected))
		},
			Entry("default values with nil input",
				nil,
				models.Sidecars{},
			),
			Entry("default values with empty input",
				&AppManifestSideCars{},
				models.Sidecars{},
			),
			Entry("one sidecar with only name",
				&AppManifestSideCars{
					AppManifestSideCar{
						Name: "test-name",
					},
				},
				models.Sidecars{
					{
						Name:         "test-name",
						ProcessTypes: []models.ProcessType{},
					},
				},
			),
			Entry("one sidecar with only command",
				&AppManifestSideCars{
					AppManifestSideCar{
						Command: "test-command",
					},
				},
				models.Sidecars{
					{
						Command:      "test-command",
						ProcessTypes: []models.ProcessType{},
					},
				},
			),
			Entry("one sidecar with only process types",
				&AppManifestSideCars{
					AppManifestSideCar{
						ProcessTypes: []string{"web", "worker"},
					},
				},
				models.Sidecars{
					{
						ProcessTypes: []models.ProcessType{models.Web, models.Worker},
					},
				},
			),
		)
	})
})

// Helper function to create a pointer to a uint value.
func uintPtr(i uint) *uint {
	return &i
}
