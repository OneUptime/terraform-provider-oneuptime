---
page_title: "oneuptime_security_event Data Source - oneuptime"
subcategory: "Other"
description: |-
  API endpoints for Security Event
---

# oneuptime_security_event (Data Source)

API endpoints for Security Event

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one security event may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_security_event" "example" {
  primary_entity_id = "example-primary-entity-id"
}

# Or by id:
data "oneuptime_security_event" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `activity_name` (String) Activity.
- `attributes` (String) Attributes.
- `category_name` (String) Category.
- `category_uid` (Number) Category UID.
- `class_name` (String) Event Class.
- `class_uid` (Number) Class UID.
- `event_uid` (String) Event UID.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `message` (String) Message.
- `primary_entity_id` (String) Source ID.
- `primary_entity_type` (String) Source Type.
- `principal_host` (String) Principal Host.
- `principal_ip` (String) Principal IP.
- `principal_process` (String) Principal Process.
- `principal_user` (String) Principal User.
- `product_name` (String) Product.
- `rule_id` (String) Rule ID.
- `rule_name` (String) Rule Name.
- `severity_id` (Number) Severity ID.
- `severity_name` (String) Severity.
- `status_name` (String) Status.
- `target_host` (String) Target Host.
- `target_ip` (String) Target IP.
- `target_port` (Number) Target Port.
- `target_resource` (String) Target Resource.
- `target_user` (String) Target User.
- `time` (String) Time.
- `vendor_name` (String) Vendor.

### Read-Only

- `attribute_keys` (Set of String) Attribute Keys.
- `mitre_tactics` (Set of String) MITRE Tactics.
- `mitre_techniques` (Set of String) MITRE Techniques.
- `observables` (Set of String) Observables.
- `project_id` (String) Project ID.
