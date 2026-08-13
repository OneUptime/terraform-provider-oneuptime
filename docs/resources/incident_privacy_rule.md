---
page_title: "oneuptime_incident_privacy_rule Resource - oneuptime"
subcategory: "Incidents"
description: |-
  Configure rules for automatically marking matching incidents as private
---

# oneuptime_incident_privacy_rule (Resource)

Configure rules for automatically marking matching incidents as private

## Example Usage

```terraform
resource "oneuptime_incident_privacy_rule" "example" {
  name = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Name of this incident privacy rule..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Description of this incident privacy rule..
- `is_enabled` (Bool) Whether this rule is enabled..
- `monitors` (Set) Only trigger for incidents from these monitors. Leave empty to match incidents from any monitor...
- `incident_severities` (Set) Only trigger for incidents with these severities. Leave empty to match incidents of any severity...
- `incident_labels` (Set) Only trigger for incidents that have at least one of these labels. Leave empty to match regardless of incident labels...
- `monitor_labels` (Set) Only trigger for incidents from monitors that have at least one of these labels. Leave empty to match regardless of monitor labels...
- `incident_title_pattern` (String) Regex (case-insensitive) matched against the incident title. Leave empty to match any title...
- `incident_description_pattern` (String) Regex (case-insensitive) matched against the incident description. Leave empty to match any description...
- `monitor_name_pattern` (String) Regex (case-insensitive) matched against any of the incident's monitor names. Leave empty to match any monitor...
- `monitor_description_pattern` (String) Regex (case-insensitive) matched against any of the incident's monitor descriptions. Leave empty to match any description...
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
terraform import oneuptime_incident_privacy_rule.example <id>
```
