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
var _ datasource.DataSource = &ScheduledMaintenanceTemplateDataSource{}

func NewScheduledMaintenanceTemplateDataSource() datasource.DataSource {
    return &ScheduledMaintenanceTemplateDataSource{}
}

// ScheduledMaintenanceTemplateDataSource defines the data source implementation.
type ScheduledMaintenanceTemplateDataSource struct {
    client *Client
}

// ScheduledMaintenanceTemplateDataSourceModel describes the data source data model.
type ScheduledMaintenanceTemplateDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    TemplateName types.String `tfsdk:"template_name"`
    TemplateDescription types.String `tfsdk:"template_description"`
    Title types.String `tfsdk:"title"`
    Description types.String `tfsdk:"description"`
    Slug types.String `tfsdk:"slug"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    Monitors types.Set `tfsdk:"monitors"`
    Hosts types.Set `tfsdk:"hosts"`
    KubernetesClusters types.Set `tfsdk:"kubernetes_clusters"`
    DockerHosts types.Set `tfsdk:"docker_hosts"`
    PodmanHosts types.Set `tfsdk:"podman_hosts"`
    Services types.Set `tfsdk:"services"`
    StatusPages types.Set `tfsdk:"status_pages"`
    Labels types.Set `tfsdk:"labels"`
    ChangeMonitorStatusToId types.String `tfsdk:"change_monitor_status_to_id"`
    FirstEventScheduledAt types.String `tfsdk:"first_event_scheduled_at"`
    FirstEventStartsAt types.String `tfsdk:"first_event_starts_at"`
    FirstEventEndsAt types.String `tfsdk:"first_event_ends_at"`
    RecurringInterval types.String `tfsdk:"recurring_interval"`
    IsRecurringEvent types.Bool `tfsdk:"is_recurring_event"`
    ScheduleNextEventAt types.String `tfsdk:"schedule_next_event_at"`
    ShouldStatusPageSubscribersBeNotifiedOnEventCreated types.Bool `tfsdk:"should_status_page_subscribers_be_notified_on_event_created"`
    ShouldStatusPageSubscribersBeNotifiedWhenEventChangedToOngoing types.Bool `tfsdk:"should_status_page_subscribers_be_notified_when_event_changed_to_ongoing"`
    ShouldStatusPageSubscribersBeNotifiedWhenEventChangedToEnded types.Bool `tfsdk:"should_status_page_subscribers_be_notified_when_event_changed_to_ended"`
    CustomFields types.String `tfsdk:"custom_fields"`
    SendSubscriberNotificationsOnBeforeTheEvent types.String `tfsdk:"send_subscriber_notifications_on_before_the_event"`
}

func (d *ScheduledMaintenanceTemplateDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_scheduled_maintenance_template"
}

func (d *ScheduledMaintenanceTemplateDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Manage scheduled maintenance templates for your project Look up an existing scheduled maintenance template by `id`, or by any of its other arguments (`change_monitor_status_to_id`, `created_by_user_id`, `description`, ...): each one set must match, and exactly one scheduled maintenance template may match them all.",

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
            "template_name": schema.StringAttribute{
                MarkdownDescription: "Name of the Scheduled Maintenance Template.",
                Optional: true,
                Computed: true,
            },
            "template_description": schema.StringAttribute{
                MarkdownDescription: "Description of the Scheduled Maintenance Template.",
                Optional: true,
                Computed: true,
            },
            "title": schema.StringAttribute{
                MarkdownDescription: "Title of this scheduled event.",
                Optional: true,
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Description of this scheduled event that will show up on Status Page. This is a markdown field.",
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
                MarkdownDescription: "List of hosts to pre-populate on scheduled maintenance events created from this template. IDs of `oneuptime_host` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "kubernetes_clusters": schema.SetAttribute{
                MarkdownDescription: "List of Kubernetes clusters to pre-populate on scheduled maintenance events created from this template. IDs of `oneuptime_kubernetes_cluster` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "docker_hosts": schema.SetAttribute{
                MarkdownDescription: "List of Docker hosts to pre-populate on scheduled maintenance events created from this template. IDs of `oneuptime_docker_host` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "podman_hosts": schema.SetAttribute{
                MarkdownDescription: "List of Podman hosts to pre-populate on scheduled maintenance events created from this template. IDs of `oneuptime_podman_host` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "services": schema.SetAttribute{
                MarkdownDescription: "List of services to pre-populate on scheduled maintenance events created from this template. IDs of `oneuptime_service` resources.",
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
            "change_monitor_status_to_id": schema.StringAttribute{
                MarkdownDescription: "Relation to Monitor Status Object ID. All monitors connected to this incident will be changed to this status when the event is ongoing. The ID of a `oneuptime_monitor_status`.",
                Optional: true,
                Computed: true,
            },
            "first_event_scheduled_at": schema.StringAttribute{
                MarkdownDescription: "When would you like to schedule the first event?",
                Computed: true,
            },
            "first_event_starts_at": schema.StringAttribute{
                MarkdownDescription: "When does the first event start?",
                Computed: true,
            },
            "first_event_ends_at": schema.StringAttribute{
                MarkdownDescription: "When does the first event end?",
                Computed: true,
            },
            "recurring_interval": schema.StringAttribute{
                MarkdownDescription: "How often should this event recur? A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "is_recurring_event": schema.BoolAttribute{
                MarkdownDescription: "Is this a recurring event?",
                Optional: true,
                Computed: true,
            },
            "schedule_next_event_at": schema.StringAttribute{
                MarkdownDescription: "When is the next event scheduled?",
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
            "send_subscriber_notifications_on_before_the_event": schema.StringAttribute{
                MarkdownDescription: "Should subscribers be notified before the event? A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
        },
    }
}

func (d *ScheduledMaintenanceTemplateDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ScheduledMaintenanceTemplateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data ScheduledMaintenanceTemplateDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.TemplateName.IsNull() && !data.TemplateName.IsUnknown() {
        filters["templateName"] = data.TemplateName.ValueString()
        filterNames = append(filterNames, "template_name = "+fmt.Sprintf("%q", data.TemplateName.ValueString()))
    }
    if !data.TemplateDescription.IsNull() && !data.TemplateDescription.IsUnknown() {
        filters["templateDescription"] = data.TemplateDescription.ValueString()
        filterNames = append(filterNames, "template_description = "+fmt.Sprintf("%q", data.TemplateDescription.ValueString()))
    }
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
    if !data.ChangeMonitorStatusToId.IsNull() && !data.ChangeMonitorStatusToId.IsUnknown() {
        filters["changeMonitorStatusToId"] = data.ChangeMonitorStatusToId.ValueString()
        filterNames = append(filterNames, "change_monitor_status_to_id = "+fmt.Sprintf("%q", data.ChangeMonitorStatusToId.ValueString()))
    }
    if !data.IsRecurringEvent.IsNull() && !data.IsRecurringEvent.IsUnknown() {
        filters["isRecurringEvent"] = data.IsRecurringEvent.ValueBool()
        filterNames = append(filterNames, "is_recurring_event = "+fmt.Sprintf("%t", data.IsRecurringEvent.ValueBool()))
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

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the scheduled maintenance template up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the scheduled maintenance template up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "templateName": true,
        "templateDescription": true,
        "title": true,
        "description": true,
        "slug": true,
        "createdByUserId": true,
        "monitors": true,
        "hosts": true,
        "kubernetesClusters": true,
        "dockerHosts": true,
        "podmanHosts": true,
        "services": true,
        "statusPages": true,
        "labels": true,
        "changeMonitorStatusToId": true,
        "firstEventScheduledAt": true,
        "firstEventStartsAt": true,
        "firstEventEndsAt": true,
        "recurringInterval": true,
        "isRecurringEvent": true,
        "scheduleNextEventAt": true,
        "shouldStatusPageSubscribersBeNotifiedOnEventCreated": true,
        "shouldStatusPageSubscribersBeNotifiedWhenEventChangedToOngoing": true,
        "shouldStatusPageSubscribersBeNotifiedWhenEventChangedToEnded": true,
        "customFields": true,
        "sendSubscriberNotificationsOnBeforeTheEvent": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/scheduled-maintenance-template/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read scheduled_maintenance_template, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No scheduled maintenance template found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read scheduled_maintenance_template: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/scheduled-maintenance-template/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list scheduled_maintenance_template, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list scheduled_maintenance_template: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No scheduled maintenance template matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one scheduled maintenance template matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for scheduled_maintenance_template.")
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
    if obj, ok := item["templateName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TemplateName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TemplateName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TemplateName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TemplateName = types.StringValue(string(jsonBytes))
        } else {
            data.TemplateName = types.StringNull()
        }
    } else if val, ok := item["templateName"].(string); ok {
        data.TemplateName = types.StringValue(val)
    } else {
        data.TemplateName = types.StringNull()
    }
    if obj, ok := item["templateDescription"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TemplateDescription = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TemplateDescription = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TemplateDescription = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TemplateDescription = types.StringValue(string(jsonBytes))
        } else {
            data.TemplateDescription = types.StringNull()
        }
    } else if val, ok := item["templateDescription"].(string); ok {
        data.TemplateDescription = types.StringValue(val)
    } else {
        data.TemplateDescription = types.StringNull()
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
    if obj, ok := item["firstEventScheduledAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FirstEventScheduledAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.FirstEventScheduledAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.FirstEventScheduledAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.FirstEventScheduledAt = types.StringValue(string(jsonBytes))
        } else {
            data.FirstEventScheduledAt = types.StringNull()
        }
    } else if val, ok := item["firstEventScheduledAt"].(string); ok {
        data.FirstEventScheduledAt = types.StringValue(val)
    } else {
        data.FirstEventScheduledAt = types.StringNull()
    }
    if obj, ok := item["firstEventStartsAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FirstEventStartsAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.FirstEventStartsAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.FirstEventStartsAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.FirstEventStartsAt = types.StringValue(string(jsonBytes))
        } else {
            data.FirstEventStartsAt = types.StringNull()
        }
    } else if val, ok := item["firstEventStartsAt"].(string); ok {
        data.FirstEventStartsAt = types.StringValue(val)
    } else {
        data.FirstEventStartsAt = types.StringNull()
    }
    if obj, ok := item["firstEventEndsAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FirstEventEndsAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.FirstEventEndsAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.FirstEventEndsAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.FirstEventEndsAt = types.StringValue(string(jsonBytes))
        } else {
            data.FirstEventEndsAt = types.StringNull()
        }
    } else if val, ok := item["firstEventEndsAt"].(string); ok {
        data.FirstEventEndsAt = types.StringValue(val)
    } else {
        data.FirstEventEndsAt = types.StringNull()
    }
    if obj, ok := item["recurringInterval"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RecurringInterval = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RecurringInterval = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RecurringInterval = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RecurringInterval = types.StringValue(string(jsonBytes))
        } else {
            data.RecurringInterval = types.StringNull()
        }
    } else if val, ok := item["recurringInterval"].(string); ok {
        data.RecurringInterval = types.StringValue(val)
    } else {
        data.RecurringInterval = types.StringNull()
    }
    if val, ok := item["isRecurringEvent"].(bool); ok {
        data.IsRecurringEvent = types.BoolValue(val)
    } else {
        data.IsRecurringEvent = types.BoolNull()
    }
    if obj, ok := item["scheduleNextEventAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ScheduleNextEventAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ScheduleNextEventAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ScheduleNextEventAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ScheduleNextEventAt = types.StringValue(string(jsonBytes))
        } else {
            data.ScheduleNextEventAt = types.StringNull()
        }
    } else if val, ok := item["scheduleNextEventAt"].(string); ok {
        data.ScheduleNextEventAt = types.StringValue(val)
    } else {
        data.ScheduleNextEventAt = types.StringNull()
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

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
