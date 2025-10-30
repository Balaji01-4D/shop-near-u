import { apiRequest } from './apiClient'
import type { NearbyShop, ShopDetails, ShopProduct, SubscribedShop } from '../types/shop'

type NearbyShopsResponse = NearbyShop[]

type ShopDetailsResponse = ShopDetails

export async function fetchNearbyShops({
  lat,
  lon,
  radius = 5000,
  limit = 10,
}: {
  lat: number
  lon: number
  radius?: number
  limit?: number
}) {
  const params = new URLSearchParams({
    lat: String(lat),
    lon: String(lon),
    radius: String(radius),
    limit: String(limit),
  })

  const { data } = await apiRequest<NearbyShopsResponse>(`shops?${params.toString()}`, {
    method: 'GET',
  })

  return data ?? []
}

export async function fetchShopDetails(shopId: number) {
  const { data } = await apiRequest<ShopDetailsResponse>(`shops/${shopId}`, {
    method: 'GET',
  })
  return data
}

export async function subscribeToShop(shopId: number) {
  const { data } = await apiRequest(`shops/${shopId}/subscribe`, {
    method: 'POST',
  })
  return data
}

export async function unsubscribeFromShop(shopId: number) {
  const { data } = await apiRequest(`shops/${shopId}/unsubscribe`, {
    method: 'POST',
  })
  return data
}

// Get shop products from the real API endpoint
export async function fetchShopProducts(shopId: number): Promise<ShopProduct[]> {
  try {
    const { data } = await apiRequest<ShopProduct[]>(`shops/${shopId}/products`, {
      method: 'GET',
    })
    return data ?? []
  } catch (error) {
    // If there's an error (like shop not found or no products), return empty array
    console.warn(`Failed to fetch products for shop ${shopId}:`, error)
    return []
  }
}

// Get user's subscribed shops
export async function fetchUserSubscriptions(): Promise<SubscribedShop[]> {
  try {
    const { data } = await apiRequest<SubscribedShop[]>('user/subscriptions', {
      method: 'GET',
    })
    return data ?? []
  } catch (error) {
    console.warn('Failed to fetch user subscriptions:', error)
    return []
  }
}
