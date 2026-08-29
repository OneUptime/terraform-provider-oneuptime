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
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
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
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
}

func (d *RumApplicationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_rum_application"
}

func (d *RumApplicationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Browser & mobile applications auto-discovered from OpenTelemetry RUM telemetry (browser.* / device.* resource attributes). One row per application, aggregating all end-user clients. Look up an existing rum_application by `id` or by `name`.",

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
            "slug": schema.StringAttribute{
                MarkdownDescription: "Friendly globally unique name for your object.",
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Friendly description that will help you remember.",
                Computed: true,
            },
            "app_identifier": schema.StringAttribute{
                MarkdownDescription: "Stable identifier for this application from the service.name OpenTelemetry resource attribute. Identity key for this RUM application..",
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
                Computed: true,
            },
            "labels": schema.SetAttribute{
                MarkdownDescription: "Relation to Labels Array where this object is categorized in..",
                Computed: true,
                ElementType: types.StringType,
            },
            "retain_telemetry_data_for_days": schema.NumberAttribute{
                MarkdownDescription: "Number of days to retain telemetry data for this application. Leave blank to use the project-wide default..",
                Computed: true,
            },
            "telemetry_retention_config": schema.StringAttribute{
                MarkdownDescription: "Per-pillar retention overrides for this application. Unset fields fall back to the application default, then the project's retention settings..",
                Computed: true,
            },
            "is_session_replay_enabled": schema.BoolAttribute{
                MarkdownDescription: "When enabled, the browser recorder may record and upload session replays for this application. On by default; Project.isSessionReplayAllowed must also be on. Turn it off here to stop recording for one application without affecting the rest of the project..",
                Computed: true,
            },
            "session_replay_masking_mode": schema.StringAttribute{
                MarkdownDescription: "How aggressively the recorder masks page content before it leaves the end user's device. MaskSensitiveInputsOnly (default) masks passwords and card / one-time-code fields and records everything else verbatim. MaskInputsOnly additionally masks every other input value. MaskAllText masks static page text too, producing a wireframe..",
                Computed: true,
            },
            "session_replay_mask_selectors": schema.StringAttribute{
                MarkdownDescription: "CSS selectors whose text content the recorder masks, in addition to whatever the masking mode already covers..",
                Computed: true,
            },
            "session_replay_block_selectors": schema.StringAttribute{
                MarkdownDescription: "CSS selectors the recorder excludes from the DOM snapshot entirely, so the subtree is never captured rather than captured and masked..",
                Computed: true,
            },
            "session_replay_ignore_error_patterns": schema.StringAttribute{
                MarkdownDescription: "Regex patterns matched against an uncaught error's message and source URL. Matching errors are still recorded in the session but no longer trigger an upload — the remedy for a chronically-throwing third-party tag that would otherwise convert error-triggered capture into always-on recording..",
                Computed: true,
            },
            "session_replay_trace_propagation_origins": schema.StringAttribute{
                MarkdownDescription: "Origins the recorder may inject a W3C traceparent header into, linking recordings to the backend traces of their requests without any OpenTelemetry browser setup. Empty means never inject: adding a header makes cross-origin requests preflighted, so each listed origin is an explicit statement that its API allows the traceparent header..",
                Computed: true,
            },
            "session_replay_lcp_budget_ms": schema.NumberAttribute{
                MarkdownDescription: "Largest Contentful Paint budget in milliseconds. A session whose LCP exceeds it uploads with the Performance trigger. 0 disables the trigger..",
                Computed: true,
            },
            "session_replay_long_task_budget_ms": schema.NumberAttribute{
                MarkdownDescription: "Main-thread long-task budget in milliseconds. A single task blocking longer than this uploads the session with the Performance trigger. 0 disables the trigger..",
                Computed: true,
            },
            "session_replay_slow_request_budget_ms": schema.NumberAttribute{
                MarkdownDescription: "Request duration budget in milliseconds. An instrumented request slower than this uploads the session with the Performance trigger. 0 disables the trigger..",
                Computed: true,
            },
            "session_replay_allowed_origins": schema.StringAttribute{
                MarkdownDescription: "Exact browser origins (scheme + host + port) allowed to upload session replay chunks for this application. Empty (the default) accepts any origin. Once you list an origin this becomes a strict allowlist: anything unlisted, and any request with no Origin header, is refused..",
                Computed: true,
            },
            "session_replay_consent_mode": schema.StringAttribute{
                MarkdownDescription: "NotRequired (default) uploads immediately, asserting a lawful basis that does not need a per-session grant. RequireExplicit buffers in memory and uploads nothing until the host page calls grantConsent(); set it if you need a per-session consent handshake, which most EU deployments will..",
                Computed: true,
            },
            "session_replay_capture_trigger": schema.StringAttribute{
                MarkdownDescription: "OnErrorOrFrustration (default) keeps a rolling in-memory buffer and uploads only when something actually went wrong. Always uploads every sampled session from its first event, which costs materially more and stores materially more end-user data..",
                Computed: true,
            },
            "session_replay_sample_percentage": schema.NumberAttribute{
                MarkdownDescription: "Percentage of sessions (0 to 100) recorded regardless of whether anything went wrong. 0 by default, so with the default trigger only failing sessions are recorded..",
                Computed: true,
            },
            "session_replay_capture_user_identity": schema.BoolAttribute{
                MarkdownDescription: "When enabled, the raw end-user reference supplied by the host page is stored alongside the recording, so a support engineer can find the session a named customer is complaining about. When off, only a one-way per-project HMAC of it is stored. On by default. Narrower create/update ACL than the other replay settings: this is the switch that turns a pseudonymous recording into an identified one..",
                Computed: true,
            },
            "session_replay_capture_geo": schema.BoolAttribute{
                MarkdownDescription: "When enabled, a country code is derived from the request and stored on the session. On by default. The end user's IP address is never stored either way - the country is the only geographic fact this keeps..",
                Computed: true,
            },
            "session_replay_record_canvas": schema.BoolAttribute{
                MarkdownDescription: "When enabled, canvas contents are recorded. Off by default because canvas capture is expensive on the end user's device and canvases routinely render content the text masking cannot reach..",
                Computed: true,
            },
            "session_replay_retention_in_days": schema.NumberAttribute{
                MarkdownDescription: "How long session recordings are kept for this application. Clamped to 1, 7, 14, 30 or 90 days. Defaults to 7 rather than the 15 the other telemetry pillars use, because a short retention is itself a privacy control..",
                Computed: true,
            },
            "session_replay_monthly_budget_in_gb": schema.NumberAttribute{
                MarkdownDescription: "Optional ceiling on replay bytes ingested per calendar month for this application. Once exceeded, live recorders are told to stop. Leave blank for no application-level ceiling..",
                Computed: true,
            },
            "session_replay_last_chunk_received_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "session_replay_budget_exceeded_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "is_archived": schema.BoolAttribute{
                MarkdownDescription: "Is this RUM application archived? Archived RUM applications are hidden from lists but keep collecting telemetry..",
                Computed: true,
            },
            "archived_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
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

    hasId := !data.Id.IsNull() && data.Id.ValueString() != ""
    hasName := !data.Name.IsNull() && data.Name.ValueString() != ""
    if hasId == hasName {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Exactly one of `id` or `name` must be set to look up a rum_application.",
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
        "deletedByUserId": true,
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
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No rum_application found with id %q.", data.Id.ValueString()))
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
    } else {
        listBody := map[string]interface{}{
            "query": map[string]interface{}{
                "name": data.Name.ValueString(),
            },
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
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No rum_application found with name %q.", data.Name.ValueString()))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one rum_application matches name %q. Use the id attribute to disambiguate.", data.Name.ValueString()))
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
