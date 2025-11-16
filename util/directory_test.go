package util

import (
	"fmt"
	"testing"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

func TestGetParentPathsWithRoot(t *testing.T) {
	fmt.Println(GetParentPathsWithRoot("/datasix/cdc/mysql/", true))
}

func TestCreateZkNode(t *testing.T) {
	conn, event, err := zk.Connect([]string{"127.0.0.1:2181"}, time.Second)
	if err != nil {
		t.Fatal(err)
	}

	for {
		isConnected := false
		select {
		case connEvent := <-event:
			if connEvent.State == zk.StateConnected {
				isConnected = true
				fmt.Printf("connect to zookeeper server success!")
			}
		case <-time.After(time.Second * 3):
			t.Fatalf("connect to zookeeper server timeout")
		}
		if isConnected {
			break
		}
	}
	if err = CreateZkNode("/datasix/cdc1/mysql1", conn); err != nil {
		t.Fatal(err)
	}
}
