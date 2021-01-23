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

	"github.com/go-chi/chi/middleware"
)

func deleteContact(handler *contact.MongoHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		existingContact := &contact.Contact{}

		phoneNumber := chi.URLParam(r, "phonenumber")
		if phoneNumber == "" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)

			return
		}

		err := handler.GetOne(existingContact, bson.M{"phoneNumber": phoneNumber})
		if err != nil {
			http.Error(w, fmt.Sprintf("Contact with phonenumber: %s does not exist", phoneNumber), http.StatusBadRequest)

			return
		}

		_, err = handler.RemoveOne(bson.M{"phoneNumber": phoneNumber})
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusBadRequest)

			return
		}

		_, _ = w.Write([]byte("Contact deleted"))
		w.WriteHeader(http.StatusOK)
	}
}

func updateContact(handler *contact.MongoHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		phoneNumber := chi.URLParam(r, "phonenumber")
		if phoneNumber == "" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)

			return
		}

		contact := &contact.Contact{}

		_ = json.NewDecoder(r.Body).Decode(contact)

		update := bson.D{
			{"$set", bson.D{{"phoneNumber", "Nicolas Raboy"}}},
		}

		fmt.Println(contact)

		_, err := handler.Update(update, contact)
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusBadRequest)

			return
		}

		_, _ = w.Write([]byte("Contact update successful"))
		w.WriteHeader(http.StatusOK)
	}
}

func getAllContact(handler *contact.MongoHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contacts := handler.Get(bson.M{})

		_ = json.NewEncoder(w).Encode(contacts)
	}
}

func getContact(handler *contact.MongoHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		phoneNumber := chi.URLParam(r, "phonenumber")
		if phoneNumber == "" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)

			return
		}

		contact := &contact.Contact{}

		err := handler.GetOne(contact, bson.M{"phoneNumber": phoneNumber})
		if err != nil {
			http.Error(w, fmt.Sprintf("Contact with phonenumber: %s not found", phoneNumber), 404)

			return
		}

		_ = json.NewEncoder(w).Encode(contact)
	}
}

func addContact(handler *contact.MongoHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		existingContact := &contact.Contact{}

		var contact contact.Contact
		_ = json.NewDecoder(r.Body).Decode(&contact)
		contact.CreatedOn = time.Now()

		err := handler.GetOne(existingContact, bson.M{"phoneNumber": contact.PhoneNumber})
		if err == nil {
			http.Error(w, fmt.Sprintf("Contact with phonenumber: %s already exist", contact.PhoneNumber), http.StatusBadRequest)

			return
		}

		_, err = handler.AddOne(&contact)
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusBadRequest)

			return
		}

		_, _ = w.Write([]byte("Contact created successfully"))
		w.WriteHeader(http.StatusCreated)
	}
}

func registerRoutes(handler *contact.MongoHandler, router *chi.Mux) http.Handler {
	router.Route("/contacts", func(router chi.Router) {
		router.Get("/", getAllContact(handler))                 // GET /contacts
		router.Get("/{phonenumber}", getContact(handler))       // GET /contacts/0147344454
		router.Post("/", addContact(handler))                   // POST /contacts
		router.Put("/{phonenumber}", updateContact(handler))    // PUT /contacts/0147344454
		router.Delete("/{phonenumber}", deleteContact(handler)) // DELETE /contacts/0147344454
	})

	return router
}

func main() {
	mongoDBConnection := "mongodb://localhost:27017"
	mh := contact.NewHandler(mongoDBConnection) // Create an instance of MongoHander with the connection string provided

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	routerHandler := registerRoutes(mh, router)
	log.Fatal(http.ListenAndServe(":3060", routerHandler)) // You can modify to run on a different port
}
