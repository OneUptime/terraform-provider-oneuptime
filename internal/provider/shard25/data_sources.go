package shard25

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The data sources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterDataSources("shard25", []OrderedDataSource{
        {Order: 56, New: NewCheckSettingDataSource},
        {Order: 62, New: NewDataAssetUsageDataSource},
        {Order: 201, New: NewIncidentEpisodeOwnerRuleDataSource},
        {Order: 206, New: NewIncidentMeasurementValueDataSource},
        {Order: 239, New: NewKubernetesClusterLabelRuleDataSource},
        {Order: 249, New: NewCephClusterLabelRuleDataSource},
        {Order: 315, New: NewWhatsappLogDataSource},
        {Order: 333, New: NewOnCallPolicyScheduleDataSource},
        {Order: 393, New: NewWorkspaceNotificationRuleOwningTeamDataSource},
        {Order: 396, New: NewAlertVideoCallDataSource},
        {Order: 410, New: NewCloudResourceFeedDataSource},
        {Order: 422, New: NewNetworkDeviceTeamOwnerDataSource},
        {Order: 469, New: NewStorageArrayDataSource},
        {Order: 480, New: NewServerlessFunctionDataSource},
        {Order: 508, New: NewRumApplicationTeamOwnerDataSource},
        {Order: 518, New: NewMonitorLogDataSource},
        {Order: 537, New: NewWhatsappLogLegacyDataSource},
    }, []LegacyDataSourceAlias{
        {Order: 315, Name: "whatsapp_log", LegacyName: "whats_app_log", New: NewWhatsappLogDataSource, NewLegacy: NewWhatsappLogLegacyDataSource},
    })
}
