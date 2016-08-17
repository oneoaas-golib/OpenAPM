package udp_listener

import (
	"log"
	"net"
	"sync"
	"time"

	"github.com/corego/vgo/mecury/agent"
)

type UdpListener struct {
	ServiceAddress string
	// UDPPacketSize is deprecated, it's only here for legacy support
	// we now always create 1 max size buffer and then copy only what we need
	// into the in channel
	// see https://github.com/influxdata/telegraf/pull/992
	UDPPacketSize          int `toml:"udp_packet_size"`
	AllowedPendingMessages int

	sync.Mutex
	wg sync.WaitGroup

	in   chan []byte
	done chan struct{}
	// drops tracks the number of dropped metrics.
	drops int
	// malformed tracks the number of malformed packets
	malformed int

	parser agent.Parser

	// Keep the accumulator in this struct
	acc agent.Accumulator

	listener *net.UDPConn
}

// UDP packet limit, see
// https://en.wikipedia.org/wiki/User_Datagram_Protocol#Packet_structure
const UDP_MAX_PACKET_SIZE int = 64 * 1024

var dropwarn = "ERROR: udp_listener message queue full. " +
	"We have dropped %d messages so far. " +
	"You may want to increase allowed_pending_messages in the config\n"

var malformedwarn = "WARNING: udp_listener has received %d malformed packets" +
	" thus far."

const sampleConfig = `
  ## Address and port to host UDP listener on
  service_address = ":8092"

  ## Number of UDP messages allowed to queue up. Once filled, the
  ## UDP listener will start dropping packets.
  allowed_pending_messages = 10000

  ## Data format to consume.
  ## Each data format has it's own unique set of configuration options, read
  ## more about them here:
  ## https://github.com/influxdata/telegraf/blob/master/docs/DATA_FORMATS_INPUT.md
  data_format = "influx"
`

func (u *UdpListener) SampleConfig() string {
	return sampleConfig
}

func (u *UdpListener) Description() string {
	return "Generic UDP listener"
}

// All the work is done in the Start() function, so this is just a dummy
// function.
func (u *UdpListener) Gather(_ agent.Accumulator) error {
	return nil
}

func (u *UdpListener) SetParser(parser agent.Parser) {
	u.parser = parser
}

func (u *UdpListener) Start(acc agent.Accumulator) error {
	u.Lock()
	defer u.Unlock()

	u.acc = acc
	u.in = make(chan []byte, u.AllowedPendingMessages)
	u.done = make(chan struct{})

	u.wg.Add(2)
	go u.udpListen()
	go u.udpParser()

	log.Printf("Started UDP listener service on %s\n", u.ServiceAddress)
	return nil
}

func (u *UdpListener) Stop() {
	u.Lock()
	defer u.Unlock()
	close(u.done)
	u.wg.Wait()
	u.listener.Close()
	close(u.in)
	log.Println("Stopped UDP listener service on ", u.ServiceAddress)
}

func (u *UdpListener) udpListen() error {
	defer u.wg.Done()
	var err error
	address, _ := net.ResolveUDPAddr("udp", u.ServiceAddress)
	u.listener, err = net.ListenUDP("udp", address)
	if err != nil {
		log.Fatalf("ERROR: ListenUDP - %s", err)
	}
	log.Println("UDP server listening on: ", u.listener.LocalAddr().String())

	buf := make([]byte, UDP_MAX_PACKET_SIZE)
	for {
		select {
		case <-u.done:
			return nil
		default:
			u.listener.SetReadDeadline(time.Now().Add(time.Second))
			n, _, err := u.listener.ReadFromUDP(buf)
			if err != nil {
				if err, ok := err.(net.Error); ok && err.Timeout() {
				} else {
					log.Printf("ERROR: %s\n", err.Error())
				}
				continue
			}
			bufCopy := make([]byte, n)
			copy(bufCopy, buf[:n])

			select {
			case u.in <- bufCopy:
			default:
				u.drops++
				if u.drops == 1 || u.drops%u.AllowedPendingMessages == 0 {
					log.Printf(dropwarn, u.drops)
				}
			}
		}
	}
}

func (u *UdpListener) udpParser() error {
	defer u.wg.Done()

	var packet []byte
	var metrics []agent.Metric
	var err error
	for {
		select {
		case <-u.done:
			if len(u.in) == 0 {
				return nil
			}
		case packet = <-u.in:
			metrics, err = u.parser.Parse(packet)
			if err == nil {
				for _, m := range metrics {
					u.acc.AddFields(m.Name(), m.Fields(), m.Tags(), m.Time())
				}
			} else {
				u.malformed++
				if u.malformed == 1 || u.malformed%1000 == 0 {
					log.Printf(malformedwarn, u.malformed)
				}
			}
		}
	}
}

func init() {
	agent.AddInput("udp_listener", &UdpListener{})
}