package shard16

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The data sources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterDataSources("shard16", []OrderedDataSource{
        {Order: 22, New: NewRuleExceptionDataSource},
        {Order: 28, New: NewSecurityEventConnectionRunDataSource},
        {Order: 32, New: NewSecurityCaseNoteDataSource},
        {Order: 43, New: NewSecurityPlaybookDataSource},
        {Order: 87, New: NewProtectedEntityDataSource},
        {Order: 107, New: NewTracePipelineDataSource},
        {Order: 161, New: NewAlertStateDataSource},
        {Order: 171, New: NewAlertEpisodeDataSource},
        {Order: 196, New: NewIncidentEpisodeTeamOwnerDataSource},
        {Order: 213, New: NewIncidentEpisodeLabelRuleDataSource},
        {Order: 231, New: NewIncidentReminderRuleDataSource},
        {Order: 266, New: NewWorkflowOwnerRuleDataSource},
        {Order: 287, New: NewScheduledMaintenanceEventDataSource},
        {Order: 353, New: NewOnCallPolicyOwnerRuleDataSource},
        {Order: 373, New: NewAiConversationDataSource},
        {Order: 441, New: NewNetworkDeviceDiscoveryScanDataSource},
        {Order: 522, New: NewSavedProductAnalyticsViewDataSource},
    }, nil)
}
