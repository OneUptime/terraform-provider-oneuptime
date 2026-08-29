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
var _ datasource.DataSource = &AlertGroupingRuleDataSource{}

func NewAlertGroupingRuleDataSource() datasource.DataSource {
    return &AlertGroupingRuleDataSource{}
}

// AlertGroupingRuleDataSource defines the data source implementation.
type AlertGroupingRuleDataSource struct {
    client *Client
}

// AlertGroupingRuleDataSourceModel describes the data source data model.
type AlertGroupingRuleDataSourceModel struct {
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
    GroupByAlertLabels types.Bool `tfsdk:"group_by_alert_labels"`
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
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *AlertGroupingRuleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_alert_grouping_rule"
}

func (d *AlertGroupingRuleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Configure rules for automatically grouping related alerts into episodes Look up an existing alert_grouping_rule by `id` or by `name`.",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Look up by unique identifier. Exactly one of `id` or `name` must be set.",
                Optional: true,
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.",
                Optional: true,
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
                MarkdownDescription: "Description of this alert grouping rule.",
                Computed: true,
            },
            "priority": schema.NumberAttribute{
                MarkdownDescription: "Priority of this rule. Lower number = higher priority. Rules are evaluated in priority order..",
                Computed: true,
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether this rule is enabled.",
                Computed: true,
            },
            "match_criteria": schema.StringAttribute{
                MarkdownDescription: "JSON object defining the criteria for matching alerts to this rule.",
                Computed: true,
            },
            "monitors": schema.SetAttribute{
                MarkdownDescription: "Only group alerts from these monitors. Leave empty to match alerts from any monitor..",
                Computed: true,
                ElementType: types.StringType,
            },
            "alert_severities": schema.SetAttribute{
                MarkdownDescription: "Only group alerts with these severities. Leave empty to match alerts of any severity..",
                Computed: true,
                ElementType: types.StringType,
            },
            "alert_labels": schema.SetAttribute{
                MarkdownDescription: "Only group alerts that have at least one of these labels. Leave empty to match alerts regardless of alert labels..",
                Computed: true,
                ElementType: types.StringType,
            },
            "monitor_labels": schema.SetAttribute{
                MarkdownDescription: "Only group alerts from monitors that have at least one of these labels. Leave empty to match alerts regardless of monitor labels..",
                Computed: true,
                ElementType: types.StringType,
            },
            "alert_title_pattern": schema.StringAttribute{
                MarkdownDescription: "Regular expression pattern to match alert titles. Leave empty to match any title. Example: 'CPU.*high' matches titles containing 'CPU' followed by 'high'..",
                Computed: true,
            },
            "alert_description_pattern": schema.StringAttribute{
                MarkdownDescription: "Regular expression pattern to match alert descriptions. Leave empty to match any description..",
                Computed: true,
            },
            "monitor_name_pattern": schema.StringAttribute{
                MarkdownDescription: "Regular expression pattern to match monitor names. Leave empty to match any monitor name. Example: 'prod-.*' matches monitors starting with 'prod-'..",
                Computed: true,
            },
            "monitor_description_pattern": schema.StringAttribute{
                MarkdownDescription: "Regular expression pattern to match monitor descriptions. Leave empty to match any monitor description..",
                Computed: true,
            },
            "group_by_monitor": schema.BoolAttribute{
                MarkdownDescription: "When enabled, alerts from different monitors will be grouped into separate episodes. When disabled, alerts from any monitor can be grouped together..",
                Computed: true,
            },
            "group_by_severity": schema.BoolAttribute{
                MarkdownDescription: "When enabled, alerts with different severities will be grouped into separate episodes. When disabled, alerts of any severity can be grouped together..",
                Computed: true,
            },
            "group_by_alert_title": schema.BoolAttribute{
                MarkdownDescription: "When enabled, alerts with different titles will be grouped into separate episodes. When disabled, alerts with any title can be grouped together..",
                Computed: true,
            },
            "group_by_alert_labels": schema.BoolAttribute{
                MarkdownDescription: "When enabled, alerts with different sets of labels will be grouped into separate episodes (exact set match). When disabled, alert labels are ignored for grouping..",
                Computed: true,
            },
            "group_by_monitor_labels": schema.BoolAttribute{
                MarkdownDescription: "When enabled, alerts whose monitors have different sets of labels will be grouped into separate episodes (exact set match). When disabled, monitor labels are ignored for grouping..",
                Computed: true,
            },
            "enable_time_window": schema.BoolAttribute{
                MarkdownDescription: "Enable time-based grouping. When enabled, alerts are grouped within the specified time window. When disabled, all matching alerts are grouped into a single ongoing episode regardless of time..",
                Computed: true,
            },
            "time_window_minutes": schema.NumberAttribute{
                MarkdownDescription: "Rolling time window in minutes. Alerts are grouped if they arrive within this gap from the last alert..",
                Computed: true,
            },
            "group_by_fields": schema.StringAttribute{
                MarkdownDescription: "JSON object defining the fields to group alerts by (e.g., monitorId, severity).",
                Computed: true,
            },
            "episode_title_template": schema.StringAttribute{
                MarkdownDescription: "Template for generating episode titles. Supports placeholders like {{alertSeverity}}, {{monitorName}}, {{alertTitle}}, {{alertDescription}}.",
                Computed: true,
            },
            "episode_description_template": schema.StringAttribute{
                MarkdownDescription: "Template for generating episode descriptions. Supports placeholders like {{alertSeverity}}, {{monitorName}}, {{alertTitle}}, {{alertDescription}}.",
                Computed: true,
            },
            "enable_resolve_delay": schema.BoolAttribute{
                MarkdownDescription: "Enable grace period before auto-resolving episode after all alerts resolve. Helps prevent rapid state changes during alert flapping..",
                Computed: true,
            },
            "resolve_delay_minutes": schema.NumberAttribute{
                MarkdownDescription: "Grace period in minutes before auto-resolving an episode after all alerts are resolved.",
                Computed: true,
            },
            "enable_reopen_window": schema.BoolAttribute{
                MarkdownDescription: "Enable reopening recently resolved episodes instead of creating new ones. Useful when related issues recur shortly after resolution..",
                Computed: true,
            },
            "reopen_window_minutes": schema.NumberAttribute{
                MarkdownDescription: "Time window in minutes to reopen a recently resolved episode instead of creating a new one.",
                Computed: true,
            },
            "enable_inactivity_timeout": schema.BoolAttribute{
                MarkdownDescription: "Enable auto-resolving episodes after a period of inactivity. Helps automatically close episodes when no new alerts arrive..",
                Computed: true,
            },
            "inactivity_timeout_minutes": schema.NumberAttribute{
                MarkdownDescription: "Time in minutes after which an inactive episode will be auto-resolved.",
                Computed: true,
            },
            "on_call_duty_policies": schema.SetAttribute{
                MarkdownDescription: "List of on-call duty policies to execute for episodes created by this rule..",
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
                MarkdownDescription: "Labels to automatically apply to episodes created by this rule..",
                Computed: true,
                ElementType: types.StringType,
            },
            "episode_owner_users": schema.SetAttribute{
                MarkdownDescription: "Users to automatically add as owners to episodes created by this rule..",
                Computed: true,
                ElementType: types.StringType,
            },
            "episode_owner_teams": schema.SetAttribute{
                MarkdownDescription: "Teams to automatically add as owners to episodes created by this rule..",
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

func (d *AlertGroupingRuleDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AlertGroupingRuleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data AlertGroupingRuleDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && data.Id.ValueString() != ""
    hasName := !data.Name.IsNull() && data.Name.ValueString() != ""
    if hasId == hasName {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Exactly one of `id` or `name` must be set to look up a alert_grouping_rule.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "name": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "projectId": true,
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
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/alert-grouping-rule/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read alert_grouping_rule, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No alert_grouping_rule found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read alert_grouping_rule: %s", err))
            return
        }
        if wrapper, ok := itemResponse["data"].(map[string]interface{}); ok {
            item = wrapper
        } else {
            item = itemResponse
        }
    } else {
        listBody := map[string]interface{}{
            "query": map[string]interface{}{
                "name": data.Name.ValueString(),
            },
            "select": selectParam,
            // limit 2 is enough to detect ambiguity without paging.
            "limit": 2,
        }
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/alert-grouping-rule/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list alert_grouping_rule, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list alert_grouping_rule: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No alert_grouping_rule found with name %q.", data.Name.ValueString()))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one alert_grouping_rule matches name %q. Use the id attribute to disambiguate.", data.Name.ValueString()))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for alert_grouping_rule.")
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
    if obj, ok := item["deletedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DeletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DeletedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DeletedAt = types.StringValue(string(jsonBytes))
        } else {
            data.DeletedAt = types.StringNull()
        }
    } else if val, ok := item["deletedAt"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    } else {
        data.DeletedAt = types.StringNull()
    }
    if val, ok := item["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["version"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.Version = types.NumberValue(big.NewFloat(val))
        } else {
            data.Version = types.NumberNull()
        }
    } else {
        data.Version = types.NumberNull()
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
    if val, ok := item["alertSeverities"].([]interface{}); ok {
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
        data.AlertSeverities = types.SetValueMust(types.StringType, setItems)
    } else {
        data.AlertSeverities = types.SetNull(types.StringType)
    }
    if val, ok := item["alertLabels"].([]interface{}); ok {
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
        data.AlertLabels = types.SetValueMust(types.StringType, setItems)
    } else {
        data.AlertLabels = types.SetNull(types.StringType)
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
    if obj, ok := item["alertTitlePattern"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AlertTitlePattern = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AlertTitlePattern = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AlertTitlePattern = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AlertTitlePattern = types.StringValue(string(jsonBytes))
        } else {
            data.AlertTitlePattern = types.StringNull()
        }
    } else if val, ok := item["alertTitlePattern"].(string); ok {
        data.AlertTitlePattern = types.StringValue(val)
    } else {
        data.AlertTitlePattern = types.StringNull()
    }
    if obj, ok := item["alertDescriptionPattern"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AlertDescriptionPattern = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AlertDescriptionPattern = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AlertDescriptionPattern = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AlertDescriptionPattern = types.StringValue(string(jsonBytes))
        } else {
            data.AlertDescriptionPattern = types.StringNull()
        }
    } else if val, ok := item["alertDescriptionPattern"].(string); ok {
        data.AlertDescriptionPattern = types.StringValue(val)
    } else {
        data.AlertDescriptionPattern = types.StringNull()
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
    if val, ok := item["groupByAlertTitle"].(bool); ok {
        data.GroupByAlertTitle = types.BoolValue(val)
    } else {
        data.GroupByAlertTitle = types.BoolNull()
    }
    if val, ok := item["groupByAlertLabels"].(bool); ok {
        data.GroupByAlertLabels = types.BoolValue(val)
    } else {
        data.GroupByAlertLabels = types.BoolNull()
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

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
