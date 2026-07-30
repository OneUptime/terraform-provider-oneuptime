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
var _ datasource.DataSource = &RecommendationDismissalDataDataSource{}

func NewRecommendationDismissalDataDataSource() datasource.DataSource {
    return &RecommendationDismissalDataDataSource{}
}

// RecommendationDismissalDataDataSource defines the data source implementation.
type RecommendationDismissalDataDataSource struct {
    client *Client
}

// RecommendationDismissalDataDataSourceModel describes the data source data model.
type RecommendationDismissalDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    RecommendationType types.String `tfsdk:"recommendation_type"`
    RecommendationId types.String `tfsdk:"recommendation_id"`
    ResourceType types.String `tfsdk:"resource_type"`
    ResourceId types.String `tfsdk:"resource_id"`
    DismissalReason types.String `tfsdk:"dismissal_reason"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
}

func (d *RecommendationDismissalDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_recommendation_dismissal_data"
}

func (d *RecommendationDismissalDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "recommendation_dismissal_data data source",

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
            "recommendation_type": schema.StringAttribute{
                MarkdownDescription: "Which family of recommendation this dismissal belongs to. See the RecommendationType enum.. Permissions - Create: [Project Owner, Project Admin, Project Member, Create Recommendation Dismissal], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Recommendation Dismissal], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "recommendation_id": schema.StringAttribute{
                MarkdownDescription: "The catalog-wide id of the dismissed recommendation, for example Kubernetes:k8s-hpa-at-max-replicas.. Permissions - Create: [Project Owner, Project Admin, Project Member, Create Recommendation Dismissal], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Recommendation Dismissal], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "resource_type": schema.StringAttribute{
                MarkdownDescription: "The kind of resource this recommendation was shown on, for example Kubernetes or Docker. Empty for recommendations that are not scoped to a resource.. Permissions - Create: [Project Owner, Project Admin, Project Member, Create Recommendation Dismissal], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Recommendation Dismissal], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "resource_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "dismissal_reason": schema.StringAttribute{
                MarkdownDescription: "Optional note explaining why this recommendation was dismissed, shown to whoever finds it in the dismissed list later.. Permissions - Create: [Project Owner, Project Admin, Project Member, Create Recommendation Dismissal], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Recommendation Dismissal], Update: [Project Owner, Project Admin, Project Member, Edit Recommendation Dismissal]",
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
        },
    }
}

func (d *RecommendationDismissalDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *RecommendationDismissalDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data RecommendationDismissalDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "recommendation-dismissal" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read recommendation_dismissal_data, got error: %s", err))
        return
    }

    var recommendationDismissalDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &recommendationDismissalDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse recommendation_dismissal_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := recommendationDismissalDataResponse["data"].(map[string]interface{}); ok {
        recommendationDismissalDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := recommendationDismissalDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := recommendationDismissalDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := recommendationDismissalDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := recommendationDismissalDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := recommendationDismissalDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := recommendationDismissalDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := recommendationDismissalDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := recommendationDismissalDataResponse["recommendation_type"].(string); ok {
        data.RecommendationType = types.StringValue(val)
    }
    if val, ok := recommendationDismissalDataResponse["recommendation_id"].(string); ok {
        data.RecommendationId = types.StringValue(val)
    }
    if val, ok := recommendationDismissalDataResponse["resource_type"].(string); ok {
        data.ResourceType = types.StringValue(val)
    }
    if val, ok := recommendationDismissalDataResponse["resource_id"].(string); ok {
        data.ResourceId = types.StringValue(val)
    }
    if val, ok := recommendationDismissalDataResponse["dismissal_reason"].(string); ok {
        data.DismissalReason = types.StringValue(val)
    }
    if val, ok := recommendationDismissalDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := recommendationDismissalDataResponse["deleted_by_user_id"].(string); ok {
        data.DeletedByUserId = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
