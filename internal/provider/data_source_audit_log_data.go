package provider

import (
    "context"
    "fmt"
    "github.com/hashicorp/terraform-plugin-framework/attr"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &AuditLogDataDataSource{}

func NewAuditLogDataDataSource() datasource.DataSource {
    return &AuditLogDataDataSource{}
}

// AuditLogDataDataSource defines the data source implementation.
type AuditLogDataDataSource struct {
    client *Client
}

// AuditLogDataDataSourceModel describes the data source data model.
type AuditLogDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    ProjectId types.String `tfsdk:"project_id"`
    ResourceType types.String `tfsdk:"resource_type"`
    ResourceId types.String `tfsdk:"resource_id"`
    ResourceName types.String `tfsdk:"resource_name"`
    Action types.String `tfsdk:"action"`
    UserId types.String `tfsdk:"user_id"`
    UserName types.String `tfsdk:"user_name"`
    UserEmail types.String `tfsdk:"user_email"`
    UserType types.String `tfsdk:"user_type"`
    ApiKeyId types.String `tfsdk:"api_key_id"`
    ApiKeyName types.String `tfsdk:"api_key_name"`
    Changes types.Set `tfsdk:"changes"`
}

func (d *AuditLogDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_audit_log_data"
}

func (d *AuditLogDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "audit_log_data data source",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Identifier to filter by",
                Optional: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Name to filter by",
                Optional: true,
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "Project ID",
                Computed: true,
            },
            "resource_type": schema.StringAttribute{
                MarkdownDescription: "Resource Type",
                Computed: true,
            },
            "resource_id": schema.StringAttribute{
                MarkdownDescription: "Resource ID",
                Computed: true,
            },
            "resource_name": schema.StringAttribute{
                MarkdownDescription: "Resource Name",
                Computed: true,
            },
            "action": schema.StringAttribute{
                MarkdownDescription: "Action",
                Computed: true,
            },
            "user_id": schema.StringAttribute{
                MarkdownDescription: "User ID",
                Computed: true,
            },
            "user_name": schema.StringAttribute{
                MarkdownDescription: "User Name",
                Computed: true,
            },
            "user_email": schema.StringAttribute{
                MarkdownDescription: "User Email",
                Computed: true,
            },
            "user_type": schema.StringAttribute{
                MarkdownDescription: "User Type",
                Computed: true,
            },
            "api_key_id": schema.StringAttribute{
                MarkdownDescription: "API Key ID",
                Computed: true,
            },
            "api_key_name": schema.StringAttribute{
                MarkdownDescription: "API Key Name",
                Computed: true,
            },
            "changes": schema.SetAttribute{
                MarkdownDescription: "Changes",
                Computed: true,
                ElementType: types.StringType,
            },
        },
    }
}

func (d *AuditLogDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AuditLogDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data AuditLogDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "audit-log" + "/" + data.Id.ValueString()
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read audit_log_data, got error: %s", err))
        return
    }

    var auditLogDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &auditLogDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse audit_log_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := auditLogDataResponse["data"].(map[string]interface{}); ok {
        auditLogDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := auditLogDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := auditLogDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := auditLogDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := auditLogDataResponse["resource_type"].(string); ok {
        data.ResourceType = types.StringValue(val)
    }
    if val, ok := auditLogDataResponse["resource_id"].(string); ok {
        data.ResourceId = types.StringValue(val)
    }
    if val, ok := auditLogDataResponse["resource_name"].(string); ok {
        data.ResourceName = types.StringValue(val)
    }
    if val, ok := auditLogDataResponse["action"].(string); ok {
        data.Action = types.StringValue(val)
    }
    if val, ok := auditLogDataResponse["user_id"].(string); ok {
        data.UserId = types.StringValue(val)
    }
    if val, ok := auditLogDataResponse["user_name"].(string); ok {
        data.UserName = types.StringValue(val)
    }
    if val, ok := auditLogDataResponse["user_email"].(string); ok {
        data.UserEmail = types.StringValue(val)
    }
    if val, ok := auditLogDataResponse["user_type"].(string); ok {
        data.UserType = types.StringValue(val)
    }
    if val, ok := auditLogDataResponse["api_key_id"].(string); ok {
        data.ApiKeyId = types.StringValue(val)
    }
    if val, ok := auditLogDataResponse["api_key_name"].(string); ok {
        data.ApiKeyName = types.StringValue(val)
    }
    if val, ok := auditLogDataResponse["changes"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.Changes = setValue
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
