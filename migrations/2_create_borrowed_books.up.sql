CREATE TABLE IF NOT EXISTS borrowed_books(
    id varchar(36) NOT NULL PRIMARY KEY,
    book_id varchar(36) NOT NULL,
    user_id varchar(36) NOT NULL,
    FOREIGN KEY (book_id) REFERENCES books(bid),
    FOREIGN KEY (user_id) REFERENCES users(uid)
);