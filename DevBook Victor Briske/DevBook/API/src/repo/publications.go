package repo

import (
	"api/src/models"
	"database/sql"
)

type Pubs struct {
	db *sql.DB
}

func NewRepoPubs(db *sql.DB) *Pubs {
	return &Pubs{db}
}

func (repo Pubs) CreatePub(pub models.Publication) (uint64, error) {
	statement, erro := repo.db.Prepare("insert into publications (title, content, author_id) values (?, ?, ?)")

	if erro != nil {
		return 0, erro
	}
	defer statement.Close()

	result, err := statement.Exec(pub.Title, pub.Content, pub.AuthorID)
	if err != nil {
		return 0, err
	}

	LastInsertId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(LastInsertId), nil
}

func (repo Pubs) SearchByID(pubID uint64) (models.Publication, error) {
	tuples, err := repo.db.Query(
		"select p.*, u.nick from publications p inner join users u on u.id = p.author_id where p.id = ?",
		pubID,
	)
	if err != nil {
		return models.Publication{}, err
	}
	defer tuples.Close()

	var publication models.Publication
	if tuples.Next() {
		if err = tuples.Scan(
			&publication.ID,
			&publication.Title,
			&publication.Content,
			&publication.AuthorID,
			&publication.Likes,
			&publication.CreatedAt,
			&publication.AuthorNick,
		); err != nil {
			return models.Publication{}, err
		}
	}

	return publication, nil
}

func (repo Pubs) Search(userID uint64) ([]models.Publication, error) {
	tuples, err := repo.db.Query(`select distinct p.*, u.nick from publications p
	inner join users u on u.id = p.author_id
	inner join followers s on p.author_id = s.user_id where u.id = ? or s.follower_id = ?
	order by 1 desc`, userID, userID)

	if err != nil {
		return nil, err
	}

	defer tuples.Close()

	var publications []models.Publication

	for tuples.Next() {
		var publication models.Publication
		if err = tuples.Scan(
			&publication.ID,
			&publication.Title,
			&publication.Content,
			&publication.AuthorID,
			&publication.Likes,
			&publication.CreatedAt,
			&publication.AuthorNick,
		); err != nil {
			return nil, err
		}

		publications = append(publications, publication)
	}
	return publications, nil
}

func (repo Pubs) Update(pubID uint64, publication models.Publication) error {
	statement, err := repo.db.Prepare("UPDATE publications SET title = ?, content = ? WHERE id = ?;")
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(publication.Title, publication.Content, pubID); err != nil {
		return err
	}

	return nil
}

func (repo Pubs) Delete(pubID uint64) error {
	statement, err := repo.db.Prepare("delete from publications where id = ?")
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(pubID); err != nil {
		return err
	}

	return nil
}

func (repo Pubs) SearchByUser(userID uint64) ([]models.Publication, error) {
	tuples, err := repo.db.Query(`
	select p.*, u.nick from publications p
	join users u on u.id = p.author_id
	where p.author_id = ?`,
		userID)
	if err != nil {
		return nil, err
	}
	defer tuples.Close()

	var publications []models.Publication

	for tuples.Next() {
		var publication models.Publication
		if err = tuples.Scan(
			&publication.ID,
			&publication.Title,
			&publication.Content,
			&publication.AuthorID,
			&publication.Likes,
			&publication.CreatedAt,
			&publication.AuthorNick,
		); err != nil {
			return nil, err
		}

		publications = append(publications, publication)
	}
	return publications, nil
}

func (repo Pubs) Like(pubID uint64) error {
	statement, err := repo.db.Prepare("update publications set likes = likes + 1 where id = ?")
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(pubID); err != nil {
		return err
	}

	return nil
}

func (repo Pubs) Deslike(pubID uint64) error {
	statement, err := repo.db.Prepare("update publications set likes = case when likes > 0 then likes - 1 else likes end  where id = ?")
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(pubID); err != nil {
		return err
	}

	return nil
}
