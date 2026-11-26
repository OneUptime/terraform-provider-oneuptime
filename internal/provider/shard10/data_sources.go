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
        {Order: 113, New: NewApiKeyPermissionDataSource},
        {Order: 137, New: NewIncidentDataSource},
        {Order: 151, New: NewIncidentNoteTemplateDataSource},
        {Order: 195, New: NewIncidentEpisodeUserOwnerDataSource},
        {Order: 204, New: NewIncidentOwnerRuleOwningTeamDataSource},
        {Order: 226, New: NewServiceLevelObjectiveTeamOwnerDataSource},
        {Order: 241, New: NewHostLabelRuleDataSource},
        {Order: 326, New: NewTelegramLogDataSource},
        {Order: 332, New: NewUserNotificationLogDataSource},
        {Order: 334, New: NewScheduledMaintenanceTemplateDataSource},
        {Order: 346, New: NewOnCallScheduleTeamOwnerDataSource},
        {Order: 350, New: NewOnCallDutyPolicyTeamOwnerDataSource},
        {Order: 366, New: NewCodeRepositoryDataSource},
        {Order: 440, New: NewSnmpCredentialProfileDataSource},
        {Order: 454, New: NewNetworkTopologySuppressionDataSource},
        {Order: 464, New: NewOpenstackCloudDataSource},
        {Order: 465, New: NewDockerSwarmClusterDataSource},
        {Order: 496, New: NewCloudAccountFeedDataSource},
        {Order: 554, New: NewOpenstackCloudLegacyDataSource},
    }, []LegacyDataSourceAlias{
        {Order: 464, Name: "openstack_cloud", LegacyName: "open_stack_cloud", New: NewOpenstackCloudDataSource, NewLegacy: NewOpenstackCloudLegacyDataSource},
    })
}
