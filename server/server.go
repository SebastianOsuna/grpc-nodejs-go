package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"protobuf/server/definitions"

	"github.com/golang/protobuf/proto"
	"github.com/sebastianosuna/s2-polygon-coverer"
)

// Req represents a TCP request
type Req struct {
	connection net.Conn
	body       polygon.Polygon
}

func (r *Req) getConnection() net.Conn {
	return r.connection
}

func (r *Req) getBody() polygon.Polygon {
	return r.body
}

// InitTCPServer starts a tcp server listening for protobufer
func InitTCPServer(host string, port int) {
	println("Started ProtoBuf Server")
	connChannel := make(chan *Req)

	go func() {
		for {
			handleRequest(<-connChannel)
		}
	}()

	listenOn := fmt.Sprintf("%s:%d", host, port)
	listener, err := net.Listen("tcp", listenOn)
	checkError(err)

	for {
		if conn, err := listener.Accept(); err == nil {
			go handleProtoClient(conn, connChannel)
		} else {
			continue
		}
	}
}

func handleRequest(req *Req) {
	handlePolygonMessage(req.getBody(), req.getConnection())
}

func handleProtoClient(conn net.Conn, channel chan *Req) {
	fmt.Println("Connection established")
	data := make([]byte, 4096)
	n, err := conn.Read(data)
	checkError(err)

	protodata := new(polygon.Polygon)
	err = proto.Unmarshal(data[0:n], protodata)
	checkError(err)

	channel <- &Req{
		body:       *protodata,
		connection: conn,
	}
}

func geoJSONfromLatLng(coords []*polygon.Polygon_Coordinate) *coverer.GeoJSON {
	latlngs := make([][]float64, len(coords))
	for i, coord := range coords {
		latlngs[i] = []float64{coord.GetLng(), coord.GetLat()}
	}

	return &coverer.GeoJSON{
		FeatureType:    "Polygon",
		RawCoordinates: [][][]float64{latlngs},
	}
}

func handlePolygonMessage(p polygon.Polygon, conn net.Conn) {
	coords := p.GetCoordinates()
	geo := geoJSONfromLatLng(coords)
	cells := coverer.CoverPolygon(geo, int(p.GetLevel()))
	msg := &polygon.CellIDList{Messages: cells}
	out, err := proto.Marshal(msg)

	if err != nil {
		log.Fatal(err)
	}

	conn.Write(out)
	conn.Close()
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s\n", err.Error())
	}
}
