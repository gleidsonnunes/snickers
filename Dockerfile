FROM flavioribeiro/snickers-docker:v3

# Download snickers
RUN go get -u github.com/gleidsonnunes/snickers

# Run snickers!
RUN curl -O http://flv.io/gleidsonnunes/config.json
RUN go install github.com/gleidsonnunes/snickers
ENTRYPOINT snickers
EXPOSE 8000
