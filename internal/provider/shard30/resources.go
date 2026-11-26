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
        {Order: 40, New: NewDataStoreUserOwnerResource},
        {Order: 62, New: NewTableLabelRuleResource},
        {Order: 64, New: NewProtectedEntityResource},
        {Order: 93, New: NewEscalationRuleTeamResource},
        {Order: 137, New: NewAlertUserOwnerResource},
        {Order: 140, New: NewAlertEpisodeResource},
        {Order: 182, New: NewIncidentEpisodeLabelRuleResource},
        {Order: 191, New: NewServiceLevelObjectiveUserOwnerResource},
        {Order: 197, New: NewIncidentReminderRuleResource},
        {Order: 264, New: NewRunbookTeamOwnerResource},
        {Order: 277, New: NewStatusPageOidcResource},
        {Order: 299, New: NewOnCallScheduleUserOwnerResource},
        {Order: 303, New: NewOnCallDutyPolicyUserOwnerResource},
        {Order: 312, New: NewEscalationRuleScheduleResource},
        {Order: 322, New: NewAiAgentTeamOwnerResource},
        {Order: 324, New: NewAiAgentTaskPullRequestResource},
        {Order: 341, New: NewDashboardOidcResource},
        {Order: 345, New: NewWorkspaceNotificationRuleResource},
        {Order: 375, New: NewNetworkDeviceOwnerRuleResource},
        {Order: 386, New: NewNetworkSiteTeamOwnerResource},
        {Order: 461, New: NewEscalationRuleTeamLegacyResource},
        {Order: 470, New: NewEscalationRuleScheduleLegacyResource},
        {Order: 472, New: NewAiAgentTeamOwnerLegacyResource},
    }, []LegacyResourceAlias{
        {Order: 93, Name: "escalation_rule_team", LegacyName: "team_on_call_duty_escalation_rule", New: NewEscalationRuleTeamResource, NewLegacy: NewEscalationRuleTeamLegacyResource},
        {Order: 312, Name: "escalation_rule_schedule", LegacyName: "schedules_on_call_duty_escalation_rule", New: NewEscalationRuleScheduleResource, NewLegacy: NewEscalationRuleScheduleLegacyResource},
        {Order: 322, Name: "ai_agent_team_owner", LegacyName: "ai_agent_owner_team", New: NewAiAgentTeamOwnerResource, NewLegacy: NewAiAgentTeamOwnerLegacyResource},
    })
}
