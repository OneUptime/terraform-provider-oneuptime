package provider

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "math/big"
    "github.com/hashicorp/terraform-plugin-framework/attr"
    "sort"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &CodeRepositoryDataSource{}

func NewCodeRepositoryDataSource() datasource.DataSource {
    return &CodeRepositoryDataSource{}
}

// CodeRepositoryDataSource defines the data source implementation.
type CodeRepositoryDataSource struct {
    client *Client
}

// CodeRepositoryDataSourceModel describes the data source data model.
type CodeRepositoryDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    Slug types.String `tfsdk:"slug"`
    Description types.String `tfsdk:"description"`
    RepositoryHostedAt types.String `tfsdk:"repository_hosted_at"`
    OrganizationName types.String `tfsdk:"organization_name"`
    RepositoryName types.String `tfsdk:"repository_name"`
    MainBranchName types.String `tfsdk:"main_branch_name"`
    SetupCommand types.String `tfsdk:"setup_command"`
    BuildCommand types.String `tfsdk:"build_command"`
    TestCommand types.String `tfsdk:"test_command"`
    MaxOpenFixPullRequests types.Number `tfsdk:"max_open_fix_pull_requests"`
    RepositoryUrl types.String `tfsdk:"repository_url"`
    GitHubAppInstallationId types.String `tfsdk:"git_hub_app_installation_id"`
    GitLabProjectId types.String `tfsdk:"git_lab_project_id"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
    Labels types.Set `tfsdk:"labels"`
}

func (d *CodeRepositoryDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_code_repository"
}

func (d *CodeRepositoryDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Connect and manage code repositories from GitHub, GitLab, and other providers Look up an existing code_repository by `id` or by `name`.",

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
            "slug": schema.StringAttribute{
                MarkdownDescription: "Friendly globally unique name for your object.",
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "A description of this code repository.",
                Computed: true,
            },
            "repository_hosted_at": schema.StringAttribute{
                MarkdownDescription: "Where is this repository hosted (GitHub, GitLab, etc.).",
                Computed: true,
            },
            "organization_name": schema.StringAttribute{
                MarkdownDescription: "GitHub organization or username that owns this repository.",
                Computed: true,
            },
            "repository_name": schema.StringAttribute{
                MarkdownDescription: "The name of the repository.",
                Computed: true,
            },
            "main_branch_name": schema.StringAttribute{
                MarkdownDescription: "The name of the main/default branch.",
                Computed: true,
            },
            "setup_command": schema.StringAttribute{
                MarkdownDescription: "Command the AI fix Runner executes at the repository root to install dependencies before verifying an AI-authored fix (e.g. 'npm ci'). Runs on your Runner, in the cloned workspace, before the build and test commands. Leave empty to skip..",
                Computed: true,
            },
            "build_command": schema.StringAttribute{
                MarkdownDescription: "Command the AI fix Runner executes at the repository root to verify an AI-authored fix compiles/builds (e.g. 'npm run build'). A failure is fed back to the code agent for bounded repair attempts before the pull request opens. Leave empty to skip the build check..",
                Computed: true,
            },
            "test_command": schema.StringAttribute{
                MarkdownDescription: "Command the AI fix Runner executes at the repository root to run the test suite against an AI-authored fix (e.g. 'npm test'). A failure is fed back to the code agent for bounded repair attempts before the pull request opens. Leave empty to skip the test check..",
                Computed: true,
            },
            "max_open_fix_pull_requests": schema.NumberAttribute{
                MarkdownDescription: "Maximum AI-authored fix pull requests that may be open on this repository at the same time. At the cap, new AI fix runs are refused a repository token, so they cannot push branches or open pull requests. Unset means the default of 5; 0 blocks AI fix pull requests for this repository entirely..",
                Computed: true,
            },
            "repository_url": schema.StringAttribute{
                MarkdownDescription: "The HTTPS URL to the repository.",
                Computed: true,
            },
            "git_hub_app_installation_id": schema.StringAttribute{
                MarkdownDescription: "The GitHub App installation ID used to authenticate with this repository.",
                Computed: true,
            },
            "git_lab_project_id": schema.StringAttribute{
                MarkdownDescription: "The GitLab project ID for this repository.",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "deleted_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "labels": schema.SetAttribute{
                MarkdownDescription: "Relation to Labels Array where this object is categorized in..",
                Computed: true,
                ElementType: types.StringType,
            },
        },
    }
}

func (d *CodeRepositoryDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *CodeRepositoryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data CodeRepositoryDataSourceModel

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
            "Exactly one of `id` or `name` must be set to look up a code_repository.",
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
        "slug": true,
        "description": true,
        "repositoryHostedAt": true,
        "organizationName": true,
        "repositoryName": true,
        "mainBranchName": true,
        "setupCommand": true,
        "buildCommand": true,
        "testCommand": true,
        "maxOpenFixPullRequests": true,
        "repositoryUrl": true,
        "gitHubAppInstallationId": true,
        "gitLabProjectId": true,
        "createdByUserId": true,
        "deletedByUserId": true,
        "labels": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/code-repository/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read code_repository, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No code_repository found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read code_repository: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/code-repository/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list code_repository, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list code_repository: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No code_repository found with name %q.", data.Name.ValueString()))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one code_repository matches name %q. Use the id attribute to disambiguate.", data.Name.ValueString()))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for code_repository.")
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
    if obj, ok := item["slug"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Slug = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Slug = types.StringValue(string(jsonBytes))
        } else {
            data.Slug = types.StringNull()
        }
    } else if val, ok := item["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
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
    if obj, ok := item["repositoryHostedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RepositoryHostedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RepositoryHostedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RepositoryHostedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RepositoryHostedAt = types.StringValue(string(jsonBytes))
        } else {
            data.RepositoryHostedAt = types.StringNull()
        }
    } else if val, ok := item["repositoryHostedAt"].(string); ok {
        data.RepositoryHostedAt = types.StringValue(val)
    } else {
        data.RepositoryHostedAt = types.StringNull()
    }
    if obj, ok := item["organizationName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OrganizationName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.OrganizationName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.OrganizationName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.OrganizationName = types.StringValue(string(jsonBytes))
        } else {
            data.OrganizationName = types.StringNull()
        }
    } else if val, ok := item["organizationName"].(string); ok {
        data.OrganizationName = types.StringValue(val)
    } else {
        data.OrganizationName = types.StringNull()
    }
    if obj, ok := item["repositoryName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RepositoryName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RepositoryName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RepositoryName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RepositoryName = types.StringValue(string(jsonBytes))
        } else {
            data.RepositoryName = types.StringNull()
        }
    } else if val, ok := item["repositoryName"].(string); ok {
        data.RepositoryName = types.StringValue(val)
    } else {
        data.RepositoryName = types.StringNull()
    }
    if obj, ok := item["mainBranchName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MainBranchName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.MainBranchName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.MainBranchName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.MainBranchName = types.StringValue(string(jsonBytes))
        } else {
            data.MainBranchName = types.StringNull()
        }
    } else if val, ok := item["mainBranchName"].(string); ok {
        data.MainBranchName = types.StringValue(val)
    } else {
        data.MainBranchName = types.StringNull()
    }
    if obj, ok := item["setupCommand"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SetupCommand = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SetupCommand = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SetupCommand = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SetupCommand = types.StringValue(string(jsonBytes))
        } else {
            data.SetupCommand = types.StringNull()
        }
    } else if val, ok := item["setupCommand"].(string); ok {
        data.SetupCommand = types.StringValue(val)
    } else {
        data.SetupCommand = types.StringNull()
    }
    if obj, ok := item["buildCommand"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.BuildCommand = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.BuildCommand = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.BuildCommand = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.BuildCommand = types.StringValue(string(jsonBytes))
        } else {
            data.BuildCommand = types.StringNull()
        }
    } else if val, ok := item["buildCommand"].(string); ok {
        data.BuildCommand = types.StringValue(val)
    } else {
        data.BuildCommand = types.StringNull()
    }
    if obj, ok := item["testCommand"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TestCommand = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TestCommand = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TestCommand = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TestCommand = types.StringValue(string(jsonBytes))
        } else {
            data.TestCommand = types.StringNull()
        }
    } else if val, ok := item["testCommand"].(string); ok {
        data.TestCommand = types.StringValue(val)
    } else {
        data.TestCommand = types.StringNull()
    }
    if val, ok := item["maxOpenFixPullRequests"].(float64); ok {
        data.MaxOpenFixPullRequests = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["maxOpenFixPullRequests"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.MaxOpenFixPullRequests = types.NumberValue(big.NewFloat(val))
        } else {
            data.MaxOpenFixPullRequests = types.NumberNull()
        }
    } else {
        data.MaxOpenFixPullRequests = types.NumberNull()
    }
    if obj, ok := item["repositoryUrl"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RepositoryUrl = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RepositoryUrl = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RepositoryUrl = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RepositoryUrl = types.StringValue(string(jsonBytes))
        } else {
            data.RepositoryUrl = types.StringNull()
        }
    } else if val, ok := item["repositoryUrl"].(string); ok {
        data.RepositoryUrl = types.StringValue(val)
    } else {
        data.RepositoryUrl = types.StringNull()
    }
    if obj, ok := item["gitHubAppInstallationId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.GitHubAppInstallationId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.GitHubAppInstallationId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.GitHubAppInstallationId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.GitHubAppInstallationId = types.StringValue(string(jsonBytes))
        } else {
            data.GitHubAppInstallationId = types.StringNull()
        }
    } else if val, ok := item["gitHubAppInstallationId"].(string); ok {
        data.GitHubAppInstallationId = types.StringValue(val)
    } else {
        data.GitHubAppInstallationId = types.StringNull()
    }
    if obj, ok := item["gitLabProjectId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.GitLabProjectId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.GitLabProjectId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.GitLabProjectId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.GitLabProjectId = types.StringValue(string(jsonBytes))
        } else {
            data.GitLabProjectId = types.StringNull()
        }
    } else if val, ok := item["gitLabProjectId"].(string); ok {
        data.GitLabProjectId = types.StringValue(val)
    } else {
        data.GitLabProjectId = types.StringNull()
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
    if obj, ok := item["deletedByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeletedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DeletedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DeletedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DeletedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.DeletedByUserId = types.StringNull()
        }
    } else if val, ok := item["deletedByUserId"].(string); ok {
        data.DeletedByUserId = types.StringValue(val)
    } else {
        data.DeletedByUserId = types.StringNull()
    }
    if val, ok := item["labels"].([]interface{}); ok {
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if jsonBytes, err := json.Marshal(itemMap); err == nil {
                    setItems = append(setItems, types.StringValue(string(jsonBytes)))
                }
            } else if str, ok := item.(string); ok {
                setItems = append(setItems, types.StringValue(str))
            } else {
                setItems = append(setItems, types.StringValue(fmt.Sprintf("%v", item)))
            }
        }
        sort.Slice(setItems, func(i, j int) bool {
            return setItems[i].(types.String).ValueString() < setItems[j].(types.String).ValueString()
        })
        data.Labels = types.SetValueMust(types.StringType, setItems)
    } else {
        data.Labels = types.SetNull(types.StringType)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
