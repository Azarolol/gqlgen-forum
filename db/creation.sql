CREATE TABLE IF NOT EXISTS Posts (
	id int PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
	author varchar(100) NOT NULL,
	header varchar(200) NOT NULL,
	body varchar(2000) NOT NULL,
	allows_comments bool NOT NULL DEFAULT false
);

CREATE TABLE IF NOT EXISTS Comments (
	id int PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
	on_post bool NOT NULL DEFAULT false,
	comment_on int NOT NULL,
	author varchar(100) NOT NULL,
	body varchar(2000) NOT NULL
);