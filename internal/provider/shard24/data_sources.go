package shard24

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The data sources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterDataSources("shard24", []OrderedDataSource{
        {Order: 73, New: NewTableTeamOwnerDataSource},
        {Order: 78, New: NewTableLabelRuleDataSource},
        {Order: 100, New: NewTraceScrubRuleDataSource},
        {Order: 126, New: NewIncidentStateDataSource},
        {Order: 140, New: NewIncidentSeverityDataSource},
        {Order: 163, New: NewAlertEpisodeStateTimelineDataSource},
        {Order: 171, New: NewAlertIncidentRuleDataSource},
        {Order: 180, New: NewIncidentEpisodeDataSource},
        {Order: 222, New: NewAlertReminderRuleDataSource},
        {Order: 228, New: NewStatusPageLabelRuleDataSource},
        {Order: 259, New: NewDashboardLabelRuleDataSource},
        {Order: 284, New: NewWorkflowLogDataSource},
        {Order: 305, New: NewStatusPageOidcDataSource},
        {Order: 309, New: NewMonitorUserOwnerDataSource},
        {Order: 312, New: NewStatusPageTeamOwnerDataSource},
        {Order: 350, New: NewEscalationRuleScheduleDataSource},
        {Order: 365, New: NewAiRunDataSource},
        {Order: 371, New: NewTotpAuthDataSource},
        {Order: 388, New: NewDashboardOidcDataSource},
        {Order: 389, New: NewDashboardTeamOwnerDataSource},
        {Order: 398, New: NewOnCallUserGroupDataSource},
        {Order: 459, New: NewProxmoxClusterUserOwnerDataSource},
        {Order: 515, New: NewMetricDataSource},
        {Order: 523, New: NewProfileSampleDataSource},
        {Order: 538, New: NewEscalationRuleScheduleLegacyDataSource},
    }, []LegacyDataSourceAlias{
        {Order: 350, Name: "escalation_rule_schedule", LegacyName: "schedules_on_call_duty_escalation_rule", New: NewEscalationRuleScheduleDataSource, NewLegacy: NewEscalationRuleScheduleLegacyDataSource},
    })
}
