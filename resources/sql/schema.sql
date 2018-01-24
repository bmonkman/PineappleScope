
CREATE TABLE `temperature` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `created_date` DATE CURRENT_DATE,
    `session_id` INTEGER,
    `inner` REAL,
    `outer` REAL
    
);

CREATE TABLE `session` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `start_date` DATE,
    `name` TEXT,
    `notes` TEXT
);
