// resolver.go

package graph

import (
	"context"
	"database/sql"

	"github.com/kilianp07/graphql/graph/model"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

type Resolver struct {
	DB *sql.DB // Connexion à la base de données
	mutationResolver
	queryResolver
}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Games(ctx context.Context, page *int, genre *string, platform *string, studio *string) (*model.Games, error) {
	rows, err := r.DB.Query("SELECT id, name, genres, publicationDate, platform FROM games")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var games []*model.Game
	for rows.Next() {
		var game model.Game
		if err := rows.Scan(&game.ID, &game.Name, &game.Genres, &game.PublicationDate, &game.Platform); err != nil {
			return nil, err
		}
		games = append(games, &game)
	}

	infos := &model.Infos{
		Count:         len(games),
		Pages:         1,
		NextPage:      nil,
		PreviousPages: nil,
	}

	return &model.Games{Infos: infos, Results: games}, nil
}

func (r *queryResolver) Game(ctx context.Context, id string) (*model.Game, error) {
	var game model.Game
	err := r.DB.QueryRow("SELECT id, name, genres, publicationDate, platform FROM games WHERE id = ?", id).
		Scan(&game.ID, &game.Name, &game.Genres, &game.PublicationDate, &game.Platform)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, gqlerror.Errorf("Game not found")
		}
		return nil, err
	}

	return &game, nil
}

func (r *queryResolver) Editors(ctx context.Context) ([]*model.Editor, error) {
	rows, err := r.DB.Query("SELECT id, name FROM editors")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var editors []*model.Editor
	for rows.Next() {
		var editor model.Editor
		if err := rows.Scan(&editor.ID, &editor.Name); err != nil {
			return nil, err
		}
		editors = append(editors, &editor)
	}

	return editors, nil
}

func (r *queryResolver) Editor(ctx context.Context, id string) (*model.Editor, error) {

	var editor model.Editor
	err := r.DB.QueryRow("SELECT id, name FROM editors WHERE id = ?", id).
		Scan(&editor.ID, &editor.Name)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, gqlerror.Errorf("Editor not found")
		}
		return nil, err
	}

	return &editor, nil
}

func (r *queryResolver) Studios(ctx context.Context) ([]*model.Studio, error) {
	rows, err := r.DB.Query("SELECT id, name FROM studios")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var studios []*model.Studio
	for rows.Next() {
		var studio model.Studio
		if err := rows.Scan(&studio.ID, &studio.Name); err != nil {
			return nil, err
		}
		studios = append(studios, &studio)
	}

	return studios, nil
}

func (r *queryResolver) Studio(ctx context.Context, id string) (*model.Studio, error) {
	var studio model.Studio
	err := r.DB.QueryRow("SELECT id, name FROM studios WHERE id = ?", id).
		Scan(&studio.ID, &studio.Name)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, gqlerror.Errorf("Studio not found")
		}
		return nil, err
	}

	return &studio, nil
}
