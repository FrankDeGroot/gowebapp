package ws

import (
	"context"
	"errors"
	"log"
	"todo-app/act"
	"todo-app/db"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

func read(conn *websocket.Conn, cont chan struct{}, repo db.TaskDber) {
	for {
		select {
		case _, ok := <-cont:
			if !ok {
				return
			}
		default:
			action := act.TaskAction{}
			err := wsjson.Read(context.Background(), conn, &action)
			if err != nil {
				closeErr, ok := errors.Unwrap(errors.Unwrap(errors.Unwrap(err))).(websocket.CloseError)
				if ok && closeErr.Code == websocket.StatusGoingAway {
					return
				}
				log.Printf("Error reading websocket %v", err)
				conn.CloseNow()
				return
			}
			log.Printf("Read %v", action)
			switch action.Verb {
			case act.Post:
				savedTask, err := repo.Insert(&action.Task)
				if err != nil {
					log.Printf("Error on insert %v", err)
				}
				notify(act.Make(act.Post, savedTask))
			case act.Put:
				err := repo.Update(&action.SavedTask)
				if err != nil {
					log.Printf("Error on update %v", err)
				}
				notify(act.Make(act.Put, &action.SavedTask))
			case act.Delete:
				err := repo.Delete(action.Id)
				if err != nil {
					log.Printf("Error on delete %v", err)
				}
				notify(act.Make(act.Delete, &action.SavedTask))
			case act.Get:
				if len(action.Id) == 0 {
					tasks, err := repo.GetAll()
					if err != nil {
						log.Printf("Error on get all %v", err)
					}
					err = wsjson.Write(context.Background(), conn, tasks)
					if err != nil {
						log.Printf("Error writing get all %v", err)
					}
				} else {
					task, err := repo.GetOne(action.SavedTask.Id)
					if err != nil {
						log.Printf("Error on get %v %v", action.Id, err)
					}
					err = wsjson.Write(context.Background(), conn, task)
					if err != nil {
						log.Printf("Error writing get one %v", err)
					}
				}
			}
		}
	}
}
