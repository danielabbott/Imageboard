package database

import "log"

func GetImageTitle(db *DB, id ImageID) (string, error) {
	rows := get_statement(db, "SELECT title FROM image WHERE id = $1;").QueryRow(id)

	var title string
	err := rows.Scan(&title)
	if err != nil {
		return "", err
	}

	return title, nil
}

func GetImageTags(db *DB, id ImageID) ([]string, error) {
	rows, err := get_statement(db,
		`SELECT tag.name FROM tag INNER JOIN image_tag ON tag.id = image_tag.tag_id
		WHERE (image_tag.deleted IS FALSE) AND image_tag.image_id = $1 ORDER BY image_tag.ser ASC;`).Query(id)
	if err != nil {
		return nil, err
	}

	tags := []string{}

	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			log.Println(err)
		} else {
			tags = append(tags, name)
		}
	}

	return tags, nil
}

func RemoveImage(db *DB, id ImageID) {
	_, _ = get_statement(db, "DELETE FROM image WHERE id=$1").Exec(id)
}

func NewImage(db *DB, title string) (ImageID, error) {
	res := get_statement(db, "INSERT INTO image(title) VALUES($1) RETURNING id").QueryRow(title)

	var id int32
	err := res.Scan(&id)
	if err != nil {
		return 0, err
	}

	return ImageID(id), nil
}

func DeleteImage(db *DB, id ImageID) {
	get_statement(db, "DELETE FROM image WHERE id=$1").Exec(id)
}
