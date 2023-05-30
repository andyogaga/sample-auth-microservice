FROM alpine:latest

RUN mkdir /app

# Copy the .env file to the container
COPY .env .env

# Add the .env file as an argument to the build
ARG ENV_FILE=.env

# Use the ARG value to set environment variables during build time
ENV $(cat $ENV_FILE | xargs)

COPY /dist/brokerApp /app

EXPOSE 3001

CMD ["app/brokerApp"]