package shard25

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "github.com/hashicorp/terraform-plugin-framework/attr"
    "sort"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &StorageArrayDataSource{}

func NewStorageArrayDataSource() datasource.DataSource {
    return &StorageArrayDataSource{}
}

// StorageArrayDataSource defines the data source implementation.
type StorageArrayDataSource struct {
    client *Client
}

// StorageArrayDataSourceModel describes the data source data model.
type StorageArrayDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    CreatedByActorType types.String `tfsdk:"created_by_actor_type"`
    CreatedByActorId types.String `tfsdk:"created_by_actor_id"`
    LastModifiedAt types.String `tfsdk:"last_modified_at"`
    LastModifiedByActorType types.String `tfsdk:"last_modified_by_actor_type"`
    LastModifiedByActorId types.String `tfsdk:"last_modified_by_actor_id"`
    LastModifiedByUserId types.String `tfsdk:"last_modified_by_user_id"`
    ProjectId types.String `tfsdk:"project_id"`
    Name types.String `tfsdk:"name"`
    Slug types.String `tfsdk:"slug"`
    Description types.String `tfsdk:"description"`
    StorageSystem types.String `tfsdk:"storage_system"`
    ReportedName types.String `tfsdk:"reported_name"`
    SystemId types.String `tfsdk:"system_id"`
    OsName types.String `tfsdk:"os_name"`
    OsVersion types.String `tfsdk:"os_version"`
    OtelCollectorStatus types.String `tfsdk:"otel_collector_status"`
    AgentVersion types.String `tfsdk:"agent_version"`
    LastSeenAt types.String `tfsdk:"last_seen_at"`
    CapacityBytes types.Number `tfsdk:"capacity_bytes"`
    UsedBytes types.Number `tfsdk:"used_bytes"`
    CapacityUsedPercent types.Number `tfsdk:"capacity_used_percent"`
    DataReductionRatio types.Number `tfsdk:"data_reduction_ratio"`
    OpenAlertCount types.Number `tfsdk:"open_alert_count"`
    CriticalAlertCount types.Number `tfsdk:"critical_alert_count"`
    WarningAlertCount types.Number `tfsdk:"warning_alert_count"`
    VolumeCount types.Number `tfsdk:"volume_count"`
    HostCount types.Number `tfsdk:"host_count"`
    PodCount types.Number `tfsdk:"pod_count"`
    FileSystemCount types.Number `tfsdk:"file_system_count"`
    BucketCount types.Number `tfsdk:"bucket_count"`
    HardwareComponentCount types.Number `tfsdk:"hardware_component_count"`
    UnhealthyHardwareCount types.Number `tfsdk:"unhealthy_hardware_count"`
    HealthStatus types.Number `tfsdk:"health_status"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    IsArchived types.Bool `tfsdk:"is_archived"`
    ArchivedAt types.String `tfsdk:"archived_at"`
    ArchivedByUserId types.String `tfsdk:"archived_by_user_id"`
    Labels types.Set `tfsdk:"labels"`
    RetainTelemetryDataForDays types.Number `tfsdk:"retain_telemetry_data_for_days"`
    TelemetryRetentionConfig types.String `tfsdk:"telemetry_retention_config"`
}

func (d *StorageArrayDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_storage_array"
}

func (d *StorageArrayDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Storage arrays (Pure Storage FlashArray and FlashBlade) that are being monitored in this project. Each array is auto-discovered when the OneUptime Storage Array Agent sends metrics, or can be registered by hand. Look up an existing storage array by `id`, or by any of its other arguments (`name`, `agent_version`, `archived_by_user_id`, ...): each one set must match, and exactly one storage array may match them all.",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Look up by unique identifier. Leave unset to look up by the other arguments instead.",
                Optional: true,
                Computed: true,
            },
            "created_at": schema.StringAttribute{
                MarkdownDescription: "Date and Time when the object was created.",
                Computed: true,
            },
            "updated_at": schema.StringAttribute{
                MarkdownDescription: "Date and Time when the object was updated.",
                Computed: true,
            },
            "created_by_actor_type": schema.StringAttribute{
                MarkdownDescription: "What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.",
                Optional: true,
                Computed: true,
            },
            "created_by_actor_id": schema.StringAttribute{
                MarkdownDescription: "ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.",
                Optional: true,
                Computed: true,
            },
            "last_modified_at": schema.StringAttribute{
                MarkdownDescription: "When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.",
                Computed: true,
            },
            "last_modified_by_actor_type": schema.StringAttribute{
                MarkdownDescription: "What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.",
                Optional: true,
                Computed: true,
            },
            "last_modified_by_actor_id": schema.StringAttribute{
                MarkdownDescription: "ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.",
                Optional: true,
                Computed: true,
            },
            "last_modified_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: DataSourceRecordIDValidators("oneuptime_user"),
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Name of this storage array in OneUptime. This is the join key — it must match the storage.array.name OTel resource attribute stamped by the OneUptime Storage Array Agent.",
                Optional: true,
                Computed: true,
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Friendly globally unique name for your object.",
                Optional: true,
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Friendly description for this storage array.",
                Optional: true,
                Computed: true,
            },
            "storage_system": schema.StringAttribute{
                MarkdownDescription: "The storage platform this array runs, normalized from the storage.system OTel resource attribute (or detected from the metric names): purestorage.flasharray or purestorage.flashblade.",
                Optional: true,
                Computed: true,
            },
            "reported_name": schema.StringAttribute{
                MarkdownDescription: "The name the array reports for itself (array_name on purefa_info / purefb_info). Can differ from the OneUptime name, which comes from the agent configuration.",
                Optional: true,
                Computed: true,
            },
            "system_id": schema.StringAttribute{
                MarkdownDescription: "The array's own system identifier (system_id on purefa_info / purefb_info).",
                Optional: true,
                Computed: true,
            },
            "os_name": schema.StringAttribute{
                MarkdownDescription: "Name of the array's operating system as it reports it (os on purefa_info / purefb_info).",
                Optional: true,
                Computed: true,
            },
            "os_version": schema.StringAttribute{
                MarkdownDescription: "Version of the array's operating system as it reports it (version on purefa_info / purefb_info).",
                Optional: true,
                Computed: true,
            },
            "otel_collector_status": schema.StringAttribute{
                MarkdownDescription: "Connection status of the OTel Collector agent (connected or disconnected).",
                Optional: true,
                Computed: true,
            },
            "agent_version": schema.StringAttribute{
                MarkdownDescription: "Version of the OneUptime Storage Array Agent reporting telemetry, as self-reported via the oneuptime.agent.version resource attribute.",
                Optional: true,
                Computed: true,
            },
            "last_seen_at": schema.StringAttribute{
                MarkdownDescription: "When metrics were last received from this storage array.",
                Computed: true,
            },
            "capacity_bytes": schema.NumberAttribute{
                MarkdownDescription: "Cached usable capacity of the array in bytes (purefa_array_space_bytes{space=\"capacity\"} / purefb_array_space_bytes{type=\"array\",space=\"capacity\"}). Null until the first array metric batch arrives.",
                Optional: true,
                Computed: true,
            },
            "used_bytes": schema.NumberAttribute{
                MarkdownDescription: "Cached physical space used on the array in bytes: capacity minus empty space, or capacity times utilization when the empty series is missing. Null until the first array metric batch arrives.",
                Optional: true,
                Computed: true,
            },
            "capacity_used_percent": schema.NumberAttribute{
                MarkdownDescription: "Cached array space utilization in percent (purefa_array_space_utilization / purefb_array_space_utilization). Stored as decimal so sub-percent precision survives the round trip. Null until the first array metric batch arrives.",
                Optional: true,
                Computed: true,
            },
            "data_reduction_ratio": schema.NumberAttribute{
                MarkdownDescription: "Cached data reduction ratio of the array (purefa_array_space_data_reduction_ratio / purefb_array_space_data_reduction_ratio), for example 4.2 for 4.2:1. Null until the first array metric batch arrives.",
                Optional: true,
                Computed: true,
            },
            "open_alert_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of alerts open on the array itself (purefa_alerts_open / purefb_alerts_open series, hidden alerts excluded).",
                Optional: true,
                Computed: true,
            },
            "critical_alert_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of open array alerts with critical severity.",
                Optional: true,
                Computed: true,
            },
            "warning_alert_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of open array alerts with warning severity.",
                Optional: true,
                Computed: true,
            },
            "volume_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of volumes on the array (FlashArray).",
                Optional: true,
                Computed: true,
            },
            "host_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of hosts defined on the array (FlashArray).",
                Optional: true,
                Computed: true,
            },
            "pod_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of pods (ActiveCluster / ActiveDR replication containers) on the array (FlashArray).",
                Optional: true,
                Computed: true,
            },
            "file_system_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of file systems on the array (FlashBlade).",
                Optional: true,
                Computed: true,
            },
            "bucket_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of object store buckets on the array (FlashBlade).",
                Optional: true,
                Computed: true,
            },
            "hardware_component_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of hardware components the array reports (chassis, controllers, drive bays, power supplies, fans, ports, blades...).",
                Optional: true,
                Computed: true,
            },
            "unhealthy_hardware_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of hardware components, drives and controllers in a critical, degraded, failed or unknown state.",
                Optional: true,
                Computed: true,
            },
            "health_status": schema.NumberAttribute{
                MarkdownDescription: "Cached array health derived from the array's open alerts and hardware state: 0 = OK, 1 = Warning (a warning alert, or a degraded or unknown component), 2 = Critical (a critical alert, or a failed or critical component). Rendered as the health pill. Null until the first metric batch arrives.",
                Optional: true,
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: DataSourceRecordIDValidators("oneuptime_user"),
            },
            "is_archived": schema.BoolAttribute{
                MarkdownDescription: "Is this storage array archived? Archived storage arrays are hidden from lists but keep collecting telemetry.",
                Optional: true,
                Computed: true,
            },
            "archived_at": schema.StringAttribute{
                MarkdownDescription: "When was this storage array archived?",
                Computed: true,
            },
            "archived_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who archived this object (if this object was archived by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: DataSourceRecordIDValidators("oneuptime_user"),
            },
            "labels": schema.SetAttribute{
                MarkdownDescription: "Labels to group and filter by. IDs of `oneuptime_label` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "retain_telemetry_data_for_days": schema.NumberAttribute{
                MarkdownDescription: "Number of days to retain telemetry data for this storage array. Leave blank to use the project-wide default.",
                Optional: true,
                Computed: true,
            },
            "telemetry_retention_config": schema.StringAttribute{
                MarkdownDescription: "Per-pillar retention overrides for this storage array (logs by severity, traces by status, metrics, profiles). Unset fields fall back to the storage array default, then the project's retention settings. A JSON value: read it with `jsondecode()`.",
                Computed: true,
            },
        },
    }
}

func (d *StorageArrayDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *StorageArrayDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data StorageArrayDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.CreatedByActorType.IsNull() && !data.CreatedByActorType.IsUnknown() {
        filters["createdByActorType"] = data.CreatedByActorType.ValueString()
        filterNames = append(filterNames, "created_by_actor_type = "+fmt.Sprintf("%q", data.CreatedByActorType.ValueString()))
    }
    if !data.CreatedByActorId.IsNull() && !data.CreatedByActorId.IsUnknown() {
        filters["createdByActorId"] = data.CreatedByActorId.ValueString()
        filterNames = append(filterNames, "created_by_actor_id = "+fmt.Sprintf("%q", data.CreatedByActorId.ValueString()))
    }
    if !data.LastModifiedByActorType.IsNull() && !data.LastModifiedByActorType.IsUnknown() {
        filters["lastModifiedByActorType"] = data.LastModifiedByActorType.ValueString()
        filterNames = append(filterNames, "last_modified_by_actor_type = "+fmt.Sprintf("%q", data.LastModifiedByActorType.ValueString()))
    }
    if !data.LastModifiedByActorId.IsNull() && !data.LastModifiedByActorId.IsUnknown() {
        filters["lastModifiedByActorId"] = data.LastModifiedByActorId.ValueString()
        filterNames = append(filterNames, "last_modified_by_actor_id = "+fmt.Sprintf("%q", data.LastModifiedByActorId.ValueString()))
    }
    if !data.LastModifiedByUserId.IsNull() && !data.LastModifiedByUserId.IsUnknown() {
        filters["lastModifiedByUserId"] = data.LastModifiedByUserId.ValueString()
        filterNames = append(filterNames, "last_modified_by_user_id = "+fmt.Sprintf("%q", data.LastModifiedByUserId.ValueString()))
    }
    if !data.Name.IsNull() && !data.Name.IsUnknown() {
        filters["name"] = data.Name.ValueString()
        filterNames = append(filterNames, "name = "+fmt.Sprintf("%q", data.Name.ValueString()))
    }
    if !data.Slug.IsNull() && !data.Slug.IsUnknown() {
        filters["slug"] = data.Slug.ValueString()
        filterNames = append(filterNames, "slug = "+fmt.Sprintf("%q", data.Slug.ValueString()))
    }
    if !data.Description.IsNull() && !data.Description.IsUnknown() {
        filters["description"] = data.Description.ValueString()
        filterNames = append(filterNames, "description = "+fmt.Sprintf("%q", data.Description.ValueString()))
    }
    if !data.StorageSystem.IsNull() && !data.StorageSystem.IsUnknown() {
        filters["storageSystem"] = data.StorageSystem.ValueString()
        filterNames = append(filterNames, "storage_system = "+fmt.Sprintf("%q", data.StorageSystem.ValueString()))
    }
    if !data.ReportedName.IsNull() && !data.ReportedName.IsUnknown() {
        filters["reportedName"] = data.ReportedName.ValueString()
        filterNames = append(filterNames, "reported_name = "+fmt.Sprintf("%q", data.ReportedName.ValueString()))
    }
    if !data.SystemId.IsNull() && !data.SystemId.IsUnknown() {
        filters["systemId"] = data.SystemId.ValueString()
        filterNames = append(filterNames, "system_id = "+fmt.Sprintf("%q", data.SystemId.ValueString()))
    }
    if !data.OsName.IsNull() && !data.OsName.IsUnknown() {
        filters["osName"] = data.OsName.ValueString()
        filterNames = append(filterNames, "os_name = "+fmt.Sprintf("%q", data.OsName.ValueString()))
    }
    if !data.OsVersion.IsNull() && !data.OsVersion.IsUnknown() {
        filters["osVersion"] = data.OsVersion.ValueString()
        filterNames = append(filterNames, "os_version = "+fmt.Sprintf("%q", data.OsVersion.ValueString()))
    }
    if !data.OtelCollectorStatus.IsNull() && !data.OtelCollectorStatus.IsUnknown() {
        filters["otelCollectorStatus"] = data.OtelCollectorStatus.ValueString()
        filterNames = append(filterNames, "otel_collector_status = "+fmt.Sprintf("%q", data.OtelCollectorStatus.ValueString()))
    }
    if !data.AgentVersion.IsNull() && !data.AgentVersion.IsUnknown() {
        filters["agentVersion"] = data.AgentVersion.ValueString()
        filterNames = append(filterNames, "agent_version = "+fmt.Sprintf("%q", data.AgentVersion.ValueString()))
    }
    if !data.CapacityBytes.IsNull() && !data.CapacityBytes.IsUnknown() {
        filters["capacityBytes"] = LookupNumber(data.CapacityBytes)
        filterNames = append(filterNames, "capacity_bytes = "+data.CapacityBytes.ValueBigFloat().String())
    }
    if !data.UsedBytes.IsNull() && !data.UsedBytes.IsUnknown() {
        filters["usedBytes"] = LookupNumber(data.UsedBytes)
        filterNames = append(filterNames, "used_bytes = "+data.UsedBytes.ValueBigFloat().String())
    }
    if !data.CapacityUsedPercent.IsNull() && !data.CapacityUsedPercent.IsUnknown() {
        filters["capacityUsedPercent"] = LookupNumber(data.CapacityUsedPercent)
        filterNames = append(filterNames, "capacity_used_percent = "+data.CapacityUsedPercent.ValueBigFloat().String())
    }
    if !data.DataReductionRatio.IsNull() && !data.DataReductionRatio.IsUnknown() {
        filters["dataReductionRatio"] = LookupNumber(data.DataReductionRatio)
        filterNames = append(filterNames, "data_reduction_ratio = "+data.DataReductionRatio.ValueBigFloat().String())
    }
    if !data.OpenAlertCount.IsNull() && !data.OpenAlertCount.IsUnknown() {
        filters["openAlertCount"] = LookupNumber(data.OpenAlertCount)
        filterNames = append(filterNames, "open_alert_count = "+data.OpenAlertCount.ValueBigFloat().String())
    }
    if !data.CriticalAlertCount.IsNull() && !data.CriticalAlertCount.IsUnknown() {
        filters["criticalAlertCount"] = LookupNumber(data.CriticalAlertCount)
        filterNames = append(filterNames, "critical_alert_count = "+data.CriticalAlertCount.ValueBigFloat().String())
    }
    if !data.WarningAlertCount.IsNull() && !data.WarningAlertCount.IsUnknown() {
        filters["warningAlertCount"] = LookupNumber(data.WarningAlertCount)
        filterNames = append(filterNames, "warning_alert_count = "+data.WarningAlertCount.ValueBigFloat().String())
    }
    if !data.VolumeCount.IsNull() && !data.VolumeCount.IsUnknown() {
        filters["volumeCount"] = LookupNumber(data.VolumeCount)
        filterNames = append(filterNames, "volume_count = "+data.VolumeCount.ValueBigFloat().String())
    }
    if !data.HostCount.IsNull() && !data.HostCount.IsUnknown() {
        filters["hostCount"] = LookupNumber(data.HostCount)
        filterNames = append(filterNames, "host_count = "+data.HostCount.ValueBigFloat().String())
    }
    if !data.PodCount.IsNull() && !data.PodCount.IsUnknown() {
        filters["podCount"] = LookupNumber(data.PodCount)
        filterNames = append(filterNames, "pod_count = "+data.PodCount.ValueBigFloat().String())
    }
    if !data.FileSystemCount.IsNull() && !data.FileSystemCount.IsUnknown() {
        filters["fileSystemCount"] = LookupNumber(data.FileSystemCount)
        filterNames = append(filterNames, "file_system_count = "+data.FileSystemCount.ValueBigFloat().String())
    }
    if !data.BucketCount.IsNull() && !data.BucketCount.IsUnknown() {
        filters["bucketCount"] = LookupNumber(data.BucketCount)
        filterNames = append(filterNames, "bucket_count = "+data.BucketCount.ValueBigFloat().String())
    }
    if !data.HardwareComponentCount.IsNull() && !data.HardwareComponentCount.IsUnknown() {
        filters["hardwareComponentCount"] = LookupNumber(data.HardwareComponentCount)
        filterNames = append(filterNames, "hardware_component_count = "+data.HardwareComponentCount.ValueBigFloat().String())
    }
    if !data.UnhealthyHardwareCount.IsNull() && !data.UnhealthyHardwareCount.IsUnknown() {
        filters["unhealthyHardwareCount"] = LookupNumber(data.UnhealthyHardwareCount)
        filterNames = append(filterNames, "unhealthy_hardware_count = "+data.UnhealthyHardwareCount.ValueBigFloat().String())
    }
    if !data.HealthStatus.IsNull() && !data.HealthStatus.IsUnknown() {
        filters["healthStatus"] = LookupNumber(data.HealthStatus)
        filterNames = append(filterNames, "health_status = "+data.HealthStatus.ValueBigFloat().String())
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }
    if !data.IsArchived.IsNull() && !data.IsArchived.IsUnknown() {
        filters["isArchived"] = data.IsArchived.ValueBool()
        filterNames = append(filterNames, "is_archived = "+fmt.Sprintf("%t", data.IsArchived.ValueBool()))
    }
    if !data.ArchivedByUserId.IsNull() && !data.ArchivedByUserId.IsUnknown() {
        filters["archivedByUserId"] = data.ArchivedByUserId.ValueString()
        filterNames = append(filterNames, "archived_by_user_id = "+fmt.Sprintf("%q", data.ArchivedByUserId.ValueString()))
    }
    if !data.RetainTelemetryDataForDays.IsNull() && !data.RetainTelemetryDataForDays.IsUnknown() {
        filters["retainTelemetryDataForDays"] = LookupNumber(data.RetainTelemetryDataForDays)
        filterNames = append(filterNames, "retain_telemetry_data_for_days = "+data.RetainTelemetryDataForDays.ValueBigFloat().String())
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the storage array up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the storage array up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "createdByActorType": true,
        "createdByActorId": true,
        "lastModifiedAt": true,
        "lastModifiedByActorType": true,
        "lastModifiedByActorId": true,
        "lastModifiedByUserId": true,
        "projectId": true,
        "name": true,
        "slug": true,
        "description": true,
        "storageSystem": true,
        "reportedName": true,
        "systemId": true,
        "osName": true,
        "osVersion": true,
        "otelCollectorStatus": true,
        "agentVersion": true,
        "lastSeenAt": true,
        "capacityBytes": true,
        "usedBytes": true,
        "capacityUsedPercent": true,
        "dataReductionRatio": true,
        "openAlertCount": true,
        "criticalAlertCount": true,
        "warningAlertCount": true,
        "volumeCount": true,
        "hostCount": true,
        "podCount": true,
        "fileSystemCount": true,
        "bucketCount": true,
        "hardwareComponentCount": true,
        "unhealthyHardwareCount": true,
        "healthStatus": true,
        "createdByUserId": true,
        "isArchived": true,
        "archivedAt": true,
        "archivedByUserId": true,
        "labels": true,
        "retainTelemetryDataForDays": true,
        "telemetryRetentionConfig": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/storage-array/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read storage_array, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No storage array found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read storage_array: %s", err))
            return
        }
        if wrapper, ok := itemResponse["data"].(map[string]interface{}); ok {
            item = wrapper
        } else {
            item = itemResponse
        }
    }
    if !hasId {
        listBody := map[string]interface{}{
            "query":  filters,
            "select": selectParam,
            // limit 2 is enough to detect ambiguity without paging.
            "limit": 2,
        }
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/storage-array/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list storage_array, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list storage_array: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No storage array matches %s.", DescribeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one storage array matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", DescribeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for storage_array.")
            return
        }
        item = first
    }

    // Update the model with response data
    if obj, ok := item["_id"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Id = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Id = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Id = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Id = types.StringValue(string(jsonBytes))
        } else {
            data.Id = types.StringNull()
        }
    } else if arr, ok := item["_id"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Id = types.StringValue(string(jsonBytes))
        } else {
            data.Id = types.StringNull()
        }
    } else if val, ok := item["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    if obj, ok := item["createdAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CreatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CreatedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CreatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedAt = types.StringNull()
        }
    } else if arr, ok := item["createdAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.CreatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedAt = types.StringNull()
        }
    } else if val, ok := item["createdAt"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    } else {
        data.CreatedAt = types.StringNull()
    }
    if obj, ok := item["updatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UpdatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.UpdatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.UpdatedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.UpdatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.UpdatedAt = types.StringNull()
        }
    } else if arr, ok := item["updatedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.UpdatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.UpdatedAt = types.StringNull()
        }
    } else if val, ok := item["updatedAt"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    } else {
        data.UpdatedAt = types.StringNull()
    }
    if obj, ok := item["createdByActorType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByActorType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CreatedByActorType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CreatedByActorType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CreatedByActorType = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByActorType = types.StringNull()
        }
    } else if arr, ok := item["createdByActorType"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.CreatedByActorType = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByActorType = types.StringNull()
        }
    } else if val, ok := item["createdByActorType"].(string); ok {
        data.CreatedByActorType = types.StringValue(val)
    } else {
        data.CreatedByActorType = types.StringNull()
    }
    if obj, ok := item["createdByActorId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByActorId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CreatedByActorId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CreatedByActorId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CreatedByActorId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByActorId = types.StringNull()
        }
    } else if arr, ok := item["createdByActorId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.CreatedByActorId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByActorId = types.StringNull()
        }
    } else if val, ok := item["createdByActorId"].(string); ok {
        data.CreatedByActorId = types.StringValue(val)
    } else {
        data.CreatedByActorId = types.StringNull()
    }
    if obj, ok := item["lastModifiedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastModifiedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastModifiedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastModifiedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastModifiedAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedAt = types.StringNull()
        }
    } else if arr, ok := item["lastModifiedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LastModifiedAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedAt = types.StringNull()
        }
    } else if val, ok := item["lastModifiedAt"].(string); ok {
        data.LastModifiedAt = types.StringValue(val)
    } else {
        data.LastModifiedAt = types.StringNull()
    }
    if obj, ok := item["lastModifiedByActorType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastModifiedByActorType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastModifiedByActorType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastModifiedByActorType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastModifiedByActorType = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByActorType = types.StringNull()
        }
    } else if arr, ok := item["lastModifiedByActorType"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LastModifiedByActorType = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByActorType = types.StringNull()
        }
    } else if val, ok := item["lastModifiedByActorType"].(string); ok {
        data.LastModifiedByActorType = types.StringValue(val)
    } else {
        data.LastModifiedByActorType = types.StringNull()
    }
    if obj, ok := item["lastModifiedByActorId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastModifiedByActorId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastModifiedByActorId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastModifiedByActorId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastModifiedByActorId = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByActorId = types.StringNull()
        }
    } else if arr, ok := item["lastModifiedByActorId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LastModifiedByActorId = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByActorId = types.StringNull()
        }
    } else if val, ok := item["lastModifiedByActorId"].(string); ok {
        data.LastModifiedByActorId = types.StringValue(val)
    } else {
        data.LastModifiedByActorId = types.StringNull()
    }
    if obj, ok := item["lastModifiedByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastModifiedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastModifiedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastModifiedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastModifiedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByUserId = types.StringNull()
        }
    } else if arr, ok := item["lastModifiedByUserId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LastModifiedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByUserId = types.StringNull()
        }
    } else if val, ok := item["lastModifiedByUserId"].(string); ok {
        data.LastModifiedByUserId = types.StringValue(val)
    } else {
        data.LastModifiedByUserId = types.StringNull()
    }
    if obj, ok := item["projectId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ProjectId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ProjectId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ProjectId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ProjectId = types.StringValue(string(jsonBytes))
        } else {
            data.ProjectId = types.StringNull()
        }
    } else if arr, ok := item["projectId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.ProjectId = types.StringValue(string(jsonBytes))
        } else {
            data.ProjectId = types.StringNull()
        }
    } else if val, ok := item["projectId"].(string); ok {
        data.ProjectId = types.StringValue(val)
    } else {
        data.ProjectId = types.StringNull()
    }
    if obj, ok := item["name"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Name = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Name = types.StringValue(string(jsonBytes))
        } else {
            data.Name = types.StringNull()
        }
    } else if arr, ok := item["name"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Name = types.StringValue(string(jsonBytes))
        } else {
            data.Name = types.StringNull()
        }
    } else if val, ok := item["name"].(string); ok {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if obj, ok := item["slug"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Slug = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Slug = types.StringValue(string(jsonBytes))
        } else {
            data.Slug = types.StringNull()
        }
    } else if arr, ok := item["slug"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Slug = types.StringValue(string(jsonBytes))
        } else {
            data.Slug = types.StringNull()
        }
    } else if val, ok := item["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
    }
    if obj, ok := item["description"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Description = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Description = types.StringValue(string(jsonBytes))
        } else {
            data.Description = types.StringNull()
        }
    } else if arr, ok := item["description"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Description = types.StringValue(string(jsonBytes))
        } else {
            data.Description = types.StringNull()
        }
    } else if val, ok := item["description"].(string); ok {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if obj, ok := item["storageSystem"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StorageSystem = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.StorageSystem = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.StorageSystem = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.StorageSystem = types.StringValue(string(jsonBytes))
        } else {
            data.StorageSystem = types.StringNull()
        }
    } else if arr, ok := item["storageSystem"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.StorageSystem = types.StringValue(string(jsonBytes))
        } else {
            data.StorageSystem = types.StringNull()
        }
    } else if val, ok := item["storageSystem"].(string); ok {
        data.StorageSystem = types.StringValue(val)
    } else {
        data.StorageSystem = types.StringNull()
    }
    if obj, ok := item["reportedName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ReportedName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ReportedName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ReportedName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ReportedName = types.StringValue(string(jsonBytes))
        } else {
            data.ReportedName = types.StringNull()
        }
    } else if arr, ok := item["reportedName"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.ReportedName = types.StringValue(string(jsonBytes))
        } else {
            data.ReportedName = types.StringNull()
        }
    } else if val, ok := item["reportedName"].(string); ok {
        data.ReportedName = types.StringValue(val)
    } else {
        data.ReportedName = types.StringNull()
    }
    if obj, ok := item["systemId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SystemId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SystemId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SystemId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SystemId = types.StringValue(string(jsonBytes))
        } else {
            data.SystemId = types.StringNull()
        }
    } else if arr, ok := item["systemId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.SystemId = types.StringValue(string(jsonBytes))
        } else {
            data.SystemId = types.StringNull()
        }
    } else if val, ok := item["systemId"].(string); ok {
        data.SystemId = types.StringValue(val)
    } else {
        data.SystemId = types.StringNull()
    }
    if obj, ok := item["osName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OsName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.OsName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.OsName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.OsName = types.StringValue(string(jsonBytes))
        } else {
            data.OsName = types.StringNull()
        }
    } else if arr, ok := item["osName"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.OsName = types.StringValue(string(jsonBytes))
        } else {
            data.OsName = types.StringNull()
        }
    } else if val, ok := item["osName"].(string); ok {
        data.OsName = types.StringValue(val)
    } else {
        data.OsName = types.StringNull()
    }
    if obj, ok := item["osVersion"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OsVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.OsVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.OsVersion = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.OsVersion = types.StringValue(string(jsonBytes))
        } else {
            data.OsVersion = types.StringNull()
        }
    } else if arr, ok := item["osVersion"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.OsVersion = types.StringValue(string(jsonBytes))
        } else {
            data.OsVersion = types.StringNull()
        }
    } else if val, ok := item["osVersion"].(string); ok {
        data.OsVersion = types.StringValue(val)
    } else {
        data.OsVersion = types.StringNull()
    }
    if obj, ok := item["otelCollectorStatus"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OtelCollectorStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.OtelCollectorStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.OtelCollectorStatus = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.OtelCollectorStatus = types.StringValue(string(jsonBytes))
        } else {
            data.OtelCollectorStatus = types.StringNull()
        }
    } else if arr, ok := item["otelCollectorStatus"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.OtelCollectorStatus = types.StringValue(string(jsonBytes))
        } else {
            data.OtelCollectorStatus = types.StringNull()
        }
    } else if val, ok := item["otelCollectorStatus"].(string); ok {
        data.OtelCollectorStatus = types.StringValue(val)
    } else {
        data.OtelCollectorStatus = types.StringNull()
    }
    if obj, ok := item["agentVersion"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AgentVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AgentVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AgentVersion = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AgentVersion = types.StringValue(string(jsonBytes))
        } else {
            data.AgentVersion = types.StringNull()
        }
    } else if arr, ok := item["agentVersion"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.AgentVersion = types.StringValue(string(jsonBytes))
        } else {
            data.AgentVersion = types.StringNull()
        }
    } else if val, ok := item["agentVersion"].(string); ok {
        data.AgentVersion = types.StringValue(val)
    } else {
        data.AgentVersion = types.StringNull()
    }
    if obj, ok := item["lastSeenAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastSeenAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastSeenAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastSeenAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastSeenAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastSeenAt = types.StringNull()
        }
    } else if arr, ok := item["lastSeenAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LastSeenAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastSeenAt = types.StringNull()
        }
    } else if val, ok := item["lastSeenAt"].(string); ok {
        data.LastSeenAt = types.StringValue(val)
    } else {
        data.LastSeenAt = types.StringNull()
    }
    data.CapacityBytes = NumberFromAPI(data.CapacityBytes, item["capacityBytes"])
    data.UsedBytes = NumberFromAPI(data.UsedBytes, item["usedBytes"])
    data.CapacityUsedPercent = NumberFromAPI(data.CapacityUsedPercent, item["capacityUsedPercent"])
    data.DataReductionRatio = NumberFromAPI(data.DataReductionRatio, item["dataReductionRatio"])
    data.OpenAlertCount = NumberFromAPI(data.OpenAlertCount, item["openAlertCount"])
    data.CriticalAlertCount = NumberFromAPI(data.CriticalAlertCount, item["criticalAlertCount"])
    data.WarningAlertCount = NumberFromAPI(data.WarningAlertCount, item["warningAlertCount"])
    data.VolumeCount = NumberFromAPI(data.VolumeCount, item["volumeCount"])
    data.HostCount = NumberFromAPI(data.HostCount, item["hostCount"])
    data.PodCount = NumberFromAPI(data.PodCount, item["podCount"])
    data.FileSystemCount = NumberFromAPI(data.FileSystemCount, item["fileSystemCount"])
    data.BucketCount = NumberFromAPI(data.BucketCount, item["bucketCount"])
    data.HardwareComponentCount = NumberFromAPI(data.HardwareComponentCount, item["hardwareComponentCount"])
    data.UnhealthyHardwareCount = NumberFromAPI(data.UnhealthyHardwareCount, item["unhealthyHardwareCount"])
    data.HealthStatus = NumberFromAPI(data.HealthStatus, item["healthStatus"])
    if obj, ok := item["createdByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if arr, ok := item["createdByUserId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if val, ok := item["createdByUserId"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if val, ok := item["isArchived"].(bool); ok {
        data.IsArchived = types.BoolValue(val)
    } else {
        data.IsArchived = types.BoolNull()
    }
    if obj, ok := item["archivedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ArchivedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ArchivedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ArchivedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ArchivedAt = types.StringValue(string(jsonBytes))
        } else {
            data.ArchivedAt = types.StringNull()
        }
    } else if arr, ok := item["archivedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.ArchivedAt = types.StringValue(string(jsonBytes))
        } else {
            data.ArchivedAt = types.StringNull()
        }
    } else if val, ok := item["archivedAt"].(string); ok {
        data.ArchivedAt = types.StringValue(val)
    } else {
        data.ArchivedAt = types.StringNull()
    }
    if obj, ok := item["archivedByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ArchivedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ArchivedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ArchivedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ArchivedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.ArchivedByUserId = types.StringNull()
        }
    } else if arr, ok := item["archivedByUserId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.ArchivedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.ArchivedByUserId = types.StringNull()
        }
    } else if val, ok := item["archivedByUserId"].(string); ok {
        data.ArchivedByUserId = types.StringValue(val)
    } else {
        data.ArchivedByUserId = types.StringNull()
    }
    if val, ok := item["labels"].([]interface{}); ok {
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if jsonBytes, err := json.Marshal(itemMap); err == nil {
                    setItems = append(setItems, types.StringValue(string(jsonBytes)))
                }
            } else if str, ok := item.(string); ok {
                setItems = append(setItems, types.StringValue(str))
            } else {
                setItems = append(setItems, types.StringValue(fmt.Sprintf("%v", item)))
            }
        }
        sort.Slice(setItems, func(i, j int) bool {
            return setItems[i].(types.String).ValueString() < setItems[j].(types.String).ValueString()
        })
        data.Labels = types.SetValueMust(types.StringType, setItems)
    } else {
        data.Labels = types.SetNull(types.StringType)
    }
    data.RetainTelemetryDataForDays = NumberFromAPI(data.RetainTelemetryDataForDays, item["retainTelemetryDataForDays"])
    if obj, ok := item["telemetryRetentionConfig"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TelemetryRetentionConfig = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TelemetryRetentionConfig = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TelemetryRetentionConfig = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TelemetryRetentionConfig = types.StringValue(string(jsonBytes))
        } else {
            data.TelemetryRetentionConfig = types.StringNull()
        }
    } else if arr, ok := item["telemetryRetentionConfig"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.TelemetryRetentionConfig = types.StringValue(string(jsonBytes))
        } else {
            data.TelemetryRetentionConfig = types.StringNull()
        }
    } else if val, ok := item["telemetryRetentionConfig"].(string); ok {
        data.TelemetryRetentionConfig = types.StringValue(val)
    } else {
        data.TelemetryRetentionConfig = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
