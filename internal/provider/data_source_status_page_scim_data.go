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
var _ datasource.DataSource = &StatusPageScimDataDataSource{}

func NewStatusPageScimDataDataSource() datasource.DataSource {
    return &StatusPageScimDataDataSource{}
}

// StatusPageScimDataDataSource defines the data source implementation.
type StatusPageScimDataDataSource struct {
    client *Client
}

// StatusPageScimDataDataSourceModel describes the data source data model.
type StatusPageScimDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    StatusPageId types.String `tfsdk:"status_page_id"`
    Description types.String `tfsdk:"description"`
    BearerToken types.String `tfsdk:"bearer_token"`
    AutoProvisionUsers types.Bool `tfsdk:"auto_provision_users"`
    AutoDeprovisionUsers types.Bool `tfsdk:"auto_deprovision_users"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
}

func (d *StatusPageScimDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_status_page_scim_data"
}

func (d *StatusPageScimDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "status_page_scim_data data source",

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
            "project_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "status_page_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Friendly description to help you remember. Permissions - Create: [Project Owner, Project Admin, Create Status Page SSO], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Status Page SSO], Update: [Project Owner, Project Admin, Edit Status Page SSO]",
                Computed: true,
            },
            "bearer_token": schema.StringAttribute{
                MarkdownDescription: "Bearer token for SCIM authentication. Keep this secure.. Permissions - Create: [Project Owner, Project Admin, Create Status Page SSO], Read: [Project Owner, Project Admin, Read Status Page SSO], Update: [Project Owner, Project Admin, Edit Status Page SSO]",
                Computed: true,
            },
            "auto_provision_users": schema.BoolAttribute{
                MarkdownDescription: "Automatically create status page users when they are added via SCIM. Permissions - Create: [Project Owner, Project Admin, Create Status Page SSO], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Status Page SSO], Update: [Project Owner, Project Admin, Edit Status Page SSO]",
                Computed: true,
            },
            "auto_deprovision_users": schema.BoolAttribute{
                MarkdownDescription: "Automatically remove status page users when they are removed via SCIM. Permissions - Create: [Project Owner, Project Admin, Create Status Page SSO], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Status Page SSO], Update: [Project Owner, Project Admin, Edit Status Page SSO]",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "deleted_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
        },
    }
}

func (d *StatusPageScimDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *StatusPageScimDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data StatusPageScimDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "status-page-scim" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read status_page_scim_data, got error: %s", err))
        return
    }

    var statusPageScimDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &statusPageScimDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse status_page_scim_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := statusPageScimDataResponse["data"].(map[string]interface{}); ok {
        statusPageScimDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := statusPageScimDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := statusPageScimDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := statusPageScimDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := statusPageScimDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := statusPageScimDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := statusPageScimDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := statusPageScimDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := statusPageScimDataResponse["status_page_id"].(string); ok {
        data.StatusPageId = types.StringValue(val)
    }
    if val, ok := statusPageScimDataResponse["description"].(string); ok {
        data.Description = types.StringValue(val)
    }
    if val, ok := statusPageScimDataResponse["bearer_token"].(string); ok {
        data.BearerToken = types.StringValue(val)
    }
    if val, ok := statusPageScimDataResponse["auto_provision_users"].(bool); ok {
        data.AutoProvisionUsers = types.BoolValue(val)
    }
    if val, ok := statusPageScimDataResponse["auto_deprovision_users"].(bool); ok {
        data.AutoDeprovisionUsers = types.BoolValue(val)
    }
    if val, ok := statusPageScimDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := statusPageScimDataResponse["deleted_by_user_id"].(string); ok {
        data.DeletedByUserId = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
