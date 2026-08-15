---
page_title: "oneuptime_inventory_item Resource - oneuptime"
subcategory: "Other"
description: |-
  Catalog of everything OneUptime knows about your estate (service, host, k8s.pod, container, network device, ...), discovered from telemetry resource attributes, mirrored from inventory tables, or registered by hand.
---

# oneuptime_inventory_item (Resource)

Catalog of everything OneUptime knows about your estate (service, host, k8s.pod, container, network device, ...), discovered from telemetry resource attributes, mirrored from inventory tables, or registered by hand.

## Example Usage

```terraform
resource "oneuptime_inventory_item" "example" {
  entity_type = "Example short text"
  entity_key = "Example short text"
  source = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `entity_type` (String) The OpenTelemetry entity type (service, host, k8s.pod, container, ...)...
- `entity_key` (String) Stable identity hash derived from the entity's identifying attributes (matches the keys stamped into signal entityKeys columns)...
- `source` (String) How this row came to exist: discovered from telemetry, mirrored from a OneUptime inventory table, or created manually by a user. Determines whether stale-entity pruning applies...

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `display_name` (String) Human-readable name shown in the Inventory list...
- `description` (String) Free-text description. Primarily for manually created entities, where there are no telemetry attributes to explain what the thing is...
- `identifying_attributes` (String) The immutable identifying attribute set (the entity's identity). Descriptive attributes are deliberately excluded so they can change without changing the entity key...
- `descriptive_attributes` (String) Mutable descriptive metadata (image tag, version, IP, ...) merged last-writer-wins. Never part of the identity...
- `labels` (String) Labels observed on this entity's telemetry (e.g. promoted from oneuptime.label.* resource attributes), merged as a set union. Simple string array in v1 — a relation to the Label table is a follow-up...
- `resource_type` (String) Polymorphic pointer type to a rich typed row, if one exists (Service / Host / DockerHost / KubernetesCluster)...
- `resource_id` (String) A unique identifier for an object, represented as a UUID..
- `first_seen_at` (String) A date time object..
- `last_seen_at` (String) A date time object..
- `is_archived` (Bool) Is this item archived? Archived items are hidden from the default list but keep their identity and keep collecting telemetry...
- `custom_fields` (String) Custom fields on this item...
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `archived_at` (String) A date time object..
- `archived_by_user_id` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_inventory_item.example <id>
```
