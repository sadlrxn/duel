import { multipliers, gradients } from '../config';

export const getMultipliers = (gameMode: string, gameRows: number): any[] => {
  //@ts-ignore
  return multipliers[gameMode.toLowerCase() + '_' + gameRows].map(
    (multiplier: string, index: number) => {
      //@ts-ignore
      return { multiplier, gradient: gradients[gameRows][index] };
    }
  );
};
