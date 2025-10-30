import type { ReactNode } from 'react'
import { Header } from './Header'

export function AppLayout({ children }: { children: ReactNode }) {
  return (
    <div className="page-shell">
      <Header />
      <main className="page-content">{children}</main>
      <footer className="footer">
        <div className="footer-inner">
          <span>© {new Date().getFullYear()} ShopNearU. All rights reserved.</span>
          <span>
            Crafted for neighborhood discovery • <a href="https://github.com/balaji01-4d" target="_blank" rel="noreferrer">GitHub</a>
          </span>
        </div>
      </footer>
    </div>
  )
}
