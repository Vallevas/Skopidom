import type {
  AuditEvent,
  Building,
  Category,
  CreateItemRequest,
  CreateUserRequest,
  Item,
  ItemFilter,
  ItemPhoto,
  LoginRequest,
  LoginResponse,
  MoveToRoomRequest,
  Room,
  UpdateItemRequest,
  UpdateUserRequest,
  User,
} from './types'

const TOKEN_KEY = 'skopidom_token'

export const tokenStorage = {
  get: (): string | null => localStorage.getItem(TOKEN_KEY),
  set: (token: string): void => localStorage.setItem(TOKEN_KEY, token),
  clear: (): void => localStorage.removeItem(TOKEN_KEY),
}

export class ApiClientError extends Error {
  constructor(
    public status: number,
    public body: { error: string; detail?: string },
  ) {
    super(body.error)
  }
}

async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
  const token = tokenStorage.get()
  const response = await fetch(`/api/v1${path}`, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
      ...options.headers,
    },
  })
  if (response.status === 204) return undefined as T
  const data = await response.json()
  if (!response.ok) throw new ApiClientError(response.status, data)
  return data as T
}

export const authApi = {
  login: (body: LoginRequest) =>
    request<LoginResponse>('/auth/login', { method: 'POST', body: JSON.stringify(body) }),
}

export const itemsApi = {
  list: (filter?: ItemFilter) => {
    const params = new URLSearchParams()
    if (filter?.category_id) params.set('category_id', String(filter.category_id))
    if (filter?.room_id) params.set('room_id', String(filter.room_id))
    if (filter?.status) params.set('status', filter.status)
    if (filter?.date_from) params.set('date_from', filter.date_from)
    if (filter?.date_to) params.set('date_to', filter.date_to)
    const qs = params.toString()
    return request<Item[]>(`/items${qs ? `?${qs}` : ''}`)
  },
  getById: (id: number) => request<Item>(`/items/${id}`),
  getByBarcode: (barcode: string) =>
    request<Item>(`/items/barcode/${encodeURIComponent(barcode)}`),
  create: (body: CreateItemRequest) =>
    request<Item>('/items', { method: 'POST', body: JSON.stringify(body) }),
  update: (id: number, body: UpdateItemRequest) =>
    request<Item>(`/items/${id}`, { method: 'PATCH', body: JSON.stringify(body) }),
  moveToRoom: (id: number, body: MoveToRoomRequest) =>
    request<Item>(`/items/${id}/room`, { method: 'PATCH', body: JSON.stringify(body) }),
  dispose: (id: number) => request<void>(`/items/${id}`, { method: 'DELETE' }),
  getAuditLog: (id: number) => request<AuditEvent[]>(`/items/${id}/audit`),

  // Photos
  listPhotos: (id: number) => request<ItemPhoto[]>(`/items/${id}/photos`),
  uploadPhoto: (id: number, file: File) => {
    const form = new FormData()
    form.append('photo', file)
    const token = tokenStorage.get()
    return fetch(`/api/v1/items/${id}/photos`, {
      method: 'POST',
      headers: token ? { Authorization: `Bearer ${token}` } : {},
      body: form,
    }).then(async (res) => {
      const data = await res.json()
      if (!res.ok) throw new ApiClientError(res.status, data)
      return data as ItemPhoto
    })
  },
  deletePhoto: (itemId: number, photoId: number) =>
    request<void>(`/items/${itemId}/photos/${photoId}`, { method: 'DELETE' }),
}

export const usersApi = {
  list: () => request<User[]>('/users'),
  getById: (id: number) => request<User>(`/users/${id}`),
  create: (body: CreateUserRequest) =>
    request<User>('/users', { method: 'POST', body: JSON.stringify(body) }),
  update: (id: number, body: UpdateUserRequest) =>
    request<User>(`/users/${id}`, { method: 'PATCH', body: JSON.stringify(body) }),
  delete: (id: number) => request<void>(`/users/${id}`, { method: 'DELETE' }),
}

export const categoriesApi = {
  list: () => request<Category[]>('/categories'),
  create: (name: string) =>
    request<Category>('/categories', { method: 'POST', body: JSON.stringify({ name }) }),
  update: (id: number, name: string) =>
    request<Category>(`/categories/${id}`, { method: 'PATCH', body: JSON.stringify({ name }) }),
  delete: (id: number) => request<void>(`/categories/${id}`, { method: 'DELETE' }),
}

export const buildingsApi = {
  list: () => request<Building[]>('/buildings'),
  create: (name: string, address: string) =>
    request<Building>('/buildings', { method: 'POST', body: JSON.stringify({ name, address }) }),
  update: (id: number, name: string, address: string) =>
    request<Building>(`/buildings/${id}`, { method: 'PATCH', body: JSON.stringify({ name, address }) }),
  delete: (id: number) => request<void>(`/buildings/${id}`, { method: 'DELETE' }),
}

export const roomsApi = {
  list: (buildingId?: number) => {
    const qs = buildingId ? `?building_id=${buildingId}` : ''
    return request<Room[]>(`/rooms${qs}`)
  },
  create: (name: string, buildingId: number) =>
    request<Room>('/rooms', { method: 'POST', body: JSON.stringify({ name, building_id: buildingId }) }),
  update: (id: number, name: string, buildingId: number) =>
    request<Room>(`/rooms/${id}`, { method: 'PATCH', body: JSON.stringify({ name, building_id: buildingId }) }),
  delete: (id: number) => request<void>(`/rooms/${id}`, { method: 'DELETE' }),
}
