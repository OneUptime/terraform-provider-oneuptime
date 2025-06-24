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
var _ datasource.DataSource = &StatusPageDomainDataDataSource{}

func NewStatusPageDomainDataDataSource() datasource.DataSource {
    return &StatusPageDomainDataDataSource{}
}

// StatusPageDomainDataDataSource defines the data source implementation.
type StatusPageDomainDataDataSource struct {
    client *Client
}

// StatusPageDomainDataDataSourceModel describes the data source data model.
type StatusPageDomainDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    DomainId types.String `tfsdk:"domain_id"`
    StatusPageId types.String `tfsdk:"status_page_id"`
    Subdomain types.String `tfsdk:"subdomain"`
    FullDomain types.String `tfsdk:"full_domain"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    IsCnameVerified types.Bool `tfsdk:"is_cname_verified"`
    IsSslOrdered types.Bool `tfsdk:"is_ssl_ordered"`
    IsSslProvisioned types.Bool `tfsdk:"is_ssl_provisioned"`
    CustomCertificate types.String `tfsdk:"custom_certificate"`
    CustomCertificateKey types.String `tfsdk:"custom_certificate_key"`
    IsCustomCertificate types.Bool `tfsdk:"is_custom_certificate"`
}

func (d *StatusPageDomainDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_status_page_domain_data"
}

func (d *StatusPageDomainDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "status_page_domain_data data source",

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
                MarkdownDescription: "Version",
                Computed: true,
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "domain_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "status_page_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "subdomain": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page Domain], Read: [Project Owner, Project Admin, Project Member, Read Status Page Domain], Update: [Project Owner, Project Admin, Project Member, Edit Status Page Domain]",
                Computed: true,
            },
            "full_domain": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page Domain], Read: [Project Owner, Project Admin, Project Member, Read Status Page Domain], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "is_cname_verified": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Status Page Domain], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "is_ssl_ordered": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Status Page Domain], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "is_ssl_provisioned": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Status Page Domain], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "custom_certificate": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page Domain], Read: [Project Owner, Project Admin, Project Member, Read Status Page Domain], Update: [Project Owner, Project Admin, Project Member, Edit Status Page Domain]",
                Computed: true,
            },
            "custom_certificate_key": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page Domain], Read: [Project Owner, Project Admin, Project Member, Read Status Page Domain], Update: [Project Owner, Project Admin, Project Member, Edit Status Page Domain]",
                Computed: true,
            },
            "is_custom_certificate": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page Domain], Read: [Project Owner, Project Admin, Project Member, Read Status Page Domain], Update: [Project Owner, Project Admin, Project Member, Edit Status Page Domain]",
                Computed: true,
            },
        },
    }
}

func (d *StatusPageDomainDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *StatusPageDomainDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data StatusPageDomainDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "status-page-domain" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read status_page_domain_data, got error: %s", err))
        return
    }

    var statusPageDomainDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &statusPageDomainDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse status_page_domain_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := statusPageDomainDataResponse["data"].(map[string]interface{}); ok {
        statusPageDomainDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := statusPageDomainDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := statusPageDomainDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := statusPageDomainDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := statusPageDomainDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := statusPageDomainDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := statusPageDomainDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := statusPageDomainDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := statusPageDomainDataResponse["domain_id"].(string); ok {
        data.DomainId = types.StringValue(val)
    }
    if val, ok := statusPageDomainDataResponse["status_page_id"].(string); ok {
        data.StatusPageId = types.StringValue(val)
    }
    if val, ok := statusPageDomainDataResponse["subdomain"].(string); ok {
        data.Subdomain = types.StringValue(val)
    }
    if val, ok := statusPageDomainDataResponse["full_domain"].(string); ok {
        data.FullDomain = types.StringValue(val)
    }
    if val, ok := statusPageDomainDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := statusPageDomainDataResponse["is_cname_verified"].(bool); ok {
        data.IsCnameVerified = types.BoolValue(val)
    }
    if val, ok := statusPageDomainDataResponse["is_ssl_ordered"].(bool); ok {
        data.IsSslOrdered = types.BoolValue(val)
    }
    if val, ok := statusPageDomainDataResponse["is_ssl_provisioned"].(bool); ok {
        data.IsSslProvisioned = types.BoolValue(val)
    }
    if val, ok := statusPageDomainDataResponse["custom_certificate"].(string); ok {
        data.CustomCertificate = types.StringValue(val)
    }
    if val, ok := statusPageDomainDataResponse["custom_certificate_key"].(string); ok {
        data.CustomCertificateKey = types.StringValue(val)
    }
    if val, ok := statusPageDomainDataResponse["is_custom_certificate"].(bool); ok {
        data.IsCustomCertificate = types.BoolValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
