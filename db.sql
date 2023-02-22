DROP TABLE IF EXISTS image,tag,image_tag,comment,useracc,session;

CREATE TABLE image(
	id SERIAL PRIMARY KEY,
	-- poster_id INT NOT NULL DEFAULT 1 REFERENCES useracc(id) ON DELETE SET DEFAULT,
	title VARCHAR(64) NOT NULL
	-- -- 32-byte blake3 hash, for detecting duplicate uploads (hash of original file)
	-- file_hash BYTEA NOT NULL
);
-- CREATE INDEX idx_image_poster_id ON image(poster_id);
-- CREATE INDEX idx_image_file_hash ON image(file_hash);

CREATE TABLE tag(
	id SERIAL PRIMARY KEY,
	name VARCHAR(80) UNIQUE NOT NULL
);

CREATE TABLE image_tag(
	image_id INT NOT NULL REFERENCES image(id) ON DELETE CASCADE,
	tag_id INT NOT NULL REFERENCES tag(id) ON DELETE CASCADE,
	-- for putting images in the order they were tagged
	ser SERIAL,
	deleted BOOLEAN DEFAULT FALSE NOT NULL,

	PRIMARY KEY(image_id, tag_id)
);
CREATE INDEX idx_image_tag ON image_tag(tag_id) WHERE deleted IS FALSE;
CREATE INDEX idx_image_tag_by_tag ON image_tag(tag_id, ser DESC) WHERE deleted IS FALSE;
CREATE INDEX idx_image_tag_by_image ON image_tag(image_id, ser ASC) WHERE deleted IS FALSE;

-- CREATE TABLE comment(
-- 	id SERIAL PRIMARY KEY,
-- 	image_id INT NOT NULL,
-- 	poster_id INT NOT NULL,
-- 	comment_text VARCHAR(1024) NOT NULL
-- );
-- CREATE INDEX idx_comment_by_image_id ON comment(image_id, id desc);

-- CREATE TABLE useracc(
-- 	id SERIAL PRIMARY KEY,
-- 	name VARCHAR(20) NOT NULL,
-- 	-- 1-byte hash type plus 16-byte argon2 hash plus 16-byte salt
-- 	password BYTEA NOT NULL
-- );
-- CREATE INDEX idx_useracc_lowercase_name ON useracc(lower(name));
-- -- Default user acquires images when accounts are deleted
-- -- Password would be chosen when the website is first set up
-- INSERT INTO useracc(id,name,password) VALUES(1,'Default User', '\x00')

-- CREATE TABLE session(
-- 	-- 128
-- 	id UUID PRIMARY KEY,
-- 	user_id INT NOT NULL,
-- 	timeout TIMESTAMPTZ DEFAULT (now() + interval '30 days') NOT NULL
-- );

