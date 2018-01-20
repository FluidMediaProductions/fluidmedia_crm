ALTER TABLE contacts DROP CONSTRAINT contacts_organisation_id_fkey;
ALTER TABLE contacts ALTER COLUMN organisation_id SET DEFAULT 0;