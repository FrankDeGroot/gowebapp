package kafka

import (
	"fmt"
	"testing"
	"time"
	"todo-app/act"
	"todo-app/dto"
	"todo-app/kafka/admin"
	"todo-app/kafka/consumer"
	"todo-app/kafka/producer"

	"github.com/stretchr/testify/assert"
)

func TestProduceConsume(t *testing.T) {
	topic := fmt.Sprintf("todo%v", time.Now().Format("_2006_01_02_15_04_05"))
	t.Log(topic)
	p, err := producer.Open(topic)
	assert.NoError(t, err)
	defer p.Close()
	c, err := consumer.Open(topic, topic)
	assert.NoError(t, err)
	defer c.Close()
	defer admin.DeleteTopic(topic)
	pTask := &dto.SavedTask{
		Id: "123",
		Task: dto.Task{
			Description: "test",
			Done:        false,
		},
	}
	err = p.Produce(act.Make(act.Post, pTask))
	assert.NoError(t, err)
	cTask, err := c.Consume()
	assert.NoError(t, err)
	assert.Equal(t, pTask.Id, cTask.Id)
	assert.Equal(t, pTask.Description, cTask.Description)
}
