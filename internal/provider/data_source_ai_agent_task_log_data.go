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
var _ datasource.DataSource = &AiAgentTaskLogDataDataSource{}

func NewAiAgentTaskLogDataDataSource() datasource.DataSource {
    return &AiAgentTaskLogDataDataSource{}
}

// AiAgentTaskLogDataDataSource defines the data source implementation.
type AiAgentTaskLogDataDataSource struct {
    client *Client
}

// AiAgentTaskLogDataDataSourceModel describes the data source data model.
type AiAgentTaskLogDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    AiAgentTaskId types.String `tfsdk:"ai_agent_task_id"`
    AiAgentId types.String `tfsdk:"ai_agent_id"`
    Severity types.String `tfsdk:"severity"`
    Message types.String `tfsdk:"message"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *AiAgentTaskLogDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_ai_agent_task_log_data"
}

func (d *AiAgentTaskLogDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "ai_agent_task_log_data data source",

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
            "ai_agent_task_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "ai_agent_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "severity": schema.StringAttribute{
                MarkdownDescription: "Severity level of this log entry (e.g., Information, Warning, Error).. Permissions - Create: [Project Owner, Project Admin, Project Member, Create AI Agent Task], Read: [Project Owner, Project Admin, Project Member, Viewer, Read AI Agent Task, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Edit AI Agent Task]",
                Computed: true,
            },
            "message": schema.StringAttribute{
                MarkdownDescription: "The log message content.. Permissions - Create: [Project Owner, Project Admin, Project Member, Create AI Agent Task], Read: [Project Owner, Project Admin, Project Member, Viewer, Read AI Agent Task, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Edit AI Agent Task]",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
        },
    }
}

func (d *AiAgentTaskLogDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AiAgentTaskLogDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data AiAgentTaskLogDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "ai-agent-task-log" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read ai_agent_task_log_data, got error: %s", err))
        return
    }

    var aiAgentTaskLogDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &aiAgentTaskLogDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse ai_agent_task_log_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := aiAgentTaskLogDataResponse["data"].(map[string]interface{}); ok {
        aiAgentTaskLogDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := aiAgentTaskLogDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := aiAgentTaskLogDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := aiAgentTaskLogDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := aiAgentTaskLogDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := aiAgentTaskLogDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := aiAgentTaskLogDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := aiAgentTaskLogDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := aiAgentTaskLogDataResponse["ai_agent_task_id"].(string); ok {
        data.AiAgentTaskId = types.StringValue(val)
    }
    if val, ok := aiAgentTaskLogDataResponse["ai_agent_id"].(string); ok {
        data.AiAgentId = types.StringValue(val)
    }
    if val, ok := aiAgentTaskLogDataResponse["severity"].(string); ok {
        data.Severity = types.StringValue(val)
    }
    if val, ok := aiAgentTaskLogDataResponse["message"].(string); ok {
        data.Message = types.StringValue(val)
    }
    if val, ok := aiAgentTaskLogDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
