import { isXGiftPlan } from './xgift-products.ts'
export const EMAIL_PRODUCTS = {
  'gmail-account': { productId: 'gmail', name: 'Gmail 邮箱账号', kind: 'cdk', label: '账号资料', duration: '1 个账号' },
  // The 18-month row is retired supply: it stays so historical orders keep their
  // claim-link delivery, and the shop never renders it.
  'gemini-pro-18m': { productId: 'gemini', name: 'Gemini Pro 18 个月会员', kind: 'link', label: '领取链接', duration: '18 个月', legacy: true },
  'gemini-plus-6m': { productId: 'gemini', name: 'Gemini Plus 6 个月会员', kind: 'cdk', label: '账号资料', duration: '6 个月' },
  'duolingo-12m': { productId: 'duolingo', name: '多邻国 1 年会员', kind: 'cdk', label: '兑换 CDK', duration: '1 年' },
  'elevenlabs-creator-3m': { productId: 'elevenlabs', name: 'ElevenLabs Creator 3 个月', kind: 'cdk', label: '账号资料', duration: '3 个月' },
  'muse-account': { productId: 'muse', name: 'Muse AI 高级账号', kind: 'cdk', label: '账号资料', duration: '长期' },
  // The supplier calls this an onboarding hand-off rather than an account, so
  // the label stays neutral: credentials and invite links both read correctly.
  'notion-business-3m': { productId: 'notion', name: 'Notion Business 3 个月', kind: 'cdk', label: '开通资料', duration: '3 个月' },
  // One manually delivered claim link: payment only reserves the order, staff
  // paste the link in the workspace. `manual` keeps the copy honest.
  'devin-max-1m': { productId: 'devin', name: 'Devin Max 1 个月', kind: 'link', label: '领取链接', duration: '1 个月', manual: true },
} as const
export type EmailPlan = keyof typeof EMAIL_PRODUCTS
export const EMAIL_PRODUCT_IDS: ReadonlySet<string> = new Set(Object.values(EMAIL_PRODUCTS).map(product => product.productId))
export function emailProduct(plan: string) {
  return Object.hasOwn(EMAIL_PRODUCTS, plan) ? EMAIL_PRODUCTS[plan as EmailPlan] : null
}
export function productNameForPlan(plan: string) {
  return isXGiftPlan(plan) ? 'X Premium' : emailProduct(plan)?.name ?? (plan === 'codex-sms-once' ? 'Codex 接码' : 'ChatGPT')
}
export function productNameForId(id: string) {
  return ({ 'x-premium': 'X Premium', chatgpt: 'ChatGPT', 'codex-sms': 'Codex 接码', gemini: 'Gemini Plus', duolingo: '多邻国', gmail: 'Gmail 账号', elevenlabs: 'ElevenLabs', muse: 'Muse AI', notion: 'Notion', devin: 'Devin Max' } as Record<string, string>)[id] ?? id
}
// Manual plans are pasted by staff after payment; their copy must not promise
// automatic delivery. `manual` only exists on those entries.
export function isManualDelivery(plan: string) {
  const product = emailProduct(plan)
  return product !== null && 'manual' in product && product.manual === true
}
export type DeliverySummary = { automatic?: boolean; value?: string; instructions?: string; kind: 'link' | 'cdk'; status: 'waiting' | 'draft' | 'sending' | 'unknown' | 'sent'; sentAt: number | null }
export type AdminDelivery = DeliverySummary & { value: string; instructions: string; version: number; firstAttemptAt: number | null; providerId: string | null; preview: { subject: string; html: string; text: string } | null }
