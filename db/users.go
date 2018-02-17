package db

import "github.com/fluidmediaproductions/fluidmedia_crm/model"

func (p *pgDb) SelectUsers() ([]*model.User, error) {
	users := make([]*model.User, 0)
	if err := p.dbConn.Select(&users, "SELECT * FROM users ORDER BY id"); err != nil {
		return nil, err
	}
	return users, nil
}

func (p *pgDb) SelectUser(id int) (*model.User, error) {
	var user model.User
	if err := p.dbConn.Get(&user, "SELECT * FROM users WHERE id=?", id); err != nil {
		return nil, err
	}
	return &user, nil
}

func (p *pgDb) UpdateUser(user *model.User) error {
	_, err := p.dbConn.NamedExec("UPDATE users SET name=:name, email=:email," +
	" isAdmin=:isadmin, phone=:phone, disabled=:disabled, login=:login, totp_secret=:totp_secret WHERE id=:id", user)
	return err
}

func (p *pgDb) UpdateUserPass(user *model.User) error {
	_, err := p.dbConn.NamedExec("UPDATE users SET pass=:pass WHERE id=:id", user)
	return err
}

func (p *pgDb) NewUser() (int, error) {
	res, err := p.dbConn.NamedExec("INSERT INTO users (name, email, isAdmin, phone, disabled," +
		" login, pass, totp_secret)" +
		" VALUES (:name, :email, :isadmin, :phone, :disabled,:login, :pass, :totp_secret)", &model.User{})
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (p *pgDb) DeleteUser(id int) error {
	_, err := p.dbConn.Exec("DELETE FROM users WHERE id=?", id)
	return err
}
