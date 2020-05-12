package sse

import (
	"bytes"
	"culture/models"
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"time"
)

//SseClient SseClient
type SseClient chan []byte

// Streamer receives events and broadcasts them to all connected clients.
// Streamer is a http.Handler. Clients making a request to this handler receive
// a stream of Server-Sent Events, which can be handled via JavaScript.
// See the linked technical specification for details.
type Streamer struct {
	Event         chan []byte
	Clients       map[SseClient]bool
	Connecting    chan SseClient
	Disconnecting chan SseClient
	BufSize       uint
}

//S S
var S *Streamer
var ticker = time.NewTicker(time.Second * 50)

// Init returns a new initialized SSE Streamer
func init() {
	S = &Streamer{
		Event:         make(chan []byte, 1),
		Clients:       make(map[SseClient]bool),
		Connecting:    make(chan SseClient),
		Disconnecting: make(chan SseClient),
		BufSize:       2,
	}
	S.run()
	go func() {
		for {
			select {
			case <-ticker.C:
				S.SendString("", "message", "ping")
			}
		}
	}()
}

// run starts a goroutine to handle client connects and broadcast events.
func (s *Streamer) run() {
	go func() {
		for {
			select {
			case cl := <-s.Connecting:
				s.Clients[cl] = true
				go initStatistics()

			case cl := <-s.Disconnecting:
				delete(s.Clients, cl)

			case event := <-s.Event:
				for cl := range s.Clients {
					// TODO: non-blocking broadcast
					//select {
					//case cl <- event: // Try to send event to client
					//default:
					//	fmt.Println("Channel full. Discarding value")
					//}
					cl <- event
				}
			}
		}
	}()
}

// SetBufSize sets the event buffer size for new clients.
func (s *Streamer) SetBufSize(size uint) {
	s.BufSize = size
}

func format(id, event string, dataLen int) (p []byte) {
	// calc length
	l := 6 // data\n\n
	if len(event) > 0 {
		l += 6 + len(event) + 1 // event:{event}\n
	}
	if dataLen > 0 {
		l += 1 + dataLen // :{data}
	}

	// build
	p = make([]byte, l)
	i := 0
	if len(event) > 0 {
		copy(p, "event:")
		i += 6 + copy(p[6:], event)
		p[i] = '\n'
		i++
	}
	i += copy(p[i:], "data")
	if dataLen > 0 {
		p[i] = ':'
		i += 1 + dataLen
	}
	copy(p[i:], "\n\n")

	// TODO: id

	return
}

// SendBytes sends an event with the given byte slice interpreted as a string
// as the data value to all connected clients.
// If the id or event string is empty, no id / event type is send.
func (s *Streamer) SendBytes(id, event string, data []byte) {
	dataLen := len(data)
	lfCount := 0

	// We must sent a "data:{data}\n" for each line
	if dataLen > 0 {
		lfCount = bytes.Count(data, []byte("\n"))
		if lfCount > 0 {
			dataLen += 5 * lfCount // data:
		}
	}

	p := format(id, event, dataLen)

	// fill in data lines
	start := 0
	ins := len(p) - (2 + dataLen)
	for i := 0; lfCount > 0; i++ {
		if data[i] == '\n' {
			copy(p[ins:], data[start:i])
			ins += i - start
			copy(p[ins:], "\ndata:")
			ins += 6

			start = i + 1
			lfCount--
		}
	}
	copy(p[ins:], data[start:])

	s.Event <- p
}

// SendInt sends an event with the given int as the data value to all connected
// clients.
// If the id or event string is empty, no id / event type is send.
func (s *Streamer) SendInt(id, event string, data int64) {
	const maxIntToStrLen = 20 // '-' + 19 digits

	p := format(id, event, maxIntToStrLen)
	p = strconv.AppendInt(p[:len(p)-(maxIntToStrLen+2)], data, 10)

	// Re-add \n\n at the end
	p = p[:len(p)+2]
	p[len(p)-2] = '\n'
	p[len(p)-1] = '\n'

	s.Event <- p
}

// SendJSON sends an event with the given data encoded as JSON to all connected
// clients.
// If the id or event string is empty, no id / event type is send.
func (s *Streamer) SendJSON(id, event string, v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	p := format(id, event, len(data))
	copy(p[len(p)-(2+len(data)):], data) // fill in data
	s.Event <- p
	return nil
}

// SendString sends an event with the given data string to all connected
// clients.
// If the id or event string is empty, no id / event type is send.
func (s *Streamer) SendString(id, event, data string) {
	dataLen := len(data)
	lfCount := 0

	// We must sent a "data:{data}\n" for each line
	if dataLen > 0 {
		lfCount = strings.Count(data, "\n")
		if lfCount > 0 {
			dataLen += 5 * lfCount // data:
		}
	}

	p := format(id, event, dataLen)

	// fill in data lines
	start := 0
	ins := len(p) - (2 + dataLen)
	for i := 0; lfCount > 0; i++ {
		if data[i] == '\n' {
			copy(p[ins:], data[start:i])
			ins += i - start
			copy(p[ins:], "\ndata:")
			ins += 6

			start = i + 1
			lfCount--
		}
	}
	copy(p[ins:], data[start:])

	s.Event <- p
}

// SendUint sends an event with the given unsigned int as the data value to all
// connected clients.
// If the id or event string is empty, no id / event type is send.
func (s *Streamer) SendUint(id, event string, data uint64) {
	const maxUintToStrLen = 20

	p := format(id, event, maxUintToStrLen)
	p = strconv.AppendUint(p[:len(p)-(maxUintToStrLen+2)], data, 10)

	// Re-add \n\n at the end
	p = p[:len(p)+2]
	p[len(p)-2] = '\n'
	p[len(p)-1] = '\n'

	s.Event <- p
}

func initStatistics() {
	go projectsCount()
	go postStatistics()
	go mapStatistics(1)
	go mapStatistics(2)
}

func projectsCount() {
	projects := &models.Projects{}
	if count, err := projects.Count(); err == nil {
		S.SendInt("", "projects_count", int64(count))
	}
}

//MapStatistics MapStatistics
func MapStatistics() {
	go mapStatistics(1)
	go mapStatistics(2)
}

//ProjectsCount ProjectsCount
func ProjectsCount() {
	projectsCount()
}

func mapStatistics(level int) {
	statistics := &models.MapStatistics{
		Level: level,
	}
	err := statistics.Scan()
	if err != nil {
		log.Printf("mapStatistics failed : %v", err)
		return
	}
	if err = S.SendJSON("", "map", statistics); err != nil {
		log.Printf("send mapStatistics failed : %v", err)
	}
}

//PostStatistics PostStatistics
func PostStatistics() {
	postStatistics()
}

func postStatistics() {
	statistics := &models.PostStatistics{
		Posts: &models.Posts{},
	}
	err := statistics.Statistics()
	if err != nil {
		log.Printf("statistics posts failed : %v", err)
		return
	}
	S.SendInt("", "posts_count", statistics.Count)
	S.SendInt("", "today_posts_count", statistics.TodayCount)
	if err := S.SendJSON("", "posts", statistics.Posts); err != nil {
		log.Printf("send statistics posts failed : %v", err)
	}
}
