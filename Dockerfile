FROM jess/inkscape

RUN apt-get update && apt-get install -y python-lxml python-numpy wget pstoedit unzip

RUN wget https://golang.org/dl/go1.16.2.linux-amd64.tar.gz && rm -rf /usr/local/go && tar -C /usr/local -xzf go1.16.2.linux-amd64.tar.gz && ln -s /usr/local/go/bin/go /usr/local/bin/go && go version

RUN wget "https://fonts.google.com/download?family=Fira%20Sans|Lobster|Pacifico|Permanent%20Marker|Roboto|Ceviche%20One|Press%20Start%202P|UnifrakturMaguntia|Zilla%20Slab%20Highlight|Londrina%20Shadow" -O /usr/share/fonts/fonts.zip
RUN cd /usr/share/fonts && unzip fonts.zip

WORKDIR /app
ADD go.mod go.sum main.go index.html template.svg ./
RUN go build -o main main.go

ENTRYPOINT ["/app/main"]
