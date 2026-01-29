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
var _ datasource.DataSource = &IncidentGroupingRuleDataDataSource{}

func NewIncidentGroupingRuleDataDataSource() datasource.DataSource {
    return &IncidentGroupingRuleDataDataSource{}
}

// IncidentGroupingRuleDataDataSource defines the data source implementation.
type IncidentGroupingRuleDataDataSource struct {
    client *Client
}

// IncidentGroupingRuleDataDataSourceModel describes the data source data model.
type IncidentGroupingRuleDataDataSourceModel struct {
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
    GroupByService types.Bool `tfsdk:"group_by_service"`
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
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *IncidentGroupingRuleDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_incident_grouping_rule_data"
}

func (d *IncidentGroupingRuleDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "incident_grouping_rule_data data source",

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
                MarkdownDescription: "Description of this incident grouping rule. Permissions - Create: [Project Owner, Project Admin, Create Incident Grouping Rule], Read: [Project Owner, Project Admin, Project Member, Read Incident Grouping Rule], Update: [Project Owner, Project Admin, Edit Incident Grouping Rule]",
                Computed: true,
            },
            "priority": schema.NumberAttribute{
                MarkdownDescription: "Priority of this rule. Lower number = higher priority. Rules are evaluated in priority order.. Permissions - Create: [Project Owner, Project Admin, Create Incident Grouping Rule], Read: [Project Owner, Project Admin, Project Member, Read Incident Grouping Rule], Update: [Project Owner, Project Admin, Edit Incident Grouping Rule]",
                Computed: true,
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether this rule is enabled. Permissions - Create: [Project Owner, Project Admin, Create Incident Grouping Rule], Read: [Project Owner, Project Admin, Project Member, Read Incident Grouping Rule], Update: [Project Owner, Project Admin, Edit Incident Grouping Rule]",
                Computed: true,
            },
            "match_criteria": schema.StringAttribute{
                MarkdownDescription: "JSON object defining the criteria for matching incidents to this rule. Permissions - Create: [Project Owner, Project Admin, Create Incident Grouping Rule], Read: [Project Owner, Project Admin, Project Member, Read Incident Grouping Rule], Update: [Project Owner, Project Admin, Edit Incident Grouping Rule]",
                Computed: true,
            },
            "monitors": schema.SetAttribute{
                MarkdownDescription: "Only group incidents from these monitors. Leave empty to match incidents from any monitor.. Permissions - Create: [Project Owner, Project Admin, Create Incident Grouping Rule], Read: [Project Owner, Project Admin, Project Member, Read Incident Grouping Rule], Update: [Project Owner, Project Admin, Edit Incident Grouping Rule]",
                Computed: true,
                ElementType: types.StringType,
            },
            "incident_severities": schema.SetAttribute{
                MarkdownDescription: "Only group incidents with these severities. Leave empty to match incidents of any severity.. Permissions - Create: [Project Owner, Project Admin, Create Incident Grouping Rule], Read: [Project Owner, Project Admin, Project Member, Read Incident Grouping Rule], Update: [Project Owner, Project Admin, Edit Incident Grouping Rule]",
                Computed: true,
                ElementType: types.StringType,
            },
            "incident_labels": schema.SetAttribute{
                MarkdownDescription: "Only group incidents that have at least one of these labels. Leave empty to match incidents regardless of incident labels.. Permissions - Create: [Project Owner, Project Admin, Create Incident Grouping Rule], Read: [Project Owner, Project Admin, Project Member, Read Incident Grouping Rule], Update: [Project Owner, Project Admin, Edit Incident Grouping Rule]",
                Computed: true,
                ElementType: types.StringType,
            },
            "monitor_labels": schema.SetAttribute{
                MarkdownDescription: "Only group incidents from monitors that have at least one of these labels. Leave empty to match incidents regardless of monitor labels.. Permissions - Create: [Project Owner, Project Admin, Create Incident Grouping Rule], Read: [Project Owner, Project Admin, Project Member, Read Incident Grouping Rule], Update: [Project Owner, Project Admin, Edit Incident Grouping Rule]",
                Computed: true,
                ElementType: types.StringType,
            },
            "incident_title_pattern": schema.StringAttribute{
                MarkdownDescription: "Regular expression pattern to match incident titles. Leave empty to match any title. Example: 'CPU.*high' matches titles containing 'CPU' followed by 'high'.. Permissions - Create: [Project Owner, Project Admin, Create Incident Grouping Rule], Read: [Project Owner, Project Admin, Project Member, Read Incident Grouping Rule], Update: [Project Owner, Project Admin, Edit Incident Grouping Rule]",
                Computed: true,
            },
            "incident_description_pattern": schema.StringAttribute{
                MarkdownDescription: "Regular expression pattern to match incident descriptions. Leave empty to match any description.. Permissions - Create: [Project Owner, Project Admin, Create Incident Grouping Rule], Read: [Project Owner, Project Admin, Project Member, Read Incident Grouping Rule], Update: [Project Owner, Project Admin, Edit Incident Grouping Rule]",
                Computed: true,
            },
            "monitor_name_pattern": schema.StringAttribute{
                MarkdownDescription: "Regular expression pattern to match monitor names. Leave empty to match any monitor name. Example: 'prod-.*' matches monitors starting with 'prod-'.. Permissions - Create: [Project Owner, Project Admin, Create Incident Grouping Rule], Read: [Project Owner, Project Admin, Project Member, Read Incident Grouping Rule], Update: [Project Owner, Project Admin, Edit Incident Grouping Rule]",
                Computed: true,
            },
            "monitor_description_pattern": schema.StringAttribute{
                MarkdownDescription: "Regular expression pattern to match monitor descriptions. Leave empty to match any monitor description.. Permissions - Create: [Project Owner, Project Admin, Create Incident Grouping Rule], Read: [Project Owner, Project Admin, Project Member, Read Incident Grouping Rule], Update: [Project Owner, Project Admin, Edit Incident Grouping Rule]",
                Computed: true,
            },
            "group_by_monitor": schema.BoolAttribute{
                MarkdownDescription: "When enabled, incidents from different monitors will be grouped into separate episodes. When disabled, incidents from any monitor can be grouped together.. Permissions - Create: [Project Owner, Project Admin, Create Incident Grouping Rule], Read: [Project Owner, Project Admin, Project Member, Read Incident Grouping Rule], Update: [Project Owner, Project Admin, Edit Incident Grouping Rule]",
                Computed: true,
            },
            "group_by_severity": schema.BoolAttribute{
                MarkdownDescription: "When enabled, incidents with different severities will be grouped into separate episodes. When disabled, incidents of any severity can be grouped together.. Permissions - Create: [Project Owner, Project Admin, Create Incident Grouping Rule], Read: [Project Owner, Project Admin, Project Member, Read Incident Grouping Rule], Update: [Project Owner, Project Admin, Edit Incident Grouping Rule]",
                Computed: true,
            },
            "group_by_incident_title": schema.BoolAttribute{
                MarkdownDescription: "When enabled, incidents with different titles will be grouped into separate episodes. When disabled, incidents with any title can be grouped together.. Permissions - Create: [Project Owner, Project Admin, Create Incident Grouping Rule], Read: [Project Owner, Project Admin, Project Member, Read Incident Grouping Rule], Update: [Project Owner, Project Admin, Edit Incident Grouping Rule]",
                Computed: true,
            },
            "group_by_service": schema.BoolAttribute{
                MarkdownDescription: "When enabled, incidents from monitors belonging to different services will be grouped into separate episodes. When disabled, incidents can be grouped together regardless of which service the monitor belongs to.. Permissions - Create: [Project Owner, Project Admin, Create Incident Grouping Rule], Read: [Project Owner, Project Admin, Project Member, Read Incident Grouping Rule], Update: [Project Owner, Project Admin, Edit Incident Grouping Rule]",
                Computed: true,
            },
            "enable_time_window": schema.BoolAttribute{
                MarkdownDescription: "Enable time-based grouping. When enabled, incidents are grouped within the specified time window. When disabled, all matching incidents are grouped into a single ongoing episode regardless of time.. Permissions - Create: [Project Owner, Project Admin, Create Incident Grouping Rule], Read: [Project Owner, Project Admin, Project Member, Read Incident Grouping Rule], Update: [Project Owner, Project Admin, Edit Incident Grouping Rule]",
                Computed: true,
            },
            "time_window_minutes": schema.NumberAttribute{
                MarkdownDescription: "Rolling time window in minutes. Incidents are grouped if they arrive within this gap from the last incident.. Permissions - Create: [Project Owner, Project Admin, Create Incident Grouping Rule], Read: [Project Owner, Project Admin, Project Member, Read Incident Grouping Rule], Update: [Project Owner, Project Admin, Edit Incident Grouping Rule]",
                Computed: true,
            },
            "group_by_fields": schema.StringAttribute{
                MarkdownDescription: "JSON object defining the fields to group incidents by (e.g., monitorId, severity). Permissions - Create: [Project Owner, Project Admin, Create Incident Grouping Rule], Read: [Project Owner, Project Admin, Project Member, Read Incident Grouping Rule], Update: [Project Owner, Project Admin, Edit Incident Grouping Rule]",
                Computed: true,
            },
            "episode_title_template": schema.StringAttribute{
                MarkdownDescription: "Template for generating episode titles. Supports placeholders like {{incidentSeverity}}, {{monitorName}}, {{incidentTitle}}, {{incidentDescription}}. Permissions - Create: [Project Owner, Project Admin, Create Incident Grouping Rule], Read: [Project Owner, Project Admin, Project Member, Read Incident Grouping Rule], Update: [Project Owner, Project Admin, Edit Incident Grouping Rule]",
                Computed: true,
            },
            "episode_description_template": schema.StringAttribute{
                MarkdownDescription: "Template for generating episode descriptions. Supports placeholders like {{incidentSeverity}}, {{monitorName}}, {{incidentTitle}}, {{incidentDescription}}. Permissions - Create: [Project Owner, Project Admin, Create Incident Grouping Rule], Read: [Project Owner, Project Admin, Project Member, Read Incident Grouping Rule], Update: [Project Owner, Project Admin, Edit Incident Grouping Rule]",
                Computed: true,
            },
            "enable_resolve_delay": schema.BoolAttribute{
                MarkdownDescription: "Enable grace period before auto-resolving episode after all incidents resolve. Helps prevent rapid state changes during incident flapping.. Permissions - Create: [Project Owner, Project Admin, Create Incident Grouping Rule], Read: [Project Owner, Project Admin, Project Member, Read Incident Grouping Rule], Update: [Project Owner, Project Admin, Edit Incident Grouping Rule]",
                Computed: true,
            },
            "resolve_delay_minutes": schema.NumberAttribute{
                MarkdownDescription: "Grace period in minutes before auto-resolving an episode after all incidents are resolved. Permissions - Create: [Project Owner, Project Admin, Create Incident Grouping Rule], Read: [Project Owner, Project Admin, Project Member, Read Incident Grouping Rule], Update: [Project Owner, Project Admin, Edit Incident Grouping Rule]",
                Computed: true,
            },
            "enable_reopen_window": schema.BoolAttribute{
                MarkdownDescription: "Enable reopening recently resolved episodes instead of creating new ones. Useful when related issues recur shortly after resolution.. Permissions - Create: [Project Owner, Project Admin, Create Incident Grouping Rule], Read: [Project Owner, Project Admin, Project Member, Read Incident Grouping Rule], Update: [Project Owner, Project Admin, Edit Incident Grouping Rule]",
                Computed: true,
            },
            "reopen_window_minutes": schema.NumberAttribute{
                MarkdownDescription: "Time window in minutes to reopen a recently resolved episode instead of creating a new one. Permissions - Create: [Project Owner, Project Admin, Create Incident Grouping Rule], Read: [Project Owner, Project Admin, Project Member, Read Incident Grouping Rule], Update: [Project Owner, Project Admin, Edit Incident Grouping Rule]",
                Computed: true,
            },
            "enable_inactivity_timeout": schema.BoolAttribute{
                MarkdownDescription: "Enable auto-resolving episodes after a period of inactivity. Helps automatically close episodes when no new incidents arrive.. Permissions - Create: [Project Owner, Project Admin, Create Incident Grouping Rule], Read: [Project Owner, Project Admin, Project Member, Read Incident Grouping Rule], Update: [Project Owner, Project Admin, Edit Incident Grouping Rule]",
                Computed: true,
            },
            "inactivity_timeout_minutes": schema.NumberAttribute{
                MarkdownDescription: "Time in minutes after which an inactive episode will be auto-resolved. Permissions - Create: [Project Owner, Project Admin, Create Incident Grouping Rule], Read: [Project Owner, Project Admin, Project Member, Read Incident Grouping Rule], Update: [Project Owner, Project Admin, Edit Incident Grouping Rule]",
                Computed: true,
            },
            "on_call_duty_policies": schema.SetAttribute{
                MarkdownDescription: "List of on-call duty policies to execute for episodes created by this rule.. Permissions - Create: [Project Owner, Project Admin, Create Incident Grouping Rule], Read: [Project Owner, Project Admin, Project Member, Read Incident Grouping Rule], Update: [Project Owner, Project Admin, Edit Incident Grouping Rule]",
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
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
        },
    }
}

func (d *IncidentGroupingRuleDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *IncidentGroupingRuleDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data IncidentGroupingRuleDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "incident-grouping-rule" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read incident_grouping_rule_data, got error: %s", err))
        return
    }

    var incidentGroupingRuleDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &incidentGroupingRuleDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse incident_grouping_rule_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := incidentGroupingRuleDataResponse["data"].(map[string]interface{}); ok {
        incidentGroupingRuleDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := incidentGroupingRuleDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := incidentGroupingRuleDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := incidentGroupingRuleDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := incidentGroupingRuleDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := incidentGroupingRuleDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := incidentGroupingRuleDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := incidentGroupingRuleDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := incidentGroupingRuleDataResponse["description"].(string); ok {
        data.Description = types.StringValue(val)
    }
    if val, ok := incidentGroupingRuleDataResponse["priority"].(float64); ok {
        data.Priority = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := incidentGroupingRuleDataResponse["is_enabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    }
    if val, ok := incidentGroupingRuleDataResponse["match_criteria"].(string); ok {
        data.MatchCriteria = types.StringValue(val)
    }
    if val, ok := incidentGroupingRuleDataResponse["monitors"].([]interface{}); ok {
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
    if val, ok := incidentGroupingRuleDataResponse["incident_severities"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.IncidentSeverities = setValue
    }
    if val, ok := incidentGroupingRuleDataResponse["incident_labels"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.IncidentLabels = setValue
    }
    if val, ok := incidentGroupingRuleDataResponse["monitor_labels"].([]interface{}); ok {
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
    if val, ok := incidentGroupingRuleDataResponse["incident_title_pattern"].(string); ok {
        data.IncidentTitlePattern = types.StringValue(val)
    }
    if val, ok := incidentGroupingRuleDataResponse["incident_description_pattern"].(string); ok {
        data.IncidentDescriptionPattern = types.StringValue(val)
    }
    if val, ok := incidentGroupingRuleDataResponse["monitor_name_pattern"].(string); ok {
        data.MonitorNamePattern = types.StringValue(val)
    }
    if val, ok := incidentGroupingRuleDataResponse["monitor_description_pattern"].(string); ok {
        data.MonitorDescriptionPattern = types.StringValue(val)
    }
    if val, ok := incidentGroupingRuleDataResponse["group_by_monitor"].(bool); ok {
        data.GroupByMonitor = types.BoolValue(val)
    }
    if val, ok := incidentGroupingRuleDataResponse["group_by_severity"].(bool); ok {
        data.GroupBySeverity = types.BoolValue(val)
    }
    if val, ok := incidentGroupingRuleDataResponse["group_by_incident_title"].(bool); ok {
        data.GroupByIncidentTitle = types.BoolValue(val)
    }
    if val, ok := incidentGroupingRuleDataResponse["group_by_service"].(bool); ok {
        data.GroupByService = types.BoolValue(val)
    }
    if val, ok := incidentGroupingRuleDataResponse["enable_time_window"].(bool); ok {
        data.EnableTimeWindow = types.BoolValue(val)
    }
    if val, ok := incidentGroupingRuleDataResponse["time_window_minutes"].(float64); ok {
        data.TimeWindowMinutes = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := incidentGroupingRuleDataResponse["group_by_fields"].(string); ok {
        data.GroupByFields = types.StringValue(val)
    }
    if val, ok := incidentGroupingRuleDataResponse["episode_title_template"].(string); ok {
        data.EpisodeTitleTemplate = types.StringValue(val)
    }
    if val, ok := incidentGroupingRuleDataResponse["episode_description_template"].(string); ok {
        data.EpisodeDescriptionTemplate = types.StringValue(val)
    }
    if val, ok := incidentGroupingRuleDataResponse["enable_resolve_delay"].(bool); ok {
        data.EnableResolveDelay = types.BoolValue(val)
    }
    if val, ok := incidentGroupingRuleDataResponse["resolve_delay_minutes"].(float64); ok {
        data.ResolveDelayMinutes = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := incidentGroupingRuleDataResponse["enable_reopen_window"].(bool); ok {
        data.EnableReopenWindow = types.BoolValue(val)
    }
    if val, ok := incidentGroupingRuleDataResponse["reopen_window_minutes"].(float64); ok {
        data.ReopenWindowMinutes = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := incidentGroupingRuleDataResponse["enable_inactivity_timeout"].(bool); ok {
        data.EnableInactivityTimeout = types.BoolValue(val)
    }
    if val, ok := incidentGroupingRuleDataResponse["inactivity_timeout_minutes"].(float64); ok {
        data.InactivityTimeoutMinutes = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := incidentGroupingRuleDataResponse["on_call_duty_policies"].([]interface{}); ok {
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
    if val, ok := incidentGroupingRuleDataResponse["default_assign_to_user_id"].(string); ok {
        data.DefaultAssignToUserId = types.StringValue(val)
    }
    if val, ok := incidentGroupingRuleDataResponse["default_assign_to_team_id"].(string); ok {
        data.DefaultAssignToTeamId = types.StringValue(val)
    }
    if val, ok := incidentGroupingRuleDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
