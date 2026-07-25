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
var _ datasource.DataSource = &NetworkDeviceDataDataSource{}

func NewNetworkDeviceDataDataSource() datasource.DataSource {
    return &NetworkDeviceDataDataSource{}
}

// NetworkDeviceDataDataSource defines the data source implementation.
type NetworkDeviceDataDataSource struct {
    client *Client
}

// NetworkDeviceDataDataSourceModel describes the data source data model.
type NetworkDeviceDataDataSourceModel struct {
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
    CurrentMonitorStatusId types.String `tfsdk:"current_monitor_status_id"`
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

func (d *NetworkDeviceDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_network_device_data"
}

func (d *NetworkDeviceDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "network_device_data data source",

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
                MarkdownDescription: "Friendly globally unique name for your object. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Device], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Friendly description for this network device. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Network Device], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Device], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Network Device]",
                Computed: true,
            },
            "hostname": schema.StringAttribute{
                MarkdownDescription: "IP address or hostname the probe polls; also matches SNMP trap sources. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Network Device], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Device], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Network Device]",
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
            "current_monitor_status_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "snmp_version": schema.StringAttribute{
                MarkdownDescription: "SNMP version to use when polling this device (V1, V2c, V3). Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Network Device], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Device], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Network Device]",
                Computed: true,
            },
            "snmp_community_string": schema.StringAttribute{
                MarkdownDescription: "Community string used for SNMP v1/v2c polling. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Network Device], Read: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Read Network Device], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Network Device]",
                Computed: true,
            },
            "snmp_port": schema.NumberAttribute{
                MarkdownDescription: "UDP port used for SNMP polling. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Network Device], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Device], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Network Device]",
                Computed: true,
            },
            "snmp_v3_auth": schema.StringAttribute{
                MarkdownDescription: "Deprecated: SNMP v3 auth is now stored in the snmpV3* columns below. Retained for reading legacy devices.. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Network Device], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Device], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Network Device]",
                Computed: true,
            },
            "snmp_v3_security_level": schema.StringAttribute{
                MarkdownDescription: "SNMP v3 security level: noAuthNoPriv, authNoPriv, or authPriv. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Network Device], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Device], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Network Device]",
                Computed: true,
            },
            "snmp_v3_username": schema.StringAttribute{
                MarkdownDescription: "Security name (username) used for SNMP v3 polling. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Network Device], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Device], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Network Device]",
                Computed: true,
            },
            "snmp_v3_auth_protocol": schema.StringAttribute{
                MarkdownDescription: "SNMP v3 authentication protocol: MD5, SHA, SHA256, or SHA512. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Network Device], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Device], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Network Device]",
                Computed: true,
            },
            "snmp_v3_auth_key": schema.StringAttribute{
                MarkdownDescription: "SNMP v3 authentication passphrase. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Network Device], Read: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Read Network Device], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Network Device]",
                Computed: true,
            },
            "snmp_v3_priv_protocol": schema.StringAttribute{
                MarkdownDescription: "SNMP v3 privacy (encryption) protocol: DES, AES, or AES256. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Network Device], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Device], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Network Device]",
                Computed: true,
            },
            "snmp_v3_priv_key": schema.StringAttribute{
                MarkdownDescription: "SNMP v3 privacy (encryption) passphrase. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Network Device], Read: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Read Network Device], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Network Device]",
                Computed: true,
            },
            "is_polling_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether the assigned probe polls this device on a schedule. Disable to pause SNMP polling without deleting the device.. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Network Device], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Device], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Network Device]",
                Computed: true,
            },
            "polling_interval_in_minutes": schema.NumberAttribute{
                MarkdownDescription: "How often, in minutes, the assigned probe polls this device via SNMP. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Network Device], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Device], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Network Device]",
                Computed: true,
            },
            "walk_interfaces": schema.BoolAttribute{
                MarkdownDescription: "Walk the IF-MIB interface tables on each poll to inventory interfaces, bandwidth, and errors. Also collects LLDP/CDP neighbors for the topology graph.. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Network Device], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Device], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Network Device]",
                Computed: true,
            },
            "collect_endpoints": schema.BoolAttribute{
                MarkdownDescription: "Also walk the device's ARP cache and bridge forwarding database on each poll to discover endpoints (laptops, printers, POS terminals) attached to it. Strictly opt-in: costs extra SNMP table walks per poll. Only meaningful when Walk Interfaces is on.. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Network Device], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Device], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Network Device]",
                Computed: true,
            },
            "snmp_oids": schema.StringAttribute{
                MarkdownDescription: "SNMP OIDs (CPU, memory, temperature, or any custom OID) collected on each poll. Values are recorded as metrics and can be alerted on through monitor criteria.. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Network Device], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Device], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Network Device]",
                Computed: true,
            },
            "next_poll_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "last_walk_log": schema.StringAttribute{
                MarkdownDescription: "The previous poll's raw walk response. Kept so interface rates (bandwidth, utilization, errors/sec) can be computed as counter deltas between polls. Managed by the server.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Read Network Device], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "sys_descr": schema.StringAttribute{
                MarkdownDescription: "System description (sysDescr) enriched from SNMP walks of this device. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Device], Update: [Project Owner, Project Admin, Edit Network Device]",
                Computed: true,
            },
            "sys_name": schema.StringAttribute{
                MarkdownDescription: "System name (sysName) enriched from SNMP walks of this device. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Device], Update: [Project Owner, Project Admin, Edit Network Device]",
                Computed: true,
            },
            "sys_object_id": schema.StringAttribute{
                MarkdownDescription: "sysObjectID — the vendor's registered OID for this device model, enriched from SNMP walks. Used to fingerprint the vendor and suggest an OID template.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Device], Update: [Project Owner, Project Admin, Edit Network Device]",
                Computed: true,
            },
            "sys_location": schema.StringAttribute{
                MarkdownDescription: "System location (sysLocation) enriched from SNMP walks of this device. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Device], Update: [Project Owner, Project Admin, Edit Network Device]",
                Computed: true,
            },
            "sys_contact": schema.StringAttribute{
                MarkdownDescription: "System contact (sysContact) enriched from SNMP walks of this device. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Device], Update: [Project Owner, Project Admin, Edit Network Device]",
                Computed: true,
            },
            "vendor": schema.StringAttribute{
                MarkdownDescription: "Hardware vendor, from ENTITY-MIB or derived from sysObjectID. Managed by the probe.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Device], Update: [Project Owner, Project Admin, Edit Network Device]",
                Computed: true,
            },
            "device_model": schema.StringAttribute{
                MarkdownDescription: "Hardware model from ENTITY-MIB (entPhysicalModelName). Managed by the probe.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Device], Update: [Project Owner, Project Admin, Edit Network Device]",
                Computed: true,
            },
            "serial_number": schema.StringAttribute{
                MarkdownDescription: "Chassis serial number from ENTITY-MIB (entPhysicalSerialNum). Managed by the probe.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Device], Update: [Project Owner, Project Admin, Edit Network Device]",
                Computed: true,
            },
            "firmware_version": schema.StringAttribute{
                MarkdownDescription: "Firmware revision from ENTITY-MIB (entPhysicalFirmwareRev). Managed by the probe.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Device], Update: [Project Owner, Project Admin, Edit Network Device]",
                Computed: true,
            },
            "software_version": schema.StringAttribute{
                MarkdownDescription: "Operating system / software revision from ENTITY-MIB (entPhysicalSoftwareRev). Managed by the probe.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Device], Update: [Project Owner, Project Admin, Edit Network Device]",
                Computed: true,
            },
            "last_rebooted_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "cdp_neighbors": schema.StringAttribute{
                MarkdownDescription: "CDP neighbors discovered on the last SNMP walk, complementing LLDP for the topology graph. Managed by the probe.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Device], Update: [Project Owner, Project Admin, Edit Network Device]",
                Computed: true,
            },
            "lldp_neighbors": schema.StringAttribute{
                MarkdownDescription: "LLDP neighbors discovered on the last SNMP walk, used to build the network topology graph. Managed by the probe.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Device], Update: [Project Owner, Project Admin, Edit Network Device]",
                Computed: true,
            },
            "last_seen_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "interfaces_total": schema.NumberAttribute{
                MarkdownDescription: "Cached total count of interfaces on this device. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Device], Update: [Project Owner, Project Admin, Edit Network Device]",
                Computed: true,
            },
            "interfaces_up": schema.NumberAttribute{
                MarkdownDescription: "Cached count of operationally up interfaces on this device. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Device], Update: [Project Owner, Project Admin, Edit Network Device]",
                Computed: true,
            },
            "interfaces_down": schema.NumberAttribute{
                MarkdownDescription: "Cached count of operationally down interfaces on this device. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Device], Update: [Project Owner, Project Admin, Edit Network Device]",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "is_archived": schema.BoolAttribute{
                MarkdownDescription: "Is this network device archived? Archived network devices are hidden from lists but keep collecting telemetry.. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Network Device], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Device], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Network Device]",
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
                MarkdownDescription: "Relation to Labels Array where this object is categorized in.. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Network Device], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Device], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Network Device]",
                Computed: true,
                ElementType: types.StringType,
            },
        },
    }
}

func (d *NetworkDeviceDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *NetworkDeviceDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data NetworkDeviceDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "network-device" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read network_device_data, got error: %s", err))
        return
    }

    var networkDeviceDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &networkDeviceDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse network_device_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := networkDeviceDataResponse["data"].(map[string]interface{}); ok {
        networkDeviceDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := networkDeviceDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := networkDeviceDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := networkDeviceDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := networkDeviceDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := networkDeviceDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := networkDeviceDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := networkDeviceDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := networkDeviceDataResponse["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    }
    if val, ok := networkDeviceDataResponse["description"].(string); ok {
        data.Description = types.StringValue(val)
    }
    if val, ok := networkDeviceDataResponse["hostname"].(string); ok {
        data.Hostname = types.StringValue(val)
    }
    if val, ok := networkDeviceDataResponse["probe_id"].(string); ok {
        data.ProbeId = types.StringValue(val)
    }
    if val, ok := networkDeviceDataResponse["site_id"].(string); ok {
        data.SiteId = types.StringValue(val)
    }
    if val, ok := networkDeviceDataResponse["current_monitor_status_id"].(string); ok {
        data.CurrentMonitorStatusId = types.StringValue(val)
    }
    if val, ok := networkDeviceDataResponse["snmp_version"].(string); ok {
        data.SnmpVersion = types.StringValue(val)
    }
    if val, ok := networkDeviceDataResponse["snmp_community_string"].(string); ok {
        data.SnmpCommunityString = types.StringValue(val)
    }
    if val, ok := networkDeviceDataResponse["snmp_port"].(float64); ok {
        data.SnmpPort = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := networkDeviceDataResponse["snmp_v3_auth"].(string); ok {
        data.SnmpV3Auth = types.StringValue(val)
    }
    if val, ok := networkDeviceDataResponse["snmp_v3_security_level"].(string); ok {
        data.SnmpV3SecurityLevel = types.StringValue(val)
    }
    if val, ok := networkDeviceDataResponse["snmp_v3_username"].(string); ok {
        data.SnmpV3Username = types.StringValue(val)
    }
    if val, ok := networkDeviceDataResponse["snmp_v3_auth_protocol"].(string); ok {
        data.SnmpV3AuthProtocol = types.StringValue(val)
    }
    if val, ok := networkDeviceDataResponse["snmp_v3_auth_key"].(string); ok {
        data.SnmpV3AuthKey = types.StringValue(val)
    }
    if val, ok := networkDeviceDataResponse["snmp_v3_priv_protocol"].(string); ok {
        data.SnmpV3PrivProtocol = types.StringValue(val)
    }
    if val, ok := networkDeviceDataResponse["snmp_v3_priv_key"].(string); ok {
        data.SnmpV3PrivKey = types.StringValue(val)
    }
    if val, ok := networkDeviceDataResponse["is_polling_enabled"].(bool); ok {
        data.IsPollingEnabled = types.BoolValue(val)
    }
    if val, ok := networkDeviceDataResponse["polling_interval_in_minutes"].(float64); ok {
        data.PollingIntervalInMinutes = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := networkDeviceDataResponse["walk_interfaces"].(bool); ok {
        data.WalkInterfaces = types.BoolValue(val)
    }
    if val, ok := networkDeviceDataResponse["collect_endpoints"].(bool); ok {
        data.CollectEndpoints = types.BoolValue(val)
    }
    if val, ok := networkDeviceDataResponse["snmp_oids"].(string); ok {
        data.SnmpOids = types.StringValue(val)
    }
    if val, ok := networkDeviceDataResponse["next_poll_at"].(string); ok {
        data.NextPollAt = types.StringValue(val)
    }
    if val, ok := networkDeviceDataResponse["last_walk_log"].(string); ok {
        data.LastWalkLog = types.StringValue(val)
    }
    if val, ok := networkDeviceDataResponse["sys_descr"].(string); ok {
        data.SysDescr = types.StringValue(val)
    }
    if val, ok := networkDeviceDataResponse["sys_name"].(string); ok {
        data.SysName = types.StringValue(val)
    }
    if val, ok := networkDeviceDataResponse["sys_object_id"].(string); ok {
        data.SysObjectId = types.StringValue(val)
    }
    if val, ok := networkDeviceDataResponse["sys_location"].(string); ok {
        data.SysLocation = types.StringValue(val)
    }
    if val, ok := networkDeviceDataResponse["sys_contact"].(string); ok {
        data.SysContact = types.StringValue(val)
    }
    if val, ok := networkDeviceDataResponse["vendor"].(string); ok {
        data.Vendor = types.StringValue(val)
    }
    if val, ok := networkDeviceDataResponse["device_model"].(string); ok {
        data.DeviceModel = types.StringValue(val)
    }
    if val, ok := networkDeviceDataResponse["serial_number"].(string); ok {
        data.SerialNumber = types.StringValue(val)
    }
    if val, ok := networkDeviceDataResponse["firmware_version"].(string); ok {
        data.FirmwareVersion = types.StringValue(val)
    }
    if val, ok := networkDeviceDataResponse["software_version"].(string); ok {
        data.SoftwareVersion = types.StringValue(val)
    }
    if val, ok := networkDeviceDataResponse["last_rebooted_at"].(string); ok {
        data.LastRebootedAt = types.StringValue(val)
    }
    if val, ok := networkDeviceDataResponse["cdp_neighbors"].(string); ok {
        data.CdpNeighbors = types.StringValue(val)
    }
    if val, ok := networkDeviceDataResponse["lldp_neighbors"].(string); ok {
        data.LldpNeighbors = types.StringValue(val)
    }
    if val, ok := networkDeviceDataResponse["last_seen_at"].(string); ok {
        data.LastSeenAt = types.StringValue(val)
    }
    if val, ok := networkDeviceDataResponse["interfaces_total"].(float64); ok {
        data.InterfacesTotal = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := networkDeviceDataResponse["interfaces_up"].(float64); ok {
        data.InterfacesUp = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := networkDeviceDataResponse["interfaces_down"].(float64); ok {
        data.InterfacesDown = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := networkDeviceDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := networkDeviceDataResponse["is_archived"].(bool); ok {
        data.IsArchived = types.BoolValue(val)
    }
    if val, ok := networkDeviceDataResponse["archived_at"].(string); ok {
        data.ArchivedAt = types.StringValue(val)
    }
    if val, ok := networkDeviceDataResponse["archived_by_user_id"].(string); ok {
        data.ArchivedByUserId = types.StringValue(val)
    }
    if val, ok := networkDeviceDataResponse["deleted_by_user_id"].(string); ok {
        data.DeletedByUserId = types.StringValue(val)
    }
    if val, ok := networkDeviceDataResponse["labels"].([]interface{}); ok {
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

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
