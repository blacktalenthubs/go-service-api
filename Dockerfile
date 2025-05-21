FROM golang:1.20

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Create a test script
RUN echo 'package main; import ("database/sql"; "fmt"; _ "github.com/lib/pq"); func main() { db, err := sql.Open("postgres", "host=consultancy-postgres port=5432 user=postgres password=postgres dbname=consultancy sslmode=disable"); if err != nil { fmt.Printf("Open error: %v\n", err); return; }; err = db.Ping(); if err != nil { fmt.Printf("Ping error: %v\n", err); return; }; fmt.Println("Connected successfully!"); }' > /app/test.go

CMD ["go", "run", "test.go"]