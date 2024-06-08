import { useContext } from 'react';
import { CrashContext } from './Provider';

const useCrash = () => {
  const crash = useContext(CrashContext);

  return crash;
};

export default useCrash;
