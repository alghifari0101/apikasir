package models

type DailyReport struct {
	TotalRevenue   int               `json:"total_revenue"`
	TotalTransaksi int               `json:"total_transaksi"`
	ProdukTerlaris BestSellerProduct `json:"produk_terlaris"`
}

type BestSellerProduct struct {
	Nama       string `json:"nama"`
	QtyTerjual int    `json:"qty_terjual"`
}

type ReportRequest struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type ReportResponse struct {
	Transactions []Transaction `json:"transactions"`
	TotalAmount  int           `json:"total_amount"`
}
