package repository

const (
	// queries получить список песен с фильтрами
	getSongsQuery = `
		SELECT id, group_name, song_name, release_date, text, link, created_at, updated_at
		FROM songs
		WHERE ($1 = '' OR group_name ILIKE '%' || $1 || '%')
		AND ($2 = '' OR song_name ILIKE '%' || $2 || '%')
		AND ($3::timestamp IS NULL OR release_date >= $3)
		AND ($4::timestamp IS NULL OR release_date <= $4)
		AND ($5 = '' OR text ILIKE '%' || $5 || '%')
		AND ($6 = '' OR link ILIKE '%' || $6 || '%')
		ORDER BY created_at DESC
		LIMIT $7 OFFSET $8`

	// queries счетчик количества песен с фильтрами
	countSongsQuery = `
		SELECT COUNT(*)
		FROM songs
		WHERE ($1 = '' OR group_name ILIKE '%' || $1 || '%')
		AND ($2 = '' OR song_name ILIKE '%' || $2 || '%')
		AND ($3::timestamp IS NULL OR release_date >= $3)
		AND ($4::timestamp IS NULL OR release_date <= $4)
		AND ($5 = '' OR text ILIKE '%' || $5 || '%')
		AND ($6 = '' OR link ILIKE '%' || $6 || '%')`

	// queries получить песню по id
	getSongByIDQuery = `
		SELECT id, group_name, song_name, release_date, text, link, created_at, updated_at
		FROM songs
		WHERE id = $1`

	// queries создать песню
	createSongQuery = `
		INSERT INTO songs (group_name, song_name, release_date, text, link)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, group_name, song_name, release_date, text, link, created_at, updated_at`

	// update обновить песню
	updateSongQuery = `
		UPDATE songs
		SET group_name = $1, 
			song_name = $2, 
			release_date = $3,
			text = $4,
			link = $5,
			updated_at = NOW()
		WHERE id = $6
		RETURNING id, group_name, song_name, release_date, text, link, created_at, updated_at`

	// delete удалить песню
	deleteSongQuery = `DELETE FROM songs WHERE id = $1`

	// queries проверить существование песни
	checkSongExistsQuery = `
		SELECT EXISTS(
			SELECT 1 FROM songs 
			WHERE group_name = $1 
			AND song_name = $2 
			AND id != $3
		)`

	// queries проверить существование песни
	checkSongExistsForCreateQuery = `
		SELECT EXISTS(
			SELECT 1 FROM songs 
			WHERE group_name = $1 
			AND song_name = $2
		)`
)
