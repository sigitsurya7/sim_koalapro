package repository

type UpdateMasterPenggunaRequest struct {
	IDPengguna *int64
	Telegram   *string
	Jenis      *string
	Active     *bool
}
