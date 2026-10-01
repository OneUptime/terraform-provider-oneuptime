---
page_title: "oneuptime_queue Data Source - oneuptime"
subcategory: "Other"
description: |-
  Message queues, topics and subscriptions this project's applications publish to and consume from. Each queue is discovered from the OpenTelemetry messaging spans of instrumented applications and from broker metrics sent by an OpenTelemetry Collector, or added manually.
---

# oneuptime_queue (Data Source)

Message queues, topics and subscriptions this project's applications publish to and consume from. Each queue is discovered from the OpenTelemetry messaging spans of instrumented applications and from broker metrics sent by an OpenTelemetry Collector, or added manually. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_queue" "by_name" {
  name = "example-queue"
}

data "oneuptime_queue" "by_id" {
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
- `project_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `slug` (String) Friendly globally unique name for your object.. Computed.
- `description` (String) Friendly description for this queue.. Computed.
- `queue_identifier` (String) Stable identity of this queue within the project: the messaging system's identity family, the broker scope (the Azure Service Bus / Event Hubs namespace, empty for every other system) and the destination, canonicalized and joined with '|' (e.g. kafka||orders.created, servicebus|orders-prod|orders). An ActiveMQ queue keys as jms, like the JMS clients that use it, so both are one queue. Computed by the server for manually added queues... Computed.
- `messaging_system` (String) The broker this queue lives on, as an OpenTelemetry messaging.system value, e.g. kafka, rabbitmq, activemq, jms, aws_sqs, aws.sns, gcp_pubsub, servicebus, eventhubs, pulsar, rocketmq, nats or bullmq. A queue first seen through JMS is refined to its broker (activemq) once that broker's metrics are seen... Computed.
- `destination_name` (String) The queue, topic or subscription name as applications and the broker name it (the OpenTelemetry messaging.destination.name value), normalized the way discovery normalizes it: an SQS queue URL or an SNS topic ARN becomes its name, a Pub/Sub resource path its id, a Pulsar short name its persistent://public/default/ topic, and each UUID {uuid}... Computed.
- `broker_scope` (String) The Azure Service Bus or Event Hubs namespace this queue lives in, lowercased (the first part of <namespace>.servicebus.windows.net). Part of the queue's identity, because two namespaces can each hold a queue of the same name. Empty for every other messaging system... Computed.
- `broker_address` (String) The broker address (host[:port]) last reported for this queue by its telemetry, with any credentials removed. Display only - it is not part of the queue's identity... Computed.
- `discovery_source` (String) How this queue was first discovered: traces (messaging spans of instrumented applications), broker-metrics (the broker's own metrics) or manual... Computed.
- `last_seen_at` (String) A date time object.. Computed.
- `broker_metrics_last_seen_at` (String) A date time object.. Computed.
- `auto_archived_at` (String) A date time object.. Computed.
- `manually_restored_at` (String) A date time object.. Computed.
- `automatic_assignments` (String) Label and owner ids that label rules or owner rules attached automatically. Maintained by OneUptime... Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `is_archived` (Bool) Is this queue archived? Archived queues are hidden from lists but their telemetry is still collected... Computed.
- `archived_at` (String) A date time object.. Computed.
- `archived_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `labels` (Set) Relation to Labels Array where this object is categorized in... Computed.
