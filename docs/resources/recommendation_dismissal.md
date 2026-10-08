---
page_title: "oneuptime_recommendation_dismissal Resource - oneuptime"
subcategory: "Other"
description: |-
  Recommendations your team has dismissed. Dismissing hides a recommendation for everyone on the project until it is restored; it never deletes anything that was already created from it.
---

# oneuptime_recommendation_dismissal (Resource)

Recommendations your team has dismissed. Dismissing hides a recommendation for everyone on the project until it is restored; it never deletes anything that was already created from it.

## Example Usage

```terraform
resource "oneuptime_recommendation_dismissal" "example" {
  recommendation_type = "Example short text"
  recommendation_id   = "Example short text"
}
```

## Schema

### Required

- `recommendation_id` (String) The catalog-wide id of the dismissed recommendation, for example Kubernetes:k8s-hpa-at-max-replicas.
- `recommendation_type` (String) Which family of recommendation this dismissal belongs to. See the RecommendationType enum.

### Optional

- `dismissal_reason` (String) Optional note explaining why this recommendation was dismissed, shown to whoever finds it in the dismissed list later.
- `resource_id` (String) ID of the resource this recommendation was shown on. Polymorphic — it points at whichever table Resource Type names — so it carries no foreign key.
- `resource_type` (String) The kind of resource this recommendation was shown on, for example Kubernetes or Docker. Empty for recommendations that are not scoped to a resource.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing recommendation dismissal by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_recommendation_dismissal.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_recommendation_dismissal.example <id>
```
