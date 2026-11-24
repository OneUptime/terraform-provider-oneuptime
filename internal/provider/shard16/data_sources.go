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
        {Order: 40, New: NewSecurityPlaybookDataSource},
        {Order: 82, New: NewProtectedEntityDataSource},
        {Order: 97, New: NewTracePipelineDataSource},
        {Order: 151, New: NewAlertStateDataSource},
        {Order: 161, New: NewAlertEpisodeDataSource},
        {Order: 186, New: NewIncidentEpisodeTeamOwnerDataSource},
        {Order: 203, New: NewIncidentEpisodeLabelRuleDataSource},
        {Order: 221, New: NewIncidentReminderRuleDataSource},
        {Order: 256, New: NewWorkflowOwnerRuleDataSource},
        {Order: 277, New: NewScheduledMaintenanceEventDataSource},
        {Order: 343, New: NewOnCallPolicyOwnerRuleDataSource},
        {Order: 363, New: NewAiConversationDataSource},
        {Order: 431, New: NewNetworkDeviceDiscoveryScanDataSource},
        {Order: 512, New: NewSavedProductAnalyticsViewDataSource},
    }, nil)
}
