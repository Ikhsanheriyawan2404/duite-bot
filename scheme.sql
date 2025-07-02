CREATE TABLE users (
  id BIGSERIAL PRIMARY KEY,
  uuid CHAR(36) NOT NULL UNIQUE,
  chat_id BIGINT NOT NULL,
  name VARCHAR(100) NOT NULL,
  is_paid TIMESTAMP NULL,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE categories (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  type VARCHAR(10) NOT NULL CHECK (type IN ('INCOME', 'EXPENSE')),
  parent_id BIGINT REFERENCES categories(id),
  user_id BIGINT REFERENCES users(id),
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

INSERT INTO categories (name, type, parent_id, user_id) VALUES
('Gaji & Pendapatan', 'INCOME', NULL, NULL),
('Pendapatan Sampingan', 'INCOME', NULL, NULL),
('Investasi', 'INCOME', NULL, NULL),
('Pengeluaran Harian', 'EXPENSE', NULL, NULL),
('Tagihan & Cicilan', 'EXPENSE', NULL, NULL),
('Tabungan & Dana Darurat', 'EXPENSE', NULL, NULL);

CREATE TABLE transactions (
  id BIGSERIAL PRIMARY KEY,
  transaction_type VARCHAR(20) NOT NULL,
  amount NUMERIC(15,2) NOT NULL,
  category_id BIGINT REFERENCES categories(id),
  description TEXT,
  transaction_date TIMESTAMPTZ NOT NULL,
  chat_id BIGINT,
  original_text TEXT,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);
