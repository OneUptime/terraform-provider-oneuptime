---
page_title: "oneuptime_queue_owner_rule Resource - oneuptime"
subcategory: "Other"
description: |-
  Configure rules for automatically assigning owner users and teams when matching queues are created
---

# oneuptime_queue_owner_rule (Resource)

Configure rules for automatically assigning owner users and teams when matching queues are created

## Example Usage

```terraform
resource "oneuptime_queue_owner_rule" "example" {
  name = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Name of this queue owner rule..

### Optional

- `criteria` (String) Versioned conditions that determine whether this rule matches a resource...
- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Description of this queue owner rule..
- `is_enabled` (Bool) Whether this rule is enabled..
- `notify_owners` (Bool) Send notifications to owner users and teams when they are added by this rule..
- `message_queue_labels` (Set) Only trigger for queues that have at least one of these labels. Leave empty to match regardless of labels...
- `message_queue_name_pattern` (String) Regex (case-insensitive) matched against the queue name. Discovered queues are named after their destination (e.g. orders.created), so ^orders\. matches every queue whose name starts with orders. - use the messaging system pattern to match by broker. Leave empty to match any name...
- `message_queue_description_pattern` (String) Regex (case-insensitive) matched against the queue description. Leave empty to match any description...
- `message_queue_system_pattern` (String) Regex (case-insensitive) matched against the queue's messaging system - both its OpenTelemetry messaging.system value (kafka, rabbitmq, aws_sqs, servicebus, ...) and its display name (Apache Kafka, RabbitMQ, Amazon SQS, Azure Service Bus, ...). ^kafka$ matches every Kafka topic. Leave empty to match any system...
- `owner_users` (Set) Users to add as owners on the queue when this rule matches...
- `owner_teams` (Set) Teams to add as owners on the queue when this rule matches...

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
terraform import oneuptime_queue_owner_rule.example <id>
```
