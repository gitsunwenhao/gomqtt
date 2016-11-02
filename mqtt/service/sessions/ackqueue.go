package sessions

import (
	"errors"
	"fmt"
	"math"
	"sync"

	"github.com/aiyun/gomqtt/mqtt/protocol"
)

var (
	errQueueFull   = errors.New("queue full")
	errQueueEmpty  = errors.New("queue empty")
	errWaitMessage = errors.New("Invalid message to wait for ack")
	errAckMessage  = errors.New("Invalid message for acking")
)

type ackmsg struct {
	// Message type of the message waiting for ack
	Mtype protocol.PacketType

	// Current state of the ack-waiting message
	State protocol.PacketType

	// Packet ID of the message. Every message that require ack'ing must have a valid
	// packet ID. Messages that have message I
	Pktid uint16

	// Slice containing the message bytes
	Msgbuf []byte

	// Slice containing the ack message bytes
	Ackbuf []byte

	// When ack cycle completes, call this function
	OnComplete interface{}
}

// Ackqueue is a growing queue implemented based on a ring buffer. As the buffer
// gets full, it will auto-grow.
//
// Ackqueue is used to store messages that are waiting for acks to come back. There
// are a few scenarios in which acks are required.
//   1. Client sends SUBSCRIBE message to server, waits for SUBACK.
//   2. Client sends UNSUBSCRIBE message to server, waits for UNSUBACK.
//   3. Client sends PUBLISH QoS 1 message to server, waits for PUBACK.
//   4. Server sends PUBLISH QoS 1 message to client, waits for PUBACK.
//   5. Client sends PUBLISH QoS 2 message to server, waits for PUBREC.
//   6. Server sends PUBREC message to client, waits for PUBREL.
//   7. Client sends PUBREL message to server, waits for PUBCOMP.
//   8. Server sends PUBLISH QoS 2 message to client, waits for PUBREC.
//   9. Client sends PUBREC message to server, waits for PUBREL.
//   10. Server sends PUBREL message to client, waits for PUBCOMP.
//   11. Client sends PINGREQ message to server, waits for PINGRESP.
type Ackqueue struct {
	size  int64
	mask  int64
	count int64
	head  int64
	tail  int64

	ping ackmsg
	ring []ackmsg
	emap map[uint16]int64

	ackdone []ackmsg

	mu sync.Mutex
}

func newAckqueue(n int) *Ackqueue {
	m := int64(n)
	if !powerOfTwo64(m) {
		m = roundUpPowerOfTwo64(m)
	}

	return &Ackqueue{
		size:    m,
		mask:    m - 1,
		count:   0,
		head:    0,
		tail:    0,
		ring:    make([]ackmsg, m),
		emap:    make(map[uint16]int64, m),
		ackdone: make([]ackmsg, 0),
	}
}

// Wait copies the message into a waiting queue, and waits for the corresponding
// ack message to be received.
func (q *Ackqueue) Wait(msg protocol.Packet, onComplete interface{}) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	switch msg := msg.(type) {
	case *protocol.PublishPacket:
		if msg.QoS() == protocol.QosAtMostOnce {
			//return fmt.Errorf("QoS 0 messages don't require ack")
			return errWaitMessage
		}

		q.insert(msg.PacketID(), msg, onComplete)

	case *protocol.SubscribePacket:
		q.insert(msg.PacketID(), msg, onComplete)

	case *protocol.UnsubscribePacket:
		q.insert(msg.PacketID(), msg, onComplete)

	case *protocol.PingreqPacket:
		q.ping = ackmsg{
			Mtype:      protocol.PINGREQ,
			State:      protocol.RESERVED,
			OnComplete: onComplete,
		}

	default:
		return errWaitMessage
	}

	return nil
}

// Ack takes the ack message supplied and updates the status of messages waiting.
func (q *Ackqueue) Ack(msg protocol.Packet) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	switch msg.Type() {
	case protocol.PUBACK, protocol.PUBREC, protocol.PUBREL, protocol.PUBCOMP, protocol.SUBACK, protocol.UNSUBACK:
		// Check to see if the message w/ the same packet ID is in the queue
		i, ok := q.emap[msg.PacketID()]
		if ok {
			// If message w/ the packet ID exists, update the message state and copy
			// the ack message
			q.ring[i].State = msg.Type()

			ml := msg.Len()
			q.ring[i].Ackbuf = make([]byte, ml)

			_, err := msg.Encode(q.ring[i].Ackbuf)
			if err != nil {
				return err
			}
			//glog.Debugf("Acked: %v", msg)
			//} else {
			//glog.Debugf("Cannot ack %s message with packet ID %d", msg.Type(), msg.PacketId())
		}

	case protocol.PINGRESP:
		if q.ping.Mtype == protocol.PINGREQ {
			q.ping.State = protocol.PINGRESP
		}

	default:
		return errAckMessage
	}

	return nil
}

// Acked returns the list of messages that have completed the ack cycle.
func (q *Ackqueue) acked() []ackmsg {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.ackdone = q.ackdone[0:0]

	if q.ping.State == protocol.PINGRESP {
		q.ackdone = append(q.ackdone, q.ping)
		q.ping = ackmsg{}
	}

FORNOTEMPTY:
	for !q.empty() {
		switch q.ring[q.head].State {
		case protocol.PUBACK, protocol.PUBREL, protocol.PUBCOMP, protocol.SUBACK, protocol.UNSUBACK:
			q.ackdone = append(q.ackdone, q.ring[q.head])
			q.removeHead()

		default:
			break FORNOTEMPTY
		}
	}

	return q.ackdone
}

func (q *Ackqueue) insert(pktid uint16, msg protocol.Packet, onComplete interface{}) error {
	if q.full() {
		q.grow()
	}

	if _, ok := q.emap[pktid]; !ok {
		// message length
		ml := msg.Len()

		// ackmsg
		am := ackmsg{
			Mtype:      msg.Type(),
			State:      protocol.RESERVED,
			Pktid:      msg.PacketID(),
			Msgbuf:     make([]byte, ml),
			OnComplete: onComplete,
		}

		if _, err := msg.Encode(am.Msgbuf); err != nil {
			return err
		}

		q.ring[q.tail] = am
		q.emap[pktid] = q.tail
		q.tail = q.increment(q.tail)
		q.count++
	} else {
		// If packet w/ pktid already exist, then this must be a PUBLISH message
		// Other message types should never send with the same packet ID
		pm, ok := msg.(*protocol.PublishPacket)
		if !ok {
			return fmt.Errorf("ack/insert: duplicate packet ID for %s message", msg.Name())
		}

		// If this is a publish message, then the DUP flag must be set. This is the
		// only scenario in which we will receive duplicate messages.
		if pm.Dup() {
			return fmt.Errorf("ack/insert: duplicate packet ID for PUBLISH message, but DUP flag is not set")
		}

		// Since it's a dup, there's really nothing we need to do. Moving on...
	}

	return nil
}

func (q *Ackqueue) removeHead() error {
	if q.empty() {
		return errQueueEmpty
	}

	it := q.ring[q.head]
	// set this to empty ackmsg{} to ensure GC will collect the buffer
	q.ring[q.head] = ackmsg{}
	q.head = q.increment(q.head)
	q.count--
	delete(q.emap, it.Pktid)

	return nil
}

func (q *Ackqueue) grow() {
	if math.MaxInt64/2 < q.size {
		panic("new size will overflow int64")
	}

	newsize := q.size << 1
	newmask := newsize - 1
	newring := make([]ackmsg, newsize)

	if q.tail > q.head {
		copy(newring, q.ring[q.head:q.tail])
	} else {
		copy(newring, q.ring[q.head:])
		copy(newring[q.size-q.head:], q.ring[:q.tail])
	}

	q.size = newsize
	q.mask = newmask
	q.ring = newring
	q.head = 0
	q.tail = q.count

	q.emap = make(map[uint16]int64, q.size)

	for i := int64(0); i < q.tail; i++ {
		q.emap[q.ring[i].Pktid] = i
	}
}

func (q *Ackqueue) len() int {
	return int(q.count)
}

func (q *Ackqueue) cap() int {
	return int(q.size)
}

func (q *Ackqueue) index(n int64) int64 {
	return n & q.mask
}

func (q *Ackqueue) full() bool {
	return q.count == q.size
}

func (q *Ackqueue) empty() bool {
	return q.count == 0
}

func (q *Ackqueue) increment(n int64) int64 {
	return q.index(n + 1)
}

func powerOfTwo64(n int64) bool {
	return n != 0 && (n&(n-1)) == 0
}

func roundUpPowerOfTwo64(n int64) int64 {
	n--
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	n |= n >> 32
	n++

	return n
}
