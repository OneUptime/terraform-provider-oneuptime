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
var _ datasource.DataSource = &AlertEpisodeUserOwnerDataDataSource{}

func NewAlertEpisodeUserOwnerDataDataSource() datasource.DataSource {
    return &AlertEpisodeUserOwnerDataDataSource{}
}

// AlertEpisodeUserOwnerDataDataSource defines the data source implementation.
type AlertEpisodeUserOwnerDataDataSource struct {
    client *Client
}

// AlertEpisodeUserOwnerDataDataSourceModel describes the data source data model.
type AlertEpisodeUserOwnerDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    UserId types.String `tfsdk:"user_id"`
    AlertEpisodeId types.String `tfsdk:"alert_episode_id"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
    IsOwnerNotified types.Bool `tfsdk:"is_owner_notified"`
}

func (d *AlertEpisodeUserOwnerDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_alert_episode_user_owner_data"
}

func (d *AlertEpisodeUserOwnerDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "alert_episode_user_owner_data data source",

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
            "user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "alert_episode_id": schema.StringAttribute{
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
                MarkdownDescription: "Are owners notified of this resource ownership?. Permissions - Create: [Project Owner, Project Admin, Project Member, Alert Admin, Alert Member, Create Alert Episode User Owner], Read: [Project Owner, Project Admin, Project Member, Viewer, Alert Admin, Alert Member, Alert Viewer, Read Alert Episode User Owner, Read All Project Resources], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
        },
    }
}

func (d *AlertEpisodeUserOwnerDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AlertEpisodeUserOwnerDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data AlertEpisodeUserOwnerDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "alert-episode-owner-user" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read alert_episode_user_owner_data, got error: %s", err))
        return
    }

    var alertEpisodeUserOwnerDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &alertEpisodeUserOwnerDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse alert_episode_user_owner_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := alertEpisodeUserOwnerDataResponse["data"].(map[string]interface{}); ok {
        alertEpisodeUserOwnerDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := alertEpisodeUserOwnerDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := alertEpisodeUserOwnerDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := alertEpisodeUserOwnerDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := alertEpisodeUserOwnerDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := alertEpisodeUserOwnerDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := alertEpisodeUserOwnerDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := alertEpisodeUserOwnerDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := alertEpisodeUserOwnerDataResponse["user_id"].(string); ok {
        data.UserId = types.StringValue(val)
    }
    if val, ok := alertEpisodeUserOwnerDataResponse["alert_episode_id"].(string); ok {
        data.AlertEpisodeId = types.StringValue(val)
    }
    if val, ok := alertEpisodeUserOwnerDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := alertEpisodeUserOwnerDataResponse["deleted_by_user_id"].(string); ok {
        data.DeletedByUserId = types.StringValue(val)
    }
    if val, ok := alertEpisodeUserOwnerDataResponse["is_owner_notified"].(bool); ok {
        data.IsOwnerNotified = types.BoolValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
