
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(50) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL
	);

	CREATE TABLE IF NOT EXISTS notes (
		id SERIAL PRIMARY KEY,
		text TEXT NOT NULL,
		user_id INT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);
    insert into users (username, password) values ('admin', '$2a$12$aRVYd0RtCP6/qGPa8Oc0Nunlr/V8DnqMRyXTXCmtfvGoD3ek9xfOe');