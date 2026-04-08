n ?= 5
ip ?= "localhost"
types ?= "bpm"
typea ?= "vmi"
portudp ?= 8080
porttcp ?= 8000

.PHONY: broker sensor atuador cliente compose_sensor compose_atuador compose_cliente compose_broker

compose_sensor:
	for i in $$(seq 1 $(N)); do \
		docker compose run -d sensor ./app $(types) $$i $(ip); \
	done

compose_atuador:
	for i in $$(seq 1 $(N)); do \
		docker compose run -d atuador ./app $(typea) $$i $(ip); \
	done

compose_cliente:
	for i in $$(seq 1 $(N)); do \
		docker compose run cliente ./app $$i; $(ip)\
	done

compose_broker:
	docker compose up broker

sensor:
	cd sensor && for i in $$(seq 1 $(N)); do \
		docker run -d sensor ./app $(types) $$i $(ip); \
	done

atuador:
	cd atuador && for i in $$(seq 1 $(N)); do \
		docker run -d atuador ./app $(typea) $$i $(ip); \
	done

cliente:
	cd cliente && for i in $$(seq 1 $(N)); do \
		docker run cliente ./app $$i; $(ip)\
	done

broker:
	cd broker && docker run -p $(portudp):$(portudp) -p $(porttcp):$(porttcp) broker
