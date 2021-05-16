build:
	@mkdir -p bin
	go build -o bin/srv42 .

clean:
	@rm bin/*

compile: clean
	@for os in darwin linux; do \
		for arch in amd arm; do \
			GOOS=$$os GOARCH=$${arch}64 go build -o bin/srv42-$$os-$${arch}64 .; \
			tar czf bin/srv42-$$os-$${arch}64.tar.gz bin/srv42-$$os-$${arch}64; \
			sha256sum bin/srv42-$$os-$${arch}64.tar.gz > bin/srv42-$$os-$${arch}64.tar.gz.sha256; \
			echo $$os $${arch}64 compiled; \
		done \
	done

brew-sha256:
	@for os in darwin linux; do \
		for arch in amd arm; do \
			echo "\"$$os-$${arch}64\" => \"$$(sha256sum bin/srv42-$$os-$${arch}64.tar.gz | egrep -o '^\w+')\""; \
		done \
	done
