package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/mnmtanish/go-graphiql"

	"github.com/charliekenney23/hello-go-graphql/model"
	"github.com/graphql-go/graphql"
)

// Todos - shared Todo instances
var Todos []model.Todo

func init() {
	t1 := *(&model.TodoBuffer{Title: "Wash the dog", IsComplete: false}).NewTodo()
	t2 := *(&model.TodoBuffer{Title: "Buy some milk", IsComplete: false}).NewTodo()
	t3 := *(&model.TodoBuffer{Title: "Do homework", IsComplete: false}).NewTodo()
	Todos = append(Todos, t1, t2, t3)
}

func main() {
	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		res := executeQuery(r.URL.Query().Get("query"), schema)
		json.NewEncoder(w).Encode(res)
	})

	http.HandleFunc("/graphiql", graphiql.ServeGraphiQL)

	http.ListenAndServe(":8080", nil)

	fmt.Println("listening at http://localhost:8080/graphql")
}

func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	res := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})

	if len(res.Errors) > 0 {
		fmt.Printf("Wrong result, expected: %v\n", res.Errors)
	}

	return res
}

var todoType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Todo",
	Fields: graphql.Fields{
		"id":         &graphql.Field{Type: graphql.String},
		"title":      &graphql.Field{Type: graphql.String},
		"isComplete": &graphql.Field{Type: graphql.Boolean},
	},
})

var rootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootMutation",
	Fields: graphql.Fields{

		/**
		 * Create Todo mutation
		 *
		 * createTodo($title: String!) {
		 *   // ...Todo
		 * }
		 */
		"createTodo": &graphql.Field{
			Type:        todoType,
			Description: "Create a new todo",
			Args: graphql.FieldConfigArgument{
				"title": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				title, _ := params.Args["title"].(string)
				buf := model.TodoBuffer{Title: title, IsComplete: false}

				newTodo := *buf.NewTodo()

				Todos = append(Todos, newTodo)

				return newTodo, nil
			},
		},

		/**
		 * Update Todo mutation
		 *
		 * updateTodo($id: String!, $title: String, $isComplete: Boolean) {
		 *   // ...Todo
		 * }
		 */
		"updateTodo": &graphql.Field{
			Type:        todoType,
			Description: "Update existing todo",
			Args: graphql.FieldConfigArgument{
				"id":         &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"title":      &graphql.ArgumentConfig{Type: graphql.String},
				"isComplete": &graphql.ArgumentConfig{Type: graphql.Boolean},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				var ind *int

				id, _ := params.Args["id"].(string)
				// get index of original todo
				for i := 0; i < len(Todos); i++ {
					if Todos[i].ID == id {
						tmp := i
						ind = &tmp
					}
				}

				// exit if not found
				if ind == nil {
					return nil, errors.New("Todo not found")
				}

				for {
					isComplete, ok := params.Args["isComplete"].(bool)
					if ok {
						fmt.Println("updating isComplete")
						Todos[*ind].IsComplete = isComplete
					}
					break
				}

				for {
					title, ok := params.Args["title"].(string)
					if ok {
						fmt.Println("updating title")
						Todos[*ind].Title = title
					}
					break
				}

				return Todos[*ind], nil
			},
		},
	},
})

var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{

		/* Todo Query
		 *
		 * todo($id: String!) {
		 *   // ...Todo
		 * }
		 */
		"todo": &graphql.Field{
			Type:        todoType,
			Description: "Get a Todo by ID",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{Type: graphql.String},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				id, _ := params.Args["id"].(string)

				for _, todo := range Todos {
					if todo.ID == id {
						return todo, nil
					}
				}

				return nil, errors.New("Todo does not exist")
			},
		},

		"lastTodo": &graphql.Field{
			Type:        todoType,
			Description: "Get the latest Todo",
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				if len(Todos) == 0 {
					return nil, errors.New("No Todos to show")
				}

				return Todos[len(Todos)-1], nil
			},
		},

		"todos": &graphql.Field{
			Type:        graphql.NewList(todoType),
			Description: "Get all Todos",
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				return Todos, nil
			},
		},
	},
})

var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    rootQuery,
	Mutation: rootMutation,
})
