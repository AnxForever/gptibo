#!/usr/bin/env bash
# Push this checkout to the deployment host and run the server-side rollout.
#
#   npm run build
#   ANX_HOST=korea-seoul ./services/storefront/scripts/deploy.sh [release]
#
# The host is expected to already hold the compose stack, its .secrets and the
# anx-deploy.sh rollout script. This helper only transfers source and built
# frontend output; secrets, .env files and build artifacts are excluded so a
# deploy can never overwrite the server's own configuration.
set -euo pipefail

HOST=${ANX_HOST:-korea-seoul}
DIR=${ANX_DIR:-/opt/anx}
RELEASE=${1:-}

root="$(cd "$(dirname "$0")/../../.." && pwd)"
cd "$root"

if [ ! -f dist/index.html ]; then
  echo "dist/index.html is missing — run 'npm run build' first." >&2
  exit 1
fi

echo "== pushing services/storefront and dist to $HOST:$DIR"
tar czf - \
  --exclude='.secrets' \
  --exclude='.env' \
  --exclude='storefront.exe' \
  services/storefront dist \
  | ssh "$HOST" "mkdir -p '$DIR' && tar xzf - -C '$DIR'"

echo "== running the rollout on $HOST"
# shellcheck disable=SC2029
ssh "$HOST" "anx-deploy.sh $RELEASE"
