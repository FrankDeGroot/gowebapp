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
	p, err := producer.Connect(topic)
	assert.NoError(t, err)
	defer p.Close()
	c, err := consumer.Connect(topic, topic)
	assert.NoError(t, err)
	defer c.Close()
	defer admin.DeleteTopic(topic)
	pTodo := &dto.SavedTodo{
		Id: "123",
		Todo: dto.Todo{
			Description: "test",
			Done:        false,
		},
	}
	err = p.Produce(act.Make(act.Add, pTodo))
	assert.NoError(t, err)
	cTodo, err := c.Consume()
	assert.NoError(t, err)
	assert.Equal(t, pTodo.Id, cTodo.Id)
	assert.Equal(t, pTodo.Description, cTodo.Description)
}
