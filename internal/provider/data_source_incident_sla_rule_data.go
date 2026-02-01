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
var _ datasource.DataSource = &IncidentSlaRuleDataDataSource{}

func NewIncidentSlaRuleDataDataSource() datasource.DataSource {
    return &IncidentSlaRuleDataDataSource{}
}

// IncidentSlaRuleDataDataSource defines the data source implementation.
type IncidentSlaRuleDataDataSource struct {
    client *Client
}

// IncidentSlaRuleDataDataSourceModel describes the data source data model.
type IncidentSlaRuleDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    Description types.String `tfsdk:"description"`
    Order types.Number `tfsdk:"order"`
    IsEnabled types.Bool `tfsdk:"is_enabled"`
    ResponseTimeInMinutes types.Number `tfsdk:"response_time_in_minutes"`
    ResolutionTimeInMinutes types.Number `tfsdk:"resolution_time_in_minutes"`
    AtRiskThresholdInPercentage types.Number `tfsdk:"at_risk_threshold_in_percentage"`
    InternalNoteReminderIntervalInMinutes types.Number `tfsdk:"internal_note_reminder_interval_in_minutes"`
    PublicNoteReminderIntervalInMinutes types.Number `tfsdk:"public_note_reminder_interval_in_minutes"`
    InternalNoteReminderTemplate types.String `tfsdk:"internal_note_reminder_template"`
    PublicNoteReminderTemplate types.String `tfsdk:"public_note_reminder_template"`
    Monitors types.Set `tfsdk:"monitors"`
    IncidentSeverities types.Set `tfsdk:"incident_severities"`
    IncidentLabels types.Set `tfsdk:"incident_labels"`
    MonitorLabels types.Set `tfsdk:"monitor_labels"`
    IncidentTitlePattern types.String `tfsdk:"incident_title_pattern"`
    IncidentDescriptionPattern types.String `tfsdk:"incident_description_pattern"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *IncidentSlaRuleDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_incident_sla_rule_data"
}

func (d *IncidentSlaRuleDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "incident_sla_rule_data data source",

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
                MarkdownDescription: "Description of this SLA rule. Permissions - Create: [Project Owner, Project Admin, Create Incident SLA Rule], Read: [Project Owner, Project Admin, Project Member, Read Incident SLA Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Incident SLA Rule]",
                Computed: true,
            },
            "order": schema.NumberAttribute{
                MarkdownDescription: "Order/priority of this rule. Rules are evaluated in order (lowest first). First matching rule wins.. Permissions - Create: [Project Owner, Project Admin, Create Incident SLA Rule], Read: [Project Owner, Project Admin, Project Member, Read Incident SLA Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Incident SLA Rule]",
                Computed: true,
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether this SLA rule is enabled. Permissions - Create: [Project Owner, Project Admin, Create Incident SLA Rule], Read: [Project Owner, Project Admin, Project Member, Read Incident SLA Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Incident SLA Rule]",
                Computed: true,
            },
            "response_time_in_minutes": schema.NumberAttribute{
                MarkdownDescription: "Target response time in minutes. This is the maximum time allowed before the incident must be acknowledged.. Permissions - Create: [Project Owner, Project Admin, Create Incident SLA Rule], Read: [Project Owner, Project Admin, Project Member, Read Incident SLA Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Incident SLA Rule]",
                Computed: true,
            },
            "resolution_time_in_minutes": schema.NumberAttribute{
                MarkdownDescription: "Target resolution time in minutes. This is the maximum time allowed before the incident must be resolved.. Permissions - Create: [Project Owner, Project Admin, Create Incident SLA Rule], Read: [Project Owner, Project Admin, Project Member, Read Incident SLA Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Incident SLA Rule]",
                Computed: true,
            },
            "at_risk_threshold_in_percentage": schema.NumberAttribute{
                MarkdownDescription: "Percentage of the deadline at which the SLA status changes to At Risk. For example, 80 means the status becomes At Risk when 80% of the time has elapsed.. Permissions - Create: [Project Owner, Project Admin, Create Incident SLA Rule], Read: [Project Owner, Project Admin, Project Member, Read Incident SLA Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Incident SLA Rule]",
                Computed: true,
            },
            "internal_note_reminder_interval_in_minutes": schema.NumberAttribute{
                MarkdownDescription: "How often (in minutes) to automatically post internal notes to unresolved incidents. Internal notes are only visible to your team. For example, set to 30 to remind your team every 30 minutes to provide an update. Leave empty to disable.. Permissions - Create: [Project Owner, Project Admin, Create Incident SLA Rule], Read: [Project Owner, Project Admin, Project Member, Read Incident SLA Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Incident SLA Rule]",
                Computed: true,
            },
            "public_note_reminder_interval_in_minutes": schema.NumberAttribute{
                MarkdownDescription: "How often (in minutes) to automatically post public notes to unresolved incidents. Public notes are visible to external stakeholders on your status page. For example, set to 60 to post a status update every hour. Leave empty to disable.. Permissions - Create: [Project Owner, Project Admin, Create Incident SLA Rule], Read: [Project Owner, Project Admin, Project Member, Read Incident SLA Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Incident SLA Rule]",
                Computed: true,
            },
            "internal_note_reminder_template": schema.StringAttribute{
                MarkdownDescription: "The content of the automatic internal note posted to your team. Use variables like {{incidentTitle}}, {{elapsedTime}}, {{slaStatus}}, {{timeToResolutionDeadline}} to include dynamic incident data. If left empty, a default template will be used.. Permissions - Create: [Project Owner, Project Admin, Create Incident SLA Rule], Read: [Project Owner, Project Admin, Project Member, Read Incident SLA Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Incident SLA Rule]",
                Computed: true,
            },
            "public_note_reminder_template": schema.StringAttribute{
                MarkdownDescription: "The content of the automatic public note shown on your status page. Use variables like {{incidentTitle}}, {{elapsedTime}}, {{slaStatus}}, {{timeToResolutionDeadline}} to include dynamic incident data. If left empty, a default template will be used.. Permissions - Create: [Project Owner, Project Admin, Create Incident SLA Rule], Read: [Project Owner, Project Admin, Project Member, Read Incident SLA Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Incident SLA Rule]",
                Computed: true,
            },
            "monitors": schema.SetAttribute{
                MarkdownDescription: "Only apply this SLA rule to incidents affecting these monitors. Leave empty to match incidents from any monitor.. Permissions - Create: [Project Owner, Project Admin, Create Incident SLA Rule], Read: [Project Owner, Project Admin, Project Member, Read Incident SLA Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Incident SLA Rule]",
                Computed: true,
                ElementType: types.StringType,
            },
            "incident_severities": schema.SetAttribute{
                MarkdownDescription: "Only apply this SLA rule to incidents with these severities. Leave empty to match incidents of any severity.. Permissions - Create: [Project Owner, Project Admin, Create Incident SLA Rule], Read: [Project Owner, Project Admin, Project Member, Read Incident SLA Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Incident SLA Rule]",
                Computed: true,
                ElementType: types.StringType,
            },
            "incident_labels": schema.SetAttribute{
                MarkdownDescription: "Only apply this SLA rule to incidents that have at least one of these labels. Leave empty to match incidents regardless of labels.. Permissions - Create: [Project Owner, Project Admin, Create Incident SLA Rule], Read: [Project Owner, Project Admin, Project Member, Read Incident SLA Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Incident SLA Rule]",
                Computed: true,
                ElementType: types.StringType,
            },
            "monitor_labels": schema.SetAttribute{
                MarkdownDescription: "Only apply this SLA rule to incidents from monitors that have at least one of these labels. Leave empty to match incidents regardless of monitor labels.. Permissions - Create: [Project Owner, Project Admin, Create Incident SLA Rule], Read: [Project Owner, Project Admin, Project Member, Read Incident SLA Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Incident SLA Rule]",
                Computed: true,
                ElementType: types.StringType,
            },
            "incident_title_pattern": schema.StringAttribute{
                MarkdownDescription: "Regular expression pattern to match incident titles. Leave empty to match any title. Example: 'CPU.*high' matches titles containing 'CPU' followed by 'high'.. Permissions - Create: [Project Owner, Project Admin, Create Incident SLA Rule], Read: [Project Owner, Project Admin, Project Member, Read Incident SLA Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Incident SLA Rule]",
                Computed: true,
            },
            "incident_description_pattern": schema.StringAttribute{
                MarkdownDescription: "Regular expression pattern to match incident descriptions. Leave empty to match any description.. Permissions - Create: [Project Owner, Project Admin, Create Incident SLA Rule], Read: [Project Owner, Project Admin, Project Member, Read Incident SLA Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Incident SLA Rule]",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
        },
    }
}

func (d *IncidentSlaRuleDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *IncidentSlaRuleDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data IncidentSlaRuleDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "incident-sla-rule" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read incident_sla_rule_data, got error: %s", err))
        return
    }

    var incidentSlaRuleDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &incidentSlaRuleDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse incident_sla_rule_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := incidentSlaRuleDataResponse["data"].(map[string]interface{}); ok {
        incidentSlaRuleDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := incidentSlaRuleDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := incidentSlaRuleDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := incidentSlaRuleDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := incidentSlaRuleDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := incidentSlaRuleDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := incidentSlaRuleDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := incidentSlaRuleDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := incidentSlaRuleDataResponse["description"].(string); ok {
        data.Description = types.StringValue(val)
    }
    if val, ok := incidentSlaRuleDataResponse["order"].(float64); ok {
        data.Order = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := incidentSlaRuleDataResponse["is_enabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    }
    if val, ok := incidentSlaRuleDataResponse["response_time_in_minutes"].(float64); ok {
        data.ResponseTimeInMinutes = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := incidentSlaRuleDataResponse["resolution_time_in_minutes"].(float64); ok {
        data.ResolutionTimeInMinutes = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := incidentSlaRuleDataResponse["at_risk_threshold_in_percentage"].(float64); ok {
        data.AtRiskThresholdInPercentage = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := incidentSlaRuleDataResponse["internal_note_reminder_interval_in_minutes"].(float64); ok {
        data.InternalNoteReminderIntervalInMinutes = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := incidentSlaRuleDataResponse["public_note_reminder_interval_in_minutes"].(float64); ok {
        data.PublicNoteReminderIntervalInMinutes = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := incidentSlaRuleDataResponse["internal_note_reminder_template"].(string); ok {
        data.InternalNoteReminderTemplate = types.StringValue(val)
    }
    if val, ok := incidentSlaRuleDataResponse["public_note_reminder_template"].(string); ok {
        data.PublicNoteReminderTemplate = types.StringValue(val)
    }
    if val, ok := incidentSlaRuleDataResponse["monitors"].([]interface{}); ok {
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
    if val, ok := incidentSlaRuleDataResponse["incident_severities"].([]interface{}); ok {
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
    if val, ok := incidentSlaRuleDataResponse["incident_labels"].([]interface{}); ok {
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
    if val, ok := incidentSlaRuleDataResponse["monitor_labels"].([]interface{}); ok {
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
    if val, ok := incidentSlaRuleDataResponse["incident_title_pattern"].(string); ok {
        data.IncidentTitlePattern = types.StringValue(val)
    }
    if val, ok := incidentSlaRuleDataResponse["incident_description_pattern"].(string); ok {
        data.IncidentDescriptionPattern = types.StringValue(val)
    }
    if val, ok := incidentSlaRuleDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
