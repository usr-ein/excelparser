package xl

type RawWorkbook struct {
	Name   string     `json:"name"`
	Sheets []RawSheet `json:"sheets"`
}

type Workbook struct {
	Name   string  `json:"name"`
	Sheets []Sheet `json:"sheets"`
}
