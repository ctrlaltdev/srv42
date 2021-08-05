_NAME=srv42

build:
	go build -o bin/$(_NAME) .

clean:
	rm bin/$(_NAME)*

compile: clean
	@for os in darwin linux; do \
		for arch in amd64 arm64; do \
			cd bin ; \
			GOOS=$$os GOARCH=$${arch} go build -o $(_NAME)-$$os-$${arch} ../; \
			tar czf $(_NAME)-$$os-$${arch}.tar.gz $(_NAME)-$$os-$${arch}; \
			sha256sum $(_NAME)-$$os-$${arch}.tar.gz > $(_NAME)-$$os-$${arch}.tar.gz.sha256; \
			echo $$os $${arch} compiled; \
			cd ../ ; \
		done \
	done

brew-sha256:
	@for os in darwin linux; do \
		for arch in amd64 arm64; do \
			cd bin ; \
			echo "\"$$os-$${arch}\" => \"$$(sha256sum $(_NAME)-$$os-$${arch}.tar.gz | egrep -o '^\w+')\""; \
			cd ../ ; \
		done \
	done
