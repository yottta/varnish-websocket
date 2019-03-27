package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type requestHandler struct {
	id int

	dataMutex *sync.Mutex
	data      []string

	connectedClientsMutex *sync.Mutex
	connectedClients      map[int]*websocket.Conn
}

func main() {
	router := setupRoutes()
	log.Fatal(http.ListenAndServe(":8080", router))
}

func setupRoutes() *mux.Router {
	handler := &requestHandler{
		connectedClientsMutex: &sync.Mutex{},
		dataMutex:             &sync.Mutex{},
		connectedClients:      make(map[int]*websocket.Conn),
		data:                  []string{"product1", "product2"},
	}

	router := mux.NewRouter()
	router.HandleFunc("/product", handler.handleGetData).Methods(http.MethodGet)
	router.HandleFunc("/product", handler.handleAddData).Methods(http.MethodPost)
	router.HandleFunc("/ws", handler.wsEndpoint)

	return router
}

func (h *requestHandler) handleGetData(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	dat, err := json.Marshal(h.products())
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"products":%s}`, dat)
}

func (h *requestHandler) handleAddData(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	defer r.Body.Close()

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("Data added %s", data)
	h.data = append(h.data, string(data))
	h.writeToWebSockets(data)
}

func (h *requestHandler) wsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, "%s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("Client successfully connected to %d\n", h.id)
	h.register(ws)
}

func (h *requestHandler) addConnection(conn *websocket.Conn) int {
	h.connectedClientsMutex.Lock()
	defer h.connectedClientsMutex.Unlock()
	connID := len(h.connectedClients)

	h.connectedClients[connID] = conn
	return connID
}

func (h *requestHandler) removeConnection(connID int) {
	h.connectedClientsMutex.Lock()
	defer h.connectedClientsMutex.Unlock()

	h.connectedClients[connID].Close()
	delete(h.connectedClients, connID)
}

func (h *requestHandler) register(conn *websocket.Conn) {
	connID := h.addConnection(conn)
	defer h.removeConnection(connID)
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(string(p))

		newMsg := fmt.Sprintf("Message received from the client: %s", p)
		if err := conn.WriteMessage(messageType, []byte(newMsg)); err != nil {
			log.Println(err)
			return
		}
	}
}

func (h *requestHandler) writeToWebSockets(data []byte) {
	for connID, conn := range h.connectedClients {
		if conn == nil {
			continue
		}
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Printf("couldn't send message to %d connection", connID)
		}
	}
}

func (h *requestHandler) products() []string {
	h.dataMutex.Lock()
	defer h.dataMutex.Unlock()

	res := make([]string, len(h.data))
	copy(res, h.data)
	return res
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
