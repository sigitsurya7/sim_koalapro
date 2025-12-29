package repository

import (
	"context"
	"database/sql"
	"time"

	"koalbot_api/internal/model"
)

type PenggunaDetailRepository struct {
	db *sql.DB
}

func NewPenggunaDetailRepository(db *sql.DB) *PenggunaDetailRepository {
	return &PenggunaDetailRepository{db: db}
}

func (r *PenggunaDetailRepository) GetByPenggunaID(ctx context.Context, penggunaID int64) (model.PenggunaDetail, error) {
	var detail model.PenggunaDetail
	var avatar sql.NullString
	var registeredAt sql.NullTime
	var coupons sql.NullString
	var freeDeals sql.NullString
	var pushCategories sql.NullString

	err := r.db.QueryRowContext(ctx, `
		SELECT id, pengguna_id, avatar, first_name, last_name, nickname,
			balance, balance_version, bonus, gender, email, email_verified,
			phone, phone_verified, phone_prefix, receive_news, receive_sms,
			receive_notification, country, country_name, currency, birthday,
			activate, password_is_set, tutorial, coupons, free_deals, blocked,
			agree_risk, agreed, status_group, docs_verified, registered_at,
			status_by_deposit, status_id, deposits_sum, push_notification_categories,
			preserve_name, registration_country_iso
		FROM t_pengguna_detail
		WHERE pengguna_id = $1
	`, penggunaID).Scan(
		&detail.ID,
		&detail.PenggunaID,
		&avatar,
		&detail.FirstName,
		&detail.LastName,
		&detail.Nickname,
		&detail.Balance,
		&detail.BalanceVersion,
		&detail.Bonus,
		&detail.Gender,
		&detail.Email,
		&detail.EmailVerified,
		&detail.Phone,
		&detail.PhoneVerified,
		&detail.PhonePrefix,
		&detail.ReceiveNews,
		&detail.ReceiveSMS,
		&detail.ReceiveNotification,
		&detail.Country,
		&detail.CountryName,
		&detail.Currency,
		&detail.Birthday,
		&detail.Activate,
		&detail.PasswordIsSet,
		&detail.Tutorial,
		&coupons,
		&freeDeals,
		&detail.Blocked,
		&detail.AgreeRisk,
		&detail.Agreed,
		&detail.StatusGroup,
		&detail.DocsVerified,
		&registeredAt,
		&detail.StatusByDeposit,
		&detail.StatusID,
		&detail.DepositsSum,
		&pushCategories,
		&detail.PreserveName,
		&detail.RegistrationCountryISO,
	)
	if err != nil {
		return model.PenggunaDetail{}, err
	}

	if avatar.Valid {
		val := avatar.String
		detail.Avatar = &val
	}
	if registeredAt.Valid {
		val := registeredAt.Time
		detail.RegisteredAt = &val
	}
	if coupons.Valid {
		detail.Coupons = []byte(coupons.String)
	}
	if freeDeals.Valid {
		detail.FreeDeals = []byte(freeDeals.String)
	}
	if pushCategories.Valid {
		detail.PushNotificationCategories = []byte(pushCategories.String)
	}

	return detail, nil
}

func (r *PenggunaDetailRepository) Upsert(ctx context.Context, detail model.PenggunaDetail) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO t_pengguna_detail (
			id, pengguna_id, avatar, first_name, last_name, nickname,
			balance, balance_version, bonus, gender, email, email_verified,
			phone, phone_verified, phone_prefix, receive_news, receive_sms,
			receive_notification, country, country_name, currency, birthday,
			activate, password_is_set, tutorial, coupons, free_deals, blocked,
			agree_risk, agreed, status_group, docs_verified, registered_at,
			status_by_deposit, status_id, deposits_sum, push_notification_categories,
			preserve_name, registration_country_iso
		) VALUES (
			$1,$2,$3,$4,$5,$6,
			$7,$8,$9,$10,$11,$12,
			$13,$14,$15,$16,$17,
			$18,$19,$20,$21,$22,
			$23,$24,$25,$26,$27,$28,
			$29,$30,$31,$32,$33,
			$34,$35,$36,$37,$38,$39
		)
		ON CONFLICT (id) DO UPDATE SET
			pengguna_id = EXCLUDED.pengguna_id,
			avatar = EXCLUDED.avatar,
			first_name = EXCLUDED.first_name,
			last_name = EXCLUDED.last_name,
			nickname = EXCLUDED.nickname,
			balance = EXCLUDED.balance,
			balance_version = EXCLUDED.balance_version,
			bonus = EXCLUDED.bonus,
			gender = EXCLUDED.gender,
			email = EXCLUDED.email,
			email_verified = EXCLUDED.email_verified,
			phone = EXCLUDED.phone,
			phone_verified = EXCLUDED.phone_verified,
			phone_prefix = EXCLUDED.phone_prefix,
			receive_news = EXCLUDED.receive_news,
			receive_sms = EXCLUDED.receive_sms,
			receive_notification = EXCLUDED.receive_notification,
			country = EXCLUDED.country,
			country_name = EXCLUDED.country_name,
			currency = EXCLUDED.currency,
			birthday = EXCLUDED.birthday,
			activate = EXCLUDED.activate,
			password_is_set = EXCLUDED.password_is_set,
			tutorial = EXCLUDED.tutorial,
			coupons = EXCLUDED.coupons,
			free_deals = EXCLUDED.free_deals,
			blocked = EXCLUDED.blocked,
			agree_risk = EXCLUDED.agree_risk,
			agreed = EXCLUDED.agreed,
			status_group = EXCLUDED.status_group,
			docs_verified = EXCLUDED.docs_verified,
			registered_at = EXCLUDED.registered_at,
			status_by_deposit = EXCLUDED.status_by_deposit,
			status_id = EXCLUDED.status_id,
			deposits_sum = EXCLUDED.deposits_sum,
			push_notification_categories = EXCLUDED.push_notification_categories,
			preserve_name = EXCLUDED.preserve_name,
			registration_country_iso = EXCLUDED.registration_country_iso
	`,
		detail.ID,
		detail.PenggunaID,
		nullableString(detail.Avatar),
		detail.FirstName,
		detail.LastName,
		detail.Nickname,
		detail.Balance,
		detail.BalanceVersion,
		detail.Bonus,
		detail.Gender,
		detail.Email,
		detail.EmailVerified,
		detail.Phone,
		detail.PhoneVerified,
		detail.PhonePrefix,
		detail.ReceiveNews,
		detail.ReceiveSMS,
		detail.ReceiveNotification,
		detail.Country,
		detail.CountryName,
		detail.Currency,
		detail.Birthday,
		detail.Activate,
		detail.PasswordIsSet,
		detail.Tutorial,
		nullableJSON(detail.Coupons),
		nullableJSON(detail.FreeDeals),
		detail.Blocked,
		detail.AgreeRisk,
		detail.Agreed,
		detail.StatusGroup,
		detail.DocsVerified,
		nullableTime(detail.RegisteredAt),
		detail.StatusByDeposit,
		detail.StatusID,
		detail.DepositsSum,
		nullableJSON(detail.PushNotificationCategories),
		detail.PreserveName,
		detail.RegistrationCountryISO,
	)
	return err
}

func nullableString(val *string) sql.NullString {
	if val == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *val, Valid: true}
}

func nullableJSON(val []byte) sql.NullString {
	if len(val) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{String: string(val), Valid: true}
}

func nullableTime(val *time.Time) sql.NullTime {
	if val == nil {
		return sql.NullTime{}
	}
	return sql.NullTime{Time: *val, Valid: true}
}
