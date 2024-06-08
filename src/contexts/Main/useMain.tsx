import { useContext } from 'react';
import { MainContext } from './Provider';

const useMain = () => {
  const main = useContext(MainContext);

  return main;
};

export default useMain;
