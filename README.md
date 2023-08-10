# nubl9 recrutation task

## Usage

```go
docker build -t random .
docker run -p 8080:8080 random
```

In the browser type:

```html
http://127.0.0.1:8080/random/mean?requests={r}&length={l}
```

Where:

- r is the number of concurrent requests
- l is the length of requested set
