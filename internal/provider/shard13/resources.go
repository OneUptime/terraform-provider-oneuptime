package shard13

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The resources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterResources("shard13", []OrderedResource{
        {Order: 30, New: NewSecurityCaseEvidenceResource},
        {Order: 115, New: NewIncidentMemberResource},
        {Order: 119, New: NewPostmortemStatusResource},
        {Order: 143, New: NewAlertGroupingRuleResource},
        {Order: 162, New: NewIncidentEpisodeInternalNoteResource},
        {Order: 163, New: NewIncidentEpisodeFeedResource},
        {Order: 180, New: NewIncidentMeasurementResource},
        {Order: 214, New: NewDockerSwarmClusterOwnerRuleResource},
        {Order: 278, New: NewScheduledMaintenanceTeamOwnerResource},
        {Order: 286, New: NewRecommendationDismissalResource},
        {Order: 389, New: NewNetworkDeviceLinkRuleResource},
        {Order: 394, New: NewDockerHostTeamOwnerResource},
        {Order: 430, New: NewCloudAccountResource},
        {Order: 438, New: NewDatabaseResource},
        {Order: 446, New: NewQueueTeamOwnerResource},
    }, nil)
}
