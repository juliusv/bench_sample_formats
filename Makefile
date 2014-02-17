proto: data.proto
	protoc --proto_path=. --cpp_out=cpp_proto_bench --go_out=go_proto_bench data.proto

bench_go_proto: proto
	make -C go_proto_bench

bench_cpp_proto: proto
	make -C cpp_proto_bench

bench_go_gob:
	make -C go_gob_bench
