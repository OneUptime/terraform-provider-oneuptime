package provider

import (
    "context"
    "fmt"
    "math/big"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &FileDataDataSource{}

func NewFileDataDataSource() datasource.DataSource {
    return &FileDataDataSource{}
}

// FileDataDataSource defines the data source implementation.
type FileDataDataSource struct {
    client *Client
}

// FileDataDataSourceModel describes the data source data model.
type FileDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    File types.String `tfsdk:"file"`
    FileType types.String `tfsdk:"file_type"`
    Slug types.String `tfsdk:"slug"`
    IsPublic types.String `tfsdk:"is_public"`
}

func (d *FileDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_file_data"
}

func (d *FileDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "file_data data source",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Identifier to filter by",
                Optional: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Name to filter by",
                Computed: true,
            },
            "created_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "updated_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "deleted_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "version": schema.NumberAttribute{
                MarkdownDescription: "Object version",
                Computed: true,
            },
            "file": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Logged in User], Read: [Logged in User], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "file_type": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Logged in User], Read: [Logged in User], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Logged in User], Read: [Logged in User], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "is_public": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Logged in User], Read: [Logged in User], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
        },
    }
}

func (d *FileDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *FileDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data FileDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "file" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read file_data, got error: %s", err))
        return
    }

    var fileDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &fileDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse file_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := fileDataResponse["data"].(map[string]interface{}); ok {
        fileDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := fileDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := fileDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := fileDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := fileDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := fileDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := fileDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := fileDataResponse["file"].(string); ok {
        data.File = types.StringValue(val)
    }
    if val, ok := fileDataResponse["file_type"].(string); ok {
        data.FileType = types.StringValue(val)
    }
    if val, ok := fileDataResponse["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    }
    if val, ok := fileDataResponse["is_public"].(string); ok {
        data.IsPublic = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
