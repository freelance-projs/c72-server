ALTER TABLE lending
ADD COLUMN num_lending INT UNSIGNED NOT NULL DEFAULT 0 AFTER `department`;
