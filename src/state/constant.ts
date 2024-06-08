let ws = "ws";
let HOST = "application:8080";
if (typeof window !== "undefined") {
  if (window.location.hostname !== "localhost") {
    HOST = window.location.host;
   // ws = "wss";
  }
}

// console.log(window.location.hostname)

export const LOW_ROOM_GAME_URL = `${ws}://${HOST}/api/jackpot/low/ws`;
export const HIGH_ROOM_GAME_URL = `${ws}://${HOST}/api/jackpot/high/ws`;
export const COINFLIP_GAME_URL = `${ws}://${HOST}/api/coinflip/ws`;
export const PAYMENT_URL = `${ws}://${HOST}/api/pay/ws`;

export const SERVER_URL = `${ws}://${HOST}/api/ws`;

export const SOCKET_RECONNECTION_TIMEOUT = 10000;
