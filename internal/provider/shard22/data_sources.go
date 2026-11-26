package shard22

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The data sources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterDataSources("shard22", []OrderedDataSource{
        {Order: 1, New: NewProbeDataSource},
        {Order: 11, New: NewTeamComplianceSettingDataSource},
        {Order: 39, New: NewSecurityCaseStatusDataSource},
        {Order: 42, New: NewSecurityDetectionDataSource},
        {Order: 47, New: NewDataStoreAccessEntryDataSource},
        {Order: 165, New: NewAlertStateTimelineDataSource},
        {Order: 179, New: NewAlertOnCallRuleDataSource},
        {Order: 188, New: NewAlertEpisodePrivacyRuleDataSource},
        {Order: 206, New: NewIncidentLabelRuleDataSource},
        {Order: 254, New: NewOpenstackCloudLabelRuleDataSource},
        {Order: 255, New: NewDockerSwarmClusterLabelRuleDataSource},
        {Order: 260, New: NewStorageArrayOwnerRuleDataSource},
        {Order: 281, New: NewStatusPageSubscriberNotificationTemplateDataSource},
        {Order: 301, New: NewRunbookUserOwnerDataSource},
        {Order: 318, New: NewMonitorTeamOwnerDataSource},
        {Order: 339, New: NewMonitorGroupDataSource},
        {Order: 345, New: NewOnCallScheduleLayerUserDataSource},
        {Order: 379, New: NewLlmLogDataSource},
        {Order: 438, New: NewNetworkDeviceRoleDataSource},
        {Order: 447, New: NewNetworkSiteUserOwnerDataSource},
        {Order: 466, New: NewProxmoxClusterTeamOwnerDataSource},
        {Order: 494, New: NewServerlessFunctionOwnerRuleDataSource},
        {Order: 534, New: NewAuditLogDataSource},
        {Order: 541, New: NewOpenstackCloudLabelRuleLegacyDataSource},
    }, []LegacyDataSourceAlias{
        {Order: 254, Name: "openstack_cloud_label_rule", LegacyName: "open_stack_cloud_label_rule", New: NewOpenstackCloudLabelRuleDataSource, NewLegacy: NewOpenstackCloudLabelRuleLegacyDataSource},
    })
}
