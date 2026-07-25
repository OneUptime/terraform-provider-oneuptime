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
var _ datasource.DataSource = &SloBurnRateRuleDataDataSource{}

func NewSloBurnRateRuleDataDataSource() datasource.DataSource {
    return &SloBurnRateRuleDataDataSource{}
}

// SloBurnRateRuleDataDataSource defines the data source implementation.
type SloBurnRateRuleDataDataSource struct {
    client *Client
}

// SloBurnRateRuleDataDataSourceModel describes the data source data model.
type SloBurnRateRuleDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    ServiceLevelObjectiveId types.String `tfsdk:"service_level_objective_id"`
    IsEnabled types.Bool `tfsdk:"is_enabled"`
    BurnRateThreshold types.Number `tfsdk:"burn_rate_threshold"`
    LongWindowInMinutes types.Number `tfsdk:"long_window_in_minutes"`
    ShortWindowInMinutes types.Number `tfsdk:"short_window_in_minutes"`
    MinimumSampleCount types.Number `tfsdk:"minimum_sample_count"`
    RefireSuppressionMinutes types.Number `tfsdk:"refire_suppression_minutes"`
    AlertSeverityId types.String `tfsdk:"alert_severity_id"`
    OnCallDutyPolicies types.Set `tfsdk:"on_call_duty_policies"`
    LastAlertCreatedAt types.String `tfsdk:"last_alert_created_at"`
    LastAlertResolvedAt types.String `tfsdk:"last_alert_resolved_at"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *SloBurnRateRuleDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_slo_burn_rate_rule_data"
}

func (d *SloBurnRateRuleDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "slo_burn_rate_rule_data data source",

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
            "service_level_objective_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether this burn rate rule is enabled. Permissions - Create: [Project Owner, Project Admin, Create SLO Burn Rate Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Read SLO Burn Rate Rule], Update: [Project Owner, Project Admin, Edit SLO Burn Rate Rule]",
                Computed: true,
            },
            "burn_rate_threshold": schema.NumberAttribute{
                MarkdownDescription: "Alert when the burn rate in both the long and short windows is at or above this threshold (e.g. 14.4).. Permissions - Create: [Project Owner, Project Admin, Create SLO Burn Rate Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Read SLO Burn Rate Rule], Update: [Project Owner, Project Admin, Edit SLO Burn Rate Rule]",
                Computed: true,
            },
            "long_window_in_minutes": schema.NumberAttribute{
                MarkdownDescription: "Length of the long lookback window in minutes (e.g. 60). The alert fires when both windows exceed the threshold and resolves when the long window drops below it.. Permissions - Create: [Project Owner, Project Admin, Create SLO Burn Rate Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Read SLO Burn Rate Rule], Update: [Project Owner, Project Admin, Edit SLO Burn Rate Rule]",
                Computed: true,
            },
            "short_window_in_minutes": schema.NumberAttribute{
                MarkdownDescription: "Length of the short lookback window in minutes (e.g. 5). Guards against alerting on burn that has already stopped.. Permissions - Create: [Project Owner, Project Admin, Create SLO Burn Rate Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Read SLO Burn Rate Rule], Update: [Project Owner, Project Admin, Edit SLO Burn Rate Rule]",
                Computed: true,
            },
            "minimum_sample_count": schema.NumberAttribute{
                MarkdownDescription: "For event-based SLIs only: skip this rule when the long window has fewer than this many total events. Prevents noisy alerts on low traffic.. Permissions - Create: [Project Owner, Project Admin, Create SLO Burn Rate Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Read SLO Burn Rate Rule], Update: [Project Owner, Project Admin, Edit SLO Burn Rate Rule]",
                Computed: true,
            },
            "refire_suppression_minutes": schema.NumberAttribute{
                MarkdownDescription: "Minimum number of minutes after an alert resolves before this rule can fire again. Defaults to the long window length when not set.. Permissions - Create: [Project Owner, Project Admin, Create SLO Burn Rate Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Read SLO Burn Rate Rule], Update: [Project Owner, Project Admin, Edit SLO Burn Rate Rule]",
                Computed: true,
            },
            "alert_severity_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "on_call_duty_policies": schema.SetAttribute{
                MarkdownDescription: "On-call duty policies attached to alerts created by this burn rate rule.. Permissions - Create: [Project Owner, Project Admin, Create SLO Burn Rate Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Read SLO Burn Rate Rule], Update: [Project Owner, Project Admin, Edit SLO Burn Rate Rule]",
                Computed: true,
                ElementType: types.StringType,
            },
            "last_alert_created_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "last_alert_resolved_at": schema.StringAttribute{
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

func (d *SloBurnRateRuleDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SloBurnRateRuleDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data SloBurnRateRuleDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "service-level-objective-burn-rate-rule" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read slo_burn_rate_rule_data, got error: %s", err))
        return
    }

    var sloBurnRateRuleDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &sloBurnRateRuleDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse slo_burn_rate_rule_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := sloBurnRateRuleDataResponse["data"].(map[string]interface{}); ok {
        sloBurnRateRuleDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := sloBurnRateRuleDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := sloBurnRateRuleDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := sloBurnRateRuleDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := sloBurnRateRuleDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := sloBurnRateRuleDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := sloBurnRateRuleDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := sloBurnRateRuleDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := sloBurnRateRuleDataResponse["service_level_objective_id"].(string); ok {
        data.ServiceLevelObjectiveId = types.StringValue(val)
    }
    if val, ok := sloBurnRateRuleDataResponse["is_enabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    }
    if val, ok := sloBurnRateRuleDataResponse["burn_rate_threshold"].(float64); ok {
        data.BurnRateThreshold = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := sloBurnRateRuleDataResponse["long_window_in_minutes"].(float64); ok {
        data.LongWindowInMinutes = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := sloBurnRateRuleDataResponse["short_window_in_minutes"].(float64); ok {
        data.ShortWindowInMinutes = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := sloBurnRateRuleDataResponse["minimum_sample_count"].(float64); ok {
        data.MinimumSampleCount = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := sloBurnRateRuleDataResponse["refire_suppression_minutes"].(float64); ok {
        data.RefireSuppressionMinutes = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := sloBurnRateRuleDataResponse["alert_severity_id"].(string); ok {
        data.AlertSeverityId = types.StringValue(val)
    }
    if val, ok := sloBurnRateRuleDataResponse["on_call_duty_policies"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.OnCallDutyPolicies = setValue
    }
    if val, ok := sloBurnRateRuleDataResponse["last_alert_created_at"].(string); ok {
        data.LastAlertCreatedAt = types.StringValue(val)
    }
    if val, ok := sloBurnRateRuleDataResponse["last_alert_resolved_at"].(string); ok {
        data.LastAlertResolvedAt = types.StringValue(val)
    }
    if val, ok := sloBurnRateRuleDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
