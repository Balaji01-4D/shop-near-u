import { useCallback, useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { ShopMap } from '../components/map/ShopMap'
import { fetchNearbyShops, fetchShopDetails, subscribeToShop, unsubscribeFromShop } from '../services/shopApi'
import { useAuth } from '../context/AuthContext'
import { updatePageTitle, pageTitles } from '../utils/pageTitle'
import { ApiError } from '../services/apiClient'
import type { NearbyShop, ShopDetails } from '../types/shop'

const DEFAULT_COORDS = { lat: 13.0827, lon: 80.2707 }

function formatDistance(distance: number) {
  if (!Number.isFinite(distance)) {
    return null
  }
  const kilometers = distance / 1000
  if (kilometers < 1) {
    return `${(kilometers * 1000).toFixed(0)} m`
  }
  return `${kilometers.toFixed(1)} km`
}

export function ShopsPage() {
  const { user } = useAuth()
  const [coords, setCoords] = useState<{ lat: number; lon: number } | null>(null)
  const [locationMessage, setLocationMessage] = useState<string | null>(null)

  // Set page title
  useEffect(() => {
    updatePageTitle(pageTitles.shops)
  }, [])

  const [shops, setShops] = useState<NearbyShop[]>([])
  const [shopMeta, setShopMeta] = useState<Record<number, ShopDetails>>({})
  const [shopsLoading, setShopsLoading] = useState(false)
  const [shopsError, setShopsError] = useState<string | null>(null)
  const [subscriptionLoading, setSubscriptionLoading] = useState<Record<number, boolean>>({})

  const [searchTerm, setSearchTerm] = useState('')
  const [viewMode, setViewMode] = useState<'list' | 'map'>('list')
  const [radius, setRadius] = useState(5000)

  const loadNearby = useCallback(
    async (latitude: number, longitude: number) => {
      setShopsLoading(true)
      setShopsError(null)
      try {
        const results = await fetchNearbyShops({ lat: latitude, lon: longitude, radius, limit: 20 })
        setShops(results)

        if (user) {
          const extrasEntries: Array<[number, ShopDetails]> = []
          await Promise.all(
            results.map(async (shop) => {
              try {
                const details = await fetchShopDetails(shop.id)
                if (details) {
                  extrasEntries.push([shop.id, details])
                }
              } catch (error) {
                if (error instanceof ApiError && error.status === 401) {
                  return
                }
              }
            }),
          )
          setShopMeta(Object.fromEntries(extrasEntries))
        } else {
          setShopMeta({})
        }
      } catch (error) {
        const message = error instanceof Error ? error.message : 'Unable to fetch nearby shops right now.'
        setShopsError(message)
      } finally {
        setShopsLoading(false)
      }
    },
    [user, radius],
  )

  const handleSubscription = async (shopId: number, isCurrentlySubscribed: boolean) => {
    if (!user) {
      alert('Please log in to subscribe to shops')
      return
    }

    try {
      setSubscriptionLoading(prev => ({ ...prev, [shopId]: true }))
      
      if (isCurrentlySubscribed) {
        await unsubscribeFromShop(shopId)
      } else {
        await subscribeToShop(shopId)
      }

      // Update shop meta with new subscription status
      setShopMeta(prev => ({
        ...prev,
        [shopId]: {
          ...prev[shopId],
          is_subscribed: !isCurrentlySubscribed,
          subscriber_count: isCurrentlySubscribed 
            ? Math.max(0, prev[shopId]?.subscriber_count - 1 || 0)
            : (prev[shopId]?.subscriber_count || 0) + 1
        }
      }))
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Failed to update subscription'
      alert(message)
    } finally {
      setSubscriptionLoading(prev => ({ ...prev, [shopId]: false }))
    }
  }

  useEffect(() => {
    let cancelled = false

    const requestBrowserLocation = () => {
      if (!('geolocation' in navigator)) {
        setLocationMessage('Location access is not available in this browser. Showing results near central Chennai as a fallback.')
        setCoords(DEFAULT_COORDS)
        void loadNearby(DEFAULT_COORDS.lat, DEFAULT_COORDS.lon)
        return
      }

      navigator.geolocation.getCurrentPosition(
        (position) => {
          if (cancelled) {
            return
          }
          const { latitude, longitude } = position.coords
          setCoords({ lat: latitude, lon: longitude })
          setLocationMessage('Using your device location for precise shop results.')
          void loadNearby(latitude, longitude)
        },
        () => {
          if (cancelled) {
            return
          }
          setLocationMessage('Unable to access your location. Showing a curated list near central Chennai. Enable location for personalised results.')
          setCoords(DEFAULT_COORDS)
          void loadNearby(DEFAULT_COORDS.lat, DEFAULT_COORDS.lon)
        },
        { enableHighAccuracy: false, timeout: 5000 },
      )
    }

    requestBrowserLocation()

    return () => {
      cancelled = true
    }
  }, [loadNearby])

  const handleUseLocation = () => {
    if (!('geolocation' in navigator)) {
      setLocationMessage('Location services are not supported in this browser.')
      return
    }

    setLocationMessage('Detecting your location…')
    navigator.geolocation.getCurrentPosition(
      (position) => {
        const { latitude, longitude } = position.coords
        setCoords({ lat: latitude, lon: longitude })
        setLocationMessage('Using your device location for precise shop results.')
        void loadNearby(latitude, longitude)
      },
      () => {
        setLocationMessage('Unable to use your location. Continuing with fallback results.')
      },
    )
  }

  const filteredShops = shops.filter((shop) => {
    const keyword = searchTerm.trim().toLowerCase()
    if (!keyword) return true
    return shop.name.toLowerCase().includes(keyword) || shop.address.toLowerCase().includes(keyword)
  })

  return (
    <div className="shops-page">
      <div className="page-header">
        <div>
          <h1>Nearby Shops</h1>
          <p className="text-muted">Discover local stores around your area with real-time availability and distance information.</p>
        </div>
      </div>

      {locationMessage ? <div className="alert">{locationMessage}</div> : null}

      <div className="modern-search-section">
        <div className="search-container">
          <div className="search-input-wrapper">
            <svg className="search-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <circle cx="11" cy="11" r="8"></circle>
              <path d="m21 21-4.35-4.35"></path>
            </svg>
            <input
              type="search"
              placeholder="Search shops by name or location..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="modern-search-input"
            />
            {searchTerm && (
              <button 
                type="button" 
                className="clear-search"
                onClick={() => setSearchTerm('')}
              >
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                  <line x1="18" y1="6" x2="6" y2="18"></line>
                  <line x1="6" y1="6" x2="18" y2="18"></line>
                </svg>
              </button>
            )}
          </div>
        </div>

        <div className="filter-controls">
          <div className="radius-control">
            <label htmlFor="radius-select" className="control-label">
              <svg className="control-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <circle cx="12" cy="12" r="10"></circle>
                <circle cx="12" cy="12" r="3"></circle>
              </svg>
              Range
            </label>
            <select
              id="radius-select"
              value={radius}
              onChange={(e) => setRadius(Number(e.target.value))}
              className="modern-select"
            >
              <option value={1000}>1 km</option>
              <option value={2000}>2 km</option>
              <option value={5000}>5 km</option>
              <option value={10000}>10 km</option>
              <option value={20000}>20 km</option>
            </select>
          </div>

          <button type="button" className="location-btn" onClick={handleUseLocation}>
            <svg className="location-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"></path>
              <circle cx="12" cy="10" r="3"></circle>
            </svg>
            Use My Location
          </button>

          <div className="view-toggle">
            <label className="control-label">View</label>
            <div className="toggle-group">
              <button
                type="button"
                className={`toggle-option ${viewMode === 'list' ? 'active' : ''}`}
                onClick={() => setViewMode('list')}
              >
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                  <line x1="8" y1="6" x2="21" y2="6"></line>
                  <line x1="8" y1="12" x2="21" y2="12"></line>
                  <line x1="8" y1="18" x2="21" y2="18"></line>
                  <line x1="3" y1="6" x2="3.01" y2="6"></line>
                  <line x1="3" y1="12" x2="3.01" y2="12"></line>
                  <line x1="3" y1="18" x2="3.01" y2="18"></line>
                </svg>
                List
              </button>
              <button
                type="button"
                className={`toggle-option ${viewMode === 'map' ? 'active' : ''}`}
                onClick={() => setViewMode('map')}
              >
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                  <polygon points="1 6 1 22 8 18 16 22 23 18 23 2 16 6 8 2 1 6"></polygon>
                  <line x1="8" y1="2" x2="8" y2="18"></line>
                  <line x1="16" y1="6" x2="16" y2="22"></line>
                </svg>
                Map
              </button>
            </div>
          </div>
        </div>
      </div>

      {shopsLoading ? (
        <div className="loading-state">
          <div className="loader-ring" />
          <p>Finding shops near you...</p>
        </div>
      ) : null}

      {shopsError ? <div className="alert">{shopsError}</div> : null}

      {!shopsLoading && !shopsError && filteredShops.length === 0 ? (
        <div className="empty-state">
          <h3>No shops found</h3>
          <p className="text-muted">Try expanding your search radius or check your location settings.</p>
        </div>
      ) : null}

      {viewMode === 'list' && filteredShops.length > 0 ? (
        <div className="shops-grid">
          {filteredShops.map((shop) => {
            const distanceText = formatDistance(shop.distance)
            const details = shopMeta[shop.id]
            return (
              <article key={shop.id} className="shop-card">
                <div className="shop-card-header">
                  <h3>{shop.name}</h3>
                  {distanceText ? <span className="distance-pill">{distanceText} away</span> : null}
                </div>
                <p className="shop-address">{shop.address}</p>
                {details ? (
                  <div className="shop-stats">
                    <span className="stat">
                      👥 {details.subscriber_count} subscriber{details.subscriber_count === 1 ? '' : 's'}
                    </span>
                    {details.is_subscribed ? <span className="subscribed-badge">Subscribed</span> : null}
                  </div>
                ) : null}
                <div className="shop-actions">
                  <Link to={`/shops/${shop.id}`} className="btn btn-primary">
                    View Shop
                  </Link>
                  {user && details ? (
                    <button 
                      type="button" 
                      className={`btn ${details.is_subscribed ? 'btn-unsubscribe' : 'btn-subscribe'}`}
                      onClick={() => handleSubscription(shop.id, details.is_subscribed)}
                      disabled={subscriptionLoading[shop.id]}
                    >
                      {subscriptionLoading[shop.id] ? (
                        '...'
                      ) : details.is_subscribed ? (
                        'Unsubscribe'
                      ) : (
                        'Subscribe'
                      )}
                    </button>
                  ) : null}
                  <button type="button" className="btn btn-ghost">
                    Get Directions
                  </button>
                </div>
              </article>
            )
          })}
        </div>
      ) : null}

      {viewMode === 'map' && coords ? (
        <div className="map-container">
          <ShopMap shops={filteredShops} userCoords={coords} fallbackCenter={DEFAULT_COORDS} height={500} />
          <div className="map-info">
            <p>📍 Showing {filteredShops.length} shop{filteredShops.length === 1 ? '' : 's'} within {radius / 1000}km</p>
          </div>
        </div>
      ) : null}
    </div>
  )
}
