package provider

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "math/big"
    "github.com/hashicorp/terraform-plugin-framework/attr"
    "sort"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &AlertEpisodeDataSource{}

func NewAlertEpisodeDataSource() datasource.DataSource {
    return &AlertEpisodeDataSource{}
}

// AlertEpisodeDataSource defines the data source implementation.
type AlertEpisodeDataSource struct {
    client *Client
}

// AlertEpisodeDataSourceModel describes the data source data model.
type AlertEpisodeDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    Title types.String `tfsdk:"title"`
    Description types.String `tfsdk:"description"`
    EpisodeNumber types.Number `tfsdk:"episode_number"`
    EpisodeNumberWithPrefix types.String `tfsdk:"episode_number_with_prefix"`
    CurrentAlertStateId types.String `tfsdk:"current_alert_state_id"`
    AlertSeverityId types.String `tfsdk:"alert_severity_id"`
    RootCause types.String `tfsdk:"root_cause"`
    LastAlertAddedAt types.String `tfsdk:"last_alert_added_at"`
    ResolvedAt types.String `tfsdk:"resolved_at"`
    AllAlertsResolvedAt types.String `tfsdk:"all_alerts_resolved_at"`
    AssignedToUserId types.String `tfsdk:"assigned_to_user_id"`
    AssignedToTeamId types.String `tfsdk:"assigned_to_team_id"`
    AlertGroupingRuleId types.String `tfsdk:"alert_grouping_rule_id"`
    OnCallDutyPolicies types.Set `tfsdk:"on_call_duty_policies"`
    IsOnCallPolicyExecuted types.Bool `tfsdk:"is_on_call_policy_executed"`
    AlertCount types.Number `tfsdk:"alert_count"`
    TitleTemplate types.String `tfsdk:"title_template"`
    DescriptionTemplate types.String `tfsdk:"description_template"`
    IsManuallyCreated types.Bool `tfsdk:"is_manually_created"`
    Labels types.Set `tfsdk:"labels"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    IsOwnerNotifiedOfEpisodeCreation types.Bool `tfsdk:"is_owner_notified_of_episode_creation"`
    GroupingKey types.String `tfsdk:"grouping_key"`
    RemediationNotes types.String `tfsdk:"remediation_notes"`
    PostUpdatesToWorkspaceChannels types.String `tfsdk:"post_updates_to_workspace_channels"`
    IsPrivate types.Bool `tfsdk:"is_private"`
}

func (d *AlertEpisodeDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_alert_episode"
}

func (d *AlertEpisodeDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Manage alert episodes (groups of related alerts) for your project Look up an existing alert episode by `id`, or by any of its other arguments (`alert_count`, `alert_grouping_rule_id`, `alert_severity_id`, ...): each one set must match, and exactly one alert episode may match them all.",

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
            "title": schema.StringAttribute{
                MarkdownDescription: "Title of this alert episode.",
                Optional: true,
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Description of this alert episode. This is in markdown format.",
                Optional: true,
                Computed: true,
            },
            "episode_number": schema.NumberAttribute{
                MarkdownDescription: "Auto-incrementing episode number per project.",
                Optional: true,
                Computed: true,
            },
            "episode_number_with_prefix": schema.StringAttribute{
                MarkdownDescription: "Episode number with prefix (e.g., 'AE-42' or '#42').",
                Optional: true,
                Computed: true,
            },
            "current_alert_state_id": schema.StringAttribute{
                MarkdownDescription: "Current Alert State ID. The ID of a `oneuptime_alert_state`.",
                Optional: true,
                Computed: true,
            },
            "alert_severity_id": schema.StringAttribute{
                MarkdownDescription: "Alert Severity ID. The ID of a `oneuptime_alert_severity`.",
                Optional: true,
                Computed: true,
            },
            "root_cause": schema.StringAttribute{
                MarkdownDescription: "User-documented root cause of this episode.",
                Optional: true,
                Computed: true,
            },
            "last_alert_added_at": schema.StringAttribute{
                MarkdownDescription: "When the last alert was added to this episode.",
                Computed: true,
            },
            "resolved_at": schema.StringAttribute{
                MarkdownDescription: "When this episode was resolved.",
                Computed: true,
            },
            "all_alerts_resolved_at": schema.StringAttribute{
                MarkdownDescription: "When all alerts in this episode were first detected as resolved. Used for resolve delay calculation.",
                Computed: true,
            },
            "assigned_to_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who is assigned to this episode. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "assigned_to_team_id": schema.StringAttribute{
                MarkdownDescription: "Team ID that is assigned to this episode. The ID of a `oneuptime_team`.",
                Optional: true,
                Computed: true,
            },
            "alert_grouping_rule_id": schema.StringAttribute{
                MarkdownDescription: "Alert Grouping Rule ID that created this episode. The ID of a `oneuptime_alert_grouping_rule`.",
                Optional: true,
                Computed: true,
            },
            "on_call_duty_policies": schema.SetAttribute{
                MarkdownDescription: "List of on-call duty policies to execute for this episode. IDs of `oneuptime_on_call_policy` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "is_on_call_policy_executed": schema.BoolAttribute{
                MarkdownDescription: "Whether the on-call policy has been executed for this episode.",
                Optional: true,
                Computed: true,
            },
            "alert_count": schema.NumberAttribute{
                MarkdownDescription: "Denormalized count of alerts in this episode.",
                Optional: true,
                Computed: true,
            },
            "title_template": schema.StringAttribute{
                MarkdownDescription: "Template used to generate the episode title. Stored for dynamic variable updates.",
                Optional: true,
                Computed: true,
            },
            "description_template": schema.StringAttribute{
                MarkdownDescription: "Template used to generate the episode description. Stored for dynamic variable updates.",
                Optional: true,
                Computed: true,
            },
            "is_manually_created": schema.BoolAttribute{
                MarkdownDescription: "Whether this episode was manually created vs auto-created by a rule.",
                Optional: true,
                Computed: true,
            },
            "labels": schema.SetAttribute{
                MarkdownDescription: "Relation to Labels Array where this object is categorized in. IDs of `oneuptime_label` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "is_owner_notified_of_episode_creation": schema.BoolAttribute{
                MarkdownDescription: "Are owners notified when this episode is created?",
                Optional: true,
                Computed: true,
            },
            "grouping_key": schema.StringAttribute{
                MarkdownDescription: "Key used for grouping alerts into this episode. Generated from groupByFields of the matching rule. When a private alert opened the episode, its title is in the key only as a keyed hash.",
                Optional: true,
                Computed: true,
            },
            "remediation_notes": schema.StringAttribute{
                MarkdownDescription: "User-documented remediation steps and notes for this episode.",
                Optional: true,
                Computed: true,
            },
            "post_updates_to_workspace_channels": schema.StringAttribute{
                MarkdownDescription: "Workspace channels to post episode updates to (e.g., Slack, Microsoft Teams). A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "is_private": schema.BoolAttribute{
                MarkdownDescription: "If true, this alert episode is only visible to its owners (users in 'owner users' and members of 'owner teams'), project admins, and project owners.",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *AlertEpisodeDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AlertEpisodeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data AlertEpisodeDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.Title.IsNull() && !data.Title.IsUnknown() {
        filters["title"] = data.Title.ValueString()
        filterNames = append(filterNames, "title = "+fmt.Sprintf("%q", data.Title.ValueString()))
    }
    if !data.Description.IsNull() && !data.Description.IsUnknown() {
        filters["description"] = data.Description.ValueString()
        filterNames = append(filterNames, "description = "+fmt.Sprintf("%q", data.Description.ValueString()))
    }
    if !data.EpisodeNumber.IsNull() && !data.EpisodeNumber.IsUnknown() {
        filters["episodeNumber"] = lookupNumber(data.EpisodeNumber)
        filterNames = append(filterNames, "episode_number = "+data.EpisodeNumber.ValueBigFloat().String())
    }
    if !data.EpisodeNumberWithPrefix.IsNull() && !data.EpisodeNumberWithPrefix.IsUnknown() {
        filters["episodeNumberWithPrefix"] = data.EpisodeNumberWithPrefix.ValueString()
        filterNames = append(filterNames, "episode_number_with_prefix = "+fmt.Sprintf("%q", data.EpisodeNumberWithPrefix.ValueString()))
    }
    if !data.CurrentAlertStateId.IsNull() && !data.CurrentAlertStateId.IsUnknown() {
        filters["currentAlertStateId"] = data.CurrentAlertStateId.ValueString()
        filterNames = append(filterNames, "current_alert_state_id = "+fmt.Sprintf("%q", data.CurrentAlertStateId.ValueString()))
    }
    if !data.AlertSeverityId.IsNull() && !data.AlertSeverityId.IsUnknown() {
        filters["alertSeverityId"] = data.AlertSeverityId.ValueString()
        filterNames = append(filterNames, "alert_severity_id = "+fmt.Sprintf("%q", data.AlertSeverityId.ValueString()))
    }
    if !data.RootCause.IsNull() && !data.RootCause.IsUnknown() {
        filters["rootCause"] = data.RootCause.ValueString()
        filterNames = append(filterNames, "root_cause = "+fmt.Sprintf("%q", data.RootCause.ValueString()))
    }
    if !data.AssignedToUserId.IsNull() && !data.AssignedToUserId.IsUnknown() {
        filters["assignedToUserId"] = data.AssignedToUserId.ValueString()
        filterNames = append(filterNames, "assigned_to_user_id = "+fmt.Sprintf("%q", data.AssignedToUserId.ValueString()))
    }
    if !data.AssignedToTeamId.IsNull() && !data.AssignedToTeamId.IsUnknown() {
        filters["assignedToTeamId"] = data.AssignedToTeamId.ValueString()
        filterNames = append(filterNames, "assigned_to_team_id = "+fmt.Sprintf("%q", data.AssignedToTeamId.ValueString()))
    }
    if !data.AlertGroupingRuleId.IsNull() && !data.AlertGroupingRuleId.IsUnknown() {
        filters["alertGroupingRuleId"] = data.AlertGroupingRuleId.ValueString()
        filterNames = append(filterNames, "alert_grouping_rule_id = "+fmt.Sprintf("%q", data.AlertGroupingRuleId.ValueString()))
    }
    if !data.IsOnCallPolicyExecuted.IsNull() && !data.IsOnCallPolicyExecuted.IsUnknown() {
        filters["isOnCallPolicyExecuted"] = data.IsOnCallPolicyExecuted.ValueBool()
        filterNames = append(filterNames, "is_on_call_policy_executed = "+fmt.Sprintf("%t", data.IsOnCallPolicyExecuted.ValueBool()))
    }
    if !data.AlertCount.IsNull() && !data.AlertCount.IsUnknown() {
        filters["alertCount"] = lookupNumber(data.AlertCount)
        filterNames = append(filterNames, "alert_count = "+data.AlertCount.ValueBigFloat().String())
    }
    if !data.TitleTemplate.IsNull() && !data.TitleTemplate.IsUnknown() {
        filters["titleTemplate"] = data.TitleTemplate.ValueString()
        filterNames = append(filterNames, "title_template = "+fmt.Sprintf("%q", data.TitleTemplate.ValueString()))
    }
    if !data.DescriptionTemplate.IsNull() && !data.DescriptionTemplate.IsUnknown() {
        filters["descriptionTemplate"] = data.DescriptionTemplate.ValueString()
        filterNames = append(filterNames, "description_template = "+fmt.Sprintf("%q", data.DescriptionTemplate.ValueString()))
    }
    if !data.IsManuallyCreated.IsNull() && !data.IsManuallyCreated.IsUnknown() {
        filters["isManuallyCreated"] = data.IsManuallyCreated.ValueBool()
        filterNames = append(filterNames, "is_manually_created = "+fmt.Sprintf("%t", data.IsManuallyCreated.ValueBool()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }
    if !data.IsOwnerNotifiedOfEpisodeCreation.IsNull() && !data.IsOwnerNotifiedOfEpisodeCreation.IsUnknown() {
        filters["isOwnerNotifiedOfEpisodeCreation"] = data.IsOwnerNotifiedOfEpisodeCreation.ValueBool()
        filterNames = append(filterNames, "is_owner_notified_of_episode_creation = "+fmt.Sprintf("%t", data.IsOwnerNotifiedOfEpisodeCreation.ValueBool()))
    }
    if !data.GroupingKey.IsNull() && !data.GroupingKey.IsUnknown() {
        filters["groupingKey"] = data.GroupingKey.ValueString()
        filterNames = append(filterNames, "grouping_key = "+fmt.Sprintf("%q", data.GroupingKey.ValueString()))
    }
    if !data.RemediationNotes.IsNull() && !data.RemediationNotes.IsUnknown() {
        filters["remediationNotes"] = data.RemediationNotes.ValueString()
        filterNames = append(filterNames, "remediation_notes = "+fmt.Sprintf("%q", data.RemediationNotes.ValueString()))
    }
    if !data.IsPrivate.IsNull() && !data.IsPrivate.IsUnknown() {
        filters["isPrivate"] = data.IsPrivate.ValueBool()
        filterNames = append(filterNames, "is_private = "+fmt.Sprintf("%t", data.IsPrivate.ValueBool()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the alert episode up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the alert episode up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "title": true,
        "description": true,
        "episodeNumber": true,
        "episodeNumberWithPrefix": true,
        "currentAlertStateId": true,
        "alertSeverityId": true,
        "rootCause": true,
        "lastAlertAddedAt": true,
        "resolvedAt": true,
        "allAlertsResolvedAt": true,
        "assignedToUserId": true,
        "assignedToTeamId": true,
        "alertGroupingRuleId": true,
        "onCallDutyPolicies": true,
        "isOnCallPolicyExecuted": true,
        "alertCount": true,
        "titleTemplate": true,
        "descriptionTemplate": true,
        "isManuallyCreated": true,
        "labels": true,
        "createdByUserId": true,
        "isOwnerNotifiedOfEpisodeCreation": true,
        "groupingKey": true,
        "remediationNotes": true,
        "postUpdatesToWorkspaceChannels": true,
        "isPrivate": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/alert-episode/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read alert_episode, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No alert episode found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read alert_episode: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/alert-episode/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list alert_episode, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list alert_episode: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No alert episode matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one alert episode matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for alert_episode.")
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
    if obj, ok := item["title"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Title = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Title = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Title = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Title = types.StringValue(string(jsonBytes))
        } else {
            data.Title = types.StringNull()
        }
    } else if val, ok := item["title"].(string); ok {
        data.Title = types.StringValue(val)
    } else {
        data.Title = types.StringNull()
    }
    if obj, ok := item["description"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Description = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Description = types.StringValue(string(jsonBytes))
        } else {
            data.Description = types.StringNull()
        }
    } else if val, ok := item["description"].(string); ok {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if val, ok := item["episodeNumber"].(float64); ok {
        data.EpisodeNumber = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["episodeNumber"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.EpisodeNumber = types.NumberValue(big.NewFloat(val))
        } else {
            data.EpisodeNumber = types.NumberNull()
        }
    } else {
        data.EpisodeNumber = types.NumberNull()
    }
    if obj, ok := item["episodeNumberWithPrefix"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EpisodeNumberWithPrefix = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.EpisodeNumberWithPrefix = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.EpisodeNumberWithPrefix = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.EpisodeNumberWithPrefix = types.StringValue(string(jsonBytes))
        } else {
            data.EpisodeNumberWithPrefix = types.StringNull()
        }
    } else if val, ok := item["episodeNumberWithPrefix"].(string); ok {
        data.EpisodeNumberWithPrefix = types.StringValue(val)
    } else {
        data.EpisodeNumberWithPrefix = types.StringNull()
    }
    if obj, ok := item["currentAlertStateId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CurrentAlertStateId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CurrentAlertStateId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CurrentAlertStateId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CurrentAlertStateId = types.StringValue(string(jsonBytes))
        } else {
            data.CurrentAlertStateId = types.StringNull()
        }
    } else if val, ok := item["currentAlertStateId"].(string); ok {
        data.CurrentAlertStateId = types.StringValue(val)
    } else {
        data.CurrentAlertStateId = types.StringNull()
    }
    if obj, ok := item["alertSeverityId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AlertSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AlertSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AlertSeverityId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AlertSeverityId = types.StringValue(string(jsonBytes))
        } else {
            data.AlertSeverityId = types.StringNull()
        }
    } else if val, ok := item["alertSeverityId"].(string); ok {
        data.AlertSeverityId = types.StringValue(val)
    } else {
        data.AlertSeverityId = types.StringNull()
    }
    if obj, ok := item["rootCause"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RootCause = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RootCause = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RootCause = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RootCause = types.StringValue(string(jsonBytes))
        } else {
            data.RootCause = types.StringNull()
        }
    } else if val, ok := item["rootCause"].(string); ok {
        data.RootCause = types.StringValue(val)
    } else {
        data.RootCause = types.StringNull()
    }
    if obj, ok := item["lastAlertAddedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastAlertAddedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastAlertAddedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastAlertAddedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastAlertAddedAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastAlertAddedAt = types.StringNull()
        }
    } else if val, ok := item["lastAlertAddedAt"].(string); ok {
        data.LastAlertAddedAt = types.StringValue(val)
    } else {
        data.LastAlertAddedAt = types.StringNull()
    }
    if obj, ok := item["resolvedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ResolvedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ResolvedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ResolvedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ResolvedAt = types.StringValue(string(jsonBytes))
        } else {
            data.ResolvedAt = types.StringNull()
        }
    } else if val, ok := item["resolvedAt"].(string); ok {
        data.ResolvedAt = types.StringValue(val)
    } else {
        data.ResolvedAt = types.StringNull()
    }
    if obj, ok := item["allAlertsResolvedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AllAlertsResolvedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AllAlertsResolvedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AllAlertsResolvedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AllAlertsResolvedAt = types.StringValue(string(jsonBytes))
        } else {
            data.AllAlertsResolvedAt = types.StringNull()
        }
    } else if val, ok := item["allAlertsResolvedAt"].(string); ok {
        data.AllAlertsResolvedAt = types.StringValue(val)
    } else {
        data.AllAlertsResolvedAt = types.StringNull()
    }
    if obj, ok := item["assignedToUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AssignedToUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AssignedToUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AssignedToUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AssignedToUserId = types.StringValue(string(jsonBytes))
        } else {
            data.AssignedToUserId = types.StringNull()
        }
    } else if val, ok := item["assignedToUserId"].(string); ok {
        data.AssignedToUserId = types.StringValue(val)
    } else {
        data.AssignedToUserId = types.StringNull()
    }
    if obj, ok := item["assignedToTeamId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AssignedToTeamId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AssignedToTeamId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AssignedToTeamId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AssignedToTeamId = types.StringValue(string(jsonBytes))
        } else {
            data.AssignedToTeamId = types.StringNull()
        }
    } else if val, ok := item["assignedToTeamId"].(string); ok {
        data.AssignedToTeamId = types.StringValue(val)
    } else {
        data.AssignedToTeamId = types.StringNull()
    }
    if obj, ok := item["alertGroupingRuleId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AlertGroupingRuleId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AlertGroupingRuleId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AlertGroupingRuleId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AlertGroupingRuleId = types.StringValue(string(jsonBytes))
        } else {
            data.AlertGroupingRuleId = types.StringNull()
        }
    } else if val, ok := item["alertGroupingRuleId"].(string); ok {
        data.AlertGroupingRuleId = types.StringValue(val)
    } else {
        data.AlertGroupingRuleId = types.StringNull()
    }
    if val, ok := item["onCallDutyPolicies"].([]interface{}); ok {
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if jsonBytes, err := json.Marshal(itemMap); err == nil {
                    setItems = append(setItems, types.StringValue(string(jsonBytes)))
                }
            } else if str, ok := item.(string); ok {
                setItems = append(setItems, types.StringValue(str))
            } else {
                setItems = append(setItems, types.StringValue(fmt.Sprintf("%v", item)))
            }
        }
        sort.Slice(setItems, func(i, j int) bool {
            return setItems[i].(types.String).ValueString() < setItems[j].(types.String).ValueString()
        })
        data.OnCallDutyPolicies = types.SetValueMust(types.StringType, setItems)
    } else {
        data.OnCallDutyPolicies = types.SetNull(types.StringType)
    }
    if val, ok := item["isOnCallPolicyExecuted"].(bool); ok {
        data.IsOnCallPolicyExecuted = types.BoolValue(val)
    } else {
        data.IsOnCallPolicyExecuted = types.BoolNull()
    }
    if val, ok := item["alertCount"].(float64); ok {
        data.AlertCount = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["alertCount"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.AlertCount = types.NumberValue(big.NewFloat(val))
        } else {
            data.AlertCount = types.NumberNull()
        }
    } else {
        data.AlertCount = types.NumberNull()
    }
    if obj, ok := item["titleTemplate"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TitleTemplate = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TitleTemplate = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TitleTemplate = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TitleTemplate = types.StringValue(string(jsonBytes))
        } else {
            data.TitleTemplate = types.StringNull()
        }
    } else if val, ok := item["titleTemplate"].(string); ok {
        data.TitleTemplate = types.StringValue(val)
    } else {
        data.TitleTemplate = types.StringNull()
    }
    if obj, ok := item["descriptionTemplate"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DescriptionTemplate = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DescriptionTemplate = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DescriptionTemplate = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DescriptionTemplate = types.StringValue(string(jsonBytes))
        } else {
            data.DescriptionTemplate = types.StringNull()
        }
    } else if val, ok := item["descriptionTemplate"].(string); ok {
        data.DescriptionTemplate = types.StringValue(val)
    } else {
        data.DescriptionTemplate = types.StringNull()
    }
    if val, ok := item["isManuallyCreated"].(bool); ok {
        data.IsManuallyCreated = types.BoolValue(val)
    } else {
        data.IsManuallyCreated = types.BoolNull()
    }
    if val, ok := item["labels"].([]interface{}); ok {
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if jsonBytes, err := json.Marshal(itemMap); err == nil {
                    setItems = append(setItems, types.StringValue(string(jsonBytes)))
                }
            } else if str, ok := item.(string); ok {
                setItems = append(setItems, types.StringValue(str))
            } else {
                setItems = append(setItems, types.StringValue(fmt.Sprintf("%v", item)))
            }
        }
        sort.Slice(setItems, func(i, j int) bool {
            return setItems[i].(types.String).ValueString() < setItems[j].(types.String).ValueString()
        })
        data.Labels = types.SetValueMust(types.StringType, setItems)
    } else {
        data.Labels = types.SetNull(types.StringType)
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
    if val, ok := item["isOwnerNotifiedOfEpisodeCreation"].(bool); ok {
        data.IsOwnerNotifiedOfEpisodeCreation = types.BoolValue(val)
    } else {
        data.IsOwnerNotifiedOfEpisodeCreation = types.BoolNull()
    }
    if obj, ok := item["groupingKey"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.GroupingKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.GroupingKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.GroupingKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.GroupingKey = types.StringValue(string(jsonBytes))
        } else {
            data.GroupingKey = types.StringNull()
        }
    } else if val, ok := item["groupingKey"].(string); ok {
        data.GroupingKey = types.StringValue(val)
    } else {
        data.GroupingKey = types.StringNull()
    }
    if obj, ok := item["remediationNotes"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RemediationNotes = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RemediationNotes = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RemediationNotes = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RemediationNotes = types.StringValue(string(jsonBytes))
        } else {
            data.RemediationNotes = types.StringNull()
        }
    } else if val, ok := item["remediationNotes"].(string); ok {
        data.RemediationNotes = types.StringValue(val)
    } else {
        data.RemediationNotes = types.StringNull()
    }
    if obj, ok := item["postUpdatesToWorkspaceChannels"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PostUpdatesToWorkspaceChannels = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.PostUpdatesToWorkspaceChannels = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.PostUpdatesToWorkspaceChannels = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.PostUpdatesToWorkspaceChannels = types.StringValue(string(jsonBytes))
        } else {
            data.PostUpdatesToWorkspaceChannels = types.StringNull()
        }
    } else if val, ok := item["postUpdatesToWorkspaceChannels"].(string); ok {
        data.PostUpdatesToWorkspaceChannels = types.StringValue(val)
    } else {
        data.PostUpdatesToWorkspaceChannels = types.StringNull()
    }
    if val, ok := item["isPrivate"].(bool); ok {
        data.IsPrivate = types.BoolValue(val)
    } else {
        data.IsPrivate = types.BoolNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
