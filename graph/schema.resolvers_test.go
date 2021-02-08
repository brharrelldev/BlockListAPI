package graph

import (
	"fmt"
	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/brharrelldev/BlockListAPI/database"
	"github.com/brharrelldev/BlockListAPI/graph/generated"
	"reflect"
	"testing"
)


func TestQueryResolver_GetIPDetails(t *testing.T) {

	var resp struct {
		GetIpDetails map[string]interface{}
	}

	expected := struct {
		GetIpDetails map[string]interface{}
	}{
		GetIpDetails: map[string]interface{}{
			"created_at":   "09-14-1979",
			"id":           "11111",
			"ip_address":   "0.0.0.0",
			"response":     "test",
			"response_cde": "teast",
			"updated_at":   "test",
		},
	}

	db, err := database.NewBlocklistDB("../data/blocklist.db")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(db.DB.Query("SELECT * FROM BlockList"))

	srv := client.New(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: &Resolver{
			DB: db,
		},
	})))

	srv.MustPost(`query GetIpDetail{
  GetIpDetails(ip: "0.0.0.0"){
      id
      response
      created_at
      updated_at
      response_code
      ip_address
  }
}`, &resp)

	t.Run("testing GetIpDetails", func(t *testing.T) {
		if !reflect.DeepEqual(resp.GetIpDetails["ip_address"], expected.GetIpDetails["ip_address"]) {
			t.Errorf("got %v expected %v", resp, expected)
		}
	})

}

func TestMutationResolver_Enqueue(t *testing.T) {

	db, err := database.NewBlocklistDB("../data/blocklist.db")
	if err != nil {
		fmt.Println(err)
	}

	resp := make(map[string]interface{})
	expected := map[string][]map[string]interface{}{
		"enqueue": {
			{
				"ip":      "104.218.235.41",
				"message": "timestamp updated",
			},
		},
	}

	srv := client.New(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: &Resolver{
			DB: db,
		},
	})))

	srv.MustPost(`mutation{
		enqueue(ips: ["104.218.235.41", "50.7.10.170"]){
		ip
		message
	}
	}`, &resp)

	t.Run("enqueue test", func(t *testing.T) {
		covertedResp := resp["enqueue"].([]interface{})[0].(interface{}).(map[string]interface{})
		covertedExpected := expected["enqueue"][0]

		if !reflect.DeepEqual(covertedResp["ip"], covertedExpected["ip"]) {
			t.Errorf("error got %v expected %v", resp, expected)
		}

	})

}
