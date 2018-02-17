package db

import "github.com/fluidmediaproductions/fluidmedia_crm/model"

func (p *pgDb) SelectOrganisations() ([]*model.Organisation, error) {
	organisation := make([]*model.Organisation, 0)
	if err := p.dbConn.Select(&organisation, "SELECT * FROM organisations ORDER BY id"); err != nil {
		return nil, err
	}
	return organisation, nil
}

func (p *pgDb) SelectOrganisation(id int) (*model.Organisation, error) {
	var organisation model.Organisation
	if err := p.dbConn.Get(&organisation, "SELECT * FROM organisations WHERE id=?", id); err != nil {
		return nil, err
	}
	return &organisation, nil
}

func (p *pgDb) UpdateOrganisation(organisation *model.Organisation) error {
	_, err := p.dbConn.NamedExec("UPDATE organisations SET name=:name, email=:email," +
		" image=:image, phone=:phone, website=:website, twitter=:twitter, youtube=:youtube, instagram=:instagram," +
		" facebook=:facebook, address=:address, description=:description WHERE id=:id", organisation)
	return err
}

func (p *pgDb) NewOrganisation() (int, error) {
	res, err := p.dbConn.NamedExec("INSERT INTO organisations (name, email, image, phone," +
		" website, twitter, youtube, instagram, facebook, address, description)" +
		" VALUES (:name, :email, :image, :phone, :website, :twitter, :youtube, :instagram, :facebook, :address," +
		" :description)", &model.Organisation{})
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (p *pgDb) DeleteOrganisation(id int) error {
	_, err := p.dbConn.Exec("DELETE FROM organisations WHERE id=?", id)
	return err
}
