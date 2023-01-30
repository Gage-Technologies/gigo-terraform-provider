package provider

import (
	"context"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"os"
	"reflect"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func agentResource() *schema.Resource {
	return &schema.Resource{
		Description: "Use this resource to associate an agent.",
		CreateContext: func(c context.Context, resourceData *schema.ResourceData, i interface{}) diag.Diagnostics {
			// create a new snowflake node using 1022 as the node id
			// this is a reserved node id in the gigo system for provisioners
			sf, err := snowflake.NewNode(1022)
			if err != nil {
				return diag.FromErr(err)
			}

			// generate a new id from the snowflake node
			// NOTE: we switched from uuid to snowflake to maintain
			// compatibility with the gigo systems existing id system
			resourceData.SetId(sf.Generate().String())

			// This should be a real authentication token!
			err = resourceData.Set("token", uuid.NewString())
			if err != nil {
				return diag.FromErr(err)
			}
			return updateInitScript(resourceData, i)
		},
		ReadWithoutTimeout: func(c context.Context, resourceData *schema.ResourceData, i interface{}) diag.Diagnostics {
			err := resourceData.Set("token", uuid.NewString())
			if err != nil {
				return diag.FromErr(err)
			}
			return updateInitScript(resourceData, i)
		},
		DeleteContext: func(c context.Context, rd *schema.ResourceData, i interface{}) diag.Diagnostics {
			return nil
		},
		Schema: map[string]*schema.Schema{
			"init_script": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Run this script on startup of an instance to initialize the agent.",
			},
			"arch": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				Description:  `The architecture the agent will run on. Must be one of: "amd64", "armv7", "arm64".`,
				ValidateFunc: validation.StringInSlice([]string{"amd64", "armv7", "arm64"}, false),
			},
			"os": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				Description:  `The operating system the agent will run on. Must be one of: "linux", "darwin", or "windows".`,
				ValidateFunc: validation.StringInSlice([]string{"linux", "darwin", "windows"}, false),
			},
			"token": {
				ForceNew:    true,
				Sensitive:   true,
				Description: `Set the environment variable "GIGO_AGENT_TOKEN" with this token to authenticate an agent.`,
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

// updateInitScript fetches parameters from a "gigo_agent" to produce the
// agent script from environment variables.
func updateInitScript(resourceData *schema.ResourceData, i interface{}) diag.Diagnostics {
	config, valid := i.(config)
	if !valid {
		return diag.Errorf("config was unexpected type %q", reflect.TypeOf(i).String())
	}
	operatingSystem, valid := resourceData.Get("os").(string)
	if !valid {
		return diag.Errorf("os was unexpected type %q", reflect.TypeOf(resourceData.Get("os")))
	}
	arch, valid := resourceData.Get("arch").(string)
	if !valid {
		return diag.Errorf("arch was unexpected type %q", reflect.TypeOf(resourceData.Get("arch")))
	}
	accessURL, err := config.URL.Parse("/")
	if err != nil {
		return diag.Errorf("parse access url: %s", err)
	}
	script := os.Getenv(fmt.Sprintf("GIGO_AGENT_SCRIPT_%s_%s", operatingSystem, arch))
	if script != "" {
		script = strings.ReplaceAll(script, "${ACCESS_URL}", accessURL.String())
	}
	err = resourceData.Set("init_script", script)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
