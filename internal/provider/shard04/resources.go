package shard04

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The resources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterResources("shard04", []OrderedResource{
        {Order: 36, New: NewSecurityPlaybookRunResource},
        {Order: 82, New: NewApiKeyPermissionResource},
        {Order: 90, New: NewEscalationRuleUserResource},
        {Order: 135, New: NewAlertNoteTemplateResource},
        {Order: 177, New: NewIncidentEpisodePrivacyRuleResource},
        {Order: 179, New: NewIncidentSlaRuleResource},
        {Order: 185, New: NewSloBurnRateRuleResource},
        {Order: 215, New: NewProxmoxClusterLabelRuleResource},
        {Order: 218, New: NewIotFleetOwnerRuleResource},
        {Order: 265, New: NewExternalMcpServerResource},
        {Order: 281, New: NewStatusPageUserOwnerResource},
        {Order: 288, New: NewMonitorGroupTeamOwnerResource},
        {Order: 305, New: NewIncomingCallPolicyUserOwnerResource},
        {Order: 325, New: NewTeamOrLabelBudgetResource},
        {Order: 339, New: NewDashboardUserOwnerResource},
        {Order: 399, New: NewProxmoxClusterResource},
        {Order: 432, New: NewCloudConnectionResource},
        {Order: 458, New: NewEscalationRuleUserLegacyResource},
        {Order: 461, New: NewIotFleetOwnerRuleLegacyResource},
    }, []LegacyResourceAlias{
        {Order: 90, Name: "escalation_rule_user", LegacyName: "users_on_call_duty_escalation_rule", New: NewEscalationRuleUserResource, NewLegacy: NewEscalationRuleUserLegacyResource},
        {Order: 218, Name: "iot_fleet_owner_rule", LegacyName: "io_t_fleet_owner_rule", New: NewIotFleetOwnerRuleResource, NewLegacy: NewIotFleetOwnerRuleLegacyResource},
    })
}
