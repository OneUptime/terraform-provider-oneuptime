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
var _ datasource.DataSource = &ExceptionInstanceDataDataSource{}

func NewExceptionInstanceDataDataSource() datasource.DataSource {
    return &ExceptionInstanceDataDataSource{}
}

// ExceptionInstanceDataDataSource defines the data source implementation.
type ExceptionInstanceDataDataSource struct {
    client *Client
}

// ExceptionInstanceDataDataSourceModel describes the data source data model.
type ExceptionInstanceDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
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
    Attributes types.String `tfsdk:"attributes"`
}

func (d *ExceptionInstanceDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_exception_instance_data"
}

func (d *ExceptionInstanceDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "exception_instance_data data source",

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
            "attributes": schema.StringAttribute{
                MarkdownDescription: "Attributes",
                Computed: true,
            },
        },
    }
}

func (d *ExceptionInstanceDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ExceptionInstanceDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data ExceptionInstanceDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "exceptions" + "/" + data.Id.ValueString()
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read exception_instance_data, got error: %s", err))
        return
    }

    var exceptionInstanceDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &exceptionInstanceDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse exception_instance_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := exceptionInstanceDataResponse["data"].(map[string]interface{}); ok {
        exceptionInstanceDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := exceptionInstanceDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := exceptionInstanceDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := exceptionInstanceDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := exceptionInstanceDataResponse["service_id"].(string); ok {
        data.ServiceId = types.StringValue(val)
    }
    if val, ok := exceptionInstanceDataResponse["time"].(string); ok {
        data.Time = types.StringValue(val)
    }
    if val, ok := exceptionInstanceDataResponse["time_unix_nano"].(float64); ok {
        data.TimeUnixNano = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := exceptionInstanceDataResponse["exception_type"].(string); ok {
        data.ExceptionType = types.StringValue(val)
    }
    if val, ok := exceptionInstanceDataResponse["stack_trace"].(string); ok {
        data.StackTrace = types.StringValue(val)
    }
    if val, ok := exceptionInstanceDataResponse["message"].(string); ok {
        data.Message = types.StringValue(val)
    }
    if val, ok := exceptionInstanceDataResponse["span_status_code"].(float64); ok {
        data.SpanStatusCode = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := exceptionInstanceDataResponse["escaped"].(bool); ok {
        data.Escaped = types.BoolValue(val)
    }
    if val, ok := exceptionInstanceDataResponse["trace_id"].(string); ok {
        data.TraceId = types.StringValue(val)
    }
    if val, ok := exceptionInstanceDataResponse["span_id"].(string); ok {
        data.SpanId = types.StringValue(val)
    }
    if val, ok := exceptionInstanceDataResponse["fingerprint"].(string); ok {
        data.Fingerprint = types.StringValue(val)
    }
    if val, ok := exceptionInstanceDataResponse["span_name"].(string); ok {
        data.SpanName = types.StringValue(val)
    }
    if val, ok := exceptionInstanceDataResponse["attributes"].(string); ok {
        data.Attributes = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
