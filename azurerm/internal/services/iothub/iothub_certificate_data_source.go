package iothub

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"time"
)

func dataSourceArmIoTHubCertificate() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIotHubCertificateRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.IoTHubName,
			},

			"iothub_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.IoTHubName,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"subject": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"expiry": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"thumbprint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"verified": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"updated": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"certificate_content": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceIotHubCertificateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.CertificatesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	iothubName := d.Get("iothub_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	resp, err := client.Get(ctx, resourceGroup, iothubName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: IoT Hub Certficate %q (IoT Hub %q) was not found", name, iothubName)
		}

		return fmt.Errorf("Error retrieving IoT Certificate %q (IoT Hub %q): %+v", name, iothubName, err)
	}

	d.SetId(*resp.ID)
	d.Set("name", resp.Name)
	d.Set("iothub_name", iothubName)
	d.Set("resource_group_name", resourceGroup)

	if props := resp.Properties; props != nil {
		d.Set("subject", props.Subject)
		d.Set("expiry", props.Expiry)
		d.Set("thumbprint", props.Thumbprint)
		d.Set("verified", props.IsVerified)
		d.Set("created", props.Created)
		d.Set("updated", props.Updated)
		d.Set("certificate_content", props.Certificate)
	}

	return nil
}
