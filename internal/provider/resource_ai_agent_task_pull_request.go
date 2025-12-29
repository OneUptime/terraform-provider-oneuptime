package provider

import (
    "context"
    "fmt"
    "github.com/hashicorp/terraform-plugin-framework/path"
    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
    "math/big"
    "net/http"
    "encoding/json"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &AiAgentTaskPullRequestResource{}
var _ resource.ResourceWithImportState = &AiAgentTaskPullRequestResource{}

func NewAiAgentTaskPullRequestResource() resource.Resource {
    return &AiAgentTaskPullRequestResource{}
}

// AiAgentTaskPullRequestResource defines the resource implementation.
type AiAgentTaskPullRequestResource struct {
    client *Client
}

// AiAgentTaskPullRequestResourceModel describes the resource data model.
type AiAgentTaskPullRequestResourceModel struct {
    Id types.String `tfsdk:"id"`
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
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (r *AiAgentTaskPullRequestResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_ai_agent_task_pull_request"
}

func (r *AiAgentTaskPullRequestResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "ai_agent_task_pull_request resource",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Unique identifier for the resource",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Optional: true,
                Computed: true,
            },
            "ai_agent_task_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Required: true,
            },
            "ai_agent_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Required: true,
            },
            "code_repository_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Optional: true,
                Computed: true,
            },
            "title": schema.StringAttribute{
                MarkdownDescription: "Title of the pull request.. Permissions - Create: [Project Owner, Project Admin, Project Member, Create AI Agent Task], Read: [Project Owner, Project Admin, Project Member, Read AI Agent Task], Update: [Project Owner, Project Admin, Project Member, Edit AI Agent Task]",
                Required: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Description/body of the pull request.. Permissions - Create: [Project Owner, Project Admin, Project Member, Create AI Agent Task], Read: [Project Owner, Project Admin, Project Member, Read AI Agent Task], Update: [Project Owner, Project Admin, Project Member, Edit AI Agent Task]",
                Optional: true,
                Computed: true,
            },
            "pull_request_url": schema.StringAttribute{
                MarkdownDescription: "URL to the pull request on the hosting platform.. Permissions - Create: [Project Owner, Project Admin, Project Member, Create AI Agent Task], Read: [Project Owner, Project Admin, Project Member, Read AI Agent Task], Update: [Project Owner, Project Admin, Project Member, Edit AI Agent Task]",
                Optional: true,
                Computed: true,
            },
            "pull_request_id": schema.NumberAttribute{
                MarkdownDescription: "The unique ID of the pull request from the hosting platform.. Permissions - Create: [Project Owner, Project Admin, Project Member, Create AI Agent Task], Read: [Project Owner, Project Admin, Project Member, Read AI Agent Task], Update: [No access - you don't have permission for this operation]",
                Optional: true,
                Computed: true,
            },
            "pull_request_number": schema.NumberAttribute{
                MarkdownDescription: "The pull request number (e.g., #123).. Permissions - Create: [Project Owner, Project Admin, Project Member, Create AI Agent Task], Read: [Project Owner, Project Admin, Project Member, Read AI Agent Task], Update: [No access - you don't have permission for this operation]",
                Optional: true,
                Computed: true,
            },
            "pull_request_state": schema.StringAttribute{
                MarkdownDescription: "Current state of the pull request (open, closed, merged).. Permissions - Create: [Project Owner, Project Admin, Project Member, Create AI Agent Task], Read: [Project Owner, Project Admin, Project Member, Read AI Agent Task], Update: [Project Owner, Project Admin, Project Member, Edit AI Agent Task]",
                Required: true,
            },
            "head_ref_name": schema.StringAttribute{
                MarkdownDescription: "The branch name of the pull request (source branch).. Permissions - Create: [Project Owner, Project Admin, Project Member, Create AI Agent Task], Read: [Project Owner, Project Admin, Project Member, Read AI Agent Task], Update: [No access - you don't have permission for this operation]",
                Optional: true,
                Computed: true,
            },
            "base_ref_name": schema.StringAttribute{
                MarkdownDescription: "The target branch for the pull request.. Permissions - Create: [Project Owner, Project Admin, Project Member, Create AI Agent Task], Read: [Project Owner, Project Admin, Project Member, Read AI Agent Task], Update: [No access - you don't have permission for this operation]",
                Optional: true,
                Computed: true,
            },
            "repo_organization_name": schema.StringAttribute{
                MarkdownDescription: "Organization or username that owns the repository.. Permissions - Create: [Project Owner, Project Admin, Project Member, Create AI Agent Task], Read: [Project Owner, Project Admin, Project Member, Read AI Agent Task], Update: [No access - you don't have permission for this operation]",
                Optional: true,
                Computed: true,
            },
            "repo_name": schema.StringAttribute{
                MarkdownDescription: "Name of the repository.. Permissions - Create: [Project Owner, Project Admin, Project Member, Create AI Agent Task], Read: [Project Owner, Project Admin, Project Member, Read AI Agent Task], Update: [No access - you don't have permission for this operation]",
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
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
        },
    }
}

func (r *AiAgentTaskPullRequestResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
    // Prevent panic if the provider has not been configured.
    if req.ProviderData == nil {
        return
    }

    client, ok := req.ProviderData.(*Client)

    if !ok {
        resp.Diagnostics.AddError(
            "Unexpected Resource Configure Type",
            fmt.Sprintf("Expected *Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
        )

        return
    }

    r.client = client
}


func (r *AiAgentTaskPullRequestResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var data AiAgentTaskPullRequestResourceModel

    // Read Terraform plan data into the model
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }



    // Create API request body
    aiAgentTaskPullRequestRequest := map[string]interface{}{
        "data": map[string]interface{}{
        "projectId": data.ProjectId.ValueString(),
        "aiAgentTaskId": data.AiAgentTaskId.ValueString(),
        "aiAgentId": data.AiAgentId.ValueString(),
        "codeRepositoryId": data.CodeRepositoryId.ValueString(),
        "title": data.Title.ValueString(),
        "description": data.Description.ValueString(),
        "pullRequestUrl": data.PullRequestUrl.ValueString(),
        "pullRequestId": data.PullRequestId.ValueBigFloat(),
        "pullRequestNumber": data.PullRequestNumber.ValueBigFloat(),
        "pullRequestState": data.PullRequestState.ValueString(),
        "headRefName": data.HeadRefName.ValueString(),
        "baseRefName": data.BaseRefName.ValueString(),
        "repoOrganizationName": data.RepoOrganizationName.ValueString(),
        "repoName": data.RepoName.ValueString(),
        },
    }

    // Make API call
    httpResp, err := r.client.Post("/ai-agent-task-pull-request", aiAgentTaskPullRequestRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create ai_agent_task_pull_request, got error: %s", err))
        return
    }

    var aiAgentTaskPullRequestResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &aiAgentTaskPullRequestResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse ai_agent_task_pull_request response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := aiAgentTaskPullRequestResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = aiAgentTaskPullRequestResponse
    }

    if obj, ok := dataMap["id"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Id = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.Id = types.StringValue(val)
        } else {
            data.Id = types.StringNull()
        }
    } else if val, ok := dataMap["id"].(string); ok && val != "" {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    if obj, ok := dataMap["projectId"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok {
            data.ProjectId = types.StringValue(val)
        } else {
            data.ProjectId = types.StringNull()
        }
    } else if val, ok := dataMap["projectId"].(string); ok {
        data.ProjectId = types.StringValue(val)
    } else {
        data.ProjectId = types.StringNull()
    }
    if obj, ok := dataMap["aiAgentTaskId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AiAgentTaskId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.AiAgentTaskId = types.StringValue(val)
        } else {
            data.AiAgentTaskId = types.StringNull()
        }
    } else if val, ok := dataMap["aiAgentTaskId"].(string); ok && val != "" {
        data.AiAgentTaskId = types.StringValue(val)
    } else {
        data.AiAgentTaskId = types.StringNull()
    }
    if obj, ok := dataMap["aiAgentId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AiAgentId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.AiAgentId = types.StringValue(val)
        } else {
            data.AiAgentId = types.StringNull()
        }
    } else if val, ok := dataMap["aiAgentId"].(string); ok && val != "" {
        data.AiAgentId = types.StringValue(val)
    } else {
        data.AiAgentId = types.StringNull()
    }
    if obj, ok := dataMap["codeRepositoryId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CodeRepositoryId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.CodeRepositoryId = types.StringValue(val)
        } else {
            data.CodeRepositoryId = types.StringNull()
        }
    } else if val, ok := dataMap["codeRepositoryId"].(string); ok && val != "" {
        data.CodeRepositoryId = types.StringValue(val)
    } else {
        data.CodeRepositoryId = types.StringNull()
    }
    if obj, ok := dataMap["title"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Title = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.Title = types.StringValue(val)
        } else {
            data.Title = types.StringNull()
        }
    } else if val, ok := dataMap["title"].(string); ok && val != "" {
        data.Title = types.StringValue(val)
    } else {
        data.Title = types.StringNull()
    }
    if obj, ok := dataMap["description"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.Description = types.StringValue(val)
        } else {
            data.Description = types.StringNull()
        }
    } else if val, ok := dataMap["description"].(string); ok && val != "" {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if obj, ok := dataMap["pullRequestUrl"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PullRequestUrl = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.PullRequestUrl = types.StringValue(val)
        } else {
            data.PullRequestUrl = types.StringNull()
        }
    } else if val, ok := dataMap["pullRequestUrl"].(string); ok && val != "" {
        data.PullRequestUrl = types.StringValue(val)
    } else {
        data.PullRequestUrl = types.StringNull()
    }
    if val, ok := dataMap["pullRequestId"].(float64); ok {
        data.PullRequestId = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["pullRequestId"].(int); ok {
        data.PullRequestId = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["pullRequestId"].(int64); ok {
        data.PullRequestId = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["pullRequestId"] == nil {
        data.PullRequestId = types.NumberNull()
    }
    if val, ok := dataMap["pullRequestNumber"].(float64); ok {
        data.PullRequestNumber = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["pullRequestNumber"].(int); ok {
        data.PullRequestNumber = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["pullRequestNumber"].(int64); ok {
        data.PullRequestNumber = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["pullRequestNumber"] == nil {
        data.PullRequestNumber = types.NumberNull()
    }
    if obj, ok := dataMap["pullRequestState"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PullRequestState = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.PullRequestState = types.StringValue(val)
        } else {
            data.PullRequestState = types.StringNull()
        }
    } else if val, ok := dataMap["pullRequestState"].(string); ok && val != "" {
        data.PullRequestState = types.StringValue(val)
    } else {
        data.PullRequestState = types.StringNull()
    }
    if obj, ok := dataMap["headRefName"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.HeadRefName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.HeadRefName = types.StringValue(val)
        } else {
            data.HeadRefName = types.StringNull()
        }
    } else if val, ok := dataMap["headRefName"].(string); ok && val != "" {
        data.HeadRefName = types.StringValue(val)
    } else {
        data.HeadRefName = types.StringNull()
    }
    if obj, ok := dataMap["baseRefName"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.BaseRefName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.BaseRefName = types.StringValue(val)
        } else {
            data.BaseRefName = types.StringNull()
        }
    } else if val, ok := dataMap["baseRefName"].(string); ok && val != "" {
        data.BaseRefName = types.StringValue(val)
    } else {
        data.BaseRefName = types.StringNull()
    }
    if obj, ok := dataMap["repoOrganizationName"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RepoOrganizationName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.RepoOrganizationName = types.StringValue(val)
        } else {
            data.RepoOrganizationName = types.StringNull()
        }
    } else if val, ok := dataMap["repoOrganizationName"].(string); ok && val != "" {
        data.RepoOrganizationName = types.StringValue(val)
    } else {
        data.RepoOrganizationName = types.StringNull()
    }
    if obj, ok := dataMap["repoName"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RepoName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.RepoName = types.StringValue(val)
        } else {
            data.RepoName = types.StringNull()
        }
    } else if val, ok := dataMap["repoName"].(string); ok && val != "" {
        data.RepoName = types.StringValue(val)
    } else {
        data.RepoName = types.StringNull()
    }
    if val, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        if jsonBytes, err := json.Marshal(val); err == nil {
            data.CreatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedAt = types.StringNull()
        }
    } else if val, ok := dataMap["createdAt"].(string); ok && val != "" {
        data.CreatedAt = types.StringValue(val)
    } else {
        data.CreatedAt = types.StringNull()
    }
    if val, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        if jsonBytes, err := json.Marshal(val); err == nil {
            data.UpdatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.UpdatedAt = types.StringNull()
        }
    } else if val, ok := dataMap["updatedAt"].(string); ok && val != "" {
        data.UpdatedAt = types.StringValue(val)
    } else {
        data.UpdatedAt = types.StringNull()
    }
    if val, ok := dataMap["deletedAt"].(map[string]interface{}); ok {
        if jsonBytes, err := json.Marshal(val); err == nil {
            data.DeletedAt = types.StringValue(string(jsonBytes))
        } else {
            data.DeletedAt = types.StringNull()
        }
    } else if val, ok := dataMap["deletedAt"].(string); ok && val != "" {
        data.DeletedAt = types.StringValue(val)
    } else {
        data.DeletedAt = types.StringNull()
    }
    if val, ok := dataMap["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["version"].(int); ok {
        data.Version = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["version"].(int64); ok {
        data.Version = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["version"] == nil {
        data.Version = types.NumberNull()
    }
    if obj, ok := dataMap["createdByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.CreatedByUserId = types.StringValue(val)
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["createdByUserId"].(string); ok && val != "" {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "created a resource")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AiAgentTaskPullRequestResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var data AiAgentTaskPullRequestResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create select parameter to get full object
    selectParam := map[string]interface{}{
        "projectId": true,
        "aiAgentTaskId": true,
        "aiAgentId": true,
        "codeRepositoryId": true,
        "title": true,
        "description": true,
        "pullRequestUrl": true,
        "pullRequestId": true,
        "pullRequestNumber": true,
        "pullRequestState": true,
        "headRefName": true,
        "baseRefName": true,
        "repoOrganizationName": true,
        "repoName": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "createdByUserId": true,
        "_id": true,
    }

    // Make API call with select parameter
    httpResp, err := r.client.PostWithSelect("/ai-agent-task-pull-request/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read ai_agent_task_pull_request, got error: %s", err))
        return
    }

    if httpResp.StatusCode == http.StatusNotFound {
        resp.State.RemoveResource(ctx)
        return
    }

    var aiAgentTaskPullRequestResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &aiAgentTaskPullRequestResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse ai_agent_task_pull_request response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := aiAgentTaskPullRequestResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = aiAgentTaskPullRequestResponse
    }

    if obj, ok := dataMap["id"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Id = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.Id = types.StringValue(val)
        } else {
            data.Id = types.StringNull()
        }
    } else if val, ok := dataMap["id"].(string); ok && val != "" {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    if obj, ok := dataMap["projectId"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok {
            data.ProjectId = types.StringValue(val)
        } else {
            data.ProjectId = types.StringNull()
        }
    } else if val, ok := dataMap["projectId"].(string); ok {
        data.ProjectId = types.StringValue(val)
    } else {
        data.ProjectId = types.StringNull()
    }
    if obj, ok := dataMap["aiAgentTaskId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AiAgentTaskId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.AiAgentTaskId = types.StringValue(val)
        } else {
            data.AiAgentTaskId = types.StringNull()
        }
    } else if val, ok := dataMap["aiAgentTaskId"].(string); ok && val != "" {
        data.AiAgentTaskId = types.StringValue(val)
    } else {
        data.AiAgentTaskId = types.StringNull()
    }
    if obj, ok := dataMap["aiAgentId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AiAgentId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.AiAgentId = types.StringValue(val)
        } else {
            data.AiAgentId = types.StringNull()
        }
    } else if val, ok := dataMap["aiAgentId"].(string); ok && val != "" {
        data.AiAgentId = types.StringValue(val)
    } else {
        data.AiAgentId = types.StringNull()
    }
    if obj, ok := dataMap["codeRepositoryId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CodeRepositoryId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.CodeRepositoryId = types.StringValue(val)
        } else {
            data.CodeRepositoryId = types.StringNull()
        }
    } else if val, ok := dataMap["codeRepositoryId"].(string); ok && val != "" {
        data.CodeRepositoryId = types.StringValue(val)
    } else {
        data.CodeRepositoryId = types.StringNull()
    }
    if obj, ok := dataMap["title"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Title = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.Title = types.StringValue(val)
        } else {
            data.Title = types.StringNull()
        }
    } else if val, ok := dataMap["title"].(string); ok && val != "" {
        data.Title = types.StringValue(val)
    } else {
        data.Title = types.StringNull()
    }
    if obj, ok := dataMap["description"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.Description = types.StringValue(val)
        } else {
            data.Description = types.StringNull()
        }
    } else if val, ok := dataMap["description"].(string); ok && val != "" {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if obj, ok := dataMap["pullRequestUrl"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PullRequestUrl = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.PullRequestUrl = types.StringValue(val)
        } else {
            data.PullRequestUrl = types.StringNull()
        }
    } else if val, ok := dataMap["pullRequestUrl"].(string); ok && val != "" {
        data.PullRequestUrl = types.StringValue(val)
    } else {
        data.PullRequestUrl = types.StringNull()
    }
    if val, ok := dataMap["pullRequestId"].(float64); ok {
        data.PullRequestId = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["pullRequestId"].(int); ok {
        data.PullRequestId = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["pullRequestId"].(int64); ok {
        data.PullRequestId = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["pullRequestId"] == nil {
        data.PullRequestId = types.NumberNull()
    }
    if val, ok := dataMap["pullRequestNumber"].(float64); ok {
        data.PullRequestNumber = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["pullRequestNumber"].(int); ok {
        data.PullRequestNumber = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["pullRequestNumber"].(int64); ok {
        data.PullRequestNumber = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["pullRequestNumber"] == nil {
        data.PullRequestNumber = types.NumberNull()
    }
    if obj, ok := dataMap["pullRequestState"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PullRequestState = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.PullRequestState = types.StringValue(val)
        } else {
            data.PullRequestState = types.StringNull()
        }
    } else if val, ok := dataMap["pullRequestState"].(string); ok && val != "" {
        data.PullRequestState = types.StringValue(val)
    } else {
        data.PullRequestState = types.StringNull()
    }
    if obj, ok := dataMap["headRefName"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.HeadRefName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.HeadRefName = types.StringValue(val)
        } else {
            data.HeadRefName = types.StringNull()
        }
    } else if val, ok := dataMap["headRefName"].(string); ok && val != "" {
        data.HeadRefName = types.StringValue(val)
    } else {
        data.HeadRefName = types.StringNull()
    }
    if obj, ok := dataMap["baseRefName"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.BaseRefName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.BaseRefName = types.StringValue(val)
        } else {
            data.BaseRefName = types.StringNull()
        }
    } else if val, ok := dataMap["baseRefName"].(string); ok && val != "" {
        data.BaseRefName = types.StringValue(val)
    } else {
        data.BaseRefName = types.StringNull()
    }
    if obj, ok := dataMap["repoOrganizationName"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RepoOrganizationName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.RepoOrganizationName = types.StringValue(val)
        } else {
            data.RepoOrganizationName = types.StringNull()
        }
    } else if val, ok := dataMap["repoOrganizationName"].(string); ok && val != "" {
        data.RepoOrganizationName = types.StringValue(val)
    } else {
        data.RepoOrganizationName = types.StringNull()
    }
    if obj, ok := dataMap["repoName"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RepoName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.RepoName = types.StringValue(val)
        } else {
            data.RepoName = types.StringNull()
        }
    } else if val, ok := dataMap["repoName"].(string); ok && val != "" {
        data.RepoName = types.StringValue(val)
    } else {
        data.RepoName = types.StringNull()
    }
    if val, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        if jsonBytes, err := json.Marshal(val); err == nil {
            data.CreatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedAt = types.StringNull()
        }
    } else if val, ok := dataMap["createdAt"].(string); ok && val != "" {
        data.CreatedAt = types.StringValue(val)
    } else {
        data.CreatedAt = types.StringNull()
    }
    if val, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        if jsonBytes, err := json.Marshal(val); err == nil {
            data.UpdatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.UpdatedAt = types.StringNull()
        }
    } else if val, ok := dataMap["updatedAt"].(string); ok && val != "" {
        data.UpdatedAt = types.StringValue(val)
    } else {
        data.UpdatedAt = types.StringNull()
    }
    if val, ok := dataMap["deletedAt"].(map[string]interface{}); ok {
        if jsonBytes, err := json.Marshal(val); err == nil {
            data.DeletedAt = types.StringValue(string(jsonBytes))
        } else {
            data.DeletedAt = types.StringNull()
        }
    } else if val, ok := dataMap["deletedAt"].(string); ok && val != "" {
        data.DeletedAt = types.StringValue(val)
    } else {
        data.DeletedAt = types.StringNull()
    }
    if val, ok := dataMap["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["version"].(int); ok {
        data.Version = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["version"].(int64); ok {
        data.Version = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["version"] == nil {
        data.Version = types.NumberNull()
    }
    if obj, ok := dataMap["createdByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.CreatedByUserId = types.StringValue(val)
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["createdByUserId"].(string); ok && val != "" {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AiAgentTaskPullRequestResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var data AiAgentTaskPullRequestResourceModel
    var state AiAgentTaskPullRequestResourceModel

    // Read Terraform current state data to get the ID
    resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Read Terraform plan data to get the new values
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Use the ID from the current state
    data.Id = state.Id

    // Create API request body
    aiAgentTaskPullRequestRequest := map[string]interface{}{
        "data": map[string]interface{}{},
    }
    requestDataMap := aiAgentTaskPullRequestRequest["data"].(map[string]interface{})

    if !data.Title.IsUnknown() && !state.Title.IsUnknown() && !data.Title.Equal(state.Title) {
        requestDataMap["title"] = data.Title.ValueString()
    }
    if !data.Description.IsUnknown() && !state.Description.IsUnknown() && !data.Description.Equal(state.Description) {
        requestDataMap["description"] = data.Description.ValueString()
    }
    if !data.PullRequestUrl.IsUnknown() && !state.PullRequestUrl.IsUnknown() && !data.PullRequestUrl.Equal(state.PullRequestUrl) {
        requestDataMap["pullRequestUrl"] = data.PullRequestUrl.ValueString()
    }
    if !data.PullRequestState.IsUnknown() && !state.PullRequestState.IsUnknown() && !data.PullRequestState.Equal(state.PullRequestState) {
        requestDataMap["pullRequestState"] = data.PullRequestState.ValueString()
    }

    // Make API call
    httpResp, err := r.client.Put("/ai-agent-task-pull-request/" + data.Id.ValueString() + "", aiAgentTaskPullRequestRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update ai_agent_task_pull_request, got error: %s", err))
        return
    }

    // Parse the update response
    var aiAgentTaskPullRequestResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &aiAgentTaskPullRequestResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse ai_agent_task_pull_request response, got error: %s", err))
        return
    }

    // After successful update, fetch the current state by calling Read with select parameter
    selectParam := map[string]interface{}{
        "projectId": true,
        "aiAgentTaskId": true,
        "aiAgentId": true,
        "codeRepositoryId": true,
        "title": true,
        "description": true,
        "pullRequestUrl": true,
        "pullRequestId": true,
        "pullRequestNumber": true,
        "pullRequestState": true,
        "headRefName": true,
        "baseRefName": true,
        "repoOrganizationName": true,
        "repoName": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "createdByUserId": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect("/ai-agent-task-pull-request/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read ai_agent_task_pull_request after update, got error: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse ai_agent_task_pull_request read response, got error: %s", err))
        return
    }

    // Update the model with response data from the Read operation
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := readResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = readResponse
    }

    if obj, ok := dataMap["id"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Id = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.Id = types.StringValue(val)
        } else {
            data.Id = types.StringNull()
        }
    } else if val, ok := dataMap["id"].(string); ok && val != "" {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    if obj, ok := dataMap["projectId"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok {
            data.ProjectId = types.StringValue(val)
        } else {
            data.ProjectId = types.StringNull()
        }
    } else if val, ok := dataMap["projectId"].(string); ok {
        data.ProjectId = types.StringValue(val)
    } else {
        data.ProjectId = types.StringNull()
    }
    if obj, ok := dataMap["aiAgentTaskId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AiAgentTaskId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.AiAgentTaskId = types.StringValue(val)
        } else {
            data.AiAgentTaskId = types.StringNull()
        }
    } else if val, ok := dataMap["aiAgentTaskId"].(string); ok && val != "" {
        data.AiAgentTaskId = types.StringValue(val)
    } else {
        data.AiAgentTaskId = types.StringNull()
    }
    if obj, ok := dataMap["aiAgentId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AiAgentId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.AiAgentId = types.StringValue(val)
        } else {
            data.AiAgentId = types.StringNull()
        }
    } else if val, ok := dataMap["aiAgentId"].(string); ok && val != "" {
        data.AiAgentId = types.StringValue(val)
    } else {
        data.AiAgentId = types.StringNull()
    }
    if obj, ok := dataMap["codeRepositoryId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CodeRepositoryId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.CodeRepositoryId = types.StringValue(val)
        } else {
            data.CodeRepositoryId = types.StringNull()
        }
    } else if val, ok := dataMap["codeRepositoryId"].(string); ok && val != "" {
        data.CodeRepositoryId = types.StringValue(val)
    } else {
        data.CodeRepositoryId = types.StringNull()
    }
    if obj, ok := dataMap["title"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Title = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.Title = types.StringValue(val)
        } else {
            data.Title = types.StringNull()
        }
    } else if val, ok := dataMap["title"].(string); ok && val != "" {
        data.Title = types.StringValue(val)
    } else {
        data.Title = types.StringNull()
    }
    if obj, ok := dataMap["description"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.Description = types.StringValue(val)
        } else {
            data.Description = types.StringNull()
        }
    } else if val, ok := dataMap["description"].(string); ok && val != "" {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if obj, ok := dataMap["pullRequestUrl"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PullRequestUrl = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.PullRequestUrl = types.StringValue(val)
        } else {
            data.PullRequestUrl = types.StringNull()
        }
    } else if val, ok := dataMap["pullRequestUrl"].(string); ok && val != "" {
        data.PullRequestUrl = types.StringValue(val)
    } else {
        data.PullRequestUrl = types.StringNull()
    }
    if val, ok := dataMap["pullRequestId"].(float64); ok {
        data.PullRequestId = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["pullRequestId"].(int); ok {
        data.PullRequestId = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["pullRequestId"].(int64); ok {
        data.PullRequestId = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["pullRequestId"] == nil {
        data.PullRequestId = types.NumberNull()
    }
    if val, ok := dataMap["pullRequestNumber"].(float64); ok {
        data.PullRequestNumber = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["pullRequestNumber"].(int); ok {
        data.PullRequestNumber = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["pullRequestNumber"].(int64); ok {
        data.PullRequestNumber = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["pullRequestNumber"] == nil {
        data.PullRequestNumber = types.NumberNull()
    }
    if obj, ok := dataMap["pullRequestState"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PullRequestState = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.PullRequestState = types.StringValue(val)
        } else {
            data.PullRequestState = types.StringNull()
        }
    } else if val, ok := dataMap["pullRequestState"].(string); ok && val != "" {
        data.PullRequestState = types.StringValue(val)
    } else {
        data.PullRequestState = types.StringNull()
    }
    if obj, ok := dataMap["headRefName"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.HeadRefName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.HeadRefName = types.StringValue(val)
        } else {
            data.HeadRefName = types.StringNull()
        }
    } else if val, ok := dataMap["headRefName"].(string); ok && val != "" {
        data.HeadRefName = types.StringValue(val)
    } else {
        data.HeadRefName = types.StringNull()
    }
    if obj, ok := dataMap["baseRefName"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.BaseRefName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.BaseRefName = types.StringValue(val)
        } else {
            data.BaseRefName = types.StringNull()
        }
    } else if val, ok := dataMap["baseRefName"].(string); ok && val != "" {
        data.BaseRefName = types.StringValue(val)
    } else {
        data.BaseRefName = types.StringNull()
    }
    if obj, ok := dataMap["repoOrganizationName"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RepoOrganizationName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.RepoOrganizationName = types.StringValue(val)
        } else {
            data.RepoOrganizationName = types.StringNull()
        }
    } else if val, ok := dataMap["repoOrganizationName"].(string); ok && val != "" {
        data.RepoOrganizationName = types.StringValue(val)
    } else {
        data.RepoOrganizationName = types.StringNull()
    }
    if obj, ok := dataMap["repoName"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RepoName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.RepoName = types.StringValue(val)
        } else {
            data.RepoName = types.StringNull()
        }
    } else if val, ok := dataMap["repoName"].(string); ok && val != "" {
        data.RepoName = types.StringValue(val)
    } else {
        data.RepoName = types.StringNull()
    }
    if val, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        if jsonBytes, err := json.Marshal(val); err == nil {
            data.CreatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedAt = types.StringNull()
        }
    } else if val, ok := dataMap["createdAt"].(string); ok && val != "" {
        data.CreatedAt = types.StringValue(val)
    } else {
        data.CreatedAt = types.StringNull()
    }
    if val, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        if jsonBytes, err := json.Marshal(val); err == nil {
            data.UpdatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.UpdatedAt = types.StringNull()
        }
    } else if val, ok := dataMap["updatedAt"].(string); ok && val != "" {
        data.UpdatedAt = types.StringValue(val)
    } else {
        data.UpdatedAt = types.StringNull()
    }
    if val, ok := dataMap["deletedAt"].(map[string]interface{}); ok {
        if jsonBytes, err := json.Marshal(val); err == nil {
            data.DeletedAt = types.StringValue(string(jsonBytes))
        } else {
            data.DeletedAt = types.StringNull()
        }
    } else if val, ok := dataMap["deletedAt"].(string); ok && val != "" {
        data.DeletedAt = types.StringValue(val)
    } else {
        data.DeletedAt = types.StringNull()
    }
    if val, ok := dataMap["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["version"].(int); ok {
        data.Version = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["version"].(int64); ok {
        data.Version = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["version"] == nil {
        data.Version = types.NumberNull()
    }
    if obj, ok := dataMap["createdByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.CreatedByUserId = types.StringValue(val)
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["createdByUserId"].(string); ok && val != "" {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AiAgentTaskPullRequestResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var data AiAgentTaskPullRequestResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Make API call
    _, err := r.client.Delete("/ai-agent-task-pull-request/" + data.Id.ValueString() + "")
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete ai_agent_task_pull_request, got error: %s", err))
        return
    }
}


func (r *AiAgentTaskPullRequestResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Helper method to convert Terraform map to Go interface{}
func (r *AiAgentTaskPullRequestResource) convertTerraformMapToInterface(terraformMap types.Map) interface{} {
    if terraformMap.IsNull() || terraformMap.IsUnknown() {
        return nil
    }
    
    result := make(map[string]string)
    terraformMap.ElementsAs(context.Background(), &result, false)
    
    // Convert map[string]string to map[string]interface{}
    interfaceResult := make(map[string]interface{})
    for key, value := range result {
        interfaceResult[key] = value
    }
    
    return interfaceResult
}

// Helper method to convert Terraform list to Go interface{}
func (r *AiAgentTaskPullRequestResource) convertTerraformListToInterface(terraformList types.List) interface{} {
    if terraformList.IsNull() || terraformList.IsUnknown() {
        return nil
    }
    
    var stringList []string
    terraformList.ElementsAs(context.Background(), &stringList, false)
    
    // Convert string array to OneUptime format with _id fields
    var result []interface{}
    for _, str := range stringList {
        if str != "" {
            result = append(result, map[string]interface{}{
                "_id": str,
            })
        }
    }
    return result
}

// Helper method to parse JSON field for complex objects
func (r *AiAgentTaskPullRequestResource) parseJSONField(terraformString types.String) interface{} {
    if terraformString.IsNull() || terraformString.IsUnknown() || terraformString.ValueString() == "" {
        return nil
    }
    
    var result interface{}
    if err := json.Unmarshal([]byte(terraformString.ValueString()), &result); err != nil {
        // If JSON parsing fails, return the raw string
        return terraformString.ValueString()
    }
    
    return result
}
