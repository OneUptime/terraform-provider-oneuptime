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
var _ resource.Resource = &NetworkDeviceResource{}
var _ resource.ResourceWithImportState = &NetworkDeviceResource{}

func NewNetworkDeviceResource() resource.Resource {
    return &NetworkDeviceResource{}
}

// NetworkDeviceResource defines the resource implementation.
type NetworkDeviceResource struct {
    client *Client
}

// NetworkDeviceResourceModel describes the resource data model.
type NetworkDeviceResourceModel struct {
    Id types.String `tfsdk:"id"`
    ProjectId types.String `tfsdk:"project_id"`
    Name types.String `tfsdk:"name"`
    Description types.String `tfsdk:"description"`
    Hostname types.String `tfsdk:"hostname"`
    ProbeId types.String `tfsdk:"probe_id"`
    SiteId types.String `tfsdk:"site_id"`
    MonitoringMethod types.String `tfsdk:"monitoring_method"`
    DeviceRole types.String `tfsdk:"device_role"`
    MonitorId types.String `tfsdk:"monitor_id"`
    SnmpVersion types.String `tfsdk:"snmp_version"`
    SnmpCommunityString types.String `tfsdk:"snmp_community_string"`
    SnmpPort types.Number `tfsdk:"snmp_port"`
    SnmpV3Auth JSONSubsetValue `tfsdk:"snmp_v3_auth"`
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
    SnmpOids JSONSubsetValue `tfsdk:"snmp_oids"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    IsArchived types.Bool `tfsdk:"is_archived"`
    Labels types.Set `tfsdk:"labels"`
    CurrentMonitorStatusId types.String `tfsdk:"current_monitor_status_id"`
    NextPollAt RFC3339Value `tfsdk:"next_poll_at"`
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
    LastRebootedAt RFC3339Value `tfsdk:"last_rebooted_at"`
    CdpNeighbors JSONSubsetValue `tfsdk:"cdp_neighbors"`
    LldpNeighbors JSONSubsetValue `tfsdk:"lldp_neighbors"`
    LastSeenAt RFC3339Value `tfsdk:"last_seen_at"`
    InterfacesTotal types.Number `tfsdk:"interfaces_total"`
    InterfacesUp types.Number `tfsdk:"interfaces_up"`
    InterfacesDown types.Number `tfsdk:"interfaces_down"`
    CreatedAt RFC3339Value `tfsdk:"created_at"`
    UpdatedAt RFC3339Value `tfsdk:"updated_at"`
    DeletedAt RFC3339Value `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    Slug types.String `tfsdk:"slug"`
    LastWalkLog JSONSubsetValue `tfsdk:"last_walk_log"`
    ArchivedAt RFC3339Value `tfsdk:"archived_at"`
    ArchivedByUserId types.String `tfsdk:"archived_by_user_id"`
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
}

func (r *NetworkDeviceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_network_device"
}

func (r *NetworkDeviceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Network Devices (routers, switches, firewalls) that are being monitored in this project via SNMP polling and traps.",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Unique identifier for the resource",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Friendly name for this network device.",
                Required: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Friendly description for this network device.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "hostname": schema.StringAttribute{
                MarkdownDescription: "IP address or hostname the probe polls; also matches SNMP trap sources.",
                Required: true,
            },
            "probe_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "site_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "monitoring_method": schema.StringAttribute{
                MarkdownDescription: "How this device's health is established: SNMP (an assigned probe walks it on a schedule) or Monitor (no polling — the linked monitor's status is the device's status). Devices created before this existed are SNMP..",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "device_role": schema.StringAttribute{
                MarkdownDescription: "What this device does on the network — router, switch, access point and so on. Left empty, the role is worked out from the device's own SNMP identity. Set it when there is no SNMP to read: a ping-only device has no identity to classify, and the role decides both the shape it is drawn with and where it sits in the topology hierarchy..",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "monitor_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "snmp_version": schema.StringAttribute{
                MarkdownDescription: "SNMP version to use when polling this device (V1, V2c, V3).",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "snmp_community_string": schema.StringAttribute{
                MarkdownDescription: "Community string used for SNMP v1/v2c polling.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "snmp_port": schema.NumberAttribute{
                MarkdownDescription: "UDP port used for SNMP polling.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "snmp_v3_auth": schema.StringAttribute{
                MarkdownDescription: "Deprecated: SNMP v3 auth is now stored in the snmpV3* columns below. Retained for reading legacy devices..",
                CustomType: JSONSubsetType{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: []validator.String{
                    JSONEnvelopeValidator(),
                },
            },
            "snmp_v3_security_level": schema.StringAttribute{
                MarkdownDescription: "SNMP v3 security level: noAuthNoPriv, authNoPriv, or authPriv.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "snmp_v3_username": schema.StringAttribute{
                MarkdownDescription: "Security name (username) used for SNMP v3 polling.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "snmp_v3_auth_protocol": schema.StringAttribute{
                MarkdownDescription: "SNMP v3 authentication protocol: MD5, SHA, SHA256, or SHA512.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "snmp_v3_auth_key": schema.StringAttribute{
                MarkdownDescription: "SNMP v3 authentication passphrase.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "snmp_v3_priv_protocol": schema.StringAttribute{
                MarkdownDescription: "SNMP v3 privacy (encryption) protocol: DES, AES, or AES256.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "snmp_v3_priv_key": schema.StringAttribute{
                MarkdownDescription: "SNMP v3 privacy (encryption) passphrase.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "is_polling_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether the assigned probe polls this device on a schedule. Disable to pause SNMP polling without deleting the device..",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(true),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "polling_interval_in_minutes": schema.NumberAttribute{
                MarkdownDescription: "How often, in minutes, the assigned probe polls this device via SNMP.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "walk_interfaces": schema.BoolAttribute{
                MarkdownDescription: "Walk the IF-MIB interface tables on each poll to inventory interfaces, bandwidth, and errors. Also collects LLDP/CDP neighbors for the topology graph..",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(true),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "collect_endpoints": schema.BoolAttribute{
                MarkdownDescription: "Also walk the device's ARP cache and bridge forwarding database on each poll to discover endpoints (laptops, printers, POS terminals) attached to it. Strictly opt-in: costs extra SNMP table walks per poll. Only meaningful when Walk Interfaces is on..",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "snmp_oids": schema.StringAttribute{
                MarkdownDescription: "SNMP OIDs (CPU, memory, temperature, or any custom OID) collected on each poll. Values are recorded as metrics and can be alerted on through monitor criteria..",
                CustomType: JSONSubsetType{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: []validator.String{
                    JSONEnvelopeValidator(),
                },
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                    stringplanmodifier.RequiresReplace(),
                },
            },
            "is_archived": schema.BoolAttribute{
                MarkdownDescription: "Is this network device archived? Archived network devices are hidden from lists but keep collecting telemetry..",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "labels": schema.SetAttribute{
                MarkdownDescription: "Relation to Labels Array where this object is categorized in..",
                Optional: true,
                Computed: true,
                ElementType: types.StringType,
                PlanModifiers: []planmodifier.Set{
                    setplanmodifier.UseStateForUnknown(),
                },
            },
            "current_monitor_status_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "next_poll_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                CustomType: RFC3339Type{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "sys_descr": schema.StringAttribute{
                MarkdownDescription: "System description (sysDescr) enriched from SNMP walks of this device.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "sys_name": schema.StringAttribute{
                MarkdownDescription: "System name (sysName) enriched from SNMP walks of this device.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "sys_object_id": schema.StringAttribute{
                MarkdownDescription: "sysObjectID — the vendor's registered OID for this device model, enriched from SNMP walks. Used to fingerprint the vendor and suggest an OID template..",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "sys_location": schema.StringAttribute{
                MarkdownDescription: "System location (sysLocation) enriched from SNMP walks of this device.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "sys_contact": schema.StringAttribute{
                MarkdownDescription: "System contact (sysContact) enriched from SNMP walks of this device.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "vendor": schema.StringAttribute{
                MarkdownDescription: "Hardware vendor, from ENTITY-MIB or derived from sysObjectID. Managed by the probe..",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "device_model": schema.StringAttribute{
                MarkdownDescription: "Hardware model from ENTITY-MIB (entPhysicalModelName). Managed by the probe..",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "serial_number": schema.StringAttribute{
                MarkdownDescription: "Chassis serial number from ENTITY-MIB (entPhysicalSerialNum). Managed by the probe..",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "firmware_version": schema.StringAttribute{
                MarkdownDescription: "Firmware revision from ENTITY-MIB (entPhysicalFirmwareRev). Managed by the probe..",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "software_version": schema.StringAttribute{
                MarkdownDescription: "Operating system / software revision from ENTITY-MIB (entPhysicalSoftwareRev). Managed by the probe..",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "last_rebooted_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                CustomType: RFC3339Type{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "cdp_neighbors": schema.StringAttribute{
                MarkdownDescription: "CDP neighbors discovered on the last SNMP walk, complementing LLDP for the topology graph. Managed by the probe..",
                CustomType: JSONSubsetType{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: []validator.String{
                    JSONEnvelopeValidator(),
                },
            },
            "lldp_neighbors": schema.StringAttribute{
                MarkdownDescription: "LLDP neighbors discovered on the last SNMP walk, used to build the network topology graph. Managed by the probe..",
                CustomType: JSONSubsetType{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: []validator.String{
                    JSONEnvelopeValidator(),
                },
            },
            "last_seen_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                CustomType: RFC3339Type{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "interfaces_total": schema.NumberAttribute{
                MarkdownDescription: "Cached total count of interfaces on this device.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "interfaces_up": schema.NumberAttribute{
                MarkdownDescription: "Cached count of operationally up interfaces on this device.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "interfaces_down": schema.NumberAttribute{
                MarkdownDescription: "Cached count of operationally down interfaces on this device.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "created_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "updated_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "deleted_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "version": schema.NumberAttribute{
                MarkdownDescription: "Object version",
                Computed: true,
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Friendly globally unique name for your object.",
                Computed: true,
            },
            "last_walk_log": schema.StringAttribute{
                MarkdownDescription: "The previous poll's raw walk response. Kept so interface rates (bandwidth, utilization, errors/sec) can be computed as counter deltas between polls. Managed by the server..",
                CustomType: JSONSubsetType{},
                Computed: true,
            },
            "archived_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                CustomType: RFC3339Type{},
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
        },
    }
}

func (r *NetworkDeviceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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


func (r *NetworkDeviceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var data NetworkDeviceResourceModel

    // Read Terraform plan data into the model
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }



    // Create API request body. Unset (null/unknown) optional fields are
    // omitted so server-side defaults apply instead of being overwritten
    // with zero values.
    networkDeviceRequest := map[string]interface{}{
        "data": map[string]interface{}{},
    }
    requestDataMap := networkDeviceRequest["data"].(map[string]interface{})

    if !data.Name.IsNull() && !data.Name.IsUnknown() {
        requestDataMap["name"] = data.Name.ValueString()
    }
    if !data.Description.IsNull() && !data.Description.IsUnknown() {
        requestDataMap["description"] = data.Description.ValueString()
    }
    if !data.Hostname.IsNull() && !data.Hostname.IsUnknown() {
        requestDataMap["hostname"] = data.Hostname.ValueString()
    }
    if !data.ProbeId.IsNull() && !data.ProbeId.IsUnknown() {
        requestDataMap["probeId"] = data.ProbeId.ValueString()
    }
    if !data.SiteId.IsNull() && !data.SiteId.IsUnknown() {
        requestDataMap["siteId"] = data.SiteId.ValueString()
    }
    if !data.MonitoringMethod.IsNull() && !data.MonitoringMethod.IsUnknown() {
        requestDataMap["monitoringMethod"] = data.MonitoringMethod.ValueString()
    }
    if !data.DeviceRole.IsNull() && !data.DeviceRole.IsUnknown() {
        requestDataMap["deviceRole"] = data.DeviceRole.ValueString()
    }
    if !data.MonitorId.IsNull() && !data.MonitorId.IsUnknown() {
        requestDataMap["monitorId"] = data.MonitorId.ValueString()
    }
    if !data.SnmpVersion.IsNull() && !data.SnmpVersion.IsUnknown() {
        requestDataMap["snmpVersion"] = data.SnmpVersion.ValueString()
    }
    if !data.SnmpCommunityString.IsNull() && !data.SnmpCommunityString.IsUnknown() {
        requestDataMap["snmpCommunityString"] = data.SnmpCommunityString.ValueString()
    }
    if !data.SnmpPort.IsNull() && !data.SnmpPort.IsUnknown() {
        requestDataMap["snmpPort"] = r.bigFloatToFloat64(data.SnmpPort.ValueBigFloat())
    }
    if parsedSnmpV3Auth := r.parseJSONField(data.SnmpV3Auth); parsedSnmpV3Auth != nil {
        requestDataMap["snmpV3Auth"] = parsedSnmpV3Auth
    }
    if !data.SnmpV3SecurityLevel.IsNull() && !data.SnmpV3SecurityLevel.IsUnknown() {
        requestDataMap["snmpV3SecurityLevel"] = data.SnmpV3SecurityLevel.ValueString()
    }
    if !data.SnmpV3Username.IsNull() && !data.SnmpV3Username.IsUnknown() {
        requestDataMap["snmpV3Username"] = data.SnmpV3Username.ValueString()
    }
    if !data.SnmpV3AuthProtocol.IsNull() && !data.SnmpV3AuthProtocol.IsUnknown() {
        requestDataMap["snmpV3AuthProtocol"] = data.SnmpV3AuthProtocol.ValueString()
    }
    if !data.SnmpV3AuthKey.IsNull() && !data.SnmpV3AuthKey.IsUnknown() {
        requestDataMap["snmpV3AuthKey"] = data.SnmpV3AuthKey.ValueString()
    }
    if !data.SnmpV3PrivProtocol.IsNull() && !data.SnmpV3PrivProtocol.IsUnknown() {
        requestDataMap["snmpV3PrivProtocol"] = data.SnmpV3PrivProtocol.ValueString()
    }
    if !data.SnmpV3PrivKey.IsNull() && !data.SnmpV3PrivKey.IsUnknown() {
        requestDataMap["snmpV3PrivKey"] = data.SnmpV3PrivKey.ValueString()
    }
    if !data.IsPollingEnabled.IsNull() && !data.IsPollingEnabled.IsUnknown() {
        requestDataMap["isPollingEnabled"] = data.IsPollingEnabled.ValueBool()
    }
    if !data.PollingIntervalInMinutes.IsNull() && !data.PollingIntervalInMinutes.IsUnknown() {
        requestDataMap["pollingIntervalInMinutes"] = r.bigFloatToFloat64(data.PollingIntervalInMinutes.ValueBigFloat())
    }
    if !data.WalkInterfaces.IsNull() && !data.WalkInterfaces.IsUnknown() {
        requestDataMap["walkInterfaces"] = data.WalkInterfaces.ValueBool()
    }
    if !data.CollectEndpoints.IsNull() && !data.CollectEndpoints.IsUnknown() {
        requestDataMap["collectEndpoints"] = data.CollectEndpoints.ValueBool()
    }
    if parsedSnmpOids := r.parseJSONField(data.SnmpOids); parsedSnmpOids != nil {
        requestDataMap["snmpOids"] = parsedSnmpOids
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        requestDataMap["createdByUserId"] = data.CreatedByUserId.ValueString()
    }
    if !data.IsArchived.IsNull() && !data.IsArchived.IsUnknown() {
        requestDataMap["isArchived"] = data.IsArchived.ValueBool()
    }
    if !data.Labels.IsNull() && !data.Labels.IsUnknown() {
        requestDataMap["labels"] = r.convertTerraformSetToInterface(data.Labels)
    }

    // Make API call
    httpResp, err := r.client.Post(ctx, "/network-device", networkDeviceRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create network_device, got error: %s", err))
        return
    }

    var networkDeviceResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &networkDeviceResponse)
    if err != nil {
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to create network_device: %s", err))
        return
    }

    // Extract the new resource id from the create response.
    createdId := ""
    if wrapper, ok := networkDeviceResponse["data"].(map[string]interface{}); ok {
        if val, ok := wrapper["_id"].(string); ok {
            createdId = val
        }
    } else if val, ok := networkDeviceResponse["_id"].(string); ok {
        createdId = val
    }
    if createdId == "" {
        resp.Diagnostics.AddError("OneUptime API Error", "Create response for network_device did not contain an id. This is a bug in the provider or the API; please report it.")
        return
    }
    data.Id = types.StringValue(createdId)

    /*
     * The server has committed the row. Persist what we know to state BEFORE
     * the read-back: if the read-back fails and we return without setting
     * state, Terraform never learns the resource exists and the created
     * network_device is orphaned server-side — never refreshed, never
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
        "hostname": true,
        "probeId": true,
        "siteId": true,
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
        "createdByUserId": true,
        "isArchived": true,
        "labels": true,
        "currentMonitorStatusId": true,
        "nextPollAt": true,
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
        "interfacesTotal": true,
        "interfacesUp": true,
        "interfacesDown": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "slug": true,
        "lastWalkLog": true,
        "archivedAt": true,
        "archivedByUserId": true,
        "deletedByUserId": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect(ctx, "/network-device/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        /*
         * State already owns the id, so the resource is tracked and the next
         * refresh reconciles the remaining attributes. Warn rather than
         * error: erroring here would strand a real resource.
         */
        resp.Diagnostics.AddWarning("Read After Create Failed", fmt.Sprintf("Created network_device but could not read it back; state is incomplete until the next refresh: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddWarning("Read After Create Failed", fmt.Sprintf("Created network_device but could not parse the read-back response; state is incomplete until the next refresh: %s", err))
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
    } else if val, ok := dataMap["description"].(string); ok {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if obj, ok := dataMap["hostname"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Hostname = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Hostname = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Hostname = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Hostname = types.StringValue(string(jsonBytes))
            } else {
                data.Hostname = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Hostname = types.StringValue(string(jsonBytes))
            } else {
                data.Hostname = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Hostname = types.StringValue(string(jsonBytes))
        } else {
            data.Hostname = types.StringNull()
        }
    } else if val, ok := dataMap["hostname"].(string); ok {
        data.Hostname = types.StringValue(val)
    } else {
        data.Hostname = types.StringNull()
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
    if obj, ok := dataMap["siteId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SiteId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SiteId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SiteId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SiteId = types.StringValue(string(jsonBytes))
            } else {
                data.SiteId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SiteId = types.StringValue(string(jsonBytes))
            } else {
                data.SiteId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SiteId = types.StringValue(string(jsonBytes))
        } else {
            data.SiteId = types.StringNull()
        }
    } else if val, ok := dataMap["siteId"].(string); ok {
        data.SiteId = types.StringValue(val)
    } else {
        data.SiteId = types.StringNull()
    }
    if obj, ok := dataMap["monitoringMethod"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MonitoringMethod = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.MonitoringMethod = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.MonitoringMethod = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.MonitoringMethod = types.StringValue(string(jsonBytes))
            } else {
                data.MonitoringMethod = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.MonitoringMethod = types.StringValue(string(jsonBytes))
            } else {
                data.MonitoringMethod = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.MonitoringMethod = types.StringValue(string(jsonBytes))
        } else {
            data.MonitoringMethod = types.StringNull()
        }
    } else if val, ok := dataMap["monitoringMethod"].(string); ok {
        data.MonitoringMethod = types.StringValue(val)
    } else {
        data.MonitoringMethod = types.StringNull()
    }
    if obj, ok := dataMap["deviceRole"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeviceRole = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DeviceRole = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DeviceRole = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DeviceRole = types.StringValue(string(jsonBytes))
            } else {
                data.DeviceRole = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DeviceRole = types.StringValue(string(jsonBytes))
            } else {
                data.DeviceRole = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DeviceRole = types.StringValue(string(jsonBytes))
        } else {
            data.DeviceRole = types.StringNull()
        }
    } else if val, ok := dataMap["deviceRole"].(string); ok {
        data.DeviceRole = types.StringValue(val)
    } else {
        data.DeviceRole = types.StringNull()
    }
    if obj, ok := dataMap["monitorId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MonitorId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.MonitorId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.MonitorId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.MonitorId = types.StringValue(string(jsonBytes))
            } else {
                data.MonitorId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.MonitorId = types.StringValue(string(jsonBytes))
            } else {
                data.MonitorId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.MonitorId = types.StringValue(string(jsonBytes))
        } else {
            data.MonitorId = types.StringNull()
        }
    } else if val, ok := dataMap["monitorId"].(string); ok {
        data.MonitorId = types.StringValue(val)
    } else {
        data.MonitorId = types.StringNull()
    }
    if obj, ok := dataMap["snmpVersion"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SnmpVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SnmpVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SnmpVersion = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SnmpVersion = types.StringValue(string(jsonBytes))
            } else {
                data.SnmpVersion = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SnmpVersion = types.StringValue(string(jsonBytes))
            } else {
                data.SnmpVersion = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SnmpVersion = types.StringValue(string(jsonBytes))
        } else {
            data.SnmpVersion = types.StringNull()
        }
    } else if val, ok := dataMap["snmpVersion"].(string); ok {
        data.SnmpVersion = types.StringValue(val)
    } else {
        data.SnmpVersion = types.StringNull()
    }
    if obj, ok := dataMap["snmpCommunityString"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SnmpCommunityString = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SnmpCommunityString = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SnmpCommunityString = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SnmpCommunityString = types.StringValue(string(jsonBytes))
            } else {
                data.SnmpCommunityString = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SnmpCommunityString = types.StringValue(string(jsonBytes))
            } else {
                data.SnmpCommunityString = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SnmpCommunityString = types.StringValue(string(jsonBytes))
        } else {
            data.SnmpCommunityString = types.StringNull()
        }
    } else if val, ok := dataMap["snmpCommunityString"].(string); ok {
        data.SnmpCommunityString = types.StringValue(val)
    } else {
        data.SnmpCommunityString = types.StringNull()
    }
    if val, ok := dataMap["snmpPort"].(float64); ok {
        data.SnmpPort = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["snmpPort"].(int); ok {
        data.SnmpPort = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["snmpPort"].(int64); ok {
        data.SnmpPort = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["snmpPort"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.SnmpPort = types.NumberValue(big.NewFloat(val))
        } else {
            data.SnmpPort = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.SnmpPort = types.NumberNull()
    }
    if obj, ok := dataMap["snmpV3Auth"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SnmpV3Auth = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SnmpV3Auth = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SnmpV3Auth = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SnmpV3Auth = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.SnmpV3Auth = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SnmpV3Auth = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.SnmpV3Auth = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SnmpV3Auth = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.SnmpV3Auth = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["snmpV3Auth"].(string); ok {
        data.SnmpV3Auth = NewJSONSubsetValue(val)
    } else {
        data.SnmpV3Auth = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["snmpV3SecurityLevel"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SnmpV3SecurityLevel = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SnmpV3SecurityLevel = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SnmpV3SecurityLevel = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SnmpV3SecurityLevel = types.StringValue(string(jsonBytes))
            } else {
                data.SnmpV3SecurityLevel = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SnmpV3SecurityLevel = types.StringValue(string(jsonBytes))
            } else {
                data.SnmpV3SecurityLevel = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SnmpV3SecurityLevel = types.StringValue(string(jsonBytes))
        } else {
            data.SnmpV3SecurityLevel = types.StringNull()
        }
    } else if val, ok := dataMap["snmpV3SecurityLevel"].(string); ok {
        data.SnmpV3SecurityLevel = types.StringValue(val)
    } else {
        data.SnmpV3SecurityLevel = types.StringNull()
    }
    if obj, ok := dataMap["snmpV3Username"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SnmpV3Username = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SnmpV3Username = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SnmpV3Username = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SnmpV3Username = types.StringValue(string(jsonBytes))
            } else {
                data.SnmpV3Username = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SnmpV3Username = types.StringValue(string(jsonBytes))
            } else {
                data.SnmpV3Username = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SnmpV3Username = types.StringValue(string(jsonBytes))
        } else {
            data.SnmpV3Username = types.StringNull()
        }
    } else if val, ok := dataMap["snmpV3Username"].(string); ok {
        data.SnmpV3Username = types.StringValue(val)
    } else {
        data.SnmpV3Username = types.StringNull()
    }
    if obj, ok := dataMap["snmpV3AuthProtocol"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SnmpV3AuthProtocol = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SnmpV3AuthProtocol = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SnmpV3AuthProtocol = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SnmpV3AuthProtocol = types.StringValue(string(jsonBytes))
            } else {
                data.SnmpV3AuthProtocol = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SnmpV3AuthProtocol = types.StringValue(string(jsonBytes))
            } else {
                data.SnmpV3AuthProtocol = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SnmpV3AuthProtocol = types.StringValue(string(jsonBytes))
        } else {
            data.SnmpV3AuthProtocol = types.StringNull()
        }
    } else if val, ok := dataMap["snmpV3AuthProtocol"].(string); ok {
        data.SnmpV3AuthProtocol = types.StringValue(val)
    } else {
        data.SnmpV3AuthProtocol = types.StringNull()
    }
    if obj, ok := dataMap["snmpV3AuthKey"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SnmpV3AuthKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SnmpV3AuthKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SnmpV3AuthKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SnmpV3AuthKey = types.StringValue(string(jsonBytes))
            } else {
                data.SnmpV3AuthKey = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SnmpV3AuthKey = types.StringValue(string(jsonBytes))
            } else {
                data.SnmpV3AuthKey = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SnmpV3AuthKey = types.StringValue(string(jsonBytes))
        } else {
            data.SnmpV3AuthKey = types.StringNull()
        }
    } else if val, ok := dataMap["snmpV3AuthKey"].(string); ok {
        data.SnmpV3AuthKey = types.StringValue(val)
    } else {
        data.SnmpV3AuthKey = types.StringNull()
    }
    if obj, ok := dataMap["snmpV3PrivProtocol"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SnmpV3PrivProtocol = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SnmpV3PrivProtocol = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SnmpV3PrivProtocol = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SnmpV3PrivProtocol = types.StringValue(string(jsonBytes))
            } else {
                data.SnmpV3PrivProtocol = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SnmpV3PrivProtocol = types.StringValue(string(jsonBytes))
            } else {
                data.SnmpV3PrivProtocol = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SnmpV3PrivProtocol = types.StringValue(string(jsonBytes))
        } else {
            data.SnmpV3PrivProtocol = types.StringNull()
        }
    } else if val, ok := dataMap["snmpV3PrivProtocol"].(string); ok {
        data.SnmpV3PrivProtocol = types.StringValue(val)
    } else {
        data.SnmpV3PrivProtocol = types.StringNull()
    }
    if obj, ok := dataMap["snmpV3PrivKey"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SnmpV3PrivKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SnmpV3PrivKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SnmpV3PrivKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SnmpV3PrivKey = types.StringValue(string(jsonBytes))
            } else {
                data.SnmpV3PrivKey = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SnmpV3PrivKey = types.StringValue(string(jsonBytes))
            } else {
                data.SnmpV3PrivKey = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SnmpV3PrivKey = types.StringValue(string(jsonBytes))
        } else {
            data.SnmpV3PrivKey = types.StringNull()
        }
    } else if val, ok := dataMap["snmpV3PrivKey"].(string); ok {
        data.SnmpV3PrivKey = types.StringValue(val)
    } else {
        data.SnmpV3PrivKey = types.StringNull()
    }
    if val, ok := dataMap["isPollingEnabled"].(bool); ok {
        data.IsPollingEnabled = types.BoolValue(val)
    }
    if val, ok := dataMap["pollingIntervalInMinutes"].(float64); ok {
        data.PollingIntervalInMinutes = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["pollingIntervalInMinutes"].(int); ok {
        data.PollingIntervalInMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["pollingIntervalInMinutes"].(int64); ok {
        data.PollingIntervalInMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["pollingIntervalInMinutes"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.PollingIntervalInMinutes = types.NumberValue(big.NewFloat(val))
        } else {
            data.PollingIntervalInMinutes = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.PollingIntervalInMinutes = types.NumberNull()
    }
    if val, ok := dataMap["walkInterfaces"].(bool); ok {
        data.WalkInterfaces = types.BoolValue(val)
    }
    if val, ok := dataMap["collectEndpoints"].(bool); ok {
        data.CollectEndpoints = types.BoolValue(val)
    }
    if obj, ok := dataMap["snmpOids"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SnmpOids = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SnmpOids = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SnmpOids = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SnmpOids = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.SnmpOids = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SnmpOids = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.SnmpOids = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SnmpOids = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.SnmpOids = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["snmpOids"].(string); ok {
        data.SnmpOids = NewJSONSubsetValue(val)
    } else {
        data.SnmpOids = NewJSONSubsetNull()
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
    if obj, ok := dataMap["currentMonitorStatusId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CurrentMonitorStatusId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CurrentMonitorStatusId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CurrentMonitorStatusId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CurrentMonitorStatusId = types.StringValue(string(jsonBytes))
            } else {
                data.CurrentMonitorStatusId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CurrentMonitorStatusId = types.StringValue(string(jsonBytes))
            } else {
                data.CurrentMonitorStatusId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CurrentMonitorStatusId = types.StringValue(string(jsonBytes))
        } else {
            data.CurrentMonitorStatusId = types.StringNull()
        }
    } else if val, ok := dataMap["currentMonitorStatusId"].(string); ok {
        data.CurrentMonitorStatusId = types.StringValue(val)
    } else {
        data.CurrentMonitorStatusId = types.StringNull()
    }
    if obj, ok := dataMap["nextPollAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.NextPollAt = NewRFC3339Value(val)
        } else {
            data.NextPollAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["nextPollAt"].(string); ok && val != "" {
        data.NextPollAt = NewRFC3339Value(val)
    } else {
        data.NextPollAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["sysDescr"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SysDescr = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SysDescr = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SysDescr = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SysDescr = types.StringValue(string(jsonBytes))
            } else {
                data.SysDescr = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SysDescr = types.StringValue(string(jsonBytes))
            } else {
                data.SysDescr = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SysDescr = types.StringValue(string(jsonBytes))
        } else {
            data.SysDescr = types.StringNull()
        }
    } else if val, ok := dataMap["sysDescr"].(string); ok {
        data.SysDescr = types.StringValue(val)
    } else {
        data.SysDescr = types.StringNull()
    }
    if obj, ok := dataMap["sysName"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SysName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SysName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SysName = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SysName = types.StringValue(string(jsonBytes))
            } else {
                data.SysName = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SysName = types.StringValue(string(jsonBytes))
            } else {
                data.SysName = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SysName = types.StringValue(string(jsonBytes))
        } else {
            data.SysName = types.StringNull()
        }
    } else if val, ok := dataMap["sysName"].(string); ok {
        data.SysName = types.StringValue(val)
    } else {
        data.SysName = types.StringNull()
    }
    if obj, ok := dataMap["sysObjectId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SysObjectId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SysObjectId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SysObjectId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SysObjectId = types.StringValue(string(jsonBytes))
            } else {
                data.SysObjectId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SysObjectId = types.StringValue(string(jsonBytes))
            } else {
                data.SysObjectId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SysObjectId = types.StringValue(string(jsonBytes))
        } else {
            data.SysObjectId = types.StringNull()
        }
    } else if val, ok := dataMap["sysObjectId"].(string); ok {
        data.SysObjectId = types.StringValue(val)
    } else {
        data.SysObjectId = types.StringNull()
    }
    if obj, ok := dataMap["sysLocation"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SysLocation = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SysLocation = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SysLocation = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SysLocation = types.StringValue(string(jsonBytes))
            } else {
                data.SysLocation = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SysLocation = types.StringValue(string(jsonBytes))
            } else {
                data.SysLocation = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SysLocation = types.StringValue(string(jsonBytes))
        } else {
            data.SysLocation = types.StringNull()
        }
    } else if val, ok := dataMap["sysLocation"].(string); ok {
        data.SysLocation = types.StringValue(val)
    } else {
        data.SysLocation = types.StringNull()
    }
    if obj, ok := dataMap["sysContact"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SysContact = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SysContact = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SysContact = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SysContact = types.StringValue(string(jsonBytes))
            } else {
                data.SysContact = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SysContact = types.StringValue(string(jsonBytes))
            } else {
                data.SysContact = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SysContact = types.StringValue(string(jsonBytes))
        } else {
            data.SysContact = types.StringNull()
        }
    } else if val, ok := dataMap["sysContact"].(string); ok {
        data.SysContact = types.StringValue(val)
    } else {
        data.SysContact = types.StringNull()
    }
    if obj, ok := dataMap["vendor"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Vendor = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Vendor = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Vendor = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Vendor = types.StringValue(string(jsonBytes))
            } else {
                data.Vendor = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Vendor = types.StringValue(string(jsonBytes))
            } else {
                data.Vendor = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Vendor = types.StringValue(string(jsonBytes))
        } else {
            data.Vendor = types.StringNull()
        }
    } else if val, ok := dataMap["vendor"].(string); ok {
        data.Vendor = types.StringValue(val)
    } else {
        data.Vendor = types.StringNull()
    }
    if obj, ok := dataMap["deviceModel"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeviceModel = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DeviceModel = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DeviceModel = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DeviceModel = types.StringValue(string(jsonBytes))
            } else {
                data.DeviceModel = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DeviceModel = types.StringValue(string(jsonBytes))
            } else {
                data.DeviceModel = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DeviceModel = types.StringValue(string(jsonBytes))
        } else {
            data.DeviceModel = types.StringNull()
        }
    } else if val, ok := dataMap["deviceModel"].(string); ok {
        data.DeviceModel = types.StringValue(val)
    } else {
        data.DeviceModel = types.StringNull()
    }
    if obj, ok := dataMap["serialNumber"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SerialNumber = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SerialNumber = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SerialNumber = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SerialNumber = types.StringValue(string(jsonBytes))
            } else {
                data.SerialNumber = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SerialNumber = types.StringValue(string(jsonBytes))
            } else {
                data.SerialNumber = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SerialNumber = types.StringValue(string(jsonBytes))
        } else {
            data.SerialNumber = types.StringNull()
        }
    } else if val, ok := dataMap["serialNumber"].(string); ok {
        data.SerialNumber = types.StringValue(val)
    } else {
        data.SerialNumber = types.StringNull()
    }
    if obj, ok := dataMap["firmwareVersion"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FirmwareVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.FirmwareVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.FirmwareVersion = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.FirmwareVersion = types.StringValue(string(jsonBytes))
            } else {
                data.FirmwareVersion = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.FirmwareVersion = types.StringValue(string(jsonBytes))
            } else {
                data.FirmwareVersion = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.FirmwareVersion = types.StringValue(string(jsonBytes))
        } else {
            data.FirmwareVersion = types.StringNull()
        }
    } else if val, ok := dataMap["firmwareVersion"].(string); ok {
        data.FirmwareVersion = types.StringValue(val)
    } else {
        data.FirmwareVersion = types.StringNull()
    }
    if obj, ok := dataMap["softwareVersion"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SoftwareVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SoftwareVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SoftwareVersion = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SoftwareVersion = types.StringValue(string(jsonBytes))
            } else {
                data.SoftwareVersion = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SoftwareVersion = types.StringValue(string(jsonBytes))
            } else {
                data.SoftwareVersion = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SoftwareVersion = types.StringValue(string(jsonBytes))
        } else {
            data.SoftwareVersion = types.StringNull()
        }
    } else if val, ok := dataMap["softwareVersion"].(string); ok {
        data.SoftwareVersion = types.StringValue(val)
    } else {
        data.SoftwareVersion = types.StringNull()
    }
    if obj, ok := dataMap["lastRebootedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.LastRebootedAt = NewRFC3339Value(val)
        } else {
            data.LastRebootedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["lastRebootedAt"].(string); ok && val != "" {
        data.LastRebootedAt = NewRFC3339Value(val)
    } else {
        data.LastRebootedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["cdpNeighbors"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CdpNeighbors = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CdpNeighbors = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CdpNeighbors = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CdpNeighbors = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.CdpNeighbors = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CdpNeighbors = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.CdpNeighbors = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CdpNeighbors = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.CdpNeighbors = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["cdpNeighbors"].(string); ok {
        data.CdpNeighbors = NewJSONSubsetValue(val)
    } else {
        data.CdpNeighbors = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["lldpNeighbors"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LldpNeighbors = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.LldpNeighbors = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.LldpNeighbors = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.LldpNeighbors = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.LldpNeighbors = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.LldpNeighbors = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.LldpNeighbors = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.LldpNeighbors = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.LldpNeighbors = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["lldpNeighbors"].(string); ok {
        data.LldpNeighbors = NewJSONSubsetValue(val)
    } else {
        data.LldpNeighbors = NewJSONSubsetNull()
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
    if val, ok := dataMap["interfacesTotal"].(float64); ok {
        data.InterfacesTotal = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["interfacesTotal"].(int); ok {
        data.InterfacesTotal = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["interfacesTotal"].(int64); ok {
        data.InterfacesTotal = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["interfacesTotal"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.InterfacesTotal = types.NumberValue(big.NewFloat(val))
        } else {
            data.InterfacesTotal = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.InterfacesTotal = types.NumberNull()
    }
    if val, ok := dataMap["interfacesUp"].(float64); ok {
        data.InterfacesUp = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["interfacesUp"].(int); ok {
        data.InterfacesUp = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["interfacesUp"].(int64); ok {
        data.InterfacesUp = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["interfacesUp"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.InterfacesUp = types.NumberValue(big.NewFloat(val))
        } else {
            data.InterfacesUp = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.InterfacesUp = types.NumberNull()
    }
    if val, ok := dataMap["interfacesDown"].(float64); ok {
        data.InterfacesDown = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["interfacesDown"].(int); ok {
        data.InterfacesDown = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["interfacesDown"].(int64); ok {
        data.InterfacesDown = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["interfacesDown"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.InterfacesDown = types.NumberValue(big.NewFloat(val))
        } else {
            data.InterfacesDown = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.InterfacesDown = types.NumberNull()
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
    if obj, ok := dataMap["deletedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.DeletedAt = NewRFC3339Value(val)
        } else {
            data.DeletedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["deletedAt"].(string); ok && val != "" {
        data.DeletedAt = NewRFC3339Value(val)
    } else {
        data.DeletedAt = NewRFC3339Null()
    }
    if val, ok := dataMap["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["version"].(int); ok {
        data.Version = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["version"].(int64); ok {
        data.Version = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["version"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.Version = types.NumberValue(big.NewFloat(val))
        } else {
            data.Version = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.Version = types.NumberNull()
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
    } else if val, ok := dataMap["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
    }
    if obj, ok := dataMap["lastWalkLog"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastWalkLog = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.LastWalkLog = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.LastWalkLog = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.LastWalkLog = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.LastWalkLog = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.LastWalkLog = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.LastWalkLog = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.LastWalkLog = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.LastWalkLog = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["lastWalkLog"].(string); ok {
        data.LastWalkLog = NewJSONSubsetValue(val)
    } else {
        data.LastWalkLog = NewJSONSubsetNull()
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
    } else if val, ok := dataMap["archivedByUserId"].(string); ok {
        data.ArchivedByUserId = types.StringValue(val)
    } else {
        data.ArchivedByUserId = types.StringNull()
    }
    if obj, ok := dataMap["deletedByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeletedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DeletedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DeletedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DeletedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.DeletedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DeletedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.DeletedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DeletedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.DeletedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["deletedByUserId"].(string); ok {
        data.DeletedByUserId = types.StringValue(val)
    } else {
        data.DeletedByUserId = types.StringNull()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    // The read response is authoritative, but never let it clobber the id we just received.
    data.Id = types.StringValue(createdId)

    // Write logs using the tflog package
    tflog.Trace(ctx, "created a resource")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NetworkDeviceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var data NetworkDeviceResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create select parameter to get full object
    selectParam := map[string]interface{}{
        "projectId": true,
        "name": true,
        "description": true,
        "hostname": true,
        "probeId": true,
        "siteId": true,
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
        "createdByUserId": true,
        "isArchived": true,
        "labels": true,
        "currentMonitorStatusId": true,
        "nextPollAt": true,
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
        "interfacesTotal": true,
        "interfacesUp": true,
        "interfacesDown": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "slug": true,
        "lastWalkLog": true,
        "archivedAt": true,
        "archivedByUserId": true,
        "deletedByUserId": true,
        "_id": true,
    }

    // Make API call with select parameter
    httpResp, err := r.client.PostWithSelect(ctx, "/network-device/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read network_device, got error: %s", err))
        return
    }

    if httpResp.StatusCode == http.StatusNotFound {
        resp.State.RemoveResource(ctx)
        return
    }

    var networkDeviceResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &networkDeviceResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse network_device response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := networkDeviceResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = networkDeviceResponse
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
    } else if val, ok := dataMap["description"].(string); ok {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if obj, ok := dataMap["hostname"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Hostname = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Hostname = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Hostname = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Hostname = types.StringValue(string(jsonBytes))
            } else {
                data.Hostname = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Hostname = types.StringValue(string(jsonBytes))
            } else {
                data.Hostname = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Hostname = types.StringValue(string(jsonBytes))
        } else {
            data.Hostname = types.StringNull()
        }
    } else if val, ok := dataMap["hostname"].(string); ok {
        data.Hostname = types.StringValue(val)
    } else {
        data.Hostname = types.StringNull()
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
    if obj, ok := dataMap["siteId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SiteId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SiteId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SiteId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SiteId = types.StringValue(string(jsonBytes))
            } else {
                data.SiteId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SiteId = types.StringValue(string(jsonBytes))
            } else {
                data.SiteId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SiteId = types.StringValue(string(jsonBytes))
        } else {
            data.SiteId = types.StringNull()
        }
    } else if val, ok := dataMap["siteId"].(string); ok {
        data.SiteId = types.StringValue(val)
    } else {
        data.SiteId = types.StringNull()
    }
    if obj, ok := dataMap["monitoringMethod"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MonitoringMethod = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.MonitoringMethod = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.MonitoringMethod = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.MonitoringMethod = types.StringValue(string(jsonBytes))
            } else {
                data.MonitoringMethod = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.MonitoringMethod = types.StringValue(string(jsonBytes))
            } else {
                data.MonitoringMethod = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.MonitoringMethod = types.StringValue(string(jsonBytes))
        } else {
            data.MonitoringMethod = types.StringNull()
        }
    } else if val, ok := dataMap["monitoringMethod"].(string); ok {
        data.MonitoringMethod = types.StringValue(val)
    } else {
        data.MonitoringMethod = types.StringNull()
    }
    if obj, ok := dataMap["deviceRole"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeviceRole = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DeviceRole = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DeviceRole = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DeviceRole = types.StringValue(string(jsonBytes))
            } else {
                data.DeviceRole = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DeviceRole = types.StringValue(string(jsonBytes))
            } else {
                data.DeviceRole = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DeviceRole = types.StringValue(string(jsonBytes))
        } else {
            data.DeviceRole = types.StringNull()
        }
    } else if val, ok := dataMap["deviceRole"].(string); ok {
        data.DeviceRole = types.StringValue(val)
    } else {
        data.DeviceRole = types.StringNull()
    }
    if obj, ok := dataMap["monitorId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MonitorId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.MonitorId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.MonitorId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.MonitorId = types.StringValue(string(jsonBytes))
            } else {
                data.MonitorId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.MonitorId = types.StringValue(string(jsonBytes))
            } else {
                data.MonitorId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.MonitorId = types.StringValue(string(jsonBytes))
        } else {
            data.MonitorId = types.StringNull()
        }
    } else if val, ok := dataMap["monitorId"].(string); ok {
        data.MonitorId = types.StringValue(val)
    } else {
        data.MonitorId = types.StringNull()
    }
    if obj, ok := dataMap["snmpVersion"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SnmpVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SnmpVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SnmpVersion = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SnmpVersion = types.StringValue(string(jsonBytes))
            } else {
                data.SnmpVersion = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SnmpVersion = types.StringValue(string(jsonBytes))
            } else {
                data.SnmpVersion = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SnmpVersion = types.StringValue(string(jsonBytes))
        } else {
            data.SnmpVersion = types.StringNull()
        }
    } else if val, ok := dataMap["snmpVersion"].(string); ok {
        data.SnmpVersion = types.StringValue(val)
    } else {
        data.SnmpVersion = types.StringNull()
    }
    if obj, ok := dataMap["snmpCommunityString"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SnmpCommunityString = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SnmpCommunityString = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SnmpCommunityString = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SnmpCommunityString = types.StringValue(string(jsonBytes))
            } else {
                data.SnmpCommunityString = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SnmpCommunityString = types.StringValue(string(jsonBytes))
            } else {
                data.SnmpCommunityString = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SnmpCommunityString = types.StringValue(string(jsonBytes))
        } else {
            data.SnmpCommunityString = types.StringNull()
        }
    } else if val, ok := dataMap["snmpCommunityString"].(string); ok {
        data.SnmpCommunityString = types.StringValue(val)
    } else {
        data.SnmpCommunityString = types.StringNull()
    }
    if val, ok := dataMap["snmpPort"].(float64); ok {
        data.SnmpPort = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["snmpPort"].(int); ok {
        data.SnmpPort = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["snmpPort"].(int64); ok {
        data.SnmpPort = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["snmpPort"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.SnmpPort = types.NumberValue(big.NewFloat(val))
        } else {
            data.SnmpPort = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.SnmpPort = types.NumberNull()
    }
    if obj, ok := dataMap["snmpV3Auth"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SnmpV3Auth = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SnmpV3Auth = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SnmpV3Auth = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SnmpV3Auth = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.SnmpV3Auth = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SnmpV3Auth = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.SnmpV3Auth = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SnmpV3Auth = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.SnmpV3Auth = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["snmpV3Auth"].(string); ok {
        data.SnmpV3Auth = NewJSONSubsetValue(val)
    } else {
        data.SnmpV3Auth = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["snmpV3SecurityLevel"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SnmpV3SecurityLevel = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SnmpV3SecurityLevel = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SnmpV3SecurityLevel = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SnmpV3SecurityLevel = types.StringValue(string(jsonBytes))
            } else {
                data.SnmpV3SecurityLevel = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SnmpV3SecurityLevel = types.StringValue(string(jsonBytes))
            } else {
                data.SnmpV3SecurityLevel = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SnmpV3SecurityLevel = types.StringValue(string(jsonBytes))
        } else {
            data.SnmpV3SecurityLevel = types.StringNull()
        }
    } else if val, ok := dataMap["snmpV3SecurityLevel"].(string); ok {
        data.SnmpV3SecurityLevel = types.StringValue(val)
    } else {
        data.SnmpV3SecurityLevel = types.StringNull()
    }
    if obj, ok := dataMap["snmpV3Username"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SnmpV3Username = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SnmpV3Username = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SnmpV3Username = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SnmpV3Username = types.StringValue(string(jsonBytes))
            } else {
                data.SnmpV3Username = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SnmpV3Username = types.StringValue(string(jsonBytes))
            } else {
                data.SnmpV3Username = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SnmpV3Username = types.StringValue(string(jsonBytes))
        } else {
            data.SnmpV3Username = types.StringNull()
        }
    } else if val, ok := dataMap["snmpV3Username"].(string); ok {
        data.SnmpV3Username = types.StringValue(val)
    } else {
        data.SnmpV3Username = types.StringNull()
    }
    if obj, ok := dataMap["snmpV3AuthProtocol"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SnmpV3AuthProtocol = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SnmpV3AuthProtocol = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SnmpV3AuthProtocol = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SnmpV3AuthProtocol = types.StringValue(string(jsonBytes))
            } else {
                data.SnmpV3AuthProtocol = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SnmpV3AuthProtocol = types.StringValue(string(jsonBytes))
            } else {
                data.SnmpV3AuthProtocol = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SnmpV3AuthProtocol = types.StringValue(string(jsonBytes))
        } else {
            data.SnmpV3AuthProtocol = types.StringNull()
        }
    } else if val, ok := dataMap["snmpV3AuthProtocol"].(string); ok {
        data.SnmpV3AuthProtocol = types.StringValue(val)
    } else {
        data.SnmpV3AuthProtocol = types.StringNull()
    }
    if obj, ok := dataMap["snmpV3AuthKey"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SnmpV3AuthKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SnmpV3AuthKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SnmpV3AuthKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SnmpV3AuthKey = types.StringValue(string(jsonBytes))
            } else {
                data.SnmpV3AuthKey = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SnmpV3AuthKey = types.StringValue(string(jsonBytes))
            } else {
                data.SnmpV3AuthKey = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SnmpV3AuthKey = types.StringValue(string(jsonBytes))
        } else {
            data.SnmpV3AuthKey = types.StringNull()
        }
    } else if val, ok := dataMap["snmpV3AuthKey"].(string); ok {
        data.SnmpV3AuthKey = types.StringValue(val)
    } else {
        data.SnmpV3AuthKey = types.StringNull()
    }
    if obj, ok := dataMap["snmpV3PrivProtocol"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SnmpV3PrivProtocol = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SnmpV3PrivProtocol = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SnmpV3PrivProtocol = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SnmpV3PrivProtocol = types.StringValue(string(jsonBytes))
            } else {
                data.SnmpV3PrivProtocol = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SnmpV3PrivProtocol = types.StringValue(string(jsonBytes))
            } else {
                data.SnmpV3PrivProtocol = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SnmpV3PrivProtocol = types.StringValue(string(jsonBytes))
        } else {
            data.SnmpV3PrivProtocol = types.StringNull()
        }
    } else if val, ok := dataMap["snmpV3PrivProtocol"].(string); ok {
        data.SnmpV3PrivProtocol = types.StringValue(val)
    } else {
        data.SnmpV3PrivProtocol = types.StringNull()
    }
    if obj, ok := dataMap["snmpV3PrivKey"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SnmpV3PrivKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SnmpV3PrivKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SnmpV3PrivKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SnmpV3PrivKey = types.StringValue(string(jsonBytes))
            } else {
                data.SnmpV3PrivKey = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SnmpV3PrivKey = types.StringValue(string(jsonBytes))
            } else {
                data.SnmpV3PrivKey = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SnmpV3PrivKey = types.StringValue(string(jsonBytes))
        } else {
            data.SnmpV3PrivKey = types.StringNull()
        }
    } else if val, ok := dataMap["snmpV3PrivKey"].(string); ok {
        data.SnmpV3PrivKey = types.StringValue(val)
    } else {
        data.SnmpV3PrivKey = types.StringNull()
    }
    if val, ok := dataMap["isPollingEnabled"].(bool); ok {
        data.IsPollingEnabled = types.BoolValue(val)
    }
    if val, ok := dataMap["pollingIntervalInMinutes"].(float64); ok {
        data.PollingIntervalInMinutes = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["pollingIntervalInMinutes"].(int); ok {
        data.PollingIntervalInMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["pollingIntervalInMinutes"].(int64); ok {
        data.PollingIntervalInMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["pollingIntervalInMinutes"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.PollingIntervalInMinutes = types.NumberValue(big.NewFloat(val))
        } else {
            data.PollingIntervalInMinutes = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.PollingIntervalInMinutes = types.NumberNull()
    }
    if val, ok := dataMap["walkInterfaces"].(bool); ok {
        data.WalkInterfaces = types.BoolValue(val)
    }
    if val, ok := dataMap["collectEndpoints"].(bool); ok {
        data.CollectEndpoints = types.BoolValue(val)
    }
    if obj, ok := dataMap["snmpOids"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SnmpOids = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SnmpOids = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SnmpOids = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SnmpOids = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.SnmpOids = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SnmpOids = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.SnmpOids = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SnmpOids = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.SnmpOids = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["snmpOids"].(string); ok {
        data.SnmpOids = NewJSONSubsetValue(val)
    } else {
        data.SnmpOids = NewJSONSubsetNull()
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
    if obj, ok := dataMap["currentMonitorStatusId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CurrentMonitorStatusId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CurrentMonitorStatusId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CurrentMonitorStatusId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CurrentMonitorStatusId = types.StringValue(string(jsonBytes))
            } else {
                data.CurrentMonitorStatusId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CurrentMonitorStatusId = types.StringValue(string(jsonBytes))
            } else {
                data.CurrentMonitorStatusId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CurrentMonitorStatusId = types.StringValue(string(jsonBytes))
        } else {
            data.CurrentMonitorStatusId = types.StringNull()
        }
    } else if val, ok := dataMap["currentMonitorStatusId"].(string); ok {
        data.CurrentMonitorStatusId = types.StringValue(val)
    } else {
        data.CurrentMonitorStatusId = types.StringNull()
    }
    if obj, ok := dataMap["nextPollAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.NextPollAt = NewRFC3339Value(val)
        } else {
            data.NextPollAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["nextPollAt"].(string); ok && val != "" {
        data.NextPollAt = NewRFC3339Value(val)
    } else {
        data.NextPollAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["sysDescr"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SysDescr = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SysDescr = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SysDescr = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SysDescr = types.StringValue(string(jsonBytes))
            } else {
                data.SysDescr = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SysDescr = types.StringValue(string(jsonBytes))
            } else {
                data.SysDescr = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SysDescr = types.StringValue(string(jsonBytes))
        } else {
            data.SysDescr = types.StringNull()
        }
    } else if val, ok := dataMap["sysDescr"].(string); ok {
        data.SysDescr = types.StringValue(val)
    } else {
        data.SysDescr = types.StringNull()
    }
    if obj, ok := dataMap["sysName"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SysName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SysName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SysName = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SysName = types.StringValue(string(jsonBytes))
            } else {
                data.SysName = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SysName = types.StringValue(string(jsonBytes))
            } else {
                data.SysName = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SysName = types.StringValue(string(jsonBytes))
        } else {
            data.SysName = types.StringNull()
        }
    } else if val, ok := dataMap["sysName"].(string); ok {
        data.SysName = types.StringValue(val)
    } else {
        data.SysName = types.StringNull()
    }
    if obj, ok := dataMap["sysObjectId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SysObjectId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SysObjectId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SysObjectId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SysObjectId = types.StringValue(string(jsonBytes))
            } else {
                data.SysObjectId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SysObjectId = types.StringValue(string(jsonBytes))
            } else {
                data.SysObjectId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SysObjectId = types.StringValue(string(jsonBytes))
        } else {
            data.SysObjectId = types.StringNull()
        }
    } else if val, ok := dataMap["sysObjectId"].(string); ok {
        data.SysObjectId = types.StringValue(val)
    } else {
        data.SysObjectId = types.StringNull()
    }
    if obj, ok := dataMap["sysLocation"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SysLocation = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SysLocation = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SysLocation = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SysLocation = types.StringValue(string(jsonBytes))
            } else {
                data.SysLocation = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SysLocation = types.StringValue(string(jsonBytes))
            } else {
                data.SysLocation = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SysLocation = types.StringValue(string(jsonBytes))
        } else {
            data.SysLocation = types.StringNull()
        }
    } else if val, ok := dataMap["sysLocation"].(string); ok {
        data.SysLocation = types.StringValue(val)
    } else {
        data.SysLocation = types.StringNull()
    }
    if obj, ok := dataMap["sysContact"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SysContact = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SysContact = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SysContact = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SysContact = types.StringValue(string(jsonBytes))
            } else {
                data.SysContact = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SysContact = types.StringValue(string(jsonBytes))
            } else {
                data.SysContact = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SysContact = types.StringValue(string(jsonBytes))
        } else {
            data.SysContact = types.StringNull()
        }
    } else if val, ok := dataMap["sysContact"].(string); ok {
        data.SysContact = types.StringValue(val)
    } else {
        data.SysContact = types.StringNull()
    }
    if obj, ok := dataMap["vendor"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Vendor = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Vendor = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Vendor = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Vendor = types.StringValue(string(jsonBytes))
            } else {
                data.Vendor = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Vendor = types.StringValue(string(jsonBytes))
            } else {
                data.Vendor = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Vendor = types.StringValue(string(jsonBytes))
        } else {
            data.Vendor = types.StringNull()
        }
    } else if val, ok := dataMap["vendor"].(string); ok {
        data.Vendor = types.StringValue(val)
    } else {
        data.Vendor = types.StringNull()
    }
    if obj, ok := dataMap["deviceModel"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeviceModel = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DeviceModel = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DeviceModel = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DeviceModel = types.StringValue(string(jsonBytes))
            } else {
                data.DeviceModel = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DeviceModel = types.StringValue(string(jsonBytes))
            } else {
                data.DeviceModel = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DeviceModel = types.StringValue(string(jsonBytes))
        } else {
            data.DeviceModel = types.StringNull()
        }
    } else if val, ok := dataMap["deviceModel"].(string); ok {
        data.DeviceModel = types.StringValue(val)
    } else {
        data.DeviceModel = types.StringNull()
    }
    if obj, ok := dataMap["serialNumber"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SerialNumber = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SerialNumber = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SerialNumber = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SerialNumber = types.StringValue(string(jsonBytes))
            } else {
                data.SerialNumber = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SerialNumber = types.StringValue(string(jsonBytes))
            } else {
                data.SerialNumber = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SerialNumber = types.StringValue(string(jsonBytes))
        } else {
            data.SerialNumber = types.StringNull()
        }
    } else if val, ok := dataMap["serialNumber"].(string); ok {
        data.SerialNumber = types.StringValue(val)
    } else {
        data.SerialNumber = types.StringNull()
    }
    if obj, ok := dataMap["firmwareVersion"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FirmwareVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.FirmwareVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.FirmwareVersion = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.FirmwareVersion = types.StringValue(string(jsonBytes))
            } else {
                data.FirmwareVersion = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.FirmwareVersion = types.StringValue(string(jsonBytes))
            } else {
                data.FirmwareVersion = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.FirmwareVersion = types.StringValue(string(jsonBytes))
        } else {
            data.FirmwareVersion = types.StringNull()
        }
    } else if val, ok := dataMap["firmwareVersion"].(string); ok {
        data.FirmwareVersion = types.StringValue(val)
    } else {
        data.FirmwareVersion = types.StringNull()
    }
    if obj, ok := dataMap["softwareVersion"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SoftwareVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SoftwareVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SoftwareVersion = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SoftwareVersion = types.StringValue(string(jsonBytes))
            } else {
                data.SoftwareVersion = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SoftwareVersion = types.StringValue(string(jsonBytes))
            } else {
                data.SoftwareVersion = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SoftwareVersion = types.StringValue(string(jsonBytes))
        } else {
            data.SoftwareVersion = types.StringNull()
        }
    } else if val, ok := dataMap["softwareVersion"].(string); ok {
        data.SoftwareVersion = types.StringValue(val)
    } else {
        data.SoftwareVersion = types.StringNull()
    }
    if obj, ok := dataMap["lastRebootedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.LastRebootedAt = NewRFC3339Value(val)
        } else {
            data.LastRebootedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["lastRebootedAt"].(string); ok && val != "" {
        data.LastRebootedAt = NewRFC3339Value(val)
    } else {
        data.LastRebootedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["cdpNeighbors"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CdpNeighbors = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CdpNeighbors = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CdpNeighbors = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CdpNeighbors = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.CdpNeighbors = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CdpNeighbors = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.CdpNeighbors = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CdpNeighbors = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.CdpNeighbors = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["cdpNeighbors"].(string); ok {
        data.CdpNeighbors = NewJSONSubsetValue(val)
    } else {
        data.CdpNeighbors = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["lldpNeighbors"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LldpNeighbors = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.LldpNeighbors = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.LldpNeighbors = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.LldpNeighbors = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.LldpNeighbors = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.LldpNeighbors = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.LldpNeighbors = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.LldpNeighbors = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.LldpNeighbors = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["lldpNeighbors"].(string); ok {
        data.LldpNeighbors = NewJSONSubsetValue(val)
    } else {
        data.LldpNeighbors = NewJSONSubsetNull()
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
    if val, ok := dataMap["interfacesTotal"].(float64); ok {
        data.InterfacesTotal = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["interfacesTotal"].(int); ok {
        data.InterfacesTotal = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["interfacesTotal"].(int64); ok {
        data.InterfacesTotal = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["interfacesTotal"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.InterfacesTotal = types.NumberValue(big.NewFloat(val))
        } else {
            data.InterfacesTotal = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.InterfacesTotal = types.NumberNull()
    }
    if val, ok := dataMap["interfacesUp"].(float64); ok {
        data.InterfacesUp = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["interfacesUp"].(int); ok {
        data.InterfacesUp = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["interfacesUp"].(int64); ok {
        data.InterfacesUp = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["interfacesUp"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.InterfacesUp = types.NumberValue(big.NewFloat(val))
        } else {
            data.InterfacesUp = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.InterfacesUp = types.NumberNull()
    }
    if val, ok := dataMap["interfacesDown"].(float64); ok {
        data.InterfacesDown = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["interfacesDown"].(int); ok {
        data.InterfacesDown = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["interfacesDown"].(int64); ok {
        data.InterfacesDown = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["interfacesDown"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.InterfacesDown = types.NumberValue(big.NewFloat(val))
        } else {
            data.InterfacesDown = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.InterfacesDown = types.NumberNull()
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
    if obj, ok := dataMap["deletedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.DeletedAt = NewRFC3339Value(val)
        } else {
            data.DeletedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["deletedAt"].(string); ok && val != "" {
        data.DeletedAt = NewRFC3339Value(val)
    } else {
        data.DeletedAt = NewRFC3339Null()
    }
    if val, ok := dataMap["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["version"].(int); ok {
        data.Version = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["version"].(int64); ok {
        data.Version = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["version"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.Version = types.NumberValue(big.NewFloat(val))
        } else {
            data.Version = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.Version = types.NumberNull()
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
    } else if val, ok := dataMap["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
    }
    if obj, ok := dataMap["lastWalkLog"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastWalkLog = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.LastWalkLog = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.LastWalkLog = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.LastWalkLog = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.LastWalkLog = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.LastWalkLog = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.LastWalkLog = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.LastWalkLog = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.LastWalkLog = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["lastWalkLog"].(string); ok {
        data.LastWalkLog = NewJSONSubsetValue(val)
    } else {
        data.LastWalkLog = NewJSONSubsetNull()
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
    } else if val, ok := dataMap["archivedByUserId"].(string); ok {
        data.ArchivedByUserId = types.StringValue(val)
    } else {
        data.ArchivedByUserId = types.StringNull()
    }
    if obj, ok := dataMap["deletedByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeletedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DeletedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DeletedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DeletedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.DeletedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DeletedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.DeletedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DeletedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.DeletedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["deletedByUserId"].(string); ok {
        data.DeletedByUserId = types.StringValue(val)
    } else {
        data.DeletedByUserId = types.StringNull()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NetworkDeviceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var data NetworkDeviceResourceModel
    var state NetworkDeviceResourceModel

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

    // Create API request body
    networkDeviceRequest := map[string]interface{}{
        "data": map[string]interface{}{},
    }
    requestDataMap := networkDeviceRequest["data"].(map[string]interface{})

    if !data.Name.IsUnknown() && !state.Name.IsUnknown() && !data.Name.Equal(state.Name) {
        requestDataMap["name"] = data.Name.ValueString()
    }
    if !data.Description.IsUnknown() && !state.Description.IsUnknown() && !data.Description.Equal(state.Description) {
        requestDataMap["description"] = data.Description.ValueString()
    }
    if !data.Hostname.IsUnknown() && !state.Hostname.IsUnknown() && !data.Hostname.Equal(state.Hostname) {
        requestDataMap["hostname"] = data.Hostname.ValueString()
    }
    if !data.ProbeId.IsUnknown() && !state.ProbeId.IsUnknown() && !data.ProbeId.Equal(state.ProbeId) {
        requestDataMap["probeId"] = data.ProbeId.ValueString()
    }
    if !data.SiteId.IsUnknown() && !state.SiteId.IsUnknown() && !data.SiteId.Equal(state.SiteId) {
        requestDataMap["siteId"] = data.SiteId.ValueString()
    }
    if !data.CurrentMonitorStatusId.IsUnknown() && !state.CurrentMonitorStatusId.IsUnknown() && !data.CurrentMonitorStatusId.Equal(state.CurrentMonitorStatusId) {
        requestDataMap["currentMonitorStatusId"] = data.CurrentMonitorStatusId.ValueString()
    }
    if !data.MonitoringMethod.IsUnknown() && !state.MonitoringMethod.IsUnknown() && !data.MonitoringMethod.Equal(state.MonitoringMethod) {
        requestDataMap["monitoringMethod"] = data.MonitoringMethod.ValueString()
    }
    if !data.DeviceRole.IsUnknown() && !state.DeviceRole.IsUnknown() && !data.DeviceRole.Equal(state.DeviceRole) {
        requestDataMap["deviceRole"] = data.DeviceRole.ValueString()
    }
    if !data.MonitorId.IsUnknown() && !state.MonitorId.IsUnknown() && !data.MonitorId.Equal(state.MonitorId) {
        requestDataMap["monitorId"] = data.MonitorId.ValueString()
    }
    if !data.SnmpVersion.IsUnknown() && !state.SnmpVersion.IsUnknown() && !data.SnmpVersion.Equal(state.SnmpVersion) {
        requestDataMap["snmpVersion"] = data.SnmpVersion.ValueString()
    }
    if !data.SnmpCommunityString.IsUnknown() && !state.SnmpCommunityString.IsUnknown() && !data.SnmpCommunityString.Equal(state.SnmpCommunityString) {
        requestDataMap["snmpCommunityString"] = data.SnmpCommunityString.ValueString()
    }
    if !data.SnmpPort.IsUnknown() && !state.SnmpPort.IsUnknown() && !data.SnmpPort.Equal(state.SnmpPort) {
        requestDataMap["snmpPort"] = r.bigFloatToFloat64(data.SnmpPort.ValueBigFloat())
    }
    if !data.SnmpV3Auth.IsUnknown() && !state.SnmpV3Auth.IsUnknown() && !data.SnmpV3Auth.Equal(state.SnmpV3Auth) {
        var snmpv3authData interface{}
        if err := json.Unmarshal([]byte(data.SnmpV3Auth.ValueString()), &snmpv3authData); err == nil {
            requestDataMap["snmpV3Auth"] = snmpv3authData
        } else {
            requestDataMap["snmpV3Auth"] = data.SnmpV3Auth.ValueString()
        }
    }
    if !data.SnmpV3SecurityLevel.IsUnknown() && !state.SnmpV3SecurityLevel.IsUnknown() && !data.SnmpV3SecurityLevel.Equal(state.SnmpV3SecurityLevel) {
        requestDataMap["snmpV3SecurityLevel"] = data.SnmpV3SecurityLevel.ValueString()
    }
    if !data.SnmpV3Username.IsUnknown() && !state.SnmpV3Username.IsUnknown() && !data.SnmpV3Username.Equal(state.SnmpV3Username) {
        requestDataMap["snmpV3Username"] = data.SnmpV3Username.ValueString()
    }
    if !data.SnmpV3AuthProtocol.IsUnknown() && !state.SnmpV3AuthProtocol.IsUnknown() && !data.SnmpV3AuthProtocol.Equal(state.SnmpV3AuthProtocol) {
        requestDataMap["snmpV3AuthProtocol"] = data.SnmpV3AuthProtocol.ValueString()
    }
    if !data.SnmpV3AuthKey.IsUnknown() && !state.SnmpV3AuthKey.IsUnknown() && !data.SnmpV3AuthKey.Equal(state.SnmpV3AuthKey) {
        requestDataMap["snmpV3AuthKey"] = data.SnmpV3AuthKey.ValueString()
    }
    if !data.SnmpV3PrivProtocol.IsUnknown() && !state.SnmpV3PrivProtocol.IsUnknown() && !data.SnmpV3PrivProtocol.Equal(state.SnmpV3PrivProtocol) {
        requestDataMap["snmpV3PrivProtocol"] = data.SnmpV3PrivProtocol.ValueString()
    }
    if !data.SnmpV3PrivKey.IsUnknown() && !state.SnmpV3PrivKey.IsUnknown() && !data.SnmpV3PrivKey.Equal(state.SnmpV3PrivKey) {
        requestDataMap["snmpV3PrivKey"] = data.SnmpV3PrivKey.ValueString()
    }
    if !data.IsPollingEnabled.IsUnknown() && !state.IsPollingEnabled.IsUnknown() && !data.IsPollingEnabled.Equal(state.IsPollingEnabled) {
        requestDataMap["isPollingEnabled"] = data.IsPollingEnabled.ValueBool()
    }
    if !data.PollingIntervalInMinutes.IsUnknown() && !state.PollingIntervalInMinutes.IsUnknown() && !data.PollingIntervalInMinutes.Equal(state.PollingIntervalInMinutes) {
        requestDataMap["pollingIntervalInMinutes"] = r.bigFloatToFloat64(data.PollingIntervalInMinutes.ValueBigFloat())
    }
    if !data.WalkInterfaces.IsUnknown() && !state.WalkInterfaces.IsUnknown() && !data.WalkInterfaces.Equal(state.WalkInterfaces) {
        requestDataMap["walkInterfaces"] = data.WalkInterfaces.ValueBool()
    }
    if !data.CollectEndpoints.IsUnknown() && !state.CollectEndpoints.IsUnknown() && !data.CollectEndpoints.Equal(state.CollectEndpoints) {
        requestDataMap["collectEndpoints"] = data.CollectEndpoints.ValueBool()
    }
    if !data.SnmpOids.IsUnknown() && !state.SnmpOids.IsUnknown() && !data.SnmpOids.Equal(state.SnmpOids) {
        var snmpoidsData interface{}
        if err := json.Unmarshal([]byte(data.SnmpOids.ValueString()), &snmpoidsData); err == nil {
            requestDataMap["snmpOids"] = snmpoidsData
        } else {
            requestDataMap["snmpOids"] = data.SnmpOids.ValueString()
        }
    }
    if !data.NextPollAt.IsUnknown() && !state.NextPollAt.IsUnknown() && !data.NextPollAt.Equal(state.NextPollAt) {
        requestDataMap["nextPollAt"] = data.NextPollAt.ValueString()
    }
    if !data.SysDescr.IsUnknown() && !state.SysDescr.IsUnknown() && !data.SysDescr.Equal(state.SysDescr) {
        requestDataMap["sysDescr"] = data.SysDescr.ValueString()
    }
    if !data.SysName.IsUnknown() && !state.SysName.IsUnknown() && !data.SysName.Equal(state.SysName) {
        requestDataMap["sysName"] = data.SysName.ValueString()
    }
    if !data.SysObjectId.IsUnknown() && !state.SysObjectId.IsUnknown() && !data.SysObjectId.Equal(state.SysObjectId) {
        requestDataMap["sysObjectId"] = data.SysObjectId.ValueString()
    }
    if !data.SysLocation.IsUnknown() && !state.SysLocation.IsUnknown() && !data.SysLocation.Equal(state.SysLocation) {
        requestDataMap["sysLocation"] = data.SysLocation.ValueString()
    }
    if !data.SysContact.IsUnknown() && !state.SysContact.IsUnknown() && !data.SysContact.Equal(state.SysContact) {
        requestDataMap["sysContact"] = data.SysContact.ValueString()
    }
    if !data.Vendor.IsUnknown() && !state.Vendor.IsUnknown() && !data.Vendor.Equal(state.Vendor) {
        requestDataMap["vendor"] = data.Vendor.ValueString()
    }
    if !data.DeviceModel.IsUnknown() && !state.DeviceModel.IsUnknown() && !data.DeviceModel.Equal(state.DeviceModel) {
        requestDataMap["deviceModel"] = data.DeviceModel.ValueString()
    }
    if !data.SerialNumber.IsUnknown() && !state.SerialNumber.IsUnknown() && !data.SerialNumber.Equal(state.SerialNumber) {
        requestDataMap["serialNumber"] = data.SerialNumber.ValueString()
    }
    if !data.FirmwareVersion.IsUnknown() && !state.FirmwareVersion.IsUnknown() && !data.FirmwareVersion.Equal(state.FirmwareVersion) {
        requestDataMap["firmwareVersion"] = data.FirmwareVersion.ValueString()
    }
    if !data.SoftwareVersion.IsUnknown() && !state.SoftwareVersion.IsUnknown() && !data.SoftwareVersion.Equal(state.SoftwareVersion) {
        requestDataMap["softwareVersion"] = data.SoftwareVersion.ValueString()
    }
    if !data.LastRebootedAt.IsUnknown() && !state.LastRebootedAt.IsUnknown() && !data.LastRebootedAt.Equal(state.LastRebootedAt) {
        requestDataMap["lastRebootedAt"] = data.LastRebootedAt.ValueString()
    }
    if !data.CdpNeighbors.IsUnknown() && !state.CdpNeighbors.IsUnknown() && !data.CdpNeighbors.Equal(state.CdpNeighbors) {
        var cdpneighborsData interface{}
        if err := json.Unmarshal([]byte(data.CdpNeighbors.ValueString()), &cdpneighborsData); err == nil {
            requestDataMap["cdpNeighbors"] = cdpneighborsData
        } else {
            requestDataMap["cdpNeighbors"] = data.CdpNeighbors.ValueString()
        }
    }
    if !data.LldpNeighbors.IsUnknown() && !state.LldpNeighbors.IsUnknown() && !data.LldpNeighbors.Equal(state.LldpNeighbors) {
        var lldpneighborsData interface{}
        if err := json.Unmarshal([]byte(data.LldpNeighbors.ValueString()), &lldpneighborsData); err == nil {
            requestDataMap["lldpNeighbors"] = lldpneighborsData
        } else {
            requestDataMap["lldpNeighbors"] = data.LldpNeighbors.ValueString()
        }
    }
    if !data.LastSeenAt.IsUnknown() && !state.LastSeenAt.IsUnknown() && !data.LastSeenAt.Equal(state.LastSeenAt) {
        requestDataMap["lastSeenAt"] = data.LastSeenAt.ValueString()
    }
    if !data.InterfacesTotal.IsUnknown() && !state.InterfacesTotal.IsUnknown() && !data.InterfacesTotal.Equal(state.InterfacesTotal) {
        requestDataMap["interfacesTotal"] = r.bigFloatToFloat64(data.InterfacesTotal.ValueBigFloat())
    }
    if !data.InterfacesUp.IsUnknown() && !state.InterfacesUp.IsUnknown() && !data.InterfacesUp.Equal(state.InterfacesUp) {
        requestDataMap["interfacesUp"] = r.bigFloatToFloat64(data.InterfacesUp.ValueBigFloat())
    }
    if !data.InterfacesDown.IsUnknown() && !state.InterfacesDown.IsUnknown() && !data.InterfacesDown.Equal(state.InterfacesDown) {
        requestDataMap["interfacesDown"] = r.bigFloatToFloat64(data.InterfacesDown.ValueBigFloat())
    }
    if !data.IsArchived.IsUnknown() && !state.IsArchived.IsUnknown() && !data.IsArchived.Equal(state.IsArchived) {
        requestDataMap["isArchived"] = data.IsArchived.ValueBool()
    }
    if !data.Labels.IsUnknown() && !state.Labels.IsUnknown() && !data.Labels.Equal(state.Labels) {
        requestDataMap["labels"] = r.convertTerraformSetToInterface(data.Labels)
    }

    // Only call the API when there are changed fields to send. An empty
    // update body is rejected by the API; state is still refreshed below so
    // this method never writes unverified plan values into state.
    if len(networkDeviceRequest["data"].(map[string]interface{})) > 0 {
        httpResp, err := r.client.Put(ctx, "/network-device/" + data.Id.ValueString() + "", networkDeviceRequest)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update network_device, got error: %s", err))
            return
        }

        // Parse the update response
        var networkDeviceResponse map[string]interface{}
        err = r.client.ParseResponse(httpResp, &networkDeviceResponse)
        if err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to update network_device: %s", err))
            return
        }
        _ = networkDeviceResponse
    }

    // After successful update, fetch the current state by calling Read with select parameter
    selectParam := map[string]interface{}{
        "projectId": true,
        "name": true,
        "description": true,
        "hostname": true,
        "probeId": true,
        "siteId": true,
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
        "createdByUserId": true,
        "isArchived": true,
        "labels": true,
        "currentMonitorStatusId": true,
        "nextPollAt": true,
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
        "interfacesTotal": true,
        "interfacesUp": true,
        "interfacesDown": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "slug": true,
        "lastWalkLog": true,
        "archivedAt": true,
        "archivedByUserId": true,
        "deletedByUserId": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect(ctx, "/network-device/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read network_device after update, got error: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read network_device after update: %s", err))
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
    } else if val, ok := dataMap["description"].(string); ok {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if obj, ok := dataMap["hostname"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Hostname = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Hostname = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Hostname = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Hostname = types.StringValue(string(jsonBytes))
            } else {
                data.Hostname = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Hostname = types.StringValue(string(jsonBytes))
            } else {
                data.Hostname = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Hostname = types.StringValue(string(jsonBytes))
        } else {
            data.Hostname = types.StringNull()
        }
    } else if val, ok := dataMap["hostname"].(string); ok {
        data.Hostname = types.StringValue(val)
    } else {
        data.Hostname = types.StringNull()
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
    if obj, ok := dataMap["siteId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SiteId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SiteId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SiteId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SiteId = types.StringValue(string(jsonBytes))
            } else {
                data.SiteId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SiteId = types.StringValue(string(jsonBytes))
            } else {
                data.SiteId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SiteId = types.StringValue(string(jsonBytes))
        } else {
            data.SiteId = types.StringNull()
        }
    } else if val, ok := dataMap["siteId"].(string); ok {
        data.SiteId = types.StringValue(val)
    } else {
        data.SiteId = types.StringNull()
    }
    if obj, ok := dataMap["monitoringMethod"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MonitoringMethod = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.MonitoringMethod = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.MonitoringMethod = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.MonitoringMethod = types.StringValue(string(jsonBytes))
            } else {
                data.MonitoringMethod = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.MonitoringMethod = types.StringValue(string(jsonBytes))
            } else {
                data.MonitoringMethod = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.MonitoringMethod = types.StringValue(string(jsonBytes))
        } else {
            data.MonitoringMethod = types.StringNull()
        }
    } else if val, ok := dataMap["monitoringMethod"].(string); ok {
        data.MonitoringMethod = types.StringValue(val)
    } else {
        data.MonitoringMethod = types.StringNull()
    }
    if obj, ok := dataMap["deviceRole"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeviceRole = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DeviceRole = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DeviceRole = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DeviceRole = types.StringValue(string(jsonBytes))
            } else {
                data.DeviceRole = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DeviceRole = types.StringValue(string(jsonBytes))
            } else {
                data.DeviceRole = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DeviceRole = types.StringValue(string(jsonBytes))
        } else {
            data.DeviceRole = types.StringNull()
        }
    } else if val, ok := dataMap["deviceRole"].(string); ok {
        data.DeviceRole = types.StringValue(val)
    } else {
        data.DeviceRole = types.StringNull()
    }
    if obj, ok := dataMap["monitorId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MonitorId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.MonitorId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.MonitorId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.MonitorId = types.StringValue(string(jsonBytes))
            } else {
                data.MonitorId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.MonitorId = types.StringValue(string(jsonBytes))
            } else {
                data.MonitorId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.MonitorId = types.StringValue(string(jsonBytes))
        } else {
            data.MonitorId = types.StringNull()
        }
    } else if val, ok := dataMap["monitorId"].(string); ok {
        data.MonitorId = types.StringValue(val)
    } else {
        data.MonitorId = types.StringNull()
    }
    if obj, ok := dataMap["snmpVersion"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SnmpVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SnmpVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SnmpVersion = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SnmpVersion = types.StringValue(string(jsonBytes))
            } else {
                data.SnmpVersion = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SnmpVersion = types.StringValue(string(jsonBytes))
            } else {
                data.SnmpVersion = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SnmpVersion = types.StringValue(string(jsonBytes))
        } else {
            data.SnmpVersion = types.StringNull()
        }
    } else if val, ok := dataMap["snmpVersion"].(string); ok {
        data.SnmpVersion = types.StringValue(val)
    } else {
        data.SnmpVersion = types.StringNull()
    }
    if obj, ok := dataMap["snmpCommunityString"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SnmpCommunityString = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SnmpCommunityString = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SnmpCommunityString = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SnmpCommunityString = types.StringValue(string(jsonBytes))
            } else {
                data.SnmpCommunityString = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SnmpCommunityString = types.StringValue(string(jsonBytes))
            } else {
                data.SnmpCommunityString = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SnmpCommunityString = types.StringValue(string(jsonBytes))
        } else {
            data.SnmpCommunityString = types.StringNull()
        }
    } else if val, ok := dataMap["snmpCommunityString"].(string); ok {
        data.SnmpCommunityString = types.StringValue(val)
    } else {
        data.SnmpCommunityString = types.StringNull()
    }
    if val, ok := dataMap["snmpPort"].(float64); ok {
        data.SnmpPort = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["snmpPort"].(int); ok {
        data.SnmpPort = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["snmpPort"].(int64); ok {
        data.SnmpPort = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["snmpPort"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.SnmpPort = types.NumberValue(big.NewFloat(val))
        } else {
            data.SnmpPort = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.SnmpPort = types.NumberNull()
    }
    if obj, ok := dataMap["snmpV3Auth"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SnmpV3Auth = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SnmpV3Auth = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SnmpV3Auth = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SnmpV3Auth = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.SnmpV3Auth = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SnmpV3Auth = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.SnmpV3Auth = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SnmpV3Auth = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.SnmpV3Auth = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["snmpV3Auth"].(string); ok {
        data.SnmpV3Auth = NewJSONSubsetValue(val)
    } else {
        data.SnmpV3Auth = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["snmpV3SecurityLevel"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SnmpV3SecurityLevel = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SnmpV3SecurityLevel = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SnmpV3SecurityLevel = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SnmpV3SecurityLevel = types.StringValue(string(jsonBytes))
            } else {
                data.SnmpV3SecurityLevel = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SnmpV3SecurityLevel = types.StringValue(string(jsonBytes))
            } else {
                data.SnmpV3SecurityLevel = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SnmpV3SecurityLevel = types.StringValue(string(jsonBytes))
        } else {
            data.SnmpV3SecurityLevel = types.StringNull()
        }
    } else if val, ok := dataMap["snmpV3SecurityLevel"].(string); ok {
        data.SnmpV3SecurityLevel = types.StringValue(val)
    } else {
        data.SnmpV3SecurityLevel = types.StringNull()
    }
    if obj, ok := dataMap["snmpV3Username"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SnmpV3Username = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SnmpV3Username = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SnmpV3Username = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SnmpV3Username = types.StringValue(string(jsonBytes))
            } else {
                data.SnmpV3Username = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SnmpV3Username = types.StringValue(string(jsonBytes))
            } else {
                data.SnmpV3Username = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SnmpV3Username = types.StringValue(string(jsonBytes))
        } else {
            data.SnmpV3Username = types.StringNull()
        }
    } else if val, ok := dataMap["snmpV3Username"].(string); ok {
        data.SnmpV3Username = types.StringValue(val)
    } else {
        data.SnmpV3Username = types.StringNull()
    }
    if obj, ok := dataMap["snmpV3AuthProtocol"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SnmpV3AuthProtocol = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SnmpV3AuthProtocol = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SnmpV3AuthProtocol = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SnmpV3AuthProtocol = types.StringValue(string(jsonBytes))
            } else {
                data.SnmpV3AuthProtocol = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SnmpV3AuthProtocol = types.StringValue(string(jsonBytes))
            } else {
                data.SnmpV3AuthProtocol = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SnmpV3AuthProtocol = types.StringValue(string(jsonBytes))
        } else {
            data.SnmpV3AuthProtocol = types.StringNull()
        }
    } else if val, ok := dataMap["snmpV3AuthProtocol"].(string); ok {
        data.SnmpV3AuthProtocol = types.StringValue(val)
    } else {
        data.SnmpV3AuthProtocol = types.StringNull()
    }
    if obj, ok := dataMap["snmpV3AuthKey"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SnmpV3AuthKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SnmpV3AuthKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SnmpV3AuthKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SnmpV3AuthKey = types.StringValue(string(jsonBytes))
            } else {
                data.SnmpV3AuthKey = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SnmpV3AuthKey = types.StringValue(string(jsonBytes))
            } else {
                data.SnmpV3AuthKey = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SnmpV3AuthKey = types.StringValue(string(jsonBytes))
        } else {
            data.SnmpV3AuthKey = types.StringNull()
        }
    } else if val, ok := dataMap["snmpV3AuthKey"].(string); ok {
        data.SnmpV3AuthKey = types.StringValue(val)
    } else {
        data.SnmpV3AuthKey = types.StringNull()
    }
    if obj, ok := dataMap["snmpV3PrivProtocol"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SnmpV3PrivProtocol = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SnmpV3PrivProtocol = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SnmpV3PrivProtocol = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SnmpV3PrivProtocol = types.StringValue(string(jsonBytes))
            } else {
                data.SnmpV3PrivProtocol = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SnmpV3PrivProtocol = types.StringValue(string(jsonBytes))
            } else {
                data.SnmpV3PrivProtocol = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SnmpV3PrivProtocol = types.StringValue(string(jsonBytes))
        } else {
            data.SnmpV3PrivProtocol = types.StringNull()
        }
    } else if val, ok := dataMap["snmpV3PrivProtocol"].(string); ok {
        data.SnmpV3PrivProtocol = types.StringValue(val)
    } else {
        data.SnmpV3PrivProtocol = types.StringNull()
    }
    if obj, ok := dataMap["snmpV3PrivKey"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SnmpV3PrivKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SnmpV3PrivKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SnmpV3PrivKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SnmpV3PrivKey = types.StringValue(string(jsonBytes))
            } else {
                data.SnmpV3PrivKey = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SnmpV3PrivKey = types.StringValue(string(jsonBytes))
            } else {
                data.SnmpV3PrivKey = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SnmpV3PrivKey = types.StringValue(string(jsonBytes))
        } else {
            data.SnmpV3PrivKey = types.StringNull()
        }
    } else if val, ok := dataMap["snmpV3PrivKey"].(string); ok {
        data.SnmpV3PrivKey = types.StringValue(val)
    } else {
        data.SnmpV3PrivKey = types.StringNull()
    }
    if val, ok := dataMap["isPollingEnabled"].(bool); ok {
        data.IsPollingEnabled = types.BoolValue(val)
    }
    if val, ok := dataMap["pollingIntervalInMinutes"].(float64); ok {
        data.PollingIntervalInMinutes = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["pollingIntervalInMinutes"].(int); ok {
        data.PollingIntervalInMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["pollingIntervalInMinutes"].(int64); ok {
        data.PollingIntervalInMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["pollingIntervalInMinutes"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.PollingIntervalInMinutes = types.NumberValue(big.NewFloat(val))
        } else {
            data.PollingIntervalInMinutes = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.PollingIntervalInMinutes = types.NumberNull()
    }
    if val, ok := dataMap["walkInterfaces"].(bool); ok {
        data.WalkInterfaces = types.BoolValue(val)
    }
    if val, ok := dataMap["collectEndpoints"].(bool); ok {
        data.CollectEndpoints = types.BoolValue(val)
    }
    if obj, ok := dataMap["snmpOids"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SnmpOids = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SnmpOids = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SnmpOids = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SnmpOids = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.SnmpOids = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SnmpOids = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.SnmpOids = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SnmpOids = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.SnmpOids = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["snmpOids"].(string); ok {
        data.SnmpOids = NewJSONSubsetValue(val)
    } else {
        data.SnmpOids = NewJSONSubsetNull()
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
    if obj, ok := dataMap["currentMonitorStatusId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CurrentMonitorStatusId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CurrentMonitorStatusId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CurrentMonitorStatusId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CurrentMonitorStatusId = types.StringValue(string(jsonBytes))
            } else {
                data.CurrentMonitorStatusId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CurrentMonitorStatusId = types.StringValue(string(jsonBytes))
            } else {
                data.CurrentMonitorStatusId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CurrentMonitorStatusId = types.StringValue(string(jsonBytes))
        } else {
            data.CurrentMonitorStatusId = types.StringNull()
        }
    } else if val, ok := dataMap["currentMonitorStatusId"].(string); ok {
        data.CurrentMonitorStatusId = types.StringValue(val)
    } else {
        data.CurrentMonitorStatusId = types.StringNull()
    }
    if obj, ok := dataMap["nextPollAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.NextPollAt = NewRFC3339Value(val)
        } else {
            data.NextPollAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["nextPollAt"].(string); ok && val != "" {
        data.NextPollAt = NewRFC3339Value(val)
    } else {
        data.NextPollAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["sysDescr"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SysDescr = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SysDescr = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SysDescr = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SysDescr = types.StringValue(string(jsonBytes))
            } else {
                data.SysDescr = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SysDescr = types.StringValue(string(jsonBytes))
            } else {
                data.SysDescr = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SysDescr = types.StringValue(string(jsonBytes))
        } else {
            data.SysDescr = types.StringNull()
        }
    } else if val, ok := dataMap["sysDescr"].(string); ok {
        data.SysDescr = types.StringValue(val)
    } else {
        data.SysDescr = types.StringNull()
    }
    if obj, ok := dataMap["sysName"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SysName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SysName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SysName = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SysName = types.StringValue(string(jsonBytes))
            } else {
                data.SysName = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SysName = types.StringValue(string(jsonBytes))
            } else {
                data.SysName = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SysName = types.StringValue(string(jsonBytes))
        } else {
            data.SysName = types.StringNull()
        }
    } else if val, ok := dataMap["sysName"].(string); ok {
        data.SysName = types.StringValue(val)
    } else {
        data.SysName = types.StringNull()
    }
    if obj, ok := dataMap["sysObjectId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SysObjectId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SysObjectId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SysObjectId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SysObjectId = types.StringValue(string(jsonBytes))
            } else {
                data.SysObjectId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SysObjectId = types.StringValue(string(jsonBytes))
            } else {
                data.SysObjectId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SysObjectId = types.StringValue(string(jsonBytes))
        } else {
            data.SysObjectId = types.StringNull()
        }
    } else if val, ok := dataMap["sysObjectId"].(string); ok {
        data.SysObjectId = types.StringValue(val)
    } else {
        data.SysObjectId = types.StringNull()
    }
    if obj, ok := dataMap["sysLocation"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SysLocation = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SysLocation = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SysLocation = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SysLocation = types.StringValue(string(jsonBytes))
            } else {
                data.SysLocation = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SysLocation = types.StringValue(string(jsonBytes))
            } else {
                data.SysLocation = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SysLocation = types.StringValue(string(jsonBytes))
        } else {
            data.SysLocation = types.StringNull()
        }
    } else if val, ok := dataMap["sysLocation"].(string); ok {
        data.SysLocation = types.StringValue(val)
    } else {
        data.SysLocation = types.StringNull()
    }
    if obj, ok := dataMap["sysContact"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SysContact = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SysContact = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SysContact = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SysContact = types.StringValue(string(jsonBytes))
            } else {
                data.SysContact = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SysContact = types.StringValue(string(jsonBytes))
            } else {
                data.SysContact = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SysContact = types.StringValue(string(jsonBytes))
        } else {
            data.SysContact = types.StringNull()
        }
    } else if val, ok := dataMap["sysContact"].(string); ok {
        data.SysContact = types.StringValue(val)
    } else {
        data.SysContact = types.StringNull()
    }
    if obj, ok := dataMap["vendor"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Vendor = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Vendor = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Vendor = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Vendor = types.StringValue(string(jsonBytes))
            } else {
                data.Vendor = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Vendor = types.StringValue(string(jsonBytes))
            } else {
                data.Vendor = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Vendor = types.StringValue(string(jsonBytes))
        } else {
            data.Vendor = types.StringNull()
        }
    } else if val, ok := dataMap["vendor"].(string); ok {
        data.Vendor = types.StringValue(val)
    } else {
        data.Vendor = types.StringNull()
    }
    if obj, ok := dataMap["deviceModel"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeviceModel = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DeviceModel = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DeviceModel = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DeviceModel = types.StringValue(string(jsonBytes))
            } else {
                data.DeviceModel = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DeviceModel = types.StringValue(string(jsonBytes))
            } else {
                data.DeviceModel = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DeviceModel = types.StringValue(string(jsonBytes))
        } else {
            data.DeviceModel = types.StringNull()
        }
    } else if val, ok := dataMap["deviceModel"].(string); ok {
        data.DeviceModel = types.StringValue(val)
    } else {
        data.DeviceModel = types.StringNull()
    }
    if obj, ok := dataMap["serialNumber"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SerialNumber = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SerialNumber = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SerialNumber = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SerialNumber = types.StringValue(string(jsonBytes))
            } else {
                data.SerialNumber = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SerialNumber = types.StringValue(string(jsonBytes))
            } else {
                data.SerialNumber = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SerialNumber = types.StringValue(string(jsonBytes))
        } else {
            data.SerialNumber = types.StringNull()
        }
    } else if val, ok := dataMap["serialNumber"].(string); ok {
        data.SerialNumber = types.StringValue(val)
    } else {
        data.SerialNumber = types.StringNull()
    }
    if obj, ok := dataMap["firmwareVersion"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FirmwareVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.FirmwareVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.FirmwareVersion = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.FirmwareVersion = types.StringValue(string(jsonBytes))
            } else {
                data.FirmwareVersion = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.FirmwareVersion = types.StringValue(string(jsonBytes))
            } else {
                data.FirmwareVersion = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.FirmwareVersion = types.StringValue(string(jsonBytes))
        } else {
            data.FirmwareVersion = types.StringNull()
        }
    } else if val, ok := dataMap["firmwareVersion"].(string); ok {
        data.FirmwareVersion = types.StringValue(val)
    } else {
        data.FirmwareVersion = types.StringNull()
    }
    if obj, ok := dataMap["softwareVersion"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SoftwareVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SoftwareVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SoftwareVersion = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SoftwareVersion = types.StringValue(string(jsonBytes))
            } else {
                data.SoftwareVersion = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SoftwareVersion = types.StringValue(string(jsonBytes))
            } else {
                data.SoftwareVersion = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SoftwareVersion = types.StringValue(string(jsonBytes))
        } else {
            data.SoftwareVersion = types.StringNull()
        }
    } else if val, ok := dataMap["softwareVersion"].(string); ok {
        data.SoftwareVersion = types.StringValue(val)
    } else {
        data.SoftwareVersion = types.StringNull()
    }
    if obj, ok := dataMap["lastRebootedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.LastRebootedAt = NewRFC3339Value(val)
        } else {
            data.LastRebootedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["lastRebootedAt"].(string); ok && val != "" {
        data.LastRebootedAt = NewRFC3339Value(val)
    } else {
        data.LastRebootedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["cdpNeighbors"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CdpNeighbors = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CdpNeighbors = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CdpNeighbors = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CdpNeighbors = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.CdpNeighbors = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CdpNeighbors = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.CdpNeighbors = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CdpNeighbors = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.CdpNeighbors = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["cdpNeighbors"].(string); ok {
        data.CdpNeighbors = NewJSONSubsetValue(val)
    } else {
        data.CdpNeighbors = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["lldpNeighbors"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LldpNeighbors = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.LldpNeighbors = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.LldpNeighbors = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.LldpNeighbors = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.LldpNeighbors = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.LldpNeighbors = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.LldpNeighbors = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.LldpNeighbors = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.LldpNeighbors = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["lldpNeighbors"].(string); ok {
        data.LldpNeighbors = NewJSONSubsetValue(val)
    } else {
        data.LldpNeighbors = NewJSONSubsetNull()
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
    if val, ok := dataMap["interfacesTotal"].(float64); ok {
        data.InterfacesTotal = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["interfacesTotal"].(int); ok {
        data.InterfacesTotal = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["interfacesTotal"].(int64); ok {
        data.InterfacesTotal = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["interfacesTotal"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.InterfacesTotal = types.NumberValue(big.NewFloat(val))
        } else {
            data.InterfacesTotal = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.InterfacesTotal = types.NumberNull()
    }
    if val, ok := dataMap["interfacesUp"].(float64); ok {
        data.InterfacesUp = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["interfacesUp"].(int); ok {
        data.InterfacesUp = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["interfacesUp"].(int64); ok {
        data.InterfacesUp = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["interfacesUp"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.InterfacesUp = types.NumberValue(big.NewFloat(val))
        } else {
            data.InterfacesUp = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.InterfacesUp = types.NumberNull()
    }
    if val, ok := dataMap["interfacesDown"].(float64); ok {
        data.InterfacesDown = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["interfacesDown"].(int); ok {
        data.InterfacesDown = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["interfacesDown"].(int64); ok {
        data.InterfacesDown = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["interfacesDown"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.InterfacesDown = types.NumberValue(big.NewFloat(val))
        } else {
            data.InterfacesDown = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.InterfacesDown = types.NumberNull()
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
    if obj, ok := dataMap["deletedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.DeletedAt = NewRFC3339Value(val)
        } else {
            data.DeletedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["deletedAt"].(string); ok && val != "" {
        data.DeletedAt = NewRFC3339Value(val)
    } else {
        data.DeletedAt = NewRFC3339Null()
    }
    if val, ok := dataMap["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["version"].(int); ok {
        data.Version = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["version"].(int64); ok {
        data.Version = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["version"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.Version = types.NumberValue(big.NewFloat(val))
        } else {
            data.Version = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.Version = types.NumberNull()
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
    } else if val, ok := dataMap["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
    }
    if obj, ok := dataMap["lastWalkLog"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastWalkLog = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.LastWalkLog = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.LastWalkLog = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.LastWalkLog = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.LastWalkLog = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.LastWalkLog = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.LastWalkLog = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.LastWalkLog = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.LastWalkLog = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["lastWalkLog"].(string); ok {
        data.LastWalkLog = NewJSONSubsetValue(val)
    } else {
        data.LastWalkLog = NewJSONSubsetNull()
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
    } else if val, ok := dataMap["archivedByUserId"].(string); ok {
        data.ArchivedByUserId = types.StringValue(val)
    } else {
        data.ArchivedByUserId = types.StringNull()
    }
    if obj, ok := dataMap["deletedByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeletedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DeletedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DeletedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DeletedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.DeletedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DeletedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.DeletedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DeletedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.DeletedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["deletedByUserId"].(string); ok {
        data.DeletedByUserId = types.StringValue(val)
    } else {
        data.DeletedByUserId = types.StringNull()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    data.Id = state.Id

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NetworkDeviceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var data NetworkDeviceResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Make API call
    httpResp, err := r.client.Delete(ctx, "/network-device/" + data.Id.ValueString() + "")
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete network_device, got error: %s", err))
        return
    }

    // A failed delete must keep the resource in state — silently dropping it
    // orphans real infrastructure. 404 means it is already gone.
    if httpResp.StatusCode >= 400 && httpResp.StatusCode != http.StatusNotFound {
        err = r.client.ParseResponse(httpResp, nil)
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to delete network_device: %s", err))
        return
    }
    if httpResp.Body != nil {
        httpResp.Body.Close()
    }
}


func (r *NetworkDeviceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Helper method to convert Terraform map to Go interface{}
func (r *NetworkDeviceResource) convertTerraformMapToInterface(terraformMap types.Map) interface{} {
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
func (r *NetworkDeviceResource) convertTerraformListToInterface(terraformList types.List) interface{} {
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
func (r *NetworkDeviceResource) convertTerraformSetToInterface(terraformSet types.Set) interface{} {
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
func (r *NetworkDeviceResource) parseJSONField(terraformString basetypes.StringValuable) interface{} {
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
func (r *NetworkDeviceResource) normalizeURLWrappers(value interface{}) interface{} {
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

func (r *NetworkDeviceResource) normalizeURLString(value string) string {
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
func (r *NetworkDeviceResource) bigFloatToFloat64(bf *big.Float) interface{} {
    if bf == nil {
        return nil
    }
    f, _ := bf.Float64()
    return f
}

// Helper method to check if a type string is a valid OneUptime ObjectType.
// The registry itself lives in objecttypes.go, shared across the package.
func (r *NetworkDeviceResource) isValidOneUptimeObjectType(typeStr string) bool {
    return validOneUptimeObjectTypes[typeStr]
}
