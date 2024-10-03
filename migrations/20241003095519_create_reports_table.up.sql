CREATE TABLE reports (
      id CHAR(36) NOT NULL,
      user_id CHAR(36) NOT NULL,
      report_text TEXT NOT NULL,
      created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
      PRIMARY KEY (id),
      FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
