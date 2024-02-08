package provider

import (
	"bytes"
	"regexp"
	"testing"
	"text/template"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

type UserResourceModelPtr struct {
	HclType          string
	HclObjectName    string
	ObjectName       string
	UserName         *string
	PresentationName *string
	Origin           *string
	Id               *string
	Labels           *string
	Annotations      *string
	CreatedAt        *string
	UpdatedAt        *string
}

func hclResourceUser(urmp *UserResourceModelPtr) string {
	if urmp != nil {
		s := `
		{{.HclType}} "cloudfoundry_user" {{.HclObjectName}} {
			{{- if .UserName}}
				username = "{{.UserName}}"
			{{- end -}}
			{{if .Id}}
				id = "{{.Id}}"
			{{- end -}}
			{{if .PresentationName}}
				presentation_name = "{{.PresentationName}}"
			{{- end -}}
			{{if .Origin}}
				origin = "{{.Origin}}"
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
		tmpl, err := template.New("resource_user").Parse(s)
		if err != nil {
			panic(err)
		}
		buf := new(bytes.Buffer)
		err = tmpl.Execute(buf, urmp)
		if err != nil {
			panic(err)
		}
		return buf.String()
	}
	return urmp.HclType + ` "cloudfoundry_user "` + urmp.HclObjectName + ` {}`
}

func TestUserResource_Configure(t *testing.T) {
	t.Parallel()
	t.Run("happy path - create/update/import/delete user", func(t *testing.T) {
		resourceName := "cloudfoundry_user.us"
		cfg := getCFHomeConf()
		rec := cfg.SetupVCR(t, "fixtures/resource_user_crud")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider(nil) + hclResourceUser(&UserResourceModelPtr{
						HclType:       hclObjectResource,
						HclObjectName: "us",
						Id:            strtostrptr("tf-test"),
						Labels:        strtostrptr(testCreateLabel),
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckNoResourceAttr(resourceName, "annotations"),
						resource.TestCheckResourceAttr(resourceName, "presentation_name", "tf-test"),
						resource.TestCheckResourceAttr(resourceName, "labels.purpose", "testing"),
					),
				},
				{
					Config: hclProvider(nil) + hclResourceUser(&UserResourceModelPtr{
						HclType:       hclObjectResource,
						HclObjectName: "us",
						Id:            strtostrptr("tf-test"),
						Labels:        strtostrptr(testUpdateLabel),
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "labels.purpose", "production"),
						resource.TestCheckResourceAttr(resourceName, "labels.%", "2"),
					),
				},
				{
					ResourceName:      resourceName,
					ImportStateIdFunc: getIdForImport(resourceName),
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	})
	t.Run("error path - create user with existing id", func(t *testing.T) {
		cfg := getCFHomeConf()
		rec := cfg.SetupVCR(t, "fixtures/resource_user_invalid")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider(nil) + hclResourceUser(&UserResourceModelPtr{
						HclType:       hclObjectResource,
						HclObjectName: "rs_invalid",
						Id:            strtostrptr(testUserGUID),
						Labels:        strtostrptr(testCreateLabel),
					}),
					ExpectError: regexp.MustCompile(`API Error Registering User`),
				},
			},
		})
	})
}
