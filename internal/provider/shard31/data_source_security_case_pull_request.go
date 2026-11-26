package shard31

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
    "context"
    "encoding/json"
    "fmt"
    "net/http"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &SecurityCasePullRequestDataSource{}

func NewSecurityCasePullRequestDataSource() datasource.DataSource {
    return &SecurityCasePullRequestDataSource{}
}

// SecurityCasePullRequestDataSource defines the data source implementation.
type SecurityCasePullRequestDataSource struct {
    client *Client
}

// SecurityCasePullRequestDataSourceModel describes the data source data model.
type SecurityCasePullRequestDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    CreatedByActorType types.String `tfsdk:"created_by_actor_type"`
    CreatedByActorId types.String `tfsdk:"created_by_actor_id"`
    LastModifiedAt types.String `tfsdk:"last_modified_at"`
    LastModifiedByActorType types.String `tfsdk:"last_modified_by_actor_type"`
    LastModifiedByActorId types.String `tfsdk:"last_modified_by_actor_id"`
    LastModifiedByUserId types.String `tfsdk:"last_modified_by_user_id"`
    ProjectId types.String `tfsdk:"project_id"`
    SecurityCaseId types.String `tfsdk:"security_case_id"`
    PullRequestUrl types.String `tfsdk:"pull_request_url"`
    CodeRepositoryId types.String `tfsdk:"code_repository_id"`
    PullRequestNumber types.Number `tfsdk:"pull_request_number"`
    RepoOrganizationName types.String `tfsdk:"repo_organization_name"`
    RepoName types.String `tfsdk:"repo_name"`
    PullRequestState types.String `tfsdk:"pull_request_state"`
    MergeCommitSha types.String `tfsdk:"merge_commit_sha"`
    MergedAt types.String `tfsdk:"merged_at"`
    DeployedAt types.String `tfsdk:"deployed_at"`
    DeployedVersion types.String `tfsdk:"deployed_version"`
    DeployedCommitSha types.String `tfsdk:"deployed_commit_sha"`
    FollowedAt types.String `tfsdk:"followed_at"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *SecurityCasePullRequestDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_security_case_pull_request"
}

func (d *SecurityCasePullRequestDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "A pull request a person linked to a security case as its fix, which OneUptime follows to its deploy. Look up an existing security case pull request by `id`, or by any of its other arguments (`code_repository_id`, `created_by_actor_id`, `created_by_actor_type`, ...): each one set must match, and exactly one security case pull request may match them all.",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Look up by unique identifier. Leave unset to look up by the other arguments instead.",
                Optional: true,
                Computed: true,
            },
            "created_at": schema.StringAttribute{
                MarkdownDescription: "Date and Time when the object was created.",
                Computed: true,
            },
            "updated_at": schema.StringAttribute{
                MarkdownDescription: "Date and Time when the object was updated.",
                Computed: true,
            },
            "created_by_actor_type": schema.StringAttribute{
                MarkdownDescription: "What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.",
                Optional: true,
                Computed: true,
            },
            "created_by_actor_id": schema.StringAttribute{
                MarkdownDescription: "ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.",
                Optional: true,
                Computed: true,
            },
            "last_modified_at": schema.StringAttribute{
                MarkdownDescription: "When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.",
                Computed: true,
            },
            "last_modified_by_actor_type": schema.StringAttribute{
                MarkdownDescription: "What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.",
                Optional: true,
                Computed: true,
            },
            "last_modified_by_actor_id": schema.StringAttribute{
                MarkdownDescription: "ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.",
                Optional: true,
                Computed: true,
            },
            "last_modified_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: DataSourceRecordIDValidators("oneuptime_user"),
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "ID of the project this pull request belongs to. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "security_case_id": schema.StringAttribute{
                MarkdownDescription: "ID of the case this pull request fixes. The ID of a `oneuptime_security_case`.",
                Optional: true,
                Computed: true,
                Validators: RecordIDValidators("oneuptime_security_case"),
            },
            "pull_request_url": schema.StringAttribute{
                MarkdownDescription: "The pull request's page, on a repository the project connected under Code: GitHub, GitLab, Forgejo or Gitea.",
                Optional: true,
                Computed: true,
            },
            "code_repository_id": schema.StringAttribute{
                MarkdownDescription: "ID of the repository the pull request is in. Found by OneUptime from the address. The ID of a `oneuptime_code_repository`.",
                Optional: true,
                Computed: true,
                Validators: RecordIDValidators("oneuptime_code_repository"),
            },
            "pull_request_number": schema.NumberAttribute{
                MarkdownDescription: "The pull request's number in its repository. Written by OneUptime.",
                Optional: true,
                Computed: true,
            },
            "repo_organization_name": schema.StringAttribute{
                MarkdownDescription: "Organization, user or group that owns the repository. Written by OneUptime.",
                Optional: true,
                Computed: true,
            },
            "repo_name": schema.StringAttribute{
                MarkdownDescription: "Name of the repository. Written by OneUptime.",
                Optional: true,
                Computed: true,
            },
            "pull_request_state": schema.StringAttribute{
                MarkdownDescription: "Whether the pull request is open, merged or closed, as OneUptime last read it from its host.",
                Optional: true,
                Computed: true,
            },
            "merge_commit_sha": schema.StringAttribute{
                MarkdownDescription: "The commit the pull request merged as: what a deploy must contain to ship it. Written by OneUptime.",
                Optional: true,
                Computed: true,
            },
            "merged_at": schema.StringAttribute{
                MarkdownDescription: "When the pull request merged. Written by OneUptime.",
                Computed: true,
            },
            "deployed_at": schema.StringAttribute{
                MarkdownDescription: "When a recorded deploy that contains the pull request's merge commit started. Written by OneUptime.",
                Computed: true,
            },
            "deployed_version": schema.StringAttribute{
                MarkdownDescription: "The version that deploy shipped, when it named one. Written by OneUptime.",
                Optional: true,
                Computed: true,
            },
            "deployed_commit_sha": schema.StringAttribute{
                MarkdownDescription: "The commit that deploy ran: the merge commit or a later one. Written by OneUptime.",
                Optional: true,
                Computed: true,
            },
            "followed_at": schema.StringAttribute{
                MarkdownDescription: "When OneUptime last read the pull request's state or looked for its deploy.",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the person who linked the pull request to the case. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: DataSourceRecordIDValidators("oneuptime_user"),
            },
        },
    }
}

func (d *SecurityCasePullRequestDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SecurityCasePullRequestDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data SecurityCasePullRequestDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.CreatedByActorType.IsNull() && !data.CreatedByActorType.IsUnknown() {
        filters["createdByActorType"] = data.CreatedByActorType.ValueString()
        filterNames = append(filterNames, "created_by_actor_type = "+fmt.Sprintf("%q", data.CreatedByActorType.ValueString()))
    }
    if !data.CreatedByActorId.IsNull() && !data.CreatedByActorId.IsUnknown() {
        filters["createdByActorId"] = data.CreatedByActorId.ValueString()
        filterNames = append(filterNames, "created_by_actor_id = "+fmt.Sprintf("%q", data.CreatedByActorId.ValueString()))
    }
    if !data.LastModifiedByActorType.IsNull() && !data.LastModifiedByActorType.IsUnknown() {
        filters["lastModifiedByActorType"] = data.LastModifiedByActorType.ValueString()
        filterNames = append(filterNames, "last_modified_by_actor_type = "+fmt.Sprintf("%q", data.LastModifiedByActorType.ValueString()))
    }
    if !data.LastModifiedByActorId.IsNull() && !data.LastModifiedByActorId.IsUnknown() {
        filters["lastModifiedByActorId"] = data.LastModifiedByActorId.ValueString()
        filterNames = append(filterNames, "last_modified_by_actor_id = "+fmt.Sprintf("%q", data.LastModifiedByActorId.ValueString()))
    }
    if !data.LastModifiedByUserId.IsNull() && !data.LastModifiedByUserId.IsUnknown() {
        filters["lastModifiedByUserId"] = data.LastModifiedByUserId.ValueString()
        filterNames = append(filterNames, "last_modified_by_user_id = "+fmt.Sprintf("%q", data.LastModifiedByUserId.ValueString()))
    }
    if !data.SecurityCaseId.IsNull() && !data.SecurityCaseId.IsUnknown() {
        filters["securityCaseId"] = data.SecurityCaseId.ValueString()
        filterNames = append(filterNames, "security_case_id = "+fmt.Sprintf("%q", data.SecurityCaseId.ValueString()))
    }
    if !data.PullRequestUrl.IsNull() && !data.PullRequestUrl.IsUnknown() {
        filters["pullRequestUrl"] = data.PullRequestUrl.ValueString()
        filterNames = append(filterNames, "pull_request_url = "+fmt.Sprintf("%q", data.PullRequestUrl.ValueString()))
    }
    if !data.CodeRepositoryId.IsNull() && !data.CodeRepositoryId.IsUnknown() {
        filters["codeRepositoryId"] = data.CodeRepositoryId.ValueString()
        filterNames = append(filterNames, "code_repository_id = "+fmt.Sprintf("%q", data.CodeRepositoryId.ValueString()))
    }
    if !data.PullRequestNumber.IsNull() && !data.PullRequestNumber.IsUnknown() {
        filters["pullRequestNumber"] = LookupNumber(data.PullRequestNumber)
        filterNames = append(filterNames, "pull_request_number = "+data.PullRequestNumber.ValueBigFloat().String())
    }
    if !data.RepoOrganizationName.IsNull() && !data.RepoOrganizationName.IsUnknown() {
        filters["repoOrganizationName"] = data.RepoOrganizationName.ValueString()
        filterNames = append(filterNames, "repo_organization_name = "+fmt.Sprintf("%q", data.RepoOrganizationName.ValueString()))
    }
    if !data.RepoName.IsNull() && !data.RepoName.IsUnknown() {
        filters["repoName"] = data.RepoName.ValueString()
        filterNames = append(filterNames, "repo_name = "+fmt.Sprintf("%q", data.RepoName.ValueString()))
    }
    if !data.PullRequestState.IsNull() && !data.PullRequestState.IsUnknown() {
        filters["pullRequestState"] = data.PullRequestState.ValueString()
        filterNames = append(filterNames, "pull_request_state = "+fmt.Sprintf("%q", data.PullRequestState.ValueString()))
    }
    if !data.MergeCommitSha.IsNull() && !data.MergeCommitSha.IsUnknown() {
        filters["mergeCommitSha"] = data.MergeCommitSha.ValueString()
        filterNames = append(filterNames, "merge_commit_sha = "+fmt.Sprintf("%q", data.MergeCommitSha.ValueString()))
    }
    if !data.DeployedVersion.IsNull() && !data.DeployedVersion.IsUnknown() {
        filters["deployedVersion"] = data.DeployedVersion.ValueString()
        filterNames = append(filterNames, "deployed_version = "+fmt.Sprintf("%q", data.DeployedVersion.ValueString()))
    }
    if !data.DeployedCommitSha.IsNull() && !data.DeployedCommitSha.IsUnknown() {
        filters["deployedCommitSha"] = data.DeployedCommitSha.ValueString()
        filterNames = append(filterNames, "deployed_commit_sha = "+fmt.Sprintf("%q", data.DeployedCommitSha.ValueString()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the security case pull request up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the security case pull request up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "createdByActorType": true,
        "createdByActorId": true,
        "lastModifiedAt": true,
        "lastModifiedByActorType": true,
        "lastModifiedByActorId": true,
        "lastModifiedByUserId": true,
        "projectId": true,
        "securityCaseId": true,
        "pullRequestUrl": true,
        "codeRepositoryId": true,
        "pullRequestNumber": true,
        "repoOrganizationName": true,
        "repoName": true,
        "pullRequestState": true,
        "mergeCommitSha": true,
        "mergedAt": true,
        "deployedAt": true,
        "deployedVersion": true,
        "deployedCommitSha": true,
        "followedAt": true,
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/security-case-pull-request/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read security_case_pull_request, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No security case pull request found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read security_case_pull_request: %s", err))
            return
        }
        if wrapper, ok := itemResponse["data"].(map[string]interface{}); ok {
            item = wrapper
        } else {
            item = itemResponse
        }
    }
    if !hasId {
        listBody := map[string]interface{}{
            "query":  filters,
            "select": selectParam,
            // limit 2 is enough to detect ambiguity without paging.
            "limit": 2,
        }
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/security-case-pull-request/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list security_case_pull_request, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list security_case_pull_request: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No security case pull request matches %s.", DescribeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one security case pull request matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", DescribeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for security_case_pull_request.")
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
    } else if arr, ok := item["_id"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Id = types.StringValue(string(jsonBytes))
        } else {
            data.Id = types.StringNull()
        }
    } else if val, ok := item["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
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
    } else if arr, ok := item["createdAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
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
    } else if arr, ok := item["updatedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.UpdatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.UpdatedAt = types.StringNull()
        }
    } else if val, ok := item["updatedAt"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    } else {
        data.UpdatedAt = types.StringNull()
    }
    if obj, ok := item["createdByActorType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByActorType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CreatedByActorType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CreatedByActorType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CreatedByActorType = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByActorType = types.StringNull()
        }
    } else if arr, ok := item["createdByActorType"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.CreatedByActorType = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByActorType = types.StringNull()
        }
    } else if val, ok := item["createdByActorType"].(string); ok {
        data.CreatedByActorType = types.StringValue(val)
    } else {
        data.CreatedByActorType = types.StringNull()
    }
    if obj, ok := item["createdByActorId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByActorId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CreatedByActorId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CreatedByActorId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CreatedByActorId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByActorId = types.StringNull()
        }
    } else if arr, ok := item["createdByActorId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.CreatedByActorId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByActorId = types.StringNull()
        }
    } else if val, ok := item["createdByActorId"].(string); ok {
        data.CreatedByActorId = types.StringValue(val)
    } else {
        data.CreatedByActorId = types.StringNull()
    }
    if obj, ok := item["lastModifiedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastModifiedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastModifiedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastModifiedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastModifiedAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedAt = types.StringNull()
        }
    } else if arr, ok := item["lastModifiedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LastModifiedAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedAt = types.StringNull()
        }
    } else if val, ok := item["lastModifiedAt"].(string); ok {
        data.LastModifiedAt = types.StringValue(val)
    } else {
        data.LastModifiedAt = types.StringNull()
    }
    if obj, ok := item["lastModifiedByActorType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastModifiedByActorType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastModifiedByActorType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastModifiedByActorType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastModifiedByActorType = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByActorType = types.StringNull()
        }
    } else if arr, ok := item["lastModifiedByActorType"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LastModifiedByActorType = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByActorType = types.StringNull()
        }
    } else if val, ok := item["lastModifiedByActorType"].(string); ok {
        data.LastModifiedByActorType = types.StringValue(val)
    } else {
        data.LastModifiedByActorType = types.StringNull()
    }
    if obj, ok := item["lastModifiedByActorId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastModifiedByActorId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastModifiedByActorId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastModifiedByActorId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastModifiedByActorId = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByActorId = types.StringNull()
        }
    } else if arr, ok := item["lastModifiedByActorId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LastModifiedByActorId = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByActorId = types.StringNull()
        }
    } else if val, ok := item["lastModifiedByActorId"].(string); ok {
        data.LastModifiedByActorId = types.StringValue(val)
    } else {
        data.LastModifiedByActorId = types.StringNull()
    }
    if obj, ok := item["lastModifiedByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastModifiedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastModifiedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastModifiedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastModifiedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByUserId = types.StringNull()
        }
    } else if arr, ok := item["lastModifiedByUserId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LastModifiedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByUserId = types.StringNull()
        }
    } else if val, ok := item["lastModifiedByUserId"].(string); ok {
        data.LastModifiedByUserId = types.StringValue(val)
    } else {
        data.LastModifiedByUserId = types.StringNull()
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
    } else if arr, ok := item["projectId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.ProjectId = types.StringValue(string(jsonBytes))
        } else {
            data.ProjectId = types.StringNull()
        }
    } else if val, ok := item["projectId"].(string); ok {
        data.ProjectId = types.StringValue(val)
    } else {
        data.ProjectId = types.StringNull()
    }
    if obj, ok := item["securityCaseId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SecurityCaseId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SecurityCaseId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SecurityCaseId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SecurityCaseId = types.StringValue(string(jsonBytes))
        } else {
            data.SecurityCaseId = types.StringNull()
        }
    } else if arr, ok := item["securityCaseId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.SecurityCaseId = types.StringValue(string(jsonBytes))
        } else {
            data.SecurityCaseId = types.StringNull()
        }
    } else if val, ok := item["securityCaseId"].(string); ok {
        data.SecurityCaseId = types.StringValue(val)
    } else {
        data.SecurityCaseId = types.StringNull()
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
    } else if arr, ok := item["pullRequestUrl"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.PullRequestUrl = types.StringValue(string(jsonBytes))
        } else {
            data.PullRequestUrl = types.StringNull()
        }
    } else if val, ok := item["pullRequestUrl"].(string); ok {
        data.PullRequestUrl = types.StringValue(val)
    } else {
        data.PullRequestUrl = types.StringNull()
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
    } else if arr, ok := item["codeRepositoryId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.CodeRepositoryId = types.StringValue(string(jsonBytes))
        } else {
            data.CodeRepositoryId = types.StringNull()
        }
    } else if val, ok := item["codeRepositoryId"].(string); ok {
        data.CodeRepositoryId = types.StringValue(val)
    } else {
        data.CodeRepositoryId = types.StringNull()
    }
    data.PullRequestNumber = NumberFromAPI(data.PullRequestNumber, item["pullRequestNumber"])
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
    } else if arr, ok := item["repoOrganizationName"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
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
    } else if arr, ok := item["repoName"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.RepoName = types.StringValue(string(jsonBytes))
        } else {
            data.RepoName = types.StringNull()
        }
    } else if val, ok := item["repoName"].(string); ok {
        data.RepoName = types.StringValue(val)
    } else {
        data.RepoName = types.StringNull()
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
    } else if arr, ok := item["pullRequestState"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.PullRequestState = types.StringValue(string(jsonBytes))
        } else {
            data.PullRequestState = types.StringNull()
        }
    } else if val, ok := item["pullRequestState"].(string); ok {
        data.PullRequestState = types.StringValue(val)
    } else {
        data.PullRequestState = types.StringNull()
    }
    if obj, ok := item["mergeCommitSha"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MergeCommitSha = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.MergeCommitSha = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.MergeCommitSha = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.MergeCommitSha = types.StringValue(string(jsonBytes))
        } else {
            data.MergeCommitSha = types.StringNull()
        }
    } else if arr, ok := item["mergeCommitSha"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.MergeCommitSha = types.StringValue(string(jsonBytes))
        } else {
            data.MergeCommitSha = types.StringNull()
        }
    } else if val, ok := item["mergeCommitSha"].(string); ok {
        data.MergeCommitSha = types.StringValue(val)
    } else {
        data.MergeCommitSha = types.StringNull()
    }
    if obj, ok := item["mergedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MergedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.MergedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.MergedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.MergedAt = types.StringValue(string(jsonBytes))
        } else {
            data.MergedAt = types.StringNull()
        }
    } else if arr, ok := item["mergedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.MergedAt = types.StringValue(string(jsonBytes))
        } else {
            data.MergedAt = types.StringNull()
        }
    } else if val, ok := item["mergedAt"].(string); ok {
        data.MergedAt = types.StringValue(val)
    } else {
        data.MergedAt = types.StringNull()
    }
    if obj, ok := item["deployedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeployedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DeployedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DeployedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DeployedAt = types.StringValue(string(jsonBytes))
        } else {
            data.DeployedAt = types.StringNull()
        }
    } else if arr, ok := item["deployedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.DeployedAt = types.StringValue(string(jsonBytes))
        } else {
            data.DeployedAt = types.StringNull()
        }
    } else if val, ok := item["deployedAt"].(string); ok {
        data.DeployedAt = types.StringValue(val)
    } else {
        data.DeployedAt = types.StringNull()
    }
    if obj, ok := item["deployedVersion"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeployedVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DeployedVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DeployedVersion = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DeployedVersion = types.StringValue(string(jsonBytes))
        } else {
            data.DeployedVersion = types.StringNull()
        }
    } else if arr, ok := item["deployedVersion"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.DeployedVersion = types.StringValue(string(jsonBytes))
        } else {
            data.DeployedVersion = types.StringNull()
        }
    } else if val, ok := item["deployedVersion"].(string); ok {
        data.DeployedVersion = types.StringValue(val)
    } else {
        data.DeployedVersion = types.StringNull()
    }
    if obj, ok := item["deployedCommitSha"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeployedCommitSha = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DeployedCommitSha = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DeployedCommitSha = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DeployedCommitSha = types.StringValue(string(jsonBytes))
        } else {
            data.DeployedCommitSha = types.StringNull()
        }
    } else if arr, ok := item["deployedCommitSha"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.DeployedCommitSha = types.StringValue(string(jsonBytes))
        } else {
            data.DeployedCommitSha = types.StringNull()
        }
    } else if val, ok := item["deployedCommitSha"].(string); ok {
        data.DeployedCommitSha = types.StringValue(val)
    } else {
        data.DeployedCommitSha = types.StringNull()
    }
    if obj, ok := item["followedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FollowedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.FollowedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.FollowedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.FollowedAt = types.StringValue(string(jsonBytes))
        } else {
            data.FollowedAt = types.StringNull()
        }
    } else if arr, ok := item["followedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.FollowedAt = types.StringValue(string(jsonBytes))
        } else {
            data.FollowedAt = types.StringNull()
        }
    } else if val, ok := item["followedAt"].(string); ok {
        data.FollowedAt = types.StringValue(val)
    } else {
        data.FollowedAt = types.StringNull()
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
    } else if arr, ok := item["createdByUserId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
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
