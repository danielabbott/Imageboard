package database

import (
	"errors"
	"fmt"
	"strings"
)

func RemoveTag(db *DB, image_id ImageID, tag string) {
	tag = strings.ToLower(strings.Trim(tag, " "))

	get_statement(db, "UPDATE image_tag SET deleted=TRUE WHERE image_id=$1 AND tag_id = "+
		"(SELECT id FROM tag WHERE name=$2)").Exec(image_id, tag)
}

func validate_tag_to_add(tag string) (string, bool) {
	if len(tag) == 0 {
		return "", false
	}

	// Only allow non-control ASCII characters (no tabs or newlines)
	for i := 0; i < len(tag); i++ {
		c := tag[i]
		if c > '~' || c == ',' || !(c >= ' ' && c <= '~') {
			return "", false
		}
	}

	new_tag := strings.Trim(tag, " ")
	if len(new_tag) == 0 {
		return "", false
	}

	if new_tag[0] == '-' {
		return "", false
	}

	return new_tag, true
}

func AddTag(db *DB, image_id ImageID, tag_ string) error {
	tag, valid := validate_tag_to_add(tag_)
	if !valid {
		return errors.New("Invalid tag")
	}

	_, err := get_statement(db, "INSERT INTO tag(name) VALUES(LOWER($1)) ON CONFLICT DO NOTHING").Exec(tag)

	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = get_statement(db, "INSERT INTO image_tag(image_id,tag_id) "+
		"VALUES($1, (SELECT id FROM tag WHERE name=LOWER($2))) "+
		"ON CONFLICT(image_id,tag_id) DO "+
		"UPDATE SET deleted = FALSE").Exec(image_id, tag)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
