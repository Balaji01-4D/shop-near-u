import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'
import { fetchUserSubscriptions, fetchShopProducts, unsubscribeFromShop } from '../services/shopApi'
import { updatePageTitle, pageTitles } from '../utils/pageTitle'
import type { SubscribedShop, ShopProduct } from '../types/shop'

export function SubscriptionsPage() {
  const { user } = useAuth()
  const [subscribedShops, setSubscribedShops] = useState<SubscribedShop[]>([])
  const [products, setProducts] = useState<ShopProduct[]>([])
  const [selectedShopId, setSelectedShopId] = useState<number | 'all'>('all')
  const [loading, setLoading] = useState(true)
  const [productsLoading, setProductsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  // Set page title
  useEffect(() => {
    updatePageTitle(pageTitles.subscriptions)
  }, [])

  useEffect(() => {
    if (!user) {
      setLoading(false)
      return
    }

    const loadSubscriptions = async () => {
      try {
        setLoading(true)
        setError(null)

        // Get user's subscribed shops directly from the API
        const subscribedShopsData = await fetchUserSubscriptions()
        setSubscribedShops(subscribedShopsData)

      } catch (error) {
        const message = error instanceof Error ? error.message : 'Unable to load your subscriptions.'
        setError(message)
      } finally {
        setLoading(false)
      }
    }

    void loadSubscriptions()
  }, [user])

  // Load products when shop selection changes
  useEffect(() => {
    if (!user || subscribedShops.length === 0) return

    const loadProducts = async () => {
      try {
        setProductsLoading(true)
        const allProducts: ShopProduct[] = []

        if (selectedShopId === 'all') {
          // Load products from all subscribed shops
          await Promise.all(
            subscribedShops.map(async (shop) => {
              try {
                const shopProducts = await fetchShopProducts(shop.id)
                const productsWithShopInfo = shopProducts.map(product => ({
                  ...product,
                  shop: {
                    id: shop.id,
                    name: shop.name
                  }
                }))
                allProducts.push(...productsWithShopInfo)
              } catch (error) {
                console.warn(`Failed to get products for shop ${shop.id}:`, error)
              }
            })
          )
        } else {
          // Load products from selected shop only
          const selectedShop = subscribedShops.find(shop => shop.id === selectedShopId)
          if (selectedShop) {
            try {
              const shopProducts = await fetchShopProducts(selectedShop.id)
              const productsWithShopInfo = shopProducts.map(product => ({
                ...product,
                shop: {
                  id: selectedShop.id,
                  name: selectedShop.name
                }
              }))
              allProducts.push(...productsWithShopInfo)
            } catch (error) {
              console.warn(`Failed to get products for shop ${selectedShop.id}:`, error)
            }
          }
        }

        // Sort products by creation date (newest first)
        allProducts.sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime())
        setProducts(allProducts)

      } catch (error) {
        console.warn('Failed to load products:', error)
      } finally {
        setProductsLoading(false)
      }
    }

    void loadProducts()
  }, [user, subscribedShops, selectedShopId])

  const handleUnsubscribe = async (shopId: number) => {
    try {
      await unsubscribeFromShop(shopId)
      // Remove the shop from the list
      setSubscribedShops(prev => prev.filter(shop => shop.id !== shopId))
      // If the unsubscribed shop was selected, switch to 'all'
      if (selectedShopId === shopId) {
        setSelectedShopId('all')
      }
    } catch (error) {
      alert('Failed to unsubscribe from shop')
    }
  }

  const getSelectedShopName = () => {
    if (selectedShopId === 'all') return 'All Subscriptions'
    const shop = subscribedShops.find(s => s.id === selectedShopId)
    return shop ? shop.name : 'Unknown Shop'
  }

  if (!user) {
    return (
      <div className="subscriptions-page">
        <div className="page-header">
          <h1>Your Subscriptions</h1>
          <p className="text-muted">Subscribe to your favorite local shops to get updates on new products and offers.</p>
        </div>

        <div className="auth-prompt">
          <div className="auth-card">
            <h2>Sign in to view subscriptions</h2>
            <p>Create an account or sign in to subscribe to shops and get personalized updates.</p>
            <div className="auth-actions">
              <Link to="/login" className="btn btn-primary">
                Sign In
              </Link>
              <Link to="/register" className="btn btn-ghost">
                Create Account
              </Link>
            </div>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="subscriptions-page">
      <div className="page-header">
        <div>
          <h1>Your Subscriptions</h1>
          <p className="text-muted">
            You're subscribed to {subscribedShops.length} shop{subscribedShops.length === 1 ? '' : 's'}. 
            Select a shop to view their latest products.
          </p>
        </div>
        <Link to="/shops" className="btn btn-primary">
          Discover More Shops
        </Link>
      </div>

      {loading ? (
        <div className="loading-state">
          <div className="loader-ring" />
          <p>Loading your subscriptions...</p>
        </div>
      ) : null}

      {error ? <div className="alert">{error}</div> : null}

      {!loading && !error && subscribedShops.length === 0 ? (
        <div className="empty-state">
          <div className="empty-icon">📋</div>
          <h3>No subscriptions yet</h3>
          <p className="text-muted">
            Start subscribing to local shops to see their latest products and offers here. 
            It's like a personalized feed for your neighborhood shopping.
          </p>
          <Link to="/shops" className="btn btn-primary">
            Browse Nearby Shops
          </Link>
        </div>
      ) : null}

      {subscribedShops.length > 0 ? (
        <>
          {/* YouTube-style subscription bar */}
          <div className="subscription-bar">
            <div className="subscription-tabs">
              <button
                type="button"
                className={`subscription-tab ${selectedShopId === 'all' ? 'active' : ''}`}
                onClick={() => setSelectedShopId('all')}
              >
                <div className="tab-icon">🏪</div>
                <span>All Shops</span>
                <div className="tab-count">{subscribedShops.length}</div>
              </button>
              
              {subscribedShops.map((shop) => (
                <button
                  key={shop.id}
                  type="button"
                  className={`subscription-tab ${selectedShopId === shop.id ? 'active' : ''}`}
                  onClick={() => setSelectedShopId(shop.id)}
                >
                  <div className="tab-avatar">
                    {shop.name.charAt(0).toUpperCase()}
                  </div>
                  <span>{shop.name}</span>
                  <div className="tab-count">{shop.subscriber_count}</div>
                  <button
                    type="button"
                    className="unsubscribe-btn"
                    onClick={(e) => {
                      e.stopPropagation()
                      handleUnsubscribe(shop.id)
                    }}
                    title="Unsubscribe"
                  >
                    ×
                  </button>
                </button>
              ))}
            </div>
          </div>

          {/* Content area */}
          <div className="subscription-content">
            <div className="content-header">
              <h2>{getSelectedShopName()}</h2>
              <p className="text-muted">
                {selectedShopId === 'all' 
                  ? `Latest products from all ${subscribedShops.length} subscribed shops`
                  : `Latest products from ${getSelectedShopName()}`
                }
              </p>
            </div>

            {productsLoading ? (
              <div className="loading-state">
                <div className="loader-ring" />
                <p>Loading products...</p>
              </div>
            ) : null}

            {!productsLoading && products.length > 0 ? (
              <div className="products-grid modern-cards">
                {products.map((product) => {
                  const hasDiscount = product.discount > 0
                  
                  return (
                    <article key={`${product.shop_id}-${product.id}`} className="product-card-modern">
                      <div className="product-image-container">
                        <img 
                          src={product.catalog_product.image_url || '/api/placeholder/200/200'} 
                          alt={product.catalog_product.name}
                          className="product-image"
                          onError={(e) => {
                            const target = e.target as HTMLImageElement
                            target.src = 'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMjAwIiBoZWlnaHQ9IjIwMCIgdmlld0JveD0iMCAwIDIwMCAyMDAiIGZpbGw9Im5vbmUiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+PHJlY3Qgd2lkdGg9IjIwMCIgaGVpZ2h0PSIyMDAiIGZpbGw9IiNmMWY1ZjkiLz48cGF0aCBkPSJNMTAwIDEwMGMxMS4wNDYgMCAyMC00LjQ3NyAyMC0xMHMtOC45NTQtMTAtMjAtMTAtMjAgNC40NzctMjAgMTAgOC45NTQgMTAgMjAgMTB6IiBmaWxsPSIjY2JkNWUxIi8+PC9zdmc+'
                          }}
                        />
                        {hasDiscount && (
                          <div className="discount-badge-modern">
                            {product.discount}% OFF
                          </div>
                        )}
                        <div className="product-actions">
                          <button className="action-btn edit-btn" title="View Details">
                            <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
                              <path d="M12 4.5C7 4.5 2.73 7.61 1 12c1.73 4.39 6 7.5 11 7.5s9.27-3.11 11-7.5c-1.73-4.39-6-7.5-11-7.5zM12 17c-2.76 0-5-2.24-5-5s2.24-5 5-5 5 2.24 5 5-2.24 5-5 5zm0-8c-1.66 0-3 1.34-3 3s1.34 3 3 3 3-1.34 3-3-1.34-3-3-3z"/>
                            </svg>
                          </button>
                          <button className="action-btn cart-btn" title="Add to Cart">
                            <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
                              <path d="M7 18c-1.1 0-2 .9-2 2s.9 2 2 2 2-.9 2-2-.9-2-2-2zM1 2v2h2l3.6 7.59-1.35 2.45c-.16.28-.25.61-.25.96 0 1.1.9 2 2 2h12v-2H7.42c-.14 0-.25-.11-.25-.25l.03-.12L8.1 13h7.45c.75 0 1.41-.41 1.75-1.03L21.7 4H5.21l-.94-2H1zm16 16c-1.1 0-2 .9-2 2s.9 2 2 2 2-.9 2-2-.9-2-2-2z"/>
                            </svg>
                          </button>
                        </div>
                      </div>
                      
                      <div className="product-header">
                        <h3 className="product-title">{product.catalog_product.name}</h3>
                        <div className="product-meta">
                          <span className="product-category-tag">{product.catalog_product.category}</span>
                          <span className="product-brand-tag">{product.catalog_product.brand}</span>
                        </div>
                      </div>

                      <div className="product-pricing">
                        <div className="pricing-row">
                          <span className="label">Price</span>
                          <span className="value">₹{product.price.toFixed(2)}</span>
                        </div>
                        <div className="pricing-row">
                          <span className="label">Stock</span>
                          <span className="value stock-value">{product.stock} units</span>
                        </div>
                        <div className="pricing-row">
                          <span className="label">Status</span>
                          <span className={`value status-badge ${product.is_available ? 'available' : 'unavailable'}`}>
                            {product.is_available ? 'AVAILABLE' : 'OUT OF STOCK'}
                          </span>
                        </div>
                      </div>

                      {selectedShopId === 'all' && (
                        <div className="product-footer">
                          <span className="shop-attribution">From {product.shop?.name}</span>
                        </div>
                      )}
                    </article>
                  )
                })}
              </div>
            ) : !productsLoading && products.length === 0 && subscribedShops.length > 0 ? (
              <div className="empty-products">
                <div className="empty-icon">📦</div>
                <h3>No products available</h3>
                <p className="text-muted">
                  {selectedShopId === 'all'
                    ? "Your subscribed shops haven't added any products yet. Check back later for updates!"
                    : `${getSelectedShopName()} hasn't added any products yet. Check back later!`
                  }
                </p>
              </div>
            ) : null}
          </div>
        </>
      ) : null}
    </div>
  )
}
