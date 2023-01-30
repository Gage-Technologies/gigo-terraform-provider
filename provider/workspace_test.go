package provider_test

import (
	"testing"

	"github.com/gage-technologies/gigo-terraform-provider/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/require"
)

func TestWorkspace(t *testing.T) {
	t.Setenv("GIGO_WORKSPACE_OWNER", "owner123")
	t.Setenv("GIGO_WORKSPACE_OWNER_EMAIL", "owner123@example.com")
	t.Setenv("GIGO_WORKSPACE_DISK", "50Gi")
	t.Setenv("GIGO_WORKSPACE_CPU", "8")
	t.Setenv("GIGO_WORKSPACE_MEM", "16G")

	resource.Test(t, resource.TestCase{
		Providers: map[string]*schema.Provider{
			"gigo": provider.New(),
		},
		IsUnitTest: true,
		Steps: []resource.TestStep{{
			Config: `
			provider "gigo" {
				url = "https://example.com:8080"
			}
			data "gigo_workspace" "me" {
			}`,
			Check: func(state *terraform.State) error {
				require.Len(t, state.Modules, 1)
				require.Len(t, state.Modules[0].Resources, 1)
				resource := state.Modules[0].Resources["data.gigo_workspace.me"]
				require.NotNil(t, resource)

				attribs := resource.Primary.Attributes
				value := attribs["transition"]
				require.NotNil(t, value)
				t.Log(value)
				require.Equal(t, "8080", attribs["access_port"])
				require.Equal(t, "owner123", attribs["owner"])
				require.Equal(t, "owner123@example.com", attribs["owner_email"])
				require.Equal(t, "50Gi", attribs["disk"])
				require.Equal(t, "8", attribs["cpu"])
				require.Equal(t, "16G", attribs["mem"])
				return nil
			},
		}},
	})
	resource.Test(t, resource.TestCase{
		Providers: map[string]*schema.Provider{
			"gigo": provider.New(),
		},
		IsUnitTest: true,
		Steps: []resource.TestStep{{
			Config: `
			provider "gigo" {
				url = "https://example.com:8080"
			}
			data "gigo_workspace" "me" {
			}`,
			Check: func(state *terraform.State) error {
				require.Len(t, state.Modules, 1)
				require.Len(t, state.Modules[0].Resources, 1)
				resource := state.Modules[0].Resources["data.gigo_workspace.me"]
				require.NotNil(t, resource)

				attribs := resource.Primary.Attributes
				value := attribs["transition"]
				require.NotNil(t, value)
				t.Log(value)
				require.Equal(t, "https://example.com:8080", attribs["access_url"])
				require.Equal(t, "owner123", attribs["owner"])
				require.Equal(t, "owner123@example.com", attribs["owner_email"])
				require.Equal(t, "50Gi", attribs["disk"])
				require.Equal(t, "8", attribs["cpu"])
				require.Equal(t, "16G", attribs["mem"])
				return nil
			},
		}},
	})
}
