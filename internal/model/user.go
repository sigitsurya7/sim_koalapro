package model

import "time"

type User struct {
	ID        int64
	UID       string
	Username  string
	Password  string
	Role      string
	Active    bool
	CreatedAt time.Time
	CreatedBy *string
	UpdatedAt *time.Time
	UpdatedBy *string
	DeletedAt *time.Time
	LastSeen  *time.Time
}

type MasterPengguna struct {
	ID         int64
	UUID       string
	IDPengguna int64
	Telegram   *string
	Jenis      string
	Active     bool
	CreatedAt  time.Time
	UpdatedAt  *time.Time
	DeletedAt  *time.Time
}

type PenggunaDetail struct {
	ID                         int64
	PenggunaID                 int64
	Avatar                     *string
	FirstName                  string
	LastName                   string
	Nickname                   string
	Balance                    float64
	BalanceVersion             int64
	Bonus                      float64
	Gender                     string
	Email                      string
	EmailVerified              bool
	Phone                      string
	PhoneVerified              bool
	PhonePrefix                string
	ReceiveNews                bool
	ReceiveSMS                 bool
	ReceiveNotification        bool
	Country                    string
	CountryName                string
	Currency                   string
	Birthday                   string
	Activate                   bool
	PasswordIsSet              bool
	Tutorial                   bool
	Coupons                    []byte
	FreeDeals                  []byte
	Blocked                    bool
	AgreeRisk                  bool
	Agreed                     bool
	StatusGroup                string
	DocsVerified               bool
	RegisteredAt               *time.Time
	StatusByDeposit            string
	StatusID                   int
	DepositsSum                float64
	PushNotificationCategories []byte
	PreserveName               bool
	RegistrationCountryISO     string
}
