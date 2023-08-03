CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL ,
    markdown TEXT NOT NULL
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY ,
    email TEXT UNIQUE  NOT NULL ,
    display_name TEXT NOT NULL
);

CREATE TABLE comments (
    id SERIAL PRIMARY KEY ,
    user_id INT,
    post_id INT,
    markdown TEXT NOT NULL
);