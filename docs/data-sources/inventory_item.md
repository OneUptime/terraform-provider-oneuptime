---
page_title: "oneuptime_inventory_item Data Source - oneuptime"
subcategory: "Other"
description: |-
  Catalog of everything OneUptime knows about your estate (service, host, k8s.pod, container, network device, ...), discovered from telemetry resource attributes, mirrored from inventory tables, or registered by hand.
---

# oneuptime_inventory_item (Data Source)

Catalog of everything OneUptime knows about your estate (service, host, k8s.pod, container, network device, ...), discovered from telemetry resource attributes, mirrored from inventory tables, or registered by hand.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one inventory item may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_inventory_item" "example" {
  entity_type = "example-entity-type"
}

# Or by id:
data "oneuptime_inventory_item" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `archived_by_user_id` (String) User ID who archived this object (if this object was archived by a User). The ID of a `oneuptime_user` (see the data source).
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Free-text description. Primarily for manually created entities, where there are no telemetry attributes to explain what the thing is.
- `display_name` (String) Human-readable name shown in the Inventory list.
- `entity_key` (String) Stable identity hash derived from the entity's identifying attributes (matches the keys stamped into signal entityKeys columns).
- `entity_type` (String) The OpenTelemetry entity type (service, host, k8s.pod, container, ...).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `inventory_status` (String) Current heartbeat status: live, recent, stale, never seen, or not tracked.
- `is_archived` (Boolean) Is this item archived? Archived items are hidden from the default list but keep their identity and keep collecting telemetry.
- `resource_id` (String) Polymorphic pointer id to the rich typed row named by resourceType, if any.
- `resource_type` (String) Polymorphic pointer type to a rich typed row, if one exists (Service / Host / DockerHost / KubernetesCluster).
- `source` (String) How this row came to exist: discovered from telemetry, mirrored from a OneUptime inventory table, or created manually by a user. Determines whether stale-entity pruning applies.

### Read-Only

- `archived_at` (String) When this item was archived.
- `created_at` (String) Date and Time when the object was created.
- `custom_fields` (String) Custom fields on this item. A JSON value: write it with `jsonencode()`.
- `descriptive_attributes` (String) Mutable descriptive metadata (image tag, version, IP, ...) merged last-writer-wins. Never part of the identity. A JSON value: write it with `jsonencode()`.
- `first_seen_at` (String) When this entity was first observed in telemetry.
- `identifying_attributes` (String) The immutable identifying attribute set (the entity's identity). Descriptive attributes are deliberately excluded so they can change without changing the entity key. A JSON value: write it with `jsonencode()`.
- `labels` (String) Labels observed on this entity's telemetry (e.g. promoted from oneuptime.label.* resource attributes), merged as a set union. Simple string array in v1 — a relation to the Label table is a follow-up. A JSON value: write it with `jsonencode()`.
- `last_seen_at` (String) Most recent time this entity was observed in telemetry (bumped, throttled). Drives staleness pruning.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
