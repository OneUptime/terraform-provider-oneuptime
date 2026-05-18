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
var _ datasource.DataSource = &IncidentSlaDataDataSource{}

func NewIncidentSlaDataDataSource() datasource.DataSource {
    return &IncidentSlaDataDataSource{}
}

// IncidentSlaDataDataSource defines the data source implementation.
type IncidentSlaDataDataSource struct {
    client *Client
}

// IncidentSlaDataDataSourceModel describes the data source data model.
type IncidentSlaDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    IncidentId types.String `tfsdk:"incident_id"`
    IncidentSlaRuleId types.String `tfsdk:"incident_sla_rule_id"`
    ResponseDeadline types.String `tfsdk:"response_deadline"`
    ResolutionDeadline types.String `tfsdk:"resolution_deadline"`
    Status types.String `tfsdk:"status"`
    RespondedAt types.String `tfsdk:"responded_at"`
    ResolvedAt types.String `tfsdk:"resolved_at"`
    LastInternalNoteReminderSentAt types.String `tfsdk:"last_internal_note_reminder_sent_at"`
    LastPublicNoteReminderSentAt types.String `tfsdk:"last_public_note_reminder_sent_at"`
    BreachNotificationSentAt types.String `tfsdk:"breach_notification_sent_at"`
    SlaStartedAt types.String `tfsdk:"sla_started_at"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *IncidentSlaDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_incident_sla_data"
}

func (d *IncidentSlaDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "incident_sla_data data source",

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
            "incident_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "incident_sla_rule_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "response_deadline": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "resolution_deadline": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "status": schema.StringAttribute{
                MarkdownDescription: "Current SLA status (On Track, At Risk, Breached, Met). Permissions - Create: [Project Owner, Project Admin, Create Incident SLA], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident SLA], Update: [Project Owner, Project Admin, Edit Incident SLA]",
                Computed: true,
            },
            "responded_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "resolved_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "last_internal_note_reminder_sent_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "last_public_note_reminder_sent_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "breach_notification_sent_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "sla_started_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
        },
    }
}

func (d *IncidentSlaDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *IncidentSlaDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data IncidentSlaDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "incident-sla" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read incident_sla_data, got error: %s", err))
        return
    }

    var incidentSlaDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &incidentSlaDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse incident_sla_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := incidentSlaDataResponse["data"].(map[string]interface{}); ok {
        incidentSlaDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := incidentSlaDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := incidentSlaDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := incidentSlaDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := incidentSlaDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := incidentSlaDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := incidentSlaDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := incidentSlaDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := incidentSlaDataResponse["incident_id"].(string); ok {
        data.IncidentId = types.StringValue(val)
    }
    if val, ok := incidentSlaDataResponse["incident_sla_rule_id"].(string); ok {
        data.IncidentSlaRuleId = types.StringValue(val)
    }
    if val, ok := incidentSlaDataResponse["response_deadline"].(string); ok {
        data.ResponseDeadline = types.StringValue(val)
    }
    if val, ok := incidentSlaDataResponse["resolution_deadline"].(string); ok {
        data.ResolutionDeadline = types.StringValue(val)
    }
    if val, ok := incidentSlaDataResponse["status"].(string); ok {
        data.Status = types.StringValue(val)
    }
    if val, ok := incidentSlaDataResponse["responded_at"].(string); ok {
        data.RespondedAt = types.StringValue(val)
    }
    if val, ok := incidentSlaDataResponse["resolved_at"].(string); ok {
        data.ResolvedAt = types.StringValue(val)
    }
    if val, ok := incidentSlaDataResponse["last_internal_note_reminder_sent_at"].(string); ok {
        data.LastInternalNoteReminderSentAt = types.StringValue(val)
    }
    if val, ok := incidentSlaDataResponse["last_public_note_reminder_sent_at"].(string); ok {
        data.LastPublicNoteReminderSentAt = types.StringValue(val)
    }
    if val, ok := incidentSlaDataResponse["breach_notification_sent_at"].(string); ok {
        data.BreachNotificationSentAt = types.StringValue(val)
    }
    if val, ok := incidentSlaDataResponse["sla_started_at"].(string); ok {
        data.SlaStartedAt = types.StringValue(val)
    }
    if val, ok := incidentSlaDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
