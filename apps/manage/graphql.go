package manage

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/thepatrick/lunch/model"
	"github.com/thepatrick/lunch/support"

	"github.com/graphql-go/graphql"
)

var placeType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Place",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					m := p.Source.(model.Place)
					return m.ID.Hex(), nil
				},
			},
			"team_id": &graphql.Field{
				Type: graphql.String,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"last_visited": &graphql.Field{
				Type: graphql.DateTime,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					m := p.Source.(model.Place)
					if m.LastVisited.IsZero() {
						return nil, nil
					}
					return m.LastVisited, nil
				},
			},
			"last_skipped": &graphql.Field{
				Type: graphql.DateTime,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					m := p.Source.(model.Place)
					if m.LastSkipped.IsZero() {
						return nil, nil
					}
					return m.LastSkipped, nil
				},
			},
			"skip_count": &graphql.Field{
				Type: graphql.Int,
			},
			"visit_count": &graphql.Field{
				Type: graphql.Int,
			},
		},
	},
)

var queryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"placeList": &graphql.Field{
				Type:        graphql.NewList(placeType),
				Description: "List of places",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					// p.Context.Value("bob")
					app := p.Context.Value("app").(App)
					session := p.Context.Value("session").(validSession)

					return app.places.AllPlaces(session.user.Team.ID)
				},
			},
			"place": &graphql.Field{
				Type: placeType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					// p.Context.Value("bob")
					app := p.Context.Value("app").(App)
					session := p.Context.Value("session").(validSession)

					return app.places.FindByID(session.user.Team.ID, p.Args["id"].(string))
				},
			},
		},
	},
)

func (app App) placesGraphql() http.HandlerFunc {
	s, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: queryType,
	})
	if err != nil {
		log.Fatalf("failed to create schema, error: %v", err)
	}
	Schema := s

	globalContext := context.WithValue(context.Background(), "app", app)

	return app.withValidSession(func(w http.ResponseWriter, r *http.Request, session validSession) {
		result := graphql.Do(graphql.Params{
			Schema:        Schema,
			RequestString: r.URL.Query().Get("query"),
			Context:       context.WithValue(globalContext, "session", session),
		})
		if len(result.Errors) > 0 {
			log.Printf("wrong result, unexpected errors: %v", result.Errors)
			support.ErrorWithJSON(w, "Unexpected errors", http.StatusInternalServerError)
			return
		}

		respBody, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			log.Fatal(err)
			support.ErrorWithJSON(w, "Failed to generate JSON", http.StatusInternalServerError)
			return
		}

		support.ResponseWithJSON(w, respBody, http.StatusOK)

		// 	allPlaces, err := app.places.AllPlaces(session.user.Team.ID)
		// 	if err != nil {
		// 		log.Printf("Failed to get all places: %v\n", err)
		// 		support.ErrorWithJSON(w, "Failed to get places", http.StatusInternalServerError)
		// 		return
		// 	}

		// 	respBody, err := json.MarshalIndent(allPlaces, "", "  ")
		// 	if err != nil {
		// 		log.Fatal(err)
		// 		support.ErrorWithJSON(w, "Failed to generate JSON", http.StatusInternalServerError)
		// 		return
		// 	}

		// 	support.ResponseWithJSON(w, respBody, http.StatusOK)
	})
}
