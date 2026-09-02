---
page_title: "oneuptime_oid_collection_template Resource - oneuptime"
subcategory: "Other"
description: |-
  A reusable set of SNMP health OIDs. Every network device linked to a template collects its OIDs, and editing the template changes what every linked device collects on its next poll.
---

# oneuptime_oid_collection_template (Resource)

A reusable set of SNMP health OIDs. Every network device linked to a template collects its OIDs, and editing the template changes what every linked device collects on its next poll.

## Example Usage

```terraform
resource "oneuptime_oid_collection_template" "example" {
  name = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) The device type this template describes. Devices linked to it all collect the same OIDs...

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Friendly description that will help you remember..
- `oids` (String) SNMP OIDs (CPU, memory, temperature, or any custom OID) collected by every device linked to this template. You do not need OIDs for interfaces - bits in/out, errors, utilization and up/down are walked for every port automatically...
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `slug` (String) Friendly globally unique name for your object..
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_oid_collection_template.example <id>
```
