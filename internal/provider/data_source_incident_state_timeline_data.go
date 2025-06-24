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
var _ datasource.DataSource = &IncidentStateTimelineDataDataSource{}

func NewIncidentStateTimelineDataDataSource() datasource.DataSource {
    return &IncidentStateTimelineDataDataSource{}
}

// IncidentStateTimelineDataDataSource defines the data source implementation.
type IncidentStateTimelineDataDataSource struct {
    client *Client
}

// IncidentStateTimelineDataDataSourceModel describes the data source data model.
type IncidentStateTimelineDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    IncidentId types.String `tfsdk:"incident_id"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    IncidentStateId types.String `tfsdk:"incident_state_id"`
    IsStatusPageSubscribersNotified types.Bool `tfsdk:"is_status_page_subscribers_notified"`
    ShouldStatusPageSubscribersBeNotified types.Bool `tfsdk:"should_status_page_subscribers_be_notified"`
    IsOwnerNotified types.Bool `tfsdk:"is_owner_notified"`
    StateChangeLog types.String `tfsdk:"state_change_log"`
    RootCause types.String `tfsdk:"root_cause"`
    EndsAt types.String `tfsdk:"ends_at"`
    StartsAt types.String `tfsdk:"starts_at"`
}

func (d *IncidentStateTimelineDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_incident_state_timeline_data"
}

func (d *IncidentStateTimelineDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "incident_state_timeline_data data source",

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
            "incident_state_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "is_status_page_subscribers_notified": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Incident State Timeline], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "should_status_page_subscribers_be_notified": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Incident State Timeline], Read: [Project Owner, Project Admin, Project Member, Read Incident State Timeline], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "is_owner_notified": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Incident State Timeline], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "state_change_log": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Incident State Timeline], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "root_cause": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Incident State Timeline], Read: [Project Owner, Project Admin, Project Member, Read Incident State Timeline], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "ends_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "starts_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
        },
    }
}

func (d *IncidentStateTimelineDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *IncidentStateTimelineDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data IncidentStateTimelineDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "incident-state-timeline" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read incident_state_timeline_data, got error: %s", err))
        return
    }

    var incidentStateTimelineDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &incidentStateTimelineDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse incident_state_timeline_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := incidentStateTimelineDataResponse["data"].(map[string]interface{}); ok {
        incidentStateTimelineDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := incidentStateTimelineDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := incidentStateTimelineDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := incidentStateTimelineDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := incidentStateTimelineDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := incidentStateTimelineDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := incidentStateTimelineDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := incidentStateTimelineDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := incidentStateTimelineDataResponse["incident_id"].(string); ok {
        data.IncidentId = types.StringValue(val)
    }
    if val, ok := incidentStateTimelineDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := incidentStateTimelineDataResponse["incident_state_id"].(string); ok {
        data.IncidentStateId = types.StringValue(val)
    }
    if val, ok := incidentStateTimelineDataResponse["is_status_page_subscribers_notified"].(bool); ok {
        data.IsStatusPageSubscribersNotified = types.BoolValue(val)
    }
    if val, ok := incidentStateTimelineDataResponse["should_status_page_subscribers_be_notified"].(bool); ok {
        data.ShouldStatusPageSubscribersBeNotified = types.BoolValue(val)
    }
    if val, ok := incidentStateTimelineDataResponse["is_owner_notified"].(bool); ok {
        data.IsOwnerNotified = types.BoolValue(val)
    }
    if val, ok := incidentStateTimelineDataResponse["state_change_log"].(string); ok {
        data.StateChangeLog = types.StringValue(val)
    }
    if val, ok := incidentStateTimelineDataResponse["root_cause"].(string); ok {
        data.RootCause = types.StringValue(val)
    }
    if val, ok := incidentStateTimelineDataResponse["ends_at"].(string); ok {
        data.EndsAt = types.StringValue(val)
    }
    if val, ok := incidentStateTimelineDataResponse["starts_at"].(string); ok {
        data.StartsAt = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
