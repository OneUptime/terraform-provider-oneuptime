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
  from_entity_key   = "Example short text"
  to_entity_key     = "Example short text"
  relationship_type = "Example short text"
  source            = "Example short text"
}
```

## Schema

### Required

- `from_entity_key` (String) Stable identity key of the source entity of this edge.
- `relationship_type` (String) The inferred relationship (runs-on, member-of, hosted-on, part-of, instance-of).
- `source` (String) Whether this edge was derived from telemetry or drawn manually by a user. Determines whether stale-edge pruning applies.
- `to_entity_key` (String) Stable identity key of the target entity of this edge.

### Optional

- `avg_duration_ms` (Number) Average call duration in milliseconds over this edge in the most recent computation window (depends-on edges only).
- `call_count` (Number) Calls observed over this edge in the most recent computation window (depends-on edges only).
- `error_count` (Number) Errored calls observed over this edge in the most recent computation window (depends-on edges only).
- `first_seen_at` (String) When this relationship was first observed in telemetry.
- `last_seen_at` (String) Most recent time this relationship was observed (bumped, throttled). Drives staleness pruning.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing inventory item relationship by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_inventory_item_relationship.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_inventory_item_relationship.example <id>
```
