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
var _ datasource.DataSource = &IncidentDataDataSource{}

func NewIncidentDataDataSource() datasource.DataSource {
    return &IncidentDataDataSource{}
}

// IncidentDataDataSource defines the data source implementation.
type IncidentDataDataSource struct {
    client *Client
}

// IncidentDataDataSourceModel describes the data source data model.
type IncidentDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    Title types.String `tfsdk:"title"`
    Description types.String `tfsdk:"description"`
    Slug types.String `tfsdk:"slug"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    Monitors types.List `tfsdk:"monitors"`
    OnCallDutyPolicies types.List `tfsdk:"on_call_duty_policies"`
    Labels types.List `tfsdk:"labels"`
    CurrentIncidentStateId types.String `tfsdk:"current_incident_state_id"`
    IncidentSeverityId types.String `tfsdk:"incident_severity_id"`
    ChangeMonitorStatusToId types.String `tfsdk:"change_monitor_status_to_id"`
    IsStatusPageSubscribersNotifiedOnIncidentCreated types.Bool `tfsdk:"is_status_page_subscribers_notified_on_incident_created"`
    ShouldStatusPageSubscribersBeNotifiedOnIncidentCreated types.Bool `tfsdk:"should_status_page_subscribers_be_notified_on_incident_created"`
    CustomFields types.String `tfsdk:"custom_fields"`
    IsOwnerNotifiedOfResourceCreation types.Bool `tfsdk:"is_owner_notified_of_resource_creation"`
    RootCause types.String `tfsdk:"root_cause"`
    CreatedStateLog types.String `tfsdk:"created_state_log"`
    CreatedCriteriaId types.String `tfsdk:"created_criteria_id"`
    CreatedIncidentTemplateId types.String `tfsdk:"created_incident_template_id"`
    CreatedByProbeId types.String `tfsdk:"created_by_probe_id"`
    IsCreatedAutomatically types.Bool `tfsdk:"is_created_automatically"`
    RemediationNotes types.String `tfsdk:"remediation_notes"`
    TelemetryQuery types.String `tfsdk:"telemetry_query"`
    IncidentNumber types.Number `tfsdk:"incident_number"`
    IsVisibleOnStatusPage types.Bool `tfsdk:"is_visible_on_status_page"`
}

func (d *IncidentDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_incident_data"
}

func (d *IncidentDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "incident_data data source",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Identifier to filter by",
                Optional: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Name to filter by",
                Optional: true,
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
                MarkdownDescription: "Version",
                Computed: true,
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "title": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Incident], Read: [Project Owner, Project Admin, Project Member, Read Incident], Update: [Project Owner, Project Admin, Project Member, Edit Incident]",
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Incident], Read: [Project Owner, Project Admin, Project Member, Read Incident], Update: [Project Owner, Project Admin, Project Member, Edit Incident]",
                Computed: true,
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Incident], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "monitors": schema.ListAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Incident], Read: [Project Owner, Project Admin, Project Member, Read Incident], Update: [Project Owner, Project Admin, Project Member, Edit Incident]",
                Computed: true,
                ElementType: types.StringType,
            },
            "on_call_duty_policies": schema.ListAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Incident], Read: [Project Owner, Project Admin, Project Member, Read Incident], Update: [Project Owner, Project Admin, Project Member, Edit Incident]",
                Computed: true,
                ElementType: types.StringType,
            },
            "labels": schema.ListAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Incident], Read: [Project Owner, Project Admin, Project Member, Read Incident], Update: [Project Owner, Project Admin, Project Member, Edit Incident]",
                Computed: true,
                ElementType: types.StringType,
            },
            "current_incident_state_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "incident_severity_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "change_monitor_status_to_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "is_status_page_subscribers_notified_on_incident_created": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Incident], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "should_status_page_subscribers_be_notified_on_incident_created": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Incident], Read: [Project Owner, Project Admin, Project Member, Read Incident], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "custom_fields": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Incident], Read: [Project Owner, Project Admin, Project Member, Read Incident], Update: [Project Owner, Project Admin, Project Member, Edit Incident]",
                Computed: true,
            },
            "is_owner_notified_of_resource_creation": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Incident], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "root_cause": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Incident], Read: [Project Owner, Project Admin, Project Member, Read Incident], Update: [Project Owner, Project Admin, Project Member, Edit Incident]",
                Computed: true,
            },
            "created_state_log": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Incident], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "created_criteria_id": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Incident], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "created_incident_template_id": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Incident], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "created_by_probe_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "is_created_automatically": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Incident], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "remediation_notes": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Incident], Read: [Project Owner, Project Admin, Project Member, Read Incident], Update: [Project Owner, Project Admin, Project Member, Edit Incident]",
                Computed: true,
            },
            "telemetry_query": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Incident], Read: [Project Owner, Project Admin, Project Member, Read Incident], Update: [Project Owner, Project Admin, Project Member, Edit Incident]",
                Computed: true,
            },
            "incident_number": schema.NumberAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Incident], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "is_visible_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Incident], Read: [Project Owner, Project Admin, Project Member, Read Incident], Update: [Project Owner, Project Admin, Project Member, Edit Incident]",
                Computed: true,
            },
        },
    }
}

func (d *IncidentDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *IncidentDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data IncidentDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "incident" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read incident_data, got error: %s", err))
        return
    }

    var incidentDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &incidentDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse incident_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := incidentDataResponse["data"].(map[string]interface{}); ok {
        incidentDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := incidentDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := incidentDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := incidentDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := incidentDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := incidentDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := incidentDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := incidentDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := incidentDataResponse["title"].(string); ok {
        data.Title = types.StringValue(val)
    }
    if val, ok := incidentDataResponse["description"].(string); ok {
        data.Description = types.StringValue(val)
    }
    if val, ok := incidentDataResponse["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    }
    if val, ok := incidentDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := incidentDataResponse["monitors"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        listValue, _ := types.ListValue(types.StringType, elements)
        data.Monitors = listValue
    }
    if val, ok := incidentDataResponse["on_call_duty_policies"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        listValue, _ := types.ListValue(types.StringType, elements)
        data.OnCallDutyPolicies = listValue
    }
    if val, ok := incidentDataResponse["labels"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        listValue, _ := types.ListValue(types.StringType, elements)
        data.Labels = listValue
    }
    if val, ok := incidentDataResponse["current_incident_state_id"].(string); ok {
        data.CurrentIncidentStateId = types.StringValue(val)
    }
    if val, ok := incidentDataResponse["incident_severity_id"].(string); ok {
        data.IncidentSeverityId = types.StringValue(val)
    }
    if val, ok := incidentDataResponse["change_monitor_status_to_id"].(string); ok {
        data.ChangeMonitorStatusToId = types.StringValue(val)
    }
    if val, ok := incidentDataResponse["is_status_page_subscribers_notified_on_incident_created"].(bool); ok {
        data.IsStatusPageSubscribersNotifiedOnIncidentCreated = types.BoolValue(val)
    }
    if val, ok := incidentDataResponse["should_status_page_subscribers_be_notified_on_incident_created"].(bool); ok {
        data.ShouldStatusPageSubscribersBeNotifiedOnIncidentCreated = types.BoolValue(val)
    }
    if val, ok := incidentDataResponse["custom_fields"].(string); ok {
        data.CustomFields = types.StringValue(val)
    }
    if val, ok := incidentDataResponse["is_owner_notified_of_resource_creation"].(bool); ok {
        data.IsOwnerNotifiedOfResourceCreation = types.BoolValue(val)
    }
    if val, ok := incidentDataResponse["root_cause"].(string); ok {
        data.RootCause = types.StringValue(val)
    }
    if val, ok := incidentDataResponse["created_state_log"].(string); ok {
        data.CreatedStateLog = types.StringValue(val)
    }
    if val, ok := incidentDataResponse["created_criteria_id"].(string); ok {
        data.CreatedCriteriaId = types.StringValue(val)
    }
    if val, ok := incidentDataResponse["created_incident_template_id"].(string); ok {
        data.CreatedIncidentTemplateId = types.StringValue(val)
    }
    if val, ok := incidentDataResponse["created_by_probe_id"].(string); ok {
        data.CreatedByProbeId = types.StringValue(val)
    }
    if val, ok := incidentDataResponse["is_created_automatically"].(bool); ok {
        data.IsCreatedAutomatically = types.BoolValue(val)
    }
    if val, ok := incidentDataResponse["remediation_notes"].(string); ok {
        data.RemediationNotes = types.StringValue(val)
    }
    if val, ok := incidentDataResponse["telemetry_query"].(string); ok {
        data.TelemetryQuery = types.StringValue(val)
    }
    if val, ok := incidentDataResponse["incident_number"].(float64); ok {
        data.IncidentNumber = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := incidentDataResponse["is_visible_on_status_page"].(bool); ok {
        data.IsVisibleOnStatusPage = types.BoolValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
