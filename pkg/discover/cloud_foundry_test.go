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
					Instances: ptrTo(uint(42)),
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
		DescribeTable("validate the correctness of the parsing logic", func(cfProcessTypes []AppProcessType, expected []models.ProcessType) {
			result := parseProcessTypes(cfProcessTypes)
			Expect(result).To(Equal(expected))
		},
			Entry("default values with nil input",
				nil,
				[]models.ProcessType{},
			),
			Entry("default values with empty input",
				[]AppProcessType{},
				[]models.ProcessType{},
			),
			Entry("with web type",
				[]AppProcessType{Web},
				[]models.ProcessType{models.Web},
			),
			Entry("with worker type",
				[]AppProcessType{Worker},
				[]models.ProcessType{models.Worker},
			),
			Entry("with multiple type",
				[]AppProcessType{"web", "worker"},
				[]models.ProcessType{models.Web, models.Worker},
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
				nil,
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
						ProcessTypes: []AppProcessType{"web", "worker"},
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

var _ = Describe("Parse Routes", func() {

	When("parsing the route information", func() {
		DescribeTable("validate the correctness of the parsing logic for the route specification", func(app AppManifest, expected models.RouteSpec) {
			result := parseRouteSpec(app.Routes, app.RandomRoute, app.NoRoute)
			Expect(result).To(Equal(expected))
		},
			Entry("when routes are nil, no-route and random-route are false", AppManifest{}, models.RouteSpec{}),
			Entry("when routes are empty, no-route and random-route are false", AppManifest{Routes: &AppManifestRoutes{}}, models.RouteSpec{Routes: models.Routes{}}),
			Entry("when routes are not empty, no-route and random-route are false",
				AppManifest{
					Routes: &AppManifestRoutes{{Route: "foo.bar"}}},
				models.RouteSpec{
					Routes: models.Routes{{Route: "foo.bar"}},
				}),
			Entry("when routes are nil, no-route is true and random-route is false",
				AppManifest{
					NoRoute: true,
				},
				models.RouteSpec{
					NoRoute: true,
				}),
			Entry("when routes have one entry and no-route is true",
				AppManifest{
					NoRoute: true,
					Routes:  &AppManifestRoutes{{Route: "foo.bar"}}},
				models.RouteSpec{
					NoRoute: true,
				}),
			Entry("when routes are nil, no-route is false and random-route is true",
				AppManifest{
					RandomRoute: true,
				},
				models.RouteSpec{
					RandomRoute: true,
				}),
			Entry("when routes have two entries, no-route and random-route are false",
				AppManifest{
					Routes: &AppManifestRoutes{{Route: "foo.bar"}, {Route: "bar.foo"}}},
				models.RouteSpec{
					Routes: models.Routes{{Route: "foo.bar"}, {Route: "bar.foo"}}},
			),
		)

		DescribeTable("validate the correctness of the parsing logic of the route structure", func(routes AppManifestRoutes, expected models.Routes) {
			result := parseRoutes(routes)
			Expect(result).To(Equal(expected))
		},
			Entry("when routes are nil", nil, nil),
			Entry("when routes are empty", AppManifestRoutes{}, models.Routes{}),
			Entry("when routes contain one element with only route field defined", AppManifestRoutes{{Route: "foo.bar"}}, models.Routes{{Route: "foo.bar"}}),
			Entry("when routes contain one element with only protocol field defined", AppManifestRoutes{{Protocol: HTTP2}}, models.Routes{{Protocol: models.HTTP2RouteProtocol}}),
			Entry("when routes contain one element with only options field defined with round-robin load balancing",
				AppManifestRoutes{
					{Options: &AppRouteOptions{LoadBalancing: "round-robin"}}},
				models.Routes{
					{Options: models.RouteOptions{LoadBalancing: models.RoundRobinLoadBalancingType}}}),
			Entry("when routes contain one element with only options field defined with least-connection load balancing",
				AppManifestRoutes{
					{Options: &AppRouteOptions{LoadBalancing: "least-connection"}}},
				models.Routes{
					{Options: models.RouteOptions{LoadBalancing: models.LeastConnectionLoadBalancingType}}}),
			Entry("when routes contain one element with all fields populated",
				AppManifestRoutes{
					{
						Route:    "foo.bar",
						Protocol: TCP,
						Options:  &AppRouteOptions{LoadBalancing: "least-connection"},
					}},
				models.Routes{
					{
						Route:    "foo.bar",
						Protocol: models.TCPRouteProtocol,
						Options:  models.RouteOptions{LoadBalancing: models.LeastConnectionLoadBalancingType}}}),
			Entry("when routes contain two elements",
				AppManifestRoutes{
					{
						Route:    "foo.bar",
						Protocol: TCP,
						Options:  &AppRouteOptions{LoadBalancing: "round-robin"},
					},
					{
						Route:    "bar.foo",
						Protocol: HTTP1,
					}},
				models.Routes{
					{
						Route:    "foo.bar",
						Protocol: models.TCPRouteProtocol,
						Options:  models.RouteOptions{LoadBalancing: models.RoundRobinLoadBalancingType}},
					{
						Route:    "bar.foo",
						Protocol: models.HTTPRouteProtocol,
					}}),
		)
	})

})

var _ = Describe("parse Services", func() {
	When("parsing the service information", func() {
		DescribeTable("validate the correctness of the parsing logic", func(services AppManifestServices, expected models.Services) {
			result := parseServices(&services)
			Expect(result).To(Equal(expected))
		},
			Entry("when services are nil", nil, models.Services{}),
			Entry("when services are empty", AppManifestServices{}, models.Services{}),
			Entry("when one service is provided with only name populated", AppManifestServices{{Name: "foo"}}, models.Services{{Name: "foo"}}),
			Entry("when one service is provided with parameters provided",
				AppManifestServices{
					{Parameters: map[string]interface{}{"foo": "bar"}},
				},
				models.Services{
					{Parameters: map[string]interface{}{"foo": "bar"}},
				}),
			Entry("when one service is provided with binding name provided", AppManifestServices{{BindingName: "foo_service"}}, models.Services{{BindingName: "foo_service"}}),
			Entry("when one service is provided with name, parameters and binding name are provided",
				AppManifestServices{
					{
						Name:        "foo_name",
						Parameters:  map[string]interface{}{"foo": "bar"},
						BindingName: "foo_service",
					},
				},
				models.Services{
					{
						Name:        "foo_name",
						Parameters:  map[string]interface{}{"foo": "bar"},
						BindingName: "foo_service",
					},
				}),
			Entry("when two services are provided with a unique name populated for each one",
				AppManifestServices{
					{Name: "foo"},
					{Name: "bar"},
				},
				models.Services{
					{Name: "foo"},
					{Name: "bar"},
				}),
		)
	})
})

var _ = Describe("parse metadata", func() {
	When("parsing the metadata information", func() {
		DescribeTable("validate the correctness of the parsing logic", func(metadata Metadata, version, space string, expected models.Metadata) {
			result, err := Discover(AppManifest{Metadata: &metadata}, version, space)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Metadata).To(Equal(expected))
		},

			Entry("when metadata is nil and version and space are empty", nil, "", "", models.Metadata{Version: "1"}),
			Entry("when empty metadata, version and space", Metadata{}, "", "", models.Metadata{Version: "1"}),
			Entry("when version is provided", Metadata{}, "2", "", models.Metadata{Version: "2"}),
			Entry("when space is provided", Metadata{}, "", "default", models.Metadata{Version: "1", Space: "default"}),
			Entry("when labels are provided", Metadata{Labels: map[string]*string{"foo": ptrTo("bar")}}, "", "", models.Metadata{Version: "1", Labels: map[string]*string{"foo": ptrTo("bar")}}),
			Entry("when annotations are provided", Metadata{Annotations: map[string]*string{"bar": ptrTo("foo")}}, "", "", models.Metadata{Version: "1", Annotations: map[string]*string{"bar": ptrTo("foo")}}),
			Entry("when all fields are provided",
				Metadata{
					Labels:      map[string]*string{"foo": ptrTo("bar")},
					Annotations: map[string]*string{"bar": ptrTo("foo")}},
				"2",
				"default",
				models.Metadata{
					Labels:      map[string]*string{"foo": ptrTo("bar")},
					Annotations: map[string]*string{"bar": ptrTo("foo")},
					Version:     "2",
					Space:       "default",
				}),
		)
	})

})
var _ = Describe("Parse Application", func() {
	When("parsing the application information", func() {
		DescribeTable("validate the correctness of the parsing logic", func(app AppManifest, version, space string, expected models.Application) {
			result, err := Discover(app, version, space)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(expected))
		},
			Entry("when app is empty",
				AppManifest{},
				"",
				"",
				models.Application{
					Metadata:  models.Metadata{Version: "1"},
					Timeout:   60,
					Instances: 1,
				},
			),
			Entry("when timeout is set",
				AppManifest{
					AppManifestProcess: AppManifestProcess{Timeout: 30},
				},
				"",
				"",
				models.Application{
					Metadata:  models.Metadata{Version: "1"},
					Timeout:   30,
					Instances: 1,
				},
			),
			Entry("when instances is set",
				AppManifest{
					AppManifestProcess: AppManifestProcess{Instances: ptrTo(uint(2))},
				},
				"",
				"",
				models.Application{
					Metadata:  models.Metadata{Version: "1"},
					Timeout:   60,
					Instances: 2,
				},
			),
			Entry("when buildpacks are set",
				AppManifest{
					Buildpacks: []string{"foo", "bar"},
				},
				"",
				"",
				models.Application{
					Metadata:   models.Metadata{Version: "1"},
					Timeout:    60,
					Instances:  1,
					BuildPacks: []string{"foo", "bar"},
				},
			),
			Entry("when environment values are set",
				AppManifest{
					Env: map[string]string{"foo": "bar"},
				},
				"",
				"",
				models.Application{
					Metadata:  models.Metadata{Version: "1"},
					Timeout:   60,
					Instances: 1,
					Env:       map[string]string{"foo": "bar"},
				},
			),
			Entry("when all fields are set",
				AppManifest{
					Name:       "foo",
					Buildpacks: []string{"foo", "bar"},
					Docker: &AppManifestDocker{
						Image:    "foo.bar:latest",
						Username: "foo@bar.org",
					},
					RandomRoute: true,
					Routes: &AppManifestRoutes{
						{
							Route:    "foo.bar.org",
							Protocol: HTTP2,
							Options:  &AppRouteOptions{LoadBalancing: "least-connection"},
						},
					},
					Env: map[string]string{"foo": "bar"},
					Services: &AppManifestServices{
						{
							Name:        "foo",
							BindingName: "foo_service",
							Parameters:  map[string]interface{}{"foo": "bar"},
						},
					},
					Sidecars: &AppManifestSideCars{
						{
							Name:         "foo_sidecar",
							ProcessTypes: []AppProcessType{Web, Worker},
							Command:      "echo hello world",
							Memory:       "2G",
						},
					},
					Stack: "docker",
					Metadata: &Metadata{
						Labels:      map[string]*string{"foo": ptrTo("label")},
						Annotations: map[string]*string{"bar": ptrTo("annotation")},
					},
					AppManifestProcess: AppManifestProcess{
						Timeout:   100,
						Instances: ptrTo(uint(5)),
					},
					Processes: &AppManifestProcesses{
						{
							Type:                             Web,
							Command:                          "sleep 100",
							DiskQuota:                        "100M",
							HealthCheckType:                  Http,
							HealthCheckHTTPEndpoint:          "/health",
							HealthCheckInvocationTimeout:     10,
							HealthCheckInterval:              60,
							ReadinessHealthCheckType:         Port,
							ReadinessHealthCheckHttpEndpoint: "localhost:8443",
							ReadinessHealthInvocationTimeout: 99,
							ReadinessHealthCheckInterval:     15,
							Instances:                        ptrTo(uint(2)),
							LogRateLimitPerSecond:            "30k",
							Memory:                           "2G",
							Timeout:                          120,
							Lifecycle:                        "container",
						},
					},
				},
				"2",
				"default",
				models.Application{
					Metadata: models.Metadata{
						Version:     "2",
						Name:        "foo",
						Labels:      map[string]*string{"foo": ptrTo("label")},
						Annotations: map[string]*string{"bar": ptrTo("annotation")},
						Space:       "default",
					},
					BuildPacks: []string{"foo", "bar"},
					Stack:      "docker",
					Timeout:    100,
					Instances:  5,
					Env:        map[string]string{"foo": "bar"},
					Routes: models.RouteSpec{
						RandomRoute: true,
						Routes: models.Routes{
							{
								Route:    "foo.bar.org",
								Protocol: models.HTTP2RouteProtocol,
								Options: models.RouteOptions{
									LoadBalancing: models.LeastConnectionLoadBalancingType,
								},
							},
						},
					},
					Docker: models.Docker{
						Image:    "foo.bar:latest",
						Username: "foo@bar.org",
					},
					Services: models.Services{
						{
							Name:        "foo",
							BindingName: "foo_service",
							Parameters:  map[string]interface{}{"foo": "bar"},
						},
					},
					Sidecars: models.Sidecars{
						{
							Name:         "foo_sidecar",
							ProcessTypes: []models.ProcessType{models.Web, models.Worker},
							Command:      "echo hello world",
							Memory:       "2G",
						},
					},
					Processes: models.Processes{
						{
							Type:         models.Web,
							Command:      "sleep 100",
							DiskQuota:    "100M",
							Instances:    2,
							LogRateLimit: "30k",
							Memory:       "2G",
							Lifecycle:    "container",
							HealthCheck: models.ProbeSpec{
								Endpoint: "/health",
								Timeout:  10,
								Interval: 60,
								Type:     models.HTTPProbeType,
							},
							ReadinessCheck: models.ProbeSpec{
								Endpoint: "localhost:8443",
								Timeout:  99,
								Interval: 15,
								Type:     models.PortProbeType,
							},
						},
					},
				},
			),
		)
	})
})

var _ = Describe("Parse docker", func() {
	When("parsing the docker information", func() {
		DescribeTable("validate the correctness of the parsing logic", func(docker AppManifestDocker, expected models.Docker) {
			result := parseDocker(&docker)
			Expect(result).To(Equal(expected))
		},
			Entry("when docker is nil", nil, models.Docker{}),
			Entry("when docker is empty", AppManifestDocker{}, models.Docker{}),
			Entry("when docker image is populated", AppManifestDocker{Image: "python3:latest"}, models.Docker{Image: "python3:latest"}),
			Entry("when docker username is populated", AppManifestDocker{Username: "foo@bar.org"}, models.Docker{Username: "foo@bar.org"}),
			Entry("when docker image and username are populated",
				AppManifestDocker{
					Image:    "python3:latest",
					Username: "foo@bar.org"},
				models.Docker{Image: "python3:latest",
					Username: "foo@bar.org"}),
		)
	})
})

// Helper function to create a pointer of a given type
func ptrTo[T comparable](t T) *T {
	return &t
}
