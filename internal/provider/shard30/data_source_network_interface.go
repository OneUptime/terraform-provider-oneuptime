package shard30

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
    "context"
    "encoding/json"
    "fmt"
    "net/http"

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
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    CreatedByActorType types.String `tfsdk:"created_by_actor_type"`
    CreatedByActorId types.String `tfsdk:"created_by_actor_id"`
    LastModifiedAt types.String `tfsdk:"last_modified_at"`
    LastModifiedByActorType types.String `tfsdk:"last_modified_by_actor_type"`
    LastModifiedByActorId types.String `tfsdk:"last_modified_by_actor_id"`
    LastModifiedByUserId types.String `tfsdk:"last_modified_by_user_id"`
    ProjectId types.String `tfsdk:"project_id"`
    NetworkDeviceId types.String `tfsdk:"network_device_id"`
    InterfaceIndex types.Number `tfsdk:"interface_index"`
    Name types.String `tfsdk:"name"`
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
        MarkdownDescription: "Interfaces discovered on Network Devices via SNMP walks. Rows are upserted by the server; users can toggle per-interface monitoring. Look up an existing network interface by `id`, or by any of its other arguments (`name`, `alias`, `created_by_actor_id`, ...): each one set must match, and exactly one network interface may match them all.",

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
            "created_by_actor_type": schema.StringAttribute{
                MarkdownDescription: "What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.",
                Optional: true,
                Computed: true,
            },
            "created_by_actor_id": schema.StringAttribute{
                MarkdownDescription: "ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.",
                Optional: true,
                Computed: true,
            },
            "last_modified_at": schema.StringAttribute{
                MarkdownDescription: "When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.",
                Computed: true,
            },
            "last_modified_by_actor_type": schema.StringAttribute{
                MarkdownDescription: "What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.",
                Optional: true,
                Computed: true,
            },
            "last_modified_by_actor_id": schema.StringAttribute{
                MarkdownDescription: "ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.",
                Optional: true,
                Computed: true,
            },
            "last_modified_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: DataSourceRecordIDValidators("oneuptime_user"),
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "network_device_id": schema.StringAttribute{
                MarkdownDescription: "ID of the Network Device this interface was discovered on. The ID of a `oneuptime_network_device`.",
                Optional: true,
                Computed: true,
                Validators: RecordIDValidators("oneuptime_network_device"),
            },
            "interface_index": schema.NumberAttribute{
                MarkdownDescription: "SNMP ifIndex of this interface on the device.",
                Optional: true,
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Interface name (ifName / ifDescr) from SNMP.",
                Optional: true,
                Computed: true,
            },
            "alias": schema.StringAttribute{
                MarkdownDescription: "Interface alias (ifAlias) from SNMP.",
                Optional: true,
                Computed: true,
            },
            "mac_address": schema.StringAttribute{
                MarkdownDescription: "Physical address (ifPhysAddress) from SNMP, colon-separated hex.",
                Optional: true,
                Computed: true,
            },
            "interface_type": schema.NumberAttribute{
                MarkdownDescription: "IANAifType number (ifType) from SNMP — 6 = ethernetCsmacd, 24 = softwareLoopback.",
                Optional: true,
                Computed: true,
            },
            "is_monitored": schema.BoolAttribute{
                MarkdownDescription: "Include this interface in down/utilization/error alerting.",
                Optional: true,
                Computed: true,
            },
            "is_operationally_up": schema.BoolAttribute{
                MarkdownDescription: "Operational status (ifOperStatus) from the last SNMP walk.",
                Optional: true,
                Computed: true,
            },
            "is_administratively_up": schema.BoolAttribute{
                MarkdownDescription: "Administrative status (ifAdminStatus) from the last SNMP walk.",
                Optional: true,
                Computed: true,
            },
            "speed_in_mbps": schema.NumberAttribute{
                MarkdownDescription: "Negotiated interface speed in Mbps. Stored as decimal so 10G+ links don't overflow integers.",
                Optional: true,
                Computed: true,
            },
            "in_rate_mbps": schema.NumberAttribute{
                MarkdownDescription: "Inbound throughput in Mbps, computed from SNMP counters between two polls. Written again when it moves by a tenth or more (and by more than 1% of the link speed), and at least every 15 minutes.",
                Optional: true,
                Computed: true,
            },
            "out_rate_mbps": schema.NumberAttribute{
                MarkdownDescription: "Outbound throughput in Mbps, computed from SNMP counters between two polls. Written again when it moves by a tenth or more (and by more than 1% of the link speed), and at least every 15 minutes.",
                Optional: true,
                Computed: true,
            },
            "utilization_percent": schema.NumberAttribute{
                MarkdownDescription: "Utilization as a percent of interface speed (max of in/out) between two polls. Written again when it moves by a tenth or more (and by more than a point), crosses 80%, and at least every 15 minutes.",
                Optional: true,
                Computed: true,
            },
            "errors_per_second": schema.NumberAttribute{
                MarkdownDescription: "Error rate (in + out errors per second) computed from SNMP counters between two polls. Written again when errors start or stop, when it moves by a tenth or more, and at least every 15 minutes.",
                Optional: true,
                Computed: true,
            },
            "last_seen_at": schema.StringAttribute{
                MarkdownDescription: "When an SNMP walk last saw this interface, to within 15 minutes: a walk that finds an interface unchanged writes it again only every 15 minutes.",
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

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.CreatedByActorType.IsNull() && !data.CreatedByActorType.IsUnknown() {
        filters["createdByActorType"] = data.CreatedByActorType.ValueString()
        filterNames = append(filterNames, "created_by_actor_type = "+fmt.Sprintf("%q", data.CreatedByActorType.ValueString()))
    }
    if !data.CreatedByActorId.IsNull() && !data.CreatedByActorId.IsUnknown() {
        filters["createdByActorId"] = data.CreatedByActorId.ValueString()
        filterNames = append(filterNames, "created_by_actor_id = "+fmt.Sprintf("%q", data.CreatedByActorId.ValueString()))
    }
    if !data.LastModifiedByActorType.IsNull() && !data.LastModifiedByActorType.IsUnknown() {
        filters["lastModifiedByActorType"] = data.LastModifiedByActorType.ValueString()
        filterNames = append(filterNames, "last_modified_by_actor_type = "+fmt.Sprintf("%q", data.LastModifiedByActorType.ValueString()))
    }
    if !data.LastModifiedByActorId.IsNull() && !data.LastModifiedByActorId.IsUnknown() {
        filters["lastModifiedByActorId"] = data.LastModifiedByActorId.ValueString()
        filterNames = append(filterNames, "last_modified_by_actor_id = "+fmt.Sprintf("%q", data.LastModifiedByActorId.ValueString()))
    }
    if !data.LastModifiedByUserId.IsNull() && !data.LastModifiedByUserId.IsUnknown() {
        filters["lastModifiedByUserId"] = data.LastModifiedByUserId.ValueString()
        filterNames = append(filterNames, "last_modified_by_user_id = "+fmt.Sprintf("%q", data.LastModifiedByUserId.ValueString()))
    }
    if !data.NetworkDeviceId.IsNull() && !data.NetworkDeviceId.IsUnknown() {
        filters["networkDeviceId"] = data.NetworkDeviceId.ValueString()
        filterNames = append(filterNames, "network_device_id = "+fmt.Sprintf("%q", data.NetworkDeviceId.ValueString()))
    }
    if !data.InterfaceIndex.IsNull() && !data.InterfaceIndex.IsUnknown() {
        filters["interfaceIndex"] = LookupNumber(data.InterfaceIndex)
        filterNames = append(filterNames, "interface_index = "+data.InterfaceIndex.ValueBigFloat().String())
    }
    if !data.Name.IsNull() && !data.Name.IsUnknown() {
        filters["name"] = data.Name.ValueString()
        filterNames = append(filterNames, "name = "+fmt.Sprintf("%q", data.Name.ValueString()))
    }
    if !data.Alias.IsNull() && !data.Alias.IsUnknown() {
        filters["alias"] = data.Alias.ValueString()
        filterNames = append(filterNames, "alias = "+fmt.Sprintf("%q", data.Alias.ValueString()))
    }
    if !data.MacAddress.IsNull() && !data.MacAddress.IsUnknown() {
        filters["macAddress"] = data.MacAddress.ValueString()
        filterNames = append(filterNames, "mac_address = "+fmt.Sprintf("%q", data.MacAddress.ValueString()))
    }
    if !data.InterfaceType.IsNull() && !data.InterfaceType.IsUnknown() {
        filters["interfaceType"] = LookupNumber(data.InterfaceType)
        filterNames = append(filterNames, "interface_type = "+data.InterfaceType.ValueBigFloat().String())
    }
    if !data.IsMonitored.IsNull() && !data.IsMonitored.IsUnknown() {
        filters["isMonitored"] = data.IsMonitored.ValueBool()
        filterNames = append(filterNames, "is_monitored = "+fmt.Sprintf("%t", data.IsMonitored.ValueBool()))
    }
    if !data.IsOperationallyUp.IsNull() && !data.IsOperationallyUp.IsUnknown() {
        filters["isOperationallyUp"] = data.IsOperationallyUp.ValueBool()
        filterNames = append(filterNames, "is_operationally_up = "+fmt.Sprintf("%t", data.IsOperationallyUp.ValueBool()))
    }
    if !data.IsAdministrativelyUp.IsNull() && !data.IsAdministrativelyUp.IsUnknown() {
        filters["isAdministrativelyUp"] = data.IsAdministrativelyUp.ValueBool()
        filterNames = append(filterNames, "is_administratively_up = "+fmt.Sprintf("%t", data.IsAdministrativelyUp.ValueBool()))
    }
    if !data.SpeedInMbps.IsNull() && !data.SpeedInMbps.IsUnknown() {
        filters["speedInMbps"] = LookupNumber(data.SpeedInMbps)
        filterNames = append(filterNames, "speed_in_mbps = "+data.SpeedInMbps.ValueBigFloat().String())
    }
    if !data.InRateMbps.IsNull() && !data.InRateMbps.IsUnknown() {
        filters["inRateMbps"] = LookupNumber(data.InRateMbps)
        filterNames = append(filterNames, "in_rate_mbps = "+data.InRateMbps.ValueBigFloat().String())
    }
    if !data.OutRateMbps.IsNull() && !data.OutRateMbps.IsUnknown() {
        filters["outRateMbps"] = LookupNumber(data.OutRateMbps)
        filterNames = append(filterNames, "out_rate_mbps = "+data.OutRateMbps.ValueBigFloat().String())
    }
    if !data.UtilizationPercent.IsNull() && !data.UtilizationPercent.IsUnknown() {
        filters["utilizationPercent"] = LookupNumber(data.UtilizationPercent)
        filterNames = append(filterNames, "utilization_percent = "+data.UtilizationPercent.ValueBigFloat().String())
    }
    if !data.ErrorsPerSecond.IsNull() && !data.ErrorsPerSecond.IsUnknown() {
        filters["errorsPerSecond"] = LookupNumber(data.ErrorsPerSecond)
        filterNames = append(filterNames, "errors_per_second = "+data.ErrorsPerSecond.ValueBigFloat().String())
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the network interface up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the network interface up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "createdByActorType": true,
        "createdByActorId": true,
        "lastModifiedAt": true,
        "lastModifiedByActorType": true,
        "lastModifiedByActorId": true,
        "lastModifiedByUserId": true,
        "projectId": true,
        "networkDeviceId": true,
        "interfaceIndex": true,
        "name": true,
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
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No network interface found with id %q.", data.Id.ValueString()))
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
    }
    if !hasId {
        listBody := map[string]interface{}{
            "query":  filters,
            "select": selectParam,
            // limit 2 is enough to detect ambiguity without paging.
            "limit": 2,
        }
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/network-interface/get-list", listBody)
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
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No network interface matches %s.", DescribeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one network interface matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", DescribeLookup(filterNames)))
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
    } else if arr, ok := item["_id"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
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
    } else if arr, ok := item["createdAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
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
    } else if arr, ok := item["updatedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.UpdatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.UpdatedAt = types.StringNull()
        }
    } else if val, ok := item["updatedAt"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    } else {
        data.UpdatedAt = types.StringNull()
    }
    if obj, ok := item["createdByActorType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByActorType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CreatedByActorType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CreatedByActorType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CreatedByActorType = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByActorType = types.StringNull()
        }
    } else if arr, ok := item["createdByActorType"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.CreatedByActorType = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByActorType = types.StringNull()
        }
    } else if val, ok := item["createdByActorType"].(string); ok {
        data.CreatedByActorType = types.StringValue(val)
    } else {
        data.CreatedByActorType = types.StringNull()
    }
    if obj, ok := item["createdByActorId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByActorId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CreatedByActorId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CreatedByActorId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CreatedByActorId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByActorId = types.StringNull()
        }
    } else if arr, ok := item["createdByActorId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.CreatedByActorId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByActorId = types.StringNull()
        }
    } else if val, ok := item["createdByActorId"].(string); ok {
        data.CreatedByActorId = types.StringValue(val)
    } else {
        data.CreatedByActorId = types.StringNull()
    }
    if obj, ok := item["lastModifiedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastModifiedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastModifiedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastModifiedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastModifiedAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedAt = types.StringNull()
        }
    } else if arr, ok := item["lastModifiedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LastModifiedAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedAt = types.StringNull()
        }
    } else if val, ok := item["lastModifiedAt"].(string); ok {
        data.LastModifiedAt = types.StringValue(val)
    } else {
        data.LastModifiedAt = types.StringNull()
    }
    if obj, ok := item["lastModifiedByActorType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastModifiedByActorType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastModifiedByActorType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastModifiedByActorType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastModifiedByActorType = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByActorType = types.StringNull()
        }
    } else if arr, ok := item["lastModifiedByActorType"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LastModifiedByActorType = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByActorType = types.StringNull()
        }
    } else if val, ok := item["lastModifiedByActorType"].(string); ok {
        data.LastModifiedByActorType = types.StringValue(val)
    } else {
        data.LastModifiedByActorType = types.StringNull()
    }
    if obj, ok := item["lastModifiedByActorId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastModifiedByActorId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastModifiedByActorId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastModifiedByActorId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastModifiedByActorId = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByActorId = types.StringNull()
        }
    } else if arr, ok := item["lastModifiedByActorId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LastModifiedByActorId = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByActorId = types.StringNull()
        }
    } else if val, ok := item["lastModifiedByActorId"].(string); ok {
        data.LastModifiedByActorId = types.StringValue(val)
    } else {
        data.LastModifiedByActorId = types.StringNull()
    }
    if obj, ok := item["lastModifiedByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastModifiedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastModifiedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastModifiedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastModifiedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByUserId = types.StringNull()
        }
    } else if arr, ok := item["lastModifiedByUserId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LastModifiedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByUserId = types.StringNull()
        }
    } else if val, ok := item["lastModifiedByUserId"].(string); ok {
        data.LastModifiedByUserId = types.StringValue(val)
    } else {
        data.LastModifiedByUserId = types.StringNull()
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
    } else if arr, ok := item["projectId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
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
    } else if arr, ok := item["networkDeviceId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.NetworkDeviceId = types.StringValue(string(jsonBytes))
        } else {
            data.NetworkDeviceId = types.StringNull()
        }
    } else if val, ok := item["networkDeviceId"].(string); ok {
        data.NetworkDeviceId = types.StringValue(val)
    } else {
        data.NetworkDeviceId = types.StringNull()
    }
    data.InterfaceIndex = NumberFromAPI(data.InterfaceIndex, item["interfaceIndex"])
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
    } else if arr, ok := item["name"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Name = types.StringValue(string(jsonBytes))
        } else {
            data.Name = types.StringNull()
        }
    } else if val, ok := item["name"].(string); ok {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
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
    } else if arr, ok := item["alias"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
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
    } else if arr, ok := item["macAddress"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.MacAddress = types.StringValue(string(jsonBytes))
        } else {
            data.MacAddress = types.StringNull()
        }
    } else if val, ok := item["macAddress"].(string); ok {
        data.MacAddress = types.StringValue(val)
    } else {
        data.MacAddress = types.StringNull()
    }
    data.InterfaceType = NumberFromAPI(data.InterfaceType, item["interfaceType"])
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
    data.SpeedInMbps = NumberFromAPI(data.SpeedInMbps, item["speedInMbps"])
    data.InRateMbps = NumberFromAPI(data.InRateMbps, item["inRateMbps"])
    data.OutRateMbps = NumberFromAPI(data.OutRateMbps, item["outRateMbps"])
    data.UtilizationPercent = NumberFromAPI(data.UtilizationPercent, item["utilizationPercent"])
    data.ErrorsPerSecond = NumberFromAPI(data.ErrorsPerSecond, item["errorsPerSecond"])
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
    } else if arr, ok := item["lastSeenAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
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
