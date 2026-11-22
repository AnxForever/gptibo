import { isXGiftPlan } from '../../shared/xgift-products.ts'
import { emailProduct } from '../../shared/email-products.ts'
import { confirmedRefundLabel, type PlatformOrderStatus } from '../../shared/platform-order-status.ts'

export const PLANS = [
  { id: 'x-premium-3m', productId: 'x-premium', name: 'Premium · 3 个月', price: 35, renewalOnly: false, subtitle: '三个月，探索更多。', description: '轻松开始，体验 Premium', features: ['无需账号密码', '付款前验证账号', '订单进度可查'] },
  { id: 'x-premium-6m', productId: 'x-premium', name: 'Premium · 6 个月', price: 55, renewalOnly: false, subtitle: '六个月，继续精彩。', description: '更长时间，持续探索', features: ['无需账号密码', '付款前验证账号', '订单进度可查'] },
  { id: 'plus', productId: 'chatgpt', name: 'Plus', price: 139, renewalOnly: false, subtitle: '日常灵感，持续在线。', description: '适合日常问答、写作与学习', features: ['写作、学习与灵感探索', '图片理解与创意表达', '个人账号使用体验'] },
  { id: 'pro-5x', productId: 'chatgpt', name: 'Pro 5×', price: 699, renewalOnly: false, subtitle: '选择你的 Pro 订阅。', description: 'Pro 5× 档位，按月订阅', features: ['Pro 5× 订阅', '个人账号使用', '订单进度可查'] },
  // Keep the existing ID for supplier mappings and historical orders.
  { id: 'pro-20x', productId: 'chatgpt', name: 'Pro 10×', price: 1400, renewalOnly: false, subtitle: '更多额度，继续你的探索。', description: 'Pro 10× 档位，支持新开与续费', features: ['Pro 10× 订阅', '支持新开与续费', '订单进度可查'] },
  { id: 'pro-25x', productId: 'chatgpt', name: 'Pro 25×', price: 4500, renewalOnly: false, subtitle: '更高额度，继续你的探索。', description: 'Pro 25× 档位，支持新开与自动充值', features: ['Pro 25× 订阅', '支持新开 · 自动充值', '订单进度可查'] },
  { id: 'codex-sms-once', productId: 'codex-sms', name: '单次接码', price: 10, renewalOnly: false, subtitle: '一次接码，清楚明了。', description: '接收一次短信验证码', features: ['按次购买', '订单内查看验证码', '无需提交 Session'] },
  { id: 'gemini-plus-6m', productId: 'gemini', name: '6 个月会员', price: 13.9, renewalOnly: false, subtitle: '让灵感，多一点时间。', description: '6 个月会员 · 账号资料邮件交付', features: ['6 个月会员', '账号资料邮件交付', '自动发货 · 订单页可查看'] },
  { id: 'duolingo-12m', productId: 'duolingo', name: '1 年会员', price: 19.9, renewalOnly: false, subtitle: '和 Duo 一起，把学习继续。', description: '1 年会员 · 兑换 CDK 邮件交付', features: ['1 年会员', '邮件兑换 CDK', '自动发货 · 订单页可领取'] },
  { id: 'elevenlabs-creator-3m', productId: 'elevenlabs', name: 'Creator · 3 个月', price: 59, renewalOnly: false, subtitle: '把文字，变成声音。', description: 'Creator 会员 3 个月 · 现成稳定账号，账号资料邮件交付', features: ['Creator 会员 3 个月 · 131K 积分/月', '现成稳定账号 · 邮件交付', '付款后 30 分钟内改密 · 不支持质保'] },
  { id: 'muse-account', productId: 'muse', name: '高级账号', price: 15.9, renewalOnly: false, subtitle: '把灵感，变成作品。', description: 'Muse AI 高级账号 · 账号资料邮件交付', features: ['现成高级账号', '账号资料邮件交付', '付款后 30 分钟内改密 · 不支持质保'] },
  { id: 'notion-business-3m', productId: 'notion', name: 'Business · 3 个月', price: 13.9, renewalOnly: false, subtitle: '三个月，把想法收好。', description: 'Notion Business 3 个月 · 开通资料邮件交付', features: ['Business 工作空间 3 个月', '开通资料邮件交付', '自动发货 · 订单页可查看'] },
  { id: 'devin-max-1m', productId: 'devin', name: '1 个月', price: 240, renewalOnly: false, subtitle: '你的 AI 工程师，替你写代码。', description: 'Devin Max 1 个月 · 专属开通链接人工交付', features: ['Devin Max 会员 1 个月', '专属开通链接 · 人工发货', '仅 1 份 · 不支持售后质保'] },
  { id: 'gmail-account', productId: 'gmail', name: '美国长效白号 · 带 2FA', price: 18.8, renewalOnly: false, subtitle: '一个邮箱，让连接开始。', description: '个人 Gmail 白号 · 账号资料自动邮件交付', features: ['美国个人白号 · 带 2FA', '付款后自动邮件发货', '购买当天首次登录质保'] },
] as const
// Retired plans remain readable without silently changing historical purchases.
export const ORDER_PLANS = [...PLANS, { id: 'pro', productId: 'chatgpt', name: 'Pro', price: 800, renewalOnly: false }] as const
export type PlanId = typeof ORDER_PLANS[number]['id']
export type OrderStatus = 'pending' | 'session' | 'processing' | 'success' | 'failed' | 'cancelled'
export type DemoOrder = { id: string; plan: PlanId; price: number; createdAt: number; status: OrderStatus; startedAt?: number; outcome?: 'success' | 'failed'; smsPhone?: string; customerEmail?: string | null; refundedCents?: number; platformStatus?: PlatformOrderStatus | null; refundHold?: boolean; redemptionBlocked?: boolean }
export type OrderAction = 'pay' | 'acquire-phone' | 'submit-success' | 'submit-failed' | 'settle' | 'retry' | 'cancel'
export const ORDER_LABELS: Record<OrderStatus, string> = { pending: '待付款', session: '待提交凭据', processing: '充值中', success: '已完成', failed: '需重新提交', cancelled: '已关闭' }
export function orderLabel(order: DemoOrder & { manualFulfillment?: boolean; fulfillment?: { paymentExpired?: boolean; rechargeStatus?: string | null } }): string {
  const refundLabel = confirmedRefundLabel(Math.max(order.refundedCents ?? 0, order.platformStatus?.refundedCents ?? 0), Math.round(order.price * 100))
  if (refundLabel) return refundLabel
  if (order.refundHold) return '售后处理中'
  if (order.redemptionBlocked && order.status !== 'success') return '已禁止兑换'
  if (['session', 'processing'].includes(order.status) && order.fulfillment?.rechargeStatus === 'unknown') return '提交待确认'
  if (order.status === 'cancelled' && order.fulfillment?.paymentExpired) return '已过期'
  if (isXGiftPlan(order.plan)) return ({ pending: '待付款', session: '等待充值', processing: '充值处理中', success: '已完成', failed: '充值未完成', cancelled: '已关闭' })[order.status]
  if (emailProduct(order.plan)) return ({ pending: '待付款', session: '待邮件发货', processing: '发货处理中', success: '已发货', failed: '发货待核对', cancelled: '已关闭' })[order.status]
  if (order.manualFulfillment && order.status === 'session') return '待人工开通'
  if (order.manualFulfillment && order.status === 'processing') return '人工开通中'
  if (order.plan !== 'codex-sms-once') return ORDER_LABELS[order.status]
  return { pending: '待付款', session: order.smsPhone ? '待查询结果' : '待获取手机号', processing: '查询中', success: '已收到验证码', failed: '接码超时', cancelled: '已关闭' }[order.status]
}
export const STORAGE_KEY = 'ymx-demo-orders-v1'
export const DEMO_SESSION = 'RECHARGE_TEST_SESSION_NOT_A_REAL_TOKEN'
export const DEMO_SMS_PHONE = '+1 202-555-0147'
export const PROCESS_MS = 4500
export function makeOrder(plan: PlanId): DemoOrder {
  const currentPlan = PLANS.find(p => p.id === plan)
  if (!currentPlan) throw new Error('This plan is no longer available for new orders')
  return { id: `YMX-${crypto.randomUUID().slice(0, 8).toUpperCase()}`, plan, price: currentPlan.price, createdAt: Date.now(), status: 'pending' }
}
export function transitionOrder(order: DemoOrder, action: OrderAction, now = Date.now()): DemoOrder {
  if (action === 'pay' && order.status === 'pending') return { ...order, status: 'session' }
  if (action === 'acquire-phone' && order.plan === 'codex-sms-once' && order.status === 'session' && !order.smsPhone) return { ...order, smsPhone: DEMO_SMS_PHONE }
  if (action.startsWith('submit-') && order.status === 'session' && (order.plan !== 'codex-sms-once' || order.smsPhone === DEMO_SMS_PHONE)) return { ...order, status: 'processing', startedAt: now, outcome: action === 'submit-success' ? 'success' : 'failed' }
  if (action === 'settle' && order.status === 'processing' && order.startedAt !== undefined && now - order.startedAt >= PROCESS_MS) return { ...order, status: order.outcome === 'failed' ? 'failed' : 'success' }
  if (action === 'retry' && order.status === 'failed') return { ...order, status: 'session', startedAt: undefined, outcome: undefined }
  if (action === 'cancel' && order.status === 'pending') return { ...order, status: 'cancelled' }
  return order
}
export function readOrders(): DemoOrder[] {
  try {
    const data: unknown = JSON.parse(localStorage.getItem(STORAGE_KEY) ?? '[]')
    if (!Array.isArray(data)) return []
    return data.filter((o): o is DemoOrder => !!o && typeof o === 'object' && /^YMX-[A-F0-9]{8}$/.test(o.id) && ORDER_PLANS.some(p => p.id === o.plan) && Number.isFinite(o.price) && Number.isFinite(o.createdAt) && Object.hasOwn(ORDER_LABELS, o.status) && (o.smsPhone === undefined || (o.plan === 'codex-sms-once' && o.smsPhone === DEMO_SMS_PHONE)) && (o.status !== 'processing' || (Number.isFinite(o.startedAt) && ['success', 'failed'].includes(o.outcome)))).slice(0, 30).map(o => {
      // Older SMS demos skipped number acquisition. Keep their completed/running
      // history usable with the same clearly labeled fixture, without recharging.
      const restored = o.plan === 'codex-sms-once' && !o.smsPhone && ['processing', 'success', 'failed'].includes(o.status) ? { ...o, smsPhone: DEMO_SMS_PHONE } : o
      return transitionOrder(restored, 'settle')
    })
  } catch { return [] }
}
