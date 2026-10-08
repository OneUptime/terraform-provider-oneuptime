---
page_title: "oneuptime_host Data Source - oneuptime"
subcategory: "Other"
description: |-
  Hosts that are being monitored in this project. Each host is auto-discovered when an OTel Collector reports the host.name resource attribute, or can be manually registered.
---

# oneuptime_host (Data Source)

Hosts that are being monitored in this project. Each host is auto-discovered when an OTel Collector reports the host.name resource attribute, or can be manually registered.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one host may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_host" "example" {
  name = "Example host"
}

# Or by id:
data "oneuptime_host" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `agent_version` (String) Version of the OneUptime agent reporting telemetry on this host, as self-reported via the oneuptime.agent.version resource attribute.
- `ai_access_last_error` (String) The most recent failure OneUptime AI hit while running a command on this host, kept until the next successful command. Set by the server.
- `ai_remediation_mode` (String) Disabled: AI never proposes or runs a change on this host. RequireApproval: AI composes a command plan and a human approves it with one click before anything runs. Automatic: safe changes (SafeWrite) run without a human; a riskier change is proposed for approval unless the host's allowlist names its exact shape. BypassApproval: every change the policy allows — safe AND riskier — runs on its own, except what always needs a human. In EVERY mode: Denied commands never run, commands the policy marks requiresHuman always ask, and the agent itself refuses every write unless it was started with ONEUPTIME_AI_ALLOW_WRITES=true (and then only on the targets ONEUPTIME_AI_WRITE_TARGETS allows, never its protected targets). Anyone who may edit the host can lower the mode; raising it needs Project Owner, Project Admin or Edit Auto Remediation Rule.
- `archived_by_user_id` (String) User ID who archived this object (if this object was archived by a User). The ID of a `oneuptime_user` (see the data source).
- `cloud_account_id` (String) Last-seen value of the cloud.account.id OpenTelemetry resource attribute.
- `cloud_platform` (String) Last-seen value of the cloud.platform OpenTelemetry resource attribute, e.g. aws_ec2, gcp_compute_engine.
- `cloud_provider` (String) Last-seen value of the cloud.provider OpenTelemetry resource attribute, e.g. aws, gcp, azure.
- `cloud_region` (String) Last-seen value of the cloud.region OpenTelemetry resource attribute, e.g. us-east-1.
- `container_runtime` (String) Container runtime detected on this host, if any (e.g. docker, containerd).
- `cpu_cores` (Number) Logical CPU core count, sourced from system.cpu.logical.count metric.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `deployment_environment` (String) Last-seen value of the deployment.environment.name (or deployment.environment) OpenTelemetry resource attribute, e.g. production, staging.
- `description` (String) Friendly description for this host.
- `docker_host_id` (String) Optional FK to the DockerHost record for this same host when it is also running the Docker runtime. The ID of a `oneuptime_docker_host`.
- `host_arch` (String) CPU architecture from the OTel host.arch resource attribute.
- `host_id` (String) Stable host identifier reported by the OTel host.id resource attribute.
- `host_identifier` (String) Unique identifier for this host, sourced from the host.name OTel resource attribute.
- `host_ip_addresses` (String) Comma-separated list of every IP address reported by the OTel host.ip resource attribute, in the order the collector reported them, deduplicated. The Hosts list shows the most routable one (IPv4, non-loopback, non-link-local) first; the host detail page groups them all by category.
- `host_type` (String) Cloud-instance class reported by the OTel host.type resource attribute.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_ai_investigation_enabled` (Boolean) When on, OneUptime AI runs read-only commands (systemctl status, journalctl, df, free, uptime, ps, ss) on this host, through its Host AI agent, while investigating incidents and alerts linked to it, and uses their output, with secret values redacted, as evidence. Nothing is ever changed by an investigation. On by default. Anyone who may edit the host can turn it on or off.
- `is_archived` (Boolean) Is this host archived? Archived hosts are hidden from lists but keep collecting telemetry.
- `kubernetes_cluster_id` (String) Optional FK to the KubernetesCluster this host belongs to when k8s.cluster.name is reported. The ID of a `oneuptime_kubernetes_cluster`.
- `name` (String) Friendly name for this host.
- `os_type` (String) Operating system type of the host.
- `os_version` (String) Operating system version of the host.
- `otel_collector_status` (String) Connection status of the OTel Collector reporting on this host (connected or disconnected).
- `process_count` (Number) Most recent process count from system.processes.count metric.
- `proxmox_cluster_id` (String) Optional FK to the ProxmoxCluster this host runs inside (as a guest VM). The ID of a `oneuptime_proxmox_cluster`.
- `retain_telemetry_data_for_days` (Number) Number of days to retain telemetry data for this host. Leave blank to use the project-wide default.
- `runtime_name` (String) Last-seen value of the process.runtime.name OpenTelemetry resource attribute.
- `runtime_version` (String) Last-seen value of the process.runtime.version OpenTelemetry resource attribute.
- `slug` (String) Friendly globally unique name for your object.
- `total_memory_bytes` (Number) Total physical memory in bytes, sourced from system.memory.usage metric (sum of all states).

### Read-Only

- `ai_access_configured_at` (String) When OneUptime AI access to this host was first configured by anyone saving an AI access setting. Set by the server; never cleared, so a Host AI agent that registers later never overwrites a setting an operator chose.
- `ai_access_last_verified_at` (String) When a command from OneUptime AI last succeeded on this host through its Host AI agent. Set by the server.
- `ai_command_allowlist` (String) Optional JSON array of command patterns that Automatic mode may run on this host without approval even though they are riskier changes. Each pattern is one command line for this host's agent (systemctl, journalctl and the other host programs it may run) and is compared with the command word by word: * stands for exactly one word (a name, an id), never for extra words or flags, and every flag the command uses must be written out in the pattern. At most 50 patterns of at most 500 characters each; a pattern that is not one valid write command for this host is refused. Destructive commands (Denied tier) never run regardless, and a command that always needs a human still asks. Adding a pattern needs Project Owner, Project Admin or Edit Auto Remediation Rule; anyone who may edit the host can remove patterns or clear the list. A JSON value: write it with `jsonencode()`.
- `archived_at` (String) When was this host archived?
- `created_at` (String) Date and Time when the object was created.
- `labels` (Set of String) Relation to Labels Array where this object is categorized in. IDs of `oneuptime_label` resources.
- `last_seen_at` (String) When telemetry was last received from this host.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `telemetry_retention_config` (String) Per-pillar retention overrides for this host (logs by severity, traces by status, metrics, profiles). Unset fields fall back to the host default, then the project's retention settings. A JSON value: write it with `jsonencode()`.
- `updated_at` (String) Date and Time when the object was updated.
