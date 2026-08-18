# PostgreSQL Disaster Recovery Runbook

Runbook ini menggunakan `pg_dump` custom format pada PVC `a-radius-postgres-backups`. PVC lokal membantu pemulihan operasional, tetapi **bukan pengganti salinan off-cluster**. Untuk disaster recovery terhadap kehilangan cluster, salin file `.dump`, `.sha256`, dan `.metadata` ke object storage atau lokasi berbeda setelah setiap backup.

## Backup otomatis

Apply manifest dengan Kustomize:

```bash
kubectl apply -k k8s/
kubectl -n a-radius get cronjob/a-radius-postgres-backup
kubectl -n a-radius create job --from=cronjob/a-radius-postgres-backup a-radius-postgres-backup-manual
kubectl -n a-radius logs -f job/a-radius-postgres-backup-manual
```

CronJob berjalan setiap hari pukul 02:17 UTC, tidak mengizinkan job overlap, menyimpan histori tiga job berhasil dan tiga job gagal, serta menghapus backup lokal lebih lama dari 30 hari. Setiap dump divalidasi dengan `pg_restore --list` dan SHA-256.

## Menemukan backup

```bash
kubectl -n a-radius run backup-shell --rm -it --restart=Never \
  --image=postgres:17-alpine --overrides='{"spec":{"containers":[{"name":"backup-shell","image":"postgres:17-alpine","command":["sh"],"stdin":true,"tty":true,"volumeMounts":[{"name":"backups","mountPath":"/backups"}]}],"volumes":[{"name":"backups","persistentVolumeClaim":{"claimName":"a-radius-postgres-backups"}}]}}'
```

Di dalam shell, gunakan `ls -lh /backups`. Jangan mengubah nama atau isi file checksum/metadata.

## Restore terkontrol

Restore harus dilakukan pada maintenance window. Hentikan API agar tidak menulis selama restore, pilih file backup yang tepat, verifikasi checksum, dan minta persetujuan operator kedua.

```bash
kubectl -n a-radius scale deployment/a-radius-api --replicas=0
kubectl -n a-radius get job/a-radius-postgres-restore-example
kubectl -n a-radius patch job/a-radius-postgres-restore-example \
  --type=merge -p '{"spec":{"suspend":false}}'
```

Sebelum patch, ubah `BACKUP_FILE` ke nama file yang benar dan ubah `CONFIRM_RESTORE` menjadi `YES` pada manifest/job. Untuk restore destructive ke database yang sama, set `ALLOW_DROP=YES` hanya setelah backup kondisi saat ini dibuat dan diverifikasi. Alternatif yang lebih aman adalah restore ke database baru, verifikasi data, lalu lakukan cutover terencana.

```bash
kubectl -n a-radius wait --for=condition=complete job/a-radius-postgres-restore-example --timeout=30m
kubectl -n a-radius logs job/a-radius-postgres-restore-example
kubectl -n a-radius scale deployment/a-radius-api --replicas=2
kubectl -n a-radius rollout status deployment/a-radius-api --timeout=10m
```

## Verifikasi pasca-restore

Periksa jumlah tabel, status migration, satu sampel profile, revision history, audit trail, dan endpoint `/healthz`. Jalankan smoke test API dengan JWT yang sesuai dan pastikan tidak ada error pada log PostgreSQL/API.

## Target operasional

Target awal yang disarankan adalah **RPO 24 jam** dengan CronJob harian dan **RTO 30–60 menit**, bergantung pada ukuran database dan provisioning storage. Target ini harus diuji melalui restore drill berkala, bukan diasumsikan dari konfigurasi.

## Keamanan

Secret Kubernetes harus dikelola melalui External Secrets, Sealed Secrets, atau secret manager cloud; jangan commit password nyata ke Git. Backup mengandung data sensitif dan harus dienkripsi saat transit maupun saat disimpan. Akses PVC backup perlu dibatasi dan salinan off-cluster harus memiliki retensi serta lifecycle policy terpisah.
