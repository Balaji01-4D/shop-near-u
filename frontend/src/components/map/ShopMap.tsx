import { useEffect, useMemo } from 'react'
import { MapContainer, TileLayer, Marker, Popup, CircleMarker, useMap } from 'react-leaflet'
import { icon, latLngBounds } from 'leaflet'
import type { LatLngBounds, LatLngExpression, LatLngTuple } from 'leaflet'
import 'leaflet/dist/leaflet.css'
import type { NearbyShop } from '../../types/shop'

const defaultMarker = icon({
  iconUrl: new URL('leaflet/dist/images/marker-icon.png', import.meta.url).toString(),
  iconRetinaUrl: new URL('leaflet/dist/images/marker-icon-2x.png', import.meta.url).toString(),
  shadowUrl: new URL('leaflet/dist/images/marker-shadow.png', import.meta.url).toString(),
  iconSize: [25, 41],
  iconAnchor: [12, 41],
  popupAnchor: [1, -34],
  shadowSize: [41, 41],
})

type Coordinates = {
  lat: number
  lon: number
}

type ShopMapProps = {
  shops: NearbyShop[]
  userCoords: Coordinates | null
  fallbackCenter?: Coordinates
  height?: number
}

export function ShopMap({ shops, userCoords, fallbackCenter, height = 360 }: ShopMapProps) {
  const center = useMemo<LatLngExpression>(() => {
    if (userCoords) {
      return [userCoords.lat, userCoords.lon]
    }
    if (shops.length) {
      return [shops[0].latitude, shops[0].longitude]
    }
    if (fallbackCenter) {
      return [fallbackCenter.lat, fallbackCenter.lon]
    }
    return [20.5937, 78.9629]
  }, [fallbackCenter, shops, userCoords])

  const bounds = useMemo<LatLngBounds | null>(() => {
    const points: LatLngTuple[] = []
    if (userCoords) {
      points.push([userCoords.lat, userCoords.lon])
    }
    shops.forEach((shop) => {
      points.push([shop.latitude, shop.longitude])
    })
    if (points.length === 0) {
      return null
    }
    return latLngBounds(points)
  }, [shops, userCoords])

  const mapKey = useMemo(() => {
    const userKey = userCoords ? `${userCoords.lat.toFixed(4)}-${userCoords.lon.toFixed(4)}` : 'no-user'
    const shopKey = shops.map((shop) => shop.id).join('-') || 'no-shops'
    return `${userKey}-${shopKey}`
  }, [shops, userCoords])

  return (
    <MapContainer
      className="shop-map"
      center={center}
      zoom={14}
      scrollWheelZoom={false}
      style={{ height }}
      maxZoom={17}
      minZoom={5}
      key={mapKey}
    >
      <TileLayer
        attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
        url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
      />
      <BoundsUpdater bounds={bounds} center={center} />
      {userCoords ? (
        <>
          <CircleMarker
            center={[userCoords.lat, userCoords.lon]}
            radius={9}
            pathOptions={{ color: '#2563eb', fillColor: '#2563eb', fillOpacity: 0.4, weight: 2 }}
          />
          <Marker position={[userCoords.lat, userCoords.lon]} icon={defaultMarker}>
            <Popup>You are here.</Popup>
          </Marker>
        </>
      ) : null}
      {shops.map((shop) => (
        <Marker key={shop.id} position={[shop.latitude, shop.longitude]} icon={defaultMarker}>
          <Popup>
            <strong>{shop.name}</strong>
            <br />
            {shop.address}
          </Popup>
        </Marker>
      ))}
    </MapContainer>
  )
}

function BoundsUpdater({ bounds, center }: { bounds: LatLngBounds | null; center: LatLngExpression }) {
  const map = useMap()

  useEffect(() => {
    if (bounds) {
      map.fitBounds(bounds, { padding: [32, 32], maxZoom: 16 })
    } else {
      map.setView(center, 14)
    }
  }, [bounds, center, map])

  return null
}
