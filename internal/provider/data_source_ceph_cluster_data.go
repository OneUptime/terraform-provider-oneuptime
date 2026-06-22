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
var _ datasource.DataSource = &CephClusterDataDataSource{}

func NewCephClusterDataDataSource() datasource.DataSource {
    return &CephClusterDataDataSource{}
}

// CephClusterDataDataSource defines the data source implementation.
type CephClusterDataDataSource struct {
    client *Client
}

// CephClusterDataDataSourceModel describes the data source data model.
type CephClusterDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    Slug types.String `tfsdk:"slug"`
    Description types.String `tfsdk:"description"`
    Fsid types.String `tfsdk:"fsid"`
    OtelCollectorStatus types.String `tfsdk:"otel_collector_status"`
    AgentVersion types.String `tfsdk:"agent_version"`
    CephVersion types.String `tfsdk:"ceph_version"`
    LastSeenAt types.String `tfsdk:"last_seen_at"`
    MonCount types.Number `tfsdk:"mon_count"`
    OsdCount types.Number `tfsdk:"osd_count"`
    OsdUpCount types.Number `tfsdk:"osd_up_count"`
    OsdInCount types.Number `tfsdk:"osd_in_count"`
    PoolCount types.Number `tfsdk:"pool_count"`
    HealthStatus types.Number `tfsdk:"health_status"`
    CapacityUsedPercent types.Number `tfsdk:"capacity_used_percent"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    IsArchived types.Bool `tfsdk:"is_archived"`
    ArchivedAt types.String `tfsdk:"archived_at"`
    ArchivedByUserId types.String `tfsdk:"archived_by_user_id"`
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
    Labels types.Set `tfsdk:"labels"`
    RetainTelemetryDataForDays types.Number `tfsdk:"retain_telemetry_data_for_days"`
    TelemetryRetentionConfig types.String `tfsdk:"telemetry_retention_config"`
}

func (d *CephClusterDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_ceph_cluster_data"
}

func (d *CephClusterDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "ceph_cluster_data data source",

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
                MarkdownDescription: "Friendly globally unique name for your object. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Ceph Cluster], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Friendly description for this Ceph cluster. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Ceph Cluster], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Ceph Cluster], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Ceph Cluster]",
                Computed: true,
            },
            "fsid": schema.StringAttribute{
                MarkdownDescription: "Ceph cluster fsid, sourced from the ceph.cluster.fsid OTel resource attribute when known. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Ceph Cluster], Update: [Project Owner, Project Admin, Edit Ceph Cluster]",
                Computed: true,
            },
            "otel_collector_status": schema.StringAttribute{
                MarkdownDescription: "Connection status of the OTel Collector agent (connected or disconnected). Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Ceph Cluster], Update: [Project Owner, Project Admin, Edit Ceph Cluster]",
                Computed: true,
            },
            "agent_version": schema.StringAttribute{
                MarkdownDescription: "Version of the OneUptime Ceph agent reporting telemetry, as self-reported via the oneuptime.agent.version resource attribute. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Ceph Cluster], Update: [Project Owner, Project Admin, Edit Ceph Cluster]",
                Computed: true,
            },
            "ceph_version": schema.StringAttribute{
                MarkdownDescription: "Ceph version reported by this cluster. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Ceph Cluster], Update: [Project Owner, Project Admin, Edit Ceph Cluster]",
                Computed: true,
            },
            "last_seen_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "mon_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of Ceph monitors (mons) in this cluster. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Ceph Cluster], Update: [Project Owner, Project Admin, Edit Ceph Cluster]",
                Computed: true,
            },
            "osd_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of OSDs in this cluster. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Ceph Cluster], Update: [Project Owner, Project Admin, Edit Ceph Cluster]",
                Computed: true,
            },
            "osd_up_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of OSDs that are up (ceph_osd_up == 1) in this cluster. Rendered as 'X up / Y in / Z total' next to osdCount.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Ceph Cluster], Update: [Project Owner, Project Admin, Edit Ceph Cluster]",
                Computed: true,
            },
            "osd_in_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of OSDs that are in the cluster (ceph_osd_in == 1). Rendered as 'X up / Y in / Z total' next to osdCount.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Ceph Cluster], Update: [Project Owner, Project Admin, Edit Ceph Cluster]",
                Computed: true,
            },
            "pool_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of pools in this cluster. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Ceph Cluster], Update: [Project Owner, Project Admin, Edit Ceph Cluster]",
                Computed: true,
            },
            "health_status": schema.NumberAttribute{
                MarkdownDescription: "Cached latest ceph_health_status value: 0 = HEALTH_OK, 1 = HEALTH_WARN, 2 = HEALTH_ERR. Rendered as the OK/WARN/ERR health pill. Null until the first metric batch arrives.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Ceph Cluster], Update: [Project Owner, Project Admin, Edit Ceph Cluster]",
                Computed: true,
            },
            "capacity_used_percent": schema.NumberAttribute{
                MarkdownDescription: "Cached cluster capacity usage percent (ceph_cluster_total_used_bytes / ceph_cluster_total_bytes * 100). Stored as decimal so sub-percent precision survives the round trip. Null until both series appear in one metric batch.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Ceph Cluster], Update: [Project Owner, Project Admin, Edit Ceph Cluster]",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "is_archived": schema.BoolAttribute{
                MarkdownDescription: "Is this Ceph cluster archived? Archived Ceph clusters are hidden from lists but keep collecting telemetry.. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Ceph Cluster], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Ceph Cluster], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Ceph Cluster]",
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
                MarkdownDescription: "Relation to Labels Array where this object is categorized in.. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Ceph Cluster], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Ceph Cluster], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Ceph Cluster]",
                Computed: true,
                ElementType: types.StringType,
            },
            "retain_telemetry_data_for_days": schema.NumberAttribute{
                MarkdownDescription: "Number of days to retain telemetry data for this Ceph cluster. Leave blank to use the project-wide default.. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Ceph Cluster], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Ceph Cluster], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Ceph Cluster]",
                Computed: true,
            },
            "telemetry_retention_config": schema.StringAttribute{
                MarkdownDescription: "Per-pillar retention overrides for this Ceph cluster (logs by severity, traces by status, metrics, profiles). Unset fields fall back to the Ceph cluster default, then the project's retention settings.. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Ceph Cluster], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Ceph Cluster], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Ceph Cluster]",
                Computed: true,
            },
        },
    }
}

func (d *CephClusterDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *CephClusterDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data CephClusterDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "ceph-cluster" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read ceph_cluster_data, got error: %s", err))
        return
    }

    var cephClusterDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &cephClusterDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse ceph_cluster_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := cephClusterDataResponse["data"].(map[string]interface{}); ok {
        cephClusterDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := cephClusterDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := cephClusterDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := cephClusterDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := cephClusterDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := cephClusterDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := cephClusterDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := cephClusterDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := cephClusterDataResponse["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    }
    if val, ok := cephClusterDataResponse["description"].(string); ok {
        data.Description = types.StringValue(val)
    }
    if val, ok := cephClusterDataResponse["fsid"].(string); ok {
        data.Fsid = types.StringValue(val)
    }
    if val, ok := cephClusterDataResponse["otel_collector_status"].(string); ok {
        data.OtelCollectorStatus = types.StringValue(val)
    }
    if val, ok := cephClusterDataResponse["agent_version"].(string); ok {
        data.AgentVersion = types.StringValue(val)
    }
    if val, ok := cephClusterDataResponse["ceph_version"].(string); ok {
        data.CephVersion = types.StringValue(val)
    }
    if val, ok := cephClusterDataResponse["last_seen_at"].(string); ok {
        data.LastSeenAt = types.StringValue(val)
    }
    if val, ok := cephClusterDataResponse["mon_count"].(float64); ok {
        data.MonCount = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := cephClusterDataResponse["osd_count"].(float64); ok {
        data.OsdCount = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := cephClusterDataResponse["osd_up_count"].(float64); ok {
        data.OsdUpCount = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := cephClusterDataResponse["osd_in_count"].(float64); ok {
        data.OsdInCount = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := cephClusterDataResponse["pool_count"].(float64); ok {
        data.PoolCount = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := cephClusterDataResponse["health_status"].(float64); ok {
        data.HealthStatus = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := cephClusterDataResponse["capacity_used_percent"].(float64); ok {
        data.CapacityUsedPercent = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := cephClusterDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := cephClusterDataResponse["is_archived"].(bool); ok {
        data.IsArchived = types.BoolValue(val)
    }
    if val, ok := cephClusterDataResponse["archived_at"].(string); ok {
        data.ArchivedAt = types.StringValue(val)
    }
    if val, ok := cephClusterDataResponse["archived_by_user_id"].(string); ok {
        data.ArchivedByUserId = types.StringValue(val)
    }
    if val, ok := cephClusterDataResponse["deleted_by_user_id"].(string); ok {
        data.DeletedByUserId = types.StringValue(val)
    }
    if val, ok := cephClusterDataResponse["labels"].([]interface{}); ok {
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
    if val, ok := cephClusterDataResponse["retain_telemetry_data_for_days"].(float64); ok {
        data.RetainTelemetryDataForDays = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := cephClusterDataResponse["telemetry_retention_config"].(string); ok {
        data.TelemetryRetentionConfig = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
