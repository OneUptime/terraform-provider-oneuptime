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
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &RuleExceptionResource{}
var _ resource.ResourceWithImportState = &RuleExceptionResource{}

func NewRuleExceptionResource() resource.Resource {
    return &RuleExceptionResource{}
}

// RuleExceptionResource defines the resource implementation.
type RuleExceptionResource struct {
    client *Client
}

// RuleExceptionResourceModel describes the resource data model.
type RuleExceptionResourceModel struct {
    Id types.String `tfsdk:"id"`
    ProjectId types.String `tfsdk:"project_id"`
    DetectionRuleId types.String `tfsdk:"detection_rule_id"`
    Kind types.String `tfsdk:"kind"`
    EntityType types.String `tfsdk:"entity_type"`
    FieldName types.String `tfsdk:"field_name"`
    Value types.String `tfsdk:"value"`
    StartsAt RFC3339Value `tfsdk:"starts_at"`
    EndsAt RFC3339Value `tfsdk:"ends_at"`
    Reason types.String `tfsdk:"reason"`
    ExpiresAt RFC3339Value `tfsdk:"expires_at"`
    CreatedAt RFC3339Value `tfsdk:"created_at"`
    UpdatedAt RFC3339Value `tfsdk:"updated_at"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (r *RuleExceptionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_rule_exception"
}

func (r *RuleExceptionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = r.schemaDefinition()
}

func (r *RuleExceptionResource) schemaDefinition() schema.Schema {
    return schema.Schema{
        MarkdownDescription: "What a detection rule skips: an entity, a value or a time window, with who added it, why, and when it expires.",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Unique identifier for the resource.",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "ID of the project this exception belongs to. The ID of a `oneuptime_project`.",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "detection_rule_id": schema.StringAttribute{
                MarkdownDescription: "ID of the detection rule that skips what this names. The ID of a `oneuptime_detection_rule`.",
                Required: true,
            },
            "kind": schema.StringAttribute{
                MarkdownDescription: "What the rule skips: Entity (a user, host or IP address), Value (one event field's value) or Time Window (every event between two times).",
                Required: true,
            },
            "entity_type": schema.StringAttribute{
                MarkdownDescription: "For an Entity: User, Host or IP Address. The rule skips every event that names it, as who acted, as who was acted on, or among its observables.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "field_name": schema.StringAttribute{
                MarkdownDescription: "For a Value: the event field to match, named as in a Sigma rule - principalUser, processName, an attribute key such as http.user_agent.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "value": schema.StringAttribute{
                MarkdownDescription: "For an Entity or a Value: what the rule skips, matched whatever the case. * stands for any run of characters (svc-* skips every service account), and an IP address may be a range such as 10.20.0.0/16.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "starts_at": schema.StringAttribute{
                MarkdownDescription: "For a Time Window: the first moment the rule skips.",
                CustomType: RFC3339Type{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "ends_at": schema.StringAttribute{
                MarkdownDescription: "For a Time Window: when the rule sees events again. A window is at most 30 days.",
                CustomType: RFC3339Type{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "reason": schema.StringAttribute{
                MarkdownDescription: "Why it is safe for the rule to skip this, for whoever reads the list later.",
                Required: true,
            },
            "expires_at": schema.StringAttribute{
                MarkdownDescription: "For an Entity or a Value: when the rule sees it again. Empty: until someone deletes the exception.",
                CustomType: RFC3339Type{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "created_at": schema.StringAttribute{
                MarkdownDescription: "Date and Time when the object was created.",
                CustomType: RFC3339Type{},
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "updated_at": schema.StringAttribute{
                MarkdownDescription: "Date and Time when the object was updated.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the person who added the exception. The ID of a `oneuptime_user` (see the data source).",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
        },
    }
}

func (r *RuleExceptionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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


func (r *RuleExceptionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var data RuleExceptionResourceModel

    // Read Terraform plan data into the model
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // What the configuration sets, and what Terraform planned (see keepPlannedValues).
    var config RuleExceptionResourceModel
    resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
    if resp.Diagnostics.HasError() {
        return
    }
    plan := data



    // Create API request body. Unset (null/unknown) optional fields are
    // omitted so server-side defaults apply instead of being overwritten
    // with zero values.
    ruleExceptionRequest := map[string]interface{}{
        "data": map[string]interface{}{},
    }
    requestDataMap := ruleExceptionRequest["data"].(map[string]interface{})

    if !data.DetectionRuleId.IsNull() && !data.DetectionRuleId.IsUnknown() {
        requestDataMap["detectionRuleId"] = data.DetectionRuleId.ValueString()
    }
    if !data.Kind.IsNull() && !data.Kind.IsUnknown() {
        requestDataMap["kind"] = data.Kind.ValueString()
    }
    if !data.EntityType.IsNull() && !data.EntityType.IsUnknown() {
        requestDataMap["entityType"] = data.EntityType.ValueString()
    }
    if !data.FieldName.IsNull() && !data.FieldName.IsUnknown() {
        requestDataMap["fieldName"] = data.FieldName.ValueString()
    }
    if !data.Value.IsNull() && !data.Value.IsUnknown() {
        requestDataMap["value"] = data.Value.ValueString()
    }
    if !data.StartsAt.IsNull() && !data.StartsAt.IsUnknown() {
        requestDataMap["startsAt"] = data.StartsAt.ValueString()
    }
    if !data.EndsAt.IsNull() && !data.EndsAt.IsUnknown() {
        requestDataMap["endsAt"] = data.EndsAt.ValueString()
    }
    if !data.Reason.IsNull() && !data.Reason.IsUnknown() {
        requestDataMap["reason"] = data.Reason.ValueString()
    }
    if !data.ExpiresAt.IsNull() && !data.ExpiresAt.IsUnknown() {
        requestDataMap["expiresAt"] = data.ExpiresAt.ValueString()
    }

    // Make API call
    httpResp, err := r.client.Post(ctx, "/detection-rule-exception", ruleExceptionRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create rule_exception, got error: %s", err))
        return
    }

    var ruleExceptionResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &ruleExceptionResponse)
    if err != nil {
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to create rule_exception: %s", err))
        return
    }

    // Extract the new resource id from the create response.
    createdId := ""
    if wrapper, ok := ruleExceptionResponse["data"].(map[string]interface{}); ok {
        if val, ok := wrapper["_id"].(string); ok {
            createdId = val
        }
    } else if val, ok := ruleExceptionResponse["_id"].(string); ok {
        createdId = val
    }
    if createdId == "" {
        resp.Diagnostics.AddError("OneUptime API Error", "Create response for rule_exception did not contain an id. This is a bug in the provider or the API; please report it.")
        return
    }
    data.Id = types.StringValue(createdId)

    /*
     * The server has committed the row. Persist what we know to state BEFORE
     * the read-back: if the read-back fails and we return without setting
     * state, Terraform never learns the resource exists and the created
     * rule_exception is orphaned server-side — never refreshed, never
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
        "detectionRuleId": true,
        "kind": true,
        "entityType": true,
        "fieldName": true,
        "value": true,
        "startsAt": true,
        "endsAt": true,
        "reason": true,
        "expiresAt": true,
        "createdAt": true,
        "updatedAt": true,
        "createdByUserId": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect(ctx, "/detection-rule-exception/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        /*
         * State already owns the id, so the resource is tracked and the next
         * refresh reconciles the remaining attributes. Warn rather than
         * error: erroring here would strand a real resource.
         */
        resp.Diagnostics.AddWarning("Read After Create Failed", fmt.Sprintf("Created rule_exception but could not read it back; state is incomplete until the next refresh: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddWarning("Read After Create Failed", fmt.Sprintf("Created rule_exception but could not parse the read-back response; state is incomplete until the next refresh: %s", err))
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
    if obj, ok := dataMap["detectionRuleId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DetectionRuleId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DetectionRuleId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DetectionRuleId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DetectionRuleId = types.StringValue(string(jsonBytes))
            } else {
                data.DetectionRuleId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DetectionRuleId = types.StringValue(string(jsonBytes))
            } else {
                data.DetectionRuleId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DetectionRuleId = types.StringValue(string(jsonBytes))
        } else {
            data.DetectionRuleId = types.StringNull()
        }
    } else if val, ok := dataMap["detectionRuleId"].(string); ok {
        data.DetectionRuleId = types.StringValue(val)
    } else {
        data.DetectionRuleId = types.StringNull()
    }
    if obj, ok := dataMap["kind"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Kind = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Kind = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Kind = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Kind = types.StringValue(string(jsonBytes))
            } else {
                data.Kind = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Kind = types.StringValue(string(jsonBytes))
            } else {
                data.Kind = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Kind = types.StringValue(string(jsonBytes))
        } else {
            data.Kind = types.StringNull()
        }
    } else if val, ok := dataMap["kind"].(string); ok {
        data.Kind = types.StringValue(val)
    } else {
        data.Kind = types.StringNull()
    }
    if obj, ok := dataMap["entityType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EntityType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.EntityType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.EntityType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.EntityType = types.StringValue(string(jsonBytes))
            } else {
                data.EntityType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.EntityType = types.StringValue(string(jsonBytes))
            } else {
                data.EntityType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.EntityType = types.StringValue(string(jsonBytes))
        } else {
            data.EntityType = types.StringNull()
        }
    } else if val, ok := dataMap["entityType"].(string); ok {
        data.EntityType = types.StringValue(val)
    } else {
        data.EntityType = types.StringNull()
    }
    if obj, ok := dataMap["fieldName"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FieldName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.FieldName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.FieldName = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.FieldName = types.StringValue(string(jsonBytes))
            } else {
                data.FieldName = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.FieldName = types.StringValue(string(jsonBytes))
            } else {
                data.FieldName = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.FieldName = types.StringValue(string(jsonBytes))
        } else {
            data.FieldName = types.StringNull()
        }
    } else if val, ok := dataMap["fieldName"].(string); ok {
        data.FieldName = types.StringValue(val)
    } else {
        data.FieldName = types.StringNull()
    }
    if obj, ok := dataMap["value"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Value = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Value = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Value = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Value = types.StringValue(string(jsonBytes))
            } else {
                data.Value = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Value = types.StringValue(string(jsonBytes))
            } else {
                data.Value = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Value = types.StringValue(string(jsonBytes))
        } else {
            data.Value = types.StringNull()
        }
    } else if val, ok := dataMap["value"].(string); ok {
        data.Value = types.StringValue(val)
    } else {
        data.Value = types.StringNull()
    }
    if obj, ok := dataMap["startsAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.StartsAt = NewRFC3339Value(val)
        } else {
            data.StartsAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["startsAt"].(string); ok && val != "" {
        data.StartsAt = NewRFC3339Value(val)
    } else {
        data.StartsAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["endsAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.EndsAt = NewRFC3339Value(val)
        } else {
            data.EndsAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["endsAt"].(string); ok && val != "" {
        data.EndsAt = NewRFC3339Value(val)
    } else {
        data.EndsAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["reason"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Reason = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Reason = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Reason = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Reason = types.StringValue(string(jsonBytes))
            } else {
                data.Reason = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Reason = types.StringValue(string(jsonBytes))
            } else {
                data.Reason = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Reason = types.StringValue(string(jsonBytes))
        } else {
            data.Reason = types.StringNull()
        }
    } else if val, ok := dataMap["reason"].(string); ok {
        data.Reason = types.StringValue(val)
    } else {
        data.Reason = types.StringNull()
    }
    if obj, ok := dataMap["expiresAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.ExpiresAt = NewRFC3339Value(val)
        } else {
            data.ExpiresAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["expiresAt"].(string); ok && val != "" {
        data.ExpiresAt = NewRFC3339Value(val)
    } else {
        data.ExpiresAt = NewRFC3339Null()
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
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    // The read response is authoritative, but never let it clobber the id we just received.
    data.Id = types.StringValue(createdId)

    // Unconfigured attributes keep their planned value; see keepPlannedValues.
    r.keepPlannedValues(&data, &plan, &config)

    // Write logs using the tflog package
    tflog.Trace(ctx, "created a resource")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RuleExceptionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var data RuleExceptionResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create select parameter to get full object
    selectParam := map[string]interface{}{
        "projectId": true,
        "detectionRuleId": true,
        "kind": true,
        "entityType": true,
        "fieldName": true,
        "value": true,
        "startsAt": true,
        "endsAt": true,
        "reason": true,
        "expiresAt": true,
        "createdAt": true,
        "updatedAt": true,
        "createdByUserId": true,
        "_id": true,
    }

    // Make API call with select parameter
    httpResp, err := r.client.PostWithSelect(ctx, "/detection-rule-exception/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read rule_exception, got error: %s", err))
        return
    }

    if httpResp.StatusCode == http.StatusNotFound {
        resp.State.RemoveResource(ctx)
        return
    }

    var ruleExceptionResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &ruleExceptionResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse rule_exception response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := ruleExceptionResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = ruleExceptionResponse
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
    if obj, ok := dataMap["detectionRuleId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DetectionRuleId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DetectionRuleId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DetectionRuleId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DetectionRuleId = types.StringValue(string(jsonBytes))
            } else {
                data.DetectionRuleId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DetectionRuleId = types.StringValue(string(jsonBytes))
            } else {
                data.DetectionRuleId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DetectionRuleId = types.StringValue(string(jsonBytes))
        } else {
            data.DetectionRuleId = types.StringNull()
        }
    } else if val, ok := dataMap["detectionRuleId"].(string); ok {
        data.DetectionRuleId = types.StringValue(val)
    } else {
        data.DetectionRuleId = types.StringNull()
    }
    if obj, ok := dataMap["kind"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Kind = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Kind = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Kind = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Kind = types.StringValue(string(jsonBytes))
            } else {
                data.Kind = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Kind = types.StringValue(string(jsonBytes))
            } else {
                data.Kind = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Kind = types.StringValue(string(jsonBytes))
        } else {
            data.Kind = types.StringNull()
        }
    } else if val, ok := dataMap["kind"].(string); ok {
        data.Kind = types.StringValue(val)
    } else {
        data.Kind = types.StringNull()
    }
    if obj, ok := dataMap["entityType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EntityType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.EntityType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.EntityType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.EntityType = types.StringValue(string(jsonBytes))
            } else {
                data.EntityType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.EntityType = types.StringValue(string(jsonBytes))
            } else {
                data.EntityType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.EntityType = types.StringValue(string(jsonBytes))
        } else {
            data.EntityType = types.StringNull()
        }
    } else if val, ok := dataMap["entityType"].(string); ok {
        data.EntityType = types.StringValue(val)
    } else {
        data.EntityType = types.StringNull()
    }
    if obj, ok := dataMap["fieldName"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FieldName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.FieldName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.FieldName = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.FieldName = types.StringValue(string(jsonBytes))
            } else {
                data.FieldName = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.FieldName = types.StringValue(string(jsonBytes))
            } else {
                data.FieldName = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.FieldName = types.StringValue(string(jsonBytes))
        } else {
            data.FieldName = types.StringNull()
        }
    } else if val, ok := dataMap["fieldName"].(string); ok {
        data.FieldName = types.StringValue(val)
    } else {
        data.FieldName = types.StringNull()
    }
    if obj, ok := dataMap["value"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Value = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Value = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Value = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Value = types.StringValue(string(jsonBytes))
            } else {
                data.Value = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Value = types.StringValue(string(jsonBytes))
            } else {
                data.Value = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Value = types.StringValue(string(jsonBytes))
        } else {
            data.Value = types.StringNull()
        }
    } else if val, ok := dataMap["value"].(string); ok {
        data.Value = types.StringValue(val)
    } else {
        data.Value = types.StringNull()
    }
    if obj, ok := dataMap["startsAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.StartsAt = NewRFC3339Value(val)
        } else {
            data.StartsAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["startsAt"].(string); ok && val != "" {
        data.StartsAt = NewRFC3339Value(val)
    } else {
        data.StartsAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["endsAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.EndsAt = NewRFC3339Value(val)
        } else {
            data.EndsAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["endsAt"].(string); ok && val != "" {
        data.EndsAt = NewRFC3339Value(val)
    } else {
        data.EndsAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["reason"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Reason = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Reason = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Reason = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Reason = types.StringValue(string(jsonBytes))
            } else {
                data.Reason = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Reason = types.StringValue(string(jsonBytes))
            } else {
                data.Reason = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Reason = types.StringValue(string(jsonBytes))
        } else {
            data.Reason = types.StringNull()
        }
    } else if val, ok := dataMap["reason"].(string); ok {
        data.Reason = types.StringValue(val)
    } else {
        data.Reason = types.StringNull()
    }
    if obj, ok := dataMap["expiresAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.ExpiresAt = NewRFC3339Value(val)
        } else {
            data.ExpiresAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["expiresAt"].(string); ok && val != "" {
        data.ExpiresAt = NewRFC3339Value(val)
    } else {
        data.ExpiresAt = NewRFC3339Null()
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
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RuleExceptionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var data RuleExceptionResourceModel
    var state RuleExceptionResourceModel

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

    // What the configuration sets, and what Terraform planned (see keepPlannedValues).
    var config RuleExceptionResourceModel
    resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
    if resp.Diagnostics.HasError() {
        return
    }
    plan := data

    // Create API request body
    ruleExceptionRequest := map[string]interface{}{
        "data": map[string]interface{}{},
    }
    requestDataMap := ruleExceptionRequest["data"].(map[string]interface{})

    if !data.DetectionRuleId.IsUnknown() && !state.DetectionRuleId.IsUnknown() && !data.DetectionRuleId.Equal(state.DetectionRuleId) {
        requestDataMap["detectionRuleId"] = data.DetectionRuleId.ValueString()
    }
    if !data.Kind.IsUnknown() && !state.Kind.IsUnknown() && !data.Kind.Equal(state.Kind) {
        requestDataMap["kind"] = data.Kind.ValueString()
    }
    if !data.EntityType.IsUnknown() && !state.EntityType.IsUnknown() && !data.EntityType.Equal(state.EntityType) {
        requestDataMap["entityType"] = data.EntityType.ValueString()
    }
    if !data.FieldName.IsUnknown() && !state.FieldName.IsUnknown() && !data.FieldName.Equal(state.FieldName) {
        requestDataMap["fieldName"] = data.FieldName.ValueString()
    }
    if !data.Value.IsUnknown() && !state.Value.IsUnknown() && !data.Value.Equal(state.Value) {
        requestDataMap["value"] = data.Value.ValueString()
    }
    if !data.StartsAt.IsUnknown() && !state.StartsAt.IsUnknown() && !data.StartsAt.Equal(state.StartsAt) {
        requestDataMap["startsAt"] = data.StartsAt.ValueString()
    }
    if !data.EndsAt.IsUnknown() && !state.EndsAt.IsUnknown() && !data.EndsAt.Equal(state.EndsAt) {
        requestDataMap["endsAt"] = data.EndsAt.ValueString()
    }
    if !data.Reason.IsUnknown() && !state.Reason.IsUnknown() && !data.Reason.Equal(state.Reason) {
        requestDataMap["reason"] = data.Reason.ValueString()
    }
    if !data.ExpiresAt.IsUnknown() && !state.ExpiresAt.IsUnknown() && !data.ExpiresAt.Equal(state.ExpiresAt) {
        requestDataMap["expiresAt"] = data.ExpiresAt.ValueString()
    }

    // Only call the API when there are changed fields to send. An empty
    // update body is rejected by the API; state is still refreshed below so
    // this method never writes unverified plan values into state.
    if len(ruleExceptionRequest["data"].(map[string]interface{})) > 0 {
        httpResp, err := r.client.Put(ctx, "/detection-rule-exception/" + data.Id.ValueString() + "", ruleExceptionRequest)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update rule_exception, got error: %s", err))
            return
        }

        // Parse the update response
        var ruleExceptionResponse map[string]interface{}
        err = r.client.ParseResponse(httpResp, &ruleExceptionResponse)
        if err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to update rule_exception: %s", err))
            return
        }
        _ = ruleExceptionResponse
    }

    // After successful update, fetch the current state by calling Read with select parameter
    selectParam := map[string]interface{}{
        "projectId": true,
        "detectionRuleId": true,
        "kind": true,
        "entityType": true,
        "fieldName": true,
        "value": true,
        "startsAt": true,
        "endsAt": true,
        "reason": true,
        "expiresAt": true,
        "createdAt": true,
        "updatedAt": true,
        "createdByUserId": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect(ctx, "/detection-rule-exception/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read rule_exception after update, got error: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read rule_exception after update: %s", err))
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
    if obj, ok := dataMap["detectionRuleId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DetectionRuleId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DetectionRuleId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DetectionRuleId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DetectionRuleId = types.StringValue(string(jsonBytes))
            } else {
                data.DetectionRuleId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DetectionRuleId = types.StringValue(string(jsonBytes))
            } else {
                data.DetectionRuleId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DetectionRuleId = types.StringValue(string(jsonBytes))
        } else {
            data.DetectionRuleId = types.StringNull()
        }
    } else if val, ok := dataMap["detectionRuleId"].(string); ok {
        data.DetectionRuleId = types.StringValue(val)
    } else {
        data.DetectionRuleId = types.StringNull()
    }
    if obj, ok := dataMap["kind"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Kind = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Kind = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Kind = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Kind = types.StringValue(string(jsonBytes))
            } else {
                data.Kind = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Kind = types.StringValue(string(jsonBytes))
            } else {
                data.Kind = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Kind = types.StringValue(string(jsonBytes))
        } else {
            data.Kind = types.StringNull()
        }
    } else if val, ok := dataMap["kind"].(string); ok {
        data.Kind = types.StringValue(val)
    } else {
        data.Kind = types.StringNull()
    }
    if obj, ok := dataMap["entityType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EntityType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.EntityType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.EntityType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.EntityType = types.StringValue(string(jsonBytes))
            } else {
                data.EntityType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.EntityType = types.StringValue(string(jsonBytes))
            } else {
                data.EntityType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.EntityType = types.StringValue(string(jsonBytes))
        } else {
            data.EntityType = types.StringNull()
        }
    } else if val, ok := dataMap["entityType"].(string); ok {
        data.EntityType = types.StringValue(val)
    } else {
        data.EntityType = types.StringNull()
    }
    if obj, ok := dataMap["fieldName"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FieldName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.FieldName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.FieldName = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.FieldName = types.StringValue(string(jsonBytes))
            } else {
                data.FieldName = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.FieldName = types.StringValue(string(jsonBytes))
            } else {
                data.FieldName = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.FieldName = types.StringValue(string(jsonBytes))
        } else {
            data.FieldName = types.StringNull()
        }
    } else if val, ok := dataMap["fieldName"].(string); ok {
        data.FieldName = types.StringValue(val)
    } else {
        data.FieldName = types.StringNull()
    }
    if obj, ok := dataMap["value"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Value = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Value = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Value = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Value = types.StringValue(string(jsonBytes))
            } else {
                data.Value = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Value = types.StringValue(string(jsonBytes))
            } else {
                data.Value = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Value = types.StringValue(string(jsonBytes))
        } else {
            data.Value = types.StringNull()
        }
    } else if val, ok := dataMap["value"].(string); ok {
        data.Value = types.StringValue(val)
    } else {
        data.Value = types.StringNull()
    }
    if obj, ok := dataMap["startsAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.StartsAt = NewRFC3339Value(val)
        } else {
            data.StartsAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["startsAt"].(string); ok && val != "" {
        data.StartsAt = NewRFC3339Value(val)
    } else {
        data.StartsAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["endsAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.EndsAt = NewRFC3339Value(val)
        } else {
            data.EndsAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["endsAt"].(string); ok && val != "" {
        data.EndsAt = NewRFC3339Value(val)
    } else {
        data.EndsAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["reason"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Reason = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Reason = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Reason = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Reason = types.StringValue(string(jsonBytes))
            } else {
                data.Reason = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Reason = types.StringValue(string(jsonBytes))
            } else {
                data.Reason = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Reason = types.StringValue(string(jsonBytes))
        } else {
            data.Reason = types.StringNull()
        }
    } else if val, ok := dataMap["reason"].(string); ok {
        data.Reason = types.StringValue(val)
    } else {
        data.Reason = types.StringNull()
    }
    if obj, ok := dataMap["expiresAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.ExpiresAt = NewRFC3339Value(val)
        } else {
            data.ExpiresAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["expiresAt"].(string); ok && val != "" {
        data.ExpiresAt = NewRFC3339Value(val)
    } else {
        data.ExpiresAt = NewRFC3339Null()
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
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    data.Id = state.Id

    // Unconfigured attributes keep their planned value; see keepPlannedValues.
    r.keepPlannedValues(&data, &plan, &config)

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RuleExceptionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var data RuleExceptionResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Make API call
    httpResp, err := r.client.Delete(ctx, "/detection-rule-exception/" + data.Id.ValueString() + "")
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete rule_exception, got error: %s", err))
        return
    }

    // A failed delete must keep the resource in state — silently dropping it
    // orphans real infrastructure. 404 means it is already gone.
    if httpResp.StatusCode >= 400 && httpResp.StatusCode != http.StatusNotFound {
        err = r.client.ParseResponse(httpResp, nil)
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to delete rule_exception: %s", err))
        return
    }
    if httpResp.Body != nil {
        httpResp.Body.Close()
    }
}


func (r *RuleExceptionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}


// keepPlannedValues puts back, after a create or an update, the planned value
// of each optional attribute the configuration leaves out. The server keeps
// some of these up to date on its own (when it last checked a heartbeat, the
// status a probe last reported...), so the value read back after the write can
// already differ from the plan, and Terraform would fail the apply with
// "Provider produced inconsistent result after apply". The next refresh reads
// the server's value, which is never a diff for an attribute nobody configured.
func (r *RuleExceptionResource) keepPlannedValues(data *RuleExceptionResourceModel, plan *RuleExceptionResourceModel, config *RuleExceptionResourceModel) {
    if config.EntityType.IsNull() && !plan.EntityType.IsUnknown() {
        data.EntityType = plan.EntityType
    }
    if config.FieldName.IsNull() && !plan.FieldName.IsUnknown() {
        data.FieldName = plan.FieldName
    }
    if config.Value.IsNull() && !plan.Value.IsUnknown() {
        data.Value = plan.Value
    }
    if config.StartsAt.IsNull() && !plan.StartsAt.IsUnknown() {
        data.StartsAt = plan.StartsAt
    }
    if config.EndsAt.IsNull() && !plan.EndsAt.IsUnknown() {
        data.EndsAt = plan.EndsAt
    }
    if config.ExpiresAt.IsNull() && !plan.ExpiresAt.IsUnknown() {
        data.ExpiresAt = plan.ExpiresAt
    }
}

// Helper method to convert Terraform map to Go interface{}
func (r *RuleExceptionResource) convertTerraformMapToInterface(terraformMap types.Map) interface{} {
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
func (r *RuleExceptionResource) convertTerraformListToInterface(terraformList types.List) interface{} {
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
func (r *RuleExceptionResource) convertTerraformSetToInterface(terraformSet types.Set) interface{} {
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
func (r *RuleExceptionResource) parseJSONField(terraformString basetypes.StringValuable) interface{} {
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
func (r *RuleExceptionResource) normalizeURLWrappers(value interface{}) interface{} {
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

func (r *RuleExceptionResource) normalizeURLString(value string) string {
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
func (r *RuleExceptionResource) bigFloatToFloat64(bf *big.Float) interface{} {
    if bf == nil {
        return nil
    }
    f, _ := bf.Float64()
    return f
}

// Helper method to check if a type string is a valid OneUptime ObjectType.
// The registry itself lives in objecttypes.go, shared across the package.
func (r *RuleExceptionResource) isValidOneUptimeObjectType(typeStr string) bool {
    return validOneUptimeObjectTypes[typeStr]
}
