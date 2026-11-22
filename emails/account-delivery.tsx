import { emailBrand } from './brand.ts'
import * as React from 'react'
import { Body, Container, Head, Heading, Html, Img, Preview, Section, Text } from 'react-email'

// The generic account mail. Suppliers describe their payloads in different
// words — a colon-joined credential, a plain account line, an invite link — so
// this template describes the shape loosely instead of promising a format.
export default function AccountDeliveryEmail() {
  return <Html lang="zh-CN"><Head /><Preview>您的账号资料已准备好，请妥善保管。</Preview>
    <Body style={{ margin: 0, padding: '40px 16px', background: '#f5f5f7', color: '#202124', fontFamily: '-apple-system,BlinkMacSystemFont,"Segoe UI","PingFang SC",Arial,sans-serif' }}>
      <Container style={{ maxWidth: '560px' }}><Img src={emailBrand.logo} width="32" height="32" alt={emailBrand.name} style={{ marginBottom: '24px' }} />
        <Section style={{ padding: '32px 24px', background: '#fff', borderRadius: '20px', borderTop: '4px solid #6b6f7a' }}>
          <Text style={{ wordBreak: 'break-all', fontSize: '12px', color: '#72757c' }}>{emailBrand.name} · 账号交付</Text>
          <Heading as="h1" style={{ fontSize: '26px' }}>您的账号资料已发货。</Heading>
          <Text style={{ fontSize: '14px', lineHeight: '24px' }}>下方为<strong>交付原文</strong>，请整段复制，不要按标点拆开。若其中包含网址，请从 <strong>https://</strong> 开始整段复制到浏览器打开；若为冒号或竖线分隔的多段内容，各段依次可能是邮箱、邮箱密码、账号密码、2FA 密钥与验证链接。请勿转发给他人。</Text>
          <Section style={{ background: '#f4f4f6', padding: '20px', borderRadius: '12px' }}><Text style={{ fontSize: '14px', lineHeight: '24px', whiteSpace: 'pre-wrap', wordBreak: 'break-all', fontFamily: 'monospace' }}>__RECHARGE_DELIVERY_VALUE__</Text></Section>
          <Text style={{ fontSize: '14px', fontWeight: 600 }}>使用说明</Text><Text style={{ fontSize: '14px', lineHeight: '24px', whiteSpace: 'pre-line' }}>__RECHARGE_DELIVERY_INSTRUCTIONS__</Text>
          <Text style={{ fontSize: '12px', color: '#72757c', wordBreak: 'break-all' }}>订单号：__RECHARGE_DELIVERY_ORDER__</Text>
          <Text style={{ fontSize: '13px', lineHeight: '22px' }}>请尽快登录并修改账号密码，超过页面写明的时限不再受理售后。如需帮助，请回到 <a href={emailBrand.orders}>我的订单</a>，通过原订单联系支持。邮件发出表示资料已交付，实际登录结果请自行核对。</Text>
        </Section><Text style={{ wordBreak: 'break-all', fontSize: '11px', color: '#858790' }}>{emailBrand.name} · 独立服务商，此邮件仅用于商品交付，请勿直接回复。</Text>
      </Container>
    </Body></Html>
}
