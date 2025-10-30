export type NearbyShop = {
  id: number
  name: string
  address: string
  latitude: number
  longitude: number
  distance: number
}

export type SubscribedShop = {
  id: number
  name: string
  owner_name: string
  type: string
  address: string
  latitude: number
  longitude: number
  subscriber_count: number
  is_open: boolean
}

export type ShopDetails = {
  id: number
  name: string
  subscriber_count: number
  is_subscribed: boolean
}

export type ShopProduct = {
  id: number
  shop_id: number
  catalog_id: number
  price: number
  stock: number
  is_available: boolean
  discount: number
  created_at: string
  updated_at: string
  catalog_product: {
    id: number
    name: string
    brand: string
    category: string
    description: string
    image_url: string
  }
  shop: {
    id: number
    name: string
  }
}
