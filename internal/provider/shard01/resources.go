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
        {Order: 61, New: NewPipelineOwnerRuleResource},
        {Order: 103, New: NewMonitorStatusResource},
        {Order: 171, New: NewIncidentOnCallRuleOwningTeamResource},
        {Order: 206, New: NewHostOwnerRuleResource},
        {Order: 223, New: NewIotFleetLabelRuleResource},
        {Order: 268, New: NewPostmortemRuleResource},
        {Order: 330, New: NewExceptionResource},
        {Order: 370, New: NewKubernetesClusterUserOwnerResource},
        {Order: 372, New: NewNetworkDeviceResource},
        {Order: 466, New: NewIotFleetLabelRuleLegacyResource},
    }, []LegacyResourceAlias{
        {Order: 223, Name: "iot_fleet_label_rule", LegacyName: "io_t_fleet_label_rule", New: NewIotFleetLabelRuleResource, NewLegacy: NewIotFleetLabelRuleLegacyResource},
    })
}
