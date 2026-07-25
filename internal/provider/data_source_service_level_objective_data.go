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
var _ datasource.DataSource = &ServiceLevelObjectiveDataDataSource{}

func NewServiceLevelObjectiveDataDataSource() datasource.DataSource {
    return &ServiceLevelObjectiveDataDataSource{}
}

// ServiceLevelObjectiveDataDataSource defines the data source implementation.
type ServiceLevelObjectiveDataDataSource struct {
    client *Client
}

// ServiceLevelObjectiveDataDataSourceModel describes the data source data model.
type ServiceLevelObjectiveDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    Description types.String `tfsdk:"description"`
    Slug types.String `tfsdk:"slug"`
    Labels types.Set `tfsdk:"labels"`
    IsEnabled types.Bool `tfsdk:"is_enabled"`
    SliType types.String `tfsdk:"sli_type"`
    MultiMonitorMode types.String `tfsdk:"multi_monitor_mode"`
    Monitors types.Set `tfsdk:"monitors"`
    DowntimeMonitorStatuses types.Set `tfsdk:"downtime_monitor_statuses"`
    MetricQueryConfig types.String `tfsdk:"metric_query_config"`
    TargetPercentage types.Number `tfsdk:"target_percentage"`
    WindowType types.String `tfsdk:"window_type"`
    WindowDays types.Number `tfsdk:"window_days"`
    Timezone types.String `tfsdk:"timezone"`
    AtRiskThresholdPercentage types.Number `tfsdk:"at_risk_threshold_percentage"`
    CurrentSliPercentage types.Number `tfsdk:"current_sli_percentage"`
    ErrorBudgetRemainingPercentage types.Number `tfsdk:"error_budget_remaining_percentage"`
    ErrorBudgetRemainingSeconds types.Number `tfsdk:"error_budget_remaining_seconds"`
    ErrorBudgetTotalSeconds types.Number `tfsdk:"error_budget_total_seconds"`
    CurrentBurnRate types.Number `tfsdk:"current_burn_rate"`
    SloStatus types.String `tfsdk:"slo_status"`
    StatusChangeNotificationSentAt types.String `tfsdk:"status_change_notification_sent_at"`
    LastEvaluatedAt types.String `tfsdk:"last_evaluated_at"`
    NextEvaluationAt types.String `tfsdk:"next_evaluation_at"`
    LastAccumulatedBucketEndAt types.String `tfsdk:"last_accumulated_bucket_end_at"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *ServiceLevelObjectiveDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_service_level_objective_data"
}

func (d *ServiceLevelObjectiveDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "service_level_objective_data data source",

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
            "description": schema.StringAttribute{
                MarkdownDescription: "Description of this Service Level Objective. Permissions - Create: [Project Owner, Project Admin, Create Service Level Objective], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Service Level Objective], Update: [Project Owner, Project Admin, Edit Service Level Objective]",
                Computed: true,
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Friendly globally unique name for your object. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Service Level Objective], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "labels": schema.SetAttribute{
                MarkdownDescription: "Relation to Labels Array where this object is categorized in.. Permissions - Create: [Project Owner, Project Admin, Create Service Level Objective], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Service Level Objective], Update: [Project Owner, Project Admin, Edit Service Level Objective]",
                Computed: true,
                ElementType: types.StringType,
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether this Service Level Objective is enabled. Disabled SLOs are not evaluated.. Permissions - Create: [Project Owner, Project Admin, Create Service Level Objective], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Service Level Objective], Update: [Project Owner, Project Admin, Edit Service Level Objective]",
                Computed: true,
            },
            "sli_type": schema.StringAttribute{
                MarkdownDescription: "Type of Service Level Indicator this objective measures (Monitor Uptime or Metric). Permissions - Create: [Project Owner, Project Admin, Create Service Level Objective], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Service Level Objective], Update: [Project Owner, Project Admin, Edit Service Level Objective]",
                Computed: true,
            },
            "multi_monitor_mode": schema.StringAttribute{
                MarkdownDescription: "How downtime is counted when multiple monitors are attached. 'Any Monitor Down' counts time when any monitor is down. 'Monitor Seconds Average' averages downtime across monitors.. Permissions - Create: [Project Owner, Project Admin, Create Service Level Objective], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Service Level Objective], Update: [Project Owner, Project Admin, Edit Service Level Objective]",
                Computed: true,
            },
            "monitors": schema.SetAttribute{
                MarkdownDescription: "Monitors whose uptime is measured by this Service Level Objective (for Monitor Uptime SLIs).. Permissions - Create: [Project Owner, Project Admin, Create Service Level Objective], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Service Level Objective], Update: [Project Owner, Project Admin, Edit Service Level Objective]",
                Computed: true,
                ElementType: types.StringType,
            },
            "downtime_monitor_statuses": schema.SetAttribute{
                MarkdownDescription: "List of monitor statuses that are considered as \"down\" for this Service Level Objective.. Permissions - Create: [Project Owner, Project Admin, Create Service Level Objective], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Service Level Objective], Update: [Project Owner, Project Admin, Edit Service Level Objective]",
                Computed: true,
                ElementType: types.StringType,
            },
            "metric_query_config": schema.StringAttribute{
                MarkdownDescription: "Query configuration for Metric SLIs: metric name, good-event predicate and optional attribute filters.. Permissions - Create: [Project Owner, Project Admin, Create Service Level Objective], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Service Level Objective], Update: [Project Owner, Project Admin, Edit Service Level Objective]",
                Computed: true,
            },
            "target_percentage": schema.NumberAttribute{
                MarkdownDescription: "Target of this Service Level Objective as a percentage (e.g. 99.9). Must be less than 100.. Permissions - Create: [Project Owner, Project Admin, Create Service Level Objective], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Service Level Objective], Update: [Project Owner, Project Admin, Edit Service Level Objective]",
                Computed: true,
            },
            "window_type": schema.StringAttribute{
                MarkdownDescription: "Type of compliance window for this objective (Rolling or Calendar Month). Permissions - Create: [Project Owner, Project Admin, Create Service Level Objective], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Service Level Objective], Update: [Project Owner, Project Admin, Edit Service Level Objective]",
                Computed: true,
            },
            "window_days": schema.NumberAttribute{
                MarkdownDescription: "Length of the rolling compliance window in days (e.g. 7, 28, 30 or 90). Ignored for Calendar Month windows.. Permissions - Create: [Project Owner, Project Admin, Create Service Level Objective], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Service Level Objective], Update: [Project Owner, Project Admin, Edit Service Level Objective]",
                Computed: true,
            },
            "timezone": schema.StringAttribute{
                MarkdownDescription: "IANA timezone (e.g. America/New_York) used for Calendar Month window boundaries. Defaults to UTC when not set.. Permissions - Create: [Project Owner, Project Admin, Create Service Level Objective], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Service Level Objective], Update: [Project Owner, Project Admin, Edit Service Level Objective]",
                Computed: true,
            },
            "at_risk_threshold_percentage": schema.NumberAttribute{
                MarkdownDescription: "Percentage of remaining error budget at which the SLO status changes to At Risk. For example, 20 means the status becomes At Risk when less than 20% of the error budget remains.. Permissions - Create: [Project Owner, Project Admin, Create Service Level Objective], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Service Level Objective], Update: [Project Owner, Project Admin, Edit Service Level Objective]",
                Computed: true,
            },
            "current_sli_percentage": schema.NumberAttribute{
                MarkdownDescription: "Current Service Level Indicator over the compliance window, as a percentage. Computed by the worker.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Service Level Objective], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "error_budget_remaining_percentage": schema.NumberAttribute{
                MarkdownDescription: "Percentage of the error budget that remains. Can be negative when the budget is exhausted. Computed by the worker.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Service Level Objective], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "error_budget_remaining_seconds": schema.NumberAttribute{
                MarkdownDescription: "Seconds of error budget that remain. Can be negative when the budget is exhausted. Computed by the worker.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Service Level Objective], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "error_budget_total_seconds": schema.NumberAttribute{
                MarkdownDescription: "Total seconds of error budget for the compliance window. Computed by the worker.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Service Level Objective], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "current_burn_rate": schema.NumberAttribute{
                MarkdownDescription: "Rate at which the error budget is currently being consumed. A burn rate of 1 exhausts the budget exactly at the end of the window. Computed by the worker.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Service Level Objective], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "slo_status": schema.StringAttribute{
                MarkdownDescription: "Current status of this Service Level Objective (Healthy, At Risk, Budget Exhausted, Misconfigured, Paused). Computed by the worker.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Service Level Objective], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "status_change_notification_sent_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "last_evaluated_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "next_evaluation_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "last_accumulated_bucket_end_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
        },
    }
}

func (d *ServiceLevelObjectiveDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ServiceLevelObjectiveDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data ServiceLevelObjectiveDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "service-level-objective" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read service_level_objective_data, got error: %s", err))
        return
    }

    var serviceLevelObjectiveDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &serviceLevelObjectiveDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse service_level_objective_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := serviceLevelObjectiveDataResponse["data"].(map[string]interface{}); ok {
        serviceLevelObjectiveDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := serviceLevelObjectiveDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := serviceLevelObjectiveDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := serviceLevelObjectiveDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := serviceLevelObjectiveDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := serviceLevelObjectiveDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := serviceLevelObjectiveDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := serviceLevelObjectiveDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := serviceLevelObjectiveDataResponse["description"].(string); ok {
        data.Description = types.StringValue(val)
    }
    if val, ok := serviceLevelObjectiveDataResponse["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    }
    if val, ok := serviceLevelObjectiveDataResponse["labels"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.Labels = setValue
    }
    if val, ok := serviceLevelObjectiveDataResponse["is_enabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    }
    if val, ok := serviceLevelObjectiveDataResponse["sli_type"].(string); ok {
        data.SliType = types.StringValue(val)
    }
    if val, ok := serviceLevelObjectiveDataResponse["multi_monitor_mode"].(string); ok {
        data.MultiMonitorMode = types.StringValue(val)
    }
    if val, ok := serviceLevelObjectiveDataResponse["monitors"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.Monitors = setValue
    }
    if val, ok := serviceLevelObjectiveDataResponse["downtime_monitor_statuses"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.DowntimeMonitorStatuses = setValue
    }
    if val, ok := serviceLevelObjectiveDataResponse["metric_query_config"].(string); ok {
        data.MetricQueryConfig = types.StringValue(val)
    }
    if val, ok := serviceLevelObjectiveDataResponse["target_percentage"].(float64); ok {
        data.TargetPercentage = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := serviceLevelObjectiveDataResponse["window_type"].(string); ok {
        data.WindowType = types.StringValue(val)
    }
    if val, ok := serviceLevelObjectiveDataResponse["window_days"].(float64); ok {
        data.WindowDays = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := serviceLevelObjectiveDataResponse["timezone"].(string); ok {
        data.Timezone = types.StringValue(val)
    }
    if val, ok := serviceLevelObjectiveDataResponse["at_risk_threshold_percentage"].(float64); ok {
        data.AtRiskThresholdPercentage = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := serviceLevelObjectiveDataResponse["current_sli_percentage"].(float64); ok {
        data.CurrentSliPercentage = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := serviceLevelObjectiveDataResponse["error_budget_remaining_percentage"].(float64); ok {
        data.ErrorBudgetRemainingPercentage = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := serviceLevelObjectiveDataResponse["error_budget_remaining_seconds"].(float64); ok {
        data.ErrorBudgetRemainingSeconds = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := serviceLevelObjectiveDataResponse["error_budget_total_seconds"].(float64); ok {
        data.ErrorBudgetTotalSeconds = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := serviceLevelObjectiveDataResponse["current_burn_rate"].(float64); ok {
        data.CurrentBurnRate = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := serviceLevelObjectiveDataResponse["slo_status"].(string); ok {
        data.SloStatus = types.StringValue(val)
    }
    if val, ok := serviceLevelObjectiveDataResponse["status_change_notification_sent_at"].(string); ok {
        data.StatusChangeNotificationSentAt = types.StringValue(val)
    }
    if val, ok := serviceLevelObjectiveDataResponse["last_evaluated_at"].(string); ok {
        data.LastEvaluatedAt = types.StringValue(val)
    }
    if val, ok := serviceLevelObjectiveDataResponse["next_evaluation_at"].(string); ok {
        data.NextEvaluationAt = types.StringValue(val)
    }
    if val, ok := serviceLevelObjectiveDataResponse["last_accumulated_bucket_end_at"].(string); ok {
        data.LastAccumulatedBucketEndAt = types.StringValue(val)
    }
    if val, ok := serviceLevelObjectiveDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
