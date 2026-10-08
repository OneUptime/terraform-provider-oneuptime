---
page_title: "oneuptime_ai_insight Data Source - oneuptime"
subcategory: "Other"
description: |-
  A preventive finding from OneUptime AI's deterministic telemetry sensors — new or spiking exceptions, error-log spikes, trace-latency regressions and metric drift — surfaced in a quiet insights inbox that never pages and never opens incidents.
---

# oneuptime_ai_insight (Data Source)

A preventive finding from OneUptime AI's deterministic telemetry sensors — new or spiking exceptions, error-log spikes, trace-latency regressions and metric drift — surfaced in a quiet insights inbox that never pages and never opens incidents.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one ai insight may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_ai_insight" "example" {
  insight_type = "example-insight-type"
}

# Or by id:
data "oneuptime_ai_insight" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `classification` (String) AI triage verdict: code-fault, user-error, expected-denial, infrastructure or unknown. Automatic fix pull requests are only opened for code-fault.
- `detail_markdown` (String) The deterministic evidence rendered as markdown: real counts, baselines and multipliers written by the detector at detect time.
- `fingerprint` (String) The detector's stable dedupe key for this finding. Recurring detections refresh the existing non-terminal insight with the same fingerprint.
- `fix_ai_run_id` (String) The AI agent fix task queued for this insight (a CodeFix AIRun that opens a pull request, ready for review).
- `human_verdict` (String) The one-click human verdict on this insight (Confirmed or Dismissed). Null until a user weighs in.
- `human_verdict_by_user_id` (String) The user who recorded (or last changed) the human verdict.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `insight_type` (String) Which deterministic detector produced this insight: NewException, ExceptionSpike, ErrorLogSpike, TraceLatencyRegression or MetricDrift.
- `metric_name` (String) The drifting metric's name (for MetricDrift insights).
- `occurrence_count` (Number) How many scanner ticks have detected this finding. Incremented on each dedupe refresh.
- `service_name` (String) Name of the telemetry service this insight is about.
- `severity` (String) How urgent this insight is (High, Medium or Low), assigned deterministically by the detector.
- `status` (String) Lifecycle of the insight. Detected is the defensive initial state — the scanner routes to ActionRequired or FixOpened in the same tick; Resolved and Dismissed are human actions.
- `telemetry_exception_id` (String) The telemetry exception behind this insight (for NewException and ExceptionSpike insights).
- `telemetry_service_id` (String) ID of the telemetry service this insight is about.
- `title` (String) One-line human-readable summary of the finding.
- `trace_id` (String) A representative slow trace (for TraceLatencyRegression insights).
- `triage_ai_run_id` (String) The budgeted, read-only AI triage run enqueued for this insight (an Investigation AIRun).
- `triage_summary_markdown` (String) The AI triage analysis for this insight: probable root cause, blast radius and suggested action, with citations.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `evidence` (String) The deterministic evidence computed at detect time: counts, baselines, multipliers and (for latency insights) span-tree findings. A JSON value: write it with `jsonencode()`.
- `first_seen_at` (String) When this finding was first detected.
- `human_verdict_at` (String) When the human verdict was recorded (or last changed).
- `last_seen_at` (String) When this finding was most recently re-detected by the scanner.
- `project_id` (String) ID of the project this insight belongs to. The ID of a `oneuptime_project`.
- `triage_completed_at` (String) When the AI triage analysis completed.
- `updated_at` (String) Date and Time when the object was updated.
