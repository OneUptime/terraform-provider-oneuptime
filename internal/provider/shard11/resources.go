package shard11

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The resources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterResources("shard11", []OrderedResource{
        {Order: 27, New: NewSecurityCaseResource},
        {Order: 38, New: NewDataStoreFeedResource},
        {Order: 44, New: NewCheckSettingResource},
        {Order: 57, New: NewTableOwnerRuleResource},
        {Order: 110, New: NewIncidentTemplateTeamOwnerResource},
        {Order: 128, New: NewAlertFeedResource},
        {Order: 176, New: NewIncidentEpisodeOwnerRuleResource},
        {Order: 297, New: NewOnCallDutyPolicyFeedResource},
        {Order: 345, New: NewAlertVideoCallResource},
        {Order: 367, New: NewDockerHostResource},
        {Order: 413, New: NewCephClusterTeamOwnerResource},
        {Order: 417, New: NewStorageArrayUserOwnerResource},
        {Order: 443, New: NewDatabaseLabelRuleResource},
        {Order: 445, New: NewQueueResource},
    }, nil)
}
