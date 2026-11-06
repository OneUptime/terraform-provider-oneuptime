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
var _ datasource.DataSource = &CloudCostRecommendationDataSource{}

func NewCloudCostRecommendationDataSource() datasource.DataSource {
    return &CloudCostRecommendationDataSource{}
}

// CloudCostRecommendationDataSource defines the data source implementation.
type CloudCostRecommendationDataSource struct {
    client *Client
}

// CloudCostRecommendationDataSourceModel describes the data source data model.
type CloudCostRecommendationDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    CloudAccountId types.String `tfsdk:"cloud_account_id"`
    CloudResourceId types.String `tfsdk:"cloud_resource_id"`
    ProviderValue types.String `tfsdk:"provider_value"`
    RecommendationId types.String `tfsdk:"recommendation_id"`
    Region types.String `tfsdk:"region"`
    ResourceId types.String `tfsdk:"resource_id"`
    ResourceType types.String `tfsdk:"resource_type"`
    ResourceName types.String `tfsdk:"resource_name"`
    Kind types.String `tfsdk:"kind"`
    Action types.String `tfsdk:"action"`
    CurrentConfiguration types.String `tfsdk:"current_configuration"`
    RecommendedConfiguration types.String `tfsdk:"recommended_configuration"`
    MonthlySavings types.Number `tfsdk:"monthly_savings"`
    MonthlyCost types.Number `tfsdk:"monthly_cost"`
    Currency types.String `tfsdk:"currency"`
    Effort types.String `tfsdk:"effort"`
    Source types.String `tfsdk:"source"`
    LastSeenAt types.String `tfsdk:"last_seen_at"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *CloudCostRecommendationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_cloud_cost_recommendation"
}

func (d *CloudCostRecommendationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Ways to spend less in your connected cloud accounts - rightsizing, idle resources, newer generations and commitments - with what each saves a month, from AWS Cost Optimization Hub, Azure Advisor and Google Cloud's Recommender. Written by the cloud connections; not user-editable. Look up an existing cloud cost recommendation by `id`, or by any of its other arguments (`action`, `cloud_account_id`, `cloud_resource_id`, ...): each one set must match, and exactly one cloud cost recommendation may match them all.",

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
                MarkdownDescription: "ID of the project this recommendation belongs to. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "cloud_account_id": schema.StringAttribute{
                MarkdownDescription: "ID of the account the recommendation is for. The ID of a `oneuptime_cloud_account`.",
                Optional: true,
                Computed: true,
                Validators: recordIDValidators("oneuptime_cloud_account"),
            },
            "cloud_resource_id": schema.StringAttribute{
                MarkdownDescription: "ID of the resource the recommendation is about, when OneUptime holds it. The ID of a `oneuptime_cloud_resource`.",
                Optional: true,
                Computed: true,
                Validators: recordIDValidators("oneuptime_cloud_resource"),
            },
            "provider_value": schema.StringAttribute{
                MarkdownDescription: "Which cloud the account is in: 'aws', 'azure' or 'gcp'.",
                Optional: true,
                Computed: true,
            },
            "recommendation_id": schema.StringAttribute{
                MarkdownDescription: "The recommender's own id for the recommendation, kept from one read to the next.",
                Optional: true,
                Computed: true,
            },
            "region": schema.StringAttribute{
                MarkdownDescription: "The Region the resource is in; empty when none is given.",
                Optional: true,
                Computed: true,
            },
            "resource_id": schema.StringAttribute{
                MarkdownDescription: "The resource the recommendation is about, as its cloud names it: an ARN, an ARM id, a full resource name.",
                Optional: true,
                Computed: true,
            },
            "resource_type": schema.StringAttribute{
                MarkdownDescription: "What kind of resource it is, in the recommender's words.",
                Optional: true,
                Computed: true,
            },
            "resource_name": schema.StringAttribute{
                MarkdownDescription: "What the resource is called, when the recommender says.",
                Optional: true,
                Computed: true,
            },
            "kind": schema.StringAttribute{
                MarkdownDescription: "What it asks: rightsize, stop or delete something idle, upgrade to a newer generation, commit to use, or other.",
                Optional: true,
                Computed: true,
            },
            "action": schema.StringAttribute{
                MarkdownDescription: "What to do, in the recommender's words: 'Rightsize', 'Delete the unattached disk'.",
                Optional: true,
                Computed: true,
            },
            "current_configuration": schema.StringAttribute{
                MarkdownDescription: "What the resource is now, in the recommender's words.",
                Optional: true,
                Computed: true,
            },
            "recommended_configuration": schema.StringAttribute{
                MarkdownDescription: "What the recommender would make it.",
                Optional: true,
                Computed: true,
            },
            "monthly_savings": schema.NumberAttribute{
                MarkdownDescription: "What doing it saves a month, as the recommender estimates it, in the recommendation's currency.",
                Optional: true,
                Computed: true,
            },
            "monthly_cost": schema.NumberAttribute{
                MarkdownDescription: "What the resource costs a month now, when the recommender says.",
                Optional: true,
                Computed: true,
            },
            "currency": schema.StringAttribute{
                MarkdownDescription: "The currency the figures are in.",
                Optional: true,
                Computed: true,
            },
            "effort": schema.StringAttribute{
                MarkdownDescription: "How much work it is, in the recommender's words: 'Low', 'Medium', 'High'.",
                Optional: true,
                Computed: true,
            },
            "source": schema.StringAttribute{
                MarkdownDescription: "Which recommender made it: AWS Cost Optimization Hub, Azure Advisor or Google Cloud Recommender.",
                Optional: true,
                Computed: true,
            },
            "last_seen_at": schema.StringAttribute{
                MarkdownDescription: "When the account's connection last read the recommendation.",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the user who created this row (the connection writes as root).",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *CloudCostRecommendationDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *CloudCostRecommendationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data CloudCostRecommendationDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.CloudAccountId.IsNull() && !data.CloudAccountId.IsUnknown() {
        filters["cloudAccountId"] = data.CloudAccountId.ValueString()
        filterNames = append(filterNames, "cloud_account_id = "+fmt.Sprintf("%q", data.CloudAccountId.ValueString()))
    }
    if !data.CloudResourceId.IsNull() && !data.CloudResourceId.IsUnknown() {
        filters["cloudResourceId"] = data.CloudResourceId.ValueString()
        filterNames = append(filterNames, "cloud_resource_id = "+fmt.Sprintf("%q", data.CloudResourceId.ValueString()))
    }
    if !data.ProviderValue.IsNull() && !data.ProviderValue.IsUnknown() {
        filters["provider"] = data.ProviderValue.ValueString()
        filterNames = append(filterNames, "provider_value = "+fmt.Sprintf("%q", data.ProviderValue.ValueString()))
    }
    if !data.RecommendationId.IsNull() && !data.RecommendationId.IsUnknown() {
        filters["recommendationId"] = data.RecommendationId.ValueString()
        filterNames = append(filterNames, "recommendation_id = "+fmt.Sprintf("%q", data.RecommendationId.ValueString()))
    }
    if !data.Region.IsNull() && !data.Region.IsUnknown() {
        filters["region"] = data.Region.ValueString()
        filterNames = append(filterNames, "region = "+fmt.Sprintf("%q", data.Region.ValueString()))
    }
    if !data.ResourceId.IsNull() && !data.ResourceId.IsUnknown() {
        filters["resourceId"] = data.ResourceId.ValueString()
        filterNames = append(filterNames, "resource_id = "+fmt.Sprintf("%q", data.ResourceId.ValueString()))
    }
    if !data.ResourceType.IsNull() && !data.ResourceType.IsUnknown() {
        filters["resourceType"] = data.ResourceType.ValueString()
        filterNames = append(filterNames, "resource_type = "+fmt.Sprintf("%q", data.ResourceType.ValueString()))
    }
    if !data.ResourceName.IsNull() && !data.ResourceName.IsUnknown() {
        filters["resourceName"] = data.ResourceName.ValueString()
        filterNames = append(filterNames, "resource_name = "+fmt.Sprintf("%q", data.ResourceName.ValueString()))
    }
    if !data.Kind.IsNull() && !data.Kind.IsUnknown() {
        filters["kind"] = data.Kind.ValueString()
        filterNames = append(filterNames, "kind = "+fmt.Sprintf("%q", data.Kind.ValueString()))
    }
    if !data.Action.IsNull() && !data.Action.IsUnknown() {
        filters["action"] = data.Action.ValueString()
        filterNames = append(filterNames, "action = "+fmt.Sprintf("%q", data.Action.ValueString()))
    }
    if !data.CurrentConfiguration.IsNull() && !data.CurrentConfiguration.IsUnknown() {
        filters["currentConfiguration"] = data.CurrentConfiguration.ValueString()
        filterNames = append(filterNames, "current_configuration = "+fmt.Sprintf("%q", data.CurrentConfiguration.ValueString()))
    }
    if !data.RecommendedConfiguration.IsNull() && !data.RecommendedConfiguration.IsUnknown() {
        filters["recommendedConfiguration"] = data.RecommendedConfiguration.ValueString()
        filterNames = append(filterNames, "recommended_configuration = "+fmt.Sprintf("%q", data.RecommendedConfiguration.ValueString()))
    }
    if !data.MonthlySavings.IsNull() && !data.MonthlySavings.IsUnknown() {
        filters["monthlySavings"] = lookupNumber(data.MonthlySavings)
        filterNames = append(filterNames, "monthly_savings = "+data.MonthlySavings.ValueBigFloat().String())
    }
    if !data.MonthlyCost.IsNull() && !data.MonthlyCost.IsUnknown() {
        filters["monthlyCost"] = lookupNumber(data.MonthlyCost)
        filterNames = append(filterNames, "monthly_cost = "+data.MonthlyCost.ValueBigFloat().String())
    }
    if !data.Currency.IsNull() && !data.Currency.IsUnknown() {
        filters["currency"] = data.Currency.ValueString()
        filterNames = append(filterNames, "currency = "+fmt.Sprintf("%q", data.Currency.ValueString()))
    }
    if !data.Effort.IsNull() && !data.Effort.IsUnknown() {
        filters["effort"] = data.Effort.ValueString()
        filterNames = append(filterNames, "effort = "+fmt.Sprintf("%q", data.Effort.ValueString()))
    }
    if !data.Source.IsNull() && !data.Source.IsUnknown() {
        filters["source"] = data.Source.ValueString()
        filterNames = append(filterNames, "source = "+fmt.Sprintf("%q", data.Source.ValueString()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the cloud cost recommendation up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the cloud cost recommendation up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "cloudAccountId": true,
        "cloudResourceId": true,
        "provider": true,
        "recommendationId": true,
        "region": true,
        "resourceId": true,
        "resourceType": true,
        "resourceName": true,
        "kind": true,
        "action": true,
        "currentConfiguration": true,
        "recommendedConfiguration": true,
        "monthlySavings": true,
        "monthlyCost": true,
        "currency": true,
        "effort": true,
        "source": true,
        "lastSeenAt": true,
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/cloud-cost-recommendation/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read cloud_cost_recommendation, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No cloud cost recommendation found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read cloud_cost_recommendation: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/cloud-cost-recommendation/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list cloud_cost_recommendation, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list cloud_cost_recommendation: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No cloud cost recommendation matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one cloud cost recommendation matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for cloud_cost_recommendation.")
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
    if obj, ok := item["cloudAccountId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CloudAccountId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CloudAccountId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CloudAccountId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CloudAccountId = types.StringValue(string(jsonBytes))
        } else {
            data.CloudAccountId = types.StringNull()
        }
    } else if val, ok := item["cloudAccountId"].(string); ok {
        data.CloudAccountId = types.StringValue(val)
    } else {
        data.CloudAccountId = types.StringNull()
    }
    if obj, ok := item["cloudResourceId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CloudResourceId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CloudResourceId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CloudResourceId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CloudResourceId = types.StringValue(string(jsonBytes))
        } else {
            data.CloudResourceId = types.StringNull()
        }
    } else if val, ok := item["cloudResourceId"].(string); ok {
        data.CloudResourceId = types.StringValue(val)
    } else {
        data.CloudResourceId = types.StringNull()
    }
    if obj, ok := item["provider"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ProviderValue = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ProviderValue = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ProviderValue = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ProviderValue = types.StringValue(string(jsonBytes))
        } else {
            data.ProviderValue = types.StringNull()
        }
    } else if val, ok := item["provider"].(string); ok {
        data.ProviderValue = types.StringValue(val)
    } else {
        data.ProviderValue = types.StringNull()
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
    if obj, ok := item["region"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Region = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Region = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Region = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Region = types.StringValue(string(jsonBytes))
        } else {
            data.Region = types.StringNull()
        }
    } else if val, ok := item["region"].(string); ok {
        data.Region = types.StringValue(val)
    } else {
        data.Region = types.StringNull()
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
    if obj, ok := item["resourceName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ResourceName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ResourceName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ResourceName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ResourceName = types.StringValue(string(jsonBytes))
        } else {
            data.ResourceName = types.StringNull()
        }
    } else if val, ok := item["resourceName"].(string); ok {
        data.ResourceName = types.StringValue(val)
    } else {
        data.ResourceName = types.StringNull()
    }
    if obj, ok := item["kind"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Kind = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Kind = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Kind = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Kind = types.StringValue(string(jsonBytes))
        } else {
            data.Kind = types.StringNull()
        }
    } else if val, ok := item["kind"].(string); ok {
        data.Kind = types.StringValue(val)
    } else {
        data.Kind = types.StringNull()
    }
    if obj, ok := item["action"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Action = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Action = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Action = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Action = types.StringValue(string(jsonBytes))
        } else {
            data.Action = types.StringNull()
        }
    } else if val, ok := item["action"].(string); ok {
        data.Action = types.StringValue(val)
    } else {
        data.Action = types.StringNull()
    }
    if obj, ok := item["currentConfiguration"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CurrentConfiguration = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CurrentConfiguration = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CurrentConfiguration = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CurrentConfiguration = types.StringValue(string(jsonBytes))
        } else {
            data.CurrentConfiguration = types.StringNull()
        }
    } else if val, ok := item["currentConfiguration"].(string); ok {
        data.CurrentConfiguration = types.StringValue(val)
    } else {
        data.CurrentConfiguration = types.StringNull()
    }
    if obj, ok := item["recommendedConfiguration"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RecommendedConfiguration = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RecommendedConfiguration = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RecommendedConfiguration = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RecommendedConfiguration = types.StringValue(string(jsonBytes))
        } else {
            data.RecommendedConfiguration = types.StringNull()
        }
    } else if val, ok := item["recommendedConfiguration"].(string); ok {
        data.RecommendedConfiguration = types.StringValue(val)
    } else {
        data.RecommendedConfiguration = types.StringNull()
    }
    data.MonthlySavings = numberFromAPI(data.MonthlySavings, item["monthlySavings"])
    data.MonthlyCost = numberFromAPI(data.MonthlyCost, item["monthlyCost"])
    if obj, ok := item["currency"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Currency = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Currency = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Currency = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Currency = types.StringValue(string(jsonBytes))
        } else {
            data.Currency = types.StringNull()
        }
    } else if val, ok := item["currency"].(string); ok {
        data.Currency = types.StringValue(val)
    } else {
        data.Currency = types.StringNull()
    }
    if obj, ok := item["effort"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Effort = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Effort = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Effort = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Effort = types.StringValue(string(jsonBytes))
        } else {
            data.Effort = types.StringNull()
        }
    } else if val, ok := item["effort"].(string); ok {
        data.Effort = types.StringValue(val)
    } else {
        data.Effort = types.StringNull()
    }
    if obj, ok := item["source"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Source = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Source = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Source = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Source = types.StringValue(string(jsonBytes))
        } else {
            data.Source = types.StringNull()
        }
    } else if val, ok := item["source"].(string); ok {
        data.Source = types.StringValue(val)
    } else {
        data.Source = types.StringNull()
    }
    if obj, ok := item["lastSeenAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastSeenAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastSeenAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastSeenAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastSeenAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastSeenAt = types.StringNull()
        }
    } else if val, ok := item["lastSeenAt"].(string); ok {
        data.LastSeenAt = types.StringValue(val)
    } else {
        data.LastSeenAt = types.StringNull()
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
