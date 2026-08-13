---
page_title: "oneuptime_exception Data Source - oneuptime"
subcategory: "Telemetry & Dashboards"
description: |-
  List of all Telemetry Exceptions created for the telemetry service for this OneUptime project and it's status.
---

# oneuptime_exception (Data Source)

List of all Telemetry Exceptions created for the telemetry service for this OneUptime project and it's status. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_exception" "by_name" {
  name = "example-exception"
}

data "oneuptime_exception" "by_id" {
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
- `primary_entity_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `primary_entity_type` (String) Resource type that produced this exception (e.g. OpenTelemetry service, Host, DockerHost, KubernetesCluster, or Unknown for unattributed telemetry)... Computed.
- `message` (String) Exception message that was thrown by the telemetry service.. Computed.
- `stack_trace` (String) Stack trace of the exception that was thrown by the telemetry service.. Computed.
- `exception_type` (String) Type of the exception that was thrown by the telemetry service.. Computed.
- `fingerprint` (String) Finger print of the exception that was thrown by the telemetry service.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `marked_as_resolved_at` (String) A date time object.. Computed.
- `marked_as_archived_at` (String) A date time object.. Computed.
- `first_seen_at` (String) A date time object.. Computed.
- `last_seen_at` (String) A date time object.. Computed.
- `assign_to_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `assign_to_team_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `marked_as_resolved_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `marked_as_archived_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `is_resolved` (Bool) Is this exception resolved?.. Computed.
- `is_archived` (Bool) Is this exception archived?.. Computed.
- `occurance_count` (Number) Number of times this exception has occurred.. Computed.
- `first_seen_in_release` (String) The service version / release in which this exception was first observed.. Computed.
- `last_seen_in_release` (String) The most recent service version / release in which this exception was observed.. Computed.
- `environment` (String) Deployment environment from deployment.environment resource attribute.. Computed.
- `unhandled` (Bool) True when at least one occurrence of this exception escaped its span scope (was unhandled, per OTel exception.escaped).. Computed.
- `ai_classification` (String) AI triage verdict for this exception group (code-fault, user-error, expected-denial, infrastructure).. Computed.
- `ai_fix_declined_at` (String) A date time object.. Computed.
