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
var _ datasource.DataSource = &ServiceCodeRepositoryDataDataSource{}

func NewServiceCodeRepositoryDataDataSource() datasource.DataSource {
    return &ServiceCodeRepositoryDataDataSource{}
}

// ServiceCodeRepositoryDataDataSource defines the data source implementation.
type ServiceCodeRepositoryDataDataSource struct {
    client *Client
}

// ServiceCodeRepositoryDataDataSourceModel describes the data source data model.
type ServiceCodeRepositoryDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    ServiceId types.String `tfsdk:"service_id"`
    CodeRepositoryId types.String `tfsdk:"code_repository_id"`
    ServicePathInRepository types.String `tfsdk:"service_path_in_repository"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
    EnableAutomaticImprovements types.Bool `tfsdk:"enable_automatic_improvements"`
    MaxOpenPullRequests types.Number `tfsdk:"max_open_pull_requests"`
    RestrictedImprovementActions types.String `tfsdk:"restricted_improvement_actions"`
}

func (d *ServiceCodeRepositoryDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_service_code_repository_data"
}

func (d *ServiceCodeRepositoryDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "service_code_repository_data data source",

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
            "service_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "code_repository_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "service_path_in_repository": schema.StringAttribute{
                MarkdownDescription: "The path in the repository where the service code lives (e.g., /services/api or /src/backend). Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Manager, Create Service Code Repository], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Manager, Read Service Code Repository, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Settings Manager, Edit Service Code Repository]",
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
            "enable_automatic_improvements": schema.BoolAttribute{
                MarkdownDescription: "Enable OneUptime to automatically create pull requests to improve the code for this service.. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Manager, Create Service Code Repository], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Manager, Read Service Code Repository, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Settings Manager, Edit Service Code Repository]",
                Computed: true,
            },
            "max_open_pull_requests": schema.NumberAttribute{
                MarkdownDescription: "Maximum number of open pull requests that OneUptime can create for this service at any given time.. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Manager, Create Service Code Repository], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Manager, Read Service Code Repository, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Settings Manager, Edit Service Code Repository]",
                Computed: true,
            },
            "restricted_improvement_actions": schema.StringAttribute{
                MarkdownDescription: "Restrict code improvements to only these actions. If empty, all improvement actions are allowed.. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Manager, Create Service Code Repository], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Manager, Read Service Code Repository, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Settings Manager, Edit Service Code Repository]",
                Computed: true,
            },
        },
    }
}

func (d *ServiceCodeRepositoryDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ServiceCodeRepositoryDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data ServiceCodeRepositoryDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "service-code-repository" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read service_code_repository_data, got error: %s", err))
        return
    }

    var serviceCodeRepositoryDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &serviceCodeRepositoryDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse service_code_repository_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := serviceCodeRepositoryDataResponse["data"].(map[string]interface{}); ok {
        serviceCodeRepositoryDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := serviceCodeRepositoryDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := serviceCodeRepositoryDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := serviceCodeRepositoryDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := serviceCodeRepositoryDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := serviceCodeRepositoryDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := serviceCodeRepositoryDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := serviceCodeRepositoryDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := serviceCodeRepositoryDataResponse["service_id"].(string); ok {
        data.ServiceId = types.StringValue(val)
    }
    if val, ok := serviceCodeRepositoryDataResponse["code_repository_id"].(string); ok {
        data.CodeRepositoryId = types.StringValue(val)
    }
    if val, ok := serviceCodeRepositoryDataResponse["service_path_in_repository"].(string); ok {
        data.ServicePathInRepository = types.StringValue(val)
    }
    if val, ok := serviceCodeRepositoryDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := serviceCodeRepositoryDataResponse["deleted_by_user_id"].(string); ok {
        data.DeletedByUserId = types.StringValue(val)
    }
    if val, ok := serviceCodeRepositoryDataResponse["enable_automatic_improvements"].(bool); ok {
        data.EnableAutomaticImprovements = types.BoolValue(val)
    }
    if val, ok := serviceCodeRepositoryDataResponse["max_open_pull_requests"].(float64); ok {
        data.MaxOpenPullRequests = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := serviceCodeRepositoryDataResponse["restricted_improvement_actions"].(string); ok {
        data.RestrictedImprovementActions = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
