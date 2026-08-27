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
var _ datasource.DataSource = &DetectionRuleDataSource{}

func NewDetectionRuleDataSource() datasource.DataSource {
    return &DetectionRuleDataSource{}
}

// DetectionRuleDataSource defines the data source implementation.
type DetectionRuleDataSource struct {
    client *Client
}

// DetectionRuleDataSourceModel describes the data source data model.
type DetectionRuleDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    Description types.String `tfsdk:"description"`
    SigmaRuleYaml types.String `tfsdk:"sigma_rule_yaml"`
    IsEnabled types.Bool `tfsdk:"is_enabled"`
    EvaluationIntervalInMinutes types.Number `tfsdk:"evaluation_interval_in_minutes"`
    GroupByField types.String `tfsdk:"group_by_field"`
    DistinctCountField types.String `tfsdk:"distinct_count_field"`
    MatchCountThreshold types.Number `tfsdk:"match_count_threshold"`
    ShouldCreateAlert types.Bool `tfsdk:"should_create_alert"`
    ShouldWriteDetectionFinding types.Bool `tfsdk:"should_write_detection_finding"`
    ShouldCreateIncident types.Bool `tfsdk:"should_create_incident"`
    AlertSeverityId types.String `tfsdk:"alert_severity_id"`
    IncidentSeverityId types.String `tfsdk:"incident_severity_id"`
    LastEvaluatedAt types.String `tfsdk:"last_evaluated_at"`
    LastMatchAt types.String `tfsdk:"last_match_at"`
    LastError types.String `tfsdk:"last_error"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
}

func (d *DetectionRuleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_detection_rule"
}

func (d *DetectionRuleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Sigma detection rules evaluated against security events. Matches create alerts and detection findings. Look up an existing detection_rule by `id` or by `name`.",

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
            "description": schema.StringAttribute{
                MarkdownDescription: "Description of what this detection rule looks for..",
                Computed: true,
            },
            "sigma_rule_yaml": schema.StringAttribute{
                MarkdownDescription: "The Sigma rule to evaluate, in YAML. detection selections and condition are compiled to a ClickHouse query over security events..",
                Computed: true,
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether this detection rule is evaluated..",
                Computed: true,
            },
            "evaluation_interval_in_minutes": schema.NumberAttribute{
                MarkdownDescription: "How often this rule is evaluated, in minutes. The evaluation window covers the time since the previous evaluation..",
                Computed: true,
            },
            "group_by_field": schema.StringAttribute{
                MarkdownDescription: "Optional security-event field (e.g. principalHost, principalUser) to group matches by. One alert is opened per distinct value; empty groups all matches into one alert..",
                Computed: true,
            },
            "distinct_count_field": schema.StringAttribute{
                MarkdownDescription: "Optional security-event field (e.g. principalUser, principalIp) whose distinct values are counted instead of raw matching events. The match count threshold then applies to that distinct count. Empty values are not counted. Names that are not typed event columns are looked up in the event's attributes map..",
                Computed: true,
            },
            "match_count_threshold": schema.NumberAttribute{
                MarkdownDescription: "Fire only when a group's count — distinct values when a distinct count field is set, matching events otherwise — reaches this number within one evaluation window. 1 fires on any match..",
                Computed: true,
            },
            "should_create_alert": schema.BoolAttribute{
                MarkdownDescription: "Whether matches open OneUptime alerts..",
                Computed: true,
            },
            "should_write_detection_finding": schema.BoolAttribute{
                MarkdownDescription: "Whether matches also write a Detection Finding security event back into the events table..",
                Computed: true,
            },
            "should_create_incident": schema.BoolAttribute{
                MarkdownDescription: "Whether matches also open OneUptime incidents. Off by default: incidents drive on-call, SLAs and status pages, so opt in per rule..",
                Computed: true,
            },
            "alert_severity_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "incident_severity_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "last_evaluated_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "last_match_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "last_error": schema.StringAttribute{
                MarkdownDescription: "The most recent evaluation error, if any. Cleared on the next successful evaluation..",
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

func (d *DetectionRuleDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DetectionRuleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data DetectionRuleDataSourceModel

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
            "Exactly one of `id` or `name` must be set to look up a detection_rule.",
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
        "description": true,
        "sigmaRuleYaml": true,
        "isEnabled": true,
        "evaluationIntervalInMinutes": true,
        "groupByField": true,
        "distinctCountField": true,
        "matchCountThreshold": true,
        "shouldCreateAlert": true,
        "shouldWriteDetectionFinding": true,
        "shouldCreateIncident": true,
        "alertSeverityId": true,
        "incidentSeverityId": true,
        "lastEvaluatedAt": true,
        "lastMatchAt": true,
        "lastError": true,
        "createdByUserId": true,
        "deletedByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/detection-rule/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read detection_rule, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No detection_rule found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read detection_rule: %s", err))
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
        httpResp, err := d.client.Post(ctx, "/detection-rule/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list detection_rule, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list detection_rule: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No detection_rule found with name %q.", data.Name.ValueString()))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one detection_rule matches name %q. Use the id attribute to disambiguate.", data.Name.ValueString()))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for detection_rule.")
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
    if obj, ok := item["sigmaRuleYaml"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SigmaRuleYaml = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SigmaRuleYaml = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SigmaRuleYaml = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SigmaRuleYaml = types.StringValue(string(jsonBytes))
        } else {
            data.SigmaRuleYaml = types.StringNull()
        }
    } else if val, ok := item["sigmaRuleYaml"].(string); ok {
        data.SigmaRuleYaml = types.StringValue(val)
    } else {
        data.SigmaRuleYaml = types.StringNull()
    }
    if val, ok := item["isEnabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    } else {
        data.IsEnabled = types.BoolNull()
    }
    if val, ok := item["evaluationIntervalInMinutes"].(float64); ok {
        data.EvaluationIntervalInMinutes = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["evaluationIntervalInMinutes"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.EvaluationIntervalInMinutes = types.NumberValue(big.NewFloat(val))
        } else {
            data.EvaluationIntervalInMinutes = types.NumberNull()
        }
    } else {
        data.EvaluationIntervalInMinutes = types.NumberNull()
    }
    if obj, ok := item["groupByField"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.GroupByField = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.GroupByField = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.GroupByField = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.GroupByField = types.StringValue(string(jsonBytes))
        } else {
            data.GroupByField = types.StringNull()
        }
    } else if val, ok := item["groupByField"].(string); ok {
        data.GroupByField = types.StringValue(val)
    } else {
        data.GroupByField = types.StringNull()
    }
    if obj, ok := item["distinctCountField"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DistinctCountField = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DistinctCountField = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DistinctCountField = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DistinctCountField = types.StringValue(string(jsonBytes))
        } else {
            data.DistinctCountField = types.StringNull()
        }
    } else if val, ok := item["distinctCountField"].(string); ok {
        data.DistinctCountField = types.StringValue(val)
    } else {
        data.DistinctCountField = types.StringNull()
    }
    if val, ok := item["matchCountThreshold"].(float64); ok {
        data.MatchCountThreshold = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["matchCountThreshold"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.MatchCountThreshold = types.NumberValue(big.NewFloat(val))
        } else {
            data.MatchCountThreshold = types.NumberNull()
        }
    } else {
        data.MatchCountThreshold = types.NumberNull()
    }
    if val, ok := item["shouldCreateAlert"].(bool); ok {
        data.ShouldCreateAlert = types.BoolValue(val)
    } else {
        data.ShouldCreateAlert = types.BoolNull()
    }
    if val, ok := item["shouldWriteDetectionFinding"].(bool); ok {
        data.ShouldWriteDetectionFinding = types.BoolValue(val)
    } else {
        data.ShouldWriteDetectionFinding = types.BoolNull()
    }
    if val, ok := item["shouldCreateIncident"].(bool); ok {
        data.ShouldCreateIncident = types.BoolValue(val)
    } else {
        data.ShouldCreateIncident = types.BoolNull()
    }
    if obj, ok := item["alertSeverityId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AlertSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AlertSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AlertSeverityId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AlertSeverityId = types.StringValue(string(jsonBytes))
        } else {
            data.AlertSeverityId = types.StringNull()
        }
    } else if val, ok := item["alertSeverityId"].(string); ok {
        data.AlertSeverityId = types.StringValue(val)
    } else {
        data.AlertSeverityId = types.StringNull()
    }
    if obj, ok := item["incidentSeverityId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncidentSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.IncidentSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.IncidentSeverityId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.IncidentSeverityId = types.StringValue(string(jsonBytes))
        } else {
            data.IncidentSeverityId = types.StringNull()
        }
    } else if val, ok := item["incidentSeverityId"].(string); ok {
        data.IncidentSeverityId = types.StringValue(val)
    } else {
        data.IncidentSeverityId = types.StringNull()
    }
    if obj, ok := item["lastEvaluatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastEvaluatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastEvaluatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastEvaluatedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastEvaluatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastEvaluatedAt = types.StringNull()
        }
    } else if val, ok := item["lastEvaluatedAt"].(string); ok {
        data.LastEvaluatedAt = types.StringValue(val)
    } else {
        data.LastEvaluatedAt = types.StringNull()
    }
    if obj, ok := item["lastMatchAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastMatchAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastMatchAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastMatchAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastMatchAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastMatchAt = types.StringNull()
        }
    } else if val, ok := item["lastMatchAt"].(string); ok {
        data.LastMatchAt = types.StringValue(val)
    } else {
        data.LastMatchAt = types.StringNull()
    }
    if obj, ok := item["lastError"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastError = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastError = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastError = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastError = types.StringValue(string(jsonBytes))
        } else {
            data.LastError = types.StringNull()
        }
    } else if val, ok := item["lastError"].(string); ok {
        data.LastError = types.StringValue(val)
    } else {
        data.LastError = types.StringNull()
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
