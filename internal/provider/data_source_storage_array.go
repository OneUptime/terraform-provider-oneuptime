package provider

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "math/big"
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
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
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
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
    Labels types.Set `tfsdk:"labels"`
    RetainTelemetryDataForDays types.Number `tfsdk:"retain_telemetry_data_for_days"`
    TelemetryRetentionConfig types.String `tfsdk:"telemetry_retention_config"`
}

func (d *StorageArrayDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_storage_array"
}

func (d *StorageArrayDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Storage arrays (Pure Storage FlashArray and FlashBlade) that are being monitored in this project. Each array is auto-discovered when the OneUptime Storage Array Agent sends metrics, or can be registered by hand. Look up an existing storage_array by `id` or by `name`.",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Look up by unique identifier. Exactly one of `id` or `name` must be set.",
                Optional: true,
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.",
                Optional: true,
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
                MarkdownDescription: "Friendly globally unique name for your object.",
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Friendly description for this storage array.",
                Computed: true,
            },
            "storage_system": schema.StringAttribute{
                MarkdownDescription: "The storage platform this array runs, normalized from the storage.system OTel resource attribute (or detected from the metric names): purestorage.flasharray or purestorage.flashblade..",
                Computed: true,
            },
            "reported_name": schema.StringAttribute{
                MarkdownDescription: "The name the array reports for itself (array_name on purefa_info / purefb_info). Can differ from the OneUptime name, which comes from the agent configuration..",
                Computed: true,
            },
            "system_id": schema.StringAttribute{
                MarkdownDescription: "The array's own system identifier (system_id on purefa_info / purefb_info)..",
                Computed: true,
            },
            "os_name": schema.StringAttribute{
                MarkdownDescription: "Name of the array's operating system as it reports it (os on purefa_info / purefb_info)..",
                Computed: true,
            },
            "os_version": schema.StringAttribute{
                MarkdownDescription: "Version of the array's operating system as it reports it (version on purefa_info / purefb_info)..",
                Computed: true,
            },
            "otel_collector_status": schema.StringAttribute{
                MarkdownDescription: "Connection status of the OTel Collector agent (connected or disconnected).",
                Computed: true,
            },
            "agent_version": schema.StringAttribute{
                MarkdownDescription: "Version of the OneUptime Storage Array Agent reporting telemetry, as self-reported via the oneuptime.agent.version resource attribute.",
                Computed: true,
            },
            "last_seen_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "capacity_bytes": schema.NumberAttribute{
                MarkdownDescription: "Cached usable capacity of the array in bytes (purefa_array_space_bytes{space=\"capacity\"} / purefb_array_space_bytes{type=\"array\",space=\"capacity\"}). Null until the first array metric batch arrives..",
                Computed: true,
            },
            "used_bytes": schema.NumberAttribute{
                MarkdownDescription: "Cached physical space used on the array in bytes: capacity minus empty space, or capacity times utilization when the empty series is missing. Null until the first array metric batch arrives..",
                Computed: true,
            },
            "capacity_used_percent": schema.NumberAttribute{
                MarkdownDescription: "Cached array space utilization in percent (purefa_array_space_utilization / purefb_array_space_utilization). Stored as decimal so sub-percent precision survives the round trip. Null until the first array metric batch arrives..",
                Computed: true,
            },
            "data_reduction_ratio": schema.NumberAttribute{
                MarkdownDescription: "Cached data reduction ratio of the array (purefa_array_space_data_reduction_ratio / purefb_array_space_data_reduction_ratio), for example 4.2 for 4.2:1. Null until the first array metric batch arrives..",
                Computed: true,
            },
            "open_alert_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of alerts open on the array itself (purefa_alerts_open / purefb_alerts_open series, hidden alerts excluded)..",
                Computed: true,
            },
            "critical_alert_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of open array alerts with critical severity..",
                Computed: true,
            },
            "warning_alert_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of open array alerts with warning severity..",
                Computed: true,
            },
            "volume_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of volumes on the array (FlashArray)..",
                Computed: true,
            },
            "host_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of hosts defined on the array (FlashArray)..",
                Computed: true,
            },
            "pod_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of pods (ActiveCluster / ActiveDR replication containers) on the array (FlashArray)..",
                Computed: true,
            },
            "file_system_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of file systems on the array (FlashBlade)..",
                Computed: true,
            },
            "bucket_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of object store buckets on the array (FlashBlade)..",
                Computed: true,
            },
            "hardware_component_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of hardware components the array reports (chassis, controllers, drive bays, power supplies, fans, ports, blades...)..",
                Computed: true,
            },
            "unhealthy_hardware_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of hardware components, drives and controllers in a critical, degraded, failed or unknown state..",
                Computed: true,
            },
            "health_status": schema.NumberAttribute{
                MarkdownDescription: "Cached array health derived from the array's open alerts and hardware state: 0 = OK, 1 = Warning (a warning alert, or a degraded or unknown component), 2 = Critical (a critical alert, or a failed or critical component). Rendered as the health pill. Null until the first metric batch arrives..",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "is_archived": schema.BoolAttribute{
                MarkdownDescription: "Is this storage array archived? Archived storage arrays are hidden from lists but keep collecting telemetry..",
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
                MarkdownDescription: "Relation to Labels Array where this object is categorized in..",
                Computed: true,
                ElementType: types.StringType,
            },
            "retain_telemetry_data_for_days": schema.NumberAttribute{
                MarkdownDescription: "Number of days to retain telemetry data for this storage array. Leave blank to use the project-wide default..",
                Computed: true,
            },
            "telemetry_retention_config": schema.StringAttribute{
                MarkdownDescription: "Per-pillar retention overrides for this storage array (logs by severity, traces by status, metrics, profiles). Unset fields fall back to the storage array default, then the project's retention settings..",
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

    hasId := !data.Id.IsNull() && data.Id.ValueString() != ""
    hasName := !data.Name.IsNull() && data.Name.ValueString() != ""
    if hasId == hasName {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Exactly one of `id` or `name` must be set to look up a storage_array.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "name": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "projectId": true,
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
        "deletedByUserId": true,
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
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No storage_array found with id %q.", data.Id.ValueString()))
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
    } else {
        listBody := map[string]interface{}{
            "query": map[string]interface{}{
                "name": data.Name.ValueString(),
            },
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
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No storage_array found with name %q.", data.Name.ValueString()))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one storage_array matches name %q. Use the id attribute to disambiguate.", data.Name.ValueString()))
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
    } else if val, ok := item["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
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
    } else if val, ok := item["name"].(string); ok {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
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
    } else if val, ok := item["updatedAt"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    } else {
        data.UpdatedAt = types.StringNull()
    }
    if obj, ok := item["deletedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DeletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DeletedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DeletedAt = types.StringValue(string(jsonBytes))
        } else {
            data.DeletedAt = types.StringNull()
        }
    } else if val, ok := item["deletedAt"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    } else {
        data.DeletedAt = types.StringNull()
    }
    if val, ok := item["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["version"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.Version = types.NumberValue(big.NewFloat(val))
        } else {
            data.Version = types.NumberNull()
        }
    } else {
        data.Version = types.NumberNull()
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
    } else if val, ok := item["projectId"].(string); ok {
        data.ProjectId = types.StringValue(val)
    } else {
        data.ProjectId = types.StringNull()
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
    } else if val, ok := item["lastSeenAt"].(string); ok {
        data.LastSeenAt = types.StringValue(val)
    } else {
        data.LastSeenAt = types.StringNull()
    }
    if val, ok := item["capacityBytes"].(float64); ok {
        data.CapacityBytes = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["capacityBytes"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.CapacityBytes = types.NumberValue(big.NewFloat(val))
        } else {
            data.CapacityBytes = types.NumberNull()
        }
    } else {
        data.CapacityBytes = types.NumberNull()
    }
    if val, ok := item["usedBytes"].(float64); ok {
        data.UsedBytes = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["usedBytes"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.UsedBytes = types.NumberValue(big.NewFloat(val))
        } else {
            data.UsedBytes = types.NumberNull()
        }
    } else {
        data.UsedBytes = types.NumberNull()
    }
    if val, ok := item["capacityUsedPercent"].(float64); ok {
        data.CapacityUsedPercent = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["capacityUsedPercent"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.CapacityUsedPercent = types.NumberValue(big.NewFloat(val))
        } else {
            data.CapacityUsedPercent = types.NumberNull()
        }
    } else {
        data.CapacityUsedPercent = types.NumberNull()
    }
    if val, ok := item["dataReductionRatio"].(float64); ok {
        data.DataReductionRatio = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["dataReductionRatio"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.DataReductionRatio = types.NumberValue(big.NewFloat(val))
        } else {
            data.DataReductionRatio = types.NumberNull()
        }
    } else {
        data.DataReductionRatio = types.NumberNull()
    }
    if val, ok := item["openAlertCount"].(float64); ok {
        data.OpenAlertCount = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["openAlertCount"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.OpenAlertCount = types.NumberValue(big.NewFloat(val))
        } else {
            data.OpenAlertCount = types.NumberNull()
        }
    } else {
        data.OpenAlertCount = types.NumberNull()
    }
    if val, ok := item["criticalAlertCount"].(float64); ok {
        data.CriticalAlertCount = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["criticalAlertCount"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.CriticalAlertCount = types.NumberValue(big.NewFloat(val))
        } else {
            data.CriticalAlertCount = types.NumberNull()
        }
    } else {
        data.CriticalAlertCount = types.NumberNull()
    }
    if val, ok := item["warningAlertCount"].(float64); ok {
        data.WarningAlertCount = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["warningAlertCount"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.WarningAlertCount = types.NumberValue(big.NewFloat(val))
        } else {
            data.WarningAlertCount = types.NumberNull()
        }
    } else {
        data.WarningAlertCount = types.NumberNull()
    }
    if val, ok := item["volumeCount"].(float64); ok {
        data.VolumeCount = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["volumeCount"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.VolumeCount = types.NumberValue(big.NewFloat(val))
        } else {
            data.VolumeCount = types.NumberNull()
        }
    } else {
        data.VolumeCount = types.NumberNull()
    }
    if val, ok := item["hostCount"].(float64); ok {
        data.HostCount = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["hostCount"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.HostCount = types.NumberValue(big.NewFloat(val))
        } else {
            data.HostCount = types.NumberNull()
        }
    } else {
        data.HostCount = types.NumberNull()
    }
    if val, ok := item["podCount"].(float64); ok {
        data.PodCount = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["podCount"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.PodCount = types.NumberValue(big.NewFloat(val))
        } else {
            data.PodCount = types.NumberNull()
        }
    } else {
        data.PodCount = types.NumberNull()
    }
    if val, ok := item["fileSystemCount"].(float64); ok {
        data.FileSystemCount = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["fileSystemCount"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.FileSystemCount = types.NumberValue(big.NewFloat(val))
        } else {
            data.FileSystemCount = types.NumberNull()
        }
    } else {
        data.FileSystemCount = types.NumberNull()
    }
    if val, ok := item["bucketCount"].(float64); ok {
        data.BucketCount = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["bucketCount"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.BucketCount = types.NumberValue(big.NewFloat(val))
        } else {
            data.BucketCount = types.NumberNull()
        }
    } else {
        data.BucketCount = types.NumberNull()
    }
    if val, ok := item["hardwareComponentCount"].(float64); ok {
        data.HardwareComponentCount = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["hardwareComponentCount"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.HardwareComponentCount = types.NumberValue(big.NewFloat(val))
        } else {
            data.HardwareComponentCount = types.NumberNull()
        }
    } else {
        data.HardwareComponentCount = types.NumberNull()
    }
    if val, ok := item["unhealthyHardwareCount"].(float64); ok {
        data.UnhealthyHardwareCount = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["unhealthyHardwareCount"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.UnhealthyHardwareCount = types.NumberValue(big.NewFloat(val))
        } else {
            data.UnhealthyHardwareCount = types.NumberNull()
        }
    } else {
        data.UnhealthyHardwareCount = types.NumberNull()
    }
    if val, ok := item["healthStatus"].(float64); ok {
        data.HealthStatus = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["healthStatus"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.HealthStatus = types.NumberValue(big.NewFloat(val))
        } else {
            data.HealthStatus = types.NumberNull()
        }
    } else {
        data.HealthStatus = types.NumberNull()
    }
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
    } else if val, ok := item["archivedByUserId"].(string); ok {
        data.ArchivedByUserId = types.StringValue(val)
    } else {
        data.ArchivedByUserId = types.StringNull()
    }
    if obj, ok := item["deletedByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeletedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DeletedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DeletedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DeletedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.DeletedByUserId = types.StringNull()
        }
    } else if val, ok := item["deletedByUserId"].(string); ok {
        data.DeletedByUserId = types.StringValue(val)
    } else {
        data.DeletedByUserId = types.StringNull()
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
    if val, ok := item["retainTelemetryDataForDays"].(float64); ok {
        data.RetainTelemetryDataForDays = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["retainTelemetryDataForDays"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.RetainTelemetryDataForDays = types.NumberValue(big.NewFloat(val))
        } else {
            data.RetainTelemetryDataForDays = types.NumberNull()
        }
    } else {
        data.RetainTelemetryDataForDays = types.NumberNull()
    }
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
