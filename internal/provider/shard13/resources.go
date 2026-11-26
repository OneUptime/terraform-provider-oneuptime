package shard13

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The resources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterResources("shard13", []OrderedResource{
        {Order: 31, New: NewSecurityCaseEvidenceResource},
        {Order: 119, New: NewIncidentMemberResource},
        {Order: 123, New: NewPostmortemStatusResource},
        {Order: 147, New: NewAlertGroupingRuleResource},
        {Order: 166, New: NewIncidentEpisodeInternalNoteResource},
        {Order: 167, New: NewIncidentEpisodeFeedResource},
        {Order: 184, New: NewIncidentMeasurementResource},
        {Order: 218, New: NewDockerSwarmClusterOwnerRuleResource},
        {Order: 282, New: NewScheduledMaintenanceTeamOwnerResource},
        {Order: 290, New: NewRecommendationDismissalResource},
        {Order: 393, New: NewNetworkDeviceLinkRuleResource},
        {Order: 398, New: NewDockerHostTeamOwnerResource},
        {Order: 434, New: NewCloudAccountResource},
        {Order: 442, New: NewDatabaseResource},
        {Order: 450, New: NewQueueTeamOwnerResource},
    }, nil)
}
