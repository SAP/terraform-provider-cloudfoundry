package provider

import (
	"bytes"
	"regexp"
	"testing"
	"text/template"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

type ServiceModelPtr struct {
	HclType       string
	HclObjectName string
	Name          *string
	Id            *string
	ServiceBroker *string
	ServicePlans  *string
}

func hclService(smp *ServiceModelPtr) string {
	if smp != nil {
		s := `
		{{.HclType}} "cloudfoundry_service" {{.HclObjectName}} {
			{{- if .Name}}
				name = "{{.Name}}"
			{{- end }}
			{{if .Id}}
				id = "{{.Id}}"
			{{- end -}}
			{{if .ServiceBroker}}
				service_broker = "{{.ServiceBroker}}"
			{{- end -}}
			{{if .ServicePlans}}
				service_plans = {{.ServicePlans}}
			{{- end -}}
		}`
		tmpl, err := template.New("service").Parse(s)
		if err != nil {
			panic(err)
		}
		buf := new(bytes.Buffer)
		err = tmpl.Execute(buf, smp)
		if err != nil {
			panic(err)
		}
		return buf.String()
	}
	return smp.HclType + ` "cloudfoundry_service" ` + smp.HclObjectName + ` {}`
}

func TestDatasourceService(t *testing.T) {
	t.Parallel()
	datasourceName := "data.cloudfoundry_service.test"
	t.Run("error path - get unavailable  service", func(t *testing.T) {
		cfg := getCFHomeConf()
		rec := cfg.SetupVCR(t, "fixtures/datasource_service_invalid_servicename")
		defer stopQuietly(rec)

		// Create a Terraform configuration that uses the data source
		// and run `terraform apply`. The data source should not be found.
		resource.UnitTest(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider(nil) + hclService(&ServiceModelPtr{
						HclType:       hclObjectDataSource,
						HclObjectName: "test",
						Name:          strtostrptr("invalid-service-name"),
					}),
					ExpectError: regexp.MustCompile(`Error: Unable to find service offering in list`),
				},
			},
		})

	})

	t.Run("happy path - read service", func(t *testing.T) {
		cfg := getCFHomeConf()
		rec := cfg.SetupVCR(t, "fixtures/datasource_service")
		defer stopQuietly(rec)
		testServiceName := "xsuaa"                                // Canary service name
		testServiceGUID := "c1876449-7493-43dc-a36b-0a8215fa46ad" // Canary service GUID
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider(nil) + hclService(&ServiceModelPtr{
						HclType:       hclObjectDataSource,
						HclObjectName: "test",
						Name:          strtostrptr(testServiceName),
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(datasourceName, "name", testServiceName),
						resource.TestCheckResourceAttr(datasourceName, "id", testServiceGUID),
						resource.TestCheckResourceAttr(datasourceName, "service_broker", "xsuaa"),
						resource.TestCheckResourceAttr(datasourceName, "service_plans.%", "5"),
						resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(datasourceName, "service_plans.apiaccess", "a4545ae8-40ed-4b2b-9cab-170cb17e4e6f"),
							resource.TestCheckResourceAttr(datasourceName, "service_plans.application", "432bd9db-20e2-4997-825f-e4a937705b87"),
							resource.TestCheckResourceAttr(datasourceName, "service_plans.broker", "11ced9f6-e4f2-466f-b86c-6b19a62e0c11"),
							resource.TestCheckResourceAttr(datasourceName, "service_plans.space", "4c15e9a6-6cbc-4d58-9243-273f57aca7d5"),
							resource.TestCheckResourceAttr(datasourceName, "service_plans.system", "0a140ebc-9d31-4114-b596-8ce920cd9ec7"),
						),
					),
				},
			},
		})
	})

}
