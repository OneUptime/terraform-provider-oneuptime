package shard13

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The data sources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterDataSources("shard13", []OrderedDataSource{
        {Order: 26, New: NewHuntRunDataSource},
        {Order: 49, New: NewDataStoreFeedDataSource},
        {Order: 55, New: NewDataSecurityIssueDataSource},
        {Order: 139, New: NewIncidentCustomFieldDataSource},
        {Order: 163, New: NewAlertFeedDataSource},
        {Order: 166, New: NewAlertInternalNoteDataSource},
        {Order: 172, New: NewAlertEpisodeMemberDataSource},
        {Order: 234, New: NewPublicNoteApprovalStepDataSource},
        {Order: 264, New: NewRunbookOwnerRuleDataSource},
        {Order: 321, New: NewScheduledMaintenanceUserOwnerDataSource},
        {Order: 328, New: NewPushNotificationLogDataSource},
        {Order: 362, New: NewProjectUserProfileDataSource},
        {Order: 383, New: NewWebAuthnCredentialDataSource},
        {Order: 387, New: NewExceptionActivityDataSource},
        {Order: 498, New: NewCloudCostRecommendationDataSource},
        {Order: 508, New: NewDatabaseTeamOwnerDataSource},
    }, nil)
}
