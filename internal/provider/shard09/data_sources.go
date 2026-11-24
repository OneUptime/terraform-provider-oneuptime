package shard09

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The data sources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterDataSources("shard09", []OrderedDataSource{
        {Order: 23, New: NewWatchlistDataSource},
        {Order: 95, New: NewAiConnectionDataSource},
        {Order: 104, New: NewApiKeyCustomRoleDataSource},
        {Order: 132, New: NewIncidentPublicNoteDataSource},
        {Order: 200, New: NewIncidentEpisodeOnCallRuleOwningTeamDataSource},
        {Order: 214, New: NewSloErrorBudgetPolicyTierDataSource},
        {Order: 301, New: NewTeamTeamOwnerDataSource},
        {Order: 327, New: NewScheduledMaintenanceNoteTemplateDataSource},
        {Order: 339, New: NewOnCallDutyPolicyFeedDataSource},
        {Order: 381, New: NewInventoryItemCustomFieldDataSource},
        {Order: 445, New: NewNetworkSiteAssignmentRuleDataSource},
        {Order: 487, New: NewCloudBudgetDataSource},
    }, nil)
}
