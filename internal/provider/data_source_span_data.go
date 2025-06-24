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
var _ datasource.DataSource = &SpanDataDataSource{}

func NewSpanDataDataSource() datasource.DataSource {
    return &SpanDataDataSource{}
}

// SpanDataDataSource defines the data source implementation.
type SpanDataDataSource struct {
    client *Client
}

// SpanDataDataSourceModel describes the data source data model.
type SpanDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
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
    Attributes types.String `tfsdk:"attributes"`
    Events types.List `tfsdk:"events"`
    Links types.String `tfsdk:"links"`
    StatusCode types.Number `tfsdk:"status_code"`
    StatusMessage types.String `tfsdk:"status_message"`
    Kind types.String `tfsdk:"kind"`
}

func (d *SpanDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_span_data"
}

func (d *SpanDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "span_data data source",

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
            "attributes": schema.StringAttribute{
                MarkdownDescription: "Attributes",
                Computed: true,
            },
            "events": schema.ListAttribute{
                MarkdownDescription: "Events",
                Computed: true,
                ElementType: types.StringType,
            },
            "links": schema.StringAttribute{
                MarkdownDescription: "Links",
                Computed: true,
            },
            "status_code": schema.NumberAttribute{
                MarkdownDescription: "Status Code",
                Computed: true,
            },
            "status_message": schema.StringAttribute{
                MarkdownDescription: "Status Message",
                Computed: true,
            },
            "kind": schema.StringAttribute{
                MarkdownDescription: "Kind",
                Computed: true,
            },
        },
    }
}

func (d *SpanDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SpanDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data SpanDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "span" + "/" + data.Id.ValueString()
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read span_data, got error: %s", err))
        return
    }

    var spanDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &spanDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse span_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := spanDataResponse["data"].(map[string]interface{}); ok {
        spanDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := spanDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := spanDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := spanDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := spanDataResponse["service_id"].(string); ok {
        data.ServiceId = types.StringValue(val)
    }
    if val, ok := spanDataResponse["start_time"].(string); ok {
        data.StartTime = types.StringValue(val)
    }
    if val, ok := spanDataResponse["end_time"].(string); ok {
        data.EndTime = types.StringValue(val)
    }
    if val, ok := spanDataResponse["start_time_unix_nano"].(float64); ok {
        data.StartTimeUnixNano = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := spanDataResponse["duration_unix_nano"].(float64); ok {
        data.DurationUnixNano = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := spanDataResponse["end_time_unix_nano"].(float64); ok {
        data.EndTimeUnixNano = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := spanDataResponse["trace_id"].(string); ok {
        data.TraceId = types.StringValue(val)
    }
    if val, ok := spanDataResponse["span_id"].(string); ok {
        data.SpanId = types.StringValue(val)
    }
    if val, ok := spanDataResponse["parent_span_id"].(string); ok {
        data.ParentSpanId = types.StringValue(val)
    }
    if val, ok := spanDataResponse["trace_state"].(string); ok {
        data.TraceState = types.StringValue(val)
    }
    if val, ok := spanDataResponse["attributes"].(string); ok {
        data.Attributes = types.StringValue(val)
    }
    if val, ok := spanDataResponse["events"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        listValue, _ := types.ListValue(types.StringType, elements)
        data.Events = listValue
    }
    if val, ok := spanDataResponse["links"].(string); ok {
        data.Links = types.StringValue(val)
    }
    if val, ok := spanDataResponse["status_code"].(float64); ok {
        data.StatusCode = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := spanDataResponse["status_message"].(string); ok {
        data.StatusMessage = types.StringValue(val)
    }
    if val, ok := spanDataResponse["kind"].(string); ok {
        data.Kind = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
