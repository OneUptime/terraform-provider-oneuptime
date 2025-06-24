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
var _ datasource.DataSource = &StatusPageGroupDataDataSource{}

func NewStatusPageGroupDataDataSource() datasource.DataSource {
    return &StatusPageGroupDataDataSource{}
}

// StatusPageGroupDataDataSource defines the data source implementation.
type StatusPageGroupDataDataSource struct {
    client *Client
}

// StatusPageGroupDataDataSourceModel describes the data source data model.
type StatusPageGroupDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    StatusPageId types.String `tfsdk:"status_page_id"`
    Slug types.String `tfsdk:"slug"`
    Description types.String `tfsdk:"description"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    Order types.Number `tfsdk:"order"`
    IsExpandedByDefault types.Bool `tfsdk:"is_expanded_by_default"`
    ShowCurrentStatus types.Bool `tfsdk:"show_current_status"`
    ShowUptimePercent types.Bool `tfsdk:"show_uptime_percent"`
    UptimePercentPrecision types.String `tfsdk:"uptime_percent_precision"`
}

func (d *StatusPageGroupDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_status_page_group_data"
}

func (d *StatusPageGroupDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "status_page_group_data data source",

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
                MarkdownDescription: "Version",
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
            "slug": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Status Page Group], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page Group], Read: [Project Owner, Project Admin, Project Member, Read Status Page Group], Update: [Project Owner, Project Admin, Project Member, Edit Status Page Group]",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "order": schema.NumberAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page Group], Read: [Project Owner, Project Admin, Project Member, Read Status Page Group], Update: [Project Owner, Project Admin, Project Member, Edit Status Page Group]",
                Computed: true,
            },
            "is_expanded_by_default": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page Group], Read: [Project Owner, Project Admin, Project Member, Read Status Page Group], Update: [Project Owner, Project Admin, Project Member, Edit Status Page Group]",
                Computed: true,
            },
            "show_current_status": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page Group], Read: [Project Owner, Project Admin, Project Member, Read Status Page Group], Update: [Project Owner, Project Admin, Project Member, Edit Status Page Group]",
                Computed: true,
            },
            "show_uptime_percent": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page Group], Read: [Project Owner, Project Admin, Project Member, Read Status Page Group], Update: [Project Owner, Project Admin, Project Member, Edit Status Page Group]",
                Computed: true,
            },
            "uptime_percent_precision": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page Group], Read: [Project Owner, Project Admin, Project Member, Read Status Page Group], Update: [Project Owner, Project Admin, Project Member, Edit Status Page Group]",
                Computed: true,
            },
        },
    }
}

func (d *StatusPageGroupDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *StatusPageGroupDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data StatusPageGroupDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "status-page-group" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read status_page_group_data, got error: %s", err))
        return
    }

    var statusPageGroupDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &statusPageGroupDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse status_page_group_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := statusPageGroupDataResponse["data"].(map[string]interface{}); ok {
        statusPageGroupDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := statusPageGroupDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := statusPageGroupDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := statusPageGroupDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := statusPageGroupDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := statusPageGroupDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := statusPageGroupDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := statusPageGroupDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := statusPageGroupDataResponse["status_page_id"].(string); ok {
        data.StatusPageId = types.StringValue(val)
    }
    if val, ok := statusPageGroupDataResponse["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    }
    if val, ok := statusPageGroupDataResponse["description"].(string); ok {
        data.Description = types.StringValue(val)
    }
    if val, ok := statusPageGroupDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := statusPageGroupDataResponse["order"].(float64); ok {
        data.Order = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := statusPageGroupDataResponse["is_expanded_by_default"].(bool); ok {
        data.IsExpandedByDefault = types.BoolValue(val)
    }
    if val, ok := statusPageGroupDataResponse["show_current_status"].(bool); ok {
        data.ShowCurrentStatus = types.BoolValue(val)
    }
    if val, ok := statusPageGroupDataResponse["show_uptime_percent"].(bool); ok {
        data.ShowUptimePercent = types.BoolValue(val)
    }
    if val, ok := statusPageGroupDataResponse["uptime_percent_precision"].(string); ok {
        data.UptimePercentPrecision = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
