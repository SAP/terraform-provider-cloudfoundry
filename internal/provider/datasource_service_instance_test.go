package provider

import (
	"bytes"
	"testing"
	"text/template"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

type DataSourceServiceInstanceModelPtr struct {
	HclType       string
	HclObjectName string
	Name          *string
	Id            *string
	Labels        *string
	Annotations   *string
	Type          *string
	Space         *string
	ServicePlan   *string
	Parameters    *string
	Credentials   *string
	Tags          *string
}

func hclServiceInstance(sip *DataSourceServiceInstanceModelPtr) string {
	if sip != nil {
		s := `
		{{.HclType}} "cloudfoundry_service_instance" {{.HclObjectName}} {
			{{- if .Name}}
				name  = "{{.Name}}"
			{{- end -}}
			{{if .Id}}
				id = "{{.Id}}"
			{{- end -}}
			{{if .Labels}}
				labels = {{.Labels}}
			{{- end -}}
			{{if .Annotations}}
				annotations = {{.Annotations}}
			{{- end -}}
			{{if .Type}}
				type = "{{.Type}}"
			{{- end -}}
			{{if .Space}}
				space = "{{.Space}}"
			{{- end }}
			{{if .ServicePlan}}
				service_plan = "{{.ServicePlan}}"
			{{- end -}}
			{{if .Parameters}}
				parameters = {{.Parameters}}
			{{- end -}}
			{{if .Credentials}}
				credentials = {{.Credentials}}
			{{- end -}}
			{{if .Tags}}
				tags = {{.Tags}}
			{{- end -}}
		}`
		tmpl, err := template.New("service_instance").Parse(s)
		if err != nil {
			panic(err)
		}
		buf := new(bytes.Buffer)
		err = tmpl.Execute(buf, sip)
		if err != nil {
			panic(err)
		}
		return buf.String()
	}
	return sip.HclType + ` "cloudfoundry_service_instance"  "` + sip.HclObjectName + ` {}`
}

func TestServiceInstanceDataSource(t *testing.T) {
	var (
		// in canary -> PerformanceTeamBLR -> tf-space-1
		testSpaceGUID           = "02c0cc92-6ecc-44b1-b7b2-096ca19ee143"
		testServiceInstance     = "tf-test-do-not-delete"
		testServiceInstanceGUID = "5e2976bb-332e-41e1-8be3-53baafea9296" // in canary -> PerformanceTeamBLR -> tf-space-1
	)
	t.Parallel()
	t.Run("happy path - read service instance user-provided", func(t *testing.T) {
		cfg := getCFHomeConf()
		dataSourceName := "data.cloudfoundry_service_instance.ds"
		rec := cfg.SetupVCR(t, "fixtures/datasource_service_instance")
		defer stopQuietly(rec)
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider(nil) + hclServiceInstance(&DataSourceServiceInstanceModelPtr{
						HclType:       hclObjectDataSource,
						HclObjectName: "ds",
						Name:          strtostrptr(testServiceInstance),
						Space:         strtostrptr(testSpaceGUID),
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(dataSourceName, "id", testServiceInstanceGUID),
						resource.TestCheckResourceAttr(dataSourceName, "name", testServiceInstance),
						resource.TestCheckResourceAttr(dataSourceName, "type", "user-provided"),
						resource.TestCheckResourceAttr(dataSourceName, "space", testSpaceGUID),
					),
				},
			},
		})
	})
}
