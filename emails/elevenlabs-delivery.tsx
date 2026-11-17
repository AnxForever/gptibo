import { emailBrand } from './brand.ts'
import * as React from 'react'
import { Body, Container, Head, Heading, Html, Img, Preview, Section, Text } from 'react-email'

export default function ElevenlabsDeliveryEmail() {
  return <Html lang="zh-CN"><Head /><Preview>您的 ElevenLabs 账号资料已准备好，请妥善保管。</Preview>
    <Body style={{ margin: 0, padding: '40px 16px', background: '#f5f5f7', color: '#202124', fontFamily: '-apple-system,BlinkMacSystemFont,"Segoe UI","PingFang SC",Arial,sans-serif' }}>
      <Container style={{ maxWidth: '560px' }}><Img src={emailBrand.logo} width="32" height="32" alt={emailBrand.name} style={{ marginBottom: '24px' }} />
        <Section style={{ padding: '32px 24px', background: '#fff', borderRadius: '20px', borderTop: '4px solid #0064e0' }}>
          <Text style={{ wordBreak: 'break-all', fontSize: '12px', color: '#72757c' }}>{emailBrand.name} · 账号交付</Text>
          <Heading as="h1" style={{ fontSize: '26px' }}>您的 ElevenLabs 账号已发货。</Heading>
          <Text style={{ fontSize: '14px', lineHeight: '24px' }}>下方四段资料依次是<strong>邮箱、密码、2FA 密钥、邮箱验证链接</strong>，以英文冒号分隔。请勿转发给他人。</Text>
          <Section style={{ background: '#f1f6fd', padding: '20px', borderRadius: '12px' }}><Text style={{ fontSize: '14px', lineHeight: '24px', whiteSpace: 'pre-wrap', wordBreak: 'break-all', fontFamily: 'monospace' }}>__RECHARGE_DELIVERY_VALUE__</Text></Section>
          <Text style={{ fontSize: '14px', fontWeight: 600 }}>使用说明</Text><Text style={{ fontSize: '14px', lineHeight: '24px', whiteSpace: 'pre-line' }}>__RECHARGE_DELIVERY_INSTRUCTIONS__</Text>
          <Text style={{ fontSize: '12px', color: '#72757c', wordBreak: 'break-all' }}>订单号：__RECHARGE_DELIVERY_ORDER__</Text>
          <Text style={{ fontSize: '13px', lineHeight: '22px' }}>首次登录后请立即修改密码。如需帮助，请回到 <a href={emailBrand.orders}>我的订单</a>，通过原订单联系支持或申请售后。邮件发出表示资料已交付，实际登录结果请自行核对。</Text>
        </Section><Text style={{ wordBreak: 'break-all', fontSize: '11px', color: '#858790' }}>{emailBrand.name} · 独立服务商，此邮件仅用于商品交付，请勿直接回复。</Text>
      </Container>
    </Body></Html>
}
