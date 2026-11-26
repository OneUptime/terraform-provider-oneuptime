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
        {Order: 32, New: NewSecurityCasePriorityResource},
        {Order: 46, New: NewComplianceSettingResource},
        {Order: 106, New: NewIncidentStateResource},
        {Order: 163, New: NewIncidentEpisodeStateTimelineResource},
        {Order: 174, New: NewIncidentPrivacyRuleResource},
        {Order: 178, New: NewIncidentEpisodeOnCallRuleResource},
        {Order: 232, New: NewWorkflowOwnerRuleResource},
        {Order: 288, New: NewScheduledMaintenanceTemplateUserOwnerResource},
        {Order: 305, New: NewOnCallPolicyOwnerRuleResource},
        {Order: 364, New: NewServiceFeedResource},
        {Order: 378, New: NewOidCollectionTemplateResource},
        {Order: 411, New: NewDockerSwarmClusterUserOwnerResource},
        {Order: 459, New: NewSavedProductAnalyticsViewResource},
    }, nil)
}
