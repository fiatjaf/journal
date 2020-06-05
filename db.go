package main

import (
	"encoding/json"

	"go.etcd.io/bbolt"
)

func save(bas []BatchAction) (errorType string, err error) {
	// do all the actions in a transaction
	err = db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("logs"))

		for _, action := range bas {
			if action.Delete {
				// delete
				err = bucket.Delete([]byte(action.Id))
				if err != nil {
					errorType = "save"
					return err
				}
			} else if action.Set != nil {
				// get id and pos
				err = action.Set.ApplyId(action.Id)
				if err != nil {
					errorType = "time"
					return err
				}

				// validate
				err = action.Set.Validate()
				if err != nil {
					errorType = "validate"
					return err
				}

				// save
				v, _ := json.Marshal(action.Set)
				err = bucket.Put([]byte(action.Id), v)
				if err != nil {
					errorType = "save"
					return err
				}
			}
		}

		return nil
	})

	return
}
