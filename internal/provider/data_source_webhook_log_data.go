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
var _ datasource.DataSource = &WebhookLogDataDataSource{}

func NewWebhookLogDataDataSource() datasource.DataSource {
    return &WebhookLogDataDataSource{}
}

// WebhookLogDataDataSource defines the data source implementation.
type WebhookLogDataDataSource struct {
    client *Client
}

// WebhookLogDataDataSourceModel describes the data source data model.
type WebhookLogDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    WebhookUrl types.String `tfsdk:"webhook_url"`
    RequestBody types.String `tfsdk:"request_body"`
    ResponseStatusCode types.Number `tfsdk:"response_status_code"`
    ResponseBody types.String `tfsdk:"response_body"`
    StatusMessage types.String `tfsdk:"status_message"`
    Status types.String `tfsdk:"status"`
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

func (d *WebhookLogDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_webhook_log_data"
}

func (d *WebhookLogDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "webhook_log_data data source",

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
            "webhook_url": schema.StringAttribute{
                MarkdownDescription: "URL the request was sent to. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Webhook Log, Read All Project Resources], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "request_body": schema.StringAttribute{
                MarkdownDescription: "JSON body that was POSTed to the webhook URL. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Webhook Log, Read All Project Resources], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "response_status_code": schema.NumberAttribute{
                MarkdownDescription: "HTTP status code returned by the webhook endpoint. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Webhook Log, Read All Project Resources], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "response_body": schema.StringAttribute{
                MarkdownDescription: "Response body returned by the webhook endpoint (truncated). Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Webhook Log, Read All Project Resources], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "status_message": schema.StringAttribute{
                MarkdownDescription: "Status Message (if any). Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Webhook Log, Read All Project Resources], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "status": schema.StringAttribute{
                MarkdownDescription: "Status of the Webhook request. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Webhook Log, Read All Project Resources], Update: [No access - you don't have permission for this operation]",
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

func (d *WebhookLogDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *WebhookLogDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data WebhookLogDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "webhook-log" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read webhook_log_data, got error: %s", err))
        return
    }

    var webhookLogDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &webhookLogDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse webhook_log_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := webhookLogDataResponse["data"].(map[string]interface{}); ok {
        webhookLogDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := webhookLogDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := webhookLogDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := webhookLogDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := webhookLogDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := webhookLogDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := webhookLogDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := webhookLogDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := webhookLogDataResponse["webhook_url"].(string); ok {
        data.WebhookUrl = types.StringValue(val)
    }
    if val, ok := webhookLogDataResponse["request_body"].(string); ok {
        data.RequestBody = types.StringValue(val)
    }
    if val, ok := webhookLogDataResponse["response_status_code"].(float64); ok {
        data.ResponseStatusCode = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := webhookLogDataResponse["response_body"].(string); ok {
        data.ResponseBody = types.StringValue(val)
    }
    if val, ok := webhookLogDataResponse["status_message"].(string); ok {
        data.StatusMessage = types.StringValue(val)
    }
    if val, ok := webhookLogDataResponse["status"].(string); ok {
        data.Status = types.StringValue(val)
    }
    if val, ok := webhookLogDataResponse["incident_id"].(string); ok {
        data.IncidentId = types.StringValue(val)
    }
    if val, ok := webhookLogDataResponse["user_id"].(string); ok {
        data.UserId = types.StringValue(val)
    }
    if val, ok := webhookLogDataResponse["alert_id"].(string); ok {
        data.AlertId = types.StringValue(val)
    }
    if val, ok := webhookLogDataResponse["monitor_id"].(string); ok {
        data.MonitorId = types.StringValue(val)
    }
    if val, ok := webhookLogDataResponse["scheduled_maintenance_id"].(string); ok {
        data.ScheduledMaintenanceId = types.StringValue(val)
    }
    if val, ok := webhookLogDataResponse["status_page_id"].(string); ok {
        data.StatusPageId = types.StringValue(val)
    }
    if val, ok := webhookLogDataResponse["status_page_announcement_id"].(string); ok {
        data.StatusPageAnnouncementId = types.StringValue(val)
    }
    if val, ok := webhookLogDataResponse["team_id"].(string); ok {
        data.TeamId = types.StringValue(val)
    }
    if val, ok := webhookLogDataResponse["on_call_duty_policy_id"].(string); ok {
        data.OnCallDutyPolicyId = types.StringValue(val)
    }
    if val, ok := webhookLogDataResponse["on_call_duty_policy_escalation_rule_id"].(string); ok {
        data.OnCallDutyPolicyEscalationRuleId = types.StringValue(val)
    }
    if val, ok := webhookLogDataResponse["on_call_duty_policy_schedule_id"].(string); ok {
        data.OnCallDutyPolicyScheduleId = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
