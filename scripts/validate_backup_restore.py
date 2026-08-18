from pathlib import Path
import yaml

root = Path('.')
backup = (root / 'scripts/pg-backup.sh').read_text()
restore = (root / 'scripts/pg-restore.sh').read_text()
manifest = list(yaml.safe_load_all((root / 'k8s/backup-restore.yaml').read_text()))
manifest = [obj for obj in manifest if obj]
keys = {(obj['kind'], obj['metadata']['name']) for obj in manifest}
assert ('PersistentVolumeClaim', 'a-radius-postgres-backups') in keys
assert ('CronJob', 'a-radius-postgres-backup') in keys
assert ('Job', 'a-radius-postgres-restore-example') in keys
assert 'sha256sum "$dump"' in backup
assert 'pg_restore --list "$tmp_dump"' in backup
assert 'RETENTION_DAYS' in backup
assert 'CONFIRM_RESTORE=YES' in restore
assert 'checksum mismatch' in restore
assert 'pg_restore --clean --if-exists' in restore
cron = next(obj for obj in manifest if obj['kind'] == 'CronJob')
assert cron['spec']['concurrencyPolicy'] == 'Forbid'
job = next(obj for obj in manifest if obj['kind'] == 'Job')
assert job['spec']['suspend'] is True
print('Backup/restore structure OK')
