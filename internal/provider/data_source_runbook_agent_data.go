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
var _ datasource.DataSource = &RunbookAgentDataDataSource{}

func NewRunbookAgentDataDataSource() datasource.DataSource {
    return &RunbookAgentDataDataSource{}
}

// RunbookAgentDataDataSource defines the data source implementation.
type RunbookAgentDataDataSource struct {
    client *Client
}

// RunbookAgentDataDataSourceModel describes the data source data model.
type RunbookAgentDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    Description types.String `tfsdk:"description"`
    Slug types.String `tfsdk:"slug"`
    Key types.String `tfsdk:"key"`
    AgentVersion types.String `tfsdk:"agent_version"`
    LastAlive types.String `tfsdk:"last_alive"`
    ConnectionStatus types.String `tfsdk:"connection_status"`
    HostInfo types.String `tfsdk:"host_info"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    Labels types.Set `tfsdk:"labels"`
}

func (d *RunbookAgentDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_runbook_agent_data"
}

func (d *RunbookAgentDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "runbook_agent_data data source",

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
                MarkdownDescription: "Optional description for this agent. Permissions - Create: [Project Owner, Project Admin, Project Member, Runbook Admin, Runbook Member, Create Runbook Agent], Read: [Project Owner, Project Admin, Project Member, Viewer, Runbook Admin, Runbook Member, Runbook Viewer, Read Runbook Agent], Update: [Project Owner, Project Admin, Runbook Admin, Edit Runbook Agent]",
                Computed: true,
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Friendly globally unique name for your object. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Runbook Admin, Runbook Member, Runbook Viewer, Read Runbook Agent], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "key": schema.StringAttribute{
                MarkdownDescription: "Secret key the agent presents on every request. Never share this key. Reset it to revoke the agent.. Permissions - Create: [Project Owner, Project Admin, Project Member, Runbook Admin, Runbook Member, Create Runbook Agent], Read: [Project Owner, Project Admin, Runbook Admin, Runbook Member, Runbook Viewer], Update: [Project Owner, Project Admin, Runbook Admin, Edit Runbook Agent]",
                Computed: true,
            },
            "agent_version": schema.StringAttribute{
                MarkdownDescription: "Version object",
                Computed: true,
            },
            "last_alive": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "connection_status": schema.StringAttribute{
                MarkdownDescription: "Connected if the agent has heartbeated recently.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Runbook Admin, Runbook Member, Runbook Viewer, Read Runbook Agent], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "host_info": schema.StringAttribute{
                MarkdownDescription: "Self-reported host info (hostname, OS, arch). Updated on each heartbeat.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Runbook Admin, Runbook Member, Runbook Viewer, Read Runbook Agent], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "labels": schema.SetAttribute{
                MarkdownDescription: "Relation to Labels Array where this object is categorized in.. Permissions - Create: [Project Owner, Project Admin, Project Member, Runbook Admin, Runbook Member, Create Runbook Agent], Read: [Project Owner, Project Admin, Project Member, Viewer, Runbook Admin, Runbook Member, Runbook Viewer, Read Runbook Agent], Update: [Project Owner, Project Admin, Project Member, Runbook Admin, Runbook Member, Edit Runbook Agent]",
                Computed: true,
                ElementType: types.StringType,
            },
        },
    }
}

func (d *RunbookAgentDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *RunbookAgentDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data RunbookAgentDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "runbook-agent" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read runbook_agent_data, got error: %s", err))
        return
    }

    var runbookAgentDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &runbookAgentDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse runbook_agent_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := runbookAgentDataResponse["data"].(map[string]interface{}); ok {
        runbookAgentDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := runbookAgentDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := runbookAgentDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := runbookAgentDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := runbookAgentDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := runbookAgentDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := runbookAgentDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := runbookAgentDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := runbookAgentDataResponse["description"].(string); ok {
        data.Description = types.StringValue(val)
    }
    if val, ok := runbookAgentDataResponse["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    }
    if val, ok := runbookAgentDataResponse["key"].(string); ok {
        data.Key = types.StringValue(val)
    }
    if val, ok := runbookAgentDataResponse["agent_version"].(string); ok {
        data.AgentVersion = types.StringValue(val)
    }
    if val, ok := runbookAgentDataResponse["last_alive"].(string); ok {
        data.LastAlive = types.StringValue(val)
    }
    if val, ok := runbookAgentDataResponse["connection_status"].(string); ok {
        data.ConnectionStatus = types.StringValue(val)
    }
    if val, ok := runbookAgentDataResponse["host_info"].(string); ok {
        data.HostInfo = types.StringValue(val)
    }
    if val, ok := runbookAgentDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := runbookAgentDataResponse["labels"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.Labels = setValue
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
