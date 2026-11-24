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
- `assign_to_user_email` (String) The email address of the person in `assign_to_user_id`, instead of their ID: someone in this project, or invited to it. Give one of the two.
- `assign_to_user_id` (String) User ID who this exception is assigned to. The ID of a `oneuptime_user` (see the data source). Or name them by email address, in `assign_to_user_email`.
- `environment` (String) Deployment environment from deployment.environment resource attribute. OneUptime keeps it up to date once the record exists: the value given here is the one the record starts with, and a later change of it is never sent.
- `exception_type` (String) Type of the exception that was thrown by the telemetry service. OneUptime keeps it up to date once the record exists: the value given here is the one the record starts with, and a later change of it is never sent.
- `fingerprint` (String) Finger print of the exception that was thrown by the telemetry service.
- `first_seen_at` (String) When the earliest occurrence of this exception happened. OneUptime keeps it up to date once the record exists: the value given here is the one the record starts with, and a later change of it is never sent.
- `first_seen_in_release` (String) The service version / release in which this exception was first observed. OneUptime keeps it up to date once the record exists: the value given here is the one the record starts with, and a later change of it is never sent.
- `is_archived` (Boolean) Is this exception archived? OneUptime changes it too, on its own: left out, it keeps the value OneUptime gives it. Defaults to `false`.
- `is_resolved` (Boolean) Is this exception resolved? OneUptime changes it too, on its own: left out, it keeps the value OneUptime gives it. Defaults to `false`.
- `last_seen_at` (String) When the latest occurrence of this exception happened. OneUptime keeps it up to date once the record exists: the value given here is the one the record starts with, and a later change of it is never sent.
- `last_seen_in_release` (String) The most recent service version / release in which this exception was observed. OneUptime keeps it up to date once the record exists: the value given here is the one the record starts with, and a later change of it is never sent.
- `message` (String) Exception message that was thrown by the telemetry service. OneUptime keeps it up to date once the record exists: the value given here is the one the record starts with, and a later change of it is never sent.
- `occurance_count` (Number) Number of times this exception has occurred. OneUptime keeps it up to date once the record exists: the value given here is the one the record starts with, and a later change of it is never sent. Defaults to `1`.
- `primary_entity_type` (String) Resource type that produced this exception (e.g. OpenTelemetry service, Host, DockerHost, KubernetesCluster, or Unknown for unattributed telemetry). Allowed values: `OpenTelemetry`, `Monitor`, `NetworkDevice`, `Alert`, `Incident`, `ScheduledMaintenance`, `ServiceLevelObjective`, `Host`, `DockerHost`, `PodmanHost`, `KubernetesCluster`, `ProxmoxCluster`, `OpenStackCloud`, `CephCluster`, `StorageArray`, `DockerSwarmCluster`, `VMwareVCenter`, `IoTDevice`, `ServerlessFunction`, `CloudResource`, `RealUserMonitor`, `DatabaseServer`, `DataStore`, `DataPipeline`, `DataProduct`, `DataReconciliation`, `DataBiConnection`, `DeviceManagementConnection`, `Unknown`.
- `resolved_in_release_after` (String) Set when the exception was resolved in the next release: the release it was resolved after. Occurrences from this release do not reopen the exception; an occurrence from any other release does. OneUptime changes it too, on its own: left out, it keeps the value OneUptime gives it.
- `stack_trace` (String) Stack trace of the exception that was thrown by the telemetry service. OneUptime keeps it up to date once the record exists: the value given here is the one the record starts with, and a later change of it is never sent.

### Read-Only

- `ai_classification` (String) AI triage verdict for this exception group (code-fault, user-error, expected-denial, infrastructure).
- `ai_explanation` (String) OneUptime AI's explanation of this exception: a summary, the likely root cause, a suggested fix, how confident it is, the suspect stack frames and next steps. A JSON value: read it with `jsondecode()`.
- `ai_explanation_generated_at` (String) When OneUptime AI last wrote its explanation of this exception.
- `ai_fix_declined_at` (String) Set when an AI-authored fix pull request for this exception was closed without merging; suppresses further automatic fix attempts.
- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `error_class` (String) Fault domain of this exception group (code-fault, user-error, expected-denial, infrastructure, unknown). Non-actionable classes are excluded from the Issues list.
- `error_class_source` (String) Where the error class came from: default (unclassified), declared (by the emitting code), ai (triage verdict) or manual (a human).
- `id` (String) Unique identifier for the resource.
- `marked_as_archived_at` (String) When this exception was last archived. Empty while it is not archived.
- `marked_as_archived_by_user_id` (String) User ID who marked this exception as archived. The ID of a `oneuptime_user` (see the data source).
- `marked_as_resolved_at` (String) When this exception was last resolved. Empty while it is unresolved.
- `marked_as_resolved_by_user_id` (String) User ID who marked this exception as resolved. The ID of a `oneuptime_user` (see the data source).
- `merged_at` (String) When this exception was merged into another.
- `merged_by_user_id` (String) ID of the user who merged this exception into another. The ID of a `oneuptime_user` (see the data source).
- `merged_into_telemetry_exception_id` (String) ID of the exception this one was merged into as the same bug. Empty for an exception that was not merged. The ID of a `oneuptime_exception`.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `regressed_at` (String) When this exception last came back after it was resolved. Empty when it has not come back since someone last resolved or reopened it.
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
