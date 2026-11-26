package shard01

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The data sources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterDataSources("shard01", []OrderedDataSource{
        {Order: 16, New: NewTraceSavedViewDataSource},
        {Order: 76, New: NewDataIssueReportDataSource},
        {Order: 79, New: NewTableOwnerRuleDataSource},
        {Order: 96, New: NewEntityCriticalityDataSource},
        {Order: 106, New: NewMetricRecordingRuleDataSource},
        {Order: 122, New: NewOnCallDutyExecutionLogDataSource},
        {Order: 144, New: NewIncidentTemplateTeamOwnerDataSource},
        {Order: 186, New: NewAlertEpisodeOnCallRuleOwningTeamDataSource},
        {Order: 218, New: NewAlertMeasurementValueDataSource},
        {Order: 237, New: NewStatusPageOwnerRuleDataSource},
        {Order: 268, New: NewDashboardOwnerRuleDataSource},
        {Order: 302, New: NewRunbookRuleDataSource},
        {Order: 361, New: NewCallAndSmsConfigDataSource},
        {Order: 376, New: NewAiRunEventDataSource},
        {Order: 405, New: NewIncidentVideoCallDataSource},
        {Order: 435, New: NewNetworkDeviceLabelRuleDataSource},
        {Order: 514, New: NewQueueUserOwnerDataSource},
        {Order: 520, New: NewRumApplicationLabelRuleDataSource},
        {Order: 531, New: NewCloudCostDataSource},
    }, nil)
}
