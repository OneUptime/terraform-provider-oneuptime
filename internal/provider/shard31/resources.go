package shard31

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The resources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterResources("shard31", []OrderedResource{
        {Order: 46, New: NewResidencyRuleResource},
        {Order: 72, New: NewLlmCostBudgetResource},
        {Order: 73, New: NewLlmModelPriceResource},
        {Order: 98, New: NewIncomingRequestSecretResource},
        {Order: 105, New: NewIncidentCustomFieldResource},
        {Order: 196, New: NewPublicNoteApprovalStepResource},
        {Order: 342, New: NewWorkspaceNotificationRuleOwningTeamResource},
        {Order: 359, New: NewCloudResourceFeedResource},
        {Order: 363, New: NewOnCallTimeLogResource},
        {Order: 427, New: NewServerlessFunctionUserOwnerResource},
    }, nil)
}
