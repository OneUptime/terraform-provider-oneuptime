package provider

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &WorkspaceNotificationLogDataSource{}

func NewWorkspaceNotificationLogDataSource() datasource.DataSource {
    return &WorkspaceNotificationLogDataSource{}
}

// WorkspaceNotificationLogDataSource defines the data source implementation.
type WorkspaceNotificationLogDataSource struct {
    client *Client
}

// WorkspaceNotificationLogDataSourceModel describes the data source data model.
type WorkspaceNotificationLogDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
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
    MonitorId types.String `tfsdk:"monitor_id"`
    AlertEpisodeId types.String `tfsdk:"alert_episode_id"`
    IncidentEpisodeId types.String `tfsdk:"incident_episode_id"`
    ScheduledMaintenanceId types.String `tfsdk:"scheduled_maintenance_id"`
    StatusPageId types.String `tfsdk:"status_page_id"`
    StatusPageAnnouncementId types.String `tfsdk:"status_page_announcement_id"`
    OnCallDutyPolicyId types.String `tfsdk:"on_call_duty_policy_id"`
    OnCallDutyPolicyEscalationRuleId types.String `tfsdk:"on_call_duty_policy_escalation_rule_id"`
    OnCallDutyPolicyScheduleId types.String `tfsdk:"on_call_duty_policy_schedule_id"`
    TeamId types.String `tfsdk:"team_id"`
}

func (d *WorkspaceNotificationLogDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_workspace_notification_log"
}

func (d *WorkspaceNotificationLogDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Logs of all workspace activities including messages, channel creation, user invitations, and button interactions for Slack and Microsoft Teams. Look up an existing workspace notification log by `id`, or by any of its other arguments (`action_type`, `alert_episode_id`, `alert_id`, ...): each one set must match, and exactly one workspace notification log may match them all.",

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
            "workspace_type": schema.StringAttribute{
                MarkdownDescription: "Type of Workspace - Slack, Microsoft Teams.",
                Optional: true,
                Computed: true,
            },
            "channel_id": schema.StringAttribute{
                MarkdownDescription: "Channel ID where the message was sent.",
                Optional: true,
                Computed: true,
            },
            "channel_name": schema.StringAttribute{
                MarkdownDescription: "Channel Name where the message was sent.",
                Optional: true,
                Computed: true,
            },
            "thread_id": schema.StringAttribute{
                MarkdownDescription: "Thread ID of the message in the channel (if any).",
                Optional: true,
                Computed: true,
            },
            "message": schema.StringAttribute{
                MarkdownDescription: "Content of the message.",
                Optional: true,
                Computed: true,
            },
            "status_message": schema.StringAttribute{
                MarkdownDescription: "Status Message (if any).",
                Optional: true,
                Computed: true,
            },
            "status": schema.StringAttribute{
                MarkdownDescription: "Status of the message.",
                Optional: true,
                Computed: true,
            },
            "action_type": schema.StringAttribute{
                MarkdownDescription: "Type of workspace action performed.",
                Optional: true,
                Computed: true,
            },
            "incident_id": schema.StringAttribute{
                MarkdownDescription: "ID of Incident associated with this message (if any). The ID of a `oneuptime_incident`.",
                Optional: true,
                Computed: true,
            },
            "user_id": schema.StringAttribute{
                MarkdownDescription: "ID of User who initiated this workspace notification (if any). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "alert_id": schema.StringAttribute{
                MarkdownDescription: "ID of Alert associated with this message (if any). The ID of a `oneuptime_alert`.",
                Optional: true,
                Computed: true,
            },
            "monitor_id": schema.StringAttribute{
                MarkdownDescription: "ID of Monitor associated with this message (if any). The ID of a `oneuptime_monitor`.",
                Optional: true,
                Computed: true,
            },
            "alert_episode_id": schema.StringAttribute{
                MarkdownDescription: "ID of Alert Episode associated with this message (if any). The ID of a `oneuptime_alert_episode`.",
                Optional: true,
                Computed: true,
            },
            "incident_episode_id": schema.StringAttribute{
                MarkdownDescription: "ID of Incident Episode associated with this message (if any). The ID of a `oneuptime_incident_episode`.",
                Optional: true,
                Computed: true,
            },
            "scheduled_maintenance_id": schema.StringAttribute{
                MarkdownDescription: "ID of Scheduled Maintenance associated with this message (if any). The ID of a `oneuptime_scheduled_maintenance_event`.",
                Optional: true,
                Computed: true,
            },
            "status_page_id": schema.StringAttribute{
                MarkdownDescription: "ID of Status Page associated with this message (if any). The ID of a `oneuptime_status_page`.",
                Optional: true,
                Computed: true,
            },
            "status_page_announcement_id": schema.StringAttribute{
                MarkdownDescription: "ID of Status Page Announcement associated with this message (if any). The ID of a `oneuptime_status_page_announcement`.",
                Optional: true,
                Computed: true,
            },
            "on_call_duty_policy_id": schema.StringAttribute{
                MarkdownDescription: "ID of On-Call Duty Policy associated with this message (if any). The ID of a `oneuptime_on_call_policy`.",
                Optional: true,
                Computed: true,
            },
            "on_call_duty_policy_escalation_rule_id": schema.StringAttribute{
                MarkdownDescription: "ID of On-Call Duty Policy Escalation Rule associated with this message (if any). The ID of a `oneuptime_escalation_rule`.",
                Optional: true,
                Computed: true,
            },
            "on_call_duty_policy_schedule_id": schema.StringAttribute{
                MarkdownDescription: "ID of On-Call Duty Policy Schedule associated with this message (if any). The ID of a `oneuptime_on_call_policy_schedule`.",
                Optional: true,
                Computed: true,
            },
            "team_id": schema.StringAttribute{
                MarkdownDescription: "ID of Team associated with this message (if any). The ID of a `oneuptime_team`.",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *WorkspaceNotificationLogDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *WorkspaceNotificationLogDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data WorkspaceNotificationLogDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.WorkspaceType.IsNull() && !data.WorkspaceType.IsUnknown() {
        filters["workspaceType"] = data.WorkspaceType.ValueString()
        filterNames = append(filterNames, "workspace_type = "+fmt.Sprintf("%q", data.WorkspaceType.ValueString()))
    }
    if !data.ChannelId.IsNull() && !data.ChannelId.IsUnknown() {
        filters["channelId"] = data.ChannelId.ValueString()
        filterNames = append(filterNames, "channel_id = "+fmt.Sprintf("%q", data.ChannelId.ValueString()))
    }
    if !data.ChannelName.IsNull() && !data.ChannelName.IsUnknown() {
        filters["channelName"] = data.ChannelName.ValueString()
        filterNames = append(filterNames, "channel_name = "+fmt.Sprintf("%q", data.ChannelName.ValueString()))
    }
    if !data.ThreadId.IsNull() && !data.ThreadId.IsUnknown() {
        filters["threadId"] = data.ThreadId.ValueString()
        filterNames = append(filterNames, "thread_id = "+fmt.Sprintf("%q", data.ThreadId.ValueString()))
    }
    if !data.Message.IsNull() && !data.Message.IsUnknown() {
        filters["message"] = data.Message.ValueString()
        filterNames = append(filterNames, "message = "+fmt.Sprintf("%q", data.Message.ValueString()))
    }
    if !data.StatusMessage.IsNull() && !data.StatusMessage.IsUnknown() {
        filters["statusMessage"] = data.StatusMessage.ValueString()
        filterNames = append(filterNames, "status_message = "+fmt.Sprintf("%q", data.StatusMessage.ValueString()))
    }
    if !data.Status.IsNull() && !data.Status.IsUnknown() {
        filters["status"] = data.Status.ValueString()
        filterNames = append(filterNames, "status = "+fmt.Sprintf("%q", data.Status.ValueString()))
    }
    if !data.ActionType.IsNull() && !data.ActionType.IsUnknown() {
        filters["actionType"] = data.ActionType.ValueString()
        filterNames = append(filterNames, "action_type = "+fmt.Sprintf("%q", data.ActionType.ValueString()))
    }
    if !data.IncidentId.IsNull() && !data.IncidentId.IsUnknown() {
        filters["incidentId"] = data.IncidentId.ValueString()
        filterNames = append(filterNames, "incident_id = "+fmt.Sprintf("%q", data.IncidentId.ValueString()))
    }
    if !data.UserId.IsNull() && !data.UserId.IsUnknown() {
        filters["userId"] = data.UserId.ValueString()
        filterNames = append(filterNames, "user_id = "+fmt.Sprintf("%q", data.UserId.ValueString()))
    }
    if !data.AlertId.IsNull() && !data.AlertId.IsUnknown() {
        filters["alertId"] = data.AlertId.ValueString()
        filterNames = append(filterNames, "alert_id = "+fmt.Sprintf("%q", data.AlertId.ValueString()))
    }
    if !data.MonitorId.IsNull() && !data.MonitorId.IsUnknown() {
        filters["monitorId"] = data.MonitorId.ValueString()
        filterNames = append(filterNames, "monitor_id = "+fmt.Sprintf("%q", data.MonitorId.ValueString()))
    }
    if !data.AlertEpisodeId.IsNull() && !data.AlertEpisodeId.IsUnknown() {
        filters["alertEpisodeId"] = data.AlertEpisodeId.ValueString()
        filterNames = append(filterNames, "alert_episode_id = "+fmt.Sprintf("%q", data.AlertEpisodeId.ValueString()))
    }
    if !data.IncidentEpisodeId.IsNull() && !data.IncidentEpisodeId.IsUnknown() {
        filters["incidentEpisodeId"] = data.IncidentEpisodeId.ValueString()
        filterNames = append(filterNames, "incident_episode_id = "+fmt.Sprintf("%q", data.IncidentEpisodeId.ValueString()))
    }
    if !data.ScheduledMaintenanceId.IsNull() && !data.ScheduledMaintenanceId.IsUnknown() {
        filters["scheduledMaintenanceId"] = data.ScheduledMaintenanceId.ValueString()
        filterNames = append(filterNames, "scheduled_maintenance_id = "+fmt.Sprintf("%q", data.ScheduledMaintenanceId.ValueString()))
    }
    if !data.StatusPageId.IsNull() && !data.StatusPageId.IsUnknown() {
        filters["statusPageId"] = data.StatusPageId.ValueString()
        filterNames = append(filterNames, "status_page_id = "+fmt.Sprintf("%q", data.StatusPageId.ValueString()))
    }
    if !data.StatusPageAnnouncementId.IsNull() && !data.StatusPageAnnouncementId.IsUnknown() {
        filters["statusPageAnnouncementId"] = data.StatusPageAnnouncementId.ValueString()
        filterNames = append(filterNames, "status_page_announcement_id = "+fmt.Sprintf("%q", data.StatusPageAnnouncementId.ValueString()))
    }
    if !data.OnCallDutyPolicyId.IsNull() && !data.OnCallDutyPolicyId.IsUnknown() {
        filters["onCallDutyPolicyId"] = data.OnCallDutyPolicyId.ValueString()
        filterNames = append(filterNames, "on_call_duty_policy_id = "+fmt.Sprintf("%q", data.OnCallDutyPolicyId.ValueString()))
    }
    if !data.OnCallDutyPolicyEscalationRuleId.IsNull() && !data.OnCallDutyPolicyEscalationRuleId.IsUnknown() {
        filters["onCallDutyPolicyEscalationRuleId"] = data.OnCallDutyPolicyEscalationRuleId.ValueString()
        filterNames = append(filterNames, "on_call_duty_policy_escalation_rule_id = "+fmt.Sprintf("%q", data.OnCallDutyPolicyEscalationRuleId.ValueString()))
    }
    if !data.OnCallDutyPolicyScheduleId.IsNull() && !data.OnCallDutyPolicyScheduleId.IsUnknown() {
        filters["onCallDutyPolicyScheduleId"] = data.OnCallDutyPolicyScheduleId.ValueString()
        filterNames = append(filterNames, "on_call_duty_policy_schedule_id = "+fmt.Sprintf("%q", data.OnCallDutyPolicyScheduleId.ValueString()))
    }
    if !data.TeamId.IsNull() && !data.TeamId.IsUnknown() {
        filters["teamId"] = data.TeamId.ValueString()
        filterNames = append(filterNames, "team_id = "+fmt.Sprintf("%q", data.TeamId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the workspace notification log up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the workspace notification log up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "workspaceType": true,
        "channelId": true,
        "channelName": true,
        "threadId": true,
        "message": true,
        "statusMessage": true,
        "status": true,
        "actionType": true,
        "incidentId": true,
        "userId": true,
        "alertId": true,
        "monitorId": true,
        "alertEpisodeId": true,
        "incidentEpisodeId": true,
        "scheduledMaintenanceId": true,
        "statusPageId": true,
        "statusPageAnnouncementId": true,
        "onCallDutyPolicyId": true,
        "onCallDutyPolicyEscalationRuleId": true,
        "onCallDutyPolicyScheduleId": true,
        "teamId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/workspace-notification-log/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read workspace_notification_log, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No workspace notification log found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read workspace_notification_log: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/workspace-notification-log/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list workspace_notification_log, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list workspace_notification_log: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No workspace notification log matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one workspace notification log matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for workspace_notification_log.")
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
    if obj, ok := item["workspaceType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.WorkspaceType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.WorkspaceType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.WorkspaceType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.WorkspaceType = types.StringValue(string(jsonBytes))
        } else {
            data.WorkspaceType = types.StringNull()
        }
    } else if val, ok := item["workspaceType"].(string); ok {
        data.WorkspaceType = types.StringValue(val)
    } else {
        data.WorkspaceType = types.StringNull()
    }
    if obj, ok := item["channelId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ChannelId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ChannelId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ChannelId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ChannelId = types.StringValue(string(jsonBytes))
        } else {
            data.ChannelId = types.StringNull()
        }
    } else if val, ok := item["channelId"].(string); ok {
        data.ChannelId = types.StringValue(val)
    } else {
        data.ChannelId = types.StringNull()
    }
    if obj, ok := item["channelName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ChannelName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ChannelName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ChannelName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ChannelName = types.StringValue(string(jsonBytes))
        } else {
            data.ChannelName = types.StringNull()
        }
    } else if val, ok := item["channelName"].(string); ok {
        data.ChannelName = types.StringValue(val)
    } else {
        data.ChannelName = types.StringNull()
    }
    if obj, ok := item["threadId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ThreadId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ThreadId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ThreadId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ThreadId = types.StringValue(string(jsonBytes))
        } else {
            data.ThreadId = types.StringNull()
        }
    } else if val, ok := item["threadId"].(string); ok {
        data.ThreadId = types.StringValue(val)
    } else {
        data.ThreadId = types.StringNull()
    }
    if obj, ok := item["message"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Message = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Message = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Message = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Message = types.StringValue(string(jsonBytes))
        } else {
            data.Message = types.StringNull()
        }
    } else if val, ok := item["message"].(string); ok {
        data.Message = types.StringValue(val)
    } else {
        data.Message = types.StringNull()
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
    if obj, ok := item["actionType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ActionType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ActionType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ActionType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ActionType = types.StringValue(string(jsonBytes))
        } else {
            data.ActionType = types.StringNull()
        }
    } else if val, ok := item["actionType"].(string); ok {
        data.ActionType = types.StringValue(val)
    } else {
        data.ActionType = types.StringNull()
    }
    if obj, ok := item["incidentId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncidentId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.IncidentId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.IncidentId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.IncidentId = types.StringValue(string(jsonBytes))
        } else {
            data.IncidentId = types.StringNull()
        }
    } else if val, ok := item["incidentId"].(string); ok {
        data.IncidentId = types.StringValue(val)
    } else {
        data.IncidentId = types.StringNull()
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
    if obj, ok := item["alertId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AlertId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AlertId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AlertId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AlertId = types.StringValue(string(jsonBytes))
        } else {
            data.AlertId = types.StringNull()
        }
    } else if val, ok := item["alertId"].(string); ok {
        data.AlertId = types.StringValue(val)
    } else {
        data.AlertId = types.StringNull()
    }
    if obj, ok := item["monitorId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MonitorId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.MonitorId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.MonitorId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.MonitorId = types.StringValue(string(jsonBytes))
        } else {
            data.MonitorId = types.StringNull()
        }
    } else if val, ok := item["monitorId"].(string); ok {
        data.MonitorId = types.StringValue(val)
    } else {
        data.MonitorId = types.StringNull()
    }
    if obj, ok := item["alertEpisodeId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AlertEpisodeId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AlertEpisodeId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AlertEpisodeId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AlertEpisodeId = types.StringValue(string(jsonBytes))
        } else {
            data.AlertEpisodeId = types.StringNull()
        }
    } else if val, ok := item["alertEpisodeId"].(string); ok {
        data.AlertEpisodeId = types.StringValue(val)
    } else {
        data.AlertEpisodeId = types.StringNull()
    }
    if obj, ok := item["incidentEpisodeId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncidentEpisodeId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.IncidentEpisodeId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.IncidentEpisodeId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.IncidentEpisodeId = types.StringValue(string(jsonBytes))
        } else {
            data.IncidentEpisodeId = types.StringNull()
        }
    } else if val, ok := item["incidentEpisodeId"].(string); ok {
        data.IncidentEpisodeId = types.StringValue(val)
    } else {
        data.IncidentEpisodeId = types.StringNull()
    }
    if obj, ok := item["scheduledMaintenanceId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ScheduledMaintenanceId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ScheduledMaintenanceId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ScheduledMaintenanceId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ScheduledMaintenanceId = types.StringValue(string(jsonBytes))
        } else {
            data.ScheduledMaintenanceId = types.StringNull()
        }
    } else if val, ok := item["scheduledMaintenanceId"].(string); ok {
        data.ScheduledMaintenanceId = types.StringValue(val)
    } else {
        data.ScheduledMaintenanceId = types.StringNull()
    }
    if obj, ok := item["statusPageId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StatusPageId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.StatusPageId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.StatusPageId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.StatusPageId = types.StringValue(string(jsonBytes))
        } else {
            data.StatusPageId = types.StringNull()
        }
    } else if val, ok := item["statusPageId"].(string); ok {
        data.StatusPageId = types.StringValue(val)
    } else {
        data.StatusPageId = types.StringNull()
    }
    if obj, ok := item["statusPageAnnouncementId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StatusPageAnnouncementId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.StatusPageAnnouncementId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.StatusPageAnnouncementId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.StatusPageAnnouncementId = types.StringValue(string(jsonBytes))
        } else {
            data.StatusPageAnnouncementId = types.StringNull()
        }
    } else if val, ok := item["statusPageAnnouncementId"].(string); ok {
        data.StatusPageAnnouncementId = types.StringValue(val)
    } else {
        data.StatusPageAnnouncementId = types.StringNull()
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
    if obj, ok := item["onCallDutyPolicyScheduleId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OnCallDutyPolicyScheduleId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.OnCallDutyPolicyScheduleId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.OnCallDutyPolicyScheduleId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.OnCallDutyPolicyScheduleId = types.StringValue(string(jsonBytes))
        } else {
            data.OnCallDutyPolicyScheduleId = types.StringNull()
        }
    } else if val, ok := item["onCallDutyPolicyScheduleId"].(string); ok {
        data.OnCallDutyPolicyScheduleId = types.StringValue(val)
    } else {
        data.OnCallDutyPolicyScheduleId = types.StringNull()
    }
    if obj, ok := item["teamId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TeamId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TeamId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TeamId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TeamId = types.StringValue(string(jsonBytes))
        } else {
            data.TeamId = types.StringNull()
        }
    } else if val, ok := item["teamId"].(string); ok {
        data.TeamId = types.StringValue(val)
    } else {
        data.TeamId = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
