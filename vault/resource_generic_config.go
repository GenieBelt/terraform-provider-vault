package vault

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"

	"github.com/hashicorp/vault/api"
)

func genericConfigResource() *schema.Resource {
	return &schema.Resource{
		Create: genericConfigResourceWrite,
		Update: genericConfigResourceWrite,
		Delete: genericConfigResourceDelete,
		Read:   genericConfigResourceRead,

		Schema: map[string]*schema.Schema{
			"path": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Full path where the generic config will be written.",
			},

			// Data is passed as JSON so that an arbitrary structure is
			// possible, rather than forcing e.g. all values to be strings.
			"data_json": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "JSON-encoded config data to write.",
				// We rebuild the attached JSON string to a simple singleline
				// string. This makes terraform not want to change when an extra
				// space is included in the JSON string. It is also necesarry
				// when allow_read is true for comparing values.
				StateFunc:    NormalizeDataJSON,
				ValidateFunc: ValidateDataJSON,
			},
		},
	}
}

func genericConfigResourceWrite(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	path := d.Get("path").(string)

	var data map[string]interface{}
	err := json.Unmarshal([]byte(d.Get("data_json").(string)), &data)
	if err != nil {
		return fmt.Errorf("data_json %#v syntax error: %s", d.Get("data_json"), err)
	}

	log.Printf("[DEBUG] Writing generic Vault config to %s", path)
	_, err = client.Logical().Write(path, data)
	if err != nil {
		return fmt.Errorf("error writing to Vault: %s", err)
	}

	d.SetId(path)

	return nil
}

func genericConfigResourceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	path := d.Id()

	log.Printf("[DEBUG] Deleting vault_generic_config from %q", path)
	_, err := client.Logical().Delete(path)
	if err != nil {
		log.Printf("[DEBUG] Error deleting %q from Vault: %q", path, err)
		return nil
	}

	return nil
}

func genericConfigResourceRead(d *schema.ResourceData, meta interface{}) error {
	path := d.Get("path").(string)
	d.SetId(path)
	return nil
}
