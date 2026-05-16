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
var _ datasource.DataSource = &AlertEpisodeMemberDataDataSource{}

func NewAlertEpisodeMemberDataDataSource() datasource.DataSource {
    return &AlertEpisodeMemberDataDataSource{}
}

// AlertEpisodeMemberDataDataSource defines the data source implementation.
type AlertEpisodeMemberDataDataSource struct {
    client *Client
}

// AlertEpisodeMemberDataDataSourceModel describes the data source data model.
type AlertEpisodeMemberDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    AlertEpisodeId types.String `tfsdk:"alert_episode_id"`
    AlertId types.String `tfsdk:"alert_id"`
    AddedAt types.String `tfsdk:"added_at"`
    AddedBy types.String `tfsdk:"added_by"`
    AddedByUserId types.String `tfsdk:"added_by_user_id"`
    MatchedRuleId types.String `tfsdk:"matched_rule_id"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
}

func (d *AlertEpisodeMemberDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_alert_episode_member_data"
}

func (d *AlertEpisodeMemberDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "alert_episode_member_data data source",

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
            "alert_episode_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "alert_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "added_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "added_by": schema.StringAttribute{
                MarkdownDescription: "How this alert was added to the episode (rule, manual, or api). Permissions - Create: [Project Owner, Project Admin, Project Member, Alert Admin, Alert Member, Create Alert Episode Member], Read: [Project Owner, Project Admin, Project Member, Viewer, Alert Admin, Alert Member, Alert Viewer, Read Alert Episode Member, Read All Project Resources], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "added_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "matched_rule_id": schema.StringAttribute{
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
        },
    }
}

func (d *AlertEpisodeMemberDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AlertEpisodeMemberDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data AlertEpisodeMemberDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "alert-episode-member" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read alert_episode_member_data, got error: %s", err))
        return
    }

    var alertEpisodeMemberDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &alertEpisodeMemberDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse alert_episode_member_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := alertEpisodeMemberDataResponse["data"].(map[string]interface{}); ok {
        alertEpisodeMemberDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := alertEpisodeMemberDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := alertEpisodeMemberDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := alertEpisodeMemberDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := alertEpisodeMemberDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := alertEpisodeMemberDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := alertEpisodeMemberDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := alertEpisodeMemberDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := alertEpisodeMemberDataResponse["alert_episode_id"].(string); ok {
        data.AlertEpisodeId = types.StringValue(val)
    }
    if val, ok := alertEpisodeMemberDataResponse["alert_id"].(string); ok {
        data.AlertId = types.StringValue(val)
    }
    if val, ok := alertEpisodeMemberDataResponse["added_at"].(string); ok {
        data.AddedAt = types.StringValue(val)
    }
    if val, ok := alertEpisodeMemberDataResponse["added_by"].(string); ok {
        data.AddedBy = types.StringValue(val)
    }
    if val, ok := alertEpisodeMemberDataResponse["added_by_user_id"].(string); ok {
        data.AddedByUserId = types.StringValue(val)
    }
    if val, ok := alertEpisodeMemberDataResponse["matched_rule_id"].(string); ok {
        data.MatchedRuleId = types.StringValue(val)
    }
    if val, ok := alertEpisodeMemberDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := alertEpisodeMemberDataResponse["deleted_by_user_id"].(string); ok {
        data.DeletedByUserId = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
