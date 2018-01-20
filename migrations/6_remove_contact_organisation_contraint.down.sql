ALTER TABLE contacts ADD CONSTRAINT contacts_organisation_id_fkey FOREIGN KEY (organisation_id) REFERENCES public.organisations (id);
ALTER TABLE contacts ALTER COLUMN organisation_id SET DEFAULT 1;