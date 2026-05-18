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
var _ datasource.DataSource = &ServiceOwnerRuleDataDataSource{}

func NewServiceOwnerRuleDataDataSource() datasource.DataSource {
    return &ServiceOwnerRuleDataDataSource{}
}

// ServiceOwnerRuleDataDataSource defines the data source implementation.
type ServiceOwnerRuleDataDataSource struct {
    client *Client
}

// ServiceOwnerRuleDataDataSourceModel describes the data source data model.
type ServiceOwnerRuleDataDataSourceModel struct {
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
    ServiceLabels types.Set `tfsdk:"service_labels"`
    ServiceNamePattern types.String `tfsdk:"service_name_pattern"`
    ServiceDescriptionPattern types.String `tfsdk:"service_description_pattern"`
    OwnerUsers types.Set `tfsdk:"owner_users"`
    OwnerTeams types.Set `tfsdk:"owner_teams"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *ServiceOwnerRuleDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_service_owner_rule_data"
}

func (d *ServiceOwnerRuleDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "service_owner_rule_data data source",

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
                MarkdownDescription: "Description of this service owner rule. Permissions - Create: [Project Owner, Project Admin, Create Service Owner Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Service Owner Rule], Update: [Project Owner, Project Admin, Edit Service Owner Rule]",
                Computed: true,
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether this rule is enabled. Permissions - Create: [Project Owner, Project Admin, Create Service Owner Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Service Owner Rule], Update: [Project Owner, Project Admin, Edit Service Owner Rule]",
                Computed: true,
            },
            "notify_owners": schema.BoolAttribute{
                MarkdownDescription: "Send notifications to owner users and teams when they are added by this rule. Permissions - Create: [Project Owner, Project Admin, Create Service Owner Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Service Owner Rule], Update: [Project Owner, Project Admin, Edit Service Owner Rule]",
                Computed: true,
            },
            "service_labels": schema.SetAttribute{
                MarkdownDescription: "Only trigger for services that have at least one of these labels. Leave empty to match regardless of labels.. Permissions - Create: [Project Owner, Project Admin, Create Service Owner Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Service Owner Rule], Update: [Project Owner, Project Admin, Edit Service Owner Rule]",
                Computed: true,
                ElementType: types.StringType,
            },
            "service_name_pattern": schema.StringAttribute{
                MarkdownDescription: "Regex (case-insensitive) matched against the service name. Leave empty to match any name.. Permissions - Create: [Project Owner, Project Admin, Create Service Owner Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Service Owner Rule], Update: [Project Owner, Project Admin, Edit Service Owner Rule]",
                Computed: true,
            },
            "service_description_pattern": schema.StringAttribute{
                MarkdownDescription: "Regex (case-insensitive) matched against the service description. Leave empty to match any description.. Permissions - Create: [Project Owner, Project Admin, Create Service Owner Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Service Owner Rule], Update: [Project Owner, Project Admin, Edit Service Owner Rule]",
                Computed: true,
            },
            "owner_users": schema.SetAttribute{
                MarkdownDescription: "Users to add as owners on the service when this rule matches.. Permissions - Create: [Project Owner, Project Admin, Create Service Owner Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Service Owner Rule], Update: [Project Owner, Project Admin, Edit Service Owner Rule]",
                Computed: true,
                ElementType: types.StringType,
            },
            "owner_teams": schema.SetAttribute{
                MarkdownDescription: "Teams to add as owners on the service when this rule matches.. Permissions - Create: [Project Owner, Project Admin, Create Service Owner Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Service Owner Rule], Update: [Project Owner, Project Admin, Edit Service Owner Rule]",
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

func (d *ServiceOwnerRuleDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ServiceOwnerRuleDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data ServiceOwnerRuleDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "service-owner-rule" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read service_owner_rule_data, got error: %s", err))
        return
    }

    var serviceOwnerRuleDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &serviceOwnerRuleDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse service_owner_rule_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := serviceOwnerRuleDataResponse["data"].(map[string]interface{}); ok {
        serviceOwnerRuleDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := serviceOwnerRuleDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := serviceOwnerRuleDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := serviceOwnerRuleDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := serviceOwnerRuleDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := serviceOwnerRuleDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := serviceOwnerRuleDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := serviceOwnerRuleDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := serviceOwnerRuleDataResponse["description"].(string); ok {
        data.Description = types.StringValue(val)
    }
    if val, ok := serviceOwnerRuleDataResponse["is_enabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    }
    if val, ok := serviceOwnerRuleDataResponse["notify_owners"].(bool); ok {
        data.NotifyOwners = types.BoolValue(val)
    }
    if val, ok := serviceOwnerRuleDataResponse["service_labels"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.ServiceLabels = setValue
    }
    if val, ok := serviceOwnerRuleDataResponse["service_name_pattern"].(string); ok {
        data.ServiceNamePattern = types.StringValue(val)
    }
    if val, ok := serviceOwnerRuleDataResponse["service_description_pattern"].(string); ok {
        data.ServiceDescriptionPattern = types.StringValue(val)
    }
    if val, ok := serviceOwnerRuleDataResponse["owner_users"].([]interface{}); ok {
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
    if val, ok := serviceOwnerRuleDataResponse["owner_teams"].([]interface{}); ok {
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
    if val, ok := serviceOwnerRuleDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
