package provider

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVMExport() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVMExportCreate,
		ReadContext:   resourceVMExportRead,
		DeleteContext: resourceVMExportDelete,

		Schema: map[string]*schema.Schema{
			"vm_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "UUID or name of the VM to export",
			},
			"output_path": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Output file path (.ova or .ovf)",
			},
			"format": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "ova",
				ForceNew:    true,
				Description: "Export format: ova (default) or ovf",
				ValidateFunc: func(val any, key string) (warns []string, errs []error) {
					v := val.(string)
					if v != "ova" && v != "ovf" {
						errs = append(errs, fmt.Errorf("%q must be ova or ovf, got: %s", key, v))
					}
					return
				},
			},
			"manifest": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				ForceNew:    true,
				Description: "Include manifest file",
			},
			"options": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Additional export options passed to VBoxManage export",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"file_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Actual output file path after export",
			},
		},

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceVMExportCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	vmID := d.Get("vm_id").(string)
	outputPath := d.Get("output_path").(string)
	format := d.Get("format").(string)
	manifest := d.Get("manifest").(bool)

	args := []string{"export", vmID, "--output", outputPath}

	if format == "ovf" {
		args = append(args, "--ovf20")
	}

	if manifest {
		args = append(args, "--manifest")
	}

	// Append additional options
	if v, ok := d.GetOk("options"); ok {
		for _, opt := range v.([]any) {
			args = append(args, opt.(string))
		}
	}

	if _, _, err := vboxRun(ctx, args...); err != nil {
		return diag.Errorf("failed to export VM %s to %s: %v", vmID, outputPath, err)
	}

	d.SetId(outputPath)
	if err := d.Set("file_path", outputPath); err != nil {
		return diag.Errorf("failed to set file_path: %v", err)
	}

	return nil
}

func resourceVMExportRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	outputPath := d.Id()
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		d.SetId("")
		return nil
	}
	if err := d.Set("file_path", outputPath); err != nil {
		return diag.Errorf("failed to set file_path: %v", err)
	}
	return nil
}

func resourceVMExportDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	outputPath := d.Id()

	// Remove the main export file
	os.Remove(outputPath) //nolint:errcheck

	// If OVF format, also clean up associated files (.vmdk, .mf, etc.)
	if strings.HasSuffix(strings.ToLower(outputPath), ".ovf") {
		dir := filepath.Dir(outputPath)
		base := strings.TrimSuffix(filepath.Base(outputPath), filepath.Ext(outputPath))
		matches, _ := filepath.Glob(filepath.Join(dir, base+"*"))
		for _, m := range matches {
			os.Remove(m) //nolint:errcheck
		}
	}

	d.SetId("")
	return nil
}
