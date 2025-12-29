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
var _ datasource.DataSource = &AiAgentTaskPullRequestDataDataSource{}

func NewAiAgentTaskPullRequestDataDataSource() datasource.DataSource {
    return &AiAgentTaskPullRequestDataDataSource{}
}

// AiAgentTaskPullRequestDataDataSource defines the data source implementation.
type AiAgentTaskPullRequestDataDataSource struct {
    client *Client
}

// AiAgentTaskPullRequestDataDataSourceModel describes the data source data model.
type AiAgentTaskPullRequestDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    AiAgentTaskId types.String `tfsdk:"ai_agent_task_id"`
    AiAgentId types.String `tfsdk:"ai_agent_id"`
    CodeRepositoryId types.String `tfsdk:"code_repository_id"`
    Title types.String `tfsdk:"title"`
    Description types.String `tfsdk:"description"`
    PullRequestUrl types.String `tfsdk:"pull_request_url"`
    PullRequestId types.Number `tfsdk:"pull_request_id"`
    PullRequestNumber types.Number `tfsdk:"pull_request_number"`
    PullRequestState types.String `tfsdk:"pull_request_state"`
    HeadRefName types.String `tfsdk:"head_ref_name"`
    BaseRefName types.String `tfsdk:"base_ref_name"`
    RepoOrganizationName types.String `tfsdk:"repo_organization_name"`
    RepoName types.String `tfsdk:"repo_name"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *AiAgentTaskPullRequestDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_ai_agent_task_pull_request_data"
}

func (d *AiAgentTaskPullRequestDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "ai_agent_task_pull_request_data data source",

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
            "code_repository_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "title": schema.StringAttribute{
                MarkdownDescription: "Title of the pull request.. Permissions - Create: [Project Owner, Project Admin, Project Member, Create AI Agent Task], Read: [Project Owner, Project Admin, Project Member, Read AI Agent Task], Update: [Project Owner, Project Admin, Project Member, Edit AI Agent Task]",
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Description/body of the pull request.. Permissions - Create: [Project Owner, Project Admin, Project Member, Create AI Agent Task], Read: [Project Owner, Project Admin, Project Member, Read AI Agent Task], Update: [Project Owner, Project Admin, Project Member, Edit AI Agent Task]",
                Computed: true,
            },
            "pull_request_url": schema.StringAttribute{
                MarkdownDescription: "URL to the pull request on the hosting platform.. Permissions - Create: [Project Owner, Project Admin, Project Member, Create AI Agent Task], Read: [Project Owner, Project Admin, Project Member, Read AI Agent Task], Update: [Project Owner, Project Admin, Project Member, Edit AI Agent Task]",
                Computed: true,
            },
            "pull_request_id": schema.NumberAttribute{
                MarkdownDescription: "The unique ID of the pull request from the hosting platform.. Permissions - Create: [Project Owner, Project Admin, Project Member, Create AI Agent Task], Read: [Project Owner, Project Admin, Project Member, Read AI Agent Task], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "pull_request_number": schema.NumberAttribute{
                MarkdownDescription: "The pull request number (e.g., #123).. Permissions - Create: [Project Owner, Project Admin, Project Member, Create AI Agent Task], Read: [Project Owner, Project Admin, Project Member, Read AI Agent Task], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "pull_request_state": schema.StringAttribute{
                MarkdownDescription: "Current state of the pull request (open, closed, merged).. Permissions - Create: [Project Owner, Project Admin, Project Member, Create AI Agent Task], Read: [Project Owner, Project Admin, Project Member, Read AI Agent Task], Update: [Project Owner, Project Admin, Project Member, Edit AI Agent Task]",
                Computed: true,
            },
            "head_ref_name": schema.StringAttribute{
                MarkdownDescription: "The branch name of the pull request (source branch).. Permissions - Create: [Project Owner, Project Admin, Project Member, Create AI Agent Task], Read: [Project Owner, Project Admin, Project Member, Read AI Agent Task], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "base_ref_name": schema.StringAttribute{
                MarkdownDescription: "The target branch for the pull request.. Permissions - Create: [Project Owner, Project Admin, Project Member, Create AI Agent Task], Read: [Project Owner, Project Admin, Project Member, Read AI Agent Task], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "repo_organization_name": schema.StringAttribute{
                MarkdownDescription: "Organization or username that owns the repository.. Permissions - Create: [Project Owner, Project Admin, Project Member, Create AI Agent Task], Read: [Project Owner, Project Admin, Project Member, Read AI Agent Task], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "repo_name": schema.StringAttribute{
                MarkdownDescription: "Name of the repository.. Permissions - Create: [Project Owner, Project Admin, Project Member, Create AI Agent Task], Read: [Project Owner, Project Admin, Project Member, Read AI Agent Task], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
        },
    }
}

func (d *AiAgentTaskPullRequestDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AiAgentTaskPullRequestDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data AiAgentTaskPullRequestDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "ai-agent-task-pull-request" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read ai_agent_task_pull_request_data, got error: %s", err))
        return
    }

    var aiAgentTaskPullRequestDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &aiAgentTaskPullRequestDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse ai_agent_task_pull_request_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := aiAgentTaskPullRequestDataResponse["data"].(map[string]interface{}); ok {
        aiAgentTaskPullRequestDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := aiAgentTaskPullRequestDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := aiAgentTaskPullRequestDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := aiAgentTaskPullRequestDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := aiAgentTaskPullRequestDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := aiAgentTaskPullRequestDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := aiAgentTaskPullRequestDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := aiAgentTaskPullRequestDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := aiAgentTaskPullRequestDataResponse["ai_agent_task_id"].(string); ok {
        data.AiAgentTaskId = types.StringValue(val)
    }
    if val, ok := aiAgentTaskPullRequestDataResponse["ai_agent_id"].(string); ok {
        data.AiAgentId = types.StringValue(val)
    }
    if val, ok := aiAgentTaskPullRequestDataResponse["code_repository_id"].(string); ok {
        data.CodeRepositoryId = types.StringValue(val)
    }
    if val, ok := aiAgentTaskPullRequestDataResponse["title"].(string); ok {
        data.Title = types.StringValue(val)
    }
    if val, ok := aiAgentTaskPullRequestDataResponse["description"].(string); ok {
        data.Description = types.StringValue(val)
    }
    if val, ok := aiAgentTaskPullRequestDataResponse["pull_request_url"].(string); ok {
        data.PullRequestUrl = types.StringValue(val)
    }
    if val, ok := aiAgentTaskPullRequestDataResponse["pull_request_id"].(float64); ok {
        data.PullRequestId = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := aiAgentTaskPullRequestDataResponse["pull_request_number"].(float64); ok {
        data.PullRequestNumber = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := aiAgentTaskPullRequestDataResponse["pull_request_state"].(string); ok {
        data.PullRequestState = types.StringValue(val)
    }
    if val, ok := aiAgentTaskPullRequestDataResponse["head_ref_name"].(string); ok {
        data.HeadRefName = types.StringValue(val)
    }
    if val, ok := aiAgentTaskPullRequestDataResponse["base_ref_name"].(string); ok {
        data.BaseRefName = types.StringValue(val)
    }
    if val, ok := aiAgentTaskPullRequestDataResponse["repo_organization_name"].(string); ok {
        data.RepoOrganizationName = types.StringValue(val)
    }
    if val, ok := aiAgentTaskPullRequestDataResponse["repo_name"].(string); ok {
        data.RepoName = types.StringValue(val)
    }
    if val, ok := aiAgentTaskPullRequestDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
