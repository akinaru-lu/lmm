package db

const createUser = `
CREATE TABLE IF NOT EXISTS user (
	id int unsigned NOT NULL AUTO_INCREMENT,
	name varchar(32) NOT NULL UNIQUE,
	password varchar(128) NOT NULL,
	guid varchar(36) NOT NULL UNIQUE,
	token varchar(36) NOT NULL UNIQUE,
	created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (id)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8;
`

const createBlog = `
CREATE TABLE IF NOT EXISTS blog (
	id int unsigned NOT NULL AUTO_INCREMENT,
	user int unsigned NOT NULL,
	title varchar(255) NOT NULL,
	text text NOT NULL,
	created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	PRIMARY KEY (id)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8;
`

const createCategory = `
CREATE TABLE IF NOT EXISTS category (
	id int unsigned NOT NULL AUTO_INCREMENT,
	user int unsigned NOT NULL,
	blog int unsigned NOT NULL,
	name varchar(32) NOT NULL,
	PRIMARY KEY (id),
	UNIQUE (user, blog)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8;
`

const createTag = `
CREATE TABLE IF NOT EXISTS tag (
	id int unsigned NOT NULL AUTO_INCREMENT,
	user int unsigned NOT NULL,
	blog int unsigned NOT NULL,
	name varchar(32) NOT NULL,
	PRIMARY KEY (id),
	UNIQUE (user, blog, name)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8;
`

const createImage = `
CREATE TABLE IF NOT EXISTS image (
	id int unsigned NOT NULL AUTO_INCREMENT,
	user int unsigned NOT NULL,
	type tinyint NOT NULL,
	url varchar(127) NOT NULL,
	created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (id),
	UNIQUE (url)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8;
`

const createProject = `
CREATE TABLE IF NOT EXISTS project (
	id int unsigned NOT NULL AUTO_INCREMENT,
	user int unsigned NOT NULL,
	icon varchar(255) NOT NULL DEFAULT "",
	name varchar(63) NOT NULL,
	url varchar(255) NOT NULL DEFAULT "",
	description varchar(1023) NOT NULL DEFAULT "",
	created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	from_date date,
	to_date date,
	PRIMARY KEY (id),
	UNIQUE (user, name)
)
`

var CreateSQL = []string{
	createUser,
	createBlog,
	createCategory,
	createTag,
	createImage,
	createProject,
}
