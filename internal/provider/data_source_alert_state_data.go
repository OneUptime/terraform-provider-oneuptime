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
var _ datasource.DataSource = &AlertStateDataDataSource{}

func NewAlertStateDataDataSource() datasource.DataSource {
    return &AlertStateDataDataSource{}
}

// AlertStateDataDataSource defines the data source implementation.
type AlertStateDataDataSource struct {
    client *Client
}

// AlertStateDataDataSourceModel describes the data source data model.
type AlertStateDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    Description types.String `tfsdk:"description"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
    Color types.String `tfsdk:"color"`
    IsCreatedState types.Bool `tfsdk:"is_created_state"`
    IsAcknowledgedState types.Bool `tfsdk:"is_acknowledged_state"`
    IsResolvedState types.Bool `tfsdk:"is_resolved_state"`
    Order types.Number `tfsdk:"order"`
}

func (d *AlertStateDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_alert_state_data"
}

func (d *AlertStateDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "alert_state_data data source",

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
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Alert State], Read: [Project Owner, Project Admin, Project Member, Read Alert State], Update: [Project Owner, Project Admin, Project Member, Edit Alert State]",
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
            "is_created_state": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Alert State], Read: [Project Owner, Project Admin, Project Member, Read Alert State], Update: [Project Owner, Project Admin, Project Member, Edit Alert State]",
                Computed: true,
            },
            "is_acknowledged_state": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Alert State], Read: [Project Owner, Project Admin, Project Member, Read Alert State], Update: [Project Owner, Project Admin, Project Member, Edit Alert State]",
                Computed: true,
            },
            "is_resolved_state": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Alert State], Read: [Project Owner, Project Admin, Project Member, Read Alert State], Update: [Project Owner, Project Admin, Project Member, Edit Alert State]",
                Computed: true,
            },
            "order": schema.NumberAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Alert State], Read: [Project Owner, Project Admin, Project Member, Read Alert State], Update: [Project Owner, Project Admin, Project Member, Edit Alert State]",
                Computed: true,
            },
        },
    }
}

func (d *AlertStateDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AlertStateDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data AlertStateDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "alert-state" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read alert_state_data, got error: %s", err))
        return
    }

    var alertStateDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &alertStateDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse alert_state_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := alertStateDataResponse["data"].(map[string]interface{}); ok {
        alertStateDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := alertStateDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := alertStateDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := alertStateDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := alertStateDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := alertStateDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := alertStateDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := alertStateDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := alertStateDataResponse["description"].(string); ok {
        data.Description = types.StringValue(val)
    }
    if val, ok := alertStateDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := alertStateDataResponse["deleted_by_user_id"].(string); ok {
        data.DeletedByUserId = types.StringValue(val)
    }
    if val, ok := alertStateDataResponse["color"].(string); ok {
        data.Color = types.StringValue(val)
    }
    if val, ok := alertStateDataResponse["is_created_state"].(bool); ok {
        data.IsCreatedState = types.BoolValue(val)
    }
    if val, ok := alertStateDataResponse["is_acknowledged_state"].(bool); ok {
        data.IsAcknowledgedState = types.BoolValue(val)
    }
    if val, ok := alertStateDataResponse["is_resolved_state"].(bool); ok {
        data.IsResolvedState = types.BoolValue(val)
    }
    if val, ok := alertStateDataResponse["order"].(float64); ok {
        data.Order = types.NumberValue(big.NewFloat(val))
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
