package iothub

import (
	"crypto/x509"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/preview/iothub/mgmt/2019-03-22-preview/devices"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"time"
)

func resourceArmIotHubCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmIotHubCertificateCreateUpdate,
		Read:   resourceArmIotHubDPSCertificateRead,
		Update: resourceArmIotHubCertificateCreateUpdate,
		Delete: resourceArmIotHubDPSCertificateDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.IoTHubName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"iothub_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.IoTHubName,
			},

			"certificate_content": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Sensitive:    true,
			},
		},
	}
}

func resourceArmIotHubCertificateCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.CertificatesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	iothubName := d.Get("iothub_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, iothubName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing IoT Hub Certificate %q (IoT Hub %q / Resource Group %q): %+v", name, iothubName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_iothub_certificate", *existing.ID)
		}
	}

	certificate := devices.CertificateBodyDescription{
		Certificate: utils.String(d.Get("certificate_content").(string)),
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, iothubName, name, certificate, ""); err != nil {
		return fmt.Errorf("Error creating/updating IoT Certificate %q (IoT Hub %q / Resource Group %q): %+v", name, iothubName, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, iothubName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving IoT Hub Certificate %q (IoT Hub %q / Resource Group %q): %+v", name, iothubName, resourceGroup, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("Cannot read IoT Hub Certificate %q (Iot Hub %q / Resource Group %q): %+v", name, iothubName, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	// TODO: Verification
	//verificationCertificate := x509.Certificate{
	//	BasicConstraintsValid: true,
	//	IsCA: false,
	//}
	//verificationCode, err := client.GenerateVerificationCode(ctx, resourceGroup, iothubName, name, *resp.Etag)
	//if err != nil {
	//	return fmt.Errorf("Error retrieving IoT Hub Certificate %q (IoT Hub %q / Resource Group %q): %+v", name, iothubName, resourceGroup, err)
	//}

	//client.Verify(ctx, resourceGroup, iothubName, name, verificationCode, *verificationCode.Etag)

	return resourceArmIotHubCertificateRead(d, meta)
}

func resourceArmIotHubCertificateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.CertificatesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	iothubName := id.Path["IotHubs"]
	name := id.Path["certificates"]

	resp, err := client.Get(ctx, resourceGroup, iothubName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving IoT Device Provisioning Service Certificate %q (Device Provisioning Service %q / Resource Group %q): %+v", name, iotDPSName, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("iothub_name", iothubName)
	d.Set("resource_group_name", resourceGroup)

	if props := resp.Properties; props != nil {
		d.Set("subject", props.Subject)
		d.Set("expiry", props.Expiry.Format(time.RFC3339))
		d.Set("thumbprint", props.Thumbprint)
		d.Set("verified", props.IsVerified)
		d.Set("created", props.Created.Format(time.RFC3339))
		d.Set("updated", props.Updated.Format(time.RFC3339))
		d.Set("certificate_content", props.Certificate)
	}

	return nil
}
