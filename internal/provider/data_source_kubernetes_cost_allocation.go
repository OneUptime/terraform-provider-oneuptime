package provider

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "math/big"
    "github.com/hashicorp/terraform-plugin-framework/attr"
    "sort"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &KubernetesCostAllocationDataSource{}

func NewKubernetesCostAllocationDataSource() datasource.DataSource {
    return &KubernetesCostAllocationDataSource{}
}

// KubernetesCostAllocationDataSource defines the data source implementation.
type KubernetesCostAllocationDataSource struct {
    client *Client
}

// KubernetesCostAllocationDataSourceModel describes the data source data model.
type KubernetesCostAllocationDataSourceModel struct {
    Id types.String `tfsdk:"id"`
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
    CpuCoreLimitAverage types.Number `tfsdk:"cpu_core_limit_average"`
    CpuCost types.Number `tfsdk:"cpu_cost"`
    GpuHours types.Number `tfsdk:"gpu_hours"`
    GpuCost types.Number `tfsdk:"gpu_cost"`
    RamByteHours types.Number `tfsdk:"ram_byte_hours"`
    RamBytesRequestAverage types.Number `tfsdk:"ram_bytes_request_average"`
    RamBytesUsageAverage types.Number `tfsdk:"ram_bytes_usage_average"`
    RamBytesLimitAverage types.Number `tfsdk:"ram_bytes_limit_average"`
    RamBytesUsageMax types.Number `tfsdk:"ram_bytes_usage_max"`
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
    ShipmentId types.String `tfsdk:"shipment_id"`
    ShipmentChunk types.Number `tfsdk:"shipment_chunk"`
}

func (d *KubernetesCostAllocationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_kubernetes_cost_allocation"
}

func (d *KubernetesCostAllocationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "API endpoints for Kubernetes Cost Allocation Look up an existing kubernetes cost allocation by `id`, or by any of its other arguments (`cluster_name`, `container_name`, `controller_kind`, ...): each one set must match, and exactly one kubernetes cost allocation may match them all.",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Look up by unique identifier. Leave unset to look up by the other arguments instead.",
                Optional: true,
                Computed: true,
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "Project ID.",
                Computed: true,
            },
            "kubernetes_cluster_id": schema.StringAttribute{
                MarkdownDescription: "Kubernetes Cluster ID.",
                Optional: true,
                Computed: true,
            },
            "cluster_name": schema.StringAttribute{
                MarkdownDescription: "Cluster Name.",
                Optional: true,
                Computed: true,
            },
            "k8s_cluster_entity_key": schema.StringAttribute{
                MarkdownDescription: "Kubernetes Cluster Entity Key.",
                Optional: true,
                Computed: true,
            },
            "window_start": schema.StringAttribute{
                MarkdownDescription: "Window Start.",
                Optional: true,
                Computed: true,
            },
            "window_end": schema.StringAttribute{
                MarkdownDescription: "Window End.",
                Optional: true,
                Computed: true,
            },
            "namespace": schema.StringAttribute{
                MarkdownDescription: "Namespace.",
                Optional: true,
                Computed: true,
            },
            "controller_kind": schema.StringAttribute{
                MarkdownDescription: "Controller Kind.",
                Optional: true,
                Computed: true,
            },
            "controller_name": schema.StringAttribute{
                MarkdownDescription: "Controller Name.",
                Optional: true,
                Computed: true,
            },
            "pod_name": schema.StringAttribute{
                MarkdownDescription: "Pod Name.",
                Optional: true,
                Computed: true,
            },
            "container_name": schema.StringAttribute{
                MarkdownDescription: "Container Name.",
                Optional: true,
                Computed: true,
            },
            "node_name": schema.StringAttribute{
                MarkdownDescription: "Node Name.",
                Optional: true,
                Computed: true,
            },
            "provider_id": schema.StringAttribute{
                MarkdownDescription: "Provider ID.",
                Optional: true,
                Computed: true,
            },
            "labels": schema.StringAttribute{
                MarkdownDescription: "Labels.",
                Optional: true,
                Computed: true,
            },
            "label_keys": schema.SetAttribute{
                MarkdownDescription: "Label Keys.",
                Computed: true,
                ElementType: types.StringType,
            },
            "cpu_core_hours": schema.NumberAttribute{
                MarkdownDescription: "CPU Core Hours.",
                Optional: true,
                Computed: true,
            },
            "cpu_core_request_average": schema.NumberAttribute{
                MarkdownDescription: "CPU Core Request Average.",
                Optional: true,
                Computed: true,
            },
            "cpu_core_usage_average": schema.NumberAttribute{
                MarkdownDescription: "CPU Core Usage Average.",
                Optional: true,
                Computed: true,
            },
            "cpu_core_limit_average": schema.NumberAttribute{
                MarkdownDescription: "CPU Core Limit Average.",
                Optional: true,
                Computed: true,
            },
            "cpu_cost": schema.NumberAttribute{
                MarkdownDescription: "CPU Cost.",
                Optional: true,
                Computed: true,
            },
            "gpu_hours": schema.NumberAttribute{
                MarkdownDescription: "GPU Hours.",
                Optional: true,
                Computed: true,
            },
            "gpu_cost": schema.NumberAttribute{
                MarkdownDescription: "GPU Cost.",
                Optional: true,
                Computed: true,
            },
            "ram_byte_hours": schema.NumberAttribute{
                MarkdownDescription: "RAM Byte Hours.",
                Optional: true,
                Computed: true,
            },
            "ram_bytes_request_average": schema.NumberAttribute{
                MarkdownDescription: "RAM Bytes Request Average.",
                Optional: true,
                Computed: true,
            },
            "ram_bytes_usage_average": schema.NumberAttribute{
                MarkdownDescription: "RAM Bytes Usage Average.",
                Optional: true,
                Computed: true,
            },
            "ram_bytes_limit_average": schema.NumberAttribute{
                MarkdownDescription: "RAM Bytes Limit Average.",
                Optional: true,
                Computed: true,
            },
            "ram_bytes_usage_max": schema.NumberAttribute{
                MarkdownDescription: "RAM Bytes Usage Max.",
                Optional: true,
                Computed: true,
            },
            "ram_cost": schema.NumberAttribute{
                MarkdownDescription: "RAM Cost.",
                Optional: true,
                Computed: true,
            },
            "pv_byte_hours": schema.NumberAttribute{
                MarkdownDescription: "PV Byte Hours.",
                Optional: true,
                Computed: true,
            },
            "pv_cost": schema.NumberAttribute{
                MarkdownDescription: "PV Cost.",
                Optional: true,
                Computed: true,
            },
            "network_cost": schema.NumberAttribute{
                MarkdownDescription: "Network Cost.",
                Optional: true,
                Computed: true,
            },
            "load_balancer_cost": schema.NumberAttribute{
                MarkdownDescription: "Load Balancer Cost.",
                Optional: true,
                Computed: true,
            },
            "shared_cost": schema.NumberAttribute{
                MarkdownDescription: "Shared Cost.",
                Optional: true,
                Computed: true,
            },
            "external_cost": schema.NumberAttribute{
                MarkdownDescription: "External Cost.",
                Optional: true,
                Computed: true,
            },
            "total_cost": schema.NumberAttribute{
                MarkdownDescription: "Total Cost.",
                Optional: true,
                Computed: true,
            },
            "cpu_efficiency": schema.NumberAttribute{
                MarkdownDescription: "CPU Efficiency.",
                Optional: true,
                Computed: true,
            },
            "ram_efficiency": schema.NumberAttribute{
                MarkdownDescription: "RAM Efficiency.",
                Optional: true,
                Computed: true,
            },
            "total_efficiency": schema.NumberAttribute{
                MarkdownDescription: "Total Efficiency.",
                Optional: true,
                Computed: true,
            },
            "currency": schema.StringAttribute{
                MarkdownDescription: "Currency.",
                Optional: true,
                Computed: true,
            },
            "shipment_id": schema.StringAttribute{
                MarkdownDescription: "Shipment ID.",
                Optional: true,
                Computed: true,
            },
            "shipment_chunk": schema.NumberAttribute{
                MarkdownDescription: "Shipment Chunk.",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *KubernetesCostAllocationDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *KubernetesCostAllocationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data KubernetesCostAllocationDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.KubernetesClusterId.IsNull() && !data.KubernetesClusterId.IsUnknown() {
        filters["kubernetesClusterId"] = data.KubernetesClusterId.ValueString()
        filterNames = append(filterNames, "kubernetes_cluster_id = "+fmt.Sprintf("%q", data.KubernetesClusterId.ValueString()))
    }
    if !data.ClusterName.IsNull() && !data.ClusterName.IsUnknown() {
        filters["clusterName"] = data.ClusterName.ValueString()
        filterNames = append(filterNames, "cluster_name = "+fmt.Sprintf("%q", data.ClusterName.ValueString()))
    }
    if !data.K8sClusterEntityKey.IsNull() && !data.K8sClusterEntityKey.IsUnknown() {
        filters["k8sClusterEntityKey"] = data.K8sClusterEntityKey.ValueString()
        filterNames = append(filterNames, "k8s_cluster_entity_key = "+fmt.Sprintf("%q", data.K8sClusterEntityKey.ValueString()))
    }
    if !data.WindowStart.IsNull() && !data.WindowStart.IsUnknown() {
        filters["windowStart"] = data.WindowStart.ValueString()
        filterNames = append(filterNames, "window_start = "+fmt.Sprintf("%q", data.WindowStart.ValueString()))
    }
    if !data.WindowEnd.IsNull() && !data.WindowEnd.IsUnknown() {
        filters["windowEnd"] = data.WindowEnd.ValueString()
        filterNames = append(filterNames, "window_end = "+fmt.Sprintf("%q", data.WindowEnd.ValueString()))
    }
    if !data.Namespace.IsNull() && !data.Namespace.IsUnknown() {
        filters["namespace"] = data.Namespace.ValueString()
        filterNames = append(filterNames, "namespace = "+fmt.Sprintf("%q", data.Namespace.ValueString()))
    }
    if !data.ControllerKind.IsNull() && !data.ControllerKind.IsUnknown() {
        filters["controllerKind"] = data.ControllerKind.ValueString()
        filterNames = append(filterNames, "controller_kind = "+fmt.Sprintf("%q", data.ControllerKind.ValueString()))
    }
    if !data.ControllerName.IsNull() && !data.ControllerName.IsUnknown() {
        filters["controllerName"] = data.ControllerName.ValueString()
        filterNames = append(filterNames, "controller_name = "+fmt.Sprintf("%q", data.ControllerName.ValueString()))
    }
    if !data.PodName.IsNull() && !data.PodName.IsUnknown() {
        filters["podName"] = data.PodName.ValueString()
        filterNames = append(filterNames, "pod_name = "+fmt.Sprintf("%q", data.PodName.ValueString()))
    }
    if !data.ContainerName.IsNull() && !data.ContainerName.IsUnknown() {
        filters["containerName"] = data.ContainerName.ValueString()
        filterNames = append(filterNames, "container_name = "+fmt.Sprintf("%q", data.ContainerName.ValueString()))
    }
    if !data.NodeName.IsNull() && !data.NodeName.IsUnknown() {
        filters["nodeName"] = data.NodeName.ValueString()
        filterNames = append(filterNames, "node_name = "+fmt.Sprintf("%q", data.NodeName.ValueString()))
    }
    if !data.ProviderId.IsNull() && !data.ProviderId.IsUnknown() {
        filters["providerId"] = data.ProviderId.ValueString()
        filterNames = append(filterNames, "provider_id = "+fmt.Sprintf("%q", data.ProviderId.ValueString()))
    }
    if !data.Labels.IsNull() && !data.Labels.IsUnknown() {
        filters["labels"] = data.Labels.ValueString()
        filterNames = append(filterNames, "labels = "+fmt.Sprintf("%q", data.Labels.ValueString()))
    }
    if !data.CpuCoreHours.IsNull() && !data.CpuCoreHours.IsUnknown() {
        filters["cpuCoreHours"] = lookupNumber(data.CpuCoreHours)
        filterNames = append(filterNames, "cpu_core_hours = "+data.CpuCoreHours.ValueBigFloat().String())
    }
    if !data.CpuCoreRequestAverage.IsNull() && !data.CpuCoreRequestAverage.IsUnknown() {
        filters["cpuCoreRequestAverage"] = lookupNumber(data.CpuCoreRequestAverage)
        filterNames = append(filterNames, "cpu_core_request_average = "+data.CpuCoreRequestAverage.ValueBigFloat().String())
    }
    if !data.CpuCoreUsageAverage.IsNull() && !data.CpuCoreUsageAverage.IsUnknown() {
        filters["cpuCoreUsageAverage"] = lookupNumber(data.CpuCoreUsageAverage)
        filterNames = append(filterNames, "cpu_core_usage_average = "+data.CpuCoreUsageAverage.ValueBigFloat().String())
    }
    if !data.CpuCoreLimitAverage.IsNull() && !data.CpuCoreLimitAverage.IsUnknown() {
        filters["cpuCoreLimitAverage"] = lookupNumber(data.CpuCoreLimitAverage)
        filterNames = append(filterNames, "cpu_core_limit_average = "+data.CpuCoreLimitAverage.ValueBigFloat().String())
    }
    if !data.CpuCost.IsNull() && !data.CpuCost.IsUnknown() {
        filters["cpuCost"] = lookupNumber(data.CpuCost)
        filterNames = append(filterNames, "cpu_cost = "+data.CpuCost.ValueBigFloat().String())
    }
    if !data.GpuHours.IsNull() && !data.GpuHours.IsUnknown() {
        filters["gpuHours"] = lookupNumber(data.GpuHours)
        filterNames = append(filterNames, "gpu_hours = "+data.GpuHours.ValueBigFloat().String())
    }
    if !data.GpuCost.IsNull() && !data.GpuCost.IsUnknown() {
        filters["gpuCost"] = lookupNumber(data.GpuCost)
        filterNames = append(filterNames, "gpu_cost = "+data.GpuCost.ValueBigFloat().String())
    }
    if !data.RamByteHours.IsNull() && !data.RamByteHours.IsUnknown() {
        filters["ramByteHours"] = lookupNumber(data.RamByteHours)
        filterNames = append(filterNames, "ram_byte_hours = "+data.RamByteHours.ValueBigFloat().String())
    }
    if !data.RamBytesRequestAverage.IsNull() && !data.RamBytesRequestAverage.IsUnknown() {
        filters["ramBytesRequestAverage"] = lookupNumber(data.RamBytesRequestAverage)
        filterNames = append(filterNames, "ram_bytes_request_average = "+data.RamBytesRequestAverage.ValueBigFloat().String())
    }
    if !data.RamBytesUsageAverage.IsNull() && !data.RamBytesUsageAverage.IsUnknown() {
        filters["ramBytesUsageAverage"] = lookupNumber(data.RamBytesUsageAverage)
        filterNames = append(filterNames, "ram_bytes_usage_average = "+data.RamBytesUsageAverage.ValueBigFloat().String())
    }
    if !data.RamBytesLimitAverage.IsNull() && !data.RamBytesLimitAverage.IsUnknown() {
        filters["ramBytesLimitAverage"] = lookupNumber(data.RamBytesLimitAverage)
        filterNames = append(filterNames, "ram_bytes_limit_average = "+data.RamBytesLimitAverage.ValueBigFloat().String())
    }
    if !data.RamBytesUsageMax.IsNull() && !data.RamBytesUsageMax.IsUnknown() {
        filters["ramBytesUsageMax"] = lookupNumber(data.RamBytesUsageMax)
        filterNames = append(filterNames, "ram_bytes_usage_max = "+data.RamBytesUsageMax.ValueBigFloat().String())
    }
    if !data.RamCost.IsNull() && !data.RamCost.IsUnknown() {
        filters["ramCost"] = lookupNumber(data.RamCost)
        filterNames = append(filterNames, "ram_cost = "+data.RamCost.ValueBigFloat().String())
    }
    if !data.PvByteHours.IsNull() && !data.PvByteHours.IsUnknown() {
        filters["pvByteHours"] = lookupNumber(data.PvByteHours)
        filterNames = append(filterNames, "pv_byte_hours = "+data.PvByteHours.ValueBigFloat().String())
    }
    if !data.PvCost.IsNull() && !data.PvCost.IsUnknown() {
        filters["pvCost"] = lookupNumber(data.PvCost)
        filterNames = append(filterNames, "pv_cost = "+data.PvCost.ValueBigFloat().String())
    }
    if !data.NetworkCost.IsNull() && !data.NetworkCost.IsUnknown() {
        filters["networkCost"] = lookupNumber(data.NetworkCost)
        filterNames = append(filterNames, "network_cost = "+data.NetworkCost.ValueBigFloat().String())
    }
    if !data.LoadBalancerCost.IsNull() && !data.LoadBalancerCost.IsUnknown() {
        filters["loadBalancerCost"] = lookupNumber(data.LoadBalancerCost)
        filterNames = append(filterNames, "load_balancer_cost = "+data.LoadBalancerCost.ValueBigFloat().String())
    }
    if !data.SharedCost.IsNull() && !data.SharedCost.IsUnknown() {
        filters["sharedCost"] = lookupNumber(data.SharedCost)
        filterNames = append(filterNames, "shared_cost = "+data.SharedCost.ValueBigFloat().String())
    }
    if !data.ExternalCost.IsNull() && !data.ExternalCost.IsUnknown() {
        filters["externalCost"] = lookupNumber(data.ExternalCost)
        filterNames = append(filterNames, "external_cost = "+data.ExternalCost.ValueBigFloat().String())
    }
    if !data.TotalCost.IsNull() && !data.TotalCost.IsUnknown() {
        filters["totalCost"] = lookupNumber(data.TotalCost)
        filterNames = append(filterNames, "total_cost = "+data.TotalCost.ValueBigFloat().String())
    }
    if !data.CpuEfficiency.IsNull() && !data.CpuEfficiency.IsUnknown() {
        filters["cpuEfficiency"] = lookupNumber(data.CpuEfficiency)
        filterNames = append(filterNames, "cpu_efficiency = "+data.CpuEfficiency.ValueBigFloat().String())
    }
    if !data.RamEfficiency.IsNull() && !data.RamEfficiency.IsUnknown() {
        filters["ramEfficiency"] = lookupNumber(data.RamEfficiency)
        filterNames = append(filterNames, "ram_efficiency = "+data.RamEfficiency.ValueBigFloat().String())
    }
    if !data.TotalEfficiency.IsNull() && !data.TotalEfficiency.IsUnknown() {
        filters["totalEfficiency"] = lookupNumber(data.TotalEfficiency)
        filterNames = append(filterNames, "total_efficiency = "+data.TotalEfficiency.ValueBigFloat().String())
    }
    if !data.Currency.IsNull() && !data.Currency.IsUnknown() {
        filters["currency"] = data.Currency.ValueString()
        filterNames = append(filterNames, "currency = "+fmt.Sprintf("%q", data.Currency.ValueString()))
    }
    if !data.ShipmentId.IsNull() && !data.ShipmentId.IsUnknown() {
        filters["shipmentId"] = data.ShipmentId.ValueString()
        filterNames = append(filterNames, "shipment_id = "+fmt.Sprintf("%q", data.ShipmentId.ValueString()))
    }
    if !data.ShipmentChunk.IsNull() && !data.ShipmentChunk.IsUnknown() {
        filters["shipmentChunk"] = lookupNumber(data.ShipmentChunk)
        filterNames = append(filterNames, "shipment_chunk = "+data.ShipmentChunk.ValueBigFloat().String())
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the kubernetes cost allocation up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the kubernetes cost allocation up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "projectId": true,
        "kubernetesClusterId": true,
        "clusterName": true,
        "k8sClusterEntityKey": true,
        "windowStart": true,
        "windowEnd": true,
        "namespace": true,
        "controllerKind": true,
        "controllerName": true,
        "podName": true,
        "containerName": true,
        "nodeName": true,
        "providerId": true,
        "labels": true,
        "labelKeys": true,
        "cpuCoreHours": true,
        "cpuCoreRequestAverage": true,
        "cpuCoreUsageAverage": true,
        "cpuCoreLimitAverage": true,
        "cpuCost": true,
        "gpuHours": true,
        "gpuCost": true,
        "ramByteHours": true,
        "ramBytesRequestAverage": true,
        "ramBytesUsageAverage": true,
        "ramBytesLimitAverage": true,
        "ramBytesUsageMax": true,
        "ramCost": true,
        "pvByteHours": true,
        "pvCost": true,
        "networkCost": true,
        "loadBalancerCost": true,
        "sharedCost": true,
        "externalCost": true,
        "totalCost": true,
        "cpuEfficiency": true,
        "ramEfficiency": true,
        "totalEfficiency": true,
        "currency": true,
        "shipmentId": true,
        "shipmentChunk": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/kubernetes-cost-allocation/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read kubernetes_cost_allocation, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No kubernetes cost allocation found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read kubernetes_cost_allocation: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/kubernetes-cost-allocation/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list kubernetes_cost_allocation, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list kubernetes_cost_allocation: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No kubernetes cost allocation matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one kubernetes cost allocation matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for kubernetes_cost_allocation.")
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
    if obj, ok := item["kubernetesClusterId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.KubernetesClusterId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.KubernetesClusterId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.KubernetesClusterId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.KubernetesClusterId = types.StringValue(string(jsonBytes))
        } else {
            data.KubernetesClusterId = types.StringNull()
        }
    } else if val, ok := item["kubernetesClusterId"].(string); ok {
        data.KubernetesClusterId = types.StringValue(val)
    } else {
        data.KubernetesClusterId = types.StringNull()
    }
    if obj, ok := item["clusterName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ClusterName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ClusterName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ClusterName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ClusterName = types.StringValue(string(jsonBytes))
        } else {
            data.ClusterName = types.StringNull()
        }
    } else if val, ok := item["clusterName"].(string); ok {
        data.ClusterName = types.StringValue(val)
    } else {
        data.ClusterName = types.StringNull()
    }
    if obj, ok := item["k8sClusterEntityKey"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.K8sClusterEntityKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.K8sClusterEntityKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.K8sClusterEntityKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.K8sClusterEntityKey = types.StringValue(string(jsonBytes))
        } else {
            data.K8sClusterEntityKey = types.StringNull()
        }
    } else if val, ok := item["k8sClusterEntityKey"].(string); ok {
        data.K8sClusterEntityKey = types.StringValue(val)
    } else {
        data.K8sClusterEntityKey = types.StringNull()
    }
    if obj, ok := item["windowStart"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.WindowStart = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.WindowStart = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.WindowStart = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.WindowStart = types.StringValue(string(jsonBytes))
        } else {
            data.WindowStart = types.StringNull()
        }
    } else if val, ok := item["windowStart"].(string); ok {
        data.WindowStart = types.StringValue(val)
    } else {
        data.WindowStart = types.StringNull()
    }
    if obj, ok := item["windowEnd"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.WindowEnd = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.WindowEnd = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.WindowEnd = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.WindowEnd = types.StringValue(string(jsonBytes))
        } else {
            data.WindowEnd = types.StringNull()
        }
    } else if val, ok := item["windowEnd"].(string); ok {
        data.WindowEnd = types.StringValue(val)
    } else {
        data.WindowEnd = types.StringNull()
    }
    if obj, ok := item["namespace"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Namespace = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Namespace = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Namespace = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Namespace = types.StringValue(string(jsonBytes))
        } else {
            data.Namespace = types.StringNull()
        }
    } else if val, ok := item["namespace"].(string); ok {
        data.Namespace = types.StringValue(val)
    } else {
        data.Namespace = types.StringNull()
    }
    if obj, ok := item["controllerKind"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ControllerKind = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ControllerKind = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ControllerKind = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ControllerKind = types.StringValue(string(jsonBytes))
        } else {
            data.ControllerKind = types.StringNull()
        }
    } else if val, ok := item["controllerKind"].(string); ok {
        data.ControllerKind = types.StringValue(val)
    } else {
        data.ControllerKind = types.StringNull()
    }
    if obj, ok := item["controllerName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ControllerName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ControllerName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ControllerName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ControllerName = types.StringValue(string(jsonBytes))
        } else {
            data.ControllerName = types.StringNull()
        }
    } else if val, ok := item["controllerName"].(string); ok {
        data.ControllerName = types.StringValue(val)
    } else {
        data.ControllerName = types.StringNull()
    }
    if obj, ok := item["podName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PodName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.PodName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.PodName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.PodName = types.StringValue(string(jsonBytes))
        } else {
            data.PodName = types.StringNull()
        }
    } else if val, ok := item["podName"].(string); ok {
        data.PodName = types.StringValue(val)
    } else {
        data.PodName = types.StringNull()
    }
    if obj, ok := item["containerName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ContainerName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ContainerName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ContainerName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ContainerName = types.StringValue(string(jsonBytes))
        } else {
            data.ContainerName = types.StringNull()
        }
    } else if val, ok := item["containerName"].(string); ok {
        data.ContainerName = types.StringValue(val)
    } else {
        data.ContainerName = types.StringNull()
    }
    if obj, ok := item["nodeName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.NodeName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.NodeName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.NodeName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.NodeName = types.StringValue(string(jsonBytes))
        } else {
            data.NodeName = types.StringNull()
        }
    } else if val, ok := item["nodeName"].(string); ok {
        data.NodeName = types.StringValue(val)
    } else {
        data.NodeName = types.StringNull()
    }
    if obj, ok := item["providerId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ProviderId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ProviderId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ProviderId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ProviderId = types.StringValue(string(jsonBytes))
        } else {
            data.ProviderId = types.StringNull()
        }
    } else if val, ok := item["providerId"].(string); ok {
        data.ProviderId = types.StringValue(val)
    } else {
        data.ProviderId = types.StringNull()
    }
    if obj, ok := item["labels"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Labels = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Labels = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Labels = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Labels = types.StringValue(string(jsonBytes))
        } else {
            data.Labels = types.StringNull()
        }
    } else if val, ok := item["labels"].(string); ok {
        data.Labels = types.StringValue(val)
    } else {
        data.Labels = types.StringNull()
    }
    if val, ok := item["labelKeys"].([]interface{}); ok {
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if jsonBytes, err := json.Marshal(itemMap); err == nil {
                    setItems = append(setItems, types.StringValue(string(jsonBytes)))
                }
            } else if str, ok := item.(string); ok {
                setItems = append(setItems, types.StringValue(str))
            } else {
                setItems = append(setItems, types.StringValue(fmt.Sprintf("%v", item)))
            }
        }
        sort.Slice(setItems, func(i, j int) bool {
            return setItems[i].(types.String).ValueString() < setItems[j].(types.String).ValueString()
        })
        data.LabelKeys = types.SetValueMust(types.StringType, setItems)
    } else {
        data.LabelKeys = types.SetNull(types.StringType)
    }
    if val, ok := item["cpuCoreHours"].(float64); ok {
        data.CpuCoreHours = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["cpuCoreHours"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.CpuCoreHours = types.NumberValue(big.NewFloat(val))
        } else {
            data.CpuCoreHours = types.NumberNull()
        }
    } else {
        data.CpuCoreHours = types.NumberNull()
    }
    if val, ok := item["cpuCoreRequestAverage"].(float64); ok {
        data.CpuCoreRequestAverage = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["cpuCoreRequestAverage"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.CpuCoreRequestAverage = types.NumberValue(big.NewFloat(val))
        } else {
            data.CpuCoreRequestAverage = types.NumberNull()
        }
    } else {
        data.CpuCoreRequestAverage = types.NumberNull()
    }
    if val, ok := item["cpuCoreUsageAverage"].(float64); ok {
        data.CpuCoreUsageAverage = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["cpuCoreUsageAverage"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.CpuCoreUsageAverage = types.NumberValue(big.NewFloat(val))
        } else {
            data.CpuCoreUsageAverage = types.NumberNull()
        }
    } else {
        data.CpuCoreUsageAverage = types.NumberNull()
    }
    if val, ok := item["cpuCoreLimitAverage"].(float64); ok {
        data.CpuCoreLimitAverage = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["cpuCoreLimitAverage"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.CpuCoreLimitAverage = types.NumberValue(big.NewFloat(val))
        } else {
            data.CpuCoreLimitAverage = types.NumberNull()
        }
    } else {
        data.CpuCoreLimitAverage = types.NumberNull()
    }
    if val, ok := item["cpuCost"].(float64); ok {
        data.CpuCost = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["cpuCost"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.CpuCost = types.NumberValue(big.NewFloat(val))
        } else {
            data.CpuCost = types.NumberNull()
        }
    } else {
        data.CpuCost = types.NumberNull()
    }
    if val, ok := item["gpuHours"].(float64); ok {
        data.GpuHours = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["gpuHours"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.GpuHours = types.NumberValue(big.NewFloat(val))
        } else {
            data.GpuHours = types.NumberNull()
        }
    } else {
        data.GpuHours = types.NumberNull()
    }
    if val, ok := item["gpuCost"].(float64); ok {
        data.GpuCost = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["gpuCost"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.GpuCost = types.NumberValue(big.NewFloat(val))
        } else {
            data.GpuCost = types.NumberNull()
        }
    } else {
        data.GpuCost = types.NumberNull()
    }
    if val, ok := item["ramByteHours"].(float64); ok {
        data.RamByteHours = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["ramByteHours"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.RamByteHours = types.NumberValue(big.NewFloat(val))
        } else {
            data.RamByteHours = types.NumberNull()
        }
    } else {
        data.RamByteHours = types.NumberNull()
    }
    if val, ok := item["ramBytesRequestAverage"].(float64); ok {
        data.RamBytesRequestAverage = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["ramBytesRequestAverage"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.RamBytesRequestAverage = types.NumberValue(big.NewFloat(val))
        } else {
            data.RamBytesRequestAverage = types.NumberNull()
        }
    } else {
        data.RamBytesRequestAverage = types.NumberNull()
    }
    if val, ok := item["ramBytesUsageAverage"].(float64); ok {
        data.RamBytesUsageAverage = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["ramBytesUsageAverage"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.RamBytesUsageAverage = types.NumberValue(big.NewFloat(val))
        } else {
            data.RamBytesUsageAverage = types.NumberNull()
        }
    } else {
        data.RamBytesUsageAverage = types.NumberNull()
    }
    if val, ok := item["ramBytesLimitAverage"].(float64); ok {
        data.RamBytesLimitAverage = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["ramBytesLimitAverage"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.RamBytesLimitAverage = types.NumberValue(big.NewFloat(val))
        } else {
            data.RamBytesLimitAverage = types.NumberNull()
        }
    } else {
        data.RamBytesLimitAverage = types.NumberNull()
    }
    if val, ok := item["ramBytesUsageMax"].(float64); ok {
        data.RamBytesUsageMax = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["ramBytesUsageMax"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.RamBytesUsageMax = types.NumberValue(big.NewFloat(val))
        } else {
            data.RamBytesUsageMax = types.NumberNull()
        }
    } else {
        data.RamBytesUsageMax = types.NumberNull()
    }
    if val, ok := item["ramCost"].(float64); ok {
        data.RamCost = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["ramCost"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.RamCost = types.NumberValue(big.NewFloat(val))
        } else {
            data.RamCost = types.NumberNull()
        }
    } else {
        data.RamCost = types.NumberNull()
    }
    if val, ok := item["pvByteHours"].(float64); ok {
        data.PvByteHours = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["pvByteHours"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.PvByteHours = types.NumberValue(big.NewFloat(val))
        } else {
            data.PvByteHours = types.NumberNull()
        }
    } else {
        data.PvByteHours = types.NumberNull()
    }
    if val, ok := item["pvCost"].(float64); ok {
        data.PvCost = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["pvCost"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.PvCost = types.NumberValue(big.NewFloat(val))
        } else {
            data.PvCost = types.NumberNull()
        }
    } else {
        data.PvCost = types.NumberNull()
    }
    if val, ok := item["networkCost"].(float64); ok {
        data.NetworkCost = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["networkCost"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.NetworkCost = types.NumberValue(big.NewFloat(val))
        } else {
            data.NetworkCost = types.NumberNull()
        }
    } else {
        data.NetworkCost = types.NumberNull()
    }
    if val, ok := item["loadBalancerCost"].(float64); ok {
        data.LoadBalancerCost = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["loadBalancerCost"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.LoadBalancerCost = types.NumberValue(big.NewFloat(val))
        } else {
            data.LoadBalancerCost = types.NumberNull()
        }
    } else {
        data.LoadBalancerCost = types.NumberNull()
    }
    if val, ok := item["sharedCost"].(float64); ok {
        data.SharedCost = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["sharedCost"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.SharedCost = types.NumberValue(big.NewFloat(val))
        } else {
            data.SharedCost = types.NumberNull()
        }
    } else {
        data.SharedCost = types.NumberNull()
    }
    if val, ok := item["externalCost"].(float64); ok {
        data.ExternalCost = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["externalCost"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.ExternalCost = types.NumberValue(big.NewFloat(val))
        } else {
            data.ExternalCost = types.NumberNull()
        }
    } else {
        data.ExternalCost = types.NumberNull()
    }
    if val, ok := item["totalCost"].(float64); ok {
        data.TotalCost = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["totalCost"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.TotalCost = types.NumberValue(big.NewFloat(val))
        } else {
            data.TotalCost = types.NumberNull()
        }
    } else {
        data.TotalCost = types.NumberNull()
    }
    if val, ok := item["cpuEfficiency"].(float64); ok {
        data.CpuEfficiency = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["cpuEfficiency"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.CpuEfficiency = types.NumberValue(big.NewFloat(val))
        } else {
            data.CpuEfficiency = types.NumberNull()
        }
    } else {
        data.CpuEfficiency = types.NumberNull()
    }
    if val, ok := item["ramEfficiency"].(float64); ok {
        data.RamEfficiency = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["ramEfficiency"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.RamEfficiency = types.NumberValue(big.NewFloat(val))
        } else {
            data.RamEfficiency = types.NumberNull()
        }
    } else {
        data.RamEfficiency = types.NumberNull()
    }
    if val, ok := item["totalEfficiency"].(float64); ok {
        data.TotalEfficiency = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["totalEfficiency"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.TotalEfficiency = types.NumberValue(big.NewFloat(val))
        } else {
            data.TotalEfficiency = types.NumberNull()
        }
    } else {
        data.TotalEfficiency = types.NumberNull()
    }
    if obj, ok := item["currency"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Currency = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Currency = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Currency = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Currency = types.StringValue(string(jsonBytes))
        } else {
            data.Currency = types.StringNull()
        }
    } else if val, ok := item["currency"].(string); ok {
        data.Currency = types.StringValue(val)
    } else {
        data.Currency = types.StringNull()
    }
    if obj, ok := item["shipmentId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ShipmentId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ShipmentId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ShipmentId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ShipmentId = types.StringValue(string(jsonBytes))
        } else {
            data.ShipmentId = types.StringNull()
        }
    } else if val, ok := item["shipmentId"].(string); ok {
        data.ShipmentId = types.StringValue(val)
    } else {
        data.ShipmentId = types.StringNull()
    }
    if val, ok := item["shipmentChunk"].(float64); ok {
        data.ShipmentChunk = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["shipmentChunk"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.ShipmentChunk = types.NumberValue(big.NewFloat(val))
        } else {
            data.ShipmentChunk = types.NumberNull()
        }
    } else {
        data.ShipmentChunk = types.NumberNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
