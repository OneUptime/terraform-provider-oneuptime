package shard30

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The data sources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterDataSources("shard30", []OrderedDataSource{
        {Order: 53, New: NewDataStoreLabelRuleDataSource},
        {Order: 56, New: NewDataSecurityIssueFeedItemDataSource},
        {Order: 118, New: NewOnCallPolicyCustomFieldDataSource},
        {Order: 121, New: NewEscalationRuleUserDataSource},
        {Order: 129, New: NewIncomingCallLogItemDataSource},
        {Order: 184, New: NewAlertLabelRuleDataSource},
        {Order: 201, New: NewIncidentOnCallRuleDataSource},
        {Order: 212, New: NewIncidentEpisodePrivacyRuleDataSource},
        {Order: 222, New: NewServiceLevelObjectiveDataSource},
        {Order: 270, New: NewMonitorStatusEventDataSource},
        {Order: 348, New: NewScheduledMaintenanceExternalCalendarDataSource},
        {Order: 357, New: NewIncomingCallPolicyUserOwnerDataSource},
        {Order: 389, New: NewDsymFileDataSource},
        {Order: 390, New: NewInventoryItemDataSource},
        {Order: 407, New: NewWorkspaceNotificationSummaryDataSource},
        {Order: 410, New: NewKubernetesClusterFeedDataSource},
        {Order: 413, New: NewCephClusterFeedDataSource},
        {Order: 439, New: NewNetworkAlertPolicyDataSource},
        {Order: 444, New: NewNetworkInterfaceDataSource},
        {Order: 539, New: NewEscalationRuleUserLegacyDataSource},
        {Order: 551, New: NewDsymFileLegacyDataSource},
    }, []LegacyDataSourceAlias{
        {Order: 121, Name: "escalation_rule_user", LegacyName: "users_on_call_duty_escalation_rule", New: NewEscalationRuleUserDataSource, NewLegacy: NewEscalationRuleUserLegacyDataSource},
        {Order: 389, Name: "dsym_file", LegacyName: "d_sym", New: NewDsymFileDataSource, NewLegacy: NewDsymFileLegacyDataSource},
    })
}
