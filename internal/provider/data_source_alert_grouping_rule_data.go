package provider

import (
    "context"
    "fmt"
    "math/big"
    "github.com/hashicorp/terraform-plugin-framework/attr"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &AlertGroupingRuleDataDataSource{}

func NewAlertGroupingRuleDataDataSource() datasource.DataSource {
    return &AlertGroupingRuleDataDataSource{}
}

// AlertGroupingRuleDataDataSource defines the data source implementation.
type AlertGroupingRuleDataDataSource struct {
    client *Client
}

// AlertGroupingRuleDataDataSourceModel describes the data source data model.
type AlertGroupingRuleDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    Description types.String `tfsdk:"description"`
    Priority types.Number `tfsdk:"priority"`
    IsEnabled types.Bool `tfsdk:"is_enabled"`
    MatchCriteria types.String `tfsdk:"match_criteria"`
    Monitors types.Set `tfsdk:"monitors"`
    AlertSeverities types.Set `tfsdk:"alert_severities"`
    AlertLabels types.Set `tfsdk:"alert_labels"`
    MonitorLabels types.Set `tfsdk:"monitor_labels"`
    AlertTitlePattern types.String `tfsdk:"alert_title_pattern"`
    AlertDescriptionPattern types.String `tfsdk:"alert_description_pattern"`
    MonitorNamePattern types.String `tfsdk:"monitor_name_pattern"`
    MonitorDescriptionPattern types.String `tfsdk:"monitor_description_pattern"`
    GroupByMonitor types.Bool `tfsdk:"group_by_monitor"`
    GroupBySeverity types.Bool `tfsdk:"group_by_severity"`
    GroupByAlertTitle types.Bool `tfsdk:"group_by_alert_title"`
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
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *AlertGroupingRuleDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_alert_grouping_rule_data"
}

func (d *AlertGroupingRuleDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "alert_grouping_rule_data data source",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Identifier to filter by",
                Optional: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Name to filter by",
                Computed: true,
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
            "description": schema.StringAttribute{
                MarkdownDescription: "Description of this alert grouping rule. Permissions - Create: [Project Owner, Project Admin, Create Alert Grouping Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Alert Admin, Alert Member, Alert Viewer, Read Alert Grouping Rule], Update: [Project Owner, Project Admin, Edit Alert Grouping Rule]",
                Computed: true,
            },
            "priority": schema.NumberAttribute{
                MarkdownDescription: "Priority of this rule. Lower number = higher priority. Rules are evaluated in priority order.. Permissions - Create: [Project Owner, Project Admin, Create Alert Grouping Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Alert Admin, Alert Member, Alert Viewer, Read Alert Grouping Rule], Update: [Project Owner, Project Admin, Edit Alert Grouping Rule]",
                Computed: true,
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether this rule is enabled. Permissions - Create: [Project Owner, Project Admin, Create Alert Grouping Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Alert Admin, Alert Member, Alert Viewer, Read Alert Grouping Rule], Update: [Project Owner, Project Admin, Edit Alert Grouping Rule]",
                Computed: true,
            },
            "match_criteria": schema.StringAttribute{
                MarkdownDescription: "JSON object defining the criteria for matching alerts to this rule. Permissions - Create: [Project Owner, Project Admin, Create Alert Grouping Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Alert Admin, Alert Member, Alert Viewer, Read Alert Grouping Rule], Update: [Project Owner, Project Admin, Edit Alert Grouping Rule]",
                Computed: true,
            },
            "monitors": schema.SetAttribute{
                MarkdownDescription: "Only group alerts from these monitors. Leave empty to match alerts from any monitor.. Permissions - Create: [Project Owner, Project Admin, Create Alert Grouping Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Alert Admin, Alert Member, Alert Viewer, Read Alert Grouping Rule], Update: [Project Owner, Project Admin, Edit Alert Grouping Rule]",
                Computed: true,
                ElementType: types.StringType,
            },
            "alert_severities": schema.SetAttribute{
                MarkdownDescription: "Only group alerts with these severities. Leave empty to match alerts of any severity.. Permissions - Create: [Project Owner, Project Admin, Create Alert Grouping Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Alert Admin, Alert Member, Alert Viewer, Read Alert Grouping Rule], Update: [Project Owner, Project Admin, Edit Alert Grouping Rule]",
                Computed: true,
                ElementType: types.StringType,
            },
            "alert_labels": schema.SetAttribute{
                MarkdownDescription: "Only group alerts that have at least one of these labels. Leave empty to match alerts regardless of alert labels.. Permissions - Create: [Project Owner, Project Admin, Create Alert Grouping Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Alert Admin, Alert Member, Alert Viewer, Read Alert Grouping Rule], Update: [Project Owner, Project Admin, Edit Alert Grouping Rule]",
                Computed: true,
                ElementType: types.StringType,
            },
            "monitor_labels": schema.SetAttribute{
                MarkdownDescription: "Only group alerts from monitors that have at least one of these labels. Leave empty to match alerts regardless of monitor labels.. Permissions - Create: [Project Owner, Project Admin, Create Alert Grouping Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Alert Admin, Alert Member, Alert Viewer, Read Alert Grouping Rule], Update: [Project Owner, Project Admin, Edit Alert Grouping Rule]",
                Computed: true,
                ElementType: types.StringType,
            },
            "alert_title_pattern": schema.StringAttribute{
                MarkdownDescription: "Regular expression pattern to match alert titles. Leave empty to match any title. Example: 'CPU.*high' matches titles containing 'CPU' followed by 'high'.. Permissions - Create: [Project Owner, Project Admin, Create Alert Grouping Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Alert Admin, Alert Member, Alert Viewer, Read Alert Grouping Rule], Update: [Project Owner, Project Admin, Edit Alert Grouping Rule]",
                Computed: true,
            },
            "alert_description_pattern": schema.StringAttribute{
                MarkdownDescription: "Regular expression pattern to match alert descriptions. Leave empty to match any description.. Permissions - Create: [Project Owner, Project Admin, Create Alert Grouping Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Alert Admin, Alert Member, Alert Viewer, Read Alert Grouping Rule], Update: [Project Owner, Project Admin, Edit Alert Grouping Rule]",
                Computed: true,
            },
            "monitor_name_pattern": schema.StringAttribute{
                MarkdownDescription: "Regular expression pattern to match monitor names. Leave empty to match any monitor name. Example: 'prod-.*' matches monitors starting with 'prod-'.. Permissions - Create: [Project Owner, Project Admin, Create Alert Grouping Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Alert Admin, Alert Member, Alert Viewer, Read Alert Grouping Rule], Update: [Project Owner, Project Admin, Edit Alert Grouping Rule]",
                Computed: true,
            },
            "monitor_description_pattern": schema.StringAttribute{
                MarkdownDescription: "Regular expression pattern to match monitor descriptions. Leave empty to match any monitor description.. Permissions - Create: [Project Owner, Project Admin, Create Alert Grouping Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Alert Admin, Alert Member, Alert Viewer, Read Alert Grouping Rule], Update: [Project Owner, Project Admin, Edit Alert Grouping Rule]",
                Computed: true,
            },
            "group_by_monitor": schema.BoolAttribute{
                MarkdownDescription: "When enabled, alerts from different monitors will be grouped into separate episodes. When disabled, alerts from any monitor can be grouped together.. Permissions - Create: [Project Owner, Project Admin, Create Alert Grouping Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Alert Admin, Alert Member, Alert Viewer, Read Alert Grouping Rule], Update: [Project Owner, Project Admin, Edit Alert Grouping Rule]",
                Computed: true,
            },
            "group_by_severity": schema.BoolAttribute{
                MarkdownDescription: "When enabled, alerts with different severities will be grouped into separate episodes. When disabled, alerts of any severity can be grouped together.. Permissions - Create: [Project Owner, Project Admin, Create Alert Grouping Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Alert Admin, Alert Member, Alert Viewer, Read Alert Grouping Rule], Update: [Project Owner, Project Admin, Edit Alert Grouping Rule]",
                Computed: true,
            },
            "group_by_alert_title": schema.BoolAttribute{
                MarkdownDescription: "When enabled, alerts with different titles will be grouped into separate episodes. When disabled, alerts with any title can be grouped together.. Permissions - Create: [Project Owner, Project Admin, Create Alert Grouping Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Alert Admin, Alert Member, Alert Viewer, Read Alert Grouping Rule], Update: [Project Owner, Project Admin, Edit Alert Grouping Rule]",
                Computed: true,
            },
            "enable_time_window": schema.BoolAttribute{
                MarkdownDescription: "Enable time-based grouping. When enabled, alerts are grouped within the specified time window. When disabled, all matching alerts are grouped into a single ongoing episode regardless of time.. Permissions - Create: [Project Owner, Project Admin, Create Alert Grouping Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Alert Admin, Alert Member, Alert Viewer, Read Alert Grouping Rule], Update: [Project Owner, Project Admin, Edit Alert Grouping Rule]",
                Computed: true,
            },
            "time_window_minutes": schema.NumberAttribute{
                MarkdownDescription: "Rolling time window in minutes. Alerts are grouped if they arrive within this gap from the last alert.. Permissions - Create: [Project Owner, Project Admin, Create Alert Grouping Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Alert Admin, Alert Member, Alert Viewer, Read Alert Grouping Rule], Update: [Project Owner, Project Admin, Edit Alert Grouping Rule]",
                Computed: true,
            },
            "group_by_fields": schema.StringAttribute{
                MarkdownDescription: "JSON object defining the fields to group alerts by (e.g., monitorId, severity). Permissions - Create: [Project Owner, Project Admin, Create Alert Grouping Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Alert Admin, Alert Member, Alert Viewer, Read Alert Grouping Rule], Update: [Project Owner, Project Admin, Edit Alert Grouping Rule]",
                Computed: true,
            },
            "episode_title_template": schema.StringAttribute{
                MarkdownDescription: "Template for generating episode titles. Supports placeholders like {{alertSeverity}}, {{monitorName}}, {{alertTitle}}, {{alertDescription}}. Permissions - Create: [Project Owner, Project Admin, Create Alert Grouping Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Alert Admin, Alert Member, Alert Viewer, Read Alert Grouping Rule], Update: [Project Owner, Project Admin, Edit Alert Grouping Rule]",
                Computed: true,
            },
            "episode_description_template": schema.StringAttribute{
                MarkdownDescription: "Template for generating episode descriptions. Supports placeholders like {{alertSeverity}}, {{monitorName}}, {{alertTitle}}, {{alertDescription}}. Permissions - Create: [Project Owner, Project Admin, Create Alert Grouping Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Alert Admin, Alert Member, Alert Viewer, Read Alert Grouping Rule], Update: [Project Owner, Project Admin, Edit Alert Grouping Rule]",
                Computed: true,
            },
            "enable_resolve_delay": schema.BoolAttribute{
                MarkdownDescription: "Enable grace period before auto-resolving episode after all alerts resolve. Helps prevent rapid state changes during alert flapping.. Permissions - Create: [Project Owner, Project Admin, Create Alert Grouping Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Alert Admin, Alert Member, Alert Viewer, Read Alert Grouping Rule], Update: [Project Owner, Project Admin, Edit Alert Grouping Rule]",
                Computed: true,
            },
            "resolve_delay_minutes": schema.NumberAttribute{
                MarkdownDescription: "Grace period in minutes before auto-resolving an episode after all alerts are resolved. Permissions - Create: [Project Owner, Project Admin, Create Alert Grouping Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Alert Admin, Alert Member, Alert Viewer, Read Alert Grouping Rule], Update: [Project Owner, Project Admin, Edit Alert Grouping Rule]",
                Computed: true,
            },
            "enable_reopen_window": schema.BoolAttribute{
                MarkdownDescription: "Enable reopening recently resolved episodes instead of creating new ones. Useful when related issues recur shortly after resolution.. Permissions - Create: [Project Owner, Project Admin, Create Alert Grouping Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Alert Admin, Alert Member, Alert Viewer, Read Alert Grouping Rule], Update: [Project Owner, Project Admin, Edit Alert Grouping Rule]",
                Computed: true,
            },
            "reopen_window_minutes": schema.NumberAttribute{
                MarkdownDescription: "Time window in minutes to reopen a recently resolved episode instead of creating a new one. Permissions - Create: [Project Owner, Project Admin, Create Alert Grouping Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Alert Admin, Alert Member, Alert Viewer, Read Alert Grouping Rule], Update: [Project Owner, Project Admin, Edit Alert Grouping Rule]",
                Computed: true,
            },
            "enable_inactivity_timeout": schema.BoolAttribute{
                MarkdownDescription: "Enable auto-resolving episodes after a period of inactivity. Helps automatically close episodes when no new alerts arrive.. Permissions - Create: [Project Owner, Project Admin, Create Alert Grouping Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Alert Admin, Alert Member, Alert Viewer, Read Alert Grouping Rule], Update: [Project Owner, Project Admin, Edit Alert Grouping Rule]",
                Computed: true,
            },
            "inactivity_timeout_minutes": schema.NumberAttribute{
                MarkdownDescription: "Time in minutes after which an inactive episode will be auto-resolved. Permissions - Create: [Project Owner, Project Admin, Create Alert Grouping Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Alert Admin, Alert Member, Alert Viewer, Read Alert Grouping Rule], Update: [Project Owner, Project Admin, Edit Alert Grouping Rule]",
                Computed: true,
            },
            "on_call_duty_policies": schema.SetAttribute{
                MarkdownDescription: "List of on-call duty policies to execute for episodes created by this rule.. Permissions - Create: [Project Owner, Project Admin, Create Alert Grouping Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Alert Admin, Alert Member, Alert Viewer, Read Alert Grouping Rule], Update: [Project Owner, Project Admin, Edit Alert Grouping Rule]",
                Computed: true,
                ElementType: types.StringType,
            },
            "default_assign_to_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "default_assign_to_team_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "episode_labels": schema.SetAttribute{
                MarkdownDescription: "Labels to automatically apply to episodes created by this rule.. Permissions - Create: [Project Owner, Project Admin, Create Alert Grouping Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Alert Admin, Alert Member, Alert Viewer, Read Alert Grouping Rule], Update: [Project Owner, Project Admin, Edit Alert Grouping Rule]",
                Computed: true,
                ElementType: types.StringType,
            },
            "episode_owner_users": schema.SetAttribute{
                MarkdownDescription: "Users to automatically add as owners to episodes created by this rule.. Permissions - Create: [Project Owner, Project Admin, Create Alert Grouping Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Alert Admin, Alert Member, Alert Viewer, Read Alert Grouping Rule], Update: [Project Owner, Project Admin, Edit Alert Grouping Rule]",
                Computed: true,
                ElementType: types.StringType,
            },
            "episode_owner_teams": schema.SetAttribute{
                MarkdownDescription: "Teams to automatically add as owners to episodes created by this rule.. Permissions - Create: [Project Owner, Project Admin, Create Alert Grouping Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Alert Admin, Alert Member, Alert Viewer, Read Alert Grouping Rule], Update: [Project Owner, Project Admin, Edit Alert Grouping Rule]",
                Computed: true,
                ElementType: types.StringType,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
        },
    }
}

func (d *AlertGroupingRuleDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AlertGroupingRuleDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data AlertGroupingRuleDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "alert-grouping-rule" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read alert_grouping_rule_data, got error: %s", err))
        return
    }

    var alertGroupingRuleDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &alertGroupingRuleDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse alert_grouping_rule_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := alertGroupingRuleDataResponse["data"].(map[string]interface{}); ok {
        alertGroupingRuleDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := alertGroupingRuleDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := alertGroupingRuleDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := alertGroupingRuleDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := alertGroupingRuleDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := alertGroupingRuleDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := alertGroupingRuleDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := alertGroupingRuleDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := alertGroupingRuleDataResponse["description"].(string); ok {
        data.Description = types.StringValue(val)
    }
    if val, ok := alertGroupingRuleDataResponse["priority"].(float64); ok {
        data.Priority = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := alertGroupingRuleDataResponse["is_enabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    }
    if val, ok := alertGroupingRuleDataResponse["match_criteria"].(string); ok {
        data.MatchCriteria = types.StringValue(val)
    }
    if val, ok := alertGroupingRuleDataResponse["monitors"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.Monitors = setValue
    }
    if val, ok := alertGroupingRuleDataResponse["alert_severities"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.AlertSeverities = setValue
    }
    if val, ok := alertGroupingRuleDataResponse["alert_labels"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.AlertLabels = setValue
    }
    if val, ok := alertGroupingRuleDataResponse["monitor_labels"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.MonitorLabels = setValue
    }
    if val, ok := alertGroupingRuleDataResponse["alert_title_pattern"].(string); ok {
        data.AlertTitlePattern = types.StringValue(val)
    }
    if val, ok := alertGroupingRuleDataResponse["alert_description_pattern"].(string); ok {
        data.AlertDescriptionPattern = types.StringValue(val)
    }
    if val, ok := alertGroupingRuleDataResponse["monitor_name_pattern"].(string); ok {
        data.MonitorNamePattern = types.StringValue(val)
    }
    if val, ok := alertGroupingRuleDataResponse["monitor_description_pattern"].(string); ok {
        data.MonitorDescriptionPattern = types.StringValue(val)
    }
    if val, ok := alertGroupingRuleDataResponse["group_by_monitor"].(bool); ok {
        data.GroupByMonitor = types.BoolValue(val)
    }
    if val, ok := alertGroupingRuleDataResponse["group_by_severity"].(bool); ok {
        data.GroupBySeverity = types.BoolValue(val)
    }
    if val, ok := alertGroupingRuleDataResponse["group_by_alert_title"].(bool); ok {
        data.GroupByAlertTitle = types.BoolValue(val)
    }
    if val, ok := alertGroupingRuleDataResponse["enable_time_window"].(bool); ok {
        data.EnableTimeWindow = types.BoolValue(val)
    }
    if val, ok := alertGroupingRuleDataResponse["time_window_minutes"].(float64); ok {
        data.TimeWindowMinutes = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := alertGroupingRuleDataResponse["group_by_fields"].(string); ok {
        data.GroupByFields = types.StringValue(val)
    }
    if val, ok := alertGroupingRuleDataResponse["episode_title_template"].(string); ok {
        data.EpisodeTitleTemplate = types.StringValue(val)
    }
    if val, ok := alertGroupingRuleDataResponse["episode_description_template"].(string); ok {
        data.EpisodeDescriptionTemplate = types.StringValue(val)
    }
    if val, ok := alertGroupingRuleDataResponse["enable_resolve_delay"].(bool); ok {
        data.EnableResolveDelay = types.BoolValue(val)
    }
    if val, ok := alertGroupingRuleDataResponse["resolve_delay_minutes"].(float64); ok {
        data.ResolveDelayMinutes = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := alertGroupingRuleDataResponse["enable_reopen_window"].(bool); ok {
        data.EnableReopenWindow = types.BoolValue(val)
    }
    if val, ok := alertGroupingRuleDataResponse["reopen_window_minutes"].(float64); ok {
        data.ReopenWindowMinutes = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := alertGroupingRuleDataResponse["enable_inactivity_timeout"].(bool); ok {
        data.EnableInactivityTimeout = types.BoolValue(val)
    }
    if val, ok := alertGroupingRuleDataResponse["inactivity_timeout_minutes"].(float64); ok {
        data.InactivityTimeoutMinutes = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := alertGroupingRuleDataResponse["on_call_duty_policies"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.OnCallDutyPolicies = setValue
    }
    if val, ok := alertGroupingRuleDataResponse["default_assign_to_user_id"].(string); ok {
        data.DefaultAssignToUserId = types.StringValue(val)
    }
    if val, ok := alertGroupingRuleDataResponse["default_assign_to_team_id"].(string); ok {
        data.DefaultAssignToTeamId = types.StringValue(val)
    }
    if val, ok := alertGroupingRuleDataResponse["episode_labels"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.EpisodeLabels = setValue
    }
    if val, ok := alertGroupingRuleDataResponse["episode_owner_users"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.EpisodeOwnerUsers = setValue
    }
    if val, ok := alertGroupingRuleDataResponse["episode_owner_teams"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.EpisodeOwnerTeams = setValue
    }
    if val, ok := alertGroupingRuleDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
