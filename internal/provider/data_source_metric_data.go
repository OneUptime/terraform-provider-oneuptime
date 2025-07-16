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
    ServiceId types.String `tfsdk:"service_id"`
    ServiceType types.String `tfsdk:"service_type"`
    AggregationTemporality types.String `tfsdk:"aggregation_temporality"`
    MetricPointType types.String `tfsdk:"metric_point_type"`
    Time types.String `tfsdk:"time"`
    StartTime types.String `tfsdk:"start_time"`
    TimeUnixNano types.Number `tfsdk:"time_unix_nano"`
    StartTimeUnixNano types.Number `tfsdk:"start_time_unix_nano"`
    Attributes types.String `tfsdk:"attributes"`
    IsMonotonic types.Bool `tfsdk:"is_monotonic"`
    CountValue types.Number `tfsdk:"count_value"`
    Sum types.Number `tfsdk:"sum"`
    Value types.Number `tfsdk:"value"`
    Min types.Number `tfsdk:"min"`
    Max types.Number `tfsdk:"max"`
    BucketCounts types.List `tfsdk:"bucket_counts"`
    ExplicitBounds types.List `tfsdk:"explicit_bounds"`
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
            "service_id": schema.StringAttribute{
                MarkdownDescription: "Service ID",
                Computed: true,
            },
            "service_type": schema.StringAttribute{
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
            "time_unix_nano": schema.NumberAttribute{
                MarkdownDescription: "Time (in Unix Nano)",
                Computed: true,
            },
            "start_time_unix_nano": schema.NumberAttribute{
                MarkdownDescription: "Start Time (in Unix Nano)",
                Computed: true,
            },
            "attributes": schema.StringAttribute{
                MarkdownDescription: "Attributes",
                Computed: true,
            },
            "is_monotonic": schema.BoolAttribute{
                MarkdownDescription: "Is Monotonic",
                Computed: true,
            },
            "count_value": schema.NumberAttribute{
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
            "bucket_counts": schema.ListAttribute{
                MarkdownDescription: "Bucket Counts",
                Computed: true,
                ElementType: types.StringType,
            },
            "explicit_bounds": schema.ListAttribute{
                MarkdownDescription: "Explicit Bonds",
                Computed: true,
                ElementType: types.StringType,
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
    if val, ok := metricDataResponse["service_id"].(string); ok {
        data.ServiceId = types.StringValue(val)
    }
    if val, ok := metricDataResponse["service_type"].(string); ok {
        data.ServiceType = types.StringValue(val)
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
    if val, ok := metricDataResponse["time_unix_nano"].(float64); ok {
        data.TimeUnixNano = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := metricDataResponse["start_time_unix_nano"].(float64); ok {
        data.StartTimeUnixNano = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := metricDataResponse["attributes"].(string); ok {
        data.Attributes = types.StringValue(val)
    }
    if val, ok := metricDataResponse["is_monotonic"].(bool); ok {
        data.IsMonotonic = types.BoolValue(val)
    }
    if val, ok := metricDataResponse["count"].(float64); ok {
        data.CountValue = types.NumberValue(big.NewFloat(val))
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
    if val, ok := metricDataResponse["bucket_counts"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        listValue, _ := types.ListValue(types.StringType, elements)
        data.BucketCounts = listValue
    }
    if val, ok := metricDataResponse["explicit_bounds"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        listValue, _ := types.ListValue(types.StringType, elements)
        data.ExplicitBounds = listValue
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
