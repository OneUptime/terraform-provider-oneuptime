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
var _ datasource.DataSource = &ProbeDataDataSource{}

func NewProbeDataDataSource() datasource.DataSource {
    return &ProbeDataDataSource{}
}

// ProbeDataDataSource defines the data source implementation.
type ProbeDataDataSource struct {
    client *Client
}

// ProbeDataDataSourceModel describes the data source data model.
type ProbeDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    Key types.String `tfsdk:"key"`
    Description types.String `tfsdk:"description"`
    Slug types.String `tfsdk:"slug"`
    ProbeVersion types.String `tfsdk:"probe_version"`
    LastAlive types.String `tfsdk:"last_alive"`
    IconFileId types.String `tfsdk:"icon_file_id"`
    ProjectId types.String `tfsdk:"project_id"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    ShouldAutoEnableProbeOnNewMonitors types.Bool `tfsdk:"should_auto_enable_probe_on_new_monitors"`
    ConnectionStatus types.String `tfsdk:"connection_status"`
    Labels types.List `tfsdk:"labels"`
}

func (d *ProbeDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_probe_data"
}

func (d *ProbeDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "probe_data data source",

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
                MarkdownDescription: "Version",
                Computed: true,
            },
            "key": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Probe], Read: [Project Owner, Project Admin], Update: [Project Owner, Project Admin, Project Member, Edit Probe]",
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Name object",
                Computed: true,
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Public], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "probe_version": schema.StringAttribute{
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
            "should_auto_enable_probe_on_new_monitors": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Edit Status Page]",
                Computed: true,
            },
            "connection_status": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Status Page], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "labels": schema.ListAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Edit Status Page]",
                Computed: true,
                ElementType: types.StringType,
            },
        },
    }
}

func (d *ProbeDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ProbeDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data ProbeDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "probe" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read probe_data, got error: %s", err))
        return
    }

    var probeDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &probeDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse probe_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := probeDataResponse["data"].(map[string]interface{}); ok {
        probeDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := probeDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := probeDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := probeDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := probeDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := probeDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := probeDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := probeDataResponse["key"].(string); ok {
        data.Key = types.StringValue(val)
    }
    if val, ok := probeDataResponse["description"].(string); ok {
        data.Description = types.StringValue(val)
    }
    if val, ok := probeDataResponse["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    }
    if val, ok := probeDataResponse["probe_version"].(string); ok {
        data.ProbeVersion = types.StringValue(val)
    }
    if val, ok := probeDataResponse["last_alive"].(string); ok {
        data.LastAlive = types.StringValue(val)
    }
    if val, ok := probeDataResponse["icon_file_id"].(string); ok {
        data.IconFileId = types.StringValue(val)
    }
    if val, ok := probeDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := probeDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := probeDataResponse["should_auto_enable_probe_on_new_monitors"].(bool); ok {
        data.ShouldAutoEnableProbeOnNewMonitors = types.BoolValue(val)
    }
    if val, ok := probeDataResponse["connection_status"].(string); ok {
        data.ConnectionStatus = types.StringValue(val)
    }
    if val, ok := probeDataResponse["labels"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        listValue, _ := types.ListValue(types.StringType, elements)
        data.Labels = listValue
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
