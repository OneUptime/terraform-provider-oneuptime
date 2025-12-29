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
var _ datasource.DataSource = &AiAgentTaskDataDataSource{}

func NewAiAgentTaskDataDataSource() datasource.DataSource {
    return &AiAgentTaskDataDataSource{}
}

// AiAgentTaskDataDataSource defines the data source implementation.
type AiAgentTaskDataDataSource struct {
    client *Client
}

// AiAgentTaskDataDataSourceModel describes the data source data model.
type AiAgentTaskDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    Description types.String `tfsdk:"description"`
    AiAgentId types.String `tfsdk:"ai_agent_id"`
    TaskType types.String `tfsdk:"task_type"`
    Status types.String `tfsdk:"status"`
    StatusMessage types.String `tfsdk:"status_message"`
    Metadata types.String `tfsdk:"metadata"`
    StartedAt types.String `tfsdk:"started_at"`
    CompletedAt types.String `tfsdk:"completed_at"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    TaskNumber types.Number `tfsdk:"task_number"`
}

func (d *AiAgentTaskDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_ai_agent_task_data"
}

func (d *AiAgentTaskDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "ai_agent_task_data data source",

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
            "description": schema.StringAttribute{
                MarkdownDescription: "Description of the AI Agent Task.. Permissions - Create: [Project Owner, Project Admin, Project Member, Create AI Agent Task], Read: [Project Owner, Project Admin, Project Member, Read AI Agent Task], Update: [Project Owner, Project Admin, Project Member, Edit AI Agent Task]",
                Computed: true,
            },
            "ai_agent_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "task_type": schema.StringAttribute{
                MarkdownDescription: "Type of task to be performed by the AI agent.. Permissions - Create: [Project Owner, Project Admin, Project Member, Create AI Agent Task], Read: [Project Owner, Project Admin, Project Member, Read AI Agent Task], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "status": schema.StringAttribute{
                MarkdownDescription: "Current status of the task.. Permissions - Create: [Project Owner, Project Admin, Project Member, Create AI Agent Task], Read: [Project Owner, Project Admin, Project Member, Read AI Agent Task], Update: [Project Owner, Project Admin, Project Member, Edit AI Agent Task]",
                Computed: true,
            },
            "status_message": schema.StringAttribute{
                MarkdownDescription: "A message describing the current status or result of the task.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read AI Agent Task], Update: [Project Owner, Project Admin, Project Member, Edit AI Agent Task]",
                Computed: true,
            },
            "metadata": schema.StringAttribute{
                MarkdownDescription: "Task-specific metadata containing context for the AI agent. Structure varies based on task type.. Permissions - Create: [Project Owner, Project Admin, Project Member, Create AI Agent Task], Read: [Project Owner, Project Admin, Project Member, Read AI Agent Task], Update: [Project Owner, Project Admin, Project Member, Edit AI Agent Task]",
                Computed: true,
            },
            "started_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "completed_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "task_number": schema.NumberAttribute{
                MarkdownDescription: "A unique, sequential number assigned to each AI Agent Task within a project.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read AI Agent Task], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
        },
    }
}

func (d *AiAgentTaskDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AiAgentTaskDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data AiAgentTaskDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "ai-agent-task" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read ai_agent_task_data, got error: %s", err))
        return
    }

    var aiAgentTaskDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &aiAgentTaskDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse ai_agent_task_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := aiAgentTaskDataResponse["data"].(map[string]interface{}); ok {
        aiAgentTaskDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := aiAgentTaskDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := aiAgentTaskDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := aiAgentTaskDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := aiAgentTaskDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := aiAgentTaskDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := aiAgentTaskDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := aiAgentTaskDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := aiAgentTaskDataResponse["description"].(string); ok {
        data.Description = types.StringValue(val)
    }
    if val, ok := aiAgentTaskDataResponse["ai_agent_id"].(string); ok {
        data.AiAgentId = types.StringValue(val)
    }
    if val, ok := aiAgentTaskDataResponse["task_type"].(string); ok {
        data.TaskType = types.StringValue(val)
    }
    if val, ok := aiAgentTaskDataResponse["status"].(string); ok {
        data.Status = types.StringValue(val)
    }
    if val, ok := aiAgentTaskDataResponse["status_message"].(string); ok {
        data.StatusMessage = types.StringValue(val)
    }
    if val, ok := aiAgentTaskDataResponse["metadata"].(string); ok {
        data.Metadata = types.StringValue(val)
    }
    if val, ok := aiAgentTaskDataResponse["started_at"].(string); ok {
        data.StartedAt = types.StringValue(val)
    }
    if val, ok := aiAgentTaskDataResponse["completed_at"].(string); ok {
        data.CompletedAt = types.StringValue(val)
    }
    if val, ok := aiAgentTaskDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := aiAgentTaskDataResponse["task_number"].(float64); ok {
        data.TaskNumber = types.NumberValue(big.NewFloat(val))
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
