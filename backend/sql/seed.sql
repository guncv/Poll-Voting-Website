-- ================================
-- 1) Insert Mock Users
-- ================================
INSERT INTO users (email, password)
VALUES
-- Passwords here are placeholder text; in a real app you'd store hashed/salted passwords.
('alice@example.com',  'hashed_password_alice'),
('bob@example.com',    'hashed_password_bob'),
('charlie@example.com','hashed_password_charlie');

-- ================================
-- 2) Insert Mock Popular Questions
-- ================================
-- We'll pretend these were the top questions on different days.
-- We'll use CURRENT_DATE for a recent one, and offset days for older ones.

-- Example #1: Pineapple on pizza?
INSERT INTO popular_questions (
  archive_date,
  question_text,
  yes_votes,
  no_votes,
  total_votes,
  created_by
)
VALUES (
  '2025-03-24',
  'Do you like pineapple on pizza?',
  30,
  20,
  50,
  1      -- created_by user_id=1 (Alice)
);

-- Example #2: Morning person?
INSERT INTO popular_questions (
  archive_date,
  question_text,
  yes_votes,
  no_votes,
  total_votes,
  created_by
)
VALUES (
  '2025-03-25',
  'Are you a morning person?',
  45,
  15,
  60,
  2      -- created_by user_id=2 (Bob)
);

-- Example #3: Would you pay for ChatGPT?
INSERT INTO popular_questions (
  archive_date,
  question_text,
  yes_votes,
  no_votes,
  total_votes,
  created_by
)
VALUES (
  CURRENT_DATE,         -- for today
  'Would you pay for ChatGPT?',
  60,
  10,
  70,
  3      -- created_by user_id=3 (Charlie)
);
