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
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    ServiceId types.String `tfsdk:"service_id"`
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
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
}

func (d *MetricPipelineRuleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_metric_pipeline_rule"
}

func (d *MetricPipelineRuleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Rules applied at metric ingest time to filter, drop, rename, enrich, redact, or sample metric data points. Look up an existing metric_pipeline_rule by `id` or by `name`.",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Look up by unique identifier. Exactly one of `id` or `name` must be set.",
                Optional: true,
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.",
                Optional: true,
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
            "service_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Description of what this rule does..",
                Computed: true,
            },
            "rule_type": schema.StringAttribute{
                MarkdownDescription: "One of: Filter, Drop, RenameMetric, RenameAttribute, AddAttribute, RemoveAttribute, RedactAttribute, Sample..",
                Computed: true,
            },
            "filter_condition": schema.StringAttribute{
                MarkdownDescription: "How to combine filters: 'All' requires every filter to match (AND), 'Any' requires at least one to match (OR)..",
                Computed: true,
            },
            "filters": schema.StringAttribute{
                MarkdownDescription: "List of filters evaluated against each metric data point. An empty list matches every data point..",
                Computed: true,
            },
            "rename_from_key": schema.StringAttribute{
                MarkdownDescription: "For RenameMetric: the existing metric name. For RenameAttribute: the existing attribute key..",
                Computed: true,
            },
            "rename_to_key": schema.StringAttribute{
                MarkdownDescription: "For RenameMetric: the new metric name. For RenameAttribute: the new attribute key..",
                Computed: true,
            },
            "add_attribute_key": schema.StringAttribute{
                MarkdownDescription: "For AddAttribute / RemoveAttribute / RedactAttribute: the attribute key to act on..",
                Computed: true,
            },
            "add_attribute_value": schema.StringAttribute{
                MarkdownDescription: "For AddAttribute: the attribute value to set..",
                Computed: true,
            },
            "redact_replacement": schema.StringAttribute{
                MarkdownDescription: "For RedactAttribute: the literal string to replace the value with. Defaults to [REDACTED]..",
                Computed: true,
            },
            "sample_percentage": schema.NumberAttribute{
                MarkdownDescription: "For Sample: percentage of matched rows to keep (0-100). 100 keeps all..",
                Computed: true,
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether this rule is active..",
                Computed: true,
            },
            "sort_order": schema.NumberAttribute{
                MarkdownDescription: "Evaluation order within its scope (service-level or project-level)..",
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

    hasId := !data.Id.IsNull() && data.Id.ValueString() != ""
    hasName := !data.Name.IsNull() && data.Name.ValueString() != ""
    if hasId == hasName {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Exactly one of `id` or `name` must be set to look up a metric_pipeline_rule.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "name": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "projectId": true,
        "serviceId": true,
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
        "deletedByUserId": true,
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
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No metric_pipeline_rule found with id %q.", data.Id.ValueString()))
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
    } else {
        listBody := map[string]interface{}{
            "query": map[string]interface{}{
                "name": data.Name.ValueString(),
            },
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
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No metric_pipeline_rule found with name %q.", data.Name.ValueString()))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one metric_pipeline_rule matches name %q. Use the id attribute to disambiguate.", data.Name.ValueString()))
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
    if obj, ok := item["deletedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DeletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DeletedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DeletedAt = types.StringValue(string(jsonBytes))
        } else {
            data.DeletedAt = types.StringNull()
        }
    } else if val, ok := item["deletedAt"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    } else {
        data.DeletedAt = types.StringNull()
    }
    if val, ok := item["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["version"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.Version = types.NumberValue(big.NewFloat(val))
        } else {
            data.Version = types.NumberNull()
        }
    } else {
        data.Version = types.NumberNull()
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
    if obj, ok := item["deletedByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeletedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DeletedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DeletedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DeletedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.DeletedByUserId = types.StringNull()
        }
    } else if val, ok := item["deletedByUserId"].(string); ok {
        data.DeletedByUserId = types.StringValue(val)
    } else {
        data.DeletedByUserId = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
