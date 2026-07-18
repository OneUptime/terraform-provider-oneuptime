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
var _ datasource.DataSource = &MetricSavedViewDataDataSource{}

func NewMetricSavedViewDataDataSource() datasource.DataSource {
    return &MetricSavedViewDataDataSource{}
}

// MetricSavedViewDataDataSource defines the data source implementation.
type MetricSavedViewDataDataSource struct {
    client *Client
}

// MetricSavedViewDataDataSourceModel describes the data source data model.
type MetricSavedViewDataDataSourceModel struct {
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
    IsDefault types.Bool `tfsdk:"is_default"`
    ViewType types.String `tfsdk:"view_type"`
}

func (d *MetricSavedViewDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_metric_saved_view_data"
}

func (d *MetricSavedViewDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "metric_saved_view_data data source",

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
                MarkdownDescription: "Serialized metrics explorer view state (search, filters, time range, page size) for this saved view.. Permissions - Create: [Project Owner, Project Admin, Project Member], Read: [Project Owner, Project Admin, Project Member, Viewer], Update: [Project Owner, Project Admin, Project Member]",
                Computed: true,
            },
            "is_default": schema.BoolAttribute{
                MarkdownDescription: "Whether this saved metric view should be applied by default.. Permissions - Create: [Project Owner, Project Admin, Project Member], Read: [Project Owner, Project Admin, Project Member, Viewer], Update: [Project Owner, Project Admin, Project Member]",
                Computed: true,
            },
            "view_type": schema.StringAttribute{
                MarkdownDescription: "Which surface this saved view belongs to ('list' or 'explorer'). Null means 'list' — rows created before this column existed all came from the metric list page.. Permissions - Create: [Project Owner, Project Admin, Project Member], Read: [Project Owner, Project Admin, Project Member, Viewer], Update: [Project Owner, Project Admin, Project Member]",
                Computed: true,
            },
        },
    }
}

func (d *MetricSavedViewDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *MetricSavedViewDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data MetricSavedViewDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "metric-saved-view" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read metric_saved_view_data, got error: %s", err))
        return
    }

    var metricSavedViewDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &metricSavedViewDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse metric_saved_view_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := metricSavedViewDataResponse["data"].(map[string]interface{}); ok {
        metricSavedViewDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := metricSavedViewDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := metricSavedViewDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := metricSavedViewDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := metricSavedViewDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := metricSavedViewDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := metricSavedViewDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := metricSavedViewDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := metricSavedViewDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := metricSavedViewDataResponse["deleted_by_user_id"].(string); ok {
        data.DeletedByUserId = types.StringValue(val)
    }
    if val, ok := metricSavedViewDataResponse["query"].(string); ok {
        data.Query = types.StringValue(val)
    }
    if val, ok := metricSavedViewDataResponse["is_default"].(bool); ok {
        data.IsDefault = types.BoolValue(val)
    }
    if val, ok := metricSavedViewDataResponse["view_type"].(string); ok {
        data.ViewType = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
