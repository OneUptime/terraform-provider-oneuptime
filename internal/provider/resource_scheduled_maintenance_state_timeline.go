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
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ScheduledMaintenanceStateTimelineResource{}
var _ resource.ResourceWithImportState = &ScheduledMaintenanceStateTimelineResource{}

func NewScheduledMaintenanceStateTimelineResource() resource.Resource {
    return &ScheduledMaintenanceStateTimelineResource{}
}

// ScheduledMaintenanceStateTimelineResource defines the resource implementation.
type ScheduledMaintenanceStateTimelineResource struct {
    client *Client
}

// ScheduledMaintenanceStateTimelineResourceModel describes the resource data model.
type ScheduledMaintenanceStateTimelineResourceModel struct {
    Id types.String `tfsdk:"id"`
    ProjectId types.String `tfsdk:"project_id"`
    ScheduledMaintenanceId types.String `tfsdk:"scheduled_maintenance_id"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    ScheduledMaintenanceStateId types.String `tfsdk:"scheduled_maintenance_state_id"`
    ShouldStatusPageSubscribersBeNotified types.Bool `tfsdk:"should_status_page_subscribers_be_notified"`
    EndsAt RFC3339Value `tfsdk:"ends_at"`
    StartsAt RFC3339Value `tfsdk:"starts_at"`
    CreatedAt RFC3339Value `tfsdk:"created_at"`
    UpdatedAt RFC3339Value `tfsdk:"updated_at"`
    DeletedAt RFC3339Value `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    SubscriberNotificationStatus types.String `tfsdk:"subscriber_notification_status"`
    SubscriberNotificationStatusMessage types.String `tfsdk:"subscriber_notification_status_message"`
    IsOwnerNotified types.Bool `tfsdk:"is_owner_notified"`
}

func (r *ScheduledMaintenanceStateTimelineResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_scheduled_maintenance_state_timeline"
}

func (r *ScheduledMaintenanceStateTimelineResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Change state of your scheduled maintenance event.",

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
            "scheduled_maintenance_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Required: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.RequiresReplace(),
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
            "scheduled_maintenance_state_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Required: true,
            },
            "should_status_page_subscribers_be_notified": schema.BoolAttribute{
                MarkdownDescription: "Should subscribers be notified about this state change?.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                    boolplanmodifier.RequiresReplace(),
                },
            },
            "ends_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                CustomType: RFC3339Type{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                    stringplanmodifier.RequiresReplace(),
                },
            },
            "starts_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                CustomType: RFC3339Type{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
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
            "subscriber_notification_status": schema.StringAttribute{
                MarkdownDescription: "Status of notification sent to subscribers about this scheduled maintenance state change.",
                Computed: true,
            },
            "subscriber_notification_status_message": schema.StringAttribute{
                MarkdownDescription: "Status message for subscriber notifications - includes success messages, failure reasons, or skip reasons.",
                Computed: true,
            },
            "is_owner_notified": schema.BoolAttribute{
                MarkdownDescription: "Are owners notified of state change?.",
                Computed: true,
            },
        },
    }
}

func (r *ScheduledMaintenanceStateTimelineResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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


func (r *ScheduledMaintenanceStateTimelineResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var data ScheduledMaintenanceStateTimelineResourceModel

    // Read Terraform plan data into the model
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }



    // Create API request body. Unset (null/unknown) optional fields are
    // omitted so server-side defaults apply instead of being overwritten
    // with zero values.
    scheduledMaintenanceStateTimelineRequest := map[string]interface{}{
        "data": map[string]interface{}{},
    }
    requestDataMap := scheduledMaintenanceStateTimelineRequest["data"].(map[string]interface{})

    if !data.ScheduledMaintenanceId.IsNull() && !data.ScheduledMaintenanceId.IsUnknown() {
        requestDataMap["scheduledMaintenanceId"] = data.ScheduledMaintenanceId.ValueString()
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        requestDataMap["createdByUserId"] = data.CreatedByUserId.ValueString()
    }
    if !data.ScheduledMaintenanceStateId.IsNull() && !data.ScheduledMaintenanceStateId.IsUnknown() {
        requestDataMap["scheduledMaintenanceStateId"] = data.ScheduledMaintenanceStateId.ValueString()
    }
    if !data.ShouldStatusPageSubscribersBeNotified.IsNull() && !data.ShouldStatusPageSubscribersBeNotified.IsUnknown() {
        requestDataMap["shouldStatusPageSubscribersBeNotified"] = data.ShouldStatusPageSubscribersBeNotified.ValueBool()
    }
    if !data.EndsAt.IsNull() && !data.EndsAt.IsUnknown() {
        requestDataMap["endsAt"] = data.EndsAt.ValueString()
    }
    if !data.StartsAt.IsNull() && !data.StartsAt.IsUnknown() {
        requestDataMap["startsAt"] = data.StartsAt.ValueString()
    }

    // Make API call
    httpResp, err := r.client.Post(ctx, "/scheduled-maintenance-state-timeline", scheduledMaintenanceStateTimelineRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create scheduled_maintenance_state_timeline, got error: %s", err))
        return
    }

    var scheduledMaintenanceStateTimelineResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &scheduledMaintenanceStateTimelineResponse)
    if err != nil {
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to create scheduled_maintenance_state_timeline: %s", err))
        return
    }

    // Extract the new resource id from the create response.
    createdId := ""
    if wrapper, ok := scheduledMaintenanceStateTimelineResponse["data"].(map[string]interface{}); ok {
        if val, ok := wrapper["_id"].(string); ok {
            createdId = val
        }
    } else if val, ok := scheduledMaintenanceStateTimelineResponse["_id"].(string); ok {
        createdId = val
    }
    if createdId == "" {
        resp.Diagnostics.AddError("OneUptime API Error", "Create response for scheduled_maintenance_state_timeline did not contain an id. This is a bug in the provider or the API; please report it.")
        return
    }
    data.Id = types.StringValue(createdId)

    /*
     * The server has committed the row. Persist what we know to state BEFORE
     * the read-back: if the read-back fails and we return without setting
     * state, Terraform never learns the resource exists and the created
     * scheduled_maintenance_state_timeline is orphaned server-side — never refreshed, never
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
        "scheduledMaintenanceId": true,
        "createdByUserId": true,
        "scheduledMaintenanceStateId": true,
        "shouldStatusPageSubscribersBeNotified": true,
        "endsAt": true,
        "startsAt": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "subscriberNotificationStatus": true,
        "subscriberNotificationStatusMessage": true,
        "isOwnerNotified": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect(ctx, "/scheduled-maintenance-state-timeline/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        /*
         * State already owns the id, so the resource is tracked and the next
         * refresh reconciles the remaining attributes. Warn rather than
         * error: erroring here would strand a real resource.
         */
        resp.Diagnostics.AddWarning("Read After Create Failed", fmt.Sprintf("Created scheduled_maintenance_state_timeline but could not read it back; state is incomplete until the next refresh: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddWarning("Read After Create Failed", fmt.Sprintf("Created scheduled_maintenance_state_timeline but could not parse the read-back response; state is incomplete until the next refresh: %s", err))
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
    if obj, ok := dataMap["scheduledMaintenanceId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ScheduledMaintenanceId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ScheduledMaintenanceId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ScheduledMaintenanceId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ScheduledMaintenanceId = types.StringValue(string(jsonBytes))
            } else {
                data.ScheduledMaintenanceId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ScheduledMaintenanceId = types.StringValue(string(jsonBytes))
            } else {
                data.ScheduledMaintenanceId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ScheduledMaintenanceId = types.StringValue(string(jsonBytes))
        } else {
            data.ScheduledMaintenanceId = types.StringNull()
        }
    } else if val, ok := dataMap["scheduledMaintenanceId"].(string); ok {
        data.ScheduledMaintenanceId = types.StringValue(val)
    } else {
        data.ScheduledMaintenanceId = types.StringNull()
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
    if obj, ok := dataMap["scheduledMaintenanceStateId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ScheduledMaintenanceStateId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ScheduledMaintenanceStateId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ScheduledMaintenanceStateId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ScheduledMaintenanceStateId = types.StringValue(string(jsonBytes))
            } else {
                data.ScheduledMaintenanceStateId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ScheduledMaintenanceStateId = types.StringValue(string(jsonBytes))
            } else {
                data.ScheduledMaintenanceStateId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ScheduledMaintenanceStateId = types.StringValue(string(jsonBytes))
        } else {
            data.ScheduledMaintenanceStateId = types.StringNull()
        }
    } else if val, ok := dataMap["scheduledMaintenanceStateId"].(string); ok {
        data.ScheduledMaintenanceStateId = types.StringValue(val)
    } else {
        data.ScheduledMaintenanceStateId = types.StringNull()
    }
    if val, ok := dataMap["shouldStatusPageSubscribersBeNotified"].(bool); ok {
        data.ShouldStatusPageSubscribersBeNotified = types.BoolValue(val)
    } else {
        data.ShouldStatusPageSubscribersBeNotified = types.BoolNull()
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
    if obj, ok := dataMap["subscriberNotificationStatus"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SubscriberNotificationStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SubscriberNotificationStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SubscriberNotificationStatus = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SubscriberNotificationStatus = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberNotificationStatus = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SubscriberNotificationStatus = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberNotificationStatus = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SubscriberNotificationStatus = types.StringValue(string(jsonBytes))
        } else {
            data.SubscriberNotificationStatus = types.StringNull()
        }
    } else if val, ok := dataMap["subscriberNotificationStatus"].(string); ok {
        data.SubscriberNotificationStatus = types.StringValue(val)
    } else {
        data.SubscriberNotificationStatus = types.StringNull()
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
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SubscriberNotificationStatusMessage = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberNotificationStatusMessage = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SubscriberNotificationStatusMessage = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberNotificationStatusMessage = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SubscriberNotificationStatusMessage = types.StringValue(string(jsonBytes))
        } else {
            data.SubscriberNotificationStatusMessage = types.StringNull()
        }
    } else if val, ok := dataMap["subscriberNotificationStatusMessage"].(string); ok {
        data.SubscriberNotificationStatusMessage = types.StringValue(val)
    } else {
        data.SubscriberNotificationStatusMessage = types.StringNull()
    }
    if val, ok := dataMap["isOwnerNotified"].(bool); ok {
        data.IsOwnerNotified = types.BoolValue(val)
    } else {
        data.IsOwnerNotified = types.BoolNull()
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

func (r *ScheduledMaintenanceStateTimelineResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var data ScheduledMaintenanceStateTimelineResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create select parameter to get full object
    selectParam := map[string]interface{}{
        "projectId": true,
        "scheduledMaintenanceId": true,
        "createdByUserId": true,
        "scheduledMaintenanceStateId": true,
        "shouldStatusPageSubscribersBeNotified": true,
        "endsAt": true,
        "startsAt": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "subscriberNotificationStatus": true,
        "subscriberNotificationStatusMessage": true,
        "isOwnerNotified": true,
        "_id": true,
    }

    // Make API call with select parameter
    httpResp, err := r.client.PostWithSelect(ctx, "/scheduled-maintenance-state-timeline/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read scheduled_maintenance_state_timeline, got error: %s", err))
        return
    }

    if httpResp.StatusCode == http.StatusNotFound {
        resp.State.RemoveResource(ctx)
        return
    }

    var scheduledMaintenanceStateTimelineResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &scheduledMaintenanceStateTimelineResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse scheduled_maintenance_state_timeline response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := scheduledMaintenanceStateTimelineResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = scheduledMaintenanceStateTimelineResponse
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
    if obj, ok := dataMap["scheduledMaintenanceId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ScheduledMaintenanceId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ScheduledMaintenanceId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ScheduledMaintenanceId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ScheduledMaintenanceId = types.StringValue(string(jsonBytes))
            } else {
                data.ScheduledMaintenanceId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ScheduledMaintenanceId = types.StringValue(string(jsonBytes))
            } else {
                data.ScheduledMaintenanceId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ScheduledMaintenanceId = types.StringValue(string(jsonBytes))
        } else {
            data.ScheduledMaintenanceId = types.StringNull()
        }
    } else if val, ok := dataMap["scheduledMaintenanceId"].(string); ok {
        data.ScheduledMaintenanceId = types.StringValue(val)
    } else {
        data.ScheduledMaintenanceId = types.StringNull()
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
    if obj, ok := dataMap["scheduledMaintenanceStateId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ScheduledMaintenanceStateId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ScheduledMaintenanceStateId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ScheduledMaintenanceStateId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ScheduledMaintenanceStateId = types.StringValue(string(jsonBytes))
            } else {
                data.ScheduledMaintenanceStateId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ScheduledMaintenanceStateId = types.StringValue(string(jsonBytes))
            } else {
                data.ScheduledMaintenanceStateId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ScheduledMaintenanceStateId = types.StringValue(string(jsonBytes))
        } else {
            data.ScheduledMaintenanceStateId = types.StringNull()
        }
    } else if val, ok := dataMap["scheduledMaintenanceStateId"].(string); ok {
        data.ScheduledMaintenanceStateId = types.StringValue(val)
    } else {
        data.ScheduledMaintenanceStateId = types.StringNull()
    }
    if val, ok := dataMap["shouldStatusPageSubscribersBeNotified"].(bool); ok {
        data.ShouldStatusPageSubscribersBeNotified = types.BoolValue(val)
    } else {
        data.ShouldStatusPageSubscribersBeNotified = types.BoolNull()
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
    if obj, ok := dataMap["subscriberNotificationStatus"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SubscriberNotificationStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SubscriberNotificationStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SubscriberNotificationStatus = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SubscriberNotificationStatus = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberNotificationStatus = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SubscriberNotificationStatus = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberNotificationStatus = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SubscriberNotificationStatus = types.StringValue(string(jsonBytes))
        } else {
            data.SubscriberNotificationStatus = types.StringNull()
        }
    } else if val, ok := dataMap["subscriberNotificationStatus"].(string); ok {
        data.SubscriberNotificationStatus = types.StringValue(val)
    } else {
        data.SubscriberNotificationStatus = types.StringNull()
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
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SubscriberNotificationStatusMessage = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberNotificationStatusMessage = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SubscriberNotificationStatusMessage = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberNotificationStatusMessage = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SubscriberNotificationStatusMessage = types.StringValue(string(jsonBytes))
        } else {
            data.SubscriberNotificationStatusMessage = types.StringNull()
        }
    } else if val, ok := dataMap["subscriberNotificationStatusMessage"].(string); ok {
        data.SubscriberNotificationStatusMessage = types.StringValue(val)
    } else {
        data.SubscriberNotificationStatusMessage = types.StringNull()
    }
    if val, ok := dataMap["isOwnerNotified"].(bool); ok {
        data.IsOwnerNotified = types.BoolValue(val)
    } else {
        data.IsOwnerNotified = types.BoolNull()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ScheduledMaintenanceStateTimelineResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var data ScheduledMaintenanceStateTimelineResourceModel
    var state ScheduledMaintenanceStateTimelineResourceModel

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
    scheduledMaintenanceStateTimelineRequest := map[string]interface{}{
        "data": map[string]interface{}{},
    }
    requestDataMap := scheduledMaintenanceStateTimelineRequest["data"].(map[string]interface{})

    if !data.ScheduledMaintenanceStateId.IsUnknown() && !state.ScheduledMaintenanceStateId.IsUnknown() && !data.ScheduledMaintenanceStateId.Equal(state.ScheduledMaintenanceStateId) {
        requestDataMap["scheduledMaintenanceStateId"] = data.ScheduledMaintenanceStateId.ValueString()
    }
    if !data.StartsAt.IsUnknown() && !state.StartsAt.IsUnknown() && !data.StartsAt.Equal(state.StartsAt) {
        requestDataMap["startsAt"] = data.StartsAt.ValueString()
    }

    // Only call the API when there are changed fields to send. An empty
    // update body is rejected by the API; state is still refreshed below so
    // this method never writes unverified plan values into state.
    if len(scheduledMaintenanceStateTimelineRequest["data"].(map[string]interface{})) > 0 {
        httpResp, err := r.client.Put(ctx, "/scheduled-maintenance-state-timeline/" + data.Id.ValueString() + "", scheduledMaintenanceStateTimelineRequest)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update scheduled_maintenance_state_timeline, got error: %s", err))
            return
        }

        // Parse the update response
        var scheduledMaintenanceStateTimelineResponse map[string]interface{}
        err = r.client.ParseResponse(httpResp, &scheduledMaintenanceStateTimelineResponse)
        if err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to update scheduled_maintenance_state_timeline: %s", err))
            return
        }
        _ = scheduledMaintenanceStateTimelineResponse
    }

    // After successful update, fetch the current state by calling Read with select parameter
    selectParam := map[string]interface{}{
        "projectId": true,
        "scheduledMaintenanceId": true,
        "createdByUserId": true,
        "scheduledMaintenanceStateId": true,
        "shouldStatusPageSubscribersBeNotified": true,
        "endsAt": true,
        "startsAt": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "subscriberNotificationStatus": true,
        "subscriberNotificationStatusMessage": true,
        "isOwnerNotified": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect(ctx, "/scheduled-maintenance-state-timeline/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read scheduled_maintenance_state_timeline after update, got error: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read scheduled_maintenance_state_timeline after update: %s", err))
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
    if obj, ok := dataMap["scheduledMaintenanceId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ScheduledMaintenanceId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ScheduledMaintenanceId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ScheduledMaintenanceId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ScheduledMaintenanceId = types.StringValue(string(jsonBytes))
            } else {
                data.ScheduledMaintenanceId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ScheduledMaintenanceId = types.StringValue(string(jsonBytes))
            } else {
                data.ScheduledMaintenanceId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ScheduledMaintenanceId = types.StringValue(string(jsonBytes))
        } else {
            data.ScheduledMaintenanceId = types.StringNull()
        }
    } else if val, ok := dataMap["scheduledMaintenanceId"].(string); ok {
        data.ScheduledMaintenanceId = types.StringValue(val)
    } else {
        data.ScheduledMaintenanceId = types.StringNull()
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
    if obj, ok := dataMap["scheduledMaintenanceStateId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ScheduledMaintenanceStateId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ScheduledMaintenanceStateId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ScheduledMaintenanceStateId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ScheduledMaintenanceStateId = types.StringValue(string(jsonBytes))
            } else {
                data.ScheduledMaintenanceStateId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ScheduledMaintenanceStateId = types.StringValue(string(jsonBytes))
            } else {
                data.ScheduledMaintenanceStateId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ScheduledMaintenanceStateId = types.StringValue(string(jsonBytes))
        } else {
            data.ScheduledMaintenanceStateId = types.StringNull()
        }
    } else if val, ok := dataMap["scheduledMaintenanceStateId"].(string); ok {
        data.ScheduledMaintenanceStateId = types.StringValue(val)
    } else {
        data.ScheduledMaintenanceStateId = types.StringNull()
    }
    if val, ok := dataMap["shouldStatusPageSubscribersBeNotified"].(bool); ok {
        data.ShouldStatusPageSubscribersBeNotified = types.BoolValue(val)
    } else {
        data.ShouldStatusPageSubscribersBeNotified = types.BoolNull()
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
    if obj, ok := dataMap["subscriberNotificationStatus"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SubscriberNotificationStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SubscriberNotificationStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SubscriberNotificationStatus = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SubscriberNotificationStatus = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberNotificationStatus = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SubscriberNotificationStatus = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberNotificationStatus = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SubscriberNotificationStatus = types.StringValue(string(jsonBytes))
        } else {
            data.SubscriberNotificationStatus = types.StringNull()
        }
    } else if val, ok := dataMap["subscriberNotificationStatus"].(string); ok {
        data.SubscriberNotificationStatus = types.StringValue(val)
    } else {
        data.SubscriberNotificationStatus = types.StringNull()
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
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SubscriberNotificationStatusMessage = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberNotificationStatusMessage = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SubscriberNotificationStatusMessage = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberNotificationStatusMessage = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SubscriberNotificationStatusMessage = types.StringValue(string(jsonBytes))
        } else {
            data.SubscriberNotificationStatusMessage = types.StringNull()
        }
    } else if val, ok := dataMap["subscriberNotificationStatusMessage"].(string); ok {
        data.SubscriberNotificationStatusMessage = types.StringValue(val)
    } else {
        data.SubscriberNotificationStatusMessage = types.StringNull()
    }
    if val, ok := dataMap["isOwnerNotified"].(bool); ok {
        data.IsOwnerNotified = types.BoolValue(val)
    } else {
        data.IsOwnerNotified = types.BoolNull()
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

func (r *ScheduledMaintenanceStateTimelineResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var data ScheduledMaintenanceStateTimelineResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Make API call
    httpResp, err := r.client.Delete(ctx, "/scheduled-maintenance-state-timeline/" + data.Id.ValueString() + "")
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete scheduled_maintenance_state_timeline, got error: %s", err))
        return
    }

    // A failed delete must keep the resource in state — silently dropping it
    // orphans real infrastructure. 404 means it is already gone.
    if httpResp.StatusCode >= 400 && httpResp.StatusCode != http.StatusNotFound {
        err = r.client.ParseResponse(httpResp, nil)
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to delete scheduled_maintenance_state_timeline: %s", err))
        return
    }
    if httpResp.Body != nil {
        httpResp.Body.Close()
    }
}


func (r *ScheduledMaintenanceStateTimelineResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Helper method to convert Terraform map to Go interface{}
func (r *ScheduledMaintenanceStateTimelineResource) convertTerraformMapToInterface(terraformMap types.Map) interface{} {
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
func (r *ScheduledMaintenanceStateTimelineResource) convertTerraformListToInterface(terraformList types.List) interface{} {
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
func (r *ScheduledMaintenanceStateTimelineResource) convertTerraformSetToInterface(terraformSet types.Set) interface{} {
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
func (r *ScheduledMaintenanceStateTimelineResource) parseJSONField(terraformString basetypes.StringValuable) interface{} {
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
func (r *ScheduledMaintenanceStateTimelineResource) normalizeURLWrappers(value interface{}) interface{} {
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

func (r *ScheduledMaintenanceStateTimelineResource) normalizeURLString(value string) string {
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
func (r *ScheduledMaintenanceStateTimelineResource) bigFloatToFloat64(bf *big.Float) interface{} {
    if bf == nil {
        return nil
    }
    f, _ := bf.Float64()
    return f
}

// Helper method to check if a type string is a valid OneUptime ObjectType.
// The registry itself lives in objecttypes.go, shared across the package.
func (r *ScheduledMaintenanceStateTimelineResource) isValidOneUptimeObjectType(typeStr string) bool {
    return validOneUptimeObjectTypes[typeStr]
}
