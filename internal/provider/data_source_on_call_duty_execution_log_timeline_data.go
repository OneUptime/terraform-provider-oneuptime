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
var _ datasource.DataSource = &OnCallDutyExecutionLogTimelineDataDataSource{}

func NewOnCallDutyExecutionLogTimelineDataDataSource() datasource.DataSource {
    return &OnCallDutyExecutionLogTimelineDataDataSource{}
}

// OnCallDutyExecutionLogTimelineDataDataSource defines the data source implementation.
type OnCallDutyExecutionLogTimelineDataDataSource struct {
    client *Client
}

// OnCallDutyExecutionLogTimelineDataDataSourceModel describes the data source data model.
type OnCallDutyExecutionLogTimelineDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    OnCallDutyPolicyId types.String `tfsdk:"on_call_duty_policy_id"`
    TriggeredByIncidentId types.String `tfsdk:"triggered_by_incident_id"`
    TriggeredByAlertId types.String `tfsdk:"triggered_by_alert_id"`
    OnCallDutyPolicyExecutionLogId types.String `tfsdk:"on_call_duty_policy_execution_log_id"`
    OnCallDutyPolicyEscalationRuleId types.String `tfsdk:"on_call_duty_policy_escalation_rule_id"`
    UserNotificationEventType types.String `tfsdk:"user_notification_event_type"`
    AlertSentToUserId types.String `tfsdk:"alert_sent_to_user_id"`
    UserBelongsToTeamId types.String `tfsdk:"user_belongs_to_team_id"`
    OnCallDutyScheduleId types.String `tfsdk:"on_call_duty_schedule_id"`
    StatusMessage types.String `tfsdk:"status_message"`
    Status types.String `tfsdk:"status"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    IsAcknowledged types.Bool `tfsdk:"is_acknowledged"`
    AcknowledgedAt types.String `tfsdk:"acknowledged_at"`
    OverridedByUserId types.String `tfsdk:"overrided_by_user_id"`
}

func (d *OnCallDutyExecutionLogTimelineDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_on_call_duty_execution_log_timeline_data"
}

func (d *OnCallDutyExecutionLogTimelineDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "on_call_duty_execution_log_timeline_data data source",

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
            "on_call_duty_policy_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "triggered_by_incident_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "triggered_by_alert_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "on_call_duty_policy_execution_log_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "on_call_duty_policy_escalation_rule_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "user_notification_event_type": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read On-Call Duty Policy Execution Log Timeline], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "alert_sent_to_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "user_belongs_to_team_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "on_call_duty_schedule_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "status_message": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read On-Call Duty Policy Execution Log Timeline], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "status": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read On-Call Duty Policy Execution Log Timeline], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "is_acknowledged": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read On-Call Duty Policy Execution Log Timeline], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "acknowledged_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "overrided_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
        },
    }
}

func (d *OnCallDutyExecutionLogTimelineDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *OnCallDutyExecutionLogTimelineDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data OnCallDutyExecutionLogTimelineDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "on-call-duty-policy-execution-log-timeline" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read on_call_duty_execution_log_timeline_data, got error: %s", err))
        return
    }

    var onCallDutyExecutionLogTimelineDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &onCallDutyExecutionLogTimelineDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse on_call_duty_execution_log_timeline_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := onCallDutyExecutionLogTimelineDataResponse["data"].(map[string]interface{}); ok {
        onCallDutyExecutionLogTimelineDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := onCallDutyExecutionLogTimelineDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := onCallDutyExecutionLogTimelineDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := onCallDutyExecutionLogTimelineDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := onCallDutyExecutionLogTimelineDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := onCallDutyExecutionLogTimelineDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := onCallDutyExecutionLogTimelineDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := onCallDutyExecutionLogTimelineDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := onCallDutyExecutionLogTimelineDataResponse["on_call_duty_policy_id"].(string); ok {
        data.OnCallDutyPolicyId = types.StringValue(val)
    }
    if val, ok := onCallDutyExecutionLogTimelineDataResponse["triggered_by_incident_id"].(string); ok {
        data.TriggeredByIncidentId = types.StringValue(val)
    }
    if val, ok := onCallDutyExecutionLogTimelineDataResponse["triggered_by_alert_id"].(string); ok {
        data.TriggeredByAlertId = types.StringValue(val)
    }
    if val, ok := onCallDutyExecutionLogTimelineDataResponse["on_call_duty_policy_execution_log_id"].(string); ok {
        data.OnCallDutyPolicyExecutionLogId = types.StringValue(val)
    }
    if val, ok := onCallDutyExecutionLogTimelineDataResponse["on_call_duty_policy_escalation_rule_id"].(string); ok {
        data.OnCallDutyPolicyEscalationRuleId = types.StringValue(val)
    }
    if val, ok := onCallDutyExecutionLogTimelineDataResponse["user_notification_event_type"].(string); ok {
        data.UserNotificationEventType = types.StringValue(val)
    }
    if val, ok := onCallDutyExecutionLogTimelineDataResponse["alert_sent_to_user_id"].(string); ok {
        data.AlertSentToUserId = types.StringValue(val)
    }
    if val, ok := onCallDutyExecutionLogTimelineDataResponse["user_belongs_to_team_id"].(string); ok {
        data.UserBelongsToTeamId = types.StringValue(val)
    }
    if val, ok := onCallDutyExecutionLogTimelineDataResponse["on_call_duty_schedule_id"].(string); ok {
        data.OnCallDutyScheduleId = types.StringValue(val)
    }
    if val, ok := onCallDutyExecutionLogTimelineDataResponse["status_message"].(string); ok {
        data.StatusMessage = types.StringValue(val)
    }
    if val, ok := onCallDutyExecutionLogTimelineDataResponse["status"].(string); ok {
        data.Status = types.StringValue(val)
    }
    if val, ok := onCallDutyExecutionLogTimelineDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := onCallDutyExecutionLogTimelineDataResponse["is_acknowledged"].(bool); ok {
        data.IsAcknowledged = types.BoolValue(val)
    }
    if val, ok := onCallDutyExecutionLogTimelineDataResponse["acknowledged_at"].(string); ok {
        data.AcknowledgedAt = types.StringValue(val)
    }
    if val, ok := onCallDutyExecutionLogTimelineDataResponse["overrided_by_user_id"].(string); ok {
        data.OverridedByUserId = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
