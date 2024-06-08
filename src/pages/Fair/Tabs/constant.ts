export const coinflipAlgorithm = `candidates.sort((a: Candidate, b: Candidate) => {
  return compareByPercentAndId(a, b);
});
const msgBuffer = new TextEncoder().encode(randomString);

const hashBuffer = await crypto.subtle.digest("SHA-256", msgBuffer);

const hashArray = Array.from(new Uint8Array(hashBuffer));

const hashHex = hashArray
  .map((b) => ("00" + b.toString(16)).slice(-2))
  .join("");

var bn = new BigNumber("0x" + hashHex.slice(-16));
var bd = new BigNumber(2 ** 64 - 1);
var outCome = bn.dividedBy(bd).toNumber();

var winnerIndex = 0;
while (outCome > candidates[winnerIndex].percent) {
  outCome -= candidates[winnerIndex].percent;
  winnerIndex++;
}

console.log();
`;

export const jackpotAlgorithm = `
import BigNumber from "bignumber.js";

interface Candidate {
  id: number;
  name: string;
  usdAmount: number;
  nftAmount: number;
}

function compare(a: Candidate, b: Candidate) {
  return a.usdAmount + a.nftAmount > b.usdAmount + b.nftAmount ? -1 : 1;
}

async function verifyFairness(candidates: Candidate[], randomString: string) {
  // Sort candidates by their bet amounts.
  candidates.sort((a: Candidate, b: Candidate) => {
    return compare(a, b);
  });

  // Calculate 'outCome' from 'randomString'
  const msgBuffer = new TextEncoder().encode(randomString);
  const hashBuffer = await crypto.subtle.digest("SHA-256", msgBuffer);
  const hashArray = Array.from(new Uint8Array(hashBuffer));
  const hashHex = hashArray
    .map((b) => ("00" + b.toString(16)).slice(-2))
    .join("");
  var bn = new BigNumber("0x" + hashHex.slice(0, 16));
  var bd = new BigNumber(2 ** 64 - 1);
  var outCome = bn.dividedBy(bd).toNumber();
  console.log("Outcome = ", outCome);

  // Calculate total amount of candidates betted
  var totalBetAmount = 0;
  candidates.forEach((candidate) => {
    totalBetAmount += candidate.usdAmount + candidate.nftAmount;
  });

  // Determine winner with 'outCome'
  var winnerIndex = 0;
  while (
    outCome >
    (candidates[winnerIndex].usdAmount + candidates[winnerIndex].nftAmount) /
      totalBetAmount
  ) {
    outCome -=
      (candidates[winnerIndex].usdAmount + candidates[winnerIndex].nftAmount) /
      totalBetAmount;
    winnerIndex++;
  }
  console.log("Winner : ", candidates[winnerIndex]);
}
`;

export const dreamtowerAlgorithm = `async function byteGenerator(
  serverSeed: string,
  clientSeed: string,
  nonce: number,
  cursor: number
) {
  const currentRound = Math.floor(cursor / 32);
  const currentRoundCursor = cursor - currentRound * 32;
  const str = serverSeed + ':' + clientSeed + ':' + nonce.toString() + ':' + currentRound.toString();
  const msgBuffer = new TextEncoder().encode(str);
  const hashBuffer = await crypto.subtle.digest('SHA-256', msgBuffer);
  const hashArray = Array.from(new Uint8Array(hashBuffer));
  return hashArray.slice(currentRoundCursor, currentRoundCursor + 4);
}

function generateRow(blocksInRow: number, starsInRow: number, shuffle: number[]) {
  var arr = [];
  for (var i = 0; i < blocksInRow; i++) {
    arr.push(i);
  }
  for (i = blocksInRow - 1; i > 0; i--) {
    const temp: number = arr[i];
    arr[i] = arr[shuffle[i]];
    arr[shuffle[i]] = temp;
  }
  return arr.slice(0, starsInRow);
}

async function generateTower(
  serverSeed: string,
  clientSeed: string,
  nonce: number,
  blocksInRow: number,
  starsInRow: number,
  count: number
) {
  var tower = [];
  var cursor = 0;
  for (var i = 0; i < count; i++) {
    var shuffle = [];
    for (var j = 0; j < blocksInRow; j++) {
      const bytes = await byteGenerator(serverSeed, clientSeed, nonce, cursor);
      const hashHex = bytes.map(b => ('00' + b.toString(16)).slice(-2)).join('');
      const rn = new BigNumber('0x' + hashHex);
      const bd = new BigNumber(2 ** 32);
      const result = rn.dividedBy(bd).toNumber() * (j + 1);
      shuffle.push(Math.floor(result));
      cursor += 4;
    }
    tower.push(generateRow(blocksInRow, starsInRow, shuffle));
  }
  console.log(tower);
}
`;

export const crashAlgorithm = `async function calculateRandomNumber(serverSeed: string, clientSeed: string) {
  const msgBuffer = new TextEncoder().encode(serverSeed + clientSeed);
  const hashBuffer = await crypto.subtle.digest('SHA-256', msgBuffer);
  const hashArray = Array.from(new Uint8Array(hashBuffer));

  const hashHex = hashArray
    .map(b => ('00' + b.toString(16)).slice(-2))
    .join('');

  return new BigNumber('0x' + hashHex.slice(0, 13));
}

export async function verifyCrashFairness(
  serverSeed: string,
  clientSeed: string
) {
  var randomNumber = await calculateRandomNumber(serverSeed, clientSeed);
  var hs = 100 / 5;
  if (randomNumber.mod(hs) === BigNumber(0)) {
    return 1;
  }

  var h = randomNumber;
  var e = new BigNumber(2 ** 52);
  return (
    Math.floor(e.minus(h).multipliedBy(100).dividedBy(e.minus(h)).toNumber()) /
    100.0
  );
}

console.log(verifyCrashFairness(serverSeed, clientSeed));
`;

export const coinflipSampleInput = `{
  "randomString": "abcdefghijklmnopq",
  "betAmount": 4827,
  "players": [
    {
      "id": 1,
      "name": "Anonymous1",
      "color": "Green"
    },
    {
      "id": 2,
      "name": "Anonymous2",
      "color": "Purple"
    }
  ]
}
`;

export const jackpotSampleInput = `{
  "randomString" : "fi(*i27kd&*ls",
  "players" : [
    {
      "id" : 1,
      "name" : "Anonymous1",
      "usdAmount" : 11395,
      "nftAmount" : 0
    },
    {
      "id" : 3,
      "name" : "Anonymous2",
      "usdAmount" : 442,
      "nftAmount" : 4352
    },
    {
      "id" : 2,
      "name" : "Anonymous3",
      "usdAmount" : 1212,
      "nftAmount" : 322
    }
  ]
}
`;

export const dreamtowerSampleInput = `{
  "clientSeed" : "4kOhs37*(*84",
  "serverSeed" : "sio298)ksjs(*724",
  "nonce" : 38,
  "difficulty": "easy"
}
`;

export const crashSampleInput = `{
  "serverSeed" : "5c240332b66182ef00a77e00e66d0080b2375c17d671beb2d34d43ba73339278",
  "clientSeed" : "Bryw516EWrm7tnDJ3etvnYeSGNkeFYaVG7gKVoTQytCs"
}`;
