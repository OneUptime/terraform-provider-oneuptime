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
var _ datasource.DataSource = &TeamDataDataSource{}

func NewTeamDataDataSource() datasource.DataSource {
    return &TeamDataDataSource{}
}

// TeamDataDataSource defines the data source implementation.
type TeamDataDataSource struct {
    client *Client
}

// TeamDataDataSourceModel describes the data source data model.
type TeamDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    Description types.String `tfsdk:"description"`
    Slug types.String `tfsdk:"slug"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    IsPermissionsEditable types.Bool `tfsdk:"is_permissions_editable"`
    IsTeamDeleteable types.Bool `tfsdk:"is_team_deleteable"`
    ShouldHaveAtLeastOneMember types.Bool `tfsdk:"should_have_at_least_one_member"`
    IsTeamEditable types.Bool `tfsdk:"is_team_editable"`
}

func (d *TeamDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_team_data"
}

func (d *TeamDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "team_data data source",

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
                MarkdownDescription: "Version",
                Computed: true,
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Create Team], Read: [Project Owner, Project Admin, Project Member, Read Teams], Update: [Project Owner, Project Admin, Edit Team]",
                Computed: true,
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Teams], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "is_permissions_editable": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Edit Team, Edit Team Permissions], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "is_team_deleteable": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Edit Team, Edit Team Permissions], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "should_have_at_least_one_member": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Edit Team, Edit Team Permissions], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "is_team_editable": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Edit Team, Edit Team Permissions], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
        },
    }
}

func (d *TeamDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *TeamDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data TeamDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "team" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read team_data, got error: %s", err))
        return
    }

    var teamDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &teamDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse team_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := teamDataResponse["data"].(map[string]interface{}); ok {
        teamDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := teamDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := teamDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := teamDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := teamDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := teamDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := teamDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := teamDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := teamDataResponse["description"].(string); ok {
        data.Description = types.StringValue(val)
    }
    if val, ok := teamDataResponse["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    }
    if val, ok := teamDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := teamDataResponse["is_permissions_editable"].(bool); ok {
        data.IsPermissionsEditable = types.BoolValue(val)
    }
    if val, ok := teamDataResponse["is_team_deleteable"].(bool); ok {
        data.IsTeamDeleteable = types.BoolValue(val)
    }
    if val, ok := teamDataResponse["should_have_at_least_one_member"].(bool); ok {
        data.ShouldHaveAtLeastOneMember = types.BoolValue(val)
    }
    if val, ok := teamDataResponse["is_team_editable"].(bool); ok {
        data.IsTeamEditable = types.BoolValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
