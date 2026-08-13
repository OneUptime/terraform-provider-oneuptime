---
page_title: "oneuptime_incident_episode_privacy_rule Resource - oneuptime"
subcategory: "Incidents"
description: |-
  Configure rules for automatically marking matching incident episodes as private
---

# oneuptime_incident_episode_privacy_rule (Resource)

Configure rules for automatically marking matching incident episodes as private

## Example Usage

```terraform
resource "oneuptime_incident_episode_privacy_rule" "example" {
  name = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Name of this incident episode privacy rule..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Description of this incident episode privacy rule..
- `is_enabled` (Bool) Whether this rule is enabled..
- `incident_severities` (Set) Only trigger for episodes with these severities. Leave empty to match episodes of any severity...
- `episode_labels` (Set) Only trigger for episodes that have at least one of these labels. Leave empty to match regardless of episode labels...
- `episode_title_pattern` (String) Regex (case-insensitive) matched against the episode title. Leave empty to match any title...
- `episode_description_pattern` (String) Regex (case-insensitive) matched against the episode description. Leave empty to match any description...
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
terraform import oneuptime_incident_episode_privacy_rule.example <id>
```
