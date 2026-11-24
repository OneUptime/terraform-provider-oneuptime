package shard30

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The data sources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterDataSources("shard30", []OrderedDataSource{
        {Order: 50, New: NewDataStoreLabelRuleDataSource},
        {Order: 53, New: NewDataSecurityIssueFeedItemDataSource},
        {Order: 108, New: NewOnCallPolicyCustomFieldDataSource},
        {Order: 111, New: NewEscalationRuleUserDataSource},
        {Order: 119, New: NewIncomingCallLogItemDataSource},
        {Order: 174, New: NewAlertLabelRuleDataSource},
        {Order: 191, New: NewIncidentOnCallRuleDataSource},
        {Order: 202, New: NewIncidentEpisodePrivacyRuleDataSource},
        {Order: 212, New: NewServiceLevelObjectiveDataSource},
        {Order: 260, New: NewMonitorStatusEventDataSource},
        {Order: 338, New: NewScheduledMaintenanceExternalCalendarDataSource},
        {Order: 347, New: NewIncomingCallPolicyUserOwnerDataSource},
        {Order: 379, New: NewDsymFileDataSource},
        {Order: 380, New: NewInventoryItemDataSource},
        {Order: 397, New: NewWorkspaceNotificationSummaryDataSource},
        {Order: 400, New: NewKubernetesClusterFeedDataSource},
        {Order: 403, New: NewCephClusterFeedDataSource},
        {Order: 429, New: NewNetworkAlertPolicyDataSource},
        {Order: 434, New: NewNetworkInterfaceDataSource},
        {Order: 529, New: NewEscalationRuleUserLegacyDataSource},
        {Order: 541, New: NewDsymFileLegacyDataSource},
    }, []LegacyDataSourceAlias{
        {Order: 111, Name: "escalation_rule_user", LegacyName: "users_on_call_duty_escalation_rule", New: NewEscalationRuleUserDataSource, NewLegacy: NewEscalationRuleUserLegacyDataSource},
        {Order: 379, Name: "dsym_file", LegacyName: "d_sym", New: NewDsymFileDataSource, NewLegacy: NewDsymFileLegacyDataSource},
    })
}
