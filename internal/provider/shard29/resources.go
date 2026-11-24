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
        {Order: 67, New: NewAuditLogExportResource},
        {Order: 81, New: NewLogRecordingRuleResource},
        {Order: 101, New: NewMonitorTemplateResource},
        {Order: 145, New: NewAlertOnCallRuleOwningTeamResource},
        {Order: 173, New: NewScheduledMaintenanceLabelRuleResource},
        {Order: 207, New: NewDockerHostLabelRuleResource},
        {Order: 240, New: NewStatusPageAnnouncementTemplateResource},
        {Order: 257, New: NewWorkflowUserOwnerResource},
        {Order: 380, New: NewPacketCaptureResource},
        {Order: 448, New: NewQueueLabelRuleResource},
    }, nil)
}
