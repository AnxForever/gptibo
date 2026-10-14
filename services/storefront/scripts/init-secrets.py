#!/usr/bin/env python3
"""Initialize independent local development secrets without overwriting files.

Compose mounts these as file-backed secrets. The application image runs as an
unprivileged user (distroless nonroot, uid 65532) while the source files are
created by root, so the files must be group/world readable -- a 0600 secret
leaves `api` and `worker` failing to start with only a generic
"storefront command failed" line. Re-running this script repairs the mode of
secrets created by earlier versions.

Create the optional auth secrets by hand when enabling compose.auth.yaml:
`.secrets/admin-totp` and `.secrets/resend-key`.
"""
import os
import secrets
from pathlib import Path

# Readable by the nonroot container user; the directory stays 0700 because
# Docker opens the file as root before bind mounting it.
MODE = 0o644

root = Path(__file__).resolve().parents[1] / '.secrets'
root.mkdir(mode=0o700, exist_ok=True)
for name in ('owner-password', 'app-password', 'auth-secret'):
    path = root / name
    try:
        fd = os.open(path, os.O_WRONLY | os.O_CREAT | os.O_EXCL, MODE)
    except FileExistsError:
        repaired = os.name == 'posix' and path.stat().st_mode & 0o777 != MODE
        if repaired:
            os.chmod(path, MODE)
        print(f'{name}: already exists, ' + (f'mode repaired to {MODE:o}' if repaired else 'unchanged'))
        continue
    with os.fdopen(fd, 'w') as output:
        output.write(secrets.token_hex(32))
    print(f'{name}: created (value hidden)')
