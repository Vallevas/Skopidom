// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/**
 * @title InventoryAudit
 * @dev Immutable audit logger for inventory item lifecycle events.
 * Once an event is logged, it cannot be modified or deleted.
 */
contract InventoryAudit {
    // Event emitted when an audit entry is created on-chain
    event AuditLogged(
        uint256 indexed eventId,
        uint256 indexed itemId,
        string action,
        address indexed actor,
        uint256 timestamp,
        string payload
    );

    // Struct representing an audit event stored on-chain
    struct AuditEvent {
        uint256 eventId;
        uint256 itemId;
        string action;
        address actor;
        uint256 timestamp;
        string payload;
    }

    // Counter for generating unique event IDs
    uint256 private _eventCounter;

    // Mapping from event ID to audit event details
    mapping(uint256 => AuditEvent) private _events;

    // Mapping from item ID to array of event IDs (for querying)
    mapping(uint256 => uint256[]) private _itemEvents;

    /**
     * @dev Constructor - no special initialization needed
     */
    constructor() {
        _eventCounter = 0;
    }

    /**
     * @dev Log a new audit event to the blockchain
     * @param itemId The ID of the inventory item
     * @param action The type of action (created, updated, moved, disposed, etc.)
     * @param payload JSON payload containing additional event data
     * @return eventId The unique ID assigned to this event
     */
    function logEvent(
        uint256 itemId,
        string calldata action,
        string calldata payload
    ) external returns (uint256) {
        _eventCounter++;
        uint256 eventId = _eventCounter;

        AuditEvent memory newEvent = AuditEvent({
            eventId: eventId,
            itemId: itemId,
            action: action,
            actor: msg.sender,
            timestamp: block.timestamp,
            payload: payload
        });

        _events[eventId] = newEvent;
        _itemEvents[itemId].push(eventId);

        emit AuditLogged(
            eventId,
            itemId,
            action,
            msg.sender,
            block.timestamp,
            payload
        );

        return eventId;
    }

    /**
     * @dev Get the total number of logged events
     * @return The total event count
     */
    function getTotalEventCount() external view returns (uint256) {
        return _eventCounter;
    }

    /**
     * @dev Get details of a specific audit event
     * @param eventId The ID of the event to retrieve
     * @return eventId The event ID
     * @return itemId The item ID
     * @return action The action type
     * @return actor The address that logged the event
     * @return timestamp When the event was logged
     * @return payload The JSON payload
     */
    function getEvent(uint256 eventId) external view returns (
        uint256,
        uint256,
        string memory,
        address,
        uint256,
        string memory
    ) {
        AuditEvent memory event = _events[eventId];
        return (
            event.eventId,
            event.itemId,
            event.action,
            event.actor,
            event.timestamp,
            event.payload
        );
    }

    /**
     * @dev Get all event IDs for a specific item
     * @param itemId The ID of the item
     * @return Array of event IDs associated with this item
     */
    function getItemEvents(uint256 itemId) external view returns (uint256[] memory) {
        return _itemEvents[itemId];
    }

    /**
     * @dev Get the number of events for a specific item
     * @param itemId The ID of the item
     * @return The count of events for this item
     */
    function getItemEventCount(uint256 itemId) external view returns (uint256) {
        return _itemEvents[itemId].length;
    }
}
