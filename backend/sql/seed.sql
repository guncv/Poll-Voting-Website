-- ================================
-- 1. Insert Mock Users
-- ================================
INSERT INTO users (username, email, password)
VALUES
('alice', 'alice@example.com', 'hashed_password_alice'),
('bob',   'bob@example.com',   'hashed_password_bob'),
('admin', 'admin@example.com', 'hashed_password_admin');

-- ================================
-- 2. Insert Mock Questions
-- ================================
INSERT INTO questions (question_text, choice1_text, choice2_text, created_by)
VALUES
('Do you pour cereal first or milk first?', 'Cereal first', 'Milk first', 1),
('Do you like pineapple on pizza?', 'Yes', 'No', 2),
('Have you paid your taxes yet?', 'Yes', 'No', 3),
('Would you kill a child for 1 million dollars?', 'Yes', 'No', 2);

-- ================================
-- 3. Insert Sample Archived Results
--    (Pretend these were top-voted in previous days)
-- ================================
INSERT INTO question_archives (question_id, archive_date, choice1_votes, choice2_votes, total_votes, was_most_popular)
VALUES
-- Example 1: The cereal question from yesterday
(1, '2025-03-26', 30, 25, 55, TRUE),

-- Example 2: The pineapple on pizza question from two days ago
(2, '2025-03-25', 40, 60, 100, TRUE),

-- Example 3: The "taxes" question was archived the same day but wasn't the top
(3, '2025-03-25', 20, 10, 30, FALSE);
