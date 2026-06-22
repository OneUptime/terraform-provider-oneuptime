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
var _ datasource.DataSource = &PodmanHostDataDataSource{}

func NewPodmanHostDataDataSource() datasource.DataSource {
    return &PodmanHostDataDataSource{}
}

// PodmanHostDataDataSource defines the data source implementation.
type PodmanHostDataDataSource struct {
    client *Client
}

// PodmanHostDataDataSourceModel describes the data source data model.
type PodmanHostDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    Slug types.String `tfsdk:"slug"`
    Description types.String `tfsdk:"description"`
    HostIdentifier types.String `tfsdk:"host_identifier"`
    OtelCollectorStatus types.String `tfsdk:"otel_collector_status"`
    AgentVersion types.String `tfsdk:"agent_version"`
    LastSeenAt types.String `tfsdk:"last_seen_at"`
    ContainersRunning types.Number `tfsdk:"containers_running"`
    ContainersStopped types.Number `tfsdk:"containers_stopped"`
    ContainersPaused types.Number `tfsdk:"containers_paused"`
    OsType types.String `tfsdk:"os_type"`
    OsVersion types.String `tfsdk:"os_version"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    IsArchived types.Bool `tfsdk:"is_archived"`
    ArchivedAt types.String `tfsdk:"archived_at"`
    ArchivedByUserId types.String `tfsdk:"archived_by_user_id"`
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
    Labels types.Set `tfsdk:"labels"`
    RetainTelemetryDataForDays types.Number `tfsdk:"retain_telemetry_data_for_days"`
    TelemetryRetentionConfig types.String `tfsdk:"telemetry_retention_config"`
}

func (d *PodmanHostDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_podman_host_data"
}

func (d *PodmanHostDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "podman_host_data data source",

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
            "slug": schema.StringAttribute{
                MarkdownDescription: "Friendly globally unique name for your object. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Podman Host], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Friendly description for this Podman host. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Podman Host], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Podman Host], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Podman Host]",
                Computed: true,
            },
            "host_identifier": schema.StringAttribute{
                MarkdownDescription: "Unique identifier for this Podman host, sourced from the host.name OTel resource attribute. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Podman Host], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Podman Host], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Podman Host]",
                Computed: true,
            },
            "otel_collector_status": schema.StringAttribute{
                MarkdownDescription: "Connection status of the OTel Collector agent (connected or disconnected). Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Podman Host], Update: [Project Owner, Project Admin, Edit Podman Host]",
                Computed: true,
            },
            "agent_version": schema.StringAttribute{
                MarkdownDescription: "Version of the OneUptime Podman agent reporting telemetry, as self-reported via the oneuptime.agent.version resource attribute. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Podman Host], Update: [Project Owner, Project Admin, Edit Podman Host]",
                Computed: true,
            },
            "last_seen_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "containers_running": schema.NumberAttribute{
                MarkdownDescription: "Cached count of running containers on this host. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Podman Host], Update: [Project Owner, Project Admin, Edit Podman Host]",
                Computed: true,
            },
            "containers_stopped": schema.NumberAttribute{
                MarkdownDescription: "Cached count of stopped containers on this host. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Podman Host], Update: [Project Owner, Project Admin, Edit Podman Host]",
                Computed: true,
            },
            "containers_paused": schema.NumberAttribute{
                MarkdownDescription: "Cached count of paused containers on this host. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Podman Host], Update: [Project Owner, Project Admin, Edit Podman Host]",
                Computed: true,
            },
            "os_type": schema.StringAttribute{
                MarkdownDescription: "Operating system type of the Podman host. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Podman Host], Update: [Project Owner, Project Admin, Edit Podman Host]",
                Computed: true,
            },
            "os_version": schema.StringAttribute{
                MarkdownDescription: "Operating system version of the Podman host. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Podman Host], Update: [Project Owner, Project Admin, Edit Podman Host]",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "is_archived": schema.BoolAttribute{
                MarkdownDescription: "Is this Podman host archived? Archived Podman hosts are hidden from lists but keep collecting telemetry.. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Podman Host], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Podman Host], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Podman Host]",
                Computed: true,
            },
            "archived_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "archived_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "deleted_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "labels": schema.SetAttribute{
                MarkdownDescription: "Relation to Labels Array where this object is categorized in.. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Podman Host], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Podman Host], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Podman Host]",
                Computed: true,
                ElementType: types.StringType,
            },
            "retain_telemetry_data_for_days": schema.NumberAttribute{
                MarkdownDescription: "Number of days to retain telemetry data for this Podman host. Leave blank to use the project-wide default.. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Podman Host], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Podman Host], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Podman Host]",
                Computed: true,
            },
            "telemetry_retention_config": schema.StringAttribute{
                MarkdownDescription: "Per-pillar retention overrides for this Podman host (logs by severity, traces by status, metrics, profiles). Unset fields fall back to the Podman host default, then the project's retention settings.. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Podman Host], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Podman Host], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Podman Host]",
                Computed: true,
            },
        },
    }
}

func (d *PodmanHostDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *PodmanHostDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data PodmanHostDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "podman-host" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read podman_host_data, got error: %s", err))
        return
    }

    var podmanHostDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &podmanHostDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse podman_host_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := podmanHostDataResponse["data"].(map[string]interface{}); ok {
        podmanHostDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := podmanHostDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := podmanHostDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := podmanHostDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := podmanHostDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := podmanHostDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := podmanHostDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := podmanHostDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := podmanHostDataResponse["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    }
    if val, ok := podmanHostDataResponse["description"].(string); ok {
        data.Description = types.StringValue(val)
    }
    if val, ok := podmanHostDataResponse["host_identifier"].(string); ok {
        data.HostIdentifier = types.StringValue(val)
    }
    if val, ok := podmanHostDataResponse["otel_collector_status"].(string); ok {
        data.OtelCollectorStatus = types.StringValue(val)
    }
    if val, ok := podmanHostDataResponse["agent_version"].(string); ok {
        data.AgentVersion = types.StringValue(val)
    }
    if val, ok := podmanHostDataResponse["last_seen_at"].(string); ok {
        data.LastSeenAt = types.StringValue(val)
    }
    if val, ok := podmanHostDataResponse["containers_running"].(float64); ok {
        data.ContainersRunning = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := podmanHostDataResponse["containers_stopped"].(float64); ok {
        data.ContainersStopped = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := podmanHostDataResponse["containers_paused"].(float64); ok {
        data.ContainersPaused = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := podmanHostDataResponse["os_type"].(string); ok {
        data.OsType = types.StringValue(val)
    }
    if val, ok := podmanHostDataResponse["os_version"].(string); ok {
        data.OsVersion = types.StringValue(val)
    }
    if val, ok := podmanHostDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := podmanHostDataResponse["is_archived"].(bool); ok {
        data.IsArchived = types.BoolValue(val)
    }
    if val, ok := podmanHostDataResponse["archived_at"].(string); ok {
        data.ArchivedAt = types.StringValue(val)
    }
    if val, ok := podmanHostDataResponse["archived_by_user_id"].(string); ok {
        data.ArchivedByUserId = types.StringValue(val)
    }
    if val, ok := podmanHostDataResponse["deleted_by_user_id"].(string); ok {
        data.DeletedByUserId = types.StringValue(val)
    }
    if val, ok := podmanHostDataResponse["labels"].([]interface{}); ok {
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
    if val, ok := podmanHostDataResponse["retain_telemetry_data_for_days"].(float64); ok {
        data.RetainTelemetryDataForDays = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := podmanHostDataResponse["telemetry_retention_config"].(string); ok {
        data.TelemetryRetentionConfig = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
