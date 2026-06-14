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
var _ datasource.DataSource = &ProxmoxClusterDataDataSource{}

func NewProxmoxClusterDataDataSource() datasource.DataSource {
    return &ProxmoxClusterDataDataSource{}
}

// ProxmoxClusterDataDataSource defines the data source implementation.
type ProxmoxClusterDataDataSource struct {
    client *Client
}

// ProxmoxClusterDataDataSourceModel describes the data source data model.
type ProxmoxClusterDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    Slug types.String `tfsdk:"slug"`
    Description types.String `tfsdk:"description"`
    OtelCollectorStatus types.String `tfsdk:"otel_collector_status"`
    AgentVersion types.String `tfsdk:"agent_version"`
    PveVersion types.String `tfsdk:"pve_version"`
    LastSeenAt types.String `tfsdk:"last_seen_at"`
    NodeCount types.Number `tfsdk:"node_count"`
    OnlineNodeCount types.Number `tfsdk:"online_node_count"`
    GuestCount types.Number `tfsdk:"guest_count"`
    StorageCount types.Number `tfsdk:"storage_count"`
    GuestsWithoutBackupCount types.Number `tfsdk:"guests_without_backup_count"`
    CephClusterId types.String `tfsdk:"ceph_cluster_id"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
    Labels types.Set `tfsdk:"labels"`
    RetainTelemetryDataForDays types.Number `tfsdk:"retain_telemetry_data_for_days"`
    TelemetryRetentionConfig types.String `tfsdk:"telemetry_retention_config"`
}

func (d *ProxmoxClusterDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_proxmox_cluster_data"
}

func (d *ProxmoxClusterDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "proxmox_cluster_data data source",

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
                MarkdownDescription: "Friendly globally unique name for your object. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Proxmox Cluster], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Friendly description for this Proxmox cluster. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Proxmox Cluster], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Proxmox Cluster], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Proxmox Cluster]",
                Computed: true,
            },
            "otel_collector_status": schema.StringAttribute{
                MarkdownDescription: "Connection status of the OTel Collector agent (connected or disconnected). Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Proxmox Cluster], Update: [Project Owner, Project Admin, Edit Proxmox Cluster]",
                Computed: true,
            },
            "agent_version": schema.StringAttribute{
                MarkdownDescription: "Version of the OneUptime Proxmox agent reporting telemetry, as self-reported via the oneuptime.agent.version resource attribute. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Proxmox Cluster], Update: [Project Owner, Project Admin, Edit Proxmox Cluster]",
                Computed: true,
            },
            "pve_version": schema.StringAttribute{
                MarkdownDescription: "Proxmox VE version reported by this cluster. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Proxmox Cluster], Update: [Project Owner, Project Admin, Edit Proxmox Cluster]",
                Computed: true,
            },
            "last_seen_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "node_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of nodes in this cluster. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Proxmox Cluster], Update: [Project Owner, Project Admin, Edit Proxmox Cluster]",
                Computed: true,
            },
            "online_node_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of nodes currently online (pve_up == 1) in this cluster. Rendered as 'Nodes X/Y online' next to nodeCount.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Proxmox Cluster], Update: [Project Owner, Project Admin, Edit Proxmox Cluster]",
                Computed: true,
            },
            "guest_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of guests (VMs and containers) in this cluster. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Proxmox Cluster], Update: [Project Owner, Project Admin, Edit Proxmox Cluster]",
                Computed: true,
            },
            "storage_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of storage pools in this cluster. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Proxmox Cluster], Update: [Project Owner, Project Admin, Edit Proxmox Cluster]",
                Computed: true,
            },
            "guests_without_backup_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of guests not covered by ANY backup job (pve_not_backed_up_total). NULL until the exporter's cluster-level backup-info collector reports. Coverage by a job is NOT the same as recent/successful backups — freshness needs the PVE task log or PBS API.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Proxmox Cluster], Update: [Project Owner, Project Admin, Edit Proxmox Cluster]",
                Computed: true,
            },
            "ceph_cluster_id": schema.StringAttribute{
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
                MarkdownDescription: "Relation to Labels Array where this object is categorized in.. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Proxmox Cluster], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Proxmox Cluster], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Proxmox Cluster]",
                Computed: true,
                ElementType: types.StringType,
            },
            "retain_telemetry_data_for_days": schema.NumberAttribute{
                MarkdownDescription: "Number of days to retain telemetry data for this Proxmox cluster. Leave blank to use the project-wide default.. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Proxmox Cluster], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Proxmox Cluster], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Proxmox Cluster]",
                Computed: true,
            },
            "telemetry_retention_config": schema.StringAttribute{
                MarkdownDescription: "Per-pillar retention overrides for this Proxmox cluster (logs by severity, traces by status, metrics, profiles). Unset fields fall back to the Proxmox cluster default, then the project's retention settings.. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Proxmox Cluster], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Proxmox Cluster], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Proxmox Cluster]",
                Computed: true,
            },
        },
    }
}

func (d *ProxmoxClusterDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ProxmoxClusterDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data ProxmoxClusterDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "proxmox-cluster" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read proxmox_cluster_data, got error: %s", err))
        return
    }

    var proxmoxClusterDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &proxmoxClusterDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse proxmox_cluster_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := proxmoxClusterDataResponse["data"].(map[string]interface{}); ok {
        proxmoxClusterDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := proxmoxClusterDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := proxmoxClusterDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := proxmoxClusterDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := proxmoxClusterDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := proxmoxClusterDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := proxmoxClusterDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := proxmoxClusterDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := proxmoxClusterDataResponse["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    }
    if val, ok := proxmoxClusterDataResponse["description"].(string); ok {
        data.Description = types.StringValue(val)
    }
    if val, ok := proxmoxClusterDataResponse["otel_collector_status"].(string); ok {
        data.OtelCollectorStatus = types.StringValue(val)
    }
    if val, ok := proxmoxClusterDataResponse["agent_version"].(string); ok {
        data.AgentVersion = types.StringValue(val)
    }
    if val, ok := proxmoxClusterDataResponse["pve_version"].(string); ok {
        data.PveVersion = types.StringValue(val)
    }
    if val, ok := proxmoxClusterDataResponse["last_seen_at"].(string); ok {
        data.LastSeenAt = types.StringValue(val)
    }
    if val, ok := proxmoxClusterDataResponse["node_count"].(float64); ok {
        data.NodeCount = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := proxmoxClusterDataResponse["online_node_count"].(float64); ok {
        data.OnlineNodeCount = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := proxmoxClusterDataResponse["guest_count"].(float64); ok {
        data.GuestCount = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := proxmoxClusterDataResponse["storage_count"].(float64); ok {
        data.StorageCount = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := proxmoxClusterDataResponse["guests_without_backup_count"].(float64); ok {
        data.GuestsWithoutBackupCount = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := proxmoxClusterDataResponse["ceph_cluster_id"].(string); ok {
        data.CephClusterId = types.StringValue(val)
    }
    if val, ok := proxmoxClusterDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := proxmoxClusterDataResponse["deleted_by_user_id"].(string); ok {
        data.DeletedByUserId = types.StringValue(val)
    }
    if val, ok := proxmoxClusterDataResponse["labels"].([]interface{}); ok {
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
    if val, ok := proxmoxClusterDataResponse["retain_telemetry_data_for_days"].(float64); ok {
        data.RetainTelemetryDataForDays = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := proxmoxClusterDataResponse["telemetry_retention_config"].(string); ok {
        data.TelemetryRetentionConfig = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
