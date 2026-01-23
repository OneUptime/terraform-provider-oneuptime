package provider

import (
    "context"
    "fmt"
    "math/big"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &WorkspaceNotificationLogDataDataSource{}

func NewWorkspaceNotificationLogDataDataSource() datasource.DataSource {
    return &WorkspaceNotificationLogDataDataSource{}
}

// WorkspaceNotificationLogDataDataSource defines the data source implementation.
type WorkspaceNotificationLogDataDataSource struct {
    client *Client
}

// WorkspaceNotificationLogDataDataSourceModel describes the data source data model.
type WorkspaceNotificationLogDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    WorkspaceType types.String `tfsdk:"workspace_type"`
    ChannelId types.String `tfsdk:"channel_id"`
    ChannelName types.String `tfsdk:"channel_name"`
    ThreadId types.String `tfsdk:"thread_id"`
    Message types.String `tfsdk:"message"`
    StatusMessage types.String `tfsdk:"status_message"`
    Status types.String `tfsdk:"status"`
    ActionType types.String `tfsdk:"action_type"`
    IncidentId types.String `tfsdk:"incident_id"`
    UserId types.String `tfsdk:"user_id"`
    AlertId types.String `tfsdk:"alert_id"`
    ScheduledMaintenanceId types.String `tfsdk:"scheduled_maintenance_id"`
    StatusPageId types.String `tfsdk:"status_page_id"`
    StatusPageAnnouncementId types.String `tfsdk:"status_page_announcement_id"`
    OnCallDutyPolicyId types.String `tfsdk:"on_call_duty_policy_id"`
    OnCallDutyPolicyEscalationRuleId types.String `tfsdk:"on_call_duty_policy_escalation_rule_id"`
    OnCallDutyPolicyScheduleId types.String `tfsdk:"on_call_duty_policy_schedule_id"`
    TeamId types.String `tfsdk:"team_id"`
}

func (d *WorkspaceNotificationLogDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_workspace_notification_log_data"
}

func (d *WorkspaceNotificationLogDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "workspace_notification_log_data data source",

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
            "workspace_type": schema.StringAttribute{
                MarkdownDescription: "Type of Workspace - Slack, Microsoft Teams. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Workspace Notification Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "channel_id": schema.StringAttribute{
                MarkdownDescription: "Channel ID where the message was sent. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Workspace Notification Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "channel_name": schema.StringAttribute{
                MarkdownDescription: "Channel Name where the message was sent. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Workspace Notification Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "thread_id": schema.StringAttribute{
                MarkdownDescription: "Thread ID of the message in the channel (if any). Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Workspace Notification Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "message": schema.StringAttribute{
                MarkdownDescription: "Content of the message. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Workspace Notification Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "status_message": schema.StringAttribute{
                MarkdownDescription: "Status Message (if any). Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Workspace Notification Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "status": schema.StringAttribute{
                MarkdownDescription: "Status of the message. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Workspace Notification Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "action_type": schema.StringAttribute{
                MarkdownDescription: "Type of workspace action performed. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Workspace Notification Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "incident_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "alert_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "scheduled_maintenance_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "status_page_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "status_page_announcement_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "on_call_duty_policy_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "on_call_duty_policy_escalation_rule_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "on_call_duty_policy_schedule_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "team_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
        },
    }
}

func (d *WorkspaceNotificationLogDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *WorkspaceNotificationLogDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data WorkspaceNotificationLogDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "workspace-notification-log" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read workspace_notification_log_data, got error: %s", err))
        return
    }

    var workspaceNotificationLogDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &workspaceNotificationLogDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse workspace_notification_log_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := workspaceNotificationLogDataResponse["data"].(map[string]interface{}); ok {
        workspaceNotificationLogDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := workspaceNotificationLogDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := workspaceNotificationLogDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := workspaceNotificationLogDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := workspaceNotificationLogDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := workspaceNotificationLogDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := workspaceNotificationLogDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := workspaceNotificationLogDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := workspaceNotificationLogDataResponse["workspace_type"].(string); ok {
        data.WorkspaceType = types.StringValue(val)
    }
    if val, ok := workspaceNotificationLogDataResponse["channel_id"].(string); ok {
        data.ChannelId = types.StringValue(val)
    }
    if val, ok := workspaceNotificationLogDataResponse["channel_name"].(string); ok {
        data.ChannelName = types.StringValue(val)
    }
    if val, ok := workspaceNotificationLogDataResponse["thread_id"].(string); ok {
        data.ThreadId = types.StringValue(val)
    }
    if val, ok := workspaceNotificationLogDataResponse["message"].(string); ok {
        data.Message = types.StringValue(val)
    }
    if val, ok := workspaceNotificationLogDataResponse["status_message"].(string); ok {
        data.StatusMessage = types.StringValue(val)
    }
    if val, ok := workspaceNotificationLogDataResponse["status"].(string); ok {
        data.Status = types.StringValue(val)
    }
    if val, ok := workspaceNotificationLogDataResponse["action_type"].(string); ok {
        data.ActionType = types.StringValue(val)
    }
    if val, ok := workspaceNotificationLogDataResponse["incident_id"].(string); ok {
        data.IncidentId = types.StringValue(val)
    }
    if val, ok := workspaceNotificationLogDataResponse["user_id"].(string); ok {
        data.UserId = types.StringValue(val)
    }
    if val, ok := workspaceNotificationLogDataResponse["alert_id"].(string); ok {
        data.AlertId = types.StringValue(val)
    }
    if val, ok := workspaceNotificationLogDataResponse["scheduled_maintenance_id"].(string); ok {
        data.ScheduledMaintenanceId = types.StringValue(val)
    }
    if val, ok := workspaceNotificationLogDataResponse["status_page_id"].(string); ok {
        data.StatusPageId = types.StringValue(val)
    }
    if val, ok := workspaceNotificationLogDataResponse["status_page_announcement_id"].(string); ok {
        data.StatusPageAnnouncementId = types.StringValue(val)
    }
    if val, ok := workspaceNotificationLogDataResponse["on_call_duty_policy_id"].(string); ok {
        data.OnCallDutyPolicyId = types.StringValue(val)
    }
    if val, ok := workspaceNotificationLogDataResponse["on_call_duty_policy_escalation_rule_id"].(string); ok {
        data.OnCallDutyPolicyEscalationRuleId = types.StringValue(val)
    }
    if val, ok := workspaceNotificationLogDataResponse["on_call_duty_policy_schedule_id"].(string); ok {
        data.OnCallDutyPolicyScheduleId = types.StringValue(val)
    }
    if val, ok := workspaceNotificationLogDataResponse["team_id"].(string); ok {
        data.TeamId = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
