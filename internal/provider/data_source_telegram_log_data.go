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
var _ datasource.DataSource = &TelegramLogDataDataSource{}

func NewTelegramLogDataDataSource() datasource.DataSource {
    return &TelegramLogDataDataSource{}
}

// TelegramLogDataDataSource defines the data source implementation.
type TelegramLogDataDataSource struct {
    client *Client
}

// TelegramLogDataDataSourceModel describes the data source data model.
type TelegramLogDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    ToChatId types.String `tfsdk:"to_chat_id"`
    FromBotUsername types.String `tfsdk:"from_bot_username"`
    MessageText types.String `tfsdk:"message_text"`
    StatusMessage types.String `tfsdk:"status_message"`
    TelegramMessageId types.String `tfsdk:"telegram_message_id"`
    Status types.String `tfsdk:"status"`
    TelegramCostInUsdCents types.Number `tfsdk:"telegram_cost_in_usd_cents"`
    IncidentId types.String `tfsdk:"incident_id"`
    UserId types.String `tfsdk:"user_id"`
    AlertId types.String `tfsdk:"alert_id"`
    MonitorId types.String `tfsdk:"monitor_id"`
    ScheduledMaintenanceId types.String `tfsdk:"scheduled_maintenance_id"`
    StatusPageId types.String `tfsdk:"status_page_id"`
    StatusPageAnnouncementId types.String `tfsdk:"status_page_announcement_id"`
    TeamId types.String `tfsdk:"team_id"`
    OnCallDutyPolicyId types.String `tfsdk:"on_call_duty_policy_id"`
    OnCallDutyPolicyEscalationRuleId types.String `tfsdk:"on_call_duty_policy_escalation_rule_id"`
    OnCallDutyPolicyScheduleId types.String `tfsdk:"on_call_duty_policy_schedule_id"`
}

func (d *TelegramLogDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_telegram_log_data"
}

func (d *TelegramLogDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "telegram_log_data data source",

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
            "to_chat_id": schema.StringAttribute{
                MarkdownDescription: "Telegram Chat ID the message was sent to. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Telegram Log, Read All Project Resources], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "from_bot_username": schema.StringAttribute{
                MarkdownDescription: "OneUptime Telegram bot username the message was sent from. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Telegram Log, Read All Project Resources], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "message_text": schema.StringAttribute{
                MarkdownDescription: "Text content of the Telegram message. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Telegram Log, Read All Project Resources], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "status_message": schema.StringAttribute{
                MarkdownDescription: "Status Message (if any). Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Telegram Log, Read All Project Resources], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "telegram_message_id": schema.StringAttribute{
                MarkdownDescription: "Message ID returned by Telegram Bot API. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Telegram Log, Read All Project Resources], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "status": schema.StringAttribute{
                MarkdownDescription: "Status of the Telegram message sent. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Telegram Log, Read All Project Resources], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "telegram_cost_in_usd_cents": schema.NumberAttribute{
                MarkdownDescription: "Telegram Message Cost in USD Cents. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Telegram Log, Read All Project Resources], Update: [No access - you don't have permission for this operation]",
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
            "monitor_id": schema.StringAttribute{
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
            "team_id": schema.StringAttribute{
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
        },
    }
}

func (d *TelegramLogDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *TelegramLogDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data TelegramLogDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "telegram-log" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read telegram_log_data, got error: %s", err))
        return
    }

    var telegramLogDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &telegramLogDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse telegram_log_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := telegramLogDataResponse["data"].(map[string]interface{}); ok {
        telegramLogDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := telegramLogDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := telegramLogDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := telegramLogDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := telegramLogDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := telegramLogDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := telegramLogDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := telegramLogDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := telegramLogDataResponse["to_chat_id"].(string); ok {
        data.ToChatId = types.StringValue(val)
    }
    if val, ok := telegramLogDataResponse["from_bot_username"].(string); ok {
        data.FromBotUsername = types.StringValue(val)
    }
    if val, ok := telegramLogDataResponse["message_text"].(string); ok {
        data.MessageText = types.StringValue(val)
    }
    if val, ok := telegramLogDataResponse["status_message"].(string); ok {
        data.StatusMessage = types.StringValue(val)
    }
    if val, ok := telegramLogDataResponse["telegram_message_id"].(string); ok {
        data.TelegramMessageId = types.StringValue(val)
    }
    if val, ok := telegramLogDataResponse["status"].(string); ok {
        data.Status = types.StringValue(val)
    }
    if val, ok := telegramLogDataResponse["telegram_cost_in_usd_cents"].(float64); ok {
        data.TelegramCostInUsdCents = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := telegramLogDataResponse["incident_id"].(string); ok {
        data.IncidentId = types.StringValue(val)
    }
    if val, ok := telegramLogDataResponse["user_id"].(string); ok {
        data.UserId = types.StringValue(val)
    }
    if val, ok := telegramLogDataResponse["alert_id"].(string); ok {
        data.AlertId = types.StringValue(val)
    }
    if val, ok := telegramLogDataResponse["monitor_id"].(string); ok {
        data.MonitorId = types.StringValue(val)
    }
    if val, ok := telegramLogDataResponse["scheduled_maintenance_id"].(string); ok {
        data.ScheduledMaintenanceId = types.StringValue(val)
    }
    if val, ok := telegramLogDataResponse["status_page_id"].(string); ok {
        data.StatusPageId = types.StringValue(val)
    }
    if val, ok := telegramLogDataResponse["status_page_announcement_id"].(string); ok {
        data.StatusPageAnnouncementId = types.StringValue(val)
    }
    if val, ok := telegramLogDataResponse["team_id"].(string); ok {
        data.TeamId = types.StringValue(val)
    }
    if val, ok := telegramLogDataResponse["on_call_duty_policy_id"].(string); ok {
        data.OnCallDutyPolicyId = types.StringValue(val)
    }
    if val, ok := telegramLogDataResponse["on_call_duty_policy_escalation_rule_id"].(string); ok {
        data.OnCallDutyPolicyEscalationRuleId = types.StringValue(val)
    }
    if val, ok := telegramLogDataResponse["on_call_duty_policy_schedule_id"].(string); ok {
        data.OnCallDutyPolicyScheduleId = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
