CREATE TABLE exams (
  id SERIAL PRIMARY KEY,
  slug text UNIQUE,
  title text,
  description text,
  duration_minutes int,
  created_at timestamptz DEFAULT now(),
  created_by INT REFERENCES users(id),
  updated_at TIMESTAMPTZ DEFAULT now();
);


CREATE TABLE questions (
id SERIAL PRIMARY KEY,
exam_id int REFERENCES exams(id) ON DELETE CASCADE,
section text,
topic text,
year int,
type text,
question_text text,
choices jsonb,
correct_answer jsonb,
marks numeric DEFAULT 1,
negative_marks numeric DEFAULT 0,
resources jsonb,
created_at timestamptz DEFAULT now()
);


CREATE TABLE tests (
id SERIAL PRIMARY KEY,
exam_id int REFERENCES exams(id),
title text,
start_time timestamptz,
duration_minutes int,
is_live boolean DEFAULT false,
config jsonb,
created_at timestamptz DEFAULT now()
updated_at timestamptz DEFAULT now()
);


CREATE TABLE test_questions (
id SERIAL PRIMARY KEY,
test_id int REFERENCES tests(id) ON DELETE CASCADE,
question_id int REFERENCES questions(id),
sequence_order int,
marks_override numeric
);


CREATE TABLE users (
id SERIAL PRIMARY KEY,
name text,
email text UNIQUE,
phone text UNIQUE,
password text,
role text DEFAULT 'student',
created_at timestamptz DEFAULT now(),
updated_at timestamptz DEFAULT now()
);


CREATE TABLE submissions (
id SERIAL PRIMARY KEY,
test_id int REFERENCES tests(id),
user_id int REFERENCES users(id),
started_at timestamptz,
submitted_at timestamptz,
status text,
score numeric,
raw_result jsonb
);


CREATE TABLE answers (
id SERIAL PRIMARY KEY,
submission_id int REFERENCES submissions(id) ON DELETE CASCADE,
question_id int REFERENCES questions(id),
answer jsonb,
is_marked_for_review boolean DEFAULT false,
time_spent_seconds int,
created_at timestamptz DEFAULT now()
);

CREATE TABLE upload_jobs (
  id SERIAL PRIMARY KEY,
  file_name TEXT NOT NULL,
  status VARCHAR(50) DEFAULT 'pending',
  total_rows INTEGER DEFAULT 0,
  processed_rows INTEGER DEFAULT 0,
  errors TEXT,
  created_at TIMESTAMPTZ DEFAULT now(),
  updated_at TIMESTAMPTZ DEFAULT now()
);