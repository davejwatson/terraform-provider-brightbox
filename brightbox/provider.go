package brightbox

import (
	"fmt"
	"log"
	"strings"

	"github.com/brightbox/gobrightbox"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

const (
	// Terraform application client credentials
	defaultClientID     = "app-dkmch"
	defaultClientSecret = "uogoelzgt0nwawb"
	appPrefix           = "app-"
	passwordEnvVar      = "BRIGHTBOX_PASSWORD"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"apiclient": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("BRIGHTBOX_CLIENT", defaultClientID),
				Description: "Brightbox Cloud API Client/OAuth Application ID",
			},
			"apisecret": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("BRIGHTBOX_CLIENT_SECRET", defaultClientSecret),
				Description: "Brightbox Cloud API Client/OAuth Application Secret",
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("BRIGHTBOX_USER_NAME", nil),
				Description: "Brightbox Cloud User Name",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc(passwordEnvVar, nil),
				Description: "Brightbox Cloud Password for User Name",
			},
			"account": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("BRIGHTBOX_ACCOUNT", nil),
				Description: "Brightbox Cloud Account to operate on",
			},
			"apiurl": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("BRIGHTBOX_API_URL", brightbox.DefaultRegionApiURL),
				Description: "Brightbox Cloud Api URL for selected Region",
			},
			"orbit_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("BRIGHTBOX_ORBIT_URL", brightbox.DefaultOrbitAuthURL),
				Description: "Brightbox Cloud Orbit URL for selected Region",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"brightbox_image":         dataSourceBrightboxImage(),
			"brightbox_database_type": dataSourceBrightboxDatabaseType(),
			"brightbox_server_group":  dataSourceBrightboxServerGroup(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"brightbox_server":          resourceBrightboxServer(),
			"brightbox_cloudip":         resourceBrightboxCloudip(),
			"brightbox_server_group":    resourceBrightboxServerGroup(),
			"brightbox_firewall_policy": resourceBrightboxFirewallPolicy(),
			"brightbox_firewall_rule":   resourceBrightboxFirewallRule(),
			"brightbox_load_balancer":   resourceBrightboxLoadBalancer(),
			"brightbox_database_server": resourceBrightboxDatabaseServer(),
			"brightbox_orbit_container": resourceBrightboxContainer(),
			"brightbox_api_client":      resourceBrightboxApiClient(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := &authdetails{
		APIClient: d.Get("apiclient").(string),
		APISecret: d.Get("apisecret").(string),
		UserName:  d.Get("username").(string),
		password:  d.Get("password").(string),
		Account:   d.Get("account").(string),
		APIURL:    d.Get("apiurl").(string),
		OrbitUrl:  d.Get("orbit_url").(string),
	}

	if strings.HasPrefix(config.APIClient, appPrefix) {
		log.Printf("[DEBUG] Detected OAuth Application. Validating User details.")
		if config.UserName == "" || config.password == "" {
			return nil,
				fmt.Errorf("User Credentials are missing. Please supply a Username and One Time Authentication code.")
		}
		if config.Account == "" {
			return nil,
				fmt.Errorf("Must specify Account with User Credentials")
		}
	} else {
		log.Printf("[DEBUG] Detected API Client.")
		if config.UserName != "" || config.password != "" {
			return nil,
				fmt.Errorf("User Credentials should be blank with an API Client")
		}
	}

	return config.Client()
}
