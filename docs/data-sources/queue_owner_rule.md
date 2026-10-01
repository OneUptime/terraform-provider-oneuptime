---
page_title: "oneuptime_queue_owner_rule Data Source - oneuptime"
subcategory: "Other"
description: |-
  Configure rules for automatically assigning owner users and teams when matching queues are created
---

# oneuptime_queue_owner_rule (Data Source)

Configure rules for automatically assigning owner users and teams when matching queues are created Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_queue_owner_rule" "by_name" {
  name = "example-queue_owner_rule"
}

data "oneuptime_queue_owner_rule" "by_id" {
  id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

- `id` (String) Look up by unique identifier. Exactly one of `id` or `name` must be set.. Computed.
- `name` (String) Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.. Computed.
- `created_at` (String) A date time object.. Computed.
- `updated_at` (String) A date time object.. Computed.
- `deleted_at` (String) A date time object.. Computed.
- `version` (Number) Object version. Computed.
- `criteria` (String) Versioned conditions that determine whether this rule matches a resource... Computed.
- `project_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `description` (String) Description of this queue owner rule.. Computed.
- `is_enabled` (Bool) Whether this rule is enabled.. Computed.
- `notify_owners` (Bool) Send notifications to owner users and teams when they are added by this rule.. Computed.
- `message_queue_labels` (Set) Only trigger for queues that have at least one of these labels. Leave empty to match regardless of labels... Computed.
- `message_queue_name_pattern` (String) Regex (case-insensitive) matched against the queue name. Discovered queues are named after their destination (e.g. orders.created), so ^orders\. matches every queue whose name starts with orders. - use the messaging system pattern to match by broker. Leave empty to match any name... Computed.
- `message_queue_description_pattern` (String) Regex (case-insensitive) matched against the queue description. Leave empty to match any description... Computed.
- `message_queue_system_pattern` (String) Regex (case-insensitive) matched against the queue's messaging system - both its OpenTelemetry messaging.system value (kafka, rabbitmq, aws_sqs, servicebus, ...) and its display name (Apache Kafka, RabbitMQ, Amazon SQS, Azure Service Bus, ...). ^kafka$ matches every Kafka topic. Leave empty to match any system... Computed.
- `owner_users` (Set) Users to add as owners on the queue when this rule matches... Computed.
- `owner_teams` (Set) Teams to add as owners on the queue when this rule matches... Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
