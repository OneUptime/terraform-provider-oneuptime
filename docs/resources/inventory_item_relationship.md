---
page_title: "oneuptime_inventory_item_relationship Resource - oneuptime"
subcategory: "Other"
description: |-
  Directed relationships between telemetry entities (runs-on, member-of, hosted-on, part-of, instance-of), inferred from resource co-occurrence.
---

# oneuptime_inventory_item_relationship (Resource)

Directed relationships between telemetry entities (runs-on, member-of, hosted-on, part-of, instance-of), inferred from resource co-occurrence.

## Example Usage

```terraform
resource "oneuptime_inventory_item_relationship" "example" {
  from_entity_key = "Example short text"
  to_entity_key = "Example short text"
  relationship_type = "Example short text"
  source = "Example short text"
}
```

## Schema

### Required

- `from_entity_key` (String) Stable identity key of the source entity of this edge...
- `to_entity_key` (String) Stable identity key of the target entity of this edge...
- `relationship_type` (String) The inferred relationship (runs-on, member-of, hosted-on, part-of, instance-of)...
- `source` (String) Whether this edge was derived from telemetry or drawn manually by a user. Determines whether stale-edge pruning applies...

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `first_seen_at` (String) A date time object..
- `last_seen_at` (String) A date time object..
- `call_count` (Number) Calls observed over this edge in the most recent computation window (depends-on edges only)...
- `error_count` (Number) Errored calls observed over this edge in the most recent computation window (depends-on edges only)...
- `avg_duration_ms` (Number) Average call duration in milliseconds over this edge in the most recent computation window (depends-on edges only)...
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_inventory_item_relationship.example <id>
```
