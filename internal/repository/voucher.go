package repository

import (
	"context"
	"errors"
	"koda-b6-backend/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type VoucherRepository struct {
	pool *pgxpool.Pool
}

func NewVoucherRepository(pool *pgxpool.Pool) *VoucherRepository {
	return &VoucherRepository{pool: pool}
}

func (r *VoucherRepository) GetByCode(ctx context.Context, code string) (*models.Voucher, error) {
	query := `
		SELECT id, code, discount_type, discount_value, min_purchase, max_discount,
		       valid_from, valid_until, usage_limit, used_count, created_at, updated_at
		FROM vouchers
		WHERE code = $1
	`
	var v models.Voucher
	err := r.pool.QueryRow(ctx, query, code).Scan(
		&v.ID, &v.Code, &v.DiscountType, &v.DiscountValue, &v.MinPurchase, &v.MaxDiscount,
		&v.ValidFrom, &v.ValidUntil, &v.UsageLimit, &v.UsedCount, &v.CreatedAt, &v.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("voucher not found")
		}
		return nil, err
	}
	return &v, nil
}

func (r *VoucherRepository) IncrementUsage(ctx context.Context, id int) error {
	query := `UPDATE vouchers SET used_count = used_count + 1, updated_at = CURRENT_TIMESTAMP WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}
