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
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ScheduledMaintenanceEventResource{}
var _ resource.ResourceWithImportState = &ScheduledMaintenanceEventResource{}

func NewScheduledMaintenanceEventResource() resource.Resource {
    return &ScheduledMaintenanceEventResource{}
}

// ScheduledMaintenanceEventResource defines the resource implementation.
type ScheduledMaintenanceEventResource struct {
    client *Client
}

// ScheduledMaintenanceEventResourceModel describes the resource data model.
type ScheduledMaintenanceEventResourceModel struct {
    Id types.String `tfsdk:"id"`
    ProjectId types.String `tfsdk:"project_id"`
    Title types.String `tfsdk:"title"`
    Description types.String `tfsdk:"description"`
    Monitors types.List `tfsdk:"monitors"`
    StatusPages types.List `tfsdk:"status_pages"`
    Labels types.List `tfsdk:"labels"`
    CurrentScheduledMaintenanceStateId types.String `tfsdk:"current_scheduled_maintenance_state_id"`
    ChangeMonitorStatusToId types.String `tfsdk:"change_monitor_status_to_id"`
    StartsAt types.Map `tfsdk:"starts_at"`
    EndsAt types.Map `tfsdk:"ends_at"`
    ShouldStatusPageSubscribersBeNotifiedOnEventCreated types.Bool `tfsdk:"should_status_page_subscribers_be_notified_on_event_created"`
    ShouldStatusPageSubscribersBeNotifiedWhenEventChangedToOngoing types.Bool `tfsdk:"should_status_page_subscribers_be_notified_when_event_changed_to_ongoing"`
    ShouldStatusPageSubscribersBeNotifiedWhenEventChangedToEnded types.Bool `tfsdk:"should_status_page_subscribers_be_notified_when_event_changed_to_ended"`
    CustomFields types.Map `tfsdk:"custom_fields"`
    IsOwnerNotifiedOfResourceCreation types.Bool `tfsdk:"is_owner_notified_of_resource_creation"`
    SendSubscriberNotificationsOnBeforeTheEvent types.Map `tfsdk:"send_subscriber_notifications_on_before_the_event"`
    NextSubscriberNotificationBeforeTheEventAt types.Map `tfsdk:"next_subscriber_notification_before_the_event_at"`
    ScheduledMaintenanceNumber types.Number `tfsdk:"scheduled_maintenance_number"`
    IsVisibleOnStatusPage types.Bool `tfsdk:"is_visible_on_status_page"`
    CreatedAt types.Map `tfsdk:"created_at"`
    UpdatedAt types.Map `tfsdk:"updated_at"`
    DeletedAt types.Map `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    Slug types.String `tfsdk:"slug"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    IsStatusPageSubscribersNotifiedOnEventScheduled types.Bool `tfsdk:"is_status_page_subscribers_notified_on_event_scheduled"`
}

func (r *ScheduledMaintenanceEventResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_scheduled_maintenance_event"
}

func (r *ScheduledMaintenanceEventResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "scheduled_maintenance_event resource",

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
                Optional: true,
            },
            "title": schema.StringAttribute{
                MarkdownDescription: "Title",
                Required: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Description",
                Optional: true,
            },
            "monitors": schema.ListAttribute{
                MarkdownDescription: "Monitors",
                Optional: true,
                ElementType: types.StringType,
            },
            "status_pages": schema.ListAttribute{
                MarkdownDescription: "Status Pages",
                Optional: true,
                ElementType: types.StringType,
            },
            "labels": schema.ListAttribute{
                MarkdownDescription: "Labels",
                Optional: true,
                ElementType: types.StringType,
            },
            "current_scheduled_maintenance_state_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Required: true,
            },
            "change_monitor_status_to_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Optional: true,
            },
            "starts_at": schema.MapAttribute{
                MarkdownDescription: "A date time object.",
                Required: true,
                ElementType: types.StringType,
            },
            "ends_at": schema.MapAttribute{
                MarkdownDescription: "A date time object.",
                Required: true,
                ElementType: types.StringType,
            },
            "should_status_page_subscribers_be_notified_on_event_created": schema.BoolAttribute{
                MarkdownDescription: "Should subscribers be notified when event is created?",
                Optional: true,
            },
            "should_status_page_subscribers_be_notified_when_event_changed_to_ongoing": schema.BoolAttribute{
                MarkdownDescription: "Should subscribers be notified when event is changed to ongoing?",
                Optional: true,
            },
            "should_status_page_subscribers_be_notified_when_event_changed_to_ended": schema.BoolAttribute{
                MarkdownDescription: "Should subscribers be notified when event is changed to ended?",
                Optional: true,
            },
            "custom_fields": schema.MapAttribute{
                MarkdownDescription: "Custom Fields",
                Optional: true,
                ElementType: types.StringType,
            },
            "is_owner_notified_of_resource_creation": schema.BoolAttribute{
                MarkdownDescription: "Are Owners Notified Of Resource Creation?",
                Optional: true,
            },
            "send_subscriber_notifications_on_before_the_event": schema.MapAttribute{
                MarkdownDescription: "Subscriber notifications before the event",
                Optional: true,
                ElementType: types.StringType,
            },
            "next_subscriber_notification_before_the_event_at": schema.MapAttribute{
                MarkdownDescription: "A date time object.",
                Optional: true,
                ElementType: types.StringType,
            },
            "scheduled_maintenance_number": schema.NumberAttribute{
                MarkdownDescription: "Scheduled Maintenance Number",
                Optional: true,
            },
            "is_visible_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Should be visible on status page?",
                Optional: true,
            },
            "created_at": schema.MapAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
                ElementType: types.StringType,
            },
            "updated_at": schema.MapAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
                ElementType: types.StringType,
            },
            "deleted_at": schema.MapAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
                ElementType: types.StringType,
            },
            "version": schema.NumberAttribute{
                MarkdownDescription: "Version",
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
            "is_status_page_subscribers_notified_on_event_scheduled": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Scheduled Maintenance], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
        },
    }
}

func (r *ScheduledMaintenanceEventResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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


func (r *ScheduledMaintenanceEventResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var data ScheduledMaintenanceEventResourceModel

    // Read Terraform plan data into the model
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create API request body
    scheduledMaintenanceEventRequest := map[string]interface{}{
        "data": map[string]interface{}{
        "projectId": data.ProjectId.ValueString(),
        "title": data.Title.ValueString(),
        "description": data.Description.ValueString(),
        "monitors": r.convertTerraformListToInterface(data.Monitors),
        "statusPages": r.convertTerraformListToInterface(data.StatusPages),
        "labels": r.convertTerraformListToInterface(data.Labels),
        "currentScheduledMaintenanceStateId": data.CurrentScheduledMaintenanceStateId.ValueString(),
        "changeMonitorStatusToId": data.ChangeMonitorStatusToId.ValueString(),
        "startsAt": r.convertTerraformMapToInterface(data.StartsAt),
        "endsAt": r.convertTerraformMapToInterface(data.EndsAt),
        "shouldStatusPageSubscribersBeNotifiedOnEventCreated": data.ShouldStatusPageSubscribersBeNotifiedOnEventCreated.ValueBool(),
        "shouldStatusPageSubscribersBeNotifiedWhenEventChangedToOngoing": data.ShouldStatusPageSubscribersBeNotifiedWhenEventChangedToOngoing.ValueBool(),
        "shouldStatusPageSubscribersBeNotifiedWhenEventChangedToEnded": data.ShouldStatusPageSubscribersBeNotifiedWhenEventChangedToEnded.ValueBool(),
        "customFields": r.convertTerraformMapToInterface(data.CustomFields),
        "isOwnerNotifiedOfResourceCreation": data.IsOwnerNotifiedOfResourceCreation.ValueBool(),
        "sendSubscriberNotificationsOnBeforeTheEvent": r.convertTerraformMapToInterface(data.SendSubscriberNotificationsOnBeforeTheEvent),
        "nextSubscriberNotificationBeforeTheEventAt": r.convertTerraformMapToInterface(data.NextSubscriberNotificationBeforeTheEventAt),
        "scheduledMaintenanceNumber": data.ScheduledMaintenanceNumber.ValueBigFloat(),
        "isVisibleOnStatusPage": data.IsVisibleOnStatusPage.ValueBool(),
        },
    }

    // Make API call
    httpResp, err := r.client.Post("/scheduled-maintenance", scheduledMaintenanceEventRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create scheduled_maintenance_event, got error: %s", err))
        return
    }

    var scheduledMaintenanceEventResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &scheduledMaintenanceEventResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse scheduled_maintenance_event response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := scheduledMaintenanceEventResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = scheduledMaintenanceEventResponse
    }

    if val, ok := dataMap["id"].(string); ok && val != "" {
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
    if val, ok := dataMap["title"].(string); ok && val != "" {
        data.Title = types.StringValue(val)
    } else {
        data.Title = types.StringNull()
    }
    if val, ok := dataMap["description"].(string); ok && val != "" {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if val, ok := dataMap["monitors"].([]interface{}); ok {
        // Convert API response list to Terraform list
        listValue, _ := types.ListValueFrom(ctx, types.StringType, val)
        data.Monitors = listValue
    } else if dataMap["monitors"] == nil {
        data.Monitors = types.ListNull(types.StringType)
    }
    if val, ok := dataMap["statusPages"].([]interface{}); ok {
        // Convert API response list to Terraform list
        listValue, _ := types.ListValueFrom(ctx, types.StringType, val)
        data.StatusPages = listValue
    } else if dataMap["statusPages"] == nil {
        data.StatusPages = types.ListNull(types.StringType)
    }
    if val, ok := dataMap["labels"].([]interface{}); ok {
        // Convert API response list to Terraform list
        listValue, _ := types.ListValueFrom(ctx, types.StringType, val)
        data.Labels = listValue
    } else if dataMap["labels"] == nil {
        data.Labels = types.ListNull(types.StringType)
    }
    if val, ok := dataMap["currentScheduledMaintenanceStateId"].(string); ok && val != "" {
        data.CurrentScheduledMaintenanceStateId = types.StringValue(val)
    } else {
        data.CurrentScheduledMaintenanceStateId = types.StringNull()
    }
    if val, ok := dataMap["changeMonitorStatusToId"].(string); ok && val != "" {
        data.ChangeMonitorStatusToId = types.StringValue(val)
    } else {
        data.ChangeMonitorStatusToId = types.StringNull()
    }
    if val, ok := dataMap["startsAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.StartsAt = mapValue
    } else if dataMap["startsAt"] == nil {
        data.StartsAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["endsAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.EndsAt = mapValue
    } else if dataMap["endsAt"] == nil {
        data.EndsAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["shouldStatusPageSubscribersBeNotifiedOnEventCreated"].(bool); ok {
        data.ShouldStatusPageSubscribersBeNotifiedOnEventCreated = types.BoolValue(val)
    } else if dataMap["shouldStatusPageSubscribersBeNotifiedOnEventCreated"] == nil {
        data.ShouldStatusPageSubscribersBeNotifiedOnEventCreated = types.BoolNull()
    }
    if val, ok := dataMap["shouldStatusPageSubscribersBeNotifiedWhenEventChangedToOngoing"].(bool); ok {
        data.ShouldStatusPageSubscribersBeNotifiedWhenEventChangedToOngoing = types.BoolValue(val)
    } else if dataMap["shouldStatusPageSubscribersBeNotifiedWhenEventChangedToOngoing"] == nil {
        data.ShouldStatusPageSubscribersBeNotifiedWhenEventChangedToOngoing = types.BoolNull()
    }
    if val, ok := dataMap["shouldStatusPageSubscribersBeNotifiedWhenEventChangedToEnded"].(bool); ok {
        data.ShouldStatusPageSubscribersBeNotifiedWhenEventChangedToEnded = types.BoolValue(val)
    } else if dataMap["shouldStatusPageSubscribersBeNotifiedWhenEventChangedToEnded"] == nil {
        data.ShouldStatusPageSubscribersBeNotifiedWhenEventChangedToEnded = types.BoolNull()
    }
    if val, ok := dataMap["customFields"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.CustomFields = mapValue
    } else if dataMap["customFields"] == nil {
        data.CustomFields = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["isOwnerNotifiedOfResourceCreation"].(bool); ok {
        data.IsOwnerNotifiedOfResourceCreation = types.BoolValue(val)
    } else if dataMap["isOwnerNotifiedOfResourceCreation"] == nil {
        data.IsOwnerNotifiedOfResourceCreation = types.BoolNull()
    }
    if val, ok := dataMap["sendSubscriberNotificationsOnBeforeTheEvent"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.SendSubscriberNotificationsOnBeforeTheEvent = mapValue
    } else if dataMap["sendSubscriberNotificationsOnBeforeTheEvent"] == nil {
        data.SendSubscriberNotificationsOnBeforeTheEvent = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["nextSubscriberNotificationBeforeTheEventAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.NextSubscriberNotificationBeforeTheEventAt = mapValue
    } else if dataMap["nextSubscriberNotificationBeforeTheEventAt"] == nil {
        data.NextSubscriberNotificationBeforeTheEventAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["scheduledMaintenanceNumber"].(float64); ok {
        data.ScheduledMaintenanceNumber = types.NumberValue(big.NewFloat(val))
    } else if dataMap["scheduledMaintenanceNumber"] == nil {
        data.ScheduledMaintenanceNumber = types.NumberNull()
    }
    if val, ok := dataMap["isVisibleOnStatusPage"].(bool); ok {
        data.IsVisibleOnStatusPage = types.BoolValue(val)
    } else if dataMap["isVisibleOnStatusPage"] == nil {
        data.IsVisibleOnStatusPage = types.BoolNull()
    }
    if val, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.CreatedAt = mapValue
    } else if dataMap["createdAt"] == nil {
        data.CreatedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.UpdatedAt = mapValue
    } else if dataMap["updatedAt"] == nil {
        data.UpdatedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["deletedAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.DeletedAt = mapValue
    } else if dataMap["deletedAt"] == nil {
        data.DeletedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if dataMap["version"] == nil {
        data.Version = types.NumberNull()
    }
    if val, ok := dataMap["slug"].(string); ok && val != "" {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
    }
    if val, ok := dataMap["createdByUserId"].(string); ok && val != "" {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if val, ok := dataMap["isStatusPageSubscribersNotifiedOnEventScheduled"].(bool); ok {
        data.IsStatusPageSubscribersNotifiedOnEventScheduled = types.BoolValue(val)
    } else if dataMap["isStatusPageSubscribersNotifiedOnEventScheduled"] == nil {
        data.IsStatusPageSubscribersNotifiedOnEventScheduled = types.BoolNull()
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

func (r *ScheduledMaintenanceEventResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var data ScheduledMaintenanceEventResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create select parameter to get full object
    selectParam := map[string]interface{}{
        "projectId": true,
        "title": true,
        "description": true,
        "monitors": true,
        "statusPages": true,
        "labels": true,
        "currentScheduledMaintenanceStateId": true,
        "changeMonitorStatusToId": true,
        "startsAt": true,
        "endsAt": true,
        "shouldStatusPageSubscribersBeNotifiedOnEventCreated": true,
        "shouldStatusPageSubscribersBeNotifiedWhenEventChangedToOngoing": true,
        "shouldStatusPageSubscribersBeNotifiedWhenEventChangedToEnded": true,
        "customFields": true,
        "isOwnerNotifiedOfResourceCreation": true,
        "sendSubscriberNotificationsOnBeforeTheEvent": true,
        "nextSubscriberNotificationBeforeTheEventAt": true,
        "scheduledMaintenanceNumber": true,
        "isVisibleOnStatusPage": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "slug": true,
        "createdByUserId": true,
        "isStatusPageSubscribersNotifiedOnEventScheduled": true,
        "_id": true,
    }

    // Make API call with select parameter
    httpResp, err := r.client.PostWithSelect("/scheduled-maintenance/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read scheduled_maintenance_event, got error: %s", err))
        return
    }

    if httpResp.StatusCode == http.StatusNotFound {
        resp.State.RemoveResource(ctx)
        return
    }

    var scheduledMaintenanceEventResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &scheduledMaintenanceEventResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse scheduled_maintenance_event response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := scheduledMaintenanceEventResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = scheduledMaintenanceEventResponse
    }

    if val, ok := dataMap["id"].(string); ok && val != "" {
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
    if val, ok := dataMap["title"].(string); ok && val != "" {
        data.Title = types.StringValue(val)
    } else {
        data.Title = types.StringNull()
    }
    if val, ok := dataMap["description"].(string); ok && val != "" {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if val, ok := dataMap["monitors"].([]interface{}); ok {
        // Convert API response list to Terraform list
        listValue, _ := types.ListValueFrom(ctx, types.StringType, val)
        data.Monitors = listValue
    } else if dataMap["monitors"] == nil {
        data.Monitors = types.ListNull(types.StringType)
    }
    if val, ok := dataMap["statusPages"].([]interface{}); ok {
        // Convert API response list to Terraform list
        listValue, _ := types.ListValueFrom(ctx, types.StringType, val)
        data.StatusPages = listValue
    } else if dataMap["statusPages"] == nil {
        data.StatusPages = types.ListNull(types.StringType)
    }
    if val, ok := dataMap["labels"].([]interface{}); ok {
        // Convert API response list to Terraform list
        listValue, _ := types.ListValueFrom(ctx, types.StringType, val)
        data.Labels = listValue
    } else if dataMap["labels"] == nil {
        data.Labels = types.ListNull(types.StringType)
    }
    if val, ok := dataMap["currentScheduledMaintenanceStateId"].(string); ok && val != "" {
        data.CurrentScheduledMaintenanceStateId = types.StringValue(val)
    } else {
        data.CurrentScheduledMaintenanceStateId = types.StringNull()
    }
    if val, ok := dataMap["changeMonitorStatusToId"].(string); ok && val != "" {
        data.ChangeMonitorStatusToId = types.StringValue(val)
    } else {
        data.ChangeMonitorStatusToId = types.StringNull()
    }
    if val, ok := dataMap["startsAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.StartsAt = mapValue
    } else if dataMap["startsAt"] == nil {
        data.StartsAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["endsAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.EndsAt = mapValue
    } else if dataMap["endsAt"] == nil {
        data.EndsAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["shouldStatusPageSubscribersBeNotifiedOnEventCreated"].(bool); ok {
        data.ShouldStatusPageSubscribersBeNotifiedOnEventCreated = types.BoolValue(val)
    } else if dataMap["shouldStatusPageSubscribersBeNotifiedOnEventCreated"] == nil {
        data.ShouldStatusPageSubscribersBeNotifiedOnEventCreated = types.BoolNull()
    }
    if val, ok := dataMap["shouldStatusPageSubscribersBeNotifiedWhenEventChangedToOngoing"].(bool); ok {
        data.ShouldStatusPageSubscribersBeNotifiedWhenEventChangedToOngoing = types.BoolValue(val)
    } else if dataMap["shouldStatusPageSubscribersBeNotifiedWhenEventChangedToOngoing"] == nil {
        data.ShouldStatusPageSubscribersBeNotifiedWhenEventChangedToOngoing = types.BoolNull()
    }
    if val, ok := dataMap["shouldStatusPageSubscribersBeNotifiedWhenEventChangedToEnded"].(bool); ok {
        data.ShouldStatusPageSubscribersBeNotifiedWhenEventChangedToEnded = types.BoolValue(val)
    } else if dataMap["shouldStatusPageSubscribersBeNotifiedWhenEventChangedToEnded"] == nil {
        data.ShouldStatusPageSubscribersBeNotifiedWhenEventChangedToEnded = types.BoolNull()
    }
    if val, ok := dataMap["customFields"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.CustomFields = mapValue
    } else if dataMap["customFields"] == nil {
        data.CustomFields = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["isOwnerNotifiedOfResourceCreation"].(bool); ok {
        data.IsOwnerNotifiedOfResourceCreation = types.BoolValue(val)
    } else if dataMap["isOwnerNotifiedOfResourceCreation"] == nil {
        data.IsOwnerNotifiedOfResourceCreation = types.BoolNull()
    }
    if val, ok := dataMap["sendSubscriberNotificationsOnBeforeTheEvent"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.SendSubscriberNotificationsOnBeforeTheEvent = mapValue
    } else if dataMap["sendSubscriberNotificationsOnBeforeTheEvent"] == nil {
        data.SendSubscriberNotificationsOnBeforeTheEvent = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["nextSubscriberNotificationBeforeTheEventAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.NextSubscriberNotificationBeforeTheEventAt = mapValue
    } else if dataMap["nextSubscriberNotificationBeforeTheEventAt"] == nil {
        data.NextSubscriberNotificationBeforeTheEventAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["scheduledMaintenanceNumber"].(float64); ok {
        data.ScheduledMaintenanceNumber = types.NumberValue(big.NewFloat(val))
    } else if dataMap["scheduledMaintenanceNumber"] == nil {
        data.ScheduledMaintenanceNumber = types.NumberNull()
    }
    if val, ok := dataMap["isVisibleOnStatusPage"].(bool); ok {
        data.IsVisibleOnStatusPage = types.BoolValue(val)
    } else if dataMap["isVisibleOnStatusPage"] == nil {
        data.IsVisibleOnStatusPage = types.BoolNull()
    }
    if val, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.CreatedAt = mapValue
    } else if dataMap["createdAt"] == nil {
        data.CreatedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.UpdatedAt = mapValue
    } else if dataMap["updatedAt"] == nil {
        data.UpdatedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["deletedAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.DeletedAt = mapValue
    } else if dataMap["deletedAt"] == nil {
        data.DeletedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if dataMap["version"] == nil {
        data.Version = types.NumberNull()
    }
    if val, ok := dataMap["slug"].(string); ok && val != "" {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
    }
    if val, ok := dataMap["createdByUserId"].(string); ok && val != "" {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if val, ok := dataMap["isStatusPageSubscribersNotifiedOnEventScheduled"].(bool); ok {
        data.IsStatusPageSubscribersNotifiedOnEventScheduled = types.BoolValue(val)
    } else if dataMap["isStatusPageSubscribersNotifiedOnEventScheduled"] == nil {
        data.IsStatusPageSubscribersNotifiedOnEventScheduled = types.BoolNull()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ScheduledMaintenanceEventResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var data ScheduledMaintenanceEventResourceModel
    var state ScheduledMaintenanceEventResourceModel

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
    scheduledMaintenanceEventRequest := map[string]interface{}{
        "data": map[string]interface{}{
        "title": data.Title.ValueString(),
        "description": data.Description.ValueString(),
        "monitors": r.convertTerraformListToInterface(data.Monitors),
        "statusPages": r.convertTerraformListToInterface(data.StatusPages),
        "labels": r.convertTerraformListToInterface(data.Labels),
        "currentScheduledMaintenanceStateId": data.CurrentScheduledMaintenanceStateId.ValueString(),
        "changeMonitorStatusToId": data.ChangeMonitorStatusToId.ValueString(),
        "startsAt": r.convertTerraformMapToInterface(data.StartsAt),
        "endsAt": r.convertTerraformMapToInterface(data.EndsAt),
        "shouldStatusPageSubscribersBeNotifiedOnEventCreated": data.ShouldStatusPageSubscribersBeNotifiedOnEventCreated.ValueBool(),
        "shouldStatusPageSubscribersBeNotifiedWhenEventChangedToOngoing": data.ShouldStatusPageSubscribersBeNotifiedWhenEventChangedToOngoing.ValueBool(),
        "shouldStatusPageSubscribersBeNotifiedWhenEventChangedToEnded": data.ShouldStatusPageSubscribersBeNotifiedWhenEventChangedToEnded.ValueBool(),
        "customFields": r.convertTerraformMapToInterface(data.CustomFields),
        "isOwnerNotifiedOfResourceCreation": data.IsOwnerNotifiedOfResourceCreation.ValueBool(),
        "sendSubscriberNotificationsOnBeforeTheEvent": r.convertTerraformMapToInterface(data.SendSubscriberNotificationsOnBeforeTheEvent),
        "nextSubscriberNotificationBeforeTheEventAt": r.convertTerraformMapToInterface(data.NextSubscriberNotificationBeforeTheEventAt),
        "scheduledMaintenanceNumber": data.ScheduledMaintenanceNumber.ValueBigFloat(),
        "isVisibleOnStatusPage": data.IsVisibleOnStatusPage.ValueBool(),
        },
    }

    // Make API call
    httpResp, err := r.client.Put("/scheduled-maintenance/" + data.Id.ValueString() + "", scheduledMaintenanceEventRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update scheduled_maintenance_event, got error: %s", err))
        return
    }

    // Parse the update response
    var scheduledMaintenanceEventResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &scheduledMaintenanceEventResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse scheduled_maintenance_event response, got error: %s", err))
        return
    }

    // After successful update, fetch the current state by calling Read with select parameter
    selectParam := map[string]interface{}{
        "projectId": true,
        "title": true,
        "description": true,
        "monitors": true,
        "statusPages": true,
        "labels": true,
        "currentScheduledMaintenanceStateId": true,
        "changeMonitorStatusToId": true,
        "startsAt": true,
        "endsAt": true,
        "shouldStatusPageSubscribersBeNotifiedOnEventCreated": true,
        "shouldStatusPageSubscribersBeNotifiedWhenEventChangedToOngoing": true,
        "shouldStatusPageSubscribersBeNotifiedWhenEventChangedToEnded": true,
        "customFields": true,
        "isOwnerNotifiedOfResourceCreation": true,
        "sendSubscriberNotificationsOnBeforeTheEvent": true,
        "nextSubscriberNotificationBeforeTheEventAt": true,
        "scheduledMaintenanceNumber": true,
        "isVisibleOnStatusPage": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "slug": true,
        "createdByUserId": true,
        "isStatusPageSubscribersNotifiedOnEventScheduled": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect("/scheduled-maintenance/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read scheduled_maintenance_event after update, got error: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse scheduled_maintenance_event read response, got error: %s", err))
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

    if val, ok := dataMap["id"].(string); ok && val != "" {
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
    if val, ok := dataMap["title"].(string); ok && val != "" {
        data.Title = types.StringValue(val)
    } else {
        data.Title = types.StringNull()
    }
    if val, ok := dataMap["description"].(string); ok && val != "" {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if val, ok := dataMap["monitors"].([]interface{}); ok {
        // Convert API response list to Terraform list
        listValue, _ := types.ListValueFrom(ctx, types.StringType, val)
        data.Monitors = listValue
    } else if dataMap["monitors"] == nil {
        data.Monitors = types.ListNull(types.StringType)
    }
    if val, ok := dataMap["statusPages"].([]interface{}); ok {
        // Convert API response list to Terraform list
        listValue, _ := types.ListValueFrom(ctx, types.StringType, val)
        data.StatusPages = listValue
    } else if dataMap["statusPages"] == nil {
        data.StatusPages = types.ListNull(types.StringType)
    }
    if val, ok := dataMap["labels"].([]interface{}); ok {
        // Convert API response list to Terraform list
        listValue, _ := types.ListValueFrom(ctx, types.StringType, val)
        data.Labels = listValue
    } else if dataMap["labels"] == nil {
        data.Labels = types.ListNull(types.StringType)
    }
    if val, ok := dataMap["currentScheduledMaintenanceStateId"].(string); ok && val != "" {
        data.CurrentScheduledMaintenanceStateId = types.StringValue(val)
    } else {
        data.CurrentScheduledMaintenanceStateId = types.StringNull()
    }
    if val, ok := dataMap["changeMonitorStatusToId"].(string); ok && val != "" {
        data.ChangeMonitorStatusToId = types.StringValue(val)
    } else {
        data.ChangeMonitorStatusToId = types.StringNull()
    }
    if val, ok := dataMap["startsAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.StartsAt = mapValue
    } else if dataMap["startsAt"] == nil {
        data.StartsAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["endsAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.EndsAt = mapValue
    } else if dataMap["endsAt"] == nil {
        data.EndsAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["shouldStatusPageSubscribersBeNotifiedOnEventCreated"].(bool); ok {
        data.ShouldStatusPageSubscribersBeNotifiedOnEventCreated = types.BoolValue(val)
    } else if dataMap["shouldStatusPageSubscribersBeNotifiedOnEventCreated"] == nil {
        data.ShouldStatusPageSubscribersBeNotifiedOnEventCreated = types.BoolNull()
    }
    if val, ok := dataMap["shouldStatusPageSubscribersBeNotifiedWhenEventChangedToOngoing"].(bool); ok {
        data.ShouldStatusPageSubscribersBeNotifiedWhenEventChangedToOngoing = types.BoolValue(val)
    } else if dataMap["shouldStatusPageSubscribersBeNotifiedWhenEventChangedToOngoing"] == nil {
        data.ShouldStatusPageSubscribersBeNotifiedWhenEventChangedToOngoing = types.BoolNull()
    }
    if val, ok := dataMap["shouldStatusPageSubscribersBeNotifiedWhenEventChangedToEnded"].(bool); ok {
        data.ShouldStatusPageSubscribersBeNotifiedWhenEventChangedToEnded = types.BoolValue(val)
    } else if dataMap["shouldStatusPageSubscribersBeNotifiedWhenEventChangedToEnded"] == nil {
        data.ShouldStatusPageSubscribersBeNotifiedWhenEventChangedToEnded = types.BoolNull()
    }
    if val, ok := dataMap["customFields"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.CustomFields = mapValue
    } else if dataMap["customFields"] == nil {
        data.CustomFields = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["isOwnerNotifiedOfResourceCreation"].(bool); ok {
        data.IsOwnerNotifiedOfResourceCreation = types.BoolValue(val)
    } else if dataMap["isOwnerNotifiedOfResourceCreation"] == nil {
        data.IsOwnerNotifiedOfResourceCreation = types.BoolNull()
    }
    if val, ok := dataMap["sendSubscriberNotificationsOnBeforeTheEvent"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.SendSubscriberNotificationsOnBeforeTheEvent = mapValue
    } else if dataMap["sendSubscriberNotificationsOnBeforeTheEvent"] == nil {
        data.SendSubscriberNotificationsOnBeforeTheEvent = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["nextSubscriberNotificationBeforeTheEventAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.NextSubscriberNotificationBeforeTheEventAt = mapValue
    } else if dataMap["nextSubscriberNotificationBeforeTheEventAt"] == nil {
        data.NextSubscriberNotificationBeforeTheEventAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["scheduledMaintenanceNumber"].(float64); ok {
        data.ScheduledMaintenanceNumber = types.NumberValue(big.NewFloat(val))
    } else if dataMap["scheduledMaintenanceNumber"] == nil {
        data.ScheduledMaintenanceNumber = types.NumberNull()
    }
    if val, ok := dataMap["isVisibleOnStatusPage"].(bool); ok {
        data.IsVisibleOnStatusPage = types.BoolValue(val)
    } else if dataMap["isVisibleOnStatusPage"] == nil {
        data.IsVisibleOnStatusPage = types.BoolNull()
    }
    if val, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.CreatedAt = mapValue
    } else if dataMap["createdAt"] == nil {
        data.CreatedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.UpdatedAt = mapValue
    } else if dataMap["updatedAt"] == nil {
        data.UpdatedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["deletedAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.DeletedAt = mapValue
    } else if dataMap["deletedAt"] == nil {
        data.DeletedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if dataMap["version"] == nil {
        data.Version = types.NumberNull()
    }
    if val, ok := dataMap["slug"].(string); ok && val != "" {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
    }
    if val, ok := dataMap["createdByUserId"].(string); ok && val != "" {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if val, ok := dataMap["isStatusPageSubscribersNotifiedOnEventScheduled"].(bool); ok {
        data.IsStatusPageSubscribersNotifiedOnEventScheduled = types.BoolValue(val)
    } else if dataMap["isStatusPageSubscribersNotifiedOnEventScheduled"] == nil {
        data.IsStatusPageSubscribersNotifiedOnEventScheduled = types.BoolNull()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ScheduledMaintenanceEventResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var data ScheduledMaintenanceEventResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Make API call
    _, err := r.client.Delete("/scheduled-maintenance/" + data.Id.ValueString() + "")
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete scheduled_maintenance_event, got error: %s", err))
        return
    }
}


func (r *ScheduledMaintenanceEventResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Helper method to convert Terraform map to Go interface{}
func (r *ScheduledMaintenanceEventResource) convertTerraformMapToInterface(terraformMap types.Map) interface{} {
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
func (r *ScheduledMaintenanceEventResource) convertTerraformListToInterface(terraformList types.List) interface{} {
    if terraformList.IsNull() || terraformList.IsUnknown() {
        return nil
    }
    
    var result []interface{}
    terraformList.ElementsAs(context.Background(), &result, false)
    return result
}
