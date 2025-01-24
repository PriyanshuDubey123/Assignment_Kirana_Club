package store

import (
	"encoding/csv"
	"log"
	"os"
	"strings"
	"sync"
)

type StoreManager struct {
	storeIds map[string]struct{}
}

var (
	once          sync.Once
	storeinstance *StoreManager
)

func NewStoreManager() (*StoreManager, error) {
	var err error
	once.Do(func() {
		storeinstance = &StoreManager{
			storeIds: make(map[string]struct{}),
		}
		err = storeinstance.LoadStoreIds()
	})
	return storeinstance, err
}

func (sm *StoreManager) LoadStoreIds() error {
	// Log to check if the function is starting
	log.Println("Loading store IDs from CSV file...")

	// Open the file
	file, err := os.Open(os.Getenv("CSVFILEPATH"))
	if err != nil {
		log.Printf("Error opening file: %v\n", err)  // Log the error if file can't be opened
		return err
	}
	defer file.Close()

	// Log after file is successfully opened
	log.Println("CSV file opened successfully.")

	// Read the CSV content
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Printf("Error reading CSV file: %v\n", err)  // Log the error if reading the file fails
		return err
	}

	// Log the number of records read from the CSV
	log.Printf("Successfully read %d records from the CSV file.\n", len(records))

	// Check if there are records and skip the header
	if len(records) > 0 {
		for _, record := range records[1:] {  // Skipping the header
			if len(record) > 2 {
				storeID := strings.TrimSpace(record[2])
				sm.storeIds[storeID] = struct{}{}
			}
		}
	}
	log.Println("All records read and stored in StoreManager.")
	return nil
}


func (sm *StoreManager) CheckStoreIDExist(store_id string) bool {
	_, exists := sm.storeIds[store_id]
	return exists
}
