DROP DATABASE IF EXISTS blog;
DROP TABLE IF EXISTS blog.posts;

CREATE DATABASE blog;

CREATE TABLE blog.posts (
  id SERIAL PRIMARY KEY,
  title TEXT NOT NULL,
  postdate DATETIME DEFAULT CURRENT_TIMESTAMP,
  link TEXT NOT NULL,
  content TEXT NOT NULL,
  UNIQUE INDEX `id_UNIQUE` (id ASC) VISIBLE);

INSERT INTO blog.posts (title, postdate, link, content)
VALUES ("Title1", '2020-21-02', "https://google/link1", "Test content1");

INSERT INTO blog.posts (title, postdate, link, content)
VALUES ("Title2", '2020-22-02', "https://google/link2", "Test content2");