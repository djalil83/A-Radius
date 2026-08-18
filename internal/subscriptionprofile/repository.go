package subscriptionprofile

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

type Repository struct{ db *sql.DB }

func NewRepository(db *sql.DB) *Repository { return &Repository{db: db} }

const profileColumns = `
 id, tenant_id, name, service_type, category, media, color, description, status,
 mikrotik_group, radius_group, rate_limit, upload_bps, download_bps, shared_users,
 vlan_id, olt_profile, ip_pool, monthly_price, active_days, commission_amount,
 commission_type, billing_cycle, auto_isolate, billing_note, version, created_by,
 updated_by, created_at, updated_at`

func scanProfile(scanner interface{ Scan(...any) error }) (Profile, error) {
	var p Profile
	var category, media, description, mt, radius, rate, olt, pool, note, createdBy, updatedBy sql.NullString
	var upload, download sql.NullInt64
	var vlan sql.NullInt64
	err := scanner.Scan(
		&p.ID, &p.TenantID, &p.Name, &p.ServiceType, &category, &media, &p.Color,
		&description, &p.Status, &mt, &radius, &rate, &upload, &download,
		&p.SharedUsers, &vlan, &olt, &pool, &p.MonthlyPrice, &p.ActiveDays,
		&p.CommissionAmount, &p.CommissionType, &p.BillingCycle, &p.AutoIsolate,
		&note, &p.Version, &createdBy, &updatedBy, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return p, err
	}
	p.Category, p.Media, p.Description = nullableString(category), nullableString(media), nullableString(description)
	p.MikrotikGroup, p.RadiusGroup, p.RateLimit = nullableString(mt), nullableString(radius), nullableString(rate)
	p.OLTProfile, p.IPPool, p.BillingNote = nullableString(olt), nullableString(pool), nullableString(note)
	p.CreatedBy, p.UpdatedBy = nullableString(createdBy), nullableString(updatedBy)
	if upload.Valid {
		p.UploadBPS = &upload.Int64
	}
	if download.Valid {
		p.DownloadBPS = &download.Int64
	}
	if vlan.Valid {
		n := int(vlan.Int64)
		p.VLANID = &n
	}
	return p, nil
}

func nullableString(v sql.NullString) *string {
	if !v.Valid {
		return nil
	}
	s := v.String
	return &s
}

func (r *Repository) List(ctx context.Context, tenantID, q, serviceType, status string, limit, offset int) ([]Profile, error) {
	where := []string{"tenant_id = $1", "deleted_at IS NULL"}
	args := []any{tenantID}
	n := 2
	if q != "" {
		where = append(where, fmt.Sprintf("(name ILIKE $%d OR COALESCE(description, '') ILIKE $%d)", n, n))
		args = append(args, "%"+q+"%")
		n++
	}
	if serviceType != "" {
		where = append(where, fmt.Sprintf("service_type = $%d", n))
		args = append(args, serviceType)
		n++
	}
	if status != "" {
		where = append(where, fmt.Sprintf("status = $%d", n))
		args = append(args, status)
		n++
	}
	args = append(args, limit, offset)
	query := fmt.Sprintf("SELECT %s FROM subscription_profiles WHERE %s ORDER BY updated_at DESC, id LIMIT $%d OFFSET $%d", profileColumns, strings.Join(where, " AND "), n, n+1)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]Profile, 0)
	for rows.Next() {
		p, err := scanProfile(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, p)
	}
	return items, rows.Err()
}

func (r *Repository) Get(ctx context.Context, tenantID, id string) (Profile, error) {
	row := r.db.QueryRowContext(ctx, fmt.Sprintf("SELECT %s FROM subscription_profiles WHERE tenant_id=$1 AND id=$2 AND deleted_at IS NULL", profileColumns), tenantID, id)
	p, err := scanProfile(row)
	if errors.Is(err, sql.ErrNoRows) {
		return p, ErrNotFound
	}
	return p, err
}

func (r *Repository) Create(ctx context.Context, tenantID, actorID string, req CreateRequest) (Profile, error) {
	const query = `INSERT INTO subscription_profiles
	(tenant_id,name,service_type,category,media,color,description,status,mikrotik_group,radius_group,rate_limit,upload_bps,download_bps,shared_users,vlan_id,olt_profile,ip_pool,monthly_price,active_days,commission_amount,commission_type,billing_cycle,auto_isolate,billing_note,created_by,updated_by)
	VALUES ($1,$2,$3,$4,$5,$6,$7,'ACTIVE',$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$25)
	RETURNING ` + profileColumns
	return scanProfile(r.db.QueryRowContext(ctx, query, tenantID, req.Name, req.ServiceType, req.Category, req.Media, req.Color, req.Description, req.MikrotikGroup, req.RadiusGroup, req.RateLimit, req.UploadBPS, req.DownloadBPS, req.SharedUsers, req.VLANID, req.OLTProfile, req.IPPool, req.MonthlyPrice, req.ActiveDays, req.CommissionAmount, req.CommissionType, req.BillingCycle, req.AutoIsolate, req.BillingNote, actorID))
}

func (r *Repository) Update(ctx context.Context, tenantID, id, actorID string, req UpdateRequest) (Profile, error) {
	const query = `UPDATE subscription_profiles SET name=$3,service_type=$4,category=$5,media=$6,color=$7,description=$8,mikrotik_group=$9,radius_group=$10,rate_limit=$11,upload_bps=$12,download_bps=$13,shared_users=$14,vlan_id=$15,olt_profile=$16,ip_pool=$17,monthly_price=$18,active_days=$19,commission_amount=$20,commission_type=$21,billing_cycle=$22,auto_isolate=$23,billing_note=$24,updated_by=$25,version=version+1 WHERE tenant_id=$1 AND id=$2 AND deleted_at IS NULL AND version=$26 RETURNING ` + profileColumns
	args := []any{tenantID, id, req.Name, req.ServiceType, req.Category, req.Media, req.Color, req.Description, req.MikrotikGroup, req.RadiusGroup, req.RateLimit, req.UploadBPS, req.DownloadBPS, req.SharedUsers, req.VLANID, req.OLTProfile, req.IPPool, req.MonthlyPrice, req.ActiveDays, req.CommissionAmount, req.CommissionType, req.BillingCycle, req.AutoIsolate, req.BillingNote, actorID, req.Version}
	p, err := scanProfile(r.db.QueryRowContext(ctx, query, args...))
	if errors.Is(err, sql.ErrNoRows) {
		return p, r.classifyUpdateMiss(ctx, tenantID, id, req.Version)
	}
	return p, err
}

func (r *Repository) classifyUpdateMiss(ctx context.Context, tenantID, id string, version int64) error {
	var current int64
	var deletedAt sql.NullTime
	err := r.db.QueryRowContext(ctx, "SELECT version, deleted_at FROM subscription_profiles WHERE tenant_id=$1 AND id=$2", tenantID, id).Scan(&current, &deletedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}
	if deletedAt.Valid || current != version {
		return ErrConflict
	}
	return ErrConflict
}

func (r *Repository) Archive(ctx context.Context, tenantID, id, actorID string, version int64) error {
	res, err := r.db.ExecContext(ctx, `UPDATE subscription_profiles SET status='ARCHIVED',deleted_at=now(),updated_by=$3,version=version+1 WHERE tenant_id=$1 AND id=$2 AND deleted_at IS NULL AND version=$4`, tenantID, id, actorID, version)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 1 {
		return nil
	}
	var current int64
	var deletedAt sql.NullTime
	err = r.db.QueryRowContext(ctx, "SELECT version, deleted_at FROM subscription_profiles WHERE tenant_id=$1 AND id=$2", tenantID, id).Scan(&current, &deletedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}
	return ErrConflict
}

func (r *Repository) Revisions(ctx context.Context, tenantID, id string, limit int) ([]Revision, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id,profile_id,version,operation,snapshot,changed_by,changed_at FROM subscription_profile_revisions WHERE tenant_id=$1 AND profile_id=$2 ORDER BY version DESC LIMIT $3`, tenantID, id, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]Revision, 0)
	for rows.Next() {
		var v Revision
		var changed sql.NullString
		if err := rows.Scan(&v.ID, &v.ProfileID, &v.Version, &v.Operation, &v.Snapshot, &changed, &v.ChangedAt); err != nil {
			return nil, err
		}
		v.ChangedBy = nullableString(changed)
		out = append(out, v)
	}
	return out, rows.Err()
}
