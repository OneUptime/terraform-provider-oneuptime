package shard31

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The resources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterResources("shard31", []OrderedResource{
        {Order: 47, New: NewResidencyRuleResource},
        {Order: 76, New: NewLlmCostBudgetResource},
        {Order: 77, New: NewLlmModelPriceResource},
        {Order: 102, New: NewIncomingRequestSecretResource},
        {Order: 109, New: NewIncidentCustomFieldResource},
        {Order: 200, New: NewPublicNoteApprovalStepResource},
        {Order: 346, New: NewWorkspaceNotificationRuleOwningTeamResource},
        {Order: 363, New: NewCloudResourceFeedResource},
        {Order: 367, New: NewOnCallTimeLogResource},
        {Order: 431, New: NewServerlessFunctionUserOwnerResource},
    }, nil)
}
