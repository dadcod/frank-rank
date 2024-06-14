migrations:
	goose sqlite3 frankrank.db -dir ./migrations up
sqlc:
	sqlc generate
run:
	go run ./cmd/*.go
templates:
	/Users/rosenpetkov/go/bin/templ generate
build-fe:
	cd fe && npm run build


# run templ generation in watch mode to detect all .templ files and 
# re-create _templ.txt files on change, then send reload event to browser. 
# Default url: http://localhost:7331
live/templ:
	/Users/rosenpetkov/go/bin/templ generate --watch --proxy="http://localhost:8080" --open-browser=true -v

# run air to detect any go file changes to re-build and re-run the server.
live/server:
	go run github.com/cosmtrek/air@v1.51.0 \
	--build.cmd "go build -o tmp/bin/main ./cmd" --build.bin "tmp/bin/main" --build.delay "100" \
	--build.exclude_dir "node_modules" \
	--build.include_ext "go" \
	--build.stop_on_error "false" \
	--misc.clean_on_exit true 

# run tailwindcss to generate the styles.css bundle in watch mode.
live/vite:
	cd ./fe && npm run watch

# watch for any js or css change in the assets/ folder, then reload the browser via templ proxy.
live/sync_assets:
	go run github.com/cosmtrek/air@v1.51.0 \
	--build.cmd "/Users/rosenpetkov/go/bin/templ generate --notify-proxy" \
	--build.bin "true" \
	--build.delay "300" \
	--build.exclude_dir "" \
	--build.include_dir "fe/dist" \
	--build.include_ext "js,css,html"

# start all 5 watch processes in parallel.
live: 
	make -j4 live/templ live/server live/vite live/sync_assets