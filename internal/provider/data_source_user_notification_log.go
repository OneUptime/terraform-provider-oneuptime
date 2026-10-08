package provider

import (
    "context"
    "encoding/json"
    "fmt"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &UserNotificationLogDataSource{}

func NewUserNotificationLogDataSource() datasource.DataSource {
    return &UserNotificationLogDataSource{}
}

// UserNotificationLogDataSource defines the data source implementation.
type UserNotificationLogDataSource struct {
    client *Client
}

// UserNotificationLogDataSourceModel describes the data source data model.
type UserNotificationLogDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    UserId types.String `tfsdk:"user_id"`
    UserBelongsToTeamId types.String `tfsdk:"user_belongs_to_team_id"`
    ProjectId types.String `tfsdk:"project_id"`
    OnCallDutyPolicyId types.String `tfsdk:"on_call_duty_policy_id"`
    OnCallDutyPolicyExecutionLogId types.String `tfsdk:"on_call_duty_policy_execution_log_id"`
    OnCallDutyPolicyEscalationRuleId types.String `tfsdk:"on_call_duty_policy_escalation_rule_id"`
    TriggeredByIncidentId types.String `tfsdk:"triggered_by_incident_id"`
    TriggeredByAlertId types.String `tfsdk:"triggered_by_alert_id"`
    TriggeredByAlertEpisodeId types.String `tfsdk:"triggered_by_alert_episode_id"`
    TriggeredByIncidentEpisodeId types.String `tfsdk:"triggered_by_incident_episode_id"`
    Status types.String `tfsdk:"status"`
    UserNotificationEventType types.String `tfsdk:"user_notification_event_type"`
    OnCallDutyPolicyExecutionLogTimelineId types.String `tfsdk:"on_call_duty_policy_execution_log_timeline_id"`
    StatusMessage types.String `tfsdk:"status_message"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    AcknowledgedByUserId types.String `tfsdk:"acknowledged_by_user_id"`
    AcknowledgedAt types.String `tfsdk:"acknowledged_at"`
    OnCallDutyScheduleId types.String `tfsdk:"on_call_duty_schedule_id"`
    OverridedByUserId types.String `tfsdk:"overrided_by_user_id"`
}

func (d *UserNotificationLogDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_user_notification_log"
}

func (d *UserNotificationLogDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Log events for user notifications Look up an existing user notification log by `id`, or by any of its other arguments (`acknowledged_by_user_id`, `created_by_user_id`, `on_call_duty_policy_escalation_rule_id`, ...): each one set must match, and exactly one user notification log may match them all.",

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
            "user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who this log belongs to. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "user_belongs_to_team_id": schema.StringAttribute{
                MarkdownDescription: "Which team did the user belong to when the alert was sent? The ID of a `oneuptime_team`.",
                Optional: true,
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
            "on_call_duty_policy_execution_log_id": schema.StringAttribute{
                MarkdownDescription: "ID of your On-Call Policy execution log which belongs to this log event. The ID of a `oneuptime_on_call_duty_execution_log`.",
                Optional: true,
                Computed: true,
            },
            "on_call_duty_policy_escalation_rule_id": schema.StringAttribute{
                MarkdownDescription: "ID of your On-Call Policy Escalation Rule which belongs to this log event. The ID of a `oneuptime_escalation_rule`.",
                Optional: true,
                Computed: true,
            },
            "triggered_by_incident_id": schema.StringAttribute{
                MarkdownDescription: "ID of the incident which triggered this on-call escalation policy. The ID of a `oneuptime_incident`.",
                Optional: true,
                Computed: true,
            },
            "triggered_by_alert_id": schema.StringAttribute{
                MarkdownDescription: "ID of the Alert which triggered this on-call escalation policy. The ID of a `oneuptime_alert`.",
                Optional: true,
                Computed: true,
            },
            "triggered_by_alert_episode_id": schema.StringAttribute{
                MarkdownDescription: "ID of the Alert Episode which triggered this on-call escalation policy. The ID of a `oneuptime_alert_episode`.",
                Optional: true,
                Computed: true,
            },
            "triggered_by_incident_episode_id": schema.StringAttribute{
                MarkdownDescription: "ID of the Incident Episode which triggered this on-call escalation policy. The ID of a `oneuptime_incident_episode`.",
                Optional: true,
                Computed: true,
            },
            "status": schema.StringAttribute{
                MarkdownDescription: "Status of this execution.",
                Optional: true,
                Computed: true,
            },
            "user_notification_event_type": schema.StringAttribute{
                MarkdownDescription: "Notification Event Type of this execution.",
                Optional: true,
                Computed: true,
            },
            "on_call_duty_policy_execution_log_timeline_id": schema.StringAttribute{
                MarkdownDescription: "ID of your On-Call Policy Execution Log where this timeline event belongs. The ID of a `oneuptime_on_call_duty_execution_log_timeline` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "status_message": schema.StringAttribute{
                MarkdownDescription: "Status message of this execution.",
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
                Computed: true,
            },
            "on_call_duty_schedule_id": schema.StringAttribute{
                MarkdownDescription: "Which schedule ID did the user belong to when the alert was sent? The ID of a `oneuptime_on_call_policy_schedule`.",
                Optional: true,
                Computed: true,
            },
            "overrided_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the user this alert would have paged, when a user override sent it to this log's user instead because they were away. Empty when no override applied. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *UserNotificationLogDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *UserNotificationLogDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data UserNotificationLogDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.UserId.IsNull() && !data.UserId.IsUnknown() {
        filters["userId"] = data.UserId.ValueString()
        filterNames = append(filterNames, "user_id = "+fmt.Sprintf("%q", data.UserId.ValueString()))
    }
    if !data.UserBelongsToTeamId.IsNull() && !data.UserBelongsToTeamId.IsUnknown() {
        filters["userBelongsToTeamId"] = data.UserBelongsToTeamId.ValueString()
        filterNames = append(filterNames, "user_belongs_to_team_id = "+fmt.Sprintf("%q", data.UserBelongsToTeamId.ValueString()))
    }
    if !data.OnCallDutyPolicyId.IsNull() && !data.OnCallDutyPolicyId.IsUnknown() {
        filters["onCallDutyPolicyId"] = data.OnCallDutyPolicyId.ValueString()
        filterNames = append(filterNames, "on_call_duty_policy_id = "+fmt.Sprintf("%q", data.OnCallDutyPolicyId.ValueString()))
    }
    if !data.OnCallDutyPolicyExecutionLogId.IsNull() && !data.OnCallDutyPolicyExecutionLogId.IsUnknown() {
        filters["onCallDutyPolicyExecutionLogId"] = data.OnCallDutyPolicyExecutionLogId.ValueString()
        filterNames = append(filterNames, "on_call_duty_policy_execution_log_id = "+fmt.Sprintf("%q", data.OnCallDutyPolicyExecutionLogId.ValueString()))
    }
    if !data.OnCallDutyPolicyEscalationRuleId.IsNull() && !data.OnCallDutyPolicyEscalationRuleId.IsUnknown() {
        filters["onCallDutyPolicyEscalationRuleId"] = data.OnCallDutyPolicyEscalationRuleId.ValueString()
        filterNames = append(filterNames, "on_call_duty_policy_escalation_rule_id = "+fmt.Sprintf("%q", data.OnCallDutyPolicyEscalationRuleId.ValueString()))
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
    if !data.UserNotificationEventType.IsNull() && !data.UserNotificationEventType.IsUnknown() {
        filters["userNotificationEventType"] = data.UserNotificationEventType.ValueString()
        filterNames = append(filterNames, "user_notification_event_type = "+fmt.Sprintf("%q", data.UserNotificationEventType.ValueString()))
    }
    if !data.OnCallDutyPolicyExecutionLogTimelineId.IsNull() && !data.OnCallDutyPolicyExecutionLogTimelineId.IsUnknown() {
        filters["onCallDutyPolicyExecutionLogTimelineId"] = data.OnCallDutyPolicyExecutionLogTimelineId.ValueString()
        filterNames = append(filterNames, "on_call_duty_policy_execution_log_timeline_id = "+fmt.Sprintf("%q", data.OnCallDutyPolicyExecutionLogTimelineId.ValueString()))
    }
    if !data.StatusMessage.IsNull() && !data.StatusMessage.IsUnknown() {
        filters["statusMessage"] = data.StatusMessage.ValueString()
        filterNames = append(filterNames, "status_message = "+fmt.Sprintf("%q", data.StatusMessage.ValueString()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }
    if !data.AcknowledgedByUserId.IsNull() && !data.AcknowledgedByUserId.IsUnknown() {
        filters["acknowledgedByUserId"] = data.AcknowledgedByUserId.ValueString()
        filterNames = append(filterNames, "acknowledged_by_user_id = "+fmt.Sprintf("%q", data.AcknowledgedByUserId.ValueString()))
    }
    if !data.OnCallDutyScheduleId.IsNull() && !data.OnCallDutyScheduleId.IsUnknown() {
        filters["onCallDutyScheduleId"] = data.OnCallDutyScheduleId.ValueString()
        filterNames = append(filterNames, "on_call_duty_schedule_id = "+fmt.Sprintf("%q", data.OnCallDutyScheduleId.ValueString()))
    }
    if !data.OverridedByUserId.IsNull() && !data.OverridedByUserId.IsUnknown() {
        filters["overridedByUserId"] = data.OverridedByUserId.ValueString()
        filterNames = append(filterNames, "overrided_by_user_id = "+fmt.Sprintf("%q", data.OverridedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the user notification log up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the user notification log up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "userId": true,
        "userBelongsToTeamId": true,
        "projectId": true,
        "onCallDutyPolicyId": true,
        "onCallDutyPolicyExecutionLogId": true,
        "onCallDutyPolicyEscalationRuleId": true,
        "triggeredByIncidentId": true,
        "triggeredByAlertId": true,
        "triggeredByAlertEpisodeId": true,
        "triggeredByIncidentEpisodeId": true,
        "status": true,
        "userNotificationEventType": true,
        "onCallDutyPolicyExecutionLogTimelineId": true,
        "statusMessage": true,
        "createdByUserId": true,
        "acknowledgedByUserId": true,
        "acknowledgedAt": true,
        "onCallDutyScheduleId": true,
        "overridedByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        // No get endpoint: find it in the list by id.
        filters["_id"] = data.Id.ValueString()
        filterNames = append(filterNames, fmt.Sprintf("id = %q", data.Id.ValueString()))
    }
    if item == nil {
        listBody := map[string]interface{}{
            "query":  filters,
            "select": selectParam,
            // limit 2 is enough to detect ambiguity without paging.
            "limit": 2,
        }
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/user-notification-log/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list user_notification_log, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list user_notification_log: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No user notification log matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one user notification log matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for user_notification_log.")
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
    if obj, ok := item["userId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.UserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.UserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.UserId = types.StringValue(string(jsonBytes))
        } else {
            data.UserId = types.StringNull()
        }
    } else if val, ok := item["userId"].(string); ok {
        data.UserId = types.StringValue(val)
    } else {
        data.UserId = types.StringNull()
    }
    if obj, ok := item["userBelongsToTeamId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UserBelongsToTeamId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.UserBelongsToTeamId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.UserBelongsToTeamId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.UserBelongsToTeamId = types.StringValue(string(jsonBytes))
        } else {
            data.UserBelongsToTeamId = types.StringNull()
        }
    } else if val, ok := item["userBelongsToTeamId"].(string); ok {
        data.UserBelongsToTeamId = types.StringValue(val)
    } else {
        data.UserBelongsToTeamId = types.StringNull()
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
    if obj, ok := item["onCallDutyPolicyExecutionLogId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OnCallDutyPolicyExecutionLogId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.OnCallDutyPolicyExecutionLogId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.OnCallDutyPolicyExecutionLogId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.OnCallDutyPolicyExecutionLogId = types.StringValue(string(jsonBytes))
        } else {
            data.OnCallDutyPolicyExecutionLogId = types.StringNull()
        }
    } else if val, ok := item["onCallDutyPolicyExecutionLogId"].(string); ok {
        data.OnCallDutyPolicyExecutionLogId = types.StringValue(val)
    } else {
        data.OnCallDutyPolicyExecutionLogId = types.StringNull()
    }
    if obj, ok := item["onCallDutyPolicyEscalationRuleId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OnCallDutyPolicyEscalationRuleId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.OnCallDutyPolicyEscalationRuleId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.OnCallDutyPolicyEscalationRuleId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.OnCallDutyPolicyEscalationRuleId = types.StringValue(string(jsonBytes))
        } else {
            data.OnCallDutyPolicyEscalationRuleId = types.StringNull()
        }
    } else if val, ok := item["onCallDutyPolicyEscalationRuleId"].(string); ok {
        data.OnCallDutyPolicyEscalationRuleId = types.StringValue(val)
    } else {
        data.OnCallDutyPolicyEscalationRuleId = types.StringNull()
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
    if obj, ok := item["onCallDutyPolicyExecutionLogTimelineId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OnCallDutyPolicyExecutionLogTimelineId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.OnCallDutyPolicyExecutionLogTimelineId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.OnCallDutyPolicyExecutionLogTimelineId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.OnCallDutyPolicyExecutionLogTimelineId = types.StringValue(string(jsonBytes))
        } else {
            data.OnCallDutyPolicyExecutionLogTimelineId = types.StringNull()
        }
    } else if val, ok := item["onCallDutyPolicyExecutionLogTimelineId"].(string); ok {
        data.OnCallDutyPolicyExecutionLogTimelineId = types.StringValue(val)
    } else {
        data.OnCallDutyPolicyExecutionLogTimelineId = types.StringNull()
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
    if obj, ok := item["onCallDutyScheduleId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OnCallDutyScheduleId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.OnCallDutyScheduleId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.OnCallDutyScheduleId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.OnCallDutyScheduleId = types.StringValue(string(jsonBytes))
        } else {
            data.OnCallDutyScheduleId = types.StringNull()
        }
    } else if val, ok := item["onCallDutyScheduleId"].(string); ok {
        data.OnCallDutyScheduleId = types.StringValue(val)
    } else {
        data.OnCallDutyScheduleId = types.StringNull()
    }
    if obj, ok := item["overridedByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OverridedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.OverridedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.OverridedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.OverridedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.OverridedByUserId = types.StringNull()
        }
    } else if val, ok := item["overridedByUserId"].(string); ok {
        data.OverridedByUserId = types.StringValue(val)
    } else {
        data.OverridedByUserId = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
