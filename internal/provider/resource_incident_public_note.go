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
    "github.com/hashicorp/terraform-plugin-framework/attr"
    "sort"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &IncidentPublicNoteResource{}
var _ resource.ResourceWithImportState = &IncidentPublicNoteResource{}

func NewIncidentPublicNoteResource() resource.Resource {
    return &IncidentPublicNoteResource{}
}

// IncidentPublicNoteResource defines the resource implementation.
type IncidentPublicNoteResource struct {
    client *Client
}

// IncidentPublicNoteResourceModel describes the resource data model.
type IncidentPublicNoteResourceModel struct {
    Id types.String `tfsdk:"id"`
    ProjectId types.String `tfsdk:"project_id"`
    IncidentId types.String `tfsdk:"incident_id"`
    Note types.String `tfsdk:"note"`
    Attachments types.List `tfsdk:"attachments"`
    SubscriberNotificationStatusMessage types.String `tfsdk:"subscriber_notification_status_message"`
    ShouldStatusPageSubscribersBeNotifiedOnNoteCreated types.Bool `tfsdk:"should_status_page_subscribers_be_notified_on_note_created"`
    PostedAt types.String `tfsdk:"posted_at"`
    PostedFromSlackMessageId types.String `tfsdk:"posted_from_slack_message_id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    SubscriberNotificationStatusOnNoteCreated types.String `tfsdk:"subscriber_notification_status_on_note_created"`
    IsOwnerNotified types.Bool `tfsdk:"is_owner_notified"`
}

func (r *IncidentPublicNoteResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_incident_public_note"
}

func (r *IncidentPublicNoteResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "incident_public_note resource",

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
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "incident_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Required: true,
            },
            "note": schema.StringAttribute{
                MarkdownDescription: "Notes in markdown. Permissions - Create: [Project Owner, Project Admin, Project Member, Create Incident Status Page Note], Read: [Project Owner, Project Admin, Project Member, Read Incident Status Page Note], Update: [Project Owner, Project Admin, Project Member, Edit Incident Status Page Note]",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "attachments": schema.ListAttribute{
                MarkdownDescription: "Files attached to this note. Permissions - Create: [Project Owner, Project Admin, Project Member, Create Incident Status Page Note], Read: [Project Owner, Project Admin, Project Member, Read Incident Status Page Note], Update: [Project Owner, Project Admin, Project Member, Edit Incident Status Page Note]",
                Optional: true,
                Computed: true,
                ElementType: types.StringType,
                PlanModifiers: []planmodifier.List{
                    listplanmodifier.UseStateForUnknown(),
                },
            },
            "subscriber_notification_status_message": schema.StringAttribute{
                MarkdownDescription: "Status message for subscriber notifications - includes success messages, failure reasons, or skip reasons. Permissions - Create: [Project Owner, Project Admin, Project Member, Create Incident Status Page Note], Read: [Project Owner, Project Admin, Project Member, Read Incident Status Page Note], Update: [Project Owner, Project Admin, Project Member, Edit Incident Status Page Note]",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "should_status_page_subscribers_be_notified_on_note_created": schema.BoolAttribute{
                MarkdownDescription: "Should subscribers be notified about this note?. Permissions - Create: [Project Owner, Project Admin, Project Member, Create Incident Status Page Note], Read: [Project Owner, Project Admin, Project Member, Read Incident Status Page Note], Update: [No access - you don't have permission for this operation]",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(true),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "posted_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "posted_from_slack_message_id": schema.StringAttribute{
                MarkdownDescription: "Unique identifier for the Slack message this note was created from (channel_id:message_ts). Used to prevent duplicate notes when multiple users react to the same message.. Permissions - Create: [Project Owner, Project Admin, Project Member, Create Incident Status Page Note], Read: [Project Owner, Project Admin, Project Member, Read Incident Status Page Note], Update: [No access - you don't have permission for this operation]",
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
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "subscriber_notification_status_on_note_created": schema.StringAttribute{
                MarkdownDescription: "Status of notification sent to subscribers about this note. Permissions - Create: [Project Owner, Project Admin, Project Member, Create Incident Status Page Note], Read: [Project Owner, Project Admin, Project Member, Read Incident Status Page Note], Update: [Project Owner, Project Admin, Project Member, Edit Incident Status Page Note]",
                Computed: true,
            },
            "is_owner_notified": schema.BoolAttribute{
                MarkdownDescription: "Are owners notified of this resource ownership?. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Incident Status Page Note], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
        },
    }
}

func (r *IncidentPublicNoteResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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


func (r *IncidentPublicNoteResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var data IncidentPublicNoteResourceModel

    // Read Terraform plan data into the model
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }



    // Create API request body
    incidentPublicNoteRequest := map[string]interface{}{
        "data": map[string]interface{}{
        "projectId": data.ProjectId.ValueString(),
        "incidentId": data.IncidentId.ValueString(),
        "note": data.Note.ValueString(),
        "attachments": r.convertTerraformListToInterface(data.Attachments),
        "subscriberNotificationStatusMessage": data.SubscriberNotificationStatusMessage.ValueString(),
        "shouldStatusPageSubscribersBeNotifiedOnNoteCreated": data.ShouldStatusPageSubscribersBeNotifiedOnNoteCreated.ValueBool(),
        "postedAt": r.parseJSONField(data.PostedAt),
        "postedFromSlackMessageId": data.PostedFromSlackMessageId.ValueString(),
        },
    }

    // Make API call
    httpResp, err := r.client.Post("/incident-public-note", incidentPublicNoteRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create incident_public_note, got error: %s", err))
        return
    }

    var incidentPublicNoteResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &incidentPublicNoteResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse incident_public_note response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := incidentPublicNoteResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = incidentPublicNoteResponse
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
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.Id = types.StringValue(string(jsonBytes))
            } else {
                data.Id = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.Id = types.StringValue(string(jsonBytes))
            } else {
                data.Id = types.StringValue(fmt.Sprintf("%v", obj["value"]))
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
    if obj, ok := dataMap["incidentId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncidentId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.IncidentId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.IncidentId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.IncidentId = types.StringValue(string(jsonBytes))
            } else {
                data.IncidentId = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.IncidentId = types.StringValue(string(jsonBytes))
            } else {
                data.IncidentId = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.IncidentId = types.StringValue(string(jsonBytes))
        } else {
            data.IncidentId = types.StringNull()
        }
    } else if val, ok := dataMap["incidentId"].(string); ok && val != "" {
        data.IncidentId = types.StringValue(val)
    } else {
        data.IncidentId = types.StringNull()
    }
    if obj, ok := dataMap["note"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Note = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Note = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Note = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.Note = types.StringValue(string(jsonBytes))
            } else {
                data.Note = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.Note = types.StringValue(string(jsonBytes))
            } else {
                data.Note = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Note = types.StringValue(string(jsonBytes))
        } else {
            data.Note = types.StringNull()
        }
    } else if val, ok := dataMap["note"].(string); ok && val != "" {
        data.Note = types.StringValue(val)
    } else {
        data.Note = types.StringNull()
    }
    if val, ok := dataMap["attachments"].([]interface{}); ok {
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
        // Sort list items by their string value to ensure consistent ordering
        // This fixes idempotency issues where server returns items in different order
        sort.Slice(listItems, func(i, j int) bool {
            iStr := listItems[i].(types.String).ValueString()
            jStr := listItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.Attachments = types.ListValueMust(types.StringType, listItems)
    } else {
        // For lists, always use empty list instead of null to match default values
        data.Attachments = types.ListValueMust(types.StringType, []attr.Value{})
    }
    if obj, ok := dataMap["subscriberNotificationStatusMessage"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SubscriberNotificationStatusMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SubscriberNotificationStatusMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SubscriberNotificationStatusMessage = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.SubscriberNotificationStatusMessage = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberNotificationStatusMessage = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.SubscriberNotificationStatusMessage = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberNotificationStatusMessage = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SubscriberNotificationStatusMessage = types.StringValue(string(jsonBytes))
        } else {
            data.SubscriberNotificationStatusMessage = types.StringNull()
        }
    } else if val, ok := dataMap["subscriberNotificationStatusMessage"].(string); ok && val != "" {
        data.SubscriberNotificationStatusMessage = types.StringValue(val)
    } else {
        data.SubscriberNotificationStatusMessage = types.StringNull()
    }
    if val, ok := dataMap["shouldStatusPageSubscribersBeNotifiedOnNoteCreated"].(bool); ok {
        data.ShouldStatusPageSubscribersBeNotifiedOnNoteCreated = types.BoolValue(val)
    }
    if obj, ok := dataMap["postedAt"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PostedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.PostedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.PostedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.PostedAt = types.StringValue(string(jsonBytes))
            } else {
                data.PostedAt = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.PostedAt = types.StringValue(string(jsonBytes))
            } else {
                data.PostedAt = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.PostedAt = types.StringValue(string(jsonBytes))
        } else {
            data.PostedAt = types.StringNull()
        }
    } else if val, ok := dataMap["postedAt"].(string); ok && val != "" {
        data.PostedAt = types.StringValue(val)
    } else {
        data.PostedAt = types.StringNull()
    }
    if obj, ok := dataMap["postedFromSlackMessageId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PostedFromSlackMessageId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.PostedFromSlackMessageId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.PostedFromSlackMessageId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.PostedFromSlackMessageId = types.StringValue(string(jsonBytes))
            } else {
                data.PostedFromSlackMessageId = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.PostedFromSlackMessageId = types.StringValue(string(jsonBytes))
            } else {
                data.PostedFromSlackMessageId = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.PostedFromSlackMessageId = types.StringValue(string(jsonBytes))
        } else {
            data.PostedFromSlackMessageId = types.StringNull()
        }
    } else if val, ok := dataMap["postedFromSlackMessageId"].(string); ok && val != "" {
        data.PostedFromSlackMessageId = types.StringValue(val)
    } else {
        data.PostedFromSlackMessageId = types.StringNull()
    }
    if obj, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.CreatedAt = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedAt = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.CreatedAt = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedAt = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedAt = types.StringNull()
        }
    } else if val, ok := dataMap["createdAt"].(string); ok && val != "" {
        data.CreatedAt = types.StringValue(val)
    } else {
        data.CreatedAt = types.StringNull()
    }
    if obj, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UpdatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.UpdatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.UpdatedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.UpdatedAt = types.StringValue(string(jsonBytes))
            } else {
                data.UpdatedAt = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.UpdatedAt = types.StringValue(string(jsonBytes))
            } else {
                data.UpdatedAt = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.UpdatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.UpdatedAt = types.StringNull()
        }
    } else if val, ok := dataMap["updatedAt"].(string); ok && val != "" {
        data.UpdatedAt = types.StringValue(val)
    } else {
        data.UpdatedAt = types.StringNull()
    }
    if obj, ok := dataMap["deletedAt"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DeletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DeletedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.DeletedAt = types.StringValue(string(jsonBytes))
            } else {
                data.DeletedAt = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.DeletedAt = types.StringValue(string(jsonBytes))
            } else {
                data.DeletedAt = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DeletedAt = types.StringValue(string(jsonBytes))
        } else {
            data.DeletedAt = types.StringNull()
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
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["createdByUserId"].(string); ok && val != "" {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if obj, ok := dataMap["subscriberNotificationStatusOnNoteCreated"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SubscriberNotificationStatusOnNoteCreated = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SubscriberNotificationStatusOnNoteCreated = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SubscriberNotificationStatusOnNoteCreated = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.SubscriberNotificationStatusOnNoteCreated = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberNotificationStatusOnNoteCreated = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.SubscriberNotificationStatusOnNoteCreated = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberNotificationStatusOnNoteCreated = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SubscriberNotificationStatusOnNoteCreated = types.StringValue(string(jsonBytes))
        } else {
            data.SubscriberNotificationStatusOnNoteCreated = types.StringNull()
        }
    } else if val, ok := dataMap["subscriberNotificationStatusOnNoteCreated"].(string); ok && val != "" {
        data.SubscriberNotificationStatusOnNoteCreated = types.StringValue(val)
    } else {
        data.SubscriberNotificationStatusOnNoteCreated = types.StringNull()
    }
    if val, ok := dataMap["isOwnerNotified"].(bool); ok {
        data.IsOwnerNotified = types.BoolValue(val)
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

func (r *IncidentPublicNoteResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var data IncidentPublicNoteResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create select parameter to get full object
    selectParam := map[string]interface{}{
        "projectId": true,
        "incidentId": true,
        "note": true,
        "attachments": true,
        "subscriberNotificationStatusMessage": true,
        "shouldStatusPageSubscribersBeNotifiedOnNoteCreated": true,
        "postedAt": true,
        "postedFromSlackMessageId": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "createdByUserId": true,
        "subscriberNotificationStatusOnNoteCreated": true,
        "isOwnerNotified": true,
        "_id": true,
    }

    // Make API call with select parameter
    httpResp, err := r.client.PostWithSelect("/incident-public-note/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read incident_public_note, got error: %s", err))
        return
    }

    if httpResp.StatusCode == http.StatusNotFound {
        resp.State.RemoveResource(ctx)
        return
    }

    var incidentPublicNoteResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &incidentPublicNoteResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse incident_public_note response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := incidentPublicNoteResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = incidentPublicNoteResponse
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
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.Id = types.StringValue(string(jsonBytes))
            } else {
                data.Id = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.Id = types.StringValue(string(jsonBytes))
            } else {
                data.Id = types.StringValue(fmt.Sprintf("%v", obj["value"]))
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
    if obj, ok := dataMap["incidentId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncidentId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.IncidentId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.IncidentId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.IncidentId = types.StringValue(string(jsonBytes))
            } else {
                data.IncidentId = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.IncidentId = types.StringValue(string(jsonBytes))
            } else {
                data.IncidentId = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.IncidentId = types.StringValue(string(jsonBytes))
        } else {
            data.IncidentId = types.StringNull()
        }
    } else if val, ok := dataMap["incidentId"].(string); ok && val != "" {
        data.IncidentId = types.StringValue(val)
    } else {
        data.IncidentId = types.StringNull()
    }
    if obj, ok := dataMap["note"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Note = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Note = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Note = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.Note = types.StringValue(string(jsonBytes))
            } else {
                data.Note = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.Note = types.StringValue(string(jsonBytes))
            } else {
                data.Note = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Note = types.StringValue(string(jsonBytes))
        } else {
            data.Note = types.StringNull()
        }
    } else if val, ok := dataMap["note"].(string); ok && val != "" {
        data.Note = types.StringValue(val)
    } else {
        data.Note = types.StringNull()
    }
    if val, ok := dataMap["attachments"].([]interface{}); ok {
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
        // Sort list items by their string value to ensure consistent ordering
        // This fixes idempotency issues where server returns items in different order
        sort.Slice(listItems, func(i, j int) bool {
            iStr := listItems[i].(types.String).ValueString()
            jStr := listItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.Attachments = types.ListValueMust(types.StringType, listItems)
    } else {
        // For lists, always use empty list instead of null to match default values
        data.Attachments = types.ListValueMust(types.StringType, []attr.Value{})
    }
    if obj, ok := dataMap["subscriberNotificationStatusMessage"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SubscriberNotificationStatusMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SubscriberNotificationStatusMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SubscriberNotificationStatusMessage = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.SubscriberNotificationStatusMessage = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberNotificationStatusMessage = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.SubscriberNotificationStatusMessage = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberNotificationStatusMessage = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SubscriberNotificationStatusMessage = types.StringValue(string(jsonBytes))
        } else {
            data.SubscriberNotificationStatusMessage = types.StringNull()
        }
    } else if val, ok := dataMap["subscriberNotificationStatusMessage"].(string); ok && val != "" {
        data.SubscriberNotificationStatusMessage = types.StringValue(val)
    } else {
        data.SubscriberNotificationStatusMessage = types.StringNull()
    }
    if val, ok := dataMap["shouldStatusPageSubscribersBeNotifiedOnNoteCreated"].(bool); ok {
        data.ShouldStatusPageSubscribersBeNotifiedOnNoteCreated = types.BoolValue(val)
    }
    if obj, ok := dataMap["postedAt"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PostedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.PostedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.PostedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.PostedAt = types.StringValue(string(jsonBytes))
            } else {
                data.PostedAt = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.PostedAt = types.StringValue(string(jsonBytes))
            } else {
                data.PostedAt = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.PostedAt = types.StringValue(string(jsonBytes))
        } else {
            data.PostedAt = types.StringNull()
        }
    } else if val, ok := dataMap["postedAt"].(string); ok && val != "" {
        data.PostedAt = types.StringValue(val)
    } else {
        data.PostedAt = types.StringNull()
    }
    if obj, ok := dataMap["postedFromSlackMessageId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PostedFromSlackMessageId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.PostedFromSlackMessageId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.PostedFromSlackMessageId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.PostedFromSlackMessageId = types.StringValue(string(jsonBytes))
            } else {
                data.PostedFromSlackMessageId = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.PostedFromSlackMessageId = types.StringValue(string(jsonBytes))
            } else {
                data.PostedFromSlackMessageId = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.PostedFromSlackMessageId = types.StringValue(string(jsonBytes))
        } else {
            data.PostedFromSlackMessageId = types.StringNull()
        }
    } else if val, ok := dataMap["postedFromSlackMessageId"].(string); ok && val != "" {
        data.PostedFromSlackMessageId = types.StringValue(val)
    } else {
        data.PostedFromSlackMessageId = types.StringNull()
    }
    if obj, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.CreatedAt = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedAt = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.CreatedAt = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedAt = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedAt = types.StringNull()
        }
    } else if val, ok := dataMap["createdAt"].(string); ok && val != "" {
        data.CreatedAt = types.StringValue(val)
    } else {
        data.CreatedAt = types.StringNull()
    }
    if obj, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UpdatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.UpdatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.UpdatedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.UpdatedAt = types.StringValue(string(jsonBytes))
            } else {
                data.UpdatedAt = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.UpdatedAt = types.StringValue(string(jsonBytes))
            } else {
                data.UpdatedAt = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.UpdatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.UpdatedAt = types.StringNull()
        }
    } else if val, ok := dataMap["updatedAt"].(string); ok && val != "" {
        data.UpdatedAt = types.StringValue(val)
    } else {
        data.UpdatedAt = types.StringNull()
    }
    if obj, ok := dataMap["deletedAt"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DeletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DeletedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.DeletedAt = types.StringValue(string(jsonBytes))
            } else {
                data.DeletedAt = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.DeletedAt = types.StringValue(string(jsonBytes))
            } else {
                data.DeletedAt = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DeletedAt = types.StringValue(string(jsonBytes))
        } else {
            data.DeletedAt = types.StringNull()
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
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["createdByUserId"].(string); ok && val != "" {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if obj, ok := dataMap["subscriberNotificationStatusOnNoteCreated"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SubscriberNotificationStatusOnNoteCreated = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SubscriberNotificationStatusOnNoteCreated = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SubscriberNotificationStatusOnNoteCreated = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.SubscriberNotificationStatusOnNoteCreated = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberNotificationStatusOnNoteCreated = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.SubscriberNotificationStatusOnNoteCreated = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberNotificationStatusOnNoteCreated = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SubscriberNotificationStatusOnNoteCreated = types.StringValue(string(jsonBytes))
        } else {
            data.SubscriberNotificationStatusOnNoteCreated = types.StringNull()
        }
    } else if val, ok := dataMap["subscriberNotificationStatusOnNoteCreated"].(string); ok && val != "" {
        data.SubscriberNotificationStatusOnNoteCreated = types.StringValue(val)
    } else {
        data.SubscriberNotificationStatusOnNoteCreated = types.StringNull()
    }
    if val, ok := dataMap["isOwnerNotified"].(bool); ok {
        data.IsOwnerNotified = types.BoolValue(val)
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IncidentPublicNoteResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var data IncidentPublicNoteResourceModel
    var state IncidentPublicNoteResourceModel

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
    incidentPublicNoteRequest := map[string]interface{}{
        "data": map[string]interface{}{},
    }
    requestDataMap := incidentPublicNoteRequest["data"].(map[string]interface{})

    if !data.Note.IsUnknown() && !state.Note.IsUnknown() && !data.Note.Equal(state.Note) {
        requestDataMap["note"] = data.Note.ValueString()
    }
    if !data.Attachments.IsUnknown() && !state.Attachments.IsUnknown() && !data.Attachments.Equal(state.Attachments) {
        requestDataMap["attachments"] = r.convertTerraformListToInterface(data.Attachments)
    }
    if !data.SubscriberNotificationStatusMessage.IsUnknown() && !state.SubscriberNotificationStatusMessage.IsUnknown() && !data.SubscriberNotificationStatusMessage.Equal(state.SubscriberNotificationStatusMessage) {
        requestDataMap["subscriberNotificationStatusMessage"] = data.SubscriberNotificationStatusMessage.ValueString()
    }
    if !data.PostedAt.IsUnknown() && !state.PostedAt.IsUnknown() && !data.PostedAt.Equal(state.PostedAt) {
        var postedatData interface{}
        if err := json.Unmarshal([]byte(data.PostedAt.ValueString()), &postedatData); err == nil {
            requestDataMap["postedAt"] = postedatData
        } else {
            requestDataMap["postedAt"] = data.PostedAt.ValueString()
        }
    }

    // Make API call
    httpResp, err := r.client.Put("/incident-public-note/" + data.Id.ValueString() + "", incidentPublicNoteRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update incident_public_note, got error: %s", err))
        return
    }

    // Parse the update response
    var incidentPublicNoteResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &incidentPublicNoteResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse incident_public_note response, got error: %s", err))
        return
    }

    // After successful update, fetch the current state by calling Read with select parameter
    selectParam := map[string]interface{}{
        "projectId": true,
        "incidentId": true,
        "note": true,
        "attachments": true,
        "subscriberNotificationStatusMessage": true,
        "shouldStatusPageSubscribersBeNotifiedOnNoteCreated": true,
        "postedAt": true,
        "postedFromSlackMessageId": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "createdByUserId": true,
        "subscriberNotificationStatusOnNoteCreated": true,
        "isOwnerNotified": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect("/incident-public-note/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read incident_public_note after update, got error: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse incident_public_note read response, got error: %s", err))
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
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.Id = types.StringValue(string(jsonBytes))
            } else {
                data.Id = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.Id = types.StringValue(string(jsonBytes))
            } else {
                data.Id = types.StringValue(fmt.Sprintf("%v", obj["value"]))
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
    if obj, ok := dataMap["incidentId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncidentId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.IncidentId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.IncidentId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.IncidentId = types.StringValue(string(jsonBytes))
            } else {
                data.IncidentId = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.IncidentId = types.StringValue(string(jsonBytes))
            } else {
                data.IncidentId = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.IncidentId = types.StringValue(string(jsonBytes))
        } else {
            data.IncidentId = types.StringNull()
        }
    } else if val, ok := dataMap["incidentId"].(string); ok && val != "" {
        data.IncidentId = types.StringValue(val)
    } else {
        data.IncidentId = types.StringNull()
    }
    if obj, ok := dataMap["note"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Note = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Note = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Note = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.Note = types.StringValue(string(jsonBytes))
            } else {
                data.Note = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.Note = types.StringValue(string(jsonBytes))
            } else {
                data.Note = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Note = types.StringValue(string(jsonBytes))
        } else {
            data.Note = types.StringNull()
        }
    } else if val, ok := dataMap["note"].(string); ok && val != "" {
        data.Note = types.StringValue(val)
    } else {
        data.Note = types.StringNull()
    }
    if val, ok := dataMap["attachments"].([]interface{}); ok {
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
        // Sort list items by their string value to ensure consistent ordering
        // This fixes idempotency issues where server returns items in different order
        sort.Slice(listItems, func(i, j int) bool {
            iStr := listItems[i].(types.String).ValueString()
            jStr := listItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.Attachments = types.ListValueMust(types.StringType, listItems)
    } else {
        // For lists, always use empty list instead of null to match default values
        data.Attachments = types.ListValueMust(types.StringType, []attr.Value{})
    }
    if obj, ok := dataMap["subscriberNotificationStatusMessage"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SubscriberNotificationStatusMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SubscriberNotificationStatusMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SubscriberNotificationStatusMessage = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.SubscriberNotificationStatusMessage = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberNotificationStatusMessage = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.SubscriberNotificationStatusMessage = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberNotificationStatusMessage = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SubscriberNotificationStatusMessage = types.StringValue(string(jsonBytes))
        } else {
            data.SubscriberNotificationStatusMessage = types.StringNull()
        }
    } else if val, ok := dataMap["subscriberNotificationStatusMessage"].(string); ok && val != "" {
        data.SubscriberNotificationStatusMessage = types.StringValue(val)
    } else {
        data.SubscriberNotificationStatusMessage = types.StringNull()
    }
    if val, ok := dataMap["shouldStatusPageSubscribersBeNotifiedOnNoteCreated"].(bool); ok {
        data.ShouldStatusPageSubscribersBeNotifiedOnNoteCreated = types.BoolValue(val)
    }
    if obj, ok := dataMap["postedAt"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PostedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.PostedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.PostedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.PostedAt = types.StringValue(string(jsonBytes))
            } else {
                data.PostedAt = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.PostedAt = types.StringValue(string(jsonBytes))
            } else {
                data.PostedAt = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.PostedAt = types.StringValue(string(jsonBytes))
        } else {
            data.PostedAt = types.StringNull()
        }
    } else if val, ok := dataMap["postedAt"].(string); ok && val != "" {
        data.PostedAt = types.StringValue(val)
    } else {
        data.PostedAt = types.StringNull()
    }
    if obj, ok := dataMap["postedFromSlackMessageId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PostedFromSlackMessageId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.PostedFromSlackMessageId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.PostedFromSlackMessageId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.PostedFromSlackMessageId = types.StringValue(string(jsonBytes))
            } else {
                data.PostedFromSlackMessageId = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.PostedFromSlackMessageId = types.StringValue(string(jsonBytes))
            } else {
                data.PostedFromSlackMessageId = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.PostedFromSlackMessageId = types.StringValue(string(jsonBytes))
        } else {
            data.PostedFromSlackMessageId = types.StringNull()
        }
    } else if val, ok := dataMap["postedFromSlackMessageId"].(string); ok && val != "" {
        data.PostedFromSlackMessageId = types.StringValue(val)
    } else {
        data.PostedFromSlackMessageId = types.StringNull()
    }
    if obj, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.CreatedAt = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedAt = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.CreatedAt = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedAt = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedAt = types.StringNull()
        }
    } else if val, ok := dataMap["createdAt"].(string); ok && val != "" {
        data.CreatedAt = types.StringValue(val)
    } else {
        data.CreatedAt = types.StringNull()
    }
    if obj, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UpdatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.UpdatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.UpdatedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.UpdatedAt = types.StringValue(string(jsonBytes))
            } else {
                data.UpdatedAt = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.UpdatedAt = types.StringValue(string(jsonBytes))
            } else {
                data.UpdatedAt = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.UpdatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.UpdatedAt = types.StringNull()
        }
    } else if val, ok := dataMap["updatedAt"].(string); ok && val != "" {
        data.UpdatedAt = types.StringValue(val)
    } else {
        data.UpdatedAt = types.StringNull()
    }
    if obj, ok := dataMap["deletedAt"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DeletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DeletedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.DeletedAt = types.StringValue(string(jsonBytes))
            } else {
                data.DeletedAt = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.DeletedAt = types.StringValue(string(jsonBytes))
            } else {
                data.DeletedAt = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DeletedAt = types.StringValue(string(jsonBytes))
        } else {
            data.DeletedAt = types.StringNull()
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
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["createdByUserId"].(string); ok && val != "" {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if obj, ok := dataMap["subscriberNotificationStatusOnNoteCreated"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SubscriberNotificationStatusOnNoteCreated = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SubscriberNotificationStatusOnNoteCreated = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SubscriberNotificationStatusOnNoteCreated = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.SubscriberNotificationStatusOnNoteCreated = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberNotificationStatusOnNoteCreated = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.SubscriberNotificationStatusOnNoteCreated = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberNotificationStatusOnNoteCreated = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SubscriberNotificationStatusOnNoteCreated = types.StringValue(string(jsonBytes))
        } else {
            data.SubscriberNotificationStatusOnNoteCreated = types.StringNull()
        }
    } else if val, ok := dataMap["subscriberNotificationStatusOnNoteCreated"].(string); ok && val != "" {
        data.SubscriberNotificationStatusOnNoteCreated = types.StringValue(val)
    } else {
        data.SubscriberNotificationStatusOnNoteCreated = types.StringNull()
    }
    if val, ok := dataMap["isOwnerNotified"].(bool); ok {
        data.IsOwnerNotified = types.BoolValue(val)
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IncidentPublicNoteResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var data IncidentPublicNoteResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Make API call
    _, err := r.client.Delete("/incident-public-note/" + data.Id.ValueString() + "")
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete incident_public_note, got error: %s", err))
        return
    }
}


func (r *IncidentPublicNoteResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Helper method to convert Terraform map to Go interface{}
func (r *IncidentPublicNoteResource) convertTerraformMapToInterface(terraformMap types.Map) interface{} {
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
func (r *IncidentPublicNoteResource) convertTerraformListToInterface(terraformList types.List) interface{} {
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
func (r *IncidentPublicNoteResource) parseJSONField(terraformString types.String) interface{} {
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

// Helper method to convert *big.Float to float64 for JSON serialization
func (r *IncidentPublicNoteResource) bigFloatToFloat64(bf *big.Float) interface{} {
    if bf == nil {
        return nil
    }
    f, _ := bf.Float64()
    return f
}

// Helper method to check if a type string is a valid OneUptime ObjectType
// Only these types should be marshalled/unmarshalled as typed wrapper objects
// This list is dynamically generated from Common/Types/JSON.ts ObjectType enum
func (r *IncidentPublicNoteResource) isValidOneUptimeObjectType(typeStr string) bool {
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
        "DashboardComponent": true,
        "DashboardViewConfig": true,
    }
    return validTypes[typeStr]
}
