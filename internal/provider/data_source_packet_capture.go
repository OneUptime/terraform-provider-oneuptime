package provider

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "math/big"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &PacketCaptureDataSource{}

func NewPacketCaptureDataSource() datasource.DataSource {
    return &PacketCaptureDataSource{}
}

// PacketCaptureDataSource defines the data source implementation.
type PacketCaptureDataSource struct {
    client *Client
}

// PacketCaptureDataSourceModel describes the data source data model.
type PacketCaptureDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    ProbeId types.String `tfsdk:"probe_id"`
    NetworkDeviceId types.String `tfsdk:"network_device_id"`
    Name types.String `tfsdk:"name"`
    InterfaceName types.String `tfsdk:"interface_name"`
    BpfFilter types.String `tfsdk:"bpf_filter"`
    MaxDurationInSeconds types.Number `tfsdk:"max_duration_in_seconds"`
    MaxPackets types.Number `tfsdk:"max_packets"`
    MaxFileSizeInMb types.Number `tfsdk:"max_file_size_in_mb"`
    Status types.String `tfsdk:"status"`
    StatusMessage types.String `tfsdk:"status_message"`
    EndReason types.String `tfsdk:"end_reason"`
    StartedAt types.String `tfsdk:"started_at"`
    CompletedAt types.String `tfsdk:"completed_at"`
    StopRequestedAt types.String `tfsdk:"stop_requested_at"`
    PacketCount types.Number `tfsdk:"packet_count"`
    FileSizeInBytes types.Number `tfsdk:"file_size_in_bytes"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *PacketCaptureDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_packet_capture"
}

func (d *PacketCaptureDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "A packet capture run on one of the project's probes, and the pcap file it produced. Captures and their files are deleted 7 days after they start. Look up an existing packet capture by `id`, or by any of its other arguments (`name`, `bpf_filter`, `created_by_user_id`, ...): each one set must match, and exactly one packet capture may match them all.",

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
            "project_id": schema.StringAttribute{
                MarkdownDescription: "ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "probe_id": schema.StringAttribute{
                MarkdownDescription: "ID of the probe that runs the capture: one of the project's own probes, with packet capture turned on. The ID of a `oneuptime_probe`.",
                Optional: true,
                Computed: true,
            },
            "network_device_id": schema.StringAttribute{
                MarkdownDescription: "ID of the network device the capture was started from, if it was started from a device's page. The ID of a `oneuptime_network_device`.",
                Optional: true,
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "The interface and the filter of the capture, written by OneUptime: eth0: host 10.0.0.5.",
                Optional: true,
                Computed: true,
            },
            "interface_name": schema.StringAttribute{
                MarkdownDescription: "The network interface the probe captures on, as the probe reported it: \"eth0\", or \"any\" for every interface at once.",
                Optional: true,
                Computed: true,
            },
            "bpf_filter": schema.StringAttribute{
                MarkdownDescription: "A BPF filter expression - the capture-filter language of tcpdump and Wireshark - that decides which packets are kept: host 10.0.0.5 and tcp port 443. Empty keeps every packet.",
                Optional: true,
                Computed: true,
            },
            "max_duration_in_seconds": schema.NumberAttribute{
                MarkdownDescription: "How long the capture runs, in seconds, at most: from 5 seconds to 30 minutes, and no longer than the probe allows. It stops earlier when it reaches its packet or file size limit.",
                Optional: true,
                Computed: true,
            },
            "max_packets": schema.NumberAttribute{
                MarkdownDescription: "The capture stops after this many packets: from 1 to 1,000,000.",
                Optional: true,
                Computed: true,
            },
            "max_file_size_in_mb": schema.NumberAttribute{
                MarkdownDescription: "The capture stops when its file reaches this many megabytes: from 1 to 25, and no more than the probe allows. The file is cut at the last whole packet.",
                Optional: true,
                Computed: true,
            },
            "status": schema.StringAttribute{
                MarkdownDescription: "Where the capture is: \"Pending\" (waiting for the probe), \"Running\" (the probe is capturing), \"Completed\" (the file is ready) or \"Failed\" (see Status Message). Managed by OneUptime and the probe.",
                Optional: true,
                Computed: true,
            },
            "status_message": schema.StringAttribute{
                MarkdownDescription: "Why a capture failed - the interface is gone, the probe may not capture, the filter did not compile - or what the capture tool said when it stopped by itself. Managed by OneUptime and the probe.",
                Optional: true,
                Computed: true,
            },
            "end_reason": schema.StringAttribute{
                MarkdownDescription: "Why a completed capture stopped: \"DurationReached\", \"PacketLimitReached\", \"FileSizeLimitReached\", \"StoppedFromDashboard\" or \"CaptureToolStopped\". Managed by the probe.",
                Optional: true,
                Computed: true,
            },
            "started_at": schema.StringAttribute{
                MarkdownDescription: "When the probe started capturing. Managed by OneUptime.",
                Computed: true,
            },
            "completed_at": schema.StringAttribute{
                MarkdownDescription: "When the capture completed or failed. Managed by OneUptime.",
                Computed: true,
            },
            "stop_requested_at": schema.StringAttribute{
                MarkdownDescription: "When someone pressed Stop on the running capture. The probe stops at its next check, within about ten seconds, and uploads what it captured. Managed by OneUptime.",
                Computed: true,
            },
            "packet_count": schema.NumberAttribute{
                MarkdownDescription: "How many packets the file holds, counted by OneUptime from the file. Managed by OneUptime.",
                Optional: true,
                Computed: true,
            },
            "file_size_in_bytes": schema.NumberAttribute{
                MarkdownDescription: "The size of the pcap file in bytes. Managed by OneUptime.",
                Optional: true,
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *PacketCaptureDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *PacketCaptureDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data PacketCaptureDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.ProbeId.IsNull() && !data.ProbeId.IsUnknown() {
        filters["probeId"] = data.ProbeId.ValueString()
        filterNames = append(filterNames, "probe_id = "+fmt.Sprintf("%q", data.ProbeId.ValueString()))
    }
    if !data.NetworkDeviceId.IsNull() && !data.NetworkDeviceId.IsUnknown() {
        filters["networkDeviceId"] = data.NetworkDeviceId.ValueString()
        filterNames = append(filterNames, "network_device_id = "+fmt.Sprintf("%q", data.NetworkDeviceId.ValueString()))
    }
    if !data.Name.IsNull() && !data.Name.IsUnknown() {
        filters["name"] = data.Name.ValueString()
        filterNames = append(filterNames, "name = "+fmt.Sprintf("%q", data.Name.ValueString()))
    }
    if !data.InterfaceName.IsNull() && !data.InterfaceName.IsUnknown() {
        filters["interfaceName"] = data.InterfaceName.ValueString()
        filterNames = append(filterNames, "interface_name = "+fmt.Sprintf("%q", data.InterfaceName.ValueString()))
    }
    if !data.BpfFilter.IsNull() && !data.BpfFilter.IsUnknown() {
        filters["bpfFilter"] = data.BpfFilter.ValueString()
        filterNames = append(filterNames, "bpf_filter = "+fmt.Sprintf("%q", data.BpfFilter.ValueString()))
    }
    if !data.MaxDurationInSeconds.IsNull() && !data.MaxDurationInSeconds.IsUnknown() {
        filters["maxDurationInSeconds"] = lookupNumber(data.MaxDurationInSeconds)
        filterNames = append(filterNames, "max_duration_in_seconds = "+data.MaxDurationInSeconds.ValueBigFloat().String())
    }
    if !data.MaxPackets.IsNull() && !data.MaxPackets.IsUnknown() {
        filters["maxPackets"] = lookupNumber(data.MaxPackets)
        filterNames = append(filterNames, "max_packets = "+data.MaxPackets.ValueBigFloat().String())
    }
    if !data.MaxFileSizeInMb.IsNull() && !data.MaxFileSizeInMb.IsUnknown() {
        filters["maxFileSizeInMB"] = lookupNumber(data.MaxFileSizeInMb)
        filterNames = append(filterNames, "max_file_size_in_mb = "+data.MaxFileSizeInMb.ValueBigFloat().String())
    }
    if !data.Status.IsNull() && !data.Status.IsUnknown() {
        filters["status"] = data.Status.ValueString()
        filterNames = append(filterNames, "status = "+fmt.Sprintf("%q", data.Status.ValueString()))
    }
    if !data.StatusMessage.IsNull() && !data.StatusMessage.IsUnknown() {
        filters["statusMessage"] = data.StatusMessage.ValueString()
        filterNames = append(filterNames, "status_message = "+fmt.Sprintf("%q", data.StatusMessage.ValueString()))
    }
    if !data.EndReason.IsNull() && !data.EndReason.IsUnknown() {
        filters["endReason"] = data.EndReason.ValueString()
        filterNames = append(filterNames, "end_reason = "+fmt.Sprintf("%q", data.EndReason.ValueString()))
    }
    if !data.PacketCount.IsNull() && !data.PacketCount.IsUnknown() {
        filters["packetCount"] = lookupNumber(data.PacketCount)
        filterNames = append(filterNames, "packet_count = "+data.PacketCount.ValueBigFloat().String())
    }
    if !data.FileSizeInBytes.IsNull() && !data.FileSizeInBytes.IsUnknown() {
        filters["fileSizeInBytes"] = lookupNumber(data.FileSizeInBytes)
        filterNames = append(filterNames, "file_size_in_bytes = "+data.FileSizeInBytes.ValueBigFloat().String())
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the packet capture up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the packet capture up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "probeId": true,
        "networkDeviceId": true,
        "name": true,
        "interfaceName": true,
        "bpfFilter": true,
        "maxDurationInSeconds": true,
        "maxPackets": true,
        "maxFileSizeInMB": true,
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

    var item map[string]interface{}
    if hasId {
        readPath := "/packet-capture/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read packet_capture, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No packet capture found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read packet_capture: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/packet-capture/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list packet_capture, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list packet_capture: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No packet capture matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one packet capture matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for packet_capture.")
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
    if obj, ok := item["probeId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ProbeId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ProbeId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ProbeId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ProbeId = types.StringValue(string(jsonBytes))
        } else {
            data.ProbeId = types.StringNull()
        }
    } else if val, ok := item["probeId"].(string); ok {
        data.ProbeId = types.StringValue(val)
    } else {
        data.ProbeId = types.StringNull()
    }
    if obj, ok := item["networkDeviceId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.NetworkDeviceId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.NetworkDeviceId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.NetworkDeviceId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.NetworkDeviceId = types.StringValue(string(jsonBytes))
        } else {
            data.NetworkDeviceId = types.StringNull()
        }
    } else if val, ok := item["networkDeviceId"].(string); ok {
        data.NetworkDeviceId = types.StringValue(val)
    } else {
        data.NetworkDeviceId = types.StringNull()
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
    if obj, ok := item["interfaceName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.InterfaceName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.InterfaceName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.InterfaceName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.InterfaceName = types.StringValue(string(jsonBytes))
        } else {
            data.InterfaceName = types.StringNull()
        }
    } else if val, ok := item["interfaceName"].(string); ok {
        data.InterfaceName = types.StringValue(val)
    } else {
        data.InterfaceName = types.StringNull()
    }
    if obj, ok := item["bpfFilter"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.BpfFilter = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.BpfFilter = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.BpfFilter = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.BpfFilter = types.StringValue(string(jsonBytes))
        } else {
            data.BpfFilter = types.StringNull()
        }
    } else if val, ok := item["bpfFilter"].(string); ok {
        data.BpfFilter = types.StringValue(val)
    } else {
        data.BpfFilter = types.StringNull()
    }
    if val, ok := item["maxDurationInSeconds"].(float64); ok {
        data.MaxDurationInSeconds = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["maxDurationInSeconds"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.MaxDurationInSeconds = types.NumberValue(big.NewFloat(val))
        } else {
            data.MaxDurationInSeconds = types.NumberNull()
        }
    } else {
        data.MaxDurationInSeconds = types.NumberNull()
    }
    if val, ok := item["maxPackets"].(float64); ok {
        data.MaxPackets = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["maxPackets"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.MaxPackets = types.NumberValue(big.NewFloat(val))
        } else {
            data.MaxPackets = types.NumberNull()
        }
    } else {
        data.MaxPackets = types.NumberNull()
    }
    if val, ok := item["maxFileSizeInMB"].(float64); ok {
        data.MaxFileSizeInMb = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["maxFileSizeInMB"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.MaxFileSizeInMb = types.NumberValue(big.NewFloat(val))
        } else {
            data.MaxFileSizeInMb = types.NumberNull()
        }
    } else {
        data.MaxFileSizeInMb = types.NumberNull()
    }
    if obj, ok := item["status"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Status = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Status = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Status = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Status = types.StringValue(string(jsonBytes))
        } else {
            data.Status = types.StringNull()
        }
    } else if val, ok := item["status"].(string); ok {
        data.Status = types.StringValue(val)
    } else {
        data.Status = types.StringNull()
    }
    if obj, ok := item["statusMessage"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StatusMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.StatusMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.StatusMessage = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.StatusMessage = types.StringValue(string(jsonBytes))
        } else {
            data.StatusMessage = types.StringNull()
        }
    } else if val, ok := item["statusMessage"].(string); ok {
        data.StatusMessage = types.StringValue(val)
    } else {
        data.StatusMessage = types.StringNull()
    }
    if obj, ok := item["endReason"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EndReason = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.EndReason = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.EndReason = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.EndReason = types.StringValue(string(jsonBytes))
        } else {
            data.EndReason = types.StringNull()
        }
    } else if val, ok := item["endReason"].(string); ok {
        data.EndReason = types.StringValue(val)
    } else {
        data.EndReason = types.StringNull()
    }
    if obj, ok := item["startedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StartedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.StartedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.StartedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.StartedAt = types.StringValue(string(jsonBytes))
        } else {
            data.StartedAt = types.StringNull()
        }
    } else if val, ok := item["startedAt"].(string); ok {
        data.StartedAt = types.StringValue(val)
    } else {
        data.StartedAt = types.StringNull()
    }
    if obj, ok := item["completedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CompletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CompletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CompletedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CompletedAt = types.StringValue(string(jsonBytes))
        } else {
            data.CompletedAt = types.StringNull()
        }
    } else if val, ok := item["completedAt"].(string); ok {
        data.CompletedAt = types.StringValue(val)
    } else {
        data.CompletedAt = types.StringNull()
    }
    if obj, ok := item["stopRequestedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StopRequestedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.StopRequestedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.StopRequestedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.StopRequestedAt = types.StringValue(string(jsonBytes))
        } else {
            data.StopRequestedAt = types.StringNull()
        }
    } else if val, ok := item["stopRequestedAt"].(string); ok {
        data.StopRequestedAt = types.StringValue(val)
    } else {
        data.StopRequestedAt = types.StringNull()
    }
    if val, ok := item["packetCount"].(float64); ok {
        data.PacketCount = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["packetCount"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.PacketCount = types.NumberValue(big.NewFloat(val))
        } else {
            data.PacketCount = types.NumberNull()
        }
    } else {
        data.PacketCount = types.NumberNull()
    }
    if val, ok := item["fileSizeInBytes"].(float64); ok {
        data.FileSizeInBytes = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["fileSizeInBytes"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.FileSizeInBytes = types.NumberValue(big.NewFloat(val))
        } else {
            data.FileSizeInBytes = types.NumberNull()
        }
    } else {
        data.FileSizeInBytes = types.NumberNull()
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

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
