import type { CSSProperties } from 'react'

/**
 * Upcoming products are shown with their real brand mark rather than a
 * placeholder: the stores we resell all have a recognisable glyph, and a
 * generic shape next to a real logo (Claude, Grok, OpenAI) reads as broken.
 * Muse ships a coloured logo; the rest are monochrome masks tinted per brand.
 */
const BRANDS: Record<string, { color: string; dark?: string; glow?: string; image?: boolean }> = {
  muse: { color: '#2f6bff', image: true },
  // Black marks need a neutral halo: tinting the glow with #0f0f0f leaves a
  // smudge behind the glyph instead of light.
  cursor: { color: '#0f0f0f', dark: '#f2f2f2', glow: '#7c8290' },
  perplexity: { color: '#1fb8cd' },
  canva: { color: '#00c4cc' },
  notion: { color: '#0f0f0f', dark: '#f2f2f2', glow: '#7c8290' },
}

export function UpcomingScene({ productId }: { productId: string }) {
  const brand = BRANDS[productId]
  if (!brand) return null
  return <div className="brand-art" data-brand={productId} style={{ '--brand': brand.color, '--brand-dark': brand.dark ?? brand.color, '--brand-glow': brand.glow ?? brand.color } as CSSProperties} aria-hidden="true">
    <span className="brand-art-glow" />
    {brand.image
      ? <img className="brand-art-logo" src={`/brands/${productId}.svg`} alt="" width={84} height={84} />
      : <span className="brand-art-logo" style={{ maskImage: `url(/brands/${productId}.svg)` }} />}
    <span className="brand-art-floor" />
  </div>
}
