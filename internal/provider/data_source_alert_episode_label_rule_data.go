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
var _ datasource.DataSource = &AlertEpisodeLabelRuleDataDataSource{}

func NewAlertEpisodeLabelRuleDataDataSource() datasource.DataSource {
    return &AlertEpisodeLabelRuleDataDataSource{}
}

// AlertEpisodeLabelRuleDataDataSource defines the data source implementation.
type AlertEpisodeLabelRuleDataDataSource struct {
    client *Client
}

// AlertEpisodeLabelRuleDataDataSourceModel describes the data source data model.
type AlertEpisodeLabelRuleDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    Description types.String `tfsdk:"description"`
    IsEnabled types.Bool `tfsdk:"is_enabled"`
    AlertSeverities types.Set `tfsdk:"alert_severities"`
    EpisodeLabels types.Set `tfsdk:"episode_labels"`
    EpisodeTitlePattern types.String `tfsdk:"episode_title_pattern"`
    EpisodeDescriptionPattern types.String `tfsdk:"episode_description_pattern"`
    LabelsToAdd types.Set `tfsdk:"labels_to_add"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *AlertEpisodeLabelRuleDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_alert_episode_label_rule_data"
}

func (d *AlertEpisodeLabelRuleDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "alert_episode_label_rule_data data source",

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
            "description": schema.StringAttribute{
                MarkdownDescription: "Description of this alert episode label rule. Permissions - Create: [Project Owner, Project Admin, Create Alert Episode Label Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Alert Admin, Alert Member, Alert Viewer, Read Alert Episode Label Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Alert Episode Label Rule]",
                Computed: true,
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether this rule is enabled. Permissions - Create: [Project Owner, Project Admin, Create Alert Episode Label Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Alert Admin, Alert Member, Alert Viewer, Read Alert Episode Label Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Alert Episode Label Rule]",
                Computed: true,
            },
            "alert_severities": schema.SetAttribute{
                MarkdownDescription: "Only trigger for episodes with these severities. Leave empty to match any severity.. Permissions - Create: [Project Owner, Project Admin, Create Alert Episode Label Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Alert Admin, Alert Member, Alert Viewer, Read Alert Episode Label Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Alert Episode Label Rule]",
                Computed: true,
                ElementType: types.StringType,
            },
            "episode_labels": schema.SetAttribute{
                MarkdownDescription: "Only trigger for episodes that already have at least one of these labels.. Permissions - Create: [Project Owner, Project Admin, Create Alert Episode Label Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Alert Admin, Alert Member, Alert Viewer, Read Alert Episode Label Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Alert Episode Label Rule]",
                Computed: true,
                ElementType: types.StringType,
            },
            "episode_title_pattern": schema.StringAttribute{
                MarkdownDescription: "Regex (case-insensitive) matched against the episode title.. Permissions - Create: [Project Owner, Project Admin, Create Alert Episode Label Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Alert Admin, Alert Member, Alert Viewer, Read Alert Episode Label Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Alert Episode Label Rule]",
                Computed: true,
            },
            "episode_description_pattern": schema.StringAttribute{
                MarkdownDescription: "Regex (case-insensitive) matched against the episode description.. Permissions - Create: [Project Owner, Project Admin, Create Alert Episode Label Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Alert Admin, Alert Member, Alert Viewer, Read Alert Episode Label Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Alert Episode Label Rule]",
                Computed: true,
            },
            "labels_to_add": schema.SetAttribute{
                MarkdownDescription: "Labels to attach to the episode when this rule matches. Already-attached labels are not duplicated.. Permissions - Create: [Project Owner, Project Admin, Create Alert Episode Label Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Alert Admin, Alert Member, Alert Viewer, Read Alert Episode Label Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Alert Episode Label Rule]",
                Computed: true,
                ElementType: types.StringType,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
        },
    }
}

func (d *AlertEpisodeLabelRuleDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AlertEpisodeLabelRuleDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data AlertEpisodeLabelRuleDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "alert-episode-label-rule" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read alert_episode_label_rule_data, got error: %s", err))
        return
    }

    var alertEpisodeLabelRuleDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &alertEpisodeLabelRuleDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse alert_episode_label_rule_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := alertEpisodeLabelRuleDataResponse["data"].(map[string]interface{}); ok {
        alertEpisodeLabelRuleDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := alertEpisodeLabelRuleDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := alertEpisodeLabelRuleDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := alertEpisodeLabelRuleDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := alertEpisodeLabelRuleDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := alertEpisodeLabelRuleDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := alertEpisodeLabelRuleDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := alertEpisodeLabelRuleDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := alertEpisodeLabelRuleDataResponse["description"].(string); ok {
        data.Description = types.StringValue(val)
    }
    if val, ok := alertEpisodeLabelRuleDataResponse["is_enabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    }
    if val, ok := alertEpisodeLabelRuleDataResponse["alert_severities"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.AlertSeverities = setValue
    }
    if val, ok := alertEpisodeLabelRuleDataResponse["episode_labels"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.EpisodeLabels = setValue
    }
    if val, ok := alertEpisodeLabelRuleDataResponse["episode_title_pattern"].(string); ok {
        data.EpisodeTitlePattern = types.StringValue(val)
    }
    if val, ok := alertEpisodeLabelRuleDataResponse["episode_description_pattern"].(string); ok {
        data.EpisodeDescriptionPattern = types.StringValue(val)
    }
    if val, ok := alertEpisodeLabelRuleDataResponse["labels_to_add"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.LabelsToAdd = setValue
    }
    if val, ok := alertEpisodeLabelRuleDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
