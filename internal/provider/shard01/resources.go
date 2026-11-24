package shard01

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The resources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterResources("shard01", []OrderedResource{
        {Order: 2, New: NewTeamResource},
        {Order: 19, New: NewLogDropFilterResource},
        {Order: 60, New: NewPipelineOwnerRuleResource},
        {Order: 99, New: NewMonitorStatusResource},
        {Order: 167, New: NewIncidentOnCallRuleOwningTeamResource},
        {Order: 202, New: NewHostOwnerRuleResource},
        {Order: 219, New: NewIotFleetLabelRuleResource},
        {Order: 264, New: NewPostmortemRuleResource},
        {Order: 326, New: NewExceptionResource},
        {Order: 366, New: NewKubernetesClusterUserOwnerResource},
        {Order: 368, New: NewNetworkDeviceResource},
        {Order: 462, New: NewIotFleetLabelRuleLegacyResource},
    }, []LegacyResourceAlias{
        {Order: 219, Name: "iot_fleet_label_rule", LegacyName: "io_t_fleet_label_rule", New: NewIotFleetLabelRuleResource, NewLegacy: NewIotFleetLabelRuleLegacyResource},
    })
}
