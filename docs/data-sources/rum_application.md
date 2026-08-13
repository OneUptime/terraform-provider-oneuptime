---
page_title: "oneuptime_rum_application Data Source - oneuptime"
subcategory: "Other"
description: |-
  Browser & mobile applications auto-discovered from OpenTelemetry RUM telemetry (browser.* / device.* resource attributes). One row per application, aggregating all end-user clients.
---

# oneuptime_rum_application (Data Source)

Browser & mobile applications auto-discovered from OpenTelemetry RUM telemetry (browser.* / device.* resource attributes). One row per application, aggregating all end-user clients. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_rum_application" "by_name" {
  name = "example-rum_application"
}

data "oneuptime_rum_application" "by_id" {
  id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

- `id` (String) Look up by unique identifier. Exactly one of `id` or `name` must be set.. Computed.
- `name` (String) Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.. Computed.
- `created_at` (String) A date time object.. Computed.
- `updated_at` (String) A date time object.. Computed.
- `deleted_at` (String) A date time object.. Computed.
- `version` (Number) Object version. Computed.
- `project_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `slug` (String) Friendly globally unique name for your object.. Computed.
- `description` (String) Friendly description that will help you remember.. Computed.
- `app_identifier` (String) Stable identifier for this application from the service.name OpenTelemetry resource attribute. Identity key for this RUM application... Computed.
- `client_type` (String) Whether this application's clients are browsers or mobile devices (browser / mobile), derived from browser.* / device.* attributes... Computed.
- `sdk_language` (String) Last-seen telemetry.sdk.language resource attribute (e.g. webjs, swift, android). Used to scope this application's client telemetry apart from a same-named backend service... Computed.
- `otel_collector_status` (String) Whether telemetry is currently being received (connected) or has gone stale (disconnected)... Computed.
- `agent_version` (String) Version of the OpenTelemetry SDK reporting this application... Computed.
- `last_seen_at` (String) A date time object.. Computed.
- `labels` (Set) Relation to Labels Array where this object is categorized in... Computed.
- `retain_telemetry_data_for_days` (Number) Number of days to retain telemetry data for this application. Leave blank to use the project-wide default... Computed.
- `telemetry_retention_config` (String) Per-pillar retention overrides for this application. Unset fields fall back to the application default, then the project's retention settings... Computed.
- `is_session_replay_enabled` (Bool) When enabled, the browser recorder may record and upload session replays for this application. On by default; Project.isSessionReplayAllowed must also be on. Turn it off here to stop recording for one application without affecting the rest of the project... Computed.
- `session_replay_masking_mode` (String) How aggressively the recorder masks page content before it leaves the end user's device. MaskSensitiveInputsOnly (default) masks passwords and card / one-time-code fields and records everything else verbatim. MaskInputsOnly additionally masks every other input value. MaskAllText masks static page text too, producing a wireframe... Computed.
- `session_replay_mask_selectors` (String) CSS selectors whose text content the recorder masks, in addition to whatever the masking mode already covers... Computed.
- `session_replay_block_selectors` (String) CSS selectors the recorder excludes from the DOM snapshot entirely, so the subtree is never captured rather than captured and masked... Computed.
- `session_replay_ignore_error_patterns` (String) Regex patterns matched against an uncaught error's message and source URL. Matching errors are still recorded in the session but no longer trigger an upload — the remedy for a chronically-throwing third-party tag that would otherwise convert error-triggered capture into always-on recording... Computed.
- `session_replay_trace_propagation_origins` (String) Origins the recorder may inject a W3C traceparent header into, linking recordings to the backend traces of their requests without any OpenTelemetry browser setup. Empty means never inject: adding a header makes cross-origin requests preflighted, so each listed origin is an explicit statement that its API allows the traceparent header... Computed.
- `session_replay_lcp_budget_ms` (Number) Largest Contentful Paint budget in milliseconds. A session whose LCP exceeds it uploads with the Performance trigger. 0 disables the trigger... Computed.
- `session_replay_long_task_budget_ms` (Number) Main-thread long-task budget in milliseconds. A single task blocking longer than this uploads the session with the Performance trigger. 0 disables the trigger... Computed.
- `session_replay_slow_request_budget_ms` (Number) Request duration budget in milliseconds. An instrumented request slower than this uploads the session with the Performance trigger. 0 disables the trigger... Computed.
- `session_replay_allowed_origins` (String) Exact browser origins (scheme + host + port) allowed to upload session replay chunks for this application. Empty (the default) accepts any origin. Once you list an origin this becomes a strict allowlist: anything unlisted, and any request with no Origin header, is refused... Computed.
- `session_replay_consent_mode` (String) NotRequired (default) uploads immediately, asserting a lawful basis that does not need a per-session grant. RequireExplicit buffers in memory and uploads nothing until the host page calls grantConsent(); set it if you need a per-session consent handshake, which most EU deployments will... Computed.
- `session_replay_capture_trigger` (String) OnErrorOrFrustration (default) keeps a rolling in-memory buffer and uploads only when something actually went wrong. Always uploads every sampled session from its first event, which costs materially more and stores materially more end-user data... Computed.
- `session_replay_sample_percentage` (Number) Percentage of sessions (0 to 100) recorded regardless of whether anything went wrong. 0 by default, so with the default trigger only failing sessions are recorded... Computed.
- `session_replay_capture_user_identity` (Bool) When enabled, the raw end-user reference supplied by the host page is stored alongside the recording, so a support engineer can find the session a named customer is complaining about. When off, only a one-way per-project HMAC of it is stored. On by default. Narrower create/update ACL than the other replay settings: this is the switch that turns a pseudonymous recording into an identified one... Computed.
- `session_replay_capture_geo` (Bool) When enabled, a country code is derived from the request and stored on the session. On by default. The end user's IP address is never stored either way - the country is the only geographic fact this keeps... Computed.
- `session_replay_record_canvas` (Bool) When enabled, canvas contents are recorded. Off by default because canvas capture is expensive on the end user's device and canvases routinely render content the text masking cannot reach... Computed.
- `session_replay_retention_in_days` (Number) How long session recordings are kept for this application. Clamped to 1, 7, 14, 30 or 90 days. Defaults to 7 rather than the 15 the other telemetry pillars use, because a short retention is itself a privacy control... Computed.
- `session_replay_monthly_budget_in_gb` (Number) Optional ceiling on replay bytes ingested per calendar month for this application. Once exceeded, live recorders are told to stop. Leave blank for no application-level ceiling... Computed.
- `session_replay_last_chunk_received_at` (String) A date time object.. Computed.
- `session_replay_budget_exceeded_at` (String) A date time object.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `is_archived` (Bool) Is this RUM application archived? Archived RUM applications are hidden from lists but keep collecting telemetry... Computed.
- `archived_at` (String) A date time object.. Computed.
- `archived_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
