package provider

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "math/big"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &MetricPipelineRuleDataSource{}

func NewMetricPipelineRuleDataSource() datasource.DataSource {
    return &MetricPipelineRuleDataSource{}
}

// MetricPipelineRuleDataSource defines the data source implementation.
type MetricPipelineRuleDataSource struct {
    client *Client
}

// MetricPipelineRuleDataSourceModel describes the data source data model.
type MetricPipelineRuleDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    ServiceId types.String `tfsdk:"service_id"`
    Name types.String `tfsdk:"name"`
    Description types.String `tfsdk:"description"`
    RuleType types.String `tfsdk:"rule_type"`
    FilterCondition types.String `tfsdk:"filter_condition"`
    Filters types.String `tfsdk:"filters"`
    RenameFromKey types.String `tfsdk:"rename_from_key"`
    RenameToKey types.String `tfsdk:"rename_to_key"`
    AddAttributeKey types.String `tfsdk:"add_attribute_key"`
    AddAttributeValue types.String `tfsdk:"add_attribute_value"`
    RedactReplacement types.String `tfsdk:"redact_replacement"`
    SamplePercentage types.Number `tfsdk:"sample_percentage"`
    IsEnabled types.Bool `tfsdk:"is_enabled"`
    SortOrder types.Number `tfsdk:"sort_order"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *MetricPipelineRuleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_metric_pipeline_rule"
}

func (d *MetricPipelineRuleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Rules applied at metric ingest time to filter, drop, rename, enrich, redact, or sample metric data points. Look up an existing metric pipeline rule by `id`, or by any of its other arguments (`name`, `add_attribute_key`, `add_attribute_value`, ...): each one set must match, and exactly one metric pipeline rule may match them all.",

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
                MarkdownDescription: "ID of the project this metric pipeline rule belongs to. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "service_id": schema.StringAttribute{
                MarkdownDescription: "Optional service ID scoping this rule. Null means the rule is project-wide. The ID of a `oneuptime_service`.",
                Optional: true,
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Friendly name for this rule.",
                Optional: true,
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Description of what this rule does.",
                Optional: true,
                Computed: true,
            },
            "rule_type": schema.StringAttribute{
                MarkdownDescription: "One of: Filter, Drop, RenameMetric, RenameAttribute, AddAttribute, RemoveAttribute, RedactAttribute, Sample.",
                Optional: true,
                Computed: true,
            },
            "filter_condition": schema.StringAttribute{
                MarkdownDescription: "How to combine filters: 'All' requires every filter to match (AND), 'Any' requires at least one to match (OR).",
                Optional: true,
                Computed: true,
            },
            "filters": schema.StringAttribute{
                MarkdownDescription: "List of filters evaluated against each metric data point. An empty list matches every data point. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "rename_from_key": schema.StringAttribute{
                MarkdownDescription: "For RenameMetric: the existing metric name. For RenameAttribute: the existing attribute key.",
                Optional: true,
                Computed: true,
            },
            "rename_to_key": schema.StringAttribute{
                MarkdownDescription: "For RenameMetric: the new metric name. For RenameAttribute: the new attribute key.",
                Optional: true,
                Computed: true,
            },
            "add_attribute_key": schema.StringAttribute{
                MarkdownDescription: "For AddAttribute / RemoveAttribute / RedactAttribute: the attribute key to act on.",
                Optional: true,
                Computed: true,
            },
            "add_attribute_value": schema.StringAttribute{
                MarkdownDescription: "For AddAttribute: the attribute value to set.",
                Optional: true,
                Computed: true,
            },
            "redact_replacement": schema.StringAttribute{
                MarkdownDescription: "For RedactAttribute: the literal string to replace the value with. Defaults to [REDACTED].",
                Optional: true,
                Computed: true,
            },
            "sample_percentage": schema.NumberAttribute{
                MarkdownDescription: "For Sample: percentage of matched rows to keep (0-100). 100 keeps all.",
                Optional: true,
                Computed: true,
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether this rule is active.",
                Optional: true,
                Computed: true,
            },
            "sort_order": schema.NumberAttribute{
                MarkdownDescription: "Where this rule is evaluated among the project's metric pipeline rules, lowest number first. A new rule is added to the end of the list. Setting a number another one already has puts it in that place, and the ones in the way move one place along to make room. In the dashboard, drag the rows to reorder them.",
                Optional: true,
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the user who created this rule. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *MetricPipelineRuleDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *MetricPipelineRuleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data MetricPipelineRuleDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.ServiceId.IsNull() && !data.ServiceId.IsUnknown() {
        filters["serviceId"] = data.ServiceId.ValueString()
        filterNames = append(filterNames, "service_id = "+fmt.Sprintf("%q", data.ServiceId.ValueString()))
    }
    if !data.Name.IsNull() && !data.Name.IsUnknown() {
        filters["name"] = data.Name.ValueString()
        filterNames = append(filterNames, "name = "+fmt.Sprintf("%q", data.Name.ValueString()))
    }
    if !data.Description.IsNull() && !data.Description.IsUnknown() {
        filters["description"] = data.Description.ValueString()
        filterNames = append(filterNames, "description = "+fmt.Sprintf("%q", data.Description.ValueString()))
    }
    if !data.RuleType.IsNull() && !data.RuleType.IsUnknown() {
        filters["ruleType"] = data.RuleType.ValueString()
        filterNames = append(filterNames, "rule_type = "+fmt.Sprintf("%q", data.RuleType.ValueString()))
    }
    if !data.FilterCondition.IsNull() && !data.FilterCondition.IsUnknown() {
        filters["filterCondition"] = data.FilterCondition.ValueString()
        filterNames = append(filterNames, "filter_condition = "+fmt.Sprintf("%q", data.FilterCondition.ValueString()))
    }
    if !data.RenameFromKey.IsNull() && !data.RenameFromKey.IsUnknown() {
        filters["renameFromKey"] = data.RenameFromKey.ValueString()
        filterNames = append(filterNames, "rename_from_key = "+fmt.Sprintf("%q", data.RenameFromKey.ValueString()))
    }
    if !data.RenameToKey.IsNull() && !data.RenameToKey.IsUnknown() {
        filters["renameToKey"] = data.RenameToKey.ValueString()
        filterNames = append(filterNames, "rename_to_key = "+fmt.Sprintf("%q", data.RenameToKey.ValueString()))
    }
    if !data.AddAttributeKey.IsNull() && !data.AddAttributeKey.IsUnknown() {
        filters["addAttributeKey"] = data.AddAttributeKey.ValueString()
        filterNames = append(filterNames, "add_attribute_key = "+fmt.Sprintf("%q", data.AddAttributeKey.ValueString()))
    }
    if !data.AddAttributeValue.IsNull() && !data.AddAttributeValue.IsUnknown() {
        filters["addAttributeValue"] = data.AddAttributeValue.ValueString()
        filterNames = append(filterNames, "add_attribute_value = "+fmt.Sprintf("%q", data.AddAttributeValue.ValueString()))
    }
    if !data.RedactReplacement.IsNull() && !data.RedactReplacement.IsUnknown() {
        filters["redactReplacement"] = data.RedactReplacement.ValueString()
        filterNames = append(filterNames, "redact_replacement = "+fmt.Sprintf("%q", data.RedactReplacement.ValueString()))
    }
    if !data.SamplePercentage.IsNull() && !data.SamplePercentage.IsUnknown() {
        filters["samplePercentage"] = lookupNumber(data.SamplePercentage)
        filterNames = append(filterNames, "sample_percentage = "+data.SamplePercentage.ValueBigFloat().String())
    }
    if !data.IsEnabled.IsNull() && !data.IsEnabled.IsUnknown() {
        filters["isEnabled"] = data.IsEnabled.ValueBool()
        filterNames = append(filterNames, "is_enabled = "+fmt.Sprintf("%t", data.IsEnabled.ValueBool()))
    }
    if !data.SortOrder.IsNull() && !data.SortOrder.IsUnknown() {
        filters["sortOrder"] = lookupNumber(data.SortOrder)
        filterNames = append(filterNames, "sort_order = "+data.SortOrder.ValueBigFloat().String())
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the metric pipeline rule up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the metric pipeline rule up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "serviceId": true,
        "name": true,
        "description": true,
        "ruleType": true,
        "filterCondition": true,
        "filters": true,
        "renameFromKey": true,
        "renameToKey": true,
        "addAttributeKey": true,
        "addAttributeValue": true,
        "redactReplacement": true,
        "samplePercentage": true,
        "isEnabled": true,
        "sortOrder": true,
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/metric-pipeline-rule/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read metric_pipeline_rule, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No metric pipeline rule found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read metric_pipeline_rule: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/metric-pipeline-rule/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list metric_pipeline_rule, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list metric_pipeline_rule: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No metric pipeline rule matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one metric pipeline rule matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for metric_pipeline_rule.")
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
    if obj, ok := item["serviceId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ServiceId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ServiceId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ServiceId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ServiceId = types.StringValue(string(jsonBytes))
        } else {
            data.ServiceId = types.StringNull()
        }
    } else if val, ok := item["serviceId"].(string); ok {
        data.ServiceId = types.StringValue(val)
    } else {
        data.ServiceId = types.StringNull()
    }
    if obj, ok := item["name"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Name = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Name = types.StringValue(string(jsonBytes))
        } else {
            data.Name = types.StringNull()
        }
    } else if val, ok := item["name"].(string); ok {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if obj, ok := item["description"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Description = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Description = types.StringValue(string(jsonBytes))
        } else {
            data.Description = types.StringNull()
        }
    } else if val, ok := item["description"].(string); ok {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if obj, ok := item["ruleType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RuleType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RuleType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RuleType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RuleType = types.StringValue(string(jsonBytes))
        } else {
            data.RuleType = types.StringNull()
        }
    } else if val, ok := item["ruleType"].(string); ok {
        data.RuleType = types.StringValue(val)
    } else {
        data.RuleType = types.StringNull()
    }
    if obj, ok := item["filterCondition"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FilterCondition = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.FilterCondition = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.FilterCondition = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.FilterCondition = types.StringValue(string(jsonBytes))
        } else {
            data.FilterCondition = types.StringNull()
        }
    } else if val, ok := item["filterCondition"].(string); ok {
        data.FilterCondition = types.StringValue(val)
    } else {
        data.FilterCondition = types.StringNull()
    }
    if obj, ok := item["filters"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Filters = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Filters = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Filters = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Filters = types.StringValue(string(jsonBytes))
        } else {
            data.Filters = types.StringNull()
        }
    } else if val, ok := item["filters"].(string); ok {
        data.Filters = types.StringValue(val)
    } else {
        data.Filters = types.StringNull()
    }
    if obj, ok := item["renameFromKey"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RenameFromKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RenameFromKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RenameFromKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RenameFromKey = types.StringValue(string(jsonBytes))
        } else {
            data.RenameFromKey = types.StringNull()
        }
    } else if val, ok := item["renameFromKey"].(string); ok {
        data.RenameFromKey = types.StringValue(val)
    } else {
        data.RenameFromKey = types.StringNull()
    }
    if obj, ok := item["renameToKey"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RenameToKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RenameToKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RenameToKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RenameToKey = types.StringValue(string(jsonBytes))
        } else {
            data.RenameToKey = types.StringNull()
        }
    } else if val, ok := item["renameToKey"].(string); ok {
        data.RenameToKey = types.StringValue(val)
    } else {
        data.RenameToKey = types.StringNull()
    }
    if obj, ok := item["addAttributeKey"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AddAttributeKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AddAttributeKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AddAttributeKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AddAttributeKey = types.StringValue(string(jsonBytes))
        } else {
            data.AddAttributeKey = types.StringNull()
        }
    } else if val, ok := item["addAttributeKey"].(string); ok {
        data.AddAttributeKey = types.StringValue(val)
    } else {
        data.AddAttributeKey = types.StringNull()
    }
    if obj, ok := item["addAttributeValue"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AddAttributeValue = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AddAttributeValue = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AddAttributeValue = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AddAttributeValue = types.StringValue(string(jsonBytes))
        } else {
            data.AddAttributeValue = types.StringNull()
        }
    } else if val, ok := item["addAttributeValue"].(string); ok {
        data.AddAttributeValue = types.StringValue(val)
    } else {
        data.AddAttributeValue = types.StringNull()
    }
    if obj, ok := item["redactReplacement"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RedactReplacement = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RedactReplacement = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RedactReplacement = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RedactReplacement = types.StringValue(string(jsonBytes))
        } else {
            data.RedactReplacement = types.StringNull()
        }
    } else if val, ok := item["redactReplacement"].(string); ok {
        data.RedactReplacement = types.StringValue(val)
    } else {
        data.RedactReplacement = types.StringNull()
    }
    if val, ok := item["samplePercentage"].(float64); ok {
        data.SamplePercentage = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["samplePercentage"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.SamplePercentage = types.NumberValue(big.NewFloat(val))
        } else {
            data.SamplePercentage = types.NumberNull()
        }
    } else {
        data.SamplePercentage = types.NumberNull()
    }
    if val, ok := item["isEnabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    } else {
        data.IsEnabled = types.BoolNull()
    }
    if val, ok := item["sortOrder"].(float64); ok {
        data.SortOrder = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["sortOrder"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.SortOrder = types.NumberValue(big.NewFloat(val))
        } else {
            data.SortOrder = types.NumberNull()
        }
    } else {
        data.SortOrder = types.NumberNull()
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
