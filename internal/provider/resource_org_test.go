package provider

import (
	"bytes"
	"testing"
	"text/template"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

type OrgResourceModelPtr struct {
	Name        *string
	Labels      *map[string]string
	Annotations *map[string]string
	Suspended   *bool
}

func hclResourceOrg(ormp *OrgResourceModelPtr) string {
	if ormp != nil {
		s := `
			resource "cloudfoundry_org" "test" {
				{{if .Name}}
					name = "{{.Name}}"
				{{- end }}
				{{if .Labels}}
					labels = "{{.Labels}}"
				{{- end }}
				{{if .Annotations}}
					annotations = "{{.Annotations}}"
				{{- end }}
				{{ if .Suspended }}
					suspended  = "{{.Suspended}}"
				{{- end }}
			}
			`
		tmpl, err := template.New("resource_cf_org").Parse(s)
		if err != nil {
			panic(err)
		}
		buf := new(bytes.Buffer)
		err = tmpl.Execute(buf, ormp)
		if err != nil {
			panic(err)
		}
		return buf.String()
	}
	return `resource "cloudfoundry_org" "test" {}`
}

func TestResourceOrg(t *testing.T) {
	t.Parallel()
	resourceName := "cloudfoundry_org.test"
	t.Run("happy path - create org", func(t *testing.T) {
		cfg := getCFHomeConf()
		rec := cfg.SetupVCR(t, "fixtures/resource_org")
		defer stopQuietly(rec)
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider(nil) + hclResourceOrg(&OrgResourceModelPtr{
						Name: strtostrptr("tf-unit-test"),
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr(resourceName, "id", regexpValidUUID),
						resource.TestMatchResourceAttr(resourceName, "created_at", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr(resourceName, "updated_at", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr(resourceName, "quota", regexpValidUUID),
					),
					//Destroy: true,
					// ImportState:       true,
					// ImportStateVerify: true,
				},
			},
		})
	})

}
