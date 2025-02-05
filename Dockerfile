FROM golang:1.23.4

LABEL Name="forum" \
    Version="Latest" \
    Maintainer="jessekuya@proton.me" \
    Contributors="Joseph Owino <https://learn.zone01kisumu.ke/git/joseowino>, Joel Amos <https://learn.zone01kisumu.ke/git/jamos>, John Eliud <https://learn.zone01kisumu.ke/git/johnOdhiambo>, Khalid Hussein <https://learn.zone01kisumu.ke/git/khussein>"

WORKDIR /app

COPY . .
RUN go mod tidy

RUN CGO_ENABLED=0 GOOG=LINUX go build -o /forum

EXPOSE 9000

CMD ["/forum"]


# Docker permission setup (optional)
# Commands to install Docker in rootless mode
# curl -fsSL https://get.docker.com/rootless | sh
# export PATH=/home/docker/bin:$PATH
# export DOCKER_HOST=unix:///run/user/10531/docker.sock
