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
var _ datasource.DataSource = &NetworkInterfaceDataDataSource{}

func NewNetworkInterfaceDataDataSource() datasource.DataSource {
    return &NetworkInterfaceDataDataSource{}
}

// NetworkInterfaceDataDataSource defines the data source implementation.
type NetworkInterfaceDataDataSource struct {
    client *Client
}

// NetworkInterfaceDataDataSourceModel describes the data source data model.
type NetworkInterfaceDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    NetworkDeviceId types.String `tfsdk:"network_device_id"`
    InterfaceIndex types.Number `tfsdk:"interface_index"`
    Alias types.String `tfsdk:"alias"`
    MacAddress types.String `tfsdk:"mac_address"`
    InterfaceType types.Number `tfsdk:"interface_type"`
    IsMonitored types.Bool `tfsdk:"is_monitored"`
    IsOperationallyUp types.Bool `tfsdk:"is_operationally_up"`
    IsAdministrativelyUp types.Bool `tfsdk:"is_administratively_up"`
    SpeedInMbps types.Number `tfsdk:"speed_in_mbps"`
    InRateMbps types.Number `tfsdk:"in_rate_mbps"`
    OutRateMbps types.Number `tfsdk:"out_rate_mbps"`
    UtilizationPercent types.Number `tfsdk:"utilization_percent"`
    ErrorsPerSecond types.Number `tfsdk:"errors_per_second"`
    LastSeenAt types.String `tfsdk:"last_seen_at"`
}

func (d *NetworkInterfaceDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_network_interface_data"
}

func (d *NetworkInterfaceDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "network_interface_data data source",

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
            "network_device_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "interface_index": schema.NumberAttribute{
                MarkdownDescription: "SNMP ifIndex of this interface on the device. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Device], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "alias": schema.StringAttribute{
                MarkdownDescription: "Interface alias (ifAlias) from SNMP. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Device], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "mac_address": schema.StringAttribute{
                MarkdownDescription: "Physical address (ifPhysAddress) from SNMP, colon-separated hex. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Device], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "interface_type": schema.NumberAttribute{
                MarkdownDescription: "IANAifType number (ifType) from SNMP — 6 = ethernetCsmacd, 24 = softwareLoopback. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Device], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "is_monitored": schema.BoolAttribute{
                MarkdownDescription: "Include this interface in down/utilization/error alerting.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Device], Update: [Project Owner, Project Admin, Project Member, Edit Network Device]",
                Computed: true,
            },
            "is_operationally_up": schema.BoolAttribute{
                MarkdownDescription: "Operational status (ifOperStatus) from the last SNMP walk. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Device], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "is_administratively_up": schema.BoolAttribute{
                MarkdownDescription: "Administrative status (ifAdminStatus) from the last SNMP walk. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Device], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "speed_in_mbps": schema.NumberAttribute{
                MarkdownDescription: "Negotiated interface speed in Mbps. Stored as decimal so 10G+ links don't overflow integers.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Device], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "in_rate_mbps": schema.NumberAttribute{
                MarkdownDescription: "Most recent inbound throughput in Mbps, computed from SNMP counters.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Device], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "out_rate_mbps": schema.NumberAttribute{
                MarkdownDescription: "Most recent outbound throughput in Mbps, computed from SNMP counters.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Device], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "utilization_percent": schema.NumberAttribute{
                MarkdownDescription: "Most recent utilization as a percent of interface speed (max of in/out).. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Device], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "errors_per_second": schema.NumberAttribute{
                MarkdownDescription: "Most recent error rate (in + out errors per second) computed from SNMP counters.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Device], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "last_seen_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
        },
    }
}

func (d *NetworkInterfaceDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *NetworkInterfaceDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data NetworkInterfaceDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "network-interface" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read network_interface_data, got error: %s", err))
        return
    }

    var networkInterfaceDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &networkInterfaceDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse network_interface_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := networkInterfaceDataResponse["data"].(map[string]interface{}); ok {
        networkInterfaceDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := networkInterfaceDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := networkInterfaceDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := networkInterfaceDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := networkInterfaceDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := networkInterfaceDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := networkInterfaceDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := networkInterfaceDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := networkInterfaceDataResponse["network_device_id"].(string); ok {
        data.NetworkDeviceId = types.StringValue(val)
    }
    if val, ok := networkInterfaceDataResponse["interface_index"].(float64); ok {
        data.InterfaceIndex = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := networkInterfaceDataResponse["alias"].(string); ok {
        data.Alias = types.StringValue(val)
    }
    if val, ok := networkInterfaceDataResponse["mac_address"].(string); ok {
        data.MacAddress = types.StringValue(val)
    }
    if val, ok := networkInterfaceDataResponse["interface_type"].(float64); ok {
        data.InterfaceType = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := networkInterfaceDataResponse["is_monitored"].(bool); ok {
        data.IsMonitored = types.BoolValue(val)
    }
    if val, ok := networkInterfaceDataResponse["is_operationally_up"].(bool); ok {
        data.IsOperationallyUp = types.BoolValue(val)
    }
    if val, ok := networkInterfaceDataResponse["is_administratively_up"].(bool); ok {
        data.IsAdministrativelyUp = types.BoolValue(val)
    }
    if val, ok := networkInterfaceDataResponse["speed_in_mbps"].(float64); ok {
        data.SpeedInMbps = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := networkInterfaceDataResponse["in_rate_mbps"].(float64); ok {
        data.InRateMbps = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := networkInterfaceDataResponse["out_rate_mbps"].(float64); ok {
        data.OutRateMbps = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := networkInterfaceDataResponse["utilization_percent"].(float64); ok {
        data.UtilizationPercent = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := networkInterfaceDataResponse["errors_per_second"].(float64); ok {
        data.ErrorsPerSecond = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := networkInterfaceDataResponse["last_seen_at"].(string); ok {
        data.LastSeenAt = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
