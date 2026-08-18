from pathlib import Path
import yaml

root = Path('k8s')
manifest = root / 'a-radius.yaml'
objects = list(yaml.safe_load_all(manifest.read_text()))
objects = [obj for obj in objects if obj]
keys = {(obj.get('kind'), obj.get('metadata', {}).get('name')) for obj in objects}
required = {
    ('Namespace', 'a-radius'),
    ('Secret', 'a-radius-secrets'),
    ('ConfigMap', 'a-radius-config'),
    ('PersistentVolumeClaim', 'a-radius-postgres-data'),
    ('Service', 'a-radius-postgres'),
    ('Deployment', 'a-radius-postgres'),
    ('Service', 'a-radius-api'),
    ('Deployment', 'a-radius-api'),
    ('Job', 'a-radius-db-migrate'),
}
missing = required - keys
assert not missing, f'missing resources: {sorted(missing)}'
assert (root / 'kustomization.yaml').exists()
for relative in [
    '../database/postgresql/migrations/0001_init.sql',
    '../database/migrations/0002_subscription_profiles.up.sql',
    '../docker/postgres/03-rbac-smoke.sql',
]:
    assert (root / relative).resolve().exists(), relative
api = next(obj for obj in objects if obj.get('kind') == 'Deployment' and obj['metadata']['name'] == 'a-radius-api')
container = api['spec']['template']['spec']['containers'][0]
assert container['readinessProbe']['httpGet']['path'] == '/healthz'
assert container['securityContext']['runAsNonRoot'] is True
job = next(obj for obj in objects if obj.get('kind') == 'Job')
assert job['spec']['template']['spec']['volumes'][0]['configMap']['name'] == 'a-radius-migrations'
print(f'Kubernetes manifest structure OK: {len(objects)} resources')
