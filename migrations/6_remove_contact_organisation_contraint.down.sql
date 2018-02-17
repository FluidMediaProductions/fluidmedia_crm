ALTER TABLE contacts ADD CONSTRAINT fk_contacts_1 FOREIGN KEY (organisation_id) REFERENCES organisations (id);
ALTER TABLE contacts ALTER COLUMN organisation_id SET DEFAULT 1;