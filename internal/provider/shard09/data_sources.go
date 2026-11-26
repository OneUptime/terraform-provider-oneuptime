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
        {Order: 105, New: NewAiConnectionDataSource},
        {Order: 114, New: NewApiKeyCustomRoleDataSource},
        {Order: 142, New: NewIncidentPublicNoteDataSource},
        {Order: 210, New: NewIncidentEpisodeOnCallRuleOwningTeamDataSource},
        {Order: 224, New: NewSloErrorBudgetPolicyTierDataSource},
        {Order: 311, New: NewTeamTeamOwnerDataSource},
        {Order: 337, New: NewScheduledMaintenanceNoteTemplateDataSource},
        {Order: 349, New: NewOnCallDutyPolicyFeedDataSource},
        {Order: 391, New: NewInventoryItemCustomFieldDataSource},
        {Order: 455, New: NewNetworkSiteAssignmentRuleDataSource},
        {Order: 497, New: NewCloudBudgetDataSource},
    }, nil)
}
