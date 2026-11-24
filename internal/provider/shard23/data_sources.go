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
        {Order: 79, New: NewDataObservabilitySettingDataSource},
        {Order: 102, New: NewLogRecordingRuleDataSource},
        {Order: 117, New: NewIncomingCallPolicyPhoneNumberDataSource},
        {Order: 144, New: NewPostmortemActionItemDataSource},
        {Order: 166, New: NewAlertEpisodeInternalNoteDataSource},
        {Order: 167, New: NewAlertEpisodeFeedDataSource},
        {Order: 223, New: NewScheduledMaintenanceReminderRuleDataSource},
        {Order: 266, New: NewStatusPageResourceDataSource},
        {Order: 276, New: NewScheduledMaintenanceStateDataSource},
        {Order: 421, New: NewNetworkDeviceDataSource},
        {Order: 423, New: NewNetworkDeviceUserOwnerDataSource},
        {Order: 443, New: NewNetworkDeviceLinkRuleDataSource},
        {Order: 446, New: NewNetworkDeviceNamingRuleDataSource},
        {Order: 494, New: NewCloudResourceOwnerRuleDataSource},
        {Order: 507, New: NewRumApplicationDataSource},
        {Order: 509, New: NewRumApplicationUserOwnerDataSource},
    }, nil)
}
