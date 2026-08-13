---
page_title: "oneuptime_kubernetes_cluster_label_rule Resource - oneuptime"
subcategory: "Other"
description: |-
  Configure rules for automatically attaching labels to Kubernetes clusters when matching Kubernetes clusters are created
---

# oneuptime_kubernetes_cluster_label_rule (Resource)

Configure rules for automatically attaching labels to Kubernetes clusters when matching Kubernetes clusters are created

## Example Usage

```terraform
resource "oneuptime_kubernetes_cluster_label_rule" "example" {
  name = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Name of this Kubernetes cluster label rule..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Description of this Kubernetes cluster label rule..
- `is_enabled` (Bool) Whether this rule is enabled..
- `kubernetes_cluster_labels` (Set) Only trigger for Kubernetes clusters that already have at least one of these labels. Leave empty to match regardless of labels...
- `kubernetes_cluster_name_pattern` (String) Regex (case-insensitive) matched against the Kubernetes cluster name. Leave empty to match any name...
- `kubernetes_cluster_description_pattern` (String) Regex (case-insensitive) matched against the Kubernetes cluster description. Leave empty to match any description...
- `labels_to_add` (Set) Labels to attach to the Kubernetes cluster when this rule matches. Already-attached labels are not duplicated...
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_kubernetes_cluster_label_rule.example <id>
```
