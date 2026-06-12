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

func (r *VoucherRepository) GetAll(ctx context.Context) ([]models.Voucher, error) {
	query := `
		SELECT id, code, discount_type, discount_value, min_purchase, max_discount,
		       valid_from, valid_until, usage_limit, used_count, created_at, updated_at
		FROM vouchers
		ORDER BY created_at DESC
	`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vouchers []models.Voucher
	for rows.Next() {
		var v models.Voucher
		err := rows.Scan(
			&v.ID, &v.Code, &v.DiscountType, &v.DiscountValue, &v.MinPurchase, &v.MaxDiscount,
			&v.ValidFrom, &v.ValidUntil, &v.UsageLimit, &v.UsedCount, &v.CreatedAt, &v.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		vouchers = append(vouchers, v)
	}
	return vouchers, nil
}

func (r *VoucherRepository) Create(ctx context.Context, v *models.Voucher) error {
	query := `
		INSERT INTO vouchers (code, discount_type, discount_value, min_purchase, max_discount, valid_from, valid_until, usage_limit)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, used_count, created_at, updated_at
	`
	err := r.pool.QueryRow(ctx, query, v.Code, v.DiscountType, v.DiscountValue, v.MinPurchase, v.MaxDiscount, v.ValidFrom, v.ValidUntil, v.UsageLimit).Scan(
		&v.ID, &v.UsedCount, &v.CreatedAt, &v.UpdatedAt,
	)
	return err
}

func (r *VoucherRepository) Update(ctx context.Context, id int, v *models.Voucher) error {
	query := `
		UPDATE vouchers 
		SET code = $1, discount_type = $2, discount_value = $3, min_purchase = $4, max_discount = $5, valid_from = $6, valid_until = $7, usage_limit = $8, updated_at = CURRENT_TIMESTAMP
		WHERE id = $9
	`
	cmdTag, err := r.pool.Exec(ctx, query, v.Code, v.DiscountType, v.DiscountValue, v.MinPurchase, v.MaxDiscount, v.ValidFrom, v.ValidUntil, v.UsageLimit, id)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return errors.New("voucher not found")
	}
	return nil
}

func (r *VoucherRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM vouchers WHERE id = $1`
	cmdTag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return errors.New("voucher not found")
	}
	return nil
}
