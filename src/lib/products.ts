import { ORDER_PLANS, PLANS, type PlanId } from './demo-orders'

export const PRODUCTS = [
  { id: 'x-premium', name: 'X Premium', category: '社交订阅', tagline: '你的账号，更多可能。', description: '选择 3 或 6 个月 Premium，继续你的精彩。', pitch: '只需 X 用户名，无需密码。', badge: '个人账号', unit: '份', quantity: '所选订阅周期', method: '账号充值', defaultPlan: 'x-premium-3m' as PlanId },
  { id: 'chatgpt', name: 'ChatGPT', category: 'AI 订阅', tagline: '你的日常灵感搭档', description: '写作、学习与深度探索，选择适合自己的订阅。', pitch: '从一个好问题，到一个新想法。在熟悉的账号里，继续你的探索。', badge: '个人账号', unit: '月', quantity: '1 个月', method: '账号直充', defaultPlan: 'plus' as PlanId },
  { id: 'gemini', name: 'Gemini Plus', category: 'AI 订阅', tagline: '6 个月，把灵感带在身边', description: '现成 Gemini 账号，Plus 会员 6 个月，付款后自动交付账号资料。', pitch: '现成账号，付款确认后自动发货，账号邮箱与密码可在订单页查看，也会发送至邮箱。请首次登录后立即修改密码，不要转发给他人。', badge: '自动发货', unit: '6 个月', quantity: '6 个月', method: '订单页与邮件交付账号资料', defaultPlan: 'gemini-plus-6m' as PlanId },
  { id: 'duolingo', name: '多邻国', category: '语言学习', tagline: '一整年，和 Duo 一起学', description: '1 年会员，付款后自动交付兑换 CDK。', pitch: '给下一年的学习留个位置。付款确认后自动发货，兑换 CDK 与说明可在订单页查看，也会发送至邮箱。', badge: '自动发货', unit: '年', quantity: '1 年', method: '订单页与邮件领取 CDK', defaultPlan: 'duolingo-12m' as PlanId },
  { id: 'elevenlabs', name: 'ElevenLabs', category: 'AI 订阅', tagline: '把文字，变成声音', description: '现成稳定账号，Creator 会员 3 个月（131K 积分/月），付款后自动交付账号资料。', pitch: '现成稳定账号，付款确认后自动发货，邮箱、邮箱密码、Creator 密码与邮箱查看链接可在订单页查看。请在付款后 30 分钟内登录并修改密码，超时不再受理售后。', badge: '自动发货', unit: '3 个月', quantity: '3 个月', method: '订单页与邮件交付账号资料', defaultPlan: 'elevenlabs-creator-3m' as PlanId },
  { id: 'codex-sms', name: 'Codex 接码', category: '短信服务', tagline: '获取手机号，查询接码结果', description: '付款后获取手机号，再查询短信接收结果。', pitch: '付款后获取手机号，再在订单中查询短信状态与验证码。', badge: '单次服务', unit: '次', quantity: '1 次', method: '短信接码', defaultPlan: 'codex-sms-once' as PlanId },
  { id: 'muse', name: 'Muse AI', category: 'AI 订阅', tagline: '把灵感，变成作品。', description: '现成 Muse AI 高级账号，付款后自动交付账号资料。', pitch: '现成账号，付款确认后自动发货，账号资料可在订单页查看，也会发送至邮箱。请付款后 30 分钟内登录并修改密码；不支持质保，请自行核对登录结果。', badge: '自动发货', unit: '长期', quantity: '1 个账号', method: '订单页与邮件交付账号资料', defaultPlan: 'muse-account' as PlanId },
  { id: 'notion', name: 'Notion', category: '效率工具', tagline: '三个月，把想法收好。', description: 'Notion Business 3 个月，付款后自动交付开通资料。', pitch: '付款确认后自动发货，开通资料可在订单页查看，也会发送至邮箱。请按邮件正文使用，不要转发给他人。', badge: '自动发货', unit: '3 个月', quantity: '3 个月', method: '订单页与邮件交付开通资料', defaultPlan: 'notion-business-3m' as PlanId },
  { id: 'devin', name: 'Devin Max', category: 'AI 订阅', tagline: '你的 AI 工程师，替你写代码', description: 'Devin Max 1 个月，付款后人工交付专属开通链接。', pitch: '仅 1 份。付款确认后，专属开通链接会发送到下单邮箱，也可以在订单页查看。链接仅限本订单使用，请勿转发给他人。', badge: '人工发货', unit: '1 个月', quantity: '1 个月', method: '订单页与邮件交付领取链接', defaultPlan: 'devin-max-1m' as PlanId },
  { id: 'gmail', name: 'Gmail 账号', category: '邮箱账号', tagline: '美国个人白号，资料自动交付', description: '美国长效个人 Gmail 白号，带 2FA 与验证链接。', pitch: '填写可正常收件的邮箱，付款确认后自动发送账号资料。请在购买当天首次登录，并提前确认登录环境与质保范围。', badge: '个人白号', unit: '个', quantity: '1 个账号', method: '自动邮件发货', defaultPlan: 'gmail-account' as PlanId },
] as const
export type ProductId = typeof PRODUCTS[number]['id']
export const SMS_STEPS = ['选择付款', '获取手机号', '查询结果'] as const
export const UPCOMING_PRODUCTS = [
  { id: 'claude', name: 'Claude', category: 'AI 订阅', description: 'Claude 订阅服务，暂未开放购买。', tags: ['套餐待公布', '交付方式待公布'] },
  { id: 'grok', name: 'Grok', category: 'AI 订阅', description: 'Grok 订阅服务，暂未开放购买。', tags: ['套餐待公布', '交付方式待公布'] },
  { id: 'cursor', name: 'Cursor', category: '开发工具', description: 'Cursor Pro 订阅，暂未开放购买。', tags: ['套餐待公布', '交付方式待公布'] },
  { id: 'perplexity', name: 'Perplexity', category: 'AI 订阅', description: 'Perplexity Pro 订阅，暂未开放购买。', tags: ['套餐待公布', '交付方式待公布'] },
  { id: 'canva', name: 'Canva', category: '设计工具', description: 'Canva Pro 与团队席位，暂未开放购买。', tags: ['套餐待公布', '交付方式待公布'] },
] as const
export function productPlans(id: ProductId) { return PLANS.filter(plan => plan.productId === id) }
export function productForPlan(planId: PlanId) { return PRODUCTS.find(product => product.id === ORDER_PLANS.find(plan => plan.id === planId)!.productId)! }
export { readStorefrontPage as readPage } from './storefront-routing'
