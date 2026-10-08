package provider

import (
    "context"
    "encoding/json"
    "fmt"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &FileDataSource{}

func NewFileDataSource() datasource.DataSource {
    return &FileDataSource{}
}

// FileDataSource defines the data source implementation.
type FileDataSource struct {
    client *Client
}

// FileDataSourceModel describes the data source data model.
type FileDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    File types.String `tfsdk:"file"`
    Name types.String `tfsdk:"name"`
    FileType types.String `tfsdk:"file_type"`
    Slug types.String `tfsdk:"slug"`
    IsPublic types.Bool `tfsdk:"is_public"`
    ImageAccessToken types.String `tfsdk:"image_access_token"`
}

func (d *FileDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_file"
}

func (d *FileDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "BLOB or File storage Look up an existing file by `id`, or by any of its other arguments (`name`, `file`, `file_type`, ...): each one set must match, and exactly one file may match them all.",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Look up by unique identifier. Leave unset to look up by the other arguments instead.",
                Optional: true,
                Computed: true,
            },
            "file": schema.StringAttribute{
                Optional: true,
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Any friendly name of this object.",
                Optional: true,
                Computed: true,
            },
            "file_type": schema.StringAttribute{
                Optional: true,
                Computed: true,
            },
            "slug": schema.StringAttribute{
                Optional: true,
                Computed: true,
            },
            "is_public": schema.BoolAttribute{
                MarkdownDescription: "Whether anyone may read the file without signing in. Set by OneUptime: every upload starts private, and a file becomes public only when a record that shows it to everyone, such as a public note or a probe's icon, is published.",
                Optional: true,
                Computed: true,
            },
            "image_access_token": schema.StringAttribute{
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *FileDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
    // Prevent panic if the provider has not been configured.
    if req.ProviderData == nil {
        return
    }

    client, ok := req.ProviderData.(*Client)

    if !ok {
        resp.Diagnostics.AddError(
            "Unexpected Data Source Configure Type",
            fmt.Sprintf("Expected *Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
        )

        return
    }

    d.client = client
}

func (d *FileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data FileDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.File.IsNull() && !data.File.IsUnknown() {
        filters["file"] = data.File.ValueString()
        filterNames = append(filterNames, "file = "+fmt.Sprintf("%q", data.File.ValueString()))
    }
    if !data.Name.IsNull() && !data.Name.IsUnknown() {
        filters["name"] = data.Name.ValueString()
        filterNames = append(filterNames, "name = "+fmt.Sprintf("%q", data.Name.ValueString()))
    }
    if !data.FileType.IsNull() && !data.FileType.IsUnknown() {
        filters["fileType"] = data.FileType.ValueString()
        filterNames = append(filterNames, "file_type = "+fmt.Sprintf("%q", data.FileType.ValueString()))
    }
    if !data.Slug.IsNull() && !data.Slug.IsUnknown() {
        filters["slug"] = data.Slug.ValueString()
        filterNames = append(filterNames, "slug = "+fmt.Sprintf("%q", data.Slug.ValueString()))
    }
    if !data.IsPublic.IsNull() && !data.IsPublic.IsUnknown() {
        filters["isPublic"] = data.IsPublic.ValueBool()
        filterNames = append(filterNames, "is_public = "+fmt.Sprintf("%t", data.IsPublic.ValueBool()))
    }
    if !data.ImageAccessToken.IsNull() && !data.ImageAccessToken.IsUnknown() {
        filters["imageAccessToken"] = data.ImageAccessToken.ValueString()
        filterNames = append(filterNames, "image_access_token = "+fmt.Sprintf("%q", data.ImageAccessToken.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the file up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the file up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "file": true,
        "name": true,
        "fileType": true,
        "slug": true,
        "isPublic": true,
        "imageAccessToken": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        // No get endpoint: find it in the list by id.
        filters["_id"] = data.Id.ValueString()
        filterNames = append(filterNames, fmt.Sprintf("id = %q", data.Id.ValueString()))
    }
    if item == nil {
        listBody := map[string]interface{}{
            "query":  filters,
            "select": selectParam,
            // limit 2 is enough to detect ambiguity without paging.
            "limit": 2,
        }
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/file/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list file, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list file: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No file matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one file matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for file.")
            return
        }
        item = first
    }

    // Update the model with response data
    if obj, ok := item["_id"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Id = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Id = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Id = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Id = types.StringValue(string(jsonBytes))
        } else {
            data.Id = types.StringNull()
        }
    } else if val, ok := item["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    if obj, ok := item["file"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.File = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.File = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.File = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.File = types.StringValue(string(jsonBytes))
        } else {
            data.File = types.StringNull()
        }
    } else if val, ok := item["file"].(string); ok {
        data.File = types.StringValue(val)
    } else {
        data.File = types.StringNull()
    }
    if obj, ok := item["name"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Name = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Name = types.StringValue(string(jsonBytes))
        } else {
            data.Name = types.StringNull()
        }
    } else if val, ok := item["name"].(string); ok {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if obj, ok := item["fileType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FileType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.FileType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.FileType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.FileType = types.StringValue(string(jsonBytes))
        } else {
            data.FileType = types.StringNull()
        }
    } else if val, ok := item["fileType"].(string); ok {
        data.FileType = types.StringValue(val)
    } else {
        data.FileType = types.StringNull()
    }
    if obj, ok := item["slug"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Slug = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Slug = types.StringValue(string(jsonBytes))
        } else {
            data.Slug = types.StringNull()
        }
    } else if val, ok := item["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
    }
    if val, ok := item["isPublic"].(bool); ok {
        data.IsPublic = types.BoolValue(val)
    } else {
        data.IsPublic = types.BoolNull()
    }
    if obj, ok := item["imageAccessToken"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ImageAccessToken = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ImageAccessToken = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ImageAccessToken = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ImageAccessToken = types.StringValue(string(jsonBytes))
        } else {
            data.ImageAccessToken = types.StringNull()
        }
    } else if val, ok := item["imageAccessToken"].(string); ok {
        data.ImageAccessToken = types.StringValue(val)
    } else {
        data.ImageAccessToken = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
