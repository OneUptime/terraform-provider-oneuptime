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
var _ datasource.DataSource = &HostLabelRuleDataDataSource{}

func NewHostLabelRuleDataDataSource() datasource.DataSource {
    return &HostLabelRuleDataDataSource{}
}

// HostLabelRuleDataDataSource defines the data source implementation.
type HostLabelRuleDataDataSource struct {
    client *Client
}

// HostLabelRuleDataDataSourceModel describes the data source data model.
type HostLabelRuleDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    Description types.String `tfsdk:"description"`
    IsEnabled types.Bool `tfsdk:"is_enabled"`
    HostLabels types.Set `tfsdk:"host_labels"`
    HostNamePattern types.String `tfsdk:"host_name_pattern"`
    HostDescriptionPattern types.String `tfsdk:"host_description_pattern"`
    LabelsToAdd types.Set `tfsdk:"labels_to_add"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *HostLabelRuleDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_host_label_rule_data"
}

func (d *HostLabelRuleDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "host_label_rule_data data source",

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
                MarkdownDescription: "Description of this host label rule. Permissions - Create: [Project Owner, Project Admin, Create Host Label Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Host Label Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Host Label Rule]",
                Computed: true,
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether this rule is enabled. Permissions - Create: [Project Owner, Project Admin, Create Host Label Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Host Label Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Host Label Rule]",
                Computed: true,
            },
            "host_labels": schema.SetAttribute{
                MarkdownDescription: "Only trigger for hosts that already have at least one of these labels. Leave empty to match regardless of labels.. Permissions - Create: [Project Owner, Project Admin, Create Host Label Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Host Label Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Host Label Rule]",
                Computed: true,
                ElementType: types.StringType,
            },
            "host_name_pattern": schema.StringAttribute{
                MarkdownDescription: "Regex (case-insensitive) matched against the host name. Leave empty to match any name.. Permissions - Create: [Project Owner, Project Admin, Create Host Label Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Host Label Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Host Label Rule]",
                Computed: true,
            },
            "host_description_pattern": schema.StringAttribute{
                MarkdownDescription: "Regex (case-insensitive) matched against the host description. Leave empty to match any description.. Permissions - Create: [Project Owner, Project Admin, Create Host Label Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Host Label Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Host Label Rule]",
                Computed: true,
            },
            "labels_to_add": schema.SetAttribute{
                MarkdownDescription: "Labels to attach to the host when this rule matches. Already-attached labels are not duplicated.. Permissions - Create: [Project Owner, Project Admin, Create Host Label Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Host Label Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Host Label Rule]",
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

func (d *HostLabelRuleDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *HostLabelRuleDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data HostLabelRuleDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "host-label-rule" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read host_label_rule_data, got error: %s", err))
        return
    }

    var hostLabelRuleDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &hostLabelRuleDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse host_label_rule_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := hostLabelRuleDataResponse["data"].(map[string]interface{}); ok {
        hostLabelRuleDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := hostLabelRuleDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := hostLabelRuleDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := hostLabelRuleDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := hostLabelRuleDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := hostLabelRuleDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := hostLabelRuleDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := hostLabelRuleDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := hostLabelRuleDataResponse["description"].(string); ok {
        data.Description = types.StringValue(val)
    }
    if val, ok := hostLabelRuleDataResponse["is_enabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    }
    if val, ok := hostLabelRuleDataResponse["host_labels"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.HostLabels = setValue
    }
    if val, ok := hostLabelRuleDataResponse["host_name_pattern"].(string); ok {
        data.HostNamePattern = types.StringValue(val)
    }
    if val, ok := hostLabelRuleDataResponse["host_description_pattern"].(string); ok {
        data.HostDescriptionPattern = types.StringValue(val)
    }
    if val, ok := hostLabelRuleDataResponse["labels_to_add"].([]interface{}); ok {
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
    if val, ok := hostLabelRuleDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
