ALTER TABLE tag
ADD COLUMN `last_used` DATETIME after `name`,
ADD COLUMN `last_washing` DATETIME after `last_used`;
