package ws

import (
	"log"
	"todo-app/act"
)

func consume(cons Consumer, consChan chan *act.TaskAction, contChan chan bool) {
	cont, ok := true, true
	for {
		select {
		case cont, ok = <-contChan:
			if !cont || !ok {
				return
			}
		default:
			task, err := cons.Consume()
			if err != nil {
				log.Printf("Error consuming task: %v\n", err)
			}
			if task == nil {
				continue
			}
			log.Printf("Consumed %v", task)
			consChan <- task
		}
	}
}
