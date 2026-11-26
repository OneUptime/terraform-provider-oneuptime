package shard03

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The resources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterResources("shard03", []OrderedResource{
        {Order: 34, New: NewSecurityCaseNoteTemplateResource},
        {Order: 55, New: NewDataIssueReportResource},
        {Order: 95, New: NewOnCallDutyExecutionLogResource},
        {Order: 187, New: NewIncidentSlaResource},
        {Order: 225, New: NewCephClusterLabelRuleResource},
        {Order: 283, New: NewScheduledMaintenanceUserOwnerResource},
        {Order: 289, New: NewScheduledMaintenanceNoteTemplateResource},
        {Order: 294, New: NewMonitorGroupResourceResource},
        {Order: 399, New: NewDockerHostUserOwnerResource},
        {Order: 419, New: NewStorageArrayResource},
        {Order: 451, New: NewQueueUserOwnerResource},
    }, nil)
}
