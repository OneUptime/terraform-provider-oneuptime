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
var _ datasource.DataSource = &CopilotPullRequestDataDataSource{}

func NewCopilotPullRequestDataDataSource() datasource.DataSource {
    return &CopilotPullRequestDataDataSource{}
}

// CopilotPullRequestDataDataSource defines the data source implementation.
type CopilotPullRequestDataDataSource struct {
    client *Client
}

// CopilotPullRequestDataDataSourceModel describes the data source data model.
type CopilotPullRequestDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    CodeRepositoryId types.String `tfsdk:"code_repository_id"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
    ServiceCatalogId types.String `tfsdk:"service_catalog_id"`
    ServiceRepositoryId types.String `tfsdk:"service_repository_id"`
    PullRequestId types.String `tfsdk:"pull_request_id"`
    CopilotPullRequestStatus types.String `tfsdk:"copilot_pull_request_status"`
    IsSetupPullRequest types.Bool `tfsdk:"is_setup_pull_request"`
}

func (d *CopilotPullRequestDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_copilot_pull_request_data"
}

func (d *CopilotPullRequestDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "copilot_pull_request_data data source",

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
                MarkdownDescription: "Version",
                Computed: true,
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "code_repository_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
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
            "service_catalog_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "service_repository_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "pull_request_id": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "copilot_pull_request_status": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "is_setup_pull_request": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
        },
    }
}

func (d *CopilotPullRequestDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *CopilotPullRequestDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data CopilotPullRequestDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "copilot-pull-request" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read copilot_pull_request_data, got error: %s", err))
        return
    }

    var copilotPullRequestDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &copilotPullRequestDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse copilot_pull_request_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := copilotPullRequestDataResponse["data"].(map[string]interface{}); ok {
        copilotPullRequestDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := copilotPullRequestDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := copilotPullRequestDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := copilotPullRequestDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := copilotPullRequestDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := copilotPullRequestDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := copilotPullRequestDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := copilotPullRequestDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := copilotPullRequestDataResponse["code_repository_id"].(string); ok {
        data.CodeRepositoryId = types.StringValue(val)
    }
    if val, ok := copilotPullRequestDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := copilotPullRequestDataResponse["deleted_by_user_id"].(string); ok {
        data.DeletedByUserId = types.StringValue(val)
    }
    if val, ok := copilotPullRequestDataResponse["service_catalog_id"].(string); ok {
        data.ServiceCatalogId = types.StringValue(val)
    }
    if val, ok := copilotPullRequestDataResponse["service_repository_id"].(string); ok {
        data.ServiceRepositoryId = types.StringValue(val)
    }
    if val, ok := copilotPullRequestDataResponse["pull_request_id"].(string); ok {
        data.PullRequestId = types.StringValue(val)
    }
    if val, ok := copilotPullRequestDataResponse["copilot_pull_request_status"].(string); ok {
        data.CopilotPullRequestStatus = types.StringValue(val)
    }
    if val, ok := copilotPullRequestDataResponse["is_setup_pull_request"].(bool); ok {
        data.IsSetupPullRequest = types.BoolValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
