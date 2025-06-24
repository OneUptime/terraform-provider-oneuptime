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
var _ datasource.DataSource = &TeamPermissionDataDataSource{}

func NewTeamPermissionDataDataSource() datasource.DataSource {
    return &TeamPermissionDataDataSource{}
}

// TeamPermissionDataDataSource defines the data source implementation.
type TeamPermissionDataDataSource struct {
    client *Client
}

// TeamPermissionDataDataSourceModel describes the data source data model.
type TeamPermissionDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    TeamId types.String `tfsdk:"team_id"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    Permission types.String `tfsdk:"permission"`
    Labels types.List `tfsdk:"labels"`
    IsBlockPermission types.Bool `tfsdk:"is_block_permission"`
}

func (d *TeamPermissionDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_team_permission_data"
}

func (d *TeamPermissionDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "team_permission_data data source",

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
            "team_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "permission": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Create Team, Edit Team Permissions], Read: [Project Owner, Project Admin, Project Member, Read Teams], Update: [Project Owner, Project Admin, Invite New Members, Edit Team Permissions, Edit Team]",
                Computed: true,
            },
            "labels": schema.ListAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Create Team, Edit Team Permissions], Read: [Project Owner, Project Admin, Project Member, Read Teams], Update: [Project Owner, Project Admin, Edit Team Permissions, Edit Team]",
                Computed: true,
                ElementType: types.StringType,
            },
            "is_block_permission": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Create Team, Edit Team Permissions], Read: [Project Owner, Project Admin, Project Member, Read Teams], Update: [Project Owner, Project Admin, Edit Team Permissions, Edit Team]",
                Computed: true,
            },
        },
    }
}

func (d *TeamPermissionDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *TeamPermissionDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data TeamPermissionDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "team-permission" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read team_permission_data, got error: %s", err))
        return
    }

    var teamPermissionDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &teamPermissionDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse team_permission_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := teamPermissionDataResponse["data"].(map[string]interface{}); ok {
        teamPermissionDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := teamPermissionDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := teamPermissionDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := teamPermissionDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := teamPermissionDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := teamPermissionDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := teamPermissionDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := teamPermissionDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := teamPermissionDataResponse["team_id"].(string); ok {
        data.TeamId = types.StringValue(val)
    }
    if val, ok := teamPermissionDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := teamPermissionDataResponse["permission"].(string); ok {
        data.Permission = types.StringValue(val)
    }
    if val, ok := teamPermissionDataResponse["labels"].([]interface{}); ok {
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
    if val, ok := teamPermissionDataResponse["is_block_permission"].(bool); ok {
        data.IsBlockPermission = types.BoolValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
