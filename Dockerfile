FROM golang:1.23.4

LABEL Name="forum" \
    Version="Latest" \
    Maintainer="jessekuya@proton.me" \
    Contributors="Joseph Owino <https://learn.zone01kisumu.ke/git/joseowino>, Joel Amos <https://learn.zone01kisumu.ke/git/jamos>, John Eliud <https://learn.zone01kisumu.ke/git/johnOdhiambo>, Khalid Hussein <https://learn.zone01kisumu.ke/git/khussein>"

WORKDIR /app

COPY . .
RUN go mod tidy

RUN go build -o /forum

EXPOSE 9000

CMD ["/forum"]
