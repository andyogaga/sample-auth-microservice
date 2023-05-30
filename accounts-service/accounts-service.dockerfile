
FROM alpine:latest

RUN mkdir /app

COPY /dist/accountsApp /app

EXPOSE 3001

CMD ["app/accountsApp"]