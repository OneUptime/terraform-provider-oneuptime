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
var _ datasource.DataSource = &OnCallDutyExecutionLogDataDataSource{}

func NewOnCallDutyExecutionLogDataDataSource() datasource.DataSource {
    return &OnCallDutyExecutionLogDataDataSource{}
}

// OnCallDutyExecutionLogDataDataSource defines the data source implementation.
type OnCallDutyExecutionLogDataDataSource struct {
    client *Client
}

// OnCallDutyExecutionLogDataDataSourceModel describes the data source data model.
type OnCallDutyExecutionLogDataDataSourceModel struct {
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
    Status types.String `tfsdk:"status"`
    StatusMessage types.String `tfsdk:"status_message"`
    UserNotificationEventType types.String `tfsdk:"user_notification_event_type"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
    AcknowledgedByUserId types.String `tfsdk:"acknowledged_by_user_id"`
    AcknowledgedAt types.String `tfsdk:"acknowledged_at"`
    AcknowledgedByTeamId types.String `tfsdk:"acknowledged_by_team_id"`
    LastExecutedEscalationRuleId types.String `tfsdk:"last_executed_escalation_rule_id"`
    TriggeredByUserId types.String `tfsdk:"triggered_by_user_id"`
}

func (d *OnCallDutyExecutionLogDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_on_call_duty_execution_log_data"
}

func (d *OnCallDutyExecutionLogDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "on_call_duty_execution_log_data data source",

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
            "status": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create On-Call Duty Policy Execution Log], Read: [Project Owner, Project Admin, Project Member, Read On-Call Duty Policy Execution Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "status_message": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create On-Call Duty Policy Execution Log], Read: [Project Owner, Project Admin, Project Member, Read On-Call Duty Policy Execution Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "user_notification_event_type": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create On-Call Duty Policy Execution Log], Read: [Project Owner, Project Admin, Project Member, Read On-Call Duty Policy Execution Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "deleted_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "acknowledged_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "acknowledged_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "acknowledged_by_team_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "last_executed_escalation_rule_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "triggered_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
        },
    }
}

func (d *OnCallDutyExecutionLogDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *OnCallDutyExecutionLogDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data OnCallDutyExecutionLogDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "on-call-duty-policy-execution-log" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read on_call_duty_execution_log_data, got error: %s", err))
        return
    }

    var onCallDutyExecutionLogDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &onCallDutyExecutionLogDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse on_call_duty_execution_log_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := onCallDutyExecutionLogDataResponse["data"].(map[string]interface{}); ok {
        onCallDutyExecutionLogDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := onCallDutyExecutionLogDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := onCallDutyExecutionLogDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := onCallDutyExecutionLogDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := onCallDutyExecutionLogDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := onCallDutyExecutionLogDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := onCallDutyExecutionLogDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := onCallDutyExecutionLogDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := onCallDutyExecutionLogDataResponse["on_call_duty_policy_id"].(string); ok {
        data.OnCallDutyPolicyId = types.StringValue(val)
    }
    if val, ok := onCallDutyExecutionLogDataResponse["triggered_by_incident_id"].(string); ok {
        data.TriggeredByIncidentId = types.StringValue(val)
    }
    if val, ok := onCallDutyExecutionLogDataResponse["triggered_by_alert_id"].(string); ok {
        data.TriggeredByAlertId = types.StringValue(val)
    }
    if val, ok := onCallDutyExecutionLogDataResponse["status"].(string); ok {
        data.Status = types.StringValue(val)
    }
    if val, ok := onCallDutyExecutionLogDataResponse["status_message"].(string); ok {
        data.StatusMessage = types.StringValue(val)
    }
    if val, ok := onCallDutyExecutionLogDataResponse["user_notification_event_type"].(string); ok {
        data.UserNotificationEventType = types.StringValue(val)
    }
    if val, ok := onCallDutyExecutionLogDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := onCallDutyExecutionLogDataResponse["deleted_by_user_id"].(string); ok {
        data.DeletedByUserId = types.StringValue(val)
    }
    if val, ok := onCallDutyExecutionLogDataResponse["acknowledged_by_user_id"].(string); ok {
        data.AcknowledgedByUserId = types.StringValue(val)
    }
    if val, ok := onCallDutyExecutionLogDataResponse["acknowledged_at"].(string); ok {
        data.AcknowledgedAt = types.StringValue(val)
    }
    if val, ok := onCallDutyExecutionLogDataResponse["acknowledged_by_team_id"].(string); ok {
        data.AcknowledgedByTeamId = types.StringValue(val)
    }
    if val, ok := onCallDutyExecutionLogDataResponse["last_executed_escalation_rule_id"].(string); ok {
        data.LastExecutedEscalationRuleId = types.StringValue(val)
    }
    if val, ok := onCallDutyExecutionLogDataResponse["triggered_by_user_id"].(string); ok {
        data.TriggeredByUserId = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
