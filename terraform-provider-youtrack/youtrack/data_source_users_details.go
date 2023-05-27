package youtrack

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func getDetailsForUsersSchema() *schema.Resource {
	return &schema.Resource{
		Read: getUsersDataSourceRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"login": {
				Type:     schema.TypeString,
				Required: true,
				Elem: schema.TypeString,
			},
			"full_name": {
				Type:     schema.TypeString,
				Optional: true,
				Elem: schema.TypeString,
			},
			"email": {
				Type:     schema.TypeString,
				Optional: true,
				Elem: schema.TypeString,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Elem: schema.TypeString,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Elem: schema.TypeString,
			},
		},
	}
}

// dataSourceRead tells Terraform how to contact our microservice and retrieve the necessary data
func getUsersDataSourceRead(d *schema.ResourceData, meta interface{}) (err error) {
	provider := meta.(ProviderClient)
	client := provider.Client

	header := make(http.Header)
	headers, exists := d.GetOk("headers")
	if exists {
		for name, value := range headers.(map[string]interface{}) {
			header.Set(name, value.(string))
		}
	}
	login := d.Get("login").(string)
	// if resourceName == "" {
	// 	return fmt.Errorf("Invalid resource type specified")
	// }
	b, err := client.doGetDetailsByName(client.BaseUrl.String(), login)
	if err != nil {
		return
	}
	outputs, err := flattenNameDetailsResponse(b)
	if err != nil {
		return
	}
	marshalData(d, outputs)

	return
}

func flattenNameDetailsResponse(b []byte) (outputs map[string]interface{}, err error) {
	var data map[string]interface{}
	err = json.Unmarshal(b, &data)
	if err != nil {
		err = fmt.Errorf("Cannot unmarshal json of API response: %v", err)
		return
	} else if data["result"] == "" {
		err = fmt.Errorf("missing result key in API response: %v", err)
		return
	}

	outputs = make(map[string]interface{})
	outputs["id"] = data["id"]
	outputs["login"] = data["login"]
	outputs["name"] = data["name"]
	outputs["full_name"] = data["fullName"]
	outputs["email"] = data["email"]
	outputs["type"] = data["$type"]

	return
}