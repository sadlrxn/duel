import { Candidate } from 'api/types/jackpot';

const randomIndices = [
  7, 61, 81, 54, 24, 77, 42, 0, 11, 8, 19, 12, 14, 49, 94, 62, 43, 84, 35, 33,
  27, 59, 31, 40, 90, 57, 50, 66, 41, 78, 32, 56, 86, 18, 16, 93, 10, 75, 47,
  65, 5, 1, 58, 30, 64, 51, 71, 28, 98, 9, 83, 37, 22, 6, 53, 15, 79, 55, 23,
  46, 21, 52, 99, 34, 97, 2, 48, 17, 89, 69, 63, 72, 80, 87, 82, 95, 20, 44, 26,
  45, 92, 88, 67, 38, 4, 13, 91, 76, 74, 85, 36, 70, 39, 25, 29, 68, 3, 73, 96,
  60
];

const DISPLAY_CARD_COUNT = 20;
const TOTAL_CARDS = 100;

export const generateCandidateData = (
  candidates: Candidate[],
  winnerId: number,
  roundId: number
) => {
  const cardCount = candidates.reduce(
    (sum: number, c) => sum + (c.count ?? 0),
    0
  );
  const candidateData = candidates
    .slice()
    // Sorts candidates to sync card display order among users
    .sort((c1, c2) => c2.id - c1.id)
    .map(c => ({
      ...c,
      count: Math.ceil(((c.count ?? 1) * 100) / cardCount)
    }))
    // Make cards count to over TOTAL_CARDS
    .reduce((accum: Candidate[], candidate) => {
      const array = Array(candidate.count).fill(candidate);
      return [...accum, ...array];
    }, []);

  const winner = candidates.find(candidate => candidate.id === winnerId);
  if (!winner) {
    return {
      candidateData: [],
      cardCount: TOTAL_CARDS,
      rotation: 0
    };
  }
  const winnerIndex = roundId % TOTAL_CARDS;
  candidateData[winnerIndex] = winner;

  const newCandidateData = candidateData
    .slice(0, TOTAL_CARDS)
    .map((c, i) => ({ ...c, randomIndex: randomIndices[i] }))
    .sort((a, b) => a.randomIndex - b.randomIndex);

  const index = randomIndices[winnerIndex];

  const degree = 360 / DISPLAY_CARD_COUNT;
  const rotation =
    degree * index +
    ((randomIndices[winnerIndex] / 100) * 0.88 - 0.44) * degree;

  return {
    candidateData: newCandidateData,
    cardCount: TOTAL_CARDS,
    rotation
  };
};
