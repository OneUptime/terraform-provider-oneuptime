package provider

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "math/big"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &NetworkSiteDataSource{}

func NewNetworkSiteDataSource() datasource.DataSource {
    return &NetworkSiteDataSource{}
}

// NetworkSiteDataSource defines the data source implementation.
type NetworkSiteDataSource struct {
    client *Client
}

// NetworkSiteDataSourceModel describes the data source data model.
type NetworkSiteDataSourceModel struct {
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
    HealthRollupPolicy types.String `tfsdk:"health_rollup_policy"`
    OfflineThresholdPercent types.Number `tfsdk:"offline_threshold_percent"`
    ShouldAlertWhenUnhealthy types.Bool `tfsdk:"should_alert_when_unhealthy"`
    AlertSeverityId types.String `tfsdk:"alert_severity_id"`
    CurrentActiveAlertId types.String `tfsdk:"current_active_alert_id"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
}

func (d *NetworkSiteDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_network_site"
}

func (d *NetworkSiteDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Self-nesting sites (Account Type -> Region / Franchisee -> Market -> Unit) that group Network Devices into a drill-down hierarchy with a persisted health rollup. Look up an existing network_site by `id` or by `name`.",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Look up by unique identifier. Exactly one of `id` or `name` must be set.",
                Optional: true,
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.",
                Optional: true,
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
                MarkdownDescription: "Friendly globally unique name for your object.",
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Friendly description for this network site.",
                Computed: true,
            },
            "site_type": schema.StringAttribute{
                MarkdownDescription: "Deprecated legacy site type string. Use the Network Site Type relation instead; this column exists only for the backfill migration and will be removed..",
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
                MarkdownDescription: "Slash-separated ancestor IDs of this site (e.g. '/rootId/childId/'). Managed by the server on parent changes; used for subtree queries and rollups..",
                Computed: true,
            },
            "depth": schema.NumberAttribute{
                MarkdownDescription: "Number of ancestors above this site (0 for root sites). Managed by the server on parent changes..",
                Computed: true,
            },
            "address": schema.StringAttribute{
                MarkdownDescription: "Street address of this site, shown on map views.",
                Computed: true,
            },
            "latitude": schema.NumberAttribute{
                MarkdownDescription: "Latitude of this site, for US and world map views.",
                Computed: true,
            },
            "longitude": schema.NumberAttribute{
                MarkdownDescription: "Longitude of this site, for US and world map views.",
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
            "health_rollup_policy": schema.StringAttribute{
                MarkdownDescription: "How this site's status is derived from the devices beneath it: WorstStatus (any device offline makes the site offline) or PercentThreshold (the share of devices that are down decides)..",
                Computed: true,
            },
            "offline_threshold_percent": schema.NumberAttribute{
                MarkdownDescription: "With the PercentThreshold rollup policy: the share of reporting devices beneath this site that must be non-operational before the site itself is marked offline. Below it (but above zero) the site is degraded..",
                Computed: true,
            },
            "should_alert_when_unhealthy": schema.BoolAttribute{
                MarkdownDescription: "When enabled, an alert opens when this site's health rollup turns non-operational and auto-resolves when it recovers..",
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

func (d *NetworkSiteDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *NetworkSiteDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data NetworkSiteDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && data.Id.ValueString() != ""
    hasName := !data.Name.IsNull() && data.Name.ValueString() != ""
    if hasId == hasName {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Exactly one of `id` or `name` must be set to look up a network_site.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "name": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "projectId": true,
        "slug": true,
        "description": true,
        "siteType": true,
        "networkSiteTypeId": true,
        "parentSiteId": true,
        "materializedPath": true,
        "depth": true,
        "address": true,
        "latitude": true,
        "longitude": true,
        "currentMonitorStatusId": true,
        "lastRollupAt": true,
        "healthRollupPolicy": true,
        "offlineThresholdPercent": true,
        "shouldAlertWhenUnhealthy": true,
        "alertSeverityId": true,
        "currentActiveAlertId": true,
        "createdByUserId": true,
        "deletedByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/network-site/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read network_site, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No network_site found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read network_site: %s", err))
            return
        }
        if wrapper, ok := itemResponse["data"].(map[string]interface{}); ok {
            item = wrapper
        } else {
            item = itemResponse
        }
    } else {
        listBody := map[string]interface{}{
            "query": map[string]interface{}{
                "name": data.Name.ValueString(),
            },
            "select": selectParam,
            // limit 2 is enough to detect ambiguity without paging.
            "limit": 2,
        }
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/network-site/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list network_site, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list network_site: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No network_site found with name %q.", data.Name.ValueString()))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one network_site matches name %q. Use the id attribute to disambiguate.", data.Name.ValueString()))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for network_site.")
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
    if obj, ok := item["createdAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CreatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CreatedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CreatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedAt = types.StringNull()
        }
    } else if val, ok := item["createdAt"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    } else {
        data.CreatedAt = types.StringNull()
    }
    if obj, ok := item["updatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UpdatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.UpdatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.UpdatedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.UpdatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.UpdatedAt = types.StringNull()
        }
    } else if val, ok := item["updatedAt"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    } else {
        data.UpdatedAt = types.StringNull()
    }
    if obj, ok := item["deletedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DeletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DeletedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DeletedAt = types.StringValue(string(jsonBytes))
        } else {
            data.DeletedAt = types.StringNull()
        }
    } else if val, ok := item["deletedAt"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    } else {
        data.DeletedAt = types.StringNull()
    }
    if val, ok := item["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["version"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.Version = types.NumberValue(big.NewFloat(val))
        } else {
            data.Version = types.NumberNull()
        }
    } else {
        data.Version = types.NumberNull()
    }
    if obj, ok := item["projectId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ProjectId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ProjectId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ProjectId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ProjectId = types.StringValue(string(jsonBytes))
        } else {
            data.ProjectId = types.StringNull()
        }
    } else if val, ok := item["projectId"].(string); ok {
        data.ProjectId = types.StringValue(val)
    } else {
        data.ProjectId = types.StringNull()
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
    if obj, ok := item["description"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Description = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Description = types.StringValue(string(jsonBytes))
        } else {
            data.Description = types.StringNull()
        }
    } else if val, ok := item["description"].(string); ok {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if obj, ok := item["siteType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SiteType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SiteType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SiteType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SiteType = types.StringValue(string(jsonBytes))
        } else {
            data.SiteType = types.StringNull()
        }
    } else if val, ok := item["siteType"].(string); ok {
        data.SiteType = types.StringValue(val)
    } else {
        data.SiteType = types.StringNull()
    }
    if obj, ok := item["networkSiteTypeId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.NetworkSiteTypeId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.NetworkSiteTypeId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.NetworkSiteTypeId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.NetworkSiteTypeId = types.StringValue(string(jsonBytes))
        } else {
            data.NetworkSiteTypeId = types.StringNull()
        }
    } else if val, ok := item["networkSiteTypeId"].(string); ok {
        data.NetworkSiteTypeId = types.StringValue(val)
    } else {
        data.NetworkSiteTypeId = types.StringNull()
    }
    if obj, ok := item["parentSiteId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ParentSiteId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ParentSiteId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ParentSiteId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ParentSiteId = types.StringValue(string(jsonBytes))
        } else {
            data.ParentSiteId = types.StringNull()
        }
    } else if val, ok := item["parentSiteId"].(string); ok {
        data.ParentSiteId = types.StringValue(val)
    } else {
        data.ParentSiteId = types.StringNull()
    }
    if obj, ok := item["materializedPath"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MaterializedPath = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.MaterializedPath = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.MaterializedPath = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.MaterializedPath = types.StringValue(string(jsonBytes))
        } else {
            data.MaterializedPath = types.StringNull()
        }
    } else if val, ok := item["materializedPath"].(string); ok {
        data.MaterializedPath = types.StringValue(val)
    } else {
        data.MaterializedPath = types.StringNull()
    }
    if val, ok := item["depth"].(float64); ok {
        data.Depth = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["depth"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.Depth = types.NumberValue(big.NewFloat(val))
        } else {
            data.Depth = types.NumberNull()
        }
    } else {
        data.Depth = types.NumberNull()
    }
    if obj, ok := item["address"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Address = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Address = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Address = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Address = types.StringValue(string(jsonBytes))
        } else {
            data.Address = types.StringNull()
        }
    } else if val, ok := item["address"].(string); ok {
        data.Address = types.StringValue(val)
    } else {
        data.Address = types.StringNull()
    }
    if val, ok := item["latitude"].(float64); ok {
        data.Latitude = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["latitude"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.Latitude = types.NumberValue(big.NewFloat(val))
        } else {
            data.Latitude = types.NumberNull()
        }
    } else {
        data.Latitude = types.NumberNull()
    }
    if val, ok := item["longitude"].(float64); ok {
        data.Longitude = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["longitude"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.Longitude = types.NumberValue(big.NewFloat(val))
        } else {
            data.Longitude = types.NumberNull()
        }
    } else {
        data.Longitude = types.NumberNull()
    }
    if obj, ok := item["currentMonitorStatusId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CurrentMonitorStatusId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CurrentMonitorStatusId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CurrentMonitorStatusId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CurrentMonitorStatusId = types.StringValue(string(jsonBytes))
        } else {
            data.CurrentMonitorStatusId = types.StringNull()
        }
    } else if val, ok := item["currentMonitorStatusId"].(string); ok {
        data.CurrentMonitorStatusId = types.StringValue(val)
    } else {
        data.CurrentMonitorStatusId = types.StringNull()
    }
    if obj, ok := item["lastRollupAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastRollupAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastRollupAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastRollupAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastRollupAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastRollupAt = types.StringNull()
        }
    } else if val, ok := item["lastRollupAt"].(string); ok {
        data.LastRollupAt = types.StringValue(val)
    } else {
        data.LastRollupAt = types.StringNull()
    }
    if obj, ok := item["healthRollupPolicy"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.HealthRollupPolicy = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.HealthRollupPolicy = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.HealthRollupPolicy = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.HealthRollupPolicy = types.StringValue(string(jsonBytes))
        } else {
            data.HealthRollupPolicy = types.StringNull()
        }
    } else if val, ok := item["healthRollupPolicy"].(string); ok {
        data.HealthRollupPolicy = types.StringValue(val)
    } else {
        data.HealthRollupPolicy = types.StringNull()
    }
    if val, ok := item["offlineThresholdPercent"].(float64); ok {
        data.OfflineThresholdPercent = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["offlineThresholdPercent"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.OfflineThresholdPercent = types.NumberValue(big.NewFloat(val))
        } else {
            data.OfflineThresholdPercent = types.NumberNull()
        }
    } else {
        data.OfflineThresholdPercent = types.NumberNull()
    }
    if val, ok := item["shouldAlertWhenUnhealthy"].(bool); ok {
        data.ShouldAlertWhenUnhealthy = types.BoolValue(val)
    } else {
        data.ShouldAlertWhenUnhealthy = types.BoolNull()
    }
    if obj, ok := item["alertSeverityId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AlertSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AlertSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AlertSeverityId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AlertSeverityId = types.StringValue(string(jsonBytes))
        } else {
            data.AlertSeverityId = types.StringNull()
        }
    } else if val, ok := item["alertSeverityId"].(string); ok {
        data.AlertSeverityId = types.StringValue(val)
    } else {
        data.AlertSeverityId = types.StringNull()
    }
    if obj, ok := item["currentActiveAlertId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CurrentActiveAlertId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CurrentActiveAlertId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CurrentActiveAlertId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CurrentActiveAlertId = types.StringValue(string(jsonBytes))
        } else {
            data.CurrentActiveAlertId = types.StringNull()
        }
    } else if val, ok := item["currentActiveAlertId"].(string); ok {
        data.CurrentActiveAlertId = types.StringValue(val)
    } else {
        data.CurrentActiveAlertId = types.StringNull()
    }
    if obj, ok := item["createdByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if val, ok := item["createdByUserId"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if obj, ok := item["deletedByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeletedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DeletedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DeletedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DeletedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.DeletedByUserId = types.StringNull()
        }
    } else if val, ok := item["deletedByUserId"].(string); ok {
        data.DeletedByUserId = types.StringValue(val)
    } else {
        data.DeletedByUserId = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
