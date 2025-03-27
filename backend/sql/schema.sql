CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE popular_questions (
    pop_q_id SERIAL PRIMARY KEY,
    archive_date DATE NOT NULL, -- date when the question was most popular
    question_text VARCHAR(255) NOT NULL,
    yes_votes INT NOT NULL,
    no_votes INT NOT NULL,
    total_votes INT NOT NULL,
    created_by INT, 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_users 
      FOREIGN KEY (created_by) 
      REFERENCES users(user_id)
      ON DELETE CASCADE
);

-- If we want EXACTLY one top question per day we can add this
-- CREATE UNIQUE INDEX unique_top_question_per_day ON popular_questions (archive_date);
