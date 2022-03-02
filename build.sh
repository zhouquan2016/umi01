go_build() {
    go build ./go03/go03.go
}

docker_build() {
    docker build -t zhqn.com:5000/umi-backend:1.0.0 ./go03
}
docker_push() {
    docker push zhqn.com:5000/umi-backend:1.0.0
}
go_build && docker_build && docker_push