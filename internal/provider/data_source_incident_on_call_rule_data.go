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
var _ datasource.DataSource = &IncidentOnCallRuleDataDataSource{}

func NewIncidentOnCallRuleDataDataSource() datasource.DataSource {
    return &IncidentOnCallRuleDataDataSource{}
}

// IncidentOnCallRuleDataDataSource defines the data source implementation.
type IncidentOnCallRuleDataDataSource struct {
    client *Client
}

// IncidentOnCallRuleDataDataSourceModel describes the data source data model.
type IncidentOnCallRuleDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    Description types.String `tfsdk:"description"`
    IsEnabled types.Bool `tfsdk:"is_enabled"`
    Monitors types.Set `tfsdk:"monitors"`
    IncidentSeverities types.Set `tfsdk:"incident_severities"`
    IncidentLabels types.Set `tfsdk:"incident_labels"`
    MonitorLabels types.Set `tfsdk:"monitor_labels"`
    IncidentTitlePattern types.String `tfsdk:"incident_title_pattern"`
    IncidentDescriptionPattern types.String `tfsdk:"incident_description_pattern"`
    MonitorNamePattern types.String `tfsdk:"monitor_name_pattern"`
    MonitorDescriptionPattern types.String `tfsdk:"monitor_description_pattern"`
    OnCallDutyPolicies types.Set `tfsdk:"on_call_duty_policies"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *IncidentOnCallRuleDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_incident_on_call_rule_data"
}

func (d *IncidentOnCallRuleDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "incident_on_call_rule_data data source",

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
                MarkdownDescription: "Description of this incident on-call rule. Permissions - Create: [Project Owner, Project Admin, Create Incident On-Call Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident On-Call Rule], Update: [Project Owner, Project Admin, Edit Incident On-Call Rule]",
                Computed: true,
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether this rule is enabled. Permissions - Create: [Project Owner, Project Admin, Create Incident On-Call Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident On-Call Rule], Update: [Project Owner, Project Admin, Edit Incident On-Call Rule]",
                Computed: true,
            },
            "monitors": schema.SetAttribute{
                MarkdownDescription: "Only trigger for incidents from these monitors. Leave empty to match incidents from any monitor.. Permissions - Create: [Project Owner, Project Admin, Create Incident On-Call Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident On-Call Rule], Update: [Project Owner, Project Admin, Edit Incident On-Call Rule]",
                Computed: true,
                ElementType: types.StringType,
            },
            "incident_severities": schema.SetAttribute{
                MarkdownDescription: "Only trigger for incidents with these severities. Leave empty to match incidents of any severity.. Permissions - Create: [Project Owner, Project Admin, Create Incident On-Call Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident On-Call Rule], Update: [Project Owner, Project Admin, Edit Incident On-Call Rule]",
                Computed: true,
                ElementType: types.StringType,
            },
            "incident_labels": schema.SetAttribute{
                MarkdownDescription: "Only trigger for incidents that have at least one of these labels. Leave empty to match incidents regardless of incident labels.. Permissions - Create: [Project Owner, Project Admin, Create Incident On-Call Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident On-Call Rule], Update: [Project Owner, Project Admin, Edit Incident On-Call Rule]",
                Computed: true,
                ElementType: types.StringType,
            },
            "monitor_labels": schema.SetAttribute{
                MarkdownDescription: "Only trigger for incidents from monitors that have at least one of these labels. Leave empty to match incidents regardless of monitor labels.. Permissions - Create: [Project Owner, Project Admin, Create Incident On-Call Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident On-Call Rule], Update: [Project Owner, Project Admin, Edit Incident On-Call Rule]",
                Computed: true,
                ElementType: types.StringType,
            },
            "incident_title_pattern": schema.StringAttribute{
                MarkdownDescription: "Regular expression pattern to match incident titles. Leave empty to match any title. Example: 'CPU.*high' matches titles containing 'CPU' followed by 'high'.. Permissions - Create: [Project Owner, Project Admin, Create Incident On-Call Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident On-Call Rule], Update: [Project Owner, Project Admin, Edit Incident On-Call Rule]",
                Computed: true,
            },
            "incident_description_pattern": schema.StringAttribute{
                MarkdownDescription: "Regular expression pattern to match incident descriptions. Leave empty to match any description.. Permissions - Create: [Project Owner, Project Admin, Create Incident On-Call Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident On-Call Rule], Update: [Project Owner, Project Admin, Edit Incident On-Call Rule]",
                Computed: true,
            },
            "monitor_name_pattern": schema.StringAttribute{
                MarkdownDescription: "Regular expression pattern to match monitor names. Leave empty to match any monitor name. Example: 'prod-.*' matches monitors starting with 'prod-'.. Permissions - Create: [Project Owner, Project Admin, Create Incident On-Call Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident On-Call Rule], Update: [Project Owner, Project Admin, Edit Incident On-Call Rule]",
                Computed: true,
            },
            "monitor_description_pattern": schema.StringAttribute{
                MarkdownDescription: "Regular expression pattern to match monitor descriptions. Leave empty to match any monitor description.. Permissions - Create: [Project Owner, Project Admin, Create Incident On-Call Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident On-Call Rule], Update: [Project Owner, Project Admin, Edit Incident On-Call Rule]",
                Computed: true,
            },
            "on_call_duty_policies": schema.SetAttribute{
                MarkdownDescription: "On-call duty policies to execute when an incident matches this rule.. Permissions - Create: [Project Owner, Project Admin, Create Incident On-Call Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident On-Call Rule], Update: [Project Owner, Project Admin, Edit Incident On-Call Rule]",
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

func (d *IncidentOnCallRuleDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *IncidentOnCallRuleDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data IncidentOnCallRuleDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "incident-on-call-rule" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read incident_on_call_rule_data, got error: %s", err))
        return
    }

    var incidentOnCallRuleDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &incidentOnCallRuleDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse incident_on_call_rule_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := incidentOnCallRuleDataResponse["data"].(map[string]interface{}); ok {
        incidentOnCallRuleDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := incidentOnCallRuleDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := incidentOnCallRuleDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := incidentOnCallRuleDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := incidentOnCallRuleDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := incidentOnCallRuleDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := incidentOnCallRuleDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := incidentOnCallRuleDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := incidentOnCallRuleDataResponse["description"].(string); ok {
        data.Description = types.StringValue(val)
    }
    if val, ok := incidentOnCallRuleDataResponse["is_enabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    }
    if val, ok := incidentOnCallRuleDataResponse["monitors"].([]interface{}); ok {
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
    if val, ok := incidentOnCallRuleDataResponse["incident_severities"].([]interface{}); ok {
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
    if val, ok := incidentOnCallRuleDataResponse["incident_labels"].([]interface{}); ok {
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
    if val, ok := incidentOnCallRuleDataResponse["monitor_labels"].([]interface{}); ok {
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
    if val, ok := incidentOnCallRuleDataResponse["incident_title_pattern"].(string); ok {
        data.IncidentTitlePattern = types.StringValue(val)
    }
    if val, ok := incidentOnCallRuleDataResponse["incident_description_pattern"].(string); ok {
        data.IncidentDescriptionPattern = types.StringValue(val)
    }
    if val, ok := incidentOnCallRuleDataResponse["monitor_name_pattern"].(string); ok {
        data.MonitorNamePattern = types.StringValue(val)
    }
    if val, ok := incidentOnCallRuleDataResponse["monitor_description_pattern"].(string); ok {
        data.MonitorDescriptionPattern = types.StringValue(val)
    }
    if val, ok := incidentOnCallRuleDataResponse["on_call_duty_policies"].([]interface{}); ok {
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
    if val, ok := incidentOnCallRuleDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
