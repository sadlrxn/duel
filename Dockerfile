FROM node:18.8-alpine as client
ARG MASTER_WALLET_PUBLIC_KEY
ARG GENERATE_SOURCEMAP
ARG NETWORK
ARG SOLANA_ENDPOINT
ARG STAGE
ARG HAPPY_HOLIDAY
WORKDIR /app
COPY package.json ./
COPY yarn.lock ./

RUN yarn install --frozen-lockfile
COPY . /app
RUN echo "REACT_APP_MASTER_WALLET_PUBLIC_KEY=$MASTER_WALLET_PUBLIC_KEY" > .env
RUN echo "GENERATE_SOURCEMAP=$GENERATE_SOURCEMAP" >> .env
RUN echo "REACT_APP_NETWORK=$NETWORK" >> .env
RUN echo "REACT_APP_SOLANA_ENDPOINT=$SOLANA_ENDPOINT" >> .env
RUN echo "REACT_APP_STAGE=$STAGE" >> .env
RUN echo "REACT_APP_HAPPY_HOLIDAY=$HAPPY_HOLIDAY" >> .env
RUN if [[ -z "$SOLANA_ENDPOINT" ]] ; then echo Argument not provided ; else echo "REACT_APP_SOLANA_ENDPOINT=$SOLANA_ENDPOINT" >> .env ; fi
RUN yarn build



FROM golang:1.21.5-alpine as build

COPY --from=client /app/build/ /app/build/
COPY --from=client /app/server/ /app/server/
COPY --from=client /app/.env /app/

WORKDIR /app/server/

RUN go mod download
RUN go build -o duel -ldflags "-s -w"
EXPOSE 8080

CMD ["/app/server/duel"]