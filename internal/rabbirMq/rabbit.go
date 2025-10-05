package rabbirmq

import (
	"log"

	"github.com/wb-go/wbf/rabbitmq"
)
"github.com/wb-go/wbf/config"

type Brocker struct {

}

func New() *Brocker {
conn, err:=rabbitmq.Connect(rabbitUri, retries, pause)
if err!=nil{
	log.Fatal("invalid connection to brocker (rebbitMQ)")
}

ch, err:=conn.Channel()
if err!=nil{
	log.Fatal("invalid open Channel (rebbitMQ)")
}
// созадние обменника
ex:=rabbit.NewExchange("exenger", "direct")
ex.Durable = true
if
}
