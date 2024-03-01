package provider

import (
	"context"
	"time"

	cfv3operation "github.com/cloudfoundry-community/go-cfclient/v3/operation"
	cfv3resource "github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Type AppType representing Schema Attribute from function Schema in go type from resource_appManifest.go file
type AppType struct {
	Name                                  types.String      `tfsdk:"name"`
	Space                                 types.String      `tfsdk:"space"`
	Org                                   types.String      `tfsdk:"org"`
	Stack                                 types.String      `tfsdk:"stack"`
	Buildpacks                            types.Set         `tfsdk:"buildpacks"`
	Path                                  types.String      `tfsdk:"path"`
	SourceCodeHash                        types.String      `tfsdk:"source_code_hash"`
	DockerImage                           types.String      `tfsdk:"docker_image"`
	DockerCredentials                     DockerCredentials `tfsdk:"docker_credentials"`
	Strategy                              types.String      `tfsdk:"strategy"`
	ServiceBindings                       []ServiceBinding  `tfsdk:"service_bindings"`
	Routes                                []Route           `tfsdk:"routes"`
	Environment                           types.Map         `tfsdk:"environment"`
	HealthCheckInterval                   types.Int64       `tfsdk:"health_check_interval"`
	ReadinessHealthCheckType              types.String      `tfsdk:"readiness_health_check_type"`
	ReadinessHealthCheckHttpEndpoint      types.String      `tfsdk:"readiness_health_check_http_endpoint"`
	ReadinessHealthCheckInvocationTimeout types.Int64       `tfsdk:"readiness_health_check_invocation_timeout"`
	ReadinessHealthCheckInterval          types.Int64       `tfsdk:"readiness_health_check_interval"`
	LogRateLimitPerSecond                 types.String      `tfsdk:"log_rate_limit_per_second"`
	NoRoute                               types.Bool        `tfsdk:"no_route"`
	RandomRoute                           types.Bool        `tfsdk:"random_route"`
	Processes                             []Process         `tfsdk:"processes"`
	Sidecars                              []Sidecar         `tfsdk:"sidecars"`
	ID                                    types.String      `tfsdk:"id"`
	CreatedAt                             types.String      `tfsdk:"created_at"`
	UpdatedAt                             types.String      `tfsdk:"updated_at"`
	Command                               types.String      `tfsdk:"command"`
	DiskQuota                             types.String      `tfsdk:"disk_quota"`
	HealthCheckHttpEndpoint               types.String      `tfsdk:"health_check_http_endpoint"`
	HealthCheckInvocationTimeout          types.Int64       `tfsdk:"health_check_invocation_timeout"`
	HealthCheckType                       types.String      `tfsdk:"health_check_type"`
	Instances                             types.Int64       `tfsdk:"instances"`
	Memory                                types.String      `tfsdk:"memory"`
	Timeout                               types.Int64       `tfsdk:"timeout"`
}

type Sidecar struct {
	Name         types.String `tfsdk:"name"`
	Command      types.String `tfsdk:"command"`
	ProcessTypes types.Set    `tfsdk:"process_types"`
	Memory       types.String `tfsdk:"memory"`
}
type Process struct {
	Type                         types.String `tfsdk:"type"`
	Command                      types.String `tfsdk:"command"`
	DiskQuota                    types.String `tfsdk:"disk_quota"`
	HealthCheckHttpEndpoint      types.String `tfsdk:"health_check_http_endpoint"`
	HealthCheckInvocationTimeout types.Int64  `tfsdk:"health_check_invocation_timeout"`
	HealthCheckType              types.String `tfsdk:"health_check_type"`
	Instances                    types.Int64  `tfsdk:"instances"`
	Memory                       types.String `tfsdk:"memory"`
	Timeout                      types.Int64  `tfsdk:"timeout"`
}

type DockerCredentials struct {
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

type ServiceBinding struct {
	ServiceInstance types.String `tfsdk:"service_instance"`
	Param           types.Map    `tfsdk:"param"`
}

type Route struct {
	Route    types.String `tfsdk:"route"`
	Protocol types.String `tfsdk:"protocol"`
}

// mapAppTypeToValues function maps AppType to cfv3resource manifest type
func (appType *AppType) mapAppTypeToValues(ctx context.Context) (*cfv3operation.AppManifest, diag.Diagnostics) {
	var diags diag.Diagnostics
	appmanifest := cfv3operation.AppManifest{}
	appmanifest.Name = appType.Name.ValueString()
	if !appType.Stack.IsNull() {
		appmanifest.Stack = appType.Stack.ValueString()
	}
	if !appType.Buildpacks.IsNull() {
		var buildpacks []string
		diags = appType.Buildpacks.ElementsAs(ctx, &buildpacks, false)
		appmanifest.Buildpacks = buildpacks
	}
	if !appType.DockerImage.IsNull() {
		appManifestDocker := cfv3operation.AppManifestDocker{
			Image: appType.DockerImage.ValueString(),
		}
		if appType.DockerCredentials != (DockerCredentials{}) {
			appManifestDocker.Username = appType.DockerCredentials.Username.ValueString()
		}
	}
	if len(appType.ServiceBindings) != 0 {
		var serviceBindings cfv3operation.AppManifestServices
		for _, serviceBinding := range appType.ServiceBindings {
			var appParam map[string]interface{}
			diags = serviceBinding.Param.ElementsAs(ctx, &appParam, false)
			serviceBindings = append(serviceBindings, cfv3operation.AppManifestService{
				Name:       serviceBinding.ServiceInstance.ValueString(),
				Parameters: appParam,
			})
		}
		appmanifest.Services = &serviceBindings
	}
	if len(appType.Routes) != 0 {
		var routes cfv3operation.AppManifestRoutes
		for _, route := range appType.Routes {
			routes = append(routes, cfv3operation.AppManifestRoute{
				Route:    route.Route.ValueString(),
				Protocol: cfv3operation.AppRouteProtocol(route.Protocol.ValueString()),
			})
		}
		appmanifest.Routes = &routes
	}
	if !appType.Environment.IsNull() {
		var env map[string]string
		diags = appType.Environment.ElementsAs(ctx, &env, false)
		appmanifest.Env = env
	}
	if !appType.HealthCheckInterval.IsNull() {
		appmanifest.HealthCheckInterval = uint(appType.HealthCheckInterval.ValueInt64())
	}
	if !appType.ReadinessHealthCheckType.IsNull() {
		appmanifest.ReadinessHealthCheckType = appType.ReadinessHealthCheckType.ValueString()
	}
	if !appType.ReadinessHealthCheckHttpEndpoint.IsNull() {
		appmanifest.ReadinessHealthCheckHttpEndpoint = appType.ReadinessHealthCheckHttpEndpoint.ValueString()
	}
	if !appType.ReadinessHealthCheckInvocationTimeout.IsNull() {
		appmanifest.ReadinessHealthInvocationTimeout = uint(appType.ReadinessHealthCheckInvocationTimeout.ValueInt64())
	}
	if !appType.ReadinessHealthCheckInterval.IsNull() {
		appmanifest.ReadinessHealthCheckInterval = uint(appType.ReadinessHealthCheckInterval.ValueInt64())
	}
	if !appType.LogRateLimitPerSecond.IsNull() {
		appmanifest.LogRateLimitPerSecond = appType.LogRateLimitPerSecond.ValueString()
	}
	if !appType.NoRoute.IsNull() {
		appmanifest.NoRoute = appType.NoRoute.ValueBool()
	}
	if !appType.RandomRoute.IsNull() {
		appmanifest.RandomRoute = appType.RandomRoute.ValueBool()
	}
	if len(appType.Processes) != 0 {
		var processes cfv3operation.AppManifestProcesses
		for _, process := range appType.Processes {
			processes = append(processes, cfv3operation.AppManifestProcess{
				Type:                         cfv3operation.AppProcessType(process.Type.ValueString()),
				Command:                      process.Command.ValueString(),
				DiskQuota:                    process.DiskQuota.ValueString(),
				HealthCheckHTTPEndpoint:      process.HealthCheckHttpEndpoint.ValueString(),
				HealthCheckInvocationTimeout: uint(process.HealthCheckInvocationTimeout.ValueInt64()),
				HealthCheckType:              cfv3operation.AppHealthCheckType(process.HealthCheckType.ValueString()),
				Instances:                    uint(process.Instances.ValueInt64()),
				Memory:                       process.Memory.ValueString(),
				Timeout:                      uint(process.Timeout.ValueInt64()),
			})
		}
		appmanifest.Processes = &processes
	}
	if len(appType.Sidecars) != 0 {
		var sidecars cfv3operation.AppManifestSideCars
		for _, sidecar := range appType.Sidecars {
			var processTypes []string
			diags = sidecar.ProcessTypes.ElementsAs(ctx, &processTypes, false)
			sidecars = append(sidecars, cfv3operation.AppManifestSideCar{
				Name:         sidecar.Name.ValueString(),
				Command:      sidecar.Command.ValueString(),
				ProcessTypes: processTypes,
				Memory:       sidecar.Memory.ValueString(),
			})
		}
		appmanifest.Sidecars = &sidecars
	}
	if !appType.Command.IsNull() {
		appmanifest.Command = appType.Command.ValueString()
	}
	if !appType.DiskQuota.IsNull() {
		appmanifest.DiskQuota = appType.DiskQuota.ValueString()
	}
	if !appType.HealthCheckHttpEndpoint.IsNull() {
		appmanifest.HealthCheckHTTPEndpoint = appType.HealthCheckHttpEndpoint.ValueString()
	}
	if !appType.HealthCheckInvocationTimeout.IsNull() {
		appmanifest.HealthCheckInvocationTimeout = uint(appType.HealthCheckInvocationTimeout.ValueInt64())
	}
	if !appType.HealthCheckType.IsNull() {
		appmanifest.HealthCheckType = cfv3operation.AppHealthCheckType(appType.HealthCheckType.ValueString())
	}
	if !appType.Instances.IsNull() {
		appmanifest.Instances = uint(appType.Instances.ValueInt64())
	}
	if !appType.Memory.IsNull() {
		appmanifest.Memory = appType.Memory.ValueString()
	}
	if !appType.Timeout.IsNull() {
		appmanifest.Timeout = uint(appType.Timeout.ValueInt64())
	}
	return &appmanifest, diags
}

// mapAppValuesToType function maps cfv3resource manifest type to AppType
func mapAppValuesToType(ctx context.Context, appManifest *cfv3operation.AppManifest, app *cfv3resource.App) (AppType, diag.Diagnostics) {
	var diags, tempDiags diag.Diagnostics
	var appType AppType
	appType.Name = types.StringValue(appManifest.Name)
	appType.Stack = types.StringValue(appManifest.Stack)
	if appManifest.Buildpacks != nil {
		var buildpacks []types.String
		for _, buildpack := range appManifest.Buildpacks {
			buildpacks = append(buildpacks, types.StringValue(buildpack))
		}
		appType.Buildpacks, tempDiags = types.SetValueFrom(ctx, types.StringType, buildpacks)
		diags = append(diags, tempDiags...)
	}
	if appManifest.Docker.Image != "" {
		appType.DockerImage = types.StringValue(appManifest.Docker.Image)
		if appManifest.Docker.Username != "" {
			appType.DockerCredentials.Username = types.StringValue(appManifest.Docker.Username)
		}
	}
	if appManifest.Services != nil {
		var serviceBindings []ServiceBinding
		for _, service := range *appManifest.Services {
			var sb ServiceBinding
			sb.ServiceInstance = types.StringValue(service.Name)
			if service.Parameters != nil {
				sb.Param, tempDiags = types.MapValueFrom(ctx, types.MapType{ElemType: types.StringType}, service.Parameters)
				diags = append(diags, tempDiags...)
			}
			serviceBindings = append(serviceBindings, sb)
		}
		appType.ServiceBindings = serviceBindings
	}
	if appManifest.Routes != nil {
		var routes []Route
		for _, route := range *appManifest.Routes {
			var r Route
			r.Route = types.StringValue(route.Route)
			if route.Protocol != "" {
				r.Protocol = types.StringValue(string(route.Protocol))
			}
			routes = append(routes, r)
		}
		appType.Routes = routes
	}
	if appManifest.Env != nil {
		appType.Environment, tempDiags = types.MapValueFrom(ctx, types.MapType{ElemType: types.StringType}, appManifest.Env)
		diags = append(diags, tempDiags...)
	}
	if appManifest.HealthCheckInterval != 0 {
		appType.HealthCheckInterval = types.Int64Value(int64(appManifest.HealthCheckInterval))
	}
	if appManifest.ReadinessHealthCheckType != "" {
		appType.ReadinessHealthCheckType = types.StringValue(appManifest.ReadinessHealthCheckType)
	}
	if appManifest.ReadinessHealthCheckHttpEndpoint != "" {
		appType.ReadinessHealthCheckHttpEndpoint = types.StringValue(appManifest.ReadinessHealthCheckHttpEndpoint)
	}
	if appManifest.ReadinessHealthInvocationTimeout != 0 {
		appType.ReadinessHealthCheckInvocationTimeout = types.Int64Value(int64(appManifest.ReadinessHealthInvocationTimeout))
	}
	if appManifest.ReadinessHealthCheckInterval != 0 {
		appType.ReadinessHealthCheckInterval = types.Int64Value(int64(appManifest.ReadinessHealthCheckInterval))
	}
	if appManifest.LogRateLimitPerSecond != "" {
		appType.LogRateLimitPerSecond = types.StringValue(appManifest.LogRateLimitPerSecond)
	}
	if appManifest.NoRoute {
		appType.NoRoute = types.BoolValue(appManifest.NoRoute)
	}
	if appManifest.RandomRoute {
		appType.RandomRoute = types.BoolValue(appManifest.RandomRoute)
	}
	if appManifest.Processes != nil {
		var processes []Process
		for _, process := range *appManifest.Processes {
			var p Process
			p.Type = types.StringValue(string(process.Type))
			p.Command = types.StringValue(process.Command)
			p.DiskQuota = types.StringValue(process.DiskQuota)
			p.HealthCheckHttpEndpoint = types.StringValue(process.HealthCheckHTTPEndpoint)
			p.HealthCheckInvocationTimeout = types.Int64Value(int64(process.HealthCheckInvocationTimeout))
			p.HealthCheckType = types.StringValue(string(process.HealthCheckType))
			p.Instances = types.Int64Value(int64(process.Instances))
			p.Memory = types.StringValue(process.Memory)
			p.Timeout = types.Int64Value(int64(process.Timeout))
			processes = append(processes, p)
		}
		appType.Processes = processes
	}
	if appManifest.Sidecars != nil {
		var sidecars []Sidecar
		for _, sidecar := range *appManifest.Sidecars {
			var s Sidecar
			s.Name = types.StringValue(sidecar.Name)
			s.Command = types.StringValue(sidecar.Command)
			var processTypes []types.String
			for _, processType := range sidecar.ProcessTypes {
				processTypes = append(processTypes, types.StringValue(processType))
			}
			s.ProcessTypes, tempDiags = types.SetValueFrom(ctx, types.StringType, processTypes)
			diags = append(diags, tempDiags...)
			s.Memory = types.StringValue(sidecar.Memory)
			sidecars = append(sidecars, s)
		}
		appType.Sidecars = sidecars
	}
	appType.ID = types.StringValue(app.GUID)
	appType.CreatedAt = types.StringValue(app.CreatedAt.Format(time.RFC3339))
	appType.UpdatedAt = types.StringValue(app.UpdatedAt.Format(time.RFC3339))
	appType.Command = types.StringValue(appManifest.Command)
	appType.DiskQuota = types.StringValue(appManifest.DiskQuota)
	appType.HealthCheckHttpEndpoint = types.StringValue(appManifest.HealthCheckHTTPEndpoint)
	appType.HealthCheckInvocationTimeout = types.Int64Value(int64(appManifest.HealthCheckInvocationTimeout))
	appType.HealthCheckType = types.StringValue(string(appManifest.HealthCheckType))
	appType.Instances = types.Int64Value(int64(appManifest.Instances))
	appType.Memory = types.StringValue(appManifest.Memory)
	appType.Timeout = types.Int64Value(int64(appManifest.Timeout))
	return appType, diags
}
