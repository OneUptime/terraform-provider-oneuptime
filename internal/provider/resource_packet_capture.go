package provider

import (
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
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/numberdefault"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/numberplanmodifier"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &PacketCaptureResource{}
var _ resource.ResourceWithImportState = &PacketCaptureResource{}

func NewPacketCaptureResource() resource.Resource {
    return &PacketCaptureResource{}
}

// PacketCaptureResource defines the resource implementation.
type PacketCaptureResource struct {
    client *Client
}

// PacketCaptureResourceModel describes the resource data model.
type PacketCaptureResourceModel struct {
    Id types.String `tfsdk:"id"`
    ProjectId types.String `tfsdk:"project_id"`
    ProbeId types.String `tfsdk:"probe_id"`
    NetworkDeviceId types.String `tfsdk:"network_device_id"`
    InterfaceName types.String `tfsdk:"interface_name"`
    BpfFilter types.String `tfsdk:"bpf_filter"`
    MaxDurationInSeconds types.Number `tfsdk:"max_duration_in_seconds"`
    MaxPackets types.Number `tfsdk:"max_packets"`
    MaxFileSizeInMb types.Number `tfsdk:"max_file_size_in_mb"`
    CreatedAt RFC3339Value `tfsdk:"created_at"`
    UpdatedAt RFC3339Value `tfsdk:"updated_at"`
    Name types.String `tfsdk:"name"`
    Status types.String `tfsdk:"status"`
    StatusMessage types.String `tfsdk:"status_message"`
    EndReason types.String `tfsdk:"end_reason"`
    StartedAt RFC3339Value `tfsdk:"started_at"`
    CompletedAt RFC3339Value `tfsdk:"completed_at"`
    StopRequestedAt RFC3339Value `tfsdk:"stop_requested_at"`
    PacketCount types.Number `tfsdk:"packet_count"`
    FileSizeInBytes types.Number `tfsdk:"file_size_in_bytes"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (r *PacketCaptureResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_packet_capture"
}

func (r *PacketCaptureResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = r.schemaDefinition()
}

func (r *PacketCaptureResource) schemaDefinition() schema.Schema {
    return schema.Schema{
        MarkdownDescription: "A packet capture run on one of the project's probes, and the pcap file it produced. Captures and their files are deleted 7 days after they start.",

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
            "probe_id": schema.StringAttribute{
                MarkdownDescription: "ID of the probe that runs the capture: one of the project's own probes, with packet capture turned on. The ID of a `oneuptime_probe`.",
                Required: true,
            },
            "network_device_id": schema.StringAttribute{
                MarkdownDescription: "ID of the network device the capture was started from, if it was started from a device's page. The ID of a `oneuptime_network_device`.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "interface_name": schema.StringAttribute{
                MarkdownDescription: "The network interface the probe captures on, as the probe reported it: \"eth0\", or \"any\" for every interface at once.",
                Required: true,
            },
            "bpf_filter": schema.StringAttribute{
                MarkdownDescription: "A BPF filter expression - the capture-filter language of tcpdump and Wireshark - that decides which packets are kept: host 10.0.0.5 and tcp port 443. Empty keeps every packet.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "max_duration_in_seconds": schema.NumberAttribute{
                MarkdownDescription: "How long the capture runs, in seconds, at most: from 5 seconds to 30 minutes, and no longer than the probe allows. It stops earlier when it reaches its packet or file size limit.",
                Optional: true,
                Computed: true,
                Default: numberdefault.StaticBigFloat(big.NewFloat(60)),
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "max_packets": schema.NumberAttribute{
                MarkdownDescription: "The capture stops after this many packets: from 1 to 1,000,000.",
                Optional: true,
                Computed: true,
                Default: numberdefault.StaticBigFloat(big.NewFloat(100000)),
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "max_file_size_in_mb": schema.NumberAttribute{
                MarkdownDescription: "The capture stops when its file reaches this many megabytes: from 1 to 25, and no more than the probe allows. The file is cut at the last whole packet.",
                Optional: true,
                Computed: true,
                Default: numberdefault.StaticBigFloat(big.NewFloat(10)),
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
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
            "name": schema.StringAttribute{
                MarkdownDescription: "The interface and the filter of the capture, written by OneUptime: eth0: host 10.0.0.5.",
                Computed: true,
            },
            "status": schema.StringAttribute{
                MarkdownDescription: "Where the capture is: \"Pending\" (waiting for the probe), \"Running\" (the probe is capturing), \"Completed\" (the file is ready) or \"Failed\" (see Status Message). Managed by OneUptime and the probe.",
                Computed: true,
            },
            "status_message": schema.StringAttribute{
                MarkdownDescription: "Why a capture failed - the interface is gone, the probe may not capture, the filter did not compile - or what the capture tool said when it stopped by itself. Managed by OneUptime and the probe.",
                Computed: true,
            },
            "end_reason": schema.StringAttribute{
                MarkdownDescription: "Why a completed capture stopped: \"DurationReached\", \"PacketLimitReached\", \"FileSizeLimitReached\", \"StoppedFromDashboard\" or \"CaptureToolStopped\". Managed by the probe.",
                Computed: true,
            },
            "started_at": schema.StringAttribute{
                MarkdownDescription: "When the probe started capturing. Managed by OneUptime.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "completed_at": schema.StringAttribute{
                MarkdownDescription: "When the capture completed or failed. Managed by OneUptime.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "stop_requested_at": schema.StringAttribute{
                MarkdownDescription: "When someone pressed Stop on the running capture. The probe stops at its next check, within about ten seconds, and uploads what it captured. Managed by OneUptime.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "packet_count": schema.NumberAttribute{
                MarkdownDescription: "How many packets the file holds, counted by OneUptime from the file. Managed by OneUptime.",
                Computed: true,
            },
            "file_size_in_bytes": schema.NumberAttribute{
                MarkdownDescription: "The size of the pcap file in bytes. Managed by OneUptime.",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
        },
    }
}

func (r *PacketCaptureResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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


func (r *PacketCaptureResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var data PacketCaptureResourceModel

    // Read Terraform plan data into the model
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // What the configuration sets, and what Terraform planned (see keepPlannedValues).
    var config PacketCaptureResourceModel
    resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
    if resp.Diagnostics.HasError() {
        return
    }
    plan := data



    // Create API request body. Unset (null/unknown) optional fields are
    // omitted so server-side defaults apply instead of being overwritten
    // with zero values.
    packetCaptureRequest := map[string]interface{}{
        "data": map[string]interface{}{},
    }
    requestDataMap := packetCaptureRequest["data"].(map[string]interface{})

    if !data.ProbeId.IsNull() && !data.ProbeId.IsUnknown() {
        requestDataMap["probeId"] = data.ProbeId.ValueString()
    }
    if !data.NetworkDeviceId.IsNull() && !data.NetworkDeviceId.IsUnknown() {
        requestDataMap["networkDeviceId"] = data.NetworkDeviceId.ValueString()
    }
    if !data.InterfaceName.IsNull() && !data.InterfaceName.IsUnknown() {
        requestDataMap["interfaceName"] = data.InterfaceName.ValueString()
    }
    if !data.BpfFilter.IsNull() && !data.BpfFilter.IsUnknown() {
        requestDataMap["bpfFilter"] = data.BpfFilter.ValueString()
    }
    if !data.MaxDurationInSeconds.IsNull() && !data.MaxDurationInSeconds.IsUnknown() {
        requestDataMap["maxDurationInSeconds"] = r.bigFloatToFloat64(data.MaxDurationInSeconds.ValueBigFloat())
    }
    if !data.MaxPackets.IsNull() && !data.MaxPackets.IsUnknown() {
        requestDataMap["maxPackets"] = r.bigFloatToFloat64(data.MaxPackets.ValueBigFloat())
    }
    if !data.MaxFileSizeInMb.IsNull() && !data.MaxFileSizeInMb.IsUnknown() {
        requestDataMap["maxFileSizeInMB"] = r.bigFloatToFloat64(data.MaxFileSizeInMb.ValueBigFloat())
    }

    // Make API call
    httpResp, err := r.client.Post(ctx, "/packet-capture", packetCaptureRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create packet_capture, got error: %s", err))
        return
    }

    var packetCaptureResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &packetCaptureResponse)
    if err != nil {
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to create packet_capture: %s", err))
        return
    }

    // Extract the new resource id from the create response.
    createdId := ""
    if wrapper, ok := packetCaptureResponse["data"].(map[string]interface{}); ok {
        if val, ok := wrapper["_id"].(string); ok {
            createdId = val
        }
    } else if val, ok := packetCaptureResponse["_id"].(string); ok {
        createdId = val
    }
    if createdId == "" {
        resp.Diagnostics.AddError("OneUptime API Error", "Create response for packet_capture did not contain an id. This is a bug in the provider or the API; please report it.")
        return
    }
    data.Id = types.StringValue(createdId)

    /*
     * The server has committed the row. Persist what we know to state BEFORE
     * the read-back: if the read-back fails and we return without setting
     * state, Terraform never learns the resource exists and the created
     * packet_capture is orphaned server-side — never refreshed, never
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
        "probeId": true,
        "networkDeviceId": true,
        "interfaceName": true,
        "bpfFilter": true,
        "maxDurationInSeconds": true,
        "maxPackets": true,
        "maxFileSizeInMB": true,
        "createdAt": true,
        "updatedAt": true,
        "name": true,
        "status": true,
        "statusMessage": true,
        "endReason": true,
        "startedAt": true,
        "completedAt": true,
        "stopRequestedAt": true,
        "packetCount": true,
        "fileSizeInBytes": true,
        "createdByUserId": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect(ctx, "/packet-capture/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        /*
         * State already owns the id, so the resource is tracked and the next
         * refresh reconciles the remaining attributes. Warn rather than
         * error: erroring here would strand a real resource.
         */
        resp.Diagnostics.AddWarning("Read After Create Failed", fmt.Sprintf("Created packet_capture but could not read it back; state is incomplete until the next refresh: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddWarning("Read After Create Failed", fmt.Sprintf("Created packet_capture but could not parse the read-back response; state is incomplete until the next refresh: %s", err))
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
    if obj, ok := dataMap["probeId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ProbeId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ProbeId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ProbeId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ProbeId = types.StringValue(string(jsonBytes))
            } else {
                data.ProbeId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ProbeId = types.StringValue(string(jsonBytes))
            } else {
                data.ProbeId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ProbeId = types.StringValue(string(jsonBytes))
        } else {
            data.ProbeId = types.StringNull()
        }
    } else if val, ok := dataMap["probeId"].(string); ok {
        data.ProbeId = types.StringValue(val)
    } else {
        data.ProbeId = types.StringNull()
    }
    if obj, ok := dataMap["networkDeviceId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.NetworkDeviceId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.NetworkDeviceId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.NetworkDeviceId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.NetworkDeviceId = types.StringValue(string(jsonBytes))
            } else {
                data.NetworkDeviceId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.NetworkDeviceId = types.StringValue(string(jsonBytes))
            } else {
                data.NetworkDeviceId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.NetworkDeviceId = types.StringValue(string(jsonBytes))
        } else {
            data.NetworkDeviceId = types.StringNull()
        }
    } else if val, ok := dataMap["networkDeviceId"].(string); ok {
        data.NetworkDeviceId = types.StringValue(val)
    } else {
        data.NetworkDeviceId = types.StringNull()
    }
    if obj, ok := dataMap["interfaceName"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.InterfaceName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.InterfaceName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.InterfaceName = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.InterfaceName = types.StringValue(string(jsonBytes))
            } else {
                data.InterfaceName = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.InterfaceName = types.StringValue(string(jsonBytes))
            } else {
                data.InterfaceName = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.InterfaceName = types.StringValue(string(jsonBytes))
        } else {
            data.InterfaceName = types.StringNull()
        }
    } else if val, ok := dataMap["interfaceName"].(string); ok {
        data.InterfaceName = types.StringValue(val)
    } else {
        data.InterfaceName = types.StringNull()
    }
    if obj, ok := dataMap["bpfFilter"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.BpfFilter = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.BpfFilter = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.BpfFilter = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.BpfFilter = types.StringValue(string(jsonBytes))
            } else {
                data.BpfFilter = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.BpfFilter = types.StringValue(string(jsonBytes))
            } else {
                data.BpfFilter = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.BpfFilter = types.StringValue(string(jsonBytes))
        } else {
            data.BpfFilter = types.StringNull()
        }
    } else if val, ok := dataMap["bpfFilter"].(string); ok {
        data.BpfFilter = types.StringValue(val)
    } else {
        data.BpfFilter = types.StringNull()
    }
    if val, ok := dataMap["maxDurationInSeconds"].(float64); ok {
        data.MaxDurationInSeconds = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["maxDurationInSeconds"].(int); ok {
        data.MaxDurationInSeconds = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["maxDurationInSeconds"].(int64); ok {
        data.MaxDurationInSeconds = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["maxDurationInSeconds"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.MaxDurationInSeconds = types.NumberValue(big.NewFloat(val))
        } else {
            data.MaxDurationInSeconds = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.MaxDurationInSeconds = types.NumberNull()
    }
    if val, ok := dataMap["maxPackets"].(float64); ok {
        data.MaxPackets = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["maxPackets"].(int); ok {
        data.MaxPackets = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["maxPackets"].(int64); ok {
        data.MaxPackets = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["maxPackets"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.MaxPackets = types.NumberValue(big.NewFloat(val))
        } else {
            data.MaxPackets = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.MaxPackets = types.NumberNull()
    }
    if val, ok := dataMap["maxFileSizeInMB"].(float64); ok {
        data.MaxFileSizeInMb = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["maxFileSizeInMB"].(int); ok {
        data.MaxFileSizeInMb = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["maxFileSizeInMB"].(int64); ok {
        data.MaxFileSizeInMb = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["maxFileSizeInMB"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.MaxFileSizeInMb = types.NumberValue(big.NewFloat(val))
        } else {
            data.MaxFileSizeInMb = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.MaxFileSizeInMb = types.NumberNull()
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
    } else if val, ok := dataMap["name"].(string); ok {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if obj, ok := dataMap["status"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Status = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Status = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Status = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Status = types.StringValue(string(jsonBytes))
            } else {
                data.Status = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Status = types.StringValue(string(jsonBytes))
            } else {
                data.Status = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Status = types.StringValue(string(jsonBytes))
        } else {
            data.Status = types.StringNull()
        }
    } else if val, ok := dataMap["status"].(string); ok {
        data.Status = types.StringValue(val)
    } else {
        data.Status = types.StringNull()
    }
    if obj, ok := dataMap["statusMessage"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StatusMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.StatusMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.StatusMessage = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.StatusMessage = types.StringValue(string(jsonBytes))
            } else {
                data.StatusMessage = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.StatusMessage = types.StringValue(string(jsonBytes))
            } else {
                data.StatusMessage = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.StatusMessage = types.StringValue(string(jsonBytes))
        } else {
            data.StatusMessage = types.StringNull()
        }
    } else if val, ok := dataMap["statusMessage"].(string); ok {
        data.StatusMessage = types.StringValue(val)
    } else {
        data.StatusMessage = types.StringNull()
    }
    if obj, ok := dataMap["endReason"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EndReason = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.EndReason = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.EndReason = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.EndReason = types.StringValue(string(jsonBytes))
            } else {
                data.EndReason = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.EndReason = types.StringValue(string(jsonBytes))
            } else {
                data.EndReason = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.EndReason = types.StringValue(string(jsonBytes))
        } else {
            data.EndReason = types.StringNull()
        }
    } else if val, ok := dataMap["endReason"].(string); ok {
        data.EndReason = types.StringValue(val)
    } else {
        data.EndReason = types.StringNull()
    }
    if obj, ok := dataMap["startedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.StartedAt = NewRFC3339Value(val)
        } else {
            data.StartedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["startedAt"].(string); ok && val != "" {
        data.StartedAt = NewRFC3339Value(val)
    } else {
        data.StartedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["completedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.CompletedAt = NewRFC3339Value(val)
        } else {
            data.CompletedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["completedAt"].(string); ok && val != "" {
        data.CompletedAt = NewRFC3339Value(val)
    } else {
        data.CompletedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["stopRequestedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.StopRequestedAt = NewRFC3339Value(val)
        } else {
            data.StopRequestedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["stopRequestedAt"].(string); ok && val != "" {
        data.StopRequestedAt = NewRFC3339Value(val)
    } else {
        data.StopRequestedAt = NewRFC3339Null()
    }
    if val, ok := dataMap["packetCount"].(float64); ok {
        data.PacketCount = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["packetCount"].(int); ok {
        data.PacketCount = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["packetCount"].(int64); ok {
        data.PacketCount = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["packetCount"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.PacketCount = types.NumberValue(big.NewFloat(val))
        } else {
            data.PacketCount = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.PacketCount = types.NumberNull()
    }
    if val, ok := dataMap["fileSizeInBytes"].(float64); ok {
        data.FileSizeInBytes = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["fileSizeInBytes"].(int); ok {
        data.FileSizeInBytes = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["fileSizeInBytes"].(int64); ok {
        data.FileSizeInBytes = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["fileSizeInBytes"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.FileSizeInBytes = types.NumberValue(big.NewFloat(val))
        } else {
            data.FileSizeInBytes = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.FileSizeInBytes = types.NumberNull()
    }
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
    } else if val, ok := dataMap["createdByUserId"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
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

    // Write logs using the tflog package
    tflog.Trace(ctx, "created a resource")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PacketCaptureResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var data PacketCaptureResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create select parameter to get full object
    selectParam := map[string]interface{}{
        "projectId": true,
        "probeId": true,
        "networkDeviceId": true,
        "interfaceName": true,
        "bpfFilter": true,
        "maxDurationInSeconds": true,
        "maxPackets": true,
        "maxFileSizeInMB": true,
        "createdAt": true,
        "updatedAt": true,
        "name": true,
        "status": true,
        "statusMessage": true,
        "endReason": true,
        "startedAt": true,
        "completedAt": true,
        "stopRequestedAt": true,
        "packetCount": true,
        "fileSizeInBytes": true,
        "createdByUserId": true,
        "_id": true,
    }

    // Make API call with select parameter
    httpResp, err := r.client.PostWithSelect(ctx, "/packet-capture/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read packet_capture, got error: %s", err))
        return
    }

    if httpResp.StatusCode == http.StatusNotFound {
        resp.State.RemoveResource(ctx)
        return
    }

    var packetCaptureResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &packetCaptureResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse packet_capture response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := packetCaptureResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = packetCaptureResponse
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
    if obj, ok := dataMap["probeId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ProbeId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ProbeId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ProbeId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ProbeId = types.StringValue(string(jsonBytes))
            } else {
                data.ProbeId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ProbeId = types.StringValue(string(jsonBytes))
            } else {
                data.ProbeId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ProbeId = types.StringValue(string(jsonBytes))
        } else {
            data.ProbeId = types.StringNull()
        }
    } else if val, ok := dataMap["probeId"].(string); ok {
        data.ProbeId = types.StringValue(val)
    } else {
        data.ProbeId = types.StringNull()
    }
    if obj, ok := dataMap["networkDeviceId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.NetworkDeviceId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.NetworkDeviceId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.NetworkDeviceId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.NetworkDeviceId = types.StringValue(string(jsonBytes))
            } else {
                data.NetworkDeviceId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.NetworkDeviceId = types.StringValue(string(jsonBytes))
            } else {
                data.NetworkDeviceId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.NetworkDeviceId = types.StringValue(string(jsonBytes))
        } else {
            data.NetworkDeviceId = types.StringNull()
        }
    } else if val, ok := dataMap["networkDeviceId"].(string); ok {
        data.NetworkDeviceId = types.StringValue(val)
    } else {
        data.NetworkDeviceId = types.StringNull()
    }
    if obj, ok := dataMap["interfaceName"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.InterfaceName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.InterfaceName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.InterfaceName = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.InterfaceName = types.StringValue(string(jsonBytes))
            } else {
                data.InterfaceName = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.InterfaceName = types.StringValue(string(jsonBytes))
            } else {
                data.InterfaceName = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.InterfaceName = types.StringValue(string(jsonBytes))
        } else {
            data.InterfaceName = types.StringNull()
        }
    } else if val, ok := dataMap["interfaceName"].(string); ok {
        data.InterfaceName = types.StringValue(val)
    } else {
        data.InterfaceName = types.StringNull()
    }
    if obj, ok := dataMap["bpfFilter"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.BpfFilter = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.BpfFilter = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.BpfFilter = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.BpfFilter = types.StringValue(string(jsonBytes))
            } else {
                data.BpfFilter = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.BpfFilter = types.StringValue(string(jsonBytes))
            } else {
                data.BpfFilter = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.BpfFilter = types.StringValue(string(jsonBytes))
        } else {
            data.BpfFilter = types.StringNull()
        }
    } else if val, ok := dataMap["bpfFilter"].(string); ok {
        data.BpfFilter = types.StringValue(val)
    } else {
        data.BpfFilter = types.StringNull()
    }
    if val, ok := dataMap["maxDurationInSeconds"].(float64); ok {
        data.MaxDurationInSeconds = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["maxDurationInSeconds"].(int); ok {
        data.MaxDurationInSeconds = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["maxDurationInSeconds"].(int64); ok {
        data.MaxDurationInSeconds = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["maxDurationInSeconds"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.MaxDurationInSeconds = types.NumberValue(big.NewFloat(val))
        } else {
            data.MaxDurationInSeconds = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.MaxDurationInSeconds = types.NumberNull()
    }
    if val, ok := dataMap["maxPackets"].(float64); ok {
        data.MaxPackets = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["maxPackets"].(int); ok {
        data.MaxPackets = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["maxPackets"].(int64); ok {
        data.MaxPackets = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["maxPackets"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.MaxPackets = types.NumberValue(big.NewFloat(val))
        } else {
            data.MaxPackets = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.MaxPackets = types.NumberNull()
    }
    if val, ok := dataMap["maxFileSizeInMB"].(float64); ok {
        data.MaxFileSizeInMb = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["maxFileSizeInMB"].(int); ok {
        data.MaxFileSizeInMb = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["maxFileSizeInMB"].(int64); ok {
        data.MaxFileSizeInMb = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["maxFileSizeInMB"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.MaxFileSizeInMb = types.NumberValue(big.NewFloat(val))
        } else {
            data.MaxFileSizeInMb = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.MaxFileSizeInMb = types.NumberNull()
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
    } else if val, ok := dataMap["name"].(string); ok {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if obj, ok := dataMap["status"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Status = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Status = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Status = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Status = types.StringValue(string(jsonBytes))
            } else {
                data.Status = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Status = types.StringValue(string(jsonBytes))
            } else {
                data.Status = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Status = types.StringValue(string(jsonBytes))
        } else {
            data.Status = types.StringNull()
        }
    } else if val, ok := dataMap["status"].(string); ok {
        data.Status = types.StringValue(val)
    } else {
        data.Status = types.StringNull()
    }
    if obj, ok := dataMap["statusMessage"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StatusMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.StatusMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.StatusMessage = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.StatusMessage = types.StringValue(string(jsonBytes))
            } else {
                data.StatusMessage = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.StatusMessage = types.StringValue(string(jsonBytes))
            } else {
                data.StatusMessage = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.StatusMessage = types.StringValue(string(jsonBytes))
        } else {
            data.StatusMessage = types.StringNull()
        }
    } else if val, ok := dataMap["statusMessage"].(string); ok {
        data.StatusMessage = types.StringValue(val)
    } else {
        data.StatusMessage = types.StringNull()
    }
    if obj, ok := dataMap["endReason"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EndReason = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.EndReason = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.EndReason = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.EndReason = types.StringValue(string(jsonBytes))
            } else {
                data.EndReason = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.EndReason = types.StringValue(string(jsonBytes))
            } else {
                data.EndReason = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.EndReason = types.StringValue(string(jsonBytes))
        } else {
            data.EndReason = types.StringNull()
        }
    } else if val, ok := dataMap["endReason"].(string); ok {
        data.EndReason = types.StringValue(val)
    } else {
        data.EndReason = types.StringNull()
    }
    if obj, ok := dataMap["startedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.StartedAt = NewRFC3339Value(val)
        } else {
            data.StartedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["startedAt"].(string); ok && val != "" {
        data.StartedAt = NewRFC3339Value(val)
    } else {
        data.StartedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["completedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.CompletedAt = NewRFC3339Value(val)
        } else {
            data.CompletedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["completedAt"].(string); ok && val != "" {
        data.CompletedAt = NewRFC3339Value(val)
    } else {
        data.CompletedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["stopRequestedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.StopRequestedAt = NewRFC3339Value(val)
        } else {
            data.StopRequestedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["stopRequestedAt"].(string); ok && val != "" {
        data.StopRequestedAt = NewRFC3339Value(val)
    } else {
        data.StopRequestedAt = NewRFC3339Null()
    }
    if val, ok := dataMap["packetCount"].(float64); ok {
        data.PacketCount = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["packetCount"].(int); ok {
        data.PacketCount = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["packetCount"].(int64); ok {
        data.PacketCount = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["packetCount"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.PacketCount = types.NumberValue(big.NewFloat(val))
        } else {
            data.PacketCount = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.PacketCount = types.NumberNull()
    }
    if val, ok := dataMap["fileSizeInBytes"].(float64); ok {
        data.FileSizeInBytes = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["fileSizeInBytes"].(int); ok {
        data.FileSizeInBytes = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["fileSizeInBytes"].(int64); ok {
        data.FileSizeInBytes = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["fileSizeInBytes"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.FileSizeInBytes = types.NumberValue(big.NewFloat(val))
        } else {
            data.FileSizeInBytes = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.FileSizeInBytes = types.NumberNull()
    }
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
    } else if val, ok := dataMap["createdByUserId"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PacketCaptureResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var data PacketCaptureResourceModel

    // Read Terraform plan data into the model
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // This resource does not have an update API endpoint.
    // Preserve the planned state.
    tflog.Trace(ctx, "updated a resource (no-op: preserving planned state)")

    // Save planned data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PacketCaptureResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var data PacketCaptureResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Make API call
    httpResp, err := r.client.Delete(ctx, "/packet-capture/" + data.Id.ValueString() + "")
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete packet_capture, got error: %s", err))
        return
    }

    // A failed delete must keep the resource in state — silently dropping it
    // orphans real infrastructure. 404 means it is already gone.
    if httpResp.StatusCode >= 400 && httpResp.StatusCode != http.StatusNotFound {
        err = r.client.ParseResponse(httpResp, nil)
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to delete packet_capture: %s", err))
        return
    }
    if httpResp.Body != nil {
        httpResp.Body.Close()
    }
}


func (r *PacketCaptureResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}


// keepPlannedValues puts back, after a create or an update, the planned value
// of each optional attribute the configuration leaves out. The server keeps
// some of these up to date on its own (when it last checked a heartbeat, the
// status a probe last reported...), so the value read back after the write can
// already differ from the plan, and Terraform would fail the apply with
// "Provider produced inconsistent result after apply". The next refresh reads
// the server's value, which is never a diff for an attribute nobody configured.
func (r *PacketCaptureResource) keepPlannedValues(data *PacketCaptureResourceModel, plan *PacketCaptureResourceModel, config *PacketCaptureResourceModel) {
    if config.NetworkDeviceId.IsNull() && !plan.NetworkDeviceId.IsUnknown() {
        data.NetworkDeviceId = plan.NetworkDeviceId
    }
    if config.BpfFilter.IsNull() && !plan.BpfFilter.IsUnknown() {
        data.BpfFilter = plan.BpfFilter
    }
    if config.MaxDurationInSeconds.IsNull() && !plan.MaxDurationInSeconds.IsUnknown() {
        data.MaxDurationInSeconds = plan.MaxDurationInSeconds
    }
    if config.MaxPackets.IsNull() && !plan.MaxPackets.IsUnknown() {
        data.MaxPackets = plan.MaxPackets
    }
    if config.MaxFileSizeInMb.IsNull() && !plan.MaxFileSizeInMb.IsUnknown() {
        data.MaxFileSizeInMb = plan.MaxFileSizeInMb
    }
}

// Helper method to convert Terraform map to Go interface{}
func (r *PacketCaptureResource) convertTerraformMapToInterface(terraformMap types.Map) interface{} {
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
func (r *PacketCaptureResource) convertTerraformListToInterface(terraformList types.List) interface{} {
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
func (r *PacketCaptureResource) convertTerraformSetToInterface(terraformSet types.Set) interface{} {
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
func (r *PacketCaptureResource) parseJSONField(terraformString basetypes.StringValuable) interface{} {
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
func (r *PacketCaptureResource) normalizeURLWrappers(value interface{}) interface{} {
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

func (r *PacketCaptureResource) normalizeURLString(value string) string {
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
func (r *PacketCaptureResource) bigFloatToFloat64(bf *big.Float) interface{} {
    if bf == nil {
        return nil
    }
    f, _ := bf.Float64()
    return f
}

// Helper method to check if a type string is a valid OneUptime ObjectType.
// The registry itself lives in objecttypes.go, shared across the package.
func (r *PacketCaptureResource) isValidOneUptimeObjectType(typeStr string) bool {
    return validOneUptimeObjectTypes[typeStr]
}
