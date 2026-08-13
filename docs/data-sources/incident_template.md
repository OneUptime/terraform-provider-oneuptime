---
page_title: "oneuptime_incident_template Data Source - oneuptime"
subcategory: "Incidents"
description: |-
  Manage incident templates for your project
---

# oneuptime_incident_template (Data Source)

Manage incident templates for your project Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_incident_template" "by_name" {
  name = "example-incident_template"
}

data "oneuptime_incident_template" "by_id" {
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
- `title` (String) Title of this incident.. Computed.
- `template_name` (String) Name of the Incident Template.. Computed.
- `template_description` (String) Description of the Incident Template.. Computed.
- `description` (String) Short description of this incident. This is in markdown and will be visible on the status page... Computed.
- `slug` (String) Friendly globally unique name for your object.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `monitors` (Set) List of monitors affected by this incident.. Computed.
- `hosts` (Set) List of hosts to pre-populate on incidents created from this template... Computed.
- `kubernetes_clusters` (Set) List of Kubernetes clusters to pre-populate on incidents created from this template... Computed.
- `docker_hosts` (Set) List of Docker hosts to pre-populate on incidents created from this template... Computed.
- `podman_hosts` (Set) List of Podman hosts to pre-populate on incidents created from this template... Computed.
- `services` (Set) List of services to pre-populate on incidents created from this template... Computed.
- `on_call_duty_policies` (Set) List of on-call duty policies affected by this incident template... Computed.
- `labels` (Set) Relation to Labels Array where this object is categorized in... Computed.
- `incident_severity_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `change_monitor_status_to_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `initial_incident_state_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `custom_fields` (String) Custom Fields on this resource... Computed.
