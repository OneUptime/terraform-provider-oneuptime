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
var _ datasource.DataSource = &ScheduledMaintenanceEventDataDataSource{}

func NewScheduledMaintenanceEventDataDataSource() datasource.DataSource {
    return &ScheduledMaintenanceEventDataDataSource{}
}

// ScheduledMaintenanceEventDataDataSource defines the data source implementation.
type ScheduledMaintenanceEventDataDataSource struct {
    client *Client
}

// ScheduledMaintenanceEventDataDataSourceModel describes the data source data model.
type ScheduledMaintenanceEventDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    Title types.String `tfsdk:"title"`
    Description types.String `tfsdk:"description"`
    Slug types.String `tfsdk:"slug"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    Monitors types.List `tfsdk:"monitors"`
    StatusPages types.List `tfsdk:"status_pages"`
    Labels types.List `tfsdk:"labels"`
    CurrentScheduledMaintenanceStateId types.String `tfsdk:"current_scheduled_maintenance_state_id"`
    ChangeMonitorStatusToId types.String `tfsdk:"change_monitor_status_to_id"`
    StartsAt types.String `tfsdk:"starts_at"`
    EndsAt types.String `tfsdk:"ends_at"`
    IsStatusPageSubscribersNotifiedOnEventScheduled types.Bool `tfsdk:"is_status_page_subscribers_notified_on_event_scheduled"`
    ShouldStatusPageSubscribersBeNotifiedOnEventCreated types.Bool `tfsdk:"should_status_page_subscribers_be_notified_on_event_created"`
    ShouldStatusPageSubscribersBeNotifiedWhenEventChangedToOngoing types.Bool `tfsdk:"should_status_page_subscribers_be_notified_when_event_changed_to_ongoing"`
    ShouldStatusPageSubscribersBeNotifiedWhenEventChangedToEnded types.Bool `tfsdk:"should_status_page_subscribers_be_notified_when_event_changed_to_ended"`
    CustomFields types.String `tfsdk:"custom_fields"`
    IsOwnerNotifiedOfResourceCreation types.Bool `tfsdk:"is_owner_notified_of_resource_creation"`
    SendSubscriberNotificationsOnBeforeTheEvent types.String `tfsdk:"send_subscriber_notifications_on_before_the_event"`
    NextSubscriberNotificationBeforeTheEventAt types.String `tfsdk:"next_subscriber_notification_before_the_event_at"`
    ScheduledMaintenanceNumber types.Number `tfsdk:"scheduled_maintenance_number"`
    IsVisibleOnStatusPage types.Bool `tfsdk:"is_visible_on_status_page"`
}

func (d *ScheduledMaintenanceEventDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_scheduled_maintenance_event_data"
}

func (d *ScheduledMaintenanceEventDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "scheduled_maintenance_event_data data source",

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
                MarkdownDescription: "Version",
                Computed: true,
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "title": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Scheduled Maintenance], Read: [Project Owner, Project Admin, Project Member, Read Scheduled Maintenance], Update: [Project Owner, Project Admin, Project Member, Edit Scheduled Maintenance]",
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Scheduled Maintenance], Read: [Project Owner, Project Admin, Project Member, Read Scheduled Maintenance], Update: [Project Owner, Project Admin, Project Member, Edit Scheduled Maintenance]",
                Computed: true,
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Scheduled Maintenance], Read: [Project Owner, Project Admin, Project Member, Read Scheduled Maintenance], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "monitors": schema.ListAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Scheduled Maintenance], Read: [Project Owner, Project Admin, Project Member, Read Scheduled Maintenance], Update: [Project Owner, Project Admin, Project Member, Edit Scheduled Maintenance]",
                Computed: true,
                ElementType: types.StringType,
            },
            "status_pages": schema.ListAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Scheduled Maintenance], Read: [Project Owner, Project Admin, Project Member, Read Scheduled Maintenance], Update: [Project Owner, Project Admin, Project Member, Edit Scheduled Maintenance]",
                Computed: true,
                ElementType: types.StringType,
            },
            "labels": schema.ListAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Scheduled Maintenance], Read: [Project Owner, Project Admin, Project Member, Read Scheduled Maintenance], Update: [Project Owner, Project Admin, Project Member, Edit Scheduled Maintenance]",
                Computed: true,
                ElementType: types.StringType,
            },
            "current_scheduled_maintenance_state_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "change_monitor_status_to_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "starts_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "ends_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "is_status_page_subscribers_notified_on_event_scheduled": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Scheduled Maintenance], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "should_status_page_subscribers_be_notified_on_event_created": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Scheduled Maintenance], Read: [Project Owner, Project Admin, Project Member, Read Scheduled Maintenance], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "should_status_page_subscribers_be_notified_when_event_changed_to_ongoing": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Scheduled Maintenance], Read: [Project Owner, Project Admin, Project Member, Read Scheduled Maintenance], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "should_status_page_subscribers_be_notified_when_event_changed_to_ended": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Scheduled Maintenance], Read: [Project Owner, Project Admin, Project Member, Read Scheduled Maintenance], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "custom_fields": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Scheduled Maintenance], Read: [Project Owner, Project Admin, Project Member, Read Scheduled Maintenance], Update: [Project Owner, Project Admin, Project Member, Edit Scheduled Maintenance]",
                Computed: true,
            },
            "is_owner_notified_of_resource_creation": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Scheduled Maintenance], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "send_subscriber_notifications_on_before_the_event": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Scheduled Maintenance], Read: [Project Owner, Project Admin, Project Member, Read Scheduled Maintenance], Update: [Project Owner, Project Admin, Project Member, Edit Scheduled Maintenance]",
                Computed: true,
            },
            "next_subscriber_notification_before_the_event_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "scheduled_maintenance_number": schema.NumberAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Scheduled Maintenance], Read: [Project Owner, Project Admin, Project Member, Read Scheduled Maintenance], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "is_visible_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Scheduled Maintenance], Read: [Project Owner, Project Admin, Project Member, Read Scheduled Maintenance], Update: [Project Owner, Project Admin, Project Member, Edit Scheduled Maintenance]",
                Computed: true,
            },
        },
    }
}

func (d *ScheduledMaintenanceEventDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ScheduledMaintenanceEventDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data ScheduledMaintenanceEventDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "scheduled-maintenance" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read scheduled_maintenance_event_data, got error: %s", err))
        return
    }

    var scheduledMaintenanceEventDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &scheduledMaintenanceEventDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse scheduled_maintenance_event_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := scheduledMaintenanceEventDataResponse["data"].(map[string]interface{}); ok {
        scheduledMaintenanceEventDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := scheduledMaintenanceEventDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceEventDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceEventDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceEventDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceEventDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceEventDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := scheduledMaintenanceEventDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceEventDataResponse["title"].(string); ok {
        data.Title = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceEventDataResponse["description"].(string); ok {
        data.Description = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceEventDataResponse["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceEventDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceEventDataResponse["monitors"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        listValue, _ := types.ListValue(types.StringType, elements)
        data.Monitors = listValue
    }
    if val, ok := scheduledMaintenanceEventDataResponse["status_pages"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        listValue, _ := types.ListValue(types.StringType, elements)
        data.StatusPages = listValue
    }
    if val, ok := scheduledMaintenanceEventDataResponse["labels"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        listValue, _ := types.ListValue(types.StringType, elements)
        data.Labels = listValue
    }
    if val, ok := scheduledMaintenanceEventDataResponse["current_scheduled_maintenance_state_id"].(string); ok {
        data.CurrentScheduledMaintenanceStateId = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceEventDataResponse["change_monitor_status_to_id"].(string); ok {
        data.ChangeMonitorStatusToId = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceEventDataResponse["starts_at"].(string); ok {
        data.StartsAt = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceEventDataResponse["ends_at"].(string); ok {
        data.EndsAt = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceEventDataResponse["is_status_page_subscribers_notified_on_event_scheduled"].(bool); ok {
        data.IsStatusPageSubscribersNotifiedOnEventScheduled = types.BoolValue(val)
    }
    if val, ok := scheduledMaintenanceEventDataResponse["should_status_page_subscribers_be_notified_on_event_created"].(bool); ok {
        data.ShouldStatusPageSubscribersBeNotifiedOnEventCreated = types.BoolValue(val)
    }
    if val, ok := scheduledMaintenanceEventDataResponse["should_status_page_subscribers_be_notified_when_event_changed_to_ongoing"].(bool); ok {
        data.ShouldStatusPageSubscribersBeNotifiedWhenEventChangedToOngoing = types.BoolValue(val)
    }
    if val, ok := scheduledMaintenanceEventDataResponse["should_status_page_subscribers_be_notified_when_event_changed_to_ended"].(bool); ok {
        data.ShouldStatusPageSubscribersBeNotifiedWhenEventChangedToEnded = types.BoolValue(val)
    }
    if val, ok := scheduledMaintenanceEventDataResponse["custom_fields"].(string); ok {
        data.CustomFields = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceEventDataResponse["is_owner_notified_of_resource_creation"].(bool); ok {
        data.IsOwnerNotifiedOfResourceCreation = types.BoolValue(val)
    }
    if val, ok := scheduledMaintenanceEventDataResponse["send_subscriber_notifications_on_before_the_event"].(string); ok {
        data.SendSubscriberNotificationsOnBeforeTheEvent = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceEventDataResponse["next_subscriber_notification_before_the_event_at"].(string); ok {
        data.NextSubscriberNotificationBeforeTheEventAt = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceEventDataResponse["scheduled_maintenance_number"].(float64); ok {
        data.ScheduledMaintenanceNumber = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := scheduledMaintenanceEventDataResponse["is_visible_on_status_page"].(bool); ok {
        data.IsVisibleOnStatusPage = types.BoolValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
