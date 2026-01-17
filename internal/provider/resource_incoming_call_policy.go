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
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/numberdefault"
    "github.com/hashicorp/terraform-plugin-framework/attr"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &IncomingCallPolicyResource{}
var _ resource.ResourceWithImportState = &IncomingCallPolicyResource{}

func NewIncomingCallPolicyResource() resource.Resource {
    return &IncomingCallPolicyResource{}
}

// IncomingCallPolicyResource defines the resource implementation.
type IncomingCallPolicyResource struct {
    client *Client
}

// IncomingCallPolicyResourceModel describes the resource data model.
type IncomingCallPolicyResourceModel struct {
    Id types.String `tfsdk:"id"`
    ProjectId types.String `tfsdk:"project_id"`
    Name types.String `tfsdk:"name"`
    Description types.String `tfsdk:"description"`
    GreetingMessage types.String `tfsdk:"greeting_message"`
    NoAnswerMessage types.String `tfsdk:"no_answer_message"`
    NoOneAvailableMessage types.String `tfsdk:"no_one_available_message"`
    IsEnabled types.Bool `tfsdk:"is_enabled"`
    RepeatPolicyIfNoOneAnswers types.Bool `tfsdk:"repeat_policy_if_no_one_answers"`
    RepeatPolicyIfNoOneAnswersTimes types.Number `tfsdk:"repeat_policy_if_no_one_answers_times"`
    Labels types.List `tfsdk:"labels"`
    ProjectCallSmsConfigId types.String `tfsdk:"project_call_sms_config_id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    Slug types.String `tfsdk:"slug"`
    RoutingPhoneNumber types.String `tfsdk:"routing_phone_number"`
    CallProviderPhoneNumberId types.String `tfsdk:"call_provider_phone_number_id"`
    PhoneNumberCountryCode types.String `tfsdk:"phone_number_country_code"`
    PhoneNumberAreaCode types.String `tfsdk:"phone_number_area_code"`
    PhoneNumberPurchasedAt types.String `tfsdk:"phone_number_purchased_at"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (r *IncomingCallPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_incoming_call_policy"
}

func (r *IncomingCallPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "incoming_call_policy resource",

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
            "name": schema.StringAttribute{
                MarkdownDescription: "Any friendly name of this policy. Permissions - Create: [Project Owner, Project Admin, Project Member, Create Incoming Call Policy], Read: [Project Owner, Project Admin, Project Member, Read Incoming Call Policy], Update: [Project Owner, Project Admin, Project Member, Edit Incoming Call Policy]",
                Required: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Friendly description that will help you remember. Permissions - Create: [Project Owner, Project Admin, Project Member, Create Incoming Call Policy], Read: [Project Owner, Project Admin, Project Member, Read Incoming Call Policy], Update: [Project Owner, Project Admin, Project Member, Edit Incoming Call Policy]",
                Optional: true,
                Computed: true,
            },
            "greeting_message": schema.StringAttribute{
                MarkdownDescription: "Custom TTS greeting message for incoming calls. Permissions - Create: [Project Owner, Project Admin, Project Member, Create Incoming Call Policy], Read: [Project Owner, Project Admin, Project Member, Read Incoming Call Policy], Update: [Project Owner, Project Admin, Project Member, Edit Incoming Call Policy]",
                Optional: true,
                Computed: true,
            },
            "no_answer_message": schema.StringAttribute{
                MarkdownDescription: "Message when escalation is exhausted and no one answers. Permissions - Create: [Project Owner, Project Admin, Project Member, Create Incoming Call Policy], Read: [Project Owner, Project Admin, Project Member, Read Incoming Call Policy], Update: [Project Owner, Project Admin, Project Member, Edit Incoming Call Policy]",
                Optional: true,
                Computed: true,
            },
            "no_one_available_message": schema.StringAttribute{
                MarkdownDescription: "Message when no one is on-call or reachable. Permissions - Create: [Project Owner, Project Admin, Project Member, Create Incoming Call Policy], Read: [Project Owner, Project Admin, Project Member, Read Incoming Call Policy], Update: [Project Owner, Project Admin, Project Member, Edit Incoming Call Policy]",
                Optional: true,
                Computed: true,
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Enable or disable this incoming call policy. Permissions - Create: [Project Owner, Project Admin, Project Member, Create Incoming Call Policy], Read: [Project Owner, Project Admin, Project Member, Read Incoming Call Policy], Update: [Project Owner, Project Admin, Project Member, Edit Incoming Call Policy]",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(true),
            },
            "repeat_policy_if_no_one_answers": schema.BoolAttribute{
                MarkdownDescription: "Restart from first rule if all fail. Permissions - Create: [Project Owner, Project Admin, Project Member, Create Incoming Call Policy], Read: [Project Owner, Project Admin, Project Member, Read Incoming Call Policy], Update: [Project Owner, Project Admin, Project Member, Edit Incoming Call Policy]",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
            },
            "repeat_policy_if_no_one_answers_times": schema.NumberAttribute{
                MarkdownDescription: "Maximum repeat attempts if no one answers. Permissions - Create: [Project Owner, Project Admin, Project Member, Create Incoming Call Policy], Read: [Project Owner, Project Admin, Project Member, Read Incoming Call Policy], Update: [Project Owner, Project Admin, Project Member, Edit Incoming Call Policy]",
                Optional: true,
                Computed: true,
                Default: numberdefault.StaticBigFloat(big.NewFloat(1)),
            },
            "labels": schema.ListAttribute{
                MarkdownDescription: "Relation to Labels Array where this object is categorized in.. Permissions - Create: [Project Owner, Project Admin, Project Member, Create Incoming Call Policy], Read: [Project Owner, Project Admin, Project Member, Read Incoming Call Policy], Update: [Project Owner, Project Admin, Project Member, Edit Incoming Call Policy]",
                Optional: true,
                Computed: true,
                ElementType: types.StringType,
            },
            "project_call_sms_config_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
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
            "slug": schema.StringAttribute{
                MarkdownDescription: "Friendly globally unique name for your object. Permissions - Create: [Project Owner, Project Admin, Project Member, Create Incoming Call Policy], Read: [Project Owner, Project Admin, Project Member, Read Incoming Call Policy], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "routing_phone_number": schema.StringAttribute{
                MarkdownDescription: "Phone object",
                Computed: true,
            },
            "call_provider_phone_number_id": schema.StringAttribute{
                MarkdownDescription: "The call provider's ID for the phone number (e.g., Twilio SID). Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Incoming Call Policy], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "phone_number_country_code": schema.StringAttribute{
                MarkdownDescription: "Country code of the phone number (US, GB, etc.). Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Incoming Call Policy], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "phone_number_area_code": schema.StringAttribute{
                MarkdownDescription: "Area code of the phone number. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Incoming Call Policy], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "phone_number_purchased_at": schema.StringAttribute{
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

func (r *IncomingCallPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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


func (r *IncomingCallPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var data IncomingCallPolicyResourceModel

    // Read Terraform plan data into the model
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }



    // Create API request body
    incomingCallPolicyRequest := map[string]interface{}{
        "data": map[string]interface{}{
        "projectId": data.ProjectId.ValueString(),
        "name": data.Name.ValueString(),
        "description": data.Description.ValueString(),
        "greetingMessage": data.GreetingMessage.ValueString(),
        "noAnswerMessage": data.NoAnswerMessage.ValueString(),
        "noOneAvailableMessage": data.NoOneAvailableMessage.ValueString(),
        "isEnabled": data.IsEnabled.ValueBool(),
        "repeatPolicyIfNoOneAnswers": data.RepeatPolicyIfNoOneAnswers.ValueBool(),
        "repeatPolicyIfNoOneAnswersTimes": data.RepeatPolicyIfNoOneAnswersTimes.ValueBigFloat(),
        "labels": r.convertTerraformListToInterface(data.Labels),
        "projectCallSMSConfigId": data.ProjectCallSmsConfigId.ValueString(),
        },
    }

    // Make API call
    httpResp, err := r.client.Post("/incoming-call-policy", incomingCallPolicyRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create incoming_call_policy, got error: %s", err))
        return
    }

    var incomingCallPolicyResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &incomingCallPolicyResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse incoming_call_policy response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := incomingCallPolicyResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = incomingCallPolicyResponse
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
    if obj, ok := dataMap["name"].(map[string]interface{}); ok {
        // Handle ObjectID and value wrapper type responses (e.g., Version with _type and value fields)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.Name = types.StringValue(val)
        } else {
            data.Name = types.StringNull()
        }
    } else if val, ok := dataMap["name"].(string); ok && val != "" {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if obj, ok := dataMap["description"].(map[string]interface{}); ok {
        // Handle ObjectID and value wrapper type responses (e.g., Version with _type and value fields)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.Description = types.StringValue(val)
        } else {
            data.Description = types.StringNull()
        }
    } else if val, ok := dataMap["description"].(string); ok && val != "" {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if obj, ok := dataMap["greetingMessage"].(map[string]interface{}); ok {
        // Handle ObjectID and value wrapper type responses (e.g., Version with _type and value fields)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.GreetingMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.GreetingMessage = types.StringValue(val)
        } else {
            data.GreetingMessage = types.StringNull()
        }
    } else if val, ok := dataMap["greetingMessage"].(string); ok && val != "" {
        data.GreetingMessage = types.StringValue(val)
    } else {
        data.GreetingMessage = types.StringNull()
    }
    if obj, ok := dataMap["noAnswerMessage"].(map[string]interface{}); ok {
        // Handle ObjectID and value wrapper type responses (e.g., Version with _type and value fields)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.NoAnswerMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.NoAnswerMessage = types.StringValue(val)
        } else {
            data.NoAnswerMessage = types.StringNull()
        }
    } else if val, ok := dataMap["noAnswerMessage"].(string); ok && val != "" {
        data.NoAnswerMessage = types.StringValue(val)
    } else {
        data.NoAnswerMessage = types.StringNull()
    }
    if obj, ok := dataMap["noOneAvailableMessage"].(map[string]interface{}); ok {
        // Handle ObjectID and value wrapper type responses (e.g., Version with _type and value fields)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.NoOneAvailableMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.NoOneAvailableMessage = types.StringValue(val)
        } else {
            data.NoOneAvailableMessage = types.StringNull()
        }
    } else if val, ok := dataMap["noOneAvailableMessage"].(string); ok && val != "" {
        data.NoOneAvailableMessage = types.StringValue(val)
    } else {
        data.NoOneAvailableMessage = types.StringNull()
    }
    if val, ok := dataMap["isEnabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    }
    if val, ok := dataMap["repeatPolicyIfNoOneAnswers"].(bool); ok {
        data.RepeatPolicyIfNoOneAnswers = types.BoolValue(val)
    }
    if val, ok := dataMap["repeatPolicyIfNoOneAnswersTimes"].(float64); ok {
        data.RepeatPolicyIfNoOneAnswersTimes = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["repeatPolicyIfNoOneAnswersTimes"].(int); ok {
        data.RepeatPolicyIfNoOneAnswersTimes = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["repeatPolicyIfNoOneAnswersTimes"].(int64); ok {
        data.RepeatPolicyIfNoOneAnswersTimes = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["repeatPolicyIfNoOneAnswersTimes"] == nil {
        data.RepeatPolicyIfNoOneAnswersTimes = types.NumberNull()
    }
    if val, ok := dataMap["labels"].([]interface{}); ok {
        // Convert API response list to Terraform list
        var listItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    listItems = append(listItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    listItems = append(listItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        listItems = append(listItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                listItems = append(listItems, types.StringValue(str))
            }
        }
        data.Labels = types.ListValueMust(types.StringType, listItems)
    } else {
        // For lists, always use empty list instead of null to match default values
        data.Labels = types.ListValueMust(types.StringType, []attr.Value{})
    }
    if obj, ok := dataMap["projectCallSMSConfigId"].(map[string]interface{}); ok {
        // Handle ObjectID and value wrapper type responses (e.g., Version with _type and value fields)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ProjectCallSmsConfigId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.ProjectCallSmsConfigId = types.StringValue(val)
        } else {
            data.ProjectCallSmsConfigId = types.StringNull()
        }
    } else if val, ok := dataMap["projectCallSMSConfigId"].(string); ok && val != "" {
        data.ProjectCallSmsConfigId = types.StringValue(val)
    } else {
        data.ProjectCallSmsConfigId = types.StringNull()
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
    if obj, ok := dataMap["slug"].(map[string]interface{}); ok {
        // Handle ObjectID and value wrapper type responses (e.g., Version with _type and value fields)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.Slug = types.StringValue(val)
        } else {
            data.Slug = types.StringNull()
        }
    } else if val, ok := dataMap["slug"].(string); ok && val != "" {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
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
    if obj, ok := dataMap["callProviderPhoneNumberId"].(map[string]interface{}); ok {
        // Handle ObjectID and value wrapper type responses (e.g., Version with _type and value fields)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CallProviderPhoneNumberId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.CallProviderPhoneNumberId = types.StringValue(val)
        } else {
            data.CallProviderPhoneNumberId = types.StringNull()
        }
    } else if val, ok := dataMap["callProviderPhoneNumberId"].(string); ok && val != "" {
        data.CallProviderPhoneNumberId = types.StringValue(val)
    } else {
        data.CallProviderPhoneNumberId = types.StringNull()
    }
    if obj, ok := dataMap["phoneNumberCountryCode"].(map[string]interface{}); ok {
        // Handle ObjectID and value wrapper type responses (e.g., Version with _type and value fields)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PhoneNumberCountryCode = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.PhoneNumberCountryCode = types.StringValue(val)
        } else {
            data.PhoneNumberCountryCode = types.StringNull()
        }
    } else if val, ok := dataMap["phoneNumberCountryCode"].(string); ok && val != "" {
        data.PhoneNumberCountryCode = types.StringValue(val)
    } else {
        data.PhoneNumberCountryCode = types.StringNull()
    }
    if obj, ok := dataMap["phoneNumberAreaCode"].(map[string]interface{}); ok {
        // Handle ObjectID and value wrapper type responses (e.g., Version with _type and value fields)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PhoneNumberAreaCode = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.PhoneNumberAreaCode = types.StringValue(val)
        } else {
            data.PhoneNumberAreaCode = types.StringNull()
        }
    } else if val, ok := dataMap["phoneNumberAreaCode"].(string); ok && val != "" {
        data.PhoneNumberAreaCode = types.StringValue(val)
    } else {
        data.PhoneNumberAreaCode = types.StringNull()
    }
    if val, ok := dataMap["phoneNumberPurchasedAt"].(map[string]interface{}); ok {
        // Check if this is a value wrapper type (e.g., Version, ObjectID with _type and value fields)
        if _, hasType := val["_type"]; hasType {
            if strVal, ok := val["value"].(string); ok && strVal != "" {
                data.PhoneNumberPurchasedAt = types.StringValue(strVal)
            } else {
                data.PhoneNumberPurchasedAt = types.StringNull()
            }
        } else if strVal, ok := val["_id"].(string); ok && strVal != "" {
            // Handle ObjectID type responses
            data.PhoneNumberPurchasedAt = types.StringValue(strVal)
        } else if strVal, ok := val["value"].(string); ok && strVal != "" {
            // Handle other value wrapper types
            data.PhoneNumberPurchasedAt = types.StringValue(strVal)
        } else {
            // Fall back to JSON marshalling for truly complex objects
            if jsonBytes, err := json.Marshal(val); err == nil {
                data.PhoneNumberPurchasedAt = types.StringValue(string(jsonBytes))
            } else {
                data.PhoneNumberPurchasedAt = types.StringNull()
            }
        }
    } else if val, ok := dataMap["phoneNumberPurchasedAt"].(string); ok && val != "" {
        data.PhoneNumberPurchasedAt = types.StringValue(val)
    } else {
        data.PhoneNumberPurchasedAt = types.StringNull()
    }
    if obj, ok := dataMap["createdByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID and value wrapper type responses (e.g., Version with _type and value fields)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.CreatedByUserId = types.StringValue(val)
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["createdByUserId"].(string); ok && val != "" {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
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

func (r *IncomingCallPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var data IncomingCallPolicyResourceModel

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
        "greetingMessage": true,
        "noAnswerMessage": true,
        "noOneAvailableMessage": true,
        "isEnabled": true,
        "repeatPolicyIfNoOneAnswers": true,
        "repeatPolicyIfNoOneAnswersTimes": true,
        "labels": true,
        "projectCallSMSConfigId": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "slug": true,
        "routingPhoneNumber": true,
        "callProviderPhoneNumberId": true,
        "phoneNumberCountryCode": true,
        "phoneNumberAreaCode": true,
        "phoneNumberPurchasedAt": true,
        "createdByUserId": true,
        "_id": true,
    }

    // Make API call with select parameter
    httpResp, err := r.client.PostWithSelect("/incoming-call-policy/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read incoming_call_policy, got error: %s", err))
        return
    }

    if httpResp.StatusCode == http.StatusNotFound {
        resp.State.RemoveResource(ctx)
        return
    }

    var incomingCallPolicyResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &incomingCallPolicyResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse incoming_call_policy response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := incomingCallPolicyResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = incomingCallPolicyResponse
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
    if obj, ok := dataMap["name"].(map[string]interface{}); ok {
        // Handle ObjectID and value wrapper type responses (e.g., Version with _type and value fields)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.Name = types.StringValue(val)
        } else {
            data.Name = types.StringNull()
        }
    } else if val, ok := dataMap["name"].(string); ok && val != "" {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if obj, ok := dataMap["description"].(map[string]interface{}); ok {
        // Handle ObjectID and value wrapper type responses (e.g., Version with _type and value fields)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.Description = types.StringValue(val)
        } else {
            data.Description = types.StringNull()
        }
    } else if val, ok := dataMap["description"].(string); ok && val != "" {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if obj, ok := dataMap["greetingMessage"].(map[string]interface{}); ok {
        // Handle ObjectID and value wrapper type responses (e.g., Version with _type and value fields)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.GreetingMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.GreetingMessage = types.StringValue(val)
        } else {
            data.GreetingMessage = types.StringNull()
        }
    } else if val, ok := dataMap["greetingMessage"].(string); ok && val != "" {
        data.GreetingMessage = types.StringValue(val)
    } else {
        data.GreetingMessage = types.StringNull()
    }
    if obj, ok := dataMap["noAnswerMessage"].(map[string]interface{}); ok {
        // Handle ObjectID and value wrapper type responses (e.g., Version with _type and value fields)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.NoAnswerMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.NoAnswerMessage = types.StringValue(val)
        } else {
            data.NoAnswerMessage = types.StringNull()
        }
    } else if val, ok := dataMap["noAnswerMessage"].(string); ok && val != "" {
        data.NoAnswerMessage = types.StringValue(val)
    } else {
        data.NoAnswerMessage = types.StringNull()
    }
    if obj, ok := dataMap["noOneAvailableMessage"].(map[string]interface{}); ok {
        // Handle ObjectID and value wrapper type responses (e.g., Version with _type and value fields)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.NoOneAvailableMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.NoOneAvailableMessage = types.StringValue(val)
        } else {
            data.NoOneAvailableMessage = types.StringNull()
        }
    } else if val, ok := dataMap["noOneAvailableMessage"].(string); ok && val != "" {
        data.NoOneAvailableMessage = types.StringValue(val)
    } else {
        data.NoOneAvailableMessage = types.StringNull()
    }
    if val, ok := dataMap["isEnabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    }
    if val, ok := dataMap["repeatPolicyIfNoOneAnswers"].(bool); ok {
        data.RepeatPolicyIfNoOneAnswers = types.BoolValue(val)
    }
    if val, ok := dataMap["repeatPolicyIfNoOneAnswersTimes"].(float64); ok {
        data.RepeatPolicyIfNoOneAnswersTimes = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["repeatPolicyIfNoOneAnswersTimes"].(int); ok {
        data.RepeatPolicyIfNoOneAnswersTimes = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["repeatPolicyIfNoOneAnswersTimes"].(int64); ok {
        data.RepeatPolicyIfNoOneAnswersTimes = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["repeatPolicyIfNoOneAnswersTimes"] == nil {
        data.RepeatPolicyIfNoOneAnswersTimes = types.NumberNull()
    }
    if val, ok := dataMap["labels"].([]interface{}); ok {
        // Convert API response list to Terraform list
        var listItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    listItems = append(listItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    listItems = append(listItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        listItems = append(listItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                listItems = append(listItems, types.StringValue(str))
            }
        }
        data.Labels = types.ListValueMust(types.StringType, listItems)
    } else {
        // For lists, always use empty list instead of null to match default values
        data.Labels = types.ListValueMust(types.StringType, []attr.Value{})
    }
    if obj, ok := dataMap["projectCallSMSConfigId"].(map[string]interface{}); ok {
        // Handle ObjectID and value wrapper type responses (e.g., Version with _type and value fields)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ProjectCallSmsConfigId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.ProjectCallSmsConfigId = types.StringValue(val)
        } else {
            data.ProjectCallSmsConfigId = types.StringNull()
        }
    } else if val, ok := dataMap["projectCallSMSConfigId"].(string); ok && val != "" {
        data.ProjectCallSmsConfigId = types.StringValue(val)
    } else {
        data.ProjectCallSmsConfigId = types.StringNull()
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
    if obj, ok := dataMap["slug"].(map[string]interface{}); ok {
        // Handle ObjectID and value wrapper type responses (e.g., Version with _type and value fields)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.Slug = types.StringValue(val)
        } else {
            data.Slug = types.StringNull()
        }
    } else if val, ok := dataMap["slug"].(string); ok && val != "" {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
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
    if obj, ok := dataMap["callProviderPhoneNumberId"].(map[string]interface{}); ok {
        // Handle ObjectID and value wrapper type responses (e.g., Version with _type and value fields)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CallProviderPhoneNumberId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.CallProviderPhoneNumberId = types.StringValue(val)
        } else {
            data.CallProviderPhoneNumberId = types.StringNull()
        }
    } else if val, ok := dataMap["callProviderPhoneNumberId"].(string); ok && val != "" {
        data.CallProviderPhoneNumberId = types.StringValue(val)
    } else {
        data.CallProviderPhoneNumberId = types.StringNull()
    }
    if obj, ok := dataMap["phoneNumberCountryCode"].(map[string]interface{}); ok {
        // Handle ObjectID and value wrapper type responses (e.g., Version with _type and value fields)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PhoneNumberCountryCode = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.PhoneNumberCountryCode = types.StringValue(val)
        } else {
            data.PhoneNumberCountryCode = types.StringNull()
        }
    } else if val, ok := dataMap["phoneNumberCountryCode"].(string); ok && val != "" {
        data.PhoneNumberCountryCode = types.StringValue(val)
    } else {
        data.PhoneNumberCountryCode = types.StringNull()
    }
    if obj, ok := dataMap["phoneNumberAreaCode"].(map[string]interface{}); ok {
        // Handle ObjectID and value wrapper type responses (e.g., Version with _type and value fields)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PhoneNumberAreaCode = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.PhoneNumberAreaCode = types.StringValue(val)
        } else {
            data.PhoneNumberAreaCode = types.StringNull()
        }
    } else if val, ok := dataMap["phoneNumberAreaCode"].(string); ok && val != "" {
        data.PhoneNumberAreaCode = types.StringValue(val)
    } else {
        data.PhoneNumberAreaCode = types.StringNull()
    }
    if val, ok := dataMap["phoneNumberPurchasedAt"].(map[string]interface{}); ok {
        // Check if this is a value wrapper type (e.g., Version, ObjectID with _type and value fields)
        if _, hasType := val["_type"]; hasType {
            if strVal, ok := val["value"].(string); ok && strVal != "" {
                data.PhoneNumberPurchasedAt = types.StringValue(strVal)
            } else {
                data.PhoneNumberPurchasedAt = types.StringNull()
            }
        } else if strVal, ok := val["_id"].(string); ok && strVal != "" {
            // Handle ObjectID type responses
            data.PhoneNumberPurchasedAt = types.StringValue(strVal)
        } else if strVal, ok := val["value"].(string); ok && strVal != "" {
            // Handle other value wrapper types
            data.PhoneNumberPurchasedAt = types.StringValue(strVal)
        } else {
            // Fall back to JSON marshalling for truly complex objects
            if jsonBytes, err := json.Marshal(val); err == nil {
                data.PhoneNumberPurchasedAt = types.StringValue(string(jsonBytes))
            } else {
                data.PhoneNumberPurchasedAt = types.StringNull()
            }
        }
    } else if val, ok := dataMap["phoneNumberPurchasedAt"].(string); ok && val != "" {
        data.PhoneNumberPurchasedAt = types.StringValue(val)
    } else {
        data.PhoneNumberPurchasedAt = types.StringNull()
    }
    if obj, ok := dataMap["createdByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID and value wrapper type responses (e.g., Version with _type and value fields)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.CreatedByUserId = types.StringValue(val)
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["createdByUserId"].(string); ok && val != "" {
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

func (r *IncomingCallPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var data IncomingCallPolicyResourceModel
    var state IncomingCallPolicyResourceModel

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
    incomingCallPolicyRequest := map[string]interface{}{
        "data": map[string]interface{}{},
    }
    requestDataMap := incomingCallPolicyRequest["data"].(map[string]interface{})

    if !data.Name.IsUnknown() && !state.Name.IsUnknown() && !data.Name.Equal(state.Name) {
        requestDataMap["name"] = data.Name.ValueString()
    }
    if !data.Description.IsUnknown() && !state.Description.IsUnknown() && !data.Description.Equal(state.Description) {
        requestDataMap["description"] = data.Description.ValueString()
    }
    if !data.GreetingMessage.IsUnknown() && !state.GreetingMessage.IsUnknown() && !data.GreetingMessage.Equal(state.GreetingMessage) {
        requestDataMap["greetingMessage"] = data.GreetingMessage.ValueString()
    }
    if !data.NoAnswerMessage.IsUnknown() && !state.NoAnswerMessage.IsUnknown() && !data.NoAnswerMessage.Equal(state.NoAnswerMessage) {
        requestDataMap["noAnswerMessage"] = data.NoAnswerMessage.ValueString()
    }
    if !data.NoOneAvailableMessage.IsUnknown() && !state.NoOneAvailableMessage.IsUnknown() && !data.NoOneAvailableMessage.Equal(state.NoOneAvailableMessage) {
        requestDataMap["noOneAvailableMessage"] = data.NoOneAvailableMessage.ValueString()
    }
    if !data.IsEnabled.IsUnknown() && !state.IsEnabled.IsUnknown() && !data.IsEnabled.Equal(state.IsEnabled) {
        requestDataMap["isEnabled"] = data.IsEnabled.ValueBool()
    }
    if !data.RepeatPolicyIfNoOneAnswers.IsUnknown() && !state.RepeatPolicyIfNoOneAnswers.IsUnknown() && !data.RepeatPolicyIfNoOneAnswers.Equal(state.RepeatPolicyIfNoOneAnswers) {
        requestDataMap["repeatPolicyIfNoOneAnswers"] = data.RepeatPolicyIfNoOneAnswers.ValueBool()
    }
    if !data.RepeatPolicyIfNoOneAnswersTimes.IsUnknown() && !state.RepeatPolicyIfNoOneAnswersTimes.IsUnknown() && !data.RepeatPolicyIfNoOneAnswersTimes.Equal(state.RepeatPolicyIfNoOneAnswersTimes) {
        requestDataMap["repeatPolicyIfNoOneAnswersTimes"] = data.RepeatPolicyIfNoOneAnswersTimes.ValueBigFloat()
    }
    if !data.Labels.IsUnknown() && !state.Labels.IsUnknown() && !data.Labels.Equal(state.Labels) {
        requestDataMap["labels"] = r.convertTerraformListToInterface(data.Labels)
    }
    if !data.ProjectCallSmsConfigId.IsUnknown() && !state.ProjectCallSmsConfigId.IsUnknown() && !data.ProjectCallSmsConfigId.Equal(state.ProjectCallSmsConfigId) {
        requestDataMap["projectCallSMSConfigId"] = data.ProjectCallSmsConfigId.ValueString()
    }

    // Make API call
    httpResp, err := r.client.Put("/incoming-call-policy/" + data.Id.ValueString() + "", incomingCallPolicyRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update incoming_call_policy, got error: %s", err))
        return
    }

    // Parse the update response
    var incomingCallPolicyResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &incomingCallPolicyResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse incoming_call_policy response, got error: %s", err))
        return
    }

    // After successful update, fetch the current state by calling Read with select parameter
    selectParam := map[string]interface{}{
        "projectId": true,
        "name": true,
        "description": true,
        "greetingMessage": true,
        "noAnswerMessage": true,
        "noOneAvailableMessage": true,
        "isEnabled": true,
        "repeatPolicyIfNoOneAnswers": true,
        "repeatPolicyIfNoOneAnswersTimes": true,
        "labels": true,
        "projectCallSMSConfigId": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "slug": true,
        "routingPhoneNumber": true,
        "callProviderPhoneNumberId": true,
        "phoneNumberCountryCode": true,
        "phoneNumberAreaCode": true,
        "phoneNumberPurchasedAt": true,
        "createdByUserId": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect("/incoming-call-policy/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read incoming_call_policy after update, got error: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse incoming_call_policy read response, got error: %s", err))
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
    if obj, ok := dataMap["name"].(map[string]interface{}); ok {
        // Handle ObjectID and value wrapper type responses (e.g., Version with _type and value fields)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.Name = types.StringValue(val)
        } else {
            data.Name = types.StringNull()
        }
    } else if val, ok := dataMap["name"].(string); ok && val != "" {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if obj, ok := dataMap["description"].(map[string]interface{}); ok {
        // Handle ObjectID and value wrapper type responses (e.g., Version with _type and value fields)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.Description = types.StringValue(val)
        } else {
            data.Description = types.StringNull()
        }
    } else if val, ok := dataMap["description"].(string); ok && val != "" {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if obj, ok := dataMap["greetingMessage"].(map[string]interface{}); ok {
        // Handle ObjectID and value wrapper type responses (e.g., Version with _type and value fields)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.GreetingMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.GreetingMessage = types.StringValue(val)
        } else {
            data.GreetingMessage = types.StringNull()
        }
    } else if val, ok := dataMap["greetingMessage"].(string); ok && val != "" {
        data.GreetingMessage = types.StringValue(val)
    } else {
        data.GreetingMessage = types.StringNull()
    }
    if obj, ok := dataMap["noAnswerMessage"].(map[string]interface{}); ok {
        // Handle ObjectID and value wrapper type responses (e.g., Version with _type and value fields)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.NoAnswerMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.NoAnswerMessage = types.StringValue(val)
        } else {
            data.NoAnswerMessage = types.StringNull()
        }
    } else if val, ok := dataMap["noAnswerMessage"].(string); ok && val != "" {
        data.NoAnswerMessage = types.StringValue(val)
    } else {
        data.NoAnswerMessage = types.StringNull()
    }
    if obj, ok := dataMap["noOneAvailableMessage"].(map[string]interface{}); ok {
        // Handle ObjectID and value wrapper type responses (e.g., Version with _type and value fields)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.NoOneAvailableMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.NoOneAvailableMessage = types.StringValue(val)
        } else {
            data.NoOneAvailableMessage = types.StringNull()
        }
    } else if val, ok := dataMap["noOneAvailableMessage"].(string); ok && val != "" {
        data.NoOneAvailableMessage = types.StringValue(val)
    } else {
        data.NoOneAvailableMessage = types.StringNull()
    }
    if val, ok := dataMap["isEnabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    }
    if val, ok := dataMap["repeatPolicyIfNoOneAnswers"].(bool); ok {
        data.RepeatPolicyIfNoOneAnswers = types.BoolValue(val)
    }
    if val, ok := dataMap["repeatPolicyIfNoOneAnswersTimes"].(float64); ok {
        data.RepeatPolicyIfNoOneAnswersTimes = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["repeatPolicyIfNoOneAnswersTimes"].(int); ok {
        data.RepeatPolicyIfNoOneAnswersTimes = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["repeatPolicyIfNoOneAnswersTimes"].(int64); ok {
        data.RepeatPolicyIfNoOneAnswersTimes = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["repeatPolicyIfNoOneAnswersTimes"] == nil {
        data.RepeatPolicyIfNoOneAnswersTimes = types.NumberNull()
    }
    if val, ok := dataMap["labels"].([]interface{}); ok {
        // Convert API response list to Terraform list
        var listItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    listItems = append(listItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    listItems = append(listItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        listItems = append(listItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                listItems = append(listItems, types.StringValue(str))
            }
        }
        data.Labels = types.ListValueMust(types.StringType, listItems)
    } else {
        // For lists, always use empty list instead of null to match default values
        data.Labels = types.ListValueMust(types.StringType, []attr.Value{})
    }
    if obj, ok := dataMap["projectCallSMSConfigId"].(map[string]interface{}); ok {
        // Handle ObjectID and value wrapper type responses (e.g., Version with _type and value fields)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ProjectCallSmsConfigId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.ProjectCallSmsConfigId = types.StringValue(val)
        } else {
            data.ProjectCallSmsConfigId = types.StringNull()
        }
    } else if val, ok := dataMap["projectCallSMSConfigId"].(string); ok && val != "" {
        data.ProjectCallSmsConfigId = types.StringValue(val)
    } else {
        data.ProjectCallSmsConfigId = types.StringNull()
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
    if obj, ok := dataMap["slug"].(map[string]interface{}); ok {
        // Handle ObjectID and value wrapper type responses (e.g., Version with _type and value fields)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.Slug = types.StringValue(val)
        } else {
            data.Slug = types.StringNull()
        }
    } else if val, ok := dataMap["slug"].(string); ok && val != "" {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
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
    if obj, ok := dataMap["callProviderPhoneNumberId"].(map[string]interface{}); ok {
        // Handle ObjectID and value wrapper type responses (e.g., Version with _type and value fields)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CallProviderPhoneNumberId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.CallProviderPhoneNumberId = types.StringValue(val)
        } else {
            data.CallProviderPhoneNumberId = types.StringNull()
        }
    } else if val, ok := dataMap["callProviderPhoneNumberId"].(string); ok && val != "" {
        data.CallProviderPhoneNumberId = types.StringValue(val)
    } else {
        data.CallProviderPhoneNumberId = types.StringNull()
    }
    if obj, ok := dataMap["phoneNumberCountryCode"].(map[string]interface{}); ok {
        // Handle ObjectID and value wrapper type responses (e.g., Version with _type and value fields)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PhoneNumberCountryCode = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.PhoneNumberCountryCode = types.StringValue(val)
        } else {
            data.PhoneNumberCountryCode = types.StringNull()
        }
    } else if val, ok := dataMap["phoneNumberCountryCode"].(string); ok && val != "" {
        data.PhoneNumberCountryCode = types.StringValue(val)
    } else {
        data.PhoneNumberCountryCode = types.StringNull()
    }
    if obj, ok := dataMap["phoneNumberAreaCode"].(map[string]interface{}); ok {
        // Handle ObjectID and value wrapper type responses (e.g., Version with _type and value fields)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PhoneNumberAreaCode = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.PhoneNumberAreaCode = types.StringValue(val)
        } else {
            data.PhoneNumberAreaCode = types.StringNull()
        }
    } else if val, ok := dataMap["phoneNumberAreaCode"].(string); ok && val != "" {
        data.PhoneNumberAreaCode = types.StringValue(val)
    } else {
        data.PhoneNumberAreaCode = types.StringNull()
    }
    if val, ok := dataMap["phoneNumberPurchasedAt"].(map[string]interface{}); ok {
        // Check if this is a value wrapper type (e.g., Version, ObjectID with _type and value fields)
        if _, hasType := val["_type"]; hasType {
            if strVal, ok := val["value"].(string); ok && strVal != "" {
                data.PhoneNumberPurchasedAt = types.StringValue(strVal)
            } else {
                data.PhoneNumberPurchasedAt = types.StringNull()
            }
        } else if strVal, ok := val["_id"].(string); ok && strVal != "" {
            // Handle ObjectID type responses
            data.PhoneNumberPurchasedAt = types.StringValue(strVal)
        } else if strVal, ok := val["value"].(string); ok && strVal != "" {
            // Handle other value wrapper types
            data.PhoneNumberPurchasedAt = types.StringValue(strVal)
        } else {
            // Fall back to JSON marshalling for truly complex objects
            if jsonBytes, err := json.Marshal(val); err == nil {
                data.PhoneNumberPurchasedAt = types.StringValue(string(jsonBytes))
            } else {
                data.PhoneNumberPurchasedAt = types.StringNull()
            }
        }
    } else if val, ok := dataMap["phoneNumberPurchasedAt"].(string); ok && val != "" {
        data.PhoneNumberPurchasedAt = types.StringValue(val)
    } else {
        data.PhoneNumberPurchasedAt = types.StringNull()
    }
    if obj, ok := dataMap["createdByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID and value wrapper type responses (e.g., Version with _type and value fields)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.CreatedByUserId = types.StringValue(val)
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["createdByUserId"].(string); ok && val != "" {
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

func (r *IncomingCallPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var data IncomingCallPolicyResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Make API call
    _, err := r.client.Delete("/incoming-call-policy/" + data.Id.ValueString() + "")
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete incoming_call_policy, got error: %s", err))
        return
    }
}


func (r *IncomingCallPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Helper method to convert Terraform map to Go interface{}
func (r *IncomingCallPolicyResource) convertTerraformMapToInterface(terraformMap types.Map) interface{} {
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
func (r *IncomingCallPolicyResource) convertTerraformListToInterface(terraformList types.List) interface{} {
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
func (r *IncomingCallPolicyResource) parseJSONField(terraformString types.String) interface{} {
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
