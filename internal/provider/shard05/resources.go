package shard05

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The resources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterResources("shard05", []OrderedResource{
        {Order: 41, New: NewDataStoreOwnerRuleResource},
        {Order: 43, New: NewCustomDataTypeResource},
        {Order: 109, New: NewIncidentTemplateResource},
        {Order: 147, New: NewAlertOwnerRuleResource},
        {Order: 225, New: NewVcenterLabelRuleResource},
        {Order: 303, New: NewOnCallScheduleOwnerRuleResource},
        {Order: 331, New: NewInventoryItemRelationshipResource},
        {Order: 412, New: NewCephClusterResource},
        {Order: 416, New: NewStorageArrayTeamOwnerResource},
        {Order: 452, New: NewRumApplicationUserOwnerResource},
        {Order: 464, New: NewVcenterLabelRuleLegacyResource},
    }, []LegacyResourceAlias{
        {Order: 225, Name: "vcenter_label_rule", LegacyName: "v_center_label_rule", New: NewVcenterLabelRuleResource, NewLegacy: NewVcenterLabelRuleLegacyResource},
    })
}
