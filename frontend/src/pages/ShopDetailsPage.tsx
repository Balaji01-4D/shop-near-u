import { useEffect, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import { fetchShopDetails, subscribeToShop, unsubscribeFromShop, fetchShopProducts } from '../services/shopApi'
import { updatePageTitle } from '../utils/pageTitle'
import type { ShopDetails, ShopProduct } from '../types/shop'
import { useAuth } from '../context/AuthContext'

export function ShopDetailsPage() {
  const { id } = useParams<{ id: string }>()
  const { user } = useAuth()
  const [shop, setShop] = useState<ShopDetails | null>(null)
  const [products, setProducts] = useState<ShopProduct[]>([])
  const [loading, setLoading] = useState(true)
  const [productsLoading, setProductsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [subscriptionLoading, setSubscriptionLoading] = useState(false)

  useEffect(() => {
    if (!id) {
      setError('Shop ID is required')
      setLoading(false)
      return
    }

    const loadShopDetails = async () => {
      try {
        setLoading(true)
        setError(null)
        const shopData = await fetchShopDetails(parseInt(id, 10))
        setShop(shopData)
        
        // Update page title with shop name
        updatePageTitle(shopData.name)
        
        // Load products for this shop
        setProductsLoading(true)
        try {
          const shopProducts = await fetchShopProducts(parseInt(id, 10))
          setProducts(shopProducts)
        } catch (productsError) {
          console.warn('Failed to load shop products:', productsError)
          // Don't show error for products, just leave empty
        } finally {
          setProductsLoading(false)
        }
        
      } catch (err) {
        const message = err instanceof Error ? err.message : 'Failed to load shop details'
        setError(message)
      } finally {
        setLoading(false)
      }
    }

    void loadShopDetails()
  }, [id])

  const handleSubscribe = async () => {
    if (!shop || !id) return

    try {
      setSubscriptionLoading(true)
      if (shop.is_subscribed) {
        await unsubscribeFromShop(parseInt(id, 10))
        setShop(prev => prev ? {
          ...prev,
          is_subscribed: false,
          subscriber_count: Math.max(0, prev.subscriber_count - 1)
        } : null)
      } else {
        await subscribeToShop(parseInt(id, 10))
        setShop(prev => prev ? {
          ...prev,
          is_subscribed: true,
          subscriber_count: prev.subscriber_count + 1
        } : null)
      }
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to update subscription'
      // You could add a toast notification here
      console.error('Subscription error:', message)
    } finally {
      setSubscriptionLoading(false)
    }
  }

  if (!user) {
    return (
      <div className="shop-details-page">
        <div className="section-block">
          <div className="card">
            <h2>Authentication Required</h2>
            <p>Please log in to view shop details.</p>
            <div style={{ display: 'flex', gap: '1rem', marginTop: '1.5rem' }}>
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

  if (loading) {
    return (
      <div className="shop-details-page">
        <div className="section-block">
          <div className="card">
            <div className="loading-state">
              <div className="spinner"></div>
              <p>Loading shop details...</p>
            </div>
          </div>
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="shop-details-page">
        <div className="section-block">
          <div className="card">
            <div className="error-state">
              <h2>Unable to Load Shop</h2>
              <p>{error}</p>
              <div style={{ display: 'flex', gap: '1rem', marginTop: '1.5rem' }}>
                <button 
                  type="button" 
                  className="btn btn-primary"
                  onClick={() => window.location.reload()}
                >
                  Try Again
                </button>
                <Link to="/shops" className="btn btn-ghost">
                  Back to Shops
                </Link>
              </div>
            </div>
          </div>
        </div>
      </div>
    )
  }

  if (!shop) {
    return (
      <div className="shop-details-page">
        <div className="section-block">
          <div className="card">
            <div className="error-state">
              <h2>Shop Not Found</h2>
              <p>The shop you're looking for doesn't exist or has been removed.</p>
              <Link to="/shops" className="btn btn-primary">
                Back to Shops
              </Link>
            </div>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="shop-details-page">
      <div className="section-block">
        <div className="breadcrumb">
          <Link to="/shops" className="breadcrumb-link">
            Shops
          </Link>
          <span className="breadcrumb-separator">›</span>
          <span className="breadcrumb-current">{shop.name}</span>
        </div>

        <div className="shop-details-card">
          <div className="shop-header">
            <div className="shop-info">
              <h1 className="shop-name">{shop.name}</h1>
              <div className="shop-stats">
                <span className="stat">
                  👥 {shop.subscriber_count} subscriber{shop.subscriber_count === 1 ? '' : 's'}
                </span>
                {shop.is_subscribed && (
                  <span className="subscribed-badge">Subscribed</span>
                )}
              </div>
            </div>
            
            <div className="shop-actions">
              {shop.is_subscribed ? (
                <button 
                  type="button" 
                  className="btn btn-ghost"
                  onClick={handleSubscribe}
                  disabled={subscriptionLoading}
                >
                  {subscriptionLoading ? 'Updating...' : 'Unsubscribe from Shop'}
                </button>
              ) : (
                <button 
                  type="button" 
                  className="btn btn-primary"
                  onClick={handleSubscribe}
                  disabled={subscriptionLoading}
                >
                  {subscriptionLoading ? 'Updating...' : 'Subscribe to Shop'}
                </button>
              )}
              <button type="button" className="btn btn-ghost">
                Get Directions
              </button>
            </div>
          </div>

          <div className="shop-content">
            <div className="section">
              <h2>About This Shop</h2>
              <p>Welcome to {shop.name}! This shop is part of our local network, offering quality products and services to the community.</p>
            </div>

            <div className="section">
              <h2>Shop Information</h2>
              <div className="info-grid">
                <div className="info-item">
                  <strong>Shop ID:</strong>
                  <span>#{shop.id}</span>
                </div>
                <div className="info-item">
                  <strong>Community:</strong>
                  <span>{shop.subscriber_count} subscriber{shop.subscriber_count === 1 ? '' : 's'}</span>
                </div>
                <div className="info-item">
                  <strong>Status:</strong>
                  <span>Active</span>
                </div>
              </div>
            </div>

            <div className="section">
              <h2>Products</h2>
              {productsLoading ? (
                <div className="loading-state">
                  <div className="loader-ring" />
                  <p>Loading products...</p>
                </div>
              ) : products.length > 0 ? (
                <div className="products-grid modern-cards">
                  {products.map((product) => {
                    const discountedPrice = product.price * (1 - product.discount / 100)
                    const hasDiscount = product.discount > 0
                    
                    return (
                      <article key={product.id} className="product-card-modern clean">
                        <div className="product-image">
                          <img 
                            src={product.catalog_product.image_url || '/api/placeholder/200/200'} 
                            alt={product.catalog_product.name}
                            onError={(e) => {
                              const target = e.target as HTMLImageElement
                              target.src = 'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMjAwIiBoZWlnaHQ9IjIwMCIgdmlld0JveD0iMCAwIDIwMCAyMDAiIGZpbGw9Im5vbmUiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+PHJlY3Qgd2lkdGg9IjIwMCIgaGVpZ2h0PSIyMDAiIGZpbGw9IiNmMWY1ZjkiLz48cGF0aCBkPSJNMTAwIDEwMGMxMS4wNDYgMCAyMC00LjQ3NyAyMC0xMHMtOC45NTQtMTAtMjAtMTAtMjAgNC40NzctMjAgMTAgOC45NTQgMTAgMjAgMTB6IiBmaWxsPSIjY2JkNWUxIi8+PC9zdmc+'
                            }}
                          />
                          {hasDiscount && (
                            <span className="discount-badge">
                              {product.discount}% OFF
                            </span>
                          )}
                        </div>

                        <div className="product-content">
                          <h3 className="product-name">{product.catalog_product.name}</h3>
                          
                          <div className="product-tags">
                            <span className="product-tag brand">{product.catalog_product.brand}</span>
                            <span className="product-tag category">{product.catalog_product.category}</span>
                          </div>

                          <div className="product-details">
                            <div className="detail-row">
                              <span className="detail-label">Price</span>
                              <div className="detail-value price">
                                {hasDiscount ? (
                                  <>
                                    <span className="current-price">₹{discountedPrice.toFixed(2)}</span>
                                    <span className="original-price">₹{product.price.toFixed(2)}</span>
                                  </>
                                ) : (
                                  <span className="current-price">₹{product.price.toFixed(2)}</span>
                                )}
                              </div>
                            </div>
                            
                            <div className="detail-row">
                              <span className="detail-label">Stock</span>
                              <span className="detail-value">{product.stock} units</span>
                            </div>
                            
                            <div className="detail-row">
                              <span className="detail-label">Status</span>
                              <span className={`detail-value status ${product.is_available && product.stock > 0 ? 'available' : 'unavailable'}`}>
                                {product.is_available && product.stock > 0 ? 'AVAILABLE' : 'OUT OF STOCK'}
                              </span>
                            </div>
                          </div>

                          <div className="product-actions">
                            <button 
                              type="button" 
                              className="btn btn-primary"
                              disabled={!product.is_available || product.stock === 0}
                            >
                              {product.is_available && product.stock > 0 ? 'Add to Cart' : 'Unavailable'}
                            </button>
                            <button type="button" className="btn btn-ghost">
                              View Details
                            </button>
                          </div>
                        </div>
                      </article>
                    )
                  })}
                </div>
              ) : (
                <div className="empty-products">
                  <div className="empty-icon">📦</div>
                  <h3>No products available</h3>
                  <p className="text-muted">
                    This shop hasn't added any products yet. Check back later for updates!
                  </p>
                </div>
              )}
            </div>

            <div className="section">
              <h2>Actions</h2>
              <div className="action-buttons">
                <Link to="/products" className="btn btn-ghost">
                  Browse All Products
                </Link>
                <Link to="/shops" className="btn btn-ghost">
                  Find Other Shops
                </Link>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}