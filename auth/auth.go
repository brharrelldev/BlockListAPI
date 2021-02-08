package auth

import (
	"fmt"
	"github.com/dgraph-io/badger/v2"
	"net/http"
)

func Authorize(db *badger.DB) func(next http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

			var res []byte
			if err := db.View(func(txn *badger.Txn) error {
				item, err := txn.Get([]byte("credentials"))

				if err != nil{
					http.Error(writer, err.Error(), http.StatusForbidden)
					return fmt.Errorf("error looking up credentialws %v", err)
				}


				res, err = item.ValueCopy(nil)
				if err != nil{
					http.Error(writer, err.Error(), http.StatusForbidden)
					return fmt.Errorf("error getting credentials %v", err)
				}



				return nil

			}); err != nil{
				return
			}

			if request.Header.Get("Authorization") != string(res){
				http.Error(writer, "invalid credentials", http.StatusForbidden)
				return
			}

			next.ServeHTTP(writer, request)
		})

	}


}
