---
page_title: "oneuptime_inventory_item Resource - oneuptime"
subcategory: "Infrastructure"
description: |-
  Catalog of everything OneUptime knows about your estate (service, host, k8s.pod, container, network device, ...), discovered from telemetry resource attributes, mirrored from inventory tables, or registered by hand.
---

# oneuptime_inventory_item (Resource)

Catalog of everything OneUptime knows about your estate (service, host, k8s.pod, container, network device, ...), discovered from telemetry resource attributes, mirrored from inventory tables, or registered by hand.

## Example Usage

```terraform
resource "oneuptime_inventory_item" "example" {
  entity_type = "service"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `entity_type` (String) The OpenTelemetry entity type (service, host, k8s.pod, container, ...). Allowed values: `service`, `service.instance`, `host`, `container`, `process`, `k8s.cluster`, `k8s.namespace`, `k8s.node`, `k8s.pod`, `k8s.deployment`, `proxmox.cluster`, `openstack.cloud`, `proxmox.node`, `openstack.hypervisor`, `proxmox.guest`, `openstack.instance`, `ceph.cluster`, `storage.array`, `vmware.vcenter`, `vmware.cluster`, `vmware.host`, `vmware.vm`, `vmware.datastore`, `docker.swarm.cluster`, `docker.swarm.node`, `docker.swarm.service`, `docker.swarm.task`, `telemetry.sdk`, `database`, `remote.service`, `database.server`, `message.queue`, `openstack.resource`, `network.device`, `cloud.resource`, `serverless.function`, `rum.application`, `iot.device`, `docker.host`, `podman.host`, `external.service`, `external.database`, `appliance`.

### Optional

- `custom_fields` (String) Custom fields on this item. A JSON value: write it with `jsonencode()`.
- `description` (String) Free-text description. Primarily for manually created entities, where there are no telemetry attributes to explain what the thing is.
- `display_name` (String) Human-readable name shown in the Inventory list.
- `is_archived` (Boolean) Is this item archived? Archived items are hidden from the default list but keep their identity and keep collecting telemetry. OneUptime changes it too, on its own: left out, it keeps the value OneUptime gives it. Defaults to `false`.

### Read-Only

- `agent_permissions_status` (String) Whether the OneUptime Android agent still lacks a grant on the device - READ_LOGS, READ_DROPBOX_DATA or usage access - or has them all.
- `archived_at` (String) When this item was archived.
- `archived_by_user_id` (String) User ID who archived this object (if this object was archived by a User). The ID of a `oneuptime_user` (see the data source).
- `created_at` (String) Date and Time when the object was created.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `descriptive_attributes` (String) Mutable descriptive metadata (image tag, version, IP, ...) merged last-writer-wins. Never part of the identity. A JSON value: read it with `jsondecode()`.
- `entity_key` (String) Stable identity hash derived from the entity's identifying attributes (matches the keys stamped into signal entityKeys columns).
- `first_seen_at` (String) When this entity was first observed in telemetry.
- `id` (String) Unique identifier for the resource.
- `identifying_attributes` (String) The immutable identifying attribute set (the entity's identity). Descriptive attributes are deliberately excluded so they can change without changing the entity key. A JSON value: read it with `jsondecode()`.
- `inventory_status` (String) Current heartbeat status: live, recent, stale, never seen, or not tracked.
- `labels` (String) Labels observed on this entity's telemetry (e.g. promoted from oneuptime.label.* resource attributes), merged as a set union. Simple string array in v1 — a relation to the Label table is a follow-up. A JSON value: read it with `jsondecode()`.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `last_seen_at` (String) Most recent time this entity was observed in telemetry (bumped, throttled). Drives staleness pruning.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `resource_id` (String) Polymorphic pointer id to the rich typed row named by resourceType, if any.
- `resource_type` (String) Polymorphic pointer type to a rich typed row, if one exists (Service / Host / DockerHost / KubernetesCluster).
- `security_patch_status` (String) An Android device's security patch: up to date when its patch level is 90 days old or newer, out of date when older.
- `source` (String) How this row came to exist: discovered from telemetry, mirrored from a OneUptime inventory table, or created manually by a user. Determines whether stale-entity pruning applies.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing inventory item by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_inventory_item.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_inventory_item.example <id>
```
