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
        {Order: 52, New: NewDataProductResource},
        {Order: 92, New: NewEscalationRuleResource},
        {Order: 130, New: NewAlertStateResource},
        {Order: 154, New: NewAlertEpisodeOnCallRuleResource},
        {Order: 239, New: NewStatusPageGroupResource},
        {Order: 257, New: NewScheduledMaintenanceFeedResource},
        {Order: 281, New: NewMonitorUserOwnerResource},
        {Order: 284, New: NewStatusPageTeamOwnerResource},
        {Order: 308, New: NewIncomingCallPolicyTeamOwnerResource},
        {Order: 320, New: NewProbeUserOwnerResource},
        {Order: 331, New: NewSourceMapResource},
        {Order: 342, New: NewDashboardTeamOwnerResource},
        {Order: 354, New: NewDockerHostFeedResource},
    }, nil)
}
