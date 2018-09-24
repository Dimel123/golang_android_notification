package main

import (
	"fmt"
	"github.com/NaySoftware/go-fcm"
	"github.com/martini-contrib/render"
	"github.com/codegangsta/martini"
	"testNotification/db/schemas"
	"gopkg.in/mgo.v2"
	"testNotification/models"
	"net/http"
	"gopkg.in/mgo.v2/bson"
	"testNotification/utils"
)

const (
	serverKey = ""
)

var NP fcm.NotificationPayload
var databaseCollections *mgo.Collection

func main() {
	setupDatabase()
	setupMartini()
	setupBackgroundNotification()
}

func setupMartini() {
	m := martini.Classic()
	m.Get("/", indexHandler)
	m.Post("/send_notification", sendNotificationHandler)
	m.Get("/send_notification_to_all", sendAllNotificationHandler)
	m.Post("/create", saveTokenHandler)
	m.Use(render.Renderer(render.Options{
		Directory:  "templates",                // Specify what path to load the templates from.
		Layout:     "layout",                   // Specify a layout template. Layouts can call {{ yield }} to render the current template.
		Extensions: []string{".tmpl", ".html"}, // Specify extensions to load for templates.
		Charset:    "UTF-8",                    // Sets encoding for json and html content-types. Default is "UTF-8".
		IndentJSON: true,                       // Output human readable JSON
	}))
	m.Run()
}

func setupDatabase() {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	databaseCollections = session.DB("FCM").C("tokens")
}

func indexHandler(rnd render.Render) {
	rnd.HTML(200, "index", nil)
}

func saveTokenHandler(rnd render.Render, r *http.Request) {

	token := utils.ParseToken(r)
	result := schemas.Schema{}
	err := databaseCollections.Find(bson.M{"account_id": token.Account_id}).One(&result)
	if (err != nil) {
		databaseCollections.Insert(token)
	} else {
		databaseCollections.UpdateId(token.Account_id, token)
	}

	rnd.HTML(200, "index", nil)
}

func sendNotificationHandler(r *http.Request, rnd render.Render) {

	notification := utils.ParseNotification(r)
	result := schemas.Schema{}
	err := databaseCollections.Find(bson.M{"account_id": notification.Account_id}).One(&result)
	if (err != nil) {
		panic(err)
	} else {
		data := map[string]string{"push_params": fmt.Sprintf(`{"event_type":"%s", "is_popup":%t, "event_details":"%s", "entity":"%s","event_uid":"%s"}`,
			notification.Event_type, notification.Is_popup, notification.Event_details, notification.Entity, notification.Event_uid) }
		sendNotificationToId(result.Token_id, data)
	}
	rnd.HTML(200, "index", nil)
}

func sendAllNotificationHandler(rnd render.Render) {

	allTokens := []schemas.Schema{}
	databaseCollections.Find(nil).All(&allTokens)

	tokens := []models.Token{}

	for _, item := range allTokens {
		token := models.Token{item.Account_id, item.Token_id}
		sendNotificationToId(token.Token_id, map[string]string{
			"push_params": `{"event_type":"transaction_add",
			"is_popup":false,
			"event_details":"transaction_rejected",
			"entity":"Transaction",
			"event_uid":"123123"}`})
		tokens = append(tokens, token)
	}
	rnd.HTML(200, "index", nil)
}

func setupBackgroundNotification() {
	NP.Title = "transaction_rejected"
	NP.Body = "transaction_add"
	NP.Sound = "mySound"
	NP.Icon = "myicon"
}

func sendNotificationToId(tokenId string, data map[string]string) {
	ids := []string{
		tokenId,
	}
	c := fcm.NewFcmClient(serverKey)
	c.NewFcmRegIdsMsg(ids, data)
	c.SetNotificationPayload(&NP)
	status, err := c.Send()
	if err == nil {
		status.PrintResults()
	} else {
		fmt.Println(err)
	}
}

