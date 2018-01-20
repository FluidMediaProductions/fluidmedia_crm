ALTER TABLE contacts ADD COLUMN organisation_id integer REFERENCES organisations (id);
UPDATE contacts SET organisation_id=1 WHERE organisation_id IS NULL;
ALTER TABLE contacts ALTER COLUMN organisation_id SET NOT NULL;