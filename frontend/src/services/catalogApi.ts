import { apiRequest } from './apiClient'
import type { CatalogProduct } from '../types/catalog'

type SuggestResponse = {
  products: CatalogProduct[]
}

export async function suggestCatalogProducts(keyword: string, limit = 8) {
  const params = new URLSearchParams()
  params.set('keyword', keyword)
  params.set('limit', String(limit))

  const { data } = await apiRequest<SuggestResponse>(`api/catalog-products/suggest?${params.toString()}`, {
    method: 'GET',
  })

  return data?.products ?? []
}
