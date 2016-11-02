package sessions

import (
	"fmt"
	"sync"

	"github.com/aiyun/gomqtt/mqtt/protocol"
)

const (
	// Queue size for the ack queue
	defaultQueueSize = 16
)

// Session is a data stucture for storing users infos
type Session struct {
	id string

	PubAckQueue *Ackqueue

	// cmsg is the CONNECT message
	Cmsg *protocol.ConnectPacket

	// Will message to publish if connect is closed unexpectedly
	Will *protocol.PublishPacket

	// Retained publish message
	Retained []*protocol.PublishPacket

	// cbuf is the CONNECT message buffer, this is for storing all the will stuff
	cbuf []byte

	// topics stores all the topis for this session/client
	topics map[string]byte

	// Serialize access to this session
	sync.Mutex
}

// Init a new session when server received the CONNECT packet
func (s *Session) Init(msg *protocol.ConnectPacket) error {
	s.Lock()
	defer s.Unlock()

	s.cbuf = make([]byte, msg.Len())
	// 这里注意！！！创建一个新的ConnectPacket是为了保证之前传入的Packet得到释放
	s.Cmsg = protocol.NewConnectPacket()

	if _, err := msg.Encode(s.cbuf); err != nil {
		return err
	}

	if _, err := s.Cmsg.Decode(s.cbuf); err != nil {
		return err
	}

	if s.Cmsg.WillFlag() {
		s.Will = protocol.NewPublishPacket()
		s.Will.SetQoS(s.Cmsg.WillQos())
		s.Will.SetTopic(s.Cmsg.WillTopic())
		s.Will.SetPayload(s.Cmsg.WillMessage())
		s.Will.SetRetain(s.Cmsg.WillRetain())
	}

	s.topics = make(map[string]byte, 1)

	s.id = string(msg.ClientId())

	return nil
}

// Update the session when server received the ConnectPacket
func (s *Session) Update(msg *protocol.ConnectPacket) error {
	s.Lock()
	defer s.Unlock()

	s.cbuf = make([]byte, msg.Len())
	s.Cmsg = protocol.NewConnectPacket()

	if _, err := msg.Encode(s.cbuf); err != nil {
		return err
	}

	if _, err := s.Cmsg.Decode(s.cbuf); err != nil {
		return err
	}

	return nil
}

// RetainMessage when server received PublishPacket
func (s *Session) RetainMessage(msg *protocol.PublishPacket) error {
	s.Lock()
	defer s.Unlock()

	return nil
}

// AddTopic to the session
func (s *Session) AddTopic(topic string, qos byte) error {
	s.Lock()
	defer s.Unlock()

	if !s.initted {
		return fmt.Errorf("Session not yet initialized")
	}

	s.topics[topic] = qos

	return nil
}

// RemoveTopic from the session
func (s *Session) RemoveTopic(topic string) error {
	s.Lock()
	defer s.Unlock()

	if !s.initted {
		return fmt.Errorf("Session not yet initialized")
	}

	delete(s.topics, topic)

	return nil
}

// Topics returns all topics in the session
func (s *Session) Topics() ([]string, []byte, error) {
	s.Lock()
	defer s.Unlock()

	if !s.initted {
		return nil, nil, fmt.Errorf("Session not yet initialized")
	}

	var (
		topics []string
		qoss   []byte
	)

	for k, v := range s.topics {
		topics = append(topics, k)
		qoss = append(qoss, v)
	}

	return topics, qoss, nil
}

// ID return the session id
func (s *Session) ID() string {
	return string(s.Cmsg.ClientId())
}
