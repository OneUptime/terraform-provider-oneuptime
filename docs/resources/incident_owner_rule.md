---
page_title: "oneuptime_incident_owner_rule Resource - oneuptime"
subcategory: "Incidents"
description: |-
  Configure rules for automatically assigning owner users and teams when matching incidents are created
---

# oneuptime_incident_owner_rule (Resource)

Configure rules for automatically assigning owner users and teams when matching incidents are created

## Example Usage

```terraform
resource "oneuptime_incident_owner_rule" "example" {
  name = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Name of this incident owner rule..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Description of this incident owner rule..
- `is_enabled` (Bool) Whether this rule is enabled..
- `notify_owners` (Bool) Send notifications to owner users and teams when they are added by this rule..
- `monitors` (Set) Only trigger for incidents from these monitors. Leave empty to match incidents from any monitor...
- `incident_severities` (Set) Only trigger for incidents with these severities. Leave empty to match incidents of any severity...
- `incident_labels` (Set) Only trigger for incidents that have at least one of these labels. Leave empty to match regardless of incident labels...
- `monitor_labels` (Set) Only trigger for incidents from monitors that have at least one of these labels. Leave empty to match regardless of monitor labels...
- `incident_title_pattern` (String) Regex (case-insensitive) matched against the incident title. Leave empty to match any title...
- `incident_description_pattern` (String) Regex (case-insensitive) matched against the incident description. Leave empty to match any description...
- `monitor_name_pattern` (String) Regex (case-insensitive) matched against any of the incident's monitor names. Leave empty to match any monitor...
- `monitor_description_pattern` (String) Regex (case-insensitive) matched against any of the incident's monitor descriptions. Leave empty to match any description...
- `owner_users` (Set) Users to add as owners on the incident when this rule matches...
- `owner_teams` (Set) Teams to add as owners on the incident when this rule matches...
- `inherit_owners_from_monitors` (Bool) When this rule matches, also assign every owner of the incident's monitors to the incident...
- `inherit_owners_from_hosts` (Bool) When this rule matches, also assign every owner of the incident's affected hosts to the incident...
- `inherit_owners_from_kubernetes_clusters` (Bool) When this rule matches, also assign every owner of the incident's affected Kubernetes clusters to the incident...
- `inherit_owners_from_docker_hosts` (Bool) When this rule matches, also assign every owner of the incident's affected Docker hosts to the incident...
- `inherit_owners_from_podman_hosts` (Bool) When this rule matches, also assign every owner of the incident's affected Podman hosts to the incident...
- `inherit_owners_from_services` (Bool) When this rule matches, also assign every owner of the incident's affected services to the incident...
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_incident_owner_rule.example <id>
```
