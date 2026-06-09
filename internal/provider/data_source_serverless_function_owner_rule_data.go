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
var _ datasource.DataSource = &ServerlessFunctionOwnerRuleDataDataSource{}

func NewServerlessFunctionOwnerRuleDataDataSource() datasource.DataSource {
    return &ServerlessFunctionOwnerRuleDataDataSource{}
}

// ServerlessFunctionOwnerRuleDataDataSource defines the data source implementation.
type ServerlessFunctionOwnerRuleDataDataSource struct {
    client *Client
}

// ServerlessFunctionOwnerRuleDataDataSourceModel describes the data source data model.
type ServerlessFunctionOwnerRuleDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    Description types.String `tfsdk:"description"`
    IsEnabled types.Bool `tfsdk:"is_enabled"`
    NotifyOwners types.Bool `tfsdk:"notify_owners"`
    MatchLabels types.Set `tfsdk:"match_labels"`
    NameRegexPattern types.String `tfsdk:"name_regex_pattern"`
    DescriptionRegexPattern types.String `tfsdk:"description_regex_pattern"`
    OwnerUsers types.Set `tfsdk:"owner_users"`
    OwnerTeams types.Set `tfsdk:"owner_teams"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *ServerlessFunctionOwnerRuleDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_serverless_function_owner_rule_data"
}

func (d *ServerlessFunctionOwnerRuleDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "serverless_function_owner_rule_data data source",

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
                MarkdownDescription: "Description of this rule. Permissions - Create: [Project Owner, Project Admin, Create Serverless Function Owner Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Serverless Function Owner Rule], Update: [Project Owner, Project Admin, Edit Serverless Function Owner Rule]",
                Computed: true,
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether this rule is enabled. Permissions - Create: [Project Owner, Project Admin, Create Serverless Function Owner Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Serverless Function Owner Rule], Update: [Project Owner, Project Admin, Edit Serverless Function Owner Rule]",
                Computed: true,
            },
            "notify_owners": schema.BoolAttribute{
                MarkdownDescription: "Send notifications to owner users and teams when they are added by this rule.. Permissions - Create: [Project Owner, Project Admin, Create Serverless Function Owner Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Serverless Function Owner Rule], Update: [Project Owner, Project Admin, Edit Serverless Function Owner Rule]",
                Computed: true,
            },
            "match_labels": schema.SetAttribute{
                MarkdownDescription: "Only trigger for functions that have at least one of these labels. Leave empty to match regardless of labels.. Permissions - Create: [Project Owner, Project Admin, Create Serverless Function Owner Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Serverless Function Owner Rule], Update: [Project Owner, Project Admin, Edit Serverless Function Owner Rule]",
                Computed: true,
                ElementType: types.StringType,
            },
            "name_regex_pattern": schema.StringAttribute{
                MarkdownDescription: "Regex (case-insensitive) matched against the function name. Leave empty to match any name.. Permissions - Create: [Project Owner, Project Admin, Create Serverless Function Owner Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Serverless Function Owner Rule], Update: [Project Owner, Project Admin, Edit Serverless Function Owner Rule]",
                Computed: true,
            },
            "description_regex_pattern": schema.StringAttribute{
                MarkdownDescription: "Regex (case-insensitive) matched against the function description. Leave empty to match any description.. Permissions - Create: [Project Owner, Project Admin, Create Serverless Function Owner Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Serverless Function Owner Rule], Update: [Project Owner, Project Admin, Edit Serverless Function Owner Rule]",
                Computed: true,
            },
            "owner_users": schema.SetAttribute{
                MarkdownDescription: "Users to add as owners on the function when this rule matches.. Permissions - Create: [Project Owner, Project Admin, Create Serverless Function Owner Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Serverless Function Owner Rule], Update: [Project Owner, Project Admin, Edit Serverless Function Owner Rule]",
                Computed: true,
                ElementType: types.StringType,
            },
            "owner_teams": schema.SetAttribute{
                MarkdownDescription: "Teams to add as owners on the function when this rule matches.. Permissions - Create: [Project Owner, Project Admin, Create Serverless Function Owner Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Serverless Function Owner Rule], Update: [Project Owner, Project Admin, Edit Serverless Function Owner Rule]",
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

func (d *ServerlessFunctionOwnerRuleDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ServerlessFunctionOwnerRuleDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data ServerlessFunctionOwnerRuleDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "serverless-function-owner-rule" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read serverless_function_owner_rule_data, got error: %s", err))
        return
    }

    var serverlessFunctionOwnerRuleDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &serverlessFunctionOwnerRuleDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse serverless_function_owner_rule_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := serverlessFunctionOwnerRuleDataResponse["data"].(map[string]interface{}); ok {
        serverlessFunctionOwnerRuleDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := serverlessFunctionOwnerRuleDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := serverlessFunctionOwnerRuleDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := serverlessFunctionOwnerRuleDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := serverlessFunctionOwnerRuleDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := serverlessFunctionOwnerRuleDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := serverlessFunctionOwnerRuleDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := serverlessFunctionOwnerRuleDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := serverlessFunctionOwnerRuleDataResponse["description"].(string); ok {
        data.Description = types.StringValue(val)
    }
    if val, ok := serverlessFunctionOwnerRuleDataResponse["is_enabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    }
    if val, ok := serverlessFunctionOwnerRuleDataResponse["notify_owners"].(bool); ok {
        data.NotifyOwners = types.BoolValue(val)
    }
    if val, ok := serverlessFunctionOwnerRuleDataResponse["match_labels"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.MatchLabels = setValue
    }
    if val, ok := serverlessFunctionOwnerRuleDataResponse["name_regex_pattern"].(string); ok {
        data.NameRegexPattern = types.StringValue(val)
    }
    if val, ok := serverlessFunctionOwnerRuleDataResponse["description_regex_pattern"].(string); ok {
        data.DescriptionRegexPattern = types.StringValue(val)
    }
    if val, ok := serverlessFunctionOwnerRuleDataResponse["owner_users"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.OwnerUsers = setValue
    }
    if val, ok := serverlessFunctionOwnerRuleDataResponse["owner_teams"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.OwnerTeams = setValue
    }
    if val, ok := serverlessFunctionOwnerRuleDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
