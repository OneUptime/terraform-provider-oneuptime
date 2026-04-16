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
var _ datasource.DataSource = &CallLogDataDataSource{}

func NewCallLogDataDataSource() datasource.DataSource {
    return &CallLogDataDataSource{}
}

// CallLogDataDataSource defines the data source implementation.
type CallLogDataDataSource struct {
    client *Client
}

// CallLogDataDataSourceModel describes the data source data model.
type CallLogDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    ToNumber types.String `tfsdk:"to_number"`
    FromNumber types.String `tfsdk:"from_number"`
    CallData types.String `tfsdk:"call_data"`
    StatusMessage types.String `tfsdk:"status_message"`
    Status types.String `tfsdk:"status"`
    CallCostInUsdCents types.Number `tfsdk:"call_cost_in_usd_cents"`
    IncidentId types.String `tfsdk:"incident_id"`
    UserId types.String `tfsdk:"user_id"`
    AlertId types.String `tfsdk:"alert_id"`
    MonitorId types.String `tfsdk:"monitor_id"`
    ScheduledMaintenanceId types.String `tfsdk:"scheduled_maintenance_id"`
    StatusPageId types.String `tfsdk:"status_page_id"`
    StatusPageAnnouncementId types.String `tfsdk:"status_page_announcement_id"`
    OnCallDutyPolicyId types.String `tfsdk:"on_call_duty_policy_id"`
    OnCallDutyPolicyEscalationRuleId types.String `tfsdk:"on_call_duty_policy_escalation_rule_id"`
    OnCallDutyPolicyScheduleId types.String `tfsdk:"on_call_duty_policy_schedule_id"`
    TeamId types.String `tfsdk:"team_id"`
}

func (d *CallLogDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_call_log_data"
}

func (d *CallLogDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "call_log_data data source",

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
            "to_number": schema.StringAttribute{
                MarkdownDescription: "Phone object",
                Computed: true,
            },
            "from_number": schema.StringAttribute{
                MarkdownDescription: "Phone object",
                Computed: true,
            },
            "call_data": schema.StringAttribute{
                MarkdownDescription: "Content of the data that was sent in the call. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Call Log, Read All Project Resources], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "status_message": schema.StringAttribute{
                MarkdownDescription: "Status Message (if any). Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Call Log, Read All Project Resources], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "status": schema.StringAttribute{
                MarkdownDescription: "Status of the Call sent. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Call Log, Read All Project Resources], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "call_cost_in_usd_cents": schema.NumberAttribute{
                MarkdownDescription: "Call Cost in USD Cents. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Call Log, Read All Project Resources], Update: [No access - you don't have permission for this operation]",
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

func (d *CallLogDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *CallLogDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data CallLogDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "call-log" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read call_log_data, got error: %s", err))
        return
    }

    var callLogDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &callLogDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse call_log_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := callLogDataResponse["data"].(map[string]interface{}); ok {
        callLogDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := callLogDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := callLogDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := callLogDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := callLogDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := callLogDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := callLogDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := callLogDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := callLogDataResponse["to_number"].(string); ok {
        data.ToNumber = types.StringValue(val)
    }
    if val, ok := callLogDataResponse["from_number"].(string); ok {
        data.FromNumber = types.StringValue(val)
    }
    if val, ok := callLogDataResponse["call_data"].(string); ok {
        data.CallData = types.StringValue(val)
    }
    if val, ok := callLogDataResponse["status_message"].(string); ok {
        data.StatusMessage = types.StringValue(val)
    }
    if val, ok := callLogDataResponse["status"].(string); ok {
        data.Status = types.StringValue(val)
    }
    if val, ok := callLogDataResponse["call_cost_in_usd_cents"].(float64); ok {
        data.CallCostInUsdCents = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := callLogDataResponse["incident_id"].(string); ok {
        data.IncidentId = types.StringValue(val)
    }
    if val, ok := callLogDataResponse["user_id"].(string); ok {
        data.UserId = types.StringValue(val)
    }
    if val, ok := callLogDataResponse["alert_id"].(string); ok {
        data.AlertId = types.StringValue(val)
    }
    if val, ok := callLogDataResponse["monitor_id"].(string); ok {
        data.MonitorId = types.StringValue(val)
    }
    if val, ok := callLogDataResponse["scheduled_maintenance_id"].(string); ok {
        data.ScheduledMaintenanceId = types.StringValue(val)
    }
    if val, ok := callLogDataResponse["status_page_id"].(string); ok {
        data.StatusPageId = types.StringValue(val)
    }
    if val, ok := callLogDataResponse["status_page_announcement_id"].(string); ok {
        data.StatusPageAnnouncementId = types.StringValue(val)
    }
    if val, ok := callLogDataResponse["on_call_duty_policy_id"].(string); ok {
        data.OnCallDutyPolicyId = types.StringValue(val)
    }
    if val, ok := callLogDataResponse["on_call_duty_policy_escalation_rule_id"].(string); ok {
        data.OnCallDutyPolicyEscalationRuleId = types.StringValue(val)
    }
    if val, ok := callLogDataResponse["on_call_duty_policy_schedule_id"].(string); ok {
        data.OnCallDutyPolicyScheduleId = types.StringValue(val)
    }
    if val, ok := callLogDataResponse["team_id"].(string); ok {
        data.TeamId = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
