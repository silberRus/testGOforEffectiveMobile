package app

// SQL запросы для инициализации и проверки состояния базы данных.
// Эти запросы используются только при запуске приложения и не связаны с бизнес-логикой.
const (
	// checkDatabaseExistsQuery проверяет существование базы данных
	checkDatabaseExistsQuery = `
		SELECT EXISTS (
			SELECT 1 FROM pg_database 
			WHERE datname = $1
		)`

	// checkTableExistsQuery проверяет существование таблицы songs
	checkTableExistsQuery = `
		SELECT EXISTS (
			SELECT 1 FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = 'songs'
		)`

	// createSongsTableQuery создает таблицу songs, если она не существует.
	// Этот запрос дублирует миграцию 000001_init_schema.up.sql и используется
	// только как резервный механизм, если миграции не удалось применить.
	// В нормальном режиме таблица создается через миграции.
	createSongsTableQuery = `
		CREATE TABLE IF NOT EXISTS songs (
			id SERIAL PRIMARY KEY,
			group_name VARCHAR(255) NOT NULL,
			song_name VARCHAR(255) NOT NULL,
			release_date TIMESTAMP NOT NULL,
			text TEXT,
			link VARCHAR(255),
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
			UNIQUE(group_name, song_name)
		)`
)
