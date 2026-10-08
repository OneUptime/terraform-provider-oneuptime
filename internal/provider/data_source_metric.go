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
var _ datasource.DataSource = &MetricDataSource{}

func NewMetricDataSource() datasource.DataSource {
    return &MetricDataSource{}
}

// MetricDataSource defines the data source implementation.
type MetricDataSource struct {
    client *Client
}

// MetricDataSourceModel describes the data source data model.
type MetricDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    ProjectId types.String `tfsdk:"project_id"`
    PrimaryEntityId types.String `tfsdk:"primary_entity_id"`
    PrimaryEntityType types.String `tfsdk:"primary_entity_type"`
    Name types.String `tfsdk:"name"`
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

func (d *MetricDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_metric"
}

func (d *MetricDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "API endpoints for Metric Look up an existing metric by `id`, or by any of its other arguments (`name`, `aggregation_temporality`, `attributes`, ...): each one set must match, and exactly one metric may match them all.",

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
            "primary_entity_id": schema.StringAttribute{
                MarkdownDescription: "Service ID.",
                Optional: true,
                Computed: true,
            },
            "primary_entity_type": schema.StringAttribute{
                MarkdownDescription: "Service Type.",
                Optional: true,
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Name.",
                Optional: true,
                Computed: true,
            },
            "aggregation_temporality": schema.StringAttribute{
                MarkdownDescription: "Aggregation Temporality.",
                Optional: true,
                Computed: true,
            },
            "metric_point_type": schema.StringAttribute{
                MarkdownDescription: "Metric Point Type.",
                Optional: true,
                Computed: true,
            },
            "time": schema.StringAttribute{
                MarkdownDescription: "Time.",
                Optional: true,
                Computed: true,
            },
            "start_time": schema.StringAttribute{
                MarkdownDescription: "Start Time.",
                Optional: true,
                Computed: true,
            },
            "time_unix_nano": schema.StringAttribute{
                MarkdownDescription: "Time (in Unix Nano).",
                Optional: true,
                Computed: true,
            },
            "start_time_unix_nano": schema.StringAttribute{
                MarkdownDescription: "Start Time (in Unix Nano).",
                Optional: true,
                Computed: true,
            },
            "attributes": schema.StringAttribute{
                MarkdownDescription: "Attributes.",
                Optional: true,
                Computed: true,
            },
            "attribute_keys": schema.SetAttribute{
                MarkdownDescription: "Attribute Keys.",
                Computed: true,
                ElementType: types.StringType,
            },
            "entity_keys": schema.SetAttribute{
                MarkdownDescription: "Entity Keys.",
                Computed: true,
                ElementType: types.StringType,
            },
            "service_entity_key": schema.StringAttribute{
                MarkdownDescription: "Service Entity Key.",
                Optional: true,
                Computed: true,
            },
            "host_entity_key": schema.StringAttribute{
                MarkdownDescription: "Host Entity Key.",
                Optional: true,
                Computed: true,
            },
            "k8s_pod_entity_key": schema.StringAttribute{
                MarkdownDescription: "Kubernetes Pod Entity Key.",
                Optional: true,
                Computed: true,
            },
            "k8s_node_entity_key": schema.StringAttribute{
                MarkdownDescription: "Kubernetes Node Entity Key.",
                Optional: true,
                Computed: true,
            },
            "k8s_cluster_entity_key": schema.StringAttribute{
                MarkdownDescription: "Kubernetes Cluster Entity Key.",
                Optional: true,
                Computed: true,
            },
            "container_entity_key": schema.StringAttribute{
                MarkdownDescription: "Container Entity Key.",
                Optional: true,
                Computed: true,
            },
            "is_monotonic": schema.BoolAttribute{
                MarkdownDescription: "Is Monotonic.",
                Optional: true,
                Computed: true,
            },
            "count_value": schema.StringAttribute{
                MarkdownDescription: "Count.",
                Optional: true,
                Computed: true,
            },
            "sum": schema.NumberAttribute{
                MarkdownDescription: "Sum.",
                Optional: true,
                Computed: true,
            },
            "value": schema.NumberAttribute{
                MarkdownDescription: "Value.",
                Optional: true,
                Computed: true,
            },
            "min": schema.NumberAttribute{
                MarkdownDescription: "Min.",
                Optional: true,
                Computed: true,
            },
            "max": schema.NumberAttribute{
                MarkdownDescription: "Max.",
                Optional: true,
                Computed: true,
            },
            "bucket_counts": schema.StringAttribute{
                MarkdownDescription: "Bucket Counts.",
                Optional: true,
                Computed: true,
            },
            "explicit_bounds": schema.StringAttribute{
                MarkdownDescription: "Explicit Bounds.",
                Optional: true,
                Computed: true,
            },
            "scale": schema.NumberAttribute{
                MarkdownDescription: "Scale.",
                Optional: true,
                Computed: true,
            },
            "zero_count": schema.StringAttribute{
                MarkdownDescription: "Zero Count.",
                Optional: true,
                Computed: true,
            },
            "positive_offset": schema.NumberAttribute{
                MarkdownDescription: "Positive Bucket Offset.",
                Optional: true,
                Computed: true,
            },
            "positive_bucket_counts": schema.StringAttribute{
                MarkdownDescription: "Positive Bucket Counts.",
                Optional: true,
                Computed: true,
            },
            "negative_offset": schema.NumberAttribute{
                MarkdownDescription: "Negative Bucket Offset.",
                Optional: true,
                Computed: true,
            },
            "negative_bucket_counts": schema.StringAttribute{
                MarkdownDescription: "Negative Bucket Counts.",
                Optional: true,
                Computed: true,
            },
            "summary_quantiles": schema.StringAttribute{
                MarkdownDescription: "Summary Quantiles.",
                Optional: true,
                Computed: true,
            },
            "summary_values": schema.StringAttribute{
                MarkdownDescription: "Summary Values.",
                Optional: true,
                Computed: true,
            },
            "trace_id": schema.StringAttribute{
                MarkdownDescription: "Trace ID.",
                Optional: true,
                Computed: true,
            },
            "span_id": schema.StringAttribute{
                MarkdownDescription: "Span ID.",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *MetricDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *MetricDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data MetricDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.PrimaryEntityId.IsNull() && !data.PrimaryEntityId.IsUnknown() {
        filters["primaryEntityId"] = data.PrimaryEntityId.ValueString()
        filterNames = append(filterNames, "primary_entity_id = "+fmt.Sprintf("%q", data.PrimaryEntityId.ValueString()))
    }
    if !data.PrimaryEntityType.IsNull() && !data.PrimaryEntityType.IsUnknown() {
        filters["primaryEntityType"] = data.PrimaryEntityType.ValueString()
        filterNames = append(filterNames, "primary_entity_type = "+fmt.Sprintf("%q", data.PrimaryEntityType.ValueString()))
    }
    if !data.Name.IsNull() && !data.Name.IsUnknown() {
        filters["name"] = data.Name.ValueString()
        filterNames = append(filterNames, "name = "+fmt.Sprintf("%q", data.Name.ValueString()))
    }
    if !data.AggregationTemporality.IsNull() && !data.AggregationTemporality.IsUnknown() {
        filters["aggregationTemporality"] = data.AggregationTemporality.ValueString()
        filterNames = append(filterNames, "aggregation_temporality = "+fmt.Sprintf("%q", data.AggregationTemporality.ValueString()))
    }
    if !data.MetricPointType.IsNull() && !data.MetricPointType.IsUnknown() {
        filters["metricPointType"] = data.MetricPointType.ValueString()
        filterNames = append(filterNames, "metric_point_type = "+fmt.Sprintf("%q", data.MetricPointType.ValueString()))
    }
    if !data.Time.IsNull() && !data.Time.IsUnknown() {
        filters["time"] = data.Time.ValueString()
        filterNames = append(filterNames, "time = "+fmt.Sprintf("%q", data.Time.ValueString()))
    }
    if !data.StartTime.IsNull() && !data.StartTime.IsUnknown() {
        filters["startTime"] = data.StartTime.ValueString()
        filterNames = append(filterNames, "start_time = "+fmt.Sprintf("%q", data.StartTime.ValueString()))
    }
    if !data.TimeUnixNano.IsNull() && !data.TimeUnixNano.IsUnknown() {
        filters["timeUnixNano"] = data.TimeUnixNano.ValueString()
        filterNames = append(filterNames, "time_unix_nano = "+fmt.Sprintf("%q", data.TimeUnixNano.ValueString()))
    }
    if !data.StartTimeUnixNano.IsNull() && !data.StartTimeUnixNano.IsUnknown() {
        filters["startTimeUnixNano"] = data.StartTimeUnixNano.ValueString()
        filterNames = append(filterNames, "start_time_unix_nano = "+fmt.Sprintf("%q", data.StartTimeUnixNano.ValueString()))
    }
    if !data.Attributes.IsNull() && !data.Attributes.IsUnknown() {
        filters["attributes"] = data.Attributes.ValueString()
        filterNames = append(filterNames, "attributes = "+fmt.Sprintf("%q", data.Attributes.ValueString()))
    }
    if !data.ServiceEntityKey.IsNull() && !data.ServiceEntityKey.IsUnknown() {
        filters["serviceEntityKey"] = data.ServiceEntityKey.ValueString()
        filterNames = append(filterNames, "service_entity_key = "+fmt.Sprintf("%q", data.ServiceEntityKey.ValueString()))
    }
    if !data.HostEntityKey.IsNull() && !data.HostEntityKey.IsUnknown() {
        filters["hostEntityKey"] = data.HostEntityKey.ValueString()
        filterNames = append(filterNames, "host_entity_key = "+fmt.Sprintf("%q", data.HostEntityKey.ValueString()))
    }
    if !data.K8sPodEntityKey.IsNull() && !data.K8sPodEntityKey.IsUnknown() {
        filters["k8sPodEntityKey"] = data.K8sPodEntityKey.ValueString()
        filterNames = append(filterNames, "k8s_pod_entity_key = "+fmt.Sprintf("%q", data.K8sPodEntityKey.ValueString()))
    }
    if !data.K8sNodeEntityKey.IsNull() && !data.K8sNodeEntityKey.IsUnknown() {
        filters["k8sNodeEntityKey"] = data.K8sNodeEntityKey.ValueString()
        filterNames = append(filterNames, "k8s_node_entity_key = "+fmt.Sprintf("%q", data.K8sNodeEntityKey.ValueString()))
    }
    if !data.K8sClusterEntityKey.IsNull() && !data.K8sClusterEntityKey.IsUnknown() {
        filters["k8sClusterEntityKey"] = data.K8sClusterEntityKey.ValueString()
        filterNames = append(filterNames, "k8s_cluster_entity_key = "+fmt.Sprintf("%q", data.K8sClusterEntityKey.ValueString()))
    }
    if !data.ContainerEntityKey.IsNull() && !data.ContainerEntityKey.IsUnknown() {
        filters["containerEntityKey"] = data.ContainerEntityKey.ValueString()
        filterNames = append(filterNames, "container_entity_key = "+fmt.Sprintf("%q", data.ContainerEntityKey.ValueString()))
    }
    if !data.IsMonotonic.IsNull() && !data.IsMonotonic.IsUnknown() {
        filters["isMonotonic"] = data.IsMonotonic.ValueBool()
        filterNames = append(filterNames, "is_monotonic = "+fmt.Sprintf("%t", data.IsMonotonic.ValueBool()))
    }
    if !data.CountValue.IsNull() && !data.CountValue.IsUnknown() {
        filters["count"] = data.CountValue.ValueString()
        filterNames = append(filterNames, "count_value = "+fmt.Sprintf("%q", data.CountValue.ValueString()))
    }
    if !data.Sum.IsNull() && !data.Sum.IsUnknown() {
        filters["sum"] = lookupNumber(data.Sum)
        filterNames = append(filterNames, "sum = "+data.Sum.ValueBigFloat().String())
    }
    if !data.Value.IsNull() && !data.Value.IsUnknown() {
        filters["value"] = lookupNumber(data.Value)
        filterNames = append(filterNames, "value = "+data.Value.ValueBigFloat().String())
    }
    if !data.Min.IsNull() && !data.Min.IsUnknown() {
        filters["min"] = lookupNumber(data.Min)
        filterNames = append(filterNames, "min = "+data.Min.ValueBigFloat().String())
    }
    if !data.Max.IsNull() && !data.Max.IsUnknown() {
        filters["max"] = lookupNumber(data.Max)
        filterNames = append(filterNames, "max = "+data.Max.ValueBigFloat().String())
    }
    if !data.BucketCounts.IsNull() && !data.BucketCounts.IsUnknown() {
        filters["bucketCounts"] = data.BucketCounts.ValueString()
        filterNames = append(filterNames, "bucket_counts = "+fmt.Sprintf("%q", data.BucketCounts.ValueString()))
    }
    if !data.ExplicitBounds.IsNull() && !data.ExplicitBounds.IsUnknown() {
        filters["explicitBounds"] = data.ExplicitBounds.ValueString()
        filterNames = append(filterNames, "explicit_bounds = "+fmt.Sprintf("%q", data.ExplicitBounds.ValueString()))
    }
    if !data.Scale.IsNull() && !data.Scale.IsUnknown() {
        filters["scale"] = lookupNumber(data.Scale)
        filterNames = append(filterNames, "scale = "+data.Scale.ValueBigFloat().String())
    }
    if !data.ZeroCount.IsNull() && !data.ZeroCount.IsUnknown() {
        filters["zeroCount"] = data.ZeroCount.ValueString()
        filterNames = append(filterNames, "zero_count = "+fmt.Sprintf("%q", data.ZeroCount.ValueString()))
    }
    if !data.PositiveOffset.IsNull() && !data.PositiveOffset.IsUnknown() {
        filters["positiveOffset"] = lookupNumber(data.PositiveOffset)
        filterNames = append(filterNames, "positive_offset = "+data.PositiveOffset.ValueBigFloat().String())
    }
    if !data.PositiveBucketCounts.IsNull() && !data.PositiveBucketCounts.IsUnknown() {
        filters["positiveBucketCounts"] = data.PositiveBucketCounts.ValueString()
        filterNames = append(filterNames, "positive_bucket_counts = "+fmt.Sprintf("%q", data.PositiveBucketCounts.ValueString()))
    }
    if !data.NegativeOffset.IsNull() && !data.NegativeOffset.IsUnknown() {
        filters["negativeOffset"] = lookupNumber(data.NegativeOffset)
        filterNames = append(filterNames, "negative_offset = "+data.NegativeOffset.ValueBigFloat().String())
    }
    if !data.NegativeBucketCounts.IsNull() && !data.NegativeBucketCounts.IsUnknown() {
        filters["negativeBucketCounts"] = data.NegativeBucketCounts.ValueString()
        filterNames = append(filterNames, "negative_bucket_counts = "+fmt.Sprintf("%q", data.NegativeBucketCounts.ValueString()))
    }
    if !data.SummaryQuantiles.IsNull() && !data.SummaryQuantiles.IsUnknown() {
        filters["summaryQuantiles"] = data.SummaryQuantiles.ValueString()
        filterNames = append(filterNames, "summary_quantiles = "+fmt.Sprintf("%q", data.SummaryQuantiles.ValueString()))
    }
    if !data.SummaryValues.IsNull() && !data.SummaryValues.IsUnknown() {
        filters["summaryValues"] = data.SummaryValues.ValueString()
        filterNames = append(filterNames, "summary_values = "+fmt.Sprintf("%q", data.SummaryValues.ValueString()))
    }
    if !data.TraceId.IsNull() && !data.TraceId.IsUnknown() {
        filters["traceId"] = data.TraceId.ValueString()
        filterNames = append(filterNames, "trace_id = "+fmt.Sprintf("%q", data.TraceId.ValueString()))
    }
    if !data.SpanId.IsNull() && !data.SpanId.IsUnknown() {
        filters["spanId"] = data.SpanId.ValueString()
        filterNames = append(filterNames, "span_id = "+fmt.Sprintf("%q", data.SpanId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the metric up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the metric up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "projectId": true,
        "primaryEntityId": true,
        "primaryEntityType": true,
        "name": true,
        "aggregationTemporality": true,
        "metricPointType": true,
        "time": true,
        "startTime": true,
        "timeUnixNano": true,
        "startTimeUnixNano": true,
        "attributes": true,
        "attributeKeys": true,
        "entityKeys": true,
        "serviceEntityKey": true,
        "hostEntityKey": true,
        "k8sPodEntityKey": true,
        "k8sNodeEntityKey": true,
        "k8sClusterEntityKey": true,
        "containerEntityKey": true,
        "isMonotonic": true,
        "count": true,
        "sum": true,
        "value": true,
        "min": true,
        "max": true,
        "bucketCounts": true,
        "explicitBounds": true,
        "scale": true,
        "zeroCount": true,
        "positiveOffset": true,
        "positiveBucketCounts": true,
        "negativeOffset": true,
        "negativeBucketCounts": true,
        "summaryQuantiles": true,
        "summaryValues": true,
        "traceId": true,
        "spanId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/metrics/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read metric, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No metric found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read metric: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/metrics/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list metric, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list metric: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No metric matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one metric matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for metric.")
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
    if obj, ok := item["primaryEntityId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PrimaryEntityId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.PrimaryEntityId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.PrimaryEntityId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.PrimaryEntityId = types.StringValue(string(jsonBytes))
        } else {
            data.PrimaryEntityId = types.StringNull()
        }
    } else if val, ok := item["primaryEntityId"].(string); ok {
        data.PrimaryEntityId = types.StringValue(val)
    } else {
        data.PrimaryEntityId = types.StringNull()
    }
    if obj, ok := item["primaryEntityType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PrimaryEntityType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.PrimaryEntityType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.PrimaryEntityType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.PrimaryEntityType = types.StringValue(string(jsonBytes))
        } else {
            data.PrimaryEntityType = types.StringNull()
        }
    } else if val, ok := item["primaryEntityType"].(string); ok {
        data.PrimaryEntityType = types.StringValue(val)
    } else {
        data.PrimaryEntityType = types.StringNull()
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
    if obj, ok := item["aggregationTemporality"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AggregationTemporality = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AggregationTemporality = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AggregationTemporality = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AggregationTemporality = types.StringValue(string(jsonBytes))
        } else {
            data.AggregationTemporality = types.StringNull()
        }
    } else if val, ok := item["aggregationTemporality"].(string); ok {
        data.AggregationTemporality = types.StringValue(val)
    } else {
        data.AggregationTemporality = types.StringNull()
    }
    if obj, ok := item["metricPointType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MetricPointType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.MetricPointType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.MetricPointType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.MetricPointType = types.StringValue(string(jsonBytes))
        } else {
            data.MetricPointType = types.StringNull()
        }
    } else if val, ok := item["metricPointType"].(string); ok {
        data.MetricPointType = types.StringValue(val)
    } else {
        data.MetricPointType = types.StringNull()
    }
    if obj, ok := item["time"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Time = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Time = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Time = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Time = types.StringValue(string(jsonBytes))
        } else {
            data.Time = types.StringNull()
        }
    } else if val, ok := item["time"].(string); ok {
        data.Time = types.StringValue(val)
    } else {
        data.Time = types.StringNull()
    }
    if obj, ok := item["startTime"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StartTime = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.StartTime = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.StartTime = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.StartTime = types.StringValue(string(jsonBytes))
        } else {
            data.StartTime = types.StringNull()
        }
    } else if val, ok := item["startTime"].(string); ok {
        data.StartTime = types.StringValue(val)
    } else {
        data.StartTime = types.StringNull()
    }
    if obj, ok := item["timeUnixNano"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TimeUnixNano = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TimeUnixNano = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TimeUnixNano = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TimeUnixNano = types.StringValue(string(jsonBytes))
        } else {
            data.TimeUnixNano = types.StringNull()
        }
    } else if val, ok := item["timeUnixNano"].(string); ok {
        data.TimeUnixNano = types.StringValue(val)
    } else {
        data.TimeUnixNano = types.StringNull()
    }
    if obj, ok := item["startTimeUnixNano"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StartTimeUnixNano = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.StartTimeUnixNano = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.StartTimeUnixNano = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.StartTimeUnixNano = types.StringValue(string(jsonBytes))
        } else {
            data.StartTimeUnixNano = types.StringNull()
        }
    } else if val, ok := item["startTimeUnixNano"].(string); ok {
        data.StartTimeUnixNano = types.StringValue(val)
    } else {
        data.StartTimeUnixNano = types.StringNull()
    }
    if obj, ok := item["attributes"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Attributes = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Attributes = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Attributes = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Attributes = types.StringValue(string(jsonBytes))
        } else {
            data.Attributes = types.StringNull()
        }
    } else if val, ok := item["attributes"].(string); ok {
        data.Attributes = types.StringValue(val)
    } else {
        data.Attributes = types.StringNull()
    }
    if val, ok := item["attributeKeys"].([]interface{}); ok {
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
        data.AttributeKeys = types.SetValueMust(types.StringType, setItems)
    } else {
        data.AttributeKeys = types.SetNull(types.StringType)
    }
    if val, ok := item["entityKeys"].([]interface{}); ok {
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
        data.EntityKeys = types.SetValueMust(types.StringType, setItems)
    } else {
        data.EntityKeys = types.SetNull(types.StringType)
    }
    if obj, ok := item["serviceEntityKey"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ServiceEntityKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ServiceEntityKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ServiceEntityKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ServiceEntityKey = types.StringValue(string(jsonBytes))
        } else {
            data.ServiceEntityKey = types.StringNull()
        }
    } else if val, ok := item["serviceEntityKey"].(string); ok {
        data.ServiceEntityKey = types.StringValue(val)
    } else {
        data.ServiceEntityKey = types.StringNull()
    }
    if obj, ok := item["hostEntityKey"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.HostEntityKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.HostEntityKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.HostEntityKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.HostEntityKey = types.StringValue(string(jsonBytes))
        } else {
            data.HostEntityKey = types.StringNull()
        }
    } else if val, ok := item["hostEntityKey"].(string); ok {
        data.HostEntityKey = types.StringValue(val)
    } else {
        data.HostEntityKey = types.StringNull()
    }
    if obj, ok := item["k8sPodEntityKey"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.K8sPodEntityKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.K8sPodEntityKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.K8sPodEntityKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.K8sPodEntityKey = types.StringValue(string(jsonBytes))
        } else {
            data.K8sPodEntityKey = types.StringNull()
        }
    } else if val, ok := item["k8sPodEntityKey"].(string); ok {
        data.K8sPodEntityKey = types.StringValue(val)
    } else {
        data.K8sPodEntityKey = types.StringNull()
    }
    if obj, ok := item["k8sNodeEntityKey"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.K8sNodeEntityKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.K8sNodeEntityKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.K8sNodeEntityKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.K8sNodeEntityKey = types.StringValue(string(jsonBytes))
        } else {
            data.K8sNodeEntityKey = types.StringNull()
        }
    } else if val, ok := item["k8sNodeEntityKey"].(string); ok {
        data.K8sNodeEntityKey = types.StringValue(val)
    } else {
        data.K8sNodeEntityKey = types.StringNull()
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
    if obj, ok := item["containerEntityKey"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ContainerEntityKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ContainerEntityKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ContainerEntityKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ContainerEntityKey = types.StringValue(string(jsonBytes))
        } else {
            data.ContainerEntityKey = types.StringNull()
        }
    } else if val, ok := item["containerEntityKey"].(string); ok {
        data.ContainerEntityKey = types.StringValue(val)
    } else {
        data.ContainerEntityKey = types.StringNull()
    }
    if val, ok := item["isMonotonic"].(bool); ok {
        data.IsMonotonic = types.BoolValue(val)
    } else {
        data.IsMonotonic = types.BoolNull()
    }
    if obj, ok := item["count"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CountValue = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CountValue = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CountValue = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CountValue = types.StringValue(string(jsonBytes))
        } else {
            data.CountValue = types.StringNull()
        }
    } else if val, ok := item["count"].(string); ok {
        data.CountValue = types.StringValue(val)
    } else {
        data.CountValue = types.StringNull()
    }
    if val, ok := item["sum"].(float64); ok {
        data.Sum = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["sum"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.Sum = types.NumberValue(big.NewFloat(val))
        } else {
            data.Sum = types.NumberNull()
        }
    } else {
        data.Sum = types.NumberNull()
    }
    if val, ok := item["value"].(float64); ok {
        data.Value = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["value"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.Value = types.NumberValue(big.NewFloat(val))
        } else {
            data.Value = types.NumberNull()
        }
    } else {
        data.Value = types.NumberNull()
    }
    if val, ok := item["min"].(float64); ok {
        data.Min = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["min"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.Min = types.NumberValue(big.NewFloat(val))
        } else {
            data.Min = types.NumberNull()
        }
    } else {
        data.Min = types.NumberNull()
    }
    if val, ok := item["max"].(float64); ok {
        data.Max = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["max"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.Max = types.NumberValue(big.NewFloat(val))
        } else {
            data.Max = types.NumberNull()
        }
    } else {
        data.Max = types.NumberNull()
    }
    if obj, ok := item["bucketCounts"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.BucketCounts = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.BucketCounts = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.BucketCounts = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.BucketCounts = types.StringValue(string(jsonBytes))
        } else {
            data.BucketCounts = types.StringNull()
        }
    } else if val, ok := item["bucketCounts"].(string); ok {
        data.BucketCounts = types.StringValue(val)
    } else {
        data.BucketCounts = types.StringNull()
    }
    if obj, ok := item["explicitBounds"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ExplicitBounds = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ExplicitBounds = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ExplicitBounds = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ExplicitBounds = types.StringValue(string(jsonBytes))
        } else {
            data.ExplicitBounds = types.StringNull()
        }
    } else if val, ok := item["explicitBounds"].(string); ok {
        data.ExplicitBounds = types.StringValue(val)
    } else {
        data.ExplicitBounds = types.StringNull()
    }
    if val, ok := item["scale"].(float64); ok {
        data.Scale = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["scale"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.Scale = types.NumberValue(big.NewFloat(val))
        } else {
            data.Scale = types.NumberNull()
        }
    } else {
        data.Scale = types.NumberNull()
    }
    if obj, ok := item["zeroCount"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ZeroCount = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ZeroCount = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ZeroCount = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ZeroCount = types.StringValue(string(jsonBytes))
        } else {
            data.ZeroCount = types.StringNull()
        }
    } else if val, ok := item["zeroCount"].(string); ok {
        data.ZeroCount = types.StringValue(val)
    } else {
        data.ZeroCount = types.StringNull()
    }
    if val, ok := item["positiveOffset"].(float64); ok {
        data.PositiveOffset = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["positiveOffset"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.PositiveOffset = types.NumberValue(big.NewFloat(val))
        } else {
            data.PositiveOffset = types.NumberNull()
        }
    } else {
        data.PositiveOffset = types.NumberNull()
    }
    if obj, ok := item["positiveBucketCounts"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PositiveBucketCounts = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.PositiveBucketCounts = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.PositiveBucketCounts = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.PositiveBucketCounts = types.StringValue(string(jsonBytes))
        } else {
            data.PositiveBucketCounts = types.StringNull()
        }
    } else if val, ok := item["positiveBucketCounts"].(string); ok {
        data.PositiveBucketCounts = types.StringValue(val)
    } else {
        data.PositiveBucketCounts = types.StringNull()
    }
    if val, ok := item["negativeOffset"].(float64); ok {
        data.NegativeOffset = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["negativeOffset"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.NegativeOffset = types.NumberValue(big.NewFloat(val))
        } else {
            data.NegativeOffset = types.NumberNull()
        }
    } else {
        data.NegativeOffset = types.NumberNull()
    }
    if obj, ok := item["negativeBucketCounts"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.NegativeBucketCounts = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.NegativeBucketCounts = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.NegativeBucketCounts = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.NegativeBucketCounts = types.StringValue(string(jsonBytes))
        } else {
            data.NegativeBucketCounts = types.StringNull()
        }
    } else if val, ok := item["negativeBucketCounts"].(string); ok {
        data.NegativeBucketCounts = types.StringValue(val)
    } else {
        data.NegativeBucketCounts = types.StringNull()
    }
    if obj, ok := item["summaryQuantiles"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SummaryQuantiles = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SummaryQuantiles = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SummaryQuantiles = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SummaryQuantiles = types.StringValue(string(jsonBytes))
        } else {
            data.SummaryQuantiles = types.StringNull()
        }
    } else if val, ok := item["summaryQuantiles"].(string); ok {
        data.SummaryQuantiles = types.StringValue(val)
    } else {
        data.SummaryQuantiles = types.StringNull()
    }
    if obj, ok := item["summaryValues"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SummaryValues = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SummaryValues = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SummaryValues = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SummaryValues = types.StringValue(string(jsonBytes))
        } else {
            data.SummaryValues = types.StringNull()
        }
    } else if val, ok := item["summaryValues"].(string); ok {
        data.SummaryValues = types.StringValue(val)
    } else {
        data.SummaryValues = types.StringNull()
    }
    if obj, ok := item["traceId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TraceId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TraceId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TraceId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TraceId = types.StringValue(string(jsonBytes))
        } else {
            data.TraceId = types.StringNull()
        }
    } else if val, ok := item["traceId"].(string); ok {
        data.TraceId = types.StringValue(val)
    } else {
        data.TraceId = types.StringNull()
    }
    if obj, ok := item["spanId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SpanId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SpanId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SpanId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SpanId = types.StringValue(string(jsonBytes))
        } else {
            data.SpanId = types.StringNull()
        }
    } else if val, ok := item["spanId"].(string); ok {
        data.SpanId = types.StringValue(val)
    } else {
        data.SpanId = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
