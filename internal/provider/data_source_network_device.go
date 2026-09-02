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
var _ datasource.DataSource = &NetworkDeviceDataSource{}

func NewNetworkDeviceDataSource() datasource.DataSource {
    return &NetworkDeviceDataSource{}
}

// NetworkDeviceDataSource defines the data source implementation.
type NetworkDeviceDataSource struct {
    client *Client
}

// NetworkDeviceDataSourceModel describes the data source data model.
type NetworkDeviceDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    Slug types.String `tfsdk:"slug"`
    Description types.String `tfsdk:"description"`
    Hostname types.String `tfsdk:"hostname"`
    ProbeId types.String `tfsdk:"probe_id"`
    SiteId types.String `tfsdk:"site_id"`
    OidTemplateId types.String `tfsdk:"oid_template_id"`
    CurrentMonitorStatusId types.String `tfsdk:"current_monitor_status_id"`
    MonitoringMethod types.String `tfsdk:"monitoring_method"`
    DeviceRole types.String `tfsdk:"device_role"`
    MonitorId types.String `tfsdk:"monitor_id"`
    SnmpVersion types.String `tfsdk:"snmp_version"`
    SnmpCommunityString types.String `tfsdk:"snmp_community_string"`
    SnmpPort types.Number `tfsdk:"snmp_port"`
    SnmpV3Auth types.String `tfsdk:"snmp_v3_auth"`
    SnmpV3SecurityLevel types.String `tfsdk:"snmp_v3_security_level"`
    SnmpV3Username types.String `tfsdk:"snmp_v3_username"`
    SnmpV3AuthProtocol types.String `tfsdk:"snmp_v3_auth_protocol"`
    SnmpV3AuthKey types.String `tfsdk:"snmp_v3_auth_key"`
    SnmpV3PrivProtocol types.String `tfsdk:"snmp_v3_priv_protocol"`
    SnmpV3PrivKey types.String `tfsdk:"snmp_v3_priv_key"`
    IsPollingEnabled types.Bool `tfsdk:"is_polling_enabled"`
    PollingIntervalInMinutes types.Number `tfsdk:"polling_interval_in_minutes"`
    WalkInterfaces types.Bool `tfsdk:"walk_interfaces"`
    CollectEndpoints types.Bool `tfsdk:"collect_endpoints"`
    SnmpOids types.String `tfsdk:"snmp_oids"`
    AutoApplyVendorHealthTemplate types.Bool `tfsdk:"auto_apply_vendor_health_template"`
    NextPollAt types.String `tfsdk:"next_poll_at"`
    LastWalkLog types.String `tfsdk:"last_walk_log"`
    SysDescr types.String `tfsdk:"sys_descr"`
    SysName types.String `tfsdk:"sys_name"`
    SysObjectId types.String `tfsdk:"sys_object_id"`
    SysLocation types.String `tfsdk:"sys_location"`
    SysContact types.String `tfsdk:"sys_contact"`
    Vendor types.String `tfsdk:"vendor"`
    DeviceModel types.String `tfsdk:"device_model"`
    SerialNumber types.String `tfsdk:"serial_number"`
    FirmwareVersion types.String `tfsdk:"firmware_version"`
    SoftwareVersion types.String `tfsdk:"software_version"`
    LastRebootedAt types.String `tfsdk:"last_rebooted_at"`
    CdpNeighbors types.String `tfsdk:"cdp_neighbors"`
    LldpNeighbors types.String `tfsdk:"lldp_neighbors"`
    LastSeenAt types.String `tfsdk:"last_seen_at"`
    LastPolledAt types.String `tfsdk:"last_polled_at"`
    IsReachable types.Bool `tfsdk:"is_reachable"`
    InterfacesTotal types.Number `tfsdk:"interfaces_total"`
    InterfacesUp types.Number `tfsdk:"interfaces_up"`
    InterfacesDown types.Number `tfsdk:"interfaces_down"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    IsArchived types.Bool `tfsdk:"is_archived"`
    ArchivedAt types.String `tfsdk:"archived_at"`
    ArchivedByUserId types.String `tfsdk:"archived_by_user_id"`
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
    Labels types.Set `tfsdk:"labels"`
}

func (d *NetworkDeviceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_network_device"
}

func (d *NetworkDeviceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Network Devices (routers, switches, firewalls) that are being monitored in this project via SNMP polling and traps. Look up an existing network_device by `id` or by `name`.",

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
                MarkdownDescription: "Friendly description for this network device.",
                Computed: true,
            },
            "hostname": schema.StringAttribute{
                MarkdownDescription: "IP address or hostname the probe polls; also matches SNMP trap sources.",
                Computed: true,
            },
            "probe_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "site_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "oid_template_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "current_monitor_status_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "monitoring_method": schema.StringAttribute{
                MarkdownDescription: "How this device's health is established: SNMP (an assigned probe walks it on a schedule) or Monitor (no polling — the linked monitor's status is the device's status). Devices created before this existed are SNMP..",
                Computed: true,
            },
            "device_role": schema.StringAttribute{
                MarkdownDescription: "What this device does on the network — router, switch, access point and so on. Left empty, the role is worked out from the device's own SNMP identity. Set it when there is no SNMP to read: a ping-only device has no identity to classify, and the role decides both the shape it is drawn with and where it sits in the topology hierarchy..",
                Computed: true,
            },
            "monitor_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "snmp_version": schema.StringAttribute{
                MarkdownDescription: "SNMP version to use when polling this device (V1, V2c, V3).",
                Computed: true,
            },
            "snmp_community_string": schema.StringAttribute{
                MarkdownDescription: "Community string used for SNMP v1/v2c polling.",
                Computed: true,
            },
            "snmp_port": schema.NumberAttribute{
                MarkdownDescription: "UDP port used for SNMP polling.",
                Computed: true,
            },
            "snmp_v3_auth": schema.StringAttribute{
                MarkdownDescription: "Deprecated: SNMP v3 auth is now stored in the snmpV3* columns below. Retained for reading legacy devices..",
                Computed: true,
            },
            "snmp_v3_security_level": schema.StringAttribute{
                MarkdownDescription: "SNMP v3 security level: noAuthNoPriv, authNoPriv, or authPriv.",
                Computed: true,
            },
            "snmp_v3_username": schema.StringAttribute{
                MarkdownDescription: "Security name (username) used for SNMP v3 polling.",
                Computed: true,
            },
            "snmp_v3_auth_protocol": schema.StringAttribute{
                MarkdownDescription: "SNMP v3 authentication protocol: MD5, SHA, SHA256, or SHA512.",
                Computed: true,
            },
            "snmp_v3_auth_key": schema.StringAttribute{
                MarkdownDescription: "SNMP v3 authentication passphrase.",
                Computed: true,
            },
            "snmp_v3_priv_protocol": schema.StringAttribute{
                MarkdownDescription: "SNMP v3 privacy (encryption) protocol: DES, AES, or AES256.",
                Computed: true,
            },
            "snmp_v3_priv_key": schema.StringAttribute{
                MarkdownDescription: "SNMP v3 privacy (encryption) passphrase.",
                Computed: true,
            },
            "is_polling_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether the assigned probe polls this device on a schedule. Disable to pause SNMP polling without deleting the device..",
                Computed: true,
            },
            "polling_interval_in_minutes": schema.NumberAttribute{
                MarkdownDescription: "How often, in minutes, the assigned probe polls this device via SNMP.",
                Computed: true,
            },
            "walk_interfaces": schema.BoolAttribute{
                MarkdownDescription: "Walk the IF-MIB interface tables on each poll to inventory interfaces, bandwidth, and errors. Also collects LLDP/CDP neighbors for the topology graph..",
                Computed: true,
            },
            "collect_endpoints": schema.BoolAttribute{
                MarkdownDescription: "Also walk the device's ARP cache and bridge forwarding database on each poll to discover endpoints (laptops, printers, POS terminals) attached to it. Strictly opt-in: costs extra SNMP table walks per poll. Only meaningful when Walk Interfaces is on..",
                Computed: true,
            },
            "snmp_oids": schema.StringAttribute{
                MarkdownDescription: "SNMP OIDs collected on each poll for this device ALONE, on top of whatever its OID Collection Template collects. Values are recorded as metrics and can be alerted on through monitor criteria. If several devices need the same OID, put it on a template instead..",
                Computed: true,
            },
            "auto_apply_vendor_health_template": schema.BoolAttribute{
                MarkdownDescription: "When the device's vendor is fingerprinted from its SNMP sysObjectID and no Health OIDs are configured yet, apply the matching vendor health template automatically on the next poll. Off by default for hand-made devices — the vendor template banner stays the manual path; auto-imported devices enable it so the zero-touch pipeline ends with health metrics, not an empty OID list..",
                Computed: true,
            },
            "next_poll_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "last_walk_log": schema.StringAttribute{
                MarkdownDescription: "The previous poll's interface counters. Kept so interface rates (bandwidth, utilization, errors/sec) can be computed as counter deltas between polls, and stores nothing else - the rest of the walk response has no reader and this column is rewritten on every poll of every device. Managed by the server..",
                Computed: true,
            },
            "sys_descr": schema.StringAttribute{
                MarkdownDescription: "System description (sysDescr) enriched from SNMP walks of this device.",
                Computed: true,
            },
            "sys_name": schema.StringAttribute{
                MarkdownDescription: "System name (sysName) enriched from SNMP walks of this device.",
                Computed: true,
            },
            "sys_object_id": schema.StringAttribute{
                MarkdownDescription: "sysObjectID — the vendor's registered OID for this device model, enriched from SNMP walks. Used to fingerprint the vendor and suggest an OID template..",
                Computed: true,
            },
            "sys_location": schema.StringAttribute{
                MarkdownDescription: "System location (sysLocation) enriched from SNMP walks of this device.",
                Computed: true,
            },
            "sys_contact": schema.StringAttribute{
                MarkdownDescription: "System contact (sysContact) enriched from SNMP walks of this device.",
                Computed: true,
            },
            "vendor": schema.StringAttribute{
                MarkdownDescription: "Hardware vendor, from ENTITY-MIB or derived from sysObjectID. Managed by the probe..",
                Computed: true,
            },
            "device_model": schema.StringAttribute{
                MarkdownDescription: "Hardware model from ENTITY-MIB (entPhysicalModelName). Managed by the probe..",
                Computed: true,
            },
            "serial_number": schema.StringAttribute{
                MarkdownDescription: "Chassis serial number from ENTITY-MIB (entPhysicalSerialNum). Managed by the probe..",
                Computed: true,
            },
            "firmware_version": schema.StringAttribute{
                MarkdownDescription: "Firmware revision from ENTITY-MIB (entPhysicalFirmwareRev). Managed by the probe..",
                Computed: true,
            },
            "software_version": schema.StringAttribute{
                MarkdownDescription: "Operating system / software revision from ENTITY-MIB (entPhysicalSoftwareRev). Managed by the probe..",
                Computed: true,
            },
            "last_rebooted_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "cdp_neighbors": schema.StringAttribute{
                MarkdownDescription: "CDP neighbors discovered on the last SNMP walk, complementing LLDP for the topology graph. Managed by the probe..",
                Computed: true,
            },
            "lldp_neighbors": schema.StringAttribute{
                MarkdownDescription: "LLDP neighbors discovered on the last SNMP walk, used to build the network topology graph. Managed by the probe..",
                Computed: true,
            },
            "last_seen_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "last_polled_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "is_reachable": schema.BoolAttribute{
                MarkdownDescription: "Whether the most recent SNMP walk reached this device. NULL means it has never been polled. This — not the age of lastSeenAt — is what the device list, the topology graph and the site rollup read, so a device whose last poll succeeded is never shown as down just because the probe is behind schedule. Managed by the probe..",
                Computed: true,
            },
            "interfaces_total": schema.NumberAttribute{
                MarkdownDescription: "Cached total count of interfaces on this device.",
                Computed: true,
            },
            "interfaces_up": schema.NumberAttribute{
                MarkdownDescription: "Cached count of operationally up interfaces on this device.",
                Computed: true,
            },
            "interfaces_down": schema.NumberAttribute{
                MarkdownDescription: "Cached count of operationally down interfaces on this device.",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "is_archived": schema.BoolAttribute{
                MarkdownDescription: "Is this network device archived? Archived network devices are hidden from lists but keep collecting telemetry..",
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
        },
    }
}

func (d *NetworkDeviceDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *NetworkDeviceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data NetworkDeviceDataSourceModel

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
            "Exactly one of `id` or `name` must be set to look up a network_device.",
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
        "hostname": true,
        "probeId": true,
        "siteId": true,
        "oidTemplateId": true,
        "currentMonitorStatusId": true,
        "monitoringMethod": true,
        "deviceRole": true,
        "monitorId": true,
        "snmpVersion": true,
        "snmpCommunityString": true,
        "snmpPort": true,
        "snmpV3Auth": true,
        "snmpV3SecurityLevel": true,
        "snmpV3Username": true,
        "snmpV3AuthProtocol": true,
        "snmpV3AuthKey": true,
        "snmpV3PrivProtocol": true,
        "snmpV3PrivKey": true,
        "isPollingEnabled": true,
        "pollingIntervalInMinutes": true,
        "walkInterfaces": true,
        "collectEndpoints": true,
        "snmpOids": true,
        "autoApplyVendorHealthTemplate": true,
        "nextPollAt": true,
        "lastWalkLog": true,
        "sysDescr": true,
        "sysName": true,
        "sysObjectId": true,
        "sysLocation": true,
        "sysContact": true,
        "vendor": true,
        "deviceModel": true,
        "serialNumber": true,
        "firmwareVersion": true,
        "softwareVersion": true,
        "lastRebootedAt": true,
        "cdpNeighbors": true,
        "lldpNeighbors": true,
        "lastSeenAt": true,
        "lastPolledAt": true,
        "isReachable": true,
        "interfacesTotal": true,
        "interfacesUp": true,
        "interfacesDown": true,
        "createdByUserId": true,
        "isArchived": true,
        "archivedAt": true,
        "archivedByUserId": true,
        "deletedByUserId": true,
        "labels": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/network-device/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read network_device, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No network_device found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read network_device: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/network-device/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list network_device, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list network_device: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No network_device found with name %q.", data.Name.ValueString()))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one network_device matches name %q. Use the id attribute to disambiguate.", data.Name.ValueString()))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for network_device.")
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
    if obj, ok := item["hostname"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Hostname = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Hostname = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Hostname = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Hostname = types.StringValue(string(jsonBytes))
        } else {
            data.Hostname = types.StringNull()
        }
    } else if val, ok := item["hostname"].(string); ok {
        data.Hostname = types.StringValue(val)
    } else {
        data.Hostname = types.StringNull()
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
    if obj, ok := item["siteId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SiteId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SiteId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SiteId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SiteId = types.StringValue(string(jsonBytes))
        } else {
            data.SiteId = types.StringNull()
        }
    } else if val, ok := item["siteId"].(string); ok {
        data.SiteId = types.StringValue(val)
    } else {
        data.SiteId = types.StringNull()
    }
    if obj, ok := item["oidTemplateId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OidTemplateId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.OidTemplateId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.OidTemplateId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.OidTemplateId = types.StringValue(string(jsonBytes))
        } else {
            data.OidTemplateId = types.StringNull()
        }
    } else if val, ok := item["oidTemplateId"].(string); ok {
        data.OidTemplateId = types.StringValue(val)
    } else {
        data.OidTemplateId = types.StringNull()
    }
    if obj, ok := item["currentMonitorStatusId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CurrentMonitorStatusId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CurrentMonitorStatusId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CurrentMonitorStatusId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CurrentMonitorStatusId = types.StringValue(string(jsonBytes))
        } else {
            data.CurrentMonitorStatusId = types.StringNull()
        }
    } else if val, ok := item["currentMonitorStatusId"].(string); ok {
        data.CurrentMonitorStatusId = types.StringValue(val)
    } else {
        data.CurrentMonitorStatusId = types.StringNull()
    }
    if obj, ok := item["monitoringMethod"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MonitoringMethod = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.MonitoringMethod = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.MonitoringMethod = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.MonitoringMethod = types.StringValue(string(jsonBytes))
        } else {
            data.MonitoringMethod = types.StringNull()
        }
    } else if val, ok := item["monitoringMethod"].(string); ok {
        data.MonitoringMethod = types.StringValue(val)
    } else {
        data.MonitoringMethod = types.StringNull()
    }
    if obj, ok := item["deviceRole"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeviceRole = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DeviceRole = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DeviceRole = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DeviceRole = types.StringValue(string(jsonBytes))
        } else {
            data.DeviceRole = types.StringNull()
        }
    } else if val, ok := item["deviceRole"].(string); ok {
        data.DeviceRole = types.StringValue(val)
    } else {
        data.DeviceRole = types.StringNull()
    }
    if obj, ok := item["monitorId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MonitorId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.MonitorId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.MonitorId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.MonitorId = types.StringValue(string(jsonBytes))
        } else {
            data.MonitorId = types.StringNull()
        }
    } else if val, ok := item["monitorId"].(string); ok {
        data.MonitorId = types.StringValue(val)
    } else {
        data.MonitorId = types.StringNull()
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
    if obj, ok := item["snmpV3Auth"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SnmpV3Auth = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SnmpV3Auth = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SnmpV3Auth = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SnmpV3Auth = types.StringValue(string(jsonBytes))
        } else {
            data.SnmpV3Auth = types.StringNull()
        }
    } else if val, ok := item["snmpV3Auth"].(string); ok {
        data.SnmpV3Auth = types.StringValue(val)
    } else {
        data.SnmpV3Auth = types.StringNull()
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
    if val, ok := item["isPollingEnabled"].(bool); ok {
        data.IsPollingEnabled = types.BoolValue(val)
    } else {
        data.IsPollingEnabled = types.BoolNull()
    }
    if val, ok := item["pollingIntervalInMinutes"].(float64); ok {
        data.PollingIntervalInMinutes = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["pollingIntervalInMinutes"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.PollingIntervalInMinutes = types.NumberValue(big.NewFloat(val))
        } else {
            data.PollingIntervalInMinutes = types.NumberNull()
        }
    } else {
        data.PollingIntervalInMinutes = types.NumberNull()
    }
    if val, ok := item["walkInterfaces"].(bool); ok {
        data.WalkInterfaces = types.BoolValue(val)
    } else {
        data.WalkInterfaces = types.BoolNull()
    }
    if val, ok := item["collectEndpoints"].(bool); ok {
        data.CollectEndpoints = types.BoolValue(val)
    } else {
        data.CollectEndpoints = types.BoolNull()
    }
    if obj, ok := item["snmpOids"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SnmpOids = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SnmpOids = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SnmpOids = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SnmpOids = types.StringValue(string(jsonBytes))
        } else {
            data.SnmpOids = types.StringNull()
        }
    } else if val, ok := item["snmpOids"].(string); ok {
        data.SnmpOids = types.StringValue(val)
    } else {
        data.SnmpOids = types.StringNull()
    }
    if val, ok := item["autoApplyVendorHealthTemplate"].(bool); ok {
        data.AutoApplyVendorHealthTemplate = types.BoolValue(val)
    } else {
        data.AutoApplyVendorHealthTemplate = types.BoolNull()
    }
    if obj, ok := item["nextPollAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.NextPollAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.NextPollAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.NextPollAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.NextPollAt = types.StringValue(string(jsonBytes))
        } else {
            data.NextPollAt = types.StringNull()
        }
    } else if val, ok := item["nextPollAt"].(string); ok {
        data.NextPollAt = types.StringValue(val)
    } else {
        data.NextPollAt = types.StringNull()
    }
    if obj, ok := item["lastWalkLog"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastWalkLog = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastWalkLog = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastWalkLog = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastWalkLog = types.StringValue(string(jsonBytes))
        } else {
            data.LastWalkLog = types.StringNull()
        }
    } else if val, ok := item["lastWalkLog"].(string); ok {
        data.LastWalkLog = types.StringValue(val)
    } else {
        data.LastWalkLog = types.StringNull()
    }
    if obj, ok := item["sysDescr"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SysDescr = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SysDescr = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SysDescr = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SysDescr = types.StringValue(string(jsonBytes))
        } else {
            data.SysDescr = types.StringNull()
        }
    } else if val, ok := item["sysDescr"].(string); ok {
        data.SysDescr = types.StringValue(val)
    } else {
        data.SysDescr = types.StringNull()
    }
    if obj, ok := item["sysName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SysName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SysName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SysName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SysName = types.StringValue(string(jsonBytes))
        } else {
            data.SysName = types.StringNull()
        }
    } else if val, ok := item["sysName"].(string); ok {
        data.SysName = types.StringValue(val)
    } else {
        data.SysName = types.StringNull()
    }
    if obj, ok := item["sysObjectId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SysObjectId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SysObjectId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SysObjectId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SysObjectId = types.StringValue(string(jsonBytes))
        } else {
            data.SysObjectId = types.StringNull()
        }
    } else if val, ok := item["sysObjectId"].(string); ok {
        data.SysObjectId = types.StringValue(val)
    } else {
        data.SysObjectId = types.StringNull()
    }
    if obj, ok := item["sysLocation"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SysLocation = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SysLocation = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SysLocation = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SysLocation = types.StringValue(string(jsonBytes))
        } else {
            data.SysLocation = types.StringNull()
        }
    } else if val, ok := item["sysLocation"].(string); ok {
        data.SysLocation = types.StringValue(val)
    } else {
        data.SysLocation = types.StringNull()
    }
    if obj, ok := item["sysContact"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SysContact = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SysContact = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SysContact = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SysContact = types.StringValue(string(jsonBytes))
        } else {
            data.SysContact = types.StringNull()
        }
    } else if val, ok := item["sysContact"].(string); ok {
        data.SysContact = types.StringValue(val)
    } else {
        data.SysContact = types.StringNull()
    }
    if obj, ok := item["vendor"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Vendor = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Vendor = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Vendor = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Vendor = types.StringValue(string(jsonBytes))
        } else {
            data.Vendor = types.StringNull()
        }
    } else if val, ok := item["vendor"].(string); ok {
        data.Vendor = types.StringValue(val)
    } else {
        data.Vendor = types.StringNull()
    }
    if obj, ok := item["deviceModel"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeviceModel = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DeviceModel = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DeviceModel = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DeviceModel = types.StringValue(string(jsonBytes))
        } else {
            data.DeviceModel = types.StringNull()
        }
    } else if val, ok := item["deviceModel"].(string); ok {
        data.DeviceModel = types.StringValue(val)
    } else {
        data.DeviceModel = types.StringNull()
    }
    if obj, ok := item["serialNumber"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SerialNumber = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SerialNumber = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SerialNumber = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SerialNumber = types.StringValue(string(jsonBytes))
        } else {
            data.SerialNumber = types.StringNull()
        }
    } else if val, ok := item["serialNumber"].(string); ok {
        data.SerialNumber = types.StringValue(val)
    } else {
        data.SerialNumber = types.StringNull()
    }
    if obj, ok := item["firmwareVersion"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FirmwareVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.FirmwareVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.FirmwareVersion = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.FirmwareVersion = types.StringValue(string(jsonBytes))
        } else {
            data.FirmwareVersion = types.StringNull()
        }
    } else if val, ok := item["firmwareVersion"].(string); ok {
        data.FirmwareVersion = types.StringValue(val)
    } else {
        data.FirmwareVersion = types.StringNull()
    }
    if obj, ok := item["softwareVersion"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SoftwareVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SoftwareVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SoftwareVersion = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SoftwareVersion = types.StringValue(string(jsonBytes))
        } else {
            data.SoftwareVersion = types.StringNull()
        }
    } else if val, ok := item["softwareVersion"].(string); ok {
        data.SoftwareVersion = types.StringValue(val)
    } else {
        data.SoftwareVersion = types.StringNull()
    }
    if obj, ok := item["lastRebootedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastRebootedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastRebootedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastRebootedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastRebootedAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastRebootedAt = types.StringNull()
        }
    } else if val, ok := item["lastRebootedAt"].(string); ok {
        data.LastRebootedAt = types.StringValue(val)
    } else {
        data.LastRebootedAt = types.StringNull()
    }
    if obj, ok := item["cdpNeighbors"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CdpNeighbors = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CdpNeighbors = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CdpNeighbors = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CdpNeighbors = types.StringValue(string(jsonBytes))
        } else {
            data.CdpNeighbors = types.StringNull()
        }
    } else if val, ok := item["cdpNeighbors"].(string); ok {
        data.CdpNeighbors = types.StringValue(val)
    } else {
        data.CdpNeighbors = types.StringNull()
    }
    if obj, ok := item["lldpNeighbors"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LldpNeighbors = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LldpNeighbors = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LldpNeighbors = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LldpNeighbors = types.StringValue(string(jsonBytes))
        } else {
            data.LldpNeighbors = types.StringNull()
        }
    } else if val, ok := item["lldpNeighbors"].(string); ok {
        data.LldpNeighbors = types.StringValue(val)
    } else {
        data.LldpNeighbors = types.StringNull()
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
    if obj, ok := item["lastPolledAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastPolledAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastPolledAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastPolledAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastPolledAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastPolledAt = types.StringNull()
        }
    } else if val, ok := item["lastPolledAt"].(string); ok {
        data.LastPolledAt = types.StringValue(val)
    } else {
        data.LastPolledAt = types.StringNull()
    }
    if val, ok := item["isReachable"].(bool); ok {
        data.IsReachable = types.BoolValue(val)
    } else {
        data.IsReachable = types.BoolNull()
    }
    if val, ok := item["interfacesTotal"].(float64); ok {
        data.InterfacesTotal = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["interfacesTotal"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.InterfacesTotal = types.NumberValue(big.NewFloat(val))
        } else {
            data.InterfacesTotal = types.NumberNull()
        }
    } else {
        data.InterfacesTotal = types.NumberNull()
    }
    if val, ok := item["interfacesUp"].(float64); ok {
        data.InterfacesUp = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["interfacesUp"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.InterfacesUp = types.NumberValue(big.NewFloat(val))
        } else {
            data.InterfacesUp = types.NumberNull()
        }
    } else {
        data.InterfacesUp = types.NumberNull()
    }
    if val, ok := item["interfacesDown"].(float64); ok {
        data.InterfacesDown = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["interfacesDown"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.InterfacesDown = types.NumberValue(big.NewFloat(val))
        } else {
            data.InterfacesDown = types.NumberNull()
        }
    } else {
        data.InterfacesDown = types.NumberNull()
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

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
