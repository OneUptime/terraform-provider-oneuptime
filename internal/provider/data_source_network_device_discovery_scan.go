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
var _ datasource.DataSource = &NetworkDeviceDiscoveryScanDataSource{}

func NewNetworkDeviceDiscoveryScanDataSource() datasource.DataSource {
    return &NetworkDeviceDiscoveryScanDataSource{}
}

// NetworkDeviceDiscoveryScanDataSource defines the data source implementation.
type NetworkDeviceDiscoveryScanDataSource struct {
    client *Client
}

// NetworkDeviceDiscoveryScanDataSourceModel describes the data source data model.
type NetworkDeviceDiscoveryScanDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    ProbeId types.String `tfsdk:"probe_id"`
    Name types.String `tfsdk:"name"`
    Cidr types.String `tfsdk:"cidr"`
    SnmpConfigs types.String `tfsdk:"snmp_configs"`
    IsSnmpEnabled types.Bool `tfsdk:"is_snmp_enabled"`
    UseShortDeviceNames types.Bool `tfsdk:"use_short_device_names"`
    SnmpVersion types.String `tfsdk:"snmp_version"`
    SnmpCommunityString types.String `tfsdk:"snmp_community_string"`
    SnmpPort types.Number `tfsdk:"snmp_port"`
    SnmpV3SecurityLevel types.String `tfsdk:"snmp_v3_security_level"`
    SnmpV3Username types.String `tfsdk:"snmp_v3_username"`
    SnmpV3AuthProtocol types.String `tfsdk:"snmp_v3_auth_protocol"`
    SnmpV3AuthKey types.String `tfsdk:"snmp_v3_auth_key"`
    SnmpV3PrivProtocol types.String `tfsdk:"snmp_v3_priv_protocol"`
    SnmpV3PrivKey types.String `tfsdk:"snmp_v3_priv_key"`
    Status types.String `tfsdk:"status"`
    StatusMessage types.String `tfsdk:"status_message"`
    IsNetbiosLookupEnabled types.Bool `tfsdk:"is_netbios_lookup_enabled"`
    DiscoveredDevices types.String `tfsdk:"discovered_devices"`
    ScannedHostCount types.Number `tfsdk:"scanned_host_count"`
    RespondedHostCount types.Number `tfsdk:"responded_host_count"`
    StartedAt types.String `tfsdk:"started_at"`
    CompletedAt types.String `tfsdk:"completed_at"`
    IsRecurring types.Bool `tfsdk:"is_recurring"`
    RescanIntervalInMinutes types.Number `tfsdk:"rescan_interval_in_minutes"`
    NextScanAt types.String `tfsdk:"next_scan_at"`
    AutoImportProcessedAt types.String `tfsdk:"auto_import_processed_at"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *NetworkDeviceDiscoveryScanDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_network_device_discovery_scan"
}

func (d *NetworkDeviceDiscoveryScanDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Network discovery scans that sweep an address space — a CIDR subnet or an octet range — from a probe and report the hosts found, so they can be imported as Network Devices. Every sweep pings; scans with Check SNMP on also query each live host over SNMP. Look up an existing network device discovery scan by `id`, or by any of its other arguments (`name`, `cidr`, `created_by_user_id`, ...): each one set must match, and exactly one network device discovery scan may match them all.",

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
                MarkdownDescription: "ID of the Probe that runs this discovery scan. The ID of a `oneuptime_probe`.",
                Optional: true,
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Optional name for this scan, so it can be told apart from other scans at a glance. Falls back to the scan target when empty.",
                Optional: true,
                Computed: true,
            },
            "cidr": schema.StringAttribute{
                MarkdownDescription: "Address space to scan, either in CIDR notation (192.168.1.0/24) or octet-range notation where any octet may be an inclusive low-high range (10.16-22.0-255.51-66).",
                Optional: true,
                Computed: true,
            },
            "snmp_configs": schema.StringAttribute{
                MarkdownDescription: "Ordered list of SNMP credential sets tried against every host in the subnet, first match wins. Each entry carries an id, an optional name, a version, a community string or the v3 credentials, and a port. When empty, the scan uses the single flattened SNMP configuration on this row. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "is_snmp_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether hosts that answer the ping sweep are then queried over SNMP. Turn it off for an ICMP-only scan, which reports every host that answers ping and asks nothing else of them.",
                Optional: true,
                Computed: true,
            },
            "use_short_device_names": schema.BoolAttribute{
                MarkdownDescription: "Name imported devices by their short hostname (the first label of a fully qualified name, e.g. 'core-sw-01' rather than 'core-sw-01.corp.example.com'). The full reverse-DNS name is still stored on the device as its DNS Name.",
                Optional: true,
                Computed: true,
            },
            "snmp_version": schema.StringAttribute{
                MarkdownDescription: "SNMP version tried against every host in the subnet (V1, V2c, V3). Ignored when Check SNMP is off.",
                Optional: true,
                Computed: true,
            },
            "snmp_community_string": schema.StringAttribute{
                MarkdownDescription: "Community string tried against every host in the subnet (SNMP v1/v2c). Ignored when Check SNMP is off.",
                Optional: true,
                Computed: true,
            },
            "snmp_port": schema.NumberAttribute{
                MarkdownDescription: "UDP port tried against every host in the subnet. Ignored when Check SNMP is off.",
                Optional: true,
                Computed: true,
            },
            "snmp_v3_security_level": schema.StringAttribute{
                MarkdownDescription: "SNMP v3 security level tried against every host: noAuthNoPriv, authNoPriv, or authPriv. Ignored when Check SNMP is off.",
                Optional: true,
                Computed: true,
            },
            "snmp_v3_username": schema.StringAttribute{
                MarkdownDescription: "SNMP v3 security name (username) tried against every host. Ignored when Check SNMP is off.",
                Optional: true,
                Computed: true,
            },
            "snmp_v3_auth_protocol": schema.StringAttribute{
                MarkdownDescription: "SNMP v3 authentication protocol: MD5, SHA, SHA256, or SHA512. Ignored when Check SNMP is off.",
                Optional: true,
                Computed: true,
            },
            "snmp_v3_auth_key": schema.StringAttribute{
                MarkdownDescription: "SNMP v3 authentication passphrase tried against every host. Ignored when Check SNMP is off.",
                Optional: true,
                Computed: true,
            },
            "snmp_v3_priv_protocol": schema.StringAttribute{
                MarkdownDescription: "SNMP v3 privacy (encryption) protocol: DES, AES, or AES256. Ignored when Check SNMP is off.",
                Optional: true,
                Computed: true,
            },
            "snmp_v3_priv_key": schema.StringAttribute{
                MarkdownDescription: "SNMP v3 privacy (encryption) passphrase tried against every host. Ignored when Check SNMP is off.",
                Optional: true,
                Computed: true,
            },
            "status": schema.StringAttribute{
                MarkdownDescription: "Status of this discovery scan: \"Pending\", \"In Progress\", \"Completed\" or \"Failed\". Managed by the scanning probe.",
                Optional: true,
                Computed: true,
            },
            "status_message": schema.StringAttribute{
                MarkdownDescription: "Details about the current status of this scan, e.g. the failure reason. Managed by the scanning probe.",
                Optional: true,
                Computed: true,
            },
            "is_netbios_lookup_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether hosts with no SNMP name and no reverse DNS record are asked for their NetBIOS name over UDP 137. Best-effort: Windows/Samba hosts that allow UDP 137 from the probe. Private addresses only; never done by global probes.",
                Optional: true,
                Computed: true,
            },
            "discovered_devices": schema.StringAttribute{
                MarkdownDescription: "Devices found by this scan: array of {ipAddress, sysName, sysDescr, isAlreadyRegistered}. Managed by the scanning probe. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "scanned_host_count": schema.NumberAttribute{
                MarkdownDescription: "Total number of host addresses swept in the subnet. Managed by the scanning probe.",
                Optional: true,
                Computed: true,
            },
            "responded_host_count": schema.NumberAttribute{
                MarkdownDescription: "Number of hosts that answered the check this scan performed: SNMP responders on a scan with Check SNMP on, hosts that answered the ping sweep on an ICMP-only one. Managed by the scanning probe.",
                Optional: true,
                Computed: true,
            },
            "started_at": schema.StringAttribute{
                MarkdownDescription: "When the scanning probe started this scan. Managed by the scanning probe.",
                Computed: true,
            },
            "completed_at": schema.StringAttribute{
                MarkdownDescription: "When the scanning probe completed (or failed) this scan. Managed by the scanning probe.",
                Computed: true,
            },
            "is_recurring": schema.BoolAttribute{
                MarkdownDescription: "Re-run this scan automatically every Rescan Interval minutes to keep discovery continuous.",
                Optional: true,
                Computed: true,
            },
            "rescan_interval_in_minutes": schema.NumberAttribute{
                MarkdownDescription: "How often a recurring scan re-runs, in minutes. Ignored unless Is Recurring is on.",
                Optional: true,
                Computed: true,
            },
            "next_scan_at": schema.StringAttribute{
                MarkdownDescription: "When a recurring scan is next due to run. Managed by the server.",
                Computed: true,
            },
            "auto_import_processed_at": schema.StringAttribute{
                MarkdownDescription: "When auto-import rules last processed this scan's results. Managed by the server: cleared when new results arrive, stamped by the worker that evaluates the rules. NULL means the current results have not been processed yet.",
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

func (d *NetworkDeviceDiscoveryScanDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *NetworkDeviceDiscoveryScanDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data NetworkDeviceDiscoveryScanDataSourceModel

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
    if !data.Name.IsNull() && !data.Name.IsUnknown() {
        filters["name"] = data.Name.ValueString()
        filterNames = append(filterNames, "name = "+fmt.Sprintf("%q", data.Name.ValueString()))
    }
    if !data.Cidr.IsNull() && !data.Cidr.IsUnknown() {
        filters["cidr"] = data.Cidr.ValueString()
        filterNames = append(filterNames, "cidr = "+fmt.Sprintf("%q", data.Cidr.ValueString()))
    }
    if !data.IsSnmpEnabled.IsNull() && !data.IsSnmpEnabled.IsUnknown() {
        filters["isSnmpEnabled"] = data.IsSnmpEnabled.ValueBool()
        filterNames = append(filterNames, "is_snmp_enabled = "+fmt.Sprintf("%t", data.IsSnmpEnabled.ValueBool()))
    }
    if !data.UseShortDeviceNames.IsNull() && !data.UseShortDeviceNames.IsUnknown() {
        filters["useShortDeviceNames"] = data.UseShortDeviceNames.ValueBool()
        filterNames = append(filterNames, "use_short_device_names = "+fmt.Sprintf("%t", data.UseShortDeviceNames.ValueBool()))
    }
    if !data.SnmpVersion.IsNull() && !data.SnmpVersion.IsUnknown() {
        filters["snmpVersion"] = data.SnmpVersion.ValueString()
        filterNames = append(filterNames, "snmp_version = "+fmt.Sprintf("%q", data.SnmpVersion.ValueString()))
    }
    if !data.SnmpCommunityString.IsNull() && !data.SnmpCommunityString.IsUnknown() {
        filters["snmpCommunityString"] = data.SnmpCommunityString.ValueString()
        filterNames = append(filterNames, "snmp_community_string = "+fmt.Sprintf("%q", data.SnmpCommunityString.ValueString()))
    }
    if !data.SnmpPort.IsNull() && !data.SnmpPort.IsUnknown() {
        filters["snmpPort"] = lookupNumber(data.SnmpPort)
        filterNames = append(filterNames, "snmp_port = "+data.SnmpPort.ValueBigFloat().String())
    }
    if !data.SnmpV3SecurityLevel.IsNull() && !data.SnmpV3SecurityLevel.IsUnknown() {
        filters["snmpV3SecurityLevel"] = data.SnmpV3SecurityLevel.ValueString()
        filterNames = append(filterNames, "snmp_v3_security_level = "+fmt.Sprintf("%q", data.SnmpV3SecurityLevel.ValueString()))
    }
    if !data.SnmpV3Username.IsNull() && !data.SnmpV3Username.IsUnknown() {
        filters["snmpV3Username"] = data.SnmpV3Username.ValueString()
        filterNames = append(filterNames, "snmp_v3_username = "+fmt.Sprintf("%q", data.SnmpV3Username.ValueString()))
    }
    if !data.SnmpV3AuthProtocol.IsNull() && !data.SnmpV3AuthProtocol.IsUnknown() {
        filters["snmpV3AuthProtocol"] = data.SnmpV3AuthProtocol.ValueString()
        filterNames = append(filterNames, "snmp_v3_auth_protocol = "+fmt.Sprintf("%q", data.SnmpV3AuthProtocol.ValueString()))
    }
    if !data.SnmpV3AuthKey.IsNull() && !data.SnmpV3AuthKey.IsUnknown() {
        filters["snmpV3AuthKey"] = data.SnmpV3AuthKey.ValueString()
        filterNames = append(filterNames, "snmp_v3_auth_key = "+fmt.Sprintf("%q", data.SnmpV3AuthKey.ValueString()))
    }
    if !data.SnmpV3PrivProtocol.IsNull() && !data.SnmpV3PrivProtocol.IsUnknown() {
        filters["snmpV3PrivProtocol"] = data.SnmpV3PrivProtocol.ValueString()
        filterNames = append(filterNames, "snmp_v3_priv_protocol = "+fmt.Sprintf("%q", data.SnmpV3PrivProtocol.ValueString()))
    }
    if !data.SnmpV3PrivKey.IsNull() && !data.SnmpV3PrivKey.IsUnknown() {
        filters["snmpV3PrivKey"] = data.SnmpV3PrivKey.ValueString()
        filterNames = append(filterNames, "snmp_v3_priv_key = "+fmt.Sprintf("%q", data.SnmpV3PrivKey.ValueString()))
    }
    if !data.Status.IsNull() && !data.Status.IsUnknown() {
        filters["status"] = data.Status.ValueString()
        filterNames = append(filterNames, "status = "+fmt.Sprintf("%q", data.Status.ValueString()))
    }
    if !data.StatusMessage.IsNull() && !data.StatusMessage.IsUnknown() {
        filters["statusMessage"] = data.StatusMessage.ValueString()
        filterNames = append(filterNames, "status_message = "+fmt.Sprintf("%q", data.StatusMessage.ValueString()))
    }
    if !data.IsNetbiosLookupEnabled.IsNull() && !data.IsNetbiosLookupEnabled.IsUnknown() {
        filters["isNetbiosLookupEnabled"] = data.IsNetbiosLookupEnabled.ValueBool()
        filterNames = append(filterNames, "is_netbios_lookup_enabled = "+fmt.Sprintf("%t", data.IsNetbiosLookupEnabled.ValueBool()))
    }
    if !data.ScannedHostCount.IsNull() && !data.ScannedHostCount.IsUnknown() {
        filters["scannedHostCount"] = lookupNumber(data.ScannedHostCount)
        filterNames = append(filterNames, "scanned_host_count = "+data.ScannedHostCount.ValueBigFloat().String())
    }
    if !data.RespondedHostCount.IsNull() && !data.RespondedHostCount.IsUnknown() {
        filters["respondedHostCount"] = lookupNumber(data.RespondedHostCount)
        filterNames = append(filterNames, "responded_host_count = "+data.RespondedHostCount.ValueBigFloat().String())
    }
    if !data.IsRecurring.IsNull() && !data.IsRecurring.IsUnknown() {
        filters["isRecurring"] = data.IsRecurring.ValueBool()
        filterNames = append(filterNames, "is_recurring = "+fmt.Sprintf("%t", data.IsRecurring.ValueBool()))
    }
    if !data.RescanIntervalInMinutes.IsNull() && !data.RescanIntervalInMinutes.IsUnknown() {
        filters["rescanIntervalInMinutes"] = lookupNumber(data.RescanIntervalInMinutes)
        filterNames = append(filterNames, "rescan_interval_in_minutes = "+data.RescanIntervalInMinutes.ValueBigFloat().String())
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the network device discovery scan up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the network device discovery scan up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "probeId": true,
        "name": true,
        "cidr": true,
        "snmpConfigs": true,
        "isSnmpEnabled": true,
        "useShortDeviceNames": true,
        "snmpVersion": true,
        "snmpCommunityString": true,
        "snmpPort": true,
        "snmpV3SecurityLevel": true,
        "snmpV3Username": true,
        "snmpV3AuthProtocol": true,
        "snmpV3AuthKey": true,
        "snmpV3PrivProtocol": true,
        "snmpV3PrivKey": true,
        "status": true,
        "statusMessage": true,
        "isNetbiosLookupEnabled": true,
        "discoveredDevices": true,
        "scannedHostCount": true,
        "respondedHostCount": true,
        "startedAt": true,
        "completedAt": true,
        "isRecurring": true,
        "rescanIntervalInMinutes": true,
        "nextScanAt": true,
        "autoImportProcessedAt": true,
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/network-device-discovery-scan/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read network_device_discovery_scan, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No network device discovery scan found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read network_device_discovery_scan: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/network-device-discovery-scan/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list network_device_discovery_scan, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list network_device_discovery_scan: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No network device discovery scan matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one network device discovery scan matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for network_device_discovery_scan.")
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
    if obj, ok := item["cidr"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Cidr = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Cidr = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Cidr = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Cidr = types.StringValue(string(jsonBytes))
        } else {
            data.Cidr = types.StringNull()
        }
    } else if val, ok := item["cidr"].(string); ok {
        data.Cidr = types.StringValue(val)
    } else {
        data.Cidr = types.StringNull()
    }
    if obj, ok := item["snmpConfigs"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SnmpConfigs = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SnmpConfigs = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SnmpConfigs = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SnmpConfigs = types.StringValue(string(jsonBytes))
        } else {
            data.SnmpConfigs = types.StringNull()
        }
    } else if val, ok := item["snmpConfigs"].(string); ok {
        data.SnmpConfigs = types.StringValue(val)
    } else {
        data.SnmpConfigs = types.StringNull()
    }
    if val, ok := item["isSnmpEnabled"].(bool); ok {
        data.IsSnmpEnabled = types.BoolValue(val)
    } else {
        data.IsSnmpEnabled = types.BoolNull()
    }
    if val, ok := item["useShortDeviceNames"].(bool); ok {
        data.UseShortDeviceNames = types.BoolValue(val)
    } else {
        data.UseShortDeviceNames = types.BoolNull()
    }
    if obj, ok := item["snmpVersion"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SnmpVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SnmpVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SnmpVersion = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SnmpVersion = types.StringValue(string(jsonBytes))
        } else {
            data.SnmpVersion = types.StringNull()
        }
    } else if val, ok := item["snmpVersion"].(string); ok {
        data.SnmpVersion = types.StringValue(val)
    } else {
        data.SnmpVersion = types.StringNull()
    }
    if obj, ok := item["snmpCommunityString"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SnmpCommunityString = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SnmpCommunityString = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SnmpCommunityString = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SnmpCommunityString = types.StringValue(string(jsonBytes))
        } else {
            data.SnmpCommunityString = types.StringNull()
        }
    } else if val, ok := item["snmpCommunityString"].(string); ok {
        data.SnmpCommunityString = types.StringValue(val)
    } else {
        data.SnmpCommunityString = types.StringNull()
    }
    if val, ok := item["snmpPort"].(float64); ok {
        data.SnmpPort = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["snmpPort"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.SnmpPort = types.NumberValue(big.NewFloat(val))
        } else {
            data.SnmpPort = types.NumberNull()
        }
    } else {
        data.SnmpPort = types.NumberNull()
    }
    if obj, ok := item["snmpV3SecurityLevel"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SnmpV3SecurityLevel = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SnmpV3SecurityLevel = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SnmpV3SecurityLevel = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SnmpV3SecurityLevel = types.StringValue(string(jsonBytes))
        } else {
            data.SnmpV3SecurityLevel = types.StringNull()
        }
    } else if val, ok := item["snmpV3SecurityLevel"].(string); ok {
        data.SnmpV3SecurityLevel = types.StringValue(val)
    } else {
        data.SnmpV3SecurityLevel = types.StringNull()
    }
    if obj, ok := item["snmpV3Username"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SnmpV3Username = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SnmpV3Username = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SnmpV3Username = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SnmpV3Username = types.StringValue(string(jsonBytes))
        } else {
            data.SnmpV3Username = types.StringNull()
        }
    } else if val, ok := item["snmpV3Username"].(string); ok {
        data.SnmpV3Username = types.StringValue(val)
    } else {
        data.SnmpV3Username = types.StringNull()
    }
    if obj, ok := item["snmpV3AuthProtocol"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SnmpV3AuthProtocol = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SnmpV3AuthProtocol = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SnmpV3AuthProtocol = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SnmpV3AuthProtocol = types.StringValue(string(jsonBytes))
        } else {
            data.SnmpV3AuthProtocol = types.StringNull()
        }
    } else if val, ok := item["snmpV3AuthProtocol"].(string); ok {
        data.SnmpV3AuthProtocol = types.StringValue(val)
    } else {
        data.SnmpV3AuthProtocol = types.StringNull()
    }
    if obj, ok := item["snmpV3AuthKey"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SnmpV3AuthKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SnmpV3AuthKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SnmpV3AuthKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SnmpV3AuthKey = types.StringValue(string(jsonBytes))
        } else {
            data.SnmpV3AuthKey = types.StringNull()
        }
    } else if val, ok := item["snmpV3AuthKey"].(string); ok {
        data.SnmpV3AuthKey = types.StringValue(val)
    } else {
        data.SnmpV3AuthKey = types.StringNull()
    }
    if obj, ok := item["snmpV3PrivProtocol"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SnmpV3PrivProtocol = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SnmpV3PrivProtocol = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SnmpV3PrivProtocol = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SnmpV3PrivProtocol = types.StringValue(string(jsonBytes))
        } else {
            data.SnmpV3PrivProtocol = types.StringNull()
        }
    } else if val, ok := item["snmpV3PrivProtocol"].(string); ok {
        data.SnmpV3PrivProtocol = types.StringValue(val)
    } else {
        data.SnmpV3PrivProtocol = types.StringNull()
    }
    if obj, ok := item["snmpV3PrivKey"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SnmpV3PrivKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SnmpV3PrivKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SnmpV3PrivKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SnmpV3PrivKey = types.StringValue(string(jsonBytes))
        } else {
            data.SnmpV3PrivKey = types.StringNull()
        }
    } else if val, ok := item["snmpV3PrivKey"].(string); ok {
        data.SnmpV3PrivKey = types.StringValue(val)
    } else {
        data.SnmpV3PrivKey = types.StringNull()
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
    if val, ok := item["isNetbiosLookupEnabled"].(bool); ok {
        data.IsNetbiosLookupEnabled = types.BoolValue(val)
    } else {
        data.IsNetbiosLookupEnabled = types.BoolNull()
    }
    if obj, ok := item["discoveredDevices"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DiscoveredDevices = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DiscoveredDevices = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DiscoveredDevices = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DiscoveredDevices = types.StringValue(string(jsonBytes))
        } else {
            data.DiscoveredDevices = types.StringNull()
        }
    } else if val, ok := item["discoveredDevices"].(string); ok {
        data.DiscoveredDevices = types.StringValue(val)
    } else {
        data.DiscoveredDevices = types.StringNull()
    }
    if val, ok := item["scannedHostCount"].(float64); ok {
        data.ScannedHostCount = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["scannedHostCount"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.ScannedHostCount = types.NumberValue(big.NewFloat(val))
        } else {
            data.ScannedHostCount = types.NumberNull()
        }
    } else {
        data.ScannedHostCount = types.NumberNull()
    }
    if val, ok := item["respondedHostCount"].(float64); ok {
        data.RespondedHostCount = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["respondedHostCount"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.RespondedHostCount = types.NumberValue(big.NewFloat(val))
        } else {
            data.RespondedHostCount = types.NumberNull()
        }
    } else {
        data.RespondedHostCount = types.NumberNull()
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
    if val, ok := item["isRecurring"].(bool); ok {
        data.IsRecurring = types.BoolValue(val)
    } else {
        data.IsRecurring = types.BoolNull()
    }
    if val, ok := item["rescanIntervalInMinutes"].(float64); ok {
        data.RescanIntervalInMinutes = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["rescanIntervalInMinutes"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.RescanIntervalInMinutes = types.NumberValue(big.NewFloat(val))
        } else {
            data.RescanIntervalInMinutes = types.NumberNull()
        }
    } else {
        data.RescanIntervalInMinutes = types.NumberNull()
    }
    if obj, ok := item["nextScanAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.NextScanAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.NextScanAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.NextScanAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.NextScanAt = types.StringValue(string(jsonBytes))
        } else {
            data.NextScanAt = types.StringNull()
        }
    } else if val, ok := item["nextScanAt"].(string); ok {
        data.NextScanAt = types.StringValue(val)
    } else {
        data.NextScanAt = types.StringNull()
    }
    if obj, ok := item["autoImportProcessedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AutoImportProcessedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AutoImportProcessedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AutoImportProcessedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AutoImportProcessedAt = types.StringValue(string(jsonBytes))
        } else {
            data.AutoImportProcessedAt = types.StringNull()
        }
    } else if val, ok := item["autoImportProcessedAt"].(string); ok {
        data.AutoImportProcessedAt = types.StringValue(val)
    } else {
        data.AutoImportProcessedAt = types.StringNull()
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
