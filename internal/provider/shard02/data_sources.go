package shard02

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The data sources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterDataSources("shard02", []OrderedDataSource{
        {Order: 19, New: NewLogPipelineProcessorDataSource},
        {Order: 77, New: NewTableUserOwnerDataSource},
        {Order: 97, New: NewSecurityEventDestinationDataSource},
        {Order: 160, New: NewFormOidcDataSource},
        {Order: 223, New: NewSloBurnRateRuleDataSource},
        {Order: 256, New: NewIotFleetOwnerRuleDataSource},
        {Order: 262, New: NewVcenterOwnerRuleDataSource},
        {Order: 271, New: NewFileDataSource},
        {Order: 323, New: NewStatusPageUserOwnerDataSource},
        {Order: 340, New: NewMonitorGroupTeamOwnerDataSource},
        {Order: 354, New: NewOnCallScheduleLabelRuleDataSource},
        {Order: 367, New: NewProbeTeamOwnerDataSource},
        {Order: 400, New: NewDashboardUserOwnerDataSource},
        {Order: 499, New: NewCloudConnectionDataSource},
        {Order: 542, New: NewIotFleetOwnerRuleLegacyDataSource},
        {Order: 544, New: NewVcenterOwnerRuleLegacyDataSource},
        {Order: 549, New: NewProbeTeamOwnerLegacyDataSource},
    }, []LegacyDataSourceAlias{
        {Order: 256, Name: "iot_fleet_owner_rule", LegacyName: "io_t_fleet_owner_rule", New: NewIotFleetOwnerRuleDataSource, NewLegacy: NewIotFleetOwnerRuleLegacyDataSource},
        {Order: 262, Name: "vcenter_owner_rule", LegacyName: "v_center_owner_rule", New: NewVcenterOwnerRuleDataSource, NewLegacy: NewVcenterOwnerRuleLegacyDataSource},
        {Order: 367, Name: "probe_team_owner", LegacyName: "probe_owner_team", New: NewProbeTeamOwnerDataSource, NewLegacy: NewProbeTeamOwnerLegacyDataSource},
    })
}
