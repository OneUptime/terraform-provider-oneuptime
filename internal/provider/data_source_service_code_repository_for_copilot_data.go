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
var _ datasource.DataSource = &ServiceCodeRepositoryForCopilotDataDataSource{}

func NewServiceCodeRepositoryForCopilotDataDataSource() datasource.DataSource {
    return &ServiceCodeRepositoryForCopilotDataDataSource{}
}

// ServiceCodeRepositoryForCopilotDataDataSource defines the data source implementation.
type ServiceCodeRepositoryForCopilotDataDataSource struct {
    client *Client
}

// ServiceCodeRepositoryForCopilotDataDataSourceModel describes the data source data model.
type ServiceCodeRepositoryForCopilotDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    ServicePathInRepository types.String `tfsdk:"service_path_in_repository"`
    LimitNumberOfOpenPullRequestsCount types.Number `tfsdk:"limit_number_of_open_pull_requests_count"`
    EnablePullRequests types.Bool `tfsdk:"enable_pull_requests"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
    CodeRepositoryId types.String `tfsdk:"code_repository_id"`
    ServiceCatalogId types.String `tfsdk:"service_catalog_id"`
}

func (d *ServiceCodeRepositoryForCopilotDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_service_code_repository_for_copilot_data"
}

func (d *ServiceCodeRepositoryForCopilotDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "service_code_repository_for_copilot_data data source",

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
            "service_path_in_repository": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Service Repository], Read: [Project Owner, Project Admin, Project Member, Project Member, Read Service Repository], Update: [Project Owner, Project Admin, Project Member, Edit Service Repository]",
                Computed: true,
            },
            "limit_number_of_open_pull_requests_count": schema.NumberAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Service Repository], Read: [Project Owner, Project Admin, Project Member, Project Member, Read Service Repository], Update: [Project Owner, Project Admin, Project Member, Edit Service Repository]",
                Computed: true,
            },
            "enable_pull_requests": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Service Repository], Read: [Project Owner, Project Admin, Project Member, Project Member, Read Service Repository], Update: [Project Owner, Project Admin, Project Member, Edit Service Repository]",
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
            "code_repository_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "service_catalog_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
        },
    }
}

func (d *ServiceCodeRepositoryForCopilotDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ServiceCodeRepositoryForCopilotDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data ServiceCodeRepositoryForCopilotDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "service-copilot-code-repository" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read service_code_repository_for_copilot_data, got error: %s", err))
        return
    }

    var serviceCodeRepositoryForCopilotDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &serviceCodeRepositoryForCopilotDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse service_code_repository_for_copilot_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := serviceCodeRepositoryForCopilotDataResponse["data"].(map[string]interface{}); ok {
        serviceCodeRepositoryForCopilotDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := serviceCodeRepositoryForCopilotDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := serviceCodeRepositoryForCopilotDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := serviceCodeRepositoryForCopilotDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := serviceCodeRepositoryForCopilotDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := serviceCodeRepositoryForCopilotDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := serviceCodeRepositoryForCopilotDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := serviceCodeRepositoryForCopilotDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := serviceCodeRepositoryForCopilotDataResponse["service_path_in_repository"].(string); ok {
        data.ServicePathInRepository = types.StringValue(val)
    }
    if val, ok := serviceCodeRepositoryForCopilotDataResponse["limit_number_of_open_pull_requests_count"].(float64); ok {
        data.LimitNumberOfOpenPullRequestsCount = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := serviceCodeRepositoryForCopilotDataResponse["enable_pull_requests"].(bool); ok {
        data.EnablePullRequests = types.BoolValue(val)
    }
    if val, ok := serviceCodeRepositoryForCopilotDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := serviceCodeRepositoryForCopilotDataResponse["deleted_by_user_id"].(string); ok {
        data.DeletedByUserId = types.StringValue(val)
    }
    if val, ok := serviceCodeRepositoryForCopilotDataResponse["code_repository_id"].(string); ok {
        data.CodeRepositoryId = types.StringValue(val)
    }
    if val, ok := serviceCodeRepositoryForCopilotDataResponse["service_catalog_id"].(string); ok {
        data.ServiceCatalogId = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
