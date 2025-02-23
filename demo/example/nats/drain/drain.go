package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/wwj31/dogactor/actor"
	"github.com/wwj31/dogactor/actor/cluster/mq"
	"github.com/wwj31/dogactor/actor/cluster/mq/nats"
	"github.com/wwj31/dogactor/demo/example/nats/msg"
	"github.com/wwj31/dogactor/logger"
	"github.com/wwj31/dogactor/tools"
)

var natsUrl = "nats://localhost:4222"

func main() {
	protoIndex := tools.NewProtoIndex(func(name string) (interface{}, bool) {
		return msg.Spawner(name)
	}, tools.EnumIdx{})

	system1, _ := actor.NewSystem(mq.WithRemote(natsUrl, nats.New()), actor.ProtoIndex(protoIndex))
	system2, _ := actor.NewSystem(mq.WithRemote(natsUrl, nats.New()), actor.ProtoIndex(protoIndex))

	name := "drainActor"
	drainActor1 := actor.BaseFunc{}
	drainActor2 := actor.BaseFunc{}
	drainActor1.Handle = func(v actor.Message) {
		drainData := v.Payload().(*msg.DrainTest)
		fmt.Println("drain actor 1 receive msg:", drainData.Data)
		time.Sleep(200 * time.Millisecond)
	}
	drainActor2.Handle = func(v actor.Message) {
		drainData := v.Payload().(*msg.DrainTest)
		fmt.Println("drain actor 2 receive msg:", drainData.Data)
		time.Sleep(100 * time.Millisecond)
	}

	_ = system1.NewActor(name, &drainActor1, actor.SetMailBoxSize(50))

	drainActor1.Init = func() {
		fmt.Println("init drain actor 1")
		drainActor1.AddTimer(tools.XUID(), tools.Now().Add(100*time.Millisecond), func(dt time.Duration) {
			_ = system2.Send("", name, "", &msg.DrainTest{Data: time.Now().String()})
		}, -1)

		drainActor1.AddTimer(tools.XUID(), tools.Now().Add(10*time.Second), func(dt time.Duration) {
			fmt.Println("drain 1")
			drainActor1.Drain(func() {
				fmt.Println("drain over drainActor1")
				_ = system2.NewActor(name, &drainActor2, actor.SetMailBoxSize(100))
			})
		})
	}

	drainActor2.Init = func() {
		drainActor2.AddTimer(tools.XUID(), tools.Now().Add(300*time.Millisecond), func(dt time.Duration) {
			_ = system1.Send("", name, "", &msg.DrainTest{Data: time.Now().String()})
		}, -1)

		drainActor2.AddTimer(tools.XUID(), time.Now().Add(10*time.Second), func(dt time.Duration) {
			fmt.Println("drain 2")
			drainActor2.Drain(func() {
				fmt.Println("drain over drainActor2")
				_ = system1.NewActor(name, &drainActor1, actor.SetMailBoxSize(100))
			})
		})
	}

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-c
	system1.Stop()
	system2.Stop()
	<-system1.Stopped
	<-system2.Stopped

	logger.Close()
}
