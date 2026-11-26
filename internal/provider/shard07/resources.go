package shard07

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The resources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterResources("shard07", []OrderedResource{
        {Order: 51, New: NewDataCheckResource},
        {Order: 98, New: NewIncomingCallPolicyEscalationRuleResource},
        {Order: 104, New: NewMonitorCustomFieldResource},
        {Order: 156, New: NewAlertEpisodeOwnerRuleResource},
        {Order: 160, New: NewIncidentEpisodeMemberResource},
        {Order: 179, New: NewIncidentEpisodeOnCallRuleOwningTeamResource},
        {Order: 315, New: NewServiceResource},
        {Order: 317, New: NewServiceUserOwnerResource},
        {Order: 344, New: NewDashboardPlaylistResource},
        {Order: 362, New: NewHostFeedResource},
        {Order: 455, New: NewRumApplicationTeamOwnerResource},
        {Order: 457, New: NewRumApplicationLabelRuleResource},
    }, nil)
}
