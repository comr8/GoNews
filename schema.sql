-- Создание таблицы
CREATE TABLE IF NOT EXISTS news (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    pubtime INTEGER DEFAULT 0,
    link TEXT NOT NULL
);

-- Очистка таблицы
-- TRUNCATE TABLE news RESTART IDENTITY;