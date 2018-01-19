package db

import "github.com/fluidmediaproductions/fluidmedia_crm/model"

func (p *pgDb) prepareOrganisationsSqlStatements() (err error) {
	if p.sqlSelectOrganisations, err = p.dbConn.Preparex("SELECT * FROM organisations ORDER BY id"); err != nil { return err}
	if p.sqlSelectOrganisation, err = p.dbConn.Preparex("SELECT * FROM organisations WHERE id=$1"); err != nil { return err }
	if p.sqlUpdateOrganisation, err = p.dbConn.PrepareNamed("UPDATE organisations SET name=:name, email=:email," +
		" image=:image, phone=:phone, website=:website, twitter=:twitter, address=:address," +
		" description=:description WHERE id=:id"); err != nil { return err }
	if p.sqlInsertOrganisation, err = p.dbConn.PrepareNamed("INSERT INTO organisations (name, email, image, phone," +
		" website, twitter, address, description)" +
		" VALUES (:name, :email, :image, :phone, :website, :twitter, :address, :description) RETURNING id"); err != nil { return err }
	if p.sqlDeleteOrganisation, err = p.dbConn.Preparex("DELETE FROM organisations WHERE id=$1"); err != nil { return err }
    return nil
}

func (p *pgDb) SelectOrganisations() ([]*model.Organisation, error){
	organisation := make([]*model.Organisation, 0)
	if err := p.sqlSelectOrganisations.Select(&organisation); err != nil {
		return nil, err
	}
	return organisation, nil
}

func (p *pgDb) SelectOrganisation(id int) (*model.Organisation, error){
	var organisation model.Organisation
	if err := p.sqlSelectOrganisation.Get(&organisation, id); err != nil {
		return nil, err
	}
	return &organisation, nil
}

func (p *pgDb) UpdateOrganisation(organisation *model.Organisation) error {
	_, err := p.sqlUpdateOrganisation.Exec(organisation)
	return err
}

func (p *pgDb) NewOrganisation() (int, error) {
	id := 0
	err := p.sqlInsertOrganisation.QueryRow(&model.Organisation{}).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (p *pgDb) DeleteOrganisation(id int) error {
	_, err := p.sqlDeleteOrganisation.Exec(id)
	return err
}