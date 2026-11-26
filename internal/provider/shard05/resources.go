package shard05

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The resources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterResources("shard05", []OrderedResource{
        {Order: 42, New: NewDataStoreOwnerRuleResource},
        {Order: 44, New: NewCustomDataTypeResource},
        {Order: 113, New: NewIncidentTemplateResource},
        {Order: 151, New: NewAlertOwnerRuleResource},
        {Order: 229, New: NewVcenterLabelRuleResource},
        {Order: 307, New: NewOnCallScheduleOwnerRuleResource},
        {Order: 335, New: NewInventoryItemRelationshipResource},
        {Order: 416, New: NewCephClusterResource},
        {Order: 420, New: NewStorageArrayTeamOwnerResource},
        {Order: 456, New: NewRumApplicationUserOwnerResource},
        {Order: 468, New: NewVcenterLabelRuleLegacyResource},
    }, []LegacyResourceAlias{
        {Order: 229, Name: "vcenter_label_rule", LegacyName: "v_center_label_rule", New: NewVcenterLabelRuleResource, NewLegacy: NewVcenterLabelRuleLegacyResource},
    })
}
