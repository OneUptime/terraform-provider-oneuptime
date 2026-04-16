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
var _ datasource.DataSource = &LogSavedViewDataDataSource{}

func NewLogSavedViewDataDataSource() datasource.DataSource {
    return &LogSavedViewDataDataSource{}
}

// LogSavedViewDataDataSource defines the data source implementation.
type LogSavedViewDataDataSource struct {
    client *Client
}

// LogSavedViewDataDataSourceModel describes the data source data model.
type LogSavedViewDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
    Query types.String `tfsdk:"query"`
    Columns types.String `tfsdk:"columns"`
    SortField types.String `tfsdk:"sort_field"`
    SortOrder types.String `tfsdk:"sort_order"`
    PageSize types.Number `tfsdk:"page_size"`
    IsDefault types.Bool `tfsdk:"is_default"`
}

func (d *LogSavedViewDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_log_saved_view_data"
}

func (d *LogSavedViewDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "log_saved_view_data data source",

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
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "deleted_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "query": schema.StringAttribute{
                MarkdownDescription: "Serialized log query for this saved view.. Permissions - Create: [Project Owner, Project Admin, Project Member], Read: [Project Owner, Project Admin, Project Member, Viewer, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member]",
                Computed: true,
            },
            "columns": schema.StringAttribute{
                MarkdownDescription: "Selected log table columns for this saved view.. Permissions - Create: [Project Owner, Project Admin, Project Member], Read: [Project Owner, Project Admin, Project Member, Viewer, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member]",
                Computed: true,
            },
            "sort_field": schema.StringAttribute{
                MarkdownDescription: "Active sort field for this saved log view.. Permissions - Create: [Project Owner, Project Admin, Project Member], Read: [Project Owner, Project Admin, Project Member, Viewer, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member]",
                Computed: true,
            },
            "sort_order": schema.StringAttribute{
                MarkdownDescription: "Sort order for this saved log view.. Permissions - Create: [Project Owner, Project Admin, Project Member], Read: [Project Owner, Project Admin, Project Member, Viewer, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member]",
                Computed: true,
            },
            "page_size": schema.NumberAttribute{
                MarkdownDescription: "Number of logs per page for this saved view.. Permissions - Create: [Project Owner, Project Admin, Project Member], Read: [Project Owner, Project Admin, Project Member, Viewer, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member]",
                Computed: true,
            },
            "is_default": schema.BoolAttribute{
                MarkdownDescription: "Whether this saved log view should be applied by default.. Permissions - Create: [Project Owner, Project Admin, Project Member], Read: [Project Owner, Project Admin, Project Member, Viewer, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member]",
                Computed: true,
            },
        },
    }
}

func (d *LogSavedViewDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *LogSavedViewDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data LogSavedViewDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "log-saved-view" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read log_saved_view_data, got error: %s", err))
        return
    }

    var logSavedViewDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &logSavedViewDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse log_saved_view_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := logSavedViewDataResponse["data"].(map[string]interface{}); ok {
        logSavedViewDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := logSavedViewDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := logSavedViewDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := logSavedViewDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := logSavedViewDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := logSavedViewDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := logSavedViewDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := logSavedViewDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := logSavedViewDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := logSavedViewDataResponse["deleted_by_user_id"].(string); ok {
        data.DeletedByUserId = types.StringValue(val)
    }
    if val, ok := logSavedViewDataResponse["query"].(string); ok {
        data.Query = types.StringValue(val)
    }
    if val, ok := logSavedViewDataResponse["columns"].(string); ok {
        data.Columns = types.StringValue(val)
    }
    if val, ok := logSavedViewDataResponse["sort_field"].(string); ok {
        data.SortField = types.StringValue(val)
    }
    if val, ok := logSavedViewDataResponse["sort_order"].(string); ok {
        data.SortOrder = types.StringValue(val)
    }
    if val, ok := logSavedViewDataResponse["page_size"].(float64); ok {
        data.PageSize = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := logSavedViewDataResponse["is_default"].(bool); ok {
        data.IsDefault = types.BoolValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
