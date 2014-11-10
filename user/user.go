package user

import (
	"github.com/WimLotz/InducoApi/datastore"
	"gopkg.in/mgo.v2/bson"
	"log"
)

type User struct {
	Id       bson.ObjectId `bson:"_id" json:"id"`
	Email    string        `bson:"email" json:"email"`
	Password string        `bson:"password" json:"password"`
}

func New() User {
	return User{}
}

func (u *User) Save() {
	_, err := datastore.UsersCollection.Upsert(bson.M{"_id": u.Id}, u)
	if err != nil {
		log.Printf("Unable to save record: %v\n", err)
	}
}

func (u *User) FetchOnEmail() error {
	err := datastore.UsersCollection.Find(bson.M{"email": u.Email}).One(u)
	if err != nil {
		log.Printf("User %v\n", err)
		return err
	}

	return nil
}

func (u *User) IsSuppliedPasswordCorrect(suppliedPassword string) bool {
	return suppliedPassword == u.Password
}
