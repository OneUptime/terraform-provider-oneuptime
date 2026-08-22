---
page_title: "oneuptime_security_event Data Source - oneuptime"
subcategory: "Other"
description: |-
  API endpoints for Security Event
---

# oneuptime_security_event (Data Source)

API endpoints for Security Event Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_security_event" "by_name" {
  name = "example-security_event"
}

data "oneuptime_security_event" "by_id" {
  id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

- `id` (String) Look up by unique identifier. Exactly one of `id` or `name` must be set.. Computed.
- `name` (String) Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.. Computed.
- `project_id` (String) Project ID. Computed.
- `primary_entity_id` (String) Source ID. Computed.
- `primary_entity_type` (String) Source Type. Computed.
- `time` (String) Time. Computed.
- `event_uid` (String) Event UID. Computed.
- `category_uid` (Number) Category UID. Computed.
- `category_name` (String) Category. Computed.
- `class_uid` (Number) Class UID. Computed.
- `class_name` (String) Event Class. Computed.
- `activity_name` (String) Activity. Computed.
- `severity_id` (Number) Severity ID. Computed.
- `severity_name` (String) Severity. Computed.
- `status_name` (String) Status. Computed.
- `message` (String) Message. Computed.
- `vendor_name` (String) Vendor. Computed.
- `product_name` (String) Product. Computed.
- `rule_id` (String) Rule ID. Computed.
- `rule_name` (String) Rule Name. Computed.
- `mitre_tactics` (Set) MITRE Tactics. Computed.
- `mitre_techniques` (Set) MITRE Techniques. Computed.
- `principal_user` (String) Principal User. Computed.
- `principal_host` (String) Principal Host. Computed.
- `principal_ip` (String) Principal IP. Computed.
- `principal_process` (String) Principal Process. Computed.
- `target_user` (String) Target User. Computed.
- `target_host` (String) Target Host. Computed.
- `target_ip` (String) Target IP. Computed.
- `target_port` (Number) Target Port. Computed.
- `target_resource` (String) Target Resource. Computed.
- `observables` (Set) Observables. Computed.
- `attributes` (String) Attributes. Computed.
- `attribute_keys` (Set) Attribute Keys. Computed.
