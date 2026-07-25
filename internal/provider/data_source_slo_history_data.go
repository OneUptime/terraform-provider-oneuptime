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
var _ datasource.DataSource = &SloHistoryDataDataSource{}

func NewSloHistoryDataDataSource() datasource.DataSource {
    return &SloHistoryDataDataSource{}
}

// SloHistoryDataDataSource defines the data source implementation.
type SloHistoryDataDataSource struct {
    client *Client
}

// SloHistoryDataDataSourceModel describes the data source data model.
type SloHistoryDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    ProjectId types.String `tfsdk:"project_id"`
    SloId types.String `tfsdk:"slo_id"`
    MetricName types.String `tfsdk:"metric_name"`
    BucketStart types.String `tfsdk:"bucket_start"`
    Value types.Number `tfsdk:"value"`
    Version types.String `tfsdk:"version"`
}

func (d *SloHistoryDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_slo_history_data"
}

func (d *SloHistoryDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "slo_history_data data source",

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
            "slo_id": schema.StringAttribute{
                MarkdownDescription: "SLO ID",
                Computed: true,
            },
            "metric_name": schema.StringAttribute{
                MarkdownDescription: "Metric Name",
                Computed: true,
            },
            "bucket_start": schema.StringAttribute{
                MarkdownDescription: "Bucket Start",
                Computed: true,
            },
            "value": schema.NumberAttribute{
                MarkdownDescription: "Value",
                Computed: true,
            },
            "version": schema.StringAttribute{
                MarkdownDescription: "Version",
                Computed: true,
            },
        },
    }
}

func (d *SloHistoryDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SloHistoryDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data SloHistoryDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "slo-history" + "/" + data.Id.ValueString()
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read slo_history_data, got error: %s", err))
        return
    }

    var sloHistoryDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &sloHistoryDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse slo_history_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := sloHistoryDataResponse["data"].(map[string]interface{}); ok {
        sloHistoryDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := sloHistoryDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := sloHistoryDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := sloHistoryDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := sloHistoryDataResponse["slo_id"].(string); ok {
        data.SloId = types.StringValue(val)
    }
    if val, ok := sloHistoryDataResponse["metric_name"].(string); ok {
        data.MetricName = types.StringValue(val)
    }
    if val, ok := sloHistoryDataResponse["bucket_start"].(string); ok {
        data.BucketStart = types.StringValue(val)
    }
    if val, ok := sloHistoryDataResponse["value"].(float64); ok {
        data.Value = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := sloHistoryDataResponse["version"].(string); ok {
        data.Version = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
