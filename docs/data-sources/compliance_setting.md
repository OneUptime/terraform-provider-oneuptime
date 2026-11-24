---
page_title: "oneuptime_compliance_setting Data Source - oneuptime"
subcategory: "Other"
description: |-
  A compliance framework, or one of its requirements, a resource turned off for its own data stores.
---

# oneuptime_compliance_setting (Data Source)

A compliance framework, or one of its requirements, a resource turned off for its own data stores.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one compliance setting may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_compliance_setting" "example" {
  resource_type = "KubernetesCluster"
}

# Or by id:
data "oneuptime_compliance_setting" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) ID of the user who created this compliance setting. The ID of a `oneuptime_user` (see the data source).
- `framework` (String) The compliance framework this setting changes for the resource: GDPR, CCPA, HIPAA, PCI DSS, SOC 2, ISO 27001, NIST CSF, LGPD or DPDP, or for a Kubernetes cluster PSS Baseline, PSS Restricted or CIS Kubernetes.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_enabled` (Boolean) Whether the framework or requirement applies to the resource. Off takes the resource's data stores out of its scope, and for an accepted Kubernetes finding, accepts it.
- `namespace` (String) For a Kubernetes cluster's PSS Baseline, PSS Restricted or CIS Kubernetes setting, the one namespace it is for. Empty for the whole resource.
- `note` (String) Why it does not apply to the resource, or why the Kubernetes finding is accepted.
- `object_kind` (String) For an accepted Kubernetes finding, the kind of the object it is accepted for: Deployment, DaemonSet, CronJob, ClusterRoleBinding and the others. Empty for a framework, a control or a namespace.
- `object_name` (String) For an accepted Kubernetes finding, the name of the object it is accepted for, as Kubernetes names it. Empty for a framework, a control or a namespace.
- `requirement` (String) One requirement of the framework, as Data Security names it (for example Art. 32 or CC6.1). Empty for the whole framework.
- `resource_id` (String) The ID of the resource that turned the framework off.
- `resource_type` (String) The kind of resource that turned the framework off: KubernetesCluster, DockerHost, PodmanHost, DatabaseServer, CloudResource or Service.
- `turned_off_by_user_id` (String) ID of the user who last turned the framework or requirement off for the resource, or accepted the finding. Empty while it applies, or when an API key turned it off. The ID of a `oneuptime_user` (see the data source).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `expires_at` (String) For an accepted Kubernetes finding, when it fails again. Empty: it stays accepted until revoked.
- `project_id` (String) ID of the project this compliance setting belongs to. The ID of a `oneuptime_project`.
- `turned_off_at` (String) When the framework or requirement was last turned off for the resource, or the finding accepted. Empty while it applies.
- `updated_at` (String) Date and Time when the object was updated.
