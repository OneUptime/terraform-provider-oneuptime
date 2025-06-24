package provider

import (
    "context"
    "fmt"
    "math/big"
    "github.com/hashicorp/terraform-plugin-framework/attr"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ApiKeyPermissionDataDataSource{}

func NewApiKeyPermissionDataDataSource() datasource.DataSource {
    return &ApiKeyPermissionDataDataSource{}
}

// ApiKeyPermissionDataDataSource defines the data source implementation.
type ApiKeyPermissionDataDataSource struct {
    client *Client
}

// ApiKeyPermissionDataDataSourceModel describes the data source data model.
type ApiKeyPermissionDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ApiKeyId types.String `tfsdk:"api_key_id"`
    ProjectId types.String `tfsdk:"project_id"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    Permission types.String `tfsdk:"permission"`
    Labels types.List `tfsdk:"labels"`
    IsBlockPermission types.Bool `tfsdk:"is_block_permission"`
}

func (d *ApiKeyPermissionDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_api_key_permission_data"
}

func (d *ApiKeyPermissionDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "api_key_permission_data data source",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Identifier to filter by",
                Optional: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Name to filter by",
                Optional: true,
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
                MarkdownDescription: "Version",
                Computed: true,
            },
            "api_key_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "permission": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Create API Key, Edit API Key Permissions], Read: [Project Owner, Project Admin, Read API Key], Update: [Project Owner, Project Admin, Edit API Key Permissions, Edit API Key]",
                Computed: true,
            },
            "labels": schema.ListAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Create API Key, Edit API Key Permissions], Read: [Project Owner, Project Admin, Read API Key], Update: [Project Owner, Project Admin, Edit API Key Permissions, Edit API Key]",
                Computed: true,
                ElementType: types.StringType,
            },
            "is_block_permission": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Create API Key, Edit API Key Permissions], Read: [Project Owner, Project Admin, Read API Key], Update: [Project Owner, Project Admin, Edit API Key Permissions, Edit API Key]",
                Computed: true,
            },
        },
    }
}

func (d *ApiKeyPermissionDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ApiKeyPermissionDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data ApiKeyPermissionDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "api-key-permission" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read api_key_permission_data, got error: %s", err))
        return
    }

    var apiKeyPermissionDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &apiKeyPermissionDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse api_key_permission_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := apiKeyPermissionDataResponse["data"].(map[string]interface{}); ok {
        apiKeyPermissionDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := apiKeyPermissionDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := apiKeyPermissionDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := apiKeyPermissionDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := apiKeyPermissionDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := apiKeyPermissionDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := apiKeyPermissionDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := apiKeyPermissionDataResponse["api_key_id"].(string); ok {
        data.ApiKeyId = types.StringValue(val)
    }
    if val, ok := apiKeyPermissionDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := apiKeyPermissionDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := apiKeyPermissionDataResponse["permission"].(string); ok {
        data.Permission = types.StringValue(val)
    }
    if val, ok := apiKeyPermissionDataResponse["labels"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        listValue, _ := types.ListValue(types.StringType, elements)
        data.Labels = listValue
    }
    if val, ok := apiKeyPermissionDataResponse["is_block_permission"].(bool); ok {
        data.IsBlockPermission = types.BoolValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
