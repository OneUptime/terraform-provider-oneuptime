package shard04

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The resources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterResources("shard04", []OrderedResource{
        {Order: 37, New: NewSecurityPlaybookRunResource},
        {Order: 86, New: NewApiKeyPermissionResource},
        {Order: 94, New: NewEscalationRuleUserResource},
        {Order: 139, New: NewAlertNoteTemplateResource},
        {Order: 181, New: NewIncidentEpisodePrivacyRuleResource},
        {Order: 183, New: NewIncidentSlaRuleResource},
        {Order: 189, New: NewSloBurnRateRuleResource},
        {Order: 219, New: NewProxmoxClusterLabelRuleResource},
        {Order: 222, New: NewIotFleetOwnerRuleResource},
        {Order: 269, New: NewExternalMcpServerResource},
        {Order: 285, New: NewStatusPageUserOwnerResource},
        {Order: 292, New: NewMonitorGroupTeamOwnerResource},
        {Order: 309, New: NewIncomingCallPolicyUserOwnerResource},
        {Order: 329, New: NewTeamOrLabelBudgetResource},
        {Order: 343, New: NewDashboardUserOwnerResource},
        {Order: 403, New: NewProxmoxClusterResource},
        {Order: 436, New: NewCloudConnectionResource},
        {Order: 462, New: NewEscalationRuleUserLegacyResource},
        {Order: 465, New: NewIotFleetOwnerRuleLegacyResource},
    }, []LegacyResourceAlias{
        {Order: 94, Name: "escalation_rule_user", LegacyName: "users_on_call_duty_escalation_rule", New: NewEscalationRuleUserResource, NewLegacy: NewEscalationRuleUserLegacyResource},
        {Order: 222, Name: "iot_fleet_owner_rule", LegacyName: "io_t_fleet_owner_rule", New: NewIotFleetOwnerRuleResource, NewLegacy: NewIotFleetOwnerRuleLegacyResource},
    })
}
