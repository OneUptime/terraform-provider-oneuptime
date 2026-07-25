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
var _ datasource.DataSource = &NetworkSiteDataDataSource{}

func NewNetworkSiteDataDataSource() datasource.DataSource {
    return &NetworkSiteDataDataSource{}
}

// NetworkSiteDataDataSource defines the data source implementation.
type NetworkSiteDataDataSource struct {
    client *Client
}

// NetworkSiteDataDataSourceModel describes the data source data model.
type NetworkSiteDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    Slug types.String `tfsdk:"slug"`
    Description types.String `tfsdk:"description"`
    SiteType types.String `tfsdk:"site_type"`
    NetworkSiteTypeId types.String `tfsdk:"network_site_type_id"`
    ParentSiteId types.String `tfsdk:"parent_site_id"`
    MaterializedPath types.String `tfsdk:"materialized_path"`
    Depth types.Number `tfsdk:"depth"`
    Address types.String `tfsdk:"address"`
    Latitude types.Number `tfsdk:"latitude"`
    Longitude types.Number `tfsdk:"longitude"`
    CurrentMonitorStatusId types.String `tfsdk:"current_monitor_status_id"`
    LastRollupAt types.String `tfsdk:"last_rollup_at"`
    ShouldAlertWhenUnhealthy types.Bool `tfsdk:"should_alert_when_unhealthy"`
    AlertSeverityId types.String `tfsdk:"alert_severity_id"`
    CurrentActiveAlertId types.String `tfsdk:"current_active_alert_id"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
}

func (d *NetworkSiteDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_network_site_data"
}

func (d *NetworkSiteDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "network_site_data data source",

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
            "slug": schema.StringAttribute{
                MarkdownDescription: "Friendly globally unique name for your object. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Site], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Friendly description for this network site. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Network Site], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Site], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Network Site]",
                Computed: true,
            },
            "site_type": schema.StringAttribute{
                MarkdownDescription: "Deprecated legacy site type string. Use the Network Site Type relation instead; this column exists only for the backfill migration and will be removed.. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Network Site], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Site], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Network Site]",
                Computed: true,
            },
            "network_site_type_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "parent_site_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "materialized_path": schema.StringAttribute{
                MarkdownDescription: "Slash-separated ancestor IDs of this site (e.g. '/rootId/childId/'). Managed by the server on parent changes; used for subtree queries and rollups.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Site], Update: [Project Owner, Project Admin, Edit Network Site]",
                Computed: true,
            },
            "depth": schema.NumberAttribute{
                MarkdownDescription: "Number of ancestors above this site (0 for root sites). Managed by the server on parent changes.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Site], Update: [Project Owner, Project Admin, Edit Network Site]",
                Computed: true,
            },
            "address": schema.StringAttribute{
                MarkdownDescription: "Street address of this site, shown on map views. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Network Site], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Site], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Network Site]",
                Computed: true,
            },
            "latitude": schema.NumberAttribute{
                MarkdownDescription: "Latitude of this site, for US and world map views. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Network Site], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Site], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Network Site]",
                Computed: true,
            },
            "longitude": schema.NumberAttribute{
                MarkdownDescription: "Longitude of this site, for US and world map views. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Network Site], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Site], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Network Site]",
                Computed: true,
            },
            "current_monitor_status_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "last_rollup_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "should_alert_when_unhealthy": schema.BoolAttribute{
                MarkdownDescription: "When enabled, an alert opens when this site's health rollup turns non-operational and auto-resolves when it recovers.. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Network Site], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Site], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Network Site]",
                Computed: true,
            },
            "alert_severity_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "current_active_alert_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
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

func (d *NetworkSiteDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *NetworkSiteDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data NetworkSiteDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "network-site" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read network_site_data, got error: %s", err))
        return
    }

    var networkSiteDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &networkSiteDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse network_site_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := networkSiteDataResponse["data"].(map[string]interface{}); ok {
        networkSiteDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := networkSiteDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := networkSiteDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := networkSiteDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := networkSiteDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := networkSiteDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := networkSiteDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := networkSiteDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := networkSiteDataResponse["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    }
    if val, ok := networkSiteDataResponse["description"].(string); ok {
        data.Description = types.StringValue(val)
    }
    if val, ok := networkSiteDataResponse["site_type"].(string); ok {
        data.SiteType = types.StringValue(val)
    }
    if val, ok := networkSiteDataResponse["network_site_type_id"].(string); ok {
        data.NetworkSiteTypeId = types.StringValue(val)
    }
    if val, ok := networkSiteDataResponse["parent_site_id"].(string); ok {
        data.ParentSiteId = types.StringValue(val)
    }
    if val, ok := networkSiteDataResponse["materialized_path"].(string); ok {
        data.MaterializedPath = types.StringValue(val)
    }
    if val, ok := networkSiteDataResponse["depth"].(float64); ok {
        data.Depth = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := networkSiteDataResponse["address"].(string); ok {
        data.Address = types.StringValue(val)
    }
    if val, ok := networkSiteDataResponse["latitude"].(float64); ok {
        data.Latitude = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := networkSiteDataResponse["longitude"].(float64); ok {
        data.Longitude = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := networkSiteDataResponse["current_monitor_status_id"].(string); ok {
        data.CurrentMonitorStatusId = types.StringValue(val)
    }
    if val, ok := networkSiteDataResponse["last_rollup_at"].(string); ok {
        data.LastRollupAt = types.StringValue(val)
    }
    if val, ok := networkSiteDataResponse["should_alert_when_unhealthy"].(bool); ok {
        data.ShouldAlertWhenUnhealthy = types.BoolValue(val)
    }
    if val, ok := networkSiteDataResponse["alert_severity_id"].(string); ok {
        data.AlertSeverityId = types.StringValue(val)
    }
    if val, ok := networkSiteDataResponse["current_active_alert_id"].(string); ok {
        data.CurrentActiveAlertId = types.StringValue(val)
    }
    if val, ok := networkSiteDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := networkSiteDataResponse["deleted_by_user_id"].(string); ok {
        data.DeletedByUserId = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
