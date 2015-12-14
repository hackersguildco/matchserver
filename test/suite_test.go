package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"gopkg.in/mgo.v2/bson"

	"github.com/cheersappio/matchserver/models"
	"github.com/cheersappio/matchserver/utils"
	"github.com/cheersappio/matchserver/ws"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	ts *httptest.Server
)

func Test(t *testing.T) {
	RegisterFailHandler(Fail)
	fmt.Println("Suite found")
	RunSpecs(t, "Api WS")
}

var _ = BeforeSuite(func() {
	initEnv()
	utils.InitLog()
	models.InitDB()
	router := mux.NewRouter()
	router.HandleFunc("/ws/{username}", ws.ServeWS)
	http.Handle("/", router)
	ts = httptest.NewServer(router)
	ws.InitSearcher()
})

var _ = AfterSuite(func() {
	ts.Close()
	models.Session.Close()
})

var _ = BeforeEach(func() {
	cleanDB()
	wsConnUser1 = createClient(ts.URL, username1)
	wsConnUser2 = createClient(ts.URL, username2)
	wsConnUser3 = createClient(ts.URL, username3)

})

var _ = AfterEach(func() {
	wsConnUser1.Close()
	wsConnUser2.Close()
	wsConnUser3.Close()
})

func cleanDB() {
	models.StrokesCollection.RemoveAll(bson.M{})
}

func initEnv() {
	path := ".env_test"
	for i := 1; ; i++ {
		if err := godotenv.Load(path); err != nil {
			if i > 3 {
				panic("Error loading .env_test file")
			} else {
				path = "../" + path
			}
		} else {
			break
		}
	}
}
