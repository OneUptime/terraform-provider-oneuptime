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
var _ datasource.DataSource = &RunbookRuleDataDataSource{}

func NewRunbookRuleDataDataSource() datasource.DataSource {
    return &RunbookRuleDataDataSource{}
}

// RunbookRuleDataDataSource defines the data source implementation.
type RunbookRuleDataDataSource struct {
    client *Client
}

// RunbookRuleDataDataSourceModel describes the data source data model.
type RunbookRuleDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    Description types.String `tfsdk:"description"`
    IsEnabled types.Bool `tfsdk:"is_enabled"`
    TriggerEntityType types.String `tfsdk:"trigger_entity_type"`
    TitlePattern types.String `tfsdk:"title_pattern"`
    DescriptionPattern types.String `tfsdk:"description_pattern"`
    Runbooks types.Set `tfsdk:"runbooks"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *RunbookRuleDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_runbook_rule_data"
}

func (d *RunbookRuleDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "runbook_rule_data data source",

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
                MarkdownDescription: "Description of this runbook rule.. Permissions - Create: [Project Owner, Project Admin, Create Runbook Rule, Runbook Admin], Read: [Project Owner, Project Admin, Project Member, Viewer, Runbook Admin, Runbook Member, Runbook Viewer, Read Runbook Rule], Update: [Project Owner, Project Admin, Edit Runbook Rule, Runbook Admin]",
                Computed: true,
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether this rule is enabled.. Permissions - Create: [Project Owner, Project Admin, Create Runbook Rule, Runbook Admin], Read: [Project Owner, Project Admin, Project Member, Viewer, Runbook Admin, Runbook Member, Runbook Viewer, Read Runbook Rule], Update: [Project Owner, Project Admin, Edit Runbook Rule, Runbook Admin]",
                Computed: true,
            },
            "trigger_entity_type": schema.StringAttribute{
                MarkdownDescription: "Entity type that triggers this rule on creation: Incident, Alert, or ScheduledMaintenance.. Permissions - Create: [Project Owner, Project Admin, Create Runbook Rule, Runbook Admin], Read: [Project Owner, Project Admin, Project Member, Viewer, Runbook Admin, Runbook Member, Runbook Viewer, Read Runbook Rule], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "title_pattern": schema.StringAttribute{
                MarkdownDescription: "Case-insensitive regex matched against the entity's title. Leave empty to match any title.. Permissions - Create: [Project Owner, Project Admin, Create Runbook Rule, Runbook Admin], Read: [Project Owner, Project Admin, Project Member, Viewer, Runbook Admin, Runbook Member, Runbook Viewer, Read Runbook Rule], Update: [Project Owner, Project Admin, Edit Runbook Rule, Runbook Admin]",
                Computed: true,
            },
            "description_pattern": schema.StringAttribute{
                MarkdownDescription: "Case-insensitive regex matched against the entity's description. Leave empty to match any description.. Permissions - Create: [Project Owner, Project Admin, Create Runbook Rule, Runbook Admin], Read: [Project Owner, Project Admin, Project Member, Viewer, Runbook Admin, Runbook Member, Runbook Viewer, Read Runbook Rule], Update: [Project Owner, Project Admin, Edit Runbook Rule, Runbook Admin]",
                Computed: true,
            },
            "runbooks": schema.SetAttribute{
                MarkdownDescription: "Runbooks to start when this rule matches. Each runbook produces its own execution.. Permissions - Create: [Project Owner, Project Admin, Create Runbook Rule, Runbook Admin], Read: [Project Owner, Project Admin, Project Member, Viewer, Runbook Admin, Runbook Member, Runbook Viewer, Read Runbook Rule], Update: [Project Owner, Project Admin, Edit Runbook Rule, Runbook Admin]",
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

func (d *RunbookRuleDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *RunbookRuleDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data RunbookRuleDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "runbook-rule" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read runbook_rule_data, got error: %s", err))
        return
    }

    var runbookRuleDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &runbookRuleDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse runbook_rule_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := runbookRuleDataResponse["data"].(map[string]interface{}); ok {
        runbookRuleDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := runbookRuleDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := runbookRuleDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := runbookRuleDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := runbookRuleDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := runbookRuleDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := runbookRuleDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := runbookRuleDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := runbookRuleDataResponse["description"].(string); ok {
        data.Description = types.StringValue(val)
    }
    if val, ok := runbookRuleDataResponse["is_enabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    }
    if val, ok := runbookRuleDataResponse["trigger_entity_type"].(string); ok {
        data.TriggerEntityType = types.StringValue(val)
    }
    if val, ok := runbookRuleDataResponse["title_pattern"].(string); ok {
        data.TitlePattern = types.StringValue(val)
    }
    if val, ok := runbookRuleDataResponse["description_pattern"].(string); ok {
        data.DescriptionPattern = types.StringValue(val)
    }
    if val, ok := runbookRuleDataResponse["runbooks"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.Runbooks = setValue
    }
    if val, ok := runbookRuleDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
