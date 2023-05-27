
package youtrack

import (

	"log"
	"net/http"
	"net/url"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

)
type Api struct {
	ApiVersion  string
	// BaseURL is the URL to the REST API endpoint for a YouTrack Project. It should
	// end is a slash. For example: https://goyt.myjetbrains.com/youtrack/api/
	BaseURL *url.URL

	// Token is the permanent token used to make authenticated requests.
	// For more information, see:
	// https://www.jetbrains.com/help/youtrack/incloud/authentication-with-permanent-token.html
	Token string

	// EnableTracing turns on extra logging, including HTTP request/response logging.
	// NOTE that the authorization token will be logged when this is enabled.
	EnableTracing bool
}
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_version": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"base_url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("Base_url", ""),
			},
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SERVICE_TOKEN", ""),
			},
			// "EnableTracing" {
			// 	Type:        bool,
			// 	Optional: true,
			// 	Default:  true,
			// },
		},
		DataSourcesMap: map[string]*schema.Resource{
			"name_users": getDetailsForUsersSchema(),
		},
		ResourcesMap: map[string]*schema.Resource{
			// "example_item": resourceItem(),
		},
		ConfigureFunc: providerConfigure,
	}
}


func (api *Api) trace(v ...interface{}) {
	if api.EnableTracing {
		log.Println(v...)
	}
}
type ProviderClient struct {
	ApiVersion  string
	Hostname string
	Client      *Client
}

func newProviderClient(apiVersion, hostname string, headers http.Header) (ProviderClient, error) {
	p := ProviderClient{
		ApiVersion: apiVersion,
		Hostname:   hostname,

	}
	p.Client = NewClient(headers, hostname, apiVersion)

	return p, nil
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	apiVersion := d.Get("api_version").(string)
	if apiVersion == "" {
		log.Println("Defaulting environment in URL config to use API default version...")
	}

	hostname := d.Get("base_url").(string)
	if hostname == "" {
		log.Println("Defaulting environment in URL config to use API default hostname...")
		hostname = "localhost"
	}
	token := d.Get("token").(string)

	h := make(http.Header)
	h.Set("Content-Type", "application/json")
	h.Set("Accept", "application/json")
	h.Set("Authorization", "Bearer "+token)

	headers, exists := d.GetOk("headers")
	if exists {
		for k, v := range headers.(map[string]interface{}) {
			h.Set(k, v.(string))
		}
	}

	return newProviderClient(apiVersion, hostname, h)
}

func marshalData(d *schema.ResourceData, vals map[string]interface{}) {
	for k, v := range vals {
		if k == "id" {
			d.SetId(v.(string))
		} else {
			str, ok := v.(string)
			if ok {
				d.Set(k, str)
			} else {
				d.Set(k, v)
			}
		}
	}
}

