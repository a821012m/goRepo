FROM golang:1.19.0-alpine
WORKDIR /line
ADD . /line
COPY appSettings.json /line/appSettings.json 
RUN cd /line && go build
EXPOSE 8080
ENTRYPOINT ./line