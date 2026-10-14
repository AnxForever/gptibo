// The first-paint bootstrap in index.html is inline, and the deployed CSP has no
// 'unsafe-inline' for scripts. Rather than weaken that, allowlist the one script
// by its SHA-256 hash. Recompute on every build so the two can never drift.
import { createHash } from 'node:crypto'
import { readFileSync, writeFileSync, existsSync } from 'node:fs'
import { fileURLToPath } from 'node:url'
import { dirname, join } from 'node:path'

const dist = join(dirname(fileURLToPath(import.meta.url)), '..', 'dist')
const htmlPath = join(dist, 'index.html')
const headersPath = join(dist, '_headers')

if (!existsSync(htmlPath) || !existsSync(headersPath)) {
  console.error('stamp-csp-hash: run this after a build (dist/index.html and dist/_headers must exist)')
  process.exit(1)
}

const html = readFileSync(htmlPath, 'utf8')
const inline = [...html.matchAll(/<script(?![^>]*\bsrc=)[^>]*>([\s\S]*?)<\/script>/g)].map(m => m[1])
if (inline.length === 0) {
  console.error('stamp-csp-hash: no inline script found in dist/index.html — has the bootstrap moved?')
  process.exit(1)
}

const hashes = inline.map(code => `'sha256-${createHash('sha256').update(code, 'utf8').digest('base64')}'`)
const headers = readFileSync(headersPath, 'utf8')
const csp = /^(\s*Content-Security-Policy:\s*)(.*)$/m
if (!csp.test(headers)) {
  console.error('stamp-csp-hash: no Content-Security-Policy line in dist/_headers')
  process.exit(1)
}

const stamped = headers.replace(csp, (_line, prefix, policy) => {
  const cleaned = policy.replace(/\s*'sha256-[^']*'/g, '').trim()
  const scriptSrc = /script-src ([^;]*)/.exec(cleaned)
  if (!scriptSrc) {
    console.error('stamp-csp-hash: the policy has no script-src directive')
    process.exit(1)
  }
  const withHashes = `script-src ${scriptSrc[1].trim()} ${hashes.join(' ')}`
  return prefix + cleaned.replace(scriptSrc[0], withHashes)
})

writeFileSync(headersPath, stamped)
console.log(`stamp-csp-hash: allowlisted ${hashes.length} inline script(s) in the deployed CSP`)
