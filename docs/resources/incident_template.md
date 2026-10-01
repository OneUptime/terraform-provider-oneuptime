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
- `status_pages` (Set) Limit incidents declared from this template to these status pages. Leave empty to reach every status page that lists the incident's monitors...
- `labels` (Set) Relation to Labels Array where this object is categorized in...
- `incident_severity_id` (String) A unique identifier for an object, represented as a UUID..
- `change_monitor_status_to_id` (String) A unique identifier for an object, represented as a UUID..
- `initial_incident_state_id` (String) A unique identifier for an object, represented as a UUID..
- `custom_fields` (String) The custom field values incidents declared from this template start with, keyed by each incident custom field's name. They are merged one field at a time under the values the request or the Declare Incident form supplies...
- `custom_field_settings` (String) How the Declare Incident form treats each incident custom field when an incident is declared from this template, keyed by the field's template variable key (variableKey). Each value is Required (asked, must be filled in), Optional (asked, may be left empty), Hidden (not asked; the field keeps this template's value) or Default. A field that is not listed, or is Default, follows its own Show on Create and Required on Create settings. Only the dashboard's Declare Incident form applies these settings: incidents created through the API are not checked against them...

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `slug` (String) Friendly globally unique name for your object..
- `is_scoped_to_status_pages` (Bool) Whether incidents declared from this template are limited to the status pages in Status Pages. Derived from Status Pages; any value sent for it is ignored...

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_incident_template.example <id>
```
