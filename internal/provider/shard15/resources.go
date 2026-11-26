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
        {Order: 88, New: NewSmtpConfigResource},
        {Order: 190, New: NewSloErrorBudgetPolicyTierResource},
        {Order: 212, New: NewPodmanHostOwnerRuleResource},
        {Order: 233, New: NewWorkflowLabelRuleResource},
        {Order: 245, New: NewStatusPageSubscriberResource},
        {Order: 260, New: NewWorkflowTeamOwnerResource},
        {Order: 295, New: NewOnCallPolicyScheduleResource},
        {Order: 304, New: NewOnCallPolicyLabelRuleResource},
        {Order: 334, New: NewInventoryItemCustomFieldResource},
        {Order: 348, New: NewIncidentVideoCallResource},
        {Order: 352, New: NewMonitorFeedResource},
        {Order: 395, New: NewNetworkSiteAssignmentRuleResource},
        {Order: 429, New: NewServerlessFunctionResource},
    }, nil)
}
