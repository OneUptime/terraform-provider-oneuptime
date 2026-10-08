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
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    Name types.String `tfsdk:"name"`
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
}

func (d *DetectionRuleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_detection_rule"
}

func (d *DetectionRuleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Sigma detection rules evaluated against security events. Matches create alerts and detection findings. Look up an existing detection rule by `id`, or by any of its other arguments (`name`, `alert_severity_id`, `created_by_user_id`, ...): each one set must match, and exactly one detection rule may match them all.",

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
                MarkdownDescription: "ID of the project this detection rule belongs to. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Friendly name for this detection rule.",
                Optional: true,
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Description of what this detection rule looks for.",
                Optional: true,
                Computed: true,
            },
            "sigma_rule_yaml": schema.StringAttribute{
                MarkdownDescription: "The Sigma rule to evaluate, in YAML. detection selections and condition are compiled to a ClickHouse query over security events.",
                Optional: true,
                Computed: true,
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether this detection rule is evaluated.",
                Optional: true,
                Computed: true,
            },
            "evaluation_interval_in_minutes": schema.NumberAttribute{
                MarkdownDescription: "How often this rule is evaluated, in minutes. The evaluation window covers the time since the previous evaluation.",
                Optional: true,
                Computed: true,
            },
            "group_by_field": schema.StringAttribute{
                MarkdownDescription: "Optional security-event field (e.g. principalHost, principalUser) to group matches by. One alert is opened per distinct value; empty groups all matches into one alert.",
                Optional: true,
                Computed: true,
            },
            "distinct_count_field": schema.StringAttribute{
                MarkdownDescription: "Optional security-event field (e.g. principalUser, principalIp) whose distinct values are counted instead of raw matching events. The match count threshold then applies to that distinct count. Empty values are not counted. Names that are not typed event columns are looked up in the event's attributes map.",
                Optional: true,
                Computed: true,
            },
            "match_count_threshold": schema.NumberAttribute{
                MarkdownDescription: "Fire only when a group's count — distinct values when a distinct count field is set, matching events otherwise — reaches this number within one evaluation window. 1 fires on any match.",
                Optional: true,
                Computed: true,
            },
            "should_create_alert": schema.BoolAttribute{
                MarkdownDescription: "Whether matches open OneUptime alerts.",
                Optional: true,
                Computed: true,
            },
            "should_write_detection_finding": schema.BoolAttribute{
                MarkdownDescription: "Whether matches also write a Detection Finding security event back into the events table.",
                Optional: true,
                Computed: true,
            },
            "should_create_incident": schema.BoolAttribute{
                MarkdownDescription: "Whether matches also open OneUptime incidents. Off by default: incidents drive on-call, SLAs and status pages, so opt in per rule.",
                Optional: true,
                Computed: true,
            },
            "alert_severity_id": schema.StringAttribute{
                MarkdownDescription: "ID of the alert severity for alerts opened by this rule. The ID of a `oneuptime_alert_severity`.",
                Optional: true,
                Computed: true,
            },
            "incident_severity_id": schema.StringAttribute{
                MarkdownDescription: "ID of the incident severity for incidents opened by this rule. The ID of a `oneuptime_incident_severity`.",
                Optional: true,
                Computed: true,
            },
            "last_evaluated_at": schema.StringAttribute{
                MarkdownDescription: "When the detection engine last evaluated this rule. Null means it has never run.",
                Computed: true,
            },
            "last_match_at": schema.StringAttribute{
                MarkdownDescription: "When this rule most recently matched security events. Null means it has never matched.",
                Computed: true,
            },
            "last_error": schema.StringAttribute{
                MarkdownDescription: "The most recent evaluation error, if any. Cleared on the next successful evaluation.",
                Optional: true,
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the user who created this detection rule. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
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

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.Name.IsNull() && !data.Name.IsUnknown() {
        filters["name"] = data.Name.ValueString()
        filterNames = append(filterNames, "name = "+fmt.Sprintf("%q", data.Name.ValueString()))
    }
    if !data.Description.IsNull() && !data.Description.IsUnknown() {
        filters["description"] = data.Description.ValueString()
        filterNames = append(filterNames, "description = "+fmt.Sprintf("%q", data.Description.ValueString()))
    }
    if !data.SigmaRuleYaml.IsNull() && !data.SigmaRuleYaml.IsUnknown() {
        filters["sigmaRuleYaml"] = data.SigmaRuleYaml.ValueString()
        filterNames = append(filterNames, "sigma_rule_yaml = "+fmt.Sprintf("%q", data.SigmaRuleYaml.ValueString()))
    }
    if !data.IsEnabled.IsNull() && !data.IsEnabled.IsUnknown() {
        filters["isEnabled"] = data.IsEnabled.ValueBool()
        filterNames = append(filterNames, "is_enabled = "+fmt.Sprintf("%t", data.IsEnabled.ValueBool()))
    }
    if !data.EvaluationIntervalInMinutes.IsNull() && !data.EvaluationIntervalInMinutes.IsUnknown() {
        filters["evaluationIntervalInMinutes"] = lookupNumber(data.EvaluationIntervalInMinutes)
        filterNames = append(filterNames, "evaluation_interval_in_minutes = "+data.EvaluationIntervalInMinutes.ValueBigFloat().String())
    }
    if !data.GroupByField.IsNull() && !data.GroupByField.IsUnknown() {
        filters["groupByField"] = data.GroupByField.ValueString()
        filterNames = append(filterNames, "group_by_field = "+fmt.Sprintf("%q", data.GroupByField.ValueString()))
    }
    if !data.DistinctCountField.IsNull() && !data.DistinctCountField.IsUnknown() {
        filters["distinctCountField"] = data.DistinctCountField.ValueString()
        filterNames = append(filterNames, "distinct_count_field = "+fmt.Sprintf("%q", data.DistinctCountField.ValueString()))
    }
    if !data.MatchCountThreshold.IsNull() && !data.MatchCountThreshold.IsUnknown() {
        filters["matchCountThreshold"] = lookupNumber(data.MatchCountThreshold)
        filterNames = append(filterNames, "match_count_threshold = "+data.MatchCountThreshold.ValueBigFloat().String())
    }
    if !data.ShouldCreateAlert.IsNull() && !data.ShouldCreateAlert.IsUnknown() {
        filters["shouldCreateAlert"] = data.ShouldCreateAlert.ValueBool()
        filterNames = append(filterNames, "should_create_alert = "+fmt.Sprintf("%t", data.ShouldCreateAlert.ValueBool()))
    }
    if !data.ShouldWriteDetectionFinding.IsNull() && !data.ShouldWriteDetectionFinding.IsUnknown() {
        filters["shouldWriteDetectionFinding"] = data.ShouldWriteDetectionFinding.ValueBool()
        filterNames = append(filterNames, "should_write_detection_finding = "+fmt.Sprintf("%t", data.ShouldWriteDetectionFinding.ValueBool()))
    }
    if !data.ShouldCreateIncident.IsNull() && !data.ShouldCreateIncident.IsUnknown() {
        filters["shouldCreateIncident"] = data.ShouldCreateIncident.ValueBool()
        filterNames = append(filterNames, "should_create_incident = "+fmt.Sprintf("%t", data.ShouldCreateIncident.ValueBool()))
    }
    if !data.AlertSeverityId.IsNull() && !data.AlertSeverityId.IsUnknown() {
        filters["alertSeverityId"] = data.AlertSeverityId.ValueString()
        filterNames = append(filterNames, "alert_severity_id = "+fmt.Sprintf("%q", data.AlertSeverityId.ValueString()))
    }
    if !data.IncidentSeverityId.IsNull() && !data.IncidentSeverityId.IsUnknown() {
        filters["incidentSeverityId"] = data.IncidentSeverityId.ValueString()
        filterNames = append(filterNames, "incident_severity_id = "+fmt.Sprintf("%q", data.IncidentSeverityId.ValueString()))
    }
    if !data.LastError.IsNull() && !data.LastError.IsUnknown() {
        filters["lastError"] = data.LastError.ValueString()
        filterNames = append(filterNames, "last_error = "+fmt.Sprintf("%q", data.LastError.ValueString()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the detection rule up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the detection rule up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "name": true,
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
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No detection rule found with id %q.", data.Id.ValueString()))
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
    }
    if !hasId {
        listBody := map[string]interface{}{
            "query":  filters,
            "select": selectParam,
            // limit 2 is enough to detect ambiguity without paging.
            "limit": 2,
        }
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/detection-rule/get-list", listBody)
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
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No detection rule matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one detection rule matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
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

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
