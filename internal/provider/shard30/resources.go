package shard30

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The resources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterResources("shard30", []OrderedResource{
        {Order: 28, New: NewSecurityCaseNoteResource},
        {Order: 39, New: NewDataStoreUserOwnerResource},
        {Order: 61, New: NewTableLabelRuleResource},
        {Order: 63, New: NewProtectedEntityResource},
        {Order: 89, New: NewEscalationRuleTeamResource},
        {Order: 133, New: NewAlertUserOwnerResource},
        {Order: 136, New: NewAlertEpisodeResource},
        {Order: 178, New: NewIncidentEpisodeLabelRuleResource},
        {Order: 187, New: NewServiceLevelObjectiveUserOwnerResource},
        {Order: 193, New: NewIncidentReminderRuleResource},
        {Order: 260, New: NewRunbookTeamOwnerResource},
        {Order: 273, New: NewStatusPageOidcResource},
        {Order: 295, New: NewOnCallScheduleUserOwnerResource},
        {Order: 299, New: NewOnCallDutyPolicyUserOwnerResource},
        {Order: 308, New: NewEscalationRuleScheduleResource},
        {Order: 318, New: NewAiAgentTeamOwnerResource},
        {Order: 320, New: NewAiAgentTaskPullRequestResource},
        {Order: 337, New: NewDashboardOidcResource},
        {Order: 341, New: NewWorkspaceNotificationRuleResource},
        {Order: 371, New: NewNetworkDeviceOwnerRuleResource},
        {Order: 382, New: NewNetworkSiteTeamOwnerResource},
        {Order: 457, New: NewEscalationRuleTeamLegacyResource},
        {Order: 466, New: NewEscalationRuleScheduleLegacyResource},
        {Order: 468, New: NewAiAgentTeamOwnerLegacyResource},
    }, []LegacyResourceAlias{
        {Order: 89, Name: "escalation_rule_team", LegacyName: "team_on_call_duty_escalation_rule", New: NewEscalationRuleTeamResource, NewLegacy: NewEscalationRuleTeamLegacyResource},
        {Order: 308, Name: "escalation_rule_schedule", LegacyName: "schedules_on_call_duty_escalation_rule", New: NewEscalationRuleScheduleResource, NewLegacy: NewEscalationRuleScheduleLegacyResource},
        {Order: 318, Name: "ai_agent_team_owner", LegacyName: "ai_agent_owner_team", New: NewAiAgentTeamOwnerResource, NewLegacy: NewAiAgentTeamOwnerLegacyResource},
    })
}
