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
var _ datasource.DataSource = &StatusPageAnnouncementDataSource{}

func NewStatusPageAnnouncementDataSource() datasource.DataSource {
    return &StatusPageAnnouncementDataSource{}
}

// StatusPageAnnouncementDataSource defines the data source implementation.
type StatusPageAnnouncementDataSource struct {
    client *Client
}

// StatusPageAnnouncementDataSourceModel describes the data source data model.
type StatusPageAnnouncementDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    StatusPages types.Set `tfsdk:"status_pages"`
    Monitors types.Set `tfsdk:"monitors"`
    Title types.String `tfsdk:"title"`
    ShowAnnouncementAt types.String `tfsdk:"show_announcement_at"`
    EndAnnouncementAt types.String `tfsdk:"end_announcement_at"`
    Description types.String `tfsdk:"description"`
    Attachments types.Set `tfsdk:"attachments"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    SubscriberNotificationStatus types.String `tfsdk:"subscriber_notification_status"`
    SubscriberNotificationStatusMessage types.String `tfsdk:"subscriber_notification_status_message"`
    SubscriberNotificationStatusOnAnnouncementUpdated types.String `tfsdk:"subscriber_notification_status_on_announcement_updated"`
    SubscriberNotificationStatusMessageOnAnnouncementUpdated types.String `tfsdk:"subscriber_notification_status_message_on_announcement_updated"`
    ShouldStatusPageSubscribersBeNotified types.Bool `tfsdk:"should_status_page_subscribers_be_notified"`
    IsOwnerNotified types.Bool `tfsdk:"is_owner_notified"`
}

func (d *StatusPageAnnouncementDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_status_page_announcement"
}

func (d *StatusPageAnnouncementDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Manage announcements on your status page Look up an existing status page announcement by `id`, or by any of its other arguments (`created_by_user_id`, `description`, `is_owner_notified`, ...): each one set must match, and exactly one status page announcement may match them all.",

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
            "status_pages": schema.SetAttribute{
                MarkdownDescription: "Status Pages to show show this announcement on. IDs of `oneuptime_status_page` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "monitors": schema.SetAttribute{
                MarkdownDescription: "List of monitors affected by this announcement. If none are selected, all subscribers will be notified. IDs of `oneuptime_monitor` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "title": schema.StringAttribute{
                MarkdownDescription: "Title of this resource.",
                Optional: true,
                Computed: true,
            },
            "show_announcement_at": schema.StringAttribute{
                MarkdownDescription: "When should this announcement be shown?",
                Computed: true,
            },
            "end_announcement_at": schema.StringAttribute{
                MarkdownDescription: "When should this announcement hidden?",
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Text of the announcement. This can be in Markdown format.",
                Optional: true,
                Computed: true,
            },
            "attachments": schema.SetAttribute{
                MarkdownDescription: "Files attached to this announcement. IDs of `oneuptime_file` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "subscriber_notification_status": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Create Status Page Announcement], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Admin, Status Page Member, Status Page Viewer, Read Status Page Announcement], Update: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Edit Status Page Announcement]",
                Optional: true,
                Computed: true,
            },
            "subscriber_notification_status_message": schema.StringAttribute{
                MarkdownDescription: "Status message for subscriber notifications - includes success messages, failure reasons, or skip reasons.",
                Optional: true,
                Computed: true,
            },
            "subscriber_notification_status_on_announcement_updated": schema.StringAttribute{
                MarkdownDescription: "Status of the notification sent to subscribers when this announcement was last updated. Empty until an update notification is requested.",
                Optional: true,
                Computed: true,
            },
            "subscriber_notification_status_message_on_announcement_updated": schema.StringAttribute{
                MarkdownDescription: "Status message for the notification sent to subscribers when this announcement was last updated - includes success messages, failure reasons, or skip reasons.",
                Optional: true,
                Computed: true,
            },
            "should_status_page_subscribers_be_notified": schema.BoolAttribute{
                MarkdownDescription: "Should subscribers be notified about this announcement?",
                Optional: true,
                Computed: true,
            },
            "is_owner_notified": schema.BoolAttribute{
                MarkdownDescription: "Are owners notified of this announcement?",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *StatusPageAnnouncementDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *StatusPageAnnouncementDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data StatusPageAnnouncementDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.Title.IsNull() && !data.Title.IsUnknown() {
        filters["title"] = data.Title.ValueString()
        filterNames = append(filterNames, "title = "+fmt.Sprintf("%q", data.Title.ValueString()))
    }
    if !data.Description.IsNull() && !data.Description.IsUnknown() {
        filters["description"] = data.Description.ValueString()
        filterNames = append(filterNames, "description = "+fmt.Sprintf("%q", data.Description.ValueString()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }
    if !data.SubscriberNotificationStatus.IsNull() && !data.SubscriberNotificationStatus.IsUnknown() {
        filters["subscriberNotificationStatus"] = data.SubscriberNotificationStatus.ValueString()
        filterNames = append(filterNames, "subscriber_notification_status = "+fmt.Sprintf("%q", data.SubscriberNotificationStatus.ValueString()))
    }
    if !data.SubscriberNotificationStatusMessage.IsNull() && !data.SubscriberNotificationStatusMessage.IsUnknown() {
        filters["subscriberNotificationStatusMessage"] = data.SubscriberNotificationStatusMessage.ValueString()
        filterNames = append(filterNames, "subscriber_notification_status_message = "+fmt.Sprintf("%q", data.SubscriberNotificationStatusMessage.ValueString()))
    }
    if !data.SubscriberNotificationStatusOnAnnouncementUpdated.IsNull() && !data.SubscriberNotificationStatusOnAnnouncementUpdated.IsUnknown() {
        filters["subscriberNotificationStatusOnAnnouncementUpdated"] = data.SubscriberNotificationStatusOnAnnouncementUpdated.ValueString()
        filterNames = append(filterNames, "subscriber_notification_status_on_announcement_updated = "+fmt.Sprintf("%q", data.SubscriberNotificationStatusOnAnnouncementUpdated.ValueString()))
    }
    if !data.SubscriberNotificationStatusMessageOnAnnouncementUpdated.IsNull() && !data.SubscriberNotificationStatusMessageOnAnnouncementUpdated.IsUnknown() {
        filters["subscriberNotificationStatusMessageOnAnnouncementUpdated"] = data.SubscriberNotificationStatusMessageOnAnnouncementUpdated.ValueString()
        filterNames = append(filterNames, "subscriber_notification_status_message_on_announcement_updated = "+fmt.Sprintf("%q", data.SubscriberNotificationStatusMessageOnAnnouncementUpdated.ValueString()))
    }
    if !data.ShouldStatusPageSubscribersBeNotified.IsNull() && !data.ShouldStatusPageSubscribersBeNotified.IsUnknown() {
        filters["shouldStatusPageSubscribersBeNotified"] = data.ShouldStatusPageSubscribersBeNotified.ValueBool()
        filterNames = append(filterNames, "should_status_page_subscribers_be_notified = "+fmt.Sprintf("%t", data.ShouldStatusPageSubscribersBeNotified.ValueBool()))
    }
    if !data.IsOwnerNotified.IsNull() && !data.IsOwnerNotified.IsUnknown() {
        filters["isOwnerNotified"] = data.IsOwnerNotified.ValueBool()
        filterNames = append(filterNames, "is_owner_notified = "+fmt.Sprintf("%t", data.IsOwnerNotified.ValueBool()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the status page announcement up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the status page announcement up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "statusPages": true,
        "monitors": true,
        "title": true,
        "showAnnouncementAt": true,
        "endAnnouncementAt": true,
        "description": true,
        "attachments": true,
        "createdByUserId": true,
        "subscriberNotificationStatus": true,
        "subscriberNotificationStatusMessage": true,
        "subscriberNotificationStatusOnAnnouncementUpdated": true,
        "subscriberNotificationStatusMessageOnAnnouncementUpdated": true,
        "shouldStatusPageSubscribersBeNotified": true,
        "isOwnerNotified": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/status-page-announcement/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read status_page_announcement, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No status page announcement found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read status_page_announcement: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/status-page-announcement/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list status_page_announcement, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list status_page_announcement: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No status page announcement matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one status page announcement matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for status_page_announcement.")
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
    if val, ok := item["statusPages"].([]interface{}); ok {
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
        data.StatusPages = types.SetValueMust(types.StringType, setItems)
    } else {
        data.StatusPages = types.SetNull(types.StringType)
    }
    if val, ok := item["monitors"].([]interface{}); ok {
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
        data.Monitors = types.SetValueMust(types.StringType, setItems)
    } else {
        data.Monitors = types.SetNull(types.StringType)
    }
    if obj, ok := item["title"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Title = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Title = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Title = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Title = types.StringValue(string(jsonBytes))
        } else {
            data.Title = types.StringNull()
        }
    } else if val, ok := item["title"].(string); ok {
        data.Title = types.StringValue(val)
    } else {
        data.Title = types.StringNull()
    }
    if obj, ok := item["showAnnouncementAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ShowAnnouncementAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ShowAnnouncementAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ShowAnnouncementAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ShowAnnouncementAt = types.StringValue(string(jsonBytes))
        } else {
            data.ShowAnnouncementAt = types.StringNull()
        }
    } else if val, ok := item["showAnnouncementAt"].(string); ok {
        data.ShowAnnouncementAt = types.StringValue(val)
    } else {
        data.ShowAnnouncementAt = types.StringNull()
    }
    if obj, ok := item["endAnnouncementAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EndAnnouncementAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.EndAnnouncementAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.EndAnnouncementAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.EndAnnouncementAt = types.StringValue(string(jsonBytes))
        } else {
            data.EndAnnouncementAt = types.StringNull()
        }
    } else if val, ok := item["endAnnouncementAt"].(string); ok {
        data.EndAnnouncementAt = types.StringValue(val)
    } else {
        data.EndAnnouncementAt = types.StringNull()
    }
    if obj, ok := item["description"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Description = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Description = types.StringValue(string(jsonBytes))
        } else {
            data.Description = types.StringNull()
        }
    } else if val, ok := item["description"].(string); ok {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
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
    if obj, ok := item["subscriberNotificationStatus"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SubscriberNotificationStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SubscriberNotificationStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SubscriberNotificationStatus = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SubscriberNotificationStatus = types.StringValue(string(jsonBytes))
        } else {
            data.SubscriberNotificationStatus = types.StringNull()
        }
    } else if val, ok := item["subscriberNotificationStatus"].(string); ok {
        data.SubscriberNotificationStatus = types.StringValue(val)
    } else {
        data.SubscriberNotificationStatus = types.StringNull()
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
    if obj, ok := item["subscriberNotificationStatusOnAnnouncementUpdated"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SubscriberNotificationStatusOnAnnouncementUpdated = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SubscriberNotificationStatusOnAnnouncementUpdated = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SubscriberNotificationStatusOnAnnouncementUpdated = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SubscriberNotificationStatusOnAnnouncementUpdated = types.StringValue(string(jsonBytes))
        } else {
            data.SubscriberNotificationStatusOnAnnouncementUpdated = types.StringNull()
        }
    } else if val, ok := item["subscriberNotificationStatusOnAnnouncementUpdated"].(string); ok {
        data.SubscriberNotificationStatusOnAnnouncementUpdated = types.StringValue(val)
    } else {
        data.SubscriberNotificationStatusOnAnnouncementUpdated = types.StringNull()
    }
    if obj, ok := item["subscriberNotificationStatusMessageOnAnnouncementUpdated"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SubscriberNotificationStatusMessageOnAnnouncementUpdated = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SubscriberNotificationStatusMessageOnAnnouncementUpdated = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SubscriberNotificationStatusMessageOnAnnouncementUpdated = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SubscriberNotificationStatusMessageOnAnnouncementUpdated = types.StringValue(string(jsonBytes))
        } else {
            data.SubscriberNotificationStatusMessageOnAnnouncementUpdated = types.StringNull()
        }
    } else if val, ok := item["subscriberNotificationStatusMessageOnAnnouncementUpdated"].(string); ok {
        data.SubscriberNotificationStatusMessageOnAnnouncementUpdated = types.StringValue(val)
    } else {
        data.SubscriberNotificationStatusMessageOnAnnouncementUpdated = types.StringNull()
    }
    if val, ok := item["shouldStatusPageSubscribersBeNotified"].(bool); ok {
        data.ShouldStatusPageSubscribersBeNotified = types.BoolValue(val)
    } else {
        data.ShouldStatusPageSubscribersBeNotified = types.BoolNull()
    }
    if val, ok := item["isOwnerNotified"].(bool); ok {
        data.IsOwnerNotified = types.BoolValue(val)
    } else {
        data.IsOwnerNotified = types.BoolNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
