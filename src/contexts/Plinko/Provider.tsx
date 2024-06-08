import React, { createContext, useState, useMemo } from 'react';

export const PlinkoContext = createContext<any>({});

export const PlinkoProvider: React.FC<React.PropsWithChildren> = ({
  children
}) => {
  const [gameMode, setGameMode] = useState('Easy');
  const [gameRows, setGameRows] = useState(8);
  const [dropSpeed, setDropSpeed] = useState<number[]>([1]);

  return (
    <PlinkoContext.Provider
      value={{
        gameMode: useMemo(() => gameMode, [gameMode]),
        gameRows: useMemo(() => gameRows, [gameRows]),
        dropSpeed: useMemo(() => dropSpeed, [dropSpeed]),
        animSpeed: useMemo(() => dropSpeed[0], [dropSpeed]),
        setGameMode,
        setGameRows,
        setDropSpeed
      }}
    >
      {children}
    </PlinkoContext.Provider>
  );
};
