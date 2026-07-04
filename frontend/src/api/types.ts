export type ApiEnvelope<T> = {
  data: T
  error: null | {
    code: string
    message: string
    details?: unknown
  }
  meta: Record<string, unknown>
}

export type SettingMap = Record<string, string>

export type Category = {
  id: number
  name: string
  slug: string
  description: string
  icon: string
  color: string
  sortOrder: number
  isVisible: boolean
  createdAt: string
  updatedAt: string
}

export type Tag = {
  id: number
  name: string
  slug: string
  color: string
  createdAt: string
  updatedAt: string
}

export type Site = {
  id: number
  categoryId: number
  category?: Category
  title: string
  url: string
  description: string
  iconUrl: string
  fallbackIcon: string
  sortOrder: number
  isPinned: boolean
  isVisible: boolean
  clickCount: number
  tags?: Tag[]
  createdAt: string
  updatedAt: string
}

export type SearchEngine = {
  id: number
  name: string
  slug: string
  searchUrl: string
  icon: string
  sortOrder: number
  isDefault: boolean
  isVisible: boolean
  createdAt: string
  updatedAt: string
}

export type User = {
  id: number
  username: string
  displayName: string
  role: string
}

export type ListMeta = {
  page?: number
  pageSize?: number
  total?: number
}
