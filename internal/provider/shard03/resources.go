package shard03

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The resources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterResources("shard03", []OrderedResource{
        {Order: 33, New: NewSecurityCaseNoteTemplateResource},
        {Order: 54, New: NewDataIssueReportResource},
        {Order: 91, New: NewOnCallDutyExecutionLogResource},
        {Order: 183, New: NewIncidentSlaResource},
        {Order: 221, New: NewCephClusterLabelRuleResource},
        {Order: 279, New: NewScheduledMaintenanceUserOwnerResource},
        {Order: 285, New: NewScheduledMaintenanceNoteTemplateResource},
        {Order: 290, New: NewMonitorGroupResourceResource},
        {Order: 395, New: NewDockerHostUserOwnerResource},
        {Order: 415, New: NewStorageArrayResource},
        {Order: 447, New: NewQueueUserOwnerResource},
    }, nil)
}
