.PHONY : 

export
local=localhost
docker=172.17.0.2

test_il: 
	go run helper/test_integration/test_integration.go \
	-server_name='${local}' 

test_id: 
	go run helper/test_integration/test_integration.go \
	-server_name='${docker}' 

