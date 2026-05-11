export type UserRole = 'admin' | 'editor'

export interface User {
  id: number
  full_name: string
  email: string
  role: UserRole
  created_at: string
  updated_at: string
}

export interface Building {
  id: number
  name: string
  address: string
}

export interface Category {
  id: number
  name: string
}

export interface Room {
  id: number
  name: string
  building_id: number
  building?: Building
}

export type ItemStatus = 'active' | 'in_repair' | 'pending_disposal' | 'disposed'

export interface ItemPhoto {
  id: number
  item_id: number
  base64_data: string
  mime_type: string
  created_at: string
}

export interface Item {
  id: number
  barcode: string
  inventory_number: string
  name: string
  category_id: number
  category?: Category
  room_id: number
  room?: Room
  description: string
  status: ItemStatus
  tx_hash?: string
  created_at: string
  updated_at: string
  pending_disposal_at?: string
  disposed_at?: string
  created_by: number
  last_edited_by: number
  creator?: User
  last_editor?: User
  photos?: ItemPhoto[]
}

export type AuditAction =
  | 'created'
  | 'updated'
  | 'disposed'
  | 'moved'
  | 'sent_to_repair'
  | 'returned_from_repair'
  | 'pending_disposal'
  | 'disposal_finalized'

// MovePayload is embedded in AuditEvent.payload when action === 'moved'.
export interface MovePayload {
  from_room_id: number
  from_room_name: string
  from_building_name: string
  to_room_id: number
  to_room_name: string
  to_building_name: string
}

export interface AuditEvent {
  id: number
  item_id: number
  actor_id: number
  actor?: User
  action: AuditAction
  payload: string
  tx_hash: string
  created_at: string
}

export interface LoginRequest {
  email: string
  password: string
}

export interface LoginResponse {
  token: string
  user: User
}

export interface CreateItemRequest {
  barcode: string
  inventory_number: string
  name: string
  category_id: number
  room_id: number
  description?: string
}

export interface UpdateItemRequest {
  description: string
}

export interface MoveToRoomRequest {
  room_id: number
}

export interface ItemFilter {
  category_id?: number
  room_id?: number
  status?: ItemStatus
  date_from?: string
  date_to?: string
}

export interface CreateUserRequest {
  full_name: string
  email: string
  password: string
  role: UserRole
}

export interface UpdateUserRequest {
  full_name?: string
  role?: UserRole
}

export interface DisposalDocument {
  id: number
  item_id: number
  filename: string
  url: string
  uploaded_at: string
  uploaded_by: number
  uploader?: User
}

export interface ApiError {
  error: string
  detail?: string
}
