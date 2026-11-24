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
        {Order: 36, New: NewSecurityCaseStatusDataSource},
        {Order: 39, New: NewSecurityDetectionDataSource},
        {Order: 44, New: NewDataStoreAccessEntryDataSource},
        {Order: 155, New: NewAlertStateTimelineDataSource},
        {Order: 169, New: NewAlertOnCallRuleDataSource},
        {Order: 178, New: NewAlertEpisodePrivacyRuleDataSource},
        {Order: 196, New: NewIncidentLabelRuleDataSource},
        {Order: 244, New: NewOpenstackCloudLabelRuleDataSource},
        {Order: 245, New: NewDockerSwarmClusterLabelRuleDataSource},
        {Order: 250, New: NewStorageArrayOwnerRuleDataSource},
        {Order: 271, New: NewStatusPageSubscriberNotificationTemplateDataSource},
        {Order: 291, New: NewRunbookUserOwnerDataSource},
        {Order: 308, New: NewMonitorTeamOwnerDataSource},
        {Order: 329, New: NewMonitorGroupDataSource},
        {Order: 335, New: NewOnCallScheduleLayerUserDataSource},
        {Order: 369, New: NewLlmLogDataSource},
        {Order: 428, New: NewNetworkDeviceRoleDataSource},
        {Order: 437, New: NewNetworkSiteUserOwnerDataSource},
        {Order: 456, New: NewProxmoxClusterTeamOwnerDataSource},
        {Order: 484, New: NewServerlessFunctionOwnerRuleDataSource},
        {Order: 524, New: NewAuditLogDataSource},
        {Order: 531, New: NewOpenstackCloudLabelRuleLegacyDataSource},
    }, []LegacyDataSourceAlias{
        {Order: 244, Name: "openstack_cloud_label_rule", LegacyName: "open_stack_cloud_label_rule", New: NewOpenstackCloudLabelRuleDataSource, NewLegacy: NewOpenstackCloudLabelRuleLegacyDataSource},
    })
}
