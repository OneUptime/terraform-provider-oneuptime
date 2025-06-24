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
var _ datasource.DataSource = &ServiceInServiceCatalogDataDataSource{}

func NewServiceInServiceCatalogDataDataSource() datasource.DataSource {
    return &ServiceInServiceCatalogDataDataSource{}
}

// ServiceInServiceCatalogDataDataSource defines the data source implementation.
type ServiceInServiceCatalogDataDataSource struct {
    client *Client
}

// ServiceInServiceCatalogDataDataSourceModel describes the data source data model.
type ServiceInServiceCatalogDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    Slug types.String `tfsdk:"slug"`
    Description types.String `tfsdk:"description"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
    Labels types.List `tfsdk:"labels"`
    ServiceColor types.String `tfsdk:"service_color"`
    ServiceLanguage types.String `tfsdk:"service_language"`
    TechStack types.String `tfsdk:"tech_stack"`
}

func (d *ServiceInServiceCatalogDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_service_in_service_catalog_data"
}

func (d *ServiceInServiceCatalogDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "service_in_service_catalog_data data source",

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
            "slug": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Project Member, Read Service Catalog], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Service Catalog], Read: [Project Owner, Project Admin, Project Member, Project Member, Read Service Catalog], Update: [Project Owner, Project Admin, Project Member, Edit Service Catalog]",
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
            "labels": schema.ListAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Service Catalog], Read: [Project Owner, Project Admin, Project Member, Project Member, Read Service Catalog], Update: [Project Owner, Project Admin, Project Member, Edit Service Catalog]",
                Computed: true,
                ElementType: types.StringType,
            },
            "service_color": schema.StringAttribute{
                MarkdownDescription: "Color object",
                Computed: true,
            },
            "service_language": schema.StringAttribute{
                MarkdownDescription: "Service Language",
                Computed: true,
            },
            "tech_stack": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Service Catalog], Read: [Project Owner, Project Admin, Project Member, Project Member, Read Service Catalog], Update: [Project Owner, Project Admin, Project Member, Edit Service Catalog]",
                Computed: true,
            },
        },
    }
}

func (d *ServiceInServiceCatalogDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ServiceInServiceCatalogDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data ServiceInServiceCatalogDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "service-catalog" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read service_in_service_catalog_data, got error: %s", err))
        return
    }

    var serviceInServiceCatalogDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &serviceInServiceCatalogDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse service_in_service_catalog_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := serviceInServiceCatalogDataResponse["data"].(map[string]interface{}); ok {
        serviceInServiceCatalogDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := serviceInServiceCatalogDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := serviceInServiceCatalogDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := serviceInServiceCatalogDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := serviceInServiceCatalogDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := serviceInServiceCatalogDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := serviceInServiceCatalogDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := serviceInServiceCatalogDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := serviceInServiceCatalogDataResponse["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    }
    if val, ok := serviceInServiceCatalogDataResponse["description"].(string); ok {
        data.Description = types.StringValue(val)
    }
    if val, ok := serviceInServiceCatalogDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := serviceInServiceCatalogDataResponse["deleted_by_user_id"].(string); ok {
        data.DeletedByUserId = types.StringValue(val)
    }
    if val, ok := serviceInServiceCatalogDataResponse["labels"].([]interface{}); ok {
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
    if val, ok := serviceInServiceCatalogDataResponse["service_color"].(string); ok {
        data.ServiceColor = types.StringValue(val)
    }
    if val, ok := serviceInServiceCatalogDataResponse["service_language"].(string); ok {
        data.ServiceLanguage = types.StringValue(val)
    }
    if val, ok := serviceInServiceCatalogDataResponse["tech_stack"].(string); ok {
        data.TechStack = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
