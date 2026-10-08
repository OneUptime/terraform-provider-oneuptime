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
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    Name types.String `tfsdk:"name"`
    Slug types.String `tfsdk:"slug"`
    Description types.String `tfsdk:"description"`
    SiteType types.String `tfsdk:"site_type"`
    NetworkSiteTypeId types.String `tfsdk:"network_site_type_id"`
    ProbeId types.String `tfsdk:"probe_id"`
    SnmpCredentialProfileId types.String `tfsdk:"snmp_credential_profile_id"`
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
}

func (d *NetworkSiteDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_network_site"
}

func (d *NetworkSiteDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Self-nesting sites (Account Type -> Region / Franchisee -> Market -> Unit) that group Network Devices into a drill-down hierarchy with a persisted health rollup. Look up an existing network site by `id`, or by any of its other arguments (`name`, `address`, `alert_severity_id`, ...): each one set must match, and exactly one network site may match them all.",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Look up by unique identifier. Leave unset to look up by the other arguments instead.",
                Optional: true,
                Computed: true,
            },
            "created_at": schema.StringAttribute{
                MarkdownDescription: "Date and Time when the object was created.",
                Computed: true,
            },
            "updated_at": schema.StringAttribute{
                MarkdownDescription: "Date and Time when the object was updated.",
                Computed: true,
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Friendly name for this network site.",
                Optional: true,
                Computed: true,
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Friendly globally unique name for your object.",
                Optional: true,
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Friendly description for this network site.",
                Optional: true,
                Computed: true,
            },
            "site_type": schema.StringAttribute{
                MarkdownDescription: "Deprecated legacy site type string. Use the Network Site Type relation instead; this column exists only for the backfill migration and will be removed.",
                Optional: true,
                Computed: true,
            },
            "network_site_type_id": schema.StringAttribute{
                MarkdownDescription: "ID of the Network Site Type this site belongs to. The ID of a `oneuptime_network_site_type`.",
                Optional: true,
                Computed: true,
            },
            "probe_id": schema.StringAttribute{
                MarkdownDescription: "ID of the probe that polls devices in this site by default. The ID of a `oneuptime_probe`.",
                Optional: true,
                Computed: true,
            },
            "snmp_credential_profile_id": schema.StringAttribute{
                MarkdownDescription: "ID of the SNMP Credential Profile devices in this site inherit. The ID of a `oneuptime_snmp_credential_profile`.",
                Optional: true,
                Computed: true,
            },
            "parent_site_id": schema.StringAttribute{
                MarkdownDescription: "ID of the parent Network Site this site is nested under (empty for root sites). The ID of a `oneuptime_network_site`.",
                Optional: true,
                Computed: true,
            },
            "materialized_path": schema.StringAttribute{
                MarkdownDescription: "Slash-separated ancestor IDs of this site (e.g. '/rootId/childId/'). Managed by the server on parent changes; used for subtree queries and rollups.",
                Optional: true,
                Computed: true,
            },
            "depth": schema.NumberAttribute{
                MarkdownDescription: "Number of ancestors above this site (0 for root sites). Managed by the server on parent changes.",
                Optional: true,
                Computed: true,
            },
            "address": schema.StringAttribute{
                MarkdownDescription: "Street address of this site, shown on map views.",
                Optional: true,
                Computed: true,
            },
            "latitude": schema.NumberAttribute{
                MarkdownDescription: "Latitude of this site, for US and world map views.",
                Optional: true,
                Computed: true,
            },
            "longitude": schema.NumberAttribute{
                MarkdownDescription: "Longitude of this site, for US and world map views.",
                Optional: true,
                Computed: true,
            },
            "current_monitor_status_id": schema.StringAttribute{
                MarkdownDescription: "Whats the current rolled-up status ID of this site? Computed from the devices and child sites below it. The ID of a `oneuptime_monitor_status`.",
                Optional: true,
                Computed: true,
            },
            "last_rollup_at": schema.StringAttribute{
                MarkdownDescription: "When the health rollup for this site was last computed.",
                Computed: true,
            },
            "health_rollup_policy": schema.StringAttribute{
                MarkdownDescription: "How this site's status is derived from the devices beneath it: WorstStatus (any device offline makes the site offline) or PercentThreshold (the share of devices that are down decides).",
                Optional: true,
                Computed: true,
            },
            "offline_threshold_percent": schema.NumberAttribute{
                MarkdownDescription: "With the PercentThreshold rollup policy: the share of reporting devices beneath this site that must be non-operational before the site itself is marked offline. Below it (but above zero) the site is degraded.",
                Optional: true,
                Computed: true,
            },
            "should_alert_when_unhealthy": schema.BoolAttribute{
                MarkdownDescription: "When enabled, an alert opens when this site's health rollup turns non-operational and auto-resolves when it recovers.",
                Optional: true,
                Computed: true,
            },
            "alert_severity_id": schema.StringAttribute{
                MarkdownDescription: "ID of the severity used for site-unhealthy alerts. The ID of a `oneuptime_alert_severity`.",
                Optional: true,
                Computed: true,
            },
            "current_active_alert_id": schema.StringAttribute{
                MarkdownDescription: "ID of the currently open site-unhealthy alert, if any. Managed by the rollup engine.",
                Optional: true,
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
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

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.Name.IsNull() && !data.Name.IsUnknown() {
        filters["name"] = data.Name.ValueString()
        filterNames = append(filterNames, "name = "+fmt.Sprintf("%q", data.Name.ValueString()))
    }
    if !data.Slug.IsNull() && !data.Slug.IsUnknown() {
        filters["slug"] = data.Slug.ValueString()
        filterNames = append(filterNames, "slug = "+fmt.Sprintf("%q", data.Slug.ValueString()))
    }
    if !data.Description.IsNull() && !data.Description.IsUnknown() {
        filters["description"] = data.Description.ValueString()
        filterNames = append(filterNames, "description = "+fmt.Sprintf("%q", data.Description.ValueString()))
    }
    if !data.SiteType.IsNull() && !data.SiteType.IsUnknown() {
        filters["siteType"] = data.SiteType.ValueString()
        filterNames = append(filterNames, "site_type = "+fmt.Sprintf("%q", data.SiteType.ValueString()))
    }
    if !data.NetworkSiteTypeId.IsNull() && !data.NetworkSiteTypeId.IsUnknown() {
        filters["networkSiteTypeId"] = data.NetworkSiteTypeId.ValueString()
        filterNames = append(filterNames, "network_site_type_id = "+fmt.Sprintf("%q", data.NetworkSiteTypeId.ValueString()))
    }
    if !data.ProbeId.IsNull() && !data.ProbeId.IsUnknown() {
        filters["probeId"] = data.ProbeId.ValueString()
        filterNames = append(filterNames, "probe_id = "+fmt.Sprintf("%q", data.ProbeId.ValueString()))
    }
    if !data.SnmpCredentialProfileId.IsNull() && !data.SnmpCredentialProfileId.IsUnknown() {
        filters["snmpCredentialProfileId"] = data.SnmpCredentialProfileId.ValueString()
        filterNames = append(filterNames, "snmp_credential_profile_id = "+fmt.Sprintf("%q", data.SnmpCredentialProfileId.ValueString()))
    }
    if !data.ParentSiteId.IsNull() && !data.ParentSiteId.IsUnknown() {
        filters["parentSiteId"] = data.ParentSiteId.ValueString()
        filterNames = append(filterNames, "parent_site_id = "+fmt.Sprintf("%q", data.ParentSiteId.ValueString()))
    }
    if !data.MaterializedPath.IsNull() && !data.MaterializedPath.IsUnknown() {
        filters["materializedPath"] = data.MaterializedPath.ValueString()
        filterNames = append(filterNames, "materialized_path = "+fmt.Sprintf("%q", data.MaterializedPath.ValueString()))
    }
    if !data.Depth.IsNull() && !data.Depth.IsUnknown() {
        filters["depth"] = lookupNumber(data.Depth)
        filterNames = append(filterNames, "depth = "+data.Depth.ValueBigFloat().String())
    }
    if !data.Address.IsNull() && !data.Address.IsUnknown() {
        filters["address"] = data.Address.ValueString()
        filterNames = append(filterNames, "address = "+fmt.Sprintf("%q", data.Address.ValueString()))
    }
    if !data.Latitude.IsNull() && !data.Latitude.IsUnknown() {
        filters["latitude"] = lookupNumber(data.Latitude)
        filterNames = append(filterNames, "latitude = "+data.Latitude.ValueBigFloat().String())
    }
    if !data.Longitude.IsNull() && !data.Longitude.IsUnknown() {
        filters["longitude"] = lookupNumber(data.Longitude)
        filterNames = append(filterNames, "longitude = "+data.Longitude.ValueBigFloat().String())
    }
    if !data.CurrentMonitorStatusId.IsNull() && !data.CurrentMonitorStatusId.IsUnknown() {
        filters["currentMonitorStatusId"] = data.CurrentMonitorStatusId.ValueString()
        filterNames = append(filterNames, "current_monitor_status_id = "+fmt.Sprintf("%q", data.CurrentMonitorStatusId.ValueString()))
    }
    if !data.HealthRollupPolicy.IsNull() && !data.HealthRollupPolicy.IsUnknown() {
        filters["healthRollupPolicy"] = data.HealthRollupPolicy.ValueString()
        filterNames = append(filterNames, "health_rollup_policy = "+fmt.Sprintf("%q", data.HealthRollupPolicy.ValueString()))
    }
    if !data.OfflineThresholdPercent.IsNull() && !data.OfflineThresholdPercent.IsUnknown() {
        filters["offlineThresholdPercent"] = lookupNumber(data.OfflineThresholdPercent)
        filterNames = append(filterNames, "offline_threshold_percent = "+data.OfflineThresholdPercent.ValueBigFloat().String())
    }
    if !data.ShouldAlertWhenUnhealthy.IsNull() && !data.ShouldAlertWhenUnhealthy.IsUnknown() {
        filters["shouldAlertWhenUnhealthy"] = data.ShouldAlertWhenUnhealthy.ValueBool()
        filterNames = append(filterNames, "should_alert_when_unhealthy = "+fmt.Sprintf("%t", data.ShouldAlertWhenUnhealthy.ValueBool()))
    }
    if !data.AlertSeverityId.IsNull() && !data.AlertSeverityId.IsUnknown() {
        filters["alertSeverityId"] = data.AlertSeverityId.ValueString()
        filterNames = append(filterNames, "alert_severity_id = "+fmt.Sprintf("%q", data.AlertSeverityId.ValueString()))
    }
    if !data.CurrentActiveAlertId.IsNull() && !data.CurrentActiveAlertId.IsUnknown() {
        filters["currentActiveAlertId"] = data.CurrentActiveAlertId.ValueString()
        filterNames = append(filterNames, "current_active_alert_id = "+fmt.Sprintf("%q", data.CurrentActiveAlertId.ValueString()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the network site up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the network site up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "name": true,
        "slug": true,
        "description": true,
        "siteType": true,
        "networkSiteTypeId": true,
        "probeId": true,
        "snmpCredentialProfileId": true,
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
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No network site found with id %q.", data.Id.ValueString()))
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
    }
    if !hasId {
        listBody := map[string]interface{}{
            "query":  filters,
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
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No network site matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one network site matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
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
    if obj, ok := item["probeId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ProbeId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ProbeId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ProbeId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ProbeId = types.StringValue(string(jsonBytes))
        } else {
            data.ProbeId = types.StringNull()
        }
    } else if val, ok := item["probeId"].(string); ok {
        data.ProbeId = types.StringValue(val)
    } else {
        data.ProbeId = types.StringNull()
    }
    if obj, ok := item["snmpCredentialProfileId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SnmpCredentialProfileId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SnmpCredentialProfileId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SnmpCredentialProfileId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SnmpCredentialProfileId = types.StringValue(string(jsonBytes))
        } else {
            data.SnmpCredentialProfileId = types.StringNull()
        }
    } else if val, ok := item["snmpCredentialProfileId"].(string); ok {
        data.SnmpCredentialProfileId = types.StringValue(val)
    } else {
        data.SnmpCredentialProfileId = types.StringNull()
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

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
