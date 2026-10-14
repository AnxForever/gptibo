#!/bin/sh
# Push this checkout to the deployment host and run the rollout there.
#
#   npm run build
#   ANX_HOST=korea-seoul ./services/storefront/scripts/deploy.sh [release]
#
# POSIX sh on purpose: the bash on a Windows dev box is sometimes a shim that
# rejects bash-only options, and this script has no need for them.
#
# The host must already hold the compose stack, its .secrets and nginx. This
# helper transfers source and built frontend output only: secrets, .env files
# and build artifacts are excluded so a deploy can never overwrite the
# server's own configuration.
set -eu

HOST=${ANX_HOST:-korea-seoul}
DIR=${ANX_DIR:-/opt/anx}
RELEASE=${1:-}

root=$(cd "$(dirname "$0")/../../.." && pwd)
cd "$root"

if [ ! -f dist/index.html ]; then
  echo "dist/index.html is missing — run 'npm run build' first." >&2
  exit 1
fi

# Fail early and clearly. The host is an ssh target, which usually means an
# alias from ~/.ssh/config; without this check a dead connection only shows up
# as a tar write error.
if ! ssh -o BatchMode=yes -o ConnectTimeout=15 "$HOST" true 2>/dev/null; then
  echo "cannot reach '$HOST' over ssh." >&2
  echo "Check the alias in ~/.ssh/config, or pass ANX_HOST=user@host." >&2
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
# sh, not ./: the executable bit does not survive every checkout on Windows.
ssh "$HOST" "cd '$DIR/services/storefront' && sh scripts/rollout.sh $RELEASE"
