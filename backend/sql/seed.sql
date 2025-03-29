-- ================================
-- 1) Insert Mock Users
-- ================================
INSERT INTO users (user_id, email, password)
VALUES
  ('11111111-1111-1111-1111-111111111111', 'alice@example.com',  'hashed_password_alice'),
  ('22222222-2222-2222-2222-222222222222', 'bob@example.com',    'hashed_password_bob'),
  ('33333333-3333-3333-3333-333333333333', 'charlie@example.com','hashed_password_charlie');

-- ================================
-- 2) Insert Mock Popular Questions
-- ================================
INSERT INTO questions (
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
  '11111111-1111-1111-1111-111111111111'
);

INSERT INTO questions (
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
  '22222222-2222-2222-2222-222222222222'
);

INSERT INTO questions (
  archive_date,
  question_text,
  yes_votes,
  no_votes,
  total_votes,
  created_by
)
VALUES (
  CURRENT_DATE,
  'Would you pay for ChatGPT?',
  60,
  10,
  70,
  '33333333-3333-3333-3333-333333333333'
);
