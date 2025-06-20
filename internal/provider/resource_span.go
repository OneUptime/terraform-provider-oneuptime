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
var _ resource.Resource = &SpanResource{}
var _ resource.ResourceWithImportState = &SpanResource{}

func NewSpanResource() resource.Resource {
    return &SpanResource{}
}

// SpanResource defines the resource implementation.
type SpanResource struct {
    client *Client
}

// SpanResourceModel describes the resource data model.
type SpanResourceModel struct {
    Id types.String `tfsdk:"id"`
    ProjectId types.String `tfsdk:"project_id"`
    ServiceId types.String `tfsdk:"service_id"`
    StartTime types.String `tfsdk:"start_time"`
    EndTime types.String `tfsdk:"end_time"`
    StartTimeUnixNano types.Number `tfsdk:"start_time_unix_nano"`
    DurationUnixNano types.Number `tfsdk:"duration_unix_nano"`
    EndTimeUnixNano types.Number `tfsdk:"end_time_unix_nano"`
    TraceId types.String `tfsdk:"trace_id"`
    SpanId types.String `tfsdk:"span_id"`
    ParentSpanId types.String `tfsdk:"parent_span_id"`
    TraceState types.String `tfsdk:"trace_state"`
    Attributes types.Map `tfsdk:"attributes"`
    Events types.List `tfsdk:"events"`
    Links types.Map `tfsdk:"links"`
    StatusCode types.Number `tfsdk:"status_code"`
    StatusMessage types.String `tfsdk:"status_message"`
    Name types.String `tfsdk:"name"`
    Kind types.String `tfsdk:"kind"`
}

func (r *SpanResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_span"
}

func (r *SpanResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "span resource",

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
            "start_time": schema.StringAttribute{
                MarkdownDescription: "Start Time",
                Computed: true,
            },
            "end_time": schema.StringAttribute{
                MarkdownDescription: "End Time",
                Computed: true,
            },
            "start_time_unix_nano": schema.NumberAttribute{
                MarkdownDescription: "Start Time in Unix Nano",
                Computed: true,
            },
            "duration_unix_nano": schema.NumberAttribute{
                MarkdownDescription: "Duration in Unix Nano",
                Computed: true,
            },
            "end_time_unix_nano": schema.NumberAttribute{
                MarkdownDescription: "End Time",
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
            "parent_span_id": schema.StringAttribute{
                MarkdownDescription: "Parent Span ID",
                Computed: true,
            },
            "trace_state": schema.StringAttribute{
                MarkdownDescription: "Trace State",
                Computed: true,
            },
            "attributes": schema.MapAttribute{
                MarkdownDescription: "Attributes",
                Computed: true,
                ElementType: types.StringType,
            },
            "events": schema.ListAttribute{
                MarkdownDescription: "Events",
                Computed: true,
                ElementType: types.StringType,
            },
            "links": schema.MapAttribute{
                MarkdownDescription: "Links",
                Computed: true,
                ElementType: types.StringType,
            },
            "status_code": schema.NumberAttribute{
                MarkdownDescription: "Status Code",
                Computed: true,
            },
            "status_message": schema.StringAttribute{
                MarkdownDescription: "Status Message",
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Name",
                Computed: true,
            },
            "kind": schema.StringAttribute{
                MarkdownDescription: "Kind",
                Computed: true,
            },
        },
    }
}

func (r *SpanResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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


func (r *SpanResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var data SpanResourceModel

    // Read Terraform plan data into the model
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create API request body
    spanRequest := map[string]interface{}{
        "data": map[string]interface{}{

        },
    }

    // Make API call
    httpResp, err := r.client.Post("/span", spanRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create span, got error: %s", err))
        return
    }

    var spanResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &spanResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse span response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := spanResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = spanResponse
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
    if val, ok := dataMap["startTime"].(string); ok && val != "" {
        data.StartTime = types.StringValue(val)
    } else {
        data.StartTime = types.StringNull()
    }
    if val, ok := dataMap["endTime"].(string); ok && val != "" {
        data.EndTime = types.StringValue(val)
    } else {
        data.EndTime = types.StringNull()
    }
    if val, ok := dataMap["startTimeUnixNano"].(float64); ok {
        data.StartTimeUnixNano = types.NumberValue(big.NewFloat(val))
    } else if dataMap["startTimeUnixNano"] == nil {
        data.StartTimeUnixNano = types.NumberNull()
    }
    if val, ok := dataMap["durationUnixNano"].(float64); ok {
        data.DurationUnixNano = types.NumberValue(big.NewFloat(val))
    } else if dataMap["durationUnixNano"] == nil {
        data.DurationUnixNano = types.NumberNull()
    }
    if val, ok := dataMap["endTimeUnixNano"].(float64); ok {
        data.EndTimeUnixNano = types.NumberValue(big.NewFloat(val))
    } else if dataMap["endTimeUnixNano"] == nil {
        data.EndTimeUnixNano = types.NumberNull()
    }
    if val, ok := dataMap["traceId"].(string); ok && val != "" {
        data.TraceId = types.StringValue(val)
    } else {
        data.TraceId = types.StringNull()
    }
    if val, ok := dataMap["spanId"].(string); ok && val != "" {
        data.SpanId = types.StringValue(val)
    } else {
        data.SpanId = types.StringNull()
    }
    if val, ok := dataMap["parentSpanId"].(string); ok && val != "" {
        data.ParentSpanId = types.StringValue(val)
    } else {
        data.ParentSpanId = types.StringNull()
    }
    if val, ok := dataMap["traceState"].(string); ok && val != "" {
        data.TraceState = types.StringValue(val)
    } else {
        data.TraceState = types.StringNull()
    }
    if val, ok := dataMap["attributes"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.Attributes = mapValue
    } else if dataMap["attributes"] == nil {
        data.Attributes = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["events"].([]interface{}); ok {
        // Convert API response list to Terraform list
        listValue, _ := types.ListValueFrom(ctx, types.StringType, val)
        data.Events = listValue
    } else if dataMap["events"] == nil {
        data.Events = types.ListNull(types.StringType)
    }
    if val, ok := dataMap["links"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.Links = mapValue
    } else if dataMap["links"] == nil {
        data.Links = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["statusCode"].(float64); ok {
        data.StatusCode = types.NumberValue(big.NewFloat(val))
    } else if dataMap["statusCode"] == nil {
        data.StatusCode = types.NumberNull()
    }
    if val, ok := dataMap["statusMessage"].(string); ok && val != "" {
        data.StatusMessage = types.StringValue(val)
    } else {
        data.StatusMessage = types.StringNull()
    }
    if val, ok := dataMap["name"].(string); ok && val != "" {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if val, ok := dataMap["kind"].(string); ok && val != "" {
        data.Kind = types.StringValue(val)
    } else {
        data.Kind = types.StringNull()
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

func (r *SpanResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var data SpanResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create select parameter to get full object
    selectParam := map[string]interface{}{
        "projectId": true,
        "serviceId": true,
        "startTime": true,
        "endTime": true,
        "startTimeUnixNano": true,
        "durationUnixNano": true,
        "endTimeUnixNano": true,
        "traceId": true,
        "spanId": true,
        "parentSpanId": true,
        "traceState": true,
        "attributes": true,
        "events": true,
        "links": true,
        "statusCode": true,
        "statusMessage": true,
        "name": true,
        "kind": true,
        "_id": true,
    }

    // Make API call with select parameter
    httpResp, err := r.client.PostWithSelect("/span/" + data.Id.ValueString() + "", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read span, got error: %s", err))
        return
    }

    if httpResp.StatusCode == http.StatusNotFound {
        resp.State.RemoveResource(ctx)
        return
    }

    var spanResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &spanResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse span response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := spanResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = spanResponse
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
    if val, ok := dataMap["startTime"].(string); ok && val != "" {
        data.StartTime = types.StringValue(val)
    } else {
        data.StartTime = types.StringNull()
    }
    if val, ok := dataMap["endTime"].(string); ok && val != "" {
        data.EndTime = types.StringValue(val)
    } else {
        data.EndTime = types.StringNull()
    }
    if val, ok := dataMap["startTimeUnixNano"].(float64); ok {
        data.StartTimeUnixNano = types.NumberValue(big.NewFloat(val))
    } else if dataMap["startTimeUnixNano"] == nil {
        data.StartTimeUnixNano = types.NumberNull()
    }
    if val, ok := dataMap["durationUnixNano"].(float64); ok {
        data.DurationUnixNano = types.NumberValue(big.NewFloat(val))
    } else if dataMap["durationUnixNano"] == nil {
        data.DurationUnixNano = types.NumberNull()
    }
    if val, ok := dataMap["endTimeUnixNano"].(float64); ok {
        data.EndTimeUnixNano = types.NumberValue(big.NewFloat(val))
    } else if dataMap["endTimeUnixNano"] == nil {
        data.EndTimeUnixNano = types.NumberNull()
    }
    if val, ok := dataMap["traceId"].(string); ok && val != "" {
        data.TraceId = types.StringValue(val)
    } else {
        data.TraceId = types.StringNull()
    }
    if val, ok := dataMap["spanId"].(string); ok && val != "" {
        data.SpanId = types.StringValue(val)
    } else {
        data.SpanId = types.StringNull()
    }
    if val, ok := dataMap["parentSpanId"].(string); ok && val != "" {
        data.ParentSpanId = types.StringValue(val)
    } else {
        data.ParentSpanId = types.StringNull()
    }
    if val, ok := dataMap["traceState"].(string); ok && val != "" {
        data.TraceState = types.StringValue(val)
    } else {
        data.TraceState = types.StringNull()
    }
    if val, ok := dataMap["attributes"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.Attributes = mapValue
    } else if dataMap["attributes"] == nil {
        data.Attributes = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["events"].([]interface{}); ok {
        // Convert API response list to Terraform list
        listValue, _ := types.ListValueFrom(ctx, types.StringType, val)
        data.Events = listValue
    } else if dataMap["events"] == nil {
        data.Events = types.ListNull(types.StringType)
    }
    if val, ok := dataMap["links"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.Links = mapValue
    } else if dataMap["links"] == nil {
        data.Links = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["statusCode"].(float64); ok {
        data.StatusCode = types.NumberValue(big.NewFloat(val))
    } else if dataMap["statusCode"] == nil {
        data.StatusCode = types.NumberNull()
    }
    if val, ok := dataMap["statusMessage"].(string); ok && val != "" {
        data.StatusMessage = types.StringValue(val)
    } else {
        data.StatusMessage = types.StringNull()
    }
    if val, ok := dataMap["name"].(string); ok && val != "" {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if val, ok := dataMap["kind"].(string); ok && val != "" {
        data.Kind = types.StringValue(val)
    } else {
        data.Kind = types.StringNull()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SpanResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var data SpanResourceModel
    var state SpanResourceModel

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
    spanRequest := map[string]interface{}{
        "data": map[string]interface{}{

        },
    }

    // Make API call
    httpResp, err := r.client.Put("/span/" + data.Id.ValueString() + "", spanRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update span, got error: %s", err))
        return
    }

    // Parse the update response
    var spanResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &spanResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse span response, got error: %s", err))
        return
    }

    // After successful update, fetch the current state by calling Read with select parameter
    selectParam := map[string]interface{}{
        "projectId": true,
        "serviceId": true,
        "startTime": true,
        "endTime": true,
        "startTimeUnixNano": true,
        "durationUnixNano": true,
        "endTimeUnixNano": true,
        "traceId": true,
        "spanId": true,
        "parentSpanId": true,
        "traceState": true,
        "attributes": true,
        "events": true,
        "links": true,
        "statusCode": true,
        "statusMessage": true,
        "name": true,
        "kind": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect("/span/" + data.Id.ValueString() + "", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read span after update, got error: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse span read response, got error: %s", err))
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
    if val, ok := dataMap["startTime"].(string); ok && val != "" {
        data.StartTime = types.StringValue(val)
    } else {
        data.StartTime = types.StringNull()
    }
    if val, ok := dataMap["endTime"].(string); ok && val != "" {
        data.EndTime = types.StringValue(val)
    } else {
        data.EndTime = types.StringNull()
    }
    if val, ok := dataMap["startTimeUnixNano"].(float64); ok {
        data.StartTimeUnixNano = types.NumberValue(big.NewFloat(val))
    } else if dataMap["startTimeUnixNano"] == nil {
        data.StartTimeUnixNano = types.NumberNull()
    }
    if val, ok := dataMap["durationUnixNano"].(float64); ok {
        data.DurationUnixNano = types.NumberValue(big.NewFloat(val))
    } else if dataMap["durationUnixNano"] == nil {
        data.DurationUnixNano = types.NumberNull()
    }
    if val, ok := dataMap["endTimeUnixNano"].(float64); ok {
        data.EndTimeUnixNano = types.NumberValue(big.NewFloat(val))
    } else if dataMap["endTimeUnixNano"] == nil {
        data.EndTimeUnixNano = types.NumberNull()
    }
    if val, ok := dataMap["traceId"].(string); ok && val != "" {
        data.TraceId = types.StringValue(val)
    } else {
        data.TraceId = types.StringNull()
    }
    if val, ok := dataMap["spanId"].(string); ok && val != "" {
        data.SpanId = types.StringValue(val)
    } else {
        data.SpanId = types.StringNull()
    }
    if val, ok := dataMap["parentSpanId"].(string); ok && val != "" {
        data.ParentSpanId = types.StringValue(val)
    } else {
        data.ParentSpanId = types.StringNull()
    }
    if val, ok := dataMap["traceState"].(string); ok && val != "" {
        data.TraceState = types.StringValue(val)
    } else {
        data.TraceState = types.StringNull()
    }
    if val, ok := dataMap["attributes"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.Attributes = mapValue
    } else if dataMap["attributes"] == nil {
        data.Attributes = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["events"].([]interface{}); ok {
        // Convert API response list to Terraform list
        listValue, _ := types.ListValueFrom(ctx, types.StringType, val)
        data.Events = listValue
    } else if dataMap["events"] == nil {
        data.Events = types.ListNull(types.StringType)
    }
    if val, ok := dataMap["links"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.Links = mapValue
    } else if dataMap["links"] == nil {
        data.Links = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["statusCode"].(float64); ok {
        data.StatusCode = types.NumberValue(big.NewFloat(val))
    } else if dataMap["statusCode"] == nil {
        data.StatusCode = types.NumberNull()
    }
    if val, ok := dataMap["statusMessage"].(string); ok && val != "" {
        data.StatusMessage = types.StringValue(val)
    } else {
        data.StatusMessage = types.StringNull()
    }
    if val, ok := dataMap["name"].(string); ok && val != "" {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if val, ok := dataMap["kind"].(string); ok && val != "" {
        data.Kind = types.StringValue(val)
    } else {
        data.Kind = types.StringNull()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SpanResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var data SpanResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Make API call
    _, err := r.client.Delete("/span/" + data.Id.ValueString() + "")
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete span, got error: %s", err))
        return
    }
}


func (r *SpanResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Helper method to convert Terraform map to Go interface{}
func (r *SpanResource) convertTerraformMapToInterface(terraformMap types.Map) interface{} {
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
func (r *SpanResource) convertTerraformListToInterface(terraformList types.List) interface{} {
    if terraformList.IsNull() || terraformList.IsUnknown() {
        return nil
    }
    
    var result []interface{}
    terraformList.ElementsAs(context.Background(), &result, false)
    return result
}
