package backend

type Surah struct {
	ID          int    `db:"id"          json:"id"`
	Name        string `db:"name"        json:"name"`
	Translation string `db:"translation" json:"translation"`
	Translated  bool   `db:"translated"  json:"translated"`
}

type Ayah struct {
	ID          int    `db:"id"          json:"id"`
	Arabic      string `db:"arabic"      json:"arabic"`
	Translation string `db:"translation" json:"translation"`
	Tafsir      string `db:"tafsir"      json:"tafsir"`
}

type Word struct {
	ID          int    `db:"id"           json:"id"`
	Ayah        int    `db:"ayah"         json:"ayah"`
	Position    int    `db:"position"     json:"position"`
	Arabic      string `db:"arabic"       json:"arabic"`
	Translation string `db:"translation"  json:"translation"`
	IsSeparator bool   `db:"is_separator" json:"isSeparator"`
}

type Choice struct {
	Text      string `db:"text"       json:"text"`
	IsCorrect bool   `db:"is_correct" json:"isCorrect"`
}
