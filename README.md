### sql

	CREATE TABLE Todo (
		id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
		created int(10),
		status ENUM('todo', 'doing', 'done'),
		title VARCHAR(255),
		description TEXT
	) ENGINE=InnoDB DEFAULT CHARSET=utf8;
