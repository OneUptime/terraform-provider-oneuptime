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
        {Order: 39, New: NewDataStoreFeedResource},
        {Order: 45, New: NewCheckSettingResource},
        {Order: 58, New: NewTableOwnerRuleResource},
        {Order: 114, New: NewIncidentTemplateTeamOwnerResource},
        {Order: 132, New: NewAlertFeedResource},
        {Order: 180, New: NewIncidentEpisodeOwnerRuleResource},
        {Order: 301, New: NewOnCallDutyPolicyFeedResource},
        {Order: 349, New: NewAlertVideoCallResource},
        {Order: 371, New: NewDockerHostResource},
        {Order: 417, New: NewCephClusterTeamOwnerResource},
        {Order: 421, New: NewStorageArrayUserOwnerResource},
        {Order: 447, New: NewDatabaseLabelRuleResource},
        {Order: 449, New: NewQueueResource},
    }, nil)
}
