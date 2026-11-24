package shard31

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The data sources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterDataSources("shard31", []OrderedDataSource{
        {Order: 65, New: NewDataPipelineRunDataSource},
        {Order: 80, New: NewDataAssetLineageDataSource},
        {Order: 81, New: NewDataAlertDataSource},
        {Order: 133, New: NewIncidentTemplateDataSource},
        {Order: 139, New: NewIncidentMemberDataSource},
        {Order: 187, New: NewIncidentEpisodeInternalNoteDataSource},
        {Order: 188, New: NewIncidentEpisodeFeedDataSource},
        {Order: 189, New: NewIncidentEpisodePublicNoteDataSource},
        {Order: 192, New: NewIncidentOnCallRuleOwningTeamDataSource},
        {Order: 198, New: NewScheduledMaintenanceLabelRuleDataSource},
        {Order: 279, New: NewScheduledEventPublicNoteDataSource},
        {Order: 310, New: NewScheduledMaintenanceTeamOwnerDataSource},
        {Order: 415, New: NewScimLogDataSource},
        {Order: 419, New: NewKubernetesClusterUserOwnerDataSource},
        {Order: 468, New: NewCephClusterUserOwnerDataSource},
    }, nil)
}
