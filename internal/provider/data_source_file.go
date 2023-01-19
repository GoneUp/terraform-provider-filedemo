package provider

import (
	"context"
	"os"
	"crypto/sha1"
	"io/ioutil"
	"encoding/hex"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var fileSchema = map[string]*schema.Schema{
	"filename": {
		Type:     schema.TypeString,
		Required: true,
	},
	"hash": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"content": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

func dataSourceFile() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Reads file data",
		ReadContext: fileRead,
		Schema:      fileSchema,
	}
}

func fileRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	//var diags diag.Diagnostics

	// If the output file doesn't exist, mark the resource for creation.
	outputPath := d.Get("filename").(string)
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		d.SetId("")
		return diag.FromErr(err)
	}

	outputContentbytes, err := ioutil.ReadFile(outputPath)
	if err != nil {
		return diag.FromErr(err)
	}

	outputChecksum := sha1.Sum(outputContentbytes)
	checksum := hex.EncodeToString(outputChecksum[:])
	if err := d.Set("hash", checksum); err != nil {
		return diag.FromErr(err)
	}

	text := string(outputContentbytes)
	if err := d.Set("content", text); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(checksum)
	return nil
}
