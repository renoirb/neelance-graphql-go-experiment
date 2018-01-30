package main

import (
	"log"
	"math"
	"net/http"

	"io/ioutil"

	"fmt"
	"time"

	graphql "github.com/neelance/graphql-go"
	"github.com/neelance/graphql-go/relay"
)

/**
 * Some simple GraphQL server written in Go
 *
 * Started from https://github.com/lpalmes/graphql-go-introduction/blob/master/main.go
 *
 * See also:
 * https://github.com/rgraphql/rgraphql-demo-server/blob/master/demo.go
 * https://github.com/neelance/graphql-go/blob/master/example/starwars/starwars.go
 * https://github.com/neelance/graphql-go/blob/master/graphql_test.go
 */

// This will be use by our handler at /graphql
var schema *graphql.Schema

// This function runs at the start of the program
func init() {

	// We get the schema from the file, rather than having the schema inline here
	// I think will lead to better organizaiton of our own code
	schemaFile, err := ioutil.ReadFile("schema.graphql")
	if err != nil {
		// We will panic if we don't find the schema.graphql file in our server
		panic(err)
	}

	// We will use graphql-go library to parse our schema from "schema.graphql"
	// and the resolver is our struct that should fullfill everything in the Query
	// from our schema
	schema, err = graphql.ParseSchema(string(schemaFile), &Resolver{})
	if err != nil {
		panic(err)
	}
}

func main() {
	// We will start a small server that reads our "lib/graphiql.html" file and
	// responds with it, so we are able to have our own graphiql
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page, err := ioutil.ReadFile("lib/graphiql.html")
		if err != nil {
			log.Fatal(err)
		}
		w.Write(page)
	}))

	// This is where our graphql server is handled, we declare "/graphql" as the route
	// where all our graphql requests will be directed to
	http.Handle("/graphql", &relay.Handler{Schema: schema})

	// We start the server by using ListenAndServe and we log if we have any error, hope not!
	fmt.Println("Listening at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

var people = []*person{
	{
		ID:   "1000",
		Name: "Luke Skywalker",
		Date: time.Date(1951, 2, 3, 4, 5, 6, 0, time.UTC),
	},
	{
		ID:   "1001",
		Name: "Leia Organa",
		Date: time.Date(1951, 2, 3, 4, 5, 6, 1, time.UTC),
	},
	{
		ID:   "1002",
		Name: "Darth Vader",
		Date: time.Date(1931, 2, 3, 4, 5, 6, 0, time.UTC),
	},
	{
		ID:   "1003",
		Name: "Han Solo",
		Date: time.Date(1946, 2, 3, 4, 5, 6, 1, time.UTC),
	},
	{
		ID:   "1004",
		Name: "Wilhuff Tarkin",
		Date: time.Date(1942, 2, 3, 4, 5, 6, 1, time.UTC),
	},
}

type DateTime struct {
	When time.Time
}

func (d DateTime) FormatOne() string {
	when := d.When
	formattedWhen := when.Format(time.RFC3339)
	return fmt.Sprintf("%v", formattedWhen)
}

//func (Time graphql.Time) DaysAgo() string {
//	when := d.When
//	formattedWhen := when.Format(time.RFC3339)
//	return fmt.Sprintf("%v", formattedWhen)
//}

// See also https://golang.org/src/time/format.go
func msToTime(ms int64) (time.Time, error) {
	return time.Unix(0, ms*int64(time.Millisecond)), nil
}

func NewDateTime(ms int64) DateTime {
	t, err := msToTime(ms)
	if err != nil {
		fmt.Printf("Error at NewDateTime: %s", err)
	}
	return DateTime{
		When: t,
	}
}

// The Resolver root on which we attach methods
type Resolver struct{}

/**
 * Data Schema
 * An entity Resolver, and methods for them
 */
type person struct {
	ID   graphql.ID
	Name string
	Date time.Time
}

type personResolver struct {
	person
}

func (r *personResolver) ID() graphql.ID {
	return r.person.ID
}

func (r *personResolver) Name() string {
	return r.person.Name
}

func (r *personResolver) Date() *graphql.Time {
	return &graphql.Time{Time: r.person.Date}
}

func (r *personResolver) Days() string {
	input := r.person.Date
	days := math.Floor(math.Abs(input.Sub(time.Now()).Hours() / 24))

	return fmt.Sprintf("%v days ago", days)
}

// AllPeople get a list of people
func (r *Resolver) AllPeople() []*personResolver {
	var result []*personResolver
	for _, person := range people {
		result = append(result, &personResolver{*person})
	}

	return result
}

func (r *Resolver) CreatePerson(args *struct {
	Name string
}) *personResolver {
	person := &person{
		Name: args.Name,
		Date: time.Now().UTC(),
	}
	people = append(people, person)

	return &personResolver{*person}
}
