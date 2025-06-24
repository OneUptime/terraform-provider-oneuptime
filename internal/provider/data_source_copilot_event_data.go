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
var _ datasource.DataSource = &CopilotEventDataDataSource{}

func NewCopilotEventDataDataSource() datasource.DataSource {
    return &CopilotEventDataDataSource{}
}

// CopilotEventDataDataSource defines the data source implementation.
type CopilotEventDataDataSource struct {
    client *Client
}

// CopilotEventDataDataSourceModel describes the data source data model.
type CopilotEventDataDataSourceModel struct {
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
    CommitHash types.String `tfsdk:"commit_hash"`
    CopilotActionType types.String `tfsdk:"copilot_action_type"`
    ServiceCatalogId types.String `tfsdk:"service_catalog_id"`
    ServiceRepositoryId types.String `tfsdk:"service_repository_id"`
    CopilotPullRequestId types.String `tfsdk:"copilot_pull_request_id"`
    CopilotActionStatus types.String `tfsdk:"copilot_action_status"`
    CopilotActionProp types.String `tfsdk:"copilot_action_prop"`
    StatusMessage types.String `tfsdk:"status_message"`
    Logs types.String `tfsdk:"logs"`
    IsPriority types.Bool `tfsdk:"is_priority"`
    StatusChangedAt types.String `tfsdk:"status_changed_at"`
}

func (d *CopilotEventDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_copilot_event_data"
}

func (d *CopilotEventDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "copilot_event_data data source",

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
            "commit_hash": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Copilot Event], Read: [Project Owner, Project Admin, Project Member, Read Copilot Event], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "copilot_action_type": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Copilot Event], Read: [Project Owner, Project Admin, Project Member, Read Copilot Event], Update: [No access - you don't have permission for this operation]",
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
            "copilot_pull_request_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "copilot_action_status": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Copilot Event], Read: [Project Owner, Project Admin, Project Member, Read Copilot Event], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "copilot_action_prop": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Copilot Event], Read: [Project Owner, Project Admin, Project Member, Read Copilot Event], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "status_message": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Copilot Event], Read: [Project Owner, Project Admin, Project Member, Read Copilot Event], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "logs": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Copilot Event], Read: [Project Owner, Project Admin, Project Member, Read Copilot Event], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "is_priority": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Copilot Event], Read: [Project Owner, Project Admin, Project Member, Read Copilot Event], Update: [Project Owner, Project Admin, Project Member, Edit Copilot Event]",
                Computed: true,
            },
            "status_changed_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
        },
    }
}

func (d *CopilotEventDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *CopilotEventDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data CopilotEventDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "copilot-action" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read copilot_event_data, got error: %s", err))
        return
    }

    var copilotEventDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &copilotEventDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse copilot_event_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := copilotEventDataResponse["data"].(map[string]interface{}); ok {
        copilotEventDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := copilotEventDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := copilotEventDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := copilotEventDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := copilotEventDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := copilotEventDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := copilotEventDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := copilotEventDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := copilotEventDataResponse["code_repository_id"].(string); ok {
        data.CodeRepositoryId = types.StringValue(val)
    }
    if val, ok := copilotEventDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := copilotEventDataResponse["deleted_by_user_id"].(string); ok {
        data.DeletedByUserId = types.StringValue(val)
    }
    if val, ok := copilotEventDataResponse["commit_hash"].(string); ok {
        data.CommitHash = types.StringValue(val)
    }
    if val, ok := copilotEventDataResponse["copilot_action_type"].(string); ok {
        data.CopilotActionType = types.StringValue(val)
    }
    if val, ok := copilotEventDataResponse["service_catalog_id"].(string); ok {
        data.ServiceCatalogId = types.StringValue(val)
    }
    if val, ok := copilotEventDataResponse["service_repository_id"].(string); ok {
        data.ServiceRepositoryId = types.StringValue(val)
    }
    if val, ok := copilotEventDataResponse["copilot_pull_request_id"].(string); ok {
        data.CopilotPullRequestId = types.StringValue(val)
    }
    if val, ok := copilotEventDataResponse["copilot_action_status"].(string); ok {
        data.CopilotActionStatus = types.StringValue(val)
    }
    if val, ok := copilotEventDataResponse["copilot_action_prop"].(string); ok {
        data.CopilotActionProp = types.StringValue(val)
    }
    if val, ok := copilotEventDataResponse["status_message"].(string); ok {
        data.StatusMessage = types.StringValue(val)
    }
    if val, ok := copilotEventDataResponse["logs"].(string); ok {
        data.Logs = types.StringValue(val)
    }
    if val, ok := copilotEventDataResponse["is_priority"].(bool); ok {
        data.IsPriority = types.BoolValue(val)
    }
    if val, ok := copilotEventDataResponse["status_changed_at"].(string); ok {
        data.StatusChangedAt = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
