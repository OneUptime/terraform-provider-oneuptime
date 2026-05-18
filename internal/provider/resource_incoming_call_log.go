package provider

import (
    "context"
    "fmt"
    "github.com/hashicorp/terraform-plugin-framework/path"
    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-framework/types/basetypes"
    "github.com/hashicorp/terraform-plugin-log/tflog"
    "math/big"
    "net/http"
    "encoding/json"
    "net/url"
    "strings"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &IncomingCallLogResource{}
var _ resource.ResourceWithImportState = &IncomingCallLogResource{}

func NewIncomingCallLogResource() resource.Resource {
    return &IncomingCallLogResource{}
}

// IncomingCallLogResource defines the resource implementation.
type IncomingCallLogResource struct {
    client *Client
}

// IncomingCallLogResourceModel describes the resource data model.
type IncomingCallLogResourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt JSONSubsetValue `tfsdk:"created_at"`
    UpdatedAt JSONSubsetValue `tfsdk:"updated_at"`
    DeletedAt JSONSubsetValue `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    IncomingCallPolicyId types.String `tfsdk:"incoming_call_policy_id"`
    CallerPhoneNumber JSONSubsetValue `tfsdk:"caller_phone_number"`
    RoutingPhoneNumber JSONSubsetValue `tfsdk:"routing_phone_number"`
    CallProviderCallId types.String `tfsdk:"call_provider_call_id"`
    Status types.String `tfsdk:"status"`
    StatusMessage types.String `tfsdk:"status_message"`
    CallDurationInSeconds types.Number `tfsdk:"call_duration_in_seconds"`
    CallCostInUsdCents types.Number `tfsdk:"call_cost_in_usd_cents"`
    IncomingCallCostInUsdCents types.Number `tfsdk:"incoming_call_cost_in_usd_cents"`
    OutgoingCallCostInUsdCents types.Number `tfsdk:"outgoing_call_cost_in_usd_cents"`
    StartedAt JSONSubsetValue `tfsdk:"started_at"`
    EndedAt JSONSubsetValue `tfsdk:"ended_at"`
    AnsweredByUserId types.String `tfsdk:"answered_by_user_id"`
    CurrentEscalationRuleOrder types.Number `tfsdk:"current_escalation_rule_order"`
    RepeatCount types.Number `tfsdk:"repeat_count"`
}

func (r *IncomingCallLogResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_incoming_call_log"
}

func (r *IncomingCallLogResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "incoming_call_log resource",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Unique identifier for the resource",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "created_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                CustomType: JSONSubsetType{},
                Computed: true,
            },
            "updated_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                CustomType: JSONSubsetType{},
                Computed: true,
            },
            "deleted_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                CustomType: JSONSubsetType{},
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
            "incoming_call_policy_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "caller_phone_number": schema.StringAttribute{
                MarkdownDescription: "Phone object",
                CustomType: JSONSubsetType{},
                Computed: true,
            },
            "routing_phone_number": schema.StringAttribute{
                MarkdownDescription: "Phone object",
                CustomType: JSONSubsetType{},
                Computed: true,
            },
            "call_provider_call_id": schema.StringAttribute{
                MarkdownDescription: "Call provider's call identifier. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Incoming Call Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "status": schema.StringAttribute{
                MarkdownDescription: "Current status of the incoming call. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Incoming Call Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "status_message": schema.StringAttribute{
                MarkdownDescription: "Additional status information. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Incoming Call Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "call_duration_in_seconds": schema.NumberAttribute{
                MarkdownDescription: "Total call duration in seconds. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Incoming Call Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "call_cost_in_usd_cents": schema.NumberAttribute{
                MarkdownDescription: "Total cost for this call in USD cents. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Incoming Call Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "incoming_call_cost_in_usd_cents": schema.NumberAttribute{
                MarkdownDescription: "Cost for incoming leg in USD cents. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Incoming Call Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "outgoing_call_cost_in_usd_cents": schema.NumberAttribute{
                MarkdownDescription: "Cost for all forwarding attempts in USD cents. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Incoming Call Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "started_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                CustomType: JSONSubsetType{},
                Computed: true,
            },
            "ended_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                CustomType: JSONSubsetType{},
                Computed: true,
            },
            "answered_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "current_escalation_rule_order": schema.NumberAttribute{
                MarkdownDescription: "The current escalation rule order being processed. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Incoming Call Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "repeat_count": schema.NumberAttribute{
                MarkdownDescription: "Number of times the policy has been repeated. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Incoming Call Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
        },
    }
}

func (r *IncomingCallLogResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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


func (r *IncomingCallLogResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var data IncomingCallLogResourceModel

    // Read Terraform plan data into the model
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }



    // Create API request body
    incomingCallLogRequest := map[string]interface{}{
        "data": map[string]interface{}{

        },
    }

    // Make API call
    httpResp, err := r.client.Post("/incoming-call-log/count", incomingCallLogRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create incoming_call_log, got error: %s", err))
        return
    }

    var incomingCallLogResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &incomingCallLogResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse incoming_call_log response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := incomingCallLogResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = incomingCallLogResponse
    }

    if obj, ok := dataMap["id"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Id = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Id = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Id = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Id = types.StringValue(string(jsonBytes))
            } else {
                data.Id = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Id = types.StringValue(string(jsonBytes))
            } else {
                data.Id = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Id = types.StringValue(string(jsonBytes))
        } else {
            data.Id = types.StringNull()
        }
    } else if val, ok := dataMap["id"].(string); ok && val != "" {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    if obj, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedAt = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedAt = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedAt = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CreatedAt = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.CreatedAt = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CreatedAt = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.CreatedAt = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedAt = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.CreatedAt = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["createdAt"].(string); ok && val != "" {
        data.CreatedAt = NewJSONSubsetValue(val)
    } else {
        data.CreatedAt = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UpdatedAt = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.UpdatedAt = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.UpdatedAt = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.UpdatedAt = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.UpdatedAt = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.UpdatedAt = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.UpdatedAt = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.UpdatedAt = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.UpdatedAt = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["updatedAt"].(string); ok && val != "" {
        data.UpdatedAt = NewJSONSubsetValue(val)
    } else {
        data.UpdatedAt = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["deletedAt"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeletedAt = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DeletedAt = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DeletedAt = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DeletedAt = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DeletedAt = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DeletedAt = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DeletedAt = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DeletedAt = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DeletedAt = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["deletedAt"].(string); ok && val != "" {
        data.DeletedAt = NewJSONSubsetValue(val)
    } else {
        data.DeletedAt = NewJSONSubsetNull()
    }
    if val, ok := dataMap["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["version"].(int); ok {
        data.Version = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["version"].(int64); ok {
        data.Version = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["version"] == nil {
        data.Version = types.NumberNull()
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
    if obj, ok := dataMap["incomingCallPolicyId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncomingCallPolicyId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.IncomingCallPolicyId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.IncomingCallPolicyId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.IncomingCallPolicyId = types.StringValue(string(jsonBytes))
            } else {
                data.IncomingCallPolicyId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.IncomingCallPolicyId = types.StringValue(string(jsonBytes))
            } else {
                data.IncomingCallPolicyId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.IncomingCallPolicyId = types.StringValue(string(jsonBytes))
        } else {
            data.IncomingCallPolicyId = types.StringNull()
        }
    } else if val, ok := dataMap["incomingCallPolicyId"].(string); ok && val != "" {
        data.IncomingCallPolicyId = types.StringValue(val)
    } else {
        data.IncomingCallPolicyId = types.StringNull()
    }
    if obj, ok := dataMap["callerPhoneNumber"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CallerPhoneNumber = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CallerPhoneNumber = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CallerPhoneNumber = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CallerPhoneNumber = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.CallerPhoneNumber = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CallerPhoneNumber = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.CallerPhoneNumber = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CallerPhoneNumber = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.CallerPhoneNumber = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["callerPhoneNumber"].(string); ok && val != "" {
        data.CallerPhoneNumber = NewJSONSubsetValue(val)
    } else {
        data.CallerPhoneNumber = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["routingPhoneNumber"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RoutingPhoneNumber = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.RoutingPhoneNumber = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.RoutingPhoneNumber = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.RoutingPhoneNumber = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.RoutingPhoneNumber = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.RoutingPhoneNumber = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.RoutingPhoneNumber = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.RoutingPhoneNumber = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.RoutingPhoneNumber = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["routingPhoneNumber"].(string); ok && val != "" {
        data.RoutingPhoneNumber = NewJSONSubsetValue(val)
    } else {
        data.RoutingPhoneNumber = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["callProviderCallId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CallProviderCallId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CallProviderCallId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CallProviderCallId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CallProviderCallId = types.StringValue(string(jsonBytes))
            } else {
                data.CallProviderCallId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CallProviderCallId = types.StringValue(string(jsonBytes))
            } else {
                data.CallProviderCallId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CallProviderCallId = types.StringValue(string(jsonBytes))
        } else {
            data.CallProviderCallId = types.StringNull()
        }
    } else if val, ok := dataMap["callProviderCallId"].(string); ok && val != "" {
        data.CallProviderCallId = types.StringValue(val)
    } else {
        data.CallProviderCallId = types.StringNull()
    }
    if obj, ok := dataMap["status"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Status = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Status = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Status = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Status = types.StringValue(string(jsonBytes))
            } else {
                data.Status = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Status = types.StringValue(string(jsonBytes))
            } else {
                data.Status = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Status = types.StringValue(string(jsonBytes))
        } else {
            data.Status = types.StringNull()
        }
    } else if val, ok := dataMap["status"].(string); ok && val != "" {
        data.Status = types.StringValue(val)
    } else {
        data.Status = types.StringNull()
    }
    if obj, ok := dataMap["statusMessage"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StatusMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.StatusMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.StatusMessage = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.StatusMessage = types.StringValue(string(jsonBytes))
            } else {
                data.StatusMessage = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.StatusMessage = types.StringValue(string(jsonBytes))
            } else {
                data.StatusMessage = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.StatusMessage = types.StringValue(string(jsonBytes))
        } else {
            data.StatusMessage = types.StringNull()
        }
    } else if val, ok := dataMap["statusMessage"].(string); ok && val != "" {
        data.StatusMessage = types.StringValue(val)
    } else {
        data.StatusMessage = types.StringNull()
    }
    if val, ok := dataMap["callDurationInSeconds"].(float64); ok {
        data.CallDurationInSeconds = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["callDurationInSeconds"].(int); ok {
        data.CallDurationInSeconds = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["callDurationInSeconds"].(int64); ok {
        data.CallDurationInSeconds = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["callDurationInSeconds"] == nil {
        data.CallDurationInSeconds = types.NumberNull()
    }
    if val, ok := dataMap["callCostInUSDCents"].(float64); ok {
        data.CallCostInUsdCents = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["callCostInUSDCents"].(int); ok {
        data.CallCostInUsdCents = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["callCostInUSDCents"].(int64); ok {
        data.CallCostInUsdCents = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["callCostInUSDCents"] == nil {
        data.CallCostInUsdCents = types.NumberNull()
    }
    if val, ok := dataMap["incomingCallCostInUSDCents"].(float64); ok {
        data.IncomingCallCostInUsdCents = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["incomingCallCostInUSDCents"].(int); ok {
        data.IncomingCallCostInUsdCents = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["incomingCallCostInUSDCents"].(int64); ok {
        data.IncomingCallCostInUsdCents = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["incomingCallCostInUSDCents"] == nil {
        data.IncomingCallCostInUsdCents = types.NumberNull()
    }
    if val, ok := dataMap["outgoingCallCostInUSDCents"].(float64); ok {
        data.OutgoingCallCostInUsdCents = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["outgoingCallCostInUSDCents"].(int); ok {
        data.OutgoingCallCostInUsdCents = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["outgoingCallCostInUSDCents"].(int64); ok {
        data.OutgoingCallCostInUsdCents = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["outgoingCallCostInUSDCents"] == nil {
        data.OutgoingCallCostInUsdCents = types.NumberNull()
    }
    if obj, ok := dataMap["startedAt"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StartedAt = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.StartedAt = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.StartedAt = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.StartedAt = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.StartedAt = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.StartedAt = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.StartedAt = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.StartedAt = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.StartedAt = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["startedAt"].(string); ok && val != "" {
        data.StartedAt = NewJSONSubsetValue(val)
    } else {
        data.StartedAt = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["endedAt"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EndedAt = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.EndedAt = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.EndedAt = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.EndedAt = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.EndedAt = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.EndedAt = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.EndedAt = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.EndedAt = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.EndedAt = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["endedAt"].(string); ok && val != "" {
        data.EndedAt = NewJSONSubsetValue(val)
    } else {
        data.EndedAt = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["answeredByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AnsweredByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AnsweredByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AnsweredByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AnsweredByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.AnsweredByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AnsweredByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.AnsweredByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AnsweredByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.AnsweredByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["answeredByUserId"].(string); ok && val != "" {
        data.AnsweredByUserId = types.StringValue(val)
    } else {
        data.AnsweredByUserId = types.StringNull()
    }
    if val, ok := dataMap["currentEscalationRuleOrder"].(float64); ok {
        data.CurrentEscalationRuleOrder = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["currentEscalationRuleOrder"].(int); ok {
        data.CurrentEscalationRuleOrder = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["currentEscalationRuleOrder"].(int64); ok {
        data.CurrentEscalationRuleOrder = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["currentEscalationRuleOrder"] == nil {
        data.CurrentEscalationRuleOrder = types.NumberNull()
    }
    if val, ok := dataMap["repeatCount"].(float64); ok {
        data.RepeatCount = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["repeatCount"].(int); ok {
        data.RepeatCount = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["repeatCount"].(int64); ok {
        data.RepeatCount = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["repeatCount"] == nil {
        data.RepeatCount = types.NumberNull()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "created a resource")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IncomingCallLogResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var data IncomingCallLogResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create select parameter to get full object
    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "projectId": true,
        "incomingCallPolicyId": true,
        "callerPhoneNumber": true,
        "routingPhoneNumber": true,
        "callProviderCallId": true,
        "status": true,
        "statusMessage": true,
        "callDurationInSeconds": true,
        "callCostInUSDCents": true,
        "incomingCallCostInUSDCents": true,
        "outgoingCallCostInUSDCents": true,
        "startedAt": true,
        "endedAt": true,
        "answeredByUserId": true,
        "currentEscalationRuleOrder": true,
        "repeatCount": true,
        "_id": true,
    }

    // Make API call with select parameter
    httpResp, err := r.client.PostWithSelect("/incoming-call-log/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read incoming_call_log, got error: %s", err))
        return
    }

    if httpResp.StatusCode == http.StatusNotFound {
        resp.State.RemoveResource(ctx)
        return
    }

    var incomingCallLogResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &incomingCallLogResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse incoming_call_log response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := incomingCallLogResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = incomingCallLogResponse
    }

    if obj, ok := dataMap["id"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Id = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Id = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Id = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Id = types.StringValue(string(jsonBytes))
            } else {
                data.Id = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Id = types.StringValue(string(jsonBytes))
            } else {
                data.Id = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Id = types.StringValue(string(jsonBytes))
        } else {
            data.Id = types.StringNull()
        }
    } else if val, ok := dataMap["id"].(string); ok && val != "" {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    if obj, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedAt = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedAt = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedAt = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CreatedAt = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.CreatedAt = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CreatedAt = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.CreatedAt = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedAt = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.CreatedAt = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["createdAt"].(string); ok && val != "" {
        data.CreatedAt = NewJSONSubsetValue(val)
    } else {
        data.CreatedAt = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UpdatedAt = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.UpdatedAt = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.UpdatedAt = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.UpdatedAt = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.UpdatedAt = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.UpdatedAt = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.UpdatedAt = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.UpdatedAt = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.UpdatedAt = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["updatedAt"].(string); ok && val != "" {
        data.UpdatedAt = NewJSONSubsetValue(val)
    } else {
        data.UpdatedAt = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["deletedAt"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeletedAt = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DeletedAt = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DeletedAt = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DeletedAt = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DeletedAt = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DeletedAt = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DeletedAt = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DeletedAt = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DeletedAt = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["deletedAt"].(string); ok && val != "" {
        data.DeletedAt = NewJSONSubsetValue(val)
    } else {
        data.DeletedAt = NewJSONSubsetNull()
    }
    if val, ok := dataMap["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["version"].(int); ok {
        data.Version = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["version"].(int64); ok {
        data.Version = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["version"] == nil {
        data.Version = types.NumberNull()
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
    if obj, ok := dataMap["incomingCallPolicyId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncomingCallPolicyId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.IncomingCallPolicyId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.IncomingCallPolicyId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.IncomingCallPolicyId = types.StringValue(string(jsonBytes))
            } else {
                data.IncomingCallPolicyId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.IncomingCallPolicyId = types.StringValue(string(jsonBytes))
            } else {
                data.IncomingCallPolicyId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.IncomingCallPolicyId = types.StringValue(string(jsonBytes))
        } else {
            data.IncomingCallPolicyId = types.StringNull()
        }
    } else if val, ok := dataMap["incomingCallPolicyId"].(string); ok && val != "" {
        data.IncomingCallPolicyId = types.StringValue(val)
    } else {
        data.IncomingCallPolicyId = types.StringNull()
    }
    if obj, ok := dataMap["callerPhoneNumber"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CallerPhoneNumber = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CallerPhoneNumber = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CallerPhoneNumber = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CallerPhoneNumber = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.CallerPhoneNumber = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CallerPhoneNumber = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.CallerPhoneNumber = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CallerPhoneNumber = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.CallerPhoneNumber = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["callerPhoneNumber"].(string); ok && val != "" {
        data.CallerPhoneNumber = NewJSONSubsetValue(val)
    } else {
        data.CallerPhoneNumber = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["routingPhoneNumber"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RoutingPhoneNumber = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.RoutingPhoneNumber = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.RoutingPhoneNumber = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.RoutingPhoneNumber = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.RoutingPhoneNumber = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.RoutingPhoneNumber = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.RoutingPhoneNumber = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.RoutingPhoneNumber = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.RoutingPhoneNumber = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["routingPhoneNumber"].(string); ok && val != "" {
        data.RoutingPhoneNumber = NewJSONSubsetValue(val)
    } else {
        data.RoutingPhoneNumber = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["callProviderCallId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CallProviderCallId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CallProviderCallId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CallProviderCallId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CallProviderCallId = types.StringValue(string(jsonBytes))
            } else {
                data.CallProviderCallId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CallProviderCallId = types.StringValue(string(jsonBytes))
            } else {
                data.CallProviderCallId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CallProviderCallId = types.StringValue(string(jsonBytes))
        } else {
            data.CallProviderCallId = types.StringNull()
        }
    } else if val, ok := dataMap["callProviderCallId"].(string); ok && val != "" {
        data.CallProviderCallId = types.StringValue(val)
    } else {
        data.CallProviderCallId = types.StringNull()
    }
    if obj, ok := dataMap["status"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Status = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Status = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Status = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Status = types.StringValue(string(jsonBytes))
            } else {
                data.Status = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Status = types.StringValue(string(jsonBytes))
            } else {
                data.Status = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Status = types.StringValue(string(jsonBytes))
        } else {
            data.Status = types.StringNull()
        }
    } else if val, ok := dataMap["status"].(string); ok && val != "" {
        data.Status = types.StringValue(val)
    } else {
        data.Status = types.StringNull()
    }
    if obj, ok := dataMap["statusMessage"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StatusMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.StatusMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.StatusMessage = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.StatusMessage = types.StringValue(string(jsonBytes))
            } else {
                data.StatusMessage = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.StatusMessage = types.StringValue(string(jsonBytes))
            } else {
                data.StatusMessage = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.StatusMessage = types.StringValue(string(jsonBytes))
        } else {
            data.StatusMessage = types.StringNull()
        }
    } else if val, ok := dataMap["statusMessage"].(string); ok && val != "" {
        data.StatusMessage = types.StringValue(val)
    } else {
        data.StatusMessage = types.StringNull()
    }
    if val, ok := dataMap["callDurationInSeconds"].(float64); ok {
        data.CallDurationInSeconds = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["callDurationInSeconds"].(int); ok {
        data.CallDurationInSeconds = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["callDurationInSeconds"].(int64); ok {
        data.CallDurationInSeconds = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["callDurationInSeconds"] == nil {
        data.CallDurationInSeconds = types.NumberNull()
    }
    if val, ok := dataMap["callCostInUSDCents"].(float64); ok {
        data.CallCostInUsdCents = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["callCostInUSDCents"].(int); ok {
        data.CallCostInUsdCents = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["callCostInUSDCents"].(int64); ok {
        data.CallCostInUsdCents = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["callCostInUSDCents"] == nil {
        data.CallCostInUsdCents = types.NumberNull()
    }
    if val, ok := dataMap["incomingCallCostInUSDCents"].(float64); ok {
        data.IncomingCallCostInUsdCents = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["incomingCallCostInUSDCents"].(int); ok {
        data.IncomingCallCostInUsdCents = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["incomingCallCostInUSDCents"].(int64); ok {
        data.IncomingCallCostInUsdCents = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["incomingCallCostInUSDCents"] == nil {
        data.IncomingCallCostInUsdCents = types.NumberNull()
    }
    if val, ok := dataMap["outgoingCallCostInUSDCents"].(float64); ok {
        data.OutgoingCallCostInUsdCents = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["outgoingCallCostInUSDCents"].(int); ok {
        data.OutgoingCallCostInUsdCents = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["outgoingCallCostInUSDCents"].(int64); ok {
        data.OutgoingCallCostInUsdCents = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["outgoingCallCostInUSDCents"] == nil {
        data.OutgoingCallCostInUsdCents = types.NumberNull()
    }
    if obj, ok := dataMap["startedAt"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StartedAt = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.StartedAt = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.StartedAt = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.StartedAt = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.StartedAt = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.StartedAt = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.StartedAt = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.StartedAt = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.StartedAt = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["startedAt"].(string); ok && val != "" {
        data.StartedAt = NewJSONSubsetValue(val)
    } else {
        data.StartedAt = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["endedAt"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EndedAt = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.EndedAt = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.EndedAt = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.EndedAt = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.EndedAt = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.EndedAt = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.EndedAt = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.EndedAt = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.EndedAt = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["endedAt"].(string); ok && val != "" {
        data.EndedAt = NewJSONSubsetValue(val)
    } else {
        data.EndedAt = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["answeredByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AnsweredByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AnsweredByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AnsweredByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AnsweredByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.AnsweredByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AnsweredByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.AnsweredByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AnsweredByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.AnsweredByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["answeredByUserId"].(string); ok && val != "" {
        data.AnsweredByUserId = types.StringValue(val)
    } else {
        data.AnsweredByUserId = types.StringNull()
    }
    if val, ok := dataMap["currentEscalationRuleOrder"].(float64); ok {
        data.CurrentEscalationRuleOrder = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["currentEscalationRuleOrder"].(int); ok {
        data.CurrentEscalationRuleOrder = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["currentEscalationRuleOrder"].(int64); ok {
        data.CurrentEscalationRuleOrder = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["currentEscalationRuleOrder"] == nil {
        data.CurrentEscalationRuleOrder = types.NumberNull()
    }
    if val, ok := dataMap["repeatCount"].(float64); ok {
        data.RepeatCount = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["repeatCount"].(int); ok {
        data.RepeatCount = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["repeatCount"].(int64); ok {
        data.RepeatCount = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["repeatCount"] == nil {
        data.RepeatCount = types.NumberNull()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IncomingCallLogResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var data IncomingCallLogResourceModel

    // Read Terraform plan data into the model
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // This resource does not have an update API endpoint.
    // Preserve the planned state.
    tflog.Trace(ctx, "updated a resource (no-op: preserving planned state)")

    // Save planned data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IncomingCallLogResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    // This resource does not have a delete API endpoint.
    // Simply remove the resource from Terraform state.
    tflog.Trace(ctx, "deleted a resource (no-op: removed from state only)")
}


func (r *IncomingCallLogResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Helper method to convert Terraform map to Go interface{}
func (r *IncomingCallLogResource) convertTerraformMapToInterface(terraformMap types.Map) interface{} {
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
func (r *IncomingCallLogResource) convertTerraformListToInterface(terraformList types.List) interface{} {
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
func (r *IncomingCallLogResource) convertTerraformSetToInterface(terraformSet types.Set) interface{} {
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
func (r *IncomingCallLogResource) parseJSONField(terraformString basetypes.StringValuable) interface{} {
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
func (r *IncomingCallLogResource) normalizeURLWrappers(value interface{}) interface{} {
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

func (r *IncomingCallLogResource) normalizeURLString(value string) string {
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
func (r *IncomingCallLogResource) bigFloatToFloat64(bf *big.Float) interface{} {
    if bf == nil {
        return nil
    }
    f, _ := bf.Float64()
    return f
}

// Helper method to check if a type string is a valid OneUptime ObjectType
// Only these types should be marshalled/unmarshalled as typed wrapper objects
// This list is dynamically generated from Common/Types/JSON.ts ObjectType enum
func (r *IncomingCallLogResource) isValidOneUptimeObjectType(typeStr string) bool {
    validTypes := map[string]bool{
        "ObjectID": true,
        "Decimal": true,
        "Name": true,
        "EqualTo": true,
        "EqualToOrNull": true,
        "MonitorSteps": true,
        "MonitorStep": true,
        "Recurring": true,
        "RestrictionTimes": true,
        "MonitorCriteria": true,
        "PositiveNumber": true,
        "MonitorCriteriaInstance": true,
        "NotEqual": true,
        "Email": true,
        "Phone": true,
        "Color": true,
        "Domain": true,
        "Version": true,
        "IP": true,
        "Route": true,
        "URL": true,
        "Permission": true,
        "Search": true,
        "MultiSearch": true,
        "GreaterThan": true,
        "GreaterThanOrEqual": true,
        "GreaterThanOrNull": true,
        "LessThanOrNull": true,
        "LessThan": true,
        "LessThanOrEqual": true,
        "Port": true,
        "Hostname": true,
        "HashedString": true,
        "DateTime": true,
        "Buffer": true,
        "InBetween": true,
        "NotNull": true,
        "IsNull": true,
        "Includes": true,
        "IncludesAll": true,
        "IncludesNone": true,
        "StartsWith": true,
        "EndsWith": true,
        "NotContains": true,
        "DashboardComponent": true,
        "DashboardViewConfig": true,
    }
    return validTypes[typeStr]
}
