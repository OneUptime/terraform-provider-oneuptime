package shard25

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The data sources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterDataSources("shard25", []OrderedDataSource{
        {Order: 61, New: NewCheckSettingDataSource},
        {Order: 67, New: NewDataAssetUsageDataSource},
        {Order: 211, New: NewIncidentEpisodeOwnerRuleDataSource},
        {Order: 216, New: NewIncidentMeasurementValueDataSource},
        {Order: 249, New: NewKubernetesClusterLabelRuleDataSource},
        {Order: 259, New: NewCephClusterLabelRuleDataSource},
        {Order: 325, New: NewWhatsappLogDataSource},
        {Order: 343, New: NewOnCallPolicyScheduleDataSource},
        {Order: 403, New: NewWorkspaceNotificationRuleOwningTeamDataSource},
        {Order: 406, New: NewAlertVideoCallDataSource},
        {Order: 420, New: NewCloudResourceFeedDataSource},
        {Order: 432, New: NewNetworkDeviceTeamOwnerDataSource},
        {Order: 479, New: NewStorageArrayDataSource},
        {Order: 490, New: NewServerlessFunctionDataSource},
        {Order: 518, New: NewRumApplicationTeamOwnerDataSource},
        {Order: 528, New: NewMonitorLogDataSource},
        {Order: 547, New: NewWhatsappLogLegacyDataSource},
    }, []LegacyDataSourceAlias{
        {Order: 325, Name: "whatsapp_log", LegacyName: "whats_app_log", New: NewWhatsappLogDataSource, NewLegacy: NewWhatsappLogLegacyDataSource},
    })
}
