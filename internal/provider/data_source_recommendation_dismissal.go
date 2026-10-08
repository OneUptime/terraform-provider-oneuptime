package provider

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &RecommendationDismissalDataSource{}

func NewRecommendationDismissalDataSource() datasource.DataSource {
    return &RecommendationDismissalDataSource{}
}

// RecommendationDismissalDataSource defines the data source implementation.
type RecommendationDismissalDataSource struct {
    client *Client
}

// RecommendationDismissalDataSourceModel describes the data source data model.
type RecommendationDismissalDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    RecommendationType types.String `tfsdk:"recommendation_type"`
    RecommendationId types.String `tfsdk:"recommendation_id"`
    ResourceType types.String `tfsdk:"resource_type"`
    ResourceId types.String `tfsdk:"resource_id"`
    DismissalReason types.String `tfsdk:"dismissal_reason"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *RecommendationDismissalDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_recommendation_dismissal"
}

func (d *RecommendationDismissalDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Recommendations your team has dismissed. Dismissing hides a recommendation for everyone on the project until it is restored; it never deletes anything that was already created from it. Look up an existing recommendation dismissal by `id`, or by any of its other arguments (`created_by_user_id`, `dismissal_reason`, `recommendation_id`, ...): each one set must match, and exactly one recommendation dismissal may match them all.",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Look up by unique identifier. Leave unset to look up by the other arguments instead.",
                Optional: true,
                Computed: true,
            },
            "created_at": schema.StringAttribute{
                MarkdownDescription: "Date and Time when the object was created.",
                Computed: true,
            },
            "updated_at": schema.StringAttribute{
                MarkdownDescription: "Date and Time when the object was updated.",
                Computed: true,
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "recommendation_type": schema.StringAttribute{
                MarkdownDescription: "Which family of recommendation this dismissal belongs to. See the RecommendationType enum.",
                Optional: true,
                Computed: true,
            },
            "recommendation_id": schema.StringAttribute{
                MarkdownDescription: "The catalog-wide id of the dismissed recommendation, for example Kubernetes:k8s-hpa-at-max-replicas.",
                Optional: true,
                Computed: true,
            },
            "resource_type": schema.StringAttribute{
                MarkdownDescription: "The kind of resource this recommendation was shown on, for example Kubernetes or Docker. Empty for recommendations that are not scoped to a resource.",
                Optional: true,
                Computed: true,
            },
            "resource_id": schema.StringAttribute{
                MarkdownDescription: "ID of the resource this recommendation was shown on. Polymorphic — it points at whichever table Resource Type names — so it carries no foreign key.",
                Optional: true,
                Computed: true,
            },
            "dismissal_reason": schema.StringAttribute{
                MarkdownDescription: "Optional note explaining why this recommendation was dismissed, shown to whoever finds it in the dismissed list later.",
                Optional: true,
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *RecommendationDismissalDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *RecommendationDismissalDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data RecommendationDismissalDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.RecommendationType.IsNull() && !data.RecommendationType.IsUnknown() {
        filters["recommendationType"] = data.RecommendationType.ValueString()
        filterNames = append(filterNames, "recommendation_type = "+fmt.Sprintf("%q", data.RecommendationType.ValueString()))
    }
    if !data.RecommendationId.IsNull() && !data.RecommendationId.IsUnknown() {
        filters["recommendationId"] = data.RecommendationId.ValueString()
        filterNames = append(filterNames, "recommendation_id = "+fmt.Sprintf("%q", data.RecommendationId.ValueString()))
    }
    if !data.ResourceType.IsNull() && !data.ResourceType.IsUnknown() {
        filters["resourceType"] = data.ResourceType.ValueString()
        filterNames = append(filterNames, "resource_type = "+fmt.Sprintf("%q", data.ResourceType.ValueString()))
    }
    if !data.ResourceId.IsNull() && !data.ResourceId.IsUnknown() {
        filters["resourceId"] = data.ResourceId.ValueString()
        filterNames = append(filterNames, "resource_id = "+fmt.Sprintf("%q", data.ResourceId.ValueString()))
    }
    if !data.DismissalReason.IsNull() && !data.DismissalReason.IsUnknown() {
        filters["dismissalReason"] = data.DismissalReason.ValueString()
        filterNames = append(filterNames, "dismissal_reason = "+fmt.Sprintf("%q", data.DismissalReason.ValueString()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the recommendation dismissal up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the recommendation dismissal up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "recommendationType": true,
        "recommendationId": true,
        "resourceType": true,
        "resourceId": true,
        "dismissalReason": true,
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/recommendation-dismissal/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read recommendation_dismissal, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No recommendation dismissal found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read recommendation_dismissal: %s", err))
            return
        }
        if wrapper, ok := itemResponse["data"].(map[string]interface{}); ok {
            item = wrapper
        } else {
            item = itemResponse
        }
    }
    if !hasId {
        listBody := map[string]interface{}{
            "query":  filters,
            "select": selectParam,
            // limit 2 is enough to detect ambiguity without paging.
            "limit": 2,
        }
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/recommendation-dismissal/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list recommendation_dismissal, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list recommendation_dismissal: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No recommendation dismissal matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one recommendation dismissal matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for recommendation_dismissal.")
            return
        }
        item = first
    }

    // Update the model with response data
    if obj, ok := item["_id"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Id = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Id = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Id = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Id = types.StringValue(string(jsonBytes))
        } else {
            data.Id = types.StringNull()
        }
    } else if val, ok := item["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    if obj, ok := item["createdAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CreatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CreatedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CreatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedAt = types.StringNull()
        }
    } else if val, ok := item["createdAt"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    } else {
        data.CreatedAt = types.StringNull()
    }
    if obj, ok := item["updatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UpdatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.UpdatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.UpdatedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.UpdatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.UpdatedAt = types.StringNull()
        }
    } else if val, ok := item["updatedAt"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    } else {
        data.UpdatedAt = types.StringNull()
    }
    if obj, ok := item["projectId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ProjectId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ProjectId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ProjectId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ProjectId = types.StringValue(string(jsonBytes))
        } else {
            data.ProjectId = types.StringNull()
        }
    } else if val, ok := item["projectId"].(string); ok {
        data.ProjectId = types.StringValue(val)
    } else {
        data.ProjectId = types.StringNull()
    }
    if obj, ok := item["recommendationType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RecommendationType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RecommendationType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RecommendationType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RecommendationType = types.StringValue(string(jsonBytes))
        } else {
            data.RecommendationType = types.StringNull()
        }
    } else if val, ok := item["recommendationType"].(string); ok {
        data.RecommendationType = types.StringValue(val)
    } else {
        data.RecommendationType = types.StringNull()
    }
    if obj, ok := item["recommendationId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RecommendationId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RecommendationId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RecommendationId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RecommendationId = types.StringValue(string(jsonBytes))
        } else {
            data.RecommendationId = types.StringNull()
        }
    } else if val, ok := item["recommendationId"].(string); ok {
        data.RecommendationId = types.StringValue(val)
    } else {
        data.RecommendationId = types.StringNull()
    }
    if obj, ok := item["resourceType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ResourceType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ResourceType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ResourceType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ResourceType = types.StringValue(string(jsonBytes))
        } else {
            data.ResourceType = types.StringNull()
        }
    } else if val, ok := item["resourceType"].(string); ok {
        data.ResourceType = types.StringValue(val)
    } else {
        data.ResourceType = types.StringNull()
    }
    if obj, ok := item["resourceId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ResourceId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ResourceId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ResourceId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ResourceId = types.StringValue(string(jsonBytes))
        } else {
            data.ResourceId = types.StringNull()
        }
    } else if val, ok := item["resourceId"].(string); ok {
        data.ResourceId = types.StringValue(val)
    } else {
        data.ResourceId = types.StringNull()
    }
    if obj, ok := item["dismissalReason"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DismissalReason = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DismissalReason = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DismissalReason = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DismissalReason = types.StringValue(string(jsonBytes))
        } else {
            data.DismissalReason = types.StringNull()
        }
    } else if val, ok := item["dismissalReason"].(string); ok {
        data.DismissalReason = types.StringValue(val)
    } else {
        data.DismissalReason = types.StringNull()
    }
    if obj, ok := item["createdByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if val, ok := item["createdByUserId"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
