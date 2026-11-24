package shard03

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
    "context"
    "fmt"
    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-framework/types/basetypes"
    "github.com/hashicorp/terraform-plugin-log/tflog"
    "math/big"
    "net/http"
    "github.com/hashicorp/terraform-plugin-framework/path"
    "encoding/json"
    "net/url"
    "strings"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
    "github.com/hashicorp/terraform-plugin-framework/attr"
    "sort"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/numberplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &StorageArrayResource{}
var _ resource.ResourceWithImportState = &StorageArrayResource{}
var _ resource.ResourceWithModifyPlan = &StorageArrayResource{}

func NewStorageArrayResource() resource.Resource {
    return &StorageArrayResource{}
}

// StorageArrayResource defines the resource implementation.
type StorageArrayResource struct {
    client *Client
}

// StorageArrayResourceModel describes the resource data model.
type StorageArrayResourceModel struct {
    Id types.String `tfsdk:"id"`
    ProjectId types.String `tfsdk:"project_id"`
    Name types.String `tfsdk:"name"`
    Description types.String `tfsdk:"description"`
    IsArchived types.Bool `tfsdk:"is_archived"`
    Labels types.Set `tfsdk:"labels"`
    RetainTelemetryDataForDays types.Number `tfsdk:"retain_telemetry_data_for_days"`
    TelemetryRetentionConfig JSONSubsetValue `tfsdk:"telemetry_retention_config"`
    CreatedAt RFC3339Value `tfsdk:"created_at"`
    UpdatedAt RFC3339Value `tfsdk:"updated_at"`
    Slug types.String `tfsdk:"slug"`
    StorageSystem types.String `tfsdk:"storage_system"`
    ReportedName types.String `tfsdk:"reported_name"`
    SystemId types.String `tfsdk:"system_id"`
    OsName types.String `tfsdk:"os_name"`
    OsVersion types.String `tfsdk:"os_version"`
    OtelCollectorStatus types.String `tfsdk:"otel_collector_status"`
    AgentVersion types.String `tfsdk:"agent_version"`
    LastSeenAt RFC3339Value `tfsdk:"last_seen_at"`
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
    ArchivedAt RFC3339Value `tfsdk:"archived_at"`
    ArchivedByUserId types.String `tfsdk:"archived_by_user_id"`
}

func (r *StorageArrayResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_storage_array"
}

func (r *StorageArrayResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = r.schemaDefinition()
}

func (r *StorageArrayResource) schemaDefinition() schema.Schema {
    return schema.Schema{
        MarkdownDescription: "Storage arrays (Pure Storage FlashArray and FlashBlade) that are being monitored in this project. Each array is auto-discovered when the OneUptime Storage Array Agent sends metrics, or can be registered by hand.",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Unique identifier for the resource.",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Name of this storage array in OneUptime. This is the join key — it must match the storage.array.name OTel resource attribute stamped by the OneUptime Storage Array Agent.",
                Required: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Friendly description for this storage array.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "is_archived": schema.BoolAttribute{
                MarkdownDescription: "Is this storage array archived? Archived storage arrays are hidden from lists but keep collecting telemetry.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "labels": schema.SetAttribute{
                MarkdownDescription: "Labels to group and filter by. OneUptime adds to it too, on its own: the entries listed here are added and taken off, and every other entry stays. Left out, it keeps the entries OneUptime gives it. IDs of `oneuptime_label` resources.",
                Optional: true,
                Computed: true,
                ElementType: types.StringType,
                PlanModifiers: []planmodifier.Set{
                    setplanmodifier.UseStateForUnknown(),
                },
                Validators: RecordIDSetValidators("oneuptime_label"),
            },
            "retain_telemetry_data_for_days": schema.NumberAttribute{
                MarkdownDescription: "Number of days to retain telemetry data for this storage array. Leave blank to use the project-wide default.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.Number{
                    NumberKeepSameAPINumber(),
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "telemetry_retention_config": schema.StringAttribute{
                MarkdownDescription: "Per-pillar retention overrides for this storage array (logs by severity, traces by status, metrics, profiles). Unset fields fall back to the storage array default, then the project's retention settings. A JSON value: write it with `jsonencode()`.",
                CustomType: JSONSubsetType{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    JSONSubsetKeepImported(),
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: []validator.String{
                    JSONEnvelopeValidator(),
                },
            },
            "created_at": schema.StringAttribute{
                MarkdownDescription: "Date and Time when the object was created.",
                CustomType: RFC3339Type{},
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "updated_at": schema.StringAttribute{
                MarkdownDescription: "Date and Time when the object was updated.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Friendly globally unique name for your object.",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "storage_system": schema.StringAttribute{
                MarkdownDescription: "The storage platform this array runs, normalized from the storage.system OTel resource attribute (or detected from the metric names): purestorage.flasharray or purestorage.flashblade.",
                Computed: true,
            },
            "reported_name": schema.StringAttribute{
                MarkdownDescription: "The name the array reports for itself (array_name on purefa_info / purefb_info). Can differ from the OneUptime name, which comes from the agent configuration.",
                Computed: true,
            },
            "system_id": schema.StringAttribute{
                MarkdownDescription: "The array's own system identifier (system_id on purefa_info / purefb_info).",
                Computed: true,
            },
            "os_name": schema.StringAttribute{
                MarkdownDescription: "Name of the array's operating system as it reports it (os on purefa_info / purefb_info).",
                Computed: true,
            },
            "os_version": schema.StringAttribute{
                MarkdownDescription: "Version of the array's operating system as it reports it (version on purefa_info / purefb_info).",
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
                MarkdownDescription: "When metrics were last received from this storage array.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "capacity_bytes": schema.NumberAttribute{
                MarkdownDescription: "Cached usable capacity of the array in bytes (purefa_array_space_bytes{space=\"capacity\"} / purefb_array_space_bytes{type=\"array\",space=\"capacity\"}). Null until the first array metric batch arrives.",
                Computed: true,
            },
            "used_bytes": schema.NumberAttribute{
                MarkdownDescription: "Cached physical space used on the array in bytes: capacity minus empty space, or capacity times utilization when the empty series is missing. Null until the first array metric batch arrives.",
                Computed: true,
            },
            "capacity_used_percent": schema.NumberAttribute{
                MarkdownDescription: "Cached array space utilization in percent (purefa_array_space_utilization / purefb_array_space_utilization). Stored as decimal so sub-percent precision survives the round trip. Null until the first array metric batch arrives.",
                Computed: true,
            },
            "data_reduction_ratio": schema.NumberAttribute{
                MarkdownDescription: "Cached data reduction ratio of the array (purefa_array_space_data_reduction_ratio / purefb_array_space_data_reduction_ratio), for example 4.2 for 4.2:1. Null until the first array metric batch arrives.",
                Computed: true,
            },
            "open_alert_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of alerts open on the array itself (purefa_alerts_open / purefb_alerts_open series, hidden alerts excluded).",
                Computed: true,
            },
            "critical_alert_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of open array alerts with critical severity.",
                Computed: true,
            },
            "warning_alert_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of open array alerts with warning severity.",
                Computed: true,
            },
            "volume_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of volumes on the array (FlashArray).",
                Computed: true,
            },
            "host_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of hosts defined on the array (FlashArray).",
                Computed: true,
            },
            "pod_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of pods (ActiveCluster / ActiveDR replication containers) on the array (FlashArray).",
                Computed: true,
            },
            "file_system_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of file systems on the array (FlashBlade).",
                Computed: true,
            },
            "bucket_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of object store buckets on the array (FlashBlade).",
                Computed: true,
            },
            "hardware_component_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of hardware components the array reports (chassis, controllers, drive bays, power supplies, fans, ports, blades...).",
                Computed: true,
            },
            "unhealthy_hardware_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of hardware components, drives and controllers in a critical, degraded, failed or unknown state.",
                Computed: true,
            },
            "health_status": schema.NumberAttribute{
                MarkdownDescription: "Cached array health derived from the array's open alerts and hardware state: 0 = OK, 1 = Warning (a warning alert, or a degraded or unknown component), 2 = Critical (a critical alert, or a failed or critical component). Rendered as the health pill. Null until the first metric batch arrives.",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "archived_at": schema.StringAttribute{
                MarkdownDescription: "When was this storage array archived?",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "archived_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who archived this object (if this object was archived by a User). The ID of a `oneuptime_user` (see the data source).",
                Computed: true,
            },
        },
    }
}

// ModifyPlan keeps the state when nothing a configuration sets changes (see
// importsettle.go).
func (r *StorageArrayResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
    // Destroying: nothing to plan.
    if req.Plan.Raw.IsNull() {
        return
    }

    // Under Crossplane, what a manifest leaves out is the record's (see crossplane.go).
    KeepLeftOutUnderCrossplane(ctx, req, resp)

    // What people and OneUptime both change: OneUptime's change is never one
    // to undo (see shared.go).
    PlanShared(ctx, req, resp, r.client, []SharedAttribute{
        {Name: "labels", Entries: true},
    })

    KeepStateWhenNothingChanges(ctx, req, resp)
}

func (r *StorageArrayResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
    // Prevent panic if the provider has not been configured.
    if req.ProviderData == nil {
        return
    }

    client, ok := req.ProviderData.(*Client)

    if !ok {
        resp.Diagnostics.AddError(
            "Unexpected Resource Configure Type",
            fmt.Sprintf("Expected *Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
        )

        return
    }

    r.client = client
}


func (r *StorageArrayResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var data StorageArrayResourceModel

    // Read Terraform plan data into the model
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // What the configuration sets, and what Terraform planned (see keepPlannedValues).
    var config StorageArrayResourceModel
    resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
    if resp.Diagnostics.HasError() {
        return
    }
    plan := data



    // Create API request body. Unset (null/unknown) optional fields are
    // omitted so server-side defaults apply instead of being overwritten
    // with zero values.
    storageArrayRequest := map[string]interface{}{
        "data": map[string]interface{}{},
    }
    requestDataMap := storageArrayRequest["data"].(map[string]interface{})

    if !data.Name.IsNull() && !data.Name.IsUnknown() {
        requestDataMap["name"] = data.Name.ValueString()
    }
    if !data.Description.IsNull() && !data.Description.IsUnknown() {
        requestDataMap["description"] = data.Description.ValueString()
    }
    if !data.IsArchived.IsNull() && !data.IsArchived.IsUnknown() {
        requestDataMap["isArchived"] = data.IsArchived.ValueBool()
    }
    if !data.Labels.IsNull() && !data.Labels.IsUnknown() {
        requestDataMap["labels"] = r.convertTerraformSetToInterface(data.Labels)
    }
    if !data.RetainTelemetryDataForDays.IsNull() && !data.RetainTelemetryDataForDays.IsUnknown() {
        requestDataMap["retainTelemetryDataForDays"] = r.bigFloatToFloat64(data.RetainTelemetryDataForDays.ValueBigFloat())
    }
    if parsedTelemetryRetentionConfig := r.parseJSONField(data.TelemetryRetentionConfig); parsedTelemetryRetentionConfig != nil {
        requestDataMap["telemetryRetentionConfig"] = parsedTelemetryRetentionConfig
    }

    // Make API call
    httpResp, err := r.client.Post(ctx, "/storage-array", storageArrayRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create storage_array, got error: %s", err))
        return
    }

    var storageArrayResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &storageArrayResponse)
    if err != nil {
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to create storage_array: %s", err))
        return
    }

    // What OneUptime left out because only it writes it, if anything (see entries.go).
    resp.Diagnostics.Append(FieldsLeftOutWarning(httpResp, "storage array")...)

    // Extract the new resource id from the create response.
    createdId := ""
    if wrapper, ok := storageArrayResponse["data"].(map[string]interface{}); ok {
        if val, ok := wrapper["_id"].(string); ok {
            createdId = val
        }
    } else if val, ok := storageArrayResponse["_id"].(string); ok {
        createdId = val
    }
    if createdId == "" {
        resp.Diagnostics.AddError("OneUptime API Error", "Create response for storage_array did not contain an id. This is a bug in the provider or the API; please report it.")
        return
    }
    data.Id = types.StringValue(createdId)


    /*
     * The server has committed the row. Persist what we know to state BEFORE
     * the read-back: if the read-back fails and we return without setting
     * state, Terraform never learns the resource exists and the created
     * storage_array is orphaned server-side — never refreshed, never
     * destroyed. Delete already refuses to drop state on failure for the
     * same reason; Create must not either.
     */
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Re-read the resource so state reflects server-normalized values.
    selectParam := map[string]interface{}{
        "projectId": true,
        "name": true,
        "description": true,
        "isArchived": true,
        "labels": true,
        "retainTelemetryDataForDays": true,
        "telemetryRetentionConfig": true,
        "createdAt": true,
        "updatedAt": true,
        "slug": true,
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
        "archivedAt": true,
        "archivedByUserId": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect(ctx, "/storage-array/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        /*
         * State already owns the id, so the resource is tracked and the next
         * refresh reconciles the remaining attributes. Warn rather than
         * error: erroring here would strand a real resource.
         */
        resp.Diagnostics.AddWarning("Read After Create Failed", fmt.Sprintf("Created storage_array but could not read it back; state is incomplete until the next refresh: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddWarning("Read After Create Failed", fmt.Sprintf("Created storage_array but could not parse the read-back response; state is incomplete until the next refresh: %s", err))
        return
    }

    // Update the model with the authoritative read response
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := readResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = readResponse
    }

    if obj, ok := dataMap["projectId"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok {
            data.ProjectId = types.StringValue(val)
        } else {
            data.ProjectId = types.StringNull()
        }
    } else if val, ok := dataMap["projectId"].(string); ok {
        data.ProjectId = types.StringValue(val)
    } else {
        data.ProjectId = types.StringNull()
    }
    if obj, ok := dataMap["name"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Name = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Name = types.StringValue(string(jsonBytes))
            } else {
                data.Name = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Name = types.StringValue(string(jsonBytes))
            } else {
                data.Name = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Name = types.StringValue(string(jsonBytes))
        } else {
            data.Name = types.StringNull()
        }
    } else if arr, ok := dataMap["name"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Name = types.StringValue(string(jsonBytes))
        } else {
            data.Name = types.StringNull()
        }
    } else if val, ok := dataMap["name"].(string); ok {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if obj, ok := dataMap["description"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Description = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Description = types.StringValue(string(jsonBytes))
            } else {
                data.Description = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Description = types.StringValue(string(jsonBytes))
            } else {
                data.Description = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Description = types.StringValue(string(jsonBytes))
        } else {
            data.Description = types.StringNull()
        }
    } else if arr, ok := dataMap["description"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Description = types.StringValue(string(jsonBytes))
        } else {
            data.Description = types.StringNull()
        }
    } else if val, ok := dataMap["description"].(string); ok {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if val, ok := dataMap["isArchived"].(bool); ok {
        data.IsArchived = types.BoolValue(val)
    }
    if val, ok := dataMap["labels"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.Labels = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.Labels = types.SetValueMust(types.StringType, []attr.Value{})
    }
    data.RetainTelemetryDataForDays = NumberFromAPI(data.RetainTelemetryDataForDays, dataMap["retainTelemetryDataForDays"])
    if obj, ok := dataMap["telemetryRetentionConfig"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TelemetryRetentionConfig = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.TelemetryRetentionConfig = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.TelemetryRetentionConfig = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.TelemetryRetentionConfig = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.TelemetryRetentionConfig = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.TelemetryRetentionConfig = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.TelemetryRetentionConfig = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.TelemetryRetentionConfig = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.TelemetryRetentionConfig = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["telemetryRetentionConfig"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.TelemetryRetentionConfig = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.TelemetryRetentionConfig = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["telemetryRetentionConfig"].(string); ok {
        data.TelemetryRetentionConfig = NewJSONSubsetValue(val)
    } else {
        data.TelemetryRetentionConfig = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.CreatedAt = NewRFC3339Value(val)
        } else {
            data.CreatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["createdAt"].(string); ok && val != "" {
        data.CreatedAt = NewRFC3339Value(val)
    } else {
        data.CreatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.UpdatedAt = NewRFC3339Value(val)
        } else {
            data.UpdatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["updatedAt"].(string); ok && val != "" {
        data.UpdatedAt = NewRFC3339Value(val)
    } else {
        data.UpdatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["slug"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Slug = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Slug = types.StringValue(string(jsonBytes))
            } else {
                data.Slug = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Slug = types.StringValue(string(jsonBytes))
            } else {
                data.Slug = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Slug = types.StringValue(string(jsonBytes))
        } else {
            data.Slug = types.StringNull()
        }
    } else if arr, ok := dataMap["slug"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Slug = types.StringValue(string(jsonBytes))
        } else {
            data.Slug = types.StringNull()
        }
    } else if val, ok := dataMap["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
    }
    if obj, ok := dataMap["storageSystem"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StorageSystem = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.StorageSystem = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.StorageSystem = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.StorageSystem = types.StringValue(string(jsonBytes))
            } else {
                data.StorageSystem = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.StorageSystem = types.StringValue(string(jsonBytes))
            } else {
                data.StorageSystem = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.StorageSystem = types.StringValue(string(jsonBytes))
        } else {
            data.StorageSystem = types.StringNull()
        }
    } else if arr, ok := dataMap["storageSystem"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.StorageSystem = types.StringValue(string(jsonBytes))
        } else {
            data.StorageSystem = types.StringNull()
        }
    } else if val, ok := dataMap["storageSystem"].(string); ok {
        data.StorageSystem = types.StringValue(val)
    } else {
        data.StorageSystem = types.StringNull()
    }
    if obj, ok := dataMap["reportedName"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ReportedName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ReportedName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ReportedName = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ReportedName = types.StringValue(string(jsonBytes))
            } else {
                data.ReportedName = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ReportedName = types.StringValue(string(jsonBytes))
            } else {
                data.ReportedName = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ReportedName = types.StringValue(string(jsonBytes))
        } else {
            data.ReportedName = types.StringNull()
        }
    } else if arr, ok := dataMap["reportedName"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.ReportedName = types.StringValue(string(jsonBytes))
        } else {
            data.ReportedName = types.StringNull()
        }
    } else if val, ok := dataMap["reportedName"].(string); ok {
        data.ReportedName = types.StringValue(val)
    } else {
        data.ReportedName = types.StringNull()
    }
    if obj, ok := dataMap["systemId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SystemId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SystemId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SystemId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SystemId = types.StringValue(string(jsonBytes))
            } else {
                data.SystemId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SystemId = types.StringValue(string(jsonBytes))
            } else {
                data.SystemId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SystemId = types.StringValue(string(jsonBytes))
        } else {
            data.SystemId = types.StringNull()
        }
    } else if arr, ok := dataMap["systemId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.SystemId = types.StringValue(string(jsonBytes))
        } else {
            data.SystemId = types.StringNull()
        }
    } else if val, ok := dataMap["systemId"].(string); ok {
        data.SystemId = types.StringValue(val)
    } else {
        data.SystemId = types.StringNull()
    }
    if obj, ok := dataMap["osName"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OsName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.OsName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.OsName = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.OsName = types.StringValue(string(jsonBytes))
            } else {
                data.OsName = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.OsName = types.StringValue(string(jsonBytes))
            } else {
                data.OsName = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.OsName = types.StringValue(string(jsonBytes))
        } else {
            data.OsName = types.StringNull()
        }
    } else if arr, ok := dataMap["osName"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.OsName = types.StringValue(string(jsonBytes))
        } else {
            data.OsName = types.StringNull()
        }
    } else if val, ok := dataMap["osName"].(string); ok {
        data.OsName = types.StringValue(val)
    } else {
        data.OsName = types.StringNull()
    }
    if obj, ok := dataMap["osVersion"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OsVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.OsVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.OsVersion = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.OsVersion = types.StringValue(string(jsonBytes))
            } else {
                data.OsVersion = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.OsVersion = types.StringValue(string(jsonBytes))
            } else {
                data.OsVersion = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.OsVersion = types.StringValue(string(jsonBytes))
        } else {
            data.OsVersion = types.StringNull()
        }
    } else if arr, ok := dataMap["osVersion"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.OsVersion = types.StringValue(string(jsonBytes))
        } else {
            data.OsVersion = types.StringNull()
        }
    } else if val, ok := dataMap["osVersion"].(string); ok {
        data.OsVersion = types.StringValue(val)
    } else {
        data.OsVersion = types.StringNull()
    }
    if obj, ok := dataMap["otelCollectorStatus"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OtelCollectorStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.OtelCollectorStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.OtelCollectorStatus = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.OtelCollectorStatus = types.StringValue(string(jsonBytes))
            } else {
                data.OtelCollectorStatus = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.OtelCollectorStatus = types.StringValue(string(jsonBytes))
            } else {
                data.OtelCollectorStatus = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.OtelCollectorStatus = types.StringValue(string(jsonBytes))
        } else {
            data.OtelCollectorStatus = types.StringNull()
        }
    } else if arr, ok := dataMap["otelCollectorStatus"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.OtelCollectorStatus = types.StringValue(string(jsonBytes))
        } else {
            data.OtelCollectorStatus = types.StringNull()
        }
    } else if val, ok := dataMap["otelCollectorStatus"].(string); ok {
        data.OtelCollectorStatus = types.StringValue(val)
    } else {
        data.OtelCollectorStatus = types.StringNull()
    }
    if obj, ok := dataMap["agentVersion"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AgentVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AgentVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AgentVersion = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AgentVersion = types.StringValue(string(jsonBytes))
            } else {
                data.AgentVersion = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AgentVersion = types.StringValue(string(jsonBytes))
            } else {
                data.AgentVersion = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AgentVersion = types.StringValue(string(jsonBytes))
        } else {
            data.AgentVersion = types.StringNull()
        }
    } else if arr, ok := dataMap["agentVersion"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.AgentVersion = types.StringValue(string(jsonBytes))
        } else {
            data.AgentVersion = types.StringNull()
        }
    } else if val, ok := dataMap["agentVersion"].(string); ok {
        data.AgentVersion = types.StringValue(val)
    } else {
        data.AgentVersion = types.StringNull()
    }
    if obj, ok := dataMap["lastSeenAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.LastSeenAt = NewRFC3339Value(val)
        } else {
            data.LastSeenAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["lastSeenAt"].(string); ok && val != "" {
        data.LastSeenAt = NewRFC3339Value(val)
    } else {
        data.LastSeenAt = NewRFC3339Null()
    }
    data.CapacityBytes = NumberFromAPI(data.CapacityBytes, dataMap["capacityBytes"])
    data.UsedBytes = NumberFromAPI(data.UsedBytes, dataMap["usedBytes"])
    data.CapacityUsedPercent = NumberFromAPI(data.CapacityUsedPercent, dataMap["capacityUsedPercent"])
    data.DataReductionRatio = NumberFromAPI(data.DataReductionRatio, dataMap["dataReductionRatio"])
    data.OpenAlertCount = NumberFromAPI(data.OpenAlertCount, dataMap["openAlertCount"])
    data.CriticalAlertCount = NumberFromAPI(data.CriticalAlertCount, dataMap["criticalAlertCount"])
    data.WarningAlertCount = NumberFromAPI(data.WarningAlertCount, dataMap["warningAlertCount"])
    data.VolumeCount = NumberFromAPI(data.VolumeCount, dataMap["volumeCount"])
    data.HostCount = NumberFromAPI(data.HostCount, dataMap["hostCount"])
    data.PodCount = NumberFromAPI(data.PodCount, dataMap["podCount"])
    data.FileSystemCount = NumberFromAPI(data.FileSystemCount, dataMap["fileSystemCount"])
    data.BucketCount = NumberFromAPI(data.BucketCount, dataMap["bucketCount"])
    data.HardwareComponentCount = NumberFromAPI(data.HardwareComponentCount, dataMap["hardwareComponentCount"])
    data.UnhealthyHardwareCount = NumberFromAPI(data.UnhealthyHardwareCount, dataMap["unhealthyHardwareCount"])
    data.HealthStatus = NumberFromAPI(data.HealthStatus, dataMap["healthStatus"])
    if obj, ok := dataMap["createdByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if arr, ok := dataMap["createdByUserId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["createdByUserId"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if obj, ok := dataMap["archivedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.ArchivedAt = NewRFC3339Value(val)
        } else {
            data.ArchivedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["archivedAt"].(string); ok && val != "" {
        data.ArchivedAt = NewRFC3339Value(val)
    } else {
        data.ArchivedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["archivedByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ArchivedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ArchivedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ArchivedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ArchivedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.ArchivedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ArchivedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.ArchivedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ArchivedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.ArchivedByUserId = types.StringNull()
        }
    } else if arr, ok := dataMap["archivedByUserId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.ArchivedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.ArchivedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["archivedByUserId"].(string); ok {
        data.ArchivedByUserId = types.StringValue(val)
    } else {
        data.ArchivedByUserId = types.StringNull()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    // The read response is authoritative, but never let it clobber the id we just received.
    data.Id = types.StringValue(createdId)

    // Unconfigured attributes keep their planned value; see keepPlannedValues.
    r.keepPlannedValues(&data, &plan, &config)

    // A list OneUptime adds to as well holds the entries the configuration declares (see entries.go).
    resp.Diagnostics.Append(r.heldEntriesAfterWrite(ctx, nil, &data, &plan, &config, true).Write(ctx, resp.Private)...)

    // Write logs using the tflog package
    tflog.Trace(ctx, "created a resource")


    // Nothing a create made was imported (see importsettle.go).
    resp.Diagnostics.Append(ForgetImportedFromServer(ctx, resp.Private)...)
    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *StorageArrayResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var data StorageArrayResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // The state before the read: the entries the configuration declared (see entries.go).
    prior := data

    // Create select parameter to get full object
    selectParam := map[string]interface{}{
        "projectId": true,
        "name": true,
        "description": true,
        "isArchived": true,
        "labels": true,
        "retainTelemetryDataForDays": true,
        "telemetryRetentionConfig": true,
        "createdAt": true,
        "updatedAt": true,
        "slug": true,
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
        "archivedAt": true,
        "archivedByUserId": true,
        "_id": true,
    }

    // Make API call with select parameter
    httpResp, err := r.client.PostWithSelect(ctx, "/storage-array/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read storage_array, got error: %s", err))
        return
    }

    if httpResp.StatusCode == http.StatusNotFound {
        resp.State.RemoveResource(ctx)
        return
    }

    var storageArrayResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &storageArrayResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse storage_array response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := storageArrayResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = storageArrayResponse
    }

    if obj, ok := dataMap["projectId"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok {
            data.ProjectId = types.StringValue(val)
        } else {
            data.ProjectId = types.StringNull()
        }
    } else if val, ok := dataMap["projectId"].(string); ok {
        data.ProjectId = types.StringValue(val)
    } else {
        data.ProjectId = types.StringNull()
    }
    if obj, ok := dataMap["name"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Name = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Name = types.StringValue(string(jsonBytes))
            } else {
                data.Name = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Name = types.StringValue(string(jsonBytes))
            } else {
                data.Name = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Name = types.StringValue(string(jsonBytes))
        } else {
            data.Name = types.StringNull()
        }
    } else if arr, ok := dataMap["name"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Name = types.StringValue(string(jsonBytes))
        } else {
            data.Name = types.StringNull()
        }
    } else if val, ok := dataMap["name"].(string); ok {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if obj, ok := dataMap["description"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Description = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Description = types.StringValue(string(jsonBytes))
            } else {
                data.Description = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Description = types.StringValue(string(jsonBytes))
            } else {
                data.Description = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Description = types.StringValue(string(jsonBytes))
        } else {
            data.Description = types.StringNull()
        }
    } else if arr, ok := dataMap["description"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Description = types.StringValue(string(jsonBytes))
        } else {
            data.Description = types.StringNull()
        }
    } else if val, ok := dataMap["description"].(string); ok {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if val, ok := dataMap["isArchived"].(bool); ok {
        data.IsArchived = types.BoolValue(val)
    }
    if val, ok := dataMap["labels"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.Labels = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.Labels = types.SetValueMust(types.StringType, []attr.Value{})
    }
    data.RetainTelemetryDataForDays = NumberFromAPI(data.RetainTelemetryDataForDays, dataMap["retainTelemetryDataForDays"])
    if obj, ok := dataMap["telemetryRetentionConfig"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TelemetryRetentionConfig = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.TelemetryRetentionConfig = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.TelemetryRetentionConfig = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.TelemetryRetentionConfig = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.TelemetryRetentionConfig = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.TelemetryRetentionConfig = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.TelemetryRetentionConfig = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.TelemetryRetentionConfig = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.TelemetryRetentionConfig = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["telemetryRetentionConfig"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.TelemetryRetentionConfig = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.TelemetryRetentionConfig = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["telemetryRetentionConfig"].(string); ok {
        data.TelemetryRetentionConfig = NewJSONSubsetValue(val)
    } else {
        data.TelemetryRetentionConfig = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.CreatedAt = NewRFC3339Value(val)
        } else {
            data.CreatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["createdAt"].(string); ok && val != "" {
        data.CreatedAt = NewRFC3339Value(val)
    } else {
        data.CreatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.UpdatedAt = NewRFC3339Value(val)
        } else {
            data.UpdatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["updatedAt"].(string); ok && val != "" {
        data.UpdatedAt = NewRFC3339Value(val)
    } else {
        data.UpdatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["slug"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Slug = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Slug = types.StringValue(string(jsonBytes))
            } else {
                data.Slug = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Slug = types.StringValue(string(jsonBytes))
            } else {
                data.Slug = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Slug = types.StringValue(string(jsonBytes))
        } else {
            data.Slug = types.StringNull()
        }
    } else if arr, ok := dataMap["slug"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Slug = types.StringValue(string(jsonBytes))
        } else {
            data.Slug = types.StringNull()
        }
    } else if val, ok := dataMap["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
    }
    if obj, ok := dataMap["storageSystem"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StorageSystem = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.StorageSystem = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.StorageSystem = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.StorageSystem = types.StringValue(string(jsonBytes))
            } else {
                data.StorageSystem = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.StorageSystem = types.StringValue(string(jsonBytes))
            } else {
                data.StorageSystem = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.StorageSystem = types.StringValue(string(jsonBytes))
        } else {
            data.StorageSystem = types.StringNull()
        }
    } else if arr, ok := dataMap["storageSystem"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.StorageSystem = types.StringValue(string(jsonBytes))
        } else {
            data.StorageSystem = types.StringNull()
        }
    } else if val, ok := dataMap["storageSystem"].(string); ok {
        data.StorageSystem = types.StringValue(val)
    } else {
        data.StorageSystem = types.StringNull()
    }
    if obj, ok := dataMap["reportedName"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ReportedName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ReportedName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ReportedName = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ReportedName = types.StringValue(string(jsonBytes))
            } else {
                data.ReportedName = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ReportedName = types.StringValue(string(jsonBytes))
            } else {
                data.ReportedName = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ReportedName = types.StringValue(string(jsonBytes))
        } else {
            data.ReportedName = types.StringNull()
        }
    } else if arr, ok := dataMap["reportedName"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.ReportedName = types.StringValue(string(jsonBytes))
        } else {
            data.ReportedName = types.StringNull()
        }
    } else if val, ok := dataMap["reportedName"].(string); ok {
        data.ReportedName = types.StringValue(val)
    } else {
        data.ReportedName = types.StringNull()
    }
    if obj, ok := dataMap["systemId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SystemId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SystemId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SystemId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SystemId = types.StringValue(string(jsonBytes))
            } else {
                data.SystemId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SystemId = types.StringValue(string(jsonBytes))
            } else {
                data.SystemId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SystemId = types.StringValue(string(jsonBytes))
        } else {
            data.SystemId = types.StringNull()
        }
    } else if arr, ok := dataMap["systemId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.SystemId = types.StringValue(string(jsonBytes))
        } else {
            data.SystemId = types.StringNull()
        }
    } else if val, ok := dataMap["systemId"].(string); ok {
        data.SystemId = types.StringValue(val)
    } else {
        data.SystemId = types.StringNull()
    }
    if obj, ok := dataMap["osName"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OsName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.OsName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.OsName = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.OsName = types.StringValue(string(jsonBytes))
            } else {
                data.OsName = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.OsName = types.StringValue(string(jsonBytes))
            } else {
                data.OsName = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.OsName = types.StringValue(string(jsonBytes))
        } else {
            data.OsName = types.StringNull()
        }
    } else if arr, ok := dataMap["osName"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.OsName = types.StringValue(string(jsonBytes))
        } else {
            data.OsName = types.StringNull()
        }
    } else if val, ok := dataMap["osName"].(string); ok {
        data.OsName = types.StringValue(val)
    } else {
        data.OsName = types.StringNull()
    }
    if obj, ok := dataMap["osVersion"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OsVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.OsVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.OsVersion = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.OsVersion = types.StringValue(string(jsonBytes))
            } else {
                data.OsVersion = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.OsVersion = types.StringValue(string(jsonBytes))
            } else {
                data.OsVersion = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.OsVersion = types.StringValue(string(jsonBytes))
        } else {
            data.OsVersion = types.StringNull()
        }
    } else if arr, ok := dataMap["osVersion"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.OsVersion = types.StringValue(string(jsonBytes))
        } else {
            data.OsVersion = types.StringNull()
        }
    } else if val, ok := dataMap["osVersion"].(string); ok {
        data.OsVersion = types.StringValue(val)
    } else {
        data.OsVersion = types.StringNull()
    }
    if obj, ok := dataMap["otelCollectorStatus"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OtelCollectorStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.OtelCollectorStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.OtelCollectorStatus = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.OtelCollectorStatus = types.StringValue(string(jsonBytes))
            } else {
                data.OtelCollectorStatus = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.OtelCollectorStatus = types.StringValue(string(jsonBytes))
            } else {
                data.OtelCollectorStatus = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.OtelCollectorStatus = types.StringValue(string(jsonBytes))
        } else {
            data.OtelCollectorStatus = types.StringNull()
        }
    } else if arr, ok := dataMap["otelCollectorStatus"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.OtelCollectorStatus = types.StringValue(string(jsonBytes))
        } else {
            data.OtelCollectorStatus = types.StringNull()
        }
    } else if val, ok := dataMap["otelCollectorStatus"].(string); ok {
        data.OtelCollectorStatus = types.StringValue(val)
    } else {
        data.OtelCollectorStatus = types.StringNull()
    }
    if obj, ok := dataMap["agentVersion"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AgentVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AgentVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AgentVersion = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AgentVersion = types.StringValue(string(jsonBytes))
            } else {
                data.AgentVersion = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AgentVersion = types.StringValue(string(jsonBytes))
            } else {
                data.AgentVersion = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AgentVersion = types.StringValue(string(jsonBytes))
        } else {
            data.AgentVersion = types.StringNull()
        }
    } else if arr, ok := dataMap["agentVersion"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.AgentVersion = types.StringValue(string(jsonBytes))
        } else {
            data.AgentVersion = types.StringNull()
        }
    } else if val, ok := dataMap["agentVersion"].(string); ok {
        data.AgentVersion = types.StringValue(val)
    } else {
        data.AgentVersion = types.StringNull()
    }
    if obj, ok := dataMap["lastSeenAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.LastSeenAt = NewRFC3339Value(val)
        } else {
            data.LastSeenAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["lastSeenAt"].(string); ok && val != "" {
        data.LastSeenAt = NewRFC3339Value(val)
    } else {
        data.LastSeenAt = NewRFC3339Null()
    }
    data.CapacityBytes = NumberFromAPI(data.CapacityBytes, dataMap["capacityBytes"])
    data.UsedBytes = NumberFromAPI(data.UsedBytes, dataMap["usedBytes"])
    data.CapacityUsedPercent = NumberFromAPI(data.CapacityUsedPercent, dataMap["capacityUsedPercent"])
    data.DataReductionRatio = NumberFromAPI(data.DataReductionRatio, dataMap["dataReductionRatio"])
    data.OpenAlertCount = NumberFromAPI(data.OpenAlertCount, dataMap["openAlertCount"])
    data.CriticalAlertCount = NumberFromAPI(data.CriticalAlertCount, dataMap["criticalAlertCount"])
    data.WarningAlertCount = NumberFromAPI(data.WarningAlertCount, dataMap["warningAlertCount"])
    data.VolumeCount = NumberFromAPI(data.VolumeCount, dataMap["volumeCount"])
    data.HostCount = NumberFromAPI(data.HostCount, dataMap["hostCount"])
    data.PodCount = NumberFromAPI(data.PodCount, dataMap["podCount"])
    data.FileSystemCount = NumberFromAPI(data.FileSystemCount, dataMap["fileSystemCount"])
    data.BucketCount = NumberFromAPI(data.BucketCount, dataMap["bucketCount"])
    data.HardwareComponentCount = NumberFromAPI(data.HardwareComponentCount, dataMap["hardwareComponentCount"])
    data.UnhealthyHardwareCount = NumberFromAPI(data.UnhealthyHardwareCount, dataMap["unhealthyHardwareCount"])
    data.HealthStatus = NumberFromAPI(data.HealthStatus, dataMap["healthStatus"])
    if obj, ok := dataMap["createdByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if arr, ok := dataMap["createdByUserId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["createdByUserId"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if obj, ok := dataMap["archivedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.ArchivedAt = NewRFC3339Value(val)
        } else {
            data.ArchivedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["archivedAt"].(string); ok && val != "" {
        data.ArchivedAt = NewRFC3339Value(val)
    } else {
        data.ArchivedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["archivedByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ArchivedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ArchivedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ArchivedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ArchivedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.ArchivedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ArchivedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.ArchivedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ArchivedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.ArchivedByUserId = types.StringNull()
        }
    } else if arr, ok := dataMap["archivedByUserId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.ArchivedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.ArchivedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["archivedByUserId"].(string); ok {
        data.ArchivedByUserId = types.StringValue(val)
    } else {
        data.ArchivedByUserId = types.StringNull()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // A list OneUptime adds to as well holds the entries the configuration declares (see entries.go).
    r.heldEntriesAfterRead(ctx, req.Private, &data, &prior)

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *StorageArrayResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var data StorageArrayResourceModel
    var state StorageArrayResourceModel

    // Read Terraform current state data to get the ID
    resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Read Terraform plan data to get the new values
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Use the ID from the current state
    data.Id = state.Id

    // The attributes this update writes from the configuration (see importsettle.go).
    writtenFromConfig := map[string]bool{
        "telemetry_retention_config": !data.TelemetryRetentionConfig.Equal(state.TelemetryRetentionConfig),
    }

    // What the configuration sets, and what Terraform planned (see keepPlannedValues).
    var config StorageArrayResourceModel
    resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
    if resp.Diagnostics.HasError() {
        return
    }
    plan := data

    // Create API request body
    storageArrayRequest := map[string]interface{}{
        "data": map[string]interface{}{},
    }
    requestDataMap := storageArrayRequest["data"].(map[string]interface{})
    // The entries this update adds to a list OneUptime adds to as well, and takes off it (see entries.go).
    listChanges := NewEntryChanges()

    if !data.Name.IsUnknown() && !state.Name.IsUnknown() && !data.Name.Equal(state.Name) {
        requestDataMap["name"] = data.Name.ValueString()
    }
    if !data.Description.IsUnknown() && !state.Description.IsUnknown() && !data.Description.Equal(state.Description) {
        requestDataMap["description"] = data.Description.ValueString()
    }
    if !data.IsArchived.IsUnknown() && !state.IsArchived.IsUnknown() && !data.IsArchived.Equal(state.IsArchived) {
        requestDataMap["isArchived"] = data.IsArchived.ValueBool()
    }
    if !data.Labels.IsUnknown() && !state.Labels.IsUnknown() && !data.Labels.Equal(state.Labels) {
        if r.client.ChangesEntriesOneByOne() {
            listChanges.Diff("labels", "labels", data.Labels, state.Labels)
        } else {
            requestDataMap["labels"] = r.convertTerraformSetToInterface(data.Labels)
        }
    }
    if !data.RetainTelemetryDataForDays.IsUnknown() && !state.RetainTelemetryDataForDays.IsUnknown() && !data.RetainTelemetryDataForDays.Equal(state.RetainTelemetryDataForDays) {
        requestDataMap["retainTelemetryDataForDays"] = r.bigFloatToFloat64(data.RetainTelemetryDataForDays.ValueBigFloat())
    }
    if !data.TelemetryRetentionConfig.IsUnknown() && !state.TelemetryRetentionConfig.IsUnknown() && !data.TelemetryRetentionConfig.Equal(state.TelemetryRetentionConfig) {
        var telemetryretentionconfigData interface{}
        if err := json.Unmarshal([]byte(data.TelemetryRetentionConfig.ValueString()), &telemetryretentionconfigData); err == nil {
            requestDataMap["telemetryRetentionConfig"] = telemetryretentionconfigData
        } else {
            requestDataMap["telemetryRetentionConfig"] = data.TelemetryRetentionConfig.ValueString()
        }
    }
    listChanges.AddTo(storageArrayRequest)

    // Only call the API when there are changed fields to send. An empty
    // update body is rejected by the API; state is still refreshed below so
    // this method never writes unverified plan values into state.
    if len(storageArrayRequest["data"].(map[string]interface{})) > 0 || listChanges.Any() {
        httpResp, err := r.client.Put(ctx, "/storage-array/" + data.Id.ValueString() + "", storageArrayRequest)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update storage_array, got error: %s", err))
            return
        }

        // Parse the update response
        var storageArrayResponse map[string]interface{}
        err = r.client.ParseResponse(httpResp, &storageArrayResponse)
        if err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to update storage_array: %s", err))
            return
        }
        _ = storageArrayResponse
    }

    // After successful update, fetch the current state by calling Read with select parameter
    selectParam := map[string]interface{}{
        "projectId": true,
        "name": true,
        "description": true,
        "isArchived": true,
        "labels": true,
        "retainTelemetryDataForDays": true,
        "telemetryRetentionConfig": true,
        "createdAt": true,
        "updatedAt": true,
        "slug": true,
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
        "archivedAt": true,
        "archivedByUserId": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect(ctx, "/storage-array/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read storage_array after update, got error: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read storage_array after update: %s", err))
        return
    }

    // Update the model with response data from the Read operation
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := readResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = readResponse
    }

    if obj, ok := dataMap["projectId"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok {
            data.ProjectId = types.StringValue(val)
        } else {
            data.ProjectId = types.StringNull()
        }
    } else if val, ok := dataMap["projectId"].(string); ok {
        data.ProjectId = types.StringValue(val)
    } else {
        data.ProjectId = types.StringNull()
    }
    if obj, ok := dataMap["name"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Name = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Name = types.StringValue(string(jsonBytes))
            } else {
                data.Name = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Name = types.StringValue(string(jsonBytes))
            } else {
                data.Name = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Name = types.StringValue(string(jsonBytes))
        } else {
            data.Name = types.StringNull()
        }
    } else if arr, ok := dataMap["name"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Name = types.StringValue(string(jsonBytes))
        } else {
            data.Name = types.StringNull()
        }
    } else if val, ok := dataMap["name"].(string); ok {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if obj, ok := dataMap["description"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Description = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Description = types.StringValue(string(jsonBytes))
            } else {
                data.Description = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Description = types.StringValue(string(jsonBytes))
            } else {
                data.Description = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Description = types.StringValue(string(jsonBytes))
        } else {
            data.Description = types.StringNull()
        }
    } else if arr, ok := dataMap["description"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Description = types.StringValue(string(jsonBytes))
        } else {
            data.Description = types.StringNull()
        }
    } else if val, ok := dataMap["description"].(string); ok {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if val, ok := dataMap["isArchived"].(bool); ok {
        data.IsArchived = types.BoolValue(val)
    }
    if val, ok := dataMap["labels"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.Labels = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.Labels = types.SetValueMust(types.StringType, []attr.Value{})
    }
    data.RetainTelemetryDataForDays = NumberFromAPI(data.RetainTelemetryDataForDays, dataMap["retainTelemetryDataForDays"])
    if obj, ok := dataMap["telemetryRetentionConfig"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TelemetryRetentionConfig = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.TelemetryRetentionConfig = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.TelemetryRetentionConfig = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.TelemetryRetentionConfig = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.TelemetryRetentionConfig = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.TelemetryRetentionConfig = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.TelemetryRetentionConfig = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.TelemetryRetentionConfig = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.TelemetryRetentionConfig = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["telemetryRetentionConfig"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.TelemetryRetentionConfig = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.TelemetryRetentionConfig = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["telemetryRetentionConfig"].(string); ok {
        data.TelemetryRetentionConfig = NewJSONSubsetValue(val)
    } else {
        data.TelemetryRetentionConfig = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.CreatedAt = NewRFC3339Value(val)
        } else {
            data.CreatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["createdAt"].(string); ok && val != "" {
        data.CreatedAt = NewRFC3339Value(val)
    } else {
        data.CreatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.UpdatedAt = NewRFC3339Value(val)
        } else {
            data.UpdatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["updatedAt"].(string); ok && val != "" {
        data.UpdatedAt = NewRFC3339Value(val)
    } else {
        data.UpdatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["slug"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Slug = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Slug = types.StringValue(string(jsonBytes))
            } else {
                data.Slug = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Slug = types.StringValue(string(jsonBytes))
            } else {
                data.Slug = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Slug = types.StringValue(string(jsonBytes))
        } else {
            data.Slug = types.StringNull()
        }
    } else if arr, ok := dataMap["slug"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Slug = types.StringValue(string(jsonBytes))
        } else {
            data.Slug = types.StringNull()
        }
    } else if val, ok := dataMap["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
    }
    if obj, ok := dataMap["storageSystem"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StorageSystem = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.StorageSystem = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.StorageSystem = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.StorageSystem = types.StringValue(string(jsonBytes))
            } else {
                data.StorageSystem = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.StorageSystem = types.StringValue(string(jsonBytes))
            } else {
                data.StorageSystem = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.StorageSystem = types.StringValue(string(jsonBytes))
        } else {
            data.StorageSystem = types.StringNull()
        }
    } else if arr, ok := dataMap["storageSystem"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.StorageSystem = types.StringValue(string(jsonBytes))
        } else {
            data.StorageSystem = types.StringNull()
        }
    } else if val, ok := dataMap["storageSystem"].(string); ok {
        data.StorageSystem = types.StringValue(val)
    } else {
        data.StorageSystem = types.StringNull()
    }
    if obj, ok := dataMap["reportedName"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ReportedName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ReportedName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ReportedName = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ReportedName = types.StringValue(string(jsonBytes))
            } else {
                data.ReportedName = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ReportedName = types.StringValue(string(jsonBytes))
            } else {
                data.ReportedName = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ReportedName = types.StringValue(string(jsonBytes))
        } else {
            data.ReportedName = types.StringNull()
        }
    } else if arr, ok := dataMap["reportedName"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.ReportedName = types.StringValue(string(jsonBytes))
        } else {
            data.ReportedName = types.StringNull()
        }
    } else if val, ok := dataMap["reportedName"].(string); ok {
        data.ReportedName = types.StringValue(val)
    } else {
        data.ReportedName = types.StringNull()
    }
    if obj, ok := dataMap["systemId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SystemId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SystemId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SystemId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SystemId = types.StringValue(string(jsonBytes))
            } else {
                data.SystemId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SystemId = types.StringValue(string(jsonBytes))
            } else {
                data.SystemId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SystemId = types.StringValue(string(jsonBytes))
        } else {
            data.SystemId = types.StringNull()
        }
    } else if arr, ok := dataMap["systemId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.SystemId = types.StringValue(string(jsonBytes))
        } else {
            data.SystemId = types.StringNull()
        }
    } else if val, ok := dataMap["systemId"].(string); ok {
        data.SystemId = types.StringValue(val)
    } else {
        data.SystemId = types.StringNull()
    }
    if obj, ok := dataMap["osName"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OsName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.OsName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.OsName = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.OsName = types.StringValue(string(jsonBytes))
            } else {
                data.OsName = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.OsName = types.StringValue(string(jsonBytes))
            } else {
                data.OsName = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.OsName = types.StringValue(string(jsonBytes))
        } else {
            data.OsName = types.StringNull()
        }
    } else if arr, ok := dataMap["osName"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.OsName = types.StringValue(string(jsonBytes))
        } else {
            data.OsName = types.StringNull()
        }
    } else if val, ok := dataMap["osName"].(string); ok {
        data.OsName = types.StringValue(val)
    } else {
        data.OsName = types.StringNull()
    }
    if obj, ok := dataMap["osVersion"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OsVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.OsVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.OsVersion = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.OsVersion = types.StringValue(string(jsonBytes))
            } else {
                data.OsVersion = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.OsVersion = types.StringValue(string(jsonBytes))
            } else {
                data.OsVersion = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.OsVersion = types.StringValue(string(jsonBytes))
        } else {
            data.OsVersion = types.StringNull()
        }
    } else if arr, ok := dataMap["osVersion"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.OsVersion = types.StringValue(string(jsonBytes))
        } else {
            data.OsVersion = types.StringNull()
        }
    } else if val, ok := dataMap["osVersion"].(string); ok {
        data.OsVersion = types.StringValue(val)
    } else {
        data.OsVersion = types.StringNull()
    }
    if obj, ok := dataMap["otelCollectorStatus"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OtelCollectorStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.OtelCollectorStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.OtelCollectorStatus = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.OtelCollectorStatus = types.StringValue(string(jsonBytes))
            } else {
                data.OtelCollectorStatus = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.OtelCollectorStatus = types.StringValue(string(jsonBytes))
            } else {
                data.OtelCollectorStatus = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.OtelCollectorStatus = types.StringValue(string(jsonBytes))
        } else {
            data.OtelCollectorStatus = types.StringNull()
        }
    } else if arr, ok := dataMap["otelCollectorStatus"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.OtelCollectorStatus = types.StringValue(string(jsonBytes))
        } else {
            data.OtelCollectorStatus = types.StringNull()
        }
    } else if val, ok := dataMap["otelCollectorStatus"].(string); ok {
        data.OtelCollectorStatus = types.StringValue(val)
    } else {
        data.OtelCollectorStatus = types.StringNull()
    }
    if obj, ok := dataMap["agentVersion"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AgentVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AgentVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AgentVersion = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AgentVersion = types.StringValue(string(jsonBytes))
            } else {
                data.AgentVersion = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AgentVersion = types.StringValue(string(jsonBytes))
            } else {
                data.AgentVersion = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AgentVersion = types.StringValue(string(jsonBytes))
        } else {
            data.AgentVersion = types.StringNull()
        }
    } else if arr, ok := dataMap["agentVersion"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.AgentVersion = types.StringValue(string(jsonBytes))
        } else {
            data.AgentVersion = types.StringNull()
        }
    } else if val, ok := dataMap["agentVersion"].(string); ok {
        data.AgentVersion = types.StringValue(val)
    } else {
        data.AgentVersion = types.StringNull()
    }
    if obj, ok := dataMap["lastSeenAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.LastSeenAt = NewRFC3339Value(val)
        } else {
            data.LastSeenAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["lastSeenAt"].(string); ok && val != "" {
        data.LastSeenAt = NewRFC3339Value(val)
    } else {
        data.LastSeenAt = NewRFC3339Null()
    }
    data.CapacityBytes = NumberFromAPI(data.CapacityBytes, dataMap["capacityBytes"])
    data.UsedBytes = NumberFromAPI(data.UsedBytes, dataMap["usedBytes"])
    data.CapacityUsedPercent = NumberFromAPI(data.CapacityUsedPercent, dataMap["capacityUsedPercent"])
    data.DataReductionRatio = NumberFromAPI(data.DataReductionRatio, dataMap["dataReductionRatio"])
    data.OpenAlertCount = NumberFromAPI(data.OpenAlertCount, dataMap["openAlertCount"])
    data.CriticalAlertCount = NumberFromAPI(data.CriticalAlertCount, dataMap["criticalAlertCount"])
    data.WarningAlertCount = NumberFromAPI(data.WarningAlertCount, dataMap["warningAlertCount"])
    data.VolumeCount = NumberFromAPI(data.VolumeCount, dataMap["volumeCount"])
    data.HostCount = NumberFromAPI(data.HostCount, dataMap["hostCount"])
    data.PodCount = NumberFromAPI(data.PodCount, dataMap["podCount"])
    data.FileSystemCount = NumberFromAPI(data.FileSystemCount, dataMap["fileSystemCount"])
    data.BucketCount = NumberFromAPI(data.BucketCount, dataMap["bucketCount"])
    data.HardwareComponentCount = NumberFromAPI(data.HardwareComponentCount, dataMap["hardwareComponentCount"])
    data.UnhealthyHardwareCount = NumberFromAPI(data.UnhealthyHardwareCount, dataMap["unhealthyHardwareCount"])
    data.HealthStatus = NumberFromAPI(data.HealthStatus, dataMap["healthStatus"])
    if obj, ok := dataMap["createdByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if arr, ok := dataMap["createdByUserId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["createdByUserId"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if obj, ok := dataMap["archivedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.ArchivedAt = NewRFC3339Value(val)
        } else {
            data.ArchivedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["archivedAt"].(string); ok && val != "" {
        data.ArchivedAt = NewRFC3339Value(val)
    } else {
        data.ArchivedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["archivedByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ArchivedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ArchivedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ArchivedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ArchivedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.ArchivedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ArchivedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.ArchivedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ArchivedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.ArchivedByUserId = types.StringNull()
        }
    } else if arr, ok := dataMap["archivedByUserId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.ArchivedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.ArchivedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["archivedByUserId"].(string); ok {
        data.ArchivedByUserId = types.StringValue(val)
    } else {
        data.ArchivedByUserId = types.StringNull()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    data.Id = state.Id

    // Unconfigured attributes keep their planned value; see keepPlannedValues.
    r.keepPlannedValues(&data, &plan, &config)

    // What this update added to a list OneUptime adds to as well is there, and what it took off is gone (see entries.go).
    resp.Diagnostics.Append(listChanges.Check("labels", data.Labels)...)
    if resp.Diagnostics.HasError() {
        return
    }

    // A list OneUptime adds to as well holds the entries the configuration declares (see entries.go).
    resp.Diagnostics.Append(r.heldEntriesAfterWrite(ctx, req.Private, &data, &plan, &config, false).Write(ctx, resp.Private)...)

    // What this update wrote is what the configuration says now (see importsettle.go).
    resp.Diagnostics.Append(KeepImportedFromServer(ctx, resp.Private, writtenFromConfig)...)

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *StorageArrayResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var data StorageArrayResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Make API call
    httpResp, err := r.client.Delete(ctx, "/storage-array/" + data.Id.ValueString() + "")
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete storage_array, got error: %s", err))
        return
    }

    // A failed delete must keep the resource in state — silently dropping it
    // orphans real infrastructure. 404 means it is already gone.
    if httpResp.StatusCode >= 400 && httpResp.StatusCode != http.StatusNotFound {
        err = r.client.ParseResponse(httpResp, nil)
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to delete storage_array: %s", err))
        return
    }
    if httpResp.Body != nil {
        httpResp.Body.Close()
    }
}


func (r *StorageArrayResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
    // These hold the record as OneUptime keeps it, keys it fills in itself
    // included, until an update writes them (see importsettle.go).
    resp.Diagnostics.Append(MarkImportedFromServer(ctx, resp.Private, []string{"telemetry_retention_config"})...)

    // An import adopts the record's entries of each list OneUptime adds to as well (see entries.go).
    if r.client.HoldsDeclaredEntries() {
        resp.Diagnostics.Append(DeclaredEntries{"labels": true}.Write(ctx, resp.Private)...)
    }
}


// keepPlannedValues puts back, after a create or an update, the planned value
// of each optional attribute the configuration leaves out. The server keeps
// some of these up to date on its own (when it last checked a heartbeat, the
// status a probe last reported...), so the value read back after the write can
// already differ from the plan, and Terraform would fail the apply with
// "Provider produced inconsistent result after apply". The next refresh reads
// the server's value, which is never a diff for an attribute nobody configured.
func (r *StorageArrayResource) keepPlannedValues(data *StorageArrayResourceModel, plan *StorageArrayResourceModel, config *StorageArrayResourceModel) {
    if config.Description.IsNull() && !plan.Description.IsUnknown() {
        data.Description = plan.Description
    }
    if config.IsArchived.IsNull() && !plan.IsArchived.IsUnknown() {
        data.IsArchived = plan.IsArchived
    }
    if config.Labels.IsNull() && !plan.Labels.IsUnknown() {
        data.Labels = plan.Labels
    }
    if config.RetainTelemetryDataForDays.IsNull() && !plan.RetainTelemetryDataForDays.IsUnknown() {
        data.RetainTelemetryDataForDays = plan.RetainTelemetryDataForDays
    }
    if config.TelemetryRetentionConfig.IsNull() && !plan.TelemetryRetentionConfig.IsUnknown() {
        data.TelemetryRetentionConfig = plan.TelemetryRetentionConfig
    }
}

// heldEntriesAfterRead holds each list OneUptime adds to as well (labels) to
// the entries the configuration declares that the record still has, after
// a read (see entries.go).
func (r *StorageArrayResource) heldEntriesAfterRead(ctx context.Context, private PrivateState, data *StorageArrayResourceModel, prior *StorageArrayResourceModel) {
    declared := ReadDeclaredEntries(ctx, private)
    data.Labels = SetEntriesAfterRead(r.client, declared, "labels", prior.Labels, data.Labels)
}

// heldEntriesAfterWrite is the state of each such list after a create or an
// update, and which of them the configuration declares from then on (see
// entries.go): the caller keeps them in the private state.
func (r *StorageArrayResource) heldEntriesAfterWrite(ctx context.Context, before PrivateState, data *StorageArrayResourceModel, plan *StorageArrayResourceModel, config *StorageArrayResourceModel, isCreate bool) DeclaredEntries {
    declared := DeclaredEntries{}
    if !isCreate {
        declared = ReadDeclaredEntries(ctx, before)
    }
    var isDeclared bool
    data.Labels, isDeclared = SetEntriesAfterWrite(r.client, isCreate, declared["labels"], config.Labels, plan.Labels, data.Labels)
    declared["labels"] = isDeclared
    return declared
}

// Helper method to convert Terraform map to Go interface{}
func (r *StorageArrayResource) convertTerraformMapToInterface(terraformMap types.Map) interface{} {
    if terraformMap.IsNull() || terraformMap.IsUnknown() {
        return nil
    }
    
    result := make(map[string]string)
    terraformMap.ElementsAs(context.Background(), &result, false)
    
    // Convert map[string]string to map[string]interface{}
    interfaceResult := make(map[string]interface{})
    for key, value := range result {
        interfaceResult[key] = value
    }
    
    return interfaceResult
}

// Helper method to convert Terraform list to Go interface{}
func (r *StorageArrayResource) convertTerraformListToInterface(terraformList types.List) interface{} {
    if terraformList.IsNull() || terraformList.IsUnknown() {
        return nil
    }

    var stringList []string
    terraformList.ElementsAs(context.Background(), &stringList, false)

    // Convert string array to OneUptime format with _id fields
    var result []interface{}
    for _, str := range stringList {
        if str != "" {
            result = append(result, map[string]interface{}{
                "_id": str,
            })
        }
    }
    return result
}

// Helper method to convert Terraform set to Go interface{}
func (r *StorageArrayResource) convertTerraformSetToInterface(terraformSet types.Set) interface{} {
    if terraformSet.IsNull() || terraformSet.IsUnknown() {
        return nil
    }

    var stringList []string
    terraformSet.ElementsAs(context.Background(), &stringList, false)

    // Convert string array to OneUptime format with _id fields
    var result []interface{}
    for _, str := range stringList {
        if str != "" {
            result = append(result, map[string]interface{}{
                "_id": str,
            })
        }
    }
    return result
}


// Helper method to parse JSON field for complex objects
func (r *StorageArrayResource) parseJSONField(terraformString basetypes.StringValuable) interface{} {
    sv, _ := terraformString.ToStringValue(context.Background())
    if sv.IsNull() || sv.IsUnknown() || sv.ValueString() == "" {
        return nil
    }

    var result interface{}
    if err := json.Unmarshal([]byte(sv.ValueString()), &result); err != nil {
        // If JSON parsing fails, return the raw string
        return sv.ValueString()
    }

    return result
}

// Normalize URL wrapper objects to avoid drift (e.g., trailing slash differences).
func (r *StorageArrayResource) normalizeURLWrappers(value interface{}) interface{} {
    switch v := value.(type) {
    case map[string]interface{}:
        if typeStr, ok := v["_type"].(string); ok && typeStr == "URL" {
            if val, ok := v["value"].(string); ok {
                v["value"] = r.normalizeURLString(val)
            }
        }
        for key, child := range v {
            v[key] = r.normalizeURLWrappers(child)
        }
        return v
    case []interface{}:
        for i, child := range v {
            v[i] = r.normalizeURLWrappers(child)
        }
        return v
    default:
        return v
    }
}

func (r *StorageArrayResource) normalizeURLString(value string) string {
    parsed, err := url.Parse(value)
    if err != nil {
        return value
    }
    if parsed.Path == "/" && parsed.RawQuery == "" && parsed.Fragment == "" {
        return strings.TrimSuffix(value, "/")
    }
    return value
}

// Helper method to convert *big.Float to float64 for JSON serialization
func (r *StorageArrayResource) bigFloatToFloat64(bf *big.Float) interface{} {
    if bf == nil {
        return nil
    }
    f, _ := bf.Float64()
    return f
}

// Helper method to check if a type string is a valid OneUptime ObjectType.
// The registry itself lives in objecttypes.go, shared across the package.
func (r *StorageArrayResource) isValidOneUptimeObjectType(typeStr string) bool {
    return ValidOneUptimeObjectTypes[typeStr]
}
