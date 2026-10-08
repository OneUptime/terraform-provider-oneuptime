package provider

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "github.com/hashicorp/terraform-plugin-framework/attr"
    "sort"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &IncidentEpisodePublicNoteDataSource{}

func NewIncidentEpisodePublicNoteDataSource() datasource.DataSource {
    return &IncidentEpisodePublicNoteDataSource{}
}

// IncidentEpisodePublicNoteDataSource defines the data source implementation.
type IncidentEpisodePublicNoteDataSource struct {
    client *Client
}

// IncidentEpisodePublicNoteDataSourceModel describes the data source data model.
type IncidentEpisodePublicNoteDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    IncidentEpisodeId types.String `tfsdk:"incident_episode_id"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    Note types.String `tfsdk:"note"`
    Attachments types.Set `tfsdk:"attachments"`
    SubscriberNotificationStatusOnNoteCreated types.String `tfsdk:"subscriber_notification_status_on_note_created"`
    SubscriberNotificationStatusMessage types.String `tfsdk:"subscriber_notification_status_message"`
    SubscriberNotificationStatusOnNoteUpdated types.String `tfsdk:"subscriber_notification_status_on_note_updated"`
    SubscriberNotificationStatusMessageOnNoteUpdated types.String `tfsdk:"subscriber_notification_status_message_on_note_updated"`
    ShouldStatusPageSubscribersBeNotifiedOnNoteCreated types.Bool `tfsdk:"should_status_page_subscribers_be_notified_on_note_created"`
    IsOwnerNotified types.Bool `tfsdk:"is_owner_notified"`
    PostedAt types.String `tfsdk:"posted_at"`
    PostedFromSlackMessageId types.String `tfsdk:"posted_from_slack_message_id"`
}

func (d *IncidentEpisodePublicNoteDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_incident_episode_public_note"
}

func (d *IncidentEpisodePublicNoteDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Manage public notes for your incident episode Look up an existing incident episode public note by `id`, or by any of its other arguments (`created_by_user_id`, `incident_episode_id`, `is_owner_notified`, ...): each one set must match, and exactly one incident episode public note may match them all.",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Look up by unique identifier. Leave unset to look up by the other arguments instead.",
                Optional: true,
                Computed: true,
            },
            "created_at": schema.StringAttribute{
                MarkdownDescription: "Date and Time when the object was created.",
                Computed: true,
            },
            "updated_at": schema.StringAttribute{
                MarkdownDescription: "Date and Time when the object was updated.",
                Computed: true,
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "incident_episode_id": schema.StringAttribute{
                MarkdownDescription: "Relation to Incident Episode ID in which this resource belongs. The ID of a `oneuptime_incident_episode`.",
                Optional: true,
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "note": schema.StringAttribute{
                MarkdownDescription: "Notes in markdown.",
                Optional: true,
                Computed: true,
            },
            "attachments": schema.SetAttribute{
                MarkdownDescription: "Files attached to this note. IDs of `oneuptime_file` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "subscriber_notification_status_on_note_created": schema.StringAttribute{
                MarkdownDescription: "Status of notification sent to subscribers about this note.",
                Optional: true,
                Computed: true,
            },
            "subscriber_notification_status_message": schema.StringAttribute{
                MarkdownDescription: "Status message for subscriber notifications - includes success messages, failure reasons, or skip reasons.",
                Optional: true,
                Computed: true,
            },
            "subscriber_notification_status_on_note_updated": schema.StringAttribute{
                MarkdownDescription: "Status of the notification sent to subscribers when this note was last updated. Empty until an update notification is requested.",
                Optional: true,
                Computed: true,
            },
            "subscriber_notification_status_message_on_note_updated": schema.StringAttribute{
                MarkdownDescription: "Status message for the notification sent to subscribers when this note was last updated - includes success messages, failure reasons, or skip reasons.",
                Optional: true,
                Computed: true,
            },
            "should_status_page_subscribers_be_notified_on_note_created": schema.BoolAttribute{
                MarkdownDescription: "Should subscribers be notified about this note? If left out, this follows the episode: true when subscribers were notified that the episode was created, false when it was created without notifying them.",
                Optional: true,
                Computed: true,
            },
            "is_owner_notified": schema.BoolAttribute{
                MarkdownDescription: "Are owners notified of this resource ownership?",
                Optional: true,
                Computed: true,
            },
            "posted_at": schema.StringAttribute{
                MarkdownDescription: "Date and time when the note was posted.",
                Computed: true,
            },
            "posted_from_slack_message_id": schema.StringAttribute{
                MarkdownDescription: "Unique identifier for the Slack message this note was created from (channel_id:message_ts). Used to prevent duplicate notes when multiple users react to the same message.",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *IncidentEpisodePublicNoteDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *IncidentEpisodePublicNoteDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data IncidentEpisodePublicNoteDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.IncidentEpisodeId.IsNull() && !data.IncidentEpisodeId.IsUnknown() {
        filters["incidentEpisodeId"] = data.IncidentEpisodeId.ValueString()
        filterNames = append(filterNames, "incident_episode_id = "+fmt.Sprintf("%q", data.IncidentEpisodeId.ValueString()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }
    if !data.Note.IsNull() && !data.Note.IsUnknown() {
        filters["note"] = data.Note.ValueString()
        filterNames = append(filterNames, "note = "+fmt.Sprintf("%q", data.Note.ValueString()))
    }
    if !data.SubscriberNotificationStatusOnNoteCreated.IsNull() && !data.SubscriberNotificationStatusOnNoteCreated.IsUnknown() {
        filters["subscriberNotificationStatusOnNoteCreated"] = data.SubscriberNotificationStatusOnNoteCreated.ValueString()
        filterNames = append(filterNames, "subscriber_notification_status_on_note_created = "+fmt.Sprintf("%q", data.SubscriberNotificationStatusOnNoteCreated.ValueString()))
    }
    if !data.SubscriberNotificationStatusMessage.IsNull() && !data.SubscriberNotificationStatusMessage.IsUnknown() {
        filters["subscriberNotificationStatusMessage"] = data.SubscriberNotificationStatusMessage.ValueString()
        filterNames = append(filterNames, "subscriber_notification_status_message = "+fmt.Sprintf("%q", data.SubscriberNotificationStatusMessage.ValueString()))
    }
    if !data.SubscriberNotificationStatusOnNoteUpdated.IsNull() && !data.SubscriberNotificationStatusOnNoteUpdated.IsUnknown() {
        filters["subscriberNotificationStatusOnNoteUpdated"] = data.SubscriberNotificationStatusOnNoteUpdated.ValueString()
        filterNames = append(filterNames, "subscriber_notification_status_on_note_updated = "+fmt.Sprintf("%q", data.SubscriberNotificationStatusOnNoteUpdated.ValueString()))
    }
    if !data.SubscriberNotificationStatusMessageOnNoteUpdated.IsNull() && !data.SubscriberNotificationStatusMessageOnNoteUpdated.IsUnknown() {
        filters["subscriberNotificationStatusMessageOnNoteUpdated"] = data.SubscriberNotificationStatusMessageOnNoteUpdated.ValueString()
        filterNames = append(filterNames, "subscriber_notification_status_message_on_note_updated = "+fmt.Sprintf("%q", data.SubscriberNotificationStatusMessageOnNoteUpdated.ValueString()))
    }
    if !data.ShouldStatusPageSubscribersBeNotifiedOnNoteCreated.IsNull() && !data.ShouldStatusPageSubscribersBeNotifiedOnNoteCreated.IsUnknown() {
        filters["shouldStatusPageSubscribersBeNotifiedOnNoteCreated"] = data.ShouldStatusPageSubscribersBeNotifiedOnNoteCreated.ValueBool()
        filterNames = append(filterNames, "should_status_page_subscribers_be_notified_on_note_created = "+fmt.Sprintf("%t", data.ShouldStatusPageSubscribersBeNotifiedOnNoteCreated.ValueBool()))
    }
    if !data.IsOwnerNotified.IsNull() && !data.IsOwnerNotified.IsUnknown() {
        filters["isOwnerNotified"] = data.IsOwnerNotified.ValueBool()
        filterNames = append(filterNames, "is_owner_notified = "+fmt.Sprintf("%t", data.IsOwnerNotified.ValueBool()))
    }
    if !data.PostedFromSlackMessageId.IsNull() && !data.PostedFromSlackMessageId.IsUnknown() {
        filters["postedFromSlackMessageId"] = data.PostedFromSlackMessageId.ValueString()
        filterNames = append(filterNames, "posted_from_slack_message_id = "+fmt.Sprintf("%q", data.PostedFromSlackMessageId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the incident episode public note up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the incident episode public note up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "incidentEpisodeId": true,
        "createdByUserId": true,
        "note": true,
        "attachments": true,
        "subscriberNotificationStatusOnNoteCreated": true,
        "subscriberNotificationStatusMessage": true,
        "subscriberNotificationStatusOnNoteUpdated": true,
        "subscriberNotificationStatusMessageOnNoteUpdated": true,
        "shouldStatusPageSubscribersBeNotifiedOnNoteCreated": true,
        "isOwnerNotified": true,
        "postedAt": true,
        "postedFromSlackMessageId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/incident-episode-public-note/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read incident_episode_public_note, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No incident episode public note found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read incident_episode_public_note: %s", err))
            return
        }
        if wrapper, ok := itemResponse["data"].(map[string]interface{}); ok {
            item = wrapper
        } else {
            item = itemResponse
        }
    }
    if !hasId {
        listBody := map[string]interface{}{
            "query":  filters,
            "select": selectParam,
            // limit 2 is enough to detect ambiguity without paging.
            "limit": 2,
        }
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/incident-episode-public-note/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list incident_episode_public_note, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list incident_episode_public_note: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No incident episode public note matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one incident episode public note matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for incident_episode_public_note.")
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
    if obj, ok := item["incidentEpisodeId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncidentEpisodeId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.IncidentEpisodeId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.IncidentEpisodeId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.IncidentEpisodeId = types.StringValue(string(jsonBytes))
        } else {
            data.IncidentEpisodeId = types.StringNull()
        }
    } else if val, ok := item["incidentEpisodeId"].(string); ok {
        data.IncidentEpisodeId = types.StringValue(val)
    } else {
        data.IncidentEpisodeId = types.StringNull()
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
    if obj, ok := item["note"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Note = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Note = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Note = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Note = types.StringValue(string(jsonBytes))
        } else {
            data.Note = types.StringNull()
        }
    } else if val, ok := item["note"].(string); ok {
        data.Note = types.StringValue(val)
    } else {
        data.Note = types.StringNull()
    }
    if val, ok := item["attachments"].([]interface{}); ok {
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
        data.Attachments = types.SetValueMust(types.StringType, setItems)
    } else {
        data.Attachments = types.SetNull(types.StringType)
    }
    if obj, ok := item["subscriberNotificationStatusOnNoteCreated"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SubscriberNotificationStatusOnNoteCreated = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SubscriberNotificationStatusOnNoteCreated = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SubscriberNotificationStatusOnNoteCreated = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SubscriberNotificationStatusOnNoteCreated = types.StringValue(string(jsonBytes))
        } else {
            data.SubscriberNotificationStatusOnNoteCreated = types.StringNull()
        }
    } else if val, ok := item["subscriberNotificationStatusOnNoteCreated"].(string); ok {
        data.SubscriberNotificationStatusOnNoteCreated = types.StringValue(val)
    } else {
        data.SubscriberNotificationStatusOnNoteCreated = types.StringNull()
    }
    if obj, ok := item["subscriberNotificationStatusMessage"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SubscriberNotificationStatusMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SubscriberNotificationStatusMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SubscriberNotificationStatusMessage = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SubscriberNotificationStatusMessage = types.StringValue(string(jsonBytes))
        } else {
            data.SubscriberNotificationStatusMessage = types.StringNull()
        }
    } else if val, ok := item["subscriberNotificationStatusMessage"].(string); ok {
        data.SubscriberNotificationStatusMessage = types.StringValue(val)
    } else {
        data.SubscriberNotificationStatusMessage = types.StringNull()
    }
    if obj, ok := item["subscriberNotificationStatusOnNoteUpdated"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SubscriberNotificationStatusOnNoteUpdated = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SubscriberNotificationStatusOnNoteUpdated = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SubscriberNotificationStatusOnNoteUpdated = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SubscriberNotificationStatusOnNoteUpdated = types.StringValue(string(jsonBytes))
        } else {
            data.SubscriberNotificationStatusOnNoteUpdated = types.StringNull()
        }
    } else if val, ok := item["subscriberNotificationStatusOnNoteUpdated"].(string); ok {
        data.SubscriberNotificationStatusOnNoteUpdated = types.StringValue(val)
    } else {
        data.SubscriberNotificationStatusOnNoteUpdated = types.StringNull()
    }
    if obj, ok := item["subscriberNotificationStatusMessageOnNoteUpdated"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SubscriberNotificationStatusMessageOnNoteUpdated = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SubscriberNotificationStatusMessageOnNoteUpdated = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SubscriberNotificationStatusMessageOnNoteUpdated = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SubscriberNotificationStatusMessageOnNoteUpdated = types.StringValue(string(jsonBytes))
        } else {
            data.SubscriberNotificationStatusMessageOnNoteUpdated = types.StringNull()
        }
    } else if val, ok := item["subscriberNotificationStatusMessageOnNoteUpdated"].(string); ok {
        data.SubscriberNotificationStatusMessageOnNoteUpdated = types.StringValue(val)
    } else {
        data.SubscriberNotificationStatusMessageOnNoteUpdated = types.StringNull()
    }
    if val, ok := item["shouldStatusPageSubscribersBeNotifiedOnNoteCreated"].(bool); ok {
        data.ShouldStatusPageSubscribersBeNotifiedOnNoteCreated = types.BoolValue(val)
    } else {
        data.ShouldStatusPageSubscribersBeNotifiedOnNoteCreated = types.BoolNull()
    }
    if val, ok := item["isOwnerNotified"].(bool); ok {
        data.IsOwnerNotified = types.BoolValue(val)
    } else {
        data.IsOwnerNotified = types.BoolNull()
    }
    if obj, ok := item["postedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PostedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.PostedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.PostedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.PostedAt = types.StringValue(string(jsonBytes))
        } else {
            data.PostedAt = types.StringNull()
        }
    } else if val, ok := item["postedAt"].(string); ok {
        data.PostedAt = types.StringValue(val)
    } else {
        data.PostedAt = types.StringNull()
    }
    if obj, ok := item["postedFromSlackMessageId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PostedFromSlackMessageId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.PostedFromSlackMessageId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.PostedFromSlackMessageId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.PostedFromSlackMessageId = types.StringValue(string(jsonBytes))
        } else {
            data.PostedFromSlackMessageId = types.StringNull()
        }
    } else if val, ok := item["postedFromSlackMessageId"].(string); ok {
        data.PostedFromSlackMessageId = types.StringValue(val)
    } else {
        data.PostedFromSlackMessageId = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
