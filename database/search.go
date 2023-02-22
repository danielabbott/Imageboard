package database

import (
	"database/sql"
	"log"
	"strconv"
	"strings"

	"github.com/lib/pq"
)

// Returns (new tag, valid, is_excluded)
// This function is for searches (not for creating tags) as it removes the exlusion symbol (-)
func validate_tag(tag string) (string, bool, bool) {
	if len(tag) == 0 {
		return "", false, false
	}

	// Only allow non-control ASCII characters (no tabs or newlines)
	for i := 0; i < len(tag); i++ {
		c := tag[i]
		if c > '~' {
			return "", false, false
		}
		if !(c >= ' ' && c <= '~') {
			return "", false, false
		}
	}

	new_tag := strings.Trim(tag, " ")

	// A hyphen before a tag means that tag is excluded from the search
	excl := false
	if new_tag[0] == '-' {
		if len(new_tag) >= 2 && new_tag[1] == '-' {
			// Tags never start with a hyphen
			return "", false, false
		}
		new_tag = new_tag[1:]
		excl = true
	}

	return new_tag, len(new_tag) > 0, excl
}

func get_tag_ids_by_use_count(db *DB, tags []string) ([]int32, error) {
	tag_count_sql :=
		`SELECT tag.id
		FROM image_tag
		RIGHT JOIN tag
		ON tag.id = image_tag.tag_id
		WHERE tag.name = ANY($1) 
		GROUP BY tag.id
		ORDER BY COUNT(image_tag.tag_id) ASC;`

	rows, err := get_statement(db, tag_count_sql).Query(pq.Array(tags))
	if err != nil {
		return []int32{}, err
	}
	defer rows.Close()

	var tag_ids []int32
	for rows.Next() {
		var id int32
		err = rows.Scan(&id)
		if err != nil {
			log.Println(err)
		} else {
			tag_ids = append(tag_ids, id)
		}
	}

	return tag_ids, nil
}

func full_search(db *DB, last_seen_image_id_or_ser_value int32, tags []string, excluded_tags []string) (*sql.Rows, error) {
	// Get number of images in databsase with each tag
	// TODO create test database with 1 million images & tags and test to see if join order effects performance
	incl_tag_ids, err := get_tag_ids_by_use_count(db, tags)
	if err != nil {
		return nil, err
	}

	if len(incl_tag_ids) < len(tags) {
		// Some tags don't have any images. There will be 0 results.
		return db.db.Query("SELECT 0,0 WHERE FALSE;")
	}

	sql_str := ""

	// Included tags

	if len(incl_tag_ids) > 0 {
		sql_str = `SELECT image_id,ser FROM image_tag WHERE (deleted = FALSE) AND (tag_id = ` +
			strconv.Itoa(int(incl_tag_ids[0])) + `)`

		for _, tag_id := range incl_tag_ids[1:] {
			sql_str =
				`SELECT A.image_id, (CASE WHEN A.ser > B.ser THEN A.ser ELSE B.ser END) as ser FROM
				(` + sql_str + `) A
				INNER JOIN
				(SELECT image_id,ser FROM image_tag WHERE (deleted = FALSE) AND (tag_id = ` +
					strconv.Itoa(int(tag_id)) + `)) B ON A.image_id = B.image_id`
		}
	} else {
		// No included tags so start with all possible images.
		sql_str = "SELECT id as image_id, id as ser FROM image"
	}

	// Excluded tags
	if len(excluded_tags) > 0 {
		sql_str = "SELECT image_id, ser FROM (" + sql_str + `) A WHERE
		A.image_id NOT IN (SELECT image_id from image_tag WHERE tag_id IN (
		SELECT id FROM tag WHERE name = ANY($1)))`
	}

	// Filter & sort

	sql_str = `SELECT image_id, ser as cmp_to_val FROM (` + sql_str + `) A`
	if last_seen_image_id_or_ser_value > -1 {
		sql_str += " WHERE A.ser < " + strconv.Itoa(int(last_seen_image_id_or_ser_value))
	}

	sql_str += " ORDER BY A.ser DESC LIMIT 48;"
	// fmt.Println(sql_str)

	if len(excluded_tags) > 0 {
		return db.db.Query(sql_str, pq.Array(excluded_tags))
	} else {
		return db.db.Query(sql_str)
	}
}

// if last_seen_image_id_or_ser_value is negative then it is ignored
// If tags is empty then gets images sorted by upload date/time
func GetImages(db *DB, last_seen_image_id_or_ser_value int32, tags_ []string) ([]ImageID, int32, error) {
	res := []ImageID{}

	tags := map[string]bool{}
	excluded_tags := map[string]bool{}

	for i := range tags_ {
		if tag, valid, is_excluded := validate_tag(tags_[i]); valid {
			if is_excluded {
				excluded_tags[tag] = true
			} else {
				tags[tag] = true
			}
		}
	}

	// Get results

	var rows *sql.Rows
	var err error

	if len(excluded_tags) == 0 && len(tags) == 0 {
		// Get most recent images with any tags

		if last_seen_image_id_or_ser_value < 0 {
			rows, err = get_statement(db,
				"SELECT id, id as cmp_to_val FROM image ORDER BY id DESC LIMIT 48;").Query()
		} else {
			rows, err = get_statement(db,
				`SELECT id, id as cmp_to_val FROM image 
					WHERE id < $1 
					ORDER BY id DESC LIMIT 48;`).Query(int32(last_seen_image_id_or_ser_value))
		}
	} else if len(excluded_tags) == 0 && len(tags) == 1 {
		// Get images with given tag, in order that this tag was added to the images

		for tag := range tags {
			// ^ Gets the tag from the map

			if last_seen_image_id_or_ser_value < 0 {
				rows, err = get_statement(db,
					`SELECT image_id, ser as cmp_to_val FROM image_tag
				WHERE (deleted IS FALSE) AND
				(tag_id = (SELECT id FROM tag WHERE name=$1))
				ORDER BY ser DESC LIMIT 48;`).Query(tag)
			} else {
				rows, err = get_statement(db,
					`SELECT image_id, ser as cmp_to_val FROM image_tag
				WHERE (deleted IS FALSE) AND
				(tag_id = (SELECT id FROM tag WHERE name=$1)) AND
				(ser < $2)
				ORDER BY ser DESC LIMIT 48;`).Query(tag, last_seen_image_id_or_ser_value)
			}
		}
	} else {
		tags_list := []string{}
		for t := range tags {
			tags_list = append(tags_list, t)
		}
		excluded_tags_list := []string{}
		for t := range excluded_tags {
			excluded_tags_list = append(excluded_tags_list, t)
		}
		rows, err = full_search(db, last_seen_image_id_or_ser_value, tags_list, excluded_tags_list)
	}

	if err != nil {
		return res, 0, err
	}
	defer rows.Close()

	// Move results into ImageSearchResults type

	lowest_image_id_or_ser_value := int32(2147483647)

	for rows.Next() {
		var id int32
		var cmp_to_val int32
		err = rows.Scan(&id, &cmp_to_val)
		if err != nil {
			log.Println(err)
		} else {
			res = append(res, ImageID(id))

			if cmp_to_val < lowest_image_id_or_ser_value {
				lowest_image_id_or_ser_value = cmp_to_val
			}
		}
	}

	return res, lowest_image_id_or_ser_value, nil
}
