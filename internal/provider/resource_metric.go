package provider

import (
    "context"
    "fmt"
    "github.com/hashicorp/terraform-plugin-framework/path"
    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
    "math/big"
    "net/http"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &MetricResource{}
var _ resource.ResourceWithImportState = &MetricResource{}

func NewMetricResource() resource.Resource {
    return &MetricResource{}
}

// MetricResource defines the resource implementation.
type MetricResource struct {
    client *Client
}

// MetricResourceModel describes the resource data model.
type MetricResourceModel struct {
    Id types.String `tfsdk:"id"`
    ProjectId types.String `tfsdk:"project_id"`
    ServiceId types.String `tfsdk:"service_id"`
    ServiceType types.String `tfsdk:"service_type"`
    Name types.String `tfsdk:"name"`
    AggregationTemporality types.String `tfsdk:"aggregation_temporality"`
    MetricPointType types.String `tfsdk:"metric_point_type"`
    Time types.String `tfsdk:"time"`
    StartTime types.String `tfsdk:"start_time"`
    TimeUnixNano types.Number `tfsdk:"time_unix_nano"`
    StartTimeUnixNano types.Number `tfsdk:"start_time_unix_nano"`
    Attributes types.Map `tfsdk:"attributes"`
    IsMonotonic types.Bool `tfsdk:"is_monotonic"`
    CountValue types.Number `tfsdk:"count_value"`
    Sum types.Number `tfsdk:"sum"`
    Value types.Number `tfsdk:"value"`
    Min types.Number `tfsdk:"min"`
    Max types.Number `tfsdk:"max"`
    BucketCounts types.List `tfsdk:"bucket_counts"`
    ExplicitBounds types.List `tfsdk:"explicit_bounds"`
}

func (r *MetricResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_metric"
}

func (r *MetricResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "metric resource",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Unique identifier for the resource",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
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
            "name": schema.StringAttribute{
                MarkdownDescription: "Name",
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
            "attributes": schema.MapAttribute{
                MarkdownDescription: "Attributes",
                Computed: true,
                ElementType: types.StringType,
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

func (r *MetricResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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


func (r *MetricResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var data MetricResourceModel

    // Read Terraform plan data into the model
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create API request body
    metricRequest := map[string]interface{}{
        "data": map[string]interface{}{

        },
    }

    // Make API call
    httpResp, err := r.client.Post("/metrics", metricRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create metric, got error: %s", err))
        return
    }

    var metricResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &metricResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse metric response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := metricResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = metricResponse
    }

    if val, ok := dataMap["id"].(string); ok && val != "" {
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
    if val, ok := dataMap["serviceId"].(string); ok && val != "" {
        data.ServiceId = types.StringValue(val)
    } else {
        data.ServiceId = types.StringNull()
    }
    if val, ok := dataMap["serviceType"].(string); ok && val != "" {
        data.ServiceType = types.StringValue(val)
    } else {
        data.ServiceType = types.StringNull()
    }
    if val, ok := dataMap["name"].(string); ok && val != "" {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if val, ok := dataMap["aggregationTemporality"].(string); ok && val != "" {
        data.AggregationTemporality = types.StringValue(val)
    } else {
        data.AggregationTemporality = types.StringNull()
    }
    if val, ok := dataMap["metricPointType"].(string); ok && val != "" {
        data.MetricPointType = types.StringValue(val)
    } else {
        data.MetricPointType = types.StringNull()
    }
    if val, ok := dataMap["time"].(string); ok && val != "" {
        data.Time = types.StringValue(val)
    } else {
        data.Time = types.StringNull()
    }
    if val, ok := dataMap["startTime"].(string); ok && val != "" {
        data.StartTime = types.StringValue(val)
    } else {
        data.StartTime = types.StringNull()
    }
    if val, ok := dataMap["timeUnixNano"].(float64); ok {
        data.TimeUnixNano = types.NumberValue(big.NewFloat(val))
    } else if dataMap["timeUnixNano"] == nil {
        data.TimeUnixNano = types.NumberNull()
    }
    if val, ok := dataMap["startTimeUnixNano"].(float64); ok {
        data.StartTimeUnixNano = types.NumberValue(big.NewFloat(val))
    } else if dataMap["startTimeUnixNano"] == nil {
        data.StartTimeUnixNano = types.NumberNull()
    }
    if val, ok := dataMap["attributes"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.Attributes = mapValue
    } else if dataMap["attributes"] == nil {
        data.Attributes = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["isMonotonic"].(bool); ok {
        data.IsMonotonic = types.BoolValue(val)
    } else if dataMap["isMonotonic"] == nil {
        data.IsMonotonic = types.BoolNull()
    }
    if val, ok := dataMap["count"].(float64); ok {
        data.CountValue = types.NumberValue(big.NewFloat(val))
    } else if dataMap["count"] == nil {
        data.CountValue = types.NumberNull()
    }
    if val, ok := dataMap["sum"].(float64); ok {
        data.Sum = types.NumberValue(big.NewFloat(val))
    } else if dataMap["sum"] == nil {
        data.Sum = types.NumberNull()
    }
    if val, ok := dataMap["value"].(float64); ok {
        data.Value = types.NumberValue(big.NewFloat(val))
    } else if dataMap["value"] == nil {
        data.Value = types.NumberNull()
    }
    if val, ok := dataMap["min"].(float64); ok {
        data.Min = types.NumberValue(big.NewFloat(val))
    } else if dataMap["min"] == nil {
        data.Min = types.NumberNull()
    }
    if val, ok := dataMap["max"].(float64); ok {
        data.Max = types.NumberValue(big.NewFloat(val))
    } else if dataMap["max"] == nil {
        data.Max = types.NumberNull()
    }
    if val, ok := dataMap["bucketCounts"].([]interface{}); ok {
        // Convert API response list to Terraform list
        listValue, _ := types.ListValueFrom(ctx, types.StringType, val)
        data.BucketCounts = listValue
    } else if dataMap["bucketCounts"] == nil {
        data.BucketCounts = types.ListNull(types.StringType)
    }
    if val, ok := dataMap["explicitBounds"].([]interface{}); ok {
        // Convert API response list to Terraform list
        listValue, _ := types.ListValueFrom(ctx, types.StringType, val)
        data.ExplicitBounds = listValue
    } else if dataMap["explicitBounds"] == nil {
        data.ExplicitBounds = types.ListNull(types.StringType)
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

func (r *MetricResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var data MetricResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create select parameter to get full object
    selectParam := map[string]interface{}{
        "projectId": true,
        "serviceId": true,
        "serviceType": true,
        "name": true,
        "aggregationTemporality": true,
        "metricPointType": true,
        "time": true,
        "startTime": true,
        "timeUnixNano": true,
        "startTimeUnixNano": true,
        "attributes": true,
        "isMonotonic": true,
        "count": true,
        "sum": true,
        "value": true,
        "min": true,
        "max": true,
        "bucketCounts": true,
        "explicitBounds": true,
        "_id": true,
    }

    // Make API call with select parameter
    httpResp, err := r.client.PostWithSelect("/metrics/" + data.Id.ValueString() + "", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read metric, got error: %s", err))
        return
    }

    if httpResp.StatusCode == http.StatusNotFound {
        resp.State.RemoveResource(ctx)
        return
    }

    var metricResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &metricResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse metric response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := metricResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = metricResponse
    }

    if val, ok := dataMap["id"].(string); ok && val != "" {
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
    if val, ok := dataMap["serviceId"].(string); ok && val != "" {
        data.ServiceId = types.StringValue(val)
    } else {
        data.ServiceId = types.StringNull()
    }
    if val, ok := dataMap["serviceType"].(string); ok && val != "" {
        data.ServiceType = types.StringValue(val)
    } else {
        data.ServiceType = types.StringNull()
    }
    if val, ok := dataMap["name"].(string); ok && val != "" {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if val, ok := dataMap["aggregationTemporality"].(string); ok && val != "" {
        data.AggregationTemporality = types.StringValue(val)
    } else {
        data.AggregationTemporality = types.StringNull()
    }
    if val, ok := dataMap["metricPointType"].(string); ok && val != "" {
        data.MetricPointType = types.StringValue(val)
    } else {
        data.MetricPointType = types.StringNull()
    }
    if val, ok := dataMap["time"].(string); ok && val != "" {
        data.Time = types.StringValue(val)
    } else {
        data.Time = types.StringNull()
    }
    if val, ok := dataMap["startTime"].(string); ok && val != "" {
        data.StartTime = types.StringValue(val)
    } else {
        data.StartTime = types.StringNull()
    }
    if val, ok := dataMap["timeUnixNano"].(float64); ok {
        data.TimeUnixNano = types.NumberValue(big.NewFloat(val))
    } else if dataMap["timeUnixNano"] == nil {
        data.TimeUnixNano = types.NumberNull()
    }
    if val, ok := dataMap["startTimeUnixNano"].(float64); ok {
        data.StartTimeUnixNano = types.NumberValue(big.NewFloat(val))
    } else if dataMap["startTimeUnixNano"] == nil {
        data.StartTimeUnixNano = types.NumberNull()
    }
    if val, ok := dataMap["attributes"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.Attributes = mapValue
    } else if dataMap["attributes"] == nil {
        data.Attributes = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["isMonotonic"].(bool); ok {
        data.IsMonotonic = types.BoolValue(val)
    } else if dataMap["isMonotonic"] == nil {
        data.IsMonotonic = types.BoolNull()
    }
    if val, ok := dataMap["count"].(float64); ok {
        data.CountValue = types.NumberValue(big.NewFloat(val))
    } else if dataMap["count"] == nil {
        data.CountValue = types.NumberNull()
    }
    if val, ok := dataMap["sum"].(float64); ok {
        data.Sum = types.NumberValue(big.NewFloat(val))
    } else if dataMap["sum"] == nil {
        data.Sum = types.NumberNull()
    }
    if val, ok := dataMap["value"].(float64); ok {
        data.Value = types.NumberValue(big.NewFloat(val))
    } else if dataMap["value"] == nil {
        data.Value = types.NumberNull()
    }
    if val, ok := dataMap["min"].(float64); ok {
        data.Min = types.NumberValue(big.NewFloat(val))
    } else if dataMap["min"] == nil {
        data.Min = types.NumberNull()
    }
    if val, ok := dataMap["max"].(float64); ok {
        data.Max = types.NumberValue(big.NewFloat(val))
    } else if dataMap["max"] == nil {
        data.Max = types.NumberNull()
    }
    if val, ok := dataMap["bucketCounts"].([]interface{}); ok {
        // Convert API response list to Terraform list
        listValue, _ := types.ListValueFrom(ctx, types.StringType, val)
        data.BucketCounts = listValue
    } else if dataMap["bucketCounts"] == nil {
        data.BucketCounts = types.ListNull(types.StringType)
    }
    if val, ok := dataMap["explicitBounds"].([]interface{}); ok {
        // Convert API response list to Terraform list
        listValue, _ := types.ListValueFrom(ctx, types.StringType, val)
        data.ExplicitBounds = listValue
    } else if dataMap["explicitBounds"] == nil {
        data.ExplicitBounds = types.ListNull(types.StringType)
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MetricResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var data MetricResourceModel
    var state MetricResourceModel

    // Read Terraform current state data to get the ID
    resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Read Terraform plan data to get the new values
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Use the ID from the current state
    data.Id = state.Id

    // Create API request body
    metricRequest := map[string]interface{}{
        "data": map[string]interface{}{

        },
    }

    // Make API call
    httpResp, err := r.client.Put("/metrics/" + data.Id.ValueString() + "", metricRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update metric, got error: %s", err))
        return
    }

    // Parse the update response
    var metricResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &metricResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse metric response, got error: %s", err))
        return
    }

    // After successful update, fetch the current state by calling Read with select parameter
    selectParam := map[string]interface{}{
        "projectId": true,
        "serviceId": true,
        "serviceType": true,
        "name": true,
        "aggregationTemporality": true,
        "metricPointType": true,
        "time": true,
        "startTime": true,
        "timeUnixNano": true,
        "startTimeUnixNano": true,
        "attributes": true,
        "isMonotonic": true,
        "count": true,
        "sum": true,
        "value": true,
        "min": true,
        "max": true,
        "bucketCounts": true,
        "explicitBounds": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect("/metrics/" + data.Id.ValueString() + "", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read metric after update, got error: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse metric read response, got error: %s", err))
        return
    }

    // Update the model with response data from the Read operation
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := readResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = readResponse
    }

    if val, ok := dataMap["id"].(string); ok && val != "" {
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
    if val, ok := dataMap["serviceId"].(string); ok && val != "" {
        data.ServiceId = types.StringValue(val)
    } else {
        data.ServiceId = types.StringNull()
    }
    if val, ok := dataMap["serviceType"].(string); ok && val != "" {
        data.ServiceType = types.StringValue(val)
    } else {
        data.ServiceType = types.StringNull()
    }
    if val, ok := dataMap["name"].(string); ok && val != "" {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if val, ok := dataMap["aggregationTemporality"].(string); ok && val != "" {
        data.AggregationTemporality = types.StringValue(val)
    } else {
        data.AggregationTemporality = types.StringNull()
    }
    if val, ok := dataMap["metricPointType"].(string); ok && val != "" {
        data.MetricPointType = types.StringValue(val)
    } else {
        data.MetricPointType = types.StringNull()
    }
    if val, ok := dataMap["time"].(string); ok && val != "" {
        data.Time = types.StringValue(val)
    } else {
        data.Time = types.StringNull()
    }
    if val, ok := dataMap["startTime"].(string); ok && val != "" {
        data.StartTime = types.StringValue(val)
    } else {
        data.StartTime = types.StringNull()
    }
    if val, ok := dataMap["timeUnixNano"].(float64); ok {
        data.TimeUnixNano = types.NumberValue(big.NewFloat(val))
    } else if dataMap["timeUnixNano"] == nil {
        data.TimeUnixNano = types.NumberNull()
    }
    if val, ok := dataMap["startTimeUnixNano"].(float64); ok {
        data.StartTimeUnixNano = types.NumberValue(big.NewFloat(val))
    } else if dataMap["startTimeUnixNano"] == nil {
        data.StartTimeUnixNano = types.NumberNull()
    }
    if val, ok := dataMap["attributes"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.Attributes = mapValue
    } else if dataMap["attributes"] == nil {
        data.Attributes = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["isMonotonic"].(bool); ok {
        data.IsMonotonic = types.BoolValue(val)
    } else if dataMap["isMonotonic"] == nil {
        data.IsMonotonic = types.BoolNull()
    }
    if val, ok := dataMap["count"].(float64); ok {
        data.CountValue = types.NumberValue(big.NewFloat(val))
    } else if dataMap["count"] == nil {
        data.CountValue = types.NumberNull()
    }
    if val, ok := dataMap["sum"].(float64); ok {
        data.Sum = types.NumberValue(big.NewFloat(val))
    } else if dataMap["sum"] == nil {
        data.Sum = types.NumberNull()
    }
    if val, ok := dataMap["value"].(float64); ok {
        data.Value = types.NumberValue(big.NewFloat(val))
    } else if dataMap["value"] == nil {
        data.Value = types.NumberNull()
    }
    if val, ok := dataMap["min"].(float64); ok {
        data.Min = types.NumberValue(big.NewFloat(val))
    } else if dataMap["min"] == nil {
        data.Min = types.NumberNull()
    }
    if val, ok := dataMap["max"].(float64); ok {
        data.Max = types.NumberValue(big.NewFloat(val))
    } else if dataMap["max"] == nil {
        data.Max = types.NumberNull()
    }
    if val, ok := dataMap["bucketCounts"].([]interface{}); ok {
        // Convert API response list to Terraform list
        listValue, _ := types.ListValueFrom(ctx, types.StringType, val)
        data.BucketCounts = listValue
    } else if dataMap["bucketCounts"] == nil {
        data.BucketCounts = types.ListNull(types.StringType)
    }
    if val, ok := dataMap["explicitBounds"].([]interface{}); ok {
        // Convert API response list to Terraform list
        listValue, _ := types.ListValueFrom(ctx, types.StringType, val)
        data.ExplicitBounds = listValue
    } else if dataMap["explicitBounds"] == nil {
        data.ExplicitBounds = types.ListNull(types.StringType)
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MetricResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var data MetricResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Make API call
    _, err := r.client.Delete("/metrics/" + data.Id.ValueString() + "")
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete metric, got error: %s", err))
        return
    }
}


func (r *MetricResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Helper method to convert Terraform map to Go interface{}
func (r *MetricResource) convertTerraformMapToInterface(terraformMap types.Map) interface{} {
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
func (r *MetricResource) convertTerraformListToInterface(terraformList types.List) interface{} {
    if terraformList.IsNull() || terraformList.IsUnknown() {
        return nil
    }
    
    var result []interface{}
    terraformList.ElementsAs(context.Background(), &result, false)
    return result
}
