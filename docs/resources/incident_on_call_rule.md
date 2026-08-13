---
page_title: "oneuptime_incident_on_call_rule Resource - oneuptime"
subcategory: "Incidents"
description: |-
  Configure rules for automatically executing on-call duty policies when matching incidents are created
---

# oneuptime_incident_on_call_rule (Resource)

Configure rules for automatically executing on-call duty policies when matching incidents are created

## Example Usage

```terraform
resource "oneuptime_incident_on_call_rule" "example" {
  name = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Name of this incident on-call rule..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Description of this incident on-call rule..
- `is_enabled` (Bool) Whether this rule is enabled..
- `monitors` (Set) Only trigger for incidents from these monitors. Leave empty to match incidents from any monitor...
- `incident_severities` (Set) Only trigger for incidents with these severities. Leave empty to match incidents of any severity...
- `incident_labels` (Set) Only trigger for incidents that have at least one of these labels. Leave empty to match incidents regardless of incident labels...
- `monitor_labels` (Set) Only trigger for incidents from monitors that have at least one of these labels. Leave empty to match incidents regardless of monitor labels...
- `incident_title_pattern` (String) Regular expression pattern to match incident titles. Leave empty to match any title. Example: 'CPU.*high' matches titles containing 'CPU' followed by 'high'...
- `incident_description_pattern` (String) Regular expression pattern to match incident descriptions. Leave empty to match any description...
- `monitor_name_pattern` (String) Regular expression pattern to match monitor names. Leave empty to match any monitor name. Example: 'prod-.*' matches monitors starting with 'prod-'...
- `monitor_description_pattern` (String) Regular expression pattern to match monitor descriptions. Leave empty to match any monitor description...
- `on_call_duty_policies` (Set) On-call duty policies to execute when an incident matches this rule...
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
terraform import oneuptime_incident_on_call_rule.example <id>
```
