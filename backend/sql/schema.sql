CREATE TABLE users (
    user_id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE questions (
    question_id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    archive_date DATE NOT NULL,
    question_text VARCHAR(255) NOT NULL,
    yes_votes INT NOT NULL,
    no_votes INT NOT NULL,
    total_votes INT NOT NULL,
    created_by uuid, 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_users 
      FOREIGN KEY (created_by) 
      REFERENCES users(user_id)
      ON DELETE CASCADE
);


-- If we want EXACTLY one top question per day we can add this
-- CREATE UNIQUE INDEX unique_top_question_per_day ON popular_questions (archive_date);
