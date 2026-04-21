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
var _ datasource.DataSource = &MetricPipelineRuleDataDataSource{}

func NewMetricPipelineRuleDataDataSource() datasource.DataSource {
    return &MetricPipelineRuleDataDataSource{}
}

// MetricPipelineRuleDataDataSource defines the data source implementation.
type MetricPipelineRuleDataDataSource struct {
    client *Client
}

// MetricPipelineRuleDataDataSourceModel describes the data source data model.
type MetricPipelineRuleDataDataSourceModel struct {
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

func (d *MetricPipelineRuleDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_metric_pipeline_rule_data"
}

func (d *MetricPipelineRuleDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "metric_pipeline_rule_data data source",

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
                MarkdownDescription: "Description of what this rule does.. Permissions - Create: [Project Owner, Project Admin, Create Metric Pipeline Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Telemetry Manager, Read Metric Pipeline Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Metric Pipeline Rule]",
                Computed: true,
            },
            "rule_type": schema.StringAttribute{
                MarkdownDescription: "One of: Filter, Drop, RenameMetric, RenameAttribute, AddAttribute, RemoveAttribute, RedactAttribute, Sample.. Permissions - Create: [Project Owner, Project Admin, Create Metric Pipeline Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Telemetry Manager, Read Metric Pipeline Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Metric Pipeline Rule]",
                Computed: true,
            },
            "filter_condition": schema.StringAttribute{
                MarkdownDescription: "How to combine filters: 'All' requires every filter to match (AND), 'Any' requires at least one to match (OR).. Permissions - Create: [Project Owner, Project Admin, Create Metric Pipeline Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Telemetry Manager, Read Metric Pipeline Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Metric Pipeline Rule]",
                Computed: true,
            },
            "filters": schema.StringAttribute{
                MarkdownDescription: "List of filters evaluated against each metric data point. An empty list matches every data point.. Permissions - Create: [Project Owner, Project Admin, Create Metric Pipeline Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Telemetry Manager, Read Metric Pipeline Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Metric Pipeline Rule]",
                Computed: true,
            },
            "rename_from_key": schema.StringAttribute{
                MarkdownDescription: "For RenameMetric: the existing metric name. For RenameAttribute: the existing attribute key.. Permissions - Create: [Project Owner, Project Admin, Create Metric Pipeline Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Telemetry Manager, Read Metric Pipeline Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Metric Pipeline Rule]",
                Computed: true,
            },
            "rename_to_key": schema.StringAttribute{
                MarkdownDescription: "For RenameMetric: the new metric name. For RenameAttribute: the new attribute key.. Permissions - Create: [Project Owner, Project Admin, Create Metric Pipeline Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Telemetry Manager, Read Metric Pipeline Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Metric Pipeline Rule]",
                Computed: true,
            },
            "add_attribute_key": schema.StringAttribute{
                MarkdownDescription: "For AddAttribute / RemoveAttribute / RedactAttribute: the attribute key to act on.. Permissions - Create: [Project Owner, Project Admin, Create Metric Pipeline Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Telemetry Manager, Read Metric Pipeline Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Metric Pipeline Rule]",
                Computed: true,
            },
            "add_attribute_value": schema.StringAttribute{
                MarkdownDescription: "For AddAttribute: the attribute value to set.. Permissions - Create: [Project Owner, Project Admin, Create Metric Pipeline Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Telemetry Manager, Read Metric Pipeline Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Metric Pipeline Rule]",
                Computed: true,
            },
            "redact_replacement": schema.StringAttribute{
                MarkdownDescription: "For RedactAttribute: the literal string to replace the value with. Defaults to [REDACTED].. Permissions - Create: [Project Owner, Project Admin, Create Metric Pipeline Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Telemetry Manager, Read Metric Pipeline Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Metric Pipeline Rule]",
                Computed: true,
            },
            "sample_percentage": schema.NumberAttribute{
                MarkdownDescription: "For Sample: percentage of matched rows to keep (0-100). 100 keeps all.. Permissions - Create: [Project Owner, Project Admin, Create Metric Pipeline Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Telemetry Manager, Read Metric Pipeline Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Metric Pipeline Rule]",
                Computed: true,
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether this rule is active.. Permissions - Create: [Project Owner, Project Admin, Create Metric Pipeline Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Telemetry Manager, Read Metric Pipeline Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Metric Pipeline Rule]",
                Computed: true,
            },
            "sort_order": schema.NumberAttribute{
                MarkdownDescription: "Evaluation order within its scope (service-level or project-level).. Permissions - Create: [Project Owner, Project Admin, Create Metric Pipeline Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Telemetry Manager, Read Metric Pipeline Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Metric Pipeline Rule]",
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

func (d *MetricPipelineRuleDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *MetricPipelineRuleDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data MetricPipelineRuleDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "metric-pipeline-rule" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read metric_pipeline_rule_data, got error: %s", err))
        return
    }

    var metricPipelineRuleDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &metricPipelineRuleDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse metric_pipeline_rule_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := metricPipelineRuleDataResponse["data"].(map[string]interface{}); ok {
        metricPipelineRuleDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := metricPipelineRuleDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := metricPipelineRuleDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := metricPipelineRuleDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := metricPipelineRuleDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := metricPipelineRuleDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := metricPipelineRuleDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := metricPipelineRuleDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := metricPipelineRuleDataResponse["service_id"].(string); ok {
        data.ServiceId = types.StringValue(val)
    }
    if val, ok := metricPipelineRuleDataResponse["description"].(string); ok {
        data.Description = types.StringValue(val)
    }
    if val, ok := metricPipelineRuleDataResponse["rule_type"].(string); ok {
        data.RuleType = types.StringValue(val)
    }
    if val, ok := metricPipelineRuleDataResponse["filter_condition"].(string); ok {
        data.FilterCondition = types.StringValue(val)
    }
    if val, ok := metricPipelineRuleDataResponse["filters"].(string); ok {
        data.Filters = types.StringValue(val)
    }
    if val, ok := metricPipelineRuleDataResponse["rename_from_key"].(string); ok {
        data.RenameFromKey = types.StringValue(val)
    }
    if val, ok := metricPipelineRuleDataResponse["rename_to_key"].(string); ok {
        data.RenameToKey = types.StringValue(val)
    }
    if val, ok := metricPipelineRuleDataResponse["add_attribute_key"].(string); ok {
        data.AddAttributeKey = types.StringValue(val)
    }
    if val, ok := metricPipelineRuleDataResponse["add_attribute_value"].(string); ok {
        data.AddAttributeValue = types.StringValue(val)
    }
    if val, ok := metricPipelineRuleDataResponse["redact_replacement"].(string); ok {
        data.RedactReplacement = types.StringValue(val)
    }
    if val, ok := metricPipelineRuleDataResponse["sample_percentage"].(float64); ok {
        data.SamplePercentage = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := metricPipelineRuleDataResponse["is_enabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    }
    if val, ok := metricPipelineRuleDataResponse["sort_order"].(float64); ok {
        data.SortOrder = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := metricPipelineRuleDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := metricPipelineRuleDataResponse["deleted_by_user_id"].(string); ok {
        data.DeletedByUserId = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
