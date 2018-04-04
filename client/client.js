const messages = require('./polygon_pb');
const fs = require('fs');
const net = require('net');
const gRPC = require('grpc');

const SERVER_HOST = '127.0.0.1:2110';

const grpcApp = gRPC.load('../definitions/polygon.proto');
const client = new grpcApp.PolygonCoverer(
  SERVER_HOST,
  gRPC.credentials.createInsecure(),
);

function makeReq(payload) {
  client.CoverPolygon(payload, (err, response) => {
    console.log(response);
  });
}

const content = fs.readFileSync('./polygon.json').toString();
const polygon = {
  id: 42,
  level: 17,
  coordinates: JSON.parse(content)
    .coordinates[0]
    .map(c => ({ lat: c[1], lng: c[0] })),
};
makeReq(polygon);
