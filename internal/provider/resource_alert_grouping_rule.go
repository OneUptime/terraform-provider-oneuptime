package provider

import (
    "context"
    "fmt"
    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-framework/types/basetypes"
    "github.com/hashicorp/terraform-plugin-log/tflog"
    "math/big"
    "net/http"
    "github.com/hashicorp/terraform-plugin-framework/path"
    "encoding/json"
    "net/url"
    "strings"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/numberdefault"
    "github.com/hashicorp/terraform-plugin-framework/attr"
    "sort"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/numberplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &AlertGroupingRuleResource{}
var _ resource.ResourceWithImportState = &AlertGroupingRuleResource{}

func NewAlertGroupingRuleResource() resource.Resource {
    return &AlertGroupingRuleResource{}
}

// AlertGroupingRuleResource defines the resource implementation.
type AlertGroupingRuleResource struct {
    client *Client
}

// AlertGroupingRuleResourceModel describes the resource data model.
type AlertGroupingRuleResourceModel struct {
    Id types.String `tfsdk:"id"`
    ProjectId types.String `tfsdk:"project_id"`
    Name types.String `tfsdk:"name"`
    Description types.String `tfsdk:"description"`
    Priority types.Number `tfsdk:"priority"`
    IsEnabled types.Bool `tfsdk:"is_enabled"`
    MatchCriteria JSONSubsetValue `tfsdk:"match_criteria"`
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
    GroupByAlertLabels types.Bool `tfsdk:"group_by_alert_labels"`
    GroupByMonitorLabels types.Bool `tfsdk:"group_by_monitor_labels"`
    EnableTimeWindow types.Bool `tfsdk:"enable_time_window"`
    TimeWindowMinutes types.Number `tfsdk:"time_window_minutes"`
    GroupByFields JSONSubsetValue `tfsdk:"group_by_fields"`
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
    CreatedAt RFC3339Value `tfsdk:"created_at"`
    UpdatedAt RFC3339Value `tfsdk:"updated_at"`
    DeletedAt RFC3339Value `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
}

func (r *AlertGroupingRuleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_alert_grouping_rule"
}

func (r *AlertGroupingRuleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Configure rules for automatically grouping related alerts into episodes",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Unique identifier for the resource",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Name of this alert grouping rule.",
                Required: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Description of this alert grouping rule.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "priority": schema.NumberAttribute{
                MarkdownDescription: "Priority of this rule. Lower number = higher priority. Rules are evaluated in priority order..",
                Optional: true,
                Computed: true,
                Default: numberdefault.StaticBigFloat(big.NewFloat(1)),
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether this rule is enabled.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(true),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "match_criteria": schema.StringAttribute{
                MarkdownDescription: "JSON object defining the criteria for matching alerts to this rule.",
                CustomType: JSONSubsetType{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: []validator.String{
                    JSONEnvelopeValidator(),
                },
            },
            "monitors": schema.SetAttribute{
                MarkdownDescription: "Only group alerts from these monitors. Leave empty to match alerts from any monitor..",
                Optional: true,
                Computed: true,
                ElementType: types.StringType,
                PlanModifiers: []planmodifier.Set{
                    setplanmodifier.UseStateForUnknown(),
                },
            },
            "alert_severities": schema.SetAttribute{
                MarkdownDescription: "Only group alerts with these severities. Leave empty to match alerts of any severity..",
                Optional: true,
                Computed: true,
                ElementType: types.StringType,
                PlanModifiers: []planmodifier.Set{
                    setplanmodifier.UseStateForUnknown(),
                },
            },
            "alert_labels": schema.SetAttribute{
                MarkdownDescription: "Only group alerts that have at least one of these labels. Leave empty to match alerts regardless of alert labels..",
                Optional: true,
                Computed: true,
                ElementType: types.StringType,
                PlanModifiers: []planmodifier.Set{
                    setplanmodifier.UseStateForUnknown(),
                },
            },
            "monitor_labels": schema.SetAttribute{
                MarkdownDescription: "Only group alerts from monitors that have at least one of these labels. Leave empty to match alerts regardless of monitor labels..",
                Optional: true,
                Computed: true,
                ElementType: types.StringType,
                PlanModifiers: []planmodifier.Set{
                    setplanmodifier.UseStateForUnknown(),
                },
            },
            "alert_title_pattern": schema.StringAttribute{
                MarkdownDescription: "Regular expression pattern to match alert titles. Leave empty to match any title. Example: 'CPU.*high' matches titles containing 'CPU' followed by 'high'..",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "alert_description_pattern": schema.StringAttribute{
                MarkdownDescription: "Regular expression pattern to match alert descriptions. Leave empty to match any description..",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "monitor_name_pattern": schema.StringAttribute{
                MarkdownDescription: "Regular expression pattern to match monitor names. Leave empty to match any monitor name. Example: 'prod-.*' matches monitors starting with 'prod-'..",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "monitor_description_pattern": schema.StringAttribute{
                MarkdownDescription: "Regular expression pattern to match monitor descriptions. Leave empty to match any monitor description..",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "group_by_monitor": schema.BoolAttribute{
                MarkdownDescription: "When enabled, alerts from different monitors will be grouped into separate episodes. When disabled, alerts from any monitor can be grouped together..",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "group_by_severity": schema.BoolAttribute{
                MarkdownDescription: "When enabled, alerts with different severities will be grouped into separate episodes. When disabled, alerts of any severity can be grouped together..",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "group_by_alert_title": schema.BoolAttribute{
                MarkdownDescription: "When enabled, alerts with different titles will be grouped into separate episodes. When disabled, alerts with any title can be grouped together..",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "group_by_alert_labels": schema.BoolAttribute{
                MarkdownDescription: "When enabled, alerts with different sets of labels will be grouped into separate episodes (exact set match). When disabled, alert labels are ignored for grouping..",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "group_by_monitor_labels": schema.BoolAttribute{
                MarkdownDescription: "When enabled, alerts whose monitors have different sets of labels will be grouped into separate episodes (exact set match). When disabled, monitor labels are ignored for grouping..",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "enable_time_window": schema.BoolAttribute{
                MarkdownDescription: "Enable time-based grouping. When enabled, alerts are grouped within the specified time window. When disabled, all matching alerts are grouped into a single ongoing episode regardless of time..",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "time_window_minutes": schema.NumberAttribute{
                MarkdownDescription: "Rolling time window in minutes. Alerts are grouped if they arrive within this gap from the last alert..",
                Optional: true,
                Computed: true,
                Default: numberdefault.StaticBigFloat(big.NewFloat(60)),
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "group_by_fields": schema.StringAttribute{
                MarkdownDescription: "JSON object defining the fields to group alerts by (e.g., monitorId, severity).",
                CustomType: JSONSubsetType{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: []validator.String{
                    JSONEnvelopeValidator(),
                },
            },
            "episode_title_template": schema.StringAttribute{
                MarkdownDescription: "Template for generating episode titles. Supports placeholders like {{alertSeverity}}, {{monitorName}}, {{alertTitle}}, {{alertDescription}}.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "episode_description_template": schema.StringAttribute{
                MarkdownDescription: "Template for generating episode descriptions. Supports placeholders like {{alertSeverity}}, {{monitorName}}, {{alertTitle}}, {{alertDescription}}.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "enable_resolve_delay": schema.BoolAttribute{
                MarkdownDescription: "Enable grace period before auto-resolving episode after all alerts resolve. Helps prevent rapid state changes during alert flapping..",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "resolve_delay_minutes": schema.NumberAttribute{
                MarkdownDescription: "Grace period in minutes before auto-resolving an episode after all alerts are resolved.",
                Optional: true,
                Computed: true,
                Default: numberdefault.StaticBigFloat(big.NewFloat(0)),
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "enable_reopen_window": schema.BoolAttribute{
                MarkdownDescription: "Enable reopening recently resolved episodes instead of creating new ones. Useful when related issues recur shortly after resolution..",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "reopen_window_minutes": schema.NumberAttribute{
                MarkdownDescription: "Time window in minutes to reopen a recently resolved episode instead of creating a new one.",
                Optional: true,
                Computed: true,
                Default: numberdefault.StaticBigFloat(big.NewFloat(0)),
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "enable_inactivity_timeout": schema.BoolAttribute{
                MarkdownDescription: "Enable auto-resolving episodes after a period of inactivity. Helps automatically close episodes when no new alerts arrive..",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "inactivity_timeout_minutes": schema.NumberAttribute{
                MarkdownDescription: "Time in minutes after which an inactive episode will be auto-resolved.",
                Optional: true,
                Computed: true,
                Default: numberdefault.StaticBigFloat(big.NewFloat(60)),
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "on_call_duty_policies": schema.SetAttribute{
                MarkdownDescription: "List of on-call duty policies to execute for episodes created by this rule..",
                Optional: true,
                Computed: true,
                ElementType: types.StringType,
                PlanModifiers: []planmodifier.Set{
                    setplanmodifier.UseStateForUnknown(),
                },
            },
            "default_assign_to_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "default_assign_to_team_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "episode_labels": schema.SetAttribute{
                MarkdownDescription: "Labels to automatically apply to episodes created by this rule..",
                Optional: true,
                Computed: true,
                ElementType: types.StringType,
                PlanModifiers: []planmodifier.Set{
                    setplanmodifier.UseStateForUnknown(),
                },
            },
            "episode_owner_users": schema.SetAttribute{
                MarkdownDescription: "Users to automatically add as owners to episodes created by this rule..",
                Optional: true,
                Computed: true,
                ElementType: types.StringType,
                PlanModifiers: []planmodifier.Set{
                    setplanmodifier.UseStateForUnknown(),
                },
            },
            "episode_owner_teams": schema.SetAttribute{
                MarkdownDescription: "Teams to automatically add as owners to episodes created by this rule..",
                Optional: true,
                Computed: true,
                ElementType: types.StringType,
                PlanModifiers: []planmodifier.Set{
                    setplanmodifier.UseStateForUnknown(),
                },
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                    stringplanmodifier.RequiresReplace(),
                },
            },
            "created_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "updated_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "deleted_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "version": schema.NumberAttribute{
                MarkdownDescription: "Object version",
                Computed: true,
            },
        },
    }
}

func (r *AlertGroupingRuleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
    // Prevent panic if the provider has not been configured.
    if req.ProviderData == nil {
        return
    }

    client, ok := req.ProviderData.(*Client)

    if !ok {
        resp.Diagnostics.AddError(
            "Unexpected Resource Configure Type",
            fmt.Sprintf("Expected *Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
        )

        return
    }

    r.client = client
}


func (r *AlertGroupingRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var data AlertGroupingRuleResourceModel

    // Read Terraform plan data into the model
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }



    // Create API request body. Unset (null/unknown) optional fields are
    // omitted so server-side defaults apply instead of being overwritten
    // with zero values.
    alertGroupingRuleRequest := map[string]interface{}{
        "data": map[string]interface{}{},
    }
    requestDataMap := alertGroupingRuleRequest["data"].(map[string]interface{})

    if !data.Name.IsNull() && !data.Name.IsUnknown() {
        requestDataMap["name"] = data.Name.ValueString()
    }
    if !data.Description.IsNull() && !data.Description.IsUnknown() {
        requestDataMap["description"] = data.Description.ValueString()
    }
    if !data.Priority.IsNull() && !data.Priority.IsUnknown() {
        requestDataMap["priority"] = r.bigFloatToFloat64(data.Priority.ValueBigFloat())
    }
    if !data.IsEnabled.IsNull() && !data.IsEnabled.IsUnknown() {
        requestDataMap["isEnabled"] = data.IsEnabled.ValueBool()
    }
    if parsedMatchCriteria := r.parseJSONField(data.MatchCriteria); parsedMatchCriteria != nil {
        requestDataMap["matchCriteria"] = parsedMatchCriteria
    }
    if !data.Monitors.IsNull() && !data.Monitors.IsUnknown() {
        requestDataMap["monitors"] = r.convertTerraformSetToInterface(data.Monitors)
    }
    if !data.AlertSeverities.IsNull() && !data.AlertSeverities.IsUnknown() {
        requestDataMap["alertSeverities"] = r.convertTerraformSetToInterface(data.AlertSeverities)
    }
    if !data.AlertLabels.IsNull() && !data.AlertLabels.IsUnknown() {
        requestDataMap["alertLabels"] = r.convertTerraformSetToInterface(data.AlertLabels)
    }
    if !data.MonitorLabels.IsNull() && !data.MonitorLabels.IsUnknown() {
        requestDataMap["monitorLabels"] = r.convertTerraformSetToInterface(data.MonitorLabels)
    }
    if !data.AlertTitlePattern.IsNull() && !data.AlertTitlePattern.IsUnknown() {
        requestDataMap["alertTitlePattern"] = data.AlertTitlePattern.ValueString()
    }
    if !data.AlertDescriptionPattern.IsNull() && !data.AlertDescriptionPattern.IsUnknown() {
        requestDataMap["alertDescriptionPattern"] = data.AlertDescriptionPattern.ValueString()
    }
    if !data.MonitorNamePattern.IsNull() && !data.MonitorNamePattern.IsUnknown() {
        requestDataMap["monitorNamePattern"] = data.MonitorNamePattern.ValueString()
    }
    if !data.MonitorDescriptionPattern.IsNull() && !data.MonitorDescriptionPattern.IsUnknown() {
        requestDataMap["monitorDescriptionPattern"] = data.MonitorDescriptionPattern.ValueString()
    }
    if !data.GroupByMonitor.IsNull() && !data.GroupByMonitor.IsUnknown() {
        requestDataMap["groupByMonitor"] = data.GroupByMonitor.ValueBool()
    }
    if !data.GroupBySeverity.IsNull() && !data.GroupBySeverity.IsUnknown() {
        requestDataMap["groupBySeverity"] = data.GroupBySeverity.ValueBool()
    }
    if !data.GroupByAlertTitle.IsNull() && !data.GroupByAlertTitle.IsUnknown() {
        requestDataMap["groupByAlertTitle"] = data.GroupByAlertTitle.ValueBool()
    }
    if !data.GroupByAlertLabels.IsNull() && !data.GroupByAlertLabels.IsUnknown() {
        requestDataMap["groupByAlertLabels"] = data.GroupByAlertLabels.ValueBool()
    }
    if !data.GroupByMonitorLabels.IsNull() && !data.GroupByMonitorLabels.IsUnknown() {
        requestDataMap["groupByMonitorLabels"] = data.GroupByMonitorLabels.ValueBool()
    }
    if !data.EnableTimeWindow.IsNull() && !data.EnableTimeWindow.IsUnknown() {
        requestDataMap["enableTimeWindow"] = data.EnableTimeWindow.ValueBool()
    }
    if !data.TimeWindowMinutes.IsNull() && !data.TimeWindowMinutes.IsUnknown() {
        requestDataMap["timeWindowMinutes"] = r.bigFloatToFloat64(data.TimeWindowMinutes.ValueBigFloat())
    }
    if parsedGroupByFields := r.parseJSONField(data.GroupByFields); parsedGroupByFields != nil {
        requestDataMap["groupByFields"] = parsedGroupByFields
    }
    if !data.EpisodeTitleTemplate.IsNull() && !data.EpisodeTitleTemplate.IsUnknown() {
        requestDataMap["episodeTitleTemplate"] = data.EpisodeTitleTemplate.ValueString()
    }
    if !data.EpisodeDescriptionTemplate.IsNull() && !data.EpisodeDescriptionTemplate.IsUnknown() {
        requestDataMap["episodeDescriptionTemplate"] = data.EpisodeDescriptionTemplate.ValueString()
    }
    if !data.EnableResolveDelay.IsNull() && !data.EnableResolveDelay.IsUnknown() {
        requestDataMap["enableResolveDelay"] = data.EnableResolveDelay.ValueBool()
    }
    if !data.ResolveDelayMinutes.IsNull() && !data.ResolveDelayMinutes.IsUnknown() {
        requestDataMap["resolveDelayMinutes"] = r.bigFloatToFloat64(data.ResolveDelayMinutes.ValueBigFloat())
    }
    if !data.EnableReopenWindow.IsNull() && !data.EnableReopenWindow.IsUnknown() {
        requestDataMap["enableReopenWindow"] = data.EnableReopenWindow.ValueBool()
    }
    if !data.ReopenWindowMinutes.IsNull() && !data.ReopenWindowMinutes.IsUnknown() {
        requestDataMap["reopenWindowMinutes"] = r.bigFloatToFloat64(data.ReopenWindowMinutes.ValueBigFloat())
    }
    if !data.EnableInactivityTimeout.IsNull() && !data.EnableInactivityTimeout.IsUnknown() {
        requestDataMap["enableInactivityTimeout"] = data.EnableInactivityTimeout.ValueBool()
    }
    if !data.InactivityTimeoutMinutes.IsNull() && !data.InactivityTimeoutMinutes.IsUnknown() {
        requestDataMap["inactivityTimeoutMinutes"] = r.bigFloatToFloat64(data.InactivityTimeoutMinutes.ValueBigFloat())
    }
    if !data.OnCallDutyPolicies.IsNull() && !data.OnCallDutyPolicies.IsUnknown() {
        requestDataMap["onCallDutyPolicies"] = r.convertTerraformSetToInterface(data.OnCallDutyPolicies)
    }
    if !data.DefaultAssignToUserId.IsNull() && !data.DefaultAssignToUserId.IsUnknown() {
        requestDataMap["defaultAssignToUserId"] = data.DefaultAssignToUserId.ValueString()
    }
    if !data.DefaultAssignToTeamId.IsNull() && !data.DefaultAssignToTeamId.IsUnknown() {
        requestDataMap["defaultAssignToTeamId"] = data.DefaultAssignToTeamId.ValueString()
    }
    if !data.EpisodeLabels.IsNull() && !data.EpisodeLabels.IsUnknown() {
        requestDataMap["episodeLabels"] = r.convertTerraformSetToInterface(data.EpisodeLabels)
    }
    if !data.EpisodeOwnerUsers.IsNull() && !data.EpisodeOwnerUsers.IsUnknown() {
        requestDataMap["episodeOwnerUsers"] = r.convertTerraformSetToInterface(data.EpisodeOwnerUsers)
    }
    if !data.EpisodeOwnerTeams.IsNull() && !data.EpisodeOwnerTeams.IsUnknown() {
        requestDataMap["episodeOwnerTeams"] = r.convertTerraformSetToInterface(data.EpisodeOwnerTeams)
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        requestDataMap["createdByUserId"] = data.CreatedByUserId.ValueString()
    }

    // Make API call
    httpResp, err := r.client.Post(ctx, "/alert-grouping-rule", alertGroupingRuleRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create alert_grouping_rule, got error: %s", err))
        return
    }

    var alertGroupingRuleResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &alertGroupingRuleResponse)
    if err != nil {
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to create alert_grouping_rule: %s", err))
        return
    }

    // Extract the new resource id from the create response.
    createdId := ""
    if wrapper, ok := alertGroupingRuleResponse["data"].(map[string]interface{}); ok {
        if val, ok := wrapper["_id"].(string); ok {
            createdId = val
        }
    } else if val, ok := alertGroupingRuleResponse["_id"].(string); ok {
        createdId = val
    }
    if createdId == "" {
        resp.Diagnostics.AddError("OneUptime API Error", "Create response for alert_grouping_rule did not contain an id. This is a bug in the provider or the API; please report it.")
        return
    }
    data.Id = types.StringValue(createdId)

    /*
     * The server has committed the row. Persist what we know to state BEFORE
     * the read-back: if the read-back fails and we return without setting
     * state, Terraform never learns the resource exists and the created
     * alert_grouping_rule is orphaned server-side — never refreshed, never
     * destroyed. Delete already refuses to drop state on failure for the
     * same reason; Create must not either.
     */
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Re-read the resource so state reflects server-normalized values.
    selectParam := map[string]interface{}{
        "projectId": true,
        "name": true,
        "description": true,
        "priority": true,
        "isEnabled": true,
        "matchCriteria": true,
        "monitors": true,
        "alertSeverities": true,
        "alertLabels": true,
        "monitorLabels": true,
        "alertTitlePattern": true,
        "alertDescriptionPattern": true,
        "monitorNamePattern": true,
        "monitorDescriptionPattern": true,
        "groupByMonitor": true,
        "groupBySeverity": true,
        "groupByAlertTitle": true,
        "groupByAlertLabels": true,
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
        "createdByUserId": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect(ctx, "/alert-grouping-rule/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        /*
         * State already owns the id, so the resource is tracked and the next
         * refresh reconciles the remaining attributes. Warn rather than
         * error: erroring here would strand a real resource.
         */
        resp.Diagnostics.AddWarning("Read After Create Failed", fmt.Sprintf("Created alert_grouping_rule but could not read it back; state is incomplete until the next refresh: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddWarning("Read After Create Failed", fmt.Sprintf("Created alert_grouping_rule but could not parse the read-back response; state is incomplete until the next refresh: %s", err))
        return
    }

    // Update the model with the authoritative read response
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := readResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = readResponse
    }

    if obj, ok := dataMap["projectId"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok {
            data.ProjectId = types.StringValue(val)
        } else {
            data.ProjectId = types.StringNull()
        }
    } else if val, ok := dataMap["projectId"].(string); ok {
        data.ProjectId = types.StringValue(val)
    } else {
        data.ProjectId = types.StringNull()
    }
    if obj, ok := dataMap["name"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Name = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Name = types.StringValue(string(jsonBytes))
            } else {
                data.Name = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Name = types.StringValue(string(jsonBytes))
            } else {
                data.Name = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Name = types.StringValue(string(jsonBytes))
        } else {
            data.Name = types.StringNull()
        }
    } else if val, ok := dataMap["name"].(string); ok {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if obj, ok := dataMap["description"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Description = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Description = types.StringValue(string(jsonBytes))
            } else {
                data.Description = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Description = types.StringValue(string(jsonBytes))
            } else {
                data.Description = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Description = types.StringValue(string(jsonBytes))
        } else {
            data.Description = types.StringNull()
        }
    } else if val, ok := dataMap["description"].(string); ok {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if val, ok := dataMap["priority"].(float64); ok {
        data.Priority = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["priority"].(int); ok {
        data.Priority = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["priority"].(int64); ok {
        data.Priority = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["priority"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.Priority = types.NumberValue(big.NewFloat(val))
        } else {
            data.Priority = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.Priority = types.NumberNull()
    }
    if val, ok := dataMap["isEnabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    }
    if obj, ok := dataMap["matchCriteria"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MatchCriteria = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.MatchCriteria = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.MatchCriteria = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.MatchCriteria = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.MatchCriteria = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.MatchCriteria = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.MatchCriteria = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.MatchCriteria = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.MatchCriteria = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["matchCriteria"].(string); ok {
        data.MatchCriteria = NewJSONSubsetValue(val)
    } else {
        data.MatchCriteria = NewJSONSubsetNull()
    }
    if val, ok := dataMap["monitors"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.Monitors = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.Monitors = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if val, ok := dataMap["alertSeverities"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.AlertSeverities = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.AlertSeverities = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if val, ok := dataMap["alertLabels"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.AlertLabels = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.AlertLabels = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if val, ok := dataMap["monitorLabels"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.MonitorLabels = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.MonitorLabels = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if obj, ok := dataMap["alertTitlePattern"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AlertTitlePattern = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AlertTitlePattern = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AlertTitlePattern = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AlertTitlePattern = types.StringValue(string(jsonBytes))
            } else {
                data.AlertTitlePattern = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AlertTitlePattern = types.StringValue(string(jsonBytes))
            } else {
                data.AlertTitlePattern = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AlertTitlePattern = types.StringValue(string(jsonBytes))
        } else {
            data.AlertTitlePattern = types.StringNull()
        }
    } else if val, ok := dataMap["alertTitlePattern"].(string); ok {
        data.AlertTitlePattern = types.StringValue(val)
    } else {
        data.AlertTitlePattern = types.StringNull()
    }
    if obj, ok := dataMap["alertDescriptionPattern"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AlertDescriptionPattern = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AlertDescriptionPattern = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AlertDescriptionPattern = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AlertDescriptionPattern = types.StringValue(string(jsonBytes))
            } else {
                data.AlertDescriptionPattern = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AlertDescriptionPattern = types.StringValue(string(jsonBytes))
            } else {
                data.AlertDescriptionPattern = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AlertDescriptionPattern = types.StringValue(string(jsonBytes))
        } else {
            data.AlertDescriptionPattern = types.StringNull()
        }
    } else if val, ok := dataMap["alertDescriptionPattern"].(string); ok {
        data.AlertDescriptionPattern = types.StringValue(val)
    } else {
        data.AlertDescriptionPattern = types.StringNull()
    }
    if obj, ok := dataMap["monitorNamePattern"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MonitorNamePattern = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.MonitorNamePattern = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.MonitorNamePattern = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.MonitorNamePattern = types.StringValue(string(jsonBytes))
            } else {
                data.MonitorNamePattern = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.MonitorNamePattern = types.StringValue(string(jsonBytes))
            } else {
                data.MonitorNamePattern = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.MonitorNamePattern = types.StringValue(string(jsonBytes))
        } else {
            data.MonitorNamePattern = types.StringNull()
        }
    } else if val, ok := dataMap["monitorNamePattern"].(string); ok {
        data.MonitorNamePattern = types.StringValue(val)
    } else {
        data.MonitorNamePattern = types.StringNull()
    }
    if obj, ok := dataMap["monitorDescriptionPattern"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MonitorDescriptionPattern = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.MonitorDescriptionPattern = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.MonitorDescriptionPattern = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.MonitorDescriptionPattern = types.StringValue(string(jsonBytes))
            } else {
                data.MonitorDescriptionPattern = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.MonitorDescriptionPattern = types.StringValue(string(jsonBytes))
            } else {
                data.MonitorDescriptionPattern = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.MonitorDescriptionPattern = types.StringValue(string(jsonBytes))
        } else {
            data.MonitorDescriptionPattern = types.StringNull()
        }
    } else if val, ok := dataMap["monitorDescriptionPattern"].(string); ok {
        data.MonitorDescriptionPattern = types.StringValue(val)
    } else {
        data.MonitorDescriptionPattern = types.StringNull()
    }
    if val, ok := dataMap["groupByMonitor"].(bool); ok {
        data.GroupByMonitor = types.BoolValue(val)
    }
    if val, ok := dataMap["groupBySeverity"].(bool); ok {
        data.GroupBySeverity = types.BoolValue(val)
    }
    if val, ok := dataMap["groupByAlertTitle"].(bool); ok {
        data.GroupByAlertTitle = types.BoolValue(val)
    }
    if val, ok := dataMap["groupByAlertLabels"].(bool); ok {
        data.GroupByAlertLabels = types.BoolValue(val)
    }
    if val, ok := dataMap["groupByMonitorLabels"].(bool); ok {
        data.GroupByMonitorLabels = types.BoolValue(val)
    }
    if val, ok := dataMap["enableTimeWindow"].(bool); ok {
        data.EnableTimeWindow = types.BoolValue(val)
    }
    if val, ok := dataMap["timeWindowMinutes"].(float64); ok {
        data.TimeWindowMinutes = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["timeWindowMinutes"].(int); ok {
        data.TimeWindowMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["timeWindowMinutes"].(int64); ok {
        data.TimeWindowMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["timeWindowMinutes"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.TimeWindowMinutes = types.NumberValue(big.NewFloat(val))
        } else {
            data.TimeWindowMinutes = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.TimeWindowMinutes = types.NumberNull()
    }
    if obj, ok := dataMap["groupByFields"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.GroupByFields = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.GroupByFields = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.GroupByFields = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.GroupByFields = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.GroupByFields = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.GroupByFields = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.GroupByFields = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.GroupByFields = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.GroupByFields = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["groupByFields"].(string); ok {
        data.GroupByFields = NewJSONSubsetValue(val)
    } else {
        data.GroupByFields = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["episodeTitleTemplate"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EpisodeTitleTemplate = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.EpisodeTitleTemplate = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.EpisodeTitleTemplate = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.EpisodeTitleTemplate = types.StringValue(string(jsonBytes))
            } else {
                data.EpisodeTitleTemplate = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.EpisodeTitleTemplate = types.StringValue(string(jsonBytes))
            } else {
                data.EpisodeTitleTemplate = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.EpisodeTitleTemplate = types.StringValue(string(jsonBytes))
        } else {
            data.EpisodeTitleTemplate = types.StringNull()
        }
    } else if val, ok := dataMap["episodeTitleTemplate"].(string); ok {
        data.EpisodeTitleTemplate = types.StringValue(val)
    } else {
        data.EpisodeTitleTemplate = types.StringNull()
    }
    if obj, ok := dataMap["episodeDescriptionTemplate"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EpisodeDescriptionTemplate = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.EpisodeDescriptionTemplate = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.EpisodeDescriptionTemplate = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.EpisodeDescriptionTemplate = types.StringValue(string(jsonBytes))
            } else {
                data.EpisodeDescriptionTemplate = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.EpisodeDescriptionTemplate = types.StringValue(string(jsonBytes))
            } else {
                data.EpisodeDescriptionTemplate = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.EpisodeDescriptionTemplate = types.StringValue(string(jsonBytes))
        } else {
            data.EpisodeDescriptionTemplate = types.StringNull()
        }
    } else if val, ok := dataMap["episodeDescriptionTemplate"].(string); ok {
        data.EpisodeDescriptionTemplate = types.StringValue(val)
    } else {
        data.EpisodeDescriptionTemplate = types.StringNull()
    }
    if val, ok := dataMap["enableResolveDelay"].(bool); ok {
        data.EnableResolveDelay = types.BoolValue(val)
    }
    if val, ok := dataMap["resolveDelayMinutes"].(float64); ok {
        data.ResolveDelayMinutes = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["resolveDelayMinutes"].(int); ok {
        data.ResolveDelayMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["resolveDelayMinutes"].(int64); ok {
        data.ResolveDelayMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["resolveDelayMinutes"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.ResolveDelayMinutes = types.NumberValue(big.NewFloat(val))
        } else {
            data.ResolveDelayMinutes = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.ResolveDelayMinutes = types.NumberNull()
    }
    if val, ok := dataMap["enableReopenWindow"].(bool); ok {
        data.EnableReopenWindow = types.BoolValue(val)
    }
    if val, ok := dataMap["reopenWindowMinutes"].(float64); ok {
        data.ReopenWindowMinutes = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["reopenWindowMinutes"].(int); ok {
        data.ReopenWindowMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["reopenWindowMinutes"].(int64); ok {
        data.ReopenWindowMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["reopenWindowMinutes"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.ReopenWindowMinutes = types.NumberValue(big.NewFloat(val))
        } else {
            data.ReopenWindowMinutes = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.ReopenWindowMinutes = types.NumberNull()
    }
    if val, ok := dataMap["enableInactivityTimeout"].(bool); ok {
        data.EnableInactivityTimeout = types.BoolValue(val)
    }
    if val, ok := dataMap["inactivityTimeoutMinutes"].(float64); ok {
        data.InactivityTimeoutMinutes = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["inactivityTimeoutMinutes"].(int); ok {
        data.InactivityTimeoutMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["inactivityTimeoutMinutes"].(int64); ok {
        data.InactivityTimeoutMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["inactivityTimeoutMinutes"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.InactivityTimeoutMinutes = types.NumberValue(big.NewFloat(val))
        } else {
            data.InactivityTimeoutMinutes = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.InactivityTimeoutMinutes = types.NumberNull()
    }
    if val, ok := dataMap["onCallDutyPolicies"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.OnCallDutyPolicies = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.OnCallDutyPolicies = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if obj, ok := dataMap["defaultAssignToUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DefaultAssignToUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DefaultAssignToUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DefaultAssignToUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DefaultAssignToUserId = types.StringValue(string(jsonBytes))
            } else {
                data.DefaultAssignToUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DefaultAssignToUserId = types.StringValue(string(jsonBytes))
            } else {
                data.DefaultAssignToUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DefaultAssignToUserId = types.StringValue(string(jsonBytes))
        } else {
            data.DefaultAssignToUserId = types.StringNull()
        }
    } else if val, ok := dataMap["defaultAssignToUserId"].(string); ok {
        data.DefaultAssignToUserId = types.StringValue(val)
    } else {
        data.DefaultAssignToUserId = types.StringNull()
    }
    if obj, ok := dataMap["defaultAssignToTeamId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DefaultAssignToTeamId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DefaultAssignToTeamId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DefaultAssignToTeamId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DefaultAssignToTeamId = types.StringValue(string(jsonBytes))
            } else {
                data.DefaultAssignToTeamId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DefaultAssignToTeamId = types.StringValue(string(jsonBytes))
            } else {
                data.DefaultAssignToTeamId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DefaultAssignToTeamId = types.StringValue(string(jsonBytes))
        } else {
            data.DefaultAssignToTeamId = types.StringNull()
        }
    } else if val, ok := dataMap["defaultAssignToTeamId"].(string); ok {
        data.DefaultAssignToTeamId = types.StringValue(val)
    } else {
        data.DefaultAssignToTeamId = types.StringNull()
    }
    if val, ok := dataMap["episodeLabels"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.EpisodeLabels = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.EpisodeLabels = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if val, ok := dataMap["episodeOwnerUsers"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.EpisodeOwnerUsers = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.EpisodeOwnerUsers = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if val, ok := dataMap["episodeOwnerTeams"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.EpisodeOwnerTeams = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.EpisodeOwnerTeams = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if obj, ok := dataMap["createdByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["createdByUserId"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if obj, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.CreatedAt = NewRFC3339Value(val)
        } else {
            data.CreatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["createdAt"].(string); ok && val != "" {
        data.CreatedAt = NewRFC3339Value(val)
    } else {
        data.CreatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.UpdatedAt = NewRFC3339Value(val)
        } else {
            data.UpdatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["updatedAt"].(string); ok && val != "" {
        data.UpdatedAt = NewRFC3339Value(val)
    } else {
        data.UpdatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["deletedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.DeletedAt = NewRFC3339Value(val)
        } else {
            data.DeletedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["deletedAt"].(string); ok && val != "" {
        data.DeletedAt = NewRFC3339Value(val)
    } else {
        data.DeletedAt = NewRFC3339Null()
    }
    if val, ok := dataMap["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["version"].(int); ok {
        data.Version = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["version"].(int64); ok {
        data.Version = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["version"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.Version = types.NumberValue(big.NewFloat(val))
        } else {
            data.Version = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.Version = types.NumberNull()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    // The read response is authoritative, but never let it clobber the id we just received.
    data.Id = types.StringValue(createdId)

    // Write logs using the tflog package
    tflog.Trace(ctx, "created a resource")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AlertGroupingRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var data AlertGroupingRuleResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create select parameter to get full object
    selectParam := map[string]interface{}{
        "projectId": true,
        "name": true,
        "description": true,
        "priority": true,
        "isEnabled": true,
        "matchCriteria": true,
        "monitors": true,
        "alertSeverities": true,
        "alertLabels": true,
        "monitorLabels": true,
        "alertTitlePattern": true,
        "alertDescriptionPattern": true,
        "monitorNamePattern": true,
        "monitorDescriptionPattern": true,
        "groupByMonitor": true,
        "groupBySeverity": true,
        "groupByAlertTitle": true,
        "groupByAlertLabels": true,
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
        "createdByUserId": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "_id": true,
    }

    // Make API call with select parameter
    httpResp, err := r.client.PostWithSelect(ctx, "/alert-grouping-rule/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read alert_grouping_rule, got error: %s", err))
        return
    }

    if httpResp.StatusCode == http.StatusNotFound {
        resp.State.RemoveResource(ctx)
        return
    }

    var alertGroupingRuleResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &alertGroupingRuleResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse alert_grouping_rule response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := alertGroupingRuleResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = alertGroupingRuleResponse
    }

    if obj, ok := dataMap["projectId"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok {
            data.ProjectId = types.StringValue(val)
        } else {
            data.ProjectId = types.StringNull()
        }
    } else if val, ok := dataMap["projectId"].(string); ok {
        data.ProjectId = types.StringValue(val)
    } else {
        data.ProjectId = types.StringNull()
    }
    if obj, ok := dataMap["name"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Name = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Name = types.StringValue(string(jsonBytes))
            } else {
                data.Name = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Name = types.StringValue(string(jsonBytes))
            } else {
                data.Name = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Name = types.StringValue(string(jsonBytes))
        } else {
            data.Name = types.StringNull()
        }
    } else if val, ok := dataMap["name"].(string); ok {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if obj, ok := dataMap["description"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Description = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Description = types.StringValue(string(jsonBytes))
            } else {
                data.Description = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Description = types.StringValue(string(jsonBytes))
            } else {
                data.Description = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Description = types.StringValue(string(jsonBytes))
        } else {
            data.Description = types.StringNull()
        }
    } else if val, ok := dataMap["description"].(string); ok {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if val, ok := dataMap["priority"].(float64); ok {
        data.Priority = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["priority"].(int); ok {
        data.Priority = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["priority"].(int64); ok {
        data.Priority = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["priority"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.Priority = types.NumberValue(big.NewFloat(val))
        } else {
            data.Priority = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.Priority = types.NumberNull()
    }
    if val, ok := dataMap["isEnabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    }
    if obj, ok := dataMap["matchCriteria"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MatchCriteria = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.MatchCriteria = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.MatchCriteria = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.MatchCriteria = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.MatchCriteria = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.MatchCriteria = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.MatchCriteria = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.MatchCriteria = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.MatchCriteria = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["matchCriteria"].(string); ok {
        data.MatchCriteria = NewJSONSubsetValue(val)
    } else {
        data.MatchCriteria = NewJSONSubsetNull()
    }
    if val, ok := dataMap["monitors"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.Monitors = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.Monitors = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if val, ok := dataMap["alertSeverities"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.AlertSeverities = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.AlertSeverities = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if val, ok := dataMap["alertLabels"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.AlertLabels = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.AlertLabels = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if val, ok := dataMap["monitorLabels"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.MonitorLabels = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.MonitorLabels = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if obj, ok := dataMap["alertTitlePattern"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AlertTitlePattern = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AlertTitlePattern = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AlertTitlePattern = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AlertTitlePattern = types.StringValue(string(jsonBytes))
            } else {
                data.AlertTitlePattern = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AlertTitlePattern = types.StringValue(string(jsonBytes))
            } else {
                data.AlertTitlePattern = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AlertTitlePattern = types.StringValue(string(jsonBytes))
        } else {
            data.AlertTitlePattern = types.StringNull()
        }
    } else if val, ok := dataMap["alertTitlePattern"].(string); ok {
        data.AlertTitlePattern = types.StringValue(val)
    } else {
        data.AlertTitlePattern = types.StringNull()
    }
    if obj, ok := dataMap["alertDescriptionPattern"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AlertDescriptionPattern = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AlertDescriptionPattern = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AlertDescriptionPattern = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AlertDescriptionPattern = types.StringValue(string(jsonBytes))
            } else {
                data.AlertDescriptionPattern = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AlertDescriptionPattern = types.StringValue(string(jsonBytes))
            } else {
                data.AlertDescriptionPattern = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AlertDescriptionPattern = types.StringValue(string(jsonBytes))
        } else {
            data.AlertDescriptionPattern = types.StringNull()
        }
    } else if val, ok := dataMap["alertDescriptionPattern"].(string); ok {
        data.AlertDescriptionPattern = types.StringValue(val)
    } else {
        data.AlertDescriptionPattern = types.StringNull()
    }
    if obj, ok := dataMap["monitorNamePattern"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MonitorNamePattern = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.MonitorNamePattern = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.MonitorNamePattern = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.MonitorNamePattern = types.StringValue(string(jsonBytes))
            } else {
                data.MonitorNamePattern = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.MonitorNamePattern = types.StringValue(string(jsonBytes))
            } else {
                data.MonitorNamePattern = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.MonitorNamePattern = types.StringValue(string(jsonBytes))
        } else {
            data.MonitorNamePattern = types.StringNull()
        }
    } else if val, ok := dataMap["monitorNamePattern"].(string); ok {
        data.MonitorNamePattern = types.StringValue(val)
    } else {
        data.MonitorNamePattern = types.StringNull()
    }
    if obj, ok := dataMap["monitorDescriptionPattern"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MonitorDescriptionPattern = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.MonitorDescriptionPattern = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.MonitorDescriptionPattern = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.MonitorDescriptionPattern = types.StringValue(string(jsonBytes))
            } else {
                data.MonitorDescriptionPattern = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.MonitorDescriptionPattern = types.StringValue(string(jsonBytes))
            } else {
                data.MonitorDescriptionPattern = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.MonitorDescriptionPattern = types.StringValue(string(jsonBytes))
        } else {
            data.MonitorDescriptionPattern = types.StringNull()
        }
    } else if val, ok := dataMap["monitorDescriptionPattern"].(string); ok {
        data.MonitorDescriptionPattern = types.StringValue(val)
    } else {
        data.MonitorDescriptionPattern = types.StringNull()
    }
    if val, ok := dataMap["groupByMonitor"].(bool); ok {
        data.GroupByMonitor = types.BoolValue(val)
    }
    if val, ok := dataMap["groupBySeverity"].(bool); ok {
        data.GroupBySeverity = types.BoolValue(val)
    }
    if val, ok := dataMap["groupByAlertTitle"].(bool); ok {
        data.GroupByAlertTitle = types.BoolValue(val)
    }
    if val, ok := dataMap["groupByAlertLabels"].(bool); ok {
        data.GroupByAlertLabels = types.BoolValue(val)
    }
    if val, ok := dataMap["groupByMonitorLabels"].(bool); ok {
        data.GroupByMonitorLabels = types.BoolValue(val)
    }
    if val, ok := dataMap["enableTimeWindow"].(bool); ok {
        data.EnableTimeWindow = types.BoolValue(val)
    }
    if val, ok := dataMap["timeWindowMinutes"].(float64); ok {
        data.TimeWindowMinutes = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["timeWindowMinutes"].(int); ok {
        data.TimeWindowMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["timeWindowMinutes"].(int64); ok {
        data.TimeWindowMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["timeWindowMinutes"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.TimeWindowMinutes = types.NumberValue(big.NewFloat(val))
        } else {
            data.TimeWindowMinutes = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.TimeWindowMinutes = types.NumberNull()
    }
    if obj, ok := dataMap["groupByFields"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.GroupByFields = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.GroupByFields = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.GroupByFields = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.GroupByFields = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.GroupByFields = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.GroupByFields = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.GroupByFields = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.GroupByFields = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.GroupByFields = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["groupByFields"].(string); ok {
        data.GroupByFields = NewJSONSubsetValue(val)
    } else {
        data.GroupByFields = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["episodeTitleTemplate"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EpisodeTitleTemplate = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.EpisodeTitleTemplate = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.EpisodeTitleTemplate = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.EpisodeTitleTemplate = types.StringValue(string(jsonBytes))
            } else {
                data.EpisodeTitleTemplate = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.EpisodeTitleTemplate = types.StringValue(string(jsonBytes))
            } else {
                data.EpisodeTitleTemplate = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.EpisodeTitleTemplate = types.StringValue(string(jsonBytes))
        } else {
            data.EpisodeTitleTemplate = types.StringNull()
        }
    } else if val, ok := dataMap["episodeTitleTemplate"].(string); ok {
        data.EpisodeTitleTemplate = types.StringValue(val)
    } else {
        data.EpisodeTitleTemplate = types.StringNull()
    }
    if obj, ok := dataMap["episodeDescriptionTemplate"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EpisodeDescriptionTemplate = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.EpisodeDescriptionTemplate = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.EpisodeDescriptionTemplate = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.EpisodeDescriptionTemplate = types.StringValue(string(jsonBytes))
            } else {
                data.EpisodeDescriptionTemplate = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.EpisodeDescriptionTemplate = types.StringValue(string(jsonBytes))
            } else {
                data.EpisodeDescriptionTemplate = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.EpisodeDescriptionTemplate = types.StringValue(string(jsonBytes))
        } else {
            data.EpisodeDescriptionTemplate = types.StringNull()
        }
    } else if val, ok := dataMap["episodeDescriptionTemplate"].(string); ok {
        data.EpisodeDescriptionTemplate = types.StringValue(val)
    } else {
        data.EpisodeDescriptionTemplate = types.StringNull()
    }
    if val, ok := dataMap["enableResolveDelay"].(bool); ok {
        data.EnableResolveDelay = types.BoolValue(val)
    }
    if val, ok := dataMap["resolveDelayMinutes"].(float64); ok {
        data.ResolveDelayMinutes = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["resolveDelayMinutes"].(int); ok {
        data.ResolveDelayMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["resolveDelayMinutes"].(int64); ok {
        data.ResolveDelayMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["resolveDelayMinutes"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.ResolveDelayMinutes = types.NumberValue(big.NewFloat(val))
        } else {
            data.ResolveDelayMinutes = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.ResolveDelayMinutes = types.NumberNull()
    }
    if val, ok := dataMap["enableReopenWindow"].(bool); ok {
        data.EnableReopenWindow = types.BoolValue(val)
    }
    if val, ok := dataMap["reopenWindowMinutes"].(float64); ok {
        data.ReopenWindowMinutes = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["reopenWindowMinutes"].(int); ok {
        data.ReopenWindowMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["reopenWindowMinutes"].(int64); ok {
        data.ReopenWindowMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["reopenWindowMinutes"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.ReopenWindowMinutes = types.NumberValue(big.NewFloat(val))
        } else {
            data.ReopenWindowMinutes = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.ReopenWindowMinutes = types.NumberNull()
    }
    if val, ok := dataMap["enableInactivityTimeout"].(bool); ok {
        data.EnableInactivityTimeout = types.BoolValue(val)
    }
    if val, ok := dataMap["inactivityTimeoutMinutes"].(float64); ok {
        data.InactivityTimeoutMinutes = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["inactivityTimeoutMinutes"].(int); ok {
        data.InactivityTimeoutMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["inactivityTimeoutMinutes"].(int64); ok {
        data.InactivityTimeoutMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["inactivityTimeoutMinutes"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.InactivityTimeoutMinutes = types.NumberValue(big.NewFloat(val))
        } else {
            data.InactivityTimeoutMinutes = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.InactivityTimeoutMinutes = types.NumberNull()
    }
    if val, ok := dataMap["onCallDutyPolicies"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.OnCallDutyPolicies = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.OnCallDutyPolicies = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if obj, ok := dataMap["defaultAssignToUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DefaultAssignToUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DefaultAssignToUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DefaultAssignToUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DefaultAssignToUserId = types.StringValue(string(jsonBytes))
            } else {
                data.DefaultAssignToUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DefaultAssignToUserId = types.StringValue(string(jsonBytes))
            } else {
                data.DefaultAssignToUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DefaultAssignToUserId = types.StringValue(string(jsonBytes))
        } else {
            data.DefaultAssignToUserId = types.StringNull()
        }
    } else if val, ok := dataMap["defaultAssignToUserId"].(string); ok {
        data.DefaultAssignToUserId = types.StringValue(val)
    } else {
        data.DefaultAssignToUserId = types.StringNull()
    }
    if obj, ok := dataMap["defaultAssignToTeamId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DefaultAssignToTeamId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DefaultAssignToTeamId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DefaultAssignToTeamId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DefaultAssignToTeamId = types.StringValue(string(jsonBytes))
            } else {
                data.DefaultAssignToTeamId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DefaultAssignToTeamId = types.StringValue(string(jsonBytes))
            } else {
                data.DefaultAssignToTeamId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DefaultAssignToTeamId = types.StringValue(string(jsonBytes))
        } else {
            data.DefaultAssignToTeamId = types.StringNull()
        }
    } else if val, ok := dataMap["defaultAssignToTeamId"].(string); ok {
        data.DefaultAssignToTeamId = types.StringValue(val)
    } else {
        data.DefaultAssignToTeamId = types.StringNull()
    }
    if val, ok := dataMap["episodeLabels"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.EpisodeLabels = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.EpisodeLabels = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if val, ok := dataMap["episodeOwnerUsers"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.EpisodeOwnerUsers = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.EpisodeOwnerUsers = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if val, ok := dataMap["episodeOwnerTeams"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.EpisodeOwnerTeams = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.EpisodeOwnerTeams = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if obj, ok := dataMap["createdByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["createdByUserId"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if obj, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.CreatedAt = NewRFC3339Value(val)
        } else {
            data.CreatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["createdAt"].(string); ok && val != "" {
        data.CreatedAt = NewRFC3339Value(val)
    } else {
        data.CreatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.UpdatedAt = NewRFC3339Value(val)
        } else {
            data.UpdatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["updatedAt"].(string); ok && val != "" {
        data.UpdatedAt = NewRFC3339Value(val)
    } else {
        data.UpdatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["deletedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.DeletedAt = NewRFC3339Value(val)
        } else {
            data.DeletedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["deletedAt"].(string); ok && val != "" {
        data.DeletedAt = NewRFC3339Value(val)
    } else {
        data.DeletedAt = NewRFC3339Null()
    }
    if val, ok := dataMap["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["version"].(int); ok {
        data.Version = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["version"].(int64); ok {
        data.Version = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["version"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.Version = types.NumberValue(big.NewFloat(val))
        } else {
            data.Version = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.Version = types.NumberNull()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AlertGroupingRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var data AlertGroupingRuleResourceModel
    var state AlertGroupingRuleResourceModel

    // Read Terraform current state data to get the ID
    resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Read Terraform plan data to get the new values
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Use the ID from the current state
    data.Id = state.Id

    // Create API request body
    alertGroupingRuleRequest := map[string]interface{}{
        "data": map[string]interface{}{},
    }
    requestDataMap := alertGroupingRuleRequest["data"].(map[string]interface{})

    if !data.Name.IsUnknown() && !state.Name.IsUnknown() && !data.Name.Equal(state.Name) {
        requestDataMap["name"] = data.Name.ValueString()
    }
    if !data.Description.IsUnknown() && !state.Description.IsUnknown() && !data.Description.Equal(state.Description) {
        requestDataMap["description"] = data.Description.ValueString()
    }
    if !data.Priority.IsUnknown() && !state.Priority.IsUnknown() && !data.Priority.Equal(state.Priority) {
        requestDataMap["priority"] = r.bigFloatToFloat64(data.Priority.ValueBigFloat())
    }
    if !data.IsEnabled.IsUnknown() && !state.IsEnabled.IsUnknown() && !data.IsEnabled.Equal(state.IsEnabled) {
        requestDataMap["isEnabled"] = data.IsEnabled.ValueBool()
    }
    if !data.MatchCriteria.IsUnknown() && !state.MatchCriteria.IsUnknown() && !data.MatchCriteria.Equal(state.MatchCriteria) {
        var matchcriteriaData interface{}
        if err := json.Unmarshal([]byte(data.MatchCriteria.ValueString()), &matchcriteriaData); err == nil {
            requestDataMap["matchCriteria"] = matchcriteriaData
        } else {
            requestDataMap["matchCriteria"] = data.MatchCriteria.ValueString()
        }
    }
    if !data.Monitors.IsUnknown() && !state.Monitors.IsUnknown() && !data.Monitors.Equal(state.Monitors) {
        requestDataMap["monitors"] = r.convertTerraformSetToInterface(data.Monitors)
    }
    if !data.AlertSeverities.IsUnknown() && !state.AlertSeverities.IsUnknown() && !data.AlertSeverities.Equal(state.AlertSeverities) {
        requestDataMap["alertSeverities"] = r.convertTerraformSetToInterface(data.AlertSeverities)
    }
    if !data.AlertLabels.IsUnknown() && !state.AlertLabels.IsUnknown() && !data.AlertLabels.Equal(state.AlertLabels) {
        requestDataMap["alertLabels"] = r.convertTerraformSetToInterface(data.AlertLabels)
    }
    if !data.MonitorLabels.IsUnknown() && !state.MonitorLabels.IsUnknown() && !data.MonitorLabels.Equal(state.MonitorLabels) {
        requestDataMap["monitorLabels"] = r.convertTerraformSetToInterface(data.MonitorLabels)
    }
    if !data.AlertTitlePattern.IsUnknown() && !state.AlertTitlePattern.IsUnknown() && !data.AlertTitlePattern.Equal(state.AlertTitlePattern) {
        requestDataMap["alertTitlePattern"] = data.AlertTitlePattern.ValueString()
    }
    if !data.AlertDescriptionPattern.IsUnknown() && !state.AlertDescriptionPattern.IsUnknown() && !data.AlertDescriptionPattern.Equal(state.AlertDescriptionPattern) {
        requestDataMap["alertDescriptionPattern"] = data.AlertDescriptionPattern.ValueString()
    }
    if !data.MonitorNamePattern.IsUnknown() && !state.MonitorNamePattern.IsUnknown() && !data.MonitorNamePattern.Equal(state.MonitorNamePattern) {
        requestDataMap["monitorNamePattern"] = data.MonitorNamePattern.ValueString()
    }
    if !data.MonitorDescriptionPattern.IsUnknown() && !state.MonitorDescriptionPattern.IsUnknown() && !data.MonitorDescriptionPattern.Equal(state.MonitorDescriptionPattern) {
        requestDataMap["monitorDescriptionPattern"] = data.MonitorDescriptionPattern.ValueString()
    }
    if !data.GroupByMonitor.IsUnknown() && !state.GroupByMonitor.IsUnknown() && !data.GroupByMonitor.Equal(state.GroupByMonitor) {
        requestDataMap["groupByMonitor"] = data.GroupByMonitor.ValueBool()
    }
    if !data.GroupBySeverity.IsUnknown() && !state.GroupBySeverity.IsUnknown() && !data.GroupBySeverity.Equal(state.GroupBySeverity) {
        requestDataMap["groupBySeverity"] = data.GroupBySeverity.ValueBool()
    }
    if !data.GroupByAlertTitle.IsUnknown() && !state.GroupByAlertTitle.IsUnknown() && !data.GroupByAlertTitle.Equal(state.GroupByAlertTitle) {
        requestDataMap["groupByAlertTitle"] = data.GroupByAlertTitle.ValueBool()
    }
    if !data.GroupByAlertLabels.IsUnknown() && !state.GroupByAlertLabels.IsUnknown() && !data.GroupByAlertLabels.Equal(state.GroupByAlertLabels) {
        requestDataMap["groupByAlertLabels"] = data.GroupByAlertLabels.ValueBool()
    }
    if !data.GroupByMonitorLabels.IsUnknown() && !state.GroupByMonitorLabels.IsUnknown() && !data.GroupByMonitorLabels.Equal(state.GroupByMonitorLabels) {
        requestDataMap["groupByMonitorLabels"] = data.GroupByMonitorLabels.ValueBool()
    }
    if !data.EnableTimeWindow.IsUnknown() && !state.EnableTimeWindow.IsUnknown() && !data.EnableTimeWindow.Equal(state.EnableTimeWindow) {
        requestDataMap["enableTimeWindow"] = data.EnableTimeWindow.ValueBool()
    }
    if !data.TimeWindowMinutes.IsUnknown() && !state.TimeWindowMinutes.IsUnknown() && !data.TimeWindowMinutes.Equal(state.TimeWindowMinutes) {
        requestDataMap["timeWindowMinutes"] = r.bigFloatToFloat64(data.TimeWindowMinutes.ValueBigFloat())
    }
    if !data.GroupByFields.IsUnknown() && !state.GroupByFields.IsUnknown() && !data.GroupByFields.Equal(state.GroupByFields) {
        var groupbyfieldsData interface{}
        if err := json.Unmarshal([]byte(data.GroupByFields.ValueString()), &groupbyfieldsData); err == nil {
            requestDataMap["groupByFields"] = groupbyfieldsData
        } else {
            requestDataMap["groupByFields"] = data.GroupByFields.ValueString()
        }
    }
    if !data.EpisodeTitleTemplate.IsUnknown() && !state.EpisodeTitleTemplate.IsUnknown() && !data.EpisodeTitleTemplate.Equal(state.EpisodeTitleTemplate) {
        requestDataMap["episodeTitleTemplate"] = data.EpisodeTitleTemplate.ValueString()
    }
    if !data.EpisodeDescriptionTemplate.IsUnknown() && !state.EpisodeDescriptionTemplate.IsUnknown() && !data.EpisodeDescriptionTemplate.Equal(state.EpisodeDescriptionTemplate) {
        requestDataMap["episodeDescriptionTemplate"] = data.EpisodeDescriptionTemplate.ValueString()
    }
    if !data.EnableResolveDelay.IsUnknown() && !state.EnableResolveDelay.IsUnknown() && !data.EnableResolveDelay.Equal(state.EnableResolveDelay) {
        requestDataMap["enableResolveDelay"] = data.EnableResolveDelay.ValueBool()
    }
    if !data.ResolveDelayMinutes.IsUnknown() && !state.ResolveDelayMinutes.IsUnknown() && !data.ResolveDelayMinutes.Equal(state.ResolveDelayMinutes) {
        requestDataMap["resolveDelayMinutes"] = r.bigFloatToFloat64(data.ResolveDelayMinutes.ValueBigFloat())
    }
    if !data.EnableReopenWindow.IsUnknown() && !state.EnableReopenWindow.IsUnknown() && !data.EnableReopenWindow.Equal(state.EnableReopenWindow) {
        requestDataMap["enableReopenWindow"] = data.EnableReopenWindow.ValueBool()
    }
    if !data.ReopenWindowMinutes.IsUnknown() && !state.ReopenWindowMinutes.IsUnknown() && !data.ReopenWindowMinutes.Equal(state.ReopenWindowMinutes) {
        requestDataMap["reopenWindowMinutes"] = r.bigFloatToFloat64(data.ReopenWindowMinutes.ValueBigFloat())
    }
    if !data.EnableInactivityTimeout.IsUnknown() && !state.EnableInactivityTimeout.IsUnknown() && !data.EnableInactivityTimeout.Equal(state.EnableInactivityTimeout) {
        requestDataMap["enableInactivityTimeout"] = data.EnableInactivityTimeout.ValueBool()
    }
    if !data.InactivityTimeoutMinutes.IsUnknown() && !state.InactivityTimeoutMinutes.IsUnknown() && !data.InactivityTimeoutMinutes.Equal(state.InactivityTimeoutMinutes) {
        requestDataMap["inactivityTimeoutMinutes"] = r.bigFloatToFloat64(data.InactivityTimeoutMinutes.ValueBigFloat())
    }
    if !data.OnCallDutyPolicies.IsUnknown() && !state.OnCallDutyPolicies.IsUnknown() && !data.OnCallDutyPolicies.Equal(state.OnCallDutyPolicies) {
        requestDataMap["onCallDutyPolicies"] = r.convertTerraformSetToInterface(data.OnCallDutyPolicies)
    }
    if !data.DefaultAssignToUserId.IsUnknown() && !state.DefaultAssignToUserId.IsUnknown() && !data.DefaultAssignToUserId.Equal(state.DefaultAssignToUserId) {
        requestDataMap["defaultAssignToUserId"] = data.DefaultAssignToUserId.ValueString()
    }
    if !data.DefaultAssignToTeamId.IsUnknown() && !state.DefaultAssignToTeamId.IsUnknown() && !data.DefaultAssignToTeamId.Equal(state.DefaultAssignToTeamId) {
        requestDataMap["defaultAssignToTeamId"] = data.DefaultAssignToTeamId.ValueString()
    }
    if !data.EpisodeLabels.IsUnknown() && !state.EpisodeLabels.IsUnknown() && !data.EpisodeLabels.Equal(state.EpisodeLabels) {
        requestDataMap["episodeLabels"] = r.convertTerraformSetToInterface(data.EpisodeLabels)
    }
    if !data.EpisodeOwnerUsers.IsUnknown() && !state.EpisodeOwnerUsers.IsUnknown() && !data.EpisodeOwnerUsers.Equal(state.EpisodeOwnerUsers) {
        requestDataMap["episodeOwnerUsers"] = r.convertTerraformSetToInterface(data.EpisodeOwnerUsers)
    }
    if !data.EpisodeOwnerTeams.IsUnknown() && !state.EpisodeOwnerTeams.IsUnknown() && !data.EpisodeOwnerTeams.Equal(state.EpisodeOwnerTeams) {
        requestDataMap["episodeOwnerTeams"] = r.convertTerraformSetToInterface(data.EpisodeOwnerTeams)
    }

    // Only call the API when there are changed fields to send. An empty
    // update body is rejected by the API; state is still refreshed below so
    // this method never writes unverified plan values into state.
    if len(alertGroupingRuleRequest["data"].(map[string]interface{})) > 0 {
        httpResp, err := r.client.Put(ctx, "/alert-grouping-rule/" + data.Id.ValueString() + "", alertGroupingRuleRequest)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update alert_grouping_rule, got error: %s", err))
            return
        }

        // Parse the update response
        var alertGroupingRuleResponse map[string]interface{}
        err = r.client.ParseResponse(httpResp, &alertGroupingRuleResponse)
        if err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to update alert_grouping_rule: %s", err))
            return
        }
        _ = alertGroupingRuleResponse
    }

    // After successful update, fetch the current state by calling Read with select parameter
    selectParam := map[string]interface{}{
        "projectId": true,
        "name": true,
        "description": true,
        "priority": true,
        "isEnabled": true,
        "matchCriteria": true,
        "monitors": true,
        "alertSeverities": true,
        "alertLabels": true,
        "monitorLabels": true,
        "alertTitlePattern": true,
        "alertDescriptionPattern": true,
        "monitorNamePattern": true,
        "monitorDescriptionPattern": true,
        "groupByMonitor": true,
        "groupBySeverity": true,
        "groupByAlertTitle": true,
        "groupByAlertLabels": true,
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
        "createdByUserId": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect(ctx, "/alert-grouping-rule/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read alert_grouping_rule after update, got error: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read alert_grouping_rule after update: %s", err))
        return
    }

    // Update the model with response data from the Read operation
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := readResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = readResponse
    }

    if obj, ok := dataMap["projectId"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok {
            data.ProjectId = types.StringValue(val)
        } else {
            data.ProjectId = types.StringNull()
        }
    } else if val, ok := dataMap["projectId"].(string); ok {
        data.ProjectId = types.StringValue(val)
    } else {
        data.ProjectId = types.StringNull()
    }
    if obj, ok := dataMap["name"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Name = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Name = types.StringValue(string(jsonBytes))
            } else {
                data.Name = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Name = types.StringValue(string(jsonBytes))
            } else {
                data.Name = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Name = types.StringValue(string(jsonBytes))
        } else {
            data.Name = types.StringNull()
        }
    } else if val, ok := dataMap["name"].(string); ok {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if obj, ok := dataMap["description"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Description = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Description = types.StringValue(string(jsonBytes))
            } else {
                data.Description = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Description = types.StringValue(string(jsonBytes))
            } else {
                data.Description = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Description = types.StringValue(string(jsonBytes))
        } else {
            data.Description = types.StringNull()
        }
    } else if val, ok := dataMap["description"].(string); ok {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if val, ok := dataMap["priority"].(float64); ok {
        data.Priority = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["priority"].(int); ok {
        data.Priority = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["priority"].(int64); ok {
        data.Priority = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["priority"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.Priority = types.NumberValue(big.NewFloat(val))
        } else {
            data.Priority = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.Priority = types.NumberNull()
    }
    if val, ok := dataMap["isEnabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    }
    if obj, ok := dataMap["matchCriteria"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MatchCriteria = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.MatchCriteria = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.MatchCriteria = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.MatchCriteria = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.MatchCriteria = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.MatchCriteria = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.MatchCriteria = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.MatchCriteria = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.MatchCriteria = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["matchCriteria"].(string); ok {
        data.MatchCriteria = NewJSONSubsetValue(val)
    } else {
        data.MatchCriteria = NewJSONSubsetNull()
    }
    if val, ok := dataMap["monitors"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.Monitors = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.Monitors = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if val, ok := dataMap["alertSeverities"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.AlertSeverities = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.AlertSeverities = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if val, ok := dataMap["alertLabels"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.AlertLabels = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.AlertLabels = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if val, ok := dataMap["monitorLabels"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.MonitorLabels = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.MonitorLabels = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if obj, ok := dataMap["alertTitlePattern"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AlertTitlePattern = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AlertTitlePattern = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AlertTitlePattern = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AlertTitlePattern = types.StringValue(string(jsonBytes))
            } else {
                data.AlertTitlePattern = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AlertTitlePattern = types.StringValue(string(jsonBytes))
            } else {
                data.AlertTitlePattern = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AlertTitlePattern = types.StringValue(string(jsonBytes))
        } else {
            data.AlertTitlePattern = types.StringNull()
        }
    } else if val, ok := dataMap["alertTitlePattern"].(string); ok {
        data.AlertTitlePattern = types.StringValue(val)
    } else {
        data.AlertTitlePattern = types.StringNull()
    }
    if obj, ok := dataMap["alertDescriptionPattern"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AlertDescriptionPattern = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AlertDescriptionPattern = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AlertDescriptionPattern = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AlertDescriptionPattern = types.StringValue(string(jsonBytes))
            } else {
                data.AlertDescriptionPattern = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AlertDescriptionPattern = types.StringValue(string(jsonBytes))
            } else {
                data.AlertDescriptionPattern = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AlertDescriptionPattern = types.StringValue(string(jsonBytes))
        } else {
            data.AlertDescriptionPattern = types.StringNull()
        }
    } else if val, ok := dataMap["alertDescriptionPattern"].(string); ok {
        data.AlertDescriptionPattern = types.StringValue(val)
    } else {
        data.AlertDescriptionPattern = types.StringNull()
    }
    if obj, ok := dataMap["monitorNamePattern"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MonitorNamePattern = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.MonitorNamePattern = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.MonitorNamePattern = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.MonitorNamePattern = types.StringValue(string(jsonBytes))
            } else {
                data.MonitorNamePattern = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.MonitorNamePattern = types.StringValue(string(jsonBytes))
            } else {
                data.MonitorNamePattern = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.MonitorNamePattern = types.StringValue(string(jsonBytes))
        } else {
            data.MonitorNamePattern = types.StringNull()
        }
    } else if val, ok := dataMap["monitorNamePattern"].(string); ok {
        data.MonitorNamePattern = types.StringValue(val)
    } else {
        data.MonitorNamePattern = types.StringNull()
    }
    if obj, ok := dataMap["monitorDescriptionPattern"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MonitorDescriptionPattern = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.MonitorDescriptionPattern = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.MonitorDescriptionPattern = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.MonitorDescriptionPattern = types.StringValue(string(jsonBytes))
            } else {
                data.MonitorDescriptionPattern = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.MonitorDescriptionPattern = types.StringValue(string(jsonBytes))
            } else {
                data.MonitorDescriptionPattern = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.MonitorDescriptionPattern = types.StringValue(string(jsonBytes))
        } else {
            data.MonitorDescriptionPattern = types.StringNull()
        }
    } else if val, ok := dataMap["monitorDescriptionPattern"].(string); ok {
        data.MonitorDescriptionPattern = types.StringValue(val)
    } else {
        data.MonitorDescriptionPattern = types.StringNull()
    }
    if val, ok := dataMap["groupByMonitor"].(bool); ok {
        data.GroupByMonitor = types.BoolValue(val)
    }
    if val, ok := dataMap["groupBySeverity"].(bool); ok {
        data.GroupBySeverity = types.BoolValue(val)
    }
    if val, ok := dataMap["groupByAlertTitle"].(bool); ok {
        data.GroupByAlertTitle = types.BoolValue(val)
    }
    if val, ok := dataMap["groupByAlertLabels"].(bool); ok {
        data.GroupByAlertLabels = types.BoolValue(val)
    }
    if val, ok := dataMap["groupByMonitorLabels"].(bool); ok {
        data.GroupByMonitorLabels = types.BoolValue(val)
    }
    if val, ok := dataMap["enableTimeWindow"].(bool); ok {
        data.EnableTimeWindow = types.BoolValue(val)
    }
    if val, ok := dataMap["timeWindowMinutes"].(float64); ok {
        data.TimeWindowMinutes = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["timeWindowMinutes"].(int); ok {
        data.TimeWindowMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["timeWindowMinutes"].(int64); ok {
        data.TimeWindowMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["timeWindowMinutes"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.TimeWindowMinutes = types.NumberValue(big.NewFloat(val))
        } else {
            data.TimeWindowMinutes = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.TimeWindowMinutes = types.NumberNull()
    }
    if obj, ok := dataMap["groupByFields"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.GroupByFields = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.GroupByFields = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.GroupByFields = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.GroupByFields = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.GroupByFields = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.GroupByFields = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.GroupByFields = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.GroupByFields = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.GroupByFields = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["groupByFields"].(string); ok {
        data.GroupByFields = NewJSONSubsetValue(val)
    } else {
        data.GroupByFields = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["episodeTitleTemplate"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EpisodeTitleTemplate = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.EpisodeTitleTemplate = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.EpisodeTitleTemplate = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.EpisodeTitleTemplate = types.StringValue(string(jsonBytes))
            } else {
                data.EpisodeTitleTemplate = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.EpisodeTitleTemplate = types.StringValue(string(jsonBytes))
            } else {
                data.EpisodeTitleTemplate = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.EpisodeTitleTemplate = types.StringValue(string(jsonBytes))
        } else {
            data.EpisodeTitleTemplate = types.StringNull()
        }
    } else if val, ok := dataMap["episodeTitleTemplate"].(string); ok {
        data.EpisodeTitleTemplate = types.StringValue(val)
    } else {
        data.EpisodeTitleTemplate = types.StringNull()
    }
    if obj, ok := dataMap["episodeDescriptionTemplate"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EpisodeDescriptionTemplate = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.EpisodeDescriptionTemplate = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.EpisodeDescriptionTemplate = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.EpisodeDescriptionTemplate = types.StringValue(string(jsonBytes))
            } else {
                data.EpisodeDescriptionTemplate = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.EpisodeDescriptionTemplate = types.StringValue(string(jsonBytes))
            } else {
                data.EpisodeDescriptionTemplate = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.EpisodeDescriptionTemplate = types.StringValue(string(jsonBytes))
        } else {
            data.EpisodeDescriptionTemplate = types.StringNull()
        }
    } else if val, ok := dataMap["episodeDescriptionTemplate"].(string); ok {
        data.EpisodeDescriptionTemplate = types.StringValue(val)
    } else {
        data.EpisodeDescriptionTemplate = types.StringNull()
    }
    if val, ok := dataMap["enableResolveDelay"].(bool); ok {
        data.EnableResolveDelay = types.BoolValue(val)
    }
    if val, ok := dataMap["resolveDelayMinutes"].(float64); ok {
        data.ResolveDelayMinutes = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["resolveDelayMinutes"].(int); ok {
        data.ResolveDelayMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["resolveDelayMinutes"].(int64); ok {
        data.ResolveDelayMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["resolveDelayMinutes"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.ResolveDelayMinutes = types.NumberValue(big.NewFloat(val))
        } else {
            data.ResolveDelayMinutes = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.ResolveDelayMinutes = types.NumberNull()
    }
    if val, ok := dataMap["enableReopenWindow"].(bool); ok {
        data.EnableReopenWindow = types.BoolValue(val)
    }
    if val, ok := dataMap["reopenWindowMinutes"].(float64); ok {
        data.ReopenWindowMinutes = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["reopenWindowMinutes"].(int); ok {
        data.ReopenWindowMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["reopenWindowMinutes"].(int64); ok {
        data.ReopenWindowMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["reopenWindowMinutes"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.ReopenWindowMinutes = types.NumberValue(big.NewFloat(val))
        } else {
            data.ReopenWindowMinutes = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.ReopenWindowMinutes = types.NumberNull()
    }
    if val, ok := dataMap["enableInactivityTimeout"].(bool); ok {
        data.EnableInactivityTimeout = types.BoolValue(val)
    }
    if val, ok := dataMap["inactivityTimeoutMinutes"].(float64); ok {
        data.InactivityTimeoutMinutes = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["inactivityTimeoutMinutes"].(int); ok {
        data.InactivityTimeoutMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["inactivityTimeoutMinutes"].(int64); ok {
        data.InactivityTimeoutMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["inactivityTimeoutMinutes"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.InactivityTimeoutMinutes = types.NumberValue(big.NewFloat(val))
        } else {
            data.InactivityTimeoutMinutes = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.InactivityTimeoutMinutes = types.NumberNull()
    }
    if val, ok := dataMap["onCallDutyPolicies"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.OnCallDutyPolicies = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.OnCallDutyPolicies = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if obj, ok := dataMap["defaultAssignToUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DefaultAssignToUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DefaultAssignToUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DefaultAssignToUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DefaultAssignToUserId = types.StringValue(string(jsonBytes))
            } else {
                data.DefaultAssignToUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DefaultAssignToUserId = types.StringValue(string(jsonBytes))
            } else {
                data.DefaultAssignToUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DefaultAssignToUserId = types.StringValue(string(jsonBytes))
        } else {
            data.DefaultAssignToUserId = types.StringNull()
        }
    } else if val, ok := dataMap["defaultAssignToUserId"].(string); ok {
        data.DefaultAssignToUserId = types.StringValue(val)
    } else {
        data.DefaultAssignToUserId = types.StringNull()
    }
    if obj, ok := dataMap["defaultAssignToTeamId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DefaultAssignToTeamId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DefaultAssignToTeamId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DefaultAssignToTeamId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DefaultAssignToTeamId = types.StringValue(string(jsonBytes))
            } else {
                data.DefaultAssignToTeamId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DefaultAssignToTeamId = types.StringValue(string(jsonBytes))
            } else {
                data.DefaultAssignToTeamId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DefaultAssignToTeamId = types.StringValue(string(jsonBytes))
        } else {
            data.DefaultAssignToTeamId = types.StringNull()
        }
    } else if val, ok := dataMap["defaultAssignToTeamId"].(string); ok {
        data.DefaultAssignToTeamId = types.StringValue(val)
    } else {
        data.DefaultAssignToTeamId = types.StringNull()
    }
    if val, ok := dataMap["episodeLabels"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.EpisodeLabels = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.EpisodeLabels = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if val, ok := dataMap["episodeOwnerUsers"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.EpisodeOwnerUsers = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.EpisodeOwnerUsers = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if val, ok := dataMap["episodeOwnerTeams"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.EpisodeOwnerTeams = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.EpisodeOwnerTeams = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if obj, ok := dataMap["createdByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["createdByUserId"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if obj, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.CreatedAt = NewRFC3339Value(val)
        } else {
            data.CreatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["createdAt"].(string); ok && val != "" {
        data.CreatedAt = NewRFC3339Value(val)
    } else {
        data.CreatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.UpdatedAt = NewRFC3339Value(val)
        } else {
            data.UpdatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["updatedAt"].(string); ok && val != "" {
        data.UpdatedAt = NewRFC3339Value(val)
    } else {
        data.UpdatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["deletedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.DeletedAt = NewRFC3339Value(val)
        } else {
            data.DeletedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["deletedAt"].(string); ok && val != "" {
        data.DeletedAt = NewRFC3339Value(val)
    } else {
        data.DeletedAt = NewRFC3339Null()
    }
    if val, ok := dataMap["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["version"].(int); ok {
        data.Version = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["version"].(int64); ok {
        data.Version = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["version"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.Version = types.NumberValue(big.NewFloat(val))
        } else {
            data.Version = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.Version = types.NumberNull()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    data.Id = state.Id

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AlertGroupingRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var data AlertGroupingRuleResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Make API call
    httpResp, err := r.client.Delete(ctx, "/alert-grouping-rule/" + data.Id.ValueString() + "")
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete alert_grouping_rule, got error: %s", err))
        return
    }

    // A failed delete must keep the resource in state — silently dropping it
    // orphans real infrastructure. 404 means it is already gone.
    if httpResp.StatusCode >= 400 && httpResp.StatusCode != http.StatusNotFound {
        err = r.client.ParseResponse(httpResp, nil)
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to delete alert_grouping_rule: %s", err))
        return
    }
    if httpResp.Body != nil {
        httpResp.Body.Close()
    }
}


func (r *AlertGroupingRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Helper method to convert Terraform map to Go interface{}
func (r *AlertGroupingRuleResource) convertTerraformMapToInterface(terraformMap types.Map) interface{} {
    if terraformMap.IsNull() || terraformMap.IsUnknown() {
        return nil
    }
    
    result := make(map[string]string)
    terraformMap.ElementsAs(context.Background(), &result, false)
    
    // Convert map[string]string to map[string]interface{}
    interfaceResult := make(map[string]interface{})
    for key, value := range result {
        interfaceResult[key] = value
    }
    
    return interfaceResult
}

// Helper method to convert Terraform list to Go interface{}
func (r *AlertGroupingRuleResource) convertTerraformListToInterface(terraformList types.List) interface{} {
    if terraformList.IsNull() || terraformList.IsUnknown() {
        return nil
    }

    var stringList []string
    terraformList.ElementsAs(context.Background(), &stringList, false)

    // Convert string array to OneUptime format with _id fields
    var result []interface{}
    for _, str := range stringList {
        if str != "" {
            result = append(result, map[string]interface{}{
                "_id": str,
            })
        }
    }
    return result
}

// Helper method to convert Terraform set to Go interface{}
func (r *AlertGroupingRuleResource) convertTerraformSetToInterface(terraformSet types.Set) interface{} {
    if terraformSet.IsNull() || terraformSet.IsUnknown() {
        return nil
    }

    var stringList []string
    terraformSet.ElementsAs(context.Background(), &stringList, false)

    // Convert string array to OneUptime format with _id fields
    var result []interface{}
    for _, str := range stringList {
        if str != "" {
            result = append(result, map[string]interface{}{
                "_id": str,
            })
        }
    }
    return result
}


// Helper method to parse JSON field for complex objects
func (r *AlertGroupingRuleResource) parseJSONField(terraformString basetypes.StringValuable) interface{} {
    sv, _ := terraformString.ToStringValue(context.Background())
    if sv.IsNull() || sv.IsUnknown() || sv.ValueString() == "" {
        return nil
    }

    var result interface{}
    if err := json.Unmarshal([]byte(sv.ValueString()), &result); err != nil {
        // If JSON parsing fails, return the raw string
        return sv.ValueString()
    }

    return result
}

// Normalize URL wrapper objects to avoid drift (e.g., trailing slash differences).
func (r *AlertGroupingRuleResource) normalizeURLWrappers(value interface{}) interface{} {
    switch v := value.(type) {
    case map[string]interface{}:
        if typeStr, ok := v["_type"].(string); ok && typeStr == "URL" {
            if val, ok := v["value"].(string); ok {
                v["value"] = r.normalizeURLString(val)
            }
        }
        for key, child := range v {
            v[key] = r.normalizeURLWrappers(child)
        }
        return v
    case []interface{}:
        for i, child := range v {
            v[i] = r.normalizeURLWrappers(child)
        }
        return v
    default:
        return v
    }
}

func (r *AlertGroupingRuleResource) normalizeURLString(value string) string {
    parsed, err := url.Parse(value)
    if err != nil {
        return value
    }
    if parsed.Path == "/" && parsed.RawQuery == "" && parsed.Fragment == "" {
        return strings.TrimSuffix(value, "/")
    }
    return value
}

// Helper method to convert *big.Float to float64 for JSON serialization
func (r *AlertGroupingRuleResource) bigFloatToFloat64(bf *big.Float) interface{} {
    if bf == nil {
        return nil
    }
    f, _ := bf.Float64()
    return f
}

// Helper method to check if a type string is a valid OneUptime ObjectType.
// The registry itself lives in objecttypes.go, shared across the package.
func (r *AlertGroupingRuleResource) isValidOneUptimeObjectType(typeStr string) bool {
    return validOneUptimeObjectTypes[typeStr]
}
