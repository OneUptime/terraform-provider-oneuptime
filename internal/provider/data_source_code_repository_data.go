package provider

import (
    "context"
    "fmt"
    "math/big"
    "github.com/hashicorp/terraform-plugin-framework/attr"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &CodeRepositoryDataDataSource{}

func NewCodeRepositoryDataDataSource() datasource.DataSource {
    return &CodeRepositoryDataDataSource{}
}

// CodeRepositoryDataDataSource defines the data source implementation.
type CodeRepositoryDataDataSource struct {
    client *Client
}

// CodeRepositoryDataDataSourceModel describes the data source data model.
type CodeRepositoryDataDataSourceModel struct {
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
    RepositoryUrl types.String `tfsdk:"repository_url"`
    GitHubAppInstallationId types.String `tfsdk:"git_hub_app_installation_id"`
    GitLabProjectId types.String `tfsdk:"git_lab_project_id"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
    Labels types.List `tfsdk:"labels"`
}

func (d *CodeRepositoryDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_code_repository_data"
}

func (d *CodeRepositoryDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "code_repository_data data source",

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
            "slug": schema.StringAttribute{
                MarkdownDescription: "Friendly globally unique name for your object. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Code Repository], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "A description of this code repository. Permissions - Create: [Project Owner, Project Admin, Project Member, Create Code Repository], Read: [Project Owner, Project Admin, Project Member, Read Code Repository], Update: [Project Owner, Project Admin, Project Member, Edit Code Repository]",
                Computed: true,
            },
            "repository_hosted_at": schema.StringAttribute{
                MarkdownDescription: "Where is this repository hosted (GitHub, GitLab, etc.). Permissions - Create: [Project Owner, Project Admin, Project Member, Create Code Repository], Read: [Project Owner, Project Admin, Project Member, Read Code Repository], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "organization_name": schema.StringAttribute{
                MarkdownDescription: "GitHub organization or username that owns this repository. Permissions - Create: [Project Owner, Project Admin, Project Member, Create Code Repository], Read: [Project Owner, Project Admin, Project Member, Read Code Repository], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "repository_name": schema.StringAttribute{
                MarkdownDescription: "The name of the repository. Permissions - Create: [Project Owner, Project Admin, Project Member, Create Code Repository], Read: [Project Owner, Project Admin, Project Member, Read Code Repository], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "main_branch_name": schema.StringAttribute{
                MarkdownDescription: "The name of the main/default branch. Permissions - Create: [Project Owner, Project Admin, Project Member, Create Code Repository], Read: [Project Owner, Project Admin, Project Member, Read Code Repository], Update: [Project Owner, Project Admin, Project Member, Edit Code Repository]",
                Computed: true,
            },
            "repository_url": schema.StringAttribute{
                MarkdownDescription: "The HTTPS URL to the repository. Permissions - Create: [Project Owner, Project Admin, Project Member, Create Code Repository], Read: [Project Owner, Project Admin, Project Member, Read Code Repository], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "git_hub_app_installation_id": schema.StringAttribute{
                MarkdownDescription: "The GitHub App installation ID used to authenticate with this repository. Permissions - Create: [Project Owner, Project Admin, Create Code Repository], Read: [Project Owner, Project Admin, Read Code Repository], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "git_lab_project_id": schema.StringAttribute{
                MarkdownDescription: "The GitLab project ID for this repository. Permissions - Create: [Project Owner, Project Admin, Create Code Repository], Read: [Project Owner, Project Admin, Read Code Repository], Update: [No access - you don't have permission for this operation]",
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
            "labels": schema.ListAttribute{
                MarkdownDescription: "Relation to Labels Array where this object is categorized in.. Permissions - Create: [Project Owner, Project Admin, Project Member, Create Code Repository], Read: [Project Owner, Project Admin, Project Member, Read Code Repository], Update: [Project Owner, Project Admin, Project Member, Edit Code Repository]",
                Computed: true,
                ElementType: types.StringType,
            },
        },
    }
}

func (d *CodeRepositoryDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *CodeRepositoryDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data CodeRepositoryDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "code-repository" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read code_repository_data, got error: %s", err))
        return
    }

    var codeRepositoryDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &codeRepositoryDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse code_repository_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := codeRepositoryDataResponse["data"].(map[string]interface{}); ok {
        codeRepositoryDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := codeRepositoryDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := codeRepositoryDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := codeRepositoryDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := codeRepositoryDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := codeRepositoryDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := codeRepositoryDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := codeRepositoryDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := codeRepositoryDataResponse["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    }
    if val, ok := codeRepositoryDataResponse["description"].(string); ok {
        data.Description = types.StringValue(val)
    }
    if val, ok := codeRepositoryDataResponse["repository_hosted_at"].(string); ok {
        data.RepositoryHostedAt = types.StringValue(val)
    }
    if val, ok := codeRepositoryDataResponse["organization_name"].(string); ok {
        data.OrganizationName = types.StringValue(val)
    }
    if val, ok := codeRepositoryDataResponse["repository_name"].(string); ok {
        data.RepositoryName = types.StringValue(val)
    }
    if val, ok := codeRepositoryDataResponse["main_branch_name"].(string); ok {
        data.MainBranchName = types.StringValue(val)
    }
    if val, ok := codeRepositoryDataResponse["repository_url"].(string); ok {
        data.RepositoryUrl = types.StringValue(val)
    }
    if val, ok := codeRepositoryDataResponse["git_hub_app_installation_id"].(string); ok {
        data.GitHubAppInstallationId = types.StringValue(val)
    }
    if val, ok := codeRepositoryDataResponse["git_lab_project_id"].(string); ok {
        data.GitLabProjectId = types.StringValue(val)
    }
    if val, ok := codeRepositoryDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := codeRepositoryDataResponse["deleted_by_user_id"].(string); ok {
        data.DeletedByUserId = types.StringValue(val)
    }
    if val, ok := codeRepositoryDataResponse["labels"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        listValue, _ := types.ListValue(types.StringType, elements)
        data.Labels = listValue
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
