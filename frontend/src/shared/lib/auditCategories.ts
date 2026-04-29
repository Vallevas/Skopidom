import type { AuditAction, AuditEvent } from '@/shared/api/types'

/**
 * Audit category represents the logical grouping of audit actions.
 * - status_log: lifecycle events (creation, movement, disposal)
 * - changelog: modifications and repairs
 */
export type AuditCategory = 'status_log' | 'changelog'

/**
 * Actions that belong to the Status Log category.
 * These track the lifecycle and location changes of items.
 */
export const STATUS_LOG_ACTIONS: readonly AuditAction[] = [
  'created',
  'moved',
  'disposed',
  'pending_disposal',
  'disposal_finalized',
] as const

/**
 * Actions that belong to the Changelog category.
 * These track modifications and repairs of items.
 */
export const CHANGELOG_ACTIONS: readonly AuditAction[] = [
  'updated',
  'sent_to_repair',
  'returned_from_repair',
] as const

/**
 * Determines which category an audit action belongs to.
 * @param action - The audit action to categorize
 * @returns The category of the action
 */
export function getAuditCategory(action: AuditAction): AuditCategory {
  if (STATUS_LOG_ACTIONS.includes(action)) {
    return 'status_log'
  }
  return 'changelog'
}

/**
 * Filters audit events by category.
 * @param events - Array of audit events to filter
 * @param category - The category to filter by
 * @returns Filtered array of events belonging to the specified category
 */
export function filterEventsByCategory(
  events: AuditEvent[],
  category: AuditCategory
): AuditEvent[] {
  const actions = category === 'status_log' ? STATUS_LOG_ACTIONS : CHANGELOG_ACTIONS
  return events.filter((e) => actions.includes(e.action))
}
