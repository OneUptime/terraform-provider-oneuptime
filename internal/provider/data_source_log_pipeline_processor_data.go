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
var _ datasource.DataSource = &LogPipelineProcessorDataDataSource{}

func NewLogPipelineProcessorDataDataSource() datasource.DataSource {
    return &LogPipelineProcessorDataDataSource{}
}

// LogPipelineProcessorDataDataSource defines the data source implementation.
type LogPipelineProcessorDataDataSource struct {
    client *Client
}

// LogPipelineProcessorDataDataSourceModel describes the data source data model.
type LogPipelineProcessorDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    LogPipelineId types.String `tfsdk:"log_pipeline_id"`
    ProcessorType types.String `tfsdk:"processor_type"`
    Configuration types.String `tfsdk:"configuration"`
    IsEnabled types.Bool `tfsdk:"is_enabled"`
    SortOrder types.Number `tfsdk:"sort_order"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *LogPipelineProcessorDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_log_pipeline_processor_data"
}

func (d *LogPipelineProcessorDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "log_pipeline_processor_data data source",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Identifier to filter by",
                Optional: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Name to filter by",
                Computed: true,
            },
            "created_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "updated_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "deleted_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "version": schema.NumberAttribute{
                MarkdownDescription: "Object version",
                Computed: true,
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "log_pipeline_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "processor_type": schema.StringAttribute{
                MarkdownDescription: "The type of processor: GrokParser, AttributeRemapper, SeverityRemapper, or CategoryProcessor.. Permissions - Create: [Project Owner, Project Admin, Create Log Pipeline Processor], Read: [Project Owner, Project Admin, Project Member, Read Log Pipeline Processor, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Log Pipeline Processor]",
                Computed: true,
            },
            "configuration": schema.StringAttribute{
                MarkdownDescription: "Processor-specific configuration as JSON (e.g., grok pattern, source/target fields, mapping rules).. Permissions - Create: [Project Owner, Project Admin, Create Log Pipeline Processor], Read: [Project Owner, Project Admin, Project Member, Read Log Pipeline Processor, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Log Pipeline Processor]",
                Computed: true,
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether this processor is active.. Permissions - Create: [Project Owner, Project Admin, Create Log Pipeline Processor], Read: [Project Owner, Project Admin, Project Member, Read Log Pipeline Processor, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Log Pipeline Processor]",
                Computed: true,
            },
            "sort_order": schema.NumberAttribute{
                MarkdownDescription: "Determines the execution order of this processor within its pipeline.. Permissions - Create: [Project Owner, Project Admin, Create Log Pipeline Processor], Read: [Project Owner, Project Admin, Project Member, Read Log Pipeline Processor, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Log Pipeline Processor]",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
        },
    }
}

func (d *LogPipelineProcessorDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *LogPipelineProcessorDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data LogPipelineProcessorDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "log-pipeline-processor" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read log_pipeline_processor_data, got error: %s", err))
        return
    }

    var logPipelineProcessorDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &logPipelineProcessorDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse log_pipeline_processor_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := logPipelineProcessorDataResponse["data"].(map[string]interface{}); ok {
        logPipelineProcessorDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := logPipelineProcessorDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := logPipelineProcessorDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := logPipelineProcessorDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := logPipelineProcessorDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := logPipelineProcessorDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := logPipelineProcessorDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := logPipelineProcessorDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := logPipelineProcessorDataResponse["log_pipeline_id"].(string); ok {
        data.LogPipelineId = types.StringValue(val)
    }
    if val, ok := logPipelineProcessorDataResponse["processor_type"].(string); ok {
        data.ProcessorType = types.StringValue(val)
    }
    if val, ok := logPipelineProcessorDataResponse["configuration"].(string); ok {
        data.Configuration = types.StringValue(val)
    }
    if val, ok := logPipelineProcessorDataResponse["is_enabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    }
    if val, ok := logPipelineProcessorDataResponse["sort_order"].(float64); ok {
        data.SortOrder = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := logPipelineProcessorDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
