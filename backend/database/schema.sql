CREATE TABLE IF NOT EXISTS tblUsers (
  id INTEGER PRIMARY KEY,
  username TEXT UNIQUE,
  email TEXT UNIQUE,
  user_password TEXT,
  joined_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS tblPosts (
  id INTEGER PRIMARY KEY,
  user_id INTEGER NOT NULL,
  title TEXT,
  body TEXT,
  parent_id INTEGER,
  created_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  post_type TEXT CHECK (post_type IN ('post', 'comment', 'reply')),
  FOREIGN KEY (user_id) REFERENCES tblUsers (id),
  FOREIGN KEY (parent_id) REFERENCES tblPosts (id)
);
CREATE INDEX IF NOT EXISTS idx_user_posts ON tblPosts (user_id);
CREATE INDEX IF NOT EXISTS idx_parent_posts ON tblPosts (parent_id);
CREATE TABLE IF NOT EXISTS tblFiles (
  id INTEGER PRIMARY KEY,
  post_id INTEGER NOT NULL,
  file_path TEXT UNIQUE,
  file_type TEXT,
  FOREIGN KEY (post_id) REFERENCES tblPosts (id)
);
CREATE TABLE IF NOT EXISTS tblReactions (
  id INTEGER PRIMARY KEY,
  reaction TEXT CHECK (reaction IN ('like', 'dislike')),
  user_id INTEGER NOT NULL,
  post_id INTEGER,
  comment_id INTEGER,
  FOREIGN KEY (user_id) REFERENCES tblUsers (id),
  FOREIGN KEY (post_id) REFERENCES tblPosts (id),
  FOREIGN KEY (comment_id) REFERENCES tblPosts (id)
);
CREATE TABLE IF NOT EXISTS tblCategories (
  id INTEGER PRIMARY KEY,
  name TEXT UNIQUE NOT NULL
);
CREATE TABLE IF NOT EXISTS tblPostCategories (
  post_id INTEGER NOT NULL,
  category_id INTEGER NOT NULL,
  PRIMARY KEY (post_id, category_id),
  FOREIGN KEY (post_id) REFERENCES tblPosts (id),
  FOREIGN KEY (category_id) REFERENCES tblCategories (id)
);
CREATE TABLE IF NOT EXISTS tblSessions (
  id INTEGER PRIMARY KEY,
  user_id INTEGER NOT NULL,
  session_token TEXT NOT NULL UNIQUE,
  expires_at TIMESTAMP NOT NULL,
  FOREIGN KEY (user_id) REFERENCES tblUsers (id)
);
CREATE TABLE IF NOT EXISTS tblComments (
  id INTEGER PRIMARY KEY,
  post_id INTEGER NOT NULL,
  user_id INTEGER NOT NULL,
  body TEXT NOT NULL,
  created_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (post_id) REFERENCES tblPosts (id),
  FOREIGN KEY (user_id) REFERENCES tblUsers (id)
);
CREATE INDEX IF NOT EXISTS idx_category_id ON tblPostCategories (category_id);
CREATE TABLE IF NOT EXISTS tblMedia (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  filename TEXT NOT NULL,
  file_path TEXT NOT NULL,
  file_type TEXT NOT NULL,
  uploader_id INTEGER NOT NULL,
  uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (uploader_id) REFERENCES tblUsers(id) ON DELETE CASCADE
);