package provider

import (
    "context"
    "fmt"
    "math/big"
    "github.com/hashicorp/terraform-plugin-framework/attr"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &KubernetesCostAllocationDataDataSource{}

func NewKubernetesCostAllocationDataDataSource() datasource.DataSource {
    return &KubernetesCostAllocationDataDataSource{}
}

// KubernetesCostAllocationDataDataSource defines the data source implementation.
type KubernetesCostAllocationDataDataSource struct {
    client *Client
}

// KubernetesCostAllocationDataDataSourceModel describes the data source data model.
type KubernetesCostAllocationDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    ProjectId types.String `tfsdk:"project_id"`
    KubernetesClusterId types.String `tfsdk:"kubernetes_cluster_id"`
    ClusterName types.String `tfsdk:"cluster_name"`
    K8sClusterEntityKey types.String `tfsdk:"k8s_cluster_entity_key"`
    WindowStart types.String `tfsdk:"window_start"`
    WindowEnd types.String `tfsdk:"window_end"`
    Namespace types.String `tfsdk:"namespace"`
    ControllerKind types.String `tfsdk:"controller_kind"`
    ControllerName types.String `tfsdk:"controller_name"`
    PodName types.String `tfsdk:"pod_name"`
    ContainerName types.String `tfsdk:"container_name"`
    NodeName types.String `tfsdk:"node_name"`
    ProviderId types.String `tfsdk:"provider_id"`
    Labels types.String `tfsdk:"labels"`
    LabelKeys types.Set `tfsdk:"label_keys"`
    CpuCoreHours types.Number `tfsdk:"cpu_core_hours"`
    CpuCoreRequestAverage types.Number `tfsdk:"cpu_core_request_average"`
    CpuCoreUsageAverage types.Number `tfsdk:"cpu_core_usage_average"`
    CpuCost types.Number `tfsdk:"cpu_cost"`
    GpuHours types.Number `tfsdk:"gpu_hours"`
    GpuCost types.Number `tfsdk:"gpu_cost"`
    RamByteHours types.Number `tfsdk:"ram_byte_hours"`
    RamBytesRequestAverage types.Number `tfsdk:"ram_bytes_request_average"`
    RamBytesUsageAverage types.Number `tfsdk:"ram_bytes_usage_average"`
    RamCost types.Number `tfsdk:"ram_cost"`
    PvByteHours types.Number `tfsdk:"pv_byte_hours"`
    PvCost types.Number `tfsdk:"pv_cost"`
    NetworkCost types.Number `tfsdk:"network_cost"`
    LoadBalancerCost types.Number `tfsdk:"load_balancer_cost"`
    SharedCost types.Number `tfsdk:"shared_cost"`
    ExternalCost types.Number `tfsdk:"external_cost"`
    TotalCost types.Number `tfsdk:"total_cost"`
    CpuEfficiency types.Number `tfsdk:"cpu_efficiency"`
    RamEfficiency types.Number `tfsdk:"ram_efficiency"`
    TotalEfficiency types.Number `tfsdk:"total_efficiency"`
    Currency types.String `tfsdk:"currency"`
}

func (d *KubernetesCostAllocationDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_kubernetes_cost_allocation_data"
}

func (d *KubernetesCostAllocationDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "kubernetes_cost_allocation_data data source",

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
            "kubernetes_cluster_id": schema.StringAttribute{
                MarkdownDescription: "Kubernetes Cluster ID",
                Computed: true,
            },
            "cluster_name": schema.StringAttribute{
                MarkdownDescription: "Cluster Name",
                Computed: true,
            },
            "k8s_cluster_entity_key": schema.StringAttribute{
                MarkdownDescription: "Kubernetes Cluster Entity Key",
                Computed: true,
            },
            "window_start": schema.StringAttribute{
                MarkdownDescription: "Window Start",
                Computed: true,
            },
            "window_end": schema.StringAttribute{
                MarkdownDescription: "Window End",
                Computed: true,
            },
            "namespace": schema.StringAttribute{
                MarkdownDescription: "Namespace",
                Computed: true,
            },
            "controller_kind": schema.StringAttribute{
                MarkdownDescription: "Controller Kind",
                Computed: true,
            },
            "controller_name": schema.StringAttribute{
                MarkdownDescription: "Controller Name",
                Computed: true,
            },
            "pod_name": schema.StringAttribute{
                MarkdownDescription: "Pod Name",
                Computed: true,
            },
            "container_name": schema.StringAttribute{
                MarkdownDescription: "Container Name",
                Computed: true,
            },
            "node_name": schema.StringAttribute{
                MarkdownDescription: "Node Name",
                Computed: true,
            },
            "provider_id": schema.StringAttribute{
                MarkdownDescription: "Provider ID",
                Computed: true,
            },
            "labels": schema.StringAttribute{
                MarkdownDescription: "Labels",
                Computed: true,
            },
            "label_keys": schema.SetAttribute{
                MarkdownDescription: "Label Keys",
                Computed: true,
                ElementType: types.StringType,
            },
            "cpu_core_hours": schema.NumberAttribute{
                MarkdownDescription: "CPU Core Hours",
                Computed: true,
            },
            "cpu_core_request_average": schema.NumberAttribute{
                MarkdownDescription: "CPU Core Request Average",
                Computed: true,
            },
            "cpu_core_usage_average": schema.NumberAttribute{
                MarkdownDescription: "CPU Core Usage Average",
                Computed: true,
            },
            "cpu_cost": schema.NumberAttribute{
                MarkdownDescription: "CPU Cost",
                Computed: true,
            },
            "gpu_hours": schema.NumberAttribute{
                MarkdownDescription: "GPU Hours",
                Computed: true,
            },
            "gpu_cost": schema.NumberAttribute{
                MarkdownDescription: "GPU Cost",
                Computed: true,
            },
            "ram_byte_hours": schema.NumberAttribute{
                MarkdownDescription: "RAM Byte Hours",
                Computed: true,
            },
            "ram_bytes_request_average": schema.NumberAttribute{
                MarkdownDescription: "RAM Bytes Request Average",
                Computed: true,
            },
            "ram_bytes_usage_average": schema.NumberAttribute{
                MarkdownDescription: "RAM Bytes Usage Average",
                Computed: true,
            },
            "ram_cost": schema.NumberAttribute{
                MarkdownDescription: "RAM Cost",
                Computed: true,
            },
            "pv_byte_hours": schema.NumberAttribute{
                MarkdownDescription: "PV Byte Hours",
                Computed: true,
            },
            "pv_cost": schema.NumberAttribute{
                MarkdownDescription: "PV Cost",
                Computed: true,
            },
            "network_cost": schema.NumberAttribute{
                MarkdownDescription: "Network Cost",
                Computed: true,
            },
            "load_balancer_cost": schema.NumberAttribute{
                MarkdownDescription: "Load Balancer Cost",
                Computed: true,
            },
            "shared_cost": schema.NumberAttribute{
                MarkdownDescription: "Shared Cost",
                Computed: true,
            },
            "external_cost": schema.NumberAttribute{
                MarkdownDescription: "External Cost",
                Computed: true,
            },
            "total_cost": schema.NumberAttribute{
                MarkdownDescription: "Total Cost",
                Computed: true,
            },
            "cpu_efficiency": schema.NumberAttribute{
                MarkdownDescription: "CPU Efficiency",
                Computed: true,
            },
            "ram_efficiency": schema.NumberAttribute{
                MarkdownDescription: "RAM Efficiency",
                Computed: true,
            },
            "total_efficiency": schema.NumberAttribute{
                MarkdownDescription: "Total Efficiency",
                Computed: true,
            },
            "currency": schema.StringAttribute{
                MarkdownDescription: "Currency",
                Computed: true,
            },
        },
    }
}

func (d *KubernetesCostAllocationDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *KubernetesCostAllocationDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data KubernetesCostAllocationDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "kubernetes-cost-allocation" + "/" + data.Id.ValueString()
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read kubernetes_cost_allocation_data, got error: %s", err))
        return
    }

    var kubernetesCostAllocationDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &kubernetesCostAllocationDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse kubernetes_cost_allocation_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := kubernetesCostAllocationDataResponse["data"].(map[string]interface{}); ok {
        kubernetesCostAllocationDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := kubernetesCostAllocationDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := kubernetesCostAllocationDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := kubernetesCostAllocationDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := kubernetesCostAllocationDataResponse["kubernetes_cluster_id"].(string); ok {
        data.KubernetesClusterId = types.StringValue(val)
    }
    if val, ok := kubernetesCostAllocationDataResponse["cluster_name"].(string); ok {
        data.ClusterName = types.StringValue(val)
    }
    if val, ok := kubernetesCostAllocationDataResponse["k8s_cluster_entity_key"].(string); ok {
        data.K8sClusterEntityKey = types.StringValue(val)
    }
    if val, ok := kubernetesCostAllocationDataResponse["window_start"].(string); ok {
        data.WindowStart = types.StringValue(val)
    }
    if val, ok := kubernetesCostAllocationDataResponse["window_end"].(string); ok {
        data.WindowEnd = types.StringValue(val)
    }
    if val, ok := kubernetesCostAllocationDataResponse["namespace"].(string); ok {
        data.Namespace = types.StringValue(val)
    }
    if val, ok := kubernetesCostAllocationDataResponse["controller_kind"].(string); ok {
        data.ControllerKind = types.StringValue(val)
    }
    if val, ok := kubernetesCostAllocationDataResponse["controller_name"].(string); ok {
        data.ControllerName = types.StringValue(val)
    }
    if val, ok := kubernetesCostAllocationDataResponse["pod_name"].(string); ok {
        data.PodName = types.StringValue(val)
    }
    if val, ok := kubernetesCostAllocationDataResponse["container_name"].(string); ok {
        data.ContainerName = types.StringValue(val)
    }
    if val, ok := kubernetesCostAllocationDataResponse["node_name"].(string); ok {
        data.NodeName = types.StringValue(val)
    }
    if val, ok := kubernetesCostAllocationDataResponse["provider_id"].(string); ok {
        data.ProviderId = types.StringValue(val)
    }
    if val, ok := kubernetesCostAllocationDataResponse["labels"].(string); ok {
        data.Labels = types.StringValue(val)
    }
    if val, ok := kubernetesCostAllocationDataResponse["label_keys"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.LabelKeys = setValue
    }
    if val, ok := kubernetesCostAllocationDataResponse["cpu_core_hours"].(float64); ok {
        data.CpuCoreHours = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := kubernetesCostAllocationDataResponse["cpu_core_request_average"].(float64); ok {
        data.CpuCoreRequestAverage = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := kubernetesCostAllocationDataResponse["cpu_core_usage_average"].(float64); ok {
        data.CpuCoreUsageAverage = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := kubernetesCostAllocationDataResponse["cpu_cost"].(float64); ok {
        data.CpuCost = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := kubernetesCostAllocationDataResponse["gpu_hours"].(float64); ok {
        data.GpuHours = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := kubernetesCostAllocationDataResponse["gpu_cost"].(float64); ok {
        data.GpuCost = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := kubernetesCostAllocationDataResponse["ram_byte_hours"].(float64); ok {
        data.RamByteHours = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := kubernetesCostAllocationDataResponse["ram_bytes_request_average"].(float64); ok {
        data.RamBytesRequestAverage = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := kubernetesCostAllocationDataResponse["ram_bytes_usage_average"].(float64); ok {
        data.RamBytesUsageAverage = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := kubernetesCostAllocationDataResponse["ram_cost"].(float64); ok {
        data.RamCost = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := kubernetesCostAllocationDataResponse["pv_byte_hours"].(float64); ok {
        data.PvByteHours = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := kubernetesCostAllocationDataResponse["pv_cost"].(float64); ok {
        data.PvCost = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := kubernetesCostAllocationDataResponse["network_cost"].(float64); ok {
        data.NetworkCost = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := kubernetesCostAllocationDataResponse["load_balancer_cost"].(float64); ok {
        data.LoadBalancerCost = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := kubernetesCostAllocationDataResponse["shared_cost"].(float64); ok {
        data.SharedCost = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := kubernetesCostAllocationDataResponse["external_cost"].(float64); ok {
        data.ExternalCost = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := kubernetesCostAllocationDataResponse["total_cost"].(float64); ok {
        data.TotalCost = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := kubernetesCostAllocationDataResponse["cpu_efficiency"].(float64); ok {
        data.CpuEfficiency = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := kubernetesCostAllocationDataResponse["ram_efficiency"].(float64); ok {
        data.RamEfficiency = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := kubernetesCostAllocationDataResponse["total_efficiency"].(float64); ok {
        data.TotalEfficiency = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := kubernetesCostAllocationDataResponse["currency"].(string); ok {
        data.Currency = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
