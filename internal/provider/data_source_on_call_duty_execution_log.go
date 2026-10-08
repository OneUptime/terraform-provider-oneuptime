package provider

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "math/big"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &OnCallDutyExecutionLogDataSource{}

func NewOnCallDutyExecutionLogDataSource() datasource.DataSource {
    return &OnCallDutyExecutionLogDataSource{}
}

// OnCallDutyExecutionLogDataSource defines the data source implementation.
type OnCallDutyExecutionLogDataSource struct {
    client *Client
}

// OnCallDutyExecutionLogDataSourceModel describes the data source data model.
type OnCallDutyExecutionLogDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    OnCallDutyPolicyId types.String `tfsdk:"on_call_duty_policy_id"`
    TriggeredByIncidentId types.String `tfsdk:"triggered_by_incident_id"`
    TriggeredByAlertId types.String `tfsdk:"triggered_by_alert_id"`
    TriggeredByAlertEpisodeId types.String `tfsdk:"triggered_by_alert_episode_id"`
    TriggeredByIncidentEpisodeId types.String `tfsdk:"triggered_by_incident_episode_id"`
    Status types.String `tfsdk:"status"`
    StatusMessage types.String `tfsdk:"status_message"`
    UserNotificationEventType types.String `tfsdk:"user_notification_event_type"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    AcknowledgedByUserId types.String `tfsdk:"acknowledged_by_user_id"`
    AcknowledgedAt types.String `tfsdk:"acknowledged_at"`
    AcknowledgedByTeamId types.String `tfsdk:"acknowledged_by_team_id"`
    LastExecutedEscalationRuleOrder types.Number `tfsdk:"last_executed_escalation_rule_order"`
    LastExecutedEscalationRuleId types.String `tfsdk:"last_executed_escalation_rule_id"`
    OnCallPolicyExecutionRepeatCount types.Number `tfsdk:"on_call_policy_execution_repeat_count"`
    ScheduleGapRetryCount types.Number `tfsdk:"schedule_gap_retry_count"`
    TriggeredByUserId types.String `tfsdk:"triggered_by_user_id"`
}

func (d *OnCallDutyExecutionLogDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_on_call_duty_execution_log"
}

func (d *OnCallDutyExecutionLogDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Logs for on-call duty policy execution. Look up an existing on call duty execution log by `id`, or by any of its other arguments (`acknowledged_by_team_id`, `acknowledged_by_user_id`, `created_by_user_id`, ...): each one set must match, and exactly one on call duty execution log may match them all.",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Look up by unique identifier. Leave unset to look up by the other arguments instead.",
                Optional: true,
                Computed: true,
            },
            "created_at": schema.StringAttribute{
                MarkdownDescription: "Date and Time when the object was created.",
                Computed: true,
            },
            "updated_at": schema.StringAttribute{
                MarkdownDescription: "Date and Time when the object was updated.",
                Computed: true,
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "on_call_duty_policy_id": schema.StringAttribute{
                MarkdownDescription: "ID of your On-Call Policy which belongs to this execution log event. The ID of a `oneuptime_on_call_policy`.",
                Optional: true,
                Computed: true,
            },
            "triggered_by_incident_id": schema.StringAttribute{
                MarkdownDescription: "ID of the incident which triggered this on-call escalation policy. The ID of a `oneuptime_incident`.",
                Optional: true,
                Computed: true,
            },
            "triggered_by_alert_id": schema.StringAttribute{
                MarkdownDescription: "ID of the incident which triggered this on-call escalation policy. The ID of a `oneuptime_alert`.",
                Optional: true,
                Computed: true,
            },
            "triggered_by_alert_episode_id": schema.StringAttribute{
                MarkdownDescription: "ID of the alert episode which triggered this on-call escalation policy. The ID of a `oneuptime_alert_episode`.",
                Optional: true,
                Computed: true,
            },
            "triggered_by_incident_episode_id": schema.StringAttribute{
                MarkdownDescription: "ID of the incident episode which triggered this on-call escalation policy. The ID of a `oneuptime_incident_episode`.",
                Optional: true,
                Computed: true,
            },
            "status": schema.StringAttribute{
                MarkdownDescription: "Status of this execution.",
                Optional: true,
                Computed: true,
            },
            "status_message": schema.StringAttribute{
                MarkdownDescription: "Status message of this execution.",
                Optional: true,
                Computed: true,
            },
            "user_notification_event_type": schema.StringAttribute{
                MarkdownDescription: "Type of event that triggered this on-call duty policy.",
                Optional: true,
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "acknowledged_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who acknowledged this object (if this object was acknowledged by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "acknowledged_at": schema.StringAttribute{
                MarkdownDescription: "When was this policy execution acknowledged?",
                Computed: true,
            },
            "acknowledged_by_team_id": schema.StringAttribute{
                MarkdownDescription: "Team ID who acknowledged this object (if this object was acknowledged by a Team). The ID of a `oneuptime_team`.",
                Optional: true,
                Computed: true,
            },
            "last_executed_escalation_rule_order": schema.NumberAttribute{
                MarkdownDescription: "Which escalation rule was executed?",
                Optional: true,
                Computed: true,
            },
            "last_executed_escalation_rule_id": schema.StringAttribute{
                MarkdownDescription: "ID of your On-Call Policy Last Executed Escalation Rule. The ID of a `oneuptime_escalation_rule`.",
                Optional: true,
                Computed: true,
            },
            "on_call_policy_execution_repeat_count": schema.NumberAttribute{
                MarkdownDescription: "How many times did we execute this on-call policy?",
                Optional: true,
                Computed: true,
            },
            "schedule_gap_retry_count": schema.NumberAttribute{
                MarkdownDescription: "How many times the current escalation rule has been re-sampled because its target schedule(s) momentarily had no on-call user.",
                Optional: true,
                Computed: true,
            },
            "triggered_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who triggered this on-call policy. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *OnCallDutyExecutionLogDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *OnCallDutyExecutionLogDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data OnCallDutyExecutionLogDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.OnCallDutyPolicyId.IsNull() && !data.OnCallDutyPolicyId.IsUnknown() {
        filters["onCallDutyPolicyId"] = data.OnCallDutyPolicyId.ValueString()
        filterNames = append(filterNames, "on_call_duty_policy_id = "+fmt.Sprintf("%q", data.OnCallDutyPolicyId.ValueString()))
    }
    if !data.TriggeredByIncidentId.IsNull() && !data.TriggeredByIncidentId.IsUnknown() {
        filters["triggeredByIncidentId"] = data.TriggeredByIncidentId.ValueString()
        filterNames = append(filterNames, "triggered_by_incident_id = "+fmt.Sprintf("%q", data.TriggeredByIncidentId.ValueString()))
    }
    if !data.TriggeredByAlertId.IsNull() && !data.TriggeredByAlertId.IsUnknown() {
        filters["triggeredByAlertId"] = data.TriggeredByAlertId.ValueString()
        filterNames = append(filterNames, "triggered_by_alert_id = "+fmt.Sprintf("%q", data.TriggeredByAlertId.ValueString()))
    }
    if !data.TriggeredByAlertEpisodeId.IsNull() && !data.TriggeredByAlertEpisodeId.IsUnknown() {
        filters["triggeredByAlertEpisodeId"] = data.TriggeredByAlertEpisodeId.ValueString()
        filterNames = append(filterNames, "triggered_by_alert_episode_id = "+fmt.Sprintf("%q", data.TriggeredByAlertEpisodeId.ValueString()))
    }
    if !data.TriggeredByIncidentEpisodeId.IsNull() && !data.TriggeredByIncidentEpisodeId.IsUnknown() {
        filters["triggeredByIncidentEpisodeId"] = data.TriggeredByIncidentEpisodeId.ValueString()
        filterNames = append(filterNames, "triggered_by_incident_episode_id = "+fmt.Sprintf("%q", data.TriggeredByIncidentEpisodeId.ValueString()))
    }
    if !data.Status.IsNull() && !data.Status.IsUnknown() {
        filters["status"] = data.Status.ValueString()
        filterNames = append(filterNames, "status = "+fmt.Sprintf("%q", data.Status.ValueString()))
    }
    if !data.StatusMessage.IsNull() && !data.StatusMessage.IsUnknown() {
        filters["statusMessage"] = data.StatusMessage.ValueString()
        filterNames = append(filterNames, "status_message = "+fmt.Sprintf("%q", data.StatusMessage.ValueString()))
    }
    if !data.UserNotificationEventType.IsNull() && !data.UserNotificationEventType.IsUnknown() {
        filters["userNotificationEventType"] = data.UserNotificationEventType.ValueString()
        filterNames = append(filterNames, "user_notification_event_type = "+fmt.Sprintf("%q", data.UserNotificationEventType.ValueString()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }
    if !data.AcknowledgedByUserId.IsNull() && !data.AcknowledgedByUserId.IsUnknown() {
        filters["acknowledgedByUserId"] = data.AcknowledgedByUserId.ValueString()
        filterNames = append(filterNames, "acknowledged_by_user_id = "+fmt.Sprintf("%q", data.AcknowledgedByUserId.ValueString()))
    }
    if !data.AcknowledgedByTeamId.IsNull() && !data.AcknowledgedByTeamId.IsUnknown() {
        filters["acknowledgedByTeamId"] = data.AcknowledgedByTeamId.ValueString()
        filterNames = append(filterNames, "acknowledged_by_team_id = "+fmt.Sprintf("%q", data.AcknowledgedByTeamId.ValueString()))
    }
    if !data.LastExecutedEscalationRuleOrder.IsNull() && !data.LastExecutedEscalationRuleOrder.IsUnknown() {
        filters["lastExecutedEscalationRuleOrder"] = lookupNumber(data.LastExecutedEscalationRuleOrder)
        filterNames = append(filterNames, "last_executed_escalation_rule_order = "+data.LastExecutedEscalationRuleOrder.ValueBigFloat().String())
    }
    if !data.LastExecutedEscalationRuleId.IsNull() && !data.LastExecutedEscalationRuleId.IsUnknown() {
        filters["lastExecutedEscalationRuleId"] = data.LastExecutedEscalationRuleId.ValueString()
        filterNames = append(filterNames, "last_executed_escalation_rule_id = "+fmt.Sprintf("%q", data.LastExecutedEscalationRuleId.ValueString()))
    }
    if !data.OnCallPolicyExecutionRepeatCount.IsNull() && !data.OnCallPolicyExecutionRepeatCount.IsUnknown() {
        filters["onCallPolicyExecutionRepeatCount"] = lookupNumber(data.OnCallPolicyExecutionRepeatCount)
        filterNames = append(filterNames, "on_call_policy_execution_repeat_count = "+data.OnCallPolicyExecutionRepeatCount.ValueBigFloat().String())
    }
    if !data.ScheduleGapRetryCount.IsNull() && !data.ScheduleGapRetryCount.IsUnknown() {
        filters["scheduleGapRetryCount"] = lookupNumber(data.ScheduleGapRetryCount)
        filterNames = append(filterNames, "schedule_gap_retry_count = "+data.ScheduleGapRetryCount.ValueBigFloat().String())
    }
    if !data.TriggeredByUserId.IsNull() && !data.TriggeredByUserId.IsUnknown() {
        filters["triggeredByUserId"] = data.TriggeredByUserId.ValueString()
        filterNames = append(filterNames, "triggered_by_user_id = "+fmt.Sprintf("%q", data.TriggeredByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the on call duty execution log up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the on call duty execution log up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "onCallDutyPolicyId": true,
        "triggeredByIncidentId": true,
        "triggeredByAlertId": true,
        "triggeredByAlertEpisodeId": true,
        "triggeredByIncidentEpisodeId": true,
        "status": true,
        "statusMessage": true,
        "userNotificationEventType": true,
        "createdByUserId": true,
        "acknowledgedByUserId": true,
        "acknowledgedAt": true,
        "acknowledgedByTeamId": true,
        "lastExecutedEscalationRuleOrder": true,
        "lastExecutedEscalationRuleId": true,
        "onCallPolicyExecutionRepeatCount": true,
        "scheduleGapRetryCount": true,
        "triggeredByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/on-call-duty-policy-execution-log/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read on_call_duty_execution_log, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No on call duty execution log found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read on_call_duty_execution_log: %s", err))
            return
        }
        if wrapper, ok := itemResponse["data"].(map[string]interface{}); ok {
            item = wrapper
        } else {
            item = itemResponse
        }
    }
    if !hasId {
        listBody := map[string]interface{}{
            "query":  filters,
            "select": selectParam,
            // limit 2 is enough to detect ambiguity without paging.
            "limit": 2,
        }
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/on-call-duty-policy-execution-log/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list on_call_duty_execution_log, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list on_call_duty_execution_log: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No on call duty execution log matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one on call duty execution log matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for on_call_duty_execution_log.")
            return
        }
        item = first
    }

    // Update the model with response data
    if obj, ok := item["_id"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Id = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Id = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Id = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Id = types.StringValue(string(jsonBytes))
        } else {
            data.Id = types.StringNull()
        }
    } else if val, ok := item["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    if obj, ok := item["createdAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CreatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CreatedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CreatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedAt = types.StringNull()
        }
    } else if val, ok := item["createdAt"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    } else {
        data.CreatedAt = types.StringNull()
    }
    if obj, ok := item["updatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UpdatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.UpdatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.UpdatedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.UpdatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.UpdatedAt = types.StringNull()
        }
    } else if val, ok := item["updatedAt"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    } else {
        data.UpdatedAt = types.StringNull()
    }
    if obj, ok := item["projectId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ProjectId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ProjectId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ProjectId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ProjectId = types.StringValue(string(jsonBytes))
        } else {
            data.ProjectId = types.StringNull()
        }
    } else if val, ok := item["projectId"].(string); ok {
        data.ProjectId = types.StringValue(val)
    } else {
        data.ProjectId = types.StringNull()
    }
    if obj, ok := item["onCallDutyPolicyId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OnCallDutyPolicyId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.OnCallDutyPolicyId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.OnCallDutyPolicyId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.OnCallDutyPolicyId = types.StringValue(string(jsonBytes))
        } else {
            data.OnCallDutyPolicyId = types.StringNull()
        }
    } else if val, ok := item["onCallDutyPolicyId"].(string); ok {
        data.OnCallDutyPolicyId = types.StringValue(val)
    } else {
        data.OnCallDutyPolicyId = types.StringNull()
    }
    if obj, ok := item["triggeredByIncidentId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TriggeredByIncidentId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TriggeredByIncidentId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TriggeredByIncidentId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TriggeredByIncidentId = types.StringValue(string(jsonBytes))
        } else {
            data.TriggeredByIncidentId = types.StringNull()
        }
    } else if val, ok := item["triggeredByIncidentId"].(string); ok {
        data.TriggeredByIncidentId = types.StringValue(val)
    } else {
        data.TriggeredByIncidentId = types.StringNull()
    }
    if obj, ok := item["triggeredByAlertId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TriggeredByAlertId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TriggeredByAlertId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TriggeredByAlertId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TriggeredByAlertId = types.StringValue(string(jsonBytes))
        } else {
            data.TriggeredByAlertId = types.StringNull()
        }
    } else if val, ok := item["triggeredByAlertId"].(string); ok {
        data.TriggeredByAlertId = types.StringValue(val)
    } else {
        data.TriggeredByAlertId = types.StringNull()
    }
    if obj, ok := item["triggeredByAlertEpisodeId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TriggeredByAlertEpisodeId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TriggeredByAlertEpisodeId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TriggeredByAlertEpisodeId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TriggeredByAlertEpisodeId = types.StringValue(string(jsonBytes))
        } else {
            data.TriggeredByAlertEpisodeId = types.StringNull()
        }
    } else if val, ok := item["triggeredByAlertEpisodeId"].(string); ok {
        data.TriggeredByAlertEpisodeId = types.StringValue(val)
    } else {
        data.TriggeredByAlertEpisodeId = types.StringNull()
    }
    if obj, ok := item["triggeredByIncidentEpisodeId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TriggeredByIncidentEpisodeId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TriggeredByIncidentEpisodeId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TriggeredByIncidentEpisodeId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TriggeredByIncidentEpisodeId = types.StringValue(string(jsonBytes))
        } else {
            data.TriggeredByIncidentEpisodeId = types.StringNull()
        }
    } else if val, ok := item["triggeredByIncidentEpisodeId"].(string); ok {
        data.TriggeredByIncidentEpisodeId = types.StringValue(val)
    } else {
        data.TriggeredByIncidentEpisodeId = types.StringNull()
    }
    if obj, ok := item["status"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Status = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Status = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Status = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Status = types.StringValue(string(jsonBytes))
        } else {
            data.Status = types.StringNull()
        }
    } else if val, ok := item["status"].(string); ok {
        data.Status = types.StringValue(val)
    } else {
        data.Status = types.StringNull()
    }
    if obj, ok := item["statusMessage"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StatusMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.StatusMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.StatusMessage = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.StatusMessage = types.StringValue(string(jsonBytes))
        } else {
            data.StatusMessage = types.StringNull()
        }
    } else if val, ok := item["statusMessage"].(string); ok {
        data.StatusMessage = types.StringValue(val)
    } else {
        data.StatusMessage = types.StringNull()
    }
    if obj, ok := item["userNotificationEventType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UserNotificationEventType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.UserNotificationEventType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.UserNotificationEventType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.UserNotificationEventType = types.StringValue(string(jsonBytes))
        } else {
            data.UserNotificationEventType = types.StringNull()
        }
    } else if val, ok := item["userNotificationEventType"].(string); ok {
        data.UserNotificationEventType = types.StringValue(val)
    } else {
        data.UserNotificationEventType = types.StringNull()
    }
    if obj, ok := item["createdByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if val, ok := item["createdByUserId"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if obj, ok := item["acknowledgedByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AcknowledgedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AcknowledgedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AcknowledgedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AcknowledgedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.AcknowledgedByUserId = types.StringNull()
        }
    } else if val, ok := item["acknowledgedByUserId"].(string); ok {
        data.AcknowledgedByUserId = types.StringValue(val)
    } else {
        data.AcknowledgedByUserId = types.StringNull()
    }
    if obj, ok := item["acknowledgedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AcknowledgedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AcknowledgedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AcknowledgedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AcknowledgedAt = types.StringValue(string(jsonBytes))
        } else {
            data.AcknowledgedAt = types.StringNull()
        }
    } else if val, ok := item["acknowledgedAt"].(string); ok {
        data.AcknowledgedAt = types.StringValue(val)
    } else {
        data.AcknowledgedAt = types.StringNull()
    }
    if obj, ok := item["acknowledgedByTeamId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AcknowledgedByTeamId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AcknowledgedByTeamId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AcknowledgedByTeamId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AcknowledgedByTeamId = types.StringValue(string(jsonBytes))
        } else {
            data.AcknowledgedByTeamId = types.StringNull()
        }
    } else if val, ok := item["acknowledgedByTeamId"].(string); ok {
        data.AcknowledgedByTeamId = types.StringValue(val)
    } else {
        data.AcknowledgedByTeamId = types.StringNull()
    }
    if val, ok := item["lastExecutedEscalationRuleOrder"].(float64); ok {
        data.LastExecutedEscalationRuleOrder = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["lastExecutedEscalationRuleOrder"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.LastExecutedEscalationRuleOrder = types.NumberValue(big.NewFloat(val))
        } else {
            data.LastExecutedEscalationRuleOrder = types.NumberNull()
        }
    } else {
        data.LastExecutedEscalationRuleOrder = types.NumberNull()
    }
    if obj, ok := item["lastExecutedEscalationRuleId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastExecutedEscalationRuleId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastExecutedEscalationRuleId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastExecutedEscalationRuleId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastExecutedEscalationRuleId = types.StringValue(string(jsonBytes))
        } else {
            data.LastExecutedEscalationRuleId = types.StringNull()
        }
    } else if val, ok := item["lastExecutedEscalationRuleId"].(string); ok {
        data.LastExecutedEscalationRuleId = types.StringValue(val)
    } else {
        data.LastExecutedEscalationRuleId = types.StringNull()
    }
    if val, ok := item["onCallPolicyExecutionRepeatCount"].(float64); ok {
        data.OnCallPolicyExecutionRepeatCount = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["onCallPolicyExecutionRepeatCount"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.OnCallPolicyExecutionRepeatCount = types.NumberValue(big.NewFloat(val))
        } else {
            data.OnCallPolicyExecutionRepeatCount = types.NumberNull()
        }
    } else {
        data.OnCallPolicyExecutionRepeatCount = types.NumberNull()
    }
    if val, ok := item["scheduleGapRetryCount"].(float64); ok {
        data.ScheduleGapRetryCount = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["scheduleGapRetryCount"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.ScheduleGapRetryCount = types.NumberValue(big.NewFloat(val))
        } else {
            data.ScheduleGapRetryCount = types.NumberNull()
        }
    } else {
        data.ScheduleGapRetryCount = types.NumberNull()
    }
    if obj, ok := item["triggeredByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TriggeredByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TriggeredByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TriggeredByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TriggeredByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.TriggeredByUserId = types.StringNull()
        }
    } else if val, ok := item["triggeredByUserId"].(string); ok {
        data.TriggeredByUserId = types.StringValue(val)
    } else {
        data.TriggeredByUserId = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
