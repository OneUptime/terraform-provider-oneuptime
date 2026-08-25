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
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    ProbeId types.String `tfsdk:"probe_id"`
    Cidr types.String `tfsdk:"cidr"`
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
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
}

func (d *NetworkDeviceDiscoveryScanDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_network_device_discovery_scan"
}

func (d *NetworkDeviceDiscoveryScanDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Network discovery scans that sweep an address space — a CIDR subnet or an octet range — via SNMP from a probe and report devices found, so they can be imported as Network Devices. Look up an existing network_device_discovery_scan by `id` or by `name`.",

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
            "probe_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "cidr": schema.StringAttribute{
                MarkdownDescription: "Address space to scan, either in CIDR notation (192.168.1.0/24) or octet-range notation where any octet may be an inclusive low-high range (10.16-22.0-255.51-66).",
                Computed: true,
            },
            "snmp_version": schema.StringAttribute{
                MarkdownDescription: "SNMP version tried against every host in the subnet (V1, V2c, V3).",
                Computed: true,
            },
            "snmp_community_string": schema.StringAttribute{
                MarkdownDescription: "Community string tried against every host in the subnet (SNMP v1/v2c).",
                Computed: true,
            },
            "snmp_port": schema.NumberAttribute{
                MarkdownDescription: "UDP port tried against every host in the subnet.",
                Computed: true,
            },
            "snmp_v3_security_level": schema.StringAttribute{
                MarkdownDescription: "SNMP v3 security level tried against every host: noAuthNoPriv, authNoPriv, or authPriv.",
                Computed: true,
            },
            "snmp_v3_username": schema.StringAttribute{
                MarkdownDescription: "SNMP v3 security name (username) tried against every host.",
                Computed: true,
            },
            "snmp_v3_auth_protocol": schema.StringAttribute{
                MarkdownDescription: "SNMP v3 authentication protocol: MD5, SHA, SHA256, or SHA512.",
                Computed: true,
            },
            "snmp_v3_auth_key": schema.StringAttribute{
                MarkdownDescription: "SNMP v3 authentication passphrase tried against every host.",
                Computed: true,
            },
            "snmp_v3_priv_protocol": schema.StringAttribute{
                MarkdownDescription: "SNMP v3 privacy (encryption) protocol: DES, AES, or AES256.",
                Computed: true,
            },
            "snmp_v3_priv_key": schema.StringAttribute{
                MarkdownDescription: "SNMP v3 privacy (encryption) passphrase tried against every host.",
                Computed: true,
            },
            "status": schema.StringAttribute{
                MarkdownDescription: "Status of this discovery scan: \"Pending\", \"In Progress\", \"Completed\" or \"Failed\". Managed by the scanning probe..",
                Computed: true,
            },
            "status_message": schema.StringAttribute{
                MarkdownDescription: "Details about the current status of this scan, e.g. the failure reason. Managed by the scanning probe..",
                Computed: true,
            },
            "discovered_devices": schema.StringAttribute{
                MarkdownDescription: "Devices found by this scan: array of {ipAddress, sysName, sysDescr, isAlreadyRegistered}. Managed by the scanning probe..",
                Computed: true,
            },
            "scanned_host_count": schema.NumberAttribute{
                MarkdownDescription: "Total number of host addresses swept in the subnet. Managed by the scanning probe..",
                Computed: true,
            },
            "responded_host_count": schema.NumberAttribute{
                MarkdownDescription: "Number of hosts that responded to SNMP during the sweep. Managed by the scanning probe..",
                Computed: true,
            },
            "started_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "completed_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "is_recurring": schema.BoolAttribute{
                MarkdownDescription: "Re-run this scan automatically every Rescan Interval minutes to keep discovery continuous..",
                Computed: true,
            },
            "rescan_interval_in_minutes": schema.NumberAttribute{
                MarkdownDescription: "How often a recurring scan re-runs, in minutes. Ignored unless Is Recurring is on..",
                Computed: true,
            },
            "next_scan_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "auto_import_processed_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
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

    hasId := !data.Id.IsNull() && data.Id.ValueString() != ""
    hasName := !data.Name.IsNull() && data.Name.ValueString() != ""
    if hasId == hasName {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Exactly one of `id` or `name` must be set to look up a network_device_discovery_scan.",
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
        "probeId": true,
        "cidr": true,
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
        "deletedByUserId": true,
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
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No network_device_discovery_scan found with id %q.", data.Id.ValueString()))
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
    } else {
        listBody := map[string]interface{}{
            "query": map[string]interface{}{
                "name": data.Name.ValueString(),
            },
            "select": selectParam,
            // limit 2 is enough to detect ambiguity without paging.
            "limit": 2,
        }
        httpResp, err := d.client.Post(ctx, "/network-device-discovery-scan/get-list", listBody)
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
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No network_device_discovery_scan found with name %q.", data.Name.ValueString()))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one network_device_discovery_scan matches name %q. Use the id attribute to disambiguate.", data.Name.ValueString()))
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

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
