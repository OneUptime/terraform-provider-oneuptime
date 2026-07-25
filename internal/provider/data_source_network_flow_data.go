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
var _ datasource.DataSource = &NetworkFlowDataDataSource{}

func NewNetworkFlowDataDataSource() datasource.DataSource {
    return &NetworkFlowDataDataSource{}
}

// NetworkFlowDataDataSource defines the data source implementation.
type NetworkFlowDataDataSource struct {
    client *Client
}

// NetworkFlowDataDataSourceModel describes the data source data model.
type NetworkFlowDataDataSourceModel struct {
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

func (d *NetworkFlowDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_network_flow_data"
}

func (d *NetworkFlowDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "network_flow_data data source",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Identifier to filter by",
                Optional: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Name to filter by",
                Optional: true,
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

func (d *NetworkFlowDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *NetworkFlowDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data NetworkFlowDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "network-flow" + "/" + data.Id.ValueString()
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read network_flow_data, got error: %s", err))
        return
    }

    var networkFlowDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &networkFlowDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse network_flow_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := networkFlowDataResponse["data"].(map[string]interface{}); ok {
        networkFlowDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := networkFlowDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := networkFlowDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := networkFlowDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := networkFlowDataResponse["network_device_id"].(string); ok {
        data.NetworkDeviceId = types.StringValue(val)
    }
    if val, ok := networkFlowDataResponse["exporter_ip"].(string); ok {
        data.ExporterIp = types.StringValue(val)
    }
    if val, ok := networkFlowDataResponse["src_ip"].(string); ok {
        data.SrcIp = types.StringValue(val)
    }
    if val, ok := networkFlowDataResponse["dst_ip"].(string); ok {
        data.DstIp = types.StringValue(val)
    }
    if val, ok := networkFlowDataResponse["src_port"].(float64); ok {
        data.SrcPort = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := networkFlowDataResponse["dst_port"].(float64); ok {
        data.DstPort = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := networkFlowDataResponse["protocol"].(float64); ok {
        data.Protocol = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := networkFlowDataResponse["input_interface_index"].(float64); ok {
        data.InputInterfaceIndex = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := networkFlowDataResponse["output_interface_index"].(float64); ok {
        data.OutputInterfaceIndex = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := networkFlowDataResponse["octets"].(string); ok {
        data.Octets = types.StringValue(val)
    }
    if val, ok := networkFlowDataResponse["packets"].(string); ok {
        data.Packets = types.StringValue(val)
    }
    if val, ok := networkFlowDataResponse["flow_start_at"].(string); ok {
        data.FlowStartAt = types.StringValue(val)
    }
    if val, ok := networkFlowDataResponse["flow_end_at"].(string); ok {
        data.FlowEndAt = types.StringValue(val)
    }
    if val, ok := networkFlowDataResponse["ingested_at"].(string); ok {
        data.IngestedAt = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
