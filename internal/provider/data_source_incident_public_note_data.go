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
var _ datasource.DataSource = &IncidentPublicNoteDataDataSource{}

func NewIncidentPublicNoteDataDataSource() datasource.DataSource {
    return &IncidentPublicNoteDataDataSource{}
}

// IncidentPublicNoteDataDataSource defines the data source implementation.
type IncidentPublicNoteDataDataSource struct {
    client *Client
}

// IncidentPublicNoteDataDataSourceModel describes the data source data model.
type IncidentPublicNoteDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    IncidentId types.String `tfsdk:"incident_id"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    Note types.String `tfsdk:"note"`
    IsStatusPageSubscribersNotifiedOnNoteCreated types.Bool `tfsdk:"is_status_page_subscribers_notified_on_note_created"`
    ShouldStatusPageSubscribersBeNotifiedOnNoteCreated types.Bool `tfsdk:"should_status_page_subscribers_be_notified_on_note_created"`
    IsOwnerNotified types.Bool `tfsdk:"is_owner_notified"`
    PostedAt types.String `tfsdk:"posted_at"`
}

func (d *IncidentPublicNoteDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_incident_public_note_data"
}

func (d *IncidentPublicNoteDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "incident_public_note_data data source",

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
            "incident_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "note": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Incident Status Page Note], Read: [Project Owner, Project Admin, Project Member, Read Incident Status Page Note], Update: [Project Owner, Project Admin, Project Member, Edit Incident Status Page Note]",
                Computed: true,
            },
            "is_status_page_subscribers_notified_on_note_created": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Incident Status Page Note], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "should_status_page_subscribers_be_notified_on_note_created": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Incident Status Page Note], Read: [Project Owner, Project Admin, Project Member, Read Incident Status Page Note], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "is_owner_notified": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Incident Status Page Note], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "posted_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
        },
    }
}

func (d *IncidentPublicNoteDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *IncidentPublicNoteDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data IncidentPublicNoteDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "incident-public-note" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read incident_public_note_data, got error: %s", err))
        return
    }

    var incidentPublicNoteDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &incidentPublicNoteDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse incident_public_note_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := incidentPublicNoteDataResponse["data"].(map[string]interface{}); ok {
        incidentPublicNoteDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := incidentPublicNoteDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := incidentPublicNoteDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := incidentPublicNoteDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := incidentPublicNoteDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := incidentPublicNoteDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := incidentPublicNoteDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := incidentPublicNoteDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := incidentPublicNoteDataResponse["incident_id"].(string); ok {
        data.IncidentId = types.StringValue(val)
    }
    if val, ok := incidentPublicNoteDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := incidentPublicNoteDataResponse["note"].(string); ok {
        data.Note = types.StringValue(val)
    }
    if val, ok := incidentPublicNoteDataResponse["is_status_page_subscribers_notified_on_note_created"].(bool); ok {
        data.IsStatusPageSubscribersNotifiedOnNoteCreated = types.BoolValue(val)
    }
    if val, ok := incidentPublicNoteDataResponse["should_status_page_subscribers_be_notified_on_note_created"].(bool); ok {
        data.ShouldStatusPageSubscribersBeNotifiedOnNoteCreated = types.BoolValue(val)
    }
    if val, ok := incidentPublicNoteDataResponse["is_owner_notified"].(bool); ok {
        data.IsOwnerNotified = types.BoolValue(val)
    }
    if val, ok := incidentPublicNoteDataResponse["posted_at"].(string); ok {
        data.PostedAt = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
