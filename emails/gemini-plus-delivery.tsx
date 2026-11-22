import { emailBrand } from './brand.ts'
import * as React from 'react'
import { Body, Container, Head, Heading, Html, Img, Preview, Section, Text } from 'react-email'

export default function GeminiPlusDeliveryEmail() {
  return <Html lang="zh-CN"><Head /><Preview>您的 Gemini Plus 账号资料已准备好，请妥善保管。</Preview>
    <Body style={{ margin: 0, padding: '40px 16px', background: '#f5f5f7', color: '#202124', fontFamily: '-apple-system,BlinkMacSystemFont,"Segoe UI","PingFang SC",Arial,sans-serif' }}>
      <Container style={{ maxWidth: '560px' }}><Img src={emailBrand.logo} width="32" height="32" alt={emailBrand.name} style={{ marginBottom: '24px' }} />
        <Section style={{ padding: '32px 24px', background: '#fff', borderRadius: '20px', borderTop: '4px solid #4888ee' }}>
          <Text style={{ wordBreak: 'break-all', fontSize: '12px', color: '#72757c' }}>{emailBrand.name} · 账号交付</Text>
          <Heading as="h1" style={{ fontSize: '26px' }}>您的 Gemini Plus 账号已发货。</Heading>
          <Text style={{ fontSize: '14px', lineHeight: '24px' }}>下方为<strong>账号资料原文</strong>，各段以英文冒号分隔，通常是<strong>邮箱、邮箱密码、账号密码与邮箱验证链接</strong>。最后一段若是网址，请从 <strong>https://</strong> 开始整段复制，不要按冒号拆开。请勿转发给他人。</Text>
          <Section style={{ background: '#f1f6fd', padding: '20px', borderRadius: '12px' }}><Text style={{ fontSize: '14px', lineHeight: '24px', whiteSpace: 'pre-wrap', wordBreak: 'break-all', fontFamily: 'monospace' }}>__RECHARGE_DELIVERY_VALUE__</Text></Section>
          <Text style={{ fontSize: '14px', fontWeight: 600 }}>使用说明</Text><Text style={{ fontSize: '14px', lineHeight: '24px', whiteSpace: 'pre-line' }}>__RECHARGE_DELIVERY_INSTRUCTIONS__</Text>
          <Text style={{ fontSize: '12px', color: '#72757c', wordBreak: 'break-all' }}>订单号：__RECHARGE_DELIVERY_ORDER__</Text>
          <Text style={{ fontSize: '13px', lineHeight: '22px' }}>请在付款后 30 分钟内登录并修改账号密码——超过 30 分钟不再受理售后。如需帮助，请回到 <a href={emailBrand.orders}>我的订单</a>，通过原订单联系支持。邮件发出表示资料已交付，实际登录结果请自行核对。</Text>
        </Section><Text style={{ wordBreak: 'break-all', fontSize: '11px', color: '#858790' }}>{emailBrand.name} · 独立服务商，此邮件仅用于商品交付，请勿直接回复。</Text>
      </Container>
    </Body></Html>
}
