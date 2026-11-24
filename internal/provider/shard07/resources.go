package shard07

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The resources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterResources("shard07", []OrderedResource{
        {Order: 50, New: NewDataCheckResource},
        {Order: 94, New: NewIncomingCallPolicyEscalationRuleResource},
        {Order: 100, New: NewMonitorCustomFieldResource},
        {Order: 152, New: NewAlertEpisodeOwnerRuleResource},
        {Order: 156, New: NewIncidentEpisodeMemberResource},
        {Order: 175, New: NewIncidentEpisodeOnCallRuleOwningTeamResource},
        {Order: 311, New: NewServiceResource},
        {Order: 313, New: NewServiceUserOwnerResource},
        {Order: 340, New: NewDashboardPlaylistResource},
        {Order: 358, New: NewHostFeedResource},
        {Order: 451, New: NewRumApplicationTeamOwnerResource},
        {Order: 453, New: NewRumApplicationLabelRuleResource},
    }, nil)
}
