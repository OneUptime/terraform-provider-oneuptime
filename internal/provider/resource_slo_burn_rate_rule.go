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
    "github.com/hashicorp/terraform-plugin-framework/attr"
    "sort"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/numberplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &SloBurnRateRuleResource{}
var _ resource.ResourceWithImportState = &SloBurnRateRuleResource{}

func NewSloBurnRateRuleResource() resource.Resource {
    return &SloBurnRateRuleResource{}
}

// SloBurnRateRuleResource defines the resource implementation.
type SloBurnRateRuleResource struct {
    client *Client
}

// SloBurnRateRuleResourceModel describes the resource data model.
type SloBurnRateRuleResourceModel struct {
    Id types.String `tfsdk:"id"`
    ProjectId types.String `tfsdk:"project_id"`
    ServiceLevelObjectiveId types.String `tfsdk:"service_level_objective_id"`
    Name types.String `tfsdk:"name"`
    IsEnabled types.Bool `tfsdk:"is_enabled"`
    BurnRateThreshold types.Number `tfsdk:"burn_rate_threshold"`
    LongWindowInMinutes types.Number `tfsdk:"long_window_in_minutes"`
    ShortWindowInMinutes types.Number `tfsdk:"short_window_in_minutes"`
    MinimumSampleCount types.Number `tfsdk:"minimum_sample_count"`
    RefireSuppressionMinutes types.Number `tfsdk:"refire_suppression_minutes"`
    AlertSeverityId types.String `tfsdk:"alert_severity_id"`
    OnCallDutyPolicies types.Set `tfsdk:"on_call_duty_policies"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    CreatedAt RFC3339Value `tfsdk:"created_at"`
    UpdatedAt RFC3339Value `tfsdk:"updated_at"`
    DeletedAt RFC3339Value `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    LastAlertCreatedAt RFC3339Value `tfsdk:"last_alert_created_at"`
    LastAlertResolvedAt RFC3339Value `tfsdk:"last_alert_resolved_at"`
}

func (r *SloBurnRateRuleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_slo_burn_rate_rule"
}

func (r *SloBurnRateRuleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Configure multi-window burn rate rules that raise alerts when a Service Level Objective consumes its error budget too quickly",

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
            "service_level_objective_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Required: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.RequiresReplace(),
                },
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Name of this burn rate rule.",
                Required: true,
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether this burn rate rule is enabled.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(true),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "burn_rate_threshold": schema.NumberAttribute{
                MarkdownDescription: "Alert when the burn rate in both the long and short windows is at or above this threshold (e.g. 14.4)..",
                Required: true,
            },
            "long_window_in_minutes": schema.NumberAttribute{
                MarkdownDescription: "Length of the long lookback window in minutes (e.g. 60). The alert fires when both windows exceed the threshold and resolves when the long window drops below it..",
                Required: true,
            },
            "short_window_in_minutes": schema.NumberAttribute{
                MarkdownDescription: "Length of the short lookback window in minutes (e.g. 5). Guards against alerting on burn that has already stopped..",
                Required: true,
            },
            "minimum_sample_count": schema.NumberAttribute{
                MarkdownDescription: "For event-based SLIs only: skip this rule when the long window has fewer than this many total events. Prevents noisy alerts on low traffic..",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "refire_suppression_minutes": schema.NumberAttribute{
                MarkdownDescription: "Minimum number of minutes after an alert resolves before this rule can fire again. Defaults to the long window length when not set..",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "alert_severity_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "on_call_duty_policies": schema.SetAttribute{
                MarkdownDescription: "On-call duty policies attached to alerts created by this burn rate rule..",
                Optional: true,
                Computed: true,
                ElementType: types.StringType,
                PlanModifiers: []planmodifier.Set{
                    setplanmodifier.UseStateForUnknown(),
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
            "last_alert_created_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "last_alert_resolved_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
        },
    }
}

func (r *SloBurnRateRuleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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


func (r *SloBurnRateRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var data SloBurnRateRuleResourceModel

    // Read Terraform plan data into the model
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }



    // Create API request body. Unset (null/unknown) optional fields are
    // omitted so server-side defaults apply instead of being overwritten
    // with zero values.
    sloBurnRateRuleRequest := map[string]interface{}{
        "data": map[string]interface{}{},
    }
    requestDataMap := sloBurnRateRuleRequest["data"].(map[string]interface{})

    if !data.ServiceLevelObjectiveId.IsNull() && !data.ServiceLevelObjectiveId.IsUnknown() {
        requestDataMap["serviceLevelObjectiveId"] = data.ServiceLevelObjectiveId.ValueString()
    }
    if !data.Name.IsNull() && !data.Name.IsUnknown() {
        requestDataMap["name"] = data.Name.ValueString()
    }
    if !data.IsEnabled.IsNull() && !data.IsEnabled.IsUnknown() {
        requestDataMap["isEnabled"] = data.IsEnabled.ValueBool()
    }
    if !data.BurnRateThreshold.IsNull() && !data.BurnRateThreshold.IsUnknown() {
        requestDataMap["burnRateThreshold"] = r.bigFloatToFloat64(data.BurnRateThreshold.ValueBigFloat())
    }
    if !data.LongWindowInMinutes.IsNull() && !data.LongWindowInMinutes.IsUnknown() {
        requestDataMap["longWindowInMinutes"] = r.bigFloatToFloat64(data.LongWindowInMinutes.ValueBigFloat())
    }
    if !data.ShortWindowInMinutes.IsNull() && !data.ShortWindowInMinutes.IsUnknown() {
        requestDataMap["shortWindowInMinutes"] = r.bigFloatToFloat64(data.ShortWindowInMinutes.ValueBigFloat())
    }
    if !data.MinimumSampleCount.IsNull() && !data.MinimumSampleCount.IsUnknown() {
        requestDataMap["minimumSampleCount"] = r.bigFloatToFloat64(data.MinimumSampleCount.ValueBigFloat())
    }
    if !data.RefireSuppressionMinutes.IsNull() && !data.RefireSuppressionMinutes.IsUnknown() {
        requestDataMap["refireSuppressionMinutes"] = r.bigFloatToFloat64(data.RefireSuppressionMinutes.ValueBigFloat())
    }
    if !data.AlertSeverityId.IsNull() && !data.AlertSeverityId.IsUnknown() {
        requestDataMap["alertSeverityId"] = data.AlertSeverityId.ValueString()
    }
    if !data.OnCallDutyPolicies.IsNull() && !data.OnCallDutyPolicies.IsUnknown() {
        requestDataMap["onCallDutyPolicies"] = r.convertTerraformSetToInterface(data.OnCallDutyPolicies)
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        requestDataMap["createdByUserId"] = data.CreatedByUserId.ValueString()
    }

    // Make API call
    httpResp, err := r.client.Post(ctx, "/service-level-objective-burn-rate-rule", sloBurnRateRuleRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create slo_burn_rate_rule, got error: %s", err))
        return
    }

    var sloBurnRateRuleResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &sloBurnRateRuleResponse)
    if err != nil {
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to create slo_burn_rate_rule: %s", err))
        return
    }

    // Extract the new resource id from the create response.
    createdId := ""
    if wrapper, ok := sloBurnRateRuleResponse["data"].(map[string]interface{}); ok {
        if val, ok := wrapper["_id"].(string); ok {
            createdId = val
        }
    } else if val, ok := sloBurnRateRuleResponse["_id"].(string); ok {
        createdId = val
    }
    if createdId == "" {
        resp.Diagnostics.AddError("OneUptime API Error", "Create response for slo_burn_rate_rule did not contain an id. This is a bug in the provider or the API; please report it.")
        return
    }
    data.Id = types.StringValue(createdId)

    /*
     * The server has committed the row. Persist what we know to state BEFORE
     * the read-back: if the read-back fails and we return without setting
     * state, Terraform never learns the resource exists and the created
     * slo_burn_rate_rule is orphaned server-side — never refreshed, never
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
        "serviceLevelObjectiveId": true,
        "name": true,
        "isEnabled": true,
        "burnRateThreshold": true,
        "longWindowInMinutes": true,
        "shortWindowInMinutes": true,
        "minimumSampleCount": true,
        "refireSuppressionMinutes": true,
        "alertSeverityId": true,
        "onCallDutyPolicies": true,
        "createdByUserId": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "lastAlertCreatedAt": true,
        "lastAlertResolvedAt": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect(ctx, "/service-level-objective-burn-rate-rule/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        /*
         * State already owns the id, so the resource is tracked and the next
         * refresh reconciles the remaining attributes. Warn rather than
         * error: erroring here would strand a real resource.
         */
        resp.Diagnostics.AddWarning("Read After Create Failed", fmt.Sprintf("Created slo_burn_rate_rule but could not read it back; state is incomplete until the next refresh: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddWarning("Read After Create Failed", fmt.Sprintf("Created slo_burn_rate_rule but could not parse the read-back response; state is incomplete until the next refresh: %s", err))
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
    if obj, ok := dataMap["serviceLevelObjectiveId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ServiceLevelObjectiveId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ServiceLevelObjectiveId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ServiceLevelObjectiveId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ServiceLevelObjectiveId = types.StringValue(string(jsonBytes))
            } else {
                data.ServiceLevelObjectiveId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ServiceLevelObjectiveId = types.StringValue(string(jsonBytes))
            } else {
                data.ServiceLevelObjectiveId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ServiceLevelObjectiveId = types.StringValue(string(jsonBytes))
        } else {
            data.ServiceLevelObjectiveId = types.StringNull()
        }
    } else if val, ok := dataMap["serviceLevelObjectiveId"].(string); ok {
        data.ServiceLevelObjectiveId = types.StringValue(val)
    } else {
        data.ServiceLevelObjectiveId = types.StringNull()
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
    if val, ok := dataMap["isEnabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    }
    if val, ok := dataMap["burnRateThreshold"].(float64); ok {
        data.BurnRateThreshold = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["burnRateThreshold"].(int); ok {
        data.BurnRateThreshold = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["burnRateThreshold"].(int64); ok {
        data.BurnRateThreshold = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["burnRateThreshold"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.BurnRateThreshold = types.NumberValue(big.NewFloat(val))
        } else {
            data.BurnRateThreshold = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.BurnRateThreshold = types.NumberNull()
    }
    if val, ok := dataMap["longWindowInMinutes"].(float64); ok {
        data.LongWindowInMinutes = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["longWindowInMinutes"].(int); ok {
        data.LongWindowInMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["longWindowInMinutes"].(int64); ok {
        data.LongWindowInMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["longWindowInMinutes"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.LongWindowInMinutes = types.NumberValue(big.NewFloat(val))
        } else {
            data.LongWindowInMinutes = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.LongWindowInMinutes = types.NumberNull()
    }
    if val, ok := dataMap["shortWindowInMinutes"].(float64); ok {
        data.ShortWindowInMinutes = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["shortWindowInMinutes"].(int); ok {
        data.ShortWindowInMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["shortWindowInMinutes"].(int64); ok {
        data.ShortWindowInMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["shortWindowInMinutes"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.ShortWindowInMinutes = types.NumberValue(big.NewFloat(val))
        } else {
            data.ShortWindowInMinutes = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.ShortWindowInMinutes = types.NumberNull()
    }
    if val, ok := dataMap["minimumSampleCount"].(float64); ok {
        data.MinimumSampleCount = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["minimumSampleCount"].(int); ok {
        data.MinimumSampleCount = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["minimumSampleCount"].(int64); ok {
        data.MinimumSampleCount = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["minimumSampleCount"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.MinimumSampleCount = types.NumberValue(big.NewFloat(val))
        } else {
            data.MinimumSampleCount = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.MinimumSampleCount = types.NumberNull()
    }
    if val, ok := dataMap["refireSuppressionMinutes"].(float64); ok {
        data.RefireSuppressionMinutes = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["refireSuppressionMinutes"].(int); ok {
        data.RefireSuppressionMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["refireSuppressionMinutes"].(int64); ok {
        data.RefireSuppressionMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["refireSuppressionMinutes"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.RefireSuppressionMinutes = types.NumberValue(big.NewFloat(val))
        } else {
            data.RefireSuppressionMinutes = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.RefireSuppressionMinutes = types.NumberNull()
    }
    if obj, ok := dataMap["alertSeverityId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AlertSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AlertSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AlertSeverityId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AlertSeverityId = types.StringValue(string(jsonBytes))
            } else {
                data.AlertSeverityId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AlertSeverityId = types.StringValue(string(jsonBytes))
            } else {
                data.AlertSeverityId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AlertSeverityId = types.StringValue(string(jsonBytes))
        } else {
            data.AlertSeverityId = types.StringNull()
        }
    } else if val, ok := dataMap["alertSeverityId"].(string); ok {
        data.AlertSeverityId = types.StringValue(val)
    } else {
        data.AlertSeverityId = types.StringNull()
    }
    if val, ok := dataMap["onCallDutyPolicies"].([]interface{}); ok {
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
        data.OnCallDutyPolicies = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.OnCallDutyPolicies = types.SetValueMust(types.StringType, []attr.Value{})
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
    if obj, ok := dataMap["lastAlertCreatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.LastAlertCreatedAt = NewRFC3339Value(val)
        } else {
            data.LastAlertCreatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["lastAlertCreatedAt"].(string); ok && val != "" {
        data.LastAlertCreatedAt = NewRFC3339Value(val)
    } else {
        data.LastAlertCreatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["lastAlertResolvedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.LastAlertResolvedAt = NewRFC3339Value(val)
        } else {
            data.LastAlertResolvedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["lastAlertResolvedAt"].(string); ok && val != "" {
        data.LastAlertResolvedAt = NewRFC3339Value(val)
    } else {
        data.LastAlertResolvedAt = NewRFC3339Null()
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

func (r *SloBurnRateRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var data SloBurnRateRuleResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create select parameter to get full object
    selectParam := map[string]interface{}{
        "projectId": true,
        "serviceLevelObjectiveId": true,
        "name": true,
        "isEnabled": true,
        "burnRateThreshold": true,
        "longWindowInMinutes": true,
        "shortWindowInMinutes": true,
        "minimumSampleCount": true,
        "refireSuppressionMinutes": true,
        "alertSeverityId": true,
        "onCallDutyPolicies": true,
        "createdByUserId": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "lastAlertCreatedAt": true,
        "lastAlertResolvedAt": true,
        "_id": true,
    }

    // Make API call with select parameter
    httpResp, err := r.client.PostWithSelect(ctx, "/service-level-objective-burn-rate-rule/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read slo_burn_rate_rule, got error: %s", err))
        return
    }

    if httpResp.StatusCode == http.StatusNotFound {
        resp.State.RemoveResource(ctx)
        return
    }

    var sloBurnRateRuleResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &sloBurnRateRuleResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse slo_burn_rate_rule response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := sloBurnRateRuleResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = sloBurnRateRuleResponse
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
    if obj, ok := dataMap["serviceLevelObjectiveId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ServiceLevelObjectiveId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ServiceLevelObjectiveId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ServiceLevelObjectiveId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ServiceLevelObjectiveId = types.StringValue(string(jsonBytes))
            } else {
                data.ServiceLevelObjectiveId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ServiceLevelObjectiveId = types.StringValue(string(jsonBytes))
            } else {
                data.ServiceLevelObjectiveId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ServiceLevelObjectiveId = types.StringValue(string(jsonBytes))
        } else {
            data.ServiceLevelObjectiveId = types.StringNull()
        }
    } else if val, ok := dataMap["serviceLevelObjectiveId"].(string); ok {
        data.ServiceLevelObjectiveId = types.StringValue(val)
    } else {
        data.ServiceLevelObjectiveId = types.StringNull()
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
    if val, ok := dataMap["isEnabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    }
    if val, ok := dataMap["burnRateThreshold"].(float64); ok {
        data.BurnRateThreshold = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["burnRateThreshold"].(int); ok {
        data.BurnRateThreshold = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["burnRateThreshold"].(int64); ok {
        data.BurnRateThreshold = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["burnRateThreshold"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.BurnRateThreshold = types.NumberValue(big.NewFloat(val))
        } else {
            data.BurnRateThreshold = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.BurnRateThreshold = types.NumberNull()
    }
    if val, ok := dataMap["longWindowInMinutes"].(float64); ok {
        data.LongWindowInMinutes = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["longWindowInMinutes"].(int); ok {
        data.LongWindowInMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["longWindowInMinutes"].(int64); ok {
        data.LongWindowInMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["longWindowInMinutes"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.LongWindowInMinutes = types.NumberValue(big.NewFloat(val))
        } else {
            data.LongWindowInMinutes = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.LongWindowInMinutes = types.NumberNull()
    }
    if val, ok := dataMap["shortWindowInMinutes"].(float64); ok {
        data.ShortWindowInMinutes = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["shortWindowInMinutes"].(int); ok {
        data.ShortWindowInMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["shortWindowInMinutes"].(int64); ok {
        data.ShortWindowInMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["shortWindowInMinutes"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.ShortWindowInMinutes = types.NumberValue(big.NewFloat(val))
        } else {
            data.ShortWindowInMinutes = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.ShortWindowInMinutes = types.NumberNull()
    }
    if val, ok := dataMap["minimumSampleCount"].(float64); ok {
        data.MinimumSampleCount = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["minimumSampleCount"].(int); ok {
        data.MinimumSampleCount = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["minimumSampleCount"].(int64); ok {
        data.MinimumSampleCount = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["minimumSampleCount"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.MinimumSampleCount = types.NumberValue(big.NewFloat(val))
        } else {
            data.MinimumSampleCount = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.MinimumSampleCount = types.NumberNull()
    }
    if val, ok := dataMap["refireSuppressionMinutes"].(float64); ok {
        data.RefireSuppressionMinutes = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["refireSuppressionMinutes"].(int); ok {
        data.RefireSuppressionMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["refireSuppressionMinutes"].(int64); ok {
        data.RefireSuppressionMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["refireSuppressionMinutes"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.RefireSuppressionMinutes = types.NumberValue(big.NewFloat(val))
        } else {
            data.RefireSuppressionMinutes = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.RefireSuppressionMinutes = types.NumberNull()
    }
    if obj, ok := dataMap["alertSeverityId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AlertSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AlertSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AlertSeverityId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AlertSeverityId = types.StringValue(string(jsonBytes))
            } else {
                data.AlertSeverityId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AlertSeverityId = types.StringValue(string(jsonBytes))
            } else {
                data.AlertSeverityId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AlertSeverityId = types.StringValue(string(jsonBytes))
        } else {
            data.AlertSeverityId = types.StringNull()
        }
    } else if val, ok := dataMap["alertSeverityId"].(string); ok {
        data.AlertSeverityId = types.StringValue(val)
    } else {
        data.AlertSeverityId = types.StringNull()
    }
    if val, ok := dataMap["onCallDutyPolicies"].([]interface{}); ok {
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
        data.OnCallDutyPolicies = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.OnCallDutyPolicies = types.SetValueMust(types.StringType, []attr.Value{})
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
    if obj, ok := dataMap["lastAlertCreatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.LastAlertCreatedAt = NewRFC3339Value(val)
        } else {
            data.LastAlertCreatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["lastAlertCreatedAt"].(string); ok && val != "" {
        data.LastAlertCreatedAt = NewRFC3339Value(val)
    } else {
        data.LastAlertCreatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["lastAlertResolvedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.LastAlertResolvedAt = NewRFC3339Value(val)
        } else {
            data.LastAlertResolvedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["lastAlertResolvedAt"].(string); ok && val != "" {
        data.LastAlertResolvedAt = NewRFC3339Value(val)
    } else {
        data.LastAlertResolvedAt = NewRFC3339Null()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SloBurnRateRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var data SloBurnRateRuleResourceModel
    var state SloBurnRateRuleResourceModel

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
    sloBurnRateRuleRequest := map[string]interface{}{
        "data": map[string]interface{}{},
    }
    requestDataMap := sloBurnRateRuleRequest["data"].(map[string]interface{})

    if !data.Name.IsUnknown() && !state.Name.IsUnknown() && !data.Name.Equal(state.Name) {
        requestDataMap["name"] = data.Name.ValueString()
    }
    if !data.IsEnabled.IsUnknown() && !state.IsEnabled.IsUnknown() && !data.IsEnabled.Equal(state.IsEnabled) {
        requestDataMap["isEnabled"] = data.IsEnabled.ValueBool()
    }
    if !data.BurnRateThreshold.IsUnknown() && !state.BurnRateThreshold.IsUnknown() && !data.BurnRateThreshold.Equal(state.BurnRateThreshold) {
        requestDataMap["burnRateThreshold"] = r.bigFloatToFloat64(data.BurnRateThreshold.ValueBigFloat())
    }
    if !data.LongWindowInMinutes.IsUnknown() && !state.LongWindowInMinutes.IsUnknown() && !data.LongWindowInMinutes.Equal(state.LongWindowInMinutes) {
        requestDataMap["longWindowInMinutes"] = r.bigFloatToFloat64(data.LongWindowInMinutes.ValueBigFloat())
    }
    if !data.ShortWindowInMinutes.IsUnknown() && !state.ShortWindowInMinutes.IsUnknown() && !data.ShortWindowInMinutes.Equal(state.ShortWindowInMinutes) {
        requestDataMap["shortWindowInMinutes"] = r.bigFloatToFloat64(data.ShortWindowInMinutes.ValueBigFloat())
    }
    if !data.MinimumSampleCount.IsUnknown() && !state.MinimumSampleCount.IsUnknown() && !data.MinimumSampleCount.Equal(state.MinimumSampleCount) {
        requestDataMap["minimumSampleCount"] = r.bigFloatToFloat64(data.MinimumSampleCount.ValueBigFloat())
    }
    if !data.RefireSuppressionMinutes.IsUnknown() && !state.RefireSuppressionMinutes.IsUnknown() && !data.RefireSuppressionMinutes.Equal(state.RefireSuppressionMinutes) {
        requestDataMap["refireSuppressionMinutes"] = r.bigFloatToFloat64(data.RefireSuppressionMinutes.ValueBigFloat())
    }
    if !data.AlertSeverityId.IsUnknown() && !state.AlertSeverityId.IsUnknown() && !data.AlertSeverityId.Equal(state.AlertSeverityId) {
        requestDataMap["alertSeverityId"] = data.AlertSeverityId.ValueString()
    }
    if !data.OnCallDutyPolicies.IsUnknown() && !state.OnCallDutyPolicies.IsUnknown() && !data.OnCallDutyPolicies.Equal(state.OnCallDutyPolicies) {
        requestDataMap["onCallDutyPolicies"] = r.convertTerraformSetToInterface(data.OnCallDutyPolicies)
    }

    // Only call the API when there are changed fields to send. An empty
    // update body is rejected by the API; state is still refreshed below so
    // this method never writes unverified plan values into state.
    if len(sloBurnRateRuleRequest["data"].(map[string]interface{})) > 0 {
        httpResp, err := r.client.Put(ctx, "/service-level-objective-burn-rate-rule/" + data.Id.ValueString() + "", sloBurnRateRuleRequest)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update slo_burn_rate_rule, got error: %s", err))
            return
        }

        // Parse the update response
        var sloBurnRateRuleResponse map[string]interface{}
        err = r.client.ParseResponse(httpResp, &sloBurnRateRuleResponse)
        if err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to update slo_burn_rate_rule: %s", err))
            return
        }
        _ = sloBurnRateRuleResponse
    }

    // After successful update, fetch the current state by calling Read with select parameter
    selectParam := map[string]interface{}{
        "projectId": true,
        "serviceLevelObjectiveId": true,
        "name": true,
        "isEnabled": true,
        "burnRateThreshold": true,
        "longWindowInMinutes": true,
        "shortWindowInMinutes": true,
        "minimumSampleCount": true,
        "refireSuppressionMinutes": true,
        "alertSeverityId": true,
        "onCallDutyPolicies": true,
        "createdByUserId": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "lastAlertCreatedAt": true,
        "lastAlertResolvedAt": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect(ctx, "/service-level-objective-burn-rate-rule/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read slo_burn_rate_rule after update, got error: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read slo_burn_rate_rule after update: %s", err))
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
    if obj, ok := dataMap["serviceLevelObjectiveId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ServiceLevelObjectiveId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ServiceLevelObjectiveId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ServiceLevelObjectiveId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ServiceLevelObjectiveId = types.StringValue(string(jsonBytes))
            } else {
                data.ServiceLevelObjectiveId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ServiceLevelObjectiveId = types.StringValue(string(jsonBytes))
            } else {
                data.ServiceLevelObjectiveId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ServiceLevelObjectiveId = types.StringValue(string(jsonBytes))
        } else {
            data.ServiceLevelObjectiveId = types.StringNull()
        }
    } else if val, ok := dataMap["serviceLevelObjectiveId"].(string); ok {
        data.ServiceLevelObjectiveId = types.StringValue(val)
    } else {
        data.ServiceLevelObjectiveId = types.StringNull()
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
    if val, ok := dataMap["isEnabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    }
    if val, ok := dataMap["burnRateThreshold"].(float64); ok {
        data.BurnRateThreshold = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["burnRateThreshold"].(int); ok {
        data.BurnRateThreshold = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["burnRateThreshold"].(int64); ok {
        data.BurnRateThreshold = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["burnRateThreshold"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.BurnRateThreshold = types.NumberValue(big.NewFloat(val))
        } else {
            data.BurnRateThreshold = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.BurnRateThreshold = types.NumberNull()
    }
    if val, ok := dataMap["longWindowInMinutes"].(float64); ok {
        data.LongWindowInMinutes = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["longWindowInMinutes"].(int); ok {
        data.LongWindowInMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["longWindowInMinutes"].(int64); ok {
        data.LongWindowInMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["longWindowInMinutes"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.LongWindowInMinutes = types.NumberValue(big.NewFloat(val))
        } else {
            data.LongWindowInMinutes = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.LongWindowInMinutes = types.NumberNull()
    }
    if val, ok := dataMap["shortWindowInMinutes"].(float64); ok {
        data.ShortWindowInMinutes = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["shortWindowInMinutes"].(int); ok {
        data.ShortWindowInMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["shortWindowInMinutes"].(int64); ok {
        data.ShortWindowInMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["shortWindowInMinutes"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.ShortWindowInMinutes = types.NumberValue(big.NewFloat(val))
        } else {
            data.ShortWindowInMinutes = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.ShortWindowInMinutes = types.NumberNull()
    }
    if val, ok := dataMap["minimumSampleCount"].(float64); ok {
        data.MinimumSampleCount = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["minimumSampleCount"].(int); ok {
        data.MinimumSampleCount = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["minimumSampleCount"].(int64); ok {
        data.MinimumSampleCount = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["minimumSampleCount"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.MinimumSampleCount = types.NumberValue(big.NewFloat(val))
        } else {
            data.MinimumSampleCount = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.MinimumSampleCount = types.NumberNull()
    }
    if val, ok := dataMap["refireSuppressionMinutes"].(float64); ok {
        data.RefireSuppressionMinutes = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["refireSuppressionMinutes"].(int); ok {
        data.RefireSuppressionMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["refireSuppressionMinutes"].(int64); ok {
        data.RefireSuppressionMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["refireSuppressionMinutes"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.RefireSuppressionMinutes = types.NumberValue(big.NewFloat(val))
        } else {
            data.RefireSuppressionMinutes = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.RefireSuppressionMinutes = types.NumberNull()
    }
    if obj, ok := dataMap["alertSeverityId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AlertSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AlertSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AlertSeverityId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AlertSeverityId = types.StringValue(string(jsonBytes))
            } else {
                data.AlertSeverityId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AlertSeverityId = types.StringValue(string(jsonBytes))
            } else {
                data.AlertSeverityId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AlertSeverityId = types.StringValue(string(jsonBytes))
        } else {
            data.AlertSeverityId = types.StringNull()
        }
    } else if val, ok := dataMap["alertSeverityId"].(string); ok {
        data.AlertSeverityId = types.StringValue(val)
    } else {
        data.AlertSeverityId = types.StringNull()
    }
    if val, ok := dataMap["onCallDutyPolicies"].([]interface{}); ok {
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
        data.OnCallDutyPolicies = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.OnCallDutyPolicies = types.SetValueMust(types.StringType, []attr.Value{})
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
    if obj, ok := dataMap["lastAlertCreatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.LastAlertCreatedAt = NewRFC3339Value(val)
        } else {
            data.LastAlertCreatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["lastAlertCreatedAt"].(string); ok && val != "" {
        data.LastAlertCreatedAt = NewRFC3339Value(val)
    } else {
        data.LastAlertCreatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["lastAlertResolvedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.LastAlertResolvedAt = NewRFC3339Value(val)
        } else {
            data.LastAlertResolvedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["lastAlertResolvedAt"].(string); ok && val != "" {
        data.LastAlertResolvedAt = NewRFC3339Value(val)
    } else {
        data.LastAlertResolvedAt = NewRFC3339Null()
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

func (r *SloBurnRateRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var data SloBurnRateRuleResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Make API call
    httpResp, err := r.client.Delete(ctx, "/service-level-objective-burn-rate-rule/" + data.Id.ValueString() + "")
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete slo_burn_rate_rule, got error: %s", err))
        return
    }

    // A failed delete must keep the resource in state — silently dropping it
    // orphans real infrastructure. 404 means it is already gone.
    if httpResp.StatusCode >= 400 && httpResp.StatusCode != http.StatusNotFound {
        err = r.client.ParseResponse(httpResp, nil)
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to delete slo_burn_rate_rule: %s", err))
        return
    }
    if httpResp.Body != nil {
        httpResp.Body.Close()
    }
}


func (r *SloBurnRateRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Helper method to convert Terraform map to Go interface{}
func (r *SloBurnRateRuleResource) convertTerraformMapToInterface(terraformMap types.Map) interface{} {
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
func (r *SloBurnRateRuleResource) convertTerraformListToInterface(terraformList types.List) interface{} {
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
func (r *SloBurnRateRuleResource) convertTerraformSetToInterface(terraformSet types.Set) interface{} {
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
func (r *SloBurnRateRuleResource) parseJSONField(terraformString basetypes.StringValuable) interface{} {
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
func (r *SloBurnRateRuleResource) normalizeURLWrappers(value interface{}) interface{} {
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

func (r *SloBurnRateRuleResource) normalizeURLString(value string) string {
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
func (r *SloBurnRateRuleResource) bigFloatToFloat64(bf *big.Float) interface{} {
    if bf == nil {
        return nil
    }
    f, _ := bf.Float64()
    return f
}

// Helper method to check if a type string is a valid OneUptime ObjectType.
// The registry itself lives in objecttypes.go, shared across the package.
func (r *SloBurnRateRuleResource) isValidOneUptimeObjectType(typeStr string) bool {
    return validOneUptimeObjectTypes[typeStr]
}
