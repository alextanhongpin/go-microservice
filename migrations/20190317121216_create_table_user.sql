-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS user (
	id BINARY(16),

	-- Email.
	email VARCHAR(255) NOT NULL UNIQUE,
	email_verified BOOLEAN NOT NULL DEFAULT 0, 

	-- Phone.
	phone_number VARCHAR(32) NOT NULL DEFAULT '',
	phone_number_verified BOOLEAN NOT NULL DEFAULT 0,

	-- Profile.
	name VARCHAR(255) NOT NULL DEFAULT '',
	family_name VARCHAR(255) NOT NULL DEFAULT '',
	given_name VARCHAR(255) NOT NULL DEFAULT '',
	middle_name VARCHAR(255) NOT NULL DEFAULT '',
	nickname VARCHAR(255) NOT NULL DEFAULT '',
	preferred_username VARCHAR(255) NOT NULL DEFAULT '',
	profile VARCHAR(2083) NOT NULL DEFAULT '' COMMENT "URL of the End-User's profile page",
	picture VARCHAR(255) NOT NULL DEFAULT '',
	website VARCHAR(2083) NOT NULL DEFAULT '' COMMENT "URL of the End-User's Web page or blog",
	gender CHAR(1) NOT NULL DEFAULT '' COMMENT "End-User's gender (m|f|o|x)",
	birthdate VARCHAR(10) NOT NULL DEFAULT '' COMMENT "End-User's birthday, represented as an ISO 8601:2004 [ISO8601?2004] YYYY-MM-DD format. The year MAY be 0000, indicating that it is omitted",
	zoneinfo VARCHAR(32) NOT NULL DEFAULT '' COMMENT "String from zoneinfo [zoneinfo] time zone database representing the End-User's time zone. For example, Europe/Paris or America/Los_Angeles",
	locale VARCHAR(35) NOT NULL DEFAULT '' COMMENT "End-User's locale, represented as a BCP47 [RFC5646] language tag. This is typically an ISO 639-1 Alpha-2 [ISO639?1] language code in lowercase and an ISO 3166-1 Alpha-2 [ISO3166?1] country code in uppercase, separated by a dash. For example, en-US or fr-CA. As a compatibility note, some implementations have used an underscore as the separator rather than a dash, for example, en_US",

	hashed_password VARCHAR(255) NOT NULL,

	-- Address.
	street_address VARCHAR(255) NOT NULL DEFAULT '' COMMENT 'Full street address component, which MAY include house number, street name, Post Office Box, and multi-line extended street address information. This field MAY contain multiple lines, separated by newlines. Newlines can be represented either as a carriage return/line feed pair ("\r\n") or as a single line feed character ("\n").',
	locality VARCHAR(58) NOT NULL DEFAULT '' COMMENT 'City or locality component.',
	region VARCHAR(56) NOT NULL DEFAULT '' COMMENT 'State, province, prefecture or region component.',
	postal_code VARCHAR(16) NOT NULL DEFAULT '' COMMENT 'Zip code or postal code component.',
	country VARCHAR(74) NOT NULL DEFAULT '' COMMENT 'Country name component.',

	-- Timestamps.
	created_at DATETIME NOT NULL DEFAULT NOW(),
	updated_at DATETIME NOT NULL DEFAULT NOW() ON UPDATE NOW(),
	deleted_at DATETIME DEFAULT NULL,
	PRIMARY KEY (id)
) ENGINE=InnoDB CHARACTER SET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS user;
