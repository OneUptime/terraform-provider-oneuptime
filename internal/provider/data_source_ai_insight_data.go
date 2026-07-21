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
var _ datasource.DataSource = &AiInsightDataDataSource{}

func NewAiInsightDataDataSource() datasource.DataSource {
    return &AiInsightDataDataSource{}
}

// AiInsightDataDataSource defines the data source implementation.
type AiInsightDataDataSource struct {
    client *Client
}

// AiInsightDataDataSourceModel describes the data source data model.
type AiInsightDataDataSourceModel struct {
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

func (d *AiInsightDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_ai_insight_data"
}

func (d *AiInsightDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "ai_insight_data data source",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Identifier to filter by",
                Optional: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Name to filter by",
                Optional: true,
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
                MarkdownDescription: "Which deterministic detector produced this insight: NewException, ExceptionSpike, ErrorLogSpike, TraceLatencyRegression or MetricDrift.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "status": schema.StringAttribute{
                MarkdownDescription: "Lifecycle of the insight. Detected is the defensive initial state — the scanner routes to ActionRequired or FixOpened in the same tick; Resolved and Dismissed are human actions.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "severity": schema.StringAttribute{
                MarkdownDescription: "How urgent this insight is (High, Medium or Low), assigned deterministically by the detector.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "classification": schema.StringAttribute{
                MarkdownDescription: "AI triage verdict: code-fault, user-error, expected-denial, infrastructure or unknown. Automatic fix pull requests are only opened for code-fault.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "fingerprint": schema.StringAttribute{
                MarkdownDescription: "The detector's stable dedupe key for this finding. Recurring detections refresh the existing non-terminal insight with the same fingerprint.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "title": schema.StringAttribute{
                MarkdownDescription: "One-line human-readable summary of the finding.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "detail_markdown": schema.StringAttribute{
                MarkdownDescription: "The deterministic evidence rendered as markdown: real counts, baselines and multipliers written by the detector at detect time.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "service_name": schema.StringAttribute{
                MarkdownDescription: "Name of the telemetry service this insight is about.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member], Update: [No access - you don't have permission for this operation]",
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
                MarkdownDescription: "A representative slow trace (for TraceLatencyRegression insights).. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "metric_name": schema.StringAttribute{
                MarkdownDescription: "The drifting metric's name (for MetricDrift insights).. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "evidence": schema.StringAttribute{
                MarkdownDescription: "The deterministic evidence computed at detect time: counts, baselines, multipliers and (for latency insights) span-tree findings.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member], Update: [No access - you don't have permission for this operation]",
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
                MarkdownDescription: "How many scanner ticks have detected this finding. Incremented on each dedupe refresh.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member], Update: [No access - you don't have permission for this operation]",
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
                MarkdownDescription: "The AI triage analysis for this insight: probable root cause, blast radius and suggested action, with citations.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "triage_completed_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "human_verdict": schema.StringAttribute{
                MarkdownDescription: "The one-click human verdict on this insight (Confirmed or Dismissed). Null until a user weighs in.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member], Update: [No access - you don't have permission for this operation]",
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

func (d *AiInsightDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AiInsightDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data AiInsightDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "ai-insight" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read ai_insight_data, got error: %s", err))
        return
    }

    var aiInsightDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &aiInsightDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse ai_insight_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := aiInsightDataResponse["data"].(map[string]interface{}); ok {
        aiInsightDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := aiInsightDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := aiInsightDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := aiInsightDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := aiInsightDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := aiInsightDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := aiInsightDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := aiInsightDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := aiInsightDataResponse["insight_type"].(string); ok {
        data.InsightType = types.StringValue(val)
    }
    if val, ok := aiInsightDataResponse["status"].(string); ok {
        data.Status = types.StringValue(val)
    }
    if val, ok := aiInsightDataResponse["severity"].(string); ok {
        data.Severity = types.StringValue(val)
    }
    if val, ok := aiInsightDataResponse["classification"].(string); ok {
        data.Classification = types.StringValue(val)
    }
    if val, ok := aiInsightDataResponse["fingerprint"].(string); ok {
        data.Fingerprint = types.StringValue(val)
    }
    if val, ok := aiInsightDataResponse["title"].(string); ok {
        data.Title = types.StringValue(val)
    }
    if val, ok := aiInsightDataResponse["detail_markdown"].(string); ok {
        data.DetailMarkdown = types.StringValue(val)
    }
    if val, ok := aiInsightDataResponse["service_name"].(string); ok {
        data.ServiceName = types.StringValue(val)
    }
    if val, ok := aiInsightDataResponse["telemetry_service_id"].(string); ok {
        data.TelemetryServiceId = types.StringValue(val)
    }
    if val, ok := aiInsightDataResponse["telemetry_exception_id"].(string); ok {
        data.TelemetryExceptionId = types.StringValue(val)
    }
    if val, ok := aiInsightDataResponse["trace_id"].(string); ok {
        data.TraceId = types.StringValue(val)
    }
    if val, ok := aiInsightDataResponse["metric_name"].(string); ok {
        data.MetricName = types.StringValue(val)
    }
    if val, ok := aiInsightDataResponse["evidence"].(string); ok {
        data.Evidence = types.StringValue(val)
    }
    if val, ok := aiInsightDataResponse["first_seen_at"].(string); ok {
        data.FirstSeenAt = types.StringValue(val)
    }
    if val, ok := aiInsightDataResponse["last_seen_at"].(string); ok {
        data.LastSeenAt = types.StringValue(val)
    }
    if val, ok := aiInsightDataResponse["occurrence_count"].(float64); ok {
        data.OccurrenceCount = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := aiInsightDataResponse["triage_ai_run_id"].(string); ok {
        data.TriageAiRunId = types.StringValue(val)
    }
    if val, ok := aiInsightDataResponse["fix_ai_run_id"].(string); ok {
        data.FixAiRunId = types.StringValue(val)
    }
    if val, ok := aiInsightDataResponse["triage_summary_markdown"].(string); ok {
        data.TriageSummaryMarkdown = types.StringValue(val)
    }
    if val, ok := aiInsightDataResponse["triage_completed_at"].(string); ok {
        data.TriageCompletedAt = types.StringValue(val)
    }
    if val, ok := aiInsightDataResponse["human_verdict"].(string); ok {
        data.HumanVerdict = types.StringValue(val)
    }
    if val, ok := aiInsightDataResponse["human_verdict_at"].(string); ok {
        data.HumanVerdictAt = types.StringValue(val)
    }
    if val, ok := aiInsightDataResponse["human_verdict_by_user_id"].(string); ok {
        data.HumanVerdictByUserId = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
