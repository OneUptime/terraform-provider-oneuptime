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
var _ datasource.DataSource = &AiRunEventDataDataSource{}

func NewAiRunEventDataDataSource() datasource.DataSource {
    return &AiRunEventDataDataSource{}
}

// AiRunEventDataDataSource defines the data source implementation.
type AiRunEventDataDataSource struct {
    client *Client
}

// AiRunEventDataDataSourceModel describes the data source data model.
type AiRunEventDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    AiRunId types.String `tfsdk:"ai_run_id"`
    UserId types.String `tfsdk:"user_id"`
    Sequence types.Number `tfsdk:"sequence"`
    EventType types.String `tfsdk:"event_type"`
    ToolName types.String `tfsdk:"tool_name"`
    ToolArguments types.String `tfsdk:"tool_arguments"`
    ResultSummary types.String `tfsdk:"result_summary"`
    CitationId types.String `tfsdk:"citation_id"`
}

func (d *AiRunEventDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_ai_run_event_data"
}

func (d *AiRunEventDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "ai_run_event_data data source",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Identifier to filter by",
                Optional: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Name to filter by",
                Optional: true,
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
            "ai_run_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "sequence": schema.NumberAttribute{
                MarkdownDescription: "Order of this event within the run.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "event_type": schema.StringAttribute{
                MarkdownDescription: "Type of event.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "tool_name": schema.StringAttribute{
                MarkdownDescription: "Name of the tool for tool-call events.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "tool_arguments": schema.StringAttribute{
                MarkdownDescription: "Validated tool arguments as executed.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "result_summary": schema.StringAttribute{
                MarkdownDescription: "Summary of the result: row count, duration, truncation and bytes sent to the LLM.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "citation_id": schema.StringAttribute{
                MarkdownDescription: "ID of the citation this event minted (e.g. C1), if it produced one.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
        },
    }
}

func (d *AiRunEventDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AiRunEventDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data AiRunEventDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "ai-run-event" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read ai_run_event_data, got error: %s", err))
        return
    }

    var aiRunEventDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &aiRunEventDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse ai_run_event_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := aiRunEventDataResponse["data"].(map[string]interface{}); ok {
        aiRunEventDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := aiRunEventDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := aiRunEventDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := aiRunEventDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := aiRunEventDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := aiRunEventDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := aiRunEventDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := aiRunEventDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := aiRunEventDataResponse["ai_run_id"].(string); ok {
        data.AiRunId = types.StringValue(val)
    }
    if val, ok := aiRunEventDataResponse["user_id"].(string); ok {
        data.UserId = types.StringValue(val)
    }
    if val, ok := aiRunEventDataResponse["sequence"].(float64); ok {
        data.Sequence = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := aiRunEventDataResponse["event_type"].(string); ok {
        data.EventType = types.StringValue(val)
    }
    if val, ok := aiRunEventDataResponse["tool_name"].(string); ok {
        data.ToolName = types.StringValue(val)
    }
    if val, ok := aiRunEventDataResponse["tool_arguments"].(string); ok {
        data.ToolArguments = types.StringValue(val)
    }
    if val, ok := aiRunEventDataResponse["result_summary"].(string); ok {
        data.ResultSummary = types.StringValue(val)
    }
    if val, ok := aiRunEventDataResponse["citation_id"].(string); ok {
        data.CitationId = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
