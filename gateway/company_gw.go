package gateway

import (
	"context"
	"database/sql"
	"thourus-api/domain/entity"
)

type CompanyGw interface {
	GetCompanyById(ctx context.Context, id int) (entity.Company, error)
	GetCompanyByUid(ctx context.Context, uid string) (entity.Company, error)
	GetSpacesInCompany(ctx context.Context, companyUid string) ([]entity.Space, error)
	CreateCompany(ctx context.Context, name string) (int, error)
	DeleteCompanyById(ctx context.Context, id int) error
	DeleteCompanyByUid(ctx context.Context, uid string) error
}

type CompanyGateway struct {
	db *sql.DB
}

func NewCompanyGateway(db *sql.DB) *CompanyGateway {
	return &CompanyGateway{
		db: db,
	}
}

func (gw *CompanyGateway) GetCompanyById(ctx context.Context, id int) (entity.Company, error) {

	const query = `
	SELECT id, uid, name FROM thourus.company c WHERE c.id = ?;
`
	company := entity.Company{}

	rows, err := gw.db.QueryContext(ctx, query, id)

	for rows.Next() {
		if err = rows.Scan(
			&company.Id,
			&company.Uid,
			&company.Name,
		); err != nil {
			return company, err
		}
	}

	return company, nil
}

func (gw *CompanyGateway) GetCompanyByUid(ctx context.Context, uid string) (entity.Company, error) {

	const query = `
	SELECT id, uid, name FROM thourus.company c WHERE c.uid = ?;
`
	company := entity.Company{}

	rows, err := gw.db.QueryContext(ctx, query, uid)

	for rows.Next() {
		if err = rows.Scan(
			&company.Id,
			&company.Uid,
			&company.Name,
		); err != nil {
			return company, err
		}
	}

	return company, nil
}

func (gw *CompanyGateway) GetSpacesInCompany(ctx context.Context, companyUid string) ([]entity.Space, error) {

	const query = `
	SELECT s.id, s.uid, s.name, c.id, c.uid, c.name FROM thourus.space s
	INNER JOIN thourus.company c ON s.company_id = c.id
	WHERE c.uid = ?;
`
	spaces := []entity.Space{}

	rows, err := gw.db.QueryContext(ctx, query, companyUid)

	for rows.Next() {
		space := entity.Space{}
		if err = rows.Scan(
			&space.Id,
			&space.Uid,
			&space.Name,
			&space.Company.Id,
			&space.Company.Uid,
			&space.Company.Name,
		); err != nil {
			return spaces, err
		}
		spaces = append(spaces, space)
	}

	return spaces, nil
}

func (gw *CompanyGateway) CreateCompany(ctx context.Context, name string) (int, error) {

	const query = `
	INSERT INTO thourus.company (name) VALUES (?);
`
	var createdCompanyId int

	res, err := gw.db.ExecContext(ctx, query, name)
	if err != nil {
		return createdCompanyId, err
	}

	createdCompanyId64, err := res.LastInsertId()
	if err != nil {
		return createdCompanyId, err
	}
	createdCompanyId = int(createdCompanyId64)

	return createdCompanyId, nil
}

func (gw *CompanyGateway) DeleteCompanyById(ctx context.Context, id int) error {

	const query = `
	DELETE FROM thourus.company c WHERE c.id = ?;
`
	_, err := gw.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (gw *CompanyGateway) DeleteCompanyByUid(ctx context.Context, uid string) error {

	const query = `
	DELETE FROM thourus.company c WHERE c.uid = ?;
`
	_, err := gw.db.ExecContext(ctx, query, uid)
	if err != nil {
		return err
	}

	return nil
}
