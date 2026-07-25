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
var _ datasource.DataSource = &NetworkEndpointDataDataSource{}

func NewNetworkEndpointDataDataSource() datasource.DataSource {
    return &NetworkEndpointDataDataSource{}
}

// NetworkEndpointDataDataSource defines the data source implementation.
type NetworkEndpointDataDataSource struct {
    client *Client
}

// NetworkEndpointDataDataSourceModel describes the data source data model.
type NetworkEndpointDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    MacAddress types.String `tfsdk:"mac_address"`
    IpAddress types.String `tfsdk:"ip_address"`
    Vendor types.String `tfsdk:"vendor"`
    Classification types.String `tfsdk:"classification"`
    AttachedNetworkDeviceId types.String `tfsdk:"attached_network_device_id"`
    AttachedInterfaceIndex types.Number `tfsdk:"attached_interface_index"`
    AttachedPortName types.String `tfsdk:"attached_port_name"`
    VlanId types.Number `tfsdk:"vlan_id"`
    SiteId types.String `tfsdk:"site_id"`
    FirstSeenAt types.String `tfsdk:"first_seen_at"`
    LastSeenAt types.String `tfsdk:"last_seen_at"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
}

func (d *NetworkEndpointDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_network_endpoint_data"
}

func (d *NetworkEndpointDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "network_endpoint_data data source",

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
            "mac_address": schema.StringAttribute{
                MarkdownDescription: "MAC address of this endpoint, colon-separated hex. One row per MAC per project.. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Network Endpoint], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Endpoint], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "ip_address": schema.StringAttribute{
                MarkdownDescription: "Last IP address seen for this endpoint in ARP tables. Managed by the server.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Endpoint], Update: [Project Owner, Project Admin, Edit Network Endpoint]",
                Computed: true,
            },
            "vendor": schema.StringAttribute{
                MarkdownDescription: "Hardware vendor derived from the MAC OUI prefix. Managed by the server.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Endpoint], Update: [Project Owner, Project Admin, Edit Network Endpoint]",
                Computed: true,
            },
            "classification": schema.StringAttribute{
                MarkdownDescription: "User-editable classification of this endpoint (POS, Kiosk, Camera, Printer, ...). Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Network Endpoint], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Endpoint], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Network Endpoint]",
                Computed: true,
            },
            "attached_network_device_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "attached_interface_index": schema.NumberAttribute{
                MarkdownDescription: "SNMP ifIndex of the switch port this endpoint was last seen on. Managed by the server.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Endpoint], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "attached_port_name": schema.StringAttribute{
                MarkdownDescription: "Name of the switch port this endpoint was last seen on. Managed by the server.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Endpoint], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "vlan_id": schema.NumberAttribute{
                MarkdownDescription: "VLAN this endpoint was last seen on, from the FDB walk. Managed by the server.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Network Endpoint], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "site_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "first_seen_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "last_seen_at": schema.StringAttribute{
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

func (d *NetworkEndpointDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *NetworkEndpointDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data NetworkEndpointDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "network-endpoint" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read network_endpoint_data, got error: %s", err))
        return
    }

    var networkEndpointDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &networkEndpointDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse network_endpoint_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := networkEndpointDataResponse["data"].(map[string]interface{}); ok {
        networkEndpointDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := networkEndpointDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := networkEndpointDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := networkEndpointDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := networkEndpointDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := networkEndpointDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := networkEndpointDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := networkEndpointDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := networkEndpointDataResponse["mac_address"].(string); ok {
        data.MacAddress = types.StringValue(val)
    }
    if val, ok := networkEndpointDataResponse["ip_address"].(string); ok {
        data.IpAddress = types.StringValue(val)
    }
    if val, ok := networkEndpointDataResponse["vendor"].(string); ok {
        data.Vendor = types.StringValue(val)
    }
    if val, ok := networkEndpointDataResponse["classification"].(string); ok {
        data.Classification = types.StringValue(val)
    }
    if val, ok := networkEndpointDataResponse["attached_network_device_id"].(string); ok {
        data.AttachedNetworkDeviceId = types.StringValue(val)
    }
    if val, ok := networkEndpointDataResponse["attached_interface_index"].(float64); ok {
        data.AttachedInterfaceIndex = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := networkEndpointDataResponse["attached_port_name"].(string); ok {
        data.AttachedPortName = types.StringValue(val)
    }
    if val, ok := networkEndpointDataResponse["vlan_id"].(float64); ok {
        data.VlanId = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := networkEndpointDataResponse["site_id"].(string); ok {
        data.SiteId = types.StringValue(val)
    }
    if val, ok := networkEndpointDataResponse["first_seen_at"].(string); ok {
        data.FirstSeenAt = types.StringValue(val)
    }
    if val, ok := networkEndpointDataResponse["last_seen_at"].(string); ok {
        data.LastSeenAt = types.StringValue(val)
    }
    if val, ok := networkEndpointDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := networkEndpointDataResponse["deleted_by_user_id"].(string); ok {
        data.DeletedByUserId = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
