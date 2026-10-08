package provider

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "math/big"
    "github.com/hashicorp/terraform-plugin-framework/attr"
    "sort"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ServiceLevelObjectiveDataSource{}

func NewServiceLevelObjectiveDataSource() datasource.DataSource {
    return &ServiceLevelObjectiveDataSource{}
}

// ServiceLevelObjectiveDataSource defines the data source implementation.
type ServiceLevelObjectiveDataSource struct {
    client *Client
}

// ServiceLevelObjectiveDataSourceModel describes the data source data model.
type ServiceLevelObjectiveDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    Name types.String `tfsdk:"name"`
    Description types.String `tfsdk:"description"`
    Slug types.String `tfsdk:"slug"`
    Labels types.Set `tfsdk:"labels"`
    IsEnabled types.Bool `tfsdk:"is_enabled"`
    IsArchived types.Bool `tfsdk:"is_archived"`
    ArchivedAt types.String `tfsdk:"archived_at"`
    ArchivedByUserId types.String `tfsdk:"archived_by_user_id"`
    SliType types.String `tfsdk:"sli_type"`
    MultiMonitorMode types.String `tfsdk:"multi_monitor_mode"`
    Monitors types.Set `tfsdk:"monitors"`
    MonitorLabels types.Set `tfsdk:"monitor_labels"`
    AutoAddedMonitors types.Set `tfsdk:"auto_added_monitors"`
    DowntimeMonitorStatuses types.Set `tfsdk:"downtime_monitor_statuses"`
    MetricQueryConfig types.String `tfsdk:"metric_query_config"`
    TargetPercentage types.Number `tfsdk:"target_percentage"`
    WindowType types.String `tfsdk:"window_type"`
    WindowDays types.Number `tfsdk:"window_days"`
    Timezone types.String `tfsdk:"timezone"`
    AtRiskThresholdPercentage types.Number `tfsdk:"at_risk_threshold_percentage"`
    CurrentSliPercentage types.Number `tfsdk:"current_sli_percentage"`
    ErrorBudgetRemainingPercentage types.Number `tfsdk:"error_budget_remaining_percentage"`
    ErrorBudgetRemainingSeconds types.Number `tfsdk:"error_budget_remaining_seconds"`
    ErrorBudgetTotalSeconds types.Number `tfsdk:"error_budget_total_seconds"`
    CurrentBurnRate types.Number `tfsdk:"current_burn_rate"`
    SloStatus types.String `tfsdk:"slo_status"`
    StatusChangeNotificationSentAt types.String `tfsdk:"status_change_notification_sent_at"`
    LastEvaluatedAt types.String `tfsdk:"last_evaluated_at"`
    NextEvaluationAt types.String `tfsdk:"next_evaluation_at"`
    LastAccumulatedBucketEndAt types.String `tfsdk:"last_accumulated_bucket_end_at"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *ServiceLevelObjectiveDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_service_level_objective"
}

func (d *ServiceLevelObjectiveDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Define Service Level Objectives (SLOs) with targets, compliance windows and error budgets, and track how much error budget remains. Look up an existing service level objective by `id`, or by any of its other arguments (`name`, `archived_by_user_id`, `at_risk_threshold_percentage`, ...): each one set must match, and exactly one service level objective may match them all.",

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
                MarkdownDescription: "Name of this Service Level Objective.",
                Optional: true,
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Description of this Service Level Objective.",
                Optional: true,
                Computed: true,
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Friendly globally unique name for your object.",
                Optional: true,
                Computed: true,
            },
            "labels": schema.SetAttribute{
                MarkdownDescription: "Relation to Labels Array where this object is categorized in. IDs of `oneuptime_label` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether this Service Level Objective is enabled. Disabled SLOs are not evaluated.",
                Optional: true,
                Computed: true,
            },
            "is_archived": schema.BoolAttribute{
                MarkdownDescription: "Archived SLOs are hidden from lists and are not evaluated.",
                Optional: true,
                Computed: true,
            },
            "archived_at": schema.StringAttribute{
                MarkdownDescription: "When this Service Level Objective was archived.",
                Computed: true,
            },
            "archived_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who archived this object (if this object was archived by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "sli_type": schema.StringAttribute{
                MarkdownDescription: "Type of Service Level Indicator this objective measures (Monitor Uptime or Metric).",
                Optional: true,
                Computed: true,
            },
            "multi_monitor_mode": schema.StringAttribute{
                MarkdownDescription: "How downtime is counted when multiple monitors are attached. 'Any Monitor Down' counts time when any monitor is down. 'Monitor Seconds Average' averages downtime across monitors.",
                Optional: true,
                Computed: true,
            },
            "monitors": schema.SetAttribute{
                MarkdownDescription: "Monitors whose uptime is measured by this Service Level Objective (for Monitor Uptime SLIs). IDs of `oneuptime_monitor` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "monitor_labels": schema.SetAttribute{
                MarkdownDescription: "Deprecated: superseded by SLO Monitor Rules and no longer read by the SLO engine. Existing labels were migrated into a monitor rule named \"Auto-add monitors with labels\". Kept only for compatibility during upgrades: labels written here to an SLO with no monitor rules are turned into that rule, and are ignored once the SLO has monitor rules. Use SLO Monitor Rules instead. IDs of `oneuptime_label` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "auto_added_monitors": schema.SetAttribute{
                MarkdownDescription: "Monitors that were attached to this SLO by its monitor rules rather than by hand. Maintained by the server. IDs of `oneuptime_monitor` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "downtime_monitor_statuses": schema.SetAttribute{
                MarkdownDescription: "List of monitor statuses that are considered as \"down\" for this Service Level Objective. IDs of `oneuptime_monitor_status` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "metric_query_config": schema.StringAttribute{
                MarkdownDescription: "Query configuration for Metric SLIs: metric name, good-event predicate and optional attribute filters. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "target_percentage": schema.NumberAttribute{
                MarkdownDescription: "Target of this Service Level Objective as a percentage (e.g. 99.9). Must be less than 100.",
                Optional: true,
                Computed: true,
            },
            "window_type": schema.StringAttribute{
                MarkdownDescription: "Type of compliance window for this objective (Rolling or Calendar Month).",
                Optional: true,
                Computed: true,
            },
            "window_days": schema.NumberAttribute{
                MarkdownDescription: "Length of the rolling compliance window in days (e.g. 7, 28, 30 or 90). Ignored for Calendar Month windows.",
                Optional: true,
                Computed: true,
            },
            "timezone": schema.StringAttribute{
                MarkdownDescription: "IANA timezone (e.g. America/New_York) used for Calendar Month window boundaries. Defaults to UTC when not set.",
                Optional: true,
                Computed: true,
            },
            "at_risk_threshold_percentage": schema.NumberAttribute{
                MarkdownDescription: "Percentage of remaining error budget at which the SLO status changes to At Risk. For example, 20 means the status becomes At Risk when less than 20% of the error budget remains.",
                Optional: true,
                Computed: true,
            },
            "current_sli_percentage": schema.NumberAttribute{
                MarkdownDescription: "Current Service Level Indicator over the compliance window, as a percentage. Computed by the worker.",
                Optional: true,
                Computed: true,
            },
            "error_budget_remaining_percentage": schema.NumberAttribute{
                MarkdownDescription: "Percentage of the error budget that remains. Can be negative when the budget is exhausted. Computed by the worker.",
                Optional: true,
                Computed: true,
            },
            "error_budget_remaining_seconds": schema.NumberAttribute{
                MarkdownDescription: "Seconds of error budget that remain. Can be negative when the budget is exhausted. Computed by the worker.",
                Optional: true,
                Computed: true,
            },
            "error_budget_total_seconds": schema.NumberAttribute{
                MarkdownDescription: "Total seconds of error budget for the compliance window. Computed by the worker.",
                Optional: true,
                Computed: true,
            },
            "current_burn_rate": schema.NumberAttribute{
                MarkdownDescription: "Rate at which the error budget is currently being consumed. A burn rate of 1 exhausts the budget exactly at the end of the window. Computed by the worker.",
                Optional: true,
                Computed: true,
            },
            "slo_status": schema.StringAttribute{
                MarkdownDescription: "Current status of this Service Level Objective (Healthy, At Risk, Budget Exhausted, Misconfigured, Paused). Computed by the worker.",
                Optional: true,
                Computed: true,
            },
            "status_change_notification_sent_at": schema.StringAttribute{
                MarkdownDescription: "The last time a status-change notification was sent to owners. Computed by the worker.",
                Computed: true,
            },
            "last_evaluated_at": schema.StringAttribute{
                MarkdownDescription: "The last time this Service Level Objective was evaluated. Computed by the worker.",
                Computed: true,
            },
            "next_evaluation_at": schema.StringAttribute{
                MarkdownDescription: "When this Service Level Objective is next due for evaluation. Computed by the worker.",
                Computed: true,
            },
            "last_accumulated_bucket_end_at": schema.StringAttribute{
                MarkdownDescription: "Accumulation cursor for Metric SLIs: end of the last bucket whose good/total counts were persisted. Computed by the worker.",
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

func (d *ServiceLevelObjectiveDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ServiceLevelObjectiveDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data ServiceLevelObjectiveDataSourceModel

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
    if !data.Description.IsNull() && !data.Description.IsUnknown() {
        filters["description"] = data.Description.ValueString()
        filterNames = append(filterNames, "description = "+fmt.Sprintf("%q", data.Description.ValueString()))
    }
    if !data.Slug.IsNull() && !data.Slug.IsUnknown() {
        filters["slug"] = data.Slug.ValueString()
        filterNames = append(filterNames, "slug = "+fmt.Sprintf("%q", data.Slug.ValueString()))
    }
    if !data.IsEnabled.IsNull() && !data.IsEnabled.IsUnknown() {
        filters["isEnabled"] = data.IsEnabled.ValueBool()
        filterNames = append(filterNames, "is_enabled = "+fmt.Sprintf("%t", data.IsEnabled.ValueBool()))
    }
    if !data.IsArchived.IsNull() && !data.IsArchived.IsUnknown() {
        filters["isArchived"] = data.IsArchived.ValueBool()
        filterNames = append(filterNames, "is_archived = "+fmt.Sprintf("%t", data.IsArchived.ValueBool()))
    }
    if !data.ArchivedByUserId.IsNull() && !data.ArchivedByUserId.IsUnknown() {
        filters["archivedByUserId"] = data.ArchivedByUserId.ValueString()
        filterNames = append(filterNames, "archived_by_user_id = "+fmt.Sprintf("%q", data.ArchivedByUserId.ValueString()))
    }
    if !data.SliType.IsNull() && !data.SliType.IsUnknown() {
        filters["sliType"] = data.SliType.ValueString()
        filterNames = append(filterNames, "sli_type = "+fmt.Sprintf("%q", data.SliType.ValueString()))
    }
    if !data.MultiMonitorMode.IsNull() && !data.MultiMonitorMode.IsUnknown() {
        filters["multiMonitorMode"] = data.MultiMonitorMode.ValueString()
        filterNames = append(filterNames, "multi_monitor_mode = "+fmt.Sprintf("%q", data.MultiMonitorMode.ValueString()))
    }
    if !data.TargetPercentage.IsNull() && !data.TargetPercentage.IsUnknown() {
        filters["targetPercentage"] = lookupNumber(data.TargetPercentage)
        filterNames = append(filterNames, "target_percentage = "+data.TargetPercentage.ValueBigFloat().String())
    }
    if !data.WindowType.IsNull() && !data.WindowType.IsUnknown() {
        filters["windowType"] = data.WindowType.ValueString()
        filterNames = append(filterNames, "window_type = "+fmt.Sprintf("%q", data.WindowType.ValueString()))
    }
    if !data.WindowDays.IsNull() && !data.WindowDays.IsUnknown() {
        filters["windowDays"] = lookupNumber(data.WindowDays)
        filterNames = append(filterNames, "window_days = "+data.WindowDays.ValueBigFloat().String())
    }
    if !data.Timezone.IsNull() && !data.Timezone.IsUnknown() {
        filters["timezone"] = data.Timezone.ValueString()
        filterNames = append(filterNames, "timezone = "+fmt.Sprintf("%q", data.Timezone.ValueString()))
    }
    if !data.AtRiskThresholdPercentage.IsNull() && !data.AtRiskThresholdPercentage.IsUnknown() {
        filters["atRiskThresholdPercentage"] = lookupNumber(data.AtRiskThresholdPercentage)
        filterNames = append(filterNames, "at_risk_threshold_percentage = "+data.AtRiskThresholdPercentage.ValueBigFloat().String())
    }
    if !data.CurrentSliPercentage.IsNull() && !data.CurrentSliPercentage.IsUnknown() {
        filters["currentSliPercentage"] = lookupNumber(data.CurrentSliPercentage)
        filterNames = append(filterNames, "current_sli_percentage = "+data.CurrentSliPercentage.ValueBigFloat().String())
    }
    if !data.ErrorBudgetRemainingPercentage.IsNull() && !data.ErrorBudgetRemainingPercentage.IsUnknown() {
        filters["errorBudgetRemainingPercentage"] = lookupNumber(data.ErrorBudgetRemainingPercentage)
        filterNames = append(filterNames, "error_budget_remaining_percentage = "+data.ErrorBudgetRemainingPercentage.ValueBigFloat().String())
    }
    if !data.ErrorBudgetRemainingSeconds.IsNull() && !data.ErrorBudgetRemainingSeconds.IsUnknown() {
        filters["errorBudgetRemainingSeconds"] = lookupNumber(data.ErrorBudgetRemainingSeconds)
        filterNames = append(filterNames, "error_budget_remaining_seconds = "+data.ErrorBudgetRemainingSeconds.ValueBigFloat().String())
    }
    if !data.ErrorBudgetTotalSeconds.IsNull() && !data.ErrorBudgetTotalSeconds.IsUnknown() {
        filters["errorBudgetTotalSeconds"] = lookupNumber(data.ErrorBudgetTotalSeconds)
        filterNames = append(filterNames, "error_budget_total_seconds = "+data.ErrorBudgetTotalSeconds.ValueBigFloat().String())
    }
    if !data.CurrentBurnRate.IsNull() && !data.CurrentBurnRate.IsUnknown() {
        filters["currentBurnRate"] = lookupNumber(data.CurrentBurnRate)
        filterNames = append(filterNames, "current_burn_rate = "+data.CurrentBurnRate.ValueBigFloat().String())
    }
    if !data.SloStatus.IsNull() && !data.SloStatus.IsUnknown() {
        filters["sloStatus"] = data.SloStatus.ValueString()
        filterNames = append(filterNames, "slo_status = "+fmt.Sprintf("%q", data.SloStatus.ValueString()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the service level objective up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the service level objective up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "name": true,
        "description": true,
        "slug": true,
        "labels": true,
        "isEnabled": true,
        "isArchived": true,
        "archivedAt": true,
        "archivedByUserId": true,
        "sliType": true,
        "multiMonitorMode": true,
        "monitors": true,
        "monitorLabels": true,
        "autoAddedMonitors": true,
        "downtimeMonitorStatuses": true,
        "metricQueryConfig": true,
        "targetPercentage": true,
        "windowType": true,
        "windowDays": true,
        "timezone": true,
        "atRiskThresholdPercentage": true,
        "currentSliPercentage": true,
        "errorBudgetRemainingPercentage": true,
        "errorBudgetRemainingSeconds": true,
        "errorBudgetTotalSeconds": true,
        "currentBurnRate": true,
        "sloStatus": true,
        "statusChangeNotificationSentAt": true,
        "lastEvaluatedAt": true,
        "nextEvaluationAt": true,
        "lastAccumulatedBucketEndAt": true,
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/service-level-objective/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read service_level_objective, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No service level objective found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read service_level_objective: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/service-level-objective/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list service_level_objective, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list service_level_objective: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No service level objective matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one service level objective matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for service_level_objective.")
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
    if val, ok := item["labels"].([]interface{}); ok {
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if jsonBytes, err := json.Marshal(itemMap); err == nil {
                    setItems = append(setItems, types.StringValue(string(jsonBytes)))
                }
            } else if str, ok := item.(string); ok {
                setItems = append(setItems, types.StringValue(str))
            } else {
                setItems = append(setItems, types.StringValue(fmt.Sprintf("%v", item)))
            }
        }
        sort.Slice(setItems, func(i, j int) bool {
            return setItems[i].(types.String).ValueString() < setItems[j].(types.String).ValueString()
        })
        data.Labels = types.SetValueMust(types.StringType, setItems)
    } else {
        data.Labels = types.SetNull(types.StringType)
    }
    if val, ok := item["isEnabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    } else {
        data.IsEnabled = types.BoolNull()
    }
    if val, ok := item["isArchived"].(bool); ok {
        data.IsArchived = types.BoolValue(val)
    } else {
        data.IsArchived = types.BoolNull()
    }
    if obj, ok := item["archivedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ArchivedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ArchivedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ArchivedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ArchivedAt = types.StringValue(string(jsonBytes))
        } else {
            data.ArchivedAt = types.StringNull()
        }
    } else if val, ok := item["archivedAt"].(string); ok {
        data.ArchivedAt = types.StringValue(val)
    } else {
        data.ArchivedAt = types.StringNull()
    }
    if obj, ok := item["archivedByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ArchivedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ArchivedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ArchivedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ArchivedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.ArchivedByUserId = types.StringNull()
        }
    } else if val, ok := item["archivedByUserId"].(string); ok {
        data.ArchivedByUserId = types.StringValue(val)
    } else {
        data.ArchivedByUserId = types.StringNull()
    }
    if obj, ok := item["sliType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SliType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SliType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SliType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SliType = types.StringValue(string(jsonBytes))
        } else {
            data.SliType = types.StringNull()
        }
    } else if val, ok := item["sliType"].(string); ok {
        data.SliType = types.StringValue(val)
    } else {
        data.SliType = types.StringNull()
    }
    if obj, ok := item["multiMonitorMode"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MultiMonitorMode = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.MultiMonitorMode = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.MultiMonitorMode = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.MultiMonitorMode = types.StringValue(string(jsonBytes))
        } else {
            data.MultiMonitorMode = types.StringNull()
        }
    } else if val, ok := item["multiMonitorMode"].(string); ok {
        data.MultiMonitorMode = types.StringValue(val)
    } else {
        data.MultiMonitorMode = types.StringNull()
    }
    if val, ok := item["monitors"].([]interface{}); ok {
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if jsonBytes, err := json.Marshal(itemMap); err == nil {
                    setItems = append(setItems, types.StringValue(string(jsonBytes)))
                }
            } else if str, ok := item.(string); ok {
                setItems = append(setItems, types.StringValue(str))
            } else {
                setItems = append(setItems, types.StringValue(fmt.Sprintf("%v", item)))
            }
        }
        sort.Slice(setItems, func(i, j int) bool {
            return setItems[i].(types.String).ValueString() < setItems[j].(types.String).ValueString()
        })
        data.Monitors = types.SetValueMust(types.StringType, setItems)
    } else {
        data.Monitors = types.SetNull(types.StringType)
    }
    if val, ok := item["monitorLabels"].([]interface{}); ok {
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if jsonBytes, err := json.Marshal(itemMap); err == nil {
                    setItems = append(setItems, types.StringValue(string(jsonBytes)))
                }
            } else if str, ok := item.(string); ok {
                setItems = append(setItems, types.StringValue(str))
            } else {
                setItems = append(setItems, types.StringValue(fmt.Sprintf("%v", item)))
            }
        }
        sort.Slice(setItems, func(i, j int) bool {
            return setItems[i].(types.String).ValueString() < setItems[j].(types.String).ValueString()
        })
        data.MonitorLabels = types.SetValueMust(types.StringType, setItems)
    } else {
        data.MonitorLabels = types.SetNull(types.StringType)
    }
    if val, ok := item["autoAddedMonitors"].([]interface{}); ok {
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if jsonBytes, err := json.Marshal(itemMap); err == nil {
                    setItems = append(setItems, types.StringValue(string(jsonBytes)))
                }
            } else if str, ok := item.(string); ok {
                setItems = append(setItems, types.StringValue(str))
            } else {
                setItems = append(setItems, types.StringValue(fmt.Sprintf("%v", item)))
            }
        }
        sort.Slice(setItems, func(i, j int) bool {
            return setItems[i].(types.String).ValueString() < setItems[j].(types.String).ValueString()
        })
        data.AutoAddedMonitors = types.SetValueMust(types.StringType, setItems)
    } else {
        data.AutoAddedMonitors = types.SetNull(types.StringType)
    }
    if val, ok := item["downtimeMonitorStatuses"].([]interface{}); ok {
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if jsonBytes, err := json.Marshal(itemMap); err == nil {
                    setItems = append(setItems, types.StringValue(string(jsonBytes)))
                }
            } else if str, ok := item.(string); ok {
                setItems = append(setItems, types.StringValue(str))
            } else {
                setItems = append(setItems, types.StringValue(fmt.Sprintf("%v", item)))
            }
        }
        sort.Slice(setItems, func(i, j int) bool {
            return setItems[i].(types.String).ValueString() < setItems[j].(types.String).ValueString()
        })
        data.DowntimeMonitorStatuses = types.SetValueMust(types.StringType, setItems)
    } else {
        data.DowntimeMonitorStatuses = types.SetNull(types.StringType)
    }
    if obj, ok := item["metricQueryConfig"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MetricQueryConfig = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.MetricQueryConfig = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.MetricQueryConfig = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.MetricQueryConfig = types.StringValue(string(jsonBytes))
        } else {
            data.MetricQueryConfig = types.StringNull()
        }
    } else if val, ok := item["metricQueryConfig"].(string); ok {
        data.MetricQueryConfig = types.StringValue(val)
    } else {
        data.MetricQueryConfig = types.StringNull()
    }
    if val, ok := item["targetPercentage"].(float64); ok {
        data.TargetPercentage = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["targetPercentage"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.TargetPercentage = types.NumberValue(big.NewFloat(val))
        } else {
            data.TargetPercentage = types.NumberNull()
        }
    } else {
        data.TargetPercentage = types.NumberNull()
    }
    if obj, ok := item["windowType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.WindowType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.WindowType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.WindowType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.WindowType = types.StringValue(string(jsonBytes))
        } else {
            data.WindowType = types.StringNull()
        }
    } else if val, ok := item["windowType"].(string); ok {
        data.WindowType = types.StringValue(val)
    } else {
        data.WindowType = types.StringNull()
    }
    if val, ok := item["windowDays"].(float64); ok {
        data.WindowDays = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["windowDays"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.WindowDays = types.NumberValue(big.NewFloat(val))
        } else {
            data.WindowDays = types.NumberNull()
        }
    } else {
        data.WindowDays = types.NumberNull()
    }
    if obj, ok := item["timezone"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Timezone = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Timezone = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Timezone = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Timezone = types.StringValue(string(jsonBytes))
        } else {
            data.Timezone = types.StringNull()
        }
    } else if val, ok := item["timezone"].(string); ok {
        data.Timezone = types.StringValue(val)
    } else {
        data.Timezone = types.StringNull()
    }
    if val, ok := item["atRiskThresholdPercentage"].(float64); ok {
        data.AtRiskThresholdPercentage = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["atRiskThresholdPercentage"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.AtRiskThresholdPercentage = types.NumberValue(big.NewFloat(val))
        } else {
            data.AtRiskThresholdPercentage = types.NumberNull()
        }
    } else {
        data.AtRiskThresholdPercentage = types.NumberNull()
    }
    if val, ok := item["currentSliPercentage"].(float64); ok {
        data.CurrentSliPercentage = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["currentSliPercentage"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.CurrentSliPercentage = types.NumberValue(big.NewFloat(val))
        } else {
            data.CurrentSliPercentage = types.NumberNull()
        }
    } else {
        data.CurrentSliPercentage = types.NumberNull()
    }
    if val, ok := item["errorBudgetRemainingPercentage"].(float64); ok {
        data.ErrorBudgetRemainingPercentage = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["errorBudgetRemainingPercentage"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.ErrorBudgetRemainingPercentage = types.NumberValue(big.NewFloat(val))
        } else {
            data.ErrorBudgetRemainingPercentage = types.NumberNull()
        }
    } else {
        data.ErrorBudgetRemainingPercentage = types.NumberNull()
    }
    if val, ok := item["errorBudgetRemainingSeconds"].(float64); ok {
        data.ErrorBudgetRemainingSeconds = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["errorBudgetRemainingSeconds"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.ErrorBudgetRemainingSeconds = types.NumberValue(big.NewFloat(val))
        } else {
            data.ErrorBudgetRemainingSeconds = types.NumberNull()
        }
    } else {
        data.ErrorBudgetRemainingSeconds = types.NumberNull()
    }
    if val, ok := item["errorBudgetTotalSeconds"].(float64); ok {
        data.ErrorBudgetTotalSeconds = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["errorBudgetTotalSeconds"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.ErrorBudgetTotalSeconds = types.NumberValue(big.NewFloat(val))
        } else {
            data.ErrorBudgetTotalSeconds = types.NumberNull()
        }
    } else {
        data.ErrorBudgetTotalSeconds = types.NumberNull()
    }
    if val, ok := item["currentBurnRate"].(float64); ok {
        data.CurrentBurnRate = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["currentBurnRate"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.CurrentBurnRate = types.NumberValue(big.NewFloat(val))
        } else {
            data.CurrentBurnRate = types.NumberNull()
        }
    } else {
        data.CurrentBurnRate = types.NumberNull()
    }
    if obj, ok := item["sloStatus"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SloStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SloStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SloStatus = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SloStatus = types.StringValue(string(jsonBytes))
        } else {
            data.SloStatus = types.StringNull()
        }
    } else if val, ok := item["sloStatus"].(string); ok {
        data.SloStatus = types.StringValue(val)
    } else {
        data.SloStatus = types.StringNull()
    }
    if obj, ok := item["statusChangeNotificationSentAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StatusChangeNotificationSentAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.StatusChangeNotificationSentAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.StatusChangeNotificationSentAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.StatusChangeNotificationSentAt = types.StringValue(string(jsonBytes))
        } else {
            data.StatusChangeNotificationSentAt = types.StringNull()
        }
    } else if val, ok := item["statusChangeNotificationSentAt"].(string); ok {
        data.StatusChangeNotificationSentAt = types.StringValue(val)
    } else {
        data.StatusChangeNotificationSentAt = types.StringNull()
    }
    if obj, ok := item["lastEvaluatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastEvaluatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastEvaluatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastEvaluatedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastEvaluatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastEvaluatedAt = types.StringNull()
        }
    } else if val, ok := item["lastEvaluatedAt"].(string); ok {
        data.LastEvaluatedAt = types.StringValue(val)
    } else {
        data.LastEvaluatedAt = types.StringNull()
    }
    if obj, ok := item["nextEvaluationAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.NextEvaluationAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.NextEvaluationAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.NextEvaluationAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.NextEvaluationAt = types.StringValue(string(jsonBytes))
        } else {
            data.NextEvaluationAt = types.StringNull()
        }
    } else if val, ok := item["nextEvaluationAt"].(string); ok {
        data.NextEvaluationAt = types.StringValue(val)
    } else {
        data.NextEvaluationAt = types.StringNull()
    }
    if obj, ok := item["lastAccumulatedBucketEndAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastAccumulatedBucketEndAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastAccumulatedBucketEndAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastAccumulatedBucketEndAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastAccumulatedBucketEndAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastAccumulatedBucketEndAt = types.StringNull()
        }
    } else if val, ok := item["lastAccumulatedBucketEndAt"].(string); ok {
        data.LastAccumulatedBucketEndAt = types.StringValue(val)
    } else {
        data.LastAccumulatedBucketEndAt = types.StringNull()
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
