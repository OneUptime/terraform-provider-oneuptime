package provider

import (
    "context"
    "fmt"
    "math/big"
    "github.com/hashicorp/terraform-plugin-framework/attr"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ScheduledMaintenanceTemplateDataDataSource{}

func NewScheduledMaintenanceTemplateDataDataSource() datasource.DataSource {
    return &ScheduledMaintenanceTemplateDataDataSource{}
}

// ScheduledMaintenanceTemplateDataDataSource defines the data source implementation.
type ScheduledMaintenanceTemplateDataDataSource struct {
    client *Client
}

// ScheduledMaintenanceTemplateDataDataSourceModel describes the data source data model.
type ScheduledMaintenanceTemplateDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
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

func (d *ScheduledMaintenanceTemplateDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_scheduled_maintenance_template_data"
}

func (d *ScheduledMaintenanceTemplateDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "scheduled_maintenance_template_data data source",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Identifier to filter by",
                Optional: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Name to filter by",
                Optional: true,
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
            "template_name": schema.StringAttribute{
                MarkdownDescription: "Name of the Scheduled Maintenance Template. Permissions - Create: [Project Owner, Project Admin, Project Member, Scheduled Maintenance Admin, Scheduled Maintenance Member, Create Scheduled Maintenance Template], Read: [Project Owner, Project Admin, Project Member, Viewer, Scheduled Maintenance Admin, Scheduled Maintenance Member, Scheduled Maintenance Viewer, Read Scheduled Maintenance Template], Update: [Project Owner, Project Admin, Project Member, Scheduled Maintenance Admin, Scheduled Maintenance Member, Edit Scheduled Maintenance Template]",
                Computed: true,
            },
            "template_description": schema.StringAttribute{
                MarkdownDescription: "Description of the Scheduled Maintenance Template. Permissions - Create: [Project Owner, Project Admin, Project Member, Scheduled Maintenance Admin, Scheduled Maintenance Member, Create Scheduled Maintenance Template], Read: [Project Owner, Project Admin, Project Member, Viewer, Scheduled Maintenance Admin, Scheduled Maintenance Member, Scheduled Maintenance Viewer, Read Scheduled Maintenance Template], Update: [Project Owner, Project Admin, Project Member, Scheduled Maintenance Admin, Scheduled Maintenance Member, Edit Scheduled Maintenance Template]",
                Computed: true,
            },
            "title": schema.StringAttribute{
                MarkdownDescription: "Title of this scheduled event.. Permissions - Create: [Project Owner, Project Admin, Project Member, Scheduled Maintenance Admin, Scheduled Maintenance Member, Create Scheduled Maintenance Template], Read: [Project Owner, Project Admin, Project Member, Viewer, Scheduled Maintenance Admin, Scheduled Maintenance Member, Scheduled Maintenance Viewer, Read Scheduled Maintenance Template], Update: [Project Owner, Project Admin, Project Member, Scheduled Maintenance Admin, Scheduled Maintenance Member, Edit Scheduled Maintenance Template]",
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Description of this scheduled event that will show up on Status Page. This is a markdown field.. Permissions - Create: [Project Owner, Project Admin, Project Member, Scheduled Maintenance Admin, Scheduled Maintenance Member, Create Scheduled Maintenance Template], Read: [Project Owner, Project Admin, Project Member, Viewer, Scheduled Maintenance Admin, Scheduled Maintenance Member, Scheduled Maintenance Viewer, Read Scheduled Maintenance Template], Update: [Project Owner, Project Admin, Project Member, Scheduled Maintenance Admin, Scheduled Maintenance Member, Edit Scheduled Maintenance Template]",
                Computed: true,
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Friendly globally unique name for your object. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Scheduled Maintenance Admin, Scheduled Maintenance Member, Scheduled Maintenance Viewer, Read Scheduled Maintenance Template], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "monitors": schema.SetAttribute{
                MarkdownDescription: "List of monitors attached to this event. Permissions - Create: [Project Owner, Project Admin, Project Member, Scheduled Maintenance Admin, Scheduled Maintenance Member, Create Scheduled Maintenance Template], Read: [Project Owner, Project Admin, Project Member, Viewer, Scheduled Maintenance Admin, Scheduled Maintenance Member, Scheduled Maintenance Viewer, Read Scheduled Maintenance Template], Update: [Project Owner, Project Admin, Project Member, Scheduled Maintenance Admin, Scheduled Maintenance Member, Edit Scheduled Maintenance Template]",
                Computed: true,
                ElementType: types.StringType,
            },
            "hosts": schema.SetAttribute{
                MarkdownDescription: "List of hosts to pre-populate on scheduled maintenance events created from this template.. Permissions - Create: [Project Owner, Project Admin, Project Member, Scheduled Maintenance Admin, Scheduled Maintenance Member, Create Scheduled Maintenance Template], Read: [Project Owner, Project Admin, Project Member, Viewer, Scheduled Maintenance Admin, Scheduled Maintenance Member, Scheduled Maintenance Viewer, Read Scheduled Maintenance Template], Update: [Project Owner, Project Admin, Project Member, Scheduled Maintenance Admin, Scheduled Maintenance Member, Edit Scheduled Maintenance Template]",
                Computed: true,
                ElementType: types.StringType,
            },
            "kubernetes_clusters": schema.SetAttribute{
                MarkdownDescription: "List of Kubernetes clusters to pre-populate on scheduled maintenance events created from this template.. Permissions - Create: [Project Owner, Project Admin, Project Member, Scheduled Maintenance Admin, Scheduled Maintenance Member, Create Scheduled Maintenance Template], Read: [Project Owner, Project Admin, Project Member, Viewer, Scheduled Maintenance Admin, Scheduled Maintenance Member, Scheduled Maintenance Viewer, Read Scheduled Maintenance Template], Update: [Project Owner, Project Admin, Project Member, Scheduled Maintenance Admin, Scheduled Maintenance Member, Edit Scheduled Maintenance Template]",
                Computed: true,
                ElementType: types.StringType,
            },
            "docker_hosts": schema.SetAttribute{
                MarkdownDescription: "List of Docker hosts to pre-populate on scheduled maintenance events created from this template.. Permissions - Create: [Project Owner, Project Admin, Project Member, Scheduled Maintenance Admin, Scheduled Maintenance Member, Create Scheduled Maintenance Template], Read: [Project Owner, Project Admin, Project Member, Viewer, Scheduled Maintenance Admin, Scheduled Maintenance Member, Scheduled Maintenance Viewer, Read Scheduled Maintenance Template], Update: [Project Owner, Project Admin, Project Member, Scheduled Maintenance Admin, Scheduled Maintenance Member, Edit Scheduled Maintenance Template]",
                Computed: true,
                ElementType: types.StringType,
            },
            "services": schema.SetAttribute{
                MarkdownDescription: "List of services to pre-populate on scheduled maintenance events created from this template.. Permissions - Create: [Project Owner, Project Admin, Project Member, Scheduled Maintenance Admin, Scheduled Maintenance Member, Create Scheduled Maintenance Template], Read: [Project Owner, Project Admin, Project Member, Viewer, Scheduled Maintenance Admin, Scheduled Maintenance Member, Scheduled Maintenance Viewer, Read Scheduled Maintenance Template], Update: [Project Owner, Project Admin, Project Member, Scheduled Maintenance Admin, Scheduled Maintenance Member, Edit Scheduled Maintenance Template]",
                Computed: true,
                ElementType: types.StringType,
            },
            "status_pages": schema.SetAttribute{
                MarkdownDescription: "List of status pages to show this event on. Permissions - Create: [Project Owner, Project Admin, Project Member, Scheduled Maintenance Admin, Scheduled Maintenance Member, Create Scheduled Maintenance Template], Read: [Project Owner, Project Admin, Project Member, Viewer, Scheduled Maintenance Admin, Scheduled Maintenance Member, Scheduled Maintenance Viewer, Read Scheduled Maintenance Template], Update: [Project Owner, Project Admin, Project Member, Scheduled Maintenance Admin, Scheduled Maintenance Member, Edit Scheduled Maintenance Template]",
                Computed: true,
                ElementType: types.StringType,
            },
            "labels": schema.SetAttribute{
                MarkdownDescription: "Relation to Labels Array where this object is categorized in.. Permissions - Create: [Project Owner, Project Admin, Project Member, Scheduled Maintenance Admin, Scheduled Maintenance Member, Create Scheduled Maintenance Template], Read: [Project Owner, Project Admin, Project Member, Viewer, Scheduled Maintenance Admin, Scheduled Maintenance Member, Scheduled Maintenance Viewer, Read Scheduled Maintenance Template], Update: [Project Owner, Project Admin, Project Member, Scheduled Maintenance Admin, Scheduled Maintenance Member, Edit Scheduled Maintenance Template]",
                Computed: true,
                ElementType: types.StringType,
            },
            "change_monitor_status_to_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "first_event_scheduled_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "first_event_starts_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "first_event_ends_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "recurring_interval": schema.StringAttribute{
                MarkdownDescription: "How often should this event recur?. Permissions - Create: [Project Owner, Project Admin, Project Member, Scheduled Maintenance Admin, Scheduled Maintenance Member, Create Scheduled Maintenance Template], Read: [Project Owner, Project Admin, Project Member, Viewer, Scheduled Maintenance Admin, Scheduled Maintenance Member, Scheduled Maintenance Viewer, Read Scheduled Maintenance Template], Update: [Project Owner, Project Admin, Project Member, Scheduled Maintenance Admin, Scheduled Maintenance Member, Edit Scheduled Maintenance Template]",
                Computed: true,
            },
            "is_recurring_event": schema.BoolAttribute{
                MarkdownDescription: "Is this a recurring event?. Permissions - Create: [Project Owner, Project Admin, Project Member, Scheduled Maintenance Admin, Scheduled Maintenance Member, Create Scheduled Maintenance Template], Read: [Project Owner, Project Admin, Project Member, Viewer, Scheduled Maintenance Admin, Scheduled Maintenance Member, Scheduled Maintenance Viewer, Read Scheduled Maintenance Template], Update: [Project Owner, Project Admin, Project Member, Scheduled Maintenance Admin, Scheduled Maintenance Member, Edit Scheduled Maintenance Note Template]",
                Computed: true,
            },
            "schedule_next_event_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "should_status_page_subscribers_be_notified_on_event_created": schema.BoolAttribute{
                MarkdownDescription: "Should subscribers be notified about this event creation?. Permissions - Create: [Project Owner, Project Admin, Project Member, Scheduled Maintenance Admin, Scheduled Maintenance Member, Create Scheduled Maintenance Template], Read: [Project Owner, Project Admin, Project Member, Viewer, Scheduled Maintenance Admin, Scheduled Maintenance Member, Scheduled Maintenance Viewer, Read Scheduled Maintenance Template], Update: [Project Owner, Project Admin, Project Member, Scheduled Maintenance Admin, Scheduled Maintenance Member, Edit Scheduled Maintenance Note Template]",
                Computed: true,
            },
            "should_status_page_subscribers_be_notified_when_event_changed_to_ongoing": schema.BoolAttribute{
                MarkdownDescription: "Should subscribers be notified about this event event is changed to ongoing?. Permissions - Create: [Project Owner, Project Admin, Project Member, Scheduled Maintenance Admin, Scheduled Maintenance Member, Create Scheduled Maintenance Template], Read: [Project Owner, Project Admin, Project Member, Viewer, Scheduled Maintenance Admin, Scheduled Maintenance Member, Scheduled Maintenance Viewer, Read Scheduled Maintenance Template], Update: [Project Owner, Project Admin, Project Member, Scheduled Maintenance Admin, Scheduled Maintenance Member, Edit Scheduled Maintenance Note Template]",
                Computed: true,
            },
            "should_status_page_subscribers_be_notified_when_event_changed_to_ended": schema.BoolAttribute{
                MarkdownDescription: "Should subscribers be notified about this event event is changed to ended?. Permissions - Create: [Project Owner, Project Admin, Project Member, Scheduled Maintenance Admin, Scheduled Maintenance Member, Create Scheduled Maintenance Template], Read: [Project Owner, Project Admin, Project Member, Viewer, Scheduled Maintenance Admin, Scheduled Maintenance Member, Scheduled Maintenance Viewer, Read Scheduled Maintenance Template], Update: [Project Owner, Project Admin, Project Member, Scheduled Maintenance Admin, Scheduled Maintenance Member, Edit Scheduled Maintenance Note Template]",
                Computed: true,
            },
            "custom_fields": schema.StringAttribute{
                MarkdownDescription: "Custom Fields on this resource.. Permissions - Create: [Project Owner, Project Admin, Project Member, Scheduled Maintenance Admin, Scheduled Maintenance Member, Create Scheduled Maintenance Template], Read: [Project Owner, Project Admin, Project Member, Viewer, Scheduled Maintenance Admin, Scheduled Maintenance Member, Scheduled Maintenance Viewer, Read Scheduled Maintenance Template], Update: [Project Owner, Project Admin, Project Member, Scheduled Maintenance Admin, Scheduled Maintenance Member, Edit Scheduled Maintenance Template]",
                Computed: true,
            },
            "send_subscriber_notifications_on_before_the_event": schema.StringAttribute{
                MarkdownDescription: "Should subscribers be notified before the event?. Permissions - Create: [Project Owner, Project Admin, Project Member, Scheduled Maintenance Admin, Scheduled Maintenance Member, Create Scheduled Maintenance Template], Read: [Project Owner, Project Admin, Project Member, Viewer, Scheduled Maintenance Admin, Scheduled Maintenance Member, Scheduled Maintenance Viewer, Read Scheduled Maintenance Template], Update: [Project Owner, Project Admin, Project Member, Scheduled Maintenance Admin, Scheduled Maintenance Member, Edit Scheduled Maintenance Template]",
                Computed: true,
            },
        },
    }
}

func (d *ScheduledMaintenanceTemplateDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ScheduledMaintenanceTemplateDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data ScheduledMaintenanceTemplateDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "scheduled-maintenance-template" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read scheduled_maintenance_template_data, got error: %s", err))
        return
    }

    var scheduledMaintenanceTemplateDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &scheduledMaintenanceTemplateDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse scheduled_maintenance_template_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := scheduledMaintenanceTemplateDataResponse["data"].(map[string]interface{}); ok {
        scheduledMaintenanceTemplateDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := scheduledMaintenanceTemplateDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceTemplateDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceTemplateDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceTemplateDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceTemplateDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceTemplateDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := scheduledMaintenanceTemplateDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceTemplateDataResponse["template_name"].(string); ok {
        data.TemplateName = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceTemplateDataResponse["template_description"].(string); ok {
        data.TemplateDescription = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceTemplateDataResponse["title"].(string); ok {
        data.Title = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceTemplateDataResponse["description"].(string); ok {
        data.Description = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceTemplateDataResponse["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceTemplateDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceTemplateDataResponse["monitors"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.Monitors = setValue
    }
    if val, ok := scheduledMaintenanceTemplateDataResponse["hosts"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.Hosts = setValue
    }
    if val, ok := scheduledMaintenanceTemplateDataResponse["kubernetes_clusters"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.KubernetesClusters = setValue
    }
    if val, ok := scheduledMaintenanceTemplateDataResponse["docker_hosts"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.DockerHosts = setValue
    }
    if val, ok := scheduledMaintenanceTemplateDataResponse["services"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.Services = setValue
    }
    if val, ok := scheduledMaintenanceTemplateDataResponse["status_pages"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.StatusPages = setValue
    }
    if val, ok := scheduledMaintenanceTemplateDataResponse["labels"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.Labels = setValue
    }
    if val, ok := scheduledMaintenanceTemplateDataResponse["change_monitor_status_to_id"].(string); ok {
        data.ChangeMonitorStatusToId = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceTemplateDataResponse["first_event_scheduled_at"].(string); ok {
        data.FirstEventScheduledAt = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceTemplateDataResponse["first_event_starts_at"].(string); ok {
        data.FirstEventStartsAt = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceTemplateDataResponse["first_event_ends_at"].(string); ok {
        data.FirstEventEndsAt = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceTemplateDataResponse["recurring_interval"].(string); ok {
        data.RecurringInterval = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceTemplateDataResponse["is_recurring_event"].(bool); ok {
        data.IsRecurringEvent = types.BoolValue(val)
    }
    if val, ok := scheduledMaintenanceTemplateDataResponse["schedule_next_event_at"].(string); ok {
        data.ScheduleNextEventAt = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceTemplateDataResponse["should_status_page_subscribers_be_notified_on_event_created"].(bool); ok {
        data.ShouldStatusPageSubscribersBeNotifiedOnEventCreated = types.BoolValue(val)
    }
    if val, ok := scheduledMaintenanceTemplateDataResponse["should_status_page_subscribers_be_notified_when_event_changed_to_ongoing"].(bool); ok {
        data.ShouldStatusPageSubscribersBeNotifiedWhenEventChangedToOngoing = types.BoolValue(val)
    }
    if val, ok := scheduledMaintenanceTemplateDataResponse["should_status_page_subscribers_be_notified_when_event_changed_to_ended"].(bool); ok {
        data.ShouldStatusPageSubscribersBeNotifiedWhenEventChangedToEnded = types.BoolValue(val)
    }
    if val, ok := scheduledMaintenanceTemplateDataResponse["custom_fields"].(string); ok {
        data.CustomFields = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceTemplateDataResponse["send_subscriber_notifications_on_before_the_event"].(string); ok {
        data.SendSubscriberNotificationsOnBeforeTheEvent = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
