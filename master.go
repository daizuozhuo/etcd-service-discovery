package discovery

import (
	"log"
	"time"

	"github.com/coreos/etcd/Godeps/_workspace/src/golang.org/x/net/context"
	"github.com/coreos/etcd/client"
	"encoding/json"
)

type Master struct {
	members  map[string]*Member
	KeysAPI  client.KeysAPI
}

// Member is a client machine
type Member struct {
	InGroup bool
	IP      string
	Name    string
	CPU    int
}

func NewMaster(endpoints []string) *Master {
	cfg := client.Config{
		Endpoints:               endpoints,
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}

	etcdClient, err := client.New(cfg)
	if err != nil {
		log.Fatal("Error: cannot connec to etcd:", err)
	}

	master := &Master{
		members: make(map[string]*Member),
		KeysAPI: client.NewKeysAPI(etcdClient),
	}
	go master.WatchWorkers()
	return master
}

func (m *Master) AddWorker(info *WorkerInfo) {
	member := &Member{
		InGroup: true,
		IP:      info.IP,
		Name:    info.Name,
		CPU:    info.CPU,
	}
	m.members[member.Name] = member
}

func (m *Master) UpdateWorker(info *WorkerInfo) {
	member := m.members[info.Name]
	member.InGroup = true
}

func (m *Master) WatchWorkers() {
	api := m.KeysAPI
	watcher := api.Watcher("workers/", &client.WatcherOptions{
		Recursive: true,
	})
	for {
		res, err := watcher.Next(context.Background())
		if err != nil {
			log.Println("Error watch workers:", err)
			break
		}
		if res.Action == "expire" {
			member, ok := m.members[res.Node.Key]
			if ok {
				member.InGroup = false
			}
		} else if res.Action == "set" || res.Action == "update"{
			info := &WorkerInfo{}
			err := json.Unmarshal([]byte(res.Node.Value), info)
			if err != nil {
				log.Print(err)
			}
			if _, ok := m.members[info.Name]; ok {
				m.UpdateWorker(info)
			} else {
				m.AddWorker(info)
			}
		} else if res.Action == "delete" {
			delete(m.members, res.Node.Key)
		}
	}

}
