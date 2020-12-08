package register

import (
	"fmt"
	"strconv"
	"testing"
)

func TestEndpoint(t *testing.T) {
	fmt.Println("----------ENDPOINT TEST-------------")
	db, err := New()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	fmt.Printf("-----------\nTesting Backup Load\n-----------\n")
	for i := 0; i < 6; i++ {
		originID := "/origin/" + strconv.Itoa(i+1)
		fmt.Printf("(%s)[REGISTRATION STATUS] = %t\n", originID, db.IsRegistered(originID))
	}

	fmt.Printf("---------------\nTesting Database Manipulation\n---------------\n")
	origin := "/origin/TESTORIGIN"
	fmt.Printf("(%s)[REGISTRATION STATUS] = %t\n", origin, db.IsRegistered(origin))

	err = db.AddOrigin(origin)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Inserted (%s) into database\n", origin)
	fmt.Printf("(%s)[REGISTRATION STATUS] = %t\n", origin, db.IsRegistered(origin))

	err = db.RemoveOrigin(origin)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Removed (%s) into database\n", origin)
	fmt.Printf("(%s)[REGISTRATION STATUS] = %t\n", origin, db.IsRegistered(origin))
}
