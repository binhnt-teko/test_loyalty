package main

import (
	"os"
	"strings"
	"context"
	"log"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/p2p"
	"fmt"
	"github.com/binhnt-teko/test_loyalty/src/config"
)

func main(){
	var cfg *config.Config
	config_file := "config_local.yaml"
	if len(os.Args) == 2 {
		 fmt.Println("Please use syntax: go run addpeers.go  configfile ")
		 config_file = os.Args[1]
	}
	cfg = config.LoadConfigFromFile(config_file)

  ctx := context.Background()

  fmt.Println("Start get nodeInfo from each client ")
	method := "admin_nodeInfo"
	clients :=  make([]*rpc.Client,0)
	enodes := make([]string,0)

	for _, node := range cfg.Networks {
				fmt.Println("Connect host: ",node.Http)
				var result p2p.NodeInfo
				client, err := rpc.DialContext(ctx, "http://"+ node.Http)
			  if err != nil {
			    log.Fatalf("Unable to connect to network:%v\n", err)
			  }
				client.CallContext(ctx, &result, method)
				enode := strings.Replace(result.Enode,"127.0.0.1",node.LocalAddr,-1)
				enode = strings.Replace(enode,"?discport=0","",-1)
				fmt.Println("Enode: ", enode)
				enodes = append(enodes,enode)
				clients = append(clients,client)
	}

	fmt.Println("Start add peers for each client ")
	method = "admin_addPeer"
	for id, client := range clients {
		 for id1, enode := range enodes {
			 if id != id1 {
				 var result  bool
				 client.CallContext(ctx, &result, method,enode)
				 if result {
					 	 fmt.Println("Add peer: ", enode)
				 }else{
					 	 fmt.Println("Failed adding peer: ", enode)
				 }
			 }
		 }
	}
}
