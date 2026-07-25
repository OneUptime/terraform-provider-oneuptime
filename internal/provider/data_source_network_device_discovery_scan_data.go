package provider

import (
    "context"
    "fmt"
    "math/big"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &NetworkDeviceDiscoveryScanDataDataSource{}

func NewNetworkDeviceDiscoveryScanDataDataSource() datasource.DataSource {
    return &NetworkDeviceDiscoveryScanDataDataSource{}
}

// NetworkDeviceDiscoveryScanDataDataSource defines the data source implementation.
type NetworkDeviceDiscoveryScanDataDataSource struct {
    client *Client
}

// NetworkDeviceDiscoveryScanDataDataSourceModel describes the data source data model.
type NetworkDeviceDiscoveryScanDataDataSourceModel struct {
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
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
}

func (d *NetworkDeviceDiscoveryScanDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_network_device_discovery_scan_data"
}

func (d *NetworkDeviceDiscoveryScanDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "network_device_discovery_scan_data data source",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Identifier to filter by",
                Optional: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Name to filter by",
                Optional: true,
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
                MarkdownDescription: "Subnet to scan in CIDR notation, e.g. 192.168.1.0/24. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Network Device Discovery Scan], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Device Discovery Scan], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "snmp_version": schema.StringAttribute{
                MarkdownDescription: "SNMP version tried against every host in the subnet (V1, V2c, V3). Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Network Device Discovery Scan], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Device Discovery Scan], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "snmp_community_string": schema.StringAttribute{
                MarkdownDescription: "Community string tried against every host in the subnet (SNMP v1/v2c). Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Network Device Discovery Scan], Read: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Read Network Device Discovery Scan], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "snmp_port": schema.NumberAttribute{
                MarkdownDescription: "UDP port tried against every host in the subnet. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Network Device Discovery Scan], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Device Discovery Scan], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "snmp_v3_security_level": schema.StringAttribute{
                MarkdownDescription: "SNMP v3 security level tried against every host: noAuthNoPriv, authNoPriv, or authPriv. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Network Device Discovery Scan], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Device Discovery Scan], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "snmp_v3_username": schema.StringAttribute{
                MarkdownDescription: "SNMP v3 security name (username) tried against every host. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Network Device Discovery Scan], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Device Discovery Scan], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "snmp_v3_auth_protocol": schema.StringAttribute{
                MarkdownDescription: "SNMP v3 authentication protocol: MD5, SHA, SHA256, or SHA512. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Network Device Discovery Scan], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Device Discovery Scan], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "snmp_v3_auth_key": schema.StringAttribute{
                MarkdownDescription: "SNMP v3 authentication passphrase tried against every host. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Network Device Discovery Scan], Read: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Read Network Device Discovery Scan], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "snmp_v3_priv_protocol": schema.StringAttribute{
                MarkdownDescription: "SNMP v3 privacy (encryption) protocol: DES, AES, or AES256. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Network Device Discovery Scan], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Device Discovery Scan], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "snmp_v3_priv_key": schema.StringAttribute{
                MarkdownDescription: "SNMP v3 privacy (encryption) passphrase tried against every host. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Network Device Discovery Scan], Read: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Read Network Device Discovery Scan], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "status": schema.StringAttribute{
                MarkdownDescription: "Status of this discovery scan: \"Pending\", \"In Progress\", \"Completed\" or \"Failed\". Managed by the scanning probe.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Device Discovery Scan], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "status_message": schema.StringAttribute{
                MarkdownDescription: "Details about the current status of this scan, e.g. the failure reason. Managed by the scanning probe.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Device Discovery Scan], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "discovered_devices": schema.StringAttribute{
                MarkdownDescription: "Devices found by this scan: array of {ipAddress, sysName, sysDescr, isAlreadyRegistered}. Managed by the scanning probe.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Device Discovery Scan], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "scanned_host_count": schema.NumberAttribute{
                MarkdownDescription: "Total number of host addresses swept in the subnet. Managed by the scanning probe.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Device Discovery Scan], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "responded_host_count": schema.NumberAttribute{
                MarkdownDescription: "Number of hosts that responded to SNMP during the sweep. Managed by the scanning probe.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Device Discovery Scan], Update: [No access - you don't have permission for this operation]",
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
                MarkdownDescription: "Re-run this scan automatically every Rescan Interval minutes to keep discovery continuous.. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Network Device Discovery Scan], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Device Discovery Scan], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Network Device Discovery Scan]",
                Computed: true,
            },
            "rescan_interval_in_minutes": schema.NumberAttribute{
                MarkdownDescription: "How often a recurring scan re-runs, in minutes. Ignored unless Is Recurring is on.. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Network Device Discovery Scan], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Device Discovery Scan], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Network Device Discovery Scan]",
                Computed: true,
            },
            "next_scan_at": schema.StringAttribute{
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

func (d *NetworkDeviceDiscoveryScanDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *NetworkDeviceDiscoveryScanDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data NetworkDeviceDiscoveryScanDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "network-device-discovery-scan" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read network_device_discovery_scan_data, got error: %s", err))
        return
    }

    var networkDeviceDiscoveryScanDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &networkDeviceDiscoveryScanDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse network_device_discovery_scan_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := networkDeviceDiscoveryScanDataResponse["data"].(map[string]interface{}); ok {
        networkDeviceDiscoveryScanDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := networkDeviceDiscoveryScanDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := networkDeviceDiscoveryScanDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := networkDeviceDiscoveryScanDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := networkDeviceDiscoveryScanDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := networkDeviceDiscoveryScanDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := networkDeviceDiscoveryScanDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := networkDeviceDiscoveryScanDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := networkDeviceDiscoveryScanDataResponse["probe_id"].(string); ok {
        data.ProbeId = types.StringValue(val)
    }
    if val, ok := networkDeviceDiscoveryScanDataResponse["cidr"].(string); ok {
        data.Cidr = types.StringValue(val)
    }
    if val, ok := networkDeviceDiscoveryScanDataResponse["snmp_version"].(string); ok {
        data.SnmpVersion = types.StringValue(val)
    }
    if val, ok := networkDeviceDiscoveryScanDataResponse["snmp_community_string"].(string); ok {
        data.SnmpCommunityString = types.StringValue(val)
    }
    if val, ok := networkDeviceDiscoveryScanDataResponse["snmp_port"].(float64); ok {
        data.SnmpPort = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := networkDeviceDiscoveryScanDataResponse["snmp_v3_security_level"].(string); ok {
        data.SnmpV3SecurityLevel = types.StringValue(val)
    }
    if val, ok := networkDeviceDiscoveryScanDataResponse["snmp_v3_username"].(string); ok {
        data.SnmpV3Username = types.StringValue(val)
    }
    if val, ok := networkDeviceDiscoveryScanDataResponse["snmp_v3_auth_protocol"].(string); ok {
        data.SnmpV3AuthProtocol = types.StringValue(val)
    }
    if val, ok := networkDeviceDiscoveryScanDataResponse["snmp_v3_auth_key"].(string); ok {
        data.SnmpV3AuthKey = types.StringValue(val)
    }
    if val, ok := networkDeviceDiscoveryScanDataResponse["snmp_v3_priv_protocol"].(string); ok {
        data.SnmpV3PrivProtocol = types.StringValue(val)
    }
    if val, ok := networkDeviceDiscoveryScanDataResponse["snmp_v3_priv_key"].(string); ok {
        data.SnmpV3PrivKey = types.StringValue(val)
    }
    if val, ok := networkDeviceDiscoveryScanDataResponse["status"].(string); ok {
        data.Status = types.StringValue(val)
    }
    if val, ok := networkDeviceDiscoveryScanDataResponse["status_message"].(string); ok {
        data.StatusMessage = types.StringValue(val)
    }
    if val, ok := networkDeviceDiscoveryScanDataResponse["discovered_devices"].(string); ok {
        data.DiscoveredDevices = types.StringValue(val)
    }
    if val, ok := networkDeviceDiscoveryScanDataResponse["scanned_host_count"].(float64); ok {
        data.ScannedHostCount = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := networkDeviceDiscoveryScanDataResponse["responded_host_count"].(float64); ok {
        data.RespondedHostCount = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := networkDeviceDiscoveryScanDataResponse["started_at"].(string); ok {
        data.StartedAt = types.StringValue(val)
    }
    if val, ok := networkDeviceDiscoveryScanDataResponse["completed_at"].(string); ok {
        data.CompletedAt = types.StringValue(val)
    }
    if val, ok := networkDeviceDiscoveryScanDataResponse["is_recurring"].(bool); ok {
        data.IsRecurring = types.BoolValue(val)
    }
    if val, ok := networkDeviceDiscoveryScanDataResponse["rescan_interval_in_minutes"].(float64); ok {
        data.RescanIntervalInMinutes = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := networkDeviceDiscoveryScanDataResponse["next_scan_at"].(string); ok {
        data.NextScanAt = types.StringValue(val)
    }
    if val, ok := networkDeviceDiscoveryScanDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := networkDeviceDiscoveryScanDataResponse["deleted_by_user_id"].(string); ok {
        data.DeletedByUserId = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
