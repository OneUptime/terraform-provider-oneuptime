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
        {Order: 71, New: NewDataIssueReportDataSource},
        {Order: 74, New: NewTableOwnerRuleDataSource},
        {Order: 86, New: NewEntityCriticalityDataSource},
        {Order: 96, New: NewMetricRecordingRuleDataSource},
        {Order: 112, New: NewOnCallDutyExecutionLogDataSource},
        {Order: 134, New: NewIncidentTemplateTeamOwnerDataSource},
        {Order: 176, New: NewAlertEpisodeOnCallRuleOwningTeamDataSource},
        {Order: 208, New: NewAlertMeasurementValueDataSource},
        {Order: 227, New: NewStatusPageOwnerRuleDataSource},
        {Order: 258, New: NewDashboardOwnerRuleDataSource},
        {Order: 292, New: NewRunbookRuleDataSource},
        {Order: 351, New: NewCallAndSmsConfigDataSource},
        {Order: 366, New: NewAiRunEventDataSource},
        {Order: 395, New: NewIncidentVideoCallDataSource},
        {Order: 425, New: NewNetworkDeviceLabelRuleDataSource},
        {Order: 504, New: NewQueueUserOwnerDataSource},
        {Order: 510, New: NewRumApplicationLabelRuleDataSource},
        {Order: 521, New: NewCloudCostDataSource},
    }, nil)
}
