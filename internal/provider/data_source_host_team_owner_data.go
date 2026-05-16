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
var _ datasource.DataSource = &HostTeamOwnerDataDataSource{}

func NewHostTeamOwnerDataDataSource() datasource.DataSource {
    return &HostTeamOwnerDataDataSource{}
}

// HostTeamOwnerDataDataSource defines the data source implementation.
type HostTeamOwnerDataDataSource struct {
    client *Client
}

// HostTeamOwnerDataDataSourceModel describes the data source data model.
type HostTeamOwnerDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    TeamId types.String `tfsdk:"team_id"`
    HostId types.String `tfsdk:"host_id"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
    IsOwnerNotified types.Bool `tfsdk:"is_owner_notified"`
}

func (d *HostTeamOwnerDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_host_team_owner_data"
}

func (d *HostTeamOwnerDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "host_team_owner_data data source",

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
            "team_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "host_id": schema.StringAttribute{
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
                MarkdownDescription: "Are owners notified of this resource ownership?. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Host Team Owner, Read All Project Resources], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
        },
    }
}

func (d *HostTeamOwnerDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *HostTeamOwnerDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data HostTeamOwnerDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "host-owner-team" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read host_team_owner_data, got error: %s", err))
        return
    }

    var hostTeamOwnerDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &hostTeamOwnerDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse host_team_owner_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := hostTeamOwnerDataResponse["data"].(map[string]interface{}); ok {
        hostTeamOwnerDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := hostTeamOwnerDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := hostTeamOwnerDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := hostTeamOwnerDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := hostTeamOwnerDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := hostTeamOwnerDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := hostTeamOwnerDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := hostTeamOwnerDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := hostTeamOwnerDataResponse["team_id"].(string); ok {
        data.TeamId = types.StringValue(val)
    }
    if val, ok := hostTeamOwnerDataResponse["host_id"].(string); ok {
        data.HostId = types.StringValue(val)
    }
    if val, ok := hostTeamOwnerDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := hostTeamOwnerDataResponse["deleted_by_user_id"].(string); ok {
        data.DeletedByUserId = types.StringValue(val)
    }
    if val, ok := hostTeamOwnerDataResponse["is_owner_notified"].(bool); ok {
        data.IsOwnerNotified = types.BoolValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
