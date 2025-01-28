package models

// Application represents an interpretation of a runtime Cloud Foundry application. This structure differs in that
// the information it contains has been processed to simplify its transformation to a Kubernetes manifest using MTA
type Application struct {
	// Metadata captures the name, labels and annotations in the application.
	Metadata Metadata `json:",inline" validate:"required"`
	// Env captures the `env` field values in the CF application manifest.
	Env map[string]string `json:"env,omitempty"`
	// Routes represent the routes that are made available by the application.
	Routes []Route `json:"route,omitempty"`
	// Services captures the `services` field values in the CF application manifest.
	Services []Service `json:"service,omitempty"`
	// Processes captures the `processes` field values in the CF application manifest.
	Processes []Process `json:"process,omitempty"`
	// Sidecars captures the `sidecars` field values in the CF application manifest.
	Sidecars []Sidecar `json:"sidecar,omitempty"`
	// Stack represents the `stack` field in the application manifest.
	// The value is captured for information purposes because it has no relevance
	// in Kubernetes.
	Stack string `json:"stack,omitempty"`
	// StartupTimeout specifies the maximum time allowed for an application to
	// respond to readiness or health checks during startup.
	// If the application does not respond within this time, the platform will mark
	// the deployment as failed. The default value is 60 seconds.
	// https://github.com/cloudfoundry/docs-dev-guide/blob/96f19d9d67f52ac7418c147d5ddaa79c957eec34/deploy-apps/large-app-deploy.html.md.erb#L35
	StartupTimeout *uint `json:"startupTimeout,omitempty"`
	// BuildPacks capture the buildpacks defined in the CF application manifest.
	BuildPacks []string `json:"buildPacks,omitempty"`
	// Docker captures the Docker specification in the CF application manifest.
	Docker *Docker `json:"docker,omitempty"`
}

type Docker struct {
	// Image represents the pullspect where the container image is located.
	Image string `json:"image" validate:"required"`
	// Username captures the username to authenticate against the container registry.
	Username string `json:"username,omitempty"`
}

// Metadata captures the name, labels and annotations in the application
type Metadata struct {
	// Name capture the `name` field int CF application manifest
	Name string `json:"name" validate:"required"`
	// Space captures the `space` where the CF application is deployed at runtime. The field is empty if the
	// application is discovered directly from the CF manifest. It is equivalent to a Namespace in Kubernetes.
	Space string `json:"space,omitempty"`
	// Labels capture the labels as defined in the `annotations` field in the CF application manifest
	Labels map[string]string `json:"labels,omitempty"`
	// Annotations capture the annotations as defined in the `labels` field in the CF application manifest
	Annotations map[string]string `json:"annotations,omitempty"`
	// Version captures the version of the manifest containing the resulting CF application manifests list retrieved via REST API.
	// Only version 1 is supported at this moment. See https://docs.cloudfoundry.org/devguide/deploy-apps/manifest-attributes.html#manifest-schema-version
	Version string `json:"version"`
}

// Routes represents a slice of Routes
type Routes []Route

// Route captures the key elements that define a Route in a string that maps to a URL structure. These values
// are captured as runtime routes, meaning that if the CF Application manifest is configured to disable all routes
// with the `no-route` value, it will translate into an empty slice.
// By default Cloud Foundry attempts to create a route for each application unless the `no-route` field is set to true.
// For further details check: https://docs.cloudfoundry.org/devguide/deploy-apps/manifest-attributes.html#no-route
// and https://docs.cloudfoundry.org/devguide/deploy-apps/manifest-attributes.html#random-route
// Example
// ---
//
//	...
//	routes:
//	- route: example.com
//	  protocol: http2
//	- route: www.example.com/foo
//	- route: tcp-example.com:1234
type Route struct {
	// Route captures the domain name, port and path of the route.
	Route string `json:"url" validate:"required"`
	// Protocol captures the protocol type: http, http2 or tcp. Note that the CF `protocol` field is only available
	// for CF deployments that use HTTP/2 routing.
	Protocol RouteProtocol `json:"protocol" validate:"required,oneof=http http2 tcp"`
}

type RouteProtocol string

const (
	HTTP  RouteProtocol = "http"
	HTTP2 RouteProtocol = "http2"
	TCP   RouteProtocol = "tcp"
)

// Services represents a slice of Service
type Services []Service

// Service contains the specification for an existing Cloud Foundry service required by the application.
// Examples:
// ---
//
//	...
//	services:
//	  - service-1
//	  - name: service-2
//	  - name: service-3
//	    parameters:
//	      key-1: value-1
//	      key-2: [value-2, value-3]
//	      key-3: ... any other kind of value ...
//	  - name: service-4
//	    binding_name: binding-1
type Service struct {
	// Name represents the name of the Cloud Foundry service required by the
	// application. This field represents the runtime name of the service, captured
	// from the 3 different cases where the service name can be listed.
	// For more information check https://docs.cloudfoundry.org/devguide/deploy-apps/manifest-attributes.html#services-block
	Name string `json:"name" validate:"required"`
	// Parameters contain the k/v relationship for the aplication to bind to the service
	Parameters map[string]interface{} `json:"parameters,omitempty"`
	// BindingName captures the name of the service to bind to.
	BindingName string `json:"bindingName,omitempty"`
}

// Processes represents a slice of Processes.
type Processes []Process

// Process represents the abstraction of the specification of a Cloud Foundry Process.
// For more information check https://docs.cloudfoundry.org/devguide/deploy-apps/manifest-attributes.html#processes
type Process struct {
	// Type captures the `type` field in the Process specification.
	// Accepted values are `web` or `worker`
	Type ProcessType `json:"type,omitempty"`
	// Command represents the command used to run the process.
	Command []string `json:"command,omitempty"`
	// DiskQuota represents the amount of persistent disk requested by the process.
	DiskQuota string `json:"disk,omitempty"`
	// Memory represents the amount of memory requested by the process.
	Memory string `json:"memory,omitempty"`
	// HealthCheck captures the health check information
	HealthCheck *Probe `json:"healthCheck,omitempty"`
	// ReadinessCheck captures the readiness check information.
	ReadinessCheck *Probe `json:"readinessCheck,omitempty"`
	// Replicas represents the number of instances for this process to run.
	Replicas uint `json:"replicas" validate:"required"`
	// LogRateLimit represents the maximum amount of logs to be captured per second.
	LogRateLimit string `json:"logRateLimit,omitempty"`
}

type Sidecars []Sidecar

// Sidecar captures the information of a Sidecar process
// https://docs.cloudfoundry.org/devguide/deploy-apps/manifest-attributes.html#sidecars
type Sidecar struct {
	// Name represents the name of the Sidecar
	Name string `json:"name" validate:"required"`
	// ProcessTypes captures the different process types defined for the sidecar.
	// Compared to a Process, which has only one type, sidecar processes can
	// accumulate more than one type.
	ProcessTypes []ProcessType `json:"processType" validate:"required,oneof=worker web"`
	// Command captures the command to run the sidecar
	Command []string `json:"command" validate:"required"`
	// Memory represents the amount of memory to allocate to the sidecar.
	// It's an optional field.
	Memory string `json:"memory,omitempty"`
}

// Probe captures the fields for managing health checks. For more information check https://docs.cloudfoundry.org/devguide/deploy-apps/healthchecks.html
type Probe struct {
	// Endpoint represents the URL location where to perform the probe check.
	Endpoint string `json:"endpoint" validate:"required"`
	// Timeout represents the number of seconds in which the probe check can be considered as timedout.
	// https://docs.cloudfoundry.org/devguide/deploy-apps/manifest-attributes.html#timeout
	Timeout uint `json:"timeout" validate:"required"`
	// Interval represents the number of seconds between probe checks.
	Interval uint `json:"interval" validate:"required"`
	// Type specifies the type of health check to perform
	Type string `json:"type" validate:"required,oneof=http tcp process"`
}

type ProcessTypes []ProcessType

type ProcessType string

const (
	// Web represents a `web` application type
	Web ProcessType = "web"
	// Worker represents a `worker` application type
	Worker ProcessType = "worker"
)
