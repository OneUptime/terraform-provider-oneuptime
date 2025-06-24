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
var _ datasource.DataSource = &StatusPageSsoDataDataSource{}

func NewStatusPageSsoDataDataSource() datasource.DataSource {
    return &StatusPageSsoDataDataSource{}
}

// StatusPageSsoDataDataSource defines the data source implementation.
type StatusPageSsoDataDataSource struct {
    client *Client
}

// StatusPageSsoDataDataSourceModel describes the data source data model.
type StatusPageSsoDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    StatusPageId types.String `tfsdk:"status_page_id"`
    Description types.String `tfsdk:"description"`
    SignatureMethod types.String `tfsdk:"signature_method"`
    DigestMethod types.String `tfsdk:"digest_method"`
    SignOnUrl types.String `tfsdk:"sign_on_url"`
    IssuerUrl types.String `tfsdk:"issuer_url"`
    PublicCertificate types.String `tfsdk:"public_certificate"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
    IsEnabled types.Bool `tfsdk:"is_enabled"`
    IsTested types.Bool `tfsdk:"is_tested"`
}

func (d *StatusPageSsoDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_status_page_sso_data"
}

func (d *StatusPageSsoDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "status_page_sso_data data source",

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
            "project_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "status_page_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Create Status Page SSO], Read: [Project Owner, Project Admin, Project User, Public, Read Status Page SSO], Update: [Project Owner, Project Admin, Edit Status Page SSO]",
                Computed: true,
            },
            "signature_method": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Create Status Page SSO], Read: [Project Owner, Project Admin, Read Status Page SSO], Update: [Project Owner, Project Admin, Edit Status Page SSO]",
                Computed: true,
            },
            "digest_method": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Create Status Page SSO], Read: [Project Owner, Project Admin, Read Status Page SSO], Update: [Project Owner, Project Admin, Edit Status Page SSO]",
                Computed: true,
            },
            "sign_on_url": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Create Status Page SSO], Read: [Project Owner, Project Admin, Read Status Page SSO, Project User, Public], Update: [Project Owner, Project Admin, Edit Status Page SSO]",
                Computed: true,
            },
            "issuer_url": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Create Status Page SSO], Read: [Project Owner, Project Admin, Read Status Page SSO], Update: [Project Owner, Project Admin, Edit Status Page SSO]",
                Computed: true,
            },
            "public_certificate": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Create Status Page SSO], Read: [Project Owner, Project Admin, Read Status Page SSO], Update: [Project Owner, Project Admin, Edit Status Page SSO]",
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
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Create Status Page SSO], Read: [Project Owner, Project Admin, Project User, Public, Read Status Page SSO], Update: [Project Owner, Project Admin, Edit Status Page SSO]",
                Computed: true,
            },
            "is_tested": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Create Status Page SSO], Read: [Project Owner, Project Admin, Read Status Page SSO], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
        },
    }
}

func (d *StatusPageSsoDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *StatusPageSsoDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data StatusPageSsoDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "status-page-sso" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read status_page_sso_data, got error: %s", err))
        return
    }

    var statusPageSsoDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &statusPageSsoDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse status_page_sso_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := statusPageSsoDataResponse["data"].(map[string]interface{}); ok {
        statusPageSsoDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := statusPageSsoDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := statusPageSsoDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := statusPageSsoDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := statusPageSsoDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := statusPageSsoDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := statusPageSsoDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := statusPageSsoDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := statusPageSsoDataResponse["status_page_id"].(string); ok {
        data.StatusPageId = types.StringValue(val)
    }
    if val, ok := statusPageSsoDataResponse["description"].(string); ok {
        data.Description = types.StringValue(val)
    }
    if val, ok := statusPageSsoDataResponse["signature_method"].(string); ok {
        data.SignatureMethod = types.StringValue(val)
    }
    if val, ok := statusPageSsoDataResponse["digest_method"].(string); ok {
        data.DigestMethod = types.StringValue(val)
    }
    if val, ok := statusPageSsoDataResponse["sign_on_url"].(string); ok {
        data.SignOnUrl = types.StringValue(val)
    }
    if val, ok := statusPageSsoDataResponse["issuer_url"].(string); ok {
        data.IssuerUrl = types.StringValue(val)
    }
    if val, ok := statusPageSsoDataResponse["public_certificate"].(string); ok {
        data.PublicCertificate = types.StringValue(val)
    }
    if val, ok := statusPageSsoDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := statusPageSsoDataResponse["deleted_by_user_id"].(string); ok {
        data.DeletedByUserId = types.StringValue(val)
    }
    if val, ok := statusPageSsoDataResponse["is_enabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    }
    if val, ok := statusPageSsoDataResponse["is_tested"].(bool); ok {
        data.IsTested = types.BoolValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
