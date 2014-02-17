default: bench

build: proto_bench.cc data.pb.cc
	g++ proto_bench.cc data.pb.cc -o proto -lprotobuf -lrt -O3 -o proto_bench

bench: build
	./proto_bench
