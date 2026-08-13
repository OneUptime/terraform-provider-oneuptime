package provider

import (
    "context"
    "fmt"
    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-framework/types/basetypes"
    "github.com/hashicorp/terraform-plugin-log/tflog"
    "math/big"
    "net/http"
    "github.com/hashicorp/terraform-plugin-framework/path"
    "encoding/json"
    "net/url"
    "strings"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/numberdefault"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
    "github.com/hashicorp/terraform-plugin-framework/attr"
    "sort"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/numberplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ServiceLevelObjectiveResource{}
var _ resource.ResourceWithImportState = &ServiceLevelObjectiveResource{}

func NewServiceLevelObjectiveResource() resource.Resource {
    return &ServiceLevelObjectiveResource{}
}

// ServiceLevelObjectiveResource defines the resource implementation.
type ServiceLevelObjectiveResource struct {
    client *Client
}

// ServiceLevelObjectiveResourceModel describes the resource data model.
type ServiceLevelObjectiveResourceModel struct {
    Id types.String `tfsdk:"id"`
    ProjectId types.String `tfsdk:"project_id"`
    Name types.String `tfsdk:"name"`
    Description types.String `tfsdk:"description"`
    Labels types.Set `tfsdk:"labels"`
    IsEnabled types.Bool `tfsdk:"is_enabled"`
    SliType types.String `tfsdk:"sli_type"`
    MultiMonitorMode types.String `tfsdk:"multi_monitor_mode"`
    Monitors types.Set `tfsdk:"monitors"`
    MonitorLabels types.Set `tfsdk:"monitor_labels"`
    MetricQueryConfig JSONSubsetValue `tfsdk:"metric_query_config"`
    TargetPercentage types.Number `tfsdk:"target_percentage"`
    WindowType types.String `tfsdk:"window_type"`
    WindowDays types.Number `tfsdk:"window_days"`
    Timezone types.String `tfsdk:"timezone"`
    AtRiskThresholdPercentage types.Number `tfsdk:"at_risk_threshold_percentage"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    CreatedAt RFC3339Value `tfsdk:"created_at"`
    UpdatedAt RFC3339Value `tfsdk:"updated_at"`
    DeletedAt RFC3339Value `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    Slug types.String `tfsdk:"slug"`
    AutoAddedMonitors types.Set `tfsdk:"auto_added_monitors"`
    DowntimeMonitorStatuses types.Set `tfsdk:"downtime_monitor_statuses"`
    CurrentSliPercentage types.Number `tfsdk:"current_sli_percentage"`
    ErrorBudgetRemainingPercentage types.Number `tfsdk:"error_budget_remaining_percentage"`
    ErrorBudgetRemainingSeconds types.Number `tfsdk:"error_budget_remaining_seconds"`
    ErrorBudgetTotalSeconds types.Number `tfsdk:"error_budget_total_seconds"`
    CurrentBurnRate types.Number `tfsdk:"current_burn_rate"`
    SloStatus types.String `tfsdk:"slo_status"`
    StatusChangeNotificationSentAt RFC3339Value `tfsdk:"status_change_notification_sent_at"`
    LastEvaluatedAt RFC3339Value `tfsdk:"last_evaluated_at"`
    NextEvaluationAt RFC3339Value `tfsdk:"next_evaluation_at"`
    LastAccumulatedBucketEndAt RFC3339Value `tfsdk:"last_accumulated_bucket_end_at"`
}

func (r *ServiceLevelObjectiveResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_service_level_objective"
}

func (r *ServiceLevelObjectiveResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Define Service Level Objectives (SLOs) with targets, compliance windows and error budgets, and track how much error budget remains.",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Unique identifier for the resource",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Name of this Service Level Objective.",
                Required: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Description of this Service Level Objective.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "labels": schema.SetAttribute{
                MarkdownDescription: "Relation to Labels Array where this object is categorized in..",
                Optional: true,
                Computed: true,
                ElementType: types.StringType,
                PlanModifiers: []planmodifier.Set{
                    setplanmodifier.UseStateForUnknown(),
                },
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether this Service Level Objective is enabled. Disabled SLOs are not evaluated..",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(true),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "sli_type": schema.StringAttribute{
                MarkdownDescription: "Type of Service Level Indicator this objective measures (Monitor Uptime or Metric).",
                Optional: true,
                Computed: true,
                Default: stringdefault.StaticString("Monitor Uptime"),
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "multi_monitor_mode": schema.StringAttribute{
                MarkdownDescription: "How downtime is counted when multiple monitors are attached. 'Any Monitor Down' counts time when any monitor is down. 'Monitor Seconds Average' averages downtime across monitors..",
                Optional: true,
                Computed: true,
                Default: stringdefault.StaticString("Any Monitor Down"),
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "monitors": schema.SetAttribute{
                MarkdownDescription: "Monitors whose uptime is measured by this Service Level Objective (for Monitor Uptime SLIs)..",
                Optional: true,
                Computed: true,
                ElementType: types.StringType,
                PlanModifiers: []planmodifier.Set{
                    setplanmodifier.UseStateForUnknown(),
                },
            },
            "monitor_labels": schema.SetAttribute{
                MarkdownDescription: "Monitor labels that automatically attach monitors to this SLO. Any monitor in the project carrying at least one of these labels is added to the Monitors list, and is removed again when it stops carrying any of them..",
                Optional: true,
                Computed: true,
                ElementType: types.StringType,
                PlanModifiers: []planmodifier.Set{
                    setplanmodifier.UseStateForUnknown(),
                },
            },
            "metric_query_config": schema.StringAttribute{
                MarkdownDescription: "Query configuration for Metric SLIs: metric name, good-event predicate and optional attribute filters..",
                CustomType: JSONSubsetType{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: []validator.String{
                    JSONEnvelopeValidator(),
                },
            },
            "target_percentage": schema.NumberAttribute{
                MarkdownDescription: "Target of this Service Level Objective as a percentage (e.g. 99.9). Must be less than 100..",
                Required: true,
            },
            "window_type": schema.StringAttribute{
                MarkdownDescription: "Type of compliance window for this objective (Rolling or Calendar Month).",
                Optional: true,
                Computed: true,
                Default: stringdefault.StaticString("Rolling"),
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "window_days": schema.NumberAttribute{
                MarkdownDescription: "Length of the rolling compliance window in days (e.g. 7, 28, 30 or 90). Ignored for Calendar Month windows..",
                Optional: true,
                Computed: true,
                Default: numberdefault.StaticBigFloat(big.NewFloat(30)),
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "timezone": schema.StringAttribute{
                MarkdownDescription: "IANA timezone (e.g. America/New_York) used for Calendar Month window boundaries. Defaults to UTC when not set..",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "at_risk_threshold_percentage": schema.NumberAttribute{
                MarkdownDescription: "Percentage of remaining error budget at which the SLO status changes to At Risk. For example, 20 means the status becomes At Risk when less than 20% of the error budget remains..",
                Optional: true,
                Computed: true,
                Default: numberdefault.StaticBigFloat(big.NewFloat(20)),
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                    stringplanmodifier.RequiresReplace(),
                },
            },
            "created_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "updated_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "deleted_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "version": schema.NumberAttribute{
                MarkdownDescription: "Object version",
                Computed: true,
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Friendly globally unique name for your object.",
                Computed: true,
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
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "last_evaluated_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "next_evaluation_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "last_accumulated_bucket_end_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
        },
    }
}

func (r *ServiceLevelObjectiveResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
    // Prevent panic if the provider has not been configured.
    if req.ProviderData == nil {
        return
    }

    client, ok := req.ProviderData.(*Client)

    if !ok {
        resp.Diagnostics.AddError(
            "Unexpected Resource Configure Type",
            fmt.Sprintf("Expected *Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
        )

        return
    }

    r.client = client
}


func (r *ServiceLevelObjectiveResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var data ServiceLevelObjectiveResourceModel

    // Read Terraform plan data into the model
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }



    // Create API request body. Unset (null/unknown) optional fields are
    // omitted so server-side defaults apply instead of being overwritten
    // with zero values.
    serviceLevelObjectiveRequest := map[string]interface{}{
        "data": map[string]interface{}{},
    }
    requestDataMap := serviceLevelObjectiveRequest["data"].(map[string]interface{})

    if !data.Name.IsNull() && !data.Name.IsUnknown() {
        requestDataMap["name"] = data.Name.ValueString()
    }
    if !data.Description.IsNull() && !data.Description.IsUnknown() {
        requestDataMap["description"] = data.Description.ValueString()
    }
    if !data.Labels.IsNull() && !data.Labels.IsUnknown() {
        requestDataMap["labels"] = r.convertTerraformSetToInterface(data.Labels)
    }
    if !data.IsEnabled.IsNull() && !data.IsEnabled.IsUnknown() {
        requestDataMap["isEnabled"] = data.IsEnabled.ValueBool()
    }
    if !data.SliType.IsNull() && !data.SliType.IsUnknown() {
        requestDataMap["sliType"] = data.SliType.ValueString()
    }
    if !data.MultiMonitorMode.IsNull() && !data.MultiMonitorMode.IsUnknown() {
        requestDataMap["multiMonitorMode"] = data.MultiMonitorMode.ValueString()
    }
    if !data.Monitors.IsNull() && !data.Monitors.IsUnknown() {
        requestDataMap["monitors"] = r.convertTerraformSetToInterface(data.Monitors)
    }
    if !data.MonitorLabels.IsNull() && !data.MonitorLabels.IsUnknown() {
        requestDataMap["monitorLabels"] = r.convertTerraformSetToInterface(data.MonitorLabels)
    }
    if parsedMetricQueryConfig := r.parseJSONField(data.MetricQueryConfig); parsedMetricQueryConfig != nil {
        requestDataMap["metricQueryConfig"] = parsedMetricQueryConfig
    }
    if !data.TargetPercentage.IsNull() && !data.TargetPercentage.IsUnknown() {
        requestDataMap["targetPercentage"] = r.bigFloatToFloat64(data.TargetPercentage.ValueBigFloat())
    }
    if !data.WindowType.IsNull() && !data.WindowType.IsUnknown() {
        requestDataMap["windowType"] = data.WindowType.ValueString()
    }
    if !data.WindowDays.IsNull() && !data.WindowDays.IsUnknown() {
        requestDataMap["windowDays"] = r.bigFloatToFloat64(data.WindowDays.ValueBigFloat())
    }
    if !data.Timezone.IsNull() && !data.Timezone.IsUnknown() {
        requestDataMap["timezone"] = data.Timezone.ValueString()
    }
    if !data.AtRiskThresholdPercentage.IsNull() && !data.AtRiskThresholdPercentage.IsUnknown() {
        requestDataMap["atRiskThresholdPercentage"] = r.bigFloatToFloat64(data.AtRiskThresholdPercentage.ValueBigFloat())
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        requestDataMap["createdByUserId"] = data.CreatedByUserId.ValueString()
    }

    // Make API call
    httpResp, err := r.client.Post(ctx, "/service-level-objective", serviceLevelObjectiveRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create service_level_objective, got error: %s", err))
        return
    }

    var serviceLevelObjectiveResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &serviceLevelObjectiveResponse)
    if err != nil {
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to create service_level_objective: %s", err))
        return
    }

    // Extract the new resource id from the create response.
    createdId := ""
    if wrapper, ok := serviceLevelObjectiveResponse["data"].(map[string]interface{}); ok {
        if val, ok := wrapper["_id"].(string); ok {
            createdId = val
        }
    } else if val, ok := serviceLevelObjectiveResponse["_id"].(string); ok {
        createdId = val
    }
    if createdId == "" {
        resp.Diagnostics.AddError("OneUptime API Error", "Create response for service_level_objective did not contain an id. This is a bug in the provider or the API; please report it.")
        return
    }
    data.Id = types.StringValue(createdId)

    /*
     * The server has committed the row. Persist what we know to state BEFORE
     * the read-back: if the read-back fails and we return without setting
     * state, Terraform never learns the resource exists and the created
     * service_level_objective is orphaned server-side — never refreshed, never
     * destroyed. Delete already refuses to drop state on failure for the
     * same reason; Create must not either.
     */
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Re-read the resource so state reflects server-normalized values.
    selectParam := map[string]interface{}{
        "projectId": true,
        "name": true,
        "description": true,
        "labels": true,
        "isEnabled": true,
        "sliType": true,
        "multiMonitorMode": true,
        "monitors": true,
        "monitorLabels": true,
        "metricQueryConfig": true,
        "targetPercentage": true,
        "windowType": true,
        "windowDays": true,
        "timezone": true,
        "atRiskThresholdPercentage": true,
        "createdByUserId": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "slug": true,
        "autoAddedMonitors": true,
        "downtimeMonitorStatuses": true,
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
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect(ctx, "/service-level-objective/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        /*
         * State already owns the id, so the resource is tracked and the next
         * refresh reconciles the remaining attributes. Warn rather than
         * error: erroring here would strand a real resource.
         */
        resp.Diagnostics.AddWarning("Read After Create Failed", fmt.Sprintf("Created service_level_objective but could not read it back; state is incomplete until the next refresh: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddWarning("Read After Create Failed", fmt.Sprintf("Created service_level_objective but could not parse the read-back response; state is incomplete until the next refresh: %s", err))
        return
    }

    // Update the model with the authoritative read response
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := readResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = readResponse
    }

    if obj, ok := dataMap["projectId"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok {
            data.ProjectId = types.StringValue(val)
        } else {
            data.ProjectId = types.StringNull()
        }
    } else if val, ok := dataMap["projectId"].(string); ok {
        data.ProjectId = types.StringValue(val)
    } else {
        data.ProjectId = types.StringNull()
    }
    if obj, ok := dataMap["name"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Name = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Name = types.StringValue(string(jsonBytes))
            } else {
                data.Name = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Name = types.StringValue(string(jsonBytes))
            } else {
                data.Name = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Name = types.StringValue(string(jsonBytes))
        } else {
            data.Name = types.StringNull()
        }
    } else if val, ok := dataMap["name"].(string); ok {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if obj, ok := dataMap["description"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Description = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Description = types.StringValue(string(jsonBytes))
            } else {
                data.Description = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Description = types.StringValue(string(jsonBytes))
            } else {
                data.Description = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Description = types.StringValue(string(jsonBytes))
        } else {
            data.Description = types.StringNull()
        }
    } else if val, ok := dataMap["description"].(string); ok {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if val, ok := dataMap["labels"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.Labels = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.Labels = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if val, ok := dataMap["isEnabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    }
    if obj, ok := dataMap["sliType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SliType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SliType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SliType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SliType = types.StringValue(string(jsonBytes))
            } else {
                data.SliType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SliType = types.StringValue(string(jsonBytes))
            } else {
                data.SliType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SliType = types.StringValue(string(jsonBytes))
        } else {
            data.SliType = types.StringNull()
        }
    } else if val, ok := dataMap["sliType"].(string); ok {
        data.SliType = types.StringValue(val)
    } else {
        data.SliType = types.StringNull()
    }
    if obj, ok := dataMap["multiMonitorMode"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MultiMonitorMode = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.MultiMonitorMode = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.MultiMonitorMode = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.MultiMonitorMode = types.StringValue(string(jsonBytes))
            } else {
                data.MultiMonitorMode = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.MultiMonitorMode = types.StringValue(string(jsonBytes))
            } else {
                data.MultiMonitorMode = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.MultiMonitorMode = types.StringValue(string(jsonBytes))
        } else {
            data.MultiMonitorMode = types.StringNull()
        }
    } else if val, ok := dataMap["multiMonitorMode"].(string); ok {
        data.MultiMonitorMode = types.StringValue(val)
    } else {
        data.MultiMonitorMode = types.StringNull()
    }
    if val, ok := dataMap["monitors"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.Monitors = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.Monitors = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if val, ok := dataMap["monitorLabels"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.MonitorLabels = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.MonitorLabels = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if obj, ok := dataMap["metricQueryConfig"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MetricQueryConfig = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.MetricQueryConfig = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.MetricQueryConfig = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.MetricQueryConfig = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.MetricQueryConfig = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.MetricQueryConfig = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.MetricQueryConfig = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.MetricQueryConfig = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.MetricQueryConfig = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["metricQueryConfig"].(string); ok {
        data.MetricQueryConfig = NewJSONSubsetValue(val)
    } else {
        data.MetricQueryConfig = NewJSONSubsetNull()
    }
    if val, ok := dataMap["targetPercentage"].(float64); ok {
        data.TargetPercentage = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["targetPercentage"].(int); ok {
        data.TargetPercentage = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["targetPercentage"].(int64); ok {
        data.TargetPercentage = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["targetPercentage"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.TargetPercentage = types.NumberValue(big.NewFloat(val))
        } else {
            data.TargetPercentage = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.TargetPercentage = types.NumberNull()
    }
    if obj, ok := dataMap["windowType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.WindowType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.WindowType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.WindowType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.WindowType = types.StringValue(string(jsonBytes))
            } else {
                data.WindowType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.WindowType = types.StringValue(string(jsonBytes))
            } else {
                data.WindowType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.WindowType = types.StringValue(string(jsonBytes))
        } else {
            data.WindowType = types.StringNull()
        }
    } else if val, ok := dataMap["windowType"].(string); ok {
        data.WindowType = types.StringValue(val)
    } else {
        data.WindowType = types.StringNull()
    }
    if val, ok := dataMap["windowDays"].(float64); ok {
        data.WindowDays = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["windowDays"].(int); ok {
        data.WindowDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["windowDays"].(int64); ok {
        data.WindowDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["windowDays"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.WindowDays = types.NumberValue(big.NewFloat(val))
        } else {
            data.WindowDays = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.WindowDays = types.NumberNull()
    }
    if obj, ok := dataMap["timezone"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Timezone = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Timezone = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Timezone = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Timezone = types.StringValue(string(jsonBytes))
            } else {
                data.Timezone = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Timezone = types.StringValue(string(jsonBytes))
            } else {
                data.Timezone = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Timezone = types.StringValue(string(jsonBytes))
        } else {
            data.Timezone = types.StringNull()
        }
    } else if val, ok := dataMap["timezone"].(string); ok {
        data.Timezone = types.StringValue(val)
    } else {
        data.Timezone = types.StringNull()
    }
    if val, ok := dataMap["atRiskThresholdPercentage"].(float64); ok {
        data.AtRiskThresholdPercentage = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["atRiskThresholdPercentage"].(int); ok {
        data.AtRiskThresholdPercentage = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["atRiskThresholdPercentage"].(int64); ok {
        data.AtRiskThresholdPercentage = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["atRiskThresholdPercentage"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.AtRiskThresholdPercentage = types.NumberValue(big.NewFloat(val))
        } else {
            data.AtRiskThresholdPercentage = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.AtRiskThresholdPercentage = types.NumberNull()
    }
    if obj, ok := dataMap["createdByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["createdByUserId"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if obj, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.CreatedAt = NewRFC3339Value(val)
        } else {
            data.CreatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["createdAt"].(string); ok && val != "" {
        data.CreatedAt = NewRFC3339Value(val)
    } else {
        data.CreatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.UpdatedAt = NewRFC3339Value(val)
        } else {
            data.UpdatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["updatedAt"].(string); ok && val != "" {
        data.UpdatedAt = NewRFC3339Value(val)
    } else {
        data.UpdatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["deletedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.DeletedAt = NewRFC3339Value(val)
        } else {
            data.DeletedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["deletedAt"].(string); ok && val != "" {
        data.DeletedAt = NewRFC3339Value(val)
    } else {
        data.DeletedAt = NewRFC3339Null()
    }
    if val, ok := dataMap["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["version"].(int); ok {
        data.Version = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["version"].(int64); ok {
        data.Version = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["version"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.Version = types.NumberValue(big.NewFloat(val))
        } else {
            data.Version = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.Version = types.NumberNull()
    }
    if obj, ok := dataMap["slug"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Slug = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Slug = types.StringValue(string(jsonBytes))
            } else {
                data.Slug = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Slug = types.StringValue(string(jsonBytes))
            } else {
                data.Slug = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Slug = types.StringValue(string(jsonBytes))
        } else {
            data.Slug = types.StringNull()
        }
    } else if val, ok := dataMap["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
    }
    if val, ok := dataMap["autoAddedMonitors"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.AutoAddedMonitors = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.AutoAddedMonitors = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if val, ok := dataMap["downtimeMonitorStatuses"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.DowntimeMonitorStatuses = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.DowntimeMonitorStatuses = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if val, ok := dataMap["currentSliPercentage"].(float64); ok {
        data.CurrentSliPercentage = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["currentSliPercentage"].(int); ok {
        data.CurrentSliPercentage = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["currentSliPercentage"].(int64); ok {
        data.CurrentSliPercentage = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["currentSliPercentage"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.CurrentSliPercentage = types.NumberValue(big.NewFloat(val))
        } else {
            data.CurrentSliPercentage = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.CurrentSliPercentage = types.NumberNull()
    }
    if val, ok := dataMap["errorBudgetRemainingPercentage"].(float64); ok {
        data.ErrorBudgetRemainingPercentage = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["errorBudgetRemainingPercentage"].(int); ok {
        data.ErrorBudgetRemainingPercentage = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["errorBudgetRemainingPercentage"].(int64); ok {
        data.ErrorBudgetRemainingPercentage = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["errorBudgetRemainingPercentage"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.ErrorBudgetRemainingPercentage = types.NumberValue(big.NewFloat(val))
        } else {
            data.ErrorBudgetRemainingPercentage = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.ErrorBudgetRemainingPercentage = types.NumberNull()
    }
    if val, ok := dataMap["errorBudgetRemainingSeconds"].(float64); ok {
        data.ErrorBudgetRemainingSeconds = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["errorBudgetRemainingSeconds"].(int); ok {
        data.ErrorBudgetRemainingSeconds = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["errorBudgetRemainingSeconds"].(int64); ok {
        data.ErrorBudgetRemainingSeconds = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["errorBudgetRemainingSeconds"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.ErrorBudgetRemainingSeconds = types.NumberValue(big.NewFloat(val))
        } else {
            data.ErrorBudgetRemainingSeconds = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.ErrorBudgetRemainingSeconds = types.NumberNull()
    }
    if val, ok := dataMap["errorBudgetTotalSeconds"].(float64); ok {
        data.ErrorBudgetTotalSeconds = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["errorBudgetTotalSeconds"].(int); ok {
        data.ErrorBudgetTotalSeconds = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["errorBudgetTotalSeconds"].(int64); ok {
        data.ErrorBudgetTotalSeconds = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["errorBudgetTotalSeconds"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.ErrorBudgetTotalSeconds = types.NumberValue(big.NewFloat(val))
        } else {
            data.ErrorBudgetTotalSeconds = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.ErrorBudgetTotalSeconds = types.NumberNull()
    }
    if val, ok := dataMap["currentBurnRate"].(float64); ok {
        data.CurrentBurnRate = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["currentBurnRate"].(int); ok {
        data.CurrentBurnRate = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["currentBurnRate"].(int64); ok {
        data.CurrentBurnRate = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["currentBurnRate"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.CurrentBurnRate = types.NumberValue(big.NewFloat(val))
        } else {
            data.CurrentBurnRate = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.CurrentBurnRate = types.NumberNull()
    }
    if obj, ok := dataMap["sloStatus"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SloStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SloStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SloStatus = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SloStatus = types.StringValue(string(jsonBytes))
            } else {
                data.SloStatus = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SloStatus = types.StringValue(string(jsonBytes))
            } else {
                data.SloStatus = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SloStatus = types.StringValue(string(jsonBytes))
        } else {
            data.SloStatus = types.StringNull()
        }
    } else if val, ok := dataMap["sloStatus"].(string); ok {
        data.SloStatus = types.StringValue(val)
    } else {
        data.SloStatus = types.StringNull()
    }
    if obj, ok := dataMap["statusChangeNotificationSentAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.StatusChangeNotificationSentAt = NewRFC3339Value(val)
        } else {
            data.StatusChangeNotificationSentAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["statusChangeNotificationSentAt"].(string); ok && val != "" {
        data.StatusChangeNotificationSentAt = NewRFC3339Value(val)
    } else {
        data.StatusChangeNotificationSentAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["lastEvaluatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.LastEvaluatedAt = NewRFC3339Value(val)
        } else {
            data.LastEvaluatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["lastEvaluatedAt"].(string); ok && val != "" {
        data.LastEvaluatedAt = NewRFC3339Value(val)
    } else {
        data.LastEvaluatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["nextEvaluationAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.NextEvaluationAt = NewRFC3339Value(val)
        } else {
            data.NextEvaluationAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["nextEvaluationAt"].(string); ok && val != "" {
        data.NextEvaluationAt = NewRFC3339Value(val)
    } else {
        data.NextEvaluationAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["lastAccumulatedBucketEndAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.LastAccumulatedBucketEndAt = NewRFC3339Value(val)
        } else {
            data.LastAccumulatedBucketEndAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["lastAccumulatedBucketEndAt"].(string); ok && val != "" {
        data.LastAccumulatedBucketEndAt = NewRFC3339Value(val)
    } else {
        data.LastAccumulatedBucketEndAt = NewRFC3339Null()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    // The read response is authoritative, but never let it clobber the id we just received.
    data.Id = types.StringValue(createdId)

    // Write logs using the tflog package
    tflog.Trace(ctx, "created a resource")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ServiceLevelObjectiveResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var data ServiceLevelObjectiveResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create select parameter to get full object
    selectParam := map[string]interface{}{
        "projectId": true,
        "name": true,
        "description": true,
        "labels": true,
        "isEnabled": true,
        "sliType": true,
        "multiMonitorMode": true,
        "monitors": true,
        "monitorLabels": true,
        "metricQueryConfig": true,
        "targetPercentage": true,
        "windowType": true,
        "windowDays": true,
        "timezone": true,
        "atRiskThresholdPercentage": true,
        "createdByUserId": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "slug": true,
        "autoAddedMonitors": true,
        "downtimeMonitorStatuses": true,
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
        "_id": true,
    }

    // Make API call with select parameter
    httpResp, err := r.client.PostWithSelect(ctx, "/service-level-objective/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read service_level_objective, got error: %s", err))
        return
    }

    if httpResp.StatusCode == http.StatusNotFound {
        resp.State.RemoveResource(ctx)
        return
    }

    var serviceLevelObjectiveResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &serviceLevelObjectiveResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse service_level_objective response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := serviceLevelObjectiveResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = serviceLevelObjectiveResponse
    }

    if obj, ok := dataMap["projectId"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok {
            data.ProjectId = types.StringValue(val)
        } else {
            data.ProjectId = types.StringNull()
        }
    } else if val, ok := dataMap["projectId"].(string); ok {
        data.ProjectId = types.StringValue(val)
    } else {
        data.ProjectId = types.StringNull()
    }
    if obj, ok := dataMap["name"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Name = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Name = types.StringValue(string(jsonBytes))
            } else {
                data.Name = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Name = types.StringValue(string(jsonBytes))
            } else {
                data.Name = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Name = types.StringValue(string(jsonBytes))
        } else {
            data.Name = types.StringNull()
        }
    } else if val, ok := dataMap["name"].(string); ok {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if obj, ok := dataMap["description"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Description = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Description = types.StringValue(string(jsonBytes))
            } else {
                data.Description = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Description = types.StringValue(string(jsonBytes))
            } else {
                data.Description = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Description = types.StringValue(string(jsonBytes))
        } else {
            data.Description = types.StringNull()
        }
    } else if val, ok := dataMap["description"].(string); ok {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if val, ok := dataMap["labels"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.Labels = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.Labels = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if val, ok := dataMap["isEnabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    }
    if obj, ok := dataMap["sliType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SliType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SliType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SliType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SliType = types.StringValue(string(jsonBytes))
            } else {
                data.SliType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SliType = types.StringValue(string(jsonBytes))
            } else {
                data.SliType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SliType = types.StringValue(string(jsonBytes))
        } else {
            data.SliType = types.StringNull()
        }
    } else if val, ok := dataMap["sliType"].(string); ok {
        data.SliType = types.StringValue(val)
    } else {
        data.SliType = types.StringNull()
    }
    if obj, ok := dataMap["multiMonitorMode"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MultiMonitorMode = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.MultiMonitorMode = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.MultiMonitorMode = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.MultiMonitorMode = types.StringValue(string(jsonBytes))
            } else {
                data.MultiMonitorMode = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.MultiMonitorMode = types.StringValue(string(jsonBytes))
            } else {
                data.MultiMonitorMode = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.MultiMonitorMode = types.StringValue(string(jsonBytes))
        } else {
            data.MultiMonitorMode = types.StringNull()
        }
    } else if val, ok := dataMap["multiMonitorMode"].(string); ok {
        data.MultiMonitorMode = types.StringValue(val)
    } else {
        data.MultiMonitorMode = types.StringNull()
    }
    if val, ok := dataMap["monitors"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.Monitors = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.Monitors = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if val, ok := dataMap["monitorLabels"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.MonitorLabels = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.MonitorLabels = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if obj, ok := dataMap["metricQueryConfig"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MetricQueryConfig = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.MetricQueryConfig = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.MetricQueryConfig = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.MetricQueryConfig = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.MetricQueryConfig = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.MetricQueryConfig = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.MetricQueryConfig = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.MetricQueryConfig = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.MetricQueryConfig = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["metricQueryConfig"].(string); ok {
        data.MetricQueryConfig = NewJSONSubsetValue(val)
    } else {
        data.MetricQueryConfig = NewJSONSubsetNull()
    }
    if val, ok := dataMap["targetPercentage"].(float64); ok {
        data.TargetPercentage = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["targetPercentage"].(int); ok {
        data.TargetPercentage = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["targetPercentage"].(int64); ok {
        data.TargetPercentage = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["targetPercentage"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.TargetPercentage = types.NumberValue(big.NewFloat(val))
        } else {
            data.TargetPercentage = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.TargetPercentage = types.NumberNull()
    }
    if obj, ok := dataMap["windowType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.WindowType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.WindowType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.WindowType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.WindowType = types.StringValue(string(jsonBytes))
            } else {
                data.WindowType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.WindowType = types.StringValue(string(jsonBytes))
            } else {
                data.WindowType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.WindowType = types.StringValue(string(jsonBytes))
        } else {
            data.WindowType = types.StringNull()
        }
    } else if val, ok := dataMap["windowType"].(string); ok {
        data.WindowType = types.StringValue(val)
    } else {
        data.WindowType = types.StringNull()
    }
    if val, ok := dataMap["windowDays"].(float64); ok {
        data.WindowDays = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["windowDays"].(int); ok {
        data.WindowDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["windowDays"].(int64); ok {
        data.WindowDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["windowDays"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.WindowDays = types.NumberValue(big.NewFloat(val))
        } else {
            data.WindowDays = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.WindowDays = types.NumberNull()
    }
    if obj, ok := dataMap["timezone"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Timezone = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Timezone = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Timezone = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Timezone = types.StringValue(string(jsonBytes))
            } else {
                data.Timezone = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Timezone = types.StringValue(string(jsonBytes))
            } else {
                data.Timezone = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Timezone = types.StringValue(string(jsonBytes))
        } else {
            data.Timezone = types.StringNull()
        }
    } else if val, ok := dataMap["timezone"].(string); ok {
        data.Timezone = types.StringValue(val)
    } else {
        data.Timezone = types.StringNull()
    }
    if val, ok := dataMap["atRiskThresholdPercentage"].(float64); ok {
        data.AtRiskThresholdPercentage = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["atRiskThresholdPercentage"].(int); ok {
        data.AtRiskThresholdPercentage = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["atRiskThresholdPercentage"].(int64); ok {
        data.AtRiskThresholdPercentage = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["atRiskThresholdPercentage"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.AtRiskThresholdPercentage = types.NumberValue(big.NewFloat(val))
        } else {
            data.AtRiskThresholdPercentage = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.AtRiskThresholdPercentage = types.NumberNull()
    }
    if obj, ok := dataMap["createdByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["createdByUserId"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if obj, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.CreatedAt = NewRFC3339Value(val)
        } else {
            data.CreatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["createdAt"].(string); ok && val != "" {
        data.CreatedAt = NewRFC3339Value(val)
    } else {
        data.CreatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.UpdatedAt = NewRFC3339Value(val)
        } else {
            data.UpdatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["updatedAt"].(string); ok && val != "" {
        data.UpdatedAt = NewRFC3339Value(val)
    } else {
        data.UpdatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["deletedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.DeletedAt = NewRFC3339Value(val)
        } else {
            data.DeletedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["deletedAt"].(string); ok && val != "" {
        data.DeletedAt = NewRFC3339Value(val)
    } else {
        data.DeletedAt = NewRFC3339Null()
    }
    if val, ok := dataMap["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["version"].(int); ok {
        data.Version = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["version"].(int64); ok {
        data.Version = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["version"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.Version = types.NumberValue(big.NewFloat(val))
        } else {
            data.Version = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.Version = types.NumberNull()
    }
    if obj, ok := dataMap["slug"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Slug = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Slug = types.StringValue(string(jsonBytes))
            } else {
                data.Slug = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Slug = types.StringValue(string(jsonBytes))
            } else {
                data.Slug = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Slug = types.StringValue(string(jsonBytes))
        } else {
            data.Slug = types.StringNull()
        }
    } else if val, ok := dataMap["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
    }
    if val, ok := dataMap["autoAddedMonitors"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.AutoAddedMonitors = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.AutoAddedMonitors = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if val, ok := dataMap["downtimeMonitorStatuses"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.DowntimeMonitorStatuses = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.DowntimeMonitorStatuses = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if val, ok := dataMap["currentSliPercentage"].(float64); ok {
        data.CurrentSliPercentage = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["currentSliPercentage"].(int); ok {
        data.CurrentSliPercentage = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["currentSliPercentage"].(int64); ok {
        data.CurrentSliPercentage = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["currentSliPercentage"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.CurrentSliPercentage = types.NumberValue(big.NewFloat(val))
        } else {
            data.CurrentSliPercentage = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.CurrentSliPercentage = types.NumberNull()
    }
    if val, ok := dataMap["errorBudgetRemainingPercentage"].(float64); ok {
        data.ErrorBudgetRemainingPercentage = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["errorBudgetRemainingPercentage"].(int); ok {
        data.ErrorBudgetRemainingPercentage = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["errorBudgetRemainingPercentage"].(int64); ok {
        data.ErrorBudgetRemainingPercentage = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["errorBudgetRemainingPercentage"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.ErrorBudgetRemainingPercentage = types.NumberValue(big.NewFloat(val))
        } else {
            data.ErrorBudgetRemainingPercentage = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.ErrorBudgetRemainingPercentage = types.NumberNull()
    }
    if val, ok := dataMap["errorBudgetRemainingSeconds"].(float64); ok {
        data.ErrorBudgetRemainingSeconds = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["errorBudgetRemainingSeconds"].(int); ok {
        data.ErrorBudgetRemainingSeconds = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["errorBudgetRemainingSeconds"].(int64); ok {
        data.ErrorBudgetRemainingSeconds = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["errorBudgetRemainingSeconds"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.ErrorBudgetRemainingSeconds = types.NumberValue(big.NewFloat(val))
        } else {
            data.ErrorBudgetRemainingSeconds = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.ErrorBudgetRemainingSeconds = types.NumberNull()
    }
    if val, ok := dataMap["errorBudgetTotalSeconds"].(float64); ok {
        data.ErrorBudgetTotalSeconds = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["errorBudgetTotalSeconds"].(int); ok {
        data.ErrorBudgetTotalSeconds = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["errorBudgetTotalSeconds"].(int64); ok {
        data.ErrorBudgetTotalSeconds = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["errorBudgetTotalSeconds"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.ErrorBudgetTotalSeconds = types.NumberValue(big.NewFloat(val))
        } else {
            data.ErrorBudgetTotalSeconds = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.ErrorBudgetTotalSeconds = types.NumberNull()
    }
    if val, ok := dataMap["currentBurnRate"].(float64); ok {
        data.CurrentBurnRate = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["currentBurnRate"].(int); ok {
        data.CurrentBurnRate = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["currentBurnRate"].(int64); ok {
        data.CurrentBurnRate = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["currentBurnRate"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.CurrentBurnRate = types.NumberValue(big.NewFloat(val))
        } else {
            data.CurrentBurnRate = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.CurrentBurnRate = types.NumberNull()
    }
    if obj, ok := dataMap["sloStatus"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SloStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SloStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SloStatus = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SloStatus = types.StringValue(string(jsonBytes))
            } else {
                data.SloStatus = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SloStatus = types.StringValue(string(jsonBytes))
            } else {
                data.SloStatus = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SloStatus = types.StringValue(string(jsonBytes))
        } else {
            data.SloStatus = types.StringNull()
        }
    } else if val, ok := dataMap["sloStatus"].(string); ok {
        data.SloStatus = types.StringValue(val)
    } else {
        data.SloStatus = types.StringNull()
    }
    if obj, ok := dataMap["statusChangeNotificationSentAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.StatusChangeNotificationSentAt = NewRFC3339Value(val)
        } else {
            data.StatusChangeNotificationSentAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["statusChangeNotificationSentAt"].(string); ok && val != "" {
        data.StatusChangeNotificationSentAt = NewRFC3339Value(val)
    } else {
        data.StatusChangeNotificationSentAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["lastEvaluatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.LastEvaluatedAt = NewRFC3339Value(val)
        } else {
            data.LastEvaluatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["lastEvaluatedAt"].(string); ok && val != "" {
        data.LastEvaluatedAt = NewRFC3339Value(val)
    } else {
        data.LastEvaluatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["nextEvaluationAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.NextEvaluationAt = NewRFC3339Value(val)
        } else {
            data.NextEvaluationAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["nextEvaluationAt"].(string); ok && val != "" {
        data.NextEvaluationAt = NewRFC3339Value(val)
    } else {
        data.NextEvaluationAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["lastAccumulatedBucketEndAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.LastAccumulatedBucketEndAt = NewRFC3339Value(val)
        } else {
            data.LastAccumulatedBucketEndAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["lastAccumulatedBucketEndAt"].(string); ok && val != "" {
        data.LastAccumulatedBucketEndAt = NewRFC3339Value(val)
    } else {
        data.LastAccumulatedBucketEndAt = NewRFC3339Null()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ServiceLevelObjectiveResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var data ServiceLevelObjectiveResourceModel
    var state ServiceLevelObjectiveResourceModel

    // Read Terraform current state data to get the ID
    resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Read Terraform plan data to get the new values
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Use the ID from the current state
    data.Id = state.Id

    // Create API request body
    serviceLevelObjectiveRequest := map[string]interface{}{
        "data": map[string]interface{}{},
    }
    requestDataMap := serviceLevelObjectiveRequest["data"].(map[string]interface{})

    if !data.Name.IsUnknown() && !state.Name.IsUnknown() && !data.Name.Equal(state.Name) {
        requestDataMap["name"] = data.Name.ValueString()
    }
    if !data.Description.IsUnknown() && !state.Description.IsUnknown() && !data.Description.Equal(state.Description) {
        requestDataMap["description"] = data.Description.ValueString()
    }
    if !data.Labels.IsUnknown() && !state.Labels.IsUnknown() && !data.Labels.Equal(state.Labels) {
        requestDataMap["labels"] = r.convertTerraformSetToInterface(data.Labels)
    }
    if !data.IsEnabled.IsUnknown() && !state.IsEnabled.IsUnknown() && !data.IsEnabled.Equal(state.IsEnabled) {
        requestDataMap["isEnabled"] = data.IsEnabled.ValueBool()
    }
    if !data.SliType.IsUnknown() && !state.SliType.IsUnknown() && !data.SliType.Equal(state.SliType) {
        requestDataMap["sliType"] = data.SliType.ValueString()
    }
    if !data.MultiMonitorMode.IsUnknown() && !state.MultiMonitorMode.IsUnknown() && !data.MultiMonitorMode.Equal(state.MultiMonitorMode) {
        requestDataMap["multiMonitorMode"] = data.MultiMonitorMode.ValueString()
    }
    if !data.Monitors.IsUnknown() && !state.Monitors.IsUnknown() && !data.Monitors.Equal(state.Monitors) {
        requestDataMap["monitors"] = r.convertTerraformSetToInterface(data.Monitors)
    }
    if !data.MonitorLabels.IsUnknown() && !state.MonitorLabels.IsUnknown() && !data.MonitorLabels.Equal(state.MonitorLabels) {
        requestDataMap["monitorLabels"] = r.convertTerraformSetToInterface(data.MonitorLabels)
    }
    if !data.MetricQueryConfig.IsUnknown() && !state.MetricQueryConfig.IsUnknown() && !data.MetricQueryConfig.Equal(state.MetricQueryConfig) {
        var metricqueryconfigData interface{}
        if err := json.Unmarshal([]byte(data.MetricQueryConfig.ValueString()), &metricqueryconfigData); err == nil {
            requestDataMap["metricQueryConfig"] = metricqueryconfigData
        } else {
            requestDataMap["metricQueryConfig"] = data.MetricQueryConfig.ValueString()
        }
    }
    if !data.TargetPercentage.IsUnknown() && !state.TargetPercentage.IsUnknown() && !data.TargetPercentage.Equal(state.TargetPercentage) {
        requestDataMap["targetPercentage"] = r.bigFloatToFloat64(data.TargetPercentage.ValueBigFloat())
    }
    if !data.WindowType.IsUnknown() && !state.WindowType.IsUnknown() && !data.WindowType.Equal(state.WindowType) {
        requestDataMap["windowType"] = data.WindowType.ValueString()
    }
    if !data.WindowDays.IsUnknown() && !state.WindowDays.IsUnknown() && !data.WindowDays.Equal(state.WindowDays) {
        requestDataMap["windowDays"] = r.bigFloatToFloat64(data.WindowDays.ValueBigFloat())
    }
    if !data.Timezone.IsUnknown() && !state.Timezone.IsUnknown() && !data.Timezone.Equal(state.Timezone) {
        requestDataMap["timezone"] = data.Timezone.ValueString()
    }
    if !data.AtRiskThresholdPercentage.IsUnknown() && !state.AtRiskThresholdPercentage.IsUnknown() && !data.AtRiskThresholdPercentage.Equal(state.AtRiskThresholdPercentage) {
        requestDataMap["atRiskThresholdPercentage"] = r.bigFloatToFloat64(data.AtRiskThresholdPercentage.ValueBigFloat())
    }

    // Only call the API when there are changed fields to send. An empty
    // update body is rejected by the API; state is still refreshed below so
    // this method never writes unverified plan values into state.
    if len(serviceLevelObjectiveRequest["data"].(map[string]interface{})) > 0 {
        httpResp, err := r.client.Put(ctx, "/service-level-objective/" + data.Id.ValueString() + "", serviceLevelObjectiveRequest)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update service_level_objective, got error: %s", err))
            return
        }

        // Parse the update response
        var serviceLevelObjectiveResponse map[string]interface{}
        err = r.client.ParseResponse(httpResp, &serviceLevelObjectiveResponse)
        if err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to update service_level_objective: %s", err))
            return
        }
        _ = serviceLevelObjectiveResponse
    }

    // After successful update, fetch the current state by calling Read with select parameter
    selectParam := map[string]interface{}{
        "projectId": true,
        "name": true,
        "description": true,
        "labels": true,
        "isEnabled": true,
        "sliType": true,
        "multiMonitorMode": true,
        "monitors": true,
        "monitorLabels": true,
        "metricQueryConfig": true,
        "targetPercentage": true,
        "windowType": true,
        "windowDays": true,
        "timezone": true,
        "atRiskThresholdPercentage": true,
        "createdByUserId": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "slug": true,
        "autoAddedMonitors": true,
        "downtimeMonitorStatuses": true,
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
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect(ctx, "/service-level-objective/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read service_level_objective after update, got error: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read service_level_objective after update: %s", err))
        return
    }

    // Update the model with response data from the Read operation
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := readResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = readResponse
    }

    if obj, ok := dataMap["projectId"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok {
            data.ProjectId = types.StringValue(val)
        } else {
            data.ProjectId = types.StringNull()
        }
    } else if val, ok := dataMap["projectId"].(string); ok {
        data.ProjectId = types.StringValue(val)
    } else {
        data.ProjectId = types.StringNull()
    }
    if obj, ok := dataMap["name"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Name = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Name = types.StringValue(string(jsonBytes))
            } else {
                data.Name = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Name = types.StringValue(string(jsonBytes))
            } else {
                data.Name = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Name = types.StringValue(string(jsonBytes))
        } else {
            data.Name = types.StringNull()
        }
    } else if val, ok := dataMap["name"].(string); ok {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if obj, ok := dataMap["description"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Description = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Description = types.StringValue(string(jsonBytes))
            } else {
                data.Description = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Description = types.StringValue(string(jsonBytes))
            } else {
                data.Description = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Description = types.StringValue(string(jsonBytes))
        } else {
            data.Description = types.StringNull()
        }
    } else if val, ok := dataMap["description"].(string); ok {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if val, ok := dataMap["labels"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.Labels = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.Labels = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if val, ok := dataMap["isEnabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    }
    if obj, ok := dataMap["sliType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SliType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SliType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SliType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SliType = types.StringValue(string(jsonBytes))
            } else {
                data.SliType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SliType = types.StringValue(string(jsonBytes))
            } else {
                data.SliType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SliType = types.StringValue(string(jsonBytes))
        } else {
            data.SliType = types.StringNull()
        }
    } else if val, ok := dataMap["sliType"].(string); ok {
        data.SliType = types.StringValue(val)
    } else {
        data.SliType = types.StringNull()
    }
    if obj, ok := dataMap["multiMonitorMode"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MultiMonitorMode = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.MultiMonitorMode = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.MultiMonitorMode = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.MultiMonitorMode = types.StringValue(string(jsonBytes))
            } else {
                data.MultiMonitorMode = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.MultiMonitorMode = types.StringValue(string(jsonBytes))
            } else {
                data.MultiMonitorMode = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.MultiMonitorMode = types.StringValue(string(jsonBytes))
        } else {
            data.MultiMonitorMode = types.StringNull()
        }
    } else if val, ok := dataMap["multiMonitorMode"].(string); ok {
        data.MultiMonitorMode = types.StringValue(val)
    } else {
        data.MultiMonitorMode = types.StringNull()
    }
    if val, ok := dataMap["monitors"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.Monitors = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.Monitors = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if val, ok := dataMap["monitorLabels"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.MonitorLabels = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.MonitorLabels = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if obj, ok := dataMap["metricQueryConfig"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MetricQueryConfig = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.MetricQueryConfig = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.MetricQueryConfig = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.MetricQueryConfig = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.MetricQueryConfig = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.MetricQueryConfig = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.MetricQueryConfig = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.MetricQueryConfig = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.MetricQueryConfig = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["metricQueryConfig"].(string); ok {
        data.MetricQueryConfig = NewJSONSubsetValue(val)
    } else {
        data.MetricQueryConfig = NewJSONSubsetNull()
    }
    if val, ok := dataMap["targetPercentage"].(float64); ok {
        data.TargetPercentage = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["targetPercentage"].(int); ok {
        data.TargetPercentage = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["targetPercentage"].(int64); ok {
        data.TargetPercentage = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["targetPercentage"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.TargetPercentage = types.NumberValue(big.NewFloat(val))
        } else {
            data.TargetPercentage = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.TargetPercentage = types.NumberNull()
    }
    if obj, ok := dataMap["windowType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.WindowType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.WindowType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.WindowType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.WindowType = types.StringValue(string(jsonBytes))
            } else {
                data.WindowType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.WindowType = types.StringValue(string(jsonBytes))
            } else {
                data.WindowType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.WindowType = types.StringValue(string(jsonBytes))
        } else {
            data.WindowType = types.StringNull()
        }
    } else if val, ok := dataMap["windowType"].(string); ok {
        data.WindowType = types.StringValue(val)
    } else {
        data.WindowType = types.StringNull()
    }
    if val, ok := dataMap["windowDays"].(float64); ok {
        data.WindowDays = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["windowDays"].(int); ok {
        data.WindowDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["windowDays"].(int64); ok {
        data.WindowDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["windowDays"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.WindowDays = types.NumberValue(big.NewFloat(val))
        } else {
            data.WindowDays = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.WindowDays = types.NumberNull()
    }
    if obj, ok := dataMap["timezone"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Timezone = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Timezone = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Timezone = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Timezone = types.StringValue(string(jsonBytes))
            } else {
                data.Timezone = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Timezone = types.StringValue(string(jsonBytes))
            } else {
                data.Timezone = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Timezone = types.StringValue(string(jsonBytes))
        } else {
            data.Timezone = types.StringNull()
        }
    } else if val, ok := dataMap["timezone"].(string); ok {
        data.Timezone = types.StringValue(val)
    } else {
        data.Timezone = types.StringNull()
    }
    if val, ok := dataMap["atRiskThresholdPercentage"].(float64); ok {
        data.AtRiskThresholdPercentage = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["atRiskThresholdPercentage"].(int); ok {
        data.AtRiskThresholdPercentage = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["atRiskThresholdPercentage"].(int64); ok {
        data.AtRiskThresholdPercentage = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["atRiskThresholdPercentage"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.AtRiskThresholdPercentage = types.NumberValue(big.NewFloat(val))
        } else {
            data.AtRiskThresholdPercentage = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.AtRiskThresholdPercentage = types.NumberNull()
    }
    if obj, ok := dataMap["createdByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["createdByUserId"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if obj, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.CreatedAt = NewRFC3339Value(val)
        } else {
            data.CreatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["createdAt"].(string); ok && val != "" {
        data.CreatedAt = NewRFC3339Value(val)
    } else {
        data.CreatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.UpdatedAt = NewRFC3339Value(val)
        } else {
            data.UpdatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["updatedAt"].(string); ok && val != "" {
        data.UpdatedAt = NewRFC3339Value(val)
    } else {
        data.UpdatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["deletedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.DeletedAt = NewRFC3339Value(val)
        } else {
            data.DeletedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["deletedAt"].(string); ok && val != "" {
        data.DeletedAt = NewRFC3339Value(val)
    } else {
        data.DeletedAt = NewRFC3339Null()
    }
    if val, ok := dataMap["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["version"].(int); ok {
        data.Version = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["version"].(int64); ok {
        data.Version = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["version"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.Version = types.NumberValue(big.NewFloat(val))
        } else {
            data.Version = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.Version = types.NumberNull()
    }
    if obj, ok := dataMap["slug"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Slug = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Slug = types.StringValue(string(jsonBytes))
            } else {
                data.Slug = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Slug = types.StringValue(string(jsonBytes))
            } else {
                data.Slug = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Slug = types.StringValue(string(jsonBytes))
        } else {
            data.Slug = types.StringNull()
        }
    } else if val, ok := dataMap["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
    }
    if val, ok := dataMap["autoAddedMonitors"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.AutoAddedMonitors = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.AutoAddedMonitors = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if val, ok := dataMap["downtimeMonitorStatuses"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.DowntimeMonitorStatuses = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.DowntimeMonitorStatuses = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if val, ok := dataMap["currentSliPercentage"].(float64); ok {
        data.CurrentSliPercentage = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["currentSliPercentage"].(int); ok {
        data.CurrentSliPercentage = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["currentSliPercentage"].(int64); ok {
        data.CurrentSliPercentage = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["currentSliPercentage"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.CurrentSliPercentage = types.NumberValue(big.NewFloat(val))
        } else {
            data.CurrentSliPercentage = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.CurrentSliPercentage = types.NumberNull()
    }
    if val, ok := dataMap["errorBudgetRemainingPercentage"].(float64); ok {
        data.ErrorBudgetRemainingPercentage = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["errorBudgetRemainingPercentage"].(int); ok {
        data.ErrorBudgetRemainingPercentage = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["errorBudgetRemainingPercentage"].(int64); ok {
        data.ErrorBudgetRemainingPercentage = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["errorBudgetRemainingPercentage"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.ErrorBudgetRemainingPercentage = types.NumberValue(big.NewFloat(val))
        } else {
            data.ErrorBudgetRemainingPercentage = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.ErrorBudgetRemainingPercentage = types.NumberNull()
    }
    if val, ok := dataMap["errorBudgetRemainingSeconds"].(float64); ok {
        data.ErrorBudgetRemainingSeconds = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["errorBudgetRemainingSeconds"].(int); ok {
        data.ErrorBudgetRemainingSeconds = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["errorBudgetRemainingSeconds"].(int64); ok {
        data.ErrorBudgetRemainingSeconds = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["errorBudgetRemainingSeconds"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.ErrorBudgetRemainingSeconds = types.NumberValue(big.NewFloat(val))
        } else {
            data.ErrorBudgetRemainingSeconds = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.ErrorBudgetRemainingSeconds = types.NumberNull()
    }
    if val, ok := dataMap["errorBudgetTotalSeconds"].(float64); ok {
        data.ErrorBudgetTotalSeconds = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["errorBudgetTotalSeconds"].(int); ok {
        data.ErrorBudgetTotalSeconds = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["errorBudgetTotalSeconds"].(int64); ok {
        data.ErrorBudgetTotalSeconds = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["errorBudgetTotalSeconds"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.ErrorBudgetTotalSeconds = types.NumberValue(big.NewFloat(val))
        } else {
            data.ErrorBudgetTotalSeconds = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.ErrorBudgetTotalSeconds = types.NumberNull()
    }
    if val, ok := dataMap["currentBurnRate"].(float64); ok {
        data.CurrentBurnRate = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["currentBurnRate"].(int); ok {
        data.CurrentBurnRate = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["currentBurnRate"].(int64); ok {
        data.CurrentBurnRate = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["currentBurnRate"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.CurrentBurnRate = types.NumberValue(big.NewFloat(val))
        } else {
            data.CurrentBurnRate = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.CurrentBurnRate = types.NumberNull()
    }
    if obj, ok := dataMap["sloStatus"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SloStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SloStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SloStatus = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SloStatus = types.StringValue(string(jsonBytes))
            } else {
                data.SloStatus = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SloStatus = types.StringValue(string(jsonBytes))
            } else {
                data.SloStatus = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SloStatus = types.StringValue(string(jsonBytes))
        } else {
            data.SloStatus = types.StringNull()
        }
    } else if val, ok := dataMap["sloStatus"].(string); ok {
        data.SloStatus = types.StringValue(val)
    } else {
        data.SloStatus = types.StringNull()
    }
    if obj, ok := dataMap["statusChangeNotificationSentAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.StatusChangeNotificationSentAt = NewRFC3339Value(val)
        } else {
            data.StatusChangeNotificationSentAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["statusChangeNotificationSentAt"].(string); ok && val != "" {
        data.StatusChangeNotificationSentAt = NewRFC3339Value(val)
    } else {
        data.StatusChangeNotificationSentAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["lastEvaluatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.LastEvaluatedAt = NewRFC3339Value(val)
        } else {
            data.LastEvaluatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["lastEvaluatedAt"].(string); ok && val != "" {
        data.LastEvaluatedAt = NewRFC3339Value(val)
    } else {
        data.LastEvaluatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["nextEvaluationAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.NextEvaluationAt = NewRFC3339Value(val)
        } else {
            data.NextEvaluationAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["nextEvaluationAt"].(string); ok && val != "" {
        data.NextEvaluationAt = NewRFC3339Value(val)
    } else {
        data.NextEvaluationAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["lastAccumulatedBucketEndAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.LastAccumulatedBucketEndAt = NewRFC3339Value(val)
        } else {
            data.LastAccumulatedBucketEndAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["lastAccumulatedBucketEndAt"].(string); ok && val != "" {
        data.LastAccumulatedBucketEndAt = NewRFC3339Value(val)
    } else {
        data.LastAccumulatedBucketEndAt = NewRFC3339Null()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    data.Id = state.Id

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ServiceLevelObjectiveResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var data ServiceLevelObjectiveResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Make API call
    httpResp, err := r.client.Delete(ctx, "/service-level-objective/" + data.Id.ValueString() + "")
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete service_level_objective, got error: %s", err))
        return
    }

    // A failed delete must keep the resource in state — silently dropping it
    // orphans real infrastructure. 404 means it is already gone.
    if httpResp.StatusCode >= 400 && httpResp.StatusCode != http.StatusNotFound {
        err = r.client.ParseResponse(httpResp, nil)
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to delete service_level_objective: %s", err))
        return
    }
    if httpResp.Body != nil {
        httpResp.Body.Close()
    }
}


func (r *ServiceLevelObjectiveResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Helper method to convert Terraform map to Go interface{}
func (r *ServiceLevelObjectiveResource) convertTerraformMapToInterface(terraformMap types.Map) interface{} {
    if terraformMap.IsNull() || terraformMap.IsUnknown() {
        return nil
    }
    
    result := make(map[string]string)
    terraformMap.ElementsAs(context.Background(), &result, false)
    
    // Convert map[string]string to map[string]interface{}
    interfaceResult := make(map[string]interface{})
    for key, value := range result {
        interfaceResult[key] = value
    }
    
    return interfaceResult
}

// Helper method to convert Terraform list to Go interface{}
func (r *ServiceLevelObjectiveResource) convertTerraformListToInterface(terraformList types.List) interface{} {
    if terraformList.IsNull() || terraformList.IsUnknown() {
        return nil
    }

    var stringList []string
    terraformList.ElementsAs(context.Background(), &stringList, false)

    // Convert string array to OneUptime format with _id fields
    var result []interface{}
    for _, str := range stringList {
        if str != "" {
            result = append(result, map[string]interface{}{
                "_id": str,
            })
        }
    }
    return result
}

// Helper method to convert Terraform set to Go interface{}
func (r *ServiceLevelObjectiveResource) convertTerraformSetToInterface(terraformSet types.Set) interface{} {
    if terraformSet.IsNull() || terraformSet.IsUnknown() {
        return nil
    }

    var stringList []string
    terraformSet.ElementsAs(context.Background(), &stringList, false)

    // Convert string array to OneUptime format with _id fields
    var result []interface{}
    for _, str := range stringList {
        if str != "" {
            result = append(result, map[string]interface{}{
                "_id": str,
            })
        }
    }
    return result
}


// Helper method to parse JSON field for complex objects
func (r *ServiceLevelObjectiveResource) parseJSONField(terraformString basetypes.StringValuable) interface{} {
    sv, _ := terraformString.ToStringValue(context.Background())
    if sv.IsNull() || sv.IsUnknown() || sv.ValueString() == "" {
        return nil
    }

    var result interface{}
    if err := json.Unmarshal([]byte(sv.ValueString()), &result); err != nil {
        // If JSON parsing fails, return the raw string
        return sv.ValueString()
    }

    return result
}

// Normalize URL wrapper objects to avoid drift (e.g., trailing slash differences).
func (r *ServiceLevelObjectiveResource) normalizeURLWrappers(value interface{}) interface{} {
    switch v := value.(type) {
    case map[string]interface{}:
        if typeStr, ok := v["_type"].(string); ok && typeStr == "URL" {
            if val, ok := v["value"].(string); ok {
                v["value"] = r.normalizeURLString(val)
            }
        }
        for key, child := range v {
            v[key] = r.normalizeURLWrappers(child)
        }
        return v
    case []interface{}:
        for i, child := range v {
            v[i] = r.normalizeURLWrappers(child)
        }
        return v
    default:
        return v
    }
}

func (r *ServiceLevelObjectiveResource) normalizeURLString(value string) string {
    parsed, err := url.Parse(value)
    if err != nil {
        return value
    }
    if parsed.Path == "/" && parsed.RawQuery == "" && parsed.Fragment == "" {
        return strings.TrimSuffix(value, "/")
    }
    return value
}

// Helper method to convert *big.Float to float64 for JSON serialization
func (r *ServiceLevelObjectiveResource) bigFloatToFloat64(bf *big.Float) interface{} {
    if bf == nil {
        return nil
    }
    f, _ := bf.Float64()
    return f
}

// Helper method to check if a type string is a valid OneUptime ObjectType.
// The registry itself lives in objecttypes.go, shared across the package.
func (r *ServiceLevelObjectiveResource) isValidOneUptimeObjectType(typeStr string) bool {
    return validOneUptimeObjectTypes[typeStr]
}
