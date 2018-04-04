package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"protobuf/server/definitions"

	coverer "github.com/sebastianosuna/s2-polygon-coverer"
	"google.golang.org/grpc"
)

var (
	host = "127.0.0.1"
	port = 2110
)

type polygonCovererServer struct {
}

func (s *polygonCovererServer) CoverPolygon(ctx context.Context, p *polygon.Polygon) (*polygon.CellIDList, error) {
	geo := geoJSONfromLatLng(p.GetCoordinates())
	cells := coverer.CoverPolygon(geo, int(p.GetLevel()))

	return &polygon.CellIDList{
		Messages: cells,
	}, nil
}

func newServer() *polygonCovererServer {
	s := &polygonCovererServer{}

	return s
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	polygon.RegisterPolygonCovererServer(grpcServer, newServer())
	grpcServer.Serve(lis)
}
