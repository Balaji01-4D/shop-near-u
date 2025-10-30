import { Link } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'
import { useEffect } from 'react'
import { setDefaultTitle } from '../utils/pageTitle'
import homePageBg from '../assets/home_page_bg.mp4'

export function HomePage() {
  const { user } = useAuth()

  useEffect(() => {
    setDefaultTitle()
  }, [])

  return (
    <div className="home-page">
      {/* Full-screen video hero section */}
      <section className="hero-video-section">
        <video 
          className="hero-video" 
          autoPlay 
          loop 
          muted 
          playsInline
          preload="metadata"
        >
          <source src={homePageBg} type="video/mp4" />
          Your browser does not support the video tag.
        </video>
        
        <div className="hero-overlay"></div>
        
        <div className="hero-content">
          <div className="hero-eyebrow">Neighborhood shopping, reimagined</div>
          <h1 className="hero-title">Your Neighbourhood. Your People. Your Trust.</h1>
          <p className="hero-subtext">Find everyday essentials, groceries, and products from trusted shops around your area.</p>

          {user ? (
            <div className="cta-section">
              <div className="cta-buttons">
                <Link to="/shops" className="btn btn-primary btn-large">
                  Explore Shops
                </Link>
                <Link to="/products" className="btn btn-secondary">
                  Browse Products
                </Link>
              </div>
            </div>
          ) : (
            <div className="cta-section">
              <div className="cta-buttons">
                <Link to="/register" className="btn btn-primary btn-large">
                  Create your account
                </Link>
                <Link to="/login" className="btn btn-secondary">
                  I already have an account
                </Link>
              </div>
            </div>
          )}
        </div>
      </section>
    </div>
  )
}