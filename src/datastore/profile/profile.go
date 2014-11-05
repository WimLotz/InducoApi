package profile

import (
	"datastore"
	"log"

	"gopkg.in/mgo.v2/bson"
)

type Profile struct {
	Id             bson.ObjectId `bson:"_id" json:"id"`
	UserId         bson.ObjectId `bson:"userId" json:"userId"`
	IsCompany      bool          `bson:"isCompany" json:"isCompany"`
	LookingForWork bool          `bson:"lookingForWork,omitempty" json:"lookingForWork,omitempty"`
	Hiring         bool          `bson:"hiring,omitempty" json:"hiring,omitempty"`
	Name           string        `bson:"name,omitempty" json:"name,omitempty"`
	Bio            string        `bson:"bio,omitempty" json:"bio,omitempty"`
	EmailAddress   string        `bson:"emailAddress,omitempty" json:"emailAddress,omitempty"`
	ContactNumber1 string        `bson:"contactNumber1,omitempty" json:"contactNumber1,omitempty"`
	ContactNumber2 string        `bson:"contactNumber2,omitempty" json:"contactNumber2,omitempty"`
	WorkExpTags    []string      `bson:"workExpTags,omitempty" json:"workExpTags,omitempty"`
	NeededExpTags  []string      `bson:"neededExpTags,omitempty" json:"neededExpTags,omitempty"`
	Found          bool          `bson:"-,omitempty" json:"-,omitempty"`
}

func New() *Profile {
	return &Profile{}
}

func (p *Profile) Save() {
	_, err := datastore.ProfileCollection.Upsert(bson.M{"_id": p.Id}, p)
	if err != nil {
		log.Printf("Unable to save record: %v\n", err)
	}
}

func (p *Profile) Fetch(userId bson.ObjectId) []Profile {
	var results []Profile
	err := datastore.ProfileCollection.Find(bson.M{"userId": userId}).All(&results)
	if err != nil {
		log.Printf("no records found: %v\n", err)
	}
	return results
}

func (p *Profile) All() []Profile {
	var results []Profile
	err := datastore.ProfileCollection.Find(nil).All(&results)
	if err != nil {
		log.Printf("No records found: %v\n", err)
	}
	return results
}
