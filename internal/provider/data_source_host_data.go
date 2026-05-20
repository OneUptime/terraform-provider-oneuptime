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
var _ datasource.DataSource = &HostDataDataSource{}

func NewHostDataDataSource() datasource.DataSource {
    return &HostDataDataSource{}
}

// HostDataDataSource defines the data source implementation.
type HostDataDataSource struct {
    client *Client
}

// HostDataDataSourceModel describes the data source data model.
type HostDataDataSourceModel struct {
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
    LastSeenAt types.String `tfsdk:"last_seen_at"`
    OsType types.String `tfsdk:"os_type"`
    OsVersion types.String `tfsdk:"os_version"`
    HostId types.String `tfsdk:"host_id"`
    HostArch types.String `tfsdk:"host_arch"`
    HostType types.String `tfsdk:"host_type"`
    HostIpAddresses types.String `tfsdk:"host_ip_addresses"`
    CpuCores types.Number `tfsdk:"cpu_cores"`
    TotalMemoryBytes types.Number `tfsdk:"total_memory_bytes"`
    ProcessCount types.Number `tfsdk:"process_count"`
    ContainerRuntime types.String `tfsdk:"container_runtime"`
    DockerHostId types.String `tfsdk:"docker_host_id"`
    KubernetesClusterId types.String `tfsdk:"kubernetes_cluster_id"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
    Labels types.Set `tfsdk:"labels"`
    RetainTelemetryDataForDays types.Number `tfsdk:"retain_telemetry_data_for_days"`
    TelemetryRetentionConfig types.String `tfsdk:"telemetry_retention_config"`
}

func (d *HostDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_host_data"
}

func (d *HostDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "host_data data source",

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
                MarkdownDescription: "Friendly globally unique name for your object. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Host], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Friendly description for this host. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Host], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Host], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Host]",
                Computed: true,
            },
            "host_identifier": schema.StringAttribute{
                MarkdownDescription: "Unique identifier for this host, sourced from the host.name OTel resource attribute. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Host], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Host], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Host]",
                Computed: true,
            },
            "otel_collector_status": schema.StringAttribute{
                MarkdownDescription: "Connection status of the OTel Collector reporting on this host (connected or disconnected). Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Host], Update: [Project Owner, Project Admin, Edit Host]",
                Computed: true,
            },
            "last_seen_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "os_type": schema.StringAttribute{
                MarkdownDescription: "Operating system type of the host. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Host], Update: [Project Owner, Project Admin, Edit Host]",
                Computed: true,
            },
            "os_version": schema.StringAttribute{
                MarkdownDescription: "Operating system version of the host. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Host], Update: [Project Owner, Project Admin, Edit Host]",
                Computed: true,
            },
            "host_id": schema.StringAttribute{
                MarkdownDescription: "Stable host identifier reported by the OTel host.id resource attribute. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Host], Update: [Project Owner, Project Admin, Edit Host]",
                Computed: true,
            },
            "host_arch": schema.StringAttribute{
                MarkdownDescription: "CPU architecture from the OTel host.arch resource attribute. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Host], Update: [Project Owner, Project Admin, Edit Host]",
                Computed: true,
            },
            "host_type": schema.StringAttribute{
                MarkdownDescription: "Cloud-instance class reported by the OTel host.type resource attribute. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Host], Update: [Project Owner, Project Admin, Edit Host]",
                Computed: true,
            },
            "host_ip_addresses": schema.StringAttribute{
                MarkdownDescription: "Comma-separated list of IP addresses reported by the OTel host.ip resource attribute. The first non-loopback IPv4 is used for display.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Host], Update: [Project Owner, Project Admin, Edit Host]",
                Computed: true,
            },
            "cpu_cores": schema.NumberAttribute{
                MarkdownDescription: "Logical CPU core count, sourced from system.cpu.logical.count metric. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Host], Update: [Project Owner, Project Admin, Edit Host]",
                Computed: true,
            },
            "total_memory_bytes": schema.NumberAttribute{
                MarkdownDescription: "Total physical memory in bytes, sourced from system.memory.usage metric (sum of all states).. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Host], Update: [Project Owner, Project Admin, Edit Host]",
                Computed: true,
            },
            "process_count": schema.NumberAttribute{
                MarkdownDescription: "Most recent process count from system.processes.count metric. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Host], Update: [Project Owner, Project Admin, Edit Host]",
                Computed: true,
            },
            "container_runtime": schema.StringAttribute{
                MarkdownDescription: "Container runtime detected on this host, if any (e.g. docker, containerd). Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Host], Update: [Project Owner, Project Admin, Edit Host]",
                Computed: true,
            },
            "docker_host_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "kubernetes_cluster_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
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
            "labels": schema.SetAttribute{
                MarkdownDescription: "Relation to Labels Array where this object is categorized in.. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Host], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Host], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Host]",
                Computed: true,
                ElementType: types.StringType,
            },
            "retain_telemetry_data_for_days": schema.NumberAttribute{
                MarkdownDescription: "Number of days to retain telemetry data for this host. Leave blank to use the project-wide default.. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Host], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Host], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Host]",
                Computed: true,
            },
            "telemetry_retention_config": schema.StringAttribute{
                MarkdownDescription: "Per-pillar retention overrides for this host (logs by severity, traces by status, metrics, profiles). Unset fields fall back to the host default, then the project's retention settings.. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Host], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Host], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Host]",
                Computed: true,
            },
        },
    }
}

func (d *HostDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *HostDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data HostDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "host" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read host_data, got error: %s", err))
        return
    }

    var hostDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &hostDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse host_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := hostDataResponse["data"].(map[string]interface{}); ok {
        hostDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := hostDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := hostDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := hostDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := hostDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := hostDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := hostDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := hostDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := hostDataResponse["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    }
    if val, ok := hostDataResponse["description"].(string); ok {
        data.Description = types.StringValue(val)
    }
    if val, ok := hostDataResponse["host_identifier"].(string); ok {
        data.HostIdentifier = types.StringValue(val)
    }
    if val, ok := hostDataResponse["otel_collector_status"].(string); ok {
        data.OtelCollectorStatus = types.StringValue(val)
    }
    if val, ok := hostDataResponse["last_seen_at"].(string); ok {
        data.LastSeenAt = types.StringValue(val)
    }
    if val, ok := hostDataResponse["os_type"].(string); ok {
        data.OsType = types.StringValue(val)
    }
    if val, ok := hostDataResponse["os_version"].(string); ok {
        data.OsVersion = types.StringValue(val)
    }
    if val, ok := hostDataResponse["host_id"].(string); ok {
        data.HostId = types.StringValue(val)
    }
    if val, ok := hostDataResponse["host_arch"].(string); ok {
        data.HostArch = types.StringValue(val)
    }
    if val, ok := hostDataResponse["host_type"].(string); ok {
        data.HostType = types.StringValue(val)
    }
    if val, ok := hostDataResponse["host_ip_addresses"].(string); ok {
        data.HostIpAddresses = types.StringValue(val)
    }
    if val, ok := hostDataResponse["cpu_cores"].(float64); ok {
        data.CpuCores = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := hostDataResponse["total_memory_bytes"].(float64); ok {
        data.TotalMemoryBytes = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := hostDataResponse["process_count"].(float64); ok {
        data.ProcessCount = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := hostDataResponse["container_runtime"].(string); ok {
        data.ContainerRuntime = types.StringValue(val)
    }
    if val, ok := hostDataResponse["docker_host_id"].(string); ok {
        data.DockerHostId = types.StringValue(val)
    }
    if val, ok := hostDataResponse["kubernetes_cluster_id"].(string); ok {
        data.KubernetesClusterId = types.StringValue(val)
    }
    if val, ok := hostDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := hostDataResponse["deleted_by_user_id"].(string); ok {
        data.DeletedByUserId = types.StringValue(val)
    }
    if val, ok := hostDataResponse["labels"].([]interface{}); ok {
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
    if val, ok := hostDataResponse["retain_telemetry_data_for_days"].(float64); ok {
        data.RetainTelemetryDataForDays = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := hostDataResponse["telemetry_retention_config"].(string); ok {
        data.TelemetryRetentionConfig = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
