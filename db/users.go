package db

import "github.com/fluidmediaproductions/fluidmedia_crm/model"

func (p *pgDb) prepareUsersSqlStatements() (err error) {
	if p.sqlSelectUsers, err = p.dbConn.Preparex("SELECT * FROM users ORDER BY id"); err != nil {
		return err
	}
	if p.sqlSelectUser, err = p.dbConn.Preparex("SELECT * FROM users WHERE id=$1"); err != nil {
		return err
	}
	if p.sqlUpdateUser, err = p.dbConn.PrepareNamed("UPDATE users SET name=:name, email=:email," +
		" isAdmin=:isadmin, phone=:phone, disabled=:disabled, login=:login, totp_secret=:totp_secret WHERE id=:id"); err != nil {
		return err
	}
	if p.sqlUpdateUserPass, err = p.dbConn.PrepareNamed("UPDATE users SET pass=:pass WHERE id=:id"); err != nil {
		return err
	}
	if p.sqlInsertUser, err = p.dbConn.PrepareNamed("INSERT INTO users (name, email, isAdmin, phone, disabled," +
		" login, pass, totp_secret)" +
		" VALUES (:name, :email, :isadmin, :phone, :disabled,:login, :pass, :totp_secret)"); err != nil {
		return err
	}
	if p.sqlDeleteUser, err = p.dbConn.Preparex("DELETE FROM users WHERE id=$1"); err != nil {
		return err
	}
	return nil
}

func (p *pgDb) SelectUsers() ([]*model.User, error) {
	users := make([]*model.User, 0)
	if err := p.sqlSelectUsers.Select(&users); err != nil {
		return nil, err
	}
	return users, nil
}

func (p *pgDb) SelectUser(id int) (*model.User, error) {
	var user model.User
	if err := p.sqlSelectUser.Get(&user, id); err != nil {
		return nil, err
	}
	return &user, nil
}

func (p *pgDb) UpdateUser(user *model.User) error {
	_, err := p.sqlUpdateUser.Exec(user)
	return err
}

func (p *pgDb) UpdateUserPass(user *model.User) error {
	_, err := p.sqlUpdateUser.Exec(user)
	return err
}

func (p *pgDb) NewUser() (int, error) {
	res, err := p.sqlInsertUser.Exec(&model.User{})
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
	_, err := p.sqlDeleteUser.Exec(id)
	return err
}
