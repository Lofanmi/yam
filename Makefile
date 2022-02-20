PWD = `pwd`
ID = `git rev-parse --short HEAD`

env5:
	docker images | grep $(ID); if [ $$? -ne 0 ]; then cd env/php5 && docker build -t $(ID) . ; fi

dev5:
	@echo "docker run --rm -it -v $(PWD):/app $(ID) /bin/bash"

build5:
	docker run --rm -it -v $(PWD):/app $(ID) /app/build5.sh

bootstrap5:
	cd bootstrap/php5 && phpize && ./configure && make && cp modules/yam.so ../../bin/yam.so

clean:
	rm -rf $(PWD)/bin/yamgo.so

clean-all:
	rm -rf $(PWD)/bin/yamgo.so
	rm -rf $(PWD)/internal/php-src/5.6.40
	docker rmi -f $(ID)

test:
//