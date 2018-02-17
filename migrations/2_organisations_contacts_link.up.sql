ALTER TABLE contacts ADD COLUMN organisation_id BIGINT UNSIGNED;
UPDATE contacts SET organisation_id=1 WHERE organisation_id IS NULL;
ALTER TABLE contacts MODIFY COLUMN organisation_id BIGINT UNSIGNED NOT NULL DEFAULT 1;
ALTER TABLE contacts ADD CONSTRAINT fk_contacts_1 FOREIGN KEY (organisation_id) REFERENCES organisations (id);