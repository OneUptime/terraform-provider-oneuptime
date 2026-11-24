package shard15

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The resources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterResources("shard15", []OrderedResource{
        {Order: 16, New: NewExceptionSavedViewResource},
        {Order: 84, New: NewSmtpConfigResource},
        {Order: 186, New: NewSloErrorBudgetPolicyTierResource},
        {Order: 208, New: NewPodmanHostOwnerRuleResource},
        {Order: 229, New: NewWorkflowLabelRuleResource},
        {Order: 241, New: NewStatusPageSubscriberResource},
        {Order: 256, New: NewWorkflowTeamOwnerResource},
        {Order: 291, New: NewOnCallPolicyScheduleResource},
        {Order: 300, New: NewOnCallPolicyLabelRuleResource},
        {Order: 330, New: NewInventoryItemCustomFieldResource},
        {Order: 344, New: NewIncidentVideoCallResource},
        {Order: 348, New: NewMonitorFeedResource},
        {Order: 391, New: NewNetworkSiteAssignmentRuleResource},
        {Order: 425, New: NewServerlessFunctionResource},
    }, nil)
}
