FROM flavioribeiro/snickers-docker:v3

# Download snickers
RUN go get -u github.com/gleidsonnunes/snickers2/snickers

# Run snickers!
RUN curl -O http://flv.io/gleidsonnunes/snickers2/config.json
RUN go install github.com/gleidsonnunes/snickers2/snickers
ENTRYPOINT snickers
EXPOSE 8000
