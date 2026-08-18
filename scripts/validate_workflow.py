from pathlib import Path

path = Path('.github/workflows/ci-cd.yml')
text = path.read_text()
required = [
    'name: CI/CD',
    'pull_request:',
    'workflow_dispatch:',
    'actions/setup-go@v5',
    'go test ./...',
    'go vet ./...',
    'go mod verify',
    'govulncheck ./...',
    'postgres:17-alpine',
    '0001_init.sql',
    '0002_subscription_profiles.up.sql',
    'docker/setup-buildx-action@v3',
    'docker/build-push-action@v6',
    'ghcr.io/${{ github.repository_owner }}/a-radius-profile-api',
    'secrets.GITHUB_TOKEN',
]
missing = [item for item in required if item not in text]
if missing:
    raise SystemExit('missing workflow fragments: ' + ', '.join(missing))
if 'push: ${{ github.event_name != \'pull_request\' }}' not in text:
    raise SystemExit('Docker push must be disabled for pull requests')
if 'permissions:\n  contents: read' not in text:
    raise SystemExit('workflow must default to read-only contents permission')
print(f'workflow structure OK: {path}')
