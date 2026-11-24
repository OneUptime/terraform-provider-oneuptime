package shard06

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The data sources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterDataSources("shard06", []OrderedDataSource{
        {Order: 10, New: NewTeamCustomRoleDataSource},
        {Order: 31, New: NewSecurityCaseFeedItemDataSource},
        {Order: 99, New: NewTraceDropFilterDataSource},
        {Order: 136, New: NewIncidentTeamOwnerDataSource},
        {Order: 149, New: NewFormSsoDataSource},
        {Order: 197, New: NewScheduledMaintenanceOwnerRuleDataSource},
        {Order: 226, New: NewMonitorLabelRuleDataSource},
        {Order: 243, New: NewProxmoxClusterLabelRuleDataSource},
        {Order: 274, New: NewStatusPagePrivateUserDataSource},
        {Order: 285, New: NewWorkflowVersionDataSource},
        {Order: 325, New: NewScheduledMaintenanceTemplateTeamOwnerDataSource},
        {Order: 374, New: NewTelemetryIngestionKeyDataSource},
        {Order: 386, New: NewDashboardPrivateUserDataSource},
        {Order: 426, New: NewNetworkDeviceAutoImportRuleDataSource},
        {Order: 442, New: NewNetworkDeviceLinkDataSource},
        {Order: 457, New: NewOpenstackCloudTeamOwnerDataSource},
        {Order: 458, New: NewDockerSwarmClusterTeamOwnerDataSource},
        {Order: 475, New: NewHostDataSource},
        {Order: 513, New: NewLogDataSource},
        {Order: 545, New: NewOpenstackCloudTeamOwnerLegacyDataSource},
    }, []LegacyDataSourceAlias{
        {Order: 457, Name: "openstack_cloud_team_owner", LegacyName: "open_stack_cloud_team_owner", New: NewOpenstackCloudTeamOwnerDataSource, NewLegacy: NewOpenstackCloudTeamOwnerLegacyDataSource},
    })
}
