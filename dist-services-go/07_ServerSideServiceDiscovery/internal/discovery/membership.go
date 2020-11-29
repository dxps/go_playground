package discovery

import (
	"log"
	"net"

	"github.com/hashicorp/serf/serf"
)

func New(handler Handler, config Config) (*Membership, error) {

	m := &Membership{
		Config:  config,
		handler: handler,
	}
	if err := m.setupSerf(); err != nil {
		return nil, err
	}
	return m, nil
}

type Membership struct {
	Config
	handler Handler
	serf    *serf.Serf
	events  chan serf.Event
}

type Config struct {
	// Node's unique identifier across the Serf cluster. If not provided, the default value is the hostname.
	NodeName string
	// The listening address for gossip traffic.
	BindAddr string
	// Tags are used for associating different KV pairs with the node.
	Tags           map[string]string
	StartJoinAddrs []string
}

// Handler is a spec for a component that handles servers joining or leaving events.
type Handler interface {
	Join(name, addr string) error
	Leave(name, addr string) error
}

func (m *Membership) setupSerf() error {

	addr, err := net.ResolveTCPAddr("tcp", m.BindAddr)
	if err != nil {
		return err
	}
	cfg := serf.DefaultConfig()
	cfg.Init()
	cfg.MemberlistConfig.BindAddr = addr.IP.String()
	cfg.MemberlistConfig.BindPort = addr.Port

	m.events = make(chan serf.Event)
	cfg.EventCh = m.events
	cfg.Tags = m.Tags
	cfg.NodeName = m.NodeName

	m.serf, err = serf.Create(cfg)
	if err != nil {
		return err
	}
	go m.eventHandler()

	if m.StartJoinAddrs != nil {
		_, err := m.serf.Join(m.StartJoinAddrs, true)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Membership) eventHandler() {

	// A continuous event loop of handling cluster membership changes.
	for e := range m.events {
		switch e.EventType() {
		case serf.EventMemberJoin:
			for _, member := range e.(serf.MemberEvent).Members {
				if m.isLocal(member) {
					continue
				}
				m.handleJoin(member)
			}
		case serf.EventMemberLeave, serf.EventMemberFailed:
			for _, member := range e.(serf.MemberEvent).Members {
				if m.isLocal(member) {
					return
				}
				m.handleLeave(member)
			}
		}
	}
}

func (m *Membership) isLocal(member serf.Member) bool {
	return m.serf.LocalMember().Name == member.Name
}

func (m *Membership) Members() []serf.Member {
	return m.serf.Members()
}

func (m *Membership) handleJoin(member serf.Member) {

	if err := m.handler.Join(member.Name, member.Tags["rpc_addr"]); err != nil {
		log.Printf("[ERROR] Failed to join: %s %s", member.Name, member.Tags["rpc_addr"])
	}
}

func (m *Membership) handleLeave(member serf.Member) {

	if err := m.handler.Leave(member.Name, member.Tags["rpc_addr"]); err != nil {
		log.Printf("[ERROR] Failed to leave: %s %s", member.Name, member.Tags["rpc_addr"])
	}
}

func (m *Membership) Leave() error {
	return m.serf.Leave()
}
