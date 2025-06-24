package provider

import (
    "context"
    "fmt"
    "math/big"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &IncidentNoteTemplateDataDataSource{}

func NewIncidentNoteTemplateDataDataSource() datasource.DataSource {
    return &IncidentNoteTemplateDataDataSource{}
}

// IncidentNoteTemplateDataDataSource defines the data source implementation.
type IncidentNoteTemplateDataDataSource struct {
    client *Client
}

// IncidentNoteTemplateDataDataSourceModel describes the data source data model.
type IncidentNoteTemplateDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    Note types.String `tfsdk:"note"`
    TemplateName types.String `tfsdk:"template_name"`
    TemplateDescription types.String `tfsdk:"template_description"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *IncidentNoteTemplateDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_incident_note_template_data"
}

func (d *IncidentNoteTemplateDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "incident_note_template_data data source",

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
            "note": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Incident Note Template], Read: [Project Owner, Project Admin, Project Member, Read Incident Note Template], Update: [Project Owner, Project Admin, Project Member, Edit Incident Note Template]",
                Computed: true,
            },
            "template_name": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Incident Note Template], Read: [Project Owner, Project Admin, Project Member, Read Incident Note Template], Update: [Project Owner, Project Admin, Project Member, Edit Incident Note Template]",
                Computed: true,
            },
            "template_description": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Incident Note Template], Read: [Project Owner, Project Admin, Project Member, Read Incident Note Template], Update: [Project Owner, Project Admin, Project Member, Edit Incident Note Template]",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
        },
    }
}

func (d *IncidentNoteTemplateDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *IncidentNoteTemplateDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data IncidentNoteTemplateDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "incident-note-template" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read incident_note_template_data, got error: %s", err))
        return
    }

    var incidentNoteTemplateDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &incidentNoteTemplateDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse incident_note_template_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := incidentNoteTemplateDataResponse["data"].(map[string]interface{}); ok {
        incidentNoteTemplateDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := incidentNoteTemplateDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := incidentNoteTemplateDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := incidentNoteTemplateDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := incidentNoteTemplateDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := incidentNoteTemplateDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := incidentNoteTemplateDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := incidentNoteTemplateDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := incidentNoteTemplateDataResponse["note"].(string); ok {
        data.Note = types.StringValue(val)
    }
    if val, ok := incidentNoteTemplateDataResponse["template_name"].(string); ok {
        data.TemplateName = types.StringValue(val)
    }
    if val, ok := incidentNoteTemplateDataResponse["template_description"].(string); ok {
        data.TemplateDescription = types.StringValue(val)
    }
    if val, ok := incidentNoteTemplateDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
