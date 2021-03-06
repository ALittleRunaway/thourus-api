package gateway

import (
	"context"
	"database/sql"
	"thourus-api/domain/entity"
)

const (
	dateFormat = "Mon, 02 Jan 2006 15:04:05"
)

type DocumentGw interface {
	GetDocumentById(ctx context.Context, documentId int) (entity.Document, error)
	GetDocumentByUid(ctx context.Context, documentUid string) (entity.Document, error)
	CreateDocument(ctx context.Context, documentName string, documentPath string, creatorId int, projectId int) (int, error)
	AddDocumentHistory(ctx context.Context, documentId int, hash string, PoW int64, initiatorIdd int) (int, error)
	RenameDocumentById(ctx context.Context, newDocumentName string, documentId int) error
	RenameDocumentByUid(ctx context.Context, newProjectName string, documentUid string) error
	DeleteDocumentById(ctx context.Context, documentId int) error
	DeleteDocumentByUid(ctx context.Context, documentUid string) error
	GetDocumentHistory(ctx context.Context, documentId int) ([]entity.DocumentHistory, error)
	ChangeDocumentStatusToPending(ctx context.Context, documentId int) error
}

type DocumentGateway struct {
	db *sql.DB
}

func NewDocumentGateway(db *sql.DB) *DocumentGateway {
	return &DocumentGateway{
		db: db,
	}
}

func (gw *DocumentGateway) GetDocumentById(ctx context.Context, documentId int) (entity.Document, error) {

	const query = `
	SELECT d.id, d.uid, d.name, d.path, u.uid,
		   p.id, p.uid, p.name,
		   s.id, s.uid, s.name
	FROM thourus.document d
	INNER JOIN thourus.project p ON d.project_id = p.id
	INNER JOIN thourus.space s ON p.space_id = s.id
	INNER JOIN thourus.user u ON d.creator_id = u.id
	WHERE d.id = ?;
`
	document := entity.Document{}

	rows, err := gw.db.QueryContext(ctx, query, documentId)

	for rows.Next() {
		if err = rows.Scan(
			&document.Id,
			&document.Uid,
			&document.Name,
			&document.Path,
			&document.Creator.Uid,
			&document.Project.Id,
			&document.Project.Uid,
			&document.Project.Name,
			&document.Space.Id,
			&document.Space.Uid,
			&document.Space.Name,
		); err != nil {
			return document, err
		}
	}

	return document, nil
}

func (gw *DocumentGateway) GetDocumentByUid(ctx context.Context, documentUid string) (entity.Document, error) {

	const query = `
	SELECT d.id, d.uid, d.name, d.path, u.uid,
		   p.id, p.uid, p.name,
		   s.id, s.uid, s.name
	FROM thourus.document d
	INNER JOIN thourus.project p ON d.project_id = p.id
	INNER JOIN thourus.space s ON p.space_id = s.id
	INNER JOIN thourus.user u ON d.creator_id = u.id
	WHERE d.uid = ?;
`
	document := entity.Document{}

	rows, err := gw.db.QueryContext(ctx, query, documentUid)

	for rows.Next() {
		if err = rows.Scan(
			&document.Id,
			&document.Uid,
			&document.Name,
			&document.Path,
			&document.Creator.Uid,
			&document.Project.Id,
			&document.Project.Uid,
			&document.Project.Name,
			&document.Space.Id,
			&document.Space.Uid,
			&document.Space.Name,
		); err != nil {
			return document, err
		}
	}

	return document, nil
}

func (gw *DocumentGateway) CreateDocument(ctx context.Context, documentName string, documentPath string, creatorId int, projectId int) (int, error) {

	const query = `
	INSERT INTO thourus.document (name, path, creator_id, status_id, project_id)
	VALUES (?, ?, ?, 1, ?);
`
	var createdDocumentId int

	res, err := gw.db.ExecContext(ctx, query, documentName, documentPath, creatorId, projectId)
	if err != nil {
		return createdDocumentId, err
	}

	createdProjectId64, err := res.LastInsertId()
	if err != nil {
		return createdDocumentId, err
	}
	createdDocumentId = int(createdProjectId64)

	return createdDocumentId, nil
}

func (gw *DocumentGateway) AddDocumentHistory(ctx context.Context, documentId int, hash string, PoW int64, initiatorIdd int) (int, error) {

	const query = `
	INSERT INTO thourus.history (document_id, hash, pow, initiator_id)
	VALUES (?, ?, ?, ?);
`
	var createdDocumentHistoryId int

	res, err := gw.db.ExecContext(ctx, query, documentId, hash, PoW, initiatorIdd)
	if err != nil {
		return createdDocumentHistoryId, err
	}

	createdProjectId64, err := res.LastInsertId()
	if err != nil {
		return createdDocumentHistoryId, err
	}
	createdDocumentHistoryId = int(createdProjectId64)

	return createdDocumentHistoryId, nil
}

func (gw *DocumentGateway) GetDocumentHistory(ctx context.Context, documentId int) ([]entity.DocumentHistory, error) {

	const query = `
	SELECT h.id, h.uid, h.document_id, d.uid, d.name, d.status_id, 
		   s.name, h.hash, h.pow, h.initiator_id, u.name, u.surname, h.date_changed
	FROM thourus.history h
	INNER JOIN thourus.user u ON h.initiator_id = u.id
	INNER JOIN thourus.document d ON h.document_id = d.id
	INNER JOIN thourus.status s ON d.status_id = s.id
	WHERE d.id = ?;
`
	historyRows := []entity.DocumentHistory{}

	rows, err := gw.db.QueryContext(ctx, query, documentId)

	for rows.Next() {
		history := entity.DocumentHistory{}
		if err = rows.Scan(
			&history.Id,
			&history.Uid,
			&history.Document.Id,
			&history.Document.Uid,
			&history.Document.Name,
			&history.Document.Status.Id,
			&history.Document.Status.Name,
			&history.Hash,
			&history.PoW,
			&history.Initiator.Id,
			&history.Initiator.Name,
			&history.Initiator.Surname,
			&history.DateChanged,
		); err != nil {
			return historyRows, err
		}

		history.DateChangedString = history.DateChanged.Format(dateFormat)
		historyRows = append(historyRows, history)
	}

	return historyRows, nil
}

func (gw *DocumentGateway) RenameDocumentById(ctx context.Context, newDocumentName string, documentId int) error {

	const query = `
	UPDATE thourus.document d SET d.name = ? WHERE d.id = ?;
`
	_, err := gw.db.ExecContext(ctx, query, newDocumentName, documentId)
	if err != nil {
		return err
	}

	return nil
}

func (gw *DocumentGateway) ChangeDocumentStatusToPending(ctx context.Context, documentId int) error {

	const query = `
	UPDATE thourus.document d SET d.status_id = 2 WHERE d.id = ?;
`
	_, err := gw.db.ExecContext(ctx, query, documentId)
	if err != nil {
		return err
	}

	return nil
}

func (gw *DocumentGateway) RenameDocumentByUid(ctx context.Context, newProjectName string, documentUid string) error {

	const query = `
	UPDATE thourus.document d SET d.name = ? WHERE d.uid = ?;
`
	_, err := gw.db.ExecContext(ctx, query, newProjectName, documentUid)
	if err != nil {
		return err
	}

	return nil
}

func (gw *DocumentGateway) DeleteDocumentById(ctx context.Context, documentId int) error {

	const query = `
	DELETE FROM thourus.document d WHERE d.id = ?;
`
	_, err := gw.db.ExecContext(ctx, query, documentId)
	if err != nil {
		return err
	}

	return nil
}

func (gw *DocumentGateway) DeleteDocumentByUid(ctx context.Context, documentUid string) error {

	const query = `
	DELETE FROM thourus.history h WHERE h.document_id = 
	(SELECT d.id FROM thourus.document d WHERE d.uid = ?);
`
	_, err := gw.db.ExecContext(ctx, query, documentUid)
	if err != nil {
		return err
	}

	const query2 = `
	DELETE FROM thourus.document d WHERE d.uid = ?;
`
	_, err = gw.db.ExecContext(ctx, query2, documentUid)
	if err != nil {
		return err
	}

	return nil
}
