package database

const ddlCreateSurah = `
CREATE TABLE IF NOT EXISTS surah (
	id          INT  NOT NULL,
	name        TEXT NOT NULL,
	translation TEXT NOT NULL,
	start       INT  NOT NULL,
	end         INT  NOT NULL,
	PRIMARY KEY (id))`

const ddlCreateAyah = `
CREATE TABLE IF NOT EXISTS ayah (
	id          INT  NOT NULL,
	translation TEXT NOT NULL,
	tafsir      TEXT NOT NULL,
	PRIMARY KEY (id))`

const ddlCreateWord = `
CREATE TABLE IF NOT EXISTS word (
	id          INT  NOT NULL,
	ayah        INT  NOT NULL,
	position    INT  NOT NULL,
	arabic      TEXT NOT NULL,
	translation TEXT NOT NULL,
	PRIMARY KEY (id),
	CONSTRAINT word_UNIQUE UNIQUE (ayah, position),
	CONSTRAINT word_ayah_FK FOREIGN KEY (ayah) REFERENCES ayah (id))`

const ddlCreateTracker = `
CREATE TABLE IF NOT EXISTS tracker (
	id        INT NOT NULL,
	last_word INT DEFAULT NULL,
	PRIMARY KEY (id),
	CONSTRAINT tracker_word_FK FOREIGN KEY (last_word) REFERENCES word (id))`
