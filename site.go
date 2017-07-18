package sitereports

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

// Site struct
type Site struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	URL         string     `json:"url"`
	Image       string     `json:"image"`
	Description string     `json:"description"`
	Pages       []Page     `json:"pages"`
	Properties  []Property `json:"properties"`
	Reports     []Report   `json:"reports"`
}

// Page struct
type Page struct {
	ID          int      `json:"id"`
	URL         string   `json:"url"`
	Name        string   `json:"name"`
	Image       string   `json:"image"`
	Description string   `json:"description"`
	Reports     []Report `json:"reports"`
}

// GetSite function
func GetSite(id string, db *bolt.DB) (Site, error) {
	var site Site
	err := db.View(func(tx *bolt.Tx) error {
		if rt := tx.Bucket([]byte("Sites")); rt != nil {
			if b := rt.Bucket([]byte(id)); b != nil {
				n := b.Get([]byte("name"))
				i := b.Get([]byte("image"))
				d := b.Get([]byte("description"))
				u := b.Get([]byte("url"))

				site.ID = string(id)
				site.Name = string(n)
				site.Image = string(i)
				site.Description = string(d)
				site.URL = string(u)
				pgs := b.Bucket([]byte("Pages"))
				if pgs != nil {
					c := pgs.Cursor()

					for k, _ := c.First(); k != nil; k, _ = c.Next() {
						page := Page{}
						json.Unmarshal(k, page)
						log.Println(page)
						site.Pages = append(site.Pages, page)
					}
				}
			} else {
				return fmt.Errorf("No Site Found: %s", id)
			}
		} else {
			return fmt.Errorf("No Sites Bucket")
		}
		return nil
	})

	return site, err
}

// Update function
func (site *Site) Update(db *bolt.DB) error {
	err := db.Update(func(tx *bolt.Tx) error {
		if rt, err := tx.CreateBucketIfNotExists([]byte("Sites")); err != nil {
			return fmt.Errorf("create sites bucket: %s", err)
		} else if b, err := rt.CreateBucketIfNotExists([]byte(site.ID)); err != nil {
			return fmt.Errorf("create site bucket: %s", err)
		} else if err = b.Put([]byte("url"), []byte(site.URL)); err != nil {
			return fmt.Errorf("site url: %s", err)
		} else if err = b.Put([]byte("name"), []byte(site.Name)); err != nil {
			return fmt.Errorf("site name: %s", err)
		} else if err = b.Put([]byte("description"), []byte(site.Description)); err != nil {
			return fmt.Errorf("site description: %s", err)
		}
		return nil
	})

	return err
}

// Add func
func (site *Site) Add(db *bolt.DB) error {
	err := db.Update(func(tx *bolt.Tx) error {
		rt, err := tx.CreateBucketIfNotExists([]byte("Sites"))
		if err != nil {
			return fmt.Errorf("Cannot Create Sites Bucket")
		}
		b, err := rt.CreateBucketIfNotExists([]byte(site.ID))
		if err != nil {
			return fmt.Errorf("Could not create bucket for site: %s", err)
		}
		err = b.Put([]byte("description"), []byte(site.Description))
		if err != nil {
			return fmt.Errorf("Could not add site description: %s", err)
		}
		err = b.Put([]byte("image"), []byte(site.Image))
		if err != nil {
			return fmt.Errorf("Could not add site image: %s", err)
		}
		err = b.Put([]byte("name"), []byte(site.Name))
		if err != nil {
			return fmt.Errorf("Could not add site name: %s", err)
		}
		err = b.Put([]byte("url"), []byte(site.URL))
		if err != nil {
			return fmt.Errorf("Could not add site url: %s", err)
		}
		return nil
	})

	return err
}

// AddPage function
func (site *Site) AddPage(page Page, db *bolt.DB) error {
	err := db.Update(func(tx *bolt.Tx) error {
		rt := tx.Bucket([]byte("Sites"))
		if rt == nil {
			return fmt.Errorf("Sites Bucket does not exist")
		}

		b := rt.Bucket([]byte(site.ID))
		if b == nil {
			return fmt.Errorf("Site Bucket does not exist: %s", site.ID)
		}
		if pgs, err := b.CreateBucketIfNotExists([]byte("Pages")); err != nil {
			log.Println(page.ID)
			if page.ID == 0 {
				id, _ := pgs.NextSequence()
				page.ID = int(id)
			}

			buf, err := json.Marshal(page)
			log.Printf("Page JSON Buffer: %s", buf)
			if err != nil {
				return fmt.Errorf("Problem marshalling page: %s", err)
			}
			fmt.Println(page)
			return pgs.Put(itob(page.ID), buf)
		}

		return nil
	})

	return err
}

// RemovePage function
func (site *Site) RemovePage(page Page, db *bolt.DB) error {
	err := db.Update(func(tx *bolt.Tx) error {
		rt := tx.Bucket([]byte("Sites"))
		if rt == nil {
			return fmt.Errorf("Sites Bucket does not exist")
		}

		b := rt.Bucket([]byte(site.ID))
		if b == nil {
			return fmt.Errorf("Site Bucket does not exist: %s", site.ID)
		}
		if pgs, err := b.CreateBucketIfNotExists([]byte("Pages")); err != nil {
			log.Println(page.ID)
			pgs.Delete(itob(page.ID))
		}

		return nil
	})

	return err
}

// GetSites func
func GetSites(db *bolt.DB) []Site {
	var sites []Site
	err := db.View(func(tx *bolt.Tx) error {
		d := tx.Bucket([]byte("Sites"))
		if d != nil {
			e := d.Cursor()

			for k, _ := e.First(); k != nil; k, _ = e.Next() {
				site, _ := GetSite(string(k), db)
				sites = append(sites, site)
			}
		}
		return nil
	})
	if err != nil {
		return []Site{}
	}
	return sites
}

// itob returns an 8-byte big endian representation of v.
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
