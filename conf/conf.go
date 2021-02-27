package conf

type Conf struct {
	BookID                 int    `json:"book_id,omitempty"`
	Login                  string `json:"login"`
	Password               string `json:"password"`
	BookPath               string `json:"book_path"`
	Format                 string `json:"format"`
	SearchByTitle          string `json:"search_by_title,omitempty"`
	NormalizedName         bool   `json:"normalized_name"`
	ShowAvailable4Download bool   `json:"show_available4_download"`
	Progress               bool   `json:"progress"`
	Verbose                bool   `json:"verbose"`
	Debug                  bool   `json:"debug"`
}
