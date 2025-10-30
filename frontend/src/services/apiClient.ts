export class ApiError extends Error {
  status: number
  details?: string

  constructor(message: string, status: number, details?: string) {
    super(message)
    this.name = 'ApiError'
    this.status = status
    this.details = details
  }
}

type ApiRequestOptions = Omit<RequestInit, 'body'> & {
  body?: unknown
}

const API_BASE_URL = (import.meta.env.VITE_API_BASE_URL as string | undefined)?.replace(/\/$/, '') ?? ''

export async function apiRequest<T>(path: string, options: ApiRequestOptions = {}) {
  const normalizedPath = path.startsWith('/') ? path : `/${path}`
  const url = `${API_BASE_URL}${normalizedPath}`
  const headers = new Headers(options.headers ?? {})

  let requestBody: BodyInit | undefined

  if (options.body instanceof FormData) {
    requestBody = options.body
  } else if (options.body !== undefined && options.body !== null) {
    headers.set('Content-Type', 'application/json')
    requestBody = JSON.stringify(options.body)
  }

  const response = await fetch(url, {
    ...options,
    headers,
    body: requestBody,
    credentials: 'include',
  })

  const text = await response.text()
  let payload: any

  if (text) {
    try {
      payload = JSON.parse(text)
    } catch (error) {
      throw new ApiError('Unexpected response from server', response.status)
    }
  }

  if (!response.ok || payload?.success === false) {
    const message = payload?.error || payload?.message || response.statusText || 'Request failed'
    throw new ApiError(message, response.status, payload?.error)
  }

  return {
    data: (payload?.data as T) ?? (undefined as T),
    message: payload?.message as string,
  }
}
