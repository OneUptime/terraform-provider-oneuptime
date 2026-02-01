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
var _ datasource.DataSource = &IncidentRoleDataDataSource{}

func NewIncidentRoleDataDataSource() datasource.DataSource {
    return &IncidentRoleDataDataSource{}
}

// IncidentRoleDataDataSource defines the data source implementation.
type IncidentRoleDataDataSource struct {
    client *Client
}

// IncidentRoleDataDataSourceModel describes the data source data model.
type IncidentRoleDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    Slug types.String `tfsdk:"slug"`
    Description types.String `tfsdk:"description"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
    Color types.String `tfsdk:"color"`
    RoleIcon types.String `tfsdk:"role_icon"`
    IsPrimaryRole types.Bool `tfsdk:"is_primary_role"`
    IsDeleteable types.Bool `tfsdk:"is_deleteable"`
    CanAssignMultipleUsers types.Bool `tfsdk:"can_assign_multiple_users"`
}

func (d *IncidentRoleDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_incident_role_data"
}

func (d *IncidentRoleDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "incident_role_data data source",

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
                MarkdownDescription: "Friendly globally unique name for your object. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Incident Role, Read All Project Resources], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Friendly description that will help you remember. Permissions - Create: [Project Owner, Project Admin, Project Member, Create Incident Role], Read: [Project Owner, Project Admin, Project Member, Read Incident Role, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Edit Incident Role]",
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
            "color": schema.StringAttribute{
                MarkdownDescription: "Color object",
                Computed: true,
            },
            "role_icon": schema.StringAttribute{
                MarkdownDescription: "Icon for this incident role (e.g., User, Shield, etc.). Permissions - Create: [Project Owner, Project Admin, Project Member, Create Incident Role], Read: [Project Owner, Project Admin, Project Member, Read Incident Role, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Edit Incident Role]",
                Computed: true,
            },
            "is_primary_role": schema.BoolAttribute{
                MarkdownDescription: "Is this the primary incident role? Primary roles like Incident Commander have special significance.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Incident Role, Read All Project Resources], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "is_deleteable": schema.BoolAttribute{
                MarkdownDescription: "Can this role be deleted? Primary roles cannot be deleted.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Incident Role, Read All Project Resources], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "can_assign_multiple_users": schema.BoolAttribute{
                MarkdownDescription: "Can multiple users be assigned to this role? If false, only one user can be assigned.. Permissions - Create: [Project Owner, Project Admin, Project Member, Create Incident Role], Read: [Project Owner, Project Admin, Project Member, Read Incident Role, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Edit Incident Role]",
                Computed: true,
            },
        },
    }
}

func (d *IncidentRoleDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *IncidentRoleDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data IncidentRoleDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "incident-role" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read incident_role_data, got error: %s", err))
        return
    }

    var incidentRoleDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &incidentRoleDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse incident_role_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := incidentRoleDataResponse["data"].(map[string]interface{}); ok {
        incidentRoleDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := incidentRoleDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := incidentRoleDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := incidentRoleDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := incidentRoleDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := incidentRoleDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := incidentRoleDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := incidentRoleDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := incidentRoleDataResponse["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    }
    if val, ok := incidentRoleDataResponse["description"].(string); ok {
        data.Description = types.StringValue(val)
    }
    if val, ok := incidentRoleDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := incidentRoleDataResponse["deleted_by_user_id"].(string); ok {
        data.DeletedByUserId = types.StringValue(val)
    }
    if val, ok := incidentRoleDataResponse["color"].(string); ok {
        data.Color = types.StringValue(val)
    }
    if val, ok := incidentRoleDataResponse["role_icon"].(string); ok {
        data.RoleIcon = types.StringValue(val)
    }
    if val, ok := incidentRoleDataResponse["is_primary_role"].(bool); ok {
        data.IsPrimaryRole = types.BoolValue(val)
    }
    if val, ok := incidentRoleDataResponse["is_deleteable"].(bool); ok {
        data.IsDeleteable = types.BoolValue(val)
    }
    if val, ok := incidentRoleDataResponse["can_assign_multiple_users"].(bool); ok {
        data.CanAssignMultipleUsers = types.BoolValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
