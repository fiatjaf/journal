package main

import (
	"encoding/json"

	"github.com/lucsky/cuid"
	"github.com/tidwall/buntdb"
	"github.com/tidwall/gjson"
)

func AppendEntry(le *LogEntry) (id string, err error) {
	lastPos, err := GetLastPosOnDate(le.Date)
	if err != nil {
		return
	}
	le.Pos = nextPos(lastPos)
	id = cuid.New()

	err = SaveBatch([]BatchAction{{id, le, false}})
	return
}

func SaveBatch(bas []BatchAction) (err error) {
	// do all the actions in a transaction
	return db.Update(func(tx *buntdb.Tx) error {
		for _, action := range bas {
			if action.Delete {
				// delete
				err = tx.Delete(action.Id)
				if err != nil {
					return err
				}
			} else if action.Set != nil {
				// validate
				err = action.Set.Validate()
				if err != nil {
					return err
				}

				// save
				v, _ := json.Marshal(action.Set)
				err = tx.Set(action.Id, string(v), nil)
				if err != nil {
					return err
				}

				return tx.Commit()
			}
		}

		return nil
	})
}

func GetLastPosOnDate(date string) (pos string, err error) {
	db.View(func(tx *buntdb.Tx) error {
		tx.AscendGreaterOrEqual("datepos", date, func(_, v string) bool {
			p := gjson.Parse(v)
			if p.Get("date").String() == date {
				pos = p.Get("pos").String()
				return true
			} else {
				return false
			}
		})
		return nil
	})
	return
}
