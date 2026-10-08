package provider

import (
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
var _ datasource.DataSource = &NetworkDeviceLinkDataSource{}

func NewNetworkDeviceLinkDataSource() datasource.DataSource {
    return &NetworkDeviceLinkDataSource{}
}

// NetworkDeviceLinkDataSource defines the data source implementation.
type NetworkDeviceLinkDataSource struct {
    client *Client
}

// NetworkDeviceLinkDataSourceModel describes the data source data model.
type NetworkDeviceLinkDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    Name types.String `tfsdk:"name"`
    FromDeviceId types.String `tfsdk:"from_device_id"`
    ToDeviceId types.String `tfsdk:"to_device_id"`
    ParentDeviceId types.String `tfsdk:"parent_device_id"`
    FromPortName types.String `tfsdk:"from_port_name"`
    ToPortName types.String `tfsdk:"to_port_name"`
    MonitorId types.String `tfsdk:"monitor_id"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *NetworkDeviceLinkDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_network_device_link"
}

func (d *NetworkDeviceLinkDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Operator-declared links between two Network Devices, for cables LLDP and CDP cannot see: a device with discovery disabled, a device that does not speak either protocol, or one monitored by ping alone. Drawn on the topology map alongside discovered links, and merged with a discovered link between the same pair rather than duplicating it. Look up an existing network device link by `id`, or by any of its other arguments (`name`, `created_by_user_id`, `from_device_id`, ...): each one set must match, and exactly one network device link may match them all.",

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
            "name": schema.StringAttribute{
                MarkdownDescription: "Friendly name for this link.",
                Optional: true,
                Computed: true,
            },
            "from_device_id": schema.StringAttribute{
                MarkdownDescription: "ID of the Network Device this link starts from. The ID of a `oneuptime_network_device`.",
                Optional: true,
                Computed: true,
            },
            "to_device_id": schema.StringAttribute{
                MarkdownDescription: "ID of the Network Device this link ends at. The ID of a `oneuptime_network_device`.",
                Optional: true,
                Computed: true,
            },
            "parent_device_id": schema.StringAttribute{
                MarkdownDescription: "ID of whichever end of this link is the parent. Must be the From Device or the To Device. Empty means the two are peers and the map infers the hierarchy. The ID of a `oneuptime_network_device`.",
                Optional: true,
                Computed: true,
            },
            "from_port_name": schema.StringAttribute{
                MarkdownDescription: "Port on the starting device, as free text. Nothing resolves it to an interface row — a hand-drawn link usually exists precisely because the port is not discoverable.",
                Optional: true,
                Computed: true,
            },
            "to_port_name": schema.StringAttribute{
                MarkdownDescription: "Port on the ending device, as free text.",
                Optional: true,
                Computed: true,
            },
            "monitor_id": schema.StringAttribute{
                MarkdownDescription: "ID of the Monitor whose status colors this link on the topology map. The ID of a `oneuptime_monitor`.",
                Optional: true,
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

func (d *NetworkDeviceLinkDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *NetworkDeviceLinkDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data NetworkDeviceLinkDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.Name.IsNull() && !data.Name.IsUnknown() {
        filters["name"] = data.Name.ValueString()
        filterNames = append(filterNames, "name = "+fmt.Sprintf("%q", data.Name.ValueString()))
    }
    if !data.FromDeviceId.IsNull() && !data.FromDeviceId.IsUnknown() {
        filters["fromDeviceId"] = data.FromDeviceId.ValueString()
        filterNames = append(filterNames, "from_device_id = "+fmt.Sprintf("%q", data.FromDeviceId.ValueString()))
    }
    if !data.ToDeviceId.IsNull() && !data.ToDeviceId.IsUnknown() {
        filters["toDeviceId"] = data.ToDeviceId.ValueString()
        filterNames = append(filterNames, "to_device_id = "+fmt.Sprintf("%q", data.ToDeviceId.ValueString()))
    }
    if !data.ParentDeviceId.IsNull() && !data.ParentDeviceId.IsUnknown() {
        filters["parentDeviceId"] = data.ParentDeviceId.ValueString()
        filterNames = append(filterNames, "parent_device_id = "+fmt.Sprintf("%q", data.ParentDeviceId.ValueString()))
    }
    if !data.FromPortName.IsNull() && !data.FromPortName.IsUnknown() {
        filters["fromPortName"] = data.FromPortName.ValueString()
        filterNames = append(filterNames, "from_port_name = "+fmt.Sprintf("%q", data.FromPortName.ValueString()))
    }
    if !data.ToPortName.IsNull() && !data.ToPortName.IsUnknown() {
        filters["toPortName"] = data.ToPortName.ValueString()
        filterNames = append(filterNames, "to_port_name = "+fmt.Sprintf("%q", data.ToPortName.ValueString()))
    }
    if !data.MonitorId.IsNull() && !data.MonitorId.IsUnknown() {
        filters["monitorId"] = data.MonitorId.ValueString()
        filterNames = append(filterNames, "monitor_id = "+fmt.Sprintf("%q", data.MonitorId.ValueString()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the network device link up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the network device link up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "name": true,
        "fromDeviceId": true,
        "toDeviceId": true,
        "parentDeviceId": true,
        "fromPortName": true,
        "toPortName": true,
        "monitorId": true,
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/network-device-link/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read network_device_link, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No network device link found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read network_device_link: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/network-device-link/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list network_device_link, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list network_device_link: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No network device link matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one network device link matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for network_device_link.")
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
    if obj, ok := item["fromDeviceId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FromDeviceId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.FromDeviceId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.FromDeviceId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.FromDeviceId = types.StringValue(string(jsonBytes))
        } else {
            data.FromDeviceId = types.StringNull()
        }
    } else if val, ok := item["fromDeviceId"].(string); ok {
        data.FromDeviceId = types.StringValue(val)
    } else {
        data.FromDeviceId = types.StringNull()
    }
    if obj, ok := item["toDeviceId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ToDeviceId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ToDeviceId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ToDeviceId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ToDeviceId = types.StringValue(string(jsonBytes))
        } else {
            data.ToDeviceId = types.StringNull()
        }
    } else if val, ok := item["toDeviceId"].(string); ok {
        data.ToDeviceId = types.StringValue(val)
    } else {
        data.ToDeviceId = types.StringNull()
    }
    if obj, ok := item["parentDeviceId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ParentDeviceId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ParentDeviceId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ParentDeviceId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ParentDeviceId = types.StringValue(string(jsonBytes))
        } else {
            data.ParentDeviceId = types.StringNull()
        }
    } else if val, ok := item["parentDeviceId"].(string); ok {
        data.ParentDeviceId = types.StringValue(val)
    } else {
        data.ParentDeviceId = types.StringNull()
    }
    if obj, ok := item["fromPortName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FromPortName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.FromPortName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.FromPortName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.FromPortName = types.StringValue(string(jsonBytes))
        } else {
            data.FromPortName = types.StringNull()
        }
    } else if val, ok := item["fromPortName"].(string); ok {
        data.FromPortName = types.StringValue(val)
    } else {
        data.FromPortName = types.StringNull()
    }
    if obj, ok := item["toPortName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ToPortName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ToPortName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ToPortName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ToPortName = types.StringValue(string(jsonBytes))
        } else {
            data.ToPortName = types.StringNull()
        }
    } else if val, ok := item["toPortName"].(string); ok {
        data.ToPortName = types.StringValue(val)
    } else {
        data.ToPortName = types.StringNull()
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
