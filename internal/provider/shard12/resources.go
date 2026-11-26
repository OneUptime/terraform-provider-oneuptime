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
        {Order: 41, New: NewDataStoreTeamOwnerResource},
        {Order: 75, New: NewMetricPipelineRuleResource},
        {Order: 82, New: NewTraceDropFilterResource},
        {Order: 128, New: NewFormSsoResource},
        {Order: 134, New: NewAlertStateTimelineResource},
        {Order: 136, New: NewAlertTeamOwnerResource},
        {Order: 192, New: NewServiceLevelObjectiveTeamOwnerResource},
        {Order: 291, New: NewMonitorGroupResource},
        {Order: 298, New: NewOnCallScheduleTeamOwnerResource},
        {Order: 302, New: NewOnCallDutyPolicyTeamOwnerResource},
        {Order: 402, New: NewPodmanHostUserOwnerResource},
    }, nil)
}
