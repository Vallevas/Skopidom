package entity

import "testing"

func TestAuditAction_Category(t *testing.T) {
	tests := []struct {
		name     string
		action   AuditAction
		expected AuditCategory
	}{
		// Status Log actions
		{
			name:     "created action should be status_log",
			action:   ActionCreated,
			expected: CategoryStatusLog,
		},
		{
			name:     "moved action should be status_log",
			action:   ActionMoved,
			expected: CategoryStatusLog,
		},
		{
			name:     "disposed action should be status_log",
			action:   ActionDisposed,
			expected: CategoryStatusLog,
		},
		{
			name:     "pending_disposal action should be status_log",
			action:   ActionPendingDisposal,
			expected: CategoryStatusLog,
		},
		{
			name:     "disposal_finalized action should be status_log",
			action:   ActionDisposalFinalized,
			expected: CategoryStatusLog,
		},
		// Changelog actions
		{
			name:     "updated action should be changelog",
			action:   ActionUpdated,
			expected: CategoryChangelog,
		},
		{
			name:     "sent_to_repair action should be changelog",
			action:   ActionSentToRepair,
			expected: CategoryChangelog,
		},
		{
			name:     "returned_from_repair action should be changelog",
			action:   ActionReturnedFromRepair,
			expected: CategoryChangelog,
		},
		// Unknown action
		{
			name:     "unknown action should default to changelog",
			action:   AuditAction("unknown_action"),
			expected: CategoryChangelog,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.action.Category()
			if got != tt.expected {
				t.Errorf("AuditAction.Category() = %v, want %v", got, tt.expected)
			}
		})
	}
}
