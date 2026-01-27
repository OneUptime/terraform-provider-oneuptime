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
var _ datasource.DataSource = &AiAgentDataDataSource{}

func NewAiAgentDataDataSource() datasource.DataSource {
    return &AiAgentDataDataSource{}
}

// AiAgentDataDataSource defines the data source implementation.
type AiAgentDataDataSource struct {
    client *Client
}

// AiAgentDataDataSourceModel describes the data source data model.
type AiAgentDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    Key types.String `tfsdk:"key"`
    Description types.String `tfsdk:"description"`
    Slug types.String `tfsdk:"slug"`
    AiAgentVersion types.String `tfsdk:"ai_agent_version"`
    LastAlive types.String `tfsdk:"last_alive"`
    IconFileId types.String `tfsdk:"icon_file_id"`
    ProjectId types.String `tfsdk:"project_id"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    ConnectionStatus types.String `tfsdk:"connection_status"`
    IsDefault types.Bool `tfsdk:"is_default"`
    Labels types.Set `tfsdk:"labels"`
}

func (d *AiAgentDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_ai_agent_data"
}

func (d *AiAgentDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "ai_agent_data data source",

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
            "key": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create AI Agent], Read: [Project Owner, Project Admin], Update: [Project Owner, Project Admin, Project Member, Edit AI Agent]",
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create AI Agent], Read: [Public], Update: [Project Owner, Project Admin, Project Member, Edit AI Agent]",
                Computed: true,
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Friendly globally unique name for your object. Permissions - Create: [No access - you don't have permission for this operation], Read: [Public], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "ai_agent_version": schema.StringAttribute{
                MarkdownDescription: "Version object",
                Computed: true,
            },
            "last_alive": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "icon_file_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "connection_status": schema.StringAttribute{
                MarkdownDescription: "Connection Status of the AI Agent. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read AI Agent], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "is_default": schema.BoolAttribute{
                MarkdownDescription: "Is this the default AI Agent for the project? When set, this agent will be used for automated tasks.. Permissions - Create: [Project Owner, Project Admin, Project Member, Create AI Agent], Read: [Public], Update: [Project Owner, Project Admin, Project Member, Edit AI Agent]",
                Computed: true,
            },
            "labels": schema.SetAttribute{
                MarkdownDescription: "Relation to Labels Array where this object is categorized in.. Permissions - Create: [Project Owner, Project Admin, Project Member, Create AI Agent], Read: [Project Owner, Project Admin, Project Member, Read AI Agent], Update: [Project Owner, Project Admin, Project Member, Edit AI Agent]",
                Computed: true,
                ElementType: types.StringType,
            },
        },
    }
}

func (d *AiAgentDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AiAgentDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data AiAgentDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "ai-agent" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read ai_agent_data, got error: %s", err))
        return
    }

    var aiAgentDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &aiAgentDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse ai_agent_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := aiAgentDataResponse["data"].(map[string]interface{}); ok {
        aiAgentDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := aiAgentDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := aiAgentDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := aiAgentDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := aiAgentDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := aiAgentDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := aiAgentDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := aiAgentDataResponse["key"].(string); ok {
        data.Key = types.StringValue(val)
    }
    if val, ok := aiAgentDataResponse["description"].(string); ok {
        data.Description = types.StringValue(val)
    }
    if val, ok := aiAgentDataResponse["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    }
    if val, ok := aiAgentDataResponse["ai_agent_version"].(string); ok {
        data.AiAgentVersion = types.StringValue(val)
    }
    if val, ok := aiAgentDataResponse["last_alive"].(string); ok {
        data.LastAlive = types.StringValue(val)
    }
    if val, ok := aiAgentDataResponse["icon_file_id"].(string); ok {
        data.IconFileId = types.StringValue(val)
    }
    if val, ok := aiAgentDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := aiAgentDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := aiAgentDataResponse["connection_status"].(string); ok {
        data.ConnectionStatus = types.StringValue(val)
    }
    if val, ok := aiAgentDataResponse["is_default"].(bool); ok {
        data.IsDefault = types.BoolValue(val)
    }
    if val, ok := aiAgentDataResponse["labels"].([]interface{}); ok {
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
