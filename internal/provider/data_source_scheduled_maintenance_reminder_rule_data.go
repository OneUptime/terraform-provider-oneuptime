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
var _ datasource.DataSource = &ScheduledMaintenanceReminderRuleDataDataSource{}

func NewScheduledMaintenanceReminderRuleDataDataSource() datasource.DataSource {
    return &ScheduledMaintenanceReminderRuleDataDataSource{}
}

// ScheduledMaintenanceReminderRuleDataDataSource defines the data source implementation.
type ScheduledMaintenanceReminderRuleDataDataSource struct {
    client *Client
}

// ScheduledMaintenanceReminderRuleDataDataSourceModel describes the data source data model.
type ScheduledMaintenanceReminderRuleDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    Description types.String `tfsdk:"description"`
    Order types.Number `tfsdk:"order"`
    IsEnabled types.Bool `tfsdk:"is_enabled"`
    ReminderIntervalInMinutes types.Number `tfsdk:"reminder_interval_in_minutes"`
    StopRemindersOnState types.String `tfsdk:"stop_reminders_on_state"`
    RemindWhileScheduled types.Bool `tfsdk:"remind_while_scheduled"`
    Labels types.Set `tfsdk:"labels"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *ScheduledMaintenanceReminderRuleDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_scheduled_maintenance_reminder_rule_data"
}

func (d *ScheduledMaintenanceReminderRuleDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "scheduled_maintenance_reminder_rule_data data source",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Identifier to filter by",
                Optional: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Name to filter by",
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
            "project_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Description of this reminder rule. Permissions - Create: [Project Owner, Project Admin, Create Scheduled Maintenance Reminder Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Scheduled Maintenance Admin, Scheduled Maintenance Member, Scheduled Maintenance Viewer, Read Scheduled Maintenance Reminder Rule], Update: [Project Owner, Project Admin, Edit Scheduled Maintenance Reminder Rule]",
                Computed: true,
            },
            "order": schema.NumberAttribute{
                MarkdownDescription: "Order/priority of this rule. Rules are evaluated in order (lowest first). First matching rule wins.. Permissions - Create: [Project Owner, Project Admin, Create Scheduled Maintenance Reminder Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Scheduled Maintenance Admin, Scheduled Maintenance Member, Scheduled Maintenance Viewer, Read Scheduled Maintenance Reminder Rule], Update: [Project Owner, Project Admin, Edit Scheduled Maintenance Reminder Rule]",
                Computed: true,
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether this reminder rule is enabled. Permissions - Create: [Project Owner, Project Admin, Create Scheduled Maintenance Reminder Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Scheduled Maintenance Admin, Scheduled Maintenance Member, Scheduled Maintenance Viewer, Read Scheduled Maintenance Reminder Rule], Update: [Project Owner, Project Admin, Edit Scheduled Maintenance Reminder Rule]",
                Computed: true,
            },
            "reminder_interval_in_minutes": schema.NumberAttribute{
                MarkdownDescription: "How often (in minutes) to remind scheduled maintenance event owners while the event is still not complete. For example, set to 30 to remind owners every 30 minutes.. Permissions - Create: [Project Owner, Project Admin, Create Scheduled Maintenance Reminder Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Scheduled Maintenance Admin, Scheduled Maintenance Member, Scheduled Maintenance Viewer, Read Scheduled Maintenance Reminder Rule], Update: [Project Owner, Project Admin, Edit Scheduled Maintenance Reminder Rule]",
                Computed: true,
            },
            "stop_reminders_on_state": schema.StringAttribute{
                MarkdownDescription: "Stop sending reminders once the scheduled maintenance event reaches this state. Select Ongoing to stop reminders when the event starts, or Completed to keep reminding until the event is completed.. Permissions - Create: [Project Owner, Project Admin, Create Scheduled Maintenance Reminder Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Scheduled Maintenance Admin, Scheduled Maintenance Member, Scheduled Maintenance Viewer, Read Scheduled Maintenance Reminder Rule], Update: [Project Owner, Project Admin, Edit Scheduled Maintenance Reminder Rule]",
                Computed: true,
            },
            "remind_while_scheduled": schema.BoolAttribute{
                MarkdownDescription: "Send reminders while the event is still scheduled (before it starts). When disabled, reminders only begin once the event has started.. Permissions - Create: [Project Owner, Project Admin, Create Scheduled Maintenance Reminder Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Scheduled Maintenance Admin, Scheduled Maintenance Member, Scheduled Maintenance Viewer, Read Scheduled Maintenance Reminder Rule], Update: [Project Owner, Project Admin, Edit Scheduled Maintenance Reminder Rule]",
                Computed: true,
            },
            "labels": schema.SetAttribute{
                MarkdownDescription: "Only apply this reminder rule to scheduled maintenance events with these labels. Leave empty to match all events.. Permissions - Create: [Project Owner, Project Admin, Create Scheduled Maintenance Reminder Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Scheduled Maintenance Admin, Scheduled Maintenance Member, Scheduled Maintenance Viewer, Read Scheduled Maintenance Reminder Rule], Update: [Project Owner, Project Admin, Edit Scheduled Maintenance Reminder Rule]",
                Computed: true,
                ElementType: types.StringType,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
        },
    }
}

func (d *ScheduledMaintenanceReminderRuleDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ScheduledMaintenanceReminderRuleDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data ScheduledMaintenanceReminderRuleDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "scheduled-maintenance-reminder-rule" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read scheduled_maintenance_reminder_rule_data, got error: %s", err))
        return
    }

    var scheduledMaintenanceReminderRuleDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &scheduledMaintenanceReminderRuleDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse scheduled_maintenance_reminder_rule_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := scheduledMaintenanceReminderRuleDataResponse["data"].(map[string]interface{}); ok {
        scheduledMaintenanceReminderRuleDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := scheduledMaintenanceReminderRuleDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceReminderRuleDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceReminderRuleDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceReminderRuleDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceReminderRuleDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceReminderRuleDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := scheduledMaintenanceReminderRuleDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceReminderRuleDataResponse["description"].(string); ok {
        data.Description = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceReminderRuleDataResponse["order"].(float64); ok {
        data.Order = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := scheduledMaintenanceReminderRuleDataResponse["is_enabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    }
    if val, ok := scheduledMaintenanceReminderRuleDataResponse["reminder_interval_in_minutes"].(float64); ok {
        data.ReminderIntervalInMinutes = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := scheduledMaintenanceReminderRuleDataResponse["stop_reminders_on_state"].(string); ok {
        data.StopRemindersOnState = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceReminderRuleDataResponse["remind_while_scheduled"].(bool); ok {
        data.RemindWhileScheduled = types.BoolValue(val)
    }
    if val, ok := scheduledMaintenanceReminderRuleDataResponse["labels"].([]interface{}); ok {
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
    if val, ok := scheduledMaintenanceReminderRuleDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
