package provider

import (
	"context"
	"net/url"
	"reflect"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type config struct {
	URL *url.URL
}

// New returns a new Terraform provider.
func New() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Description: "The URL to access Gigo Core.",
				Optional:    true,
				// The "GIGO_AGENT_URL" environment variable is used by default
				// as the Access URL when generating scripts.
				DefaultFunc: schema.EnvDefaultFunc("GIGO_AGENT_URL", "https://gigo.dev"),
				ValidateFunc: func(i interface{}, s string) ([]string, []error) {
					_, err := url.Parse(s)
					if err != nil {
						return nil, []error{err}
					}
					return nil, nil
				},
			},
		},
		ConfigureContextFunc: func(c context.Context, resourceData *schema.ResourceData) (interface{}, diag.Diagnostics) {
			rawURL, ok := resourceData.Get("url").(string)
			if !ok {
				return nil, diag.Errorf("unexpected type %q for url", reflect.TypeOf(resourceData.Get("url")).String())
			}
			if rawURL == "" {
				return nil, diag.Errorf("GIGO_AGENT_URL must not be empty; got %q", rawURL)
			}
			parsed, err := url.Parse(resourceData.Get("url").(string))
			if err != nil {
				return nil, diag.FromErr(err)
			}
			rawHost, ok := resourceData.Get("host").(string)
			if ok && rawHost != "" {
				rawPort := parsed.Port()
				if rawPort != "" && !strings.Contains(rawHost, ":") {
					rawHost += ":" + rawPort
				}
				parsed.Host = rawHost
			}
			return config{
				URL: parsed,
			}, nil
		},
		DataSourcesMap: map[string]*schema.Resource{
			"gigo_workspace":   workspaceDataSource(),
			"gigo_provisioner": provisionerDataSource(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"gigo_agent": agentResource(),
		},
	}
}

// valueAsString takes a cty.Value that may be a string or null, and converts it to a Go string,
// which will be empty if the input value was null.
// or a nil interface{}
func valueAsString(value cty.Value) string {
	if value.IsNull() {
		return ""
	}
	return value.AsString()
}

// valueAsString takes a cty.Value that may be a boolean or null, and converts it to either a Go bool
// or a nil interface{}
func valueAsBool(value cty.Value) interface{} {
	if value.IsNull() {
		return nil
	}
	return value.True()
}
