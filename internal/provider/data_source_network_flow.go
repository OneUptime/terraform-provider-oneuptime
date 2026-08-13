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
var _ datasource.DataSource = &NetworkFlowDataSource{}

func NewNetworkFlowDataSource() datasource.DataSource {
    return &NetworkFlowDataSource{}
}

// NetworkFlowDataSource defines the data source implementation.
type NetworkFlowDataSource struct {
    client *Client
}

// NetworkFlowDataSourceModel describes the data source data model.
type NetworkFlowDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    ProjectId types.String `tfsdk:"project_id"`
    NetworkDeviceId types.String `tfsdk:"network_device_id"`
    ExporterIp types.String `tfsdk:"exporter_ip"`
    SrcIp types.String `tfsdk:"src_ip"`
    DstIp types.String `tfsdk:"dst_ip"`
    SrcPort types.Number `tfsdk:"src_port"`
    DstPort types.Number `tfsdk:"dst_port"`
    Protocol types.Number `tfsdk:"protocol"`
    InputInterfaceIndex types.Number `tfsdk:"input_interface_index"`
    OutputInterfaceIndex types.Number `tfsdk:"output_interface_index"`
    Octets types.String `tfsdk:"octets"`
    Packets types.String `tfsdk:"packets"`
    FlowStartAt types.String `tfsdk:"flow_start_at"`
    FlowEndAt types.String `tfsdk:"flow_end_at"`
    IngestedAt types.String `tfsdk:"ingested_at"`
}

func (d *NetworkFlowDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_network_flow"
}

func (d *NetworkFlowDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "API endpoints for Network Flow Look up an existing network_flow by `id` or by `name`.",

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
            "project_id": schema.StringAttribute{
                MarkdownDescription: "Project ID",
                Computed: true,
            },
            "network_device_id": schema.StringAttribute{
                MarkdownDescription: "Network Device ID",
                Computed: true,
            },
            "exporter_ip": schema.StringAttribute{
                MarkdownDescription: "Exporter IP",
                Computed: true,
            },
            "src_ip": schema.StringAttribute{
                MarkdownDescription: "Source IP",
                Computed: true,
            },
            "dst_ip": schema.StringAttribute{
                MarkdownDescription: "Destination IP",
                Computed: true,
            },
            "src_port": schema.NumberAttribute{
                MarkdownDescription: "Source Port",
                Computed: true,
            },
            "dst_port": schema.NumberAttribute{
                MarkdownDescription: "Destination Port",
                Computed: true,
            },
            "protocol": schema.NumberAttribute{
                MarkdownDescription: "Protocol",
                Computed: true,
            },
            "input_interface_index": schema.NumberAttribute{
                MarkdownDescription: "Input Interface Index",
                Computed: true,
            },
            "output_interface_index": schema.NumberAttribute{
                MarkdownDescription: "Output Interface Index",
                Computed: true,
            },
            "octets": schema.StringAttribute{
                MarkdownDescription: "Octets",
                Computed: true,
            },
            "packets": schema.StringAttribute{
                MarkdownDescription: "Packets",
                Computed: true,
            },
            "flow_start_at": schema.StringAttribute{
                MarkdownDescription: "Flow Start",
                Computed: true,
            },
            "flow_end_at": schema.StringAttribute{
                MarkdownDescription: "Flow End",
                Computed: true,
            },
            "ingested_at": schema.StringAttribute{
                MarkdownDescription: "Ingested At",
                Computed: true,
            },
        },
    }
}

func (d *NetworkFlowDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *NetworkFlowDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data NetworkFlowDataSourceModel

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
            "Exactly one of `id` or `name` must be set to look up a network_flow.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "name": true,
        "projectId": true,
        "networkDeviceId": true,
        "exporterIp": true,
        "srcIp": true,
        "dstIp": true,
        "srcPort": true,
        "dstPort": true,
        "protocol": true,
        "inputInterfaceIndex": true,
        "outputInterfaceIndex": true,
        "octets": true,
        "packets": true,
        "flowStartAt": true,
        "flowEndAt": true,
        "ingestedAt": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/network-flow/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read network_flow, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No network_flow found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read network_flow: %s", err))
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
        httpResp, err := d.client.Post(ctx, "/network-flow/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list network_flow, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list network_flow: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No network_flow found with name %q.", data.Name.ValueString()))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one network_flow matches name %q. Use the id attribute to disambiguate.", data.Name.ValueString()))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for network_flow.")
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
    if obj, ok := item["exporterIp"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ExporterIp = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ExporterIp = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ExporterIp = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ExporterIp = types.StringValue(string(jsonBytes))
        } else {
            data.ExporterIp = types.StringNull()
        }
    } else if val, ok := item["exporterIp"].(string); ok {
        data.ExporterIp = types.StringValue(val)
    } else {
        data.ExporterIp = types.StringNull()
    }
    if obj, ok := item["srcIp"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SrcIp = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SrcIp = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SrcIp = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SrcIp = types.StringValue(string(jsonBytes))
        } else {
            data.SrcIp = types.StringNull()
        }
    } else if val, ok := item["srcIp"].(string); ok {
        data.SrcIp = types.StringValue(val)
    } else {
        data.SrcIp = types.StringNull()
    }
    if obj, ok := item["dstIp"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DstIp = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DstIp = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DstIp = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DstIp = types.StringValue(string(jsonBytes))
        } else {
            data.DstIp = types.StringNull()
        }
    } else if val, ok := item["dstIp"].(string); ok {
        data.DstIp = types.StringValue(val)
    } else {
        data.DstIp = types.StringNull()
    }
    if val, ok := item["srcPort"].(float64); ok {
        data.SrcPort = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["srcPort"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.SrcPort = types.NumberValue(big.NewFloat(val))
        } else {
            data.SrcPort = types.NumberNull()
        }
    } else {
        data.SrcPort = types.NumberNull()
    }
    if val, ok := item["dstPort"].(float64); ok {
        data.DstPort = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["dstPort"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.DstPort = types.NumberValue(big.NewFloat(val))
        } else {
            data.DstPort = types.NumberNull()
        }
    } else {
        data.DstPort = types.NumberNull()
    }
    if val, ok := item["protocol"].(float64); ok {
        data.Protocol = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["protocol"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.Protocol = types.NumberValue(big.NewFloat(val))
        } else {
            data.Protocol = types.NumberNull()
        }
    } else {
        data.Protocol = types.NumberNull()
    }
    if val, ok := item["inputInterfaceIndex"].(float64); ok {
        data.InputInterfaceIndex = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["inputInterfaceIndex"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.InputInterfaceIndex = types.NumberValue(big.NewFloat(val))
        } else {
            data.InputInterfaceIndex = types.NumberNull()
        }
    } else {
        data.InputInterfaceIndex = types.NumberNull()
    }
    if val, ok := item["outputInterfaceIndex"].(float64); ok {
        data.OutputInterfaceIndex = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["outputInterfaceIndex"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.OutputInterfaceIndex = types.NumberValue(big.NewFloat(val))
        } else {
            data.OutputInterfaceIndex = types.NumberNull()
        }
    } else {
        data.OutputInterfaceIndex = types.NumberNull()
    }
    if obj, ok := item["octets"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Octets = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Octets = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Octets = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Octets = types.StringValue(string(jsonBytes))
        } else {
            data.Octets = types.StringNull()
        }
    } else if val, ok := item["octets"].(string); ok {
        data.Octets = types.StringValue(val)
    } else {
        data.Octets = types.StringNull()
    }
    if obj, ok := item["packets"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Packets = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Packets = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Packets = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Packets = types.StringValue(string(jsonBytes))
        } else {
            data.Packets = types.StringNull()
        }
    } else if val, ok := item["packets"].(string); ok {
        data.Packets = types.StringValue(val)
    } else {
        data.Packets = types.StringNull()
    }
    if obj, ok := item["flowStartAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FlowStartAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.FlowStartAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.FlowStartAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.FlowStartAt = types.StringValue(string(jsonBytes))
        } else {
            data.FlowStartAt = types.StringNull()
        }
    } else if val, ok := item["flowStartAt"].(string); ok {
        data.FlowStartAt = types.StringValue(val)
    } else {
        data.FlowStartAt = types.StringNull()
    }
    if obj, ok := item["flowEndAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FlowEndAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.FlowEndAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.FlowEndAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.FlowEndAt = types.StringValue(string(jsonBytes))
        } else {
            data.FlowEndAt = types.StringNull()
        }
    } else if val, ok := item["flowEndAt"].(string); ok {
        data.FlowEndAt = types.StringValue(val)
    } else {
        data.FlowEndAt = types.StringNull()
    }
    if obj, ok := item["ingestedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IngestedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.IngestedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.IngestedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.IngestedAt = types.StringValue(string(jsonBytes))
        } else {
            data.IngestedAt = types.StringNull()
        }
    } else if val, ok := item["ingestedAt"].(string); ok {
        data.IngestedAt = types.StringValue(val)
    } else {
        data.IngestedAt = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
