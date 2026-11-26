package shard24

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The data sources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterDataSources("shard24", []OrderedDataSource{
        {Order: 78, New: NewTableTeamOwnerDataSource},
        {Order: 83, New: NewTableLabelRuleDataSource},
        {Order: 110, New: NewTraceScrubRuleDataSource},
        {Order: 136, New: NewIncidentStateDataSource},
        {Order: 150, New: NewIncidentSeverityDataSource},
        {Order: 173, New: NewAlertEpisodeStateTimelineDataSource},
        {Order: 181, New: NewAlertIncidentRuleDataSource},
        {Order: 190, New: NewIncidentEpisodeDataSource},
        {Order: 232, New: NewAlertReminderRuleDataSource},
        {Order: 238, New: NewStatusPageLabelRuleDataSource},
        {Order: 269, New: NewDashboardLabelRuleDataSource},
        {Order: 294, New: NewWorkflowLogDataSource},
        {Order: 315, New: NewStatusPageOidcDataSource},
        {Order: 319, New: NewMonitorUserOwnerDataSource},
        {Order: 322, New: NewStatusPageTeamOwnerDataSource},
        {Order: 360, New: NewEscalationRuleScheduleDataSource},
        {Order: 375, New: NewAiRunDataSource},
        {Order: 381, New: NewTotpAuthDataSource},
        {Order: 398, New: NewDashboardOidcDataSource},
        {Order: 399, New: NewDashboardTeamOwnerDataSource},
        {Order: 408, New: NewOnCallUserGroupDataSource},
        {Order: 469, New: NewProxmoxClusterUserOwnerDataSource},
        {Order: 525, New: NewMetricDataSource},
        {Order: 533, New: NewProfileSampleDataSource},
        {Order: 548, New: NewEscalationRuleScheduleLegacyDataSource},
    }, []LegacyDataSourceAlias{
        {Order: 360, Name: "escalation_rule_schedule", LegacyName: "schedules_on_call_duty_escalation_rule", New: NewEscalationRuleScheduleDataSource, NewLegacy: NewEscalationRuleScheduleLegacyDataSource},
    })
}
