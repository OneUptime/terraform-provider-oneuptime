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
var _ datasource.DataSource = &TelemetryEntityRelationshipDataDataSource{}

func NewTelemetryEntityRelationshipDataDataSource() datasource.DataSource {
    return &TelemetryEntityRelationshipDataDataSource{}
}

// TelemetryEntityRelationshipDataDataSource defines the data source implementation.
type TelemetryEntityRelationshipDataDataSource struct {
    client *Client
}

// TelemetryEntityRelationshipDataDataSourceModel describes the data source data model.
type TelemetryEntityRelationshipDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    FromEntityKey types.String `tfsdk:"from_entity_key"`
    ToEntityKey types.String `tfsdk:"to_entity_key"`
    RelationshipType types.String `tfsdk:"relationship_type"`
    FirstSeenAt types.String `tfsdk:"first_seen_at"`
    LastSeenAt types.String `tfsdk:"last_seen_at"`
    CallCount types.Number `tfsdk:"call_count"`
    ErrorCount types.Number `tfsdk:"error_count"`
    AvgDurationMs types.Number `tfsdk:"avg_duration_ms"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
}

func (d *TelemetryEntityRelationshipDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_telemetry_entity_relationship_data"
}

func (d *TelemetryEntityRelationshipDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "telemetry_entity_relationship_data data source",

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
            "from_entity_key": schema.StringAttribute{
                MarkdownDescription: "Stable identity key of the source entity of this edge.. Permissions - Create: [Project Owner, Project Admin, Telemetry Admin, Create Telemetry Service], Read: [Project Owner, Project Admin, Project Member, Viewer, Telemetry Admin, Telemetry Member, Telemetry Viewer, Read Telemetry Service], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "to_entity_key": schema.StringAttribute{
                MarkdownDescription: "Stable identity key of the target entity of this edge.. Permissions - Create: [Project Owner, Project Admin, Telemetry Admin, Create Telemetry Service], Read: [Project Owner, Project Admin, Project Member, Viewer, Telemetry Admin, Telemetry Member, Telemetry Viewer, Read Telemetry Service], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "relationship_type": schema.StringAttribute{
                MarkdownDescription: "The inferred relationship (runs-on, member-of, hosted-on, part-of, instance-of).. Permissions - Create: [Project Owner, Project Admin, Telemetry Admin, Create Telemetry Service], Read: [Project Owner, Project Admin, Project Member, Viewer, Telemetry Admin, Telemetry Member, Telemetry Viewer, Read Telemetry Service], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "first_seen_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "last_seen_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "call_count": schema.NumberAttribute{
                MarkdownDescription: "Calls observed over this edge in the most recent computation window (depends-on edges only).. Permissions - Create: [Project Owner, Project Admin, Telemetry Admin, Create Telemetry Service], Read: [Project Owner, Project Admin, Project Member, Viewer, Telemetry Admin, Telemetry Member, Telemetry Viewer, Read Telemetry Service], Update: [Project Owner, Project Admin, Telemetry Admin, Edit Telemetry Service]",
                Computed: true,
            },
            "error_count": schema.NumberAttribute{
                MarkdownDescription: "Errored calls observed over this edge in the most recent computation window (depends-on edges only).. Permissions - Create: [Project Owner, Project Admin, Telemetry Admin, Create Telemetry Service], Read: [Project Owner, Project Admin, Project Member, Viewer, Telemetry Admin, Telemetry Member, Telemetry Viewer, Read Telemetry Service], Update: [Project Owner, Project Admin, Telemetry Admin, Edit Telemetry Service]",
                Computed: true,
            },
            "avg_duration_ms": schema.NumberAttribute{
                MarkdownDescription: "Average call duration in milliseconds over this edge in the most recent computation window (depends-on edges only).. Permissions - Create: [Project Owner, Project Admin, Telemetry Admin, Create Telemetry Service], Read: [Project Owner, Project Admin, Project Member, Viewer, Telemetry Admin, Telemetry Member, Telemetry Viewer, Read Telemetry Service], Update: [Project Owner, Project Admin, Telemetry Admin, Edit Telemetry Service]",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "deleted_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
        },
    }
}

func (d *TelemetryEntityRelationshipDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *TelemetryEntityRelationshipDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data TelemetryEntityRelationshipDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "telemetry-entity-relationship" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read telemetry_entity_relationship_data, got error: %s", err))
        return
    }

    var telemetryEntityRelationshipDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &telemetryEntityRelationshipDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse telemetry_entity_relationship_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := telemetryEntityRelationshipDataResponse["data"].(map[string]interface{}); ok {
        telemetryEntityRelationshipDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := telemetryEntityRelationshipDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := telemetryEntityRelationshipDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := telemetryEntityRelationshipDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := telemetryEntityRelationshipDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := telemetryEntityRelationshipDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := telemetryEntityRelationshipDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := telemetryEntityRelationshipDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := telemetryEntityRelationshipDataResponse["from_entity_key"].(string); ok {
        data.FromEntityKey = types.StringValue(val)
    }
    if val, ok := telemetryEntityRelationshipDataResponse["to_entity_key"].(string); ok {
        data.ToEntityKey = types.StringValue(val)
    }
    if val, ok := telemetryEntityRelationshipDataResponse["relationship_type"].(string); ok {
        data.RelationshipType = types.StringValue(val)
    }
    if val, ok := telemetryEntityRelationshipDataResponse["first_seen_at"].(string); ok {
        data.FirstSeenAt = types.StringValue(val)
    }
    if val, ok := telemetryEntityRelationshipDataResponse["last_seen_at"].(string); ok {
        data.LastSeenAt = types.StringValue(val)
    }
    if val, ok := telemetryEntityRelationshipDataResponse["call_count"].(float64); ok {
        data.CallCount = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := telemetryEntityRelationshipDataResponse["error_count"].(float64); ok {
        data.ErrorCount = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := telemetryEntityRelationshipDataResponse["avg_duration_ms"].(float64); ok {
        data.AvgDurationMs = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := telemetryEntityRelationshipDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := telemetryEntityRelationshipDataResponse["deleted_by_user_id"].(string); ok {
        data.DeletedByUserId = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
