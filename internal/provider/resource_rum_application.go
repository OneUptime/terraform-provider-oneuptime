package provider

import (
    "context"
    "fmt"
    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-framework/types/basetypes"
    "github.com/hashicorp/terraform-plugin-log/tflog"
    "math/big"
    "net/http"
    "github.com/hashicorp/terraform-plugin-framework/path"
    "encoding/json"
    "net/url"
    "strings"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/numberdefault"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
    "github.com/hashicorp/terraform-plugin-framework/attr"
    "sort"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/numberplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &RumApplicationResource{}
var _ resource.ResourceWithImportState = &RumApplicationResource{}

func NewRumApplicationResource() resource.Resource {
    return &RumApplicationResource{}
}

// RumApplicationResource defines the resource implementation.
type RumApplicationResource struct {
    client *Client
}

// RumApplicationResourceModel describes the resource data model.
type RumApplicationResourceModel struct {
    Id types.String `tfsdk:"id"`
    ProjectId types.String `tfsdk:"project_id"`
    Name types.String `tfsdk:"name"`
    Description types.String `tfsdk:"description"`
    AppIdentifier types.String `tfsdk:"app_identifier"`
    Labels types.Set `tfsdk:"labels"`
    RetainTelemetryDataForDays types.Number `tfsdk:"retain_telemetry_data_for_days"`
    TelemetryRetentionConfig JSONSubsetValue `tfsdk:"telemetry_retention_config"`
    IsSessionReplayEnabled types.Bool `tfsdk:"is_session_replay_enabled"`
    SessionReplayMaskingMode types.String `tfsdk:"session_replay_masking_mode"`
    SessionReplayMaskSelectors JSONSubsetValue `tfsdk:"session_replay_mask_selectors"`
    SessionReplayBlockSelectors JSONSubsetValue `tfsdk:"session_replay_block_selectors"`
    SessionReplayIgnoreErrorPatterns JSONSubsetValue `tfsdk:"session_replay_ignore_error_patterns"`
    SessionReplayTracePropagationOrigins JSONSubsetValue `tfsdk:"session_replay_trace_propagation_origins"`
    SessionReplayLcpBudgetMs types.Number `tfsdk:"session_replay_lcp_budget_ms"`
    SessionReplayLongTaskBudgetMs types.Number `tfsdk:"session_replay_long_task_budget_ms"`
    SessionReplaySlowRequestBudgetMs types.Number `tfsdk:"session_replay_slow_request_budget_ms"`
    SessionReplayAllowedOrigins JSONSubsetValue `tfsdk:"session_replay_allowed_origins"`
    SessionReplayConsentMode types.String `tfsdk:"session_replay_consent_mode"`
    SessionReplayCaptureTrigger types.String `tfsdk:"session_replay_capture_trigger"`
    SessionReplaySamplePercentage types.Number `tfsdk:"session_replay_sample_percentage"`
    SessionReplayCaptureUserIdentity types.Bool `tfsdk:"session_replay_capture_user_identity"`
    SessionReplayCaptureGeo types.Bool `tfsdk:"session_replay_capture_geo"`
    SessionReplayRecordCanvas types.Bool `tfsdk:"session_replay_record_canvas"`
    SessionReplayRetentionInDays types.Number `tfsdk:"session_replay_retention_in_days"`
    SessionReplayMonthlyBudgetInGb types.Number `tfsdk:"session_replay_monthly_budget_in_gb"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    IsArchived types.Bool `tfsdk:"is_archived"`
    CreatedAt RFC3339Value `tfsdk:"created_at"`
    UpdatedAt RFC3339Value `tfsdk:"updated_at"`
    DeletedAt RFC3339Value `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    Slug types.String `tfsdk:"slug"`
    ClientType types.String `tfsdk:"client_type"`
    SdkLanguage types.String `tfsdk:"sdk_language"`
    OtelCollectorStatus types.String `tfsdk:"otel_collector_status"`
    AgentVersion types.String `tfsdk:"agent_version"`
    LastSeenAt RFC3339Value `tfsdk:"last_seen_at"`
    SessionReplayLastChunkReceivedAt RFC3339Value `tfsdk:"session_replay_last_chunk_received_at"`
    SessionReplayBudgetExceededAt RFC3339Value `tfsdk:"session_replay_budget_exceeded_at"`
    ArchivedAt RFC3339Value `tfsdk:"archived_at"`
    ArchivedByUserId types.String `tfsdk:"archived_by_user_id"`
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
}

func (r *RumApplicationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_rum_application"
}

func (r *RumApplicationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Browser & mobile applications auto-discovered from OpenTelemetry RUM telemetry (browser.* / device.* resource attributes). One row per application, aggregating all end-user clients.",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Unique identifier for the resource",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Friendly name for this application.",
                Required: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Friendly description that will help you remember.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "app_identifier": schema.StringAttribute{
                MarkdownDescription: "Stable identifier for this application from the service.name OpenTelemetry resource attribute. Identity key for this RUM application..",
                Required: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.RequiresReplace(),
                },
            },
            "labels": schema.SetAttribute{
                MarkdownDescription: "Relation to Labels Array where this object is categorized in..",
                Optional: true,
                Computed: true,
                ElementType: types.StringType,
                PlanModifiers: []planmodifier.Set{
                    setplanmodifier.UseStateForUnknown(),
                },
            },
            "retain_telemetry_data_for_days": schema.NumberAttribute{
                MarkdownDescription: "Number of days to retain telemetry data for this application. Leave blank to use the project-wide default..",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "telemetry_retention_config": schema.StringAttribute{
                MarkdownDescription: "Per-pillar retention overrides for this application. Unset fields fall back to the application default, then the project's retention settings..",
                CustomType: JSONSubsetType{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: []validator.String{
                    JSONEnvelopeValidator(),
                },
            },
            "is_session_replay_enabled": schema.BoolAttribute{
                MarkdownDescription: "When enabled, the browser recorder may record and upload session replays for this application. On by default; Project.isSessionReplayAllowed must also be on. Turn it off here to stop recording for one application without affecting the rest of the project..",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(true),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "session_replay_masking_mode": schema.StringAttribute{
                MarkdownDescription: "How aggressively the recorder masks page content before it leaves the end user's device. MaskSensitiveInputsOnly (default) masks passwords and card / one-time-code fields and records everything else verbatim. MaskInputsOnly additionally masks every other input value. MaskAllText masks static page text too, producing a wireframe..",
                Optional: true,
                Computed: true,
                Default: stringdefault.StaticString("MaskSensitiveInputsOnly"),
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "session_replay_mask_selectors": schema.StringAttribute{
                MarkdownDescription: "CSS selectors whose text content the recorder masks, in addition to whatever the masking mode already covers..",
                CustomType: JSONSubsetType{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: []validator.String{
                    JSONEnvelopeValidator(),
                },
            },
            "session_replay_block_selectors": schema.StringAttribute{
                MarkdownDescription: "CSS selectors the recorder excludes from the DOM snapshot entirely, so the subtree is never captured rather than captured and masked..",
                CustomType: JSONSubsetType{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: []validator.String{
                    JSONEnvelopeValidator(),
                },
            },
            "session_replay_ignore_error_patterns": schema.StringAttribute{
                MarkdownDescription: "Regex patterns matched against an uncaught error's message and source URL. Matching errors are still recorded in the session but no longer trigger an upload — the remedy for a chronically-throwing third-party tag that would otherwise convert error-triggered capture into always-on recording..",
                CustomType: JSONSubsetType{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: []validator.String{
                    JSONEnvelopeValidator(),
                },
            },
            "session_replay_trace_propagation_origins": schema.StringAttribute{
                MarkdownDescription: "Origins the recorder may inject a W3C traceparent header into, linking recordings to the backend traces of their requests without any OpenTelemetry browser setup. Empty means never inject: adding a header makes cross-origin requests preflighted, so each listed origin is an explicit statement that its API allows the traceparent header..",
                CustomType: JSONSubsetType{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: []validator.String{
                    JSONEnvelopeValidator(),
                },
            },
            "session_replay_lcp_budget_ms": schema.NumberAttribute{
                MarkdownDescription: "Largest Contentful Paint budget in milliseconds. A session whose LCP exceeds it uploads with the Performance trigger. 0 disables the trigger..",
                Optional: true,
                Computed: true,
                Default: numberdefault.StaticBigFloat(big.NewFloat(0)),
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "session_replay_long_task_budget_ms": schema.NumberAttribute{
                MarkdownDescription: "Main-thread long-task budget in milliseconds. A single task blocking longer than this uploads the session with the Performance trigger. 0 disables the trigger..",
                Optional: true,
                Computed: true,
                Default: numberdefault.StaticBigFloat(big.NewFloat(0)),
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "session_replay_slow_request_budget_ms": schema.NumberAttribute{
                MarkdownDescription: "Request duration budget in milliseconds. An instrumented request slower than this uploads the session with the Performance trigger. 0 disables the trigger..",
                Optional: true,
                Computed: true,
                Default: numberdefault.StaticBigFloat(big.NewFloat(0)),
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "session_replay_allowed_origins": schema.StringAttribute{
                MarkdownDescription: "Exact browser origins (scheme + host + port) allowed to upload session replay chunks for this application. Empty (the default) accepts any origin. Once you list an origin this becomes a strict allowlist: anything unlisted, and any request with no Origin header, is refused..",
                CustomType: JSONSubsetType{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: []validator.String{
                    JSONEnvelopeValidator(),
                },
            },
            "session_replay_consent_mode": schema.StringAttribute{
                MarkdownDescription: "NotRequired (default) uploads immediately, asserting a lawful basis that does not need a per-session grant. RequireExplicit buffers in memory and uploads nothing until the host page calls grantConsent(); set it if you need a per-session consent handshake, which most EU deployments will..",
                Optional: true,
                Computed: true,
                Default: stringdefault.StaticString("NotRequired"),
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "session_replay_capture_trigger": schema.StringAttribute{
                MarkdownDescription: "OnErrorOrFrustration (default) keeps a rolling in-memory buffer and uploads only when something actually went wrong. Always uploads every sampled session from its first event, which costs materially more and stores materially more end-user data..",
                Optional: true,
                Computed: true,
                Default: stringdefault.StaticString("OnErrorOrFrustration"),
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "session_replay_sample_percentage": schema.NumberAttribute{
                MarkdownDescription: "Percentage of sessions (0 to 100) recorded regardless of whether anything went wrong. 0 by default, so with the default trigger only failing sessions are recorded..",
                Optional: true,
                Computed: true,
                Default: numberdefault.StaticBigFloat(big.NewFloat(0)),
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "session_replay_capture_user_identity": schema.BoolAttribute{
                MarkdownDescription: "When enabled, the end-user reference supplied by the host page is stored alongside the recording - as a one-way per-project HMAC for lookup and erasure, plus the raw reference behind its own narrower column ACL - so a support engineer can find the session a named customer is complaining about. When off, the reference is never attached to a recording and neither column is stored. (It is still sent once on the policy request, which is how targeted capture matches a named user; it is not persisted.) The reference must be supplied at load time - identify() called later reaches the server only on the session's final chunk, which the header is not rebuilt from. On by default. Narrower create/update ACL than the other replay settings: this is the switch that turns a pseudonymous recording into an identified one..",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(true),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "session_replay_capture_geo": schema.BoolAttribute{
                MarkdownDescription: "When enabled, a country code is derived from the request and stored on the session. On by default. The end user's IP address is never stored either way - the country is the only geographic fact this keeps..",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(true),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "session_replay_record_canvas": schema.BoolAttribute{
                MarkdownDescription: "When enabled, canvas contents are recorded. Off by default because canvas capture is expensive on the end user's device and canvases routinely render content the text masking cannot reach..",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "session_replay_retention_in_days": schema.NumberAttribute{
                MarkdownDescription: "How long session recordings are kept for this application. Clamped to 1, 7, 14, 30 or 90 days. Defaults to 7 rather than the 15 the other telemetry pillars use, because a short retention is itself a privacy control..",
                Optional: true,
                Computed: true,
                Default: numberdefault.StaticBigFloat(big.NewFloat(7)),
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "session_replay_monthly_budget_in_gb": schema.NumberAttribute{
                MarkdownDescription: "Optional ceiling on replay bytes ingested per calendar month for this application. Once exceeded, live recorders are told to stop. Leave blank for no application-level ceiling..",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                    stringplanmodifier.RequiresReplace(),
                },
            },
            "is_archived": schema.BoolAttribute{
                MarkdownDescription: "Is this RUM application archived? Archived RUM applications are hidden from lists but keep collecting telemetry..",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "created_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "updated_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "deleted_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "version": schema.NumberAttribute{
                MarkdownDescription: "Object version",
                Computed: true,
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Friendly globally unique name for your object.",
                Computed: true,
            },
            "client_type": schema.StringAttribute{
                MarkdownDescription: "Whether this application's clients are browsers or mobile devices (browser / mobile), derived from browser.* / device.* attributes..",
                Computed: true,
            },
            "sdk_language": schema.StringAttribute{
                MarkdownDescription: "Last-seen telemetry.sdk.language resource attribute (e.g. webjs, swift, android). Used to scope this application's client telemetry apart from a same-named backend service..",
                Computed: true,
            },
            "otel_collector_status": schema.StringAttribute{
                MarkdownDescription: "Whether telemetry is currently being received (connected) or has gone stale (disconnected)..",
                Computed: true,
            },
            "agent_version": schema.StringAttribute{
                MarkdownDescription: "Version of the OpenTelemetry SDK reporting this application..",
                Computed: true,
            },
            "last_seen_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "session_replay_last_chunk_received_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "session_replay_budget_exceeded_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "archived_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "archived_by_user_id": schema.StringAttribute{
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

func (r *RumApplicationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
    // Prevent panic if the provider has not been configured.
    if req.ProviderData == nil {
        return
    }

    client, ok := req.ProviderData.(*Client)

    if !ok {
        resp.Diagnostics.AddError(
            "Unexpected Resource Configure Type",
            fmt.Sprintf("Expected *Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
        )

        return
    }

    r.client = client
}


func (r *RumApplicationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var data RumApplicationResourceModel

    // Read Terraform plan data into the model
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }



    // Create API request body. Unset (null/unknown) optional fields are
    // omitted so server-side defaults apply instead of being overwritten
    // with zero values.
    rumApplicationRequest := map[string]interface{}{
        "data": map[string]interface{}{},
    }
    requestDataMap := rumApplicationRequest["data"].(map[string]interface{})

    if !data.Name.IsNull() && !data.Name.IsUnknown() {
        requestDataMap["name"] = data.Name.ValueString()
    }
    if !data.Description.IsNull() && !data.Description.IsUnknown() {
        requestDataMap["description"] = data.Description.ValueString()
    }
    if !data.AppIdentifier.IsNull() && !data.AppIdentifier.IsUnknown() {
        requestDataMap["appIdentifier"] = data.AppIdentifier.ValueString()
    }
    if !data.Labels.IsNull() && !data.Labels.IsUnknown() {
        requestDataMap["labels"] = r.convertTerraformSetToInterface(data.Labels)
    }
    if !data.RetainTelemetryDataForDays.IsNull() && !data.RetainTelemetryDataForDays.IsUnknown() {
        requestDataMap["retainTelemetryDataForDays"] = r.bigFloatToFloat64(data.RetainTelemetryDataForDays.ValueBigFloat())
    }
    if parsedTelemetryRetentionConfig := r.parseJSONField(data.TelemetryRetentionConfig); parsedTelemetryRetentionConfig != nil {
        requestDataMap["telemetryRetentionConfig"] = parsedTelemetryRetentionConfig
    }
    if !data.IsSessionReplayEnabled.IsNull() && !data.IsSessionReplayEnabled.IsUnknown() {
        requestDataMap["isSessionReplayEnabled"] = data.IsSessionReplayEnabled.ValueBool()
    }
    if !data.SessionReplayMaskingMode.IsNull() && !data.SessionReplayMaskingMode.IsUnknown() {
        requestDataMap["sessionReplayMaskingMode"] = data.SessionReplayMaskingMode.ValueString()
    }
    if parsedSessionReplayMaskSelectors := r.parseJSONField(data.SessionReplayMaskSelectors); parsedSessionReplayMaskSelectors != nil {
        requestDataMap["sessionReplayMaskSelectors"] = parsedSessionReplayMaskSelectors
    }
    if parsedSessionReplayBlockSelectors := r.parseJSONField(data.SessionReplayBlockSelectors); parsedSessionReplayBlockSelectors != nil {
        requestDataMap["sessionReplayBlockSelectors"] = parsedSessionReplayBlockSelectors
    }
    if parsedSessionReplayIgnoreErrorPatterns := r.parseJSONField(data.SessionReplayIgnoreErrorPatterns); parsedSessionReplayIgnoreErrorPatterns != nil {
        requestDataMap["sessionReplayIgnoreErrorPatterns"] = parsedSessionReplayIgnoreErrorPatterns
    }
    if parsedSessionReplayTracePropagationOrigins := r.parseJSONField(data.SessionReplayTracePropagationOrigins); parsedSessionReplayTracePropagationOrigins != nil {
        requestDataMap["sessionReplayTracePropagationOrigins"] = parsedSessionReplayTracePropagationOrigins
    }
    if !data.SessionReplayLcpBudgetMs.IsNull() && !data.SessionReplayLcpBudgetMs.IsUnknown() {
        requestDataMap["sessionReplayLcpBudgetMs"] = r.bigFloatToFloat64(data.SessionReplayLcpBudgetMs.ValueBigFloat())
    }
    if !data.SessionReplayLongTaskBudgetMs.IsNull() && !data.SessionReplayLongTaskBudgetMs.IsUnknown() {
        requestDataMap["sessionReplayLongTaskBudgetMs"] = r.bigFloatToFloat64(data.SessionReplayLongTaskBudgetMs.ValueBigFloat())
    }
    if !data.SessionReplaySlowRequestBudgetMs.IsNull() && !data.SessionReplaySlowRequestBudgetMs.IsUnknown() {
        requestDataMap["sessionReplaySlowRequestBudgetMs"] = r.bigFloatToFloat64(data.SessionReplaySlowRequestBudgetMs.ValueBigFloat())
    }
    if parsedSessionReplayAllowedOrigins := r.parseJSONField(data.SessionReplayAllowedOrigins); parsedSessionReplayAllowedOrigins != nil {
        requestDataMap["sessionReplayAllowedOrigins"] = parsedSessionReplayAllowedOrigins
    }
    if !data.SessionReplayConsentMode.IsNull() && !data.SessionReplayConsentMode.IsUnknown() {
        requestDataMap["sessionReplayConsentMode"] = data.SessionReplayConsentMode.ValueString()
    }
    if !data.SessionReplayCaptureTrigger.IsNull() && !data.SessionReplayCaptureTrigger.IsUnknown() {
        requestDataMap["sessionReplayCaptureTrigger"] = data.SessionReplayCaptureTrigger.ValueString()
    }
    if !data.SessionReplaySamplePercentage.IsNull() && !data.SessionReplaySamplePercentage.IsUnknown() {
        requestDataMap["sessionReplaySamplePercentage"] = r.bigFloatToFloat64(data.SessionReplaySamplePercentage.ValueBigFloat())
    }
    if !data.SessionReplayCaptureUserIdentity.IsNull() && !data.SessionReplayCaptureUserIdentity.IsUnknown() {
        requestDataMap["sessionReplayCaptureUserIdentity"] = data.SessionReplayCaptureUserIdentity.ValueBool()
    }
    if !data.SessionReplayCaptureGeo.IsNull() && !data.SessionReplayCaptureGeo.IsUnknown() {
        requestDataMap["sessionReplayCaptureGeo"] = data.SessionReplayCaptureGeo.ValueBool()
    }
    if !data.SessionReplayRecordCanvas.IsNull() && !data.SessionReplayRecordCanvas.IsUnknown() {
        requestDataMap["sessionReplayRecordCanvas"] = data.SessionReplayRecordCanvas.ValueBool()
    }
    if !data.SessionReplayRetentionInDays.IsNull() && !data.SessionReplayRetentionInDays.IsUnknown() {
        requestDataMap["sessionReplayRetentionInDays"] = r.bigFloatToFloat64(data.SessionReplayRetentionInDays.ValueBigFloat())
    }
    if !data.SessionReplayMonthlyBudgetInGb.IsNull() && !data.SessionReplayMonthlyBudgetInGb.IsUnknown() {
        requestDataMap["sessionReplayMonthlyBudgetInGB"] = r.bigFloatToFloat64(data.SessionReplayMonthlyBudgetInGb.ValueBigFloat())
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        requestDataMap["createdByUserId"] = data.CreatedByUserId.ValueString()
    }
    if !data.IsArchived.IsNull() && !data.IsArchived.IsUnknown() {
        requestDataMap["isArchived"] = data.IsArchived.ValueBool()
    }

    // Make API call
    httpResp, err := r.client.Post(ctx, "/rum-application", rumApplicationRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create rum_application, got error: %s", err))
        return
    }

    var rumApplicationResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &rumApplicationResponse)
    if err != nil {
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to create rum_application: %s", err))
        return
    }

    // Extract the new resource id from the create response.
    createdId := ""
    if wrapper, ok := rumApplicationResponse["data"].(map[string]interface{}); ok {
        if val, ok := wrapper["_id"].(string); ok {
            createdId = val
        }
    } else if val, ok := rumApplicationResponse["_id"].(string); ok {
        createdId = val
    }
    if createdId == "" {
        resp.Diagnostics.AddError("OneUptime API Error", "Create response for rum_application did not contain an id. This is a bug in the provider or the API; please report it.")
        return
    }
    data.Id = types.StringValue(createdId)

    /*
     * The server has committed the row. Persist what we know to state BEFORE
     * the read-back: if the read-back fails and we return without setting
     * state, Terraform never learns the resource exists and the created
     * rum_application is orphaned server-side — never refreshed, never
     * destroyed. Delete already refuses to drop state on failure for the
     * same reason; Create must not either.
     */
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Re-read the resource so state reflects server-normalized values.
    selectParam := map[string]interface{}{
        "projectId": true,
        "name": true,
        "description": true,
        "appIdentifier": true,
        "labels": true,
        "retainTelemetryDataForDays": true,
        "telemetryRetentionConfig": true,
        "isSessionReplayEnabled": true,
        "sessionReplayMaskingMode": true,
        "sessionReplayMaskSelectors": true,
        "sessionReplayBlockSelectors": true,
        "sessionReplayIgnoreErrorPatterns": true,
        "sessionReplayTracePropagationOrigins": true,
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
        "createdByUserId": true,
        "isArchived": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "slug": true,
        "clientType": true,
        "sdkLanguage": true,
        "otelCollectorStatus": true,
        "agentVersion": true,
        "lastSeenAt": true,
        "sessionReplayLastChunkReceivedAt": true,
        "sessionReplayBudgetExceededAt": true,
        "archivedAt": true,
        "archivedByUserId": true,
        "deletedByUserId": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect(ctx, "/rum-application/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        /*
         * State already owns the id, so the resource is tracked and the next
         * refresh reconciles the remaining attributes. Warn rather than
         * error: erroring here would strand a real resource.
         */
        resp.Diagnostics.AddWarning("Read After Create Failed", fmt.Sprintf("Created rum_application but could not read it back; state is incomplete until the next refresh: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddWarning("Read After Create Failed", fmt.Sprintf("Created rum_application but could not parse the read-back response; state is incomplete until the next refresh: %s", err))
        return
    }

    // Update the model with the authoritative read response
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := readResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = readResponse
    }

    if obj, ok := dataMap["projectId"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok {
            data.ProjectId = types.StringValue(val)
        } else {
            data.ProjectId = types.StringNull()
        }
    } else if val, ok := dataMap["projectId"].(string); ok {
        data.ProjectId = types.StringValue(val)
    } else {
        data.ProjectId = types.StringNull()
    }
    if obj, ok := dataMap["name"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Name = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Name = types.StringValue(string(jsonBytes))
            } else {
                data.Name = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Name = types.StringValue(string(jsonBytes))
            } else {
                data.Name = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Name = types.StringValue(string(jsonBytes))
        } else {
            data.Name = types.StringNull()
        }
    } else if val, ok := dataMap["name"].(string); ok {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if obj, ok := dataMap["description"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Description = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Description = types.StringValue(string(jsonBytes))
            } else {
                data.Description = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Description = types.StringValue(string(jsonBytes))
            } else {
                data.Description = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Description = types.StringValue(string(jsonBytes))
        } else {
            data.Description = types.StringNull()
        }
    } else if val, ok := dataMap["description"].(string); ok {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if obj, ok := dataMap["appIdentifier"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AppIdentifier = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AppIdentifier = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AppIdentifier = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AppIdentifier = types.StringValue(string(jsonBytes))
            } else {
                data.AppIdentifier = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AppIdentifier = types.StringValue(string(jsonBytes))
            } else {
                data.AppIdentifier = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AppIdentifier = types.StringValue(string(jsonBytes))
        } else {
            data.AppIdentifier = types.StringNull()
        }
    } else if val, ok := dataMap["appIdentifier"].(string); ok {
        data.AppIdentifier = types.StringValue(val)
    } else {
        data.AppIdentifier = types.StringNull()
    }
    if val, ok := dataMap["labels"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.Labels = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.Labels = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if val, ok := dataMap["retainTelemetryDataForDays"].(float64); ok {
        data.RetainTelemetryDataForDays = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["retainTelemetryDataForDays"].(int); ok {
        data.RetainTelemetryDataForDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["retainTelemetryDataForDays"].(int64); ok {
        data.RetainTelemetryDataForDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["retainTelemetryDataForDays"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.RetainTelemetryDataForDays = types.NumberValue(big.NewFloat(val))
        } else {
            data.RetainTelemetryDataForDays = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.RetainTelemetryDataForDays = types.NumberNull()
    }
    if obj, ok := dataMap["telemetryRetentionConfig"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TelemetryRetentionConfig = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.TelemetryRetentionConfig = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.TelemetryRetentionConfig = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.TelemetryRetentionConfig = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.TelemetryRetentionConfig = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.TelemetryRetentionConfig = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.TelemetryRetentionConfig = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.TelemetryRetentionConfig = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.TelemetryRetentionConfig = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["telemetryRetentionConfig"].(string); ok {
        data.TelemetryRetentionConfig = NewJSONSubsetValue(val)
    } else {
        data.TelemetryRetentionConfig = NewJSONSubsetNull()
    }
    if val, ok := dataMap["isSessionReplayEnabled"].(bool); ok {
        data.IsSessionReplayEnabled = types.BoolValue(val)
    }
    if obj, ok := dataMap["sessionReplayMaskingMode"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SessionReplayMaskingMode = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SessionReplayMaskingMode = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SessionReplayMaskingMode = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SessionReplayMaskingMode = types.StringValue(string(jsonBytes))
            } else {
                data.SessionReplayMaskingMode = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SessionReplayMaskingMode = types.StringValue(string(jsonBytes))
            } else {
                data.SessionReplayMaskingMode = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SessionReplayMaskingMode = types.StringValue(string(jsonBytes))
        } else {
            data.SessionReplayMaskingMode = types.StringNull()
        }
    } else if val, ok := dataMap["sessionReplayMaskingMode"].(string); ok {
        data.SessionReplayMaskingMode = types.StringValue(val)
    } else {
        data.SessionReplayMaskingMode = types.StringNull()
    }
    if obj, ok := dataMap["sessionReplayMaskSelectors"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SessionReplayMaskSelectors = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SessionReplayMaskSelectors = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SessionReplayMaskSelectors = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SessionReplayMaskSelectors = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.SessionReplayMaskSelectors = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SessionReplayMaskSelectors = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.SessionReplayMaskSelectors = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SessionReplayMaskSelectors = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.SessionReplayMaskSelectors = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["sessionReplayMaskSelectors"].(string); ok {
        data.SessionReplayMaskSelectors = NewJSONSubsetValue(val)
    } else {
        data.SessionReplayMaskSelectors = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["sessionReplayBlockSelectors"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SessionReplayBlockSelectors = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SessionReplayBlockSelectors = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SessionReplayBlockSelectors = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SessionReplayBlockSelectors = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.SessionReplayBlockSelectors = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SessionReplayBlockSelectors = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.SessionReplayBlockSelectors = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SessionReplayBlockSelectors = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.SessionReplayBlockSelectors = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["sessionReplayBlockSelectors"].(string); ok {
        data.SessionReplayBlockSelectors = NewJSONSubsetValue(val)
    } else {
        data.SessionReplayBlockSelectors = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["sessionReplayIgnoreErrorPatterns"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SessionReplayIgnoreErrorPatterns = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SessionReplayIgnoreErrorPatterns = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SessionReplayIgnoreErrorPatterns = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SessionReplayIgnoreErrorPatterns = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.SessionReplayIgnoreErrorPatterns = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SessionReplayIgnoreErrorPatterns = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.SessionReplayIgnoreErrorPatterns = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SessionReplayIgnoreErrorPatterns = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.SessionReplayIgnoreErrorPatterns = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["sessionReplayIgnoreErrorPatterns"].(string); ok {
        data.SessionReplayIgnoreErrorPatterns = NewJSONSubsetValue(val)
    } else {
        data.SessionReplayIgnoreErrorPatterns = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["sessionReplayTracePropagationOrigins"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SessionReplayTracePropagationOrigins = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SessionReplayTracePropagationOrigins = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SessionReplayTracePropagationOrigins = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SessionReplayTracePropagationOrigins = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.SessionReplayTracePropagationOrigins = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SessionReplayTracePropagationOrigins = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.SessionReplayTracePropagationOrigins = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SessionReplayTracePropagationOrigins = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.SessionReplayTracePropagationOrigins = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["sessionReplayTracePropagationOrigins"].(string); ok {
        data.SessionReplayTracePropagationOrigins = NewJSONSubsetValue(val)
    } else {
        data.SessionReplayTracePropagationOrigins = NewJSONSubsetNull()
    }
    if val, ok := dataMap["sessionReplayLcpBudgetMs"].(float64); ok {
        data.SessionReplayLcpBudgetMs = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["sessionReplayLcpBudgetMs"].(int); ok {
        data.SessionReplayLcpBudgetMs = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["sessionReplayLcpBudgetMs"].(int64); ok {
        data.SessionReplayLcpBudgetMs = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["sessionReplayLcpBudgetMs"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.SessionReplayLcpBudgetMs = types.NumberValue(big.NewFloat(val))
        } else {
            data.SessionReplayLcpBudgetMs = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.SessionReplayLcpBudgetMs = types.NumberNull()
    }
    if val, ok := dataMap["sessionReplayLongTaskBudgetMs"].(float64); ok {
        data.SessionReplayLongTaskBudgetMs = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["sessionReplayLongTaskBudgetMs"].(int); ok {
        data.SessionReplayLongTaskBudgetMs = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["sessionReplayLongTaskBudgetMs"].(int64); ok {
        data.SessionReplayLongTaskBudgetMs = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["sessionReplayLongTaskBudgetMs"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.SessionReplayLongTaskBudgetMs = types.NumberValue(big.NewFloat(val))
        } else {
            data.SessionReplayLongTaskBudgetMs = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.SessionReplayLongTaskBudgetMs = types.NumberNull()
    }
    if val, ok := dataMap["sessionReplaySlowRequestBudgetMs"].(float64); ok {
        data.SessionReplaySlowRequestBudgetMs = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["sessionReplaySlowRequestBudgetMs"].(int); ok {
        data.SessionReplaySlowRequestBudgetMs = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["sessionReplaySlowRequestBudgetMs"].(int64); ok {
        data.SessionReplaySlowRequestBudgetMs = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["sessionReplaySlowRequestBudgetMs"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.SessionReplaySlowRequestBudgetMs = types.NumberValue(big.NewFloat(val))
        } else {
            data.SessionReplaySlowRequestBudgetMs = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.SessionReplaySlowRequestBudgetMs = types.NumberNull()
    }
    if obj, ok := dataMap["sessionReplayAllowedOrigins"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SessionReplayAllowedOrigins = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SessionReplayAllowedOrigins = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SessionReplayAllowedOrigins = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SessionReplayAllowedOrigins = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.SessionReplayAllowedOrigins = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SessionReplayAllowedOrigins = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.SessionReplayAllowedOrigins = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SessionReplayAllowedOrigins = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.SessionReplayAllowedOrigins = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["sessionReplayAllowedOrigins"].(string); ok {
        data.SessionReplayAllowedOrigins = NewJSONSubsetValue(val)
    } else {
        data.SessionReplayAllowedOrigins = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["sessionReplayConsentMode"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SessionReplayConsentMode = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SessionReplayConsentMode = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SessionReplayConsentMode = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SessionReplayConsentMode = types.StringValue(string(jsonBytes))
            } else {
                data.SessionReplayConsentMode = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SessionReplayConsentMode = types.StringValue(string(jsonBytes))
            } else {
                data.SessionReplayConsentMode = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SessionReplayConsentMode = types.StringValue(string(jsonBytes))
        } else {
            data.SessionReplayConsentMode = types.StringNull()
        }
    } else if val, ok := dataMap["sessionReplayConsentMode"].(string); ok {
        data.SessionReplayConsentMode = types.StringValue(val)
    } else {
        data.SessionReplayConsentMode = types.StringNull()
    }
    if obj, ok := dataMap["sessionReplayCaptureTrigger"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SessionReplayCaptureTrigger = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SessionReplayCaptureTrigger = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SessionReplayCaptureTrigger = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SessionReplayCaptureTrigger = types.StringValue(string(jsonBytes))
            } else {
                data.SessionReplayCaptureTrigger = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SessionReplayCaptureTrigger = types.StringValue(string(jsonBytes))
            } else {
                data.SessionReplayCaptureTrigger = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SessionReplayCaptureTrigger = types.StringValue(string(jsonBytes))
        } else {
            data.SessionReplayCaptureTrigger = types.StringNull()
        }
    } else if val, ok := dataMap["sessionReplayCaptureTrigger"].(string); ok {
        data.SessionReplayCaptureTrigger = types.StringValue(val)
    } else {
        data.SessionReplayCaptureTrigger = types.StringNull()
    }
    if val, ok := dataMap["sessionReplaySamplePercentage"].(float64); ok {
        data.SessionReplaySamplePercentage = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["sessionReplaySamplePercentage"].(int); ok {
        data.SessionReplaySamplePercentage = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["sessionReplaySamplePercentage"].(int64); ok {
        data.SessionReplaySamplePercentage = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["sessionReplaySamplePercentage"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.SessionReplaySamplePercentage = types.NumberValue(big.NewFloat(val))
        } else {
            data.SessionReplaySamplePercentage = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.SessionReplaySamplePercentage = types.NumberNull()
    }
    if val, ok := dataMap["sessionReplayCaptureUserIdentity"].(bool); ok {
        data.SessionReplayCaptureUserIdentity = types.BoolValue(val)
    }
    if val, ok := dataMap["sessionReplayCaptureGeo"].(bool); ok {
        data.SessionReplayCaptureGeo = types.BoolValue(val)
    }
    if val, ok := dataMap["sessionReplayRecordCanvas"].(bool); ok {
        data.SessionReplayRecordCanvas = types.BoolValue(val)
    }
    if val, ok := dataMap["sessionReplayRetentionInDays"].(float64); ok {
        data.SessionReplayRetentionInDays = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["sessionReplayRetentionInDays"].(int); ok {
        data.SessionReplayRetentionInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["sessionReplayRetentionInDays"].(int64); ok {
        data.SessionReplayRetentionInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["sessionReplayRetentionInDays"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.SessionReplayRetentionInDays = types.NumberValue(big.NewFloat(val))
        } else {
            data.SessionReplayRetentionInDays = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.SessionReplayRetentionInDays = types.NumberNull()
    }
    if val, ok := dataMap["sessionReplayMonthlyBudgetInGB"].(float64); ok {
        data.SessionReplayMonthlyBudgetInGb = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["sessionReplayMonthlyBudgetInGB"].(int); ok {
        data.SessionReplayMonthlyBudgetInGb = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["sessionReplayMonthlyBudgetInGB"].(int64); ok {
        data.SessionReplayMonthlyBudgetInGb = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["sessionReplayMonthlyBudgetInGB"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.SessionReplayMonthlyBudgetInGb = types.NumberValue(big.NewFloat(val))
        } else {
            data.SessionReplayMonthlyBudgetInGb = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.SessionReplayMonthlyBudgetInGb = types.NumberNull()
    }
    if obj, ok := dataMap["createdByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["createdByUserId"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if val, ok := dataMap["isArchived"].(bool); ok {
        data.IsArchived = types.BoolValue(val)
    }
    if obj, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.CreatedAt = NewRFC3339Value(val)
        } else {
            data.CreatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["createdAt"].(string); ok && val != "" {
        data.CreatedAt = NewRFC3339Value(val)
    } else {
        data.CreatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.UpdatedAt = NewRFC3339Value(val)
        } else {
            data.UpdatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["updatedAt"].(string); ok && val != "" {
        data.UpdatedAt = NewRFC3339Value(val)
    } else {
        data.UpdatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["deletedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.DeletedAt = NewRFC3339Value(val)
        } else {
            data.DeletedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["deletedAt"].(string); ok && val != "" {
        data.DeletedAt = NewRFC3339Value(val)
    } else {
        data.DeletedAt = NewRFC3339Null()
    }
    if val, ok := dataMap["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["version"].(int); ok {
        data.Version = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["version"].(int64); ok {
        data.Version = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["version"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.Version = types.NumberValue(big.NewFloat(val))
        } else {
            data.Version = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.Version = types.NumberNull()
    }
    if obj, ok := dataMap["slug"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Slug = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Slug = types.StringValue(string(jsonBytes))
            } else {
                data.Slug = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Slug = types.StringValue(string(jsonBytes))
            } else {
                data.Slug = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Slug = types.StringValue(string(jsonBytes))
        } else {
            data.Slug = types.StringNull()
        }
    } else if val, ok := dataMap["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
    }
    if obj, ok := dataMap["clientType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ClientType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ClientType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ClientType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ClientType = types.StringValue(string(jsonBytes))
            } else {
                data.ClientType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ClientType = types.StringValue(string(jsonBytes))
            } else {
                data.ClientType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ClientType = types.StringValue(string(jsonBytes))
        } else {
            data.ClientType = types.StringNull()
        }
    } else if val, ok := dataMap["clientType"].(string); ok {
        data.ClientType = types.StringValue(val)
    } else {
        data.ClientType = types.StringNull()
    }
    if obj, ok := dataMap["sdkLanguage"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SdkLanguage = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SdkLanguage = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SdkLanguage = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SdkLanguage = types.StringValue(string(jsonBytes))
            } else {
                data.SdkLanguage = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SdkLanguage = types.StringValue(string(jsonBytes))
            } else {
                data.SdkLanguage = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SdkLanguage = types.StringValue(string(jsonBytes))
        } else {
            data.SdkLanguage = types.StringNull()
        }
    } else if val, ok := dataMap["sdkLanguage"].(string); ok {
        data.SdkLanguage = types.StringValue(val)
    } else {
        data.SdkLanguage = types.StringNull()
    }
    if obj, ok := dataMap["otelCollectorStatus"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OtelCollectorStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.OtelCollectorStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.OtelCollectorStatus = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.OtelCollectorStatus = types.StringValue(string(jsonBytes))
            } else {
                data.OtelCollectorStatus = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.OtelCollectorStatus = types.StringValue(string(jsonBytes))
            } else {
                data.OtelCollectorStatus = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.OtelCollectorStatus = types.StringValue(string(jsonBytes))
        } else {
            data.OtelCollectorStatus = types.StringNull()
        }
    } else if val, ok := dataMap["otelCollectorStatus"].(string); ok {
        data.OtelCollectorStatus = types.StringValue(val)
    } else {
        data.OtelCollectorStatus = types.StringNull()
    }
    if obj, ok := dataMap["agentVersion"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AgentVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AgentVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AgentVersion = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AgentVersion = types.StringValue(string(jsonBytes))
            } else {
                data.AgentVersion = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AgentVersion = types.StringValue(string(jsonBytes))
            } else {
                data.AgentVersion = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AgentVersion = types.StringValue(string(jsonBytes))
        } else {
            data.AgentVersion = types.StringNull()
        }
    } else if val, ok := dataMap["agentVersion"].(string); ok {
        data.AgentVersion = types.StringValue(val)
    } else {
        data.AgentVersion = types.StringNull()
    }
    if obj, ok := dataMap["lastSeenAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.LastSeenAt = NewRFC3339Value(val)
        } else {
            data.LastSeenAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["lastSeenAt"].(string); ok && val != "" {
        data.LastSeenAt = NewRFC3339Value(val)
    } else {
        data.LastSeenAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["sessionReplayLastChunkReceivedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.SessionReplayLastChunkReceivedAt = NewRFC3339Value(val)
        } else {
            data.SessionReplayLastChunkReceivedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["sessionReplayLastChunkReceivedAt"].(string); ok && val != "" {
        data.SessionReplayLastChunkReceivedAt = NewRFC3339Value(val)
    } else {
        data.SessionReplayLastChunkReceivedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["sessionReplayBudgetExceededAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.SessionReplayBudgetExceededAt = NewRFC3339Value(val)
        } else {
            data.SessionReplayBudgetExceededAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["sessionReplayBudgetExceededAt"].(string); ok && val != "" {
        data.SessionReplayBudgetExceededAt = NewRFC3339Value(val)
    } else {
        data.SessionReplayBudgetExceededAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["archivedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.ArchivedAt = NewRFC3339Value(val)
        } else {
            data.ArchivedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["archivedAt"].(string); ok && val != "" {
        data.ArchivedAt = NewRFC3339Value(val)
    } else {
        data.ArchivedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["archivedByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ArchivedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ArchivedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ArchivedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ArchivedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.ArchivedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ArchivedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.ArchivedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ArchivedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.ArchivedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["archivedByUserId"].(string); ok {
        data.ArchivedByUserId = types.StringValue(val)
    } else {
        data.ArchivedByUserId = types.StringNull()
    }
    if obj, ok := dataMap["deletedByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeletedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DeletedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DeletedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DeletedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.DeletedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DeletedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.DeletedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DeletedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.DeletedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["deletedByUserId"].(string); ok {
        data.DeletedByUserId = types.StringValue(val)
    } else {
        data.DeletedByUserId = types.StringNull()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    // The read response is authoritative, but never let it clobber the id we just received.
    data.Id = types.StringValue(createdId)

    // Write logs using the tflog package
    tflog.Trace(ctx, "created a resource")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RumApplicationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var data RumApplicationResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create select parameter to get full object
    selectParam := map[string]interface{}{
        "projectId": true,
        "name": true,
        "description": true,
        "appIdentifier": true,
        "labels": true,
        "retainTelemetryDataForDays": true,
        "telemetryRetentionConfig": true,
        "isSessionReplayEnabled": true,
        "sessionReplayMaskingMode": true,
        "sessionReplayMaskSelectors": true,
        "sessionReplayBlockSelectors": true,
        "sessionReplayIgnoreErrorPatterns": true,
        "sessionReplayTracePropagationOrigins": true,
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
        "createdByUserId": true,
        "isArchived": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "slug": true,
        "clientType": true,
        "sdkLanguage": true,
        "otelCollectorStatus": true,
        "agentVersion": true,
        "lastSeenAt": true,
        "sessionReplayLastChunkReceivedAt": true,
        "sessionReplayBudgetExceededAt": true,
        "archivedAt": true,
        "archivedByUserId": true,
        "deletedByUserId": true,
        "_id": true,
    }

    // Make API call with select parameter
    httpResp, err := r.client.PostWithSelect(ctx, "/rum-application/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read rum_application, got error: %s", err))
        return
    }

    if httpResp.StatusCode == http.StatusNotFound {
        resp.State.RemoveResource(ctx)
        return
    }

    var rumApplicationResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &rumApplicationResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse rum_application response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := rumApplicationResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = rumApplicationResponse
    }

    if obj, ok := dataMap["projectId"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok {
            data.ProjectId = types.StringValue(val)
        } else {
            data.ProjectId = types.StringNull()
        }
    } else if val, ok := dataMap["projectId"].(string); ok {
        data.ProjectId = types.StringValue(val)
    } else {
        data.ProjectId = types.StringNull()
    }
    if obj, ok := dataMap["name"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Name = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Name = types.StringValue(string(jsonBytes))
            } else {
                data.Name = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Name = types.StringValue(string(jsonBytes))
            } else {
                data.Name = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Name = types.StringValue(string(jsonBytes))
        } else {
            data.Name = types.StringNull()
        }
    } else if val, ok := dataMap["name"].(string); ok {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if obj, ok := dataMap["description"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Description = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Description = types.StringValue(string(jsonBytes))
            } else {
                data.Description = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Description = types.StringValue(string(jsonBytes))
            } else {
                data.Description = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Description = types.StringValue(string(jsonBytes))
        } else {
            data.Description = types.StringNull()
        }
    } else if val, ok := dataMap["description"].(string); ok {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if obj, ok := dataMap["appIdentifier"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AppIdentifier = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AppIdentifier = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AppIdentifier = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AppIdentifier = types.StringValue(string(jsonBytes))
            } else {
                data.AppIdentifier = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AppIdentifier = types.StringValue(string(jsonBytes))
            } else {
                data.AppIdentifier = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AppIdentifier = types.StringValue(string(jsonBytes))
        } else {
            data.AppIdentifier = types.StringNull()
        }
    } else if val, ok := dataMap["appIdentifier"].(string); ok {
        data.AppIdentifier = types.StringValue(val)
    } else {
        data.AppIdentifier = types.StringNull()
    }
    if val, ok := dataMap["labels"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.Labels = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.Labels = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if val, ok := dataMap["retainTelemetryDataForDays"].(float64); ok {
        data.RetainTelemetryDataForDays = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["retainTelemetryDataForDays"].(int); ok {
        data.RetainTelemetryDataForDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["retainTelemetryDataForDays"].(int64); ok {
        data.RetainTelemetryDataForDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["retainTelemetryDataForDays"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.RetainTelemetryDataForDays = types.NumberValue(big.NewFloat(val))
        } else {
            data.RetainTelemetryDataForDays = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.RetainTelemetryDataForDays = types.NumberNull()
    }
    if obj, ok := dataMap["telemetryRetentionConfig"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TelemetryRetentionConfig = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.TelemetryRetentionConfig = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.TelemetryRetentionConfig = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.TelemetryRetentionConfig = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.TelemetryRetentionConfig = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.TelemetryRetentionConfig = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.TelemetryRetentionConfig = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.TelemetryRetentionConfig = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.TelemetryRetentionConfig = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["telemetryRetentionConfig"].(string); ok {
        data.TelemetryRetentionConfig = NewJSONSubsetValue(val)
    } else {
        data.TelemetryRetentionConfig = NewJSONSubsetNull()
    }
    if val, ok := dataMap["isSessionReplayEnabled"].(bool); ok {
        data.IsSessionReplayEnabled = types.BoolValue(val)
    }
    if obj, ok := dataMap["sessionReplayMaskingMode"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SessionReplayMaskingMode = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SessionReplayMaskingMode = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SessionReplayMaskingMode = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SessionReplayMaskingMode = types.StringValue(string(jsonBytes))
            } else {
                data.SessionReplayMaskingMode = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SessionReplayMaskingMode = types.StringValue(string(jsonBytes))
            } else {
                data.SessionReplayMaskingMode = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SessionReplayMaskingMode = types.StringValue(string(jsonBytes))
        } else {
            data.SessionReplayMaskingMode = types.StringNull()
        }
    } else if val, ok := dataMap["sessionReplayMaskingMode"].(string); ok {
        data.SessionReplayMaskingMode = types.StringValue(val)
    } else {
        data.SessionReplayMaskingMode = types.StringNull()
    }
    if obj, ok := dataMap["sessionReplayMaskSelectors"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SessionReplayMaskSelectors = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SessionReplayMaskSelectors = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SessionReplayMaskSelectors = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SessionReplayMaskSelectors = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.SessionReplayMaskSelectors = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SessionReplayMaskSelectors = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.SessionReplayMaskSelectors = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SessionReplayMaskSelectors = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.SessionReplayMaskSelectors = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["sessionReplayMaskSelectors"].(string); ok {
        data.SessionReplayMaskSelectors = NewJSONSubsetValue(val)
    } else {
        data.SessionReplayMaskSelectors = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["sessionReplayBlockSelectors"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SessionReplayBlockSelectors = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SessionReplayBlockSelectors = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SessionReplayBlockSelectors = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SessionReplayBlockSelectors = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.SessionReplayBlockSelectors = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SessionReplayBlockSelectors = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.SessionReplayBlockSelectors = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SessionReplayBlockSelectors = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.SessionReplayBlockSelectors = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["sessionReplayBlockSelectors"].(string); ok {
        data.SessionReplayBlockSelectors = NewJSONSubsetValue(val)
    } else {
        data.SessionReplayBlockSelectors = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["sessionReplayIgnoreErrorPatterns"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SessionReplayIgnoreErrorPatterns = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SessionReplayIgnoreErrorPatterns = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SessionReplayIgnoreErrorPatterns = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SessionReplayIgnoreErrorPatterns = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.SessionReplayIgnoreErrorPatterns = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SessionReplayIgnoreErrorPatterns = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.SessionReplayIgnoreErrorPatterns = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SessionReplayIgnoreErrorPatterns = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.SessionReplayIgnoreErrorPatterns = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["sessionReplayIgnoreErrorPatterns"].(string); ok {
        data.SessionReplayIgnoreErrorPatterns = NewJSONSubsetValue(val)
    } else {
        data.SessionReplayIgnoreErrorPatterns = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["sessionReplayTracePropagationOrigins"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SessionReplayTracePropagationOrigins = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SessionReplayTracePropagationOrigins = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SessionReplayTracePropagationOrigins = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SessionReplayTracePropagationOrigins = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.SessionReplayTracePropagationOrigins = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SessionReplayTracePropagationOrigins = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.SessionReplayTracePropagationOrigins = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SessionReplayTracePropagationOrigins = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.SessionReplayTracePropagationOrigins = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["sessionReplayTracePropagationOrigins"].(string); ok {
        data.SessionReplayTracePropagationOrigins = NewJSONSubsetValue(val)
    } else {
        data.SessionReplayTracePropagationOrigins = NewJSONSubsetNull()
    }
    if val, ok := dataMap["sessionReplayLcpBudgetMs"].(float64); ok {
        data.SessionReplayLcpBudgetMs = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["sessionReplayLcpBudgetMs"].(int); ok {
        data.SessionReplayLcpBudgetMs = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["sessionReplayLcpBudgetMs"].(int64); ok {
        data.SessionReplayLcpBudgetMs = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["sessionReplayLcpBudgetMs"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.SessionReplayLcpBudgetMs = types.NumberValue(big.NewFloat(val))
        } else {
            data.SessionReplayLcpBudgetMs = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.SessionReplayLcpBudgetMs = types.NumberNull()
    }
    if val, ok := dataMap["sessionReplayLongTaskBudgetMs"].(float64); ok {
        data.SessionReplayLongTaskBudgetMs = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["sessionReplayLongTaskBudgetMs"].(int); ok {
        data.SessionReplayLongTaskBudgetMs = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["sessionReplayLongTaskBudgetMs"].(int64); ok {
        data.SessionReplayLongTaskBudgetMs = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["sessionReplayLongTaskBudgetMs"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.SessionReplayLongTaskBudgetMs = types.NumberValue(big.NewFloat(val))
        } else {
            data.SessionReplayLongTaskBudgetMs = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.SessionReplayLongTaskBudgetMs = types.NumberNull()
    }
    if val, ok := dataMap["sessionReplaySlowRequestBudgetMs"].(float64); ok {
        data.SessionReplaySlowRequestBudgetMs = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["sessionReplaySlowRequestBudgetMs"].(int); ok {
        data.SessionReplaySlowRequestBudgetMs = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["sessionReplaySlowRequestBudgetMs"].(int64); ok {
        data.SessionReplaySlowRequestBudgetMs = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["sessionReplaySlowRequestBudgetMs"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.SessionReplaySlowRequestBudgetMs = types.NumberValue(big.NewFloat(val))
        } else {
            data.SessionReplaySlowRequestBudgetMs = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.SessionReplaySlowRequestBudgetMs = types.NumberNull()
    }
    if obj, ok := dataMap["sessionReplayAllowedOrigins"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SessionReplayAllowedOrigins = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SessionReplayAllowedOrigins = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SessionReplayAllowedOrigins = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SessionReplayAllowedOrigins = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.SessionReplayAllowedOrigins = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SessionReplayAllowedOrigins = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.SessionReplayAllowedOrigins = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SessionReplayAllowedOrigins = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.SessionReplayAllowedOrigins = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["sessionReplayAllowedOrigins"].(string); ok {
        data.SessionReplayAllowedOrigins = NewJSONSubsetValue(val)
    } else {
        data.SessionReplayAllowedOrigins = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["sessionReplayConsentMode"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SessionReplayConsentMode = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SessionReplayConsentMode = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SessionReplayConsentMode = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SessionReplayConsentMode = types.StringValue(string(jsonBytes))
            } else {
                data.SessionReplayConsentMode = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SessionReplayConsentMode = types.StringValue(string(jsonBytes))
            } else {
                data.SessionReplayConsentMode = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SessionReplayConsentMode = types.StringValue(string(jsonBytes))
        } else {
            data.SessionReplayConsentMode = types.StringNull()
        }
    } else if val, ok := dataMap["sessionReplayConsentMode"].(string); ok {
        data.SessionReplayConsentMode = types.StringValue(val)
    } else {
        data.SessionReplayConsentMode = types.StringNull()
    }
    if obj, ok := dataMap["sessionReplayCaptureTrigger"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SessionReplayCaptureTrigger = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SessionReplayCaptureTrigger = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SessionReplayCaptureTrigger = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SessionReplayCaptureTrigger = types.StringValue(string(jsonBytes))
            } else {
                data.SessionReplayCaptureTrigger = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SessionReplayCaptureTrigger = types.StringValue(string(jsonBytes))
            } else {
                data.SessionReplayCaptureTrigger = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SessionReplayCaptureTrigger = types.StringValue(string(jsonBytes))
        } else {
            data.SessionReplayCaptureTrigger = types.StringNull()
        }
    } else if val, ok := dataMap["sessionReplayCaptureTrigger"].(string); ok {
        data.SessionReplayCaptureTrigger = types.StringValue(val)
    } else {
        data.SessionReplayCaptureTrigger = types.StringNull()
    }
    if val, ok := dataMap["sessionReplaySamplePercentage"].(float64); ok {
        data.SessionReplaySamplePercentage = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["sessionReplaySamplePercentage"].(int); ok {
        data.SessionReplaySamplePercentage = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["sessionReplaySamplePercentage"].(int64); ok {
        data.SessionReplaySamplePercentage = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["sessionReplaySamplePercentage"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.SessionReplaySamplePercentage = types.NumberValue(big.NewFloat(val))
        } else {
            data.SessionReplaySamplePercentage = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.SessionReplaySamplePercentage = types.NumberNull()
    }
    if val, ok := dataMap["sessionReplayCaptureUserIdentity"].(bool); ok {
        data.SessionReplayCaptureUserIdentity = types.BoolValue(val)
    }
    if val, ok := dataMap["sessionReplayCaptureGeo"].(bool); ok {
        data.SessionReplayCaptureGeo = types.BoolValue(val)
    }
    if val, ok := dataMap["sessionReplayRecordCanvas"].(bool); ok {
        data.SessionReplayRecordCanvas = types.BoolValue(val)
    }
    if val, ok := dataMap["sessionReplayRetentionInDays"].(float64); ok {
        data.SessionReplayRetentionInDays = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["sessionReplayRetentionInDays"].(int); ok {
        data.SessionReplayRetentionInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["sessionReplayRetentionInDays"].(int64); ok {
        data.SessionReplayRetentionInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["sessionReplayRetentionInDays"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.SessionReplayRetentionInDays = types.NumberValue(big.NewFloat(val))
        } else {
            data.SessionReplayRetentionInDays = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.SessionReplayRetentionInDays = types.NumberNull()
    }
    if val, ok := dataMap["sessionReplayMonthlyBudgetInGB"].(float64); ok {
        data.SessionReplayMonthlyBudgetInGb = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["sessionReplayMonthlyBudgetInGB"].(int); ok {
        data.SessionReplayMonthlyBudgetInGb = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["sessionReplayMonthlyBudgetInGB"].(int64); ok {
        data.SessionReplayMonthlyBudgetInGb = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["sessionReplayMonthlyBudgetInGB"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.SessionReplayMonthlyBudgetInGb = types.NumberValue(big.NewFloat(val))
        } else {
            data.SessionReplayMonthlyBudgetInGb = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.SessionReplayMonthlyBudgetInGb = types.NumberNull()
    }
    if obj, ok := dataMap["createdByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["createdByUserId"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if val, ok := dataMap["isArchived"].(bool); ok {
        data.IsArchived = types.BoolValue(val)
    }
    if obj, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.CreatedAt = NewRFC3339Value(val)
        } else {
            data.CreatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["createdAt"].(string); ok && val != "" {
        data.CreatedAt = NewRFC3339Value(val)
    } else {
        data.CreatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.UpdatedAt = NewRFC3339Value(val)
        } else {
            data.UpdatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["updatedAt"].(string); ok && val != "" {
        data.UpdatedAt = NewRFC3339Value(val)
    } else {
        data.UpdatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["deletedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.DeletedAt = NewRFC3339Value(val)
        } else {
            data.DeletedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["deletedAt"].(string); ok && val != "" {
        data.DeletedAt = NewRFC3339Value(val)
    } else {
        data.DeletedAt = NewRFC3339Null()
    }
    if val, ok := dataMap["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["version"].(int); ok {
        data.Version = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["version"].(int64); ok {
        data.Version = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["version"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.Version = types.NumberValue(big.NewFloat(val))
        } else {
            data.Version = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.Version = types.NumberNull()
    }
    if obj, ok := dataMap["slug"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Slug = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Slug = types.StringValue(string(jsonBytes))
            } else {
                data.Slug = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Slug = types.StringValue(string(jsonBytes))
            } else {
                data.Slug = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Slug = types.StringValue(string(jsonBytes))
        } else {
            data.Slug = types.StringNull()
        }
    } else if val, ok := dataMap["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
    }
    if obj, ok := dataMap["clientType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ClientType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ClientType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ClientType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ClientType = types.StringValue(string(jsonBytes))
            } else {
                data.ClientType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ClientType = types.StringValue(string(jsonBytes))
            } else {
                data.ClientType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ClientType = types.StringValue(string(jsonBytes))
        } else {
            data.ClientType = types.StringNull()
        }
    } else if val, ok := dataMap["clientType"].(string); ok {
        data.ClientType = types.StringValue(val)
    } else {
        data.ClientType = types.StringNull()
    }
    if obj, ok := dataMap["sdkLanguage"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SdkLanguage = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SdkLanguage = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SdkLanguage = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SdkLanguage = types.StringValue(string(jsonBytes))
            } else {
                data.SdkLanguage = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SdkLanguage = types.StringValue(string(jsonBytes))
            } else {
                data.SdkLanguage = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SdkLanguage = types.StringValue(string(jsonBytes))
        } else {
            data.SdkLanguage = types.StringNull()
        }
    } else if val, ok := dataMap["sdkLanguage"].(string); ok {
        data.SdkLanguage = types.StringValue(val)
    } else {
        data.SdkLanguage = types.StringNull()
    }
    if obj, ok := dataMap["otelCollectorStatus"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OtelCollectorStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.OtelCollectorStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.OtelCollectorStatus = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.OtelCollectorStatus = types.StringValue(string(jsonBytes))
            } else {
                data.OtelCollectorStatus = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.OtelCollectorStatus = types.StringValue(string(jsonBytes))
            } else {
                data.OtelCollectorStatus = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.OtelCollectorStatus = types.StringValue(string(jsonBytes))
        } else {
            data.OtelCollectorStatus = types.StringNull()
        }
    } else if val, ok := dataMap["otelCollectorStatus"].(string); ok {
        data.OtelCollectorStatus = types.StringValue(val)
    } else {
        data.OtelCollectorStatus = types.StringNull()
    }
    if obj, ok := dataMap["agentVersion"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AgentVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AgentVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AgentVersion = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AgentVersion = types.StringValue(string(jsonBytes))
            } else {
                data.AgentVersion = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AgentVersion = types.StringValue(string(jsonBytes))
            } else {
                data.AgentVersion = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AgentVersion = types.StringValue(string(jsonBytes))
        } else {
            data.AgentVersion = types.StringNull()
        }
    } else if val, ok := dataMap["agentVersion"].(string); ok {
        data.AgentVersion = types.StringValue(val)
    } else {
        data.AgentVersion = types.StringNull()
    }
    if obj, ok := dataMap["lastSeenAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.LastSeenAt = NewRFC3339Value(val)
        } else {
            data.LastSeenAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["lastSeenAt"].(string); ok && val != "" {
        data.LastSeenAt = NewRFC3339Value(val)
    } else {
        data.LastSeenAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["sessionReplayLastChunkReceivedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.SessionReplayLastChunkReceivedAt = NewRFC3339Value(val)
        } else {
            data.SessionReplayLastChunkReceivedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["sessionReplayLastChunkReceivedAt"].(string); ok && val != "" {
        data.SessionReplayLastChunkReceivedAt = NewRFC3339Value(val)
    } else {
        data.SessionReplayLastChunkReceivedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["sessionReplayBudgetExceededAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.SessionReplayBudgetExceededAt = NewRFC3339Value(val)
        } else {
            data.SessionReplayBudgetExceededAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["sessionReplayBudgetExceededAt"].(string); ok && val != "" {
        data.SessionReplayBudgetExceededAt = NewRFC3339Value(val)
    } else {
        data.SessionReplayBudgetExceededAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["archivedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.ArchivedAt = NewRFC3339Value(val)
        } else {
            data.ArchivedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["archivedAt"].(string); ok && val != "" {
        data.ArchivedAt = NewRFC3339Value(val)
    } else {
        data.ArchivedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["archivedByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ArchivedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ArchivedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ArchivedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ArchivedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.ArchivedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ArchivedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.ArchivedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ArchivedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.ArchivedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["archivedByUserId"].(string); ok {
        data.ArchivedByUserId = types.StringValue(val)
    } else {
        data.ArchivedByUserId = types.StringNull()
    }
    if obj, ok := dataMap["deletedByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeletedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DeletedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DeletedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DeletedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.DeletedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DeletedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.DeletedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DeletedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.DeletedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["deletedByUserId"].(string); ok {
        data.DeletedByUserId = types.StringValue(val)
    } else {
        data.DeletedByUserId = types.StringNull()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RumApplicationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var data RumApplicationResourceModel
    var state RumApplicationResourceModel

    // Read Terraform current state data to get the ID
    resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Read Terraform plan data to get the new values
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Use the ID from the current state
    data.Id = state.Id

    // Create API request body
    rumApplicationRequest := map[string]interface{}{
        "data": map[string]interface{}{},
    }
    requestDataMap := rumApplicationRequest["data"].(map[string]interface{})

    if !data.Name.IsUnknown() && !state.Name.IsUnknown() && !data.Name.Equal(state.Name) {
        requestDataMap["name"] = data.Name.ValueString()
    }
    if !data.Description.IsUnknown() && !state.Description.IsUnknown() && !data.Description.Equal(state.Description) {
        requestDataMap["description"] = data.Description.ValueString()
    }
    if !data.Labels.IsUnknown() && !state.Labels.IsUnknown() && !data.Labels.Equal(state.Labels) {
        requestDataMap["labels"] = r.convertTerraformSetToInterface(data.Labels)
    }
    if !data.RetainTelemetryDataForDays.IsUnknown() && !state.RetainTelemetryDataForDays.IsUnknown() && !data.RetainTelemetryDataForDays.Equal(state.RetainTelemetryDataForDays) {
        requestDataMap["retainTelemetryDataForDays"] = r.bigFloatToFloat64(data.RetainTelemetryDataForDays.ValueBigFloat())
    }
    if !data.TelemetryRetentionConfig.IsUnknown() && !state.TelemetryRetentionConfig.IsUnknown() && !data.TelemetryRetentionConfig.Equal(state.TelemetryRetentionConfig) {
        var telemetryretentionconfigData interface{}
        if err := json.Unmarshal([]byte(data.TelemetryRetentionConfig.ValueString()), &telemetryretentionconfigData); err == nil {
            requestDataMap["telemetryRetentionConfig"] = telemetryretentionconfigData
        } else {
            requestDataMap["telemetryRetentionConfig"] = data.TelemetryRetentionConfig.ValueString()
        }
    }
    if !data.IsSessionReplayEnabled.IsUnknown() && !state.IsSessionReplayEnabled.IsUnknown() && !data.IsSessionReplayEnabled.Equal(state.IsSessionReplayEnabled) {
        requestDataMap["isSessionReplayEnabled"] = data.IsSessionReplayEnabled.ValueBool()
    }
    if !data.SessionReplayMaskingMode.IsUnknown() && !state.SessionReplayMaskingMode.IsUnknown() && !data.SessionReplayMaskingMode.Equal(state.SessionReplayMaskingMode) {
        requestDataMap["sessionReplayMaskingMode"] = data.SessionReplayMaskingMode.ValueString()
    }
    if !data.SessionReplayMaskSelectors.IsUnknown() && !state.SessionReplayMaskSelectors.IsUnknown() && !data.SessionReplayMaskSelectors.Equal(state.SessionReplayMaskSelectors) {
        var sessionreplaymaskselectorsData interface{}
        if err := json.Unmarshal([]byte(data.SessionReplayMaskSelectors.ValueString()), &sessionreplaymaskselectorsData); err == nil {
            requestDataMap["sessionReplayMaskSelectors"] = sessionreplaymaskselectorsData
        } else {
            requestDataMap["sessionReplayMaskSelectors"] = data.SessionReplayMaskSelectors.ValueString()
        }
    }
    if !data.SessionReplayBlockSelectors.IsUnknown() && !state.SessionReplayBlockSelectors.IsUnknown() && !data.SessionReplayBlockSelectors.Equal(state.SessionReplayBlockSelectors) {
        var sessionreplayblockselectorsData interface{}
        if err := json.Unmarshal([]byte(data.SessionReplayBlockSelectors.ValueString()), &sessionreplayblockselectorsData); err == nil {
            requestDataMap["sessionReplayBlockSelectors"] = sessionreplayblockselectorsData
        } else {
            requestDataMap["sessionReplayBlockSelectors"] = data.SessionReplayBlockSelectors.ValueString()
        }
    }
    if !data.SessionReplayIgnoreErrorPatterns.IsUnknown() && !state.SessionReplayIgnoreErrorPatterns.IsUnknown() && !data.SessionReplayIgnoreErrorPatterns.Equal(state.SessionReplayIgnoreErrorPatterns) {
        var sessionreplayignoreerrorpatternsData interface{}
        if err := json.Unmarshal([]byte(data.SessionReplayIgnoreErrorPatterns.ValueString()), &sessionreplayignoreerrorpatternsData); err == nil {
            requestDataMap["sessionReplayIgnoreErrorPatterns"] = sessionreplayignoreerrorpatternsData
        } else {
            requestDataMap["sessionReplayIgnoreErrorPatterns"] = data.SessionReplayIgnoreErrorPatterns.ValueString()
        }
    }
    if !data.SessionReplayTracePropagationOrigins.IsUnknown() && !state.SessionReplayTracePropagationOrigins.IsUnknown() && !data.SessionReplayTracePropagationOrigins.Equal(state.SessionReplayTracePropagationOrigins) {
        var sessionreplaytracepropagationoriginsData interface{}
        if err := json.Unmarshal([]byte(data.SessionReplayTracePropagationOrigins.ValueString()), &sessionreplaytracepropagationoriginsData); err == nil {
            requestDataMap["sessionReplayTracePropagationOrigins"] = sessionreplaytracepropagationoriginsData
        } else {
            requestDataMap["sessionReplayTracePropagationOrigins"] = data.SessionReplayTracePropagationOrigins.ValueString()
        }
    }
    if !data.SessionReplayLcpBudgetMs.IsUnknown() && !state.SessionReplayLcpBudgetMs.IsUnknown() && !data.SessionReplayLcpBudgetMs.Equal(state.SessionReplayLcpBudgetMs) {
        requestDataMap["sessionReplayLcpBudgetMs"] = r.bigFloatToFloat64(data.SessionReplayLcpBudgetMs.ValueBigFloat())
    }
    if !data.SessionReplayLongTaskBudgetMs.IsUnknown() && !state.SessionReplayLongTaskBudgetMs.IsUnknown() && !data.SessionReplayLongTaskBudgetMs.Equal(state.SessionReplayLongTaskBudgetMs) {
        requestDataMap["sessionReplayLongTaskBudgetMs"] = r.bigFloatToFloat64(data.SessionReplayLongTaskBudgetMs.ValueBigFloat())
    }
    if !data.SessionReplaySlowRequestBudgetMs.IsUnknown() && !state.SessionReplaySlowRequestBudgetMs.IsUnknown() && !data.SessionReplaySlowRequestBudgetMs.Equal(state.SessionReplaySlowRequestBudgetMs) {
        requestDataMap["sessionReplaySlowRequestBudgetMs"] = r.bigFloatToFloat64(data.SessionReplaySlowRequestBudgetMs.ValueBigFloat())
    }
    if !data.SessionReplayAllowedOrigins.IsUnknown() && !state.SessionReplayAllowedOrigins.IsUnknown() && !data.SessionReplayAllowedOrigins.Equal(state.SessionReplayAllowedOrigins) {
        var sessionreplayallowedoriginsData interface{}
        if err := json.Unmarshal([]byte(data.SessionReplayAllowedOrigins.ValueString()), &sessionreplayallowedoriginsData); err == nil {
            requestDataMap["sessionReplayAllowedOrigins"] = sessionreplayallowedoriginsData
        } else {
            requestDataMap["sessionReplayAllowedOrigins"] = data.SessionReplayAllowedOrigins.ValueString()
        }
    }
    if !data.SessionReplayConsentMode.IsUnknown() && !state.SessionReplayConsentMode.IsUnknown() && !data.SessionReplayConsentMode.Equal(state.SessionReplayConsentMode) {
        requestDataMap["sessionReplayConsentMode"] = data.SessionReplayConsentMode.ValueString()
    }
    if !data.SessionReplayCaptureTrigger.IsUnknown() && !state.SessionReplayCaptureTrigger.IsUnknown() && !data.SessionReplayCaptureTrigger.Equal(state.SessionReplayCaptureTrigger) {
        requestDataMap["sessionReplayCaptureTrigger"] = data.SessionReplayCaptureTrigger.ValueString()
    }
    if !data.SessionReplaySamplePercentage.IsUnknown() && !state.SessionReplaySamplePercentage.IsUnknown() && !data.SessionReplaySamplePercentage.Equal(state.SessionReplaySamplePercentage) {
        requestDataMap["sessionReplaySamplePercentage"] = r.bigFloatToFloat64(data.SessionReplaySamplePercentage.ValueBigFloat())
    }
    if !data.SessionReplayCaptureUserIdentity.IsUnknown() && !state.SessionReplayCaptureUserIdentity.IsUnknown() && !data.SessionReplayCaptureUserIdentity.Equal(state.SessionReplayCaptureUserIdentity) {
        requestDataMap["sessionReplayCaptureUserIdentity"] = data.SessionReplayCaptureUserIdentity.ValueBool()
    }
    if !data.SessionReplayCaptureGeo.IsUnknown() && !state.SessionReplayCaptureGeo.IsUnknown() && !data.SessionReplayCaptureGeo.Equal(state.SessionReplayCaptureGeo) {
        requestDataMap["sessionReplayCaptureGeo"] = data.SessionReplayCaptureGeo.ValueBool()
    }
    if !data.SessionReplayRecordCanvas.IsUnknown() && !state.SessionReplayRecordCanvas.IsUnknown() && !data.SessionReplayRecordCanvas.Equal(state.SessionReplayRecordCanvas) {
        requestDataMap["sessionReplayRecordCanvas"] = data.SessionReplayRecordCanvas.ValueBool()
    }
    if !data.SessionReplayRetentionInDays.IsUnknown() && !state.SessionReplayRetentionInDays.IsUnknown() && !data.SessionReplayRetentionInDays.Equal(state.SessionReplayRetentionInDays) {
        requestDataMap["sessionReplayRetentionInDays"] = r.bigFloatToFloat64(data.SessionReplayRetentionInDays.ValueBigFloat())
    }
    if !data.SessionReplayMonthlyBudgetInGb.IsUnknown() && !state.SessionReplayMonthlyBudgetInGb.IsUnknown() && !data.SessionReplayMonthlyBudgetInGb.Equal(state.SessionReplayMonthlyBudgetInGb) {
        requestDataMap["sessionReplayMonthlyBudgetInGB"] = r.bigFloatToFloat64(data.SessionReplayMonthlyBudgetInGb.ValueBigFloat())
    }
    if !data.IsArchived.IsUnknown() && !state.IsArchived.IsUnknown() && !data.IsArchived.Equal(state.IsArchived) {
        requestDataMap["isArchived"] = data.IsArchived.ValueBool()
    }

    // Only call the API when there are changed fields to send. An empty
    // update body is rejected by the API; state is still refreshed below so
    // this method never writes unverified plan values into state.
    if len(rumApplicationRequest["data"].(map[string]interface{})) > 0 {
        httpResp, err := r.client.Put(ctx, "/rum-application/" + data.Id.ValueString() + "", rumApplicationRequest)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update rum_application, got error: %s", err))
            return
        }

        // Parse the update response
        var rumApplicationResponse map[string]interface{}
        err = r.client.ParseResponse(httpResp, &rumApplicationResponse)
        if err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to update rum_application: %s", err))
            return
        }
        _ = rumApplicationResponse
    }

    // After successful update, fetch the current state by calling Read with select parameter
    selectParam := map[string]interface{}{
        "projectId": true,
        "name": true,
        "description": true,
        "appIdentifier": true,
        "labels": true,
        "retainTelemetryDataForDays": true,
        "telemetryRetentionConfig": true,
        "isSessionReplayEnabled": true,
        "sessionReplayMaskingMode": true,
        "sessionReplayMaskSelectors": true,
        "sessionReplayBlockSelectors": true,
        "sessionReplayIgnoreErrorPatterns": true,
        "sessionReplayTracePropagationOrigins": true,
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
        "createdByUserId": true,
        "isArchived": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "slug": true,
        "clientType": true,
        "sdkLanguage": true,
        "otelCollectorStatus": true,
        "agentVersion": true,
        "lastSeenAt": true,
        "sessionReplayLastChunkReceivedAt": true,
        "sessionReplayBudgetExceededAt": true,
        "archivedAt": true,
        "archivedByUserId": true,
        "deletedByUserId": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect(ctx, "/rum-application/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read rum_application after update, got error: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read rum_application after update: %s", err))
        return
    }

    // Update the model with response data from the Read operation
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := readResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = readResponse
    }

    if obj, ok := dataMap["projectId"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok {
            data.ProjectId = types.StringValue(val)
        } else {
            data.ProjectId = types.StringNull()
        }
    } else if val, ok := dataMap["projectId"].(string); ok {
        data.ProjectId = types.StringValue(val)
    } else {
        data.ProjectId = types.StringNull()
    }
    if obj, ok := dataMap["name"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Name = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Name = types.StringValue(string(jsonBytes))
            } else {
                data.Name = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Name = types.StringValue(string(jsonBytes))
            } else {
                data.Name = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Name = types.StringValue(string(jsonBytes))
        } else {
            data.Name = types.StringNull()
        }
    } else if val, ok := dataMap["name"].(string); ok {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if obj, ok := dataMap["description"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Description = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Description = types.StringValue(string(jsonBytes))
            } else {
                data.Description = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Description = types.StringValue(string(jsonBytes))
            } else {
                data.Description = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Description = types.StringValue(string(jsonBytes))
        } else {
            data.Description = types.StringNull()
        }
    } else if val, ok := dataMap["description"].(string); ok {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if obj, ok := dataMap["appIdentifier"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AppIdentifier = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AppIdentifier = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AppIdentifier = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AppIdentifier = types.StringValue(string(jsonBytes))
            } else {
                data.AppIdentifier = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AppIdentifier = types.StringValue(string(jsonBytes))
            } else {
                data.AppIdentifier = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AppIdentifier = types.StringValue(string(jsonBytes))
        } else {
            data.AppIdentifier = types.StringNull()
        }
    } else if val, ok := dataMap["appIdentifier"].(string); ok {
        data.AppIdentifier = types.StringValue(val)
    } else {
        data.AppIdentifier = types.StringNull()
    }
    if val, ok := dataMap["labels"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.Labels = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.Labels = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if val, ok := dataMap["retainTelemetryDataForDays"].(float64); ok {
        data.RetainTelemetryDataForDays = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["retainTelemetryDataForDays"].(int); ok {
        data.RetainTelemetryDataForDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["retainTelemetryDataForDays"].(int64); ok {
        data.RetainTelemetryDataForDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["retainTelemetryDataForDays"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.RetainTelemetryDataForDays = types.NumberValue(big.NewFloat(val))
        } else {
            data.RetainTelemetryDataForDays = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.RetainTelemetryDataForDays = types.NumberNull()
    }
    if obj, ok := dataMap["telemetryRetentionConfig"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TelemetryRetentionConfig = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.TelemetryRetentionConfig = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.TelemetryRetentionConfig = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.TelemetryRetentionConfig = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.TelemetryRetentionConfig = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.TelemetryRetentionConfig = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.TelemetryRetentionConfig = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.TelemetryRetentionConfig = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.TelemetryRetentionConfig = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["telemetryRetentionConfig"].(string); ok {
        data.TelemetryRetentionConfig = NewJSONSubsetValue(val)
    } else {
        data.TelemetryRetentionConfig = NewJSONSubsetNull()
    }
    if val, ok := dataMap["isSessionReplayEnabled"].(bool); ok {
        data.IsSessionReplayEnabled = types.BoolValue(val)
    }
    if obj, ok := dataMap["sessionReplayMaskingMode"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SessionReplayMaskingMode = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SessionReplayMaskingMode = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SessionReplayMaskingMode = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SessionReplayMaskingMode = types.StringValue(string(jsonBytes))
            } else {
                data.SessionReplayMaskingMode = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SessionReplayMaskingMode = types.StringValue(string(jsonBytes))
            } else {
                data.SessionReplayMaskingMode = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SessionReplayMaskingMode = types.StringValue(string(jsonBytes))
        } else {
            data.SessionReplayMaskingMode = types.StringNull()
        }
    } else if val, ok := dataMap["sessionReplayMaskingMode"].(string); ok {
        data.SessionReplayMaskingMode = types.StringValue(val)
    } else {
        data.SessionReplayMaskingMode = types.StringNull()
    }
    if obj, ok := dataMap["sessionReplayMaskSelectors"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SessionReplayMaskSelectors = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SessionReplayMaskSelectors = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SessionReplayMaskSelectors = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SessionReplayMaskSelectors = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.SessionReplayMaskSelectors = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SessionReplayMaskSelectors = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.SessionReplayMaskSelectors = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SessionReplayMaskSelectors = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.SessionReplayMaskSelectors = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["sessionReplayMaskSelectors"].(string); ok {
        data.SessionReplayMaskSelectors = NewJSONSubsetValue(val)
    } else {
        data.SessionReplayMaskSelectors = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["sessionReplayBlockSelectors"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SessionReplayBlockSelectors = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SessionReplayBlockSelectors = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SessionReplayBlockSelectors = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SessionReplayBlockSelectors = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.SessionReplayBlockSelectors = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SessionReplayBlockSelectors = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.SessionReplayBlockSelectors = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SessionReplayBlockSelectors = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.SessionReplayBlockSelectors = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["sessionReplayBlockSelectors"].(string); ok {
        data.SessionReplayBlockSelectors = NewJSONSubsetValue(val)
    } else {
        data.SessionReplayBlockSelectors = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["sessionReplayIgnoreErrorPatterns"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SessionReplayIgnoreErrorPatterns = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SessionReplayIgnoreErrorPatterns = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SessionReplayIgnoreErrorPatterns = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SessionReplayIgnoreErrorPatterns = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.SessionReplayIgnoreErrorPatterns = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SessionReplayIgnoreErrorPatterns = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.SessionReplayIgnoreErrorPatterns = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SessionReplayIgnoreErrorPatterns = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.SessionReplayIgnoreErrorPatterns = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["sessionReplayIgnoreErrorPatterns"].(string); ok {
        data.SessionReplayIgnoreErrorPatterns = NewJSONSubsetValue(val)
    } else {
        data.SessionReplayIgnoreErrorPatterns = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["sessionReplayTracePropagationOrigins"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SessionReplayTracePropagationOrigins = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SessionReplayTracePropagationOrigins = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SessionReplayTracePropagationOrigins = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SessionReplayTracePropagationOrigins = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.SessionReplayTracePropagationOrigins = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SessionReplayTracePropagationOrigins = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.SessionReplayTracePropagationOrigins = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SessionReplayTracePropagationOrigins = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.SessionReplayTracePropagationOrigins = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["sessionReplayTracePropagationOrigins"].(string); ok {
        data.SessionReplayTracePropagationOrigins = NewJSONSubsetValue(val)
    } else {
        data.SessionReplayTracePropagationOrigins = NewJSONSubsetNull()
    }
    if val, ok := dataMap["sessionReplayLcpBudgetMs"].(float64); ok {
        data.SessionReplayLcpBudgetMs = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["sessionReplayLcpBudgetMs"].(int); ok {
        data.SessionReplayLcpBudgetMs = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["sessionReplayLcpBudgetMs"].(int64); ok {
        data.SessionReplayLcpBudgetMs = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["sessionReplayLcpBudgetMs"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.SessionReplayLcpBudgetMs = types.NumberValue(big.NewFloat(val))
        } else {
            data.SessionReplayLcpBudgetMs = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.SessionReplayLcpBudgetMs = types.NumberNull()
    }
    if val, ok := dataMap["sessionReplayLongTaskBudgetMs"].(float64); ok {
        data.SessionReplayLongTaskBudgetMs = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["sessionReplayLongTaskBudgetMs"].(int); ok {
        data.SessionReplayLongTaskBudgetMs = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["sessionReplayLongTaskBudgetMs"].(int64); ok {
        data.SessionReplayLongTaskBudgetMs = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["sessionReplayLongTaskBudgetMs"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.SessionReplayLongTaskBudgetMs = types.NumberValue(big.NewFloat(val))
        } else {
            data.SessionReplayLongTaskBudgetMs = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.SessionReplayLongTaskBudgetMs = types.NumberNull()
    }
    if val, ok := dataMap["sessionReplaySlowRequestBudgetMs"].(float64); ok {
        data.SessionReplaySlowRequestBudgetMs = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["sessionReplaySlowRequestBudgetMs"].(int); ok {
        data.SessionReplaySlowRequestBudgetMs = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["sessionReplaySlowRequestBudgetMs"].(int64); ok {
        data.SessionReplaySlowRequestBudgetMs = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["sessionReplaySlowRequestBudgetMs"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.SessionReplaySlowRequestBudgetMs = types.NumberValue(big.NewFloat(val))
        } else {
            data.SessionReplaySlowRequestBudgetMs = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.SessionReplaySlowRequestBudgetMs = types.NumberNull()
    }
    if obj, ok := dataMap["sessionReplayAllowedOrigins"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SessionReplayAllowedOrigins = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SessionReplayAllowedOrigins = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SessionReplayAllowedOrigins = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SessionReplayAllowedOrigins = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.SessionReplayAllowedOrigins = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SessionReplayAllowedOrigins = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.SessionReplayAllowedOrigins = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SessionReplayAllowedOrigins = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.SessionReplayAllowedOrigins = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["sessionReplayAllowedOrigins"].(string); ok {
        data.SessionReplayAllowedOrigins = NewJSONSubsetValue(val)
    } else {
        data.SessionReplayAllowedOrigins = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["sessionReplayConsentMode"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SessionReplayConsentMode = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SessionReplayConsentMode = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SessionReplayConsentMode = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SessionReplayConsentMode = types.StringValue(string(jsonBytes))
            } else {
                data.SessionReplayConsentMode = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SessionReplayConsentMode = types.StringValue(string(jsonBytes))
            } else {
                data.SessionReplayConsentMode = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SessionReplayConsentMode = types.StringValue(string(jsonBytes))
        } else {
            data.SessionReplayConsentMode = types.StringNull()
        }
    } else if val, ok := dataMap["sessionReplayConsentMode"].(string); ok {
        data.SessionReplayConsentMode = types.StringValue(val)
    } else {
        data.SessionReplayConsentMode = types.StringNull()
    }
    if obj, ok := dataMap["sessionReplayCaptureTrigger"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SessionReplayCaptureTrigger = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SessionReplayCaptureTrigger = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SessionReplayCaptureTrigger = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SessionReplayCaptureTrigger = types.StringValue(string(jsonBytes))
            } else {
                data.SessionReplayCaptureTrigger = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SessionReplayCaptureTrigger = types.StringValue(string(jsonBytes))
            } else {
                data.SessionReplayCaptureTrigger = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SessionReplayCaptureTrigger = types.StringValue(string(jsonBytes))
        } else {
            data.SessionReplayCaptureTrigger = types.StringNull()
        }
    } else if val, ok := dataMap["sessionReplayCaptureTrigger"].(string); ok {
        data.SessionReplayCaptureTrigger = types.StringValue(val)
    } else {
        data.SessionReplayCaptureTrigger = types.StringNull()
    }
    if val, ok := dataMap["sessionReplaySamplePercentage"].(float64); ok {
        data.SessionReplaySamplePercentage = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["sessionReplaySamplePercentage"].(int); ok {
        data.SessionReplaySamplePercentage = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["sessionReplaySamplePercentage"].(int64); ok {
        data.SessionReplaySamplePercentage = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["sessionReplaySamplePercentage"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.SessionReplaySamplePercentage = types.NumberValue(big.NewFloat(val))
        } else {
            data.SessionReplaySamplePercentage = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.SessionReplaySamplePercentage = types.NumberNull()
    }
    if val, ok := dataMap["sessionReplayCaptureUserIdentity"].(bool); ok {
        data.SessionReplayCaptureUserIdentity = types.BoolValue(val)
    }
    if val, ok := dataMap["sessionReplayCaptureGeo"].(bool); ok {
        data.SessionReplayCaptureGeo = types.BoolValue(val)
    }
    if val, ok := dataMap["sessionReplayRecordCanvas"].(bool); ok {
        data.SessionReplayRecordCanvas = types.BoolValue(val)
    }
    if val, ok := dataMap["sessionReplayRetentionInDays"].(float64); ok {
        data.SessionReplayRetentionInDays = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["sessionReplayRetentionInDays"].(int); ok {
        data.SessionReplayRetentionInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["sessionReplayRetentionInDays"].(int64); ok {
        data.SessionReplayRetentionInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["sessionReplayRetentionInDays"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.SessionReplayRetentionInDays = types.NumberValue(big.NewFloat(val))
        } else {
            data.SessionReplayRetentionInDays = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.SessionReplayRetentionInDays = types.NumberNull()
    }
    if val, ok := dataMap["sessionReplayMonthlyBudgetInGB"].(float64); ok {
        data.SessionReplayMonthlyBudgetInGb = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["sessionReplayMonthlyBudgetInGB"].(int); ok {
        data.SessionReplayMonthlyBudgetInGb = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["sessionReplayMonthlyBudgetInGB"].(int64); ok {
        data.SessionReplayMonthlyBudgetInGb = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["sessionReplayMonthlyBudgetInGB"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.SessionReplayMonthlyBudgetInGb = types.NumberValue(big.NewFloat(val))
        } else {
            data.SessionReplayMonthlyBudgetInGb = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.SessionReplayMonthlyBudgetInGb = types.NumberNull()
    }
    if obj, ok := dataMap["createdByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["createdByUserId"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if val, ok := dataMap["isArchived"].(bool); ok {
        data.IsArchived = types.BoolValue(val)
    }
    if obj, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.CreatedAt = NewRFC3339Value(val)
        } else {
            data.CreatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["createdAt"].(string); ok && val != "" {
        data.CreatedAt = NewRFC3339Value(val)
    } else {
        data.CreatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.UpdatedAt = NewRFC3339Value(val)
        } else {
            data.UpdatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["updatedAt"].(string); ok && val != "" {
        data.UpdatedAt = NewRFC3339Value(val)
    } else {
        data.UpdatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["deletedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.DeletedAt = NewRFC3339Value(val)
        } else {
            data.DeletedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["deletedAt"].(string); ok && val != "" {
        data.DeletedAt = NewRFC3339Value(val)
    } else {
        data.DeletedAt = NewRFC3339Null()
    }
    if val, ok := dataMap["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["version"].(int); ok {
        data.Version = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["version"].(int64); ok {
        data.Version = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["version"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.Version = types.NumberValue(big.NewFloat(val))
        } else {
            data.Version = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.Version = types.NumberNull()
    }
    if obj, ok := dataMap["slug"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Slug = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Slug = types.StringValue(string(jsonBytes))
            } else {
                data.Slug = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Slug = types.StringValue(string(jsonBytes))
            } else {
                data.Slug = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Slug = types.StringValue(string(jsonBytes))
        } else {
            data.Slug = types.StringNull()
        }
    } else if val, ok := dataMap["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
    }
    if obj, ok := dataMap["clientType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ClientType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ClientType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ClientType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ClientType = types.StringValue(string(jsonBytes))
            } else {
                data.ClientType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ClientType = types.StringValue(string(jsonBytes))
            } else {
                data.ClientType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ClientType = types.StringValue(string(jsonBytes))
        } else {
            data.ClientType = types.StringNull()
        }
    } else if val, ok := dataMap["clientType"].(string); ok {
        data.ClientType = types.StringValue(val)
    } else {
        data.ClientType = types.StringNull()
    }
    if obj, ok := dataMap["sdkLanguage"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SdkLanguage = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SdkLanguage = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SdkLanguage = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SdkLanguage = types.StringValue(string(jsonBytes))
            } else {
                data.SdkLanguage = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SdkLanguage = types.StringValue(string(jsonBytes))
            } else {
                data.SdkLanguage = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SdkLanguage = types.StringValue(string(jsonBytes))
        } else {
            data.SdkLanguage = types.StringNull()
        }
    } else if val, ok := dataMap["sdkLanguage"].(string); ok {
        data.SdkLanguage = types.StringValue(val)
    } else {
        data.SdkLanguage = types.StringNull()
    }
    if obj, ok := dataMap["otelCollectorStatus"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OtelCollectorStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.OtelCollectorStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.OtelCollectorStatus = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.OtelCollectorStatus = types.StringValue(string(jsonBytes))
            } else {
                data.OtelCollectorStatus = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.OtelCollectorStatus = types.StringValue(string(jsonBytes))
            } else {
                data.OtelCollectorStatus = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.OtelCollectorStatus = types.StringValue(string(jsonBytes))
        } else {
            data.OtelCollectorStatus = types.StringNull()
        }
    } else if val, ok := dataMap["otelCollectorStatus"].(string); ok {
        data.OtelCollectorStatus = types.StringValue(val)
    } else {
        data.OtelCollectorStatus = types.StringNull()
    }
    if obj, ok := dataMap["agentVersion"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AgentVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AgentVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AgentVersion = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AgentVersion = types.StringValue(string(jsonBytes))
            } else {
                data.AgentVersion = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AgentVersion = types.StringValue(string(jsonBytes))
            } else {
                data.AgentVersion = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AgentVersion = types.StringValue(string(jsonBytes))
        } else {
            data.AgentVersion = types.StringNull()
        }
    } else if val, ok := dataMap["agentVersion"].(string); ok {
        data.AgentVersion = types.StringValue(val)
    } else {
        data.AgentVersion = types.StringNull()
    }
    if obj, ok := dataMap["lastSeenAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.LastSeenAt = NewRFC3339Value(val)
        } else {
            data.LastSeenAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["lastSeenAt"].(string); ok && val != "" {
        data.LastSeenAt = NewRFC3339Value(val)
    } else {
        data.LastSeenAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["sessionReplayLastChunkReceivedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.SessionReplayLastChunkReceivedAt = NewRFC3339Value(val)
        } else {
            data.SessionReplayLastChunkReceivedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["sessionReplayLastChunkReceivedAt"].(string); ok && val != "" {
        data.SessionReplayLastChunkReceivedAt = NewRFC3339Value(val)
    } else {
        data.SessionReplayLastChunkReceivedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["sessionReplayBudgetExceededAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.SessionReplayBudgetExceededAt = NewRFC3339Value(val)
        } else {
            data.SessionReplayBudgetExceededAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["sessionReplayBudgetExceededAt"].(string); ok && val != "" {
        data.SessionReplayBudgetExceededAt = NewRFC3339Value(val)
    } else {
        data.SessionReplayBudgetExceededAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["archivedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.ArchivedAt = NewRFC3339Value(val)
        } else {
            data.ArchivedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["archivedAt"].(string); ok && val != "" {
        data.ArchivedAt = NewRFC3339Value(val)
    } else {
        data.ArchivedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["archivedByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ArchivedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ArchivedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ArchivedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ArchivedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.ArchivedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ArchivedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.ArchivedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ArchivedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.ArchivedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["archivedByUserId"].(string); ok {
        data.ArchivedByUserId = types.StringValue(val)
    } else {
        data.ArchivedByUserId = types.StringNull()
    }
    if obj, ok := dataMap["deletedByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeletedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DeletedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DeletedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DeletedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.DeletedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DeletedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.DeletedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DeletedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.DeletedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["deletedByUserId"].(string); ok {
        data.DeletedByUserId = types.StringValue(val)
    } else {
        data.DeletedByUserId = types.StringNull()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    data.Id = state.Id

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RumApplicationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var data RumApplicationResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Make API call
    httpResp, err := r.client.Delete(ctx, "/rum-application/" + data.Id.ValueString() + "")
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete rum_application, got error: %s", err))
        return
    }

    // A failed delete must keep the resource in state — silently dropping it
    // orphans real infrastructure. 404 means it is already gone.
    if httpResp.StatusCode >= 400 && httpResp.StatusCode != http.StatusNotFound {
        err = r.client.ParseResponse(httpResp, nil)
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to delete rum_application: %s", err))
        return
    }
    if httpResp.Body != nil {
        httpResp.Body.Close()
    }
}


func (r *RumApplicationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Helper method to convert Terraform map to Go interface{}
func (r *RumApplicationResource) convertTerraformMapToInterface(terraformMap types.Map) interface{} {
    if terraformMap.IsNull() || terraformMap.IsUnknown() {
        return nil
    }
    
    result := make(map[string]string)
    terraformMap.ElementsAs(context.Background(), &result, false)
    
    // Convert map[string]string to map[string]interface{}
    interfaceResult := make(map[string]interface{})
    for key, value := range result {
        interfaceResult[key] = value
    }
    
    return interfaceResult
}

// Helper method to convert Terraform list to Go interface{}
func (r *RumApplicationResource) convertTerraformListToInterface(terraformList types.List) interface{} {
    if terraformList.IsNull() || terraformList.IsUnknown() {
        return nil
    }

    var stringList []string
    terraformList.ElementsAs(context.Background(), &stringList, false)

    // Convert string array to OneUptime format with _id fields
    var result []interface{}
    for _, str := range stringList {
        if str != "" {
            result = append(result, map[string]interface{}{
                "_id": str,
            })
        }
    }
    return result
}

// Helper method to convert Terraform set to Go interface{}
func (r *RumApplicationResource) convertTerraformSetToInterface(terraformSet types.Set) interface{} {
    if terraformSet.IsNull() || terraformSet.IsUnknown() {
        return nil
    }

    var stringList []string
    terraformSet.ElementsAs(context.Background(), &stringList, false)

    // Convert string array to OneUptime format with _id fields
    var result []interface{}
    for _, str := range stringList {
        if str != "" {
            result = append(result, map[string]interface{}{
                "_id": str,
            })
        }
    }
    return result
}


// Helper method to parse JSON field for complex objects
func (r *RumApplicationResource) parseJSONField(terraformString basetypes.StringValuable) interface{} {
    sv, _ := terraformString.ToStringValue(context.Background())
    if sv.IsNull() || sv.IsUnknown() || sv.ValueString() == "" {
        return nil
    }

    var result interface{}
    if err := json.Unmarshal([]byte(sv.ValueString()), &result); err != nil {
        // If JSON parsing fails, return the raw string
        return sv.ValueString()
    }

    return result
}

// Normalize URL wrapper objects to avoid drift (e.g., trailing slash differences).
func (r *RumApplicationResource) normalizeURLWrappers(value interface{}) interface{} {
    switch v := value.(type) {
    case map[string]interface{}:
        if typeStr, ok := v["_type"].(string); ok && typeStr == "URL" {
            if val, ok := v["value"].(string); ok {
                v["value"] = r.normalizeURLString(val)
            }
        }
        for key, child := range v {
            v[key] = r.normalizeURLWrappers(child)
        }
        return v
    case []interface{}:
        for i, child := range v {
            v[i] = r.normalizeURLWrappers(child)
        }
        return v
    default:
        return v
    }
}

func (r *RumApplicationResource) normalizeURLString(value string) string {
    parsed, err := url.Parse(value)
    if err != nil {
        return value
    }
    if parsed.Path == "/" && parsed.RawQuery == "" && parsed.Fragment == "" {
        return strings.TrimSuffix(value, "/")
    }
    return value
}

// Helper method to convert *big.Float to float64 for JSON serialization
func (r *RumApplicationResource) bigFloatToFloat64(bf *big.Float) interface{} {
    if bf == nil {
        return nil
    }
    f, _ := bf.Float64()
    return f
}

// Helper method to check if a type string is a valid OneUptime ObjectType.
// The registry itself lives in objecttypes.go, shared across the package.
func (r *RumApplicationResource) isValidOneUptimeObjectType(typeStr string) bool {
    return validOneUptimeObjectTypes[typeStr]
}
