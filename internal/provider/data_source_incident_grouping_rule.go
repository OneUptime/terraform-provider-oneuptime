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
var _ datasource.DataSource = &IncidentGroupingRuleDataSource{}

func NewIncidentGroupingRuleDataSource() datasource.DataSource {
    return &IncidentGroupingRuleDataSource{}
}

// IncidentGroupingRuleDataSource defines the data source implementation.
type IncidentGroupingRuleDataSource struct {
    client *Client
}

// IncidentGroupingRuleDataSourceModel describes the data source data model.
type IncidentGroupingRuleDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    Criteria types.String `tfsdk:"criteria"`
    ProjectId types.String `tfsdk:"project_id"`
    Name types.String `tfsdk:"name"`
    Description types.String `tfsdk:"description"`
    Priority types.Number `tfsdk:"priority"`
    IsEnabled types.Bool `tfsdk:"is_enabled"`
    MatchCriteria types.String `tfsdk:"match_criteria"`
    Monitors types.Set `tfsdk:"monitors"`
    IncidentSeverities types.Set `tfsdk:"incident_severities"`
    IncidentLabels types.Set `tfsdk:"incident_labels"`
    MonitorLabels types.Set `tfsdk:"monitor_labels"`
    IncidentTitlePattern types.String `tfsdk:"incident_title_pattern"`
    IncidentDescriptionPattern types.String `tfsdk:"incident_description_pattern"`
    MonitorNamePattern types.String `tfsdk:"monitor_name_pattern"`
    MonitorDescriptionPattern types.String `tfsdk:"monitor_description_pattern"`
    GroupByMonitor types.Bool `tfsdk:"group_by_monitor"`
    GroupBySeverity types.Bool `tfsdk:"group_by_severity"`
    GroupByIncidentTitle types.Bool `tfsdk:"group_by_incident_title"`
    GroupByIncidentLabels types.Bool `tfsdk:"group_by_incident_labels"`
    GroupByMonitorLabels types.Bool `tfsdk:"group_by_monitor_labels"`
    EnableTimeWindow types.Bool `tfsdk:"enable_time_window"`
    TimeWindowMinutes types.Number `tfsdk:"time_window_minutes"`
    GroupByFields types.String `tfsdk:"group_by_fields"`
    EpisodeTitleTemplate types.String `tfsdk:"episode_title_template"`
    EpisodeDescriptionTemplate types.String `tfsdk:"episode_description_template"`
    EnableResolveDelay types.Bool `tfsdk:"enable_resolve_delay"`
    ResolveDelayMinutes types.Number `tfsdk:"resolve_delay_minutes"`
    EnableReopenWindow types.Bool `tfsdk:"enable_reopen_window"`
    ReopenWindowMinutes types.Number `tfsdk:"reopen_window_minutes"`
    EnableInactivityTimeout types.Bool `tfsdk:"enable_inactivity_timeout"`
    InactivityTimeoutMinutes types.Number `tfsdk:"inactivity_timeout_minutes"`
    OnCallDutyPolicies types.Set `tfsdk:"on_call_duty_policies"`
    DefaultAssignToUserId types.String `tfsdk:"default_assign_to_user_id"`
    DefaultAssignToTeamId types.String `tfsdk:"default_assign_to_team_id"`
    EpisodeLabels types.Set `tfsdk:"episode_labels"`
    EpisodeOwnerUsers types.Set `tfsdk:"episode_owner_users"`
    EpisodeOwnerTeams types.Set `tfsdk:"episode_owner_teams"`
    EpisodeMemberRoles types.Set `tfsdk:"episode_member_roles"`
    EpisodeMemberRoleAssignments types.String `tfsdk:"episode_member_role_assignments"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    ShowEpisodeOnStatusPage types.Bool `tfsdk:"show_episode_on_status_page"`
}

func (d *IncidentGroupingRuleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_incident_grouping_rule"
}

func (d *IncidentGroupingRuleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Configure rules for automatically grouping related incidents into episodes Look up an existing incident grouping rule by `id`, or by any of its other arguments (`name`, `created_by_user_id`, `default_assign_to_team_id`, ...): each one set must match, and exactly one incident grouping rule may match them all.",

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
            "criteria": schema.StringAttribute{
                MarkdownDescription: "Versioned conditions that determine whether this rule matches a resource. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Name of this incident grouping rule.",
                Optional: true,
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Description of this incident grouping rule.",
                Optional: true,
                Computed: true,
            },
            "priority": schema.NumberAttribute{
                MarkdownDescription: "Where this rule sits in the list. Rules are evaluated from the top of the list down, lowest number first. A new rule is added to the end of the list. Setting a number another one already has puts it in that place, and the ones in the way move one place along to make room. In the dashboard, drag the rows to reorder them.",
                Optional: true,
                Computed: true,
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether this rule is enabled.",
                Optional: true,
                Computed: true,
            },
            "match_criteria": schema.StringAttribute{
                MarkdownDescription: "JSON object defining the criteria for matching incidents to this rule. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "monitors": schema.SetAttribute{
                MarkdownDescription: "Only group incidents from these monitors. Leave empty to match incidents from any monitor. IDs of `oneuptime_monitor` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "incident_severities": schema.SetAttribute{
                MarkdownDescription: "Only group incidents with these severities. Leave empty to match incidents of any severity. IDs of `oneuptime_incident_severity` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "incident_labels": schema.SetAttribute{
                MarkdownDescription: "Only group incidents that have at least one of these labels. Leave empty to match incidents regardless of incident labels. IDs of `oneuptime_label` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "monitor_labels": schema.SetAttribute{
                MarkdownDescription: "Only group incidents from monitors that have at least one of these labels. Leave empty to match incidents regardless of monitor labels. IDs of `oneuptime_label` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "incident_title_pattern": schema.StringAttribute{
                MarkdownDescription: "Regular expression pattern to match incident titles. Leave empty to match any title. Example: 'CPU.*high' matches titles containing 'CPU' followed by 'high'.",
                Optional: true,
                Computed: true,
            },
            "incident_description_pattern": schema.StringAttribute{
                MarkdownDescription: "Regular expression pattern to match incident descriptions. Leave empty to match any description.",
                Optional: true,
                Computed: true,
            },
            "monitor_name_pattern": schema.StringAttribute{
                MarkdownDescription: "Regular expression pattern to match monitor names. Leave empty to match any monitor name. Example: 'prod-.*' matches monitors starting with 'prod-'.",
                Optional: true,
                Computed: true,
            },
            "monitor_description_pattern": schema.StringAttribute{
                MarkdownDescription: "Regular expression pattern to match monitor descriptions. Leave empty to match any monitor description.",
                Optional: true,
                Computed: true,
            },
            "group_by_monitor": schema.BoolAttribute{
                MarkdownDescription: "When enabled, incidents from different monitors will be grouped into separate episodes. When disabled, incidents from any monitor can be grouped together.",
                Optional: true,
                Computed: true,
            },
            "group_by_severity": schema.BoolAttribute{
                MarkdownDescription: "When enabled, incidents with different severities will be grouped into separate episodes. When disabled, incidents of any severity can be grouped together.",
                Optional: true,
                Computed: true,
            },
            "group_by_incident_title": schema.BoolAttribute{
                MarkdownDescription: "When enabled, incidents with different titles will be grouped into separate episodes. When disabled, incidents with any title can be grouped together.",
                Optional: true,
                Computed: true,
            },
            "group_by_incident_labels": schema.BoolAttribute{
                MarkdownDescription: "When enabled, incidents with different sets of labels will be grouped into separate episodes (exact set match). When disabled, incident labels are ignored for grouping.",
                Optional: true,
                Computed: true,
            },
            "group_by_monitor_labels": schema.BoolAttribute{
                MarkdownDescription: "When enabled, incidents whose monitors have different sets of labels will be grouped into separate episodes (exact set match). When disabled, monitor labels are ignored for grouping.",
                Optional: true,
                Computed: true,
            },
            "enable_time_window": schema.BoolAttribute{
                MarkdownDescription: "Enable time-based grouping. When enabled, incidents are grouped within the specified time window. When disabled, all matching incidents are grouped into a single ongoing episode regardless of time.",
                Optional: true,
                Computed: true,
            },
            "time_window_minutes": schema.NumberAttribute{
                MarkdownDescription: "Rolling time window in minutes. Incidents are grouped if they arrive within this gap from the last incident.",
                Optional: true,
                Computed: true,
            },
            "group_by_fields": schema.StringAttribute{
                MarkdownDescription: "JSON object defining the fields to group incidents by (e.g., monitorId, severity). A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "episode_title_template": schema.StringAttribute{
                MarkdownDescription: "Template for generating episode titles. Supports placeholders like {{incidentSeverity}}, {{monitorName}}, {{incidentTitle}}, {{incidentDescription}}.",
                Optional: true,
                Computed: true,
            },
            "episode_description_template": schema.StringAttribute{
                MarkdownDescription: "Template for generating episode descriptions. Supports placeholders like {{incidentSeverity}}, {{monitorName}}, {{incidentTitle}}, {{incidentDescription}}.",
                Optional: true,
                Computed: true,
            },
            "enable_resolve_delay": schema.BoolAttribute{
                MarkdownDescription: "Enable grace period before auto-resolving episode after all incidents resolve. Helps prevent rapid state changes during incident flapping.",
                Optional: true,
                Computed: true,
            },
            "resolve_delay_minutes": schema.NumberAttribute{
                MarkdownDescription: "Grace period in minutes before auto-resolving an episode after all incidents are resolved.",
                Optional: true,
                Computed: true,
            },
            "enable_reopen_window": schema.BoolAttribute{
                MarkdownDescription: "Enable reopening recently resolved episodes instead of creating new ones. Useful when related issues recur shortly after resolution.",
                Optional: true,
                Computed: true,
            },
            "reopen_window_minutes": schema.NumberAttribute{
                MarkdownDescription: "Time window in minutes to reopen a recently resolved episode instead of creating a new one.",
                Optional: true,
                Computed: true,
            },
            "enable_inactivity_timeout": schema.BoolAttribute{
                MarkdownDescription: "Enable auto-resolving episodes after a period of inactivity. Helps automatically close episodes when no new incidents arrive.",
                Optional: true,
                Computed: true,
            },
            "inactivity_timeout_minutes": schema.NumberAttribute{
                MarkdownDescription: "Time in minutes after which an inactive episode will be auto-resolved.",
                Optional: true,
                Computed: true,
            },
            "on_call_duty_policies": schema.SetAttribute{
                MarkdownDescription: "List of on-call duty policies to execute for episodes created by this rule. IDs of `oneuptime_on_call_policy` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "default_assign_to_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of defaultAssignToUser. Kept for API compatibility: OneUptime does not show it anywhere. To make someone responsible for the episodes this rule opens, use episodeOwnerUsers. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "default_assign_to_team_id": schema.StringAttribute{
                MarkdownDescription: "ID of defaultAssignToTeam. Kept for API compatibility: OneUptime does not show it anywhere. To make a team responsible for the episodes this rule opens, use episodeOwnerTeams. The ID of a `oneuptime_team`.",
                Optional: true,
                Computed: true,
            },
            "episode_labels": schema.SetAttribute{
                MarkdownDescription: "Labels to automatically apply to episodes created by this rule. IDs of `oneuptime_label` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "episode_owner_users": schema.SetAttribute{
                MarkdownDescription: "Users added as owners of every episode this rule opens, and notified like any owner. Each must be a member of the project. IDs of `oneuptime_user` records.",
                Computed: true,
                ElementType: types.StringType,
            },
            "episode_owner_teams": schema.SetAttribute{
                MarkdownDescription: "Teams added as owners of every episode this rule opens, and notified like any owner. Each must be a team of the project. IDs of `oneuptime_team` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "episode_member_roles": schema.SetAttribute{
                MarkdownDescription: "Incident roles to display in the episode members form. Select the roles that can be assigned to episode members. IDs of `oneuptime_incident_role` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "episode_member_role_assignments": schema.StringAttribute{
                MarkdownDescription: "Users with specific incident roles to automatically add as members to episodes created by this rule. Each assignment includes a user ID and an incident role ID. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "show_episode_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Should episodes created by this rule be shown on the status page?",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *IncidentGroupingRuleDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *IncidentGroupingRuleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data IncidentGroupingRuleDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.Name.IsNull() && !data.Name.IsUnknown() {
        filters["name"] = data.Name.ValueString()
        filterNames = append(filterNames, "name = "+fmt.Sprintf("%q", data.Name.ValueString()))
    }
    if !data.Description.IsNull() && !data.Description.IsUnknown() {
        filters["description"] = data.Description.ValueString()
        filterNames = append(filterNames, "description = "+fmt.Sprintf("%q", data.Description.ValueString()))
    }
    if !data.Priority.IsNull() && !data.Priority.IsUnknown() {
        filters["priority"] = lookupNumber(data.Priority)
        filterNames = append(filterNames, "priority = "+data.Priority.ValueBigFloat().String())
    }
    if !data.IsEnabled.IsNull() && !data.IsEnabled.IsUnknown() {
        filters["isEnabled"] = data.IsEnabled.ValueBool()
        filterNames = append(filterNames, "is_enabled = "+fmt.Sprintf("%t", data.IsEnabled.ValueBool()))
    }
    if !data.IncidentTitlePattern.IsNull() && !data.IncidentTitlePattern.IsUnknown() {
        filters["incidentTitlePattern"] = data.IncidentTitlePattern.ValueString()
        filterNames = append(filterNames, "incident_title_pattern = "+fmt.Sprintf("%q", data.IncidentTitlePattern.ValueString()))
    }
    if !data.IncidentDescriptionPattern.IsNull() && !data.IncidentDescriptionPattern.IsUnknown() {
        filters["incidentDescriptionPattern"] = data.IncidentDescriptionPattern.ValueString()
        filterNames = append(filterNames, "incident_description_pattern = "+fmt.Sprintf("%q", data.IncidentDescriptionPattern.ValueString()))
    }
    if !data.MonitorNamePattern.IsNull() && !data.MonitorNamePattern.IsUnknown() {
        filters["monitorNamePattern"] = data.MonitorNamePattern.ValueString()
        filterNames = append(filterNames, "monitor_name_pattern = "+fmt.Sprintf("%q", data.MonitorNamePattern.ValueString()))
    }
    if !data.MonitorDescriptionPattern.IsNull() && !data.MonitorDescriptionPattern.IsUnknown() {
        filters["monitorDescriptionPattern"] = data.MonitorDescriptionPattern.ValueString()
        filterNames = append(filterNames, "monitor_description_pattern = "+fmt.Sprintf("%q", data.MonitorDescriptionPattern.ValueString()))
    }
    if !data.GroupByMonitor.IsNull() && !data.GroupByMonitor.IsUnknown() {
        filters["groupByMonitor"] = data.GroupByMonitor.ValueBool()
        filterNames = append(filterNames, "group_by_monitor = "+fmt.Sprintf("%t", data.GroupByMonitor.ValueBool()))
    }
    if !data.GroupBySeverity.IsNull() && !data.GroupBySeverity.IsUnknown() {
        filters["groupBySeverity"] = data.GroupBySeverity.ValueBool()
        filterNames = append(filterNames, "group_by_severity = "+fmt.Sprintf("%t", data.GroupBySeverity.ValueBool()))
    }
    if !data.GroupByIncidentTitle.IsNull() && !data.GroupByIncidentTitle.IsUnknown() {
        filters["groupByIncidentTitle"] = data.GroupByIncidentTitle.ValueBool()
        filterNames = append(filterNames, "group_by_incident_title = "+fmt.Sprintf("%t", data.GroupByIncidentTitle.ValueBool()))
    }
    if !data.GroupByIncidentLabels.IsNull() && !data.GroupByIncidentLabels.IsUnknown() {
        filters["groupByIncidentLabels"] = data.GroupByIncidentLabels.ValueBool()
        filterNames = append(filterNames, "group_by_incident_labels = "+fmt.Sprintf("%t", data.GroupByIncidentLabels.ValueBool()))
    }
    if !data.GroupByMonitorLabels.IsNull() && !data.GroupByMonitorLabels.IsUnknown() {
        filters["groupByMonitorLabels"] = data.GroupByMonitorLabels.ValueBool()
        filterNames = append(filterNames, "group_by_monitor_labels = "+fmt.Sprintf("%t", data.GroupByMonitorLabels.ValueBool()))
    }
    if !data.EnableTimeWindow.IsNull() && !data.EnableTimeWindow.IsUnknown() {
        filters["enableTimeWindow"] = data.EnableTimeWindow.ValueBool()
        filterNames = append(filterNames, "enable_time_window = "+fmt.Sprintf("%t", data.EnableTimeWindow.ValueBool()))
    }
    if !data.TimeWindowMinutes.IsNull() && !data.TimeWindowMinutes.IsUnknown() {
        filters["timeWindowMinutes"] = lookupNumber(data.TimeWindowMinutes)
        filterNames = append(filterNames, "time_window_minutes = "+data.TimeWindowMinutes.ValueBigFloat().String())
    }
    if !data.EpisodeTitleTemplate.IsNull() && !data.EpisodeTitleTemplate.IsUnknown() {
        filters["episodeTitleTemplate"] = data.EpisodeTitleTemplate.ValueString()
        filterNames = append(filterNames, "episode_title_template = "+fmt.Sprintf("%q", data.EpisodeTitleTemplate.ValueString()))
    }
    if !data.EpisodeDescriptionTemplate.IsNull() && !data.EpisodeDescriptionTemplate.IsUnknown() {
        filters["episodeDescriptionTemplate"] = data.EpisodeDescriptionTemplate.ValueString()
        filterNames = append(filterNames, "episode_description_template = "+fmt.Sprintf("%q", data.EpisodeDescriptionTemplate.ValueString()))
    }
    if !data.EnableResolveDelay.IsNull() && !data.EnableResolveDelay.IsUnknown() {
        filters["enableResolveDelay"] = data.EnableResolveDelay.ValueBool()
        filterNames = append(filterNames, "enable_resolve_delay = "+fmt.Sprintf("%t", data.EnableResolveDelay.ValueBool()))
    }
    if !data.ResolveDelayMinutes.IsNull() && !data.ResolveDelayMinutes.IsUnknown() {
        filters["resolveDelayMinutes"] = lookupNumber(data.ResolveDelayMinutes)
        filterNames = append(filterNames, "resolve_delay_minutes = "+data.ResolveDelayMinutes.ValueBigFloat().String())
    }
    if !data.EnableReopenWindow.IsNull() && !data.EnableReopenWindow.IsUnknown() {
        filters["enableReopenWindow"] = data.EnableReopenWindow.ValueBool()
        filterNames = append(filterNames, "enable_reopen_window = "+fmt.Sprintf("%t", data.EnableReopenWindow.ValueBool()))
    }
    if !data.ReopenWindowMinutes.IsNull() && !data.ReopenWindowMinutes.IsUnknown() {
        filters["reopenWindowMinutes"] = lookupNumber(data.ReopenWindowMinutes)
        filterNames = append(filterNames, "reopen_window_minutes = "+data.ReopenWindowMinutes.ValueBigFloat().String())
    }
    if !data.EnableInactivityTimeout.IsNull() && !data.EnableInactivityTimeout.IsUnknown() {
        filters["enableInactivityTimeout"] = data.EnableInactivityTimeout.ValueBool()
        filterNames = append(filterNames, "enable_inactivity_timeout = "+fmt.Sprintf("%t", data.EnableInactivityTimeout.ValueBool()))
    }
    if !data.InactivityTimeoutMinutes.IsNull() && !data.InactivityTimeoutMinutes.IsUnknown() {
        filters["inactivityTimeoutMinutes"] = lookupNumber(data.InactivityTimeoutMinutes)
        filterNames = append(filterNames, "inactivity_timeout_minutes = "+data.InactivityTimeoutMinutes.ValueBigFloat().String())
    }
    if !data.DefaultAssignToUserId.IsNull() && !data.DefaultAssignToUserId.IsUnknown() {
        filters["defaultAssignToUserId"] = data.DefaultAssignToUserId.ValueString()
        filterNames = append(filterNames, "default_assign_to_user_id = "+fmt.Sprintf("%q", data.DefaultAssignToUserId.ValueString()))
    }
    if !data.DefaultAssignToTeamId.IsNull() && !data.DefaultAssignToTeamId.IsUnknown() {
        filters["defaultAssignToTeamId"] = data.DefaultAssignToTeamId.ValueString()
        filterNames = append(filterNames, "default_assign_to_team_id = "+fmt.Sprintf("%q", data.DefaultAssignToTeamId.ValueString()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }
    if !data.ShowEpisodeOnStatusPage.IsNull() && !data.ShowEpisodeOnStatusPage.IsUnknown() {
        filters["showEpisodeOnStatusPage"] = data.ShowEpisodeOnStatusPage.ValueBool()
        filterNames = append(filterNames, "show_episode_on_status_page = "+fmt.Sprintf("%t", data.ShowEpisodeOnStatusPage.ValueBool()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the incident grouping rule up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the incident grouping rule up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "criteria": true,
        "projectId": true,
        "name": true,
        "description": true,
        "priority": true,
        "isEnabled": true,
        "matchCriteria": true,
        "monitors": true,
        "incidentSeverities": true,
        "incidentLabels": true,
        "monitorLabels": true,
        "incidentTitlePattern": true,
        "incidentDescriptionPattern": true,
        "monitorNamePattern": true,
        "monitorDescriptionPattern": true,
        "groupByMonitor": true,
        "groupBySeverity": true,
        "groupByIncidentTitle": true,
        "groupByIncidentLabels": true,
        "groupByMonitorLabels": true,
        "enableTimeWindow": true,
        "timeWindowMinutes": true,
        "groupByFields": true,
        "episodeTitleTemplate": true,
        "episodeDescriptionTemplate": true,
        "enableResolveDelay": true,
        "resolveDelayMinutes": true,
        "enableReopenWindow": true,
        "reopenWindowMinutes": true,
        "enableInactivityTimeout": true,
        "inactivityTimeoutMinutes": true,
        "onCallDutyPolicies": true,
        "defaultAssignToUserId": true,
        "defaultAssignToTeamId": true,
        "episodeLabels": true,
        "episodeOwnerUsers": true,
        "episodeOwnerTeams": true,
        "episodeMemberRoles": true,
        "episodeMemberRoleAssignments": true,
        "createdByUserId": true,
        "showEpisodeOnStatusPage": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/incident-grouping-rule/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read incident_grouping_rule, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No incident grouping rule found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read incident_grouping_rule: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/incident-grouping-rule/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list incident_grouping_rule, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list incident_grouping_rule: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No incident grouping rule matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one incident grouping rule matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for incident_grouping_rule.")
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
    if obj, ok := item["criteria"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Criteria = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Criteria = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Criteria = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Criteria = types.StringValue(string(jsonBytes))
        } else {
            data.Criteria = types.StringNull()
        }
    } else if val, ok := item["criteria"].(string); ok {
        data.Criteria = types.StringValue(val)
    } else {
        data.Criteria = types.StringNull()
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
    if obj, ok := item["name"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Name = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Name = types.StringValue(string(jsonBytes))
        } else {
            data.Name = types.StringNull()
        }
    } else if val, ok := item["name"].(string); ok {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
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
    if val, ok := item["priority"].(float64); ok {
        data.Priority = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["priority"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.Priority = types.NumberValue(big.NewFloat(val))
        } else {
            data.Priority = types.NumberNull()
        }
    } else {
        data.Priority = types.NumberNull()
    }
    if val, ok := item["isEnabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    } else {
        data.IsEnabled = types.BoolNull()
    }
    if obj, ok := item["matchCriteria"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MatchCriteria = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.MatchCriteria = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.MatchCriteria = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.MatchCriteria = types.StringValue(string(jsonBytes))
        } else {
            data.MatchCriteria = types.StringNull()
        }
    } else if val, ok := item["matchCriteria"].(string); ok {
        data.MatchCriteria = types.StringValue(val)
    } else {
        data.MatchCriteria = types.StringNull()
    }
    if val, ok := item["monitors"].([]interface{}); ok {
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
        data.Monitors = types.SetValueMust(types.StringType, setItems)
    } else {
        data.Monitors = types.SetNull(types.StringType)
    }
    if val, ok := item["incidentSeverities"].([]interface{}); ok {
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
        data.IncidentSeverities = types.SetValueMust(types.StringType, setItems)
    } else {
        data.IncidentSeverities = types.SetNull(types.StringType)
    }
    if val, ok := item["incidentLabels"].([]interface{}); ok {
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
        data.IncidentLabels = types.SetValueMust(types.StringType, setItems)
    } else {
        data.IncidentLabels = types.SetNull(types.StringType)
    }
    if val, ok := item["monitorLabels"].([]interface{}); ok {
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
        data.MonitorLabels = types.SetValueMust(types.StringType, setItems)
    } else {
        data.MonitorLabels = types.SetNull(types.StringType)
    }
    if obj, ok := item["incidentTitlePattern"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncidentTitlePattern = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.IncidentTitlePattern = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.IncidentTitlePattern = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.IncidentTitlePattern = types.StringValue(string(jsonBytes))
        } else {
            data.IncidentTitlePattern = types.StringNull()
        }
    } else if val, ok := item["incidentTitlePattern"].(string); ok {
        data.IncidentTitlePattern = types.StringValue(val)
    } else {
        data.IncidentTitlePattern = types.StringNull()
    }
    if obj, ok := item["incidentDescriptionPattern"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncidentDescriptionPattern = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.IncidentDescriptionPattern = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.IncidentDescriptionPattern = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.IncidentDescriptionPattern = types.StringValue(string(jsonBytes))
        } else {
            data.IncidentDescriptionPattern = types.StringNull()
        }
    } else if val, ok := item["incidentDescriptionPattern"].(string); ok {
        data.IncidentDescriptionPattern = types.StringValue(val)
    } else {
        data.IncidentDescriptionPattern = types.StringNull()
    }
    if obj, ok := item["monitorNamePattern"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MonitorNamePattern = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.MonitorNamePattern = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.MonitorNamePattern = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.MonitorNamePattern = types.StringValue(string(jsonBytes))
        } else {
            data.MonitorNamePattern = types.StringNull()
        }
    } else if val, ok := item["monitorNamePattern"].(string); ok {
        data.MonitorNamePattern = types.StringValue(val)
    } else {
        data.MonitorNamePattern = types.StringNull()
    }
    if obj, ok := item["monitorDescriptionPattern"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MonitorDescriptionPattern = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.MonitorDescriptionPattern = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.MonitorDescriptionPattern = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.MonitorDescriptionPattern = types.StringValue(string(jsonBytes))
        } else {
            data.MonitorDescriptionPattern = types.StringNull()
        }
    } else if val, ok := item["monitorDescriptionPattern"].(string); ok {
        data.MonitorDescriptionPattern = types.StringValue(val)
    } else {
        data.MonitorDescriptionPattern = types.StringNull()
    }
    if val, ok := item["groupByMonitor"].(bool); ok {
        data.GroupByMonitor = types.BoolValue(val)
    } else {
        data.GroupByMonitor = types.BoolNull()
    }
    if val, ok := item["groupBySeverity"].(bool); ok {
        data.GroupBySeverity = types.BoolValue(val)
    } else {
        data.GroupBySeverity = types.BoolNull()
    }
    if val, ok := item["groupByIncidentTitle"].(bool); ok {
        data.GroupByIncidentTitle = types.BoolValue(val)
    } else {
        data.GroupByIncidentTitle = types.BoolNull()
    }
    if val, ok := item["groupByIncidentLabels"].(bool); ok {
        data.GroupByIncidentLabels = types.BoolValue(val)
    } else {
        data.GroupByIncidentLabels = types.BoolNull()
    }
    if val, ok := item["groupByMonitorLabels"].(bool); ok {
        data.GroupByMonitorLabels = types.BoolValue(val)
    } else {
        data.GroupByMonitorLabels = types.BoolNull()
    }
    if val, ok := item["enableTimeWindow"].(bool); ok {
        data.EnableTimeWindow = types.BoolValue(val)
    } else {
        data.EnableTimeWindow = types.BoolNull()
    }
    if val, ok := item["timeWindowMinutes"].(float64); ok {
        data.TimeWindowMinutes = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["timeWindowMinutes"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.TimeWindowMinutes = types.NumberValue(big.NewFloat(val))
        } else {
            data.TimeWindowMinutes = types.NumberNull()
        }
    } else {
        data.TimeWindowMinutes = types.NumberNull()
    }
    if obj, ok := item["groupByFields"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.GroupByFields = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.GroupByFields = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.GroupByFields = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.GroupByFields = types.StringValue(string(jsonBytes))
        } else {
            data.GroupByFields = types.StringNull()
        }
    } else if val, ok := item["groupByFields"].(string); ok {
        data.GroupByFields = types.StringValue(val)
    } else {
        data.GroupByFields = types.StringNull()
    }
    if obj, ok := item["episodeTitleTemplate"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EpisodeTitleTemplate = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.EpisodeTitleTemplate = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.EpisodeTitleTemplate = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.EpisodeTitleTemplate = types.StringValue(string(jsonBytes))
        } else {
            data.EpisodeTitleTemplate = types.StringNull()
        }
    } else if val, ok := item["episodeTitleTemplate"].(string); ok {
        data.EpisodeTitleTemplate = types.StringValue(val)
    } else {
        data.EpisodeTitleTemplate = types.StringNull()
    }
    if obj, ok := item["episodeDescriptionTemplate"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EpisodeDescriptionTemplate = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.EpisodeDescriptionTemplate = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.EpisodeDescriptionTemplate = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.EpisodeDescriptionTemplate = types.StringValue(string(jsonBytes))
        } else {
            data.EpisodeDescriptionTemplate = types.StringNull()
        }
    } else if val, ok := item["episodeDescriptionTemplate"].(string); ok {
        data.EpisodeDescriptionTemplate = types.StringValue(val)
    } else {
        data.EpisodeDescriptionTemplate = types.StringNull()
    }
    if val, ok := item["enableResolveDelay"].(bool); ok {
        data.EnableResolveDelay = types.BoolValue(val)
    } else {
        data.EnableResolveDelay = types.BoolNull()
    }
    if val, ok := item["resolveDelayMinutes"].(float64); ok {
        data.ResolveDelayMinutes = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["resolveDelayMinutes"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.ResolveDelayMinutes = types.NumberValue(big.NewFloat(val))
        } else {
            data.ResolveDelayMinutes = types.NumberNull()
        }
    } else {
        data.ResolveDelayMinutes = types.NumberNull()
    }
    if val, ok := item["enableReopenWindow"].(bool); ok {
        data.EnableReopenWindow = types.BoolValue(val)
    } else {
        data.EnableReopenWindow = types.BoolNull()
    }
    if val, ok := item["reopenWindowMinutes"].(float64); ok {
        data.ReopenWindowMinutes = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["reopenWindowMinutes"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.ReopenWindowMinutes = types.NumberValue(big.NewFloat(val))
        } else {
            data.ReopenWindowMinutes = types.NumberNull()
        }
    } else {
        data.ReopenWindowMinutes = types.NumberNull()
    }
    if val, ok := item["enableInactivityTimeout"].(bool); ok {
        data.EnableInactivityTimeout = types.BoolValue(val)
    } else {
        data.EnableInactivityTimeout = types.BoolNull()
    }
    if val, ok := item["inactivityTimeoutMinutes"].(float64); ok {
        data.InactivityTimeoutMinutes = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["inactivityTimeoutMinutes"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.InactivityTimeoutMinutes = types.NumberValue(big.NewFloat(val))
        } else {
            data.InactivityTimeoutMinutes = types.NumberNull()
        }
    } else {
        data.InactivityTimeoutMinutes = types.NumberNull()
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
    if obj, ok := item["defaultAssignToUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DefaultAssignToUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DefaultAssignToUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DefaultAssignToUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DefaultAssignToUserId = types.StringValue(string(jsonBytes))
        } else {
            data.DefaultAssignToUserId = types.StringNull()
        }
    } else if val, ok := item["defaultAssignToUserId"].(string); ok {
        data.DefaultAssignToUserId = types.StringValue(val)
    } else {
        data.DefaultAssignToUserId = types.StringNull()
    }
    if obj, ok := item["defaultAssignToTeamId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DefaultAssignToTeamId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DefaultAssignToTeamId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DefaultAssignToTeamId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DefaultAssignToTeamId = types.StringValue(string(jsonBytes))
        } else {
            data.DefaultAssignToTeamId = types.StringNull()
        }
    } else if val, ok := item["defaultAssignToTeamId"].(string); ok {
        data.DefaultAssignToTeamId = types.StringValue(val)
    } else {
        data.DefaultAssignToTeamId = types.StringNull()
    }
    if val, ok := item["episodeLabels"].([]interface{}); ok {
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
        data.EpisodeLabels = types.SetValueMust(types.StringType, setItems)
    } else {
        data.EpisodeLabels = types.SetNull(types.StringType)
    }
    if val, ok := item["episodeOwnerUsers"].([]interface{}); ok {
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
        data.EpisodeOwnerUsers = types.SetValueMust(types.StringType, setItems)
    } else {
        data.EpisodeOwnerUsers = types.SetNull(types.StringType)
    }
    if val, ok := item["episodeOwnerTeams"].([]interface{}); ok {
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
        data.EpisodeOwnerTeams = types.SetValueMust(types.StringType, setItems)
    } else {
        data.EpisodeOwnerTeams = types.SetNull(types.StringType)
    }
    if val, ok := item["episodeMemberRoles"].([]interface{}); ok {
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
        data.EpisodeMemberRoles = types.SetValueMust(types.StringType, setItems)
    } else {
        data.EpisodeMemberRoles = types.SetNull(types.StringType)
    }
    if obj, ok := item["episodeMemberRoleAssignments"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EpisodeMemberRoleAssignments = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.EpisodeMemberRoleAssignments = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.EpisodeMemberRoleAssignments = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.EpisodeMemberRoleAssignments = types.StringValue(string(jsonBytes))
        } else {
            data.EpisodeMemberRoleAssignments = types.StringNull()
        }
    } else if val, ok := item["episodeMemberRoleAssignments"].(string); ok {
        data.EpisodeMemberRoleAssignments = types.StringValue(val)
    } else {
        data.EpisodeMemberRoleAssignments = types.StringNull()
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
    if val, ok := item["showEpisodeOnStatusPage"].(bool); ok {
        data.ShowEpisodeOnStatusPage = types.BoolValue(val)
    } else {
        data.ShowEpisodeOnStatusPage = types.BoolNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
