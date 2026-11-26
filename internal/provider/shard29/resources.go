package shard29

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The resources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterResources("shard29", []OrderedResource{
        {Order: 5, New: NewTeamMemberCustomFieldResource},
        {Order: 71, New: NewAuditLogExportResource},
        {Order: 85, New: NewLogRecordingRuleResource},
        {Order: 105, New: NewMonitorTemplateResource},
        {Order: 149, New: NewAlertOnCallRuleOwningTeamResource},
        {Order: 177, New: NewScheduledMaintenanceLabelRuleResource},
        {Order: 211, New: NewDockerHostLabelRuleResource},
        {Order: 244, New: NewStatusPageAnnouncementTemplateResource},
        {Order: 261, New: NewWorkflowUserOwnerResource},
        {Order: 384, New: NewPacketCaptureResource},
        {Order: 452, New: NewQueueLabelRuleResource},
    }, nil)
}
