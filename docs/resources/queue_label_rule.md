---
page_title: "oneuptime_queue_label_rule Resource - oneuptime"
subcategory: "Other"
description: |-
  Configure rules for automatically attaching labels to queues when matching queues are created
---

# oneuptime_queue_label_rule (Resource)

Configure rules for automatically attaching labels to queues when matching queues are created

## Example Usage

```terraform
resource "oneuptime_queue_label_rule" "example" {
  name        = "Example queue label rule"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Name of this queue label rule.

### Optional

- `criteria` (String) Versioned conditions that determine whether this rule matches a resource. A JSON value: write it with `jsonencode()`.
- `description` (String) Description of this queue label rule.
- `is_enabled` (Boolean) Whether this rule is enabled. Defaults to `true`.
- `labels_to_add` (Set of String) Labels to attach to the queue when this rule matches. Already-attached labels are not duplicated. IDs of `oneuptime_label` resources.
- `message_queue_description_pattern` (String) Regex (case-insensitive) matched against the queue description. Leave empty to match any description.
- `message_queue_labels` (Set of String) Only trigger for queues that already have at least one of these labels. Leave empty to match regardless of labels. IDs of `oneuptime_label` resources.
- `message_queue_name_pattern` (String) Regex (case-insensitive) matched against the queue name. Discovered queues are named after their destination (e.g. orders.created), so ^orders\. matches every queue whose name starts with orders. - use the messaging system pattern to match by broker. Leave empty to match any name.
- `message_queue_system_pattern` (String) Regex (case-insensitive) matched against the queue's messaging system - both its OpenTelemetry messaging.system value (kafka, rabbitmq, aws_sqs, servicebus, ...) and its display name (Apache Kafka, RabbitMQ, Amazon SQS, Azure Service Bus, ...). ^kafka$ matches every Kafka topic. Leave empty to match any system.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing queue label rule by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_queue_label_rule.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_queue_label_rule.example <id>
```
