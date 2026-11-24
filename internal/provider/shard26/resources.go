package shard26

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The resources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterResources("shard26", []OrderedResource{
        {Order: 4, New: NewTeamMemberResource},
        {Order: 51, New: NewDataProductResource},
        {Order: 88, New: NewEscalationRuleResource},
        {Order: 126, New: NewAlertStateResource},
        {Order: 150, New: NewAlertEpisodeOnCallRuleResource},
        {Order: 235, New: NewStatusPageGroupResource},
        {Order: 253, New: NewScheduledMaintenanceFeedResource},
        {Order: 277, New: NewMonitorUserOwnerResource},
        {Order: 280, New: NewStatusPageTeamOwnerResource},
        {Order: 304, New: NewIncomingCallPolicyTeamOwnerResource},
        {Order: 316, New: NewProbeUserOwnerResource},
        {Order: 327, New: NewSourceMapResource},
        {Order: 338, New: NewDashboardTeamOwnerResource},
        {Order: 350, New: NewDockerHostFeedResource},
    }, nil)
}
