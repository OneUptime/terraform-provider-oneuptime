package provider

import (
    "context"
    "fmt"
    "github.com/hashicorp/terraform-plugin-framework/path"
    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
    "math/big"
    "net/http"
    "encoding/json"
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
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    IncomingCallPolicyId types.String `tfsdk:"incoming_call_policy_id"`
    CallerPhoneNumber types.String `tfsdk:"caller_phone_number"`
    RoutingPhoneNumber types.String `tfsdk:"routing_phone_number"`
    CallProviderCallId types.String `tfsdk:"call_provider_call_id"`
    Status types.String `tfsdk:"status"`
    StatusMessage types.String `tfsdk:"status_message"`
    CallDurationInSeconds types.Number `tfsdk:"call_duration_in_seconds"`
    CallCostInUsdCents types.Number `tfsdk:"call_cost_in_usd_cents"`
    IncomingCallCostInUsdCents types.Number `tfsdk:"incoming_call_cost_in_usd_cents"`
    OutgoingCallCostInUsdCents types.Number `tfsdk:"outgoing_call_cost_in_usd_cents"`
    StartedAt types.String `tfsdk:"started_at"`
    EndedAt types.String `tfsdk:"ended_at"`
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
            "incoming_call_policy_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "caller_phone_number": schema.StringAttribute{
                MarkdownDescription: "Phone object",
                Computed: true,
            },
            "routing_phone_number": schema.StringAttribute{
                MarkdownDescription: "Phone object",
                Computed: true,
            },
            "call_provider_call_id": schema.StringAttribute{
                MarkdownDescription: "Call provider's call identifier. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Incoming Call Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "status": schema.StringAttribute{
                MarkdownDescription: "Current status of the incoming call. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Incoming Call Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "status_message": schema.StringAttribute{
                MarkdownDescription: "Additional status information. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Incoming Call Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "call_duration_in_seconds": schema.NumberAttribute{
                MarkdownDescription: "Total call duration in seconds. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Incoming Call Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "call_cost_in_usd_cents": schema.NumberAttribute{
                MarkdownDescription: "Total cost for this call in USD cents. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Incoming Call Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "incoming_call_cost_in_usd_cents": schema.NumberAttribute{
                MarkdownDescription: "Cost for incoming leg in USD cents. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Incoming Call Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "outgoing_call_cost_in_usd_cents": schema.NumberAttribute{
                MarkdownDescription: "Cost for all forwarding attempts in USD cents. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Incoming Call Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "started_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "ended_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "answered_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "current_escalation_rule_order": schema.NumberAttribute{
                MarkdownDescription: "The current escalation rule order being processed. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Incoming Call Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "repeat_count": schema.NumberAttribute{
                MarkdownDescription: "Number of times the policy has been repeated. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Incoming Call Log], Update: [No access - you don't have permission for this operation]",
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
        // Handle ObjectID and value wrapper type responses (e.g., Version with _type and value fields)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Id = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.Id = types.StringValue(val)
        } else {
            data.Id = types.StringNull()
        }
    } else if val, ok := dataMap["id"].(string); ok && val != "" {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    if val, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        // Check if this is a value wrapper type (e.g., Version, ObjectID with _type and value fields)
        if _, hasType := val["_type"]; hasType {
            if strVal, ok := val["value"].(string); ok && strVal != "" {
                data.CreatedAt = types.StringValue(strVal)
            } else {
                data.CreatedAt = types.StringNull()
            }
        } else if strVal, ok := val["_id"].(string); ok && strVal != "" {
            // Handle ObjectID type responses
            data.CreatedAt = types.StringValue(strVal)
        } else if strVal, ok := val["value"].(string); ok && strVal != "" {
            // Handle other value wrapper types
            data.CreatedAt = types.StringValue(strVal)
        } else {
            // Fall back to JSON marshalling for truly complex objects
            if jsonBytes, err := json.Marshal(val); err == nil {
                data.CreatedAt = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedAt = types.StringNull()
            }
        }
    } else if val, ok := dataMap["createdAt"].(string); ok && val != "" {
        data.CreatedAt = types.StringValue(val)
    } else {
        data.CreatedAt = types.StringNull()
    }
    if val, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        // Check if this is a value wrapper type (e.g., Version, ObjectID with _type and value fields)
        if _, hasType := val["_type"]; hasType {
            if strVal, ok := val["value"].(string); ok && strVal != "" {
                data.UpdatedAt = types.StringValue(strVal)
            } else {
                data.UpdatedAt = types.StringNull()
            }
        } else if strVal, ok := val["_id"].(string); ok && strVal != "" {
            // Handle ObjectID type responses
            data.UpdatedAt = types.StringValue(strVal)
        } else if strVal, ok := val["value"].(string); ok && strVal != "" {
            // Handle other value wrapper types
            data.UpdatedAt = types.StringValue(strVal)
        } else {
            // Fall back to JSON marshalling for truly complex objects
            if jsonBytes, err := json.Marshal(val); err == nil {
                data.UpdatedAt = types.StringValue(string(jsonBytes))
            } else {
                data.UpdatedAt = types.StringNull()
            }
        }
    } else if val, ok := dataMap["updatedAt"].(string); ok && val != "" {
        data.UpdatedAt = types.StringValue(val)
    } else {
        data.UpdatedAt = types.StringNull()
    }
    if val, ok := dataMap["deletedAt"].(map[string]interface{}); ok {
        // Check if this is a value wrapper type (e.g., Version, ObjectID with _type and value fields)
        if _, hasType := val["_type"]; hasType {
            if strVal, ok := val["value"].(string); ok && strVal != "" {
                data.DeletedAt = types.StringValue(strVal)
            } else {
                data.DeletedAt = types.StringNull()
            }
        } else if strVal, ok := val["_id"].(string); ok && strVal != "" {
            // Handle ObjectID type responses
            data.DeletedAt = types.StringValue(strVal)
        } else if strVal, ok := val["value"].(string); ok && strVal != "" {
            // Handle other value wrapper types
            data.DeletedAt = types.StringValue(strVal)
        } else {
            // Fall back to JSON marshalling for truly complex objects
            if jsonBytes, err := json.Marshal(val); err == nil {
                data.DeletedAt = types.StringValue(string(jsonBytes))
            } else {
                data.DeletedAt = types.StringNull()
            }
        }
    } else if val, ok := dataMap["deletedAt"].(string); ok && val != "" {
        data.DeletedAt = types.StringValue(val)
    } else {
        data.DeletedAt = types.StringNull()
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
        // Handle ObjectID and value wrapper type responses (e.g., Version with _type and value fields)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncomingCallPolicyId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.IncomingCallPolicyId = types.StringValue(val)
        } else {
            data.IncomingCallPolicyId = types.StringNull()
        }
    } else if val, ok := dataMap["incomingCallPolicyId"].(string); ok && val != "" {
        data.IncomingCallPolicyId = types.StringValue(val)
    } else {
        data.IncomingCallPolicyId = types.StringNull()
    }
    if val, ok := dataMap["callerPhoneNumber"].(map[string]interface{}); ok {
        // Check if this is a value wrapper type (e.g., Version, ObjectID with _type and value fields)
        if _, hasType := val["_type"]; hasType {
            if strVal, ok := val["value"].(string); ok && strVal != "" {
                data.CallerPhoneNumber = types.StringValue(strVal)
            } else {
                data.CallerPhoneNumber = types.StringNull()
            }
        } else if strVal, ok := val["_id"].(string); ok && strVal != "" {
            // Handle ObjectID type responses
            data.CallerPhoneNumber = types.StringValue(strVal)
        } else if strVal, ok := val["value"].(string); ok && strVal != "" {
            // Handle other value wrapper types
            data.CallerPhoneNumber = types.StringValue(strVal)
        } else {
            // Fall back to JSON marshalling for truly complex objects
            if jsonBytes, err := json.Marshal(val); err == nil {
                data.CallerPhoneNumber = types.StringValue(string(jsonBytes))
            } else {
                data.CallerPhoneNumber = types.StringNull()
            }
        }
    } else if val, ok := dataMap["callerPhoneNumber"].(string); ok && val != "" {
        data.CallerPhoneNumber = types.StringValue(val)
    } else {
        data.CallerPhoneNumber = types.StringNull()
    }
    if val, ok := dataMap["routingPhoneNumber"].(map[string]interface{}); ok {
        // Check if this is a value wrapper type (e.g., Version, ObjectID with _type and value fields)
        if _, hasType := val["_type"]; hasType {
            if strVal, ok := val["value"].(string); ok && strVal != "" {
                data.RoutingPhoneNumber = types.StringValue(strVal)
            } else {
                data.RoutingPhoneNumber = types.StringNull()
            }
        } else if strVal, ok := val["_id"].(string); ok && strVal != "" {
            // Handle ObjectID type responses
            data.RoutingPhoneNumber = types.StringValue(strVal)
        } else if strVal, ok := val["value"].(string); ok && strVal != "" {
            // Handle other value wrapper types
            data.RoutingPhoneNumber = types.StringValue(strVal)
        } else {
            // Fall back to JSON marshalling for truly complex objects
            if jsonBytes, err := json.Marshal(val); err == nil {
                data.RoutingPhoneNumber = types.StringValue(string(jsonBytes))
            } else {
                data.RoutingPhoneNumber = types.StringNull()
            }
        }
    } else if val, ok := dataMap["routingPhoneNumber"].(string); ok && val != "" {
        data.RoutingPhoneNumber = types.StringValue(val)
    } else {
        data.RoutingPhoneNumber = types.StringNull()
    }
    if obj, ok := dataMap["callProviderCallId"].(map[string]interface{}); ok {
        // Handle ObjectID and value wrapper type responses (e.g., Version with _type and value fields)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CallProviderCallId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.CallProviderCallId = types.StringValue(val)
        } else {
            data.CallProviderCallId = types.StringNull()
        }
    } else if val, ok := dataMap["callProviderCallId"].(string); ok && val != "" {
        data.CallProviderCallId = types.StringValue(val)
    } else {
        data.CallProviderCallId = types.StringNull()
    }
    if obj, ok := dataMap["status"].(map[string]interface{}); ok {
        // Handle ObjectID and value wrapper type responses (e.g., Version with _type and value fields)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Status = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.Status = types.StringValue(val)
        } else {
            data.Status = types.StringNull()
        }
    } else if val, ok := dataMap["status"].(string); ok && val != "" {
        data.Status = types.StringValue(val)
    } else {
        data.Status = types.StringNull()
    }
    if obj, ok := dataMap["statusMessage"].(map[string]interface{}); ok {
        // Handle ObjectID and value wrapper type responses (e.g., Version with _type and value fields)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StatusMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.StatusMessage = types.StringValue(val)
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
    if val, ok := dataMap["startedAt"].(map[string]interface{}); ok {
        // Check if this is a value wrapper type (e.g., Version, ObjectID with _type and value fields)
        if _, hasType := val["_type"]; hasType {
            if strVal, ok := val["value"].(string); ok && strVal != "" {
                data.StartedAt = types.StringValue(strVal)
            } else {
                data.StartedAt = types.StringNull()
            }
        } else if strVal, ok := val["_id"].(string); ok && strVal != "" {
            // Handle ObjectID type responses
            data.StartedAt = types.StringValue(strVal)
        } else if strVal, ok := val["value"].(string); ok && strVal != "" {
            // Handle other value wrapper types
            data.StartedAt = types.StringValue(strVal)
        } else {
            // Fall back to JSON marshalling for truly complex objects
            if jsonBytes, err := json.Marshal(val); err == nil {
                data.StartedAt = types.StringValue(string(jsonBytes))
            } else {
                data.StartedAt = types.StringNull()
            }
        }
    } else if val, ok := dataMap["startedAt"].(string); ok && val != "" {
        data.StartedAt = types.StringValue(val)
    } else {
        data.StartedAt = types.StringNull()
    }
    if val, ok := dataMap["endedAt"].(map[string]interface{}); ok {
        // Check if this is a value wrapper type (e.g., Version, ObjectID with _type and value fields)
        if _, hasType := val["_type"]; hasType {
            if strVal, ok := val["value"].(string); ok && strVal != "" {
                data.EndedAt = types.StringValue(strVal)
            } else {
                data.EndedAt = types.StringNull()
            }
        } else if strVal, ok := val["_id"].(string); ok && strVal != "" {
            // Handle ObjectID type responses
            data.EndedAt = types.StringValue(strVal)
        } else if strVal, ok := val["value"].(string); ok && strVal != "" {
            // Handle other value wrapper types
            data.EndedAt = types.StringValue(strVal)
        } else {
            // Fall back to JSON marshalling for truly complex objects
            if jsonBytes, err := json.Marshal(val); err == nil {
                data.EndedAt = types.StringValue(string(jsonBytes))
            } else {
                data.EndedAt = types.StringNull()
            }
        }
    } else if val, ok := dataMap["endedAt"].(string); ok && val != "" {
        data.EndedAt = types.StringValue(val)
    } else {
        data.EndedAt = types.StringNull()
    }
    if obj, ok := dataMap["answeredByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID and value wrapper type responses (e.g., Version with _type and value fields)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AnsweredByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.AnsweredByUserId = types.StringValue(val)
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
        // Handle ObjectID and value wrapper type responses (e.g., Version with _type and value fields)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Id = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.Id = types.StringValue(val)
        } else {
            data.Id = types.StringNull()
        }
    } else if val, ok := dataMap["id"].(string); ok && val != "" {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    if val, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        // Check if this is a value wrapper type (e.g., Version, ObjectID with _type and value fields)
        if _, hasType := val["_type"]; hasType {
            if strVal, ok := val["value"].(string); ok && strVal != "" {
                data.CreatedAt = types.StringValue(strVal)
            } else {
                data.CreatedAt = types.StringNull()
            }
        } else if strVal, ok := val["_id"].(string); ok && strVal != "" {
            // Handle ObjectID type responses
            data.CreatedAt = types.StringValue(strVal)
        } else if strVal, ok := val["value"].(string); ok && strVal != "" {
            // Handle other value wrapper types
            data.CreatedAt = types.StringValue(strVal)
        } else {
            // Fall back to JSON marshalling for truly complex objects
            if jsonBytes, err := json.Marshal(val); err == nil {
                data.CreatedAt = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedAt = types.StringNull()
            }
        }
    } else if val, ok := dataMap["createdAt"].(string); ok && val != "" {
        data.CreatedAt = types.StringValue(val)
    } else {
        data.CreatedAt = types.StringNull()
    }
    if val, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        // Check if this is a value wrapper type (e.g., Version, ObjectID with _type and value fields)
        if _, hasType := val["_type"]; hasType {
            if strVal, ok := val["value"].(string); ok && strVal != "" {
                data.UpdatedAt = types.StringValue(strVal)
            } else {
                data.UpdatedAt = types.StringNull()
            }
        } else if strVal, ok := val["_id"].(string); ok && strVal != "" {
            // Handle ObjectID type responses
            data.UpdatedAt = types.StringValue(strVal)
        } else if strVal, ok := val["value"].(string); ok && strVal != "" {
            // Handle other value wrapper types
            data.UpdatedAt = types.StringValue(strVal)
        } else {
            // Fall back to JSON marshalling for truly complex objects
            if jsonBytes, err := json.Marshal(val); err == nil {
                data.UpdatedAt = types.StringValue(string(jsonBytes))
            } else {
                data.UpdatedAt = types.StringNull()
            }
        }
    } else if val, ok := dataMap["updatedAt"].(string); ok && val != "" {
        data.UpdatedAt = types.StringValue(val)
    } else {
        data.UpdatedAt = types.StringNull()
    }
    if val, ok := dataMap["deletedAt"].(map[string]interface{}); ok {
        // Check if this is a value wrapper type (e.g., Version, ObjectID with _type and value fields)
        if _, hasType := val["_type"]; hasType {
            if strVal, ok := val["value"].(string); ok && strVal != "" {
                data.DeletedAt = types.StringValue(strVal)
            } else {
                data.DeletedAt = types.StringNull()
            }
        } else if strVal, ok := val["_id"].(string); ok && strVal != "" {
            // Handle ObjectID type responses
            data.DeletedAt = types.StringValue(strVal)
        } else if strVal, ok := val["value"].(string); ok && strVal != "" {
            // Handle other value wrapper types
            data.DeletedAt = types.StringValue(strVal)
        } else {
            // Fall back to JSON marshalling for truly complex objects
            if jsonBytes, err := json.Marshal(val); err == nil {
                data.DeletedAt = types.StringValue(string(jsonBytes))
            } else {
                data.DeletedAt = types.StringNull()
            }
        }
    } else if val, ok := dataMap["deletedAt"].(string); ok && val != "" {
        data.DeletedAt = types.StringValue(val)
    } else {
        data.DeletedAt = types.StringNull()
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
        // Handle ObjectID and value wrapper type responses (e.g., Version with _type and value fields)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncomingCallPolicyId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.IncomingCallPolicyId = types.StringValue(val)
        } else {
            data.IncomingCallPolicyId = types.StringNull()
        }
    } else if val, ok := dataMap["incomingCallPolicyId"].(string); ok && val != "" {
        data.IncomingCallPolicyId = types.StringValue(val)
    } else {
        data.IncomingCallPolicyId = types.StringNull()
    }
    if val, ok := dataMap["callerPhoneNumber"].(map[string]interface{}); ok {
        // Check if this is a value wrapper type (e.g., Version, ObjectID with _type and value fields)
        if _, hasType := val["_type"]; hasType {
            if strVal, ok := val["value"].(string); ok && strVal != "" {
                data.CallerPhoneNumber = types.StringValue(strVal)
            } else {
                data.CallerPhoneNumber = types.StringNull()
            }
        } else if strVal, ok := val["_id"].(string); ok && strVal != "" {
            // Handle ObjectID type responses
            data.CallerPhoneNumber = types.StringValue(strVal)
        } else if strVal, ok := val["value"].(string); ok && strVal != "" {
            // Handle other value wrapper types
            data.CallerPhoneNumber = types.StringValue(strVal)
        } else {
            // Fall back to JSON marshalling for truly complex objects
            if jsonBytes, err := json.Marshal(val); err == nil {
                data.CallerPhoneNumber = types.StringValue(string(jsonBytes))
            } else {
                data.CallerPhoneNumber = types.StringNull()
            }
        }
    } else if val, ok := dataMap["callerPhoneNumber"].(string); ok && val != "" {
        data.CallerPhoneNumber = types.StringValue(val)
    } else {
        data.CallerPhoneNumber = types.StringNull()
    }
    if val, ok := dataMap["routingPhoneNumber"].(map[string]interface{}); ok {
        // Check if this is a value wrapper type (e.g., Version, ObjectID with _type and value fields)
        if _, hasType := val["_type"]; hasType {
            if strVal, ok := val["value"].(string); ok && strVal != "" {
                data.RoutingPhoneNumber = types.StringValue(strVal)
            } else {
                data.RoutingPhoneNumber = types.StringNull()
            }
        } else if strVal, ok := val["_id"].(string); ok && strVal != "" {
            // Handle ObjectID type responses
            data.RoutingPhoneNumber = types.StringValue(strVal)
        } else if strVal, ok := val["value"].(string); ok && strVal != "" {
            // Handle other value wrapper types
            data.RoutingPhoneNumber = types.StringValue(strVal)
        } else {
            // Fall back to JSON marshalling for truly complex objects
            if jsonBytes, err := json.Marshal(val); err == nil {
                data.RoutingPhoneNumber = types.StringValue(string(jsonBytes))
            } else {
                data.RoutingPhoneNumber = types.StringNull()
            }
        }
    } else if val, ok := dataMap["routingPhoneNumber"].(string); ok && val != "" {
        data.RoutingPhoneNumber = types.StringValue(val)
    } else {
        data.RoutingPhoneNumber = types.StringNull()
    }
    if obj, ok := dataMap["callProviderCallId"].(map[string]interface{}); ok {
        // Handle ObjectID and value wrapper type responses (e.g., Version with _type and value fields)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CallProviderCallId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.CallProviderCallId = types.StringValue(val)
        } else {
            data.CallProviderCallId = types.StringNull()
        }
    } else if val, ok := dataMap["callProviderCallId"].(string); ok && val != "" {
        data.CallProviderCallId = types.StringValue(val)
    } else {
        data.CallProviderCallId = types.StringNull()
    }
    if obj, ok := dataMap["status"].(map[string]interface{}); ok {
        // Handle ObjectID and value wrapper type responses (e.g., Version with _type and value fields)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Status = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.Status = types.StringValue(val)
        } else {
            data.Status = types.StringNull()
        }
    } else if val, ok := dataMap["status"].(string); ok && val != "" {
        data.Status = types.StringValue(val)
    } else {
        data.Status = types.StringNull()
    }
    if obj, ok := dataMap["statusMessage"].(map[string]interface{}); ok {
        // Handle ObjectID and value wrapper type responses (e.g., Version with _type and value fields)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StatusMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.StatusMessage = types.StringValue(val)
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
    if val, ok := dataMap["startedAt"].(map[string]interface{}); ok {
        // Check if this is a value wrapper type (e.g., Version, ObjectID with _type and value fields)
        if _, hasType := val["_type"]; hasType {
            if strVal, ok := val["value"].(string); ok && strVal != "" {
                data.StartedAt = types.StringValue(strVal)
            } else {
                data.StartedAt = types.StringNull()
            }
        } else if strVal, ok := val["_id"].(string); ok && strVal != "" {
            // Handle ObjectID type responses
            data.StartedAt = types.StringValue(strVal)
        } else if strVal, ok := val["value"].(string); ok && strVal != "" {
            // Handle other value wrapper types
            data.StartedAt = types.StringValue(strVal)
        } else {
            // Fall back to JSON marshalling for truly complex objects
            if jsonBytes, err := json.Marshal(val); err == nil {
                data.StartedAt = types.StringValue(string(jsonBytes))
            } else {
                data.StartedAt = types.StringNull()
            }
        }
    } else if val, ok := dataMap["startedAt"].(string); ok && val != "" {
        data.StartedAt = types.StringValue(val)
    } else {
        data.StartedAt = types.StringNull()
    }
    if val, ok := dataMap["endedAt"].(map[string]interface{}); ok {
        // Check if this is a value wrapper type (e.g., Version, ObjectID with _type and value fields)
        if _, hasType := val["_type"]; hasType {
            if strVal, ok := val["value"].(string); ok && strVal != "" {
                data.EndedAt = types.StringValue(strVal)
            } else {
                data.EndedAt = types.StringNull()
            }
        } else if strVal, ok := val["_id"].(string); ok && strVal != "" {
            // Handle ObjectID type responses
            data.EndedAt = types.StringValue(strVal)
        } else if strVal, ok := val["value"].(string); ok && strVal != "" {
            // Handle other value wrapper types
            data.EndedAt = types.StringValue(strVal)
        } else {
            // Fall back to JSON marshalling for truly complex objects
            if jsonBytes, err := json.Marshal(val); err == nil {
                data.EndedAt = types.StringValue(string(jsonBytes))
            } else {
                data.EndedAt = types.StringNull()
            }
        }
    } else if val, ok := dataMap["endedAt"].(string); ok && val != "" {
        data.EndedAt = types.StringValue(val)
    } else {
        data.EndedAt = types.StringNull()
    }
    if obj, ok := dataMap["answeredByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID and value wrapper type responses (e.g., Version with _type and value fields)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AnsweredByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.AnsweredByUserId = types.StringValue(val)
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
    resp.Diagnostics.AddError(
        "Update Not Implemented",
        "This resource does not support update operations",
    )
}

func (r *IncomingCallLogResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    resp.Diagnostics.AddError(
        "Delete Not Implemented",
        "This resource does not support delete operations", 
    )
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

// Helper method to parse JSON field for complex objects
func (r *IncomingCallLogResource) parseJSONField(terraformString types.String) interface{} {
    if terraformString.IsNull() || terraformString.IsUnknown() || terraformString.ValueString() == "" {
        return nil
    }
    
    var result interface{}
    if err := json.Unmarshal([]byte(terraformString.ValueString()), &result); err != nil {
        // If JSON parsing fails, return the raw string
        return terraformString.ValueString()
    }
    
    return result
}
