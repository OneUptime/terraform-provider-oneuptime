package shard23

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The data sources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterDataSources("shard23", []OrderedDataSource{
        {Order: 18, New: NewLogPipelineDataSource},
        {Order: 36, New: NewSecurityCaseAiActionDataSource},
        {Order: 84, New: NewDataObservabilitySettingDataSource},
        {Order: 112, New: NewLogRecordingRuleDataSource},
        {Order: 127, New: NewIncomingCallPolicyPhoneNumberDataSource},
        {Order: 154, New: NewPostmortemActionItemDataSource},
        {Order: 176, New: NewAlertEpisodeInternalNoteDataSource},
        {Order: 177, New: NewAlertEpisodeFeedDataSource},
        {Order: 233, New: NewScheduledMaintenanceReminderRuleDataSource},
        {Order: 276, New: NewStatusPageResourceDataSource},
        {Order: 286, New: NewScheduledMaintenanceStateDataSource},
        {Order: 431, New: NewNetworkDeviceDataSource},
        {Order: 433, New: NewNetworkDeviceUserOwnerDataSource},
        {Order: 453, New: NewNetworkDeviceLinkRuleDataSource},
        {Order: 456, New: NewNetworkDeviceNamingRuleDataSource},
        {Order: 504, New: NewCloudResourceOwnerRuleDataSource},
        {Order: 517, New: NewRumApplicationDataSource},
        {Order: 519, New: NewRumApplicationUserOwnerDataSource},
    }, nil)
}
