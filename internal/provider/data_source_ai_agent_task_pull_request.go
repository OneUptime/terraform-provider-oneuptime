package provider

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "math/big"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &AiAgentTaskPullRequestDataSource{}

func NewAiAgentTaskPullRequestDataSource() datasource.DataSource {
    return &AiAgentTaskPullRequestDataSource{}
}

// AiAgentTaskPullRequestDataSource defines the data source implementation.
type AiAgentTaskPullRequestDataSource struct {
    client *Client
}

// AiAgentTaskPullRequestDataSourceModel describes the data source data model.
type AiAgentTaskPullRequestDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    AiRunId types.String `tfsdk:"ai_run_id"`
    AiAgentId types.String `tfsdk:"ai_agent_id"`
    CodeRepositoryId types.String `tfsdk:"code_repository_id"`
    Title types.String `tfsdk:"title"`
    Description types.String `tfsdk:"description"`
    PullRequestUrl types.String `tfsdk:"pull_request_url"`
    PullRequestId types.Number `tfsdk:"pull_request_id"`
    PullRequestNumber types.Number `tfsdk:"pull_request_number"`
    PullRequestState types.String `tfsdk:"pull_request_state"`
    CiStatus types.String `tfsdk:"ci_status"`
    CiStatusAt types.String `tfsdk:"ci_status_at"`
    RunnerVerificationStatus types.String `tfsdk:"runner_verification_status"`
    RunnerVerificationSummary types.String `tfsdk:"runner_verification_summary"`
    HeadRefName types.String `tfsdk:"head_ref_name"`
    BaseRefName types.String `tfsdk:"base_ref_name"`
    RepoOrganizationName types.String `tfsdk:"repo_organization_name"`
    RepoName types.String `tfsdk:"repo_name"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *AiAgentTaskPullRequestDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_ai_agent_task_pull_request"
}

func (d *AiAgentTaskPullRequestDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Pull requests created by AI agents during task execution. Look up an existing ai_agent_task_pull_request by `id` or by `name`.",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Look up by unique identifier. Exactly one of `id` or `name` must be set.",
                Optional: true,
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.",
                Optional: true,
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
            "ai_run_id": schema.StringAttribute{
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
                MarkdownDescription: "Title of the pull request..",
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Description/body of the pull request..",
                Computed: true,
            },
            "pull_request_url": schema.StringAttribute{
                MarkdownDescription: "URL to the pull request on the hosting platform..",
                Computed: true,
            },
            "pull_request_id": schema.NumberAttribute{
                MarkdownDescription: "The unique ID of the pull request from the hosting platform..",
                Computed: true,
            },
            "pull_request_number": schema.NumberAttribute{
                MarkdownDescription: "The pull request number (e.g., #123)..",
                Computed: true,
            },
            "pull_request_state": schema.StringAttribute{
                MarkdownDescription: "Current state of the pull request (open, closed, merged)..",
                Computed: true,
            },
            "ci_status": schema.StringAttribute{
                MarkdownDescription: "Rolled-up conclusion of the repository's own CI check runs on this pull request (Pending, Green, Red, ExpectedFailureObserved for should-fail regression-test PRs, NoCiConfigured). Null until the sync job first polls check runs. Written by AIAgent:SyncPullRequestStates — never by users..",
                Computed: true,
            },
            "ci_status_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "runner_verification_status": schema.StringAttribute{
                MarkdownDescription: "Outcome of the Runner-side build/test verification that ran against the fix BEFORE this pull request opened (Passed, Failed, Skipped when the repository has no verification commands configured). Distinct from CI Status, which mirrors the repository's own CI checks after the PR exists. Written by the Runner at record time — never by users..",
                Computed: true,
            },
            "runner_verification_summary": schema.StringAttribute{
                MarkdownDescription: "Human-readable summary of the Runner-side verification (which commands ran, what failed, how many repair attempts were used). Written by the Runner at record time..",
                Computed: true,
            },
            "head_ref_name": schema.StringAttribute{
                MarkdownDescription: "The branch name of the pull request (source branch)..",
                Computed: true,
            },
            "base_ref_name": schema.StringAttribute{
                MarkdownDescription: "The target branch for the pull request..",
                Computed: true,
            },
            "repo_organization_name": schema.StringAttribute{
                MarkdownDescription: "Organization or username that owns the repository..",
                Computed: true,
            },
            "repo_name": schema.StringAttribute{
                MarkdownDescription: "Name of the repository..",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
        },
    }
}

func (d *AiAgentTaskPullRequestDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AiAgentTaskPullRequestDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data AiAgentTaskPullRequestDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && data.Id.ValueString() != ""
    hasName := !data.Name.IsNull() && data.Name.ValueString() != ""
    if hasId == hasName {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Exactly one of `id` or `name` must be set to look up a ai_agent_task_pull_request.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "name": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "projectId": true,
        "aiRunId": true,
        "aiAgentId": true,
        "codeRepositoryId": true,
        "title": true,
        "description": true,
        "pullRequestUrl": true,
        "pullRequestId": true,
        "pullRequestNumber": true,
        "pullRequestState": true,
        "ciStatus": true,
        "ciStatusAt": true,
        "runnerVerificationStatus": true,
        "runnerVerificationSummary": true,
        "headRefName": true,
        "baseRefName": true,
        "repoOrganizationName": true,
        "repoName": true,
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/ai-agent-task-pull-request/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read ai_agent_task_pull_request, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No ai_agent_task_pull_request found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read ai_agent_task_pull_request: %s", err))
            return
        }
        if wrapper, ok := itemResponse["data"].(map[string]interface{}); ok {
            item = wrapper
        } else {
            item = itemResponse
        }
    } else {
        listBody := map[string]interface{}{
            "query": map[string]interface{}{
                "name": data.Name.ValueString(),
            },
            "select": selectParam,
            // limit 2 is enough to detect ambiguity without paging.
            "limit": 2,
        }
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/ai-agent-task-pull-request/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list ai_agent_task_pull_request, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list ai_agent_task_pull_request: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No ai_agent_task_pull_request found with name %q.", data.Name.ValueString()))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one ai_agent_task_pull_request matches name %q. Use the id attribute to disambiguate.", data.Name.ValueString()))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for ai_agent_task_pull_request.")
            return
        }
        item = first
    }

    // Update the model with response data
    if obj, ok := item["_id"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Id = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Id = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Id = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Id = types.StringValue(string(jsonBytes))
        } else {
            data.Id = types.StringNull()
        }
    } else if val, ok := item["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    if obj, ok := item["name"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Name = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Name = types.StringValue(string(jsonBytes))
        } else {
            data.Name = types.StringNull()
        }
    } else if val, ok := item["name"].(string); ok {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if obj, ok := item["createdAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CreatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CreatedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CreatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedAt = types.StringNull()
        }
    } else if val, ok := item["createdAt"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    } else {
        data.CreatedAt = types.StringNull()
    }
    if obj, ok := item["updatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UpdatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.UpdatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.UpdatedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.UpdatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.UpdatedAt = types.StringNull()
        }
    } else if val, ok := item["updatedAt"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    } else {
        data.UpdatedAt = types.StringNull()
    }
    if obj, ok := item["deletedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DeletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DeletedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DeletedAt = types.StringValue(string(jsonBytes))
        } else {
            data.DeletedAt = types.StringNull()
        }
    } else if val, ok := item["deletedAt"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    } else {
        data.DeletedAt = types.StringNull()
    }
    if val, ok := item["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["version"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.Version = types.NumberValue(big.NewFloat(val))
        } else {
            data.Version = types.NumberNull()
        }
    } else {
        data.Version = types.NumberNull()
    }
    if obj, ok := item["projectId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ProjectId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ProjectId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ProjectId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ProjectId = types.StringValue(string(jsonBytes))
        } else {
            data.ProjectId = types.StringNull()
        }
    } else if val, ok := item["projectId"].(string); ok {
        data.ProjectId = types.StringValue(val)
    } else {
        data.ProjectId = types.StringNull()
    }
    if obj, ok := item["aiRunId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AiRunId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AiRunId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AiRunId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AiRunId = types.StringValue(string(jsonBytes))
        } else {
            data.AiRunId = types.StringNull()
        }
    } else if val, ok := item["aiRunId"].(string); ok {
        data.AiRunId = types.StringValue(val)
    } else {
        data.AiRunId = types.StringNull()
    }
    if obj, ok := item["aiAgentId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AiAgentId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AiAgentId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AiAgentId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AiAgentId = types.StringValue(string(jsonBytes))
        } else {
            data.AiAgentId = types.StringNull()
        }
    } else if val, ok := item["aiAgentId"].(string); ok {
        data.AiAgentId = types.StringValue(val)
    } else {
        data.AiAgentId = types.StringNull()
    }
    if obj, ok := item["codeRepositoryId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CodeRepositoryId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CodeRepositoryId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CodeRepositoryId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CodeRepositoryId = types.StringValue(string(jsonBytes))
        } else {
            data.CodeRepositoryId = types.StringNull()
        }
    } else if val, ok := item["codeRepositoryId"].(string); ok {
        data.CodeRepositoryId = types.StringValue(val)
    } else {
        data.CodeRepositoryId = types.StringNull()
    }
    if obj, ok := item["title"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Title = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Title = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Title = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Title = types.StringValue(string(jsonBytes))
        } else {
            data.Title = types.StringNull()
        }
    } else if val, ok := item["title"].(string); ok {
        data.Title = types.StringValue(val)
    } else {
        data.Title = types.StringNull()
    }
    if obj, ok := item["description"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Description = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Description = types.StringValue(string(jsonBytes))
        } else {
            data.Description = types.StringNull()
        }
    } else if val, ok := item["description"].(string); ok {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if obj, ok := item["pullRequestUrl"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PullRequestUrl = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.PullRequestUrl = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.PullRequestUrl = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.PullRequestUrl = types.StringValue(string(jsonBytes))
        } else {
            data.PullRequestUrl = types.StringNull()
        }
    } else if val, ok := item["pullRequestUrl"].(string); ok {
        data.PullRequestUrl = types.StringValue(val)
    } else {
        data.PullRequestUrl = types.StringNull()
    }
    if val, ok := item["pullRequestId"].(float64); ok {
        data.PullRequestId = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["pullRequestId"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.PullRequestId = types.NumberValue(big.NewFloat(val))
        } else {
            data.PullRequestId = types.NumberNull()
        }
    } else {
        data.PullRequestId = types.NumberNull()
    }
    if val, ok := item["pullRequestNumber"].(float64); ok {
        data.PullRequestNumber = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["pullRequestNumber"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.PullRequestNumber = types.NumberValue(big.NewFloat(val))
        } else {
            data.PullRequestNumber = types.NumberNull()
        }
    } else {
        data.PullRequestNumber = types.NumberNull()
    }
    if obj, ok := item["pullRequestState"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PullRequestState = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.PullRequestState = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.PullRequestState = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.PullRequestState = types.StringValue(string(jsonBytes))
        } else {
            data.PullRequestState = types.StringNull()
        }
    } else if val, ok := item["pullRequestState"].(string); ok {
        data.PullRequestState = types.StringValue(val)
    } else {
        data.PullRequestState = types.StringNull()
    }
    if obj, ok := item["ciStatus"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CiStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CiStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CiStatus = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CiStatus = types.StringValue(string(jsonBytes))
        } else {
            data.CiStatus = types.StringNull()
        }
    } else if val, ok := item["ciStatus"].(string); ok {
        data.CiStatus = types.StringValue(val)
    } else {
        data.CiStatus = types.StringNull()
    }
    if obj, ok := item["ciStatusAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CiStatusAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CiStatusAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CiStatusAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CiStatusAt = types.StringValue(string(jsonBytes))
        } else {
            data.CiStatusAt = types.StringNull()
        }
    } else if val, ok := item["ciStatusAt"].(string); ok {
        data.CiStatusAt = types.StringValue(val)
    } else {
        data.CiStatusAt = types.StringNull()
    }
    if obj, ok := item["runnerVerificationStatus"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RunnerVerificationStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RunnerVerificationStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RunnerVerificationStatus = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RunnerVerificationStatus = types.StringValue(string(jsonBytes))
        } else {
            data.RunnerVerificationStatus = types.StringNull()
        }
    } else if val, ok := item["runnerVerificationStatus"].(string); ok {
        data.RunnerVerificationStatus = types.StringValue(val)
    } else {
        data.RunnerVerificationStatus = types.StringNull()
    }
    if obj, ok := item["runnerVerificationSummary"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RunnerVerificationSummary = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RunnerVerificationSummary = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RunnerVerificationSummary = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RunnerVerificationSummary = types.StringValue(string(jsonBytes))
        } else {
            data.RunnerVerificationSummary = types.StringNull()
        }
    } else if val, ok := item["runnerVerificationSummary"].(string); ok {
        data.RunnerVerificationSummary = types.StringValue(val)
    } else {
        data.RunnerVerificationSummary = types.StringNull()
    }
    if obj, ok := item["headRefName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.HeadRefName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.HeadRefName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.HeadRefName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.HeadRefName = types.StringValue(string(jsonBytes))
        } else {
            data.HeadRefName = types.StringNull()
        }
    } else if val, ok := item["headRefName"].(string); ok {
        data.HeadRefName = types.StringValue(val)
    } else {
        data.HeadRefName = types.StringNull()
    }
    if obj, ok := item["baseRefName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.BaseRefName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.BaseRefName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.BaseRefName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.BaseRefName = types.StringValue(string(jsonBytes))
        } else {
            data.BaseRefName = types.StringNull()
        }
    } else if val, ok := item["baseRefName"].(string); ok {
        data.BaseRefName = types.StringValue(val)
    } else {
        data.BaseRefName = types.StringNull()
    }
    if obj, ok := item["repoOrganizationName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RepoOrganizationName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RepoOrganizationName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RepoOrganizationName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RepoOrganizationName = types.StringValue(string(jsonBytes))
        } else {
            data.RepoOrganizationName = types.StringNull()
        }
    } else if val, ok := item["repoOrganizationName"].(string); ok {
        data.RepoOrganizationName = types.StringValue(val)
    } else {
        data.RepoOrganizationName = types.StringNull()
    }
    if obj, ok := item["repoName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RepoName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RepoName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RepoName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RepoName = types.StringValue(string(jsonBytes))
        } else {
            data.RepoName = types.StringNull()
        }
    } else if val, ok := item["repoName"].(string); ok {
        data.RepoName = types.StringValue(val)
    } else {
        data.RepoName = types.StringNull()
    }
    if obj, ok := item["createdByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if val, ok := item["createdByUserId"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
