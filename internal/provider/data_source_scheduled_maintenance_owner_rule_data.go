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
var _ datasource.DataSource = &ScheduledMaintenanceOwnerRuleDataDataSource{}

func NewScheduledMaintenanceOwnerRuleDataDataSource() datasource.DataSource {
    return &ScheduledMaintenanceOwnerRuleDataDataSource{}
}

// ScheduledMaintenanceOwnerRuleDataDataSource defines the data source implementation.
type ScheduledMaintenanceOwnerRuleDataDataSource struct {
    client *Client
}

// ScheduledMaintenanceOwnerRuleDataDataSourceModel describes the data source data model.
type ScheduledMaintenanceOwnerRuleDataDataSourceModel struct {
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
    Monitors types.Set `tfsdk:"monitors"`
    ScheduledMaintenanceLabels types.Set `tfsdk:"scheduled_maintenance_labels"`
    MonitorLabels types.Set `tfsdk:"monitor_labels"`
    TitlePattern types.String `tfsdk:"title_pattern"`
    DescriptionPattern types.String `tfsdk:"description_pattern"`
    MonitorNamePattern types.String `tfsdk:"monitor_name_pattern"`
    MonitorDescriptionPattern types.String `tfsdk:"monitor_description_pattern"`
    OwnerUsers types.Set `tfsdk:"owner_users"`
    OwnerTeams types.Set `tfsdk:"owner_teams"`
    InheritOwnersFromMonitors types.Bool `tfsdk:"inherit_owners_from_monitors"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *ScheduledMaintenanceOwnerRuleDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_scheduled_maintenance_owner_rule_data"
}

func (d *ScheduledMaintenanceOwnerRuleDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "scheduled_maintenance_owner_rule_data data source",

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
                MarkdownDescription: "Description of this scheduled maintenance owner rule. Permissions - Create: [Project Owner, Project Admin, Create Scheduled Maintenance Owner Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Scheduled Maintenance Admin, Scheduled Maintenance Member, Scheduled Maintenance Viewer, Read Scheduled Maintenance Owner Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Scheduled Maintenance Owner Rule]",
                Computed: true,
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether this rule is enabled. Permissions - Create: [Project Owner, Project Admin, Create Scheduled Maintenance Owner Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Scheduled Maintenance Admin, Scheduled Maintenance Member, Scheduled Maintenance Viewer, Read Scheduled Maintenance Owner Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Scheduled Maintenance Owner Rule]",
                Computed: true,
            },
            "notify_owners": schema.BoolAttribute{
                MarkdownDescription: "Send notifications to owner users and teams when they are added by this rule. Permissions - Create: [Project Owner, Project Admin, Create Scheduled Maintenance Owner Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Scheduled Maintenance Admin, Scheduled Maintenance Member, Scheduled Maintenance Viewer, Read Scheduled Maintenance Owner Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Scheduled Maintenance Owner Rule]",
                Computed: true,
            },
            "monitors": schema.SetAttribute{
                MarkdownDescription: "Only trigger for scheduled maintenance events on these monitors. Leave empty to match events on any monitor.. Permissions - Create: [Project Owner, Project Admin, Create Scheduled Maintenance Owner Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Scheduled Maintenance Admin, Scheduled Maintenance Member, Scheduled Maintenance Viewer, Read Scheduled Maintenance Owner Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Scheduled Maintenance Owner Rule]",
                Computed: true,
                ElementType: types.StringType,
            },
            "scheduled_maintenance_labels": schema.SetAttribute{
                MarkdownDescription: "Only trigger for events that have at least one of these labels. Leave empty to match regardless of event labels.. Permissions - Create: [Project Owner, Project Admin, Create Scheduled Maintenance Owner Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Scheduled Maintenance Admin, Scheduled Maintenance Member, Scheduled Maintenance Viewer, Read Scheduled Maintenance Owner Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Scheduled Maintenance Owner Rule]",
                Computed: true,
                ElementType: types.StringType,
            },
            "monitor_labels": schema.SetAttribute{
                MarkdownDescription: "Only trigger for events on monitors that have at least one of these labels. Leave empty to match regardless of monitor labels.. Permissions - Create: [Project Owner, Project Admin, Create Scheduled Maintenance Owner Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Scheduled Maintenance Admin, Scheduled Maintenance Member, Scheduled Maintenance Viewer, Read Scheduled Maintenance Owner Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Scheduled Maintenance Owner Rule]",
                Computed: true,
                ElementType: types.StringType,
            },
            "title_pattern": schema.StringAttribute{
                MarkdownDescription: "Regex (case-insensitive) matched against the scheduled maintenance event title. Leave empty to match any title.. Permissions - Create: [Project Owner, Project Admin, Create Scheduled Maintenance Owner Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Scheduled Maintenance Admin, Scheduled Maintenance Member, Scheduled Maintenance Viewer, Read Scheduled Maintenance Owner Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Scheduled Maintenance Owner Rule]",
                Computed: true,
            },
            "description_pattern": schema.StringAttribute{
                MarkdownDescription: "Regex (case-insensitive) matched against the scheduled maintenance event description. Leave empty to match any description.. Permissions - Create: [Project Owner, Project Admin, Create Scheduled Maintenance Owner Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Scheduled Maintenance Admin, Scheduled Maintenance Member, Scheduled Maintenance Viewer, Read Scheduled Maintenance Owner Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Scheduled Maintenance Owner Rule]",
                Computed: true,
            },
            "monitor_name_pattern": schema.StringAttribute{
                MarkdownDescription: "Regex (case-insensitive) matched against any of the event's monitor names. Leave empty to match any monitor.. Permissions - Create: [Project Owner, Project Admin, Create Scheduled Maintenance Owner Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Scheduled Maintenance Admin, Scheduled Maintenance Member, Scheduled Maintenance Viewer, Read Scheduled Maintenance Owner Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Scheduled Maintenance Owner Rule]",
                Computed: true,
            },
            "monitor_description_pattern": schema.StringAttribute{
                MarkdownDescription: "Regex (case-insensitive) matched against any of the event's monitor descriptions. Leave empty to match any description.. Permissions - Create: [Project Owner, Project Admin, Create Scheduled Maintenance Owner Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Scheduled Maintenance Admin, Scheduled Maintenance Member, Scheduled Maintenance Viewer, Read Scheduled Maintenance Owner Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Scheduled Maintenance Owner Rule]",
                Computed: true,
            },
            "owner_users": schema.SetAttribute{
                MarkdownDescription: "Users to add as owners on the scheduled maintenance event when this rule matches.. Permissions - Create: [Project Owner, Project Admin, Create Scheduled Maintenance Owner Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Scheduled Maintenance Admin, Scheduled Maintenance Member, Scheduled Maintenance Viewer, Read Scheduled Maintenance Owner Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Scheduled Maintenance Owner Rule]",
                Computed: true,
                ElementType: types.StringType,
            },
            "owner_teams": schema.SetAttribute{
                MarkdownDescription: "Teams to add as owners on the scheduled maintenance event when this rule matches.. Permissions - Create: [Project Owner, Project Admin, Create Scheduled Maintenance Owner Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Scheduled Maintenance Admin, Scheduled Maintenance Member, Scheduled Maintenance Viewer, Read Scheduled Maintenance Owner Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Scheduled Maintenance Owner Rule]",
                Computed: true,
                ElementType: types.StringType,
            },
            "inherit_owners_from_monitors": schema.BoolAttribute{
                MarkdownDescription: "When this rule matches, also assign every owner of the event's monitors to the event.. Permissions - Create: [Project Owner, Project Admin, Create Scheduled Maintenance Owner Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Scheduled Maintenance Admin, Scheduled Maintenance Member, Scheduled Maintenance Viewer, Read Scheduled Maintenance Owner Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Scheduled Maintenance Owner Rule]",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
        },
    }
}

func (d *ScheduledMaintenanceOwnerRuleDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ScheduledMaintenanceOwnerRuleDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data ScheduledMaintenanceOwnerRuleDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "scheduled-maintenance-owner-rule" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read scheduled_maintenance_owner_rule_data, got error: %s", err))
        return
    }

    var scheduledMaintenanceOwnerRuleDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &scheduledMaintenanceOwnerRuleDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse scheduled_maintenance_owner_rule_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := scheduledMaintenanceOwnerRuleDataResponse["data"].(map[string]interface{}); ok {
        scheduledMaintenanceOwnerRuleDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := scheduledMaintenanceOwnerRuleDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceOwnerRuleDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceOwnerRuleDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceOwnerRuleDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceOwnerRuleDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceOwnerRuleDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := scheduledMaintenanceOwnerRuleDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceOwnerRuleDataResponse["description"].(string); ok {
        data.Description = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceOwnerRuleDataResponse["is_enabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    }
    if val, ok := scheduledMaintenanceOwnerRuleDataResponse["notify_owners"].(bool); ok {
        data.NotifyOwners = types.BoolValue(val)
    }
    if val, ok := scheduledMaintenanceOwnerRuleDataResponse["monitors"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.Monitors = setValue
    }
    if val, ok := scheduledMaintenanceOwnerRuleDataResponse["scheduled_maintenance_labels"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.ScheduledMaintenanceLabels = setValue
    }
    if val, ok := scheduledMaintenanceOwnerRuleDataResponse["monitor_labels"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.MonitorLabels = setValue
    }
    if val, ok := scheduledMaintenanceOwnerRuleDataResponse["title_pattern"].(string); ok {
        data.TitlePattern = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceOwnerRuleDataResponse["description_pattern"].(string); ok {
        data.DescriptionPattern = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceOwnerRuleDataResponse["monitor_name_pattern"].(string); ok {
        data.MonitorNamePattern = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceOwnerRuleDataResponse["monitor_description_pattern"].(string); ok {
        data.MonitorDescriptionPattern = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceOwnerRuleDataResponse["owner_users"].([]interface{}); ok {
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
    if val, ok := scheduledMaintenanceOwnerRuleDataResponse["owner_teams"].([]interface{}); ok {
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
    if val, ok := scheduledMaintenanceOwnerRuleDataResponse["inherit_owners_from_monitors"].(bool); ok {
        data.InheritOwnersFromMonitors = types.BoolValue(val)
    }
    if val, ok := scheduledMaintenanceOwnerRuleDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
