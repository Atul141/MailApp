package handlers

import (
	"fmt"
	"net/http"

	"encoding/json"

	m "git.mailbox.com/mailbox/models"
	u "git.mailbox.com/mailbox/utils"
)

func dealersHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var dealers []m.Dealer
		dealers = append(dealers, m.Dealer{
			DealerID:   u.SPtr("2cd9c96a-edc1-4196-93ca-13d2f756d9a0"),
			DealerName: "DTDC",
		})
		dealers = append(dealers, m.Dealer{
			DealerID:   u.SPtr("ada1103c-4024-4ea4-b955-58c1c2c702b7"),
			DealerName: "Flipkart",
		})
		dealers = append(dealers, m.Dealer{
			DealerID:   u.SPtr("884a4cf3-2399-4481-b872-13233eaa3d6f"),
			DealerName: "Amazon",
		})

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin","*")
		w.Header().Set("Access-Control-Allow-Headers","Origin, X-Requested-With, Content-Type, Accept")

		marshalledRes, _ := json.Marshal(dealers)
		fmt.Fprintf(w, string(marshalledRes))
	}
}
