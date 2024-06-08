import { useContext } from 'react';
import { PlinkoContext } from './Provider';

const usePlinko = () => {
  const plinko = useContext(PlinkoContext);

  return plinko;
};

export default usePlinko;
