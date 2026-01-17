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
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &IncomingCallLogItemResource{}
var _ resource.ResourceWithImportState = &IncomingCallLogItemResource{}

func NewIncomingCallLogItemResource() resource.Resource {
    return &IncomingCallLogItemResource{}
}

// IncomingCallLogItemResource defines the resource implementation.
type IncomingCallLogItemResource struct {
    client *Client
}

// IncomingCallLogItemResourceModel describes the resource data model.
type IncomingCallLogItemResourceModel struct {
    Id types.String `tfsdk:"id"`
    ProjectId types.String `tfsdk:"project_id"`
    IncomingCallLogId types.String `tfsdk:"incoming_call_log_id"`
    IncomingCallPolicyEscalationRuleId types.String `tfsdk:"incoming_call_policy_escalation_rule_id"`
    UserId types.String `tfsdk:"user_id"`
    UserPhoneNumber types.String `tfsdk:"user_phone_number"`
    Status types.String `tfsdk:"status"`
    StatusMessage types.String `tfsdk:"status_message"`
    DialDurationInSeconds types.Number `tfsdk:"dial_duration_in_seconds"`
    CallCostInUsdCents types.Number `tfsdk:"call_cost_in_usd_cents"`
    StartedAt types.String `tfsdk:"started_at"`
    EndedAt types.String `tfsdk:"ended_at"`
    IsAnswered types.Bool `tfsdk:"is_answered"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
}

func (r *IncomingCallLogItemResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_incoming_call_log_item"
}

func (r *IncomingCallLogItemResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "incoming_call_log_item resource",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Unique identifier for the resource",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Optional: true,
                Computed: true,
            },
            "incoming_call_log_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Required: true,
            },
            "incoming_call_policy_escalation_rule_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Optional: true,
                Computed: true,
            },
            "user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Optional: true,
                Computed: true,
            },
            "user_phone_number": schema.StringAttribute{
                MarkdownDescription: "Phone object",
                Optional: true,
                Computed: true,
            },
            "status": schema.StringAttribute{
                MarkdownDescription: "Status of this dial attempt. Permissions - Create: [Project Owner, Project Admin], Read: [Project Owner, Project Admin, Project Member, Read Incoming Call Log Item], Update: [Project Owner, Project Admin]",
                Required: true,
            },
            "status_message": schema.StringAttribute{
                MarkdownDescription: "Additional status information. Permissions - Create: [Project Owner, Project Admin], Read: [Project Owner, Project Admin, Project Member, Read Incoming Call Log Item], Update: [Project Owner, Project Admin]",
                Optional: true,
                Computed: true,
            },
            "dial_duration_in_seconds": schema.NumberAttribute{
                MarkdownDescription: "How long this dial lasted in seconds. Permissions - Create: [Project Owner, Project Admin], Read: [Project Owner, Project Admin, Project Member, Read Incoming Call Log Item], Update: [Project Owner, Project Admin]",
                Optional: true,
                Computed: true,
            },
            "call_cost_in_usd_cents": schema.NumberAttribute{
                MarkdownDescription: "Cost for this dial attempt in USD cents. Permissions - Create: [Project Owner, Project Admin], Read: [Project Owner, Project Admin, Project Member, Read Incoming Call Log Item], Update: [Project Owner, Project Admin]",
                Optional: true,
                Computed: true,
            },
            "started_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Optional: true,
                Computed: true,
            },
            "ended_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Optional: true,
                Computed: true,
            },
            "is_answered": schema.BoolAttribute{
                MarkdownDescription: "Whether this user answered the call. Permissions - Create: [Project Owner, Project Admin], Read: [Project Owner, Project Admin, Project Member, Read Incoming Call Log Item], Update: [Project Owner, Project Admin]",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
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
        },
    }
}

func (r *IncomingCallLogItemResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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


func (r *IncomingCallLogItemResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var data IncomingCallLogItemResourceModel

    // Read Terraform plan data into the model
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }



    // Create API request body
    incomingCallLogItemRequest := map[string]interface{}{
        "data": map[string]interface{}{
        "projectId": data.ProjectId.ValueString(),
        "incomingCallLogId": data.IncomingCallLogId.ValueString(),
        "incomingCallPolicyEscalationRuleId": data.IncomingCallPolicyEscalationRuleId.ValueString(),
        "userId": data.UserId.ValueString(),
        "userPhoneNumber": r.parseJSONField(data.UserPhoneNumber),
        "status": data.Status.ValueString(),
        "statusMessage": data.StatusMessage.ValueString(),
        "dialDurationInSeconds": data.DialDurationInSeconds.ValueBigFloat(),
        "callCostInUSDCents": data.CallCostInUsdCents.ValueBigFloat(),
        "startedAt": r.parseJSONField(data.StartedAt),
        "endedAt": r.parseJSONField(data.EndedAt),
        "isAnswered": data.IsAnswered.ValueBool(),
        },
    }

    // Make API call
    httpResp, err := r.client.Post("/incoming-call-log-item", incomingCallLogItemRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create incoming_call_log_item, got error: %s", err))
        return
    }

    var incomingCallLogItemResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &incomingCallLogItemResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse incoming_call_log_item response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := incomingCallLogItemResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = incomingCallLogItemResponse
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
    if obj, ok := dataMap["incomingCallLogId"].(map[string]interface{}); ok {
        // Handle ObjectID and value wrapper type responses (e.g., Version with _type and value fields)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncomingCallLogId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.IncomingCallLogId = types.StringValue(val)
        } else {
            data.IncomingCallLogId = types.StringNull()
        }
    } else if val, ok := dataMap["incomingCallLogId"].(string); ok && val != "" {
        data.IncomingCallLogId = types.StringValue(val)
    } else {
        data.IncomingCallLogId = types.StringNull()
    }
    if obj, ok := dataMap["incomingCallPolicyEscalationRuleId"].(map[string]interface{}); ok {
        // Handle ObjectID and value wrapper type responses (e.g., Version with _type and value fields)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncomingCallPolicyEscalationRuleId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.IncomingCallPolicyEscalationRuleId = types.StringValue(val)
        } else {
            data.IncomingCallPolicyEscalationRuleId = types.StringNull()
        }
    } else if val, ok := dataMap["incomingCallPolicyEscalationRuleId"].(string); ok && val != "" {
        data.IncomingCallPolicyEscalationRuleId = types.StringValue(val)
    } else {
        data.IncomingCallPolicyEscalationRuleId = types.StringNull()
    }
    if obj, ok := dataMap["userId"].(map[string]interface{}); ok {
        // Handle ObjectID and value wrapper type responses (e.g., Version with _type and value fields)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.UserId = types.StringValue(val)
        } else {
            data.UserId = types.StringNull()
        }
    } else if val, ok := dataMap["userId"].(string); ok && val != "" {
        data.UserId = types.StringValue(val)
    } else {
        data.UserId = types.StringNull()
    }
    if val, ok := dataMap["userPhoneNumber"].(map[string]interface{}); ok {
        // Check if this is a value wrapper type (e.g., Version, ObjectID with _type and value fields)
        if _, hasType := val["_type"]; hasType {
            if strVal, ok := val["value"].(string); ok && strVal != "" {
                data.UserPhoneNumber = types.StringValue(strVal)
            } else {
                data.UserPhoneNumber = types.StringNull()
            }
        } else if strVal, ok := val["_id"].(string); ok && strVal != "" {
            // Handle ObjectID type responses
            data.UserPhoneNumber = types.StringValue(strVal)
        } else if strVal, ok := val["value"].(string); ok && strVal != "" {
            // Handle other value wrapper types
            data.UserPhoneNumber = types.StringValue(strVal)
        } else {
            // Fall back to JSON marshalling for truly complex objects
            if jsonBytes, err := json.Marshal(val); err == nil {
                data.UserPhoneNumber = types.StringValue(string(jsonBytes))
            } else {
                data.UserPhoneNumber = types.StringNull()
            }
        }
    } else if val, ok := dataMap["userPhoneNumber"].(string); ok && val != "" {
        data.UserPhoneNumber = types.StringValue(val)
    } else {
        data.UserPhoneNumber = types.StringNull()
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
    if val, ok := dataMap["dialDurationInSeconds"].(float64); ok {
        data.DialDurationInSeconds = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["dialDurationInSeconds"].(int); ok {
        data.DialDurationInSeconds = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["dialDurationInSeconds"].(int64); ok {
        data.DialDurationInSeconds = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["dialDurationInSeconds"] == nil {
        data.DialDurationInSeconds = types.NumberNull()
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
    if val, ok := dataMap["isAnswered"].(bool); ok {
        data.IsAnswered = types.BoolValue(val)
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

func (r *IncomingCallLogItemResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var data IncomingCallLogItemResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create select parameter to get full object
    selectParam := map[string]interface{}{
        "projectId": true,
        "incomingCallLogId": true,
        "incomingCallPolicyEscalationRuleId": true,
        "userId": true,
        "userPhoneNumber": true,
        "status": true,
        "statusMessage": true,
        "dialDurationInSeconds": true,
        "callCostInUSDCents": true,
        "startedAt": true,
        "endedAt": true,
        "isAnswered": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "_id": true,
    }

    // Make API call with select parameter
    httpResp, err := r.client.PostWithSelect("/incoming-call-log-item/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read incoming_call_log_item, got error: %s", err))
        return
    }

    if httpResp.StatusCode == http.StatusNotFound {
        resp.State.RemoveResource(ctx)
        return
    }

    var incomingCallLogItemResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &incomingCallLogItemResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse incoming_call_log_item response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := incomingCallLogItemResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = incomingCallLogItemResponse
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
    if obj, ok := dataMap["incomingCallLogId"].(map[string]interface{}); ok {
        // Handle ObjectID and value wrapper type responses (e.g., Version with _type and value fields)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncomingCallLogId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.IncomingCallLogId = types.StringValue(val)
        } else {
            data.IncomingCallLogId = types.StringNull()
        }
    } else if val, ok := dataMap["incomingCallLogId"].(string); ok && val != "" {
        data.IncomingCallLogId = types.StringValue(val)
    } else {
        data.IncomingCallLogId = types.StringNull()
    }
    if obj, ok := dataMap["incomingCallPolicyEscalationRuleId"].(map[string]interface{}); ok {
        // Handle ObjectID and value wrapper type responses (e.g., Version with _type and value fields)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncomingCallPolicyEscalationRuleId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.IncomingCallPolicyEscalationRuleId = types.StringValue(val)
        } else {
            data.IncomingCallPolicyEscalationRuleId = types.StringNull()
        }
    } else if val, ok := dataMap["incomingCallPolicyEscalationRuleId"].(string); ok && val != "" {
        data.IncomingCallPolicyEscalationRuleId = types.StringValue(val)
    } else {
        data.IncomingCallPolicyEscalationRuleId = types.StringNull()
    }
    if obj, ok := dataMap["userId"].(map[string]interface{}); ok {
        // Handle ObjectID and value wrapper type responses (e.g., Version with _type and value fields)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.UserId = types.StringValue(val)
        } else {
            data.UserId = types.StringNull()
        }
    } else if val, ok := dataMap["userId"].(string); ok && val != "" {
        data.UserId = types.StringValue(val)
    } else {
        data.UserId = types.StringNull()
    }
    if val, ok := dataMap["userPhoneNumber"].(map[string]interface{}); ok {
        // Check if this is a value wrapper type (e.g., Version, ObjectID with _type and value fields)
        if _, hasType := val["_type"]; hasType {
            if strVal, ok := val["value"].(string); ok && strVal != "" {
                data.UserPhoneNumber = types.StringValue(strVal)
            } else {
                data.UserPhoneNumber = types.StringNull()
            }
        } else if strVal, ok := val["_id"].(string); ok && strVal != "" {
            // Handle ObjectID type responses
            data.UserPhoneNumber = types.StringValue(strVal)
        } else if strVal, ok := val["value"].(string); ok && strVal != "" {
            // Handle other value wrapper types
            data.UserPhoneNumber = types.StringValue(strVal)
        } else {
            // Fall back to JSON marshalling for truly complex objects
            if jsonBytes, err := json.Marshal(val); err == nil {
                data.UserPhoneNumber = types.StringValue(string(jsonBytes))
            } else {
                data.UserPhoneNumber = types.StringNull()
            }
        }
    } else if val, ok := dataMap["userPhoneNumber"].(string); ok && val != "" {
        data.UserPhoneNumber = types.StringValue(val)
    } else {
        data.UserPhoneNumber = types.StringNull()
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
    if val, ok := dataMap["dialDurationInSeconds"].(float64); ok {
        data.DialDurationInSeconds = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["dialDurationInSeconds"].(int); ok {
        data.DialDurationInSeconds = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["dialDurationInSeconds"].(int64); ok {
        data.DialDurationInSeconds = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["dialDurationInSeconds"] == nil {
        data.DialDurationInSeconds = types.NumberNull()
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
    if val, ok := dataMap["isAnswered"].(bool); ok {
        data.IsAnswered = types.BoolValue(val)
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
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IncomingCallLogItemResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var data IncomingCallLogItemResourceModel
    var state IncomingCallLogItemResourceModel

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
    incomingCallLogItemRequest := map[string]interface{}{
        "data": map[string]interface{}{},
    }
    requestDataMap := incomingCallLogItemRequest["data"].(map[string]interface{})

    if !data.Status.IsUnknown() && !state.Status.IsUnknown() && !data.Status.Equal(state.Status) {
        requestDataMap["status"] = data.Status.ValueString()
    }
    if !data.StatusMessage.IsUnknown() && !state.StatusMessage.IsUnknown() && !data.StatusMessage.Equal(state.StatusMessage) {
        requestDataMap["statusMessage"] = data.StatusMessage.ValueString()
    }
    if !data.DialDurationInSeconds.IsUnknown() && !state.DialDurationInSeconds.IsUnknown() && !data.DialDurationInSeconds.Equal(state.DialDurationInSeconds) {
        requestDataMap["dialDurationInSeconds"] = data.DialDurationInSeconds.ValueBigFloat()
    }
    if !data.CallCostInUsdCents.IsUnknown() && !state.CallCostInUsdCents.IsUnknown() && !data.CallCostInUsdCents.Equal(state.CallCostInUsdCents) {
        requestDataMap["callCostInUSDCents"] = data.CallCostInUsdCents.ValueBigFloat()
    }
    if !data.StartedAt.IsUnknown() && !state.StartedAt.IsUnknown() && !data.StartedAt.Equal(state.StartedAt) {
        var startedatData interface{}
        if err := json.Unmarshal([]byte(data.StartedAt.ValueString()), &startedatData); err == nil {
            requestDataMap["startedAt"] = startedatData
        }
    }
    if !data.EndedAt.IsUnknown() && !state.EndedAt.IsUnknown() && !data.EndedAt.Equal(state.EndedAt) {
        var endedatData interface{}
        if err := json.Unmarshal([]byte(data.EndedAt.ValueString()), &endedatData); err == nil {
            requestDataMap["endedAt"] = endedatData
        }
    }
    if !data.IsAnswered.IsUnknown() && !state.IsAnswered.IsUnknown() && !data.IsAnswered.Equal(state.IsAnswered) {
        requestDataMap["isAnswered"] = data.IsAnswered.ValueBool()
    }

    // Make API call
    httpResp, err := r.client.Put("/incoming-call-log-item/" + data.Id.ValueString() + "", incomingCallLogItemRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update incoming_call_log_item, got error: %s", err))
        return
    }

    // Parse the update response
    var incomingCallLogItemResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &incomingCallLogItemResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse incoming_call_log_item response, got error: %s", err))
        return
    }

    // After successful update, fetch the current state by calling Read with select parameter
    selectParam := map[string]interface{}{
        "projectId": true,
        "incomingCallLogId": true,
        "incomingCallPolicyEscalationRuleId": true,
        "userId": true,
        "userPhoneNumber": true,
        "status": true,
        "statusMessage": true,
        "dialDurationInSeconds": true,
        "callCostInUSDCents": true,
        "startedAt": true,
        "endedAt": true,
        "isAnswered": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect("/incoming-call-log-item/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read incoming_call_log_item after update, got error: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse incoming_call_log_item read response, got error: %s", err))
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
    if obj, ok := dataMap["incomingCallLogId"].(map[string]interface{}); ok {
        // Handle ObjectID and value wrapper type responses (e.g., Version with _type and value fields)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncomingCallLogId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.IncomingCallLogId = types.StringValue(val)
        } else {
            data.IncomingCallLogId = types.StringNull()
        }
    } else if val, ok := dataMap["incomingCallLogId"].(string); ok && val != "" {
        data.IncomingCallLogId = types.StringValue(val)
    } else {
        data.IncomingCallLogId = types.StringNull()
    }
    if obj, ok := dataMap["incomingCallPolicyEscalationRuleId"].(map[string]interface{}); ok {
        // Handle ObjectID and value wrapper type responses (e.g., Version with _type and value fields)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncomingCallPolicyEscalationRuleId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.IncomingCallPolicyEscalationRuleId = types.StringValue(val)
        } else {
            data.IncomingCallPolicyEscalationRuleId = types.StringNull()
        }
    } else if val, ok := dataMap["incomingCallPolicyEscalationRuleId"].(string); ok && val != "" {
        data.IncomingCallPolicyEscalationRuleId = types.StringValue(val)
    } else {
        data.IncomingCallPolicyEscalationRuleId = types.StringNull()
    }
    if obj, ok := dataMap["userId"].(map[string]interface{}); ok {
        // Handle ObjectID and value wrapper type responses (e.g., Version with _type and value fields)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.UserId = types.StringValue(val)
        } else {
            data.UserId = types.StringNull()
        }
    } else if val, ok := dataMap["userId"].(string); ok && val != "" {
        data.UserId = types.StringValue(val)
    } else {
        data.UserId = types.StringNull()
    }
    if val, ok := dataMap["userPhoneNumber"].(map[string]interface{}); ok {
        // Check if this is a value wrapper type (e.g., Version, ObjectID with _type and value fields)
        if _, hasType := val["_type"]; hasType {
            if strVal, ok := val["value"].(string); ok && strVal != "" {
                data.UserPhoneNumber = types.StringValue(strVal)
            } else {
                data.UserPhoneNumber = types.StringNull()
            }
        } else if strVal, ok := val["_id"].(string); ok && strVal != "" {
            // Handle ObjectID type responses
            data.UserPhoneNumber = types.StringValue(strVal)
        } else if strVal, ok := val["value"].(string); ok && strVal != "" {
            // Handle other value wrapper types
            data.UserPhoneNumber = types.StringValue(strVal)
        } else {
            // Fall back to JSON marshalling for truly complex objects
            if jsonBytes, err := json.Marshal(val); err == nil {
                data.UserPhoneNumber = types.StringValue(string(jsonBytes))
            } else {
                data.UserPhoneNumber = types.StringNull()
            }
        }
    } else if val, ok := dataMap["userPhoneNumber"].(string); ok && val != "" {
        data.UserPhoneNumber = types.StringValue(val)
    } else {
        data.UserPhoneNumber = types.StringNull()
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
    if val, ok := dataMap["dialDurationInSeconds"].(float64); ok {
        data.DialDurationInSeconds = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["dialDurationInSeconds"].(int); ok {
        data.DialDurationInSeconds = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["dialDurationInSeconds"].(int64); ok {
        data.DialDurationInSeconds = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["dialDurationInSeconds"] == nil {
        data.DialDurationInSeconds = types.NumberNull()
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
    if val, ok := dataMap["isAnswered"].(bool); ok {
        data.IsAnswered = types.BoolValue(val)
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
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IncomingCallLogItemResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var data IncomingCallLogItemResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Make API call
    _, err := r.client.Delete("/incoming-call-log-item/" + data.Id.ValueString() + "")
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete incoming_call_log_item, got error: %s", err))
        return
    }
}


func (r *IncomingCallLogItemResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Helper method to convert Terraform map to Go interface{}
func (r *IncomingCallLogItemResource) convertTerraformMapToInterface(terraformMap types.Map) interface{} {
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
func (r *IncomingCallLogItemResource) convertTerraformListToInterface(terraformList types.List) interface{} {
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
func (r *IncomingCallLogItemResource) parseJSONField(terraformString types.String) interface{} {
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
