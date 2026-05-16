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
var _ datasource.DataSource = &RunbookAgentUserOwnerDataDataSource{}

func NewRunbookAgentUserOwnerDataDataSource() datasource.DataSource {
    return &RunbookAgentUserOwnerDataDataSource{}
}

// RunbookAgentUserOwnerDataDataSource defines the data source implementation.
type RunbookAgentUserOwnerDataDataSource struct {
    client *Client
}

// RunbookAgentUserOwnerDataDataSourceModel describes the data source data model.
type RunbookAgentUserOwnerDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    RunbookAgentId types.String `tfsdk:"runbook_agent_id"`
    UserId types.String `tfsdk:"user_id"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
    IsOwnerNotified types.Bool `tfsdk:"is_owner_notified"`
}

func (d *RunbookAgentUserOwnerDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_runbook_agent_user_owner_data"
}

func (d *RunbookAgentUserOwnerDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "runbook_agent_user_owner_data data source",

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
            "runbook_agent_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "user_id": schema.StringAttribute{
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
            "is_owner_notified": schema.BoolAttribute{
                MarkdownDescription: "Are owners notified of this resource ownership?. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Runbook Admin, Runbook Member, Runbook Viewer, Read Runbook Agent User Owner, Read All Project Resources], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
        },
    }
}

func (d *RunbookAgentUserOwnerDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *RunbookAgentUserOwnerDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data RunbookAgentUserOwnerDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "runbook-agent-owner-user" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read runbook_agent_user_owner_data, got error: %s", err))
        return
    }

    var runbookAgentUserOwnerDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &runbookAgentUserOwnerDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse runbook_agent_user_owner_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := runbookAgentUserOwnerDataResponse["data"].(map[string]interface{}); ok {
        runbookAgentUserOwnerDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := runbookAgentUserOwnerDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := runbookAgentUserOwnerDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := runbookAgentUserOwnerDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := runbookAgentUserOwnerDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := runbookAgentUserOwnerDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := runbookAgentUserOwnerDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := runbookAgentUserOwnerDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := runbookAgentUserOwnerDataResponse["runbook_agent_id"].(string); ok {
        data.RunbookAgentId = types.StringValue(val)
    }
    if val, ok := runbookAgentUserOwnerDataResponse["user_id"].(string); ok {
        data.UserId = types.StringValue(val)
    }
    if val, ok := runbookAgentUserOwnerDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := runbookAgentUserOwnerDataResponse["deleted_by_user_id"].(string); ok {
        data.DeletedByUserId = types.StringValue(val)
    }
    if val, ok := runbookAgentUserOwnerDataResponse["is_owner_notified"].(bool); ok {
        data.IsOwnerNotified = types.BoolValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
