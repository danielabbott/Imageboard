
-- Browse by tag
SELECT image_id FROM image_tag WHERE (deleted IS FALSE) AND (tag_id = 3) AND (ser < 123) 
    ORDER BY ser DESC LIMIT 48;
SELECT image_id FROM image_tag WHERE (deleted IS FALSE) AND 
    (tag_id = (SELECT id FROM tag WHERE name='1'))
    ORDER BY ser DESC LIMIT 48;
-- ser puts the images in the order they were tagged


-- Search

-- Get id & count of all tags. This is done twice, once for tags we want, once for tags we exclude
-- This is for choosing the tag order
SELECT tag_id FROM image_tag WHERE tag_id IN 
    (SELECT id FROM tag WHERE name IN ('abc', 'def'))
    AND deleted=FALSE
    GROUP BY tag_id
    ORDER BY COUNT(*) DESC;


SELECT COUNT(*) as image_count FROM image;


-- A and B

SELECT A.image_id, (CASE WHEN A.ser > B.ser THEN A.ser ELSE B.ser END) as ser
FROM image_tag A
INNER JOIN image_tag B
ON A.image_id = B.image_id
WHERE (A.deleted = FALSE) AND (A.tag_id = 1) AND (A.ser < 123)
    AND (B.deleted = FALSE) AND (B.tag_id = 2) AND (B.ser < 123);



-- A but not B

SELECT image_id, ser FROM image_tag
WHERE (deleted = FALSE) AND (tag_id = 1) AND (ser < 123)
    AND (image_id NOT IN (SELECT image_id from image_tag WHERE tag_id = 2));


-- Included tags all come before the excluded tags

-- A and B and C but not D or E

SELECT image_id, ser FROM 
    (SELECT image_id, ser FROM
        (SELECT image_id, ser FROM
            (SELECT A.image_id, (CASE WHEN A.ser > B.ser THEN A.ser ELSE B.ser END) as ser FROM
            (
                SELECT A.image_id, (CASE WHEN A.ser > B.ser THEN A.ser ELSE B.ser END) as ser FROM
                (SELECT image_id,ser FROM image_tag WHERE (deleted = FALSE) AND (tag_id = 1)) A
                INNER JOIN
                (SELECT image_id,ser FROM image_tag WHERE (deleted = FALSE) AND (tag_id = 2)) B
                ON A.image_id = B.image_id
            ) A
            INNER JOIN
            (SELECT image_id,ser FROM image_tag WHERE (deleted = FALSE) AND (tag_id = 3)) B
            ON A.image_id = B.image_id) A
        WHERE A.image_id NOT IN (SELECT image_id from image_tag WHERE tag_id = 4)) A
    WHERE A.image_id NOT IN (SELECT image_id from image_tag WHERE tag_id = 5)) A
WHERE A.ser < 123
ORDER BY A.ser DESC LIMIT 48;





-- Browse all
SELECT id FROM image WHERE id < 123 ORDER BY id DESC LIMIT 10;

-- -- Browse by user
-- SELECT id FROM image WHERE (poster_id = 1) AND (id < 123) 
--     ORDER BY id DESC LIMIT 10;

-- Get tags from image
SELECT tag.name from tag LEFT JOIN image_tag on tag.id = image_tag.tag_id
    WHERE (image_tag.deleted IS FALSE) AND image_tag.image_id = 1 ORDER BY image_tag.ser ASC;


-- For given image id, get username and title
SELECT useracc.name, image.title FROM image
    INNER JOIN useracc ON useracc.id = image.poster_id
    WHERE image.id = 1;


-- -- For given image id, get 10 comments
-- SELECT comment.id, comment.comment_text, comment.poster_id, useracc.name FROM comment
--     INNER JOIN useracc ON useracc.id = comment.poster_id
--     WHERE (comment.image_id = 1) AND (comment.id < 123)
--     ORDER BY comment.id DESC LIMIT 10;

-- -- Post comment
-- INSERT INTO comment(image_id,poster_id,comment_text) VALUES(1,1,'hello');


-- Ensure tag exists
INSERT INTO tag(name) VALUES('mytag') ON CONFLICT(name) DO NOTHING;

-- Add tag to image
INSERT INTO image_tag(image_id, tag_id) VALUES(1, 1);

-- Undelete tag on image if previous query failed because the record already exists
UPDATE image_tag SET deleted = FALSE WHERE tag_id=0 AND image_id=0

-- Remove tag from image (2 ways)
UPDATE image_tag SET deleted = True WHERE id = 0
UPDATE image_tag SET deleted = True WHERE image_id = 0 AND tag_id = 0


-- -- Check for file hash duplication
-- SELECT COUNT(*) FROM image WHERE file_hash = '\x33';


-- -- Add user
-- INSERT INTO useracc(name,password) VALUES('Joe', '\x00');

-- -- Login
-- SELECT id,password FROM useracc WHERE lower(name) = lower('joe');
-- INSERT INTO session(id,user_id) VALUES('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 1);

-- -- Get session
-- SELECT timeout,user_id FROM session WHERE id = 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11';


-- -- Clear sessions (runs once per day during time of minimum load)
-- DELETE FROM session WHERE timeout < now();

-- Clear unused tags
DELETE FROM tag WHERE id NOT IN 
    (SELECT tag_id FROM image_tag WHERE deleted=FALSE);