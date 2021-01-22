package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	contact "contact.com"
	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/bson"
)

var mh *contact.MongoHandler

func getContact(w http.ResponseWriter, r *http.Request) {
	phoneNumber := chi.URLParam(r, "phonenumber")
	if phoneNumber == "" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)

		return
	}

	contact := &contact.Contact{}

	err := mh.GetOne(contact, bson.M{"phoneNumber": phoneNumber})
	if err != nil {
		http.Error(w, fmt.Sprintf("Contact with phonenumber: %s not found", phoneNumber), 404)

		return
	}

	_ = json.NewEncoder(w).Encode(contact)
}

func getAllContact(w http.ResponseWriter, r *http.Request) {
	contacts := mh.Get(bson.M{})

	_ = json.NewEncoder(w).Encode(contacts)
}

func addContact(w http.ResponseWriter, r *http.Request) {
	existingContact := &contact.Contact{}

	var contact contact.Contact
	_ = json.NewDecoder(r.Body).Decode(&contact)
	contact.CreatedOn = time.Now()

	err := mh.GetOne(existingContact, bson.M{"phoneNumber": contact.PhoneNumber})
	if err == nil {
		http.Error(w, fmt.Sprintf("Contact with phonenumber: %s already exist", contact.PhoneNumber), http.StatusBadRequest)

		return
	}

	_, err = mh.AddOne(&contact)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)

		return
	}

	_, _ = w.Write([]byte("Contact created successfully"))
	w.WriteHeader(http.StatusCreated)
}

func deleteContact(w http.ResponseWriter, r *http.Request) {
	existingContact := &contact.Contact{}

	phoneNumber := chi.URLParam(r, "phonenumber")
	if phoneNumber == "" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)

		return
	}

	err := mh.GetOne(existingContact, bson.M{"phoneNumber": phoneNumber})
	if err != nil {
		http.Error(w, fmt.Sprintf("Contact with phonenumber: %s does not exist", phoneNumber), http.StatusBadRequest)

		return
	}

	_, err = mh.RemoveOne(bson.M{"phoneNumber": phoneNumber})
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)

		return
	}

	_, _ = w.Write([]byte("Contact deleted"))
	w.WriteHeader(http.StatusOK)
}

func updateContact(w http.ResponseWriter, r *http.Request) {
	phoneNumber := chi.URLParam(r, "phonenumber")
	if phoneNumber == "" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)

		return
	}

	contact := &contact.Contact{}

	_ = json.NewDecoder(r.Body).Decode(contact)

	_, err := mh.Update(bson.M{"phoneNumber": phoneNumber}, contact)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)

		return
	}

	_, _ = w.Write([]byte("Contact update successful"))
	w.WriteHeader(http.StatusOK)
}

func registerRoutes() http.Handler {
	router := chi.NewRouter()
	router.Route("/contacts", func(r chi.Router) {
		r.Get("/", getAllContact)                 // GET /contacts
		r.Get("/{phonenumber}", getContact)       // GET /contacts/0147344454
		r.Post("/", addContact)                   // POST /contacts
		r.Put("/{phonenumber}", updateContact)    // PUT /contacts/0147344454
		r.Delete("/{phonenumber}", deleteContact) // DELETE /contacts/0147344454
	})

	return router
}

func main() {
	mongoDBConnection := "mongodb://localhost:27017"
	mh = contact.NewHandler(mongoDBConnection) // Create an instance of MongoHander with the connection string provided
	r := registerRoutes()
	log.Fatal(http.ListenAndServe(":3060", r)) // You can modify to run on a different port
}
