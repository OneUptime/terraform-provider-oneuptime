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
        {Order: 54, New: NewDataSecurityScanDataSource},
        {Order: 143, New: NewPostmortemStatusDataSource},
        {Order: 190, New: NewIncidentGroupingRuleDataSource},
        {Order: 193, New: NewIncidentOwnerRuleDataSource},
        {Order: 241, New: NewOpenstackCloudOwnerRuleDataSource},
        {Order: 242, New: NewDockerSwarmClusterOwnerRuleDataSource},
        {Order: 247, New: NewIotFleetLabelRuleDataSource},
        {Order: 253, New: NewVcenterLabelRuleDataSource},
        {Order: 368, New: NewLlmProviderDataSource},
        {Order: 394, New: NewVideoCallConnectionDataSource},
        {Order: 432, New: NewNetworkDeviceDiagnosticDataSource},
        {Order: 479, New: NewCollectorConfigVersionDataSource},
        {Order: 499, New: NewDatabaseUserOwnerDataSource},
        {Order: 530, New: NewOpenstackCloudOwnerRuleLegacyDataSource},
        {Order: 533, New: NewIotFleetLabelRuleLegacyDataSource},
        {Order: 535, New: NewVcenterLabelRuleLegacyDataSource},
    }, []LegacyDataSourceAlias{
        {Order: 241, Name: "openstack_cloud_owner_rule", LegacyName: "open_stack_cloud_owner_rule", New: NewOpenstackCloudOwnerRuleDataSource, NewLegacy: NewOpenstackCloudOwnerRuleLegacyDataSource},
        {Order: 247, Name: "iot_fleet_label_rule", LegacyName: "io_t_fleet_label_rule", New: NewIotFleetLabelRuleDataSource, NewLegacy: NewIotFleetLabelRuleLegacyDataSource},
        {Order: 253, Name: "vcenter_label_rule", LegacyName: "v_center_label_rule", New: NewVcenterLabelRuleDataSource, NewLegacy: NewVcenterLabelRuleLegacyDataSource},
    })
}
