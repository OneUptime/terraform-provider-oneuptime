package provider

import (
    "context"
    "fmt"
    "github.com/hashicorp/terraform-plugin-framework/path"
    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-framework/types/basetypes"
    "github.com/hashicorp/terraform-plugin-log/tflog"
    "math/big"
    "net/http"
    "encoding/json"
    "net/url"
    "strings"
    "github.com/hashicorp/terraform-plugin-framework/attr"
    "sort"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &KubernetesCostAllocationResource{}
var _ resource.ResourceWithImportState = &KubernetesCostAllocationResource{}

func NewKubernetesCostAllocationResource() resource.Resource {
    return &KubernetesCostAllocationResource{}
}

// KubernetesCostAllocationResource defines the resource implementation.
type KubernetesCostAllocationResource struct {
    client *Client
}

// KubernetesCostAllocationResourceModel describes the resource data model.
type KubernetesCostAllocationResourceModel struct {
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

func (r *KubernetesCostAllocationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_kubernetes_cost_allocation"
}

func (r *KubernetesCostAllocationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "kubernetes_cost_allocation resource",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Unique identifier for the resource",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
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
            "cpu_core_limit_average": schema.NumberAttribute{
                MarkdownDescription: "CPU Core Limit Average",
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
            "ram_bytes_limit_average": schema.NumberAttribute{
                MarkdownDescription: "RAM Bytes Limit Average",
                Computed: true,
            },
            "ram_bytes_usage_max": schema.NumberAttribute{
                MarkdownDescription: "RAM Bytes Usage Max",
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
            "shipment_id": schema.StringAttribute{
                MarkdownDescription: "Shipment ID",
                Computed: true,
            },
            "shipment_chunk": schema.NumberAttribute{
                MarkdownDescription: "Shipment Chunk",
                Computed: true,
            },
        },
    }
}

func (r *KubernetesCostAllocationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
    // Prevent panic if the provider has not been configured.
    if req.ProviderData == nil {
        return
    }

    client, ok := req.ProviderData.(*Client)

    if !ok {
        resp.Diagnostics.AddError(
            "Unexpected Resource Configure Type",
            fmt.Sprintf("Expected *Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
        )

        return
    }

    r.client = client
}


func (r *KubernetesCostAllocationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var data KubernetesCostAllocationResourceModel

    // Read Terraform plan data into the model
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }



    // Create API request body
    kubernetesCostAllocationRequest := map[string]interface{}{
        "data": map[string]interface{}{

        },
    }

    // Make API call
    httpResp, err := r.client.Post("/kubernetes-cost-allocation", kubernetesCostAllocationRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create kubernetes_cost_allocation, got error: %s", err))
        return
    }

    var kubernetesCostAllocationResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &kubernetesCostAllocationResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse kubernetes_cost_allocation response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := kubernetesCostAllocationResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = kubernetesCostAllocationResponse
    }

    if obj, ok := dataMap["id"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Id = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Id = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Id = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Id = types.StringValue(string(jsonBytes))
            } else {
                data.Id = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Id = types.StringValue(string(jsonBytes))
            } else {
                data.Id = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Id = types.StringValue(string(jsonBytes))
        } else {
            data.Id = types.StringNull()
        }
    } else if val, ok := dataMap["id"].(string); ok && val != "" {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    if obj, ok := dataMap["projectId"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok {
            data.ProjectId = types.StringValue(val)
        } else {
            data.ProjectId = types.StringNull()
        }
    } else if val, ok := dataMap["projectId"].(string); ok {
        data.ProjectId = types.StringValue(val)
    } else {
        data.ProjectId = types.StringNull()
    }
    if obj, ok := dataMap["kubernetesClusterId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.KubernetesClusterId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.KubernetesClusterId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.KubernetesClusterId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.KubernetesClusterId = types.StringValue(string(jsonBytes))
            } else {
                data.KubernetesClusterId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.KubernetesClusterId = types.StringValue(string(jsonBytes))
            } else {
                data.KubernetesClusterId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.KubernetesClusterId = types.StringValue(string(jsonBytes))
        } else {
            data.KubernetesClusterId = types.StringNull()
        }
    } else if val, ok := dataMap["kubernetesClusterId"].(string); ok && val != "" {
        data.KubernetesClusterId = types.StringValue(val)
    } else {
        data.KubernetesClusterId = types.StringNull()
    }
    if obj, ok := dataMap["clusterName"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ClusterName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ClusterName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ClusterName = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ClusterName = types.StringValue(string(jsonBytes))
            } else {
                data.ClusterName = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ClusterName = types.StringValue(string(jsonBytes))
            } else {
                data.ClusterName = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ClusterName = types.StringValue(string(jsonBytes))
        } else {
            data.ClusterName = types.StringNull()
        }
    } else if val, ok := dataMap["clusterName"].(string); ok && val != "" {
        data.ClusterName = types.StringValue(val)
    } else {
        data.ClusterName = types.StringNull()
    }
    if obj, ok := dataMap["k8sClusterEntityKey"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.K8sClusterEntityKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.K8sClusterEntityKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.K8sClusterEntityKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.K8sClusterEntityKey = types.StringValue(string(jsonBytes))
            } else {
                data.K8sClusterEntityKey = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.K8sClusterEntityKey = types.StringValue(string(jsonBytes))
            } else {
                data.K8sClusterEntityKey = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.K8sClusterEntityKey = types.StringValue(string(jsonBytes))
        } else {
            data.K8sClusterEntityKey = types.StringNull()
        }
    } else if val, ok := dataMap["k8sClusterEntityKey"].(string); ok && val != "" {
        data.K8sClusterEntityKey = types.StringValue(val)
    } else {
        data.K8sClusterEntityKey = types.StringNull()
    }
    if obj, ok := dataMap["windowStart"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.WindowStart = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.WindowStart = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.WindowStart = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.WindowStart = types.StringValue(string(jsonBytes))
            } else {
                data.WindowStart = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.WindowStart = types.StringValue(string(jsonBytes))
            } else {
                data.WindowStart = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.WindowStart = types.StringValue(string(jsonBytes))
        } else {
            data.WindowStart = types.StringNull()
        }
    } else if val, ok := dataMap["windowStart"].(string); ok && val != "" {
        data.WindowStart = types.StringValue(val)
    } else {
        data.WindowStart = types.StringNull()
    }
    if obj, ok := dataMap["windowEnd"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.WindowEnd = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.WindowEnd = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.WindowEnd = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.WindowEnd = types.StringValue(string(jsonBytes))
            } else {
                data.WindowEnd = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.WindowEnd = types.StringValue(string(jsonBytes))
            } else {
                data.WindowEnd = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.WindowEnd = types.StringValue(string(jsonBytes))
        } else {
            data.WindowEnd = types.StringNull()
        }
    } else if val, ok := dataMap["windowEnd"].(string); ok && val != "" {
        data.WindowEnd = types.StringValue(val)
    } else {
        data.WindowEnd = types.StringNull()
    }
    if obj, ok := dataMap["namespace"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Namespace = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Namespace = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Namespace = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Namespace = types.StringValue(string(jsonBytes))
            } else {
                data.Namespace = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Namespace = types.StringValue(string(jsonBytes))
            } else {
                data.Namespace = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Namespace = types.StringValue(string(jsonBytes))
        } else {
            data.Namespace = types.StringNull()
        }
    } else if val, ok := dataMap["namespace"].(string); ok && val != "" {
        data.Namespace = types.StringValue(val)
    } else {
        data.Namespace = types.StringNull()
    }
    if obj, ok := dataMap["controllerKind"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ControllerKind = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ControllerKind = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ControllerKind = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ControllerKind = types.StringValue(string(jsonBytes))
            } else {
                data.ControllerKind = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ControllerKind = types.StringValue(string(jsonBytes))
            } else {
                data.ControllerKind = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ControllerKind = types.StringValue(string(jsonBytes))
        } else {
            data.ControllerKind = types.StringNull()
        }
    } else if val, ok := dataMap["controllerKind"].(string); ok && val != "" {
        data.ControllerKind = types.StringValue(val)
    } else {
        data.ControllerKind = types.StringNull()
    }
    if obj, ok := dataMap["controllerName"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ControllerName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ControllerName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ControllerName = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ControllerName = types.StringValue(string(jsonBytes))
            } else {
                data.ControllerName = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ControllerName = types.StringValue(string(jsonBytes))
            } else {
                data.ControllerName = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ControllerName = types.StringValue(string(jsonBytes))
        } else {
            data.ControllerName = types.StringNull()
        }
    } else if val, ok := dataMap["controllerName"].(string); ok && val != "" {
        data.ControllerName = types.StringValue(val)
    } else {
        data.ControllerName = types.StringNull()
    }
    if obj, ok := dataMap["podName"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PodName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.PodName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.PodName = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.PodName = types.StringValue(string(jsonBytes))
            } else {
                data.PodName = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.PodName = types.StringValue(string(jsonBytes))
            } else {
                data.PodName = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.PodName = types.StringValue(string(jsonBytes))
        } else {
            data.PodName = types.StringNull()
        }
    } else if val, ok := dataMap["podName"].(string); ok && val != "" {
        data.PodName = types.StringValue(val)
    } else {
        data.PodName = types.StringNull()
    }
    if obj, ok := dataMap["containerName"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ContainerName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ContainerName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ContainerName = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ContainerName = types.StringValue(string(jsonBytes))
            } else {
                data.ContainerName = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ContainerName = types.StringValue(string(jsonBytes))
            } else {
                data.ContainerName = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ContainerName = types.StringValue(string(jsonBytes))
        } else {
            data.ContainerName = types.StringNull()
        }
    } else if val, ok := dataMap["containerName"].(string); ok && val != "" {
        data.ContainerName = types.StringValue(val)
    } else {
        data.ContainerName = types.StringNull()
    }
    if obj, ok := dataMap["nodeName"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.NodeName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.NodeName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.NodeName = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.NodeName = types.StringValue(string(jsonBytes))
            } else {
                data.NodeName = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.NodeName = types.StringValue(string(jsonBytes))
            } else {
                data.NodeName = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.NodeName = types.StringValue(string(jsonBytes))
        } else {
            data.NodeName = types.StringNull()
        }
    } else if val, ok := dataMap["nodeName"].(string); ok && val != "" {
        data.NodeName = types.StringValue(val)
    } else {
        data.NodeName = types.StringNull()
    }
    if obj, ok := dataMap["providerId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ProviderId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ProviderId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ProviderId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ProviderId = types.StringValue(string(jsonBytes))
            } else {
                data.ProviderId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ProviderId = types.StringValue(string(jsonBytes))
            } else {
                data.ProviderId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ProviderId = types.StringValue(string(jsonBytes))
        } else {
            data.ProviderId = types.StringNull()
        }
    } else if val, ok := dataMap["providerId"].(string); ok && val != "" {
        data.ProviderId = types.StringValue(val)
    } else {
        data.ProviderId = types.StringNull()
    }
    if obj, ok := dataMap["labels"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Labels = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Labels = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Labels = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Labels = types.StringValue(string(jsonBytes))
            } else {
                data.Labels = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Labels = types.StringValue(string(jsonBytes))
            } else {
                data.Labels = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Labels = types.StringValue(string(jsonBytes))
        } else {
            data.Labels = types.StringNull()
        }
    } else if val, ok := dataMap["labels"].(string); ok && val != "" {
        data.Labels = types.StringValue(val)
    } else {
        data.Labels = types.StringNull()
    }
    if val, ok := dataMap["labelKeys"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.LabelKeys = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.LabelKeys = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if val, ok := dataMap["cpuCoreHours"].(float64); ok {
        data.CpuCoreHours = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["cpuCoreHours"].(int); ok {
        data.CpuCoreHours = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["cpuCoreHours"].(int64); ok {
        data.CpuCoreHours = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["cpuCoreHours"] == nil {
        data.CpuCoreHours = types.NumberNull()
    }
    if val, ok := dataMap["cpuCoreRequestAverage"].(float64); ok {
        data.CpuCoreRequestAverage = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["cpuCoreRequestAverage"].(int); ok {
        data.CpuCoreRequestAverage = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["cpuCoreRequestAverage"].(int64); ok {
        data.CpuCoreRequestAverage = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["cpuCoreRequestAverage"] == nil {
        data.CpuCoreRequestAverage = types.NumberNull()
    }
    if val, ok := dataMap["cpuCoreUsageAverage"].(float64); ok {
        data.CpuCoreUsageAverage = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["cpuCoreUsageAverage"].(int); ok {
        data.CpuCoreUsageAverage = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["cpuCoreUsageAverage"].(int64); ok {
        data.CpuCoreUsageAverage = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["cpuCoreUsageAverage"] == nil {
        data.CpuCoreUsageAverage = types.NumberNull()
    }
    if val, ok := dataMap["cpuCoreLimitAverage"].(float64); ok {
        data.CpuCoreLimitAverage = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["cpuCoreLimitAverage"].(int); ok {
        data.CpuCoreLimitAverage = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["cpuCoreLimitAverage"].(int64); ok {
        data.CpuCoreLimitAverage = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["cpuCoreLimitAverage"] == nil {
        data.CpuCoreLimitAverage = types.NumberNull()
    }
    if val, ok := dataMap["cpuCost"].(float64); ok {
        data.CpuCost = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["cpuCost"].(int); ok {
        data.CpuCost = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["cpuCost"].(int64); ok {
        data.CpuCost = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["cpuCost"] == nil {
        data.CpuCost = types.NumberNull()
    }
    if val, ok := dataMap["gpuHours"].(float64); ok {
        data.GpuHours = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["gpuHours"].(int); ok {
        data.GpuHours = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["gpuHours"].(int64); ok {
        data.GpuHours = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["gpuHours"] == nil {
        data.GpuHours = types.NumberNull()
    }
    if val, ok := dataMap["gpuCost"].(float64); ok {
        data.GpuCost = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["gpuCost"].(int); ok {
        data.GpuCost = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["gpuCost"].(int64); ok {
        data.GpuCost = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["gpuCost"] == nil {
        data.GpuCost = types.NumberNull()
    }
    if val, ok := dataMap["ramByteHours"].(float64); ok {
        data.RamByteHours = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["ramByteHours"].(int); ok {
        data.RamByteHours = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["ramByteHours"].(int64); ok {
        data.RamByteHours = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["ramByteHours"] == nil {
        data.RamByteHours = types.NumberNull()
    }
    if val, ok := dataMap["ramBytesRequestAverage"].(float64); ok {
        data.RamBytesRequestAverage = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["ramBytesRequestAverage"].(int); ok {
        data.RamBytesRequestAverage = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["ramBytesRequestAverage"].(int64); ok {
        data.RamBytesRequestAverage = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["ramBytesRequestAverage"] == nil {
        data.RamBytesRequestAverage = types.NumberNull()
    }
    if val, ok := dataMap["ramBytesUsageAverage"].(float64); ok {
        data.RamBytesUsageAverage = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["ramBytesUsageAverage"].(int); ok {
        data.RamBytesUsageAverage = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["ramBytesUsageAverage"].(int64); ok {
        data.RamBytesUsageAverage = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["ramBytesUsageAverage"] == nil {
        data.RamBytesUsageAverage = types.NumberNull()
    }
    if val, ok := dataMap["ramBytesLimitAverage"].(float64); ok {
        data.RamBytesLimitAverage = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["ramBytesLimitAverage"].(int); ok {
        data.RamBytesLimitAverage = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["ramBytesLimitAverage"].(int64); ok {
        data.RamBytesLimitAverage = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["ramBytesLimitAverage"] == nil {
        data.RamBytesLimitAverage = types.NumberNull()
    }
    if val, ok := dataMap["ramBytesUsageMax"].(float64); ok {
        data.RamBytesUsageMax = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["ramBytesUsageMax"].(int); ok {
        data.RamBytesUsageMax = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["ramBytesUsageMax"].(int64); ok {
        data.RamBytesUsageMax = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["ramBytesUsageMax"] == nil {
        data.RamBytesUsageMax = types.NumberNull()
    }
    if val, ok := dataMap["ramCost"].(float64); ok {
        data.RamCost = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["ramCost"].(int); ok {
        data.RamCost = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["ramCost"].(int64); ok {
        data.RamCost = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["ramCost"] == nil {
        data.RamCost = types.NumberNull()
    }
    if val, ok := dataMap["pvByteHours"].(float64); ok {
        data.PvByteHours = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["pvByteHours"].(int); ok {
        data.PvByteHours = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["pvByteHours"].(int64); ok {
        data.PvByteHours = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["pvByteHours"] == nil {
        data.PvByteHours = types.NumberNull()
    }
    if val, ok := dataMap["pvCost"].(float64); ok {
        data.PvCost = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["pvCost"].(int); ok {
        data.PvCost = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["pvCost"].(int64); ok {
        data.PvCost = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["pvCost"] == nil {
        data.PvCost = types.NumberNull()
    }
    if val, ok := dataMap["networkCost"].(float64); ok {
        data.NetworkCost = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["networkCost"].(int); ok {
        data.NetworkCost = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["networkCost"].(int64); ok {
        data.NetworkCost = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["networkCost"] == nil {
        data.NetworkCost = types.NumberNull()
    }
    if val, ok := dataMap["loadBalancerCost"].(float64); ok {
        data.LoadBalancerCost = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["loadBalancerCost"].(int); ok {
        data.LoadBalancerCost = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["loadBalancerCost"].(int64); ok {
        data.LoadBalancerCost = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["loadBalancerCost"] == nil {
        data.LoadBalancerCost = types.NumberNull()
    }
    if val, ok := dataMap["sharedCost"].(float64); ok {
        data.SharedCost = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["sharedCost"].(int); ok {
        data.SharedCost = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["sharedCost"].(int64); ok {
        data.SharedCost = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["sharedCost"] == nil {
        data.SharedCost = types.NumberNull()
    }
    if val, ok := dataMap["externalCost"].(float64); ok {
        data.ExternalCost = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["externalCost"].(int); ok {
        data.ExternalCost = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["externalCost"].(int64); ok {
        data.ExternalCost = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["externalCost"] == nil {
        data.ExternalCost = types.NumberNull()
    }
    if val, ok := dataMap["totalCost"].(float64); ok {
        data.TotalCost = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["totalCost"].(int); ok {
        data.TotalCost = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["totalCost"].(int64); ok {
        data.TotalCost = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["totalCost"] == nil {
        data.TotalCost = types.NumberNull()
    }
    if val, ok := dataMap["cpuEfficiency"].(float64); ok {
        data.CpuEfficiency = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["cpuEfficiency"].(int); ok {
        data.CpuEfficiency = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["cpuEfficiency"].(int64); ok {
        data.CpuEfficiency = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["cpuEfficiency"] == nil {
        data.CpuEfficiency = types.NumberNull()
    }
    if val, ok := dataMap["ramEfficiency"].(float64); ok {
        data.RamEfficiency = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["ramEfficiency"].(int); ok {
        data.RamEfficiency = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["ramEfficiency"].(int64); ok {
        data.RamEfficiency = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["ramEfficiency"] == nil {
        data.RamEfficiency = types.NumberNull()
    }
    if val, ok := dataMap["totalEfficiency"].(float64); ok {
        data.TotalEfficiency = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["totalEfficiency"].(int); ok {
        data.TotalEfficiency = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["totalEfficiency"].(int64); ok {
        data.TotalEfficiency = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["totalEfficiency"] == nil {
        data.TotalEfficiency = types.NumberNull()
    }
    if obj, ok := dataMap["currency"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Currency = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Currency = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Currency = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Currency = types.StringValue(string(jsonBytes))
            } else {
                data.Currency = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Currency = types.StringValue(string(jsonBytes))
            } else {
                data.Currency = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Currency = types.StringValue(string(jsonBytes))
        } else {
            data.Currency = types.StringNull()
        }
    } else if val, ok := dataMap["currency"].(string); ok && val != "" {
        data.Currency = types.StringValue(val)
    } else {
        data.Currency = types.StringNull()
    }
    if obj, ok := dataMap["shipmentId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ShipmentId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ShipmentId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ShipmentId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ShipmentId = types.StringValue(string(jsonBytes))
            } else {
                data.ShipmentId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ShipmentId = types.StringValue(string(jsonBytes))
            } else {
                data.ShipmentId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ShipmentId = types.StringValue(string(jsonBytes))
        } else {
            data.ShipmentId = types.StringNull()
        }
    } else if val, ok := dataMap["shipmentId"].(string); ok && val != "" {
        data.ShipmentId = types.StringValue(val)
    } else {
        data.ShipmentId = types.StringNull()
    }
    if val, ok := dataMap["shipmentChunk"].(float64); ok {
        data.ShipmentChunk = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["shipmentChunk"].(int); ok {
        data.ShipmentChunk = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["shipmentChunk"].(int64); ok {
        data.ShipmentChunk = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["shipmentChunk"] == nil {
        data.ShipmentChunk = types.NumberNull()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "created a resource")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *KubernetesCostAllocationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var data KubernetesCostAllocationResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create select parameter to get full object
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

    // Make API call with select parameter
    httpResp, err := r.client.PostWithSelect("/kubernetes-cost-allocation/" + data.Id.ValueString() + "", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read kubernetes_cost_allocation, got error: %s", err))
        return
    }

    if httpResp.StatusCode == http.StatusNotFound {
        resp.State.RemoveResource(ctx)
        return
    }

    var kubernetesCostAllocationResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &kubernetesCostAllocationResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse kubernetes_cost_allocation response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := kubernetesCostAllocationResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = kubernetesCostAllocationResponse
    }

    if obj, ok := dataMap["id"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Id = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Id = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Id = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Id = types.StringValue(string(jsonBytes))
            } else {
                data.Id = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Id = types.StringValue(string(jsonBytes))
            } else {
                data.Id = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Id = types.StringValue(string(jsonBytes))
        } else {
            data.Id = types.StringNull()
        }
    } else if val, ok := dataMap["id"].(string); ok && val != "" {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    if obj, ok := dataMap["projectId"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok {
            data.ProjectId = types.StringValue(val)
        } else {
            data.ProjectId = types.StringNull()
        }
    } else if val, ok := dataMap["projectId"].(string); ok {
        data.ProjectId = types.StringValue(val)
    } else {
        data.ProjectId = types.StringNull()
    }
    if obj, ok := dataMap["kubernetesClusterId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.KubernetesClusterId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.KubernetesClusterId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.KubernetesClusterId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.KubernetesClusterId = types.StringValue(string(jsonBytes))
            } else {
                data.KubernetesClusterId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.KubernetesClusterId = types.StringValue(string(jsonBytes))
            } else {
                data.KubernetesClusterId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.KubernetesClusterId = types.StringValue(string(jsonBytes))
        } else {
            data.KubernetesClusterId = types.StringNull()
        }
    } else if val, ok := dataMap["kubernetesClusterId"].(string); ok && val != "" {
        data.KubernetesClusterId = types.StringValue(val)
    } else {
        data.KubernetesClusterId = types.StringNull()
    }
    if obj, ok := dataMap["clusterName"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ClusterName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ClusterName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ClusterName = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ClusterName = types.StringValue(string(jsonBytes))
            } else {
                data.ClusterName = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ClusterName = types.StringValue(string(jsonBytes))
            } else {
                data.ClusterName = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ClusterName = types.StringValue(string(jsonBytes))
        } else {
            data.ClusterName = types.StringNull()
        }
    } else if val, ok := dataMap["clusterName"].(string); ok && val != "" {
        data.ClusterName = types.StringValue(val)
    } else {
        data.ClusterName = types.StringNull()
    }
    if obj, ok := dataMap["k8sClusterEntityKey"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.K8sClusterEntityKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.K8sClusterEntityKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.K8sClusterEntityKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.K8sClusterEntityKey = types.StringValue(string(jsonBytes))
            } else {
                data.K8sClusterEntityKey = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.K8sClusterEntityKey = types.StringValue(string(jsonBytes))
            } else {
                data.K8sClusterEntityKey = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.K8sClusterEntityKey = types.StringValue(string(jsonBytes))
        } else {
            data.K8sClusterEntityKey = types.StringNull()
        }
    } else if val, ok := dataMap["k8sClusterEntityKey"].(string); ok && val != "" {
        data.K8sClusterEntityKey = types.StringValue(val)
    } else {
        data.K8sClusterEntityKey = types.StringNull()
    }
    if obj, ok := dataMap["windowStart"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.WindowStart = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.WindowStart = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.WindowStart = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.WindowStart = types.StringValue(string(jsonBytes))
            } else {
                data.WindowStart = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.WindowStart = types.StringValue(string(jsonBytes))
            } else {
                data.WindowStart = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.WindowStart = types.StringValue(string(jsonBytes))
        } else {
            data.WindowStart = types.StringNull()
        }
    } else if val, ok := dataMap["windowStart"].(string); ok && val != "" {
        data.WindowStart = types.StringValue(val)
    } else {
        data.WindowStart = types.StringNull()
    }
    if obj, ok := dataMap["windowEnd"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.WindowEnd = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.WindowEnd = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.WindowEnd = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.WindowEnd = types.StringValue(string(jsonBytes))
            } else {
                data.WindowEnd = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.WindowEnd = types.StringValue(string(jsonBytes))
            } else {
                data.WindowEnd = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.WindowEnd = types.StringValue(string(jsonBytes))
        } else {
            data.WindowEnd = types.StringNull()
        }
    } else if val, ok := dataMap["windowEnd"].(string); ok && val != "" {
        data.WindowEnd = types.StringValue(val)
    } else {
        data.WindowEnd = types.StringNull()
    }
    if obj, ok := dataMap["namespace"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Namespace = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Namespace = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Namespace = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Namespace = types.StringValue(string(jsonBytes))
            } else {
                data.Namespace = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Namespace = types.StringValue(string(jsonBytes))
            } else {
                data.Namespace = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Namespace = types.StringValue(string(jsonBytes))
        } else {
            data.Namespace = types.StringNull()
        }
    } else if val, ok := dataMap["namespace"].(string); ok && val != "" {
        data.Namespace = types.StringValue(val)
    } else {
        data.Namespace = types.StringNull()
    }
    if obj, ok := dataMap["controllerKind"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ControllerKind = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ControllerKind = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ControllerKind = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ControllerKind = types.StringValue(string(jsonBytes))
            } else {
                data.ControllerKind = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ControllerKind = types.StringValue(string(jsonBytes))
            } else {
                data.ControllerKind = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ControllerKind = types.StringValue(string(jsonBytes))
        } else {
            data.ControllerKind = types.StringNull()
        }
    } else if val, ok := dataMap["controllerKind"].(string); ok && val != "" {
        data.ControllerKind = types.StringValue(val)
    } else {
        data.ControllerKind = types.StringNull()
    }
    if obj, ok := dataMap["controllerName"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ControllerName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ControllerName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ControllerName = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ControllerName = types.StringValue(string(jsonBytes))
            } else {
                data.ControllerName = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ControllerName = types.StringValue(string(jsonBytes))
            } else {
                data.ControllerName = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ControllerName = types.StringValue(string(jsonBytes))
        } else {
            data.ControllerName = types.StringNull()
        }
    } else if val, ok := dataMap["controllerName"].(string); ok && val != "" {
        data.ControllerName = types.StringValue(val)
    } else {
        data.ControllerName = types.StringNull()
    }
    if obj, ok := dataMap["podName"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PodName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.PodName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.PodName = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.PodName = types.StringValue(string(jsonBytes))
            } else {
                data.PodName = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.PodName = types.StringValue(string(jsonBytes))
            } else {
                data.PodName = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.PodName = types.StringValue(string(jsonBytes))
        } else {
            data.PodName = types.StringNull()
        }
    } else if val, ok := dataMap["podName"].(string); ok && val != "" {
        data.PodName = types.StringValue(val)
    } else {
        data.PodName = types.StringNull()
    }
    if obj, ok := dataMap["containerName"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ContainerName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ContainerName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ContainerName = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ContainerName = types.StringValue(string(jsonBytes))
            } else {
                data.ContainerName = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ContainerName = types.StringValue(string(jsonBytes))
            } else {
                data.ContainerName = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ContainerName = types.StringValue(string(jsonBytes))
        } else {
            data.ContainerName = types.StringNull()
        }
    } else if val, ok := dataMap["containerName"].(string); ok && val != "" {
        data.ContainerName = types.StringValue(val)
    } else {
        data.ContainerName = types.StringNull()
    }
    if obj, ok := dataMap["nodeName"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.NodeName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.NodeName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.NodeName = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.NodeName = types.StringValue(string(jsonBytes))
            } else {
                data.NodeName = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.NodeName = types.StringValue(string(jsonBytes))
            } else {
                data.NodeName = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.NodeName = types.StringValue(string(jsonBytes))
        } else {
            data.NodeName = types.StringNull()
        }
    } else if val, ok := dataMap["nodeName"].(string); ok && val != "" {
        data.NodeName = types.StringValue(val)
    } else {
        data.NodeName = types.StringNull()
    }
    if obj, ok := dataMap["providerId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ProviderId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ProviderId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ProviderId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ProviderId = types.StringValue(string(jsonBytes))
            } else {
                data.ProviderId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ProviderId = types.StringValue(string(jsonBytes))
            } else {
                data.ProviderId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ProviderId = types.StringValue(string(jsonBytes))
        } else {
            data.ProviderId = types.StringNull()
        }
    } else if val, ok := dataMap["providerId"].(string); ok && val != "" {
        data.ProviderId = types.StringValue(val)
    } else {
        data.ProviderId = types.StringNull()
    }
    if obj, ok := dataMap["labels"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Labels = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Labels = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Labels = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Labels = types.StringValue(string(jsonBytes))
            } else {
                data.Labels = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Labels = types.StringValue(string(jsonBytes))
            } else {
                data.Labels = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Labels = types.StringValue(string(jsonBytes))
        } else {
            data.Labels = types.StringNull()
        }
    } else if val, ok := dataMap["labels"].(string); ok && val != "" {
        data.Labels = types.StringValue(val)
    } else {
        data.Labels = types.StringNull()
    }
    if val, ok := dataMap["labelKeys"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.LabelKeys = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.LabelKeys = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if val, ok := dataMap["cpuCoreHours"].(float64); ok {
        data.CpuCoreHours = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["cpuCoreHours"].(int); ok {
        data.CpuCoreHours = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["cpuCoreHours"].(int64); ok {
        data.CpuCoreHours = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["cpuCoreHours"] == nil {
        data.CpuCoreHours = types.NumberNull()
    }
    if val, ok := dataMap["cpuCoreRequestAverage"].(float64); ok {
        data.CpuCoreRequestAverage = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["cpuCoreRequestAverage"].(int); ok {
        data.CpuCoreRequestAverage = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["cpuCoreRequestAverage"].(int64); ok {
        data.CpuCoreRequestAverage = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["cpuCoreRequestAverage"] == nil {
        data.CpuCoreRequestAverage = types.NumberNull()
    }
    if val, ok := dataMap["cpuCoreUsageAverage"].(float64); ok {
        data.CpuCoreUsageAverage = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["cpuCoreUsageAverage"].(int); ok {
        data.CpuCoreUsageAverage = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["cpuCoreUsageAverage"].(int64); ok {
        data.CpuCoreUsageAverage = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["cpuCoreUsageAverage"] == nil {
        data.CpuCoreUsageAverage = types.NumberNull()
    }
    if val, ok := dataMap["cpuCoreLimitAverage"].(float64); ok {
        data.CpuCoreLimitAverage = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["cpuCoreLimitAverage"].(int); ok {
        data.CpuCoreLimitAverage = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["cpuCoreLimitAverage"].(int64); ok {
        data.CpuCoreLimitAverage = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["cpuCoreLimitAverage"] == nil {
        data.CpuCoreLimitAverage = types.NumberNull()
    }
    if val, ok := dataMap["cpuCost"].(float64); ok {
        data.CpuCost = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["cpuCost"].(int); ok {
        data.CpuCost = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["cpuCost"].(int64); ok {
        data.CpuCost = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["cpuCost"] == nil {
        data.CpuCost = types.NumberNull()
    }
    if val, ok := dataMap["gpuHours"].(float64); ok {
        data.GpuHours = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["gpuHours"].(int); ok {
        data.GpuHours = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["gpuHours"].(int64); ok {
        data.GpuHours = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["gpuHours"] == nil {
        data.GpuHours = types.NumberNull()
    }
    if val, ok := dataMap["gpuCost"].(float64); ok {
        data.GpuCost = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["gpuCost"].(int); ok {
        data.GpuCost = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["gpuCost"].(int64); ok {
        data.GpuCost = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["gpuCost"] == nil {
        data.GpuCost = types.NumberNull()
    }
    if val, ok := dataMap["ramByteHours"].(float64); ok {
        data.RamByteHours = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["ramByteHours"].(int); ok {
        data.RamByteHours = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["ramByteHours"].(int64); ok {
        data.RamByteHours = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["ramByteHours"] == nil {
        data.RamByteHours = types.NumberNull()
    }
    if val, ok := dataMap["ramBytesRequestAverage"].(float64); ok {
        data.RamBytesRequestAverage = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["ramBytesRequestAverage"].(int); ok {
        data.RamBytesRequestAverage = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["ramBytesRequestAverage"].(int64); ok {
        data.RamBytesRequestAverage = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["ramBytesRequestAverage"] == nil {
        data.RamBytesRequestAverage = types.NumberNull()
    }
    if val, ok := dataMap["ramBytesUsageAverage"].(float64); ok {
        data.RamBytesUsageAverage = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["ramBytesUsageAverage"].(int); ok {
        data.RamBytesUsageAverage = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["ramBytesUsageAverage"].(int64); ok {
        data.RamBytesUsageAverage = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["ramBytesUsageAverage"] == nil {
        data.RamBytesUsageAverage = types.NumberNull()
    }
    if val, ok := dataMap["ramBytesLimitAverage"].(float64); ok {
        data.RamBytesLimitAverage = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["ramBytesLimitAverage"].(int); ok {
        data.RamBytesLimitAverage = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["ramBytesLimitAverage"].(int64); ok {
        data.RamBytesLimitAverage = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["ramBytesLimitAverage"] == nil {
        data.RamBytesLimitAverage = types.NumberNull()
    }
    if val, ok := dataMap["ramBytesUsageMax"].(float64); ok {
        data.RamBytesUsageMax = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["ramBytesUsageMax"].(int); ok {
        data.RamBytesUsageMax = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["ramBytesUsageMax"].(int64); ok {
        data.RamBytesUsageMax = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["ramBytesUsageMax"] == nil {
        data.RamBytesUsageMax = types.NumberNull()
    }
    if val, ok := dataMap["ramCost"].(float64); ok {
        data.RamCost = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["ramCost"].(int); ok {
        data.RamCost = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["ramCost"].(int64); ok {
        data.RamCost = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["ramCost"] == nil {
        data.RamCost = types.NumberNull()
    }
    if val, ok := dataMap["pvByteHours"].(float64); ok {
        data.PvByteHours = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["pvByteHours"].(int); ok {
        data.PvByteHours = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["pvByteHours"].(int64); ok {
        data.PvByteHours = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["pvByteHours"] == nil {
        data.PvByteHours = types.NumberNull()
    }
    if val, ok := dataMap["pvCost"].(float64); ok {
        data.PvCost = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["pvCost"].(int); ok {
        data.PvCost = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["pvCost"].(int64); ok {
        data.PvCost = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["pvCost"] == nil {
        data.PvCost = types.NumberNull()
    }
    if val, ok := dataMap["networkCost"].(float64); ok {
        data.NetworkCost = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["networkCost"].(int); ok {
        data.NetworkCost = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["networkCost"].(int64); ok {
        data.NetworkCost = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["networkCost"] == nil {
        data.NetworkCost = types.NumberNull()
    }
    if val, ok := dataMap["loadBalancerCost"].(float64); ok {
        data.LoadBalancerCost = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["loadBalancerCost"].(int); ok {
        data.LoadBalancerCost = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["loadBalancerCost"].(int64); ok {
        data.LoadBalancerCost = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["loadBalancerCost"] == nil {
        data.LoadBalancerCost = types.NumberNull()
    }
    if val, ok := dataMap["sharedCost"].(float64); ok {
        data.SharedCost = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["sharedCost"].(int); ok {
        data.SharedCost = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["sharedCost"].(int64); ok {
        data.SharedCost = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["sharedCost"] == nil {
        data.SharedCost = types.NumberNull()
    }
    if val, ok := dataMap["externalCost"].(float64); ok {
        data.ExternalCost = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["externalCost"].(int); ok {
        data.ExternalCost = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["externalCost"].(int64); ok {
        data.ExternalCost = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["externalCost"] == nil {
        data.ExternalCost = types.NumberNull()
    }
    if val, ok := dataMap["totalCost"].(float64); ok {
        data.TotalCost = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["totalCost"].(int); ok {
        data.TotalCost = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["totalCost"].(int64); ok {
        data.TotalCost = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["totalCost"] == nil {
        data.TotalCost = types.NumberNull()
    }
    if val, ok := dataMap["cpuEfficiency"].(float64); ok {
        data.CpuEfficiency = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["cpuEfficiency"].(int); ok {
        data.CpuEfficiency = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["cpuEfficiency"].(int64); ok {
        data.CpuEfficiency = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["cpuEfficiency"] == nil {
        data.CpuEfficiency = types.NumberNull()
    }
    if val, ok := dataMap["ramEfficiency"].(float64); ok {
        data.RamEfficiency = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["ramEfficiency"].(int); ok {
        data.RamEfficiency = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["ramEfficiency"].(int64); ok {
        data.RamEfficiency = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["ramEfficiency"] == nil {
        data.RamEfficiency = types.NumberNull()
    }
    if val, ok := dataMap["totalEfficiency"].(float64); ok {
        data.TotalEfficiency = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["totalEfficiency"].(int); ok {
        data.TotalEfficiency = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["totalEfficiency"].(int64); ok {
        data.TotalEfficiency = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["totalEfficiency"] == nil {
        data.TotalEfficiency = types.NumberNull()
    }
    if obj, ok := dataMap["currency"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Currency = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Currency = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Currency = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Currency = types.StringValue(string(jsonBytes))
            } else {
                data.Currency = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Currency = types.StringValue(string(jsonBytes))
            } else {
                data.Currency = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Currency = types.StringValue(string(jsonBytes))
        } else {
            data.Currency = types.StringNull()
        }
    } else if val, ok := dataMap["currency"].(string); ok && val != "" {
        data.Currency = types.StringValue(val)
    } else {
        data.Currency = types.StringNull()
    }
    if obj, ok := dataMap["shipmentId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ShipmentId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ShipmentId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ShipmentId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ShipmentId = types.StringValue(string(jsonBytes))
            } else {
                data.ShipmentId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ShipmentId = types.StringValue(string(jsonBytes))
            } else {
                data.ShipmentId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ShipmentId = types.StringValue(string(jsonBytes))
        } else {
            data.ShipmentId = types.StringNull()
        }
    } else if val, ok := dataMap["shipmentId"].(string); ok && val != "" {
        data.ShipmentId = types.StringValue(val)
    } else {
        data.ShipmentId = types.StringNull()
    }
    if val, ok := dataMap["shipmentChunk"].(float64); ok {
        data.ShipmentChunk = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["shipmentChunk"].(int); ok {
        data.ShipmentChunk = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["shipmentChunk"].(int64); ok {
        data.ShipmentChunk = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["shipmentChunk"] == nil {
        data.ShipmentChunk = types.NumberNull()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *KubernetesCostAllocationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var data KubernetesCostAllocationResourceModel

    // Read Terraform plan data into the model
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // This resource does not have an update API endpoint.
    // Preserve the planned state.
    tflog.Trace(ctx, "updated a resource (no-op: preserving planned state)")

    // Save planned data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *KubernetesCostAllocationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var data KubernetesCostAllocationResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Make API call
    _, err := r.client.Delete("/kubernetes-cost-allocation/" + data.Id.ValueString() + "")
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete kubernetes_cost_allocation, got error: %s", err))
        return
    }
}


func (r *KubernetesCostAllocationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Helper method to convert Terraform map to Go interface{}
func (r *KubernetesCostAllocationResource) convertTerraformMapToInterface(terraformMap types.Map) interface{} {
    if terraformMap.IsNull() || terraformMap.IsUnknown() {
        return nil
    }
    
    result := make(map[string]string)
    terraformMap.ElementsAs(context.Background(), &result, false)
    
    // Convert map[string]string to map[string]interface{}
    interfaceResult := make(map[string]interface{})
    for key, value := range result {
        interfaceResult[key] = value
    }
    
    return interfaceResult
}

// Helper method to convert Terraform list to Go interface{}
func (r *KubernetesCostAllocationResource) convertTerraformListToInterface(terraformList types.List) interface{} {
    if terraformList.IsNull() || terraformList.IsUnknown() {
        return nil
    }
    
    var stringList []string
    terraformList.ElementsAs(context.Background(), &stringList, false)
    
    // Convert string array to OneUptime format with _id fields
    var result []interface{}
    for _, str := range stringList {
        if str != "" {
            result = append(result, map[string]interface{}{
                "_id": str,
            })
        }
    }
    return result
}

// Helper method to convert Terraform set to Go interface{}
func (r *KubernetesCostAllocationResource) convertTerraformSetToInterface(terraformSet types.Set) interface{} {
    if terraformSet.IsNull() || terraformSet.IsUnknown() {
        return nil
    }
    
    var stringList []string
    terraformSet.ElementsAs(context.Background(), &stringList, false)
    
    // Convert string array to OneUptime format with _id fields
    var result []interface{}
    for _, str := range stringList {
        if str != "" {
            result = append(result, map[string]interface{}{
                "_id": str,
            })
        }
    }
    return result
}

// Helper method to parse JSON field for complex objects
func (r *KubernetesCostAllocationResource) parseJSONField(terraformString basetypes.StringValuable) interface{} {
    sv, _ := terraformString.ToStringValue(context.Background())
    if sv.IsNull() || sv.IsUnknown() || sv.ValueString() == "" {
        return nil
    }

    var result interface{}
    if err := json.Unmarshal([]byte(sv.ValueString()), &result); err != nil {
        // If JSON parsing fails, return the raw string
        return sv.ValueString()
    }

    return result
}

// Normalize URL wrapper objects to avoid drift (e.g., trailing slash differences).
func (r *KubernetesCostAllocationResource) normalizeURLWrappers(value interface{}) interface{} {
    switch v := value.(type) {
    case map[string]interface{}:
        if typeStr, ok := v["_type"].(string); ok && typeStr == "URL" {
            if val, ok := v["value"].(string); ok {
                v["value"] = r.normalizeURLString(val)
            }
        }
        for key, child := range v {
            v[key] = r.normalizeURLWrappers(child)
        }
        return v
    case []interface{}:
        for i, child := range v {
            v[i] = r.normalizeURLWrappers(child)
        }
        return v
    default:
        return v
    }
}

func (r *KubernetesCostAllocationResource) normalizeURLString(value string) string {
    parsed, err := url.Parse(value)
    if err != nil {
        return value
    }
    if parsed.Path == "/" && parsed.RawQuery == "" && parsed.Fragment == "" {
        return strings.TrimSuffix(value, "/")
    }
    return value
}

// Helper method to convert *big.Float to float64 for JSON serialization
func (r *KubernetesCostAllocationResource) bigFloatToFloat64(bf *big.Float) interface{} {
    if bf == nil {
        return nil
    }
    f, _ := bf.Float64()
    return f
}

// Helper method to check if a type string is a valid OneUptime ObjectType
// Only these types should be marshalled/unmarshalled as typed wrapper objects
// This list is dynamically generated from Common/Types/JSON.ts ObjectType enum
func (r *KubernetesCostAllocationResource) isValidOneUptimeObjectType(typeStr string) bool {
    validTypes := map[string]bool{
        "ObjectID": true,
        "Decimal": true,
        "Name": true,
        "EqualTo": true,
        "EqualToOrNull": true,
        "MonitorSteps": true,
        "MonitorStep": true,
        "Recurring": true,
        "RestrictionTimes": true,
        "MonitorCriteria": true,
        "PositiveNumber": true,
        "MonitorCriteriaInstance": true,
        "NotEqual": true,
        "Email": true,
        "Phone": true,
        "Color": true,
        "Domain": true,
        "Version": true,
        "IP": true,
        "Route": true,
        "URL": true,
        "Permission": true,
        "Search": true,
        "MultiSearch": true,
        "GreaterThan": true,
        "GreaterThanOrEqual": true,
        "GreaterThanOrNull": true,
        "LessThanOrNull": true,
        "LessThan": true,
        "LessThanOrEqual": true,
        "Port": true,
        "Hostname": true,
        "HashedString": true,
        "DateTime": true,
        "Buffer": true,
        "InBetween": true,
        "NotNull": true,
        "IsNull": true,
        "Includes": true,
        "IncludesAll": true,
        "IncludesNone": true,
        "StartsWith": true,
        "EndsWith": true,
        "NotContains": true,
        "DashboardComponent": true,
        "DashboardViewConfig": true,
    }
    return validTypes[typeStr]
}
