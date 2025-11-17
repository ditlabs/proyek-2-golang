package models

type Barang struct {
	ID             int     `json:"id"`
	KodeBarang     string  `json:"kode_barang"`
	NamaBarang     string  `json:"nama_barang"`
	Deskripsi      string  `json:"deskripsi"`
	JumlahTotal    int     `json:"jumlah_total"`
	JumlahTersedia int     `json:"jumlah_tersedia"`
	RuanganID      *int    `json:"ruangan_id"`
	Ruangan        *Ruangan `json:"ruangan,omitempty"`
}

type CreateBarangRequest struct {
	KodeBarang  string `json:"kode_barang"`
	NamaBarang  string `json:"nama_barang"`
	Deskripsi   string `json:"deskripsi"`
	JumlahTotal int    `json:"jumlah_total"`
	RuanganID   *int   `json:"ruangan_id"`
}

type UpdateBarangRequest struct {
	NamaBarang  string `json:"nama_barang"`
	Deskripsi   string `json:"deskripsi"`
	JumlahTotal int    `json:"jumlah_total"`
	RuanganID   *int   `json:"ruangan_id"`
}

