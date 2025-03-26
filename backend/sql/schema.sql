CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE questions (
    question_id SERIAL PRIMARY KEY,
    question_text TEXT NOT NULL,
    choice1_text VARCHAR(100) NOT NULL,
    choice2_text VARCHAR(100) NOT NULL,
    created_by INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_user
      FOREIGN KEY (created_by)
      REFERENCES users(user_id)
      ON DELETE CASCADE
);

CREATE TABLE question_archives (
    archive_id SERIAL PRIMARY KEY,
    question_id INT NOT NULL,
    archive_date DATE NOT NULL,         
    choice1_votes INT NOT NULL,
    choice2_votes INT NOT NULL,
    total_votes INT NOT NULL,
    was_most_popular BOOLEAN DEFAULT FALSE,
    archived_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_question
      FOREIGN KEY (question_id)
      REFERENCES questions(question_id)
      ON DELETE CASCADE
);

CREATE INDEX idx_questions_created_by ON questions (created_by);
CREATE INDEX idx_question_archives_qid ON question_archives (question_id);
