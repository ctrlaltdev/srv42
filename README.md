Catch All Test HTTP Server

You can listen to any and all incoming requests and print them to stdout, or you can serve local static files if you provide a path.

## Usage

```sh
srv42
```

To use another port than 1337:
```sh
srv42 -p 5000
```

To response another status than 202 (ignored if a path is provided):
```sh
srv42 -s 204
```

To have the full details of the request:
```sh
srv42 -v
```

To serve static files:
```sh
srv42 www/static/
```
