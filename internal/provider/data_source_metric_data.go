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
var _ datasource.DataSource = &MetricDataDataSource{}

func NewMetricDataDataSource() datasource.DataSource {
    return &MetricDataDataSource{}
}

// MetricDataDataSource defines the data source implementation.
type MetricDataDataSource struct {
    client *Client
}

// MetricDataDataSourceModel describes the data source data model.
type MetricDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    ProjectId types.String `tfsdk:"project_id"`
    PrimaryEntityId types.String `tfsdk:"primary_entity_id"`
    PrimaryEntityType types.String `tfsdk:"primary_entity_type"`
    AggregationTemporality types.String `tfsdk:"aggregation_temporality"`
    MetricPointType types.String `tfsdk:"metric_point_type"`
    Time types.String `tfsdk:"time"`
    StartTime types.String `tfsdk:"start_time"`
    TimeUnixNano types.String `tfsdk:"time_unix_nano"`
    StartTimeUnixNano types.String `tfsdk:"start_time_unix_nano"`
    Attributes types.String `tfsdk:"attributes"`
    AttributeKeys types.Set `tfsdk:"attribute_keys"`
    EntityKeys types.Set `tfsdk:"entity_keys"`
    ServiceEntityKey types.String `tfsdk:"service_entity_key"`
    HostEntityKey types.String `tfsdk:"host_entity_key"`
    K8sPodEntityKey types.String `tfsdk:"k8s_pod_entity_key"`
    K8sNodeEntityKey types.String `tfsdk:"k8s_node_entity_key"`
    K8sClusterEntityKey types.String `tfsdk:"k8s_cluster_entity_key"`
    ContainerEntityKey types.String `tfsdk:"container_entity_key"`
    IsMonotonic types.Bool `tfsdk:"is_monotonic"`
    CountValue types.String `tfsdk:"count_value"`
    Sum types.Number `tfsdk:"sum"`
    Value types.Number `tfsdk:"value"`
    Min types.Number `tfsdk:"min"`
    Max types.Number `tfsdk:"max"`
    BucketCounts types.String `tfsdk:"bucket_counts"`
    ExplicitBounds types.String `tfsdk:"explicit_bounds"`
    Scale types.Number `tfsdk:"scale"`
    ZeroCount types.String `tfsdk:"zero_count"`
    PositiveOffset types.Number `tfsdk:"positive_offset"`
    PositiveBucketCounts types.String `tfsdk:"positive_bucket_counts"`
    NegativeOffset types.Number `tfsdk:"negative_offset"`
    NegativeBucketCounts types.String `tfsdk:"negative_bucket_counts"`
    SummaryQuantiles types.String `tfsdk:"summary_quantiles"`
    SummaryValues types.String `tfsdk:"summary_values"`
    TraceId types.String `tfsdk:"trace_id"`
    SpanId types.String `tfsdk:"span_id"`
}

func (d *MetricDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_metric_data"
}

func (d *MetricDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "metric_data data source",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Identifier to filter by",
                Optional: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Name to filter by",
                Computed: true,
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "Project ID",
                Computed: true,
            },
            "primary_entity_id": schema.StringAttribute{
                MarkdownDescription: "Service ID",
                Computed: true,
            },
            "primary_entity_type": schema.StringAttribute{
                MarkdownDescription: "Service Type",
                Computed: true,
            },
            "aggregation_temporality": schema.StringAttribute{
                MarkdownDescription: "Aggregation Temporality",
                Computed: true,
            },
            "metric_point_type": schema.StringAttribute{
                MarkdownDescription: "Metric Point Type",
                Computed: true,
            },
            "time": schema.StringAttribute{
                MarkdownDescription: "Time",
                Computed: true,
            },
            "start_time": schema.StringAttribute{
                MarkdownDescription: "Start Time",
                Computed: true,
            },
            "time_unix_nano": schema.StringAttribute{
                MarkdownDescription: "Time (in Unix Nano)",
                Computed: true,
            },
            "start_time_unix_nano": schema.StringAttribute{
                MarkdownDescription: "Start Time (in Unix Nano)",
                Computed: true,
            },
            "attributes": schema.StringAttribute{
                MarkdownDescription: "Attributes",
                Computed: true,
            },
            "attribute_keys": schema.SetAttribute{
                MarkdownDescription: "Attribute Keys",
                Computed: true,
                ElementType: types.StringType,
            },
            "entity_keys": schema.SetAttribute{
                MarkdownDescription: "Entity Keys",
                Computed: true,
                ElementType: types.StringType,
            },
            "service_entity_key": schema.StringAttribute{
                MarkdownDescription: "Service Entity Key",
                Computed: true,
            },
            "host_entity_key": schema.StringAttribute{
                MarkdownDescription: "Host Entity Key",
                Computed: true,
            },
            "k8s_pod_entity_key": schema.StringAttribute{
                MarkdownDescription: "Kubernetes Pod Entity Key",
                Computed: true,
            },
            "k8s_node_entity_key": schema.StringAttribute{
                MarkdownDescription: "Kubernetes Node Entity Key",
                Computed: true,
            },
            "k8s_cluster_entity_key": schema.StringAttribute{
                MarkdownDescription: "Kubernetes Cluster Entity Key",
                Computed: true,
            },
            "container_entity_key": schema.StringAttribute{
                MarkdownDescription: "Container Entity Key",
                Computed: true,
            },
            "is_monotonic": schema.BoolAttribute{
                MarkdownDescription: "Is Monotonic",
                Computed: true,
            },
            "count_value": schema.StringAttribute{
                MarkdownDescription: "Count",
                Computed: true,
            },
            "sum": schema.NumberAttribute{
                MarkdownDescription: "Sum",
                Computed: true,
            },
            "value": schema.NumberAttribute{
                MarkdownDescription: "Value",
                Computed: true,
            },
            "min": schema.NumberAttribute{
                MarkdownDescription: "Min",
                Computed: true,
            },
            "max": schema.NumberAttribute{
                MarkdownDescription: "Max",
                Computed: true,
            },
            "bucket_counts": schema.StringAttribute{
                MarkdownDescription: "Bucket Counts",
                Computed: true,
            },
            "explicit_bounds": schema.StringAttribute{
                MarkdownDescription: "Explicit Bounds",
                Computed: true,
            },
            "scale": schema.NumberAttribute{
                MarkdownDescription: "Scale",
                Computed: true,
            },
            "zero_count": schema.StringAttribute{
                MarkdownDescription: "Zero Count",
                Computed: true,
            },
            "positive_offset": schema.NumberAttribute{
                MarkdownDescription: "Positive Bucket Offset",
                Computed: true,
            },
            "positive_bucket_counts": schema.StringAttribute{
                MarkdownDescription: "Positive Bucket Counts",
                Computed: true,
            },
            "negative_offset": schema.NumberAttribute{
                MarkdownDescription: "Negative Bucket Offset",
                Computed: true,
            },
            "negative_bucket_counts": schema.StringAttribute{
                MarkdownDescription: "Negative Bucket Counts",
                Computed: true,
            },
            "summary_quantiles": schema.StringAttribute{
                MarkdownDescription: "Summary Quantiles",
                Computed: true,
            },
            "summary_values": schema.StringAttribute{
                MarkdownDescription: "Summary Values",
                Computed: true,
            },
            "trace_id": schema.StringAttribute{
                MarkdownDescription: "Trace ID",
                Computed: true,
            },
            "span_id": schema.StringAttribute{
                MarkdownDescription: "Span ID",
                Computed: true,
            },
        },
    }
}

func (d *MetricDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *MetricDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data MetricDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "metrics" + "/" + data.Id.ValueString()
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read metric_data, got error: %s", err))
        return
    }

    var metricDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &metricDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse metric_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := metricDataResponse["data"].(map[string]interface{}); ok {
        metricDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := metricDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := metricDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := metricDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := metricDataResponse["primary_entity_id"].(string); ok {
        data.PrimaryEntityId = types.StringValue(val)
    }
    if val, ok := metricDataResponse["primary_entity_type"].(string); ok {
        data.PrimaryEntityType = types.StringValue(val)
    }
    if val, ok := metricDataResponse["aggregation_temporality"].(string); ok {
        data.AggregationTemporality = types.StringValue(val)
    }
    if val, ok := metricDataResponse["metric_point_type"].(string); ok {
        data.MetricPointType = types.StringValue(val)
    }
    if val, ok := metricDataResponse["time"].(string); ok {
        data.Time = types.StringValue(val)
    }
    if val, ok := metricDataResponse["start_time"].(string); ok {
        data.StartTime = types.StringValue(val)
    }
    if val, ok := metricDataResponse["time_unix_nano"].(string); ok {
        data.TimeUnixNano = types.StringValue(val)
    }
    if val, ok := metricDataResponse["start_time_unix_nano"].(string); ok {
        data.StartTimeUnixNano = types.StringValue(val)
    }
    if val, ok := metricDataResponse["attributes"].(string); ok {
        data.Attributes = types.StringValue(val)
    }
    if val, ok := metricDataResponse["attribute_keys"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.AttributeKeys = setValue
    }
    if val, ok := metricDataResponse["entity_keys"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.EntityKeys = setValue
    }
    if val, ok := metricDataResponse["service_entity_key"].(string); ok {
        data.ServiceEntityKey = types.StringValue(val)
    }
    if val, ok := metricDataResponse["host_entity_key"].(string); ok {
        data.HostEntityKey = types.StringValue(val)
    }
    if val, ok := metricDataResponse["k8s_pod_entity_key"].(string); ok {
        data.K8sPodEntityKey = types.StringValue(val)
    }
    if val, ok := metricDataResponse["k8s_node_entity_key"].(string); ok {
        data.K8sNodeEntityKey = types.StringValue(val)
    }
    if val, ok := metricDataResponse["k8s_cluster_entity_key"].(string); ok {
        data.K8sClusterEntityKey = types.StringValue(val)
    }
    if val, ok := metricDataResponse["container_entity_key"].(string); ok {
        data.ContainerEntityKey = types.StringValue(val)
    }
    if val, ok := metricDataResponse["is_monotonic"].(bool); ok {
        data.IsMonotonic = types.BoolValue(val)
    }
    if val, ok := metricDataResponse["count"].(string); ok {
        data.CountValue = types.StringValue(val)
    }
    if val, ok := metricDataResponse["sum"].(float64); ok {
        data.Sum = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := metricDataResponse["value"].(float64); ok {
        data.Value = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := metricDataResponse["min"].(float64); ok {
        data.Min = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := metricDataResponse["max"].(float64); ok {
        data.Max = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := metricDataResponse["bucket_counts"].(string); ok {
        data.BucketCounts = types.StringValue(val)
    }
    if val, ok := metricDataResponse["explicit_bounds"].(string); ok {
        data.ExplicitBounds = types.StringValue(val)
    }
    if val, ok := metricDataResponse["scale"].(float64); ok {
        data.Scale = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := metricDataResponse["zero_count"].(string); ok {
        data.ZeroCount = types.StringValue(val)
    }
    if val, ok := metricDataResponse["positive_offset"].(float64); ok {
        data.PositiveOffset = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := metricDataResponse["positive_bucket_counts"].(string); ok {
        data.PositiveBucketCounts = types.StringValue(val)
    }
    if val, ok := metricDataResponse["negative_offset"].(float64); ok {
        data.NegativeOffset = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := metricDataResponse["negative_bucket_counts"].(string); ok {
        data.NegativeBucketCounts = types.StringValue(val)
    }
    if val, ok := metricDataResponse["summary_quantiles"].(string); ok {
        data.SummaryQuantiles = types.StringValue(val)
    }
    if val, ok := metricDataResponse["summary_values"].(string); ok {
        data.SummaryValues = types.StringValue(val)
    }
    if val, ok := metricDataResponse["trace_id"].(string); ok {
        data.TraceId = types.StringValue(val)
    }
    if val, ok := metricDataResponse["span_id"].(string); ok {
        data.SpanId = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
