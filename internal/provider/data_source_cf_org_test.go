package provider

import (
	"bytes"
	"regexp"
	"testing"
	"text/template"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

type OrgDataSourceModelPtr struct {
	Name        *string
	Id          *string
	Labels      *map[string]string
	Annotations *map[string]string
}

func hclDataSourceOrg(odsmp *OrgDataSourceModelPtr) string {
	if odsmp != nil {
		s := `
			data "cloudfoundry_org" "ds" {
			{{- if .Name}}
				name  = "{{.Name}}"
			{{- end -}}
			{{if .Id}}
				id = "{{.Id}}"
			{{- end -}}
			{{if .Labels}}
				labels = "{{.Labels}}"
			{{- end -}}
			{{if .Annotations}}
				annotations = "{{.Annotations}}"
			{{- end }}
			}`
		tmpl, err := template.New("datasource_org").Parse(s)
		if err != nil {
			panic(err)
		}
		buf := new(bytes.Buffer)
		err = tmpl.Execute(buf, odsmp)
		if err != nil {
			panic(err)
		}
		return buf.String()
	}
	return `data "cloudfoundry_org" "ds" {}`
}
func TestOrgDataSource_Configure(t *testing.T) {
	t.Parallel()
	t.Run("error path - get unavailable datasource org", func(t *testing.T) {
		cfg := getCFHomeConf()
		rec := cfg.SetupVCR(t, "fixtures/datasource_org_invalid_orgname.yaml")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider(nil) + hclDataSourceOrg(&OrgDataSourceModelPtr{
						Name: strtostrptr("testunavailableorg"),
					}),
					ExpectError: regexp.MustCompile(`Error: Unable to find org data in list`),
				},
			},
		})
	})
	t.Run("get available datasource org", func(t *testing.T) {
		cfg := getCFHomeConf()
		rec := cfg.SetupVCR(t, "fixtures/datasource_org.yaml")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider(nil) + hclDataSourceOrg(&OrgDataSourceModelPtr{
						Name: strtostrptr("PerformanceTeamBLR"),
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("data.cloudfoundry_org.ds", "id", regexpValidUUID),
						resource.TestCheckResourceAttr("data.cloudfoundry_org.ds", labelsKey+".env", "canary"),
						resource.TestCheckResourceAttr("data.cloudfoundry_org.ds", annotationsKey+".%", "0"),
					),
				},
			},
		})
	})
}
