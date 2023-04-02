FROM alpine:latest

RUN mkdir /app

COPY /dist/brokerApp /app

EXPOSE 3001

CMD ["app/brokerApp"]