package provider

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"io/ioutil"
	"os"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceFile() *schema.Resource {
	return &schema.Resource{
		// This description is used by tPhe documentation generator and the language server.
		Description: "Create a text file and write to a local folder.",

		CreateContext: resourceCreate,
		ReadContext:   resourceRead,
		DeleteContext: resourceDelete,
		UpdateContext: resourceUpdate,

		Schema: map[string]*schema.Schema{
			"content": {
				Description: "Content of file",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    false,
			},
			"filename": {
				Description: "The path to the file that will be created",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"hash": {
				Description: "SHA1 Hash of the file",
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics

	content := d.Get("content").(string)
	filename := d.Get("filename").(string)

	dataBytes := []byte(content)
	err := os.WriteFile(filename, dataBytes, 0644)
	if err != nil {
		return diag.FromErr(err)
	}

	outputChecksum := sha1.Sum(dataBytes)
	checksum := hex.EncodeToString(outputChecksum[:])
	d.SetId(checksum)
	if d.Set("hash", checksum); err != nil {
		return diag.FromErr(err)
	  }
	

	tflog.Trace(ctx, "created a file resource with id " + d.Id())

	return diags
}

func resourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// If the output file doesn't exist, mark the resource for creation.
	outputPath := d.Get("filename").(string)
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		d.SetId("")
		return nil
	}

	// Verify that the content of the destination file matches the content we
	// expect. Otherwise, the file might have been modified externally, and we
	// must reconcile.
	outputContent, err := ioutil.ReadFile(outputPath)
	if err != nil {
		return diag.FromErr(err)
	}

	outputChecksum := sha1.Sum(outputContent)
	checksum := hex.EncodeToString(outputChecksum[:])
	if checksum != d.Id() {
		d.Set("content", string(outputContent))
		d.Set("hash", checksum)
		return nil
	}

	return nil
}

func resourceUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	outputPath := d.Get("filename").(string)
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		d.SetId("")
		return nil
	}


	//filename change -> delete, create new
	//content chagne, only set new text

	if d.HasChange("content") {
		content := d.Get("content").(string)

		dataBytes := []byte(content)
		err := os.WriteFile(outputPath, dataBytes, 0644)
		if err != nil {
			return diag.FromErr(err)
		}
	}


	return resourceRead(ctx, d, meta)
}



func resourceDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	os.Remove(d.Get("filename").(string))
	return nil
}
