package shard31

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The data sources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterDataSources("shard31", []OrderedDataSource{
        {Order: 35, New: NewSecurityCasePullRequestDataSource},
        {Order: 70, New: NewDataPipelineRunDataSource},
        {Order: 85, New: NewDataAssetLineageDataSource},
        {Order: 86, New: NewDataAlertDataSource},
        {Order: 143, New: NewIncidentTemplateDataSource},
        {Order: 149, New: NewIncidentMemberDataSource},
        {Order: 197, New: NewIncidentEpisodeInternalNoteDataSource},
        {Order: 198, New: NewIncidentEpisodeFeedDataSource},
        {Order: 199, New: NewIncidentEpisodePublicNoteDataSource},
        {Order: 202, New: NewIncidentOnCallRuleOwningTeamDataSource},
        {Order: 208, New: NewScheduledMaintenanceLabelRuleDataSource},
        {Order: 289, New: NewScheduledEventPublicNoteDataSource},
        {Order: 320, New: NewScheduledMaintenanceTeamOwnerDataSource},
        {Order: 425, New: NewScimLogDataSource},
        {Order: 429, New: NewKubernetesClusterUserOwnerDataSource},
        {Order: 478, New: NewCephClusterUserOwnerDataSource},
    }, nil)
}
