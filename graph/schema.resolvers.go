// schema.resolvers.go

package graph

import (
	"context"
	"database/sql"
	"log"
	"strconv"

	"github.com/kilianp07/graphql/graph/model"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

type mutationResolver struct {
	*Resolver
}

func (r *mutationResolver) CreateGame(ctx context.Context, input model.GameInput) (*model.Game, error) {

	result, err := r.DB.Exec(`
		INSERT INTO Game (name, publicationDate, platform)
		VALUES (?, ?, ?)
	`, input.Name, input.PublicationDate, input.Platform)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	gameID, _ := result.LastInsertId()

	for _, editorID := range input.EditorIDs {
		_, err := r.DB.Exec(`
			INSERT INTO GameEditor (gameID, editorID)
			VALUES (?, ?)
		`, gameID, editorID)

		if err != nil {
			log.Fatal(err)
			return nil, err
		}
	}

	for _, studioID := range input.StudioIDs {
		_, err := r.DB.Exec(`
			INSERT INTO GameStudio (gameID, studioID)
			VALUES (?, ?)
		`, gameID, studioID)

		if err != nil {
			log.Fatal(err)
			return nil, err
		}
	}

	var publicationDate int
	if input.PublicationDate != nil {
		publicationDate = *input.PublicationDate
	}

	id := strconv.Itoa(int(gameID))
	return &model.Game{
		ID:              &id,
		Name:            input.Name,
		Genres:          input.Genres,
		PublicationDate: &publicationDate,
		Platform:        input.Platform,
	}, nil
}

func (r *mutationResolver) UpdateGame(ctx context.Context, id string, input model.GameInput) (*model.Game, error) {
	_, err := r.DB.Exec(`
		UPDATE Game
		SET name = ?, publicationDate = ?, platform = ?
		WHERE id = ?
	`, input.Name, input.PublicationDate, input.Platform, id)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	_, err = r.DB.Exec(`
		DELETE FROM GameEditor WHERE gameID = ?
	`, id)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	_, err = r.DB.Exec(`
		DELETE FROM GameStudio WHERE gameID = ?
	`, id)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	for _, editorID := range input.EditorIDs {
		_, err := r.DB.Exec(`
			INSERT INTO GameEditor (gameID, editorID)
			VALUES (?, ?)
		`, id, editorID)

		if err != nil {
			log.Fatal(err)
			return nil, err
		}
	}

	for _, studioID := range input.StudioIDs {
		_, err := r.DB.Exec(`
			INSERT INTO GameStudio (gameID, studioID)
			VALUES (?, ?)
		`, id, studioID)

		if err != nil {
			log.Fatal(err)
			return nil, err
		}
	}

	var publicationDate int
	if input.PublicationDate != nil {
		publicationDate = *input.PublicationDate
	}

	return &model.Game{
		ID:              &id,
		Name:            input.Name,
		Genres:          input.Genres,
		PublicationDate: &publicationDate,
		Platform:        input.Platform,
	}, nil
}

func (r *mutationResolver) DeleteGame(ctx context.Context, id string) (*string, error) {
	_, err := r.DB.Exec(`
		DELETE FROM Game
		WHERE id = ?
	`, id)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	_, err = r.DB.Exec(`
		DELETE FROM GameEditor WHERE gameID = ?
	`, id)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	_, err = r.DB.Exec(`
		DELETE FROM GameStudio WHERE gameID = ?
	`, id)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &id, nil
}

func (r *mutationResolver) Mutation() MutationResolver {
	return &mutationResolver{r.Resolver}
}

func (r *mutationResolver) CreateEditor(ctx context.Context, input model.EditorInput) (*model.Editor, error) {
	var existingEditorID int
	err := r.DB.QueryRow("SELECT id FROM Editor WHERE name = ?", input.Name).Scan(&existingEditorID)
	if err == nil {
		return nil, gqlerror.Errorf("Editor with the name already exists")
	} else if err != sql.ErrNoRows {
		log.Fatal(err)
		return nil, err
	}

	result, err := r.DB.Exec("INSERT INTO Editor (name) VALUES (?)", input.Name)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	editorID, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	a := strconv.Itoa(int(editorID))
	newEditor := &model.Editor{
		ID:   &a,
		Name: input.Name,
	}

	return newEditor, nil
}

func (r *mutationResolver) CreateStudio(ctx context.Context, input model.StudioInput) (*model.Studio, error) {
	var existingStudioID int
	err := r.DB.QueryRow("SELECT id FROM Studio WHERE name = ?", input.Name).Scan(&existingStudioID)
	if err == nil {
		return nil, gqlerror.Errorf("Studio with the name already exists")
	} else if err != sql.ErrNoRows {
		log.Fatal(err)
		return nil, err
	}

	result, err := r.DB.Exec("INSERT INTO Studio (name) VALUES (?)", input.Name)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	studioID, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	b := strconv.Itoa(int(studioID))

	newStudio := &model.Studio{
		ID:   &b,
		Name: input.Name,
	}

	return newStudio, nil
}

func (r *mutationResolver) DeleteEditor(ctx context.Context, id string) (*string, error) {

	result, err := r.DB.Exec("DELETE FROM Editor WHERE id = ?", id)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, gqlerror.Errorf("Editor not found")
	}

	return &id, nil
}

func (r *mutationResolver) DeleteStudio(ctx context.Context, id string) (*string, error) {

	result, err := r.DB.Exec("DELETE FROM Studio WHERE id = ?", id)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, gqlerror.Errorf("Studio not found")
	}

	return &id, nil
}

func (r *mutationResolver) UpdateEditor(ctx context.Context, id string, input model.EditorInput) (*model.Editor, error) {
	result, err := r.DB.Exec("UPDATE Editor SET name = ? WHERE id = ?", input.Name, id)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, gqlerror.Errorf("Editor not found")
	}

	updatedEditor, err := r.Editor(ctx, id)
	if err != nil {
		return nil, err
	}

	return updatedEditor, nil
}

func (r *mutationResolver) UpdateStudio(ctx context.Context, id string, input model.StudioInput) (*model.Studio, error) {

	result, err := r.DB.Exec("UPDATE Studio SET name = ? WHERE id = ?", input.Name, id)
	if err != nil {
		return nil, err
	}

	// Vérifiez le nombre de lignes affectées par la mise à jour
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, gqlerror.Errorf("Studio not found")
	}

	// Récupérez le studio mis à jour
	updatedStudio, err := r.Studio(ctx, id)
	if err != nil {
		return nil, err
	}

	// Retournez le studio mis à jour sous forme de modèle
	return &model.Studio{
		ID:    updatedStudio.ID,
		Name:  updatedStudio.Name,
		Games: updatedStudio.Games,
	}, nil
}
