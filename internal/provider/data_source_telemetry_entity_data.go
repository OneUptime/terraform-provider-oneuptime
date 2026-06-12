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
var _ datasource.DataSource = &TelemetryEntityDataDataSource{}

func NewTelemetryEntityDataDataSource() datasource.DataSource {
    return &TelemetryEntityDataDataSource{}
}

// TelemetryEntityDataDataSource defines the data source implementation.
type TelemetryEntityDataDataSource struct {
    client *Client
}

// TelemetryEntityDataDataSourceModel describes the data source data model.
type TelemetryEntityDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    EntityType types.String `tfsdk:"entity_type"`
    EntityKey types.String `tfsdk:"entity_key"`
    DisplayName types.String `tfsdk:"display_name"`
    IdentifyingAttributes types.String `tfsdk:"identifying_attributes"`
    DescriptiveAttributes types.String `tfsdk:"descriptive_attributes"`
    Labels types.String `tfsdk:"labels"`
    ResourceType types.String `tfsdk:"resource_type"`
    ResourceId types.String `tfsdk:"resource_id"`
    FirstSeenAt types.String `tfsdk:"first_seen_at"`
    LastSeenAt types.String `tfsdk:"last_seen_at"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
}

func (d *TelemetryEntityDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_telemetry_entity_data"
}

func (d *TelemetryEntityDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "telemetry_entity_data data source",

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
            "entity_type": schema.StringAttribute{
                MarkdownDescription: "The OpenTelemetry entity type (service, host, k8s.pod, container, ...).. Permissions - Create: [Project Owner, Project Admin, Telemetry Admin, Create Telemetry Service], Read: [Project Owner, Project Admin, Project Member, Viewer, Telemetry Admin, Telemetry Member, Telemetry Viewer, Read Telemetry Service], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "entity_key": schema.StringAttribute{
                MarkdownDescription: "Stable identity hash derived from the entity's identifying attributes (matches the keys stamped into signal entityKeys columns).. Permissions - Create: [Project Owner, Project Admin, Telemetry Admin, Create Telemetry Service], Read: [Project Owner, Project Admin, Project Member, Viewer, Telemetry Admin, Telemetry Member, Telemetry Viewer, Read Telemetry Service], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "display_name": schema.StringAttribute{
                MarkdownDescription: "Human-readable name derived for the entity explorer UI.. Permissions - Create: [Project Owner, Project Admin, Telemetry Admin, Create Telemetry Service], Read: [Project Owner, Project Admin, Project Member, Viewer, Telemetry Admin, Telemetry Member, Telemetry Viewer, Read Telemetry Service], Update: [Project Owner, Project Admin, Telemetry Admin, Edit Telemetry Service]",
                Computed: true,
            },
            "identifying_attributes": schema.StringAttribute{
                MarkdownDescription: "The immutable identifying attribute set (the entity's identity). Descriptive attributes are deliberately excluded so they can change without changing the entity key.. Permissions - Create: [Project Owner, Project Admin, Telemetry Admin, Create Telemetry Service], Read: [Project Owner, Project Admin, Project Member, Viewer, Telemetry Admin, Telemetry Member, Telemetry Viewer, Read Telemetry Service], Update: [Project Owner, Project Admin, Telemetry Admin, Edit Telemetry Service]",
                Computed: true,
            },
            "descriptive_attributes": schema.StringAttribute{
                MarkdownDescription: "Mutable descriptive metadata (image tag, version, IP, ...) merged last-writer-wins. Never part of the identity.. Permissions - Create: [Project Owner, Project Admin, Telemetry Admin, Create Telemetry Service], Read: [Project Owner, Project Admin, Project Member, Viewer, Telemetry Admin, Telemetry Member, Telemetry Viewer, Read Telemetry Service], Update: [Project Owner, Project Admin, Telemetry Admin, Edit Telemetry Service]",
                Computed: true,
            },
            "labels": schema.StringAttribute{
                MarkdownDescription: "Labels observed on this entity's telemetry (e.g. promoted from oneuptime.label.* resource attributes), merged as a set union. Simple string array in v1 — a relation to the Label table is a follow-up.. Permissions - Create: [Project Owner, Project Admin, Telemetry Admin, Create Telemetry Service], Read: [Project Owner, Project Admin, Project Member, Viewer, Telemetry Admin, Telemetry Member, Telemetry Viewer, Read Telemetry Service], Update: [Project Owner, Project Admin, Telemetry Admin, Edit Telemetry Service]",
                Computed: true,
            },
            "resource_type": schema.StringAttribute{
                MarkdownDescription: "Polymorphic pointer type to a rich typed row, if one exists (Service / Host / DockerHost / KubernetesCluster).. Permissions - Create: [Project Owner, Project Admin, Telemetry Admin, Create Telemetry Service], Read: [Project Owner, Project Admin, Project Member, Viewer, Telemetry Admin, Telemetry Member, Telemetry Viewer, Read Telemetry Service], Update: [Project Owner, Project Admin, Telemetry Admin, Edit Telemetry Service]",
                Computed: true,
            },
            "resource_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
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

func (d *TelemetryEntityDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *TelemetryEntityDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data TelemetryEntityDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "telemetry-entity" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read telemetry_entity_data, got error: %s", err))
        return
    }

    var telemetryEntityDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &telemetryEntityDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse telemetry_entity_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := telemetryEntityDataResponse["data"].(map[string]interface{}); ok {
        telemetryEntityDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := telemetryEntityDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := telemetryEntityDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := telemetryEntityDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := telemetryEntityDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := telemetryEntityDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := telemetryEntityDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := telemetryEntityDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := telemetryEntityDataResponse["entity_type"].(string); ok {
        data.EntityType = types.StringValue(val)
    }
    if val, ok := telemetryEntityDataResponse["entity_key"].(string); ok {
        data.EntityKey = types.StringValue(val)
    }
    if val, ok := telemetryEntityDataResponse["display_name"].(string); ok {
        data.DisplayName = types.StringValue(val)
    }
    if val, ok := telemetryEntityDataResponse["identifying_attributes"].(string); ok {
        data.IdentifyingAttributes = types.StringValue(val)
    }
    if val, ok := telemetryEntityDataResponse["descriptive_attributes"].(string); ok {
        data.DescriptiveAttributes = types.StringValue(val)
    }
    if val, ok := telemetryEntityDataResponse["labels"].(string); ok {
        data.Labels = types.StringValue(val)
    }
    if val, ok := telemetryEntityDataResponse["resource_type"].(string); ok {
        data.ResourceType = types.StringValue(val)
    }
    if val, ok := telemetryEntityDataResponse["resource_id"].(string); ok {
        data.ResourceId = types.StringValue(val)
    }
    if val, ok := telemetryEntityDataResponse["first_seen_at"].(string); ok {
        data.FirstSeenAt = types.StringValue(val)
    }
    if val, ok := telemetryEntityDataResponse["last_seen_at"].(string); ok {
        data.LastSeenAt = types.StringValue(val)
    }
    if val, ok := telemetryEntityDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := telemetryEntityDataResponse["deleted_by_user_id"].(string); ok {
        data.DeletedByUserId = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
