package shard22

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The resources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterResources("shard22", []OrderedResource{
        {Order: 11, New: NewApiKeyResource},
        {Order: 31, New: NewSecurityCasePriorityResource},
        {Order: 45, New: NewComplianceSettingResource},
        {Order: 102, New: NewIncidentStateResource},
        {Order: 159, New: NewIncidentEpisodeStateTimelineResource},
        {Order: 170, New: NewIncidentPrivacyRuleResource},
        {Order: 174, New: NewIncidentEpisodeOnCallRuleResource},
        {Order: 228, New: NewWorkflowOwnerRuleResource},
        {Order: 284, New: NewScheduledMaintenanceTemplateUserOwnerResource},
        {Order: 301, New: NewOnCallPolicyOwnerRuleResource},
        {Order: 360, New: NewServiceFeedResource},
        {Order: 374, New: NewOidCollectionTemplateResource},
        {Order: 407, New: NewDockerSwarmClusterUserOwnerResource},
        {Order: 455, New: NewSavedProductAnalyticsViewResource},
    }, nil)
}
