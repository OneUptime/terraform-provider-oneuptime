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
var _ datasource.DataSource = &IncidentTemplateDataDataSource{}

func NewIncidentTemplateDataDataSource() datasource.DataSource {
    return &IncidentTemplateDataDataSource{}
}

// IncidentTemplateDataDataSource defines the data source implementation.
type IncidentTemplateDataDataSource struct {
    client *Client
}

// IncidentTemplateDataDataSourceModel describes the data source data model.
type IncidentTemplateDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    Title types.String `tfsdk:"title"`
    TemplateName types.String `tfsdk:"template_name"`
    TemplateDescription types.String `tfsdk:"template_description"`
    Description types.String `tfsdk:"description"`
    Slug types.String `tfsdk:"slug"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    Monitors types.Set `tfsdk:"monitors"`
    OnCallDutyPolicies types.Set `tfsdk:"on_call_duty_policies"`
    Labels types.Set `tfsdk:"labels"`
    IncidentSeverityId types.String `tfsdk:"incident_severity_id"`
    ChangeMonitorStatusToId types.String `tfsdk:"change_monitor_status_to_id"`
    InitialIncidentStateId types.String `tfsdk:"initial_incident_state_id"`
    CustomFields types.String `tfsdk:"custom_fields"`
}

func (d *IncidentTemplateDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_incident_template_data"
}

func (d *IncidentTemplateDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "incident_template_data data source",

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
                MarkdownDescription: "Object version",
                Computed: true,
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "title": schema.StringAttribute{
                MarkdownDescription: "Title of this incident. Permissions - Create: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Create Incident Template], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident Template, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Edit Incident Template]",
                Computed: true,
            },
            "template_name": schema.StringAttribute{
                MarkdownDescription: "Name of the Incident Template. Permissions - Create: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Create Incident Template], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident Template, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Edit Incident Template]",
                Computed: true,
            },
            "template_description": schema.StringAttribute{
                MarkdownDescription: "Description of the Incident Template. Permissions - Create: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Create Incident Template], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident Template, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Edit Incident Template]",
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Short description of this incident. This is in markdown and will be visible on the status page.. Permissions - Create: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Create Incident Template], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident Template, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Edit Incident Template]",
                Computed: true,
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Friendly globally unique name for your object. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident Template, Read All Project Resources], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "monitors": schema.SetAttribute{
                MarkdownDescription: "List of monitors affected by this incident. Permissions - Create: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Create Incident Template], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident Template, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Edit Incident Template]",
                Computed: true,
                ElementType: types.StringType,
            },
            "on_call_duty_policies": schema.SetAttribute{
                MarkdownDescription: "List of on-call duty policies affected by this incident template.. Permissions - Create: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Create Incident Template], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident Template, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Edit Incident Template]",
                Computed: true,
                ElementType: types.StringType,
            },
            "labels": schema.SetAttribute{
                MarkdownDescription: "Relation to Labels Array where this object is categorized in.. Permissions - Create: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Create Incident Template], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident Template, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Edit Incident Template]",
                Computed: true,
                ElementType: types.StringType,
            },
            "incident_severity_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "change_monitor_status_to_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "initial_incident_state_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "custom_fields": schema.StringAttribute{
                MarkdownDescription: "Custom Fields on this resource.. Permissions - Create: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Create Incident Template], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident Template, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Edit Incident Template]",
                Computed: true,
            },
        },
    }
}

func (d *IncidentTemplateDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *IncidentTemplateDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data IncidentTemplateDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "incident-templates" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read incident_template_data, got error: %s", err))
        return
    }

    var incidentTemplateDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &incidentTemplateDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse incident_template_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := incidentTemplateDataResponse["data"].(map[string]interface{}); ok {
        incidentTemplateDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := incidentTemplateDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := incidentTemplateDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := incidentTemplateDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := incidentTemplateDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := incidentTemplateDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := incidentTemplateDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := incidentTemplateDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := incidentTemplateDataResponse["title"].(string); ok {
        data.Title = types.StringValue(val)
    }
    if val, ok := incidentTemplateDataResponse["template_name"].(string); ok {
        data.TemplateName = types.StringValue(val)
    }
    if val, ok := incidentTemplateDataResponse["template_description"].(string); ok {
        data.TemplateDescription = types.StringValue(val)
    }
    if val, ok := incidentTemplateDataResponse["description"].(string); ok {
        data.Description = types.StringValue(val)
    }
    if val, ok := incidentTemplateDataResponse["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    }
    if val, ok := incidentTemplateDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := incidentTemplateDataResponse["monitors"].([]interface{}); ok {
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
    if val, ok := incidentTemplateDataResponse["on_call_duty_policies"].([]interface{}); ok {
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
    if val, ok := incidentTemplateDataResponse["labels"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.Labels = setValue
    }
    if val, ok := incidentTemplateDataResponse["incident_severity_id"].(string); ok {
        data.IncidentSeverityId = types.StringValue(val)
    }
    if val, ok := incidentTemplateDataResponse["change_monitor_status_to_id"].(string); ok {
        data.ChangeMonitorStatusToId = types.StringValue(val)
    }
    if val, ok := incidentTemplateDataResponse["initial_incident_state_id"].(string); ok {
        data.InitialIncidentStateId = types.StringValue(val)
    }
    if val, ok := incidentTemplateDataResponse["custom_fields"].(string); ok {
        data.CustomFields = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
