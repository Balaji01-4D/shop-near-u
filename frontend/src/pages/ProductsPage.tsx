import { useCallback, useEffect, useMemo, useState } from 'react'
import { suggestCatalogProducts } from '../services/catalogApi'
import type { CatalogProduct } from '../types/catalog'

function formatCategory(label: string) {
  if (!label) {
    return ''
  }
  return label
    .toLowerCase()
    .split(' ')
    .map((segment) => segment.charAt(0).toUpperCase() + segment.slice(1))
    .join(' ')
}

export function ProductsPage() {
  const [products, setProducts] = useState<CatalogProduct[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [searchTerm, setSearchTerm] = useState('')
  const [selectedCategory, setSelectedCategory] = useState('All categories')

  const loadProducts = useCallback(async (keyword = '', limit = 50) => {
    setLoading(true)
    setError(null)
    try {
      const results = await suggestCatalogProducts(keyword, limit)
      setProducts(results)
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Unable to load products right now.'
      setError(message)
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    void loadProducts()
  }, [loadProducts])

  const handleSearch = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    const keyword = searchTerm.trim()
    await loadProducts(keyword, 50)
  }

  const categories = useMemo(() => {
    const categorySet = new Set<string>()
    products.forEach((product) => {
      if (product.category) {
        categorySet.add(formatCategory(product.category))
      }
    })
    const dynamicCategories = Array.from(categorySet).sort((a, b) => a.localeCompare(b))
    return ['All categories', ...dynamicCategories]
  }, [products])

  const filteredProducts = useMemo(() => {
    if (selectedCategory === 'All categories') {
      return products
    }
    return products.filter((product) => formatCategory(product.category) === selectedCategory)
  }, [products, selectedCategory])

  return (
    <div className="products-page">
      <div className="search-section">
        <form onSubmit={handleSearch} className="search-form">
          <div className="search-inputs">
            <input
              type="search"
              placeholder="Search products by name, brand, or description..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="search-input"
            />
            <select
              value={selectedCategory}
              onChange={(e) => setSelectedCategory(e.target.value)}
              className="category-select"
            >
              {categories.map((category) => (
                <option key={category} value={category}>
                  {category}
                </option>
              ))}
            </select>
          </div>
          <button type="submit" className="btn btn-primary" disabled={loading}>
            {loading ? 'Searching...' : 'Search'}
          </button>
        </form>
      </div>

      <div className="results-info">
        <p className="text-muted">
          {loading
            ? 'Loading products...'
            : `Showing ${filteredProducts.length} product${filteredProducts.length === 1 ? '' : 's'}`}
          {selectedCategory !== 'All categories' ? ` in ${selectedCategory}` : ''}
        </p>
      </div>

      {error ? <div className="alert">{error}</div> : null}

      {!loading && !error && filteredProducts.length === 0 ? (
        <div className="empty-state">
          <h3>No products found</h3>
          <p className="text-muted">Try adjusting your search terms or browse different categories.</p>
        </div>
      ) : null}

      {loading ? (
        <div className="loading-state">
          <div className="loader-ring" />
          <p>Finding products...</p>
        </div>
      ) : (
        <div className="products-grid modern-cards">
          {filteredProducts.map((product) => (
            <article key={product.id} className="product-card-modern clean">
              <div className="product-image">
                {product.image_url ? (
                  <img src={product.image_url} alt={product.name} />
                ) : (
                  <div className="image-placeholder">
                    <span>📦</span>
                  </div>
                )}
              </div>
              
              <div className="product-content">
                <h3 className="product-name">{product.name}</h3>
                
                <div className="product-tags">
                  <span className="product-tag category">{formatCategory(product.category) || 'General'}</span>
                  {product.brand && <span className="product-tag brand">{product.brand}</span>}
                </div>

                {product.description && (
                  <p className="product-description">{product.description}</p>
                )}

                <div className="product-actions">
                  <button type="button" className="btn btn-primary">
                    Find in Shops
                  </button>
                  <button type="button" className="btn btn-ghost">
                    Save for Later
                  </button>
                </div>
              </div>
            </article>
          ))}
        </div>
      )}
    </div>
  )
}
