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
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    Description types.String `tfsdk:"description"`
    Slug types.String `tfsdk:"slug"`
    Labels types.Set `tfsdk:"labels"`
    IsEnabled types.Bool `tfsdk:"is_enabled"`
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
        MarkdownDescription: "Define Service Level Objectives (SLOs) with targets, compliance windows and error budgets, and track how much error budget remains. Look up an existing service_level_objective by `id` or by `name`.",

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
            "description": schema.StringAttribute{
                MarkdownDescription: "Description of this Service Level Objective.",
                Computed: true,
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Friendly globally unique name for your object.",
                Computed: true,
            },
            "labels": schema.SetAttribute{
                MarkdownDescription: "Relation to Labels Array where this object is categorized in..",
                Computed: true,
                ElementType: types.StringType,
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether this Service Level Objective is enabled. Disabled SLOs are not evaluated..",
                Computed: true,
            },
            "sli_type": schema.StringAttribute{
                MarkdownDescription: "Type of Service Level Indicator this objective measures (Monitor Uptime or Metric).",
                Computed: true,
            },
            "multi_monitor_mode": schema.StringAttribute{
                MarkdownDescription: "How downtime is counted when multiple monitors are attached. 'Any Monitor Down' counts time when any monitor is down. 'Monitor Seconds Average' averages downtime across monitors..",
                Computed: true,
            },
            "monitors": schema.SetAttribute{
                MarkdownDescription: "Monitors whose uptime is measured by this Service Level Objective (for Monitor Uptime SLIs)..",
                Computed: true,
                ElementType: types.StringType,
            },
            "monitor_labels": schema.SetAttribute{
                MarkdownDescription: "Monitor labels that automatically attach monitors to this SLO. Any monitor in the project carrying at least one of these labels is added to the Monitors list, and is removed again when it stops carrying any of them..",
                Computed: true,
                ElementType: types.StringType,
            },
            "auto_added_monitors": schema.SetAttribute{
                MarkdownDescription: "Monitors that were attached to this SLO by its label rule rather than by hand. Maintained by the server..",
                Computed: true,
                ElementType: types.StringType,
            },
            "downtime_monitor_statuses": schema.SetAttribute{
                MarkdownDescription: "List of monitor statuses that are considered as \"down\" for this Service Level Objective..",
                Computed: true,
                ElementType: types.StringType,
            },
            "metric_query_config": schema.StringAttribute{
                MarkdownDescription: "Query configuration for Metric SLIs: metric name, good-event predicate and optional attribute filters..",
                Computed: true,
            },
            "target_percentage": schema.NumberAttribute{
                MarkdownDescription: "Target of this Service Level Objective as a percentage (e.g. 99.9). Must be less than 100..",
                Computed: true,
            },
            "window_type": schema.StringAttribute{
                MarkdownDescription: "Type of compliance window for this objective (Rolling or Calendar Month).",
                Computed: true,
            },
            "window_days": schema.NumberAttribute{
                MarkdownDescription: "Length of the rolling compliance window in days (e.g. 7, 28, 30 or 90). Ignored for Calendar Month windows..",
                Computed: true,
            },
            "timezone": schema.StringAttribute{
                MarkdownDescription: "IANA timezone (e.g. America/New_York) used for Calendar Month window boundaries. Defaults to UTC when not set..",
                Computed: true,
            },
            "at_risk_threshold_percentage": schema.NumberAttribute{
                MarkdownDescription: "Percentage of remaining error budget at which the SLO status changes to At Risk. For example, 20 means the status becomes At Risk when less than 20% of the error budget remains..",
                Computed: true,
            },
            "current_sli_percentage": schema.NumberAttribute{
                MarkdownDescription: "Current Service Level Indicator over the compliance window, as a percentage. Computed by the worker..",
                Computed: true,
            },
            "error_budget_remaining_percentage": schema.NumberAttribute{
                MarkdownDescription: "Percentage of the error budget that remains. Can be negative when the budget is exhausted. Computed by the worker..",
                Computed: true,
            },
            "error_budget_remaining_seconds": schema.NumberAttribute{
                MarkdownDescription: "Seconds of error budget that remain. Can be negative when the budget is exhausted. Computed by the worker..",
                Computed: true,
            },
            "error_budget_total_seconds": schema.NumberAttribute{
                MarkdownDescription: "Total seconds of error budget for the compliance window. Computed by the worker..",
                Computed: true,
            },
            "current_burn_rate": schema.NumberAttribute{
                MarkdownDescription: "Rate at which the error budget is currently being consumed. A burn rate of 1 exhausts the budget exactly at the end of the window. Computed by the worker..",
                Computed: true,
            },
            "slo_status": schema.StringAttribute{
                MarkdownDescription: "Current status of this Service Level Objective (Healthy, At Risk, Budget Exhausted, Misconfigured, Paused). Computed by the worker..",
                Computed: true,
            },
            "status_change_notification_sent_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "last_evaluated_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "next_evaluation_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "last_accumulated_bucket_end_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
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

    hasId := !data.Id.IsNull() && data.Id.ValueString() != ""
    hasName := !data.Name.IsNull() && data.Name.ValueString() != ""
    if hasId == hasName {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Exactly one of `id` or `name` must be set to look up a service_level_objective.",
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
        "description": true,
        "slug": true,
        "labels": true,
        "isEnabled": true,
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
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No service_level_objective found with id %q.", data.Id.ValueString()))
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
    } else {
        listBody := map[string]interface{}{
            "query": map[string]interface{}{
                "name": data.Name.ValueString(),
            },
            "select": selectParam,
            // limit 2 is enough to detect ambiguity without paging.
            "limit": 2,
        }
        httpResp, err := d.client.Post(ctx, "/service-level-objective/get-list", listBody)
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
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No service_level_objective found with name %q.", data.Name.ValueString()))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one service_level_objective matches name %q. Use the id attribute to disambiguate.", data.Name.ValueString()))
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
