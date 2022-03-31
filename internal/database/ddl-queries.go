package database

const ddlCreateSurah = `
CREATE TABLE IF NOT EXISTS surah (
	id          INT UNSIGNED NOT NULL,
	name        TEXT         NOT NULL,
	translation TEXT         NOT NULL,
	PRIMARY KEY (id))`

const ddlCreateWord = `
CREATE TABLE IF NOT EXISTS word (
	id          INT UNSIGNED NOT NULL,
	surah       INT UNSIGNED NOT NULL,
	ayah        INT UNSIGNED NOT NULL,
	position    INT UNSIGNED NOT NULL,
	arabic      TEXT         NOT NULL,
	translation TEXT         NOT NULL,
	PRIMARY KEY (id),
	CONSTRAINT word_UNIQUE UNIQUE (surah, ayah, position),
	CONSTRAINT word_surah_FK FOREIGN KEY (surah) REFERENCES surah (id))`
