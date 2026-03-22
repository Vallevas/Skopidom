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

export type ItemStatus = 'active' | 'disposed'

export interface ItemPhoto {
  id: number
  item_id: number
  url: string
  created_at: string
}

export interface Item {
  id: number
  barcode: string
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
  created_by: number
  last_edited_by: number
  creator?: User
  last_editor?: User
  photos?: ItemPhoto[]
}

export type AuditAction = 'created' | 'updated' | 'disposed'

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

export interface ApiError {
  error: string
  detail?: string
}
