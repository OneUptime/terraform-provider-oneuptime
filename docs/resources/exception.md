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

- `primary_entity_id` (String) ID of the resource this exception belongs to (Service / Host / DockerHost / KubernetesCluster, or the projectId for unattributed telemetry — disambiguated by primaryEntityType).

### Optional

- `assign_to_team_id` (String) Team ID who this exception is assigned to. The ID of a `oneuptime_team`.
- `assign_to_user_id` (String) User ID who this exception is assigned to. The ID of a `oneuptime_user` (see the data source).
- `environment` (String) Deployment environment from deployment.environment resource attribute.
- `exception_type` (String) Type of the exception that was thrown by the telemetry service.
- `fingerprint` (String) Finger print of the exception that was thrown by the telemetry service.
- `first_seen_at` (String) When did this team member accept invitation.
- `first_seen_in_release` (String) The service version / release in which this exception was first observed.
- `is_archived` (Boolean) Is this exception archived? Defaults to `false`.
- `is_resolved` (Boolean) Is this exception resolved? Defaults to `false`.
- `last_seen_at` (String) When did this team member accept invitation.
- `last_seen_in_release` (String) The most recent service version / release in which this exception was observed.
- `message` (String) Exception message that was thrown by the telemetry service.
- `occurance_count` (Number) Number of times this exception has occurred. Defaults to `1`.
- `primary_entity_type` (String) Resource type that produced this exception (e.g. OpenTelemetry service, Host, DockerHost, KubernetesCluster, or Unknown for unattributed telemetry).
- `stack_trace` (String) Stack trace of the exception that was thrown by the telemetry service.

### Read-Only

- `ai_classification` (String) AI triage verdict for this exception group (code-fault, user-error, expected-denial, infrastructure).
- `ai_fix_declined_at` (String) Set when an AI-authored fix pull request for this exception was closed without merging; suppresses further automatic fix attempts.
- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `error_class` (String) Fault domain of this exception group (code-fault, user-error, expected-denial, infrastructure, unknown). Non-actionable classes are excluded from the Issues list.
- `error_class_source` (String) Where the error class came from: default (unclassified), declared (by the emitting code), ai (triage verdict) or manual (a human).
- `id` (String) Unique identifier for the resource.
- `marked_as_archived_at` (String) When did this team member accept invitation.
- `marked_as_archived_by_user_id` (String) User ID who marked this exception as archived. The ID of a `oneuptime_user` (see the data source).
- `marked_as_resolved_at` (String) When did this team member accept invitation.
- `marked_as_resolved_by_user_id` (String) User ID who marked this exception as resolved. The ID of a `oneuptime_user` (see the data source).
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `unhandled` (Boolean) True when at least one occurrence of this exception escaped its span scope (was unhandled, per OTel exception.escaped).
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing exception by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_exception.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_exception.example <id>
```
