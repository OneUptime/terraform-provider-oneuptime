---
page_title: "oneuptime_exception Resource - oneuptime"
subcategory: "Telemetry & Dashboards"
description: |-
  List of all Telemetry Exceptions created for the telemetry service for this OneUptime project and it's status.
---

# oneuptime_exception (Resource)

List of all Telemetry Exceptions created for the telemetry service for this OneUptime project and it's status.

## Example Usage

```terraform
resource "oneuptime_exception" "example" {
  primary_entity_id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

### Required

- `primary_entity_id` (String) A unique identifier for an object, represented as a UUID..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `primary_entity_type` (String) Resource type that produced this exception (e.g. OpenTelemetry service, Host, DockerHost, KubernetesCluster, or Unknown for unattributed telemetry)...
- `message` (String) Exception message that was thrown by the telemetry service..
- `stack_trace` (String) Stack trace of the exception that was thrown by the telemetry service..
- `exception_type` (String) Type of the exception that was thrown by the telemetry service..
- `fingerprint` (String) Finger print of the exception that was thrown by the telemetry service..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `marked_as_resolved_at` (String) A date time object..
- `marked_as_archived_at` (String) A date time object..
- `first_seen_at` (String) A date time object..
- `last_seen_at` (String) A date time object..
- `assign_to_user_id` (String) A unique identifier for an object, represented as a UUID..
- `assign_to_team_id` (String) A unique identifier for an object, represented as a UUID..
- `marked_as_resolved_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `marked_as_archived_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `is_resolved` (Bool) Is this exception resolved?..
- `is_archived` (Bool) Is this exception archived?..
- `occurance_count` (Number) Number of times this exception has occurred..
- `first_seen_in_release` (String) The service version / release in which this exception was first observed..
- `last_seen_in_release` (String) The most recent service version / release in which this exception was observed..
- `environment` (String) Deployment environment from deployment.environment resource attribute..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `unhandled` (Bool) True when at least one occurrence of this exception escaped its span scope (was unhandled, per OTel exception.escaped)..
- `ai_classification` (String) AI triage verdict for this exception group (code-fault, user-error, expected-denial, infrastructure)..
- `error_class` (String) Fault domain of this exception group (code-fault, user-error, expected-denial, infrastructure, unknown). Non-actionable classes are excluded from the Issues list...
- `error_class_source` (String) Where the error class came from: default (unclassified), declared (by the emitting code), ai (triage verdict) or manual (a human)..
- `ai_fix_declined_at` (String) A date time object..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_exception.example <id>
```
