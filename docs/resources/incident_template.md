---
page_title: "oneuptime_incident_template Resource - oneuptime"
subcategory: "Incidents"
description: |-
  Manage incident templates for your project
---

# oneuptime_incident_template (Resource)

Manage incident templates for your project

## Example Usage

```terraform
resource "oneuptime_incident_template" "example" {
  title = "This is an example of longer text content that might be stored in this field."
  template_name = "Example short text"
  template_description = "This is an example of longer text content that might be stored in this field."
  description = "# Heading

This is **markdown** content"
}
```

## Schema

### Required

- `title` (String) Title of this incident..
- `template_name` (String) Name of the Incident Template..
- `template_description` (String) Description of the Incident Template..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Short description of this incident. This is in markdown and will be visible on the status page...
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `monitors` (Set) List of monitors affected by this incident..
- `hosts` (Set) List of hosts to pre-populate on incidents created from this template...
- `kubernetes_clusters` (Set) List of Kubernetes clusters to pre-populate on incidents created from this template...
- `docker_hosts` (Set) List of Docker hosts to pre-populate on incidents created from this template...
- `podman_hosts` (Set) List of Podman hosts to pre-populate on incidents created from this template...
- `services` (Set) List of services to pre-populate on incidents created from this template...
- `on_call_duty_policies` (Set) List of on-call duty policies affected by this incident template...
- `labels` (Set) Relation to Labels Array where this object is categorized in...
- `incident_severity_id` (String) A unique identifier for an object, represented as a UUID..
- `change_monitor_status_to_id` (String) A unique identifier for an object, represented as a UUID..
- `initial_incident_state_id` (String) A unique identifier for an object, represented as a UUID..
- `custom_fields` (String) Custom Fields on this resource...

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `slug` (String) Friendly globally unique name for your object..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_incident_template.example <id>
```
