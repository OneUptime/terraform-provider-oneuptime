package shard12

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The resources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterResources("shard12", []OrderedResource{
        {Order: 12, New: NewLabelResource},
        {Order: 20, New: NewDetectionRuleResource},
        {Order: 40, New: NewDataStoreTeamOwnerResource},
        {Order: 71, New: NewMetricPipelineRuleResource},
        {Order: 78, New: NewTraceDropFilterResource},
        {Order: 124, New: NewFormSsoResource},
        {Order: 130, New: NewAlertStateTimelineResource},
        {Order: 132, New: NewAlertTeamOwnerResource},
        {Order: 188, New: NewServiceLevelObjectiveTeamOwnerResource},
        {Order: 287, New: NewMonitorGroupResource},
        {Order: 294, New: NewOnCallScheduleTeamOwnerResource},
        {Order: 298, New: NewOnCallDutyPolicyTeamOwnerResource},
        {Order: 398, New: NewPodmanHostUserOwnerResource},
    }, nil)
}
