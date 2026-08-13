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
var _ datasource.DataSource = &AiInsightDataSource{}

func NewAiInsightDataSource() datasource.DataSource {
    return &AiInsightDataSource{}
}

// AiInsightDataSource defines the data source implementation.
type AiInsightDataSource struct {
    client *Client
}

// AiInsightDataSourceModel describes the data source data model.
type AiInsightDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    InsightType types.String `tfsdk:"insight_type"`
    Status types.String `tfsdk:"status"`
    Severity types.String `tfsdk:"severity"`
    Classification types.String `tfsdk:"classification"`
    Fingerprint types.String `tfsdk:"fingerprint"`
    Title types.String `tfsdk:"title"`
    DetailMarkdown types.String `tfsdk:"detail_markdown"`
    ServiceName types.String `tfsdk:"service_name"`
    TelemetryServiceId types.String `tfsdk:"telemetry_service_id"`
    TelemetryExceptionId types.String `tfsdk:"telemetry_exception_id"`
    TraceId types.String `tfsdk:"trace_id"`
    MetricName types.String `tfsdk:"metric_name"`
    Evidence types.String `tfsdk:"evidence"`
    FirstSeenAt types.String `tfsdk:"first_seen_at"`
    LastSeenAt types.String `tfsdk:"last_seen_at"`
    OccurrenceCount types.Number `tfsdk:"occurrence_count"`
    TriageAiRunId types.String `tfsdk:"triage_ai_run_id"`
    FixAiRunId types.String `tfsdk:"fix_ai_run_id"`
    TriageSummaryMarkdown types.String `tfsdk:"triage_summary_markdown"`
    TriageCompletedAt types.String `tfsdk:"triage_completed_at"`
    HumanVerdict types.String `tfsdk:"human_verdict"`
    HumanVerdictAt types.String `tfsdk:"human_verdict_at"`
    HumanVerdictByUserId types.String `tfsdk:"human_verdict_by_user_id"`
}

func (d *AiInsightDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_ai_insight"
}

func (d *AiInsightDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "A preventive finding from OneUptime AI's deterministic telemetry sensors — new or spiking exceptions, error-log spikes, trace-latency regressions and metric drift — surfaced in a quiet insights inbox that never pages and never opens incidents. Look up an existing ai_insight by `id` or by `name`.",

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
            "insight_type": schema.StringAttribute{
                MarkdownDescription: "Which deterministic detector produced this insight: NewException, ExceptionSpike, ErrorLogSpike, TraceLatencyRegression or MetricDrift..",
                Computed: true,
            },
            "status": schema.StringAttribute{
                MarkdownDescription: "Lifecycle of the insight. Detected is the defensive initial state — the scanner routes to ActionRequired or FixOpened in the same tick; Resolved and Dismissed are human actions..",
                Computed: true,
            },
            "severity": schema.StringAttribute{
                MarkdownDescription: "How urgent this insight is (High, Medium or Low), assigned deterministically by the detector..",
                Computed: true,
            },
            "classification": schema.StringAttribute{
                MarkdownDescription: "AI triage verdict: code-fault, user-error, expected-denial, infrastructure or unknown. Automatic fix pull requests are only opened for code-fault..",
                Computed: true,
            },
            "fingerprint": schema.StringAttribute{
                MarkdownDescription: "The detector's stable dedupe key for this finding. Recurring detections refresh the existing non-terminal insight with the same fingerprint..",
                Computed: true,
            },
            "title": schema.StringAttribute{
                MarkdownDescription: "One-line human-readable summary of the finding..",
                Computed: true,
            },
            "detail_markdown": schema.StringAttribute{
                MarkdownDescription: "The deterministic evidence rendered as markdown: real counts, baselines and multipliers written by the detector at detect time..",
                Computed: true,
            },
            "service_name": schema.StringAttribute{
                MarkdownDescription: "Name of the telemetry service this insight is about..",
                Computed: true,
            },
            "telemetry_service_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "telemetry_exception_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "trace_id": schema.StringAttribute{
                MarkdownDescription: "A representative slow trace (for TraceLatencyRegression insights)..",
                Computed: true,
            },
            "metric_name": schema.StringAttribute{
                MarkdownDescription: "The drifting metric's name (for MetricDrift insights)..",
                Computed: true,
            },
            "evidence": schema.StringAttribute{
                MarkdownDescription: "The deterministic evidence computed at detect time: counts, baselines, multipliers and (for latency insights) span-tree findings..",
                Computed: true,
            },
            "first_seen_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "last_seen_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "occurrence_count": schema.NumberAttribute{
                MarkdownDescription: "How many scanner ticks have detected this finding. Incremented on each dedupe refresh..",
                Computed: true,
            },
            "triage_ai_run_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "fix_ai_run_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "triage_summary_markdown": schema.StringAttribute{
                MarkdownDescription: "The AI triage analysis for this insight: probable root cause, blast radius and suggested action, with citations..",
                Computed: true,
            },
            "triage_completed_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "human_verdict": schema.StringAttribute{
                MarkdownDescription: "The one-click human verdict on this insight (Confirmed or Dismissed). Null until a user weighs in..",
                Computed: true,
            },
            "human_verdict_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "human_verdict_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
        },
    }
}

func (d *AiInsightDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AiInsightDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data AiInsightDataSourceModel

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
            "Exactly one of `id` or `name` must be set to look up a ai_insight.",
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
        "insightType": true,
        "status": true,
        "severity": true,
        "classification": true,
        "fingerprint": true,
        "title": true,
        "detailMarkdown": true,
        "serviceName": true,
        "telemetryServiceId": true,
        "telemetryExceptionId": true,
        "traceId": true,
        "metricName": true,
        "evidence": true,
        "firstSeenAt": true,
        "lastSeenAt": true,
        "occurrenceCount": true,
        "triageAiRunId": true,
        "fixAiRunId": true,
        "triageSummaryMarkdown": true,
        "triageCompletedAt": true,
        "humanVerdict": true,
        "humanVerdictAt": true,
        "humanVerdictByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/ai-insight/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read ai_insight, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No ai_insight found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read ai_insight: %s", err))
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
        httpResp, err := d.client.Post(ctx, "/ai-insight/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list ai_insight, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list ai_insight: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No ai_insight found with name %q.", data.Name.ValueString()))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one ai_insight matches name %q. Use the id attribute to disambiguate.", data.Name.ValueString()))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for ai_insight.")
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
    if obj, ok := item["insightType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.InsightType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.InsightType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.InsightType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.InsightType = types.StringValue(string(jsonBytes))
        } else {
            data.InsightType = types.StringNull()
        }
    } else if val, ok := item["insightType"].(string); ok {
        data.InsightType = types.StringValue(val)
    } else {
        data.InsightType = types.StringNull()
    }
    if obj, ok := item["status"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Status = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Status = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Status = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Status = types.StringValue(string(jsonBytes))
        } else {
            data.Status = types.StringNull()
        }
    } else if val, ok := item["status"].(string); ok {
        data.Status = types.StringValue(val)
    } else {
        data.Status = types.StringNull()
    }
    if obj, ok := item["severity"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Severity = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Severity = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Severity = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Severity = types.StringValue(string(jsonBytes))
        } else {
            data.Severity = types.StringNull()
        }
    } else if val, ok := item["severity"].(string); ok {
        data.Severity = types.StringValue(val)
    } else {
        data.Severity = types.StringNull()
    }
    if obj, ok := item["classification"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Classification = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Classification = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Classification = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Classification = types.StringValue(string(jsonBytes))
        } else {
            data.Classification = types.StringNull()
        }
    } else if val, ok := item["classification"].(string); ok {
        data.Classification = types.StringValue(val)
    } else {
        data.Classification = types.StringNull()
    }
    if obj, ok := item["fingerprint"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Fingerprint = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Fingerprint = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Fingerprint = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Fingerprint = types.StringValue(string(jsonBytes))
        } else {
            data.Fingerprint = types.StringNull()
        }
    } else if val, ok := item["fingerprint"].(string); ok {
        data.Fingerprint = types.StringValue(val)
    } else {
        data.Fingerprint = types.StringNull()
    }
    if obj, ok := item["title"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Title = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Title = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Title = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Title = types.StringValue(string(jsonBytes))
        } else {
            data.Title = types.StringNull()
        }
    } else if val, ok := item["title"].(string); ok {
        data.Title = types.StringValue(val)
    } else {
        data.Title = types.StringNull()
    }
    if obj, ok := item["detailMarkdown"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DetailMarkdown = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DetailMarkdown = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DetailMarkdown = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DetailMarkdown = types.StringValue(string(jsonBytes))
        } else {
            data.DetailMarkdown = types.StringNull()
        }
    } else if val, ok := item["detailMarkdown"].(string); ok {
        data.DetailMarkdown = types.StringValue(val)
    } else {
        data.DetailMarkdown = types.StringNull()
    }
    if obj, ok := item["serviceName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ServiceName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ServiceName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ServiceName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ServiceName = types.StringValue(string(jsonBytes))
        } else {
            data.ServiceName = types.StringNull()
        }
    } else if val, ok := item["serviceName"].(string); ok {
        data.ServiceName = types.StringValue(val)
    } else {
        data.ServiceName = types.StringNull()
    }
    if obj, ok := item["telemetryServiceId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TelemetryServiceId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TelemetryServiceId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TelemetryServiceId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TelemetryServiceId = types.StringValue(string(jsonBytes))
        } else {
            data.TelemetryServiceId = types.StringNull()
        }
    } else if val, ok := item["telemetryServiceId"].(string); ok {
        data.TelemetryServiceId = types.StringValue(val)
    } else {
        data.TelemetryServiceId = types.StringNull()
    }
    if obj, ok := item["telemetryExceptionId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TelemetryExceptionId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TelemetryExceptionId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TelemetryExceptionId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TelemetryExceptionId = types.StringValue(string(jsonBytes))
        } else {
            data.TelemetryExceptionId = types.StringNull()
        }
    } else if val, ok := item["telemetryExceptionId"].(string); ok {
        data.TelemetryExceptionId = types.StringValue(val)
    } else {
        data.TelemetryExceptionId = types.StringNull()
    }
    if obj, ok := item["traceId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TraceId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TraceId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TraceId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TraceId = types.StringValue(string(jsonBytes))
        } else {
            data.TraceId = types.StringNull()
        }
    } else if val, ok := item["traceId"].(string); ok {
        data.TraceId = types.StringValue(val)
    } else {
        data.TraceId = types.StringNull()
    }
    if obj, ok := item["metricName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MetricName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.MetricName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.MetricName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.MetricName = types.StringValue(string(jsonBytes))
        } else {
            data.MetricName = types.StringNull()
        }
    } else if val, ok := item["metricName"].(string); ok {
        data.MetricName = types.StringValue(val)
    } else {
        data.MetricName = types.StringNull()
    }
    if obj, ok := item["evidence"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Evidence = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Evidence = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Evidence = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Evidence = types.StringValue(string(jsonBytes))
        } else {
            data.Evidence = types.StringNull()
        }
    } else if val, ok := item["evidence"].(string); ok {
        data.Evidence = types.StringValue(val)
    } else {
        data.Evidence = types.StringNull()
    }
    if obj, ok := item["firstSeenAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FirstSeenAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.FirstSeenAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.FirstSeenAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.FirstSeenAt = types.StringValue(string(jsonBytes))
        } else {
            data.FirstSeenAt = types.StringNull()
        }
    } else if val, ok := item["firstSeenAt"].(string); ok {
        data.FirstSeenAt = types.StringValue(val)
    } else {
        data.FirstSeenAt = types.StringNull()
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
    if val, ok := item["occurrenceCount"].(float64); ok {
        data.OccurrenceCount = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["occurrenceCount"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.OccurrenceCount = types.NumberValue(big.NewFloat(val))
        } else {
            data.OccurrenceCount = types.NumberNull()
        }
    } else {
        data.OccurrenceCount = types.NumberNull()
    }
    if obj, ok := item["triageAiRunId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TriageAiRunId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TriageAiRunId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TriageAiRunId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TriageAiRunId = types.StringValue(string(jsonBytes))
        } else {
            data.TriageAiRunId = types.StringNull()
        }
    } else if val, ok := item["triageAiRunId"].(string); ok {
        data.TriageAiRunId = types.StringValue(val)
    } else {
        data.TriageAiRunId = types.StringNull()
    }
    if obj, ok := item["fixAiRunId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FixAiRunId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.FixAiRunId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.FixAiRunId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.FixAiRunId = types.StringValue(string(jsonBytes))
        } else {
            data.FixAiRunId = types.StringNull()
        }
    } else if val, ok := item["fixAiRunId"].(string); ok {
        data.FixAiRunId = types.StringValue(val)
    } else {
        data.FixAiRunId = types.StringNull()
    }
    if obj, ok := item["triageSummaryMarkdown"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TriageSummaryMarkdown = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TriageSummaryMarkdown = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TriageSummaryMarkdown = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TriageSummaryMarkdown = types.StringValue(string(jsonBytes))
        } else {
            data.TriageSummaryMarkdown = types.StringNull()
        }
    } else if val, ok := item["triageSummaryMarkdown"].(string); ok {
        data.TriageSummaryMarkdown = types.StringValue(val)
    } else {
        data.TriageSummaryMarkdown = types.StringNull()
    }
    if obj, ok := item["triageCompletedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TriageCompletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TriageCompletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TriageCompletedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TriageCompletedAt = types.StringValue(string(jsonBytes))
        } else {
            data.TriageCompletedAt = types.StringNull()
        }
    } else if val, ok := item["triageCompletedAt"].(string); ok {
        data.TriageCompletedAt = types.StringValue(val)
    } else {
        data.TriageCompletedAt = types.StringNull()
    }
    if obj, ok := item["humanVerdict"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.HumanVerdict = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.HumanVerdict = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.HumanVerdict = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.HumanVerdict = types.StringValue(string(jsonBytes))
        } else {
            data.HumanVerdict = types.StringNull()
        }
    } else if val, ok := item["humanVerdict"].(string); ok {
        data.HumanVerdict = types.StringValue(val)
    } else {
        data.HumanVerdict = types.StringNull()
    }
    if obj, ok := item["humanVerdictAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.HumanVerdictAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.HumanVerdictAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.HumanVerdictAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.HumanVerdictAt = types.StringValue(string(jsonBytes))
        } else {
            data.HumanVerdictAt = types.StringNull()
        }
    } else if val, ok := item["humanVerdictAt"].(string); ok {
        data.HumanVerdictAt = types.StringValue(val)
    } else {
        data.HumanVerdictAt = types.StringNull()
    }
    if obj, ok := item["humanVerdictByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.HumanVerdictByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.HumanVerdictByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.HumanVerdictByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.HumanVerdictByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.HumanVerdictByUserId = types.StringNull()
        }
    } else if val, ok := item["humanVerdictByUserId"].(string); ok {
        data.HumanVerdictByUserId = types.StringValue(val)
    } else {
        data.HumanVerdictByUserId = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
