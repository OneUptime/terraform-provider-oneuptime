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
var _ datasource.DataSource = &StatusPageResourceDataSource{}

func NewStatusPageResourceDataSource() datasource.DataSource {
    return &StatusPageResourceDataSource{}
}

// StatusPageResourceDataSource defines the data source implementation.
type StatusPageResourceDataSource struct {
    client *Client
}

// StatusPageResourceDataSourceModel describes the data source data model.
type StatusPageResourceDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    StatusPageId types.String `tfsdk:"status_page_id"`
    MonitorId types.String `tfsdk:"monitor_id"`
    MonitorGroupId types.String `tfsdk:"monitor_group_id"`
    StatusPageGroupId types.String `tfsdk:"status_page_group_id"`
    StatusPageMonitorRuleId types.String `tfsdk:"status_page_monitor_rule_id"`
    DisplayName types.String `tfsdk:"display_name"`
    DisplayDescription types.String `tfsdk:"display_description"`
    DisplayTooltip types.String `tfsdk:"display_tooltip"`
    ShowCurrentStatus types.Bool `tfsdk:"show_current_status"`
    ShowUptimePercent types.Bool `tfsdk:"show_uptime_percent"`
    UptimePercentPrecision types.String `tfsdk:"uptime_percent_precision"`
    ShowStatusHistoryChart types.Bool `tfsdk:"show_status_history_chart"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    Order types.Number `tfsdk:"order"`
    RowAxisValue types.String `tfsdk:"row_axis_value"`
    ColumnAxisValue types.String `tfsdk:"column_axis_value"`
}

func (d *StatusPageResourceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_status_page_resource"
}

func (d *StatusPageResourceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Add resources like monitors to your status page Look up an existing status page resource by `id`, or by any of its other arguments (`column_axis_value`, `created_by_user_id`, `display_description`, ...): each one set must match, and exactly one status page resource may match them all.",

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
            "status_page_id": schema.StringAttribute{
                MarkdownDescription: "ID of your Status Page resource where this object belongs. The ID of a `oneuptime_status_page`.",
                Optional: true,
                Computed: true,
            },
            "monitor_id": schema.StringAttribute{
                MarkdownDescription: "Relation to Monitor ID Resource in which this object belongs. The ID of a `oneuptime_monitor`.",
                Optional: true,
                Computed: true,
            },
            "monitor_group_id": schema.StringAttribute{
                MarkdownDescription: "Relation to Monitor Group ID Resource in which this object belongs. The ID of a `oneuptime_monitor_group`.",
                Optional: true,
                Computed: true,
            },
            "status_page_group_id": schema.StringAttribute{
                MarkdownDescription: "Does this monitor belong to a status page group? The ID of a `oneuptime_status_page_group`.",
                Optional: true,
                Computed: true,
            },
            "status_page_monitor_rule_id": schema.StringAttribute{
                MarkdownDescription: "ID of the rule that added this resource, if it was added by a rule instead of by hand. The ID of a `oneuptime_status_page_monitor_rule`.",
                Optional: true,
                Computed: true,
            },
            "display_name": schema.StringAttribute{
                MarkdownDescription: "Display name of the monitor on the Status Page.",
                Optional: true,
                Computed: true,
            },
            "display_description": schema.StringAttribute{
                MarkdownDescription: "Display description of the monitor on the Status Page. This is in markdown format.",
                Optional: true,
                Computed: true,
            },
            "display_tooltip": schema.StringAttribute{
                MarkdownDescription: "Tooltip of the monitor on the Status Page.",
                Optional: true,
                Computed: true,
            },
            "show_current_status": schema.BoolAttribute{
                MarkdownDescription: "Show current status like offline, operational or degraded.",
                Optional: true,
                Computed: true,
            },
            "show_uptime_percent": schema.BoolAttribute{
                MarkdownDescription: "Show uptime percent of this monitor for the last 90 days.",
                Optional: true,
                Computed: true,
            },
            "uptime_percent_precision": schema.StringAttribute{
                MarkdownDescription: "Precision of uptime percent of this monitor for the last 90 days.",
                Optional: true,
                Computed: true,
            },
            "show_status_history_chart": schema.BoolAttribute{
                MarkdownDescription: "Show a 90 day uptime history of this monitor.",
                Optional: true,
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "order": schema.NumberAttribute{
                MarkdownDescription: "Order / Priority of this resource.",
                Optional: true,
                Computed: true,
            },
            "row_axis_value": schema.StringAttribute{
                MarkdownDescription: "Row this resource belongs to when its status page group is rendered as a grid. Should match one of the row axis values defined on the group.",
                Optional: true,
                Computed: true,
            },
            "column_axis_value": schema.StringAttribute{
                MarkdownDescription: "Column this resource belongs to when its status page group is rendered as a grid. Should match one of the column axis values defined on the group.",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *StatusPageResourceDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *StatusPageResourceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data StatusPageResourceDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.StatusPageId.IsNull() && !data.StatusPageId.IsUnknown() {
        filters["statusPageId"] = data.StatusPageId.ValueString()
        filterNames = append(filterNames, "status_page_id = "+fmt.Sprintf("%q", data.StatusPageId.ValueString()))
    }
    if !data.MonitorId.IsNull() && !data.MonitorId.IsUnknown() {
        filters["monitorId"] = data.MonitorId.ValueString()
        filterNames = append(filterNames, "monitor_id = "+fmt.Sprintf("%q", data.MonitorId.ValueString()))
    }
    if !data.MonitorGroupId.IsNull() && !data.MonitorGroupId.IsUnknown() {
        filters["monitorGroupId"] = data.MonitorGroupId.ValueString()
        filterNames = append(filterNames, "monitor_group_id = "+fmt.Sprintf("%q", data.MonitorGroupId.ValueString()))
    }
    if !data.StatusPageGroupId.IsNull() && !data.StatusPageGroupId.IsUnknown() {
        filters["statusPageGroupId"] = data.StatusPageGroupId.ValueString()
        filterNames = append(filterNames, "status_page_group_id = "+fmt.Sprintf("%q", data.StatusPageGroupId.ValueString()))
    }
    if !data.StatusPageMonitorRuleId.IsNull() && !data.StatusPageMonitorRuleId.IsUnknown() {
        filters["statusPageMonitorRuleId"] = data.StatusPageMonitorRuleId.ValueString()
        filterNames = append(filterNames, "status_page_monitor_rule_id = "+fmt.Sprintf("%q", data.StatusPageMonitorRuleId.ValueString()))
    }
    if !data.DisplayName.IsNull() && !data.DisplayName.IsUnknown() {
        filters["displayName"] = data.DisplayName.ValueString()
        filterNames = append(filterNames, "display_name = "+fmt.Sprintf("%q", data.DisplayName.ValueString()))
    }
    if !data.DisplayDescription.IsNull() && !data.DisplayDescription.IsUnknown() {
        filters["displayDescription"] = data.DisplayDescription.ValueString()
        filterNames = append(filterNames, "display_description = "+fmt.Sprintf("%q", data.DisplayDescription.ValueString()))
    }
    if !data.DisplayTooltip.IsNull() && !data.DisplayTooltip.IsUnknown() {
        filters["displayTooltip"] = data.DisplayTooltip.ValueString()
        filterNames = append(filterNames, "display_tooltip = "+fmt.Sprintf("%q", data.DisplayTooltip.ValueString()))
    }
    if !data.ShowCurrentStatus.IsNull() && !data.ShowCurrentStatus.IsUnknown() {
        filters["showCurrentStatus"] = data.ShowCurrentStatus.ValueBool()
        filterNames = append(filterNames, "show_current_status = "+fmt.Sprintf("%t", data.ShowCurrentStatus.ValueBool()))
    }
    if !data.ShowUptimePercent.IsNull() && !data.ShowUptimePercent.IsUnknown() {
        filters["showUptimePercent"] = data.ShowUptimePercent.ValueBool()
        filterNames = append(filterNames, "show_uptime_percent = "+fmt.Sprintf("%t", data.ShowUptimePercent.ValueBool()))
    }
    if !data.UptimePercentPrecision.IsNull() && !data.UptimePercentPrecision.IsUnknown() {
        filters["uptimePercentPrecision"] = data.UptimePercentPrecision.ValueString()
        filterNames = append(filterNames, "uptime_percent_precision = "+fmt.Sprintf("%q", data.UptimePercentPrecision.ValueString()))
    }
    if !data.ShowStatusHistoryChart.IsNull() && !data.ShowStatusHistoryChart.IsUnknown() {
        filters["showStatusHistoryChart"] = data.ShowStatusHistoryChart.ValueBool()
        filterNames = append(filterNames, "show_status_history_chart = "+fmt.Sprintf("%t", data.ShowStatusHistoryChart.ValueBool()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }
    if !data.Order.IsNull() && !data.Order.IsUnknown() {
        filters["order"] = lookupNumber(data.Order)
        filterNames = append(filterNames, "order = "+data.Order.ValueBigFloat().String())
    }
    if !data.RowAxisValue.IsNull() && !data.RowAxisValue.IsUnknown() {
        filters["rowAxisValue"] = data.RowAxisValue.ValueString()
        filterNames = append(filterNames, "row_axis_value = "+fmt.Sprintf("%q", data.RowAxisValue.ValueString()))
    }
    if !data.ColumnAxisValue.IsNull() && !data.ColumnAxisValue.IsUnknown() {
        filters["columnAxisValue"] = data.ColumnAxisValue.ValueString()
        filterNames = append(filterNames, "column_axis_value = "+fmt.Sprintf("%q", data.ColumnAxisValue.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the status page resource up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the status page resource up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "statusPageId": true,
        "monitorId": true,
        "monitorGroupId": true,
        "statusPageGroupId": true,
        "statusPageMonitorRuleId": true,
        "displayName": true,
        "displayDescription": true,
        "displayTooltip": true,
        "showCurrentStatus": true,
        "showUptimePercent": true,
        "uptimePercentPrecision": true,
        "showStatusHistoryChart": true,
        "createdByUserId": true,
        "order": true,
        "rowAxisValue": true,
        "columnAxisValue": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/status-page-resource/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read status_page_resource, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No status page resource found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read status_page_resource: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/status-page-resource/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list status_page_resource, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list status_page_resource: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No status page resource matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one status page resource matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for status_page_resource.")
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
    if obj, ok := item["statusPageId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StatusPageId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.StatusPageId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.StatusPageId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.StatusPageId = types.StringValue(string(jsonBytes))
        } else {
            data.StatusPageId = types.StringNull()
        }
    } else if val, ok := item["statusPageId"].(string); ok {
        data.StatusPageId = types.StringValue(val)
    } else {
        data.StatusPageId = types.StringNull()
    }
    if obj, ok := item["monitorId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MonitorId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.MonitorId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.MonitorId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.MonitorId = types.StringValue(string(jsonBytes))
        } else {
            data.MonitorId = types.StringNull()
        }
    } else if val, ok := item["monitorId"].(string); ok {
        data.MonitorId = types.StringValue(val)
    } else {
        data.MonitorId = types.StringNull()
    }
    if obj, ok := item["monitorGroupId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MonitorGroupId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.MonitorGroupId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.MonitorGroupId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.MonitorGroupId = types.StringValue(string(jsonBytes))
        } else {
            data.MonitorGroupId = types.StringNull()
        }
    } else if val, ok := item["monitorGroupId"].(string); ok {
        data.MonitorGroupId = types.StringValue(val)
    } else {
        data.MonitorGroupId = types.StringNull()
    }
    if obj, ok := item["statusPageGroupId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StatusPageGroupId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.StatusPageGroupId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.StatusPageGroupId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.StatusPageGroupId = types.StringValue(string(jsonBytes))
        } else {
            data.StatusPageGroupId = types.StringNull()
        }
    } else if val, ok := item["statusPageGroupId"].(string); ok {
        data.StatusPageGroupId = types.StringValue(val)
    } else {
        data.StatusPageGroupId = types.StringNull()
    }
    if obj, ok := item["statusPageMonitorRuleId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StatusPageMonitorRuleId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.StatusPageMonitorRuleId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.StatusPageMonitorRuleId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.StatusPageMonitorRuleId = types.StringValue(string(jsonBytes))
        } else {
            data.StatusPageMonitorRuleId = types.StringNull()
        }
    } else if val, ok := item["statusPageMonitorRuleId"].(string); ok {
        data.StatusPageMonitorRuleId = types.StringValue(val)
    } else {
        data.StatusPageMonitorRuleId = types.StringNull()
    }
    if obj, ok := item["displayName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DisplayName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DisplayName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DisplayName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DisplayName = types.StringValue(string(jsonBytes))
        } else {
            data.DisplayName = types.StringNull()
        }
    } else if val, ok := item["displayName"].(string); ok {
        data.DisplayName = types.StringValue(val)
    } else {
        data.DisplayName = types.StringNull()
    }
    if obj, ok := item["displayDescription"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DisplayDescription = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DisplayDescription = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DisplayDescription = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DisplayDescription = types.StringValue(string(jsonBytes))
        } else {
            data.DisplayDescription = types.StringNull()
        }
    } else if val, ok := item["displayDescription"].(string); ok {
        data.DisplayDescription = types.StringValue(val)
    } else {
        data.DisplayDescription = types.StringNull()
    }
    if obj, ok := item["displayTooltip"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DisplayTooltip = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DisplayTooltip = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DisplayTooltip = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DisplayTooltip = types.StringValue(string(jsonBytes))
        } else {
            data.DisplayTooltip = types.StringNull()
        }
    } else if val, ok := item["displayTooltip"].(string); ok {
        data.DisplayTooltip = types.StringValue(val)
    } else {
        data.DisplayTooltip = types.StringNull()
    }
    if val, ok := item["showCurrentStatus"].(bool); ok {
        data.ShowCurrentStatus = types.BoolValue(val)
    } else {
        data.ShowCurrentStatus = types.BoolNull()
    }
    if val, ok := item["showUptimePercent"].(bool); ok {
        data.ShowUptimePercent = types.BoolValue(val)
    } else {
        data.ShowUptimePercent = types.BoolNull()
    }
    if obj, ok := item["uptimePercentPrecision"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UptimePercentPrecision = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.UptimePercentPrecision = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.UptimePercentPrecision = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.UptimePercentPrecision = types.StringValue(string(jsonBytes))
        } else {
            data.UptimePercentPrecision = types.StringNull()
        }
    } else if val, ok := item["uptimePercentPrecision"].(string); ok {
        data.UptimePercentPrecision = types.StringValue(val)
    } else {
        data.UptimePercentPrecision = types.StringNull()
    }
    if val, ok := item["showStatusHistoryChart"].(bool); ok {
        data.ShowStatusHistoryChart = types.BoolValue(val)
    } else {
        data.ShowStatusHistoryChart = types.BoolNull()
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
    if val, ok := item["order"].(float64); ok {
        data.Order = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["order"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.Order = types.NumberValue(big.NewFloat(val))
        } else {
            data.Order = types.NumberNull()
        }
    } else {
        data.Order = types.NumberNull()
    }
    if obj, ok := item["rowAxisValue"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RowAxisValue = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RowAxisValue = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RowAxisValue = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RowAxisValue = types.StringValue(string(jsonBytes))
        } else {
            data.RowAxisValue = types.StringNull()
        }
    } else if val, ok := item["rowAxisValue"].(string); ok {
        data.RowAxisValue = types.StringValue(val)
    } else {
        data.RowAxisValue = types.StringNull()
    }
    if obj, ok := item["columnAxisValue"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ColumnAxisValue = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ColumnAxisValue = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ColumnAxisValue = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ColumnAxisValue = types.StringValue(string(jsonBytes))
        } else {
            data.ColumnAxisValue = types.StringNull()
        }
    } else if val, ok := item["columnAxisValue"].(string); ok {
        data.ColumnAxisValue = types.StringValue(val)
    } else {
        data.ColumnAxisValue = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
