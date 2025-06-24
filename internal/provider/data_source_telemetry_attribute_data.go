package provider

import (
    "context"
    "fmt"
    "github.com/hashicorp/terraform-plugin-framework/attr"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &TelemetryAttributeDataDataSource{}

func NewTelemetryAttributeDataDataSource() datasource.DataSource {
    return &TelemetryAttributeDataDataSource{}
}

// TelemetryAttributeDataDataSource defines the data source implementation.
type TelemetryAttributeDataDataSource struct {
    client *Client
}

// TelemetryAttributeDataDataSourceModel describes the data source data model.
type TelemetryAttributeDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    ProjectId types.String `tfsdk:"project_id"`
    TelemetryType types.String `tfsdk:"telemetry_type"`
    Attributes types.List `tfsdk:"attributes"`
}

func (d *TelemetryAttributeDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_telemetry_attribute_data"
}

func (d *TelemetryAttributeDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "telemetry_attribute_data data source",

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
            "telemetry_type": schema.StringAttribute{
                MarkdownDescription: "Telemetry Type",
                Computed: true,
            },
            "attributes": schema.ListAttribute{
                MarkdownDescription: "Attributes",
                Computed: true,
                ElementType: types.StringType,
            },
        },
    }
}

func (d *TelemetryAttributeDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *TelemetryAttributeDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data TelemetryAttributeDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "telemetry-attributes" + "/" + data.Id.ValueString()
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read telemetry_attribute_data, got error: %s", err))
        return
    }

    var telemetryAttributeDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &telemetryAttributeDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse telemetry_attribute_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := telemetryAttributeDataResponse["data"].(map[string]interface{}); ok {
        telemetryAttributeDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := telemetryAttributeDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := telemetryAttributeDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := telemetryAttributeDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := telemetryAttributeDataResponse["telemetry_type"].(string); ok {
        data.TelemetryType = types.StringValue(val)
    }
    if val, ok := telemetryAttributeDataResponse["attributes"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        listValue, _ := types.ListValue(types.StringType, elements)
        data.Attributes = listValue
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
