package provider

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "math/big"
    "github.com/hashicorp/terraform-plugin-framework/attr"
    "sort"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &RumApplicationDataSource{}

func NewRumApplicationDataSource() datasource.DataSource {
    return &RumApplicationDataSource{}
}

// RumApplicationDataSource defines the data source implementation.
type RumApplicationDataSource struct {
    client *Client
}

// RumApplicationDataSourceModel describes the data source data model.
type RumApplicationDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    Name types.String `tfsdk:"name"`
    Slug types.String `tfsdk:"slug"`
    Description types.String `tfsdk:"description"`
    AppIdentifier types.String `tfsdk:"app_identifier"`
    ClientType types.String `tfsdk:"client_type"`
    SdkLanguage types.String `tfsdk:"sdk_language"`
    OtelCollectorStatus types.String `tfsdk:"otel_collector_status"`
    AgentVersion types.String `tfsdk:"agent_version"`
    LastSeenAt types.String `tfsdk:"last_seen_at"`
    Labels types.Set `tfsdk:"labels"`
    RetainTelemetryDataForDays types.Number `tfsdk:"retain_telemetry_data_for_days"`
    TelemetryRetentionConfig types.String `tfsdk:"telemetry_retention_config"`
    IsSessionReplayEnabled types.Bool `tfsdk:"is_session_replay_enabled"`
    SessionReplayMaskingMode types.String `tfsdk:"session_replay_masking_mode"`
    SessionReplayMaskSelectors types.String `tfsdk:"session_replay_mask_selectors"`
    SessionReplayBlockSelectors types.String `tfsdk:"session_replay_block_selectors"`
    SessionReplayIgnoreErrorPatterns types.String `tfsdk:"session_replay_ignore_error_patterns"`
    SessionReplayTracePropagationOrigins types.String `tfsdk:"session_replay_trace_propagation_origins"`
    SessionReplaySameOriginTracePropagation types.Bool `tfsdk:"session_replay_same_origin_trace_propagation"`
    SessionReplayLcpBudgetMs types.Number `tfsdk:"session_replay_lcp_budget_ms"`
    SessionReplayLongTaskBudgetMs types.Number `tfsdk:"session_replay_long_task_budget_ms"`
    SessionReplaySlowRequestBudgetMs types.Number `tfsdk:"session_replay_slow_request_budget_ms"`
    SessionReplayAllowedOrigins types.String `tfsdk:"session_replay_allowed_origins"`
    SessionReplayConsentMode types.String `tfsdk:"session_replay_consent_mode"`
    SessionReplayCaptureTrigger types.String `tfsdk:"session_replay_capture_trigger"`
    SessionReplaySamplePercentage types.Number `tfsdk:"session_replay_sample_percentage"`
    SessionReplayCaptureUserIdentity types.Bool `tfsdk:"session_replay_capture_user_identity"`
    SessionReplayCaptureGeo types.Bool `tfsdk:"session_replay_capture_geo"`
    SessionReplayRecordCanvas types.Bool `tfsdk:"session_replay_record_canvas"`
    SessionReplayRetentionInDays types.Number `tfsdk:"session_replay_retention_in_days"`
    SessionReplayMonthlyBudgetInGb types.Number `tfsdk:"session_replay_monthly_budget_in_gb"`
    SessionReplayLastChunkReceivedAt types.String `tfsdk:"session_replay_last_chunk_received_at"`
    SessionReplayBudgetExceededAt types.String `tfsdk:"session_replay_budget_exceeded_at"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    IsArchived types.Bool `tfsdk:"is_archived"`
    ArchivedAt types.String `tfsdk:"archived_at"`
    ArchivedByUserId types.String `tfsdk:"archived_by_user_id"`
}

func (d *RumApplicationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_rum_application"
}

func (d *RumApplicationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Browser & mobile applications auto-discovered from OpenTelemetry RUM telemetry (browser.* / device.* resource attributes). One row per application, aggregating all end-user clients. Look up an existing rum application by `id`, or by any of its other arguments (`name`, `agent_version`, `app_identifier`, ...): each one set must match, and exactly one rum application may match them all.",

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
                MarkdownDescription: "ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Friendly name for this application.",
                Optional: true,
                Computed: true,
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Friendly globally unique name for your object.",
                Optional: true,
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Friendly description that will help you remember.",
                Optional: true,
                Computed: true,
            },
            "app_identifier": schema.StringAttribute{
                MarkdownDescription: "Stable identifier for this application from the service.name OpenTelemetry resource attribute. Identity key for this RUM application.",
                Optional: true,
                Computed: true,
            },
            "client_type": schema.StringAttribute{
                MarkdownDescription: "Whether this application's clients are browsers or mobile devices (browser / mobile), derived from browser.* / device.* attributes.",
                Optional: true,
                Computed: true,
            },
            "sdk_language": schema.StringAttribute{
                MarkdownDescription: "Last-seen telemetry.sdk.language resource attribute (e.g. webjs, swift, android). Used to scope this application's client telemetry apart from a same-named backend service.",
                Optional: true,
                Computed: true,
            },
            "otel_collector_status": schema.StringAttribute{
                MarkdownDescription: "Whether telemetry is currently being received (connected) or has gone stale (disconnected).",
                Optional: true,
                Computed: true,
            },
            "agent_version": schema.StringAttribute{
                MarkdownDescription: "Version of the OpenTelemetry SDK reporting this application.",
                Optional: true,
                Computed: true,
            },
            "last_seen_at": schema.StringAttribute{
                MarkdownDescription: "When telemetry was last received for this application.",
                Computed: true,
            },
            "labels": schema.SetAttribute{
                MarkdownDescription: "Relation to Labels Array where this object is categorized in. IDs of `oneuptime_label` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "retain_telemetry_data_for_days": schema.NumberAttribute{
                MarkdownDescription: "Number of days to retain telemetry data for this application. Leave blank to use the project-wide default.",
                Optional: true,
                Computed: true,
            },
            "telemetry_retention_config": schema.StringAttribute{
                MarkdownDescription: "Per-pillar retention overrides for this application. Unset fields fall back to the application default, then the project's retention settings. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "is_session_replay_enabled": schema.BoolAttribute{
                MarkdownDescription: "When enabled, the browser recorder may record and upload session replays for this application. On by default; Project.isSessionReplayAllowed must also be on. Turn it off here to stop recording for one application without affecting the rest of the project.",
                Optional: true,
                Computed: true,
            },
            "session_replay_masking_mode": schema.StringAttribute{
                MarkdownDescription: "How aggressively the recorder masks page content before it leaves the end user's device. MaskSensitiveInputsOnly (default) masks passwords and card / one-time-code fields and records everything else verbatim. MaskInputsOnly additionally masks every other input value. MaskAllText masks static page text too, producing a wireframe.",
                Optional: true,
                Computed: true,
            },
            "session_replay_mask_selectors": schema.StringAttribute{
                MarkdownDescription: "CSS selectors whose text content the recorder masks, in addition to whatever the masking mode already covers. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "session_replay_block_selectors": schema.StringAttribute{
                MarkdownDescription: "CSS selectors the recorder excludes from the DOM snapshot entirely, so the subtree is never captured rather than captured and masked. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "session_replay_ignore_error_patterns": schema.StringAttribute{
                MarkdownDescription: "Regex patterns matched against an uncaught error's message and source URL. Matching errors are still recorded in the session but no longer trigger an upload — the remedy for a chronically-throwing third-party tag that would otherwise convert error-triggered capture into always-on recording. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "session_replay_trace_propagation_origins": schema.StringAttribute{
                MarkdownDescription: "APIs on OTHER origins than the page (for example https://api.example.com) that the recorder may inject a W3C traceparent header into, linking recordings to the backend traces of their requests without any OpenTelemetry browser setup. Requests to the page's own origin need no entry here: Same-origin trace propagation covers them. Empty (the default) injects nothing cross-origin: adding a header makes a cross-origin request preflighted, so each listed origin is an explicit statement that its API allows traceparent in Access-Control-Allow-Headers. Listed origins get traceparent only, never the session id. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "session_replay_same_origin_trace_propagation": schema.BoolAttribute{
                MarkdownDescription: "When enabled, the recorder adds a W3C traceparent and a tracestate member carrying the replay session id (oneuptime=sid:<session id>) to the fetch and XHR requests the page makes to its own origin while a session is uploading, so the backend spans, logs and exceptions those requests cause link to the recording automatically, with no code in your frontend or backend. Nothing is added before consent, after consent is revoked, or before an on-error trigger fires, and a request that already carries a traceparent or tracestate keeps its own. A traceparent the recorder generates is marked sampled, so ParentBased samplers in your backend keep every browser-originated trace (to keep ratio sampling, set a remoteParentSampled ratio delegate only on the service(s) your pages call directly, never on the services they call: ratio decisions differ between language SDKs; for one rate across services, use tail sampling in an OpenTelemetry Collector). Your backend's OpenTelemetry forwards the tracestate, and with it the session id, to every service it calls, third parties included; the visitor id is never sent. On by default. Narrower create/update ACL than the other replay settings: it links recordings to backend telemetry that may name the user.",
                Optional: true,
                Computed: true,
            },
            "session_replay_lcp_budget_ms": schema.NumberAttribute{
                MarkdownDescription: "Largest Contentful Paint budget in milliseconds. A session whose LCP exceeds it uploads with the Performance trigger. 0 disables the trigger.",
                Optional: true,
                Computed: true,
            },
            "session_replay_long_task_budget_ms": schema.NumberAttribute{
                MarkdownDescription: "Main-thread long-task budget in milliseconds. A single task blocking longer than this uploads the session with the Performance trigger. 0 disables the trigger.",
                Optional: true,
                Computed: true,
            },
            "session_replay_slow_request_budget_ms": schema.NumberAttribute{
                MarkdownDescription: "Request duration budget in milliseconds. An instrumented request slower than this uploads the session with the Performance trigger. 0 disables the trigger.",
                Optional: true,
                Computed: true,
            },
            "session_replay_allowed_origins": schema.StringAttribute{
                MarkdownDescription: "Browser origins (scheme + host + port) and exact React Native app identities (app:// followed by the Android package or iOS bundle id) allowed to upload replays for this application. Empty (the default) accepts any sender. Web origins may use one leading host wildcard; app:// identities never allow wildcards. Once populated, browsers must send a listed Origin and native recorders without Origin must send a listed app identity. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "session_replay_consent_mode": schema.StringAttribute{
                MarkdownDescription: "NotRequired (default) uploads immediately, asserting a lawful basis that does not need a per-session grant. RequireExplicit buffers in memory and uploads nothing until the host page calls grantConsent(); set it if you need a per-session consent handshake, which most EU deployments will.",
                Optional: true,
                Computed: true,
            },
            "session_replay_capture_trigger": schema.StringAttribute{
                MarkdownDescription: "Always (default) uploads every sampled session from its first event, so an ordinary session is just as watchable as a broken one. OnErrorOrFrustration keeps a rolling in-memory buffer and uploads only when something actually went wrong, which costs roughly 15x less and stores far less end-user data.",
                Optional: true,
                Computed: true,
            },
            "session_replay_sample_percentage": schema.NumberAttribute{
                MarkdownDescription: "Percentage of sessions (0 to 100) eligible for recording. 100 by default, so with the default Always trigger every session is recorded. Lower it to cut storage and end-user data at rest; the decision is made once per session from a hash of the session id, so a session is never half-recorded.",
                Optional: true,
                Computed: true,
            },
            "session_replay_capture_user_identity": schema.BoolAttribute{
                MarkdownDescription: "When enabled, the end-user reference supplied by the host page is stored alongside the recording - as a one-way per-project HMAC for lookup and erasure, plus the raw reference behind its own narrower column ACL - so a support engineer can find the session a named customer is complaining about. When off, the reference is never attached to a recording and neither column is stored. (It is still sent once on the policy request, which is how targeted capture matches a named user; it is not persisted.) The reference must be supplied at load time - identify() called later reaches the server only on the session's final chunk, which the header is not rebuilt from. On by default. Narrower create/update ACL than the other replay settings: this is the switch that turns a pseudonymous recording into an identified one.",
                Optional: true,
                Computed: true,
            },
            "session_replay_capture_geo": schema.BoolAttribute{
                MarkdownDescription: "When enabled, a country code is derived from the request and stored on the session. On by default. The end user's IP address is never stored either way - the country is the only geographic fact this keeps.",
                Optional: true,
                Computed: true,
            },
            "session_replay_record_canvas": schema.BoolAttribute{
                MarkdownDescription: "When enabled, canvas contents are recorded. Off by default because canvas capture is expensive on the end user's device and canvases routinely render content the text masking cannot reach.",
                Optional: true,
                Computed: true,
            },
            "session_replay_retention_in_days": schema.NumberAttribute{
                MarkdownDescription: "How long session recordings are kept for this application. Clamped to 1, 7, 14, 30 or 90 days. Defaults to 7 rather than the 15 the other telemetry pillars use, because a short retention is itself a privacy control.",
                Optional: true,
                Computed: true,
            },
            "session_replay_monthly_budget_in_gb": schema.NumberAttribute{
                MarkdownDescription: "Optional ceiling on replay bytes ingested per calendar month for this application. Once exceeded, live recorders are told to stop. Leave blank for no application-level ceiling.",
                Optional: true,
                Computed: true,
            },
            "session_replay_last_chunk_received_at": schema.StringAttribute{
                MarkdownDescription: "When a session replay chunk was last accepted for this application.",
                Computed: true,
            },
            "session_replay_budget_exceeded_at": schema.StringAttribute{
                MarkdownDescription: "When the session replay byte budget was last hit for this application. Non-null means recorders are currently being told to stop.",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "is_archived": schema.BoolAttribute{
                MarkdownDescription: "Is this RUM application archived? Archived RUM applications are hidden from lists but keep collecting telemetry.",
                Optional: true,
                Computed: true,
            },
            "archived_at": schema.StringAttribute{
                MarkdownDescription: "When was this RUM application archived?",
                Computed: true,
            },
            "archived_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who archived this object (if this object was archived by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *RumApplicationDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *RumApplicationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data RumApplicationDataSourceModel

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
    if !data.Slug.IsNull() && !data.Slug.IsUnknown() {
        filters["slug"] = data.Slug.ValueString()
        filterNames = append(filterNames, "slug = "+fmt.Sprintf("%q", data.Slug.ValueString()))
    }
    if !data.Description.IsNull() && !data.Description.IsUnknown() {
        filters["description"] = data.Description.ValueString()
        filterNames = append(filterNames, "description = "+fmt.Sprintf("%q", data.Description.ValueString()))
    }
    if !data.AppIdentifier.IsNull() && !data.AppIdentifier.IsUnknown() {
        filters["appIdentifier"] = data.AppIdentifier.ValueString()
        filterNames = append(filterNames, "app_identifier = "+fmt.Sprintf("%q", data.AppIdentifier.ValueString()))
    }
    if !data.ClientType.IsNull() && !data.ClientType.IsUnknown() {
        filters["clientType"] = data.ClientType.ValueString()
        filterNames = append(filterNames, "client_type = "+fmt.Sprintf("%q", data.ClientType.ValueString()))
    }
    if !data.SdkLanguage.IsNull() && !data.SdkLanguage.IsUnknown() {
        filters["sdkLanguage"] = data.SdkLanguage.ValueString()
        filterNames = append(filterNames, "sdk_language = "+fmt.Sprintf("%q", data.SdkLanguage.ValueString()))
    }
    if !data.OtelCollectorStatus.IsNull() && !data.OtelCollectorStatus.IsUnknown() {
        filters["otelCollectorStatus"] = data.OtelCollectorStatus.ValueString()
        filterNames = append(filterNames, "otel_collector_status = "+fmt.Sprintf("%q", data.OtelCollectorStatus.ValueString()))
    }
    if !data.AgentVersion.IsNull() && !data.AgentVersion.IsUnknown() {
        filters["agentVersion"] = data.AgentVersion.ValueString()
        filterNames = append(filterNames, "agent_version = "+fmt.Sprintf("%q", data.AgentVersion.ValueString()))
    }
    if !data.RetainTelemetryDataForDays.IsNull() && !data.RetainTelemetryDataForDays.IsUnknown() {
        filters["retainTelemetryDataForDays"] = lookupNumber(data.RetainTelemetryDataForDays)
        filterNames = append(filterNames, "retain_telemetry_data_for_days = "+data.RetainTelemetryDataForDays.ValueBigFloat().String())
    }
    if !data.IsSessionReplayEnabled.IsNull() && !data.IsSessionReplayEnabled.IsUnknown() {
        filters["isSessionReplayEnabled"] = data.IsSessionReplayEnabled.ValueBool()
        filterNames = append(filterNames, "is_session_replay_enabled = "+fmt.Sprintf("%t", data.IsSessionReplayEnabled.ValueBool()))
    }
    if !data.SessionReplayMaskingMode.IsNull() && !data.SessionReplayMaskingMode.IsUnknown() {
        filters["sessionReplayMaskingMode"] = data.SessionReplayMaskingMode.ValueString()
        filterNames = append(filterNames, "session_replay_masking_mode = "+fmt.Sprintf("%q", data.SessionReplayMaskingMode.ValueString()))
    }
    if !data.SessionReplaySameOriginTracePropagation.IsNull() && !data.SessionReplaySameOriginTracePropagation.IsUnknown() {
        filters["sessionReplaySameOriginTracePropagation"] = data.SessionReplaySameOriginTracePropagation.ValueBool()
        filterNames = append(filterNames, "session_replay_same_origin_trace_propagation = "+fmt.Sprintf("%t", data.SessionReplaySameOriginTracePropagation.ValueBool()))
    }
    if !data.SessionReplayLcpBudgetMs.IsNull() && !data.SessionReplayLcpBudgetMs.IsUnknown() {
        filters["sessionReplayLcpBudgetMs"] = lookupNumber(data.SessionReplayLcpBudgetMs)
        filterNames = append(filterNames, "session_replay_lcp_budget_ms = "+data.SessionReplayLcpBudgetMs.ValueBigFloat().String())
    }
    if !data.SessionReplayLongTaskBudgetMs.IsNull() && !data.SessionReplayLongTaskBudgetMs.IsUnknown() {
        filters["sessionReplayLongTaskBudgetMs"] = lookupNumber(data.SessionReplayLongTaskBudgetMs)
        filterNames = append(filterNames, "session_replay_long_task_budget_ms = "+data.SessionReplayLongTaskBudgetMs.ValueBigFloat().String())
    }
    if !data.SessionReplaySlowRequestBudgetMs.IsNull() && !data.SessionReplaySlowRequestBudgetMs.IsUnknown() {
        filters["sessionReplaySlowRequestBudgetMs"] = lookupNumber(data.SessionReplaySlowRequestBudgetMs)
        filterNames = append(filterNames, "session_replay_slow_request_budget_ms = "+data.SessionReplaySlowRequestBudgetMs.ValueBigFloat().String())
    }
    if !data.SessionReplayConsentMode.IsNull() && !data.SessionReplayConsentMode.IsUnknown() {
        filters["sessionReplayConsentMode"] = data.SessionReplayConsentMode.ValueString()
        filterNames = append(filterNames, "session_replay_consent_mode = "+fmt.Sprintf("%q", data.SessionReplayConsentMode.ValueString()))
    }
    if !data.SessionReplayCaptureTrigger.IsNull() && !data.SessionReplayCaptureTrigger.IsUnknown() {
        filters["sessionReplayCaptureTrigger"] = data.SessionReplayCaptureTrigger.ValueString()
        filterNames = append(filterNames, "session_replay_capture_trigger = "+fmt.Sprintf("%q", data.SessionReplayCaptureTrigger.ValueString()))
    }
    if !data.SessionReplaySamplePercentage.IsNull() && !data.SessionReplaySamplePercentage.IsUnknown() {
        filters["sessionReplaySamplePercentage"] = lookupNumber(data.SessionReplaySamplePercentage)
        filterNames = append(filterNames, "session_replay_sample_percentage = "+data.SessionReplaySamplePercentage.ValueBigFloat().String())
    }
    if !data.SessionReplayCaptureUserIdentity.IsNull() && !data.SessionReplayCaptureUserIdentity.IsUnknown() {
        filters["sessionReplayCaptureUserIdentity"] = data.SessionReplayCaptureUserIdentity.ValueBool()
        filterNames = append(filterNames, "session_replay_capture_user_identity = "+fmt.Sprintf("%t", data.SessionReplayCaptureUserIdentity.ValueBool()))
    }
    if !data.SessionReplayCaptureGeo.IsNull() && !data.SessionReplayCaptureGeo.IsUnknown() {
        filters["sessionReplayCaptureGeo"] = data.SessionReplayCaptureGeo.ValueBool()
        filterNames = append(filterNames, "session_replay_capture_geo = "+fmt.Sprintf("%t", data.SessionReplayCaptureGeo.ValueBool()))
    }
    if !data.SessionReplayRecordCanvas.IsNull() && !data.SessionReplayRecordCanvas.IsUnknown() {
        filters["sessionReplayRecordCanvas"] = data.SessionReplayRecordCanvas.ValueBool()
        filterNames = append(filterNames, "session_replay_record_canvas = "+fmt.Sprintf("%t", data.SessionReplayRecordCanvas.ValueBool()))
    }
    if !data.SessionReplayRetentionInDays.IsNull() && !data.SessionReplayRetentionInDays.IsUnknown() {
        filters["sessionReplayRetentionInDays"] = lookupNumber(data.SessionReplayRetentionInDays)
        filterNames = append(filterNames, "session_replay_retention_in_days = "+data.SessionReplayRetentionInDays.ValueBigFloat().String())
    }
    if !data.SessionReplayMonthlyBudgetInGb.IsNull() && !data.SessionReplayMonthlyBudgetInGb.IsUnknown() {
        filters["sessionReplayMonthlyBudgetInGB"] = lookupNumber(data.SessionReplayMonthlyBudgetInGb)
        filterNames = append(filterNames, "session_replay_monthly_budget_in_gb = "+data.SessionReplayMonthlyBudgetInGb.ValueBigFloat().String())
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }
    if !data.IsArchived.IsNull() && !data.IsArchived.IsUnknown() {
        filters["isArchived"] = data.IsArchived.ValueBool()
        filterNames = append(filterNames, "is_archived = "+fmt.Sprintf("%t", data.IsArchived.ValueBool()))
    }
    if !data.ArchivedByUserId.IsNull() && !data.ArchivedByUserId.IsUnknown() {
        filters["archivedByUserId"] = data.ArchivedByUserId.ValueString()
        filterNames = append(filterNames, "archived_by_user_id = "+fmt.Sprintf("%q", data.ArchivedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the rum application up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the rum application up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "name": true,
        "slug": true,
        "description": true,
        "appIdentifier": true,
        "clientType": true,
        "sdkLanguage": true,
        "otelCollectorStatus": true,
        "agentVersion": true,
        "lastSeenAt": true,
        "labels": true,
        "retainTelemetryDataForDays": true,
        "telemetryRetentionConfig": true,
        "isSessionReplayEnabled": true,
        "sessionReplayMaskingMode": true,
        "sessionReplayMaskSelectors": true,
        "sessionReplayBlockSelectors": true,
        "sessionReplayIgnoreErrorPatterns": true,
        "sessionReplayTracePropagationOrigins": true,
        "sessionReplaySameOriginTracePropagation": true,
        "sessionReplayLcpBudgetMs": true,
        "sessionReplayLongTaskBudgetMs": true,
        "sessionReplaySlowRequestBudgetMs": true,
        "sessionReplayAllowedOrigins": true,
        "sessionReplayConsentMode": true,
        "sessionReplayCaptureTrigger": true,
        "sessionReplaySamplePercentage": true,
        "sessionReplayCaptureUserIdentity": true,
        "sessionReplayCaptureGeo": true,
        "sessionReplayRecordCanvas": true,
        "sessionReplayRetentionInDays": true,
        "sessionReplayMonthlyBudgetInGB": true,
        "sessionReplayLastChunkReceivedAt": true,
        "sessionReplayBudgetExceededAt": true,
        "createdByUserId": true,
        "isArchived": true,
        "archivedAt": true,
        "archivedByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/rum-application/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read rum_application, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No rum application found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read rum_application: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/rum-application/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list rum_application, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list rum_application: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No rum application matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one rum application matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for rum_application.")
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
    if obj, ok := item["slug"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Slug = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Slug = types.StringValue(string(jsonBytes))
        } else {
            data.Slug = types.StringNull()
        }
    } else if val, ok := item["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
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
    if obj, ok := item["appIdentifier"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AppIdentifier = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AppIdentifier = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AppIdentifier = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AppIdentifier = types.StringValue(string(jsonBytes))
        } else {
            data.AppIdentifier = types.StringNull()
        }
    } else if val, ok := item["appIdentifier"].(string); ok {
        data.AppIdentifier = types.StringValue(val)
    } else {
        data.AppIdentifier = types.StringNull()
    }
    if obj, ok := item["clientType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ClientType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ClientType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ClientType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ClientType = types.StringValue(string(jsonBytes))
        } else {
            data.ClientType = types.StringNull()
        }
    } else if val, ok := item["clientType"].(string); ok {
        data.ClientType = types.StringValue(val)
    } else {
        data.ClientType = types.StringNull()
    }
    if obj, ok := item["sdkLanguage"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SdkLanguage = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SdkLanguage = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SdkLanguage = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SdkLanguage = types.StringValue(string(jsonBytes))
        } else {
            data.SdkLanguage = types.StringNull()
        }
    } else if val, ok := item["sdkLanguage"].(string); ok {
        data.SdkLanguage = types.StringValue(val)
    } else {
        data.SdkLanguage = types.StringNull()
    }
    if obj, ok := item["otelCollectorStatus"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OtelCollectorStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.OtelCollectorStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.OtelCollectorStatus = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.OtelCollectorStatus = types.StringValue(string(jsonBytes))
        } else {
            data.OtelCollectorStatus = types.StringNull()
        }
    } else if val, ok := item["otelCollectorStatus"].(string); ok {
        data.OtelCollectorStatus = types.StringValue(val)
    } else {
        data.OtelCollectorStatus = types.StringNull()
    }
    if obj, ok := item["agentVersion"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AgentVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AgentVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AgentVersion = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AgentVersion = types.StringValue(string(jsonBytes))
        } else {
            data.AgentVersion = types.StringNull()
        }
    } else if val, ok := item["agentVersion"].(string); ok {
        data.AgentVersion = types.StringValue(val)
    } else {
        data.AgentVersion = types.StringNull()
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
    if val, ok := item["labels"].([]interface{}); ok {
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if jsonBytes, err := json.Marshal(itemMap); err == nil {
                    setItems = append(setItems, types.StringValue(string(jsonBytes)))
                }
            } else if str, ok := item.(string); ok {
                setItems = append(setItems, types.StringValue(str))
            } else {
                setItems = append(setItems, types.StringValue(fmt.Sprintf("%v", item)))
            }
        }
        sort.Slice(setItems, func(i, j int) bool {
            return setItems[i].(types.String).ValueString() < setItems[j].(types.String).ValueString()
        })
        data.Labels = types.SetValueMust(types.StringType, setItems)
    } else {
        data.Labels = types.SetNull(types.StringType)
    }
    if val, ok := item["retainTelemetryDataForDays"].(float64); ok {
        data.RetainTelemetryDataForDays = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["retainTelemetryDataForDays"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.RetainTelemetryDataForDays = types.NumberValue(big.NewFloat(val))
        } else {
            data.RetainTelemetryDataForDays = types.NumberNull()
        }
    } else {
        data.RetainTelemetryDataForDays = types.NumberNull()
    }
    if obj, ok := item["telemetryRetentionConfig"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TelemetryRetentionConfig = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TelemetryRetentionConfig = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TelemetryRetentionConfig = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TelemetryRetentionConfig = types.StringValue(string(jsonBytes))
        } else {
            data.TelemetryRetentionConfig = types.StringNull()
        }
    } else if val, ok := item["telemetryRetentionConfig"].(string); ok {
        data.TelemetryRetentionConfig = types.StringValue(val)
    } else {
        data.TelemetryRetentionConfig = types.StringNull()
    }
    if val, ok := item["isSessionReplayEnabled"].(bool); ok {
        data.IsSessionReplayEnabled = types.BoolValue(val)
    } else {
        data.IsSessionReplayEnabled = types.BoolNull()
    }
    if obj, ok := item["sessionReplayMaskingMode"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SessionReplayMaskingMode = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SessionReplayMaskingMode = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SessionReplayMaskingMode = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SessionReplayMaskingMode = types.StringValue(string(jsonBytes))
        } else {
            data.SessionReplayMaskingMode = types.StringNull()
        }
    } else if val, ok := item["sessionReplayMaskingMode"].(string); ok {
        data.SessionReplayMaskingMode = types.StringValue(val)
    } else {
        data.SessionReplayMaskingMode = types.StringNull()
    }
    if obj, ok := item["sessionReplayMaskSelectors"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SessionReplayMaskSelectors = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SessionReplayMaskSelectors = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SessionReplayMaskSelectors = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SessionReplayMaskSelectors = types.StringValue(string(jsonBytes))
        } else {
            data.SessionReplayMaskSelectors = types.StringNull()
        }
    } else if val, ok := item["sessionReplayMaskSelectors"].(string); ok {
        data.SessionReplayMaskSelectors = types.StringValue(val)
    } else {
        data.SessionReplayMaskSelectors = types.StringNull()
    }
    if obj, ok := item["sessionReplayBlockSelectors"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SessionReplayBlockSelectors = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SessionReplayBlockSelectors = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SessionReplayBlockSelectors = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SessionReplayBlockSelectors = types.StringValue(string(jsonBytes))
        } else {
            data.SessionReplayBlockSelectors = types.StringNull()
        }
    } else if val, ok := item["sessionReplayBlockSelectors"].(string); ok {
        data.SessionReplayBlockSelectors = types.StringValue(val)
    } else {
        data.SessionReplayBlockSelectors = types.StringNull()
    }
    if obj, ok := item["sessionReplayIgnoreErrorPatterns"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SessionReplayIgnoreErrorPatterns = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SessionReplayIgnoreErrorPatterns = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SessionReplayIgnoreErrorPatterns = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SessionReplayIgnoreErrorPatterns = types.StringValue(string(jsonBytes))
        } else {
            data.SessionReplayIgnoreErrorPatterns = types.StringNull()
        }
    } else if val, ok := item["sessionReplayIgnoreErrorPatterns"].(string); ok {
        data.SessionReplayIgnoreErrorPatterns = types.StringValue(val)
    } else {
        data.SessionReplayIgnoreErrorPatterns = types.StringNull()
    }
    if obj, ok := item["sessionReplayTracePropagationOrigins"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SessionReplayTracePropagationOrigins = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SessionReplayTracePropagationOrigins = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SessionReplayTracePropagationOrigins = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SessionReplayTracePropagationOrigins = types.StringValue(string(jsonBytes))
        } else {
            data.SessionReplayTracePropagationOrigins = types.StringNull()
        }
    } else if val, ok := item["sessionReplayTracePropagationOrigins"].(string); ok {
        data.SessionReplayTracePropagationOrigins = types.StringValue(val)
    } else {
        data.SessionReplayTracePropagationOrigins = types.StringNull()
    }
    if val, ok := item["sessionReplaySameOriginTracePropagation"].(bool); ok {
        data.SessionReplaySameOriginTracePropagation = types.BoolValue(val)
    } else {
        data.SessionReplaySameOriginTracePropagation = types.BoolNull()
    }
    if val, ok := item["sessionReplayLcpBudgetMs"].(float64); ok {
        data.SessionReplayLcpBudgetMs = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["sessionReplayLcpBudgetMs"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.SessionReplayLcpBudgetMs = types.NumberValue(big.NewFloat(val))
        } else {
            data.SessionReplayLcpBudgetMs = types.NumberNull()
        }
    } else {
        data.SessionReplayLcpBudgetMs = types.NumberNull()
    }
    if val, ok := item["sessionReplayLongTaskBudgetMs"].(float64); ok {
        data.SessionReplayLongTaskBudgetMs = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["sessionReplayLongTaskBudgetMs"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.SessionReplayLongTaskBudgetMs = types.NumberValue(big.NewFloat(val))
        } else {
            data.SessionReplayLongTaskBudgetMs = types.NumberNull()
        }
    } else {
        data.SessionReplayLongTaskBudgetMs = types.NumberNull()
    }
    if val, ok := item["sessionReplaySlowRequestBudgetMs"].(float64); ok {
        data.SessionReplaySlowRequestBudgetMs = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["sessionReplaySlowRequestBudgetMs"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.SessionReplaySlowRequestBudgetMs = types.NumberValue(big.NewFloat(val))
        } else {
            data.SessionReplaySlowRequestBudgetMs = types.NumberNull()
        }
    } else {
        data.SessionReplaySlowRequestBudgetMs = types.NumberNull()
    }
    if obj, ok := item["sessionReplayAllowedOrigins"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SessionReplayAllowedOrigins = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SessionReplayAllowedOrigins = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SessionReplayAllowedOrigins = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SessionReplayAllowedOrigins = types.StringValue(string(jsonBytes))
        } else {
            data.SessionReplayAllowedOrigins = types.StringNull()
        }
    } else if val, ok := item["sessionReplayAllowedOrigins"].(string); ok {
        data.SessionReplayAllowedOrigins = types.StringValue(val)
    } else {
        data.SessionReplayAllowedOrigins = types.StringNull()
    }
    if obj, ok := item["sessionReplayConsentMode"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SessionReplayConsentMode = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SessionReplayConsentMode = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SessionReplayConsentMode = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SessionReplayConsentMode = types.StringValue(string(jsonBytes))
        } else {
            data.SessionReplayConsentMode = types.StringNull()
        }
    } else if val, ok := item["sessionReplayConsentMode"].(string); ok {
        data.SessionReplayConsentMode = types.StringValue(val)
    } else {
        data.SessionReplayConsentMode = types.StringNull()
    }
    if obj, ok := item["sessionReplayCaptureTrigger"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SessionReplayCaptureTrigger = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SessionReplayCaptureTrigger = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SessionReplayCaptureTrigger = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SessionReplayCaptureTrigger = types.StringValue(string(jsonBytes))
        } else {
            data.SessionReplayCaptureTrigger = types.StringNull()
        }
    } else if val, ok := item["sessionReplayCaptureTrigger"].(string); ok {
        data.SessionReplayCaptureTrigger = types.StringValue(val)
    } else {
        data.SessionReplayCaptureTrigger = types.StringNull()
    }
    if val, ok := item["sessionReplaySamplePercentage"].(float64); ok {
        data.SessionReplaySamplePercentage = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["sessionReplaySamplePercentage"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.SessionReplaySamplePercentage = types.NumberValue(big.NewFloat(val))
        } else {
            data.SessionReplaySamplePercentage = types.NumberNull()
        }
    } else {
        data.SessionReplaySamplePercentage = types.NumberNull()
    }
    if val, ok := item["sessionReplayCaptureUserIdentity"].(bool); ok {
        data.SessionReplayCaptureUserIdentity = types.BoolValue(val)
    } else {
        data.SessionReplayCaptureUserIdentity = types.BoolNull()
    }
    if val, ok := item["sessionReplayCaptureGeo"].(bool); ok {
        data.SessionReplayCaptureGeo = types.BoolValue(val)
    } else {
        data.SessionReplayCaptureGeo = types.BoolNull()
    }
    if val, ok := item["sessionReplayRecordCanvas"].(bool); ok {
        data.SessionReplayRecordCanvas = types.BoolValue(val)
    } else {
        data.SessionReplayRecordCanvas = types.BoolNull()
    }
    if val, ok := item["sessionReplayRetentionInDays"].(float64); ok {
        data.SessionReplayRetentionInDays = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["sessionReplayRetentionInDays"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.SessionReplayRetentionInDays = types.NumberValue(big.NewFloat(val))
        } else {
            data.SessionReplayRetentionInDays = types.NumberNull()
        }
    } else {
        data.SessionReplayRetentionInDays = types.NumberNull()
    }
    if val, ok := item["sessionReplayMonthlyBudgetInGB"].(float64); ok {
        data.SessionReplayMonthlyBudgetInGb = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["sessionReplayMonthlyBudgetInGB"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.SessionReplayMonthlyBudgetInGb = types.NumberValue(big.NewFloat(val))
        } else {
            data.SessionReplayMonthlyBudgetInGb = types.NumberNull()
        }
    } else {
        data.SessionReplayMonthlyBudgetInGb = types.NumberNull()
    }
    if obj, ok := item["sessionReplayLastChunkReceivedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SessionReplayLastChunkReceivedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SessionReplayLastChunkReceivedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SessionReplayLastChunkReceivedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SessionReplayLastChunkReceivedAt = types.StringValue(string(jsonBytes))
        } else {
            data.SessionReplayLastChunkReceivedAt = types.StringNull()
        }
    } else if val, ok := item["sessionReplayLastChunkReceivedAt"].(string); ok {
        data.SessionReplayLastChunkReceivedAt = types.StringValue(val)
    } else {
        data.SessionReplayLastChunkReceivedAt = types.StringNull()
    }
    if obj, ok := item["sessionReplayBudgetExceededAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SessionReplayBudgetExceededAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SessionReplayBudgetExceededAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SessionReplayBudgetExceededAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SessionReplayBudgetExceededAt = types.StringValue(string(jsonBytes))
        } else {
            data.SessionReplayBudgetExceededAt = types.StringNull()
        }
    } else if val, ok := item["sessionReplayBudgetExceededAt"].(string); ok {
        data.SessionReplayBudgetExceededAt = types.StringValue(val)
    } else {
        data.SessionReplayBudgetExceededAt = types.StringNull()
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
    if val, ok := item["isArchived"].(bool); ok {
        data.IsArchived = types.BoolValue(val)
    } else {
        data.IsArchived = types.BoolNull()
    }
    if obj, ok := item["archivedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ArchivedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ArchivedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ArchivedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ArchivedAt = types.StringValue(string(jsonBytes))
        } else {
            data.ArchivedAt = types.StringNull()
        }
    } else if val, ok := item["archivedAt"].(string); ok {
        data.ArchivedAt = types.StringValue(val)
    } else {
        data.ArchivedAt = types.StringNull()
    }
    if obj, ok := item["archivedByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ArchivedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ArchivedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ArchivedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ArchivedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.ArchivedByUserId = types.StringNull()
        }
    } else if val, ok := item["archivedByUserId"].(string); ok {
        data.ArchivedByUserId = types.StringValue(val)
    } else {
        data.ArchivedByUserId = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
