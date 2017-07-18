package sitereports

import (
	"fmt"

	"github.com/boltdb/bolt"
)

// Property struct
type Property struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Image       string   `json:"image"`
	Description string   `json:"description"`
	Sites       []string `json:"sites"`
	Teams       []string `json:"teams"`
	Users       []string `json:"users"`
}

// GetProperty function
func GetProperty(id string, db *bolt.DB) (Property, error) {
	var property Property

	err := db.View(func(tx *bolt.Tx) error {
		if rt := tx.Bucket([]byte("Properties")); rt != nil {
			if b := rt.Bucket([]byte(id)); b != nil {
				n := b.Get([]byte("name"))
				i := b.Get([]byte("image"))
				d := b.Get([]byte("description"))

				property.ID = string(id)
				property.Name = string(n)
				property.Image = string(i)
				property.Description = string(d)
			} else {
				return fmt.Errorf("No Property Found: %s", id)
			}
		} else {
			return fmt.Errorf("No Properties Bucket")
		}
		return nil
	})

	return property, err
}

// GetProperties func
func GetProperties(db *bolt.DB) []Property {
	var properties []Property
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Properties"))
		if b != nil {
			c := b.Cursor()

			for k, _ := c.First(); k != nil; k, _ = c.Next() {
				property, _ := GetProperty(string(k), db)
				properties = append(properties, property)
			}
		}
		return nil
	})
	if err != nil {
		return []Property{}
	}
	return properties
}

// Add func
func (property *Property) Add(db *bolt.DB) error {
	err := db.Update(func(tx *bolt.Tx) error {
		rt, err := tx.CreateBucketIfNotExists([]byte("Properties"))
		if err != nil {
			return fmt.Errorf("create properties bucket for: %s", property.ID)
		}
		b, err := rt.CreateBucketIfNotExists([]byte(property.ID))
		if err != nil {
			return fmt.Errorf("create property bucket: %s", property.ID)
		}
		err = b.Put([]byte("description"), []byte(property.Description))
		if err != nil {
			return fmt.Errorf("property description: %s", err)
		}
		err = b.Put([]byte("image"), []byte(property.Image))
		if err != nil {
			return fmt.Errorf("property image: %s", err)
		}
		err = b.Put([]byte("name"), []byte(property.Name))
		if err != nil {
			return fmt.Errorf("property name: %s", err)
		}
		return nil
	})

	return err
}

// Update function
func (property *Property) Update(db *bolt.DB) error {
	err := db.Update(func(tx *bolt.Tx) error {
		rt := tx.Bucket([]byte("Properties"))
		if rt == nil {
			return fmt.Errorf("create properties bucket for: %s", property.ID)
		}
		b := rt.Bucket([]byte(property.ID))
		if b == nil {
			return fmt.Errorf("create property bucket: %s", property.ID)
		}
		err := b.Put([]byte("description"), []byte(property.Description))
		if err != nil {
			return fmt.Errorf("property description: %s", err)
		}
		err = b.Put([]byte("image"), []byte(property.Image))
		if err != nil {
			return fmt.Errorf("property image: %s", err)
		}
		err = b.Put([]byte("name"), []byte(property.Name))
		if err != nil {
			return fmt.Errorf("property name: %s", err)
		}
		return nil
	})
	return err
}
