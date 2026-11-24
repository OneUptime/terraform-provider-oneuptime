package shard10

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The data sources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterDataSources("shard10", []OrderedDataSource{
        {Order: 9, New: NewCustomRolePermissionDataSource},
        {Order: 27, New: NewSecurityEventConnectionDataSource},
        {Order: 103, New: NewApiKeyPermissionDataSource},
        {Order: 127, New: NewIncidentDataSource},
        {Order: 141, New: NewIncidentNoteTemplateDataSource},
        {Order: 185, New: NewIncidentEpisodeUserOwnerDataSource},
        {Order: 194, New: NewIncidentOwnerRuleOwningTeamDataSource},
        {Order: 216, New: NewServiceLevelObjectiveTeamOwnerDataSource},
        {Order: 231, New: NewHostLabelRuleDataSource},
        {Order: 316, New: NewTelegramLogDataSource},
        {Order: 322, New: NewUserNotificationLogDataSource},
        {Order: 324, New: NewScheduledMaintenanceTemplateDataSource},
        {Order: 336, New: NewOnCallScheduleTeamOwnerDataSource},
        {Order: 340, New: NewOnCallDutyPolicyTeamOwnerDataSource},
        {Order: 356, New: NewCodeRepositoryDataSource},
        {Order: 430, New: NewSnmpCredentialProfileDataSource},
        {Order: 444, New: NewNetworkTopologySuppressionDataSource},
        {Order: 454, New: NewOpenstackCloudDataSource},
        {Order: 455, New: NewDockerSwarmClusterDataSource},
        {Order: 486, New: NewCloudAccountFeedDataSource},
        {Order: 544, New: NewOpenstackCloudLegacyDataSource},
    }, []LegacyDataSourceAlias{
        {Order: 454, Name: "openstack_cloud", LegacyName: "open_stack_cloud", New: NewOpenstackCloudDataSource, NewLegacy: NewOpenstackCloudLegacyDataSource},
    })
}
