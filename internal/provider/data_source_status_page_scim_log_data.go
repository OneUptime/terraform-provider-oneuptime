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
var _ datasource.DataSource = &StatusPageScimLogDataDataSource{}

func NewStatusPageScimLogDataDataSource() datasource.DataSource {
    return &StatusPageScimLogDataDataSource{}
}

// StatusPageScimLogDataDataSource defines the data source implementation.
type StatusPageScimLogDataDataSource struct {
    client *Client
}

// StatusPageScimLogDataDataSourceModel describes the data source data model.
type StatusPageScimLogDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    StatusPageId types.String `tfsdk:"status_page_id"`
    StatusPageScimId types.String `tfsdk:"status_page_scim_id"`
    OperationType types.String `tfsdk:"operation_type"`
    Status types.String `tfsdk:"status"`
    StatusMessage types.String `tfsdk:"status_message"`
    LogBody types.String `tfsdk:"log_body"`
    HttpMethod types.String `tfsdk:"http_method"`
    RequestPath types.String `tfsdk:"request_path"`
    HttpStatusCode types.Number `tfsdk:"http_status_code"`
    AffectedUserEmail types.String `tfsdk:"affected_user_email"`
}

func (d *StatusPageScimLogDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_status_page_scim_log_data"
}

func (d *StatusPageScimLogDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "status_page_scim_log_data data source",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Identifier to filter by",
                Optional: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Name to filter by",
                Optional: true,
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
            "status_page_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "status_page_scim_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "operation_type": schema.StringAttribute{
                MarkdownDescription: "Type of SCIM operation (e.g., CreateUser, UpdateUser, DeleteUser, ListUsers, GetUser, BulkOperation). Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Status Page SCIM Log, Read All Project Resources], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "status": schema.StringAttribute{
                MarkdownDescription: "Status of the SCIM operation. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Status Page SCIM Log, Read All Project Resources], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "status_message": schema.StringAttribute{
                MarkdownDescription: "Short error or status description. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Status Page SCIM Log, Read All Project Resources], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "log_body": schema.StringAttribute{
                MarkdownDescription: "Detailed JSON with request/response data. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Status Page SCIM Log, Read All Project Resources], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "http_method": schema.StringAttribute{
                MarkdownDescription: "HTTP method used (GET, POST, PUT, PATCH, DELETE). Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Status Page SCIM Log, Read All Project Resources], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "request_path": schema.StringAttribute{
                MarkdownDescription: "The SCIM endpoint path. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Status Page SCIM Log, Read All Project Resources], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "http_status_code": schema.NumberAttribute{
                MarkdownDescription: "Response HTTP status code. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Status Page SCIM Log, Read All Project Resources], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "affected_user_email": schema.StringAttribute{
                MarkdownDescription: "Email object",
                Computed: true,
            },
        },
    }
}

func (d *StatusPageScimLogDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *StatusPageScimLogDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data StatusPageScimLogDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "status-page-scim-log" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read status_page_scim_log_data, got error: %s", err))
        return
    }

    var statusPageScimLogDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &statusPageScimLogDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse status_page_scim_log_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := statusPageScimLogDataResponse["data"].(map[string]interface{}); ok {
        statusPageScimLogDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := statusPageScimLogDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := statusPageScimLogDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := statusPageScimLogDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := statusPageScimLogDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := statusPageScimLogDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := statusPageScimLogDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := statusPageScimLogDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := statusPageScimLogDataResponse["status_page_id"].(string); ok {
        data.StatusPageId = types.StringValue(val)
    }
    if val, ok := statusPageScimLogDataResponse["status_page_scim_id"].(string); ok {
        data.StatusPageScimId = types.StringValue(val)
    }
    if val, ok := statusPageScimLogDataResponse["operation_type"].(string); ok {
        data.OperationType = types.StringValue(val)
    }
    if val, ok := statusPageScimLogDataResponse["status"].(string); ok {
        data.Status = types.StringValue(val)
    }
    if val, ok := statusPageScimLogDataResponse["status_message"].(string); ok {
        data.StatusMessage = types.StringValue(val)
    }
    if val, ok := statusPageScimLogDataResponse["log_body"].(string); ok {
        data.LogBody = types.StringValue(val)
    }
    if val, ok := statusPageScimLogDataResponse["http_method"].(string); ok {
        data.HttpMethod = types.StringValue(val)
    }
    if val, ok := statusPageScimLogDataResponse["request_path"].(string); ok {
        data.RequestPath = types.StringValue(val)
    }
    if val, ok := statusPageScimLogDataResponse["http_status_code"].(float64); ok {
        data.HttpStatusCode = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := statusPageScimLogDataResponse["affected_user_email"].(string); ok {
        data.AffectedUserEmail = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
