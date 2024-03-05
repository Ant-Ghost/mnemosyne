package database

import (
	"log"
	"net"
	"os"
	"path"
	"strings"
	"time"

	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
)

// Replicate implements the Replication Service of the Server.
// It creates LogStore, StableStore and Snapstore for the current Server Node
// It creates a Raft Node and Boostraps Raft Cluster with other Server Nodes
// It returns a Raft Node and a Machine(KeyValueFSM).
func Replicate() (*raft.Raft, InternalFSM) {

	// Ensuring Directory exists for storing Raft data
	raftStorePath := os.Getenv("RAFT_STORE_PATH")
	if _, err := os.Stat(raftStorePath); os.IsNotExist(err) {
		if os.Mkdir(raftStorePath, 0755) != nil {
			log.Fatal("failed to create Raft path")
		}
	}

	// Address of this node
	nodeAddress := os.Getenv("NODE_ADDRESS")

	// Name of this node
	nodeIdentifier := os.Getenv("NODE_IDENTIFIER")

	// Addresses of all the other nodes in the network
	nodeAddressesString := os.Getenv("ALL_NODE_ADDRESSES")
	nodeAddresses := strings.Split(nodeAddressesString, ",")
	if len(nodeAddresses) == 0 {
		log.Fatal("no Raft node addresses specified")
	}

	// Reconfigure Default Raft Configuration
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(nodeIdentifier)
	config.ElectionTimeout = 2 * time.Second
	config.HeartbeatTimeout = 500 * time.Millisecond

	// Creating LogStore for RAFT
	logStore, err := raftboltdb.NewBoltStore(path.Join(raftStorePath, "raft-log"))
	if err != nil {
		log.Fatal(err)
	}

	// Creating StableStore for storing the Write Logs for more data durability and integrity
	stableStore, err := raftboltdb.NewBoltStore(path.Join(raftStorePath, "raft-stable"))
	if err != nil {
		log.Fatal(err)
	}

	// Resolving the addresses of current node
	addr, err := net.ResolveTCPAddr("tcp", nodeAddress)
	if err != nil {
		log.Fatal(err)
	}

	// Creating a transport Layer for the Raft cluster based on TCP/IP communication
	transport, err := raft.NewTCPTransport(addr.String(), addr, 3, 10*time.Second, os.Stderr)
	if err != nil {
		log.Fatal(err)
	}

	// Creating a state of the Machine
	fsmStore := NewFSM()

	// Creating SnapshotStore for storing the snapshots
	snapshotStore, err := raft.NewFileSnapshotStore(path.Join(raftStorePath, "raft-snapshots"), 1, os.Stderr)
	if err != nil {
		log.Fatal(err)
	}

	// Creating a new Raft Node
	raftNode, err := raft.NewRaft(config, fsmStore, logStore, stableStore, snapshotStore, transport)
	if err != nil {
		log.Fatal(err)
	}

	// Creating a slice of Raft Servers from the given node addresses.
	var servers []raft.Server
	for _, nodeAddr := range nodeAddresses {
		s := strings.Split(nodeAddr, "@")
		if len(s) != 2 {
			log.Fatalf("invalid Raft node address: %s", nodeAddr)
		}
		servers = append(servers, raft.Server{
			ID:      raft.ServerID(s[0]),
			Address: raft.ServerAddress(s[1]),
		})
	}

	// Boostrapping Raft Cluster with Raft Servers
	f := raftNode.BootstrapCluster(raft.Configuration{
		Servers: servers,
	})
	if err := f.Error(); err != nil {
		log.Println(err)
	}

	return raftNode, fsmStore
}
