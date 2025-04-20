PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS tblUsers (
  id INTEGER PRIMARY KEY,
  username TEXT UNIQUE NOT NULL, 
  email TEXT UNIQUE NOT NULL,
  user_password TEXT NULL,
  auth_provider TEXT NOT NULL DEFAULT '',
  joined_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tblPosts (
  id INTEGER PRIMARY KEY,
  user_id INTEGER NOT NULL,
  post_title TEXT,
  body TEXT,
  parent_id INTEGER DEFAULT NULL,
  created_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  post_status TEXT DEFAULT 'visible',
  media_url TEXT DEFAULT '', 
  FOREIGN KEY (user_id) REFERENCES tblUsers (id),
  FOREIGN KEY (parent_id) REFERENCES tblPosts (id)
);

CREATE TABLE IF NOT EXISTS tblPostCategories (
  id INTEGER PRIMARY KEY,
  post_id INTEGER NOT NULL,
  category TEXT,
  FOREIGN KEY (post_id) REFERENCES tblPosts (id)
);

CREATE TABLE IF NOT EXISTS tblReactions (
  id INTEGER PRIMARY KEY,
  reaction TEXT,
  reaction_status TEXT DEFAULT 'clicked',
  user_id INTEGER NOT NULL,
  post_id INTEGER NOT NULL,
  FOREIGN KEY (user_id) REFERENCES tblUsers (id),
  FOREIGN KEY (post_id) REFERENCES tblPosts (id)
);

CREATE TABLE IF NOT EXISTS tblSessions (
  id INTEGER PRIMARY KEY,
  user_id INTEGER NOT NULL,
  session_token TEXT NOT NULL UNIQUE,
  expires_at TIMESTAMP NULL,
  FOREIGN KEY (user_id) REFERENCES tblUsers (id)
);
