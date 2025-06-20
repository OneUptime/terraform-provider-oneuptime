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
var _ resource.Resource = &ExceptionInstanceResource{}
var _ resource.ResourceWithImportState = &ExceptionInstanceResource{}

func NewExceptionInstanceResource() resource.Resource {
    return &ExceptionInstanceResource{}
}

// ExceptionInstanceResource defines the resource implementation.
type ExceptionInstanceResource struct {
    client *Client
}

// ExceptionInstanceResourceModel describes the resource data model.
type ExceptionInstanceResourceModel struct {
    Id types.String `tfsdk:"id"`
    ProjectId types.String `tfsdk:"project_id"`
    ServiceId types.String `tfsdk:"service_id"`
    Time types.String `tfsdk:"time"`
    TimeUnixNano types.Number `tfsdk:"time_unix_nano"`
    ExceptionType types.String `tfsdk:"exception_type"`
    StackTrace types.String `tfsdk:"stack_trace"`
    Message types.String `tfsdk:"message"`
    SpanStatusCode types.Number `tfsdk:"span_status_code"`
    Escaped types.Bool `tfsdk:"escaped"`
    TraceId types.String `tfsdk:"trace_id"`
    SpanId types.String `tfsdk:"span_id"`
    Fingerprint types.String `tfsdk:"fingerprint"`
    SpanName types.String `tfsdk:"span_name"`
    Attributes types.Map `tfsdk:"attributes"`
}

func (r *ExceptionInstanceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_exception_instance"
}

func (r *ExceptionInstanceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "exception_instance resource",

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
            "time": schema.StringAttribute{
                MarkdownDescription: "Time",
                Computed: true,
            },
            "time_unix_nano": schema.NumberAttribute{
                MarkdownDescription: "Time (in Unix Nano)",
                Computed: true,
            },
            "exception_type": schema.StringAttribute{
                MarkdownDescription: "Exception Type",
                Computed: true,
            },
            "stack_trace": schema.StringAttribute{
                MarkdownDescription: "Stack Trace",
                Computed: true,
            },
            "message": schema.StringAttribute{
                MarkdownDescription: "Exception Message",
                Computed: true,
            },
            "span_status_code": schema.NumberAttribute{
                MarkdownDescription: "Span Status Code",
                Computed: true,
            },
            "escaped": schema.BoolAttribute{
                MarkdownDescription: "Exception Escaped",
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
            "fingerprint": schema.StringAttribute{
                MarkdownDescription: "Fingerprint",
                Computed: true,
            },
            "span_name": schema.StringAttribute{
                MarkdownDescription: "Span Name",
                Computed: true,
            },
            "attributes": schema.MapAttribute{
                MarkdownDescription: "Attributes",
                Computed: true,
                ElementType: types.StringType,
            },
        },
    }
}

func (r *ExceptionInstanceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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


func (r *ExceptionInstanceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var data ExceptionInstanceResourceModel

    // Read Terraform plan data into the model
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create API request body
    exceptionInstanceRequest := map[string]interface{}{
        "data": map[string]interface{}{

        },
    }

    // Make API call
    httpResp, err := r.client.Post("/exceptions", exceptionInstanceRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create exception_instance, got error: %s", err))
        return
    }

    var exceptionInstanceResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &exceptionInstanceResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse exception_instance response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := exceptionInstanceResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = exceptionInstanceResponse
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
    if val, ok := dataMap["time"].(string); ok && val != "" {
        data.Time = types.StringValue(val)
    } else {
        data.Time = types.StringNull()
    }
    if val, ok := dataMap["timeUnixNano"].(float64); ok {
        data.TimeUnixNano = types.NumberValue(big.NewFloat(val))
    } else if dataMap["timeUnixNano"] == nil {
        data.TimeUnixNano = types.NumberNull()
    }
    if val, ok := dataMap["exceptionType"].(string); ok && val != "" {
        data.ExceptionType = types.StringValue(val)
    } else {
        data.ExceptionType = types.StringNull()
    }
    if val, ok := dataMap["stackTrace"].(string); ok && val != "" {
        data.StackTrace = types.StringValue(val)
    } else {
        data.StackTrace = types.StringNull()
    }
    if val, ok := dataMap["message"].(string); ok && val != "" {
        data.Message = types.StringValue(val)
    } else {
        data.Message = types.StringNull()
    }
    if val, ok := dataMap["spanStatusCode"].(float64); ok {
        data.SpanStatusCode = types.NumberValue(big.NewFloat(val))
    } else if dataMap["spanStatusCode"] == nil {
        data.SpanStatusCode = types.NumberNull()
    }
    if val, ok := dataMap["escaped"].(bool); ok {
        data.Escaped = types.BoolValue(val)
    } else if dataMap["escaped"] == nil {
        data.Escaped = types.BoolNull()
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
    if val, ok := dataMap["fingerprint"].(string); ok && val != "" {
        data.Fingerprint = types.StringValue(val)
    } else {
        data.Fingerprint = types.StringNull()
    }
    if val, ok := dataMap["spanName"].(string); ok && val != "" {
        data.SpanName = types.StringValue(val)
    } else {
        data.SpanName = types.StringNull()
    }
    if val, ok := dataMap["attributes"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.Attributes = mapValue
    } else if dataMap["attributes"] == nil {
        data.Attributes = types.MapNull(types.StringType)
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

func (r *ExceptionInstanceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var data ExceptionInstanceResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create select parameter to get full object
    selectParam := map[string]interface{}{
        "projectId": true,
        "serviceId": true,
        "time": true,
        "timeUnixNano": true,
        "exceptionType": true,
        "stackTrace": true,
        "message": true,
        "spanStatusCode": true,
        "escaped": true,
        "traceId": true,
        "spanId": true,
        "fingerprint": true,
        "spanName": true,
        "attributes": true,
        "_id": true,
    }

    // Make API call with select parameter
    httpResp, err := r.client.PostWithSelect("/exceptions/" + data.Id.ValueString() + "", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read exception_instance, got error: %s", err))
        return
    }

    if httpResp.StatusCode == http.StatusNotFound {
        resp.State.RemoveResource(ctx)
        return
    }

    var exceptionInstanceResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &exceptionInstanceResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse exception_instance response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := exceptionInstanceResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = exceptionInstanceResponse
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
    if val, ok := dataMap["time"].(string); ok && val != "" {
        data.Time = types.StringValue(val)
    } else {
        data.Time = types.StringNull()
    }
    if val, ok := dataMap["timeUnixNano"].(float64); ok {
        data.TimeUnixNano = types.NumberValue(big.NewFloat(val))
    } else if dataMap["timeUnixNano"] == nil {
        data.TimeUnixNano = types.NumberNull()
    }
    if val, ok := dataMap["exceptionType"].(string); ok && val != "" {
        data.ExceptionType = types.StringValue(val)
    } else {
        data.ExceptionType = types.StringNull()
    }
    if val, ok := dataMap["stackTrace"].(string); ok && val != "" {
        data.StackTrace = types.StringValue(val)
    } else {
        data.StackTrace = types.StringNull()
    }
    if val, ok := dataMap["message"].(string); ok && val != "" {
        data.Message = types.StringValue(val)
    } else {
        data.Message = types.StringNull()
    }
    if val, ok := dataMap["spanStatusCode"].(float64); ok {
        data.SpanStatusCode = types.NumberValue(big.NewFloat(val))
    } else if dataMap["spanStatusCode"] == nil {
        data.SpanStatusCode = types.NumberNull()
    }
    if val, ok := dataMap["escaped"].(bool); ok {
        data.Escaped = types.BoolValue(val)
    } else if dataMap["escaped"] == nil {
        data.Escaped = types.BoolNull()
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
    if val, ok := dataMap["fingerprint"].(string); ok && val != "" {
        data.Fingerprint = types.StringValue(val)
    } else {
        data.Fingerprint = types.StringNull()
    }
    if val, ok := dataMap["spanName"].(string); ok && val != "" {
        data.SpanName = types.StringValue(val)
    } else {
        data.SpanName = types.StringNull()
    }
    if val, ok := dataMap["attributes"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.Attributes = mapValue
    } else if dataMap["attributes"] == nil {
        data.Attributes = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ExceptionInstanceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var data ExceptionInstanceResourceModel
    var state ExceptionInstanceResourceModel

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
    exceptionInstanceRequest := map[string]interface{}{
        "data": map[string]interface{}{

        },
    }

    // Make API call
    httpResp, err := r.client.Put("/exceptions/" + data.Id.ValueString() + "", exceptionInstanceRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update exception_instance, got error: %s", err))
        return
    }

    // Parse the update response
    var exceptionInstanceResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &exceptionInstanceResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse exception_instance response, got error: %s", err))
        return
    }

    // After successful update, fetch the current state by calling Read with select parameter
    selectParam := map[string]interface{}{
        "projectId": true,
        "serviceId": true,
        "time": true,
        "timeUnixNano": true,
        "exceptionType": true,
        "stackTrace": true,
        "message": true,
        "spanStatusCode": true,
        "escaped": true,
        "traceId": true,
        "spanId": true,
        "fingerprint": true,
        "spanName": true,
        "attributes": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect("/exceptions/" + data.Id.ValueString() + "", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read exception_instance after update, got error: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse exception_instance read response, got error: %s", err))
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
    if val, ok := dataMap["time"].(string); ok && val != "" {
        data.Time = types.StringValue(val)
    } else {
        data.Time = types.StringNull()
    }
    if val, ok := dataMap["timeUnixNano"].(float64); ok {
        data.TimeUnixNano = types.NumberValue(big.NewFloat(val))
    } else if dataMap["timeUnixNano"] == nil {
        data.TimeUnixNano = types.NumberNull()
    }
    if val, ok := dataMap["exceptionType"].(string); ok && val != "" {
        data.ExceptionType = types.StringValue(val)
    } else {
        data.ExceptionType = types.StringNull()
    }
    if val, ok := dataMap["stackTrace"].(string); ok && val != "" {
        data.StackTrace = types.StringValue(val)
    } else {
        data.StackTrace = types.StringNull()
    }
    if val, ok := dataMap["message"].(string); ok && val != "" {
        data.Message = types.StringValue(val)
    } else {
        data.Message = types.StringNull()
    }
    if val, ok := dataMap["spanStatusCode"].(float64); ok {
        data.SpanStatusCode = types.NumberValue(big.NewFloat(val))
    } else if dataMap["spanStatusCode"] == nil {
        data.SpanStatusCode = types.NumberNull()
    }
    if val, ok := dataMap["escaped"].(bool); ok {
        data.Escaped = types.BoolValue(val)
    } else if dataMap["escaped"] == nil {
        data.Escaped = types.BoolNull()
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
    if val, ok := dataMap["fingerprint"].(string); ok && val != "" {
        data.Fingerprint = types.StringValue(val)
    } else {
        data.Fingerprint = types.StringNull()
    }
    if val, ok := dataMap["spanName"].(string); ok && val != "" {
        data.SpanName = types.StringValue(val)
    } else {
        data.SpanName = types.StringNull()
    }
    if val, ok := dataMap["attributes"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.Attributes = mapValue
    } else if dataMap["attributes"] == nil {
        data.Attributes = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ExceptionInstanceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var data ExceptionInstanceResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Make API call
    _, err := r.client.Delete("/exceptions/" + data.Id.ValueString() + "")
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete exception_instance, got error: %s", err))
        return
    }
}


func (r *ExceptionInstanceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Helper method to convert Terraform map to Go interface{}
func (r *ExceptionInstanceResource) convertTerraformMapToInterface(terraformMap types.Map) interface{} {
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
func (r *ExceptionInstanceResource) convertTerraformListToInterface(terraformList types.List) interface{} {
    if terraformList.IsNull() || terraformList.IsUnknown() {
        return nil
    }
    
    var result []interface{}
    terraformList.ElementsAs(context.Background(), &result, false)
    return result
}
