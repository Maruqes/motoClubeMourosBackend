package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Maruqes/Tokenize"
	funchooks "github.com/Maruqes/Tokenize/FuncHooks"
	functions "github.com/Maruqes/Tokenize/Functions"
	login "github.com/Maruqes/Tokenize/Login"
	types "github.com/Maruqes/Tokenize/Types"
	"github.com/Maruqes/Tokenize/UserFuncs"
	"github.com/Maruqes/Tokenize/database"
)

func testLogado(w http.ResponseWriter, r *http.Request) {
	if !login.CheckToken(r) {
		w.Write([]byte("Não Logadoo"))
		return
	}
	w.Write([]byte("Logado"))
}

func testPago(w http.ResponseWriter, r *http.Request) {
	if !login.CheckToken(r) {
		w.Write([]byte("Não Logado"))
		return
	}
	if val, err := login.IsUserActiveRequest(r); err != nil || !val {
		w.Write([]byte("Não Ativo"))
		return
	}
	w.Write([]byte("Pago"))
}

// needs to cheak all code for offline payments and user funcs
func checkIfUserHasLatePayments(id int) (bool, []int64, error) {
	endDate, err := UserFuncs.GetEndDateForUser(id)
	if err != nil {
		return false, []int64{}, err
	}

	if endDate == (database.Date{}) {
		return false, []int64{}, nil
	}

	timeDB := time.Date(endDate.Year, time.Month(endDate.Month), endDate.Day, 0, 0, 0, 0, time.UTC).Unix()

	mourosDate, _ := functions.GetMourosStartingDate()
	lastMourosDate := time.Date(time.Now().Year(), mourosDate.Month(), mourosDate.Day(), 0, 0, 0, 0, time.UTC).Unix()

	yearDifference := time.Unix(timeDB, 0).Year() - time.Unix(lastMourosDate, 0).Year()
	fmt.Println("yearDifference: ", yearDifference)
	fmt.Println("timeDB: ", time.Unix(timeDB, 0))
	fmt.Println("lastMourosDate: ", time.Unix(lastMourosDate, 0))

	if yearDifference < 0 {
		yearsLeft := []int64{}
		for i := 0; i < -yearDifference; i++ {
			yearsLeft = append(yearsLeft, time.Date(int(time.Unix(timeDB, 0).Year())+i, 1, 1, 0, 0, 0, 0, time.UTC).Unix())
		}

		return true, yearsLeft, nil
	}

	return false, []int64{}, nil
}

func checkIfUserHasLatePaymentsRequest(w http.ResponseWriter, r *http.Request) bool {
	res, numberOfYears, err := checkIfUserHasLatePayments(1)
	if err != nil {
		fmt.Println("Error checking if user has late payments")
		return false
	}
	if res {
		fmt.Println("User has late payments: ", numberOfYears)
		w.Write([]byte("User has late payments"))
		return true
	}
	return false
}

func main() {
	Tokenize.Initialize()

	date, err := functions.GetMourosStartingDate()
	if err != nil {
		fmt.Println("Error getting mouros date")
		return
	}
	fmt.Println("Mouros date: ", date)
	fmt.Println()
	fmt.Println()

	res, numberOfYears, err := checkIfUserHasLatePayments(2)
	if err != nil {
		fmt.Println("Error checking if user has late payments")
	}
	if res {
		for _, year := range numberOfYears {
			fmt.Println()
			fmt.Println(time.Unix(year, 0).Format("01/12/2006"))
		}
	}
	fmt.Println()
	fmt.Println()

	http.HandleFunc("/logado", testLogado)
	http.HandleFunc("/pago", testPago)
	// http.HandleFunc("/hasLatePayments", testPago)

	// UserFuncs.ProhibitUser(0)
	// UserFuncs.UnprohibitUser(0)
	funchooks.SetCheckout_UserFunc(checkIfUserHasLatePaymentsRequest)

	if os.Getenv("DEV") == "True" {
		http.FileServer(http.Dir("public"))
	}
	Tokenize.InitListen("4242", "/sucess", "/cancel", types.TypeOfSubscriptionValues.MourosSubscription, []types.ExtraPayments{types.ExtraPaymentsValues.Multibanco})
}
