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
  name = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Name of this queue label rule..

### Optional

- `criteria` (String) Versioned conditions that determine whether this rule matches a resource...
- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Description of this queue label rule..
- `is_enabled` (Bool) Whether this rule is enabled..
- `message_queue_labels` (Set) Only trigger for queues that already have at least one of these labels. Leave empty to match regardless of labels...
- `message_queue_name_pattern` (String) Regex (case-insensitive) matched against the queue name. Discovered queues are named after their destination (e.g. orders.created), so ^orders\. matches every queue whose name starts with orders. - use the messaging system pattern to match by broker. Leave empty to match any name...
- `message_queue_description_pattern` (String) Regex (case-insensitive) matched against the queue description. Leave empty to match any description...
- `message_queue_system_pattern` (String) Regex (case-insensitive) matched against the queue's messaging system - both its OpenTelemetry messaging.system value (kafka, rabbitmq, aws_sqs, servicebus, ...) and its display name (Apache Kafka, RabbitMQ, Amazon SQS, Azure Service Bus, ...). ^kafka$ matches every Kafka topic. Leave empty to match any system...
- `labels_to_add` (Set) Labels to attach to the queue when this rule matches. Already-attached labels are not duplicated...

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_queue_label_rule.example <id>
```
