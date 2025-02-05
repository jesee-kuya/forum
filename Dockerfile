FROM golang:1.23.5

LABEL Name="forum" \
      Maintainer="kuyajesee@proton.me" \
      Version="latest" \
      Contributors="John Odhiambo <https://learn.zone01kisumu.ke/git/johnOdhiambo>, Joseph Owino <https://learn.zone01kisumu.ke/git/jowino>, Khalid Hussein <https://learn.zone01kisumu.ke/git/khussein>, Joel Amos <https://learn.zone01kisumu.ke/git/jamos>"

WORKDIR /app

COPY . .

RUN go mod tidy

EXPOSE 9000

CMD [ "go", "run", "main.go", "9000" ]

