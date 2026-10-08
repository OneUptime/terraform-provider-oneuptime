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
var _ datasource.DataSource = &NetworkEndpointDataSource{}

func NewNetworkEndpointDataSource() datasource.DataSource {
    return &NetworkEndpointDataSource{}
}

// NetworkEndpointDataSource defines the data source implementation.
type NetworkEndpointDataSource struct {
    client *Client
}

// NetworkEndpointDataSourceModel describes the data source data model.
type NetworkEndpointDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
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
}

func (d *NetworkEndpointDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_network_endpoint"
}

func (d *NetworkEndpointDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "LAN endpoints (POS terminals, kiosks, cameras, printers) discovered via ARP and FDB walks of Network Devices. Rows are upserted by the server; users can classify them. Look up an existing network endpoint by `id`, or by any of its other arguments (`attached_interface_index`, `attached_network_device_id`, `attached_port_name`, ...): each one set must match, and exactly one network endpoint may match them all.",

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
            "mac_address": schema.StringAttribute{
                MarkdownDescription: "MAC address of this endpoint, colon-separated hex. One row per MAC per project.",
                Optional: true,
                Computed: true,
            },
            "ip_address": schema.StringAttribute{
                MarkdownDescription: "Last IP address seen for this endpoint in ARP tables. Managed by the server.",
                Optional: true,
                Computed: true,
            },
            "vendor": schema.StringAttribute{
                MarkdownDescription: "Hardware vendor derived from the MAC OUI prefix. Managed by the server.",
                Optional: true,
                Computed: true,
            },
            "classification": schema.StringAttribute{
                MarkdownDescription: "User-editable classification of this endpoint (POS, Kiosk, Camera, Printer, ...).",
                Optional: true,
                Computed: true,
            },
            "attached_network_device_id": schema.StringAttribute{
                MarkdownDescription: "ID of the Network Device this endpoint was last seen attached to. The ID of a `oneuptime_network_device`.",
                Optional: true,
                Computed: true,
            },
            "attached_interface_index": schema.NumberAttribute{
                MarkdownDescription: "SNMP ifIndex of the switch port this endpoint was last seen on. Managed by the server.",
                Optional: true,
                Computed: true,
            },
            "attached_port_name": schema.StringAttribute{
                MarkdownDescription: "Name of the switch port this endpoint was last seen on. Managed by the server.",
                Optional: true,
                Computed: true,
            },
            "vlan_id": schema.NumberAttribute{
                MarkdownDescription: "VLAN this endpoint was last seen on, from the FDB walk. Managed by the server.",
                Optional: true,
                Computed: true,
            },
            "site_id": schema.StringAttribute{
                MarkdownDescription: "ID of the Network Site this endpoint belongs to. The ID of a `oneuptime_network_site`.",
                Optional: true,
                Computed: true,
            },
            "first_seen_at": schema.StringAttribute{
                MarkdownDescription: "When this endpoint was first discovered on the network.",
                Computed: true,
            },
            "last_seen_at": schema.StringAttribute{
                MarkdownDescription: "When this endpoint was last seen in an ARP or FDB walk.",
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

func (d *NetworkEndpointDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *NetworkEndpointDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data NetworkEndpointDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.MacAddress.IsNull() && !data.MacAddress.IsUnknown() {
        filters["macAddress"] = data.MacAddress.ValueString()
        filterNames = append(filterNames, "mac_address = "+fmt.Sprintf("%q", data.MacAddress.ValueString()))
    }
    if !data.IpAddress.IsNull() && !data.IpAddress.IsUnknown() {
        filters["ipAddress"] = data.IpAddress.ValueString()
        filterNames = append(filterNames, "ip_address = "+fmt.Sprintf("%q", data.IpAddress.ValueString()))
    }
    if !data.Vendor.IsNull() && !data.Vendor.IsUnknown() {
        filters["vendor"] = data.Vendor.ValueString()
        filterNames = append(filterNames, "vendor = "+fmt.Sprintf("%q", data.Vendor.ValueString()))
    }
    if !data.Classification.IsNull() && !data.Classification.IsUnknown() {
        filters["classification"] = data.Classification.ValueString()
        filterNames = append(filterNames, "classification = "+fmt.Sprintf("%q", data.Classification.ValueString()))
    }
    if !data.AttachedNetworkDeviceId.IsNull() && !data.AttachedNetworkDeviceId.IsUnknown() {
        filters["attachedNetworkDeviceId"] = data.AttachedNetworkDeviceId.ValueString()
        filterNames = append(filterNames, "attached_network_device_id = "+fmt.Sprintf("%q", data.AttachedNetworkDeviceId.ValueString()))
    }
    if !data.AttachedInterfaceIndex.IsNull() && !data.AttachedInterfaceIndex.IsUnknown() {
        filters["attachedInterfaceIndex"] = lookupNumber(data.AttachedInterfaceIndex)
        filterNames = append(filterNames, "attached_interface_index = "+data.AttachedInterfaceIndex.ValueBigFloat().String())
    }
    if !data.AttachedPortName.IsNull() && !data.AttachedPortName.IsUnknown() {
        filters["attachedPortName"] = data.AttachedPortName.ValueString()
        filterNames = append(filterNames, "attached_port_name = "+fmt.Sprintf("%q", data.AttachedPortName.ValueString()))
    }
    if !data.VlanId.IsNull() && !data.VlanId.IsUnknown() {
        filters["vlanId"] = lookupNumber(data.VlanId)
        filterNames = append(filterNames, "vlan_id = "+data.VlanId.ValueBigFloat().String())
    }
    if !data.SiteId.IsNull() && !data.SiteId.IsUnknown() {
        filters["siteId"] = data.SiteId.ValueString()
        filterNames = append(filterNames, "site_id = "+fmt.Sprintf("%q", data.SiteId.ValueString()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the network endpoint up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the network endpoint up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "macAddress": true,
        "ipAddress": true,
        "vendor": true,
        "classification": true,
        "attachedNetworkDeviceId": true,
        "attachedInterfaceIndex": true,
        "attachedPortName": true,
        "vlanId": true,
        "siteId": true,
        "firstSeenAt": true,
        "lastSeenAt": true,
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/network-endpoint/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read network_endpoint, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No network endpoint found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read network_endpoint: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/network-endpoint/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list network_endpoint, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list network_endpoint: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No network endpoint matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one network endpoint matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for network_endpoint.")
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
    if obj, ok := item["macAddress"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MacAddress = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.MacAddress = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.MacAddress = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.MacAddress = types.StringValue(string(jsonBytes))
        } else {
            data.MacAddress = types.StringNull()
        }
    } else if val, ok := item["macAddress"].(string); ok {
        data.MacAddress = types.StringValue(val)
    } else {
        data.MacAddress = types.StringNull()
    }
    if obj, ok := item["ipAddress"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IpAddress = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.IpAddress = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.IpAddress = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.IpAddress = types.StringValue(string(jsonBytes))
        } else {
            data.IpAddress = types.StringNull()
        }
    } else if val, ok := item["ipAddress"].(string); ok {
        data.IpAddress = types.StringValue(val)
    } else {
        data.IpAddress = types.StringNull()
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
    if obj, ok := item["classification"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Classification = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Classification = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Classification = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Classification = types.StringValue(string(jsonBytes))
        } else {
            data.Classification = types.StringNull()
        }
    } else if val, ok := item["classification"].(string); ok {
        data.Classification = types.StringValue(val)
    } else {
        data.Classification = types.StringNull()
    }
    if obj, ok := item["attachedNetworkDeviceId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AttachedNetworkDeviceId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AttachedNetworkDeviceId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AttachedNetworkDeviceId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AttachedNetworkDeviceId = types.StringValue(string(jsonBytes))
        } else {
            data.AttachedNetworkDeviceId = types.StringNull()
        }
    } else if val, ok := item["attachedNetworkDeviceId"].(string); ok {
        data.AttachedNetworkDeviceId = types.StringValue(val)
    } else {
        data.AttachedNetworkDeviceId = types.StringNull()
    }
    if val, ok := item["attachedInterfaceIndex"].(float64); ok {
        data.AttachedInterfaceIndex = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["attachedInterfaceIndex"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.AttachedInterfaceIndex = types.NumberValue(big.NewFloat(val))
        } else {
            data.AttachedInterfaceIndex = types.NumberNull()
        }
    } else {
        data.AttachedInterfaceIndex = types.NumberNull()
    }
    if obj, ok := item["attachedPortName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AttachedPortName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AttachedPortName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AttachedPortName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AttachedPortName = types.StringValue(string(jsonBytes))
        } else {
            data.AttachedPortName = types.StringNull()
        }
    } else if val, ok := item["attachedPortName"].(string); ok {
        data.AttachedPortName = types.StringValue(val)
    } else {
        data.AttachedPortName = types.StringNull()
    }
    if val, ok := item["vlanId"].(float64); ok {
        data.VlanId = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["vlanId"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.VlanId = types.NumberValue(big.NewFloat(val))
        } else {
            data.VlanId = types.NumberNull()
        }
    } else {
        data.VlanId = types.NumberNull()
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
    if obj, ok := item["firstSeenAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FirstSeenAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.FirstSeenAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.FirstSeenAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.FirstSeenAt = types.StringValue(string(jsonBytes))
        } else {
            data.FirstSeenAt = types.StringNull()
        }
    } else if val, ok := item["firstSeenAt"].(string); ok {
        data.FirstSeenAt = types.StringValue(val)
    } else {
        data.FirstSeenAt = types.StringNull()
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
