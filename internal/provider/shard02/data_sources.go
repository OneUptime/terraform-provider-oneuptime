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
        {Order: 72, New: NewTableUserOwnerDataSource},
        {Order: 87, New: NewSecurityEventDestinationDataSource},
        {Order: 150, New: NewFormOidcDataSource},
        {Order: 213, New: NewSloBurnRateRuleDataSource},
        {Order: 246, New: NewIotFleetOwnerRuleDataSource},
        {Order: 252, New: NewVcenterOwnerRuleDataSource},
        {Order: 261, New: NewFileDataSource},
        {Order: 313, New: NewStatusPageUserOwnerDataSource},
        {Order: 330, New: NewMonitorGroupTeamOwnerDataSource},
        {Order: 344, New: NewOnCallScheduleLabelRuleDataSource},
        {Order: 357, New: NewProbeTeamOwnerDataSource},
        {Order: 390, New: NewDashboardUserOwnerDataSource},
        {Order: 489, New: NewCloudConnectionDataSource},
        {Order: 532, New: NewIotFleetOwnerRuleLegacyDataSource},
        {Order: 534, New: NewVcenterOwnerRuleLegacyDataSource},
        {Order: 539, New: NewProbeTeamOwnerLegacyDataSource},
    }, []LegacyDataSourceAlias{
        {Order: 246, Name: "iot_fleet_owner_rule", LegacyName: "io_t_fleet_owner_rule", New: NewIotFleetOwnerRuleDataSource, NewLegacy: NewIotFleetOwnerRuleLegacyDataSource},
        {Order: 252, Name: "vcenter_owner_rule", LegacyName: "v_center_owner_rule", New: NewVcenterOwnerRuleDataSource, NewLegacy: NewVcenterOwnerRuleLegacyDataSource},
        {Order: 357, Name: "probe_team_owner", LegacyName: "probe_owner_team", New: NewProbeTeamOwnerDataSource, NewLegacy: NewProbeTeamOwnerLegacyDataSource},
    })
}
