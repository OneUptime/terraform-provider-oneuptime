package shard03

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The data sources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterDataSources("shard03", []OrderedDataSource{
        {Order: 7, New: NewTeamPermissionDataSource},
        {Order: 58, New: NewDataSecurityAiDecisionDataSource},
        {Order: 59, New: NewDataSecurityScanDataSource},
        {Order: 153, New: NewPostmortemStatusDataSource},
        {Order: 200, New: NewIncidentGroupingRuleDataSource},
        {Order: 203, New: NewIncidentOwnerRuleDataSource},
        {Order: 251, New: NewOpenstackCloudOwnerRuleDataSource},
        {Order: 252, New: NewDockerSwarmClusterOwnerRuleDataSource},
        {Order: 257, New: NewIotFleetLabelRuleDataSource},
        {Order: 263, New: NewVcenterLabelRuleDataSource},
        {Order: 378, New: NewLlmProviderDataSource},
        {Order: 404, New: NewVideoCallConnectionDataSource},
        {Order: 442, New: NewNetworkDeviceDiagnosticDataSource},
        {Order: 489, New: NewCollectorConfigVersionDataSource},
        {Order: 509, New: NewDatabaseUserOwnerDataSource},
        {Order: 540, New: NewOpenstackCloudOwnerRuleLegacyDataSource},
        {Order: 543, New: NewIotFleetLabelRuleLegacyDataSource},
        {Order: 545, New: NewVcenterLabelRuleLegacyDataSource},
    }, []LegacyDataSourceAlias{
        {Order: 251, Name: "openstack_cloud_owner_rule", LegacyName: "open_stack_cloud_owner_rule", New: NewOpenstackCloudOwnerRuleDataSource, NewLegacy: NewOpenstackCloudOwnerRuleLegacyDataSource},
        {Order: 257, Name: "iot_fleet_label_rule", LegacyName: "io_t_fleet_label_rule", New: NewIotFleetLabelRuleDataSource, NewLegacy: NewIotFleetLabelRuleLegacyDataSource},
        {Order: 263, Name: "vcenter_label_rule", LegacyName: "v_center_label_rule", New: NewVcenterLabelRuleDataSource, NewLegacy: NewVcenterLabelRuleLegacyDataSource},
    })
}
