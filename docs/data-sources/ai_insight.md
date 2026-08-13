---
page_title: "oneuptime_ai_insight Data Source - oneuptime"
subcategory: "Other"
description: |-
  A preventive finding from OneUptime AI's deterministic telemetry sensors — new or spiking exceptions, error-log spikes, trace-latency regressions and metric drift — surfaced in a quiet insights inbox that never pages and never opens incidents.
---

# oneuptime_ai_insight (Data Source)

A preventive finding from OneUptime AI's deterministic telemetry sensors — new or spiking exceptions, error-log spikes, trace-latency regressions and metric drift — surfaced in a quiet insights inbox that never pages and never opens incidents. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_ai_insight" "by_name" {
  name = "example-ai_insight"
}

data "oneuptime_ai_insight" "by_id" {
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
- `insight_type` (String) Which deterministic detector produced this insight: NewException, ExceptionSpike, ErrorLogSpike, TraceLatencyRegression or MetricDrift... Computed.
- `status` (String) Lifecycle of the insight. Detected is the defensive initial state — the scanner routes to ActionRequired or FixOpened in the same tick; Resolved and Dismissed are human actions... Computed.
- `severity` (String) How urgent this insight is (High, Medium or Low), assigned deterministically by the detector... Computed.
- `classification` (String) AI triage verdict: code-fault, user-error, expected-denial, infrastructure or unknown. Automatic fix pull requests are only opened for code-fault... Computed.
- `fingerprint` (String) The detector's stable dedupe key for this finding. Recurring detections refresh the existing non-terminal insight with the same fingerprint... Computed.
- `title` (String) One-line human-readable summary of the finding... Computed.
- `detail_markdown` (String) The deterministic evidence rendered as markdown: real counts, baselines and multipliers written by the detector at detect time... Computed.
- `service_name` (String) Name of the telemetry service this insight is about... Computed.
- `telemetry_service_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `telemetry_exception_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `trace_id` (String) A representative slow trace (for TraceLatencyRegression insights)... Computed.
- `metric_name` (String) The drifting metric's name (for MetricDrift insights)... Computed.
- `evidence` (String) The deterministic evidence computed at detect time: counts, baselines, multipliers and (for latency insights) span-tree findings... Computed.
- `first_seen_at` (String) A date time object.. Computed.
- `last_seen_at` (String) A date time object.. Computed.
- `occurrence_count` (Number) How many scanner ticks have detected this finding. Incremented on each dedupe refresh... Computed.
- `triage_ai_run_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `fix_ai_run_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `triage_summary_markdown` (String) The AI triage analysis for this insight: probable root cause, blast radius and suggested action, with citations... Computed.
- `triage_completed_at` (String) A date time object.. Computed.
- `human_verdict` (String) The one-click human verdict on this insight (Confirmed or Dismissed). Null until a user weighs in... Computed.
- `human_verdict_at` (String) A date time object.. Computed.
- `human_verdict_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
