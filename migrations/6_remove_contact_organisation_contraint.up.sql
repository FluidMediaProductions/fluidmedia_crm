ALTER TABLE contacts DROP FOREIGN KEY fk_contacts_1;
ALTER TABLE contacts ALTER COLUMN organisation_id SET DEFAULT 0;