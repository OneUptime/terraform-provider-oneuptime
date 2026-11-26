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
        {Order: 94, New: NewAutomaticResponseActionDataSource},
        {Order: 109, New: NewTraceDropFilterDataSource},
        {Order: 146, New: NewIncidentTeamOwnerDataSource},
        {Order: 159, New: NewFormSsoDataSource},
        {Order: 207, New: NewScheduledMaintenanceOwnerRuleDataSource},
        {Order: 236, New: NewMonitorLabelRuleDataSource},
        {Order: 253, New: NewProxmoxClusterLabelRuleDataSource},
        {Order: 284, New: NewStatusPagePrivateUserDataSource},
        {Order: 295, New: NewWorkflowVersionDataSource},
        {Order: 335, New: NewScheduledMaintenanceTemplateTeamOwnerDataSource},
        {Order: 384, New: NewTelemetryIngestionKeyDataSource},
        {Order: 396, New: NewDashboardPrivateUserDataSource},
        {Order: 436, New: NewNetworkDeviceAutoImportRuleDataSource},
        {Order: 452, New: NewNetworkDeviceLinkDataSource},
        {Order: 467, New: NewOpenstackCloudTeamOwnerDataSource},
        {Order: 468, New: NewDockerSwarmClusterTeamOwnerDataSource},
        {Order: 485, New: NewHostDataSource},
        {Order: 523, New: NewLogDataSource},
        {Order: 555, New: NewOpenstackCloudTeamOwnerLegacyDataSource},
    }, []LegacyDataSourceAlias{
        {Order: 467, Name: "openstack_cloud_team_owner", LegacyName: "open_stack_cloud_team_owner", New: NewOpenstackCloudTeamOwnerDataSource, NewLegacy: NewOpenstackCloudTeamOwnerLegacyDataSource},
    })
}
