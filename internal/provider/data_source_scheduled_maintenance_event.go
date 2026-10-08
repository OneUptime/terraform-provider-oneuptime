package provider

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "math/big"
    "github.com/hashicorp/terraform-plugin-framework/attr"
    "sort"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ScheduledMaintenanceEventDataSource{}

func NewScheduledMaintenanceEventDataSource() datasource.DataSource {
    return &ScheduledMaintenanceEventDataSource{}
}

// ScheduledMaintenanceEventDataSource defines the data source implementation.
type ScheduledMaintenanceEventDataSource struct {
    client *Client
}

// ScheduledMaintenanceEventDataSourceModel describes the data source data model.
type ScheduledMaintenanceEventDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    Title types.String `tfsdk:"title"`
    Description types.String `tfsdk:"description"`
    Slug types.String `tfsdk:"slug"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    Monitors types.Set `tfsdk:"monitors"`
    Hosts types.Set `tfsdk:"hosts"`
    KubernetesClusters types.Set `tfsdk:"kubernetes_clusters"`
    DockerHosts types.Set `tfsdk:"docker_hosts"`
    PodmanHosts types.Set `tfsdk:"podman_hosts"`
    ProxmoxClusters types.Set `tfsdk:"proxmox_clusters"`
    VmwareVCenters types.Set `tfsdk:"vmware_v_centers"`
    IotFleets types.Set `tfsdk:"iot_fleets"`
    NetworkSites types.Set `tfsdk:"network_sites"`
    DockerSwarmClusters types.Set `tfsdk:"docker_swarm_clusters"`
    CephClusters types.Set `tfsdk:"ceph_clusters"`
    StorageArrays types.Set `tfsdk:"storage_arrays"`
    DatabaseServers types.Set `tfsdk:"database_servers"`
    Services types.Set `tfsdk:"services"`
    StatusPages types.Set `tfsdk:"status_pages"`
    Labels types.Set `tfsdk:"labels"`
    CurrentScheduledMaintenanceStateId types.String `tfsdk:"current_scheduled_maintenance_state_id"`
    ChangeMonitorStatusToId types.String `tfsdk:"change_monitor_status_to_id"`
    StartsAt types.String `tfsdk:"starts_at"`
    EndsAt types.String `tfsdk:"ends_at"`
    SubscriberNotificationStatusOnEventScheduled types.String `tfsdk:"subscriber_notification_status_on_event_scheduled"`
    SubscriberNotificationStatusMessage types.String `tfsdk:"subscriber_notification_status_message"`
    ShouldStatusPageSubscribersBeNotifiedOnEventCreated types.Bool `tfsdk:"should_status_page_subscribers_be_notified_on_event_created"`
    ShouldStatusPageSubscribersBeNotifiedWhenEventChangedToOngoing types.Bool `tfsdk:"should_status_page_subscribers_be_notified_when_event_changed_to_ongoing"`
    ShouldStatusPageSubscribersBeNotifiedWhenEventChangedToEnded types.Bool `tfsdk:"should_status_page_subscribers_be_notified_when_event_changed_to_ended"`
    CustomFields types.String `tfsdk:"custom_fields"`
    IsOwnerNotifiedOfResourceCreation types.Bool `tfsdk:"is_owner_notified_of_resource_creation"`
    SendSubscriberNotificationsOnBeforeTheEvent types.String `tfsdk:"send_subscriber_notifications_on_before_the_event"`
    NextSubscriberNotificationBeforeTheEventAt types.String `tfsdk:"next_subscriber_notification_before_the_event_at"`
    ScheduledMaintenanceNumber types.Number `tfsdk:"scheduled_maintenance_number"`
    ScheduledMaintenanceNumberWithPrefix types.String `tfsdk:"scheduled_maintenance_number_with_prefix"`
    IsVisibleOnStatusPage types.Bool `tfsdk:"is_visible_on_status_page"`
    EnableReminders types.Bool `tfsdk:"enable_reminders"`
    NextReminderNotificationAt types.String `tfsdk:"next_reminder_notification_at"`
    ReminderNotificationSentCount types.Number `tfsdk:"reminder_notification_sent_count"`
}

func (d *ScheduledMaintenanceEventDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_scheduled_maintenance_event"
}

func (d *ScheduledMaintenanceEventDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Manage scheduled maintenance event for your project Look up an existing scheduled maintenance event by `id`, or by any of its other arguments (`change_monitor_status_to_id`, `created_by_user_id`, `current_scheduled_maintenance_state_id`, ...): each one set must match, and exactly one scheduled maintenance event may match them all.",

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
            "title": schema.StringAttribute{
                MarkdownDescription: "Title of this scheduled event.",
                Optional: true,
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Description of this scheduled event that will show up on Status Page. This is in markdown.",
                Optional: true,
                Computed: true,
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Friendly globally unique name for your object.",
                Optional: true,
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "monitors": schema.SetAttribute{
                MarkdownDescription: "List of monitors attached to this event. IDs of `oneuptime_monitor` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "hosts": schema.SetAttribute{
                MarkdownDescription: "List of hosts affected by this event. IDs of `oneuptime_host` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "kubernetes_clusters": schema.SetAttribute{
                MarkdownDescription: "List of Kubernetes clusters affected by this event. IDs of `oneuptime_kubernetes_cluster` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "docker_hosts": schema.SetAttribute{
                MarkdownDescription: "List of Docker hosts affected by this event. IDs of `oneuptime_docker_host` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "podman_hosts": schema.SetAttribute{
                MarkdownDescription: "List of Podman hosts affected by this event. IDs of `oneuptime_podman_host` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "proxmox_clusters": schema.SetAttribute{
                MarkdownDescription: "List of Proxmox clusters affected by this event. IDs of `oneuptime_proxmox_cluster` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "vmware_v_centers": schema.SetAttribute{
                MarkdownDescription: "List of vCenters affected by this event. IDs of `oneuptime_vcenter` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "iot_fleets": schema.SetAttribute{
                MarkdownDescription: "List of IoT fleets affected by this event. IDs of `oneuptime_iot_fleet` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "network_sites": schema.SetAttribute{
                MarkdownDescription: "List of network sites affected by this event. Their descendants are covered too. IDs of `oneuptime_network_site` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "docker_swarm_clusters": schema.SetAttribute{
                MarkdownDescription: "List of Docker Swarm clusters affected by this event. IDs of `oneuptime_docker_swarm_cluster` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "ceph_clusters": schema.SetAttribute{
                MarkdownDescription: "List of Ceph clusters affected by this event. IDs of `oneuptime_ceph_cluster` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "storage_arrays": schema.SetAttribute{
                MarkdownDescription: "List of storage arrays affected by this event. IDs of `oneuptime_storage_array` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "database_servers": schema.SetAttribute{
                MarkdownDescription: "List of databases affected by this event. IDs of `oneuptime_database` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "services": schema.SetAttribute{
                MarkdownDescription: "List of services affected by this event. IDs of `oneuptime_service` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "status_pages": schema.SetAttribute{
                MarkdownDescription: "List of status pages to show this event on. IDs of `oneuptime_status_page` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "labels": schema.SetAttribute{
                MarkdownDescription: "Relation to Labels Array where this object is categorized in. IDs of `oneuptime_label` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "current_scheduled_maintenance_state_id": schema.StringAttribute{
                MarkdownDescription: "Scheduled Maintenance State ID. The state the event currently is in. The ID of a `oneuptime_scheduled_maintenance_state`.",
                Optional: true,
                Computed: true,
            },
            "change_monitor_status_to_id": schema.StringAttribute{
                MarkdownDescription: "Relation to Monitor Status Object ID. The monitors attached to this event change to this status when the event starts, and back to operational when it ends. It can be changed until the event starts. The ID of a `oneuptime_monitor_status`.",
                Optional: true,
                Computed: true,
            },
            "starts_at": schema.StringAttribute{
                MarkdownDescription: "When does this event start?",
                Computed: true,
            },
            "ends_at": schema.StringAttribute{
                MarkdownDescription: "When does this event end?",
                Computed: true,
            },
            "subscriber_notification_status_on_event_scheduled": schema.StringAttribute{
                MarkdownDescription: "Status of notification sent to subscribers when event was scheduled.",
                Optional: true,
                Computed: true,
            },
            "subscriber_notification_status_message": schema.StringAttribute{
                MarkdownDescription: "Status message for subscriber notifications when event is scheduled - includes success messages, failure reasons, or skip reasons.",
                Optional: true,
                Computed: true,
            },
            "should_status_page_subscribers_be_notified_on_event_created": schema.BoolAttribute{
                MarkdownDescription: "Should subscribers be notified about this event creation?",
                Optional: true,
                Computed: true,
            },
            "should_status_page_subscribers_be_notified_when_event_changed_to_ongoing": schema.BoolAttribute{
                MarkdownDescription: "Should subscribers be notified about this event event is changed to ongoing?",
                Optional: true,
                Computed: true,
            },
            "should_status_page_subscribers_be_notified_when_event_changed_to_ended": schema.BoolAttribute{
                MarkdownDescription: "Should subscribers be notified about this event event is changed to ended?",
                Optional: true,
                Computed: true,
            },
            "custom_fields": schema.StringAttribute{
                MarkdownDescription: "Custom Fields on this resource. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "is_owner_notified_of_resource_creation": schema.BoolAttribute{
                MarkdownDescription: "Are owners notified of when this resource is created?",
                Optional: true,
                Computed: true,
            },
            "send_subscriber_notifications_on_before_the_event": schema.StringAttribute{
                MarkdownDescription: "Should subscribers be notified before the event? A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "next_subscriber_notification_before_the_event_at": schema.StringAttribute{
                MarkdownDescription: "When will the next notification to subscribers be sent out?",
                Computed: true,
            },
            "scheduled_maintenance_number": schema.NumberAttribute{
                MarkdownDescription: "Scheduled Maintenance Number.",
                Optional: true,
                Computed: true,
            },
            "scheduled_maintenance_number_with_prefix": schema.StringAttribute{
                MarkdownDescription: "Scheduled maintenance number with prefix (e.g., 'SM-42' or '#42').",
                Optional: true,
                Computed: true,
            },
            "is_visible_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Should this incident be visible on the status page?",
                Optional: true,
                Computed: true,
            },
            "enable_reminders": schema.BoolAttribute{
                MarkdownDescription: "Should reminder notifications be sent to owners while this scheduled maintenance event is still not complete? Reminders are sent based on the reminder rules configured for this project.",
                Optional: true,
                Computed: true,
            },
            "next_reminder_notification_at": schema.StringAttribute{
                MarkdownDescription: "When will the next reminder notification be sent to owners of this scheduled maintenance event? This is set automatically based on the reminder rules configured for this project.",
                Computed: true,
            },
            "reminder_notification_sent_count": schema.NumberAttribute{
                MarkdownDescription: "How many reminder notifications have been sent to owners of this scheduled maintenance event so far.",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *ScheduledMaintenanceEventDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ScheduledMaintenanceEventDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data ScheduledMaintenanceEventDataSourceModel

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
    if !data.Slug.IsNull() && !data.Slug.IsUnknown() {
        filters["slug"] = data.Slug.ValueString()
        filterNames = append(filterNames, "slug = "+fmt.Sprintf("%q", data.Slug.ValueString()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }
    if !data.CurrentScheduledMaintenanceStateId.IsNull() && !data.CurrentScheduledMaintenanceStateId.IsUnknown() {
        filters["currentScheduledMaintenanceStateId"] = data.CurrentScheduledMaintenanceStateId.ValueString()
        filterNames = append(filterNames, "current_scheduled_maintenance_state_id = "+fmt.Sprintf("%q", data.CurrentScheduledMaintenanceStateId.ValueString()))
    }
    if !data.ChangeMonitorStatusToId.IsNull() && !data.ChangeMonitorStatusToId.IsUnknown() {
        filters["changeMonitorStatusToId"] = data.ChangeMonitorStatusToId.ValueString()
        filterNames = append(filterNames, "change_monitor_status_to_id = "+fmt.Sprintf("%q", data.ChangeMonitorStatusToId.ValueString()))
    }
    if !data.SubscriberNotificationStatusOnEventScheduled.IsNull() && !data.SubscriberNotificationStatusOnEventScheduled.IsUnknown() {
        filters["subscriberNotificationStatusOnEventScheduled"] = data.SubscriberNotificationStatusOnEventScheduled.ValueString()
        filterNames = append(filterNames, "subscriber_notification_status_on_event_scheduled = "+fmt.Sprintf("%q", data.SubscriberNotificationStatusOnEventScheduled.ValueString()))
    }
    if !data.SubscriberNotificationStatusMessage.IsNull() && !data.SubscriberNotificationStatusMessage.IsUnknown() {
        filters["subscriberNotificationStatusMessage"] = data.SubscriberNotificationStatusMessage.ValueString()
        filterNames = append(filterNames, "subscriber_notification_status_message = "+fmt.Sprintf("%q", data.SubscriberNotificationStatusMessage.ValueString()))
    }
    if !data.ShouldStatusPageSubscribersBeNotifiedOnEventCreated.IsNull() && !data.ShouldStatusPageSubscribersBeNotifiedOnEventCreated.IsUnknown() {
        filters["shouldStatusPageSubscribersBeNotifiedOnEventCreated"] = data.ShouldStatusPageSubscribersBeNotifiedOnEventCreated.ValueBool()
        filterNames = append(filterNames, "should_status_page_subscribers_be_notified_on_event_created = "+fmt.Sprintf("%t", data.ShouldStatusPageSubscribersBeNotifiedOnEventCreated.ValueBool()))
    }
    if !data.ShouldStatusPageSubscribersBeNotifiedWhenEventChangedToOngoing.IsNull() && !data.ShouldStatusPageSubscribersBeNotifiedWhenEventChangedToOngoing.IsUnknown() {
        filters["shouldStatusPageSubscribersBeNotifiedWhenEventChangedToOngoing"] = data.ShouldStatusPageSubscribersBeNotifiedWhenEventChangedToOngoing.ValueBool()
        filterNames = append(filterNames, "should_status_page_subscribers_be_notified_when_event_changed_to_ongoing = "+fmt.Sprintf("%t", data.ShouldStatusPageSubscribersBeNotifiedWhenEventChangedToOngoing.ValueBool()))
    }
    if !data.ShouldStatusPageSubscribersBeNotifiedWhenEventChangedToEnded.IsNull() && !data.ShouldStatusPageSubscribersBeNotifiedWhenEventChangedToEnded.IsUnknown() {
        filters["shouldStatusPageSubscribersBeNotifiedWhenEventChangedToEnded"] = data.ShouldStatusPageSubscribersBeNotifiedWhenEventChangedToEnded.ValueBool()
        filterNames = append(filterNames, "should_status_page_subscribers_be_notified_when_event_changed_to_ended = "+fmt.Sprintf("%t", data.ShouldStatusPageSubscribersBeNotifiedWhenEventChangedToEnded.ValueBool()))
    }
    if !data.IsOwnerNotifiedOfResourceCreation.IsNull() && !data.IsOwnerNotifiedOfResourceCreation.IsUnknown() {
        filters["isOwnerNotifiedOfResourceCreation"] = data.IsOwnerNotifiedOfResourceCreation.ValueBool()
        filterNames = append(filterNames, "is_owner_notified_of_resource_creation = "+fmt.Sprintf("%t", data.IsOwnerNotifiedOfResourceCreation.ValueBool()))
    }
    if !data.ScheduledMaintenanceNumber.IsNull() && !data.ScheduledMaintenanceNumber.IsUnknown() {
        filters["scheduledMaintenanceNumber"] = lookupNumber(data.ScheduledMaintenanceNumber)
        filterNames = append(filterNames, "scheduled_maintenance_number = "+data.ScheduledMaintenanceNumber.ValueBigFloat().String())
    }
    if !data.ScheduledMaintenanceNumberWithPrefix.IsNull() && !data.ScheduledMaintenanceNumberWithPrefix.IsUnknown() {
        filters["scheduledMaintenanceNumberWithPrefix"] = data.ScheduledMaintenanceNumberWithPrefix.ValueString()
        filterNames = append(filterNames, "scheduled_maintenance_number_with_prefix = "+fmt.Sprintf("%q", data.ScheduledMaintenanceNumberWithPrefix.ValueString()))
    }
    if !data.IsVisibleOnStatusPage.IsNull() && !data.IsVisibleOnStatusPage.IsUnknown() {
        filters["isVisibleOnStatusPage"] = data.IsVisibleOnStatusPage.ValueBool()
        filterNames = append(filterNames, "is_visible_on_status_page = "+fmt.Sprintf("%t", data.IsVisibleOnStatusPage.ValueBool()))
    }
    if !data.EnableReminders.IsNull() && !data.EnableReminders.IsUnknown() {
        filters["enableReminders"] = data.EnableReminders.ValueBool()
        filterNames = append(filterNames, "enable_reminders = "+fmt.Sprintf("%t", data.EnableReminders.ValueBool()))
    }
    if !data.ReminderNotificationSentCount.IsNull() && !data.ReminderNotificationSentCount.IsUnknown() {
        filters["reminderNotificationSentCount"] = lookupNumber(data.ReminderNotificationSentCount)
        filterNames = append(filterNames, "reminder_notification_sent_count = "+data.ReminderNotificationSentCount.ValueBigFloat().String())
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the scheduled maintenance event up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the scheduled maintenance event up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "title": true,
        "description": true,
        "slug": true,
        "createdByUserId": true,
        "monitors": true,
        "hosts": true,
        "kubernetesClusters": true,
        "dockerHosts": true,
        "podmanHosts": true,
        "proxmoxClusters": true,
        "vmwareVCenters": true,
        "iotFleets": true,
        "networkSites": true,
        "dockerSwarmClusters": true,
        "cephClusters": true,
        "storageArrays": true,
        "databaseServers": true,
        "services": true,
        "statusPages": true,
        "labels": true,
        "currentScheduledMaintenanceStateId": true,
        "changeMonitorStatusToId": true,
        "startsAt": true,
        "endsAt": true,
        "subscriberNotificationStatusOnEventScheduled": true,
        "subscriberNotificationStatusMessage": true,
        "shouldStatusPageSubscribersBeNotifiedOnEventCreated": true,
        "shouldStatusPageSubscribersBeNotifiedWhenEventChangedToOngoing": true,
        "shouldStatusPageSubscribersBeNotifiedWhenEventChangedToEnded": true,
        "customFields": true,
        "isOwnerNotifiedOfResourceCreation": true,
        "sendSubscriberNotificationsOnBeforeTheEvent": true,
        "nextSubscriberNotificationBeforeTheEventAt": true,
        "scheduledMaintenanceNumber": true,
        "scheduledMaintenanceNumberWithPrefix": true,
        "isVisibleOnStatusPage": true,
        "enableReminders": true,
        "nextReminderNotificationAt": true,
        "reminderNotificationSentCount": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/scheduled-maintenance/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read scheduled_maintenance_event, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No scheduled maintenance event found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read scheduled_maintenance_event: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/scheduled-maintenance/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list scheduled_maintenance_event, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list scheduled_maintenance_event: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No scheduled maintenance event matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one scheduled maintenance event matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for scheduled_maintenance_event.")
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
    if obj, ok := item["slug"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Slug = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Slug = types.StringValue(string(jsonBytes))
        } else {
            data.Slug = types.StringNull()
        }
    } else if val, ok := item["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
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
    if val, ok := item["hosts"].([]interface{}); ok {
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
        data.Hosts = types.SetValueMust(types.StringType, setItems)
    } else {
        data.Hosts = types.SetNull(types.StringType)
    }
    if val, ok := item["kubernetesClusters"].([]interface{}); ok {
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
        data.KubernetesClusters = types.SetValueMust(types.StringType, setItems)
    } else {
        data.KubernetesClusters = types.SetNull(types.StringType)
    }
    if val, ok := item["dockerHosts"].([]interface{}); ok {
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
        data.DockerHosts = types.SetValueMust(types.StringType, setItems)
    } else {
        data.DockerHosts = types.SetNull(types.StringType)
    }
    if val, ok := item["podmanHosts"].([]interface{}); ok {
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
        data.PodmanHosts = types.SetValueMust(types.StringType, setItems)
    } else {
        data.PodmanHosts = types.SetNull(types.StringType)
    }
    if val, ok := item["proxmoxClusters"].([]interface{}); ok {
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
        data.ProxmoxClusters = types.SetValueMust(types.StringType, setItems)
    } else {
        data.ProxmoxClusters = types.SetNull(types.StringType)
    }
    if val, ok := item["vmwareVCenters"].([]interface{}); ok {
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
        data.VmwareVCenters = types.SetValueMust(types.StringType, setItems)
    } else {
        data.VmwareVCenters = types.SetNull(types.StringType)
    }
    if val, ok := item["iotFleets"].([]interface{}); ok {
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
        data.IotFleets = types.SetValueMust(types.StringType, setItems)
    } else {
        data.IotFleets = types.SetNull(types.StringType)
    }
    if val, ok := item["networkSites"].([]interface{}); ok {
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
        data.NetworkSites = types.SetValueMust(types.StringType, setItems)
    } else {
        data.NetworkSites = types.SetNull(types.StringType)
    }
    if val, ok := item["dockerSwarmClusters"].([]interface{}); ok {
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
        data.DockerSwarmClusters = types.SetValueMust(types.StringType, setItems)
    } else {
        data.DockerSwarmClusters = types.SetNull(types.StringType)
    }
    if val, ok := item["cephClusters"].([]interface{}); ok {
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
        data.CephClusters = types.SetValueMust(types.StringType, setItems)
    } else {
        data.CephClusters = types.SetNull(types.StringType)
    }
    if val, ok := item["storageArrays"].([]interface{}); ok {
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
        data.StorageArrays = types.SetValueMust(types.StringType, setItems)
    } else {
        data.StorageArrays = types.SetNull(types.StringType)
    }
    if val, ok := item["databaseServers"].([]interface{}); ok {
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
        data.DatabaseServers = types.SetValueMust(types.StringType, setItems)
    } else {
        data.DatabaseServers = types.SetNull(types.StringType)
    }
    if val, ok := item["services"].([]interface{}); ok {
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
        data.Services = types.SetValueMust(types.StringType, setItems)
    } else {
        data.Services = types.SetNull(types.StringType)
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
    if val, ok := item["labels"].([]interface{}); ok {
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
        data.Labels = types.SetValueMust(types.StringType, setItems)
    } else {
        data.Labels = types.SetNull(types.StringType)
    }
    if obj, ok := item["currentScheduledMaintenanceStateId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CurrentScheduledMaintenanceStateId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CurrentScheduledMaintenanceStateId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CurrentScheduledMaintenanceStateId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CurrentScheduledMaintenanceStateId = types.StringValue(string(jsonBytes))
        } else {
            data.CurrentScheduledMaintenanceStateId = types.StringNull()
        }
    } else if val, ok := item["currentScheduledMaintenanceStateId"].(string); ok {
        data.CurrentScheduledMaintenanceStateId = types.StringValue(val)
    } else {
        data.CurrentScheduledMaintenanceStateId = types.StringNull()
    }
    if obj, ok := item["changeMonitorStatusToId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ChangeMonitorStatusToId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ChangeMonitorStatusToId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ChangeMonitorStatusToId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ChangeMonitorStatusToId = types.StringValue(string(jsonBytes))
        } else {
            data.ChangeMonitorStatusToId = types.StringNull()
        }
    } else if val, ok := item["changeMonitorStatusToId"].(string); ok {
        data.ChangeMonitorStatusToId = types.StringValue(val)
    } else {
        data.ChangeMonitorStatusToId = types.StringNull()
    }
    if obj, ok := item["startsAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StartsAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.StartsAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.StartsAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.StartsAt = types.StringValue(string(jsonBytes))
        } else {
            data.StartsAt = types.StringNull()
        }
    } else if val, ok := item["startsAt"].(string); ok {
        data.StartsAt = types.StringValue(val)
    } else {
        data.StartsAt = types.StringNull()
    }
    if obj, ok := item["endsAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EndsAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.EndsAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.EndsAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.EndsAt = types.StringValue(string(jsonBytes))
        } else {
            data.EndsAt = types.StringNull()
        }
    } else if val, ok := item["endsAt"].(string); ok {
        data.EndsAt = types.StringValue(val)
    } else {
        data.EndsAt = types.StringNull()
    }
    if obj, ok := item["subscriberNotificationStatusOnEventScheduled"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SubscriberNotificationStatusOnEventScheduled = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SubscriberNotificationStatusOnEventScheduled = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SubscriberNotificationStatusOnEventScheduled = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SubscriberNotificationStatusOnEventScheduled = types.StringValue(string(jsonBytes))
        } else {
            data.SubscriberNotificationStatusOnEventScheduled = types.StringNull()
        }
    } else if val, ok := item["subscriberNotificationStatusOnEventScheduled"].(string); ok {
        data.SubscriberNotificationStatusOnEventScheduled = types.StringValue(val)
    } else {
        data.SubscriberNotificationStatusOnEventScheduled = types.StringNull()
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
    if val, ok := item["shouldStatusPageSubscribersBeNotifiedOnEventCreated"].(bool); ok {
        data.ShouldStatusPageSubscribersBeNotifiedOnEventCreated = types.BoolValue(val)
    } else {
        data.ShouldStatusPageSubscribersBeNotifiedOnEventCreated = types.BoolNull()
    }
    if val, ok := item["shouldStatusPageSubscribersBeNotifiedWhenEventChangedToOngoing"].(bool); ok {
        data.ShouldStatusPageSubscribersBeNotifiedWhenEventChangedToOngoing = types.BoolValue(val)
    } else {
        data.ShouldStatusPageSubscribersBeNotifiedWhenEventChangedToOngoing = types.BoolNull()
    }
    if val, ok := item["shouldStatusPageSubscribersBeNotifiedWhenEventChangedToEnded"].(bool); ok {
        data.ShouldStatusPageSubscribersBeNotifiedWhenEventChangedToEnded = types.BoolValue(val)
    } else {
        data.ShouldStatusPageSubscribersBeNotifiedWhenEventChangedToEnded = types.BoolNull()
    }
    if obj, ok := item["customFields"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CustomFields = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CustomFields = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CustomFields = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CustomFields = types.StringValue(string(jsonBytes))
        } else {
            data.CustomFields = types.StringNull()
        }
    } else if val, ok := item["customFields"].(string); ok {
        data.CustomFields = types.StringValue(val)
    } else {
        data.CustomFields = types.StringNull()
    }
    if val, ok := item["isOwnerNotifiedOfResourceCreation"].(bool); ok {
        data.IsOwnerNotifiedOfResourceCreation = types.BoolValue(val)
    } else {
        data.IsOwnerNotifiedOfResourceCreation = types.BoolNull()
    }
    if obj, ok := item["sendSubscriberNotificationsOnBeforeTheEvent"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SendSubscriberNotificationsOnBeforeTheEvent = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SendSubscriberNotificationsOnBeforeTheEvent = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SendSubscriberNotificationsOnBeforeTheEvent = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SendSubscriberNotificationsOnBeforeTheEvent = types.StringValue(string(jsonBytes))
        } else {
            data.SendSubscriberNotificationsOnBeforeTheEvent = types.StringNull()
        }
    } else if val, ok := item["sendSubscriberNotificationsOnBeforeTheEvent"].(string); ok {
        data.SendSubscriberNotificationsOnBeforeTheEvent = types.StringValue(val)
    } else {
        data.SendSubscriberNotificationsOnBeforeTheEvent = types.StringNull()
    }
    if obj, ok := item["nextSubscriberNotificationBeforeTheEventAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.NextSubscriberNotificationBeforeTheEventAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.NextSubscriberNotificationBeforeTheEventAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.NextSubscriberNotificationBeforeTheEventAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.NextSubscriberNotificationBeforeTheEventAt = types.StringValue(string(jsonBytes))
        } else {
            data.NextSubscriberNotificationBeforeTheEventAt = types.StringNull()
        }
    } else if val, ok := item["nextSubscriberNotificationBeforeTheEventAt"].(string); ok {
        data.NextSubscriberNotificationBeforeTheEventAt = types.StringValue(val)
    } else {
        data.NextSubscriberNotificationBeforeTheEventAt = types.StringNull()
    }
    if val, ok := item["scheduledMaintenanceNumber"].(float64); ok {
        data.ScheduledMaintenanceNumber = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["scheduledMaintenanceNumber"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.ScheduledMaintenanceNumber = types.NumberValue(big.NewFloat(val))
        } else {
            data.ScheduledMaintenanceNumber = types.NumberNull()
        }
    } else {
        data.ScheduledMaintenanceNumber = types.NumberNull()
    }
    if obj, ok := item["scheduledMaintenanceNumberWithPrefix"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ScheduledMaintenanceNumberWithPrefix = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ScheduledMaintenanceNumberWithPrefix = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ScheduledMaintenanceNumberWithPrefix = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ScheduledMaintenanceNumberWithPrefix = types.StringValue(string(jsonBytes))
        } else {
            data.ScheduledMaintenanceNumberWithPrefix = types.StringNull()
        }
    } else if val, ok := item["scheduledMaintenanceNumberWithPrefix"].(string); ok {
        data.ScheduledMaintenanceNumberWithPrefix = types.StringValue(val)
    } else {
        data.ScheduledMaintenanceNumberWithPrefix = types.StringNull()
    }
    if val, ok := item["isVisibleOnStatusPage"].(bool); ok {
        data.IsVisibleOnStatusPage = types.BoolValue(val)
    } else {
        data.IsVisibleOnStatusPage = types.BoolNull()
    }
    if val, ok := item["enableReminders"].(bool); ok {
        data.EnableReminders = types.BoolValue(val)
    } else {
        data.EnableReminders = types.BoolNull()
    }
    if obj, ok := item["nextReminderNotificationAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.NextReminderNotificationAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.NextReminderNotificationAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.NextReminderNotificationAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.NextReminderNotificationAt = types.StringValue(string(jsonBytes))
        } else {
            data.NextReminderNotificationAt = types.StringNull()
        }
    } else if val, ok := item["nextReminderNotificationAt"].(string); ok {
        data.NextReminderNotificationAt = types.StringValue(val)
    } else {
        data.NextReminderNotificationAt = types.StringNull()
    }
    if val, ok := item["reminderNotificationSentCount"].(float64); ok {
        data.ReminderNotificationSentCount = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["reminderNotificationSentCount"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.ReminderNotificationSentCount = types.NumberValue(big.NewFloat(val))
        } else {
            data.ReminderNotificationSentCount = types.NumberNull()
        }
    } else {
        data.ReminderNotificationSentCount = types.NumberNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
