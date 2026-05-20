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
var _ datasource.DataSource = &ProfileSampleDataDataSource{}

func NewProfileSampleDataDataSource() datasource.DataSource {
    return &ProfileSampleDataDataSource{}
}

// ProfileSampleDataDataSource defines the data source implementation.
type ProfileSampleDataDataSource struct {
    client *Client
}

// ProfileSampleDataDataSourceModel describes the data source data model.
type ProfileSampleDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    ProjectId types.String `tfsdk:"project_id"`
    ServiceId types.String `tfsdk:"service_id"`
    ServiceType types.String `tfsdk:"service_type"`
    ProfileId types.String `tfsdk:"profile_id"`
    TraceId types.String `tfsdk:"trace_id"`
    SpanId types.String `tfsdk:"span_id"`
    Time types.String `tfsdk:"time"`
    TimeUnixNano types.Number `tfsdk:"time_unix_nano"`
    Stacktrace types.Set `tfsdk:"stacktrace"`
    StacktraceHash types.String `tfsdk:"stacktrace_hash"`
    FrameTypes types.Set `tfsdk:"frame_types"`
    Value types.Number `tfsdk:"value"`
    ProfileType types.String `tfsdk:"profile_type"`
    Labels types.String `tfsdk:"labels"`
}

func (d *ProfileSampleDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_profile_sample_data"
}

func (d *ProfileSampleDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "profile_sample_data data source",

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
            "service_type": schema.StringAttribute{
                MarkdownDescription: "Service Type",
                Computed: true,
            },
            "profile_id": schema.StringAttribute{
                MarkdownDescription: "Profile ID",
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
            "time": schema.StringAttribute{
                MarkdownDescription: "Time",
                Computed: true,
            },
            "time_unix_nano": schema.NumberAttribute{
                MarkdownDescription: "Time (in Unix Nano)",
                Computed: true,
            },
            "stacktrace": schema.SetAttribute{
                MarkdownDescription: "Stacktrace",
                Computed: true,
                ElementType: types.StringType,
            },
            "stacktrace_hash": schema.StringAttribute{
                MarkdownDescription: "Stacktrace Hash",
                Computed: true,
            },
            "frame_types": schema.SetAttribute{
                MarkdownDescription: "Frame Types",
                Computed: true,
                ElementType: types.StringType,
            },
            "value": schema.NumberAttribute{
                MarkdownDescription: "Value",
                Computed: true,
            },
            "profile_type": schema.StringAttribute{
                MarkdownDescription: "Profile Type",
                Computed: true,
            },
            "labels": schema.StringAttribute{
                MarkdownDescription: "Labels",
                Computed: true,
            },
        },
    }
}

func (d *ProfileSampleDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ProfileSampleDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data ProfileSampleDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "profile-sample" + "/" + data.Id.ValueString()
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read profile_sample_data, got error: %s", err))
        return
    }

    var profileSampleDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &profileSampleDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse profile_sample_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := profileSampleDataResponse["data"].(map[string]interface{}); ok {
        profileSampleDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := profileSampleDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := profileSampleDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := profileSampleDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := profileSampleDataResponse["service_id"].(string); ok {
        data.ServiceId = types.StringValue(val)
    }
    if val, ok := profileSampleDataResponse["service_type"].(string); ok {
        data.ServiceType = types.StringValue(val)
    }
    if val, ok := profileSampleDataResponse["profile_id"].(string); ok {
        data.ProfileId = types.StringValue(val)
    }
    if val, ok := profileSampleDataResponse["trace_id"].(string); ok {
        data.TraceId = types.StringValue(val)
    }
    if val, ok := profileSampleDataResponse["span_id"].(string); ok {
        data.SpanId = types.StringValue(val)
    }
    if val, ok := profileSampleDataResponse["time"].(string); ok {
        data.Time = types.StringValue(val)
    }
    if val, ok := profileSampleDataResponse["time_unix_nano"].(float64); ok {
        data.TimeUnixNano = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := profileSampleDataResponse["stacktrace"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.Stacktrace = setValue
    }
    if val, ok := profileSampleDataResponse["stacktrace_hash"].(string); ok {
        data.StacktraceHash = types.StringValue(val)
    }
    if val, ok := profileSampleDataResponse["frame_types"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.FrameTypes = setValue
    }
    if val, ok := profileSampleDataResponse["value"].(float64); ok {
        data.Value = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := profileSampleDataResponse["profile_type"].(string); ok {
        data.ProfileType = types.StringValue(val)
    }
    if val, ok := profileSampleDataResponse["labels"].(string); ok {
        data.Labels = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
