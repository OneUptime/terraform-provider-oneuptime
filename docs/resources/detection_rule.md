---
page_title: "oneuptime_detection_rule Resource - oneuptime"
subcategory: "Other"
description: |-
  Sigma detection rules evaluated against security events. Matches create alerts and detection findings.
---

# oneuptime_detection_rule (Resource)

Sigma detection rules evaluated against security events. Matches create alerts and detection findings.

## Example Usage

```terraform
resource "oneuptime_detection_rule" "example" {
  name = jsonencode({
    "_type": "Name",
    "value": "John Doe"
  })
  sigma_rule_yaml = "This is an example of very long text content that might be stored in this field. It can contain a lot of information, such as detailed descriptions, comments, or any other lengthy text data that needs to be stored in the database."
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Name object.
- `sigma_rule_yaml` (String) The Sigma rule to evaluate, in YAML. detection selections and condition are compiled to a ClickHouse query over security events...

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Description of what this detection rule looks for...
- `is_enabled` (Bool) Whether this detection rule is evaluated...
- `evaluation_interval_in_minutes` (Number) How often this rule is evaluated, in minutes. The evaluation window covers the time since the previous evaluation...
- `group_by_field` (String) Optional security-event field (e.g. principalHost, principalUser) to group matches by. One alert is opened per distinct value; empty groups all matches into one alert...
- `should_create_alert` (Bool) Whether matches open OneUptime alerts...
- `should_write_detection_finding` (Bool) Whether matches also write a Detection Finding security event back into the events table...
- `alert_severity_id` (String) A unique identifier for an object, represented as a UUID..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `last_evaluated_at` (String) A date time object..
- `last_match_at` (String) A date time object..
- `last_error` (String) The most recent evaluation error, if any. Cleared on the next successful evaluation...
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_detection_rule.example <id>
```
