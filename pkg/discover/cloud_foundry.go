package discover

import (
	"cf-application-discovery/pkg/models"
	"encoding/json"

	"github.com/cloudfoundry/go-cfclient/v3/operation"
)

func Discover(cfApp operation.AppManifest, version, space string) (models.Application, error) {
	appVersion := "1"
	if version != "" {
		appVersion = version

	}
	timeout := 60
	if cfApp.Timeout != 0 {
		timeout = int(cfApp.Timeout)
	}
	var instances int = 1
	if cfApp.Instances != nil {
		instances = int(*cfApp.Instances)
	}
	services := parseServices(cfApp.Services)
	routes := parseRoutes(cfApp.Name, cfApp.Routes, cfApp.NoRoute, cfApp.DefaultRoute)
	docker := parseDocker(cfApp.Docker)
	sidecars := parseSidecars(cfApp.Sidecars)
	processes, err := parseProcesses(cfApp)
	if err != nil {
		return models.Application{}, err
	}
	var labels, annotations map[string]*string

	if cfApp.Metadata != nil {
		labels = cfApp.Metadata.Labels
		annotations = cfApp.Metadata.Annotations
	}

	return models.Application{
		Metadata: models.Metadata{
			Version:     appVersion,
			Name:        cfApp.Name,
			Labels:      labels,
			Annotations: annotations,
		},
		Timeout:    timeout,
		Instances:  instances,
		BuildPacks: cfApp.Buildpacks,
		Env:        cfApp.Env,
		Stack:      cfApp.Stack,
		Services:   services,
		Routes:     routes,
		Docker:     docker,
		Sidecars:   sidecars,
		Processes:  processes,
	}, nil
}

func parseHealthCheck(cfType operation.AppHealthCheckType, cfEndpoint string, cfInterval, cfTimeout uint) models.Probe {
	t := models.PortProbeType
	if len(cfType) > 0 {
		t = models.ProbeType(cfType)
	}
	endpoint := "/"
	if len(cfEndpoint) > 0 {
		endpoint = cfEndpoint
	}
	timeout := 1
	if cfTimeout != 0 {
		timeout = int(cfTimeout)
	}
	interval := 30
	if cfInterval > 0 {
		interval = int(cfInterval)
	}
	return models.Probe{
		Type:     t,
		Endpoint: endpoint,
		Timeout:  timeout,
		Interval: interval,
	}
}

func parseReadinessHealthCheck(cfType operation.AppHealthCheckType, cfEndpoint string, cfInterval, cfTimeout uint) models.Probe {
	t := models.ProcessProbeType
	if len(cfType) > 0 {
		t = models.ProbeType(cfType)
	}
	endpoint := "/"
	if len(cfEndpoint) > 0 {
		endpoint = cfEndpoint
	}
	timeout := 1
	if cfTimeout != 0 {
		timeout = int(cfTimeout)
	}
	interval := 30
	if cfInterval > 0 {
		interval = int(cfInterval)
	}
	return models.Probe{
		Type:     t,
		Endpoint: endpoint,
		Timeout:  timeout,
		Interval: interval,
	}
}

func parseProcesses(cfApp operation.AppManifest) (models.Processes, error) {
	processes := models.Processes{}
	if cfApp.Processes == nil {
		return processes, nil
	}
	for _, cfProcess := range *cfApp.Processes {
		processes = append(processes, parseProcess(cfProcess))
	}
	if cfApp.Type != "" {
		// Type is the only mandatory field for the process.
		// https://github.com/SchemaStore/schemastore/blob/c06e2183289c50bdb0816050dfec002e5ebd8477/src/schemas/json/cloudfoundry-application-manifest.json#L280
		// If it's not defined it means there is no process spec at the application field level and we should return an empty structure
		proc, err := parseInlinedProcessSpec(cfApp)
		if err != nil {
			return nil, err
		}
		processes = append(processes, parseProcess(proc))
	}
	return processes, nil
}

func parseInlinedProcessSpec(cfApp operation.AppManifest) (operation.AppManifestProcess, error) {
	cfProc := operation.AppManifestProcess{}
	b, err := json.Marshal(cfApp)
	if err != nil {
		return cfProc, err
	}
	err = json.Unmarshal(b, &cfProc)
	return cfProc, err
}

func parseProcess(cfProcess operation.AppManifestProcess) models.Process {
	memory := "1G"
	if len(cfProcess.Memory) == 0 {
		memory = cfProcess.Memory
	}
	instances := 1
	if cfProcess.Instances != nil {
		instances = int(*cfProcess.Instances)
	}
	logRateLimit := "16K"
	if len(cfProcess.LogRateLimitPerSecond) > 0 {
		logRateLimit = cfProcess.LogRateLimitPerSecond
	}
	p := models.Process{
		Type:           models.ProcessType(cfProcess.Type),
		Command:        cfProcess.Command,
		DiskQuota:      cfProcess.Command,
		Memory:         memory,
		HealthCheck:    parseHealthCheck(cfProcess.HealthCheckType, cfProcess.HealthCheckHTTPEndpoint, cfProcess.HealthCheckInterval, cfProcess.HealthCheckInvocationTimeout),
		ReadinessCheck: parseReadinessHealthCheck(cfProcess.HealthCheckType, cfProcess.HealthCheckHTTPEndpoint, cfProcess.HealthCheckInterval, cfProcess.HealthCheckInvocationTimeout),
		Instances:      instances,
		LogRateLimit:   logRateLimit,
		//Lifecycle:      models.LifecycleType(cfProcess.Lifecycle),
	}
	return p
}

func parseProcessTypes(cfProcessTypes []string) []models.ProcessType {
	types := []models.ProcessType{}
	for _, cfType := range cfProcessTypes {
		types = append(types, models.ProcessType(cfType))
	}
	return types

}
func parseSidecars(cfSidecars *operation.AppManifestSideCars) models.Sidecars {
	sidecars := models.Sidecars{}
	if cfSidecars == nil {
		return sidecars
	}
	for _, cfSidecar := range *cfSidecars {
		pt := parseProcessTypes(cfSidecar.ProcessTypes)
		s := models.Sidecar{
			Name:         cfSidecar.Name,
			Command:      cfSidecar.Command,
			ProcessTypes: pt,
			Memory:       cfSidecar.Memory,
		}
		sidecars = append(sidecars, s)
	}
	return sidecars
}

func parseDocker(cfDocker *operation.AppManifestDocker) models.Docker {
	if cfDocker == nil {
		return models.Docker{}
	}
	return models.Docker{
		Image:    cfDocker.Image,
		Username: cfDocker.Username,
	}
}
func parseServices(cfServices *operation.AppManifestServices) models.Services {
	services := models.Services{}
	if cfServices == nil {
		return services
	}
	for _, svc := range *cfServices {
		s := models.Service{
			Name:        svc.Name,
			Parameters:  svc.Parameters,
			BindingName: svc.BindingName,
		}
		services = append(services, s)
	}
	return services
}

func parseRoutes(cfAppName string, cfRoutes *operation.AppManifestRoutes, noRoute, defaultRoute bool) models.Routes {
	routes := models.Routes{}
	if noRoute || cfRoutes == nil || (!noRoute && cfRoutes != nil && len(*cfRoutes) == 0) {
		return routes
	}
	if defaultRoute {
		r := models.Route{
			Route: cfAppName,
		}
		return append(routes, r)
	}
	for _, cfRoute := range *cfRoutes {
		r := models.Route{
			Route:    cfRoute.Route,
			Protocol: models.RouteProtocol(cfRoute.Protocol),
		}
		routes = append(routes, r)
	}
	return routes
}
