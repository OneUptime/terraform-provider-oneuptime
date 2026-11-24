package shard09

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The resources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterResources("shard09", []OrderedResource{
        {Order: 6, New: NewTeamPermissionResource},
        {Order: 24, New: NewHuntResource},
        {Order: 29, New: NewSecurityCaseTaskResource},
        {Order: 68, New: NewSecurityLogSourceResource},
        {Order: 70, New: NewLogScrubRuleResource},
        {Order: 142, New: NewAlertEpisodeFeedResource},
        {Order: 165, New: NewIncidentGroupingRuleResource},
        {Order: 195, New: NewScheduledMaintenanceReminderRuleResource},
        {Order: 197, New: NewMonitorOwnerRuleResource},
        {Order: 259, New: NewRunbookExecutionResource},
        {Order: 263, New: NewAiRuleResource},
        {Order: 312, New: NewServiceTeamOwnerResource},
        {Order: 322, New: NewLlmProviderResource},
        {Order: 442, New: NewDatabaseUserOwnerResource},
        {Order: 465, New: NewAiRuleLegacyResource},
    }, []LegacyResourceAlias{
        {Order: 263, Name: "ai_rule", LegacyName: "auto_remediation_rule", New: NewAiRuleResource, NewLegacy: NewAiRuleLegacyResource},
    })
}
