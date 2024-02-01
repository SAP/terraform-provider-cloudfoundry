package provider

import (
	"bytes"
	"testing"
	"text/template"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func hclResourceSpace(smp *SpaceModelPtr) string {
	if smp != nil {
		s := `
			resource "cloudfoundry_space" "ds" {
			{{- if .Name}}
				name = "{{.Name}}"
			{{- end -}}
			{{if .Id}}
				id = "{{.Id}}"
			{{- end -}}
			{{if .OrgId}}
				org = "{{.OrgId}}"
			{{- end -}}
			{{if .Quota}}
				quota = "{{.Quota}}"
			{{- end -}}
			{{if .AllowSSH}}
				allow_ssh = "{{.AllowSSH}}"
			{{- end -}}
			{{if .IsolationSegment}}
				isolation_segment = "{{.IsolationSegment}}"
			{{- end -}}
			{{if .CreatedAt}}
				created_at = "{{.CreatedAt}}"
			{{- end -}}
			{{if .UpdatedAt}}
				updated_at = "{{.UpdatedAt}}"
			{{- end -}}
			{{if .Labels}}
				labels = {{.Labels}}
			{{- end -}}
			{{if .Annotations}}
				annotations = {{.Annotations}}
			{{- end }}
			}`
		tmpl, err := template.New("resource_space").Parse(s)
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
	return `resource "cloudfoundry_space" "ds" {}`
}

func TestSpaceResource_Configure(t *testing.T) {
	t.Parallel()
	resourceName := "cloudfoundry_space.ds"
	t.Run("happy path - create/read/update/delete/import space", func(t *testing.T) {
		cfg := getCFHomeConf()
		rec := cfg.SetupVCR(t, "fixtures/resource_space_crud")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider(nil) + hclResourceSpace(&SpaceModelPtr{
						Name:             strtostrptr(testSpaceResourceName),
						OrgId:            strtostrptr(testOrgGUID),
						AllowSSH:         booltoboolptr(true),
						Labels:           strtostrptr(testSpaceResourceCreateLabel),
						IsolationSegment: strtostrptr(testIsolationSegmentGUID),
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr(resourceName, "id", regexpValidUUID),
						resource.TestCheckNoResourceAttr(resourceName, "quota"),
						resource.TestCheckResourceAttr(resourceName, "allow_ssh", "true"),
						resource.TestCheckResourceAttr(resourceName, "labels.purpose", "testing"),
						resource.TestCheckResourceAttr(resourceName, "isolation_segment", testIsolationSegmentGUID),
					),
				},
				{
					Config: hclProvider(nil) + hclResourceSpace(&SpaceModelPtr{
						Name:     strtostrptr(testSpaceResourceName),
						OrgId:    strtostrptr(testOrgGUID),
						AllowSSH: booltoboolptr(false),
						Labels:   strtostrptr(testSpaceResourceUpdateLabel),
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "allow_ssh", "false"),
						resource.TestCheckResourceAttr(resourceName, "labels.purpose", "production"),
						resource.TestCheckResourceAttr(resourceName, "labels.%", "2"),
						resource.TestCheckNoResourceAttr(resourceName, "isolation_segment"),
					),
				},
				{
					ResourceName:      resourceName,
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	})
}
