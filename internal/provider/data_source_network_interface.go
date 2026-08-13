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
var _ datasource.DataSource = &NetworkInterfaceDataSource{}

func NewNetworkInterfaceDataSource() datasource.DataSource {
    return &NetworkInterfaceDataSource{}
}

// NetworkInterfaceDataSource defines the data source implementation.
type NetworkInterfaceDataSource struct {
    client *Client
}

// NetworkInterfaceDataSourceModel describes the data source data model.
type NetworkInterfaceDataSourceModel struct {
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

func (d *NetworkInterfaceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_network_interface"
}

func (d *NetworkInterfaceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Interfaces discovered on Network Devices via SNMP walks. Rows are upserted by the server; users can toggle per-interface monitoring. Look up an existing network_interface by `id` or by `name`.",

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
            "network_device_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "interface_index": schema.NumberAttribute{
                MarkdownDescription: "SNMP ifIndex of this interface on the device.",
                Computed: true,
            },
            "alias": schema.StringAttribute{
                MarkdownDescription: "Interface alias (ifAlias) from SNMP.",
                Computed: true,
            },
            "mac_address": schema.StringAttribute{
                MarkdownDescription: "Physical address (ifPhysAddress) from SNMP, colon-separated hex.",
                Computed: true,
            },
            "interface_type": schema.NumberAttribute{
                MarkdownDescription: "IANAifType number (ifType) from SNMP — 6 = ethernetCsmacd, 24 = softwareLoopback.",
                Computed: true,
            },
            "is_monitored": schema.BoolAttribute{
                MarkdownDescription: "Include this interface in down/utilization/error alerting..",
                Computed: true,
            },
            "is_operationally_up": schema.BoolAttribute{
                MarkdownDescription: "Operational status (ifOperStatus) from the last SNMP walk.",
                Computed: true,
            },
            "is_administratively_up": schema.BoolAttribute{
                MarkdownDescription: "Administrative status (ifAdminStatus) from the last SNMP walk.",
                Computed: true,
            },
            "speed_in_mbps": schema.NumberAttribute{
                MarkdownDescription: "Negotiated interface speed in Mbps. Stored as decimal so 10G+ links don't overflow integers..",
                Computed: true,
            },
            "in_rate_mbps": schema.NumberAttribute{
                MarkdownDescription: "Most recent inbound throughput in Mbps, computed from SNMP counters..",
                Computed: true,
            },
            "out_rate_mbps": schema.NumberAttribute{
                MarkdownDescription: "Most recent outbound throughput in Mbps, computed from SNMP counters..",
                Computed: true,
            },
            "utilization_percent": schema.NumberAttribute{
                MarkdownDescription: "Most recent utilization as a percent of interface speed (max of in/out)..",
                Computed: true,
            },
            "errors_per_second": schema.NumberAttribute{
                MarkdownDescription: "Most recent error rate (in + out errors per second) computed from SNMP counters..",
                Computed: true,
            },
            "last_seen_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
        },
    }
}

func (d *NetworkInterfaceDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *NetworkInterfaceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data NetworkInterfaceDataSourceModel

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
            "Exactly one of `id` or `name` must be set to look up a network_interface.",
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
        "networkDeviceId": true,
        "interfaceIndex": true,
        "alias": true,
        "macAddress": true,
        "interfaceType": true,
        "isMonitored": true,
        "isOperationallyUp": true,
        "isAdministrativelyUp": true,
        "speedInMbps": true,
        "inRateMbps": true,
        "outRateMbps": true,
        "utilizationPercent": true,
        "errorsPerSecond": true,
        "lastSeenAt": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/network-interface/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read network_interface, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No network_interface found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read network_interface: %s", err))
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
        httpResp, err := d.client.Post(ctx, "/network-interface/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list network_interface, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list network_interface: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No network_interface found with name %q.", data.Name.ValueString()))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one network_interface matches name %q. Use the id attribute to disambiguate.", data.Name.ValueString()))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for network_interface.")
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
    if val, ok := item["interfaceIndex"].(float64); ok {
        data.InterfaceIndex = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["interfaceIndex"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.InterfaceIndex = types.NumberValue(big.NewFloat(val))
        } else {
            data.InterfaceIndex = types.NumberNull()
        }
    } else {
        data.InterfaceIndex = types.NumberNull()
    }
    if obj, ok := item["alias"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Alias = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Alias = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Alias = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Alias = types.StringValue(string(jsonBytes))
        } else {
            data.Alias = types.StringNull()
        }
    } else if val, ok := item["alias"].(string); ok {
        data.Alias = types.StringValue(val)
    } else {
        data.Alias = types.StringNull()
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
    if val, ok := item["interfaceType"].(float64); ok {
        data.InterfaceType = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["interfaceType"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.InterfaceType = types.NumberValue(big.NewFloat(val))
        } else {
            data.InterfaceType = types.NumberNull()
        }
    } else {
        data.InterfaceType = types.NumberNull()
    }
    if val, ok := item["isMonitored"].(bool); ok {
        data.IsMonitored = types.BoolValue(val)
    } else {
        data.IsMonitored = types.BoolNull()
    }
    if val, ok := item["isOperationallyUp"].(bool); ok {
        data.IsOperationallyUp = types.BoolValue(val)
    } else {
        data.IsOperationallyUp = types.BoolNull()
    }
    if val, ok := item["isAdministrativelyUp"].(bool); ok {
        data.IsAdministrativelyUp = types.BoolValue(val)
    } else {
        data.IsAdministrativelyUp = types.BoolNull()
    }
    if val, ok := item["speedInMbps"].(float64); ok {
        data.SpeedInMbps = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["speedInMbps"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.SpeedInMbps = types.NumberValue(big.NewFloat(val))
        } else {
            data.SpeedInMbps = types.NumberNull()
        }
    } else {
        data.SpeedInMbps = types.NumberNull()
    }
    if val, ok := item["inRateMbps"].(float64); ok {
        data.InRateMbps = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["inRateMbps"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.InRateMbps = types.NumberValue(big.NewFloat(val))
        } else {
            data.InRateMbps = types.NumberNull()
        }
    } else {
        data.InRateMbps = types.NumberNull()
    }
    if val, ok := item["outRateMbps"].(float64); ok {
        data.OutRateMbps = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["outRateMbps"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.OutRateMbps = types.NumberValue(big.NewFloat(val))
        } else {
            data.OutRateMbps = types.NumberNull()
        }
    } else {
        data.OutRateMbps = types.NumberNull()
    }
    if val, ok := item["utilizationPercent"].(float64); ok {
        data.UtilizationPercent = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["utilizationPercent"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.UtilizationPercent = types.NumberValue(big.NewFloat(val))
        } else {
            data.UtilizationPercent = types.NumberNull()
        }
    } else {
        data.UtilizationPercent = types.NumberNull()
    }
    if val, ok := item["errorsPerSecond"].(float64); ok {
        data.ErrorsPerSecond = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["errorsPerSecond"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.ErrorsPerSecond = types.NumberValue(big.NewFloat(val))
        } else {
            data.ErrorsPerSecond = types.NumberNull()
        }
    } else {
        data.ErrorsPerSecond = types.NumberNull()
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

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
