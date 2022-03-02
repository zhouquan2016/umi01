go_build() {
    go build
}

docker_build() {
    docker build -t zhqn.com:5000/umi-backend:1.0.0 .
}
docker_push() {
    docker push zhqn.com:5000/umi-backend:1.0.0
}
cd go03
go_build && docker_build && docker_push
cd ../