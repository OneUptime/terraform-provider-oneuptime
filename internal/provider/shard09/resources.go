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
        {Order: 72, New: NewSecurityLogSourceResource},
        {Order: 74, New: NewLogScrubRuleResource},
        {Order: 146, New: NewAlertEpisodeFeedResource},
        {Order: 169, New: NewIncidentGroupingRuleResource},
        {Order: 199, New: NewScheduledMaintenanceReminderRuleResource},
        {Order: 201, New: NewMonitorOwnerRuleResource},
        {Order: 263, New: NewRunbookExecutionResource},
        {Order: 267, New: NewAiRuleResource},
        {Order: 316, New: NewServiceTeamOwnerResource},
        {Order: 326, New: NewLlmProviderResource},
        {Order: 446, New: NewDatabaseUserOwnerResource},
        {Order: 469, New: NewAiRuleLegacyResource},
    }, []LegacyResourceAlias{
        {Order: 267, Name: "ai_rule", LegacyName: "auto_remediation_rule", New: NewAiRuleResource, NewLegacy: NewAiRuleLegacyResource},
    })
}
