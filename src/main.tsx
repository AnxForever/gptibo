import { StrictMode, lazy, Suspense } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import Root from './Root.tsx'
import { RechargeSitePreview } from './components/recharge-site-preview'
const SessionPrecheckPreview = import.meta.env.DEV ? lazy(() => import('./components/session-precheck-preview')) : null

// Storefront HTML stays neutral until its catalog identifies the shop.
if (location.pathname === '/admin' || location.pathname.startsWith('/admin/') || ['/join', '/join/', '/partner'].includes(location.pathname)) {
  document.title = location.pathname.startsWith('/admin') ? '管理工作台 · Anx' : '代理工作台 · Anx'
}

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    {SessionPrecheckPreview && location.pathname === '/__preview/session-precheck' ? <Suspense fallback={null}><SessionPrecheckPreview /></Suspense> : window.parent !== window && window.frameElement?.getAttribute('data-guide-scene') === 'package' ? <RechargeSitePreview /> : <Root />}
  </StrictMode>,
)
