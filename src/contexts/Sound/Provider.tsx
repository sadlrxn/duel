import React, { createContext } from 'react';
import useSound from 'use-sound';

import clickSound from 'assets/audio/button-click.wav';

import coinStartSound from 'assets/audio/coin-start.wav';
import coinEndSound from 'assets/audio/coin-end.wav';

import count1Sound from 'assets/audio/count-1.wav';
import count2Sound from 'assets/audio/count-2.wav';
import count3Sound from 'assets/audio/count-3.wav';

import tickSound from 'assets/audio/jackpot-tick.wav';
import rollendSound from 'assets/audio/jackpot-rollend.wav';
import winSound from 'assets/audio/jackpot-win.wav';

import successSound from 'assets/audio/success.wav';
import errorSound from 'assets/audio/error.wav';

import towerBreakSound from 'assets/audio/tower-break.wav';
import towerFireworksSound from 'assets/audio/tower-fireworks.wav';
import towerLastStarSound from 'assets/audio/tower-last-star.wav';
import towerSelectStarSound from 'assets/audio/tower-select-star.wav';
import towerWinSound from 'assets/audio/tower-win.wav';

import crashBoostSound from 'assets/audio/crash-rocket-boost.wav';
import crashLaunchSound from 'assets/audio/crash-launch.wav';
import crashMultipleSound from 'assets/audio/crash-100x.wav';
import crashCount4Sound from 'assets/audio/crash-count-4.mp3';
import crashCount3Sound from 'assets/audio/crash-count-3.mp3';
import crashCount2Sound from 'assets/audio/crash-count-2.mp3';
import crashCount1Sound from 'assets/audio/crash-count-1.mp3';
import crashCountEndSound from 'assets/audio/crash-count-end.mp3';
import crashCashoutSound from 'assets/audio/crash-cashout.wav';

const defaultFunction = (_param?: any) => {};

const buttonPlay: any = defaultFunction;
const coinPlay = {
  start: defaultFunction,
  end: defaultFunction
};
const countPlay = {
  c1: defaultFunction,
  c2: defaultFunction,
  c3: defaultFunction
};
const crashPlay = {
  boost: defaultFunction,
  launch: defaultFunction,
  c4: defaultFunction,
  c3: defaultFunction,
  c2: defaultFunction,
  c1: defaultFunction,
  cend: defaultFunction,
  cashout: defaultFunction,
  multi: defaultFunction
};
const jackpotPlay = {
  tick: defaultFunction,
  rollend: defaultFunction,
  win: defaultFunction
};
const towerPlay = {
  break: defaultFunction,
  fireworks: defaultFunction,
  lastStar: defaultFunction,
  selectStar: defaultFunction,
  win: defaultFunction
};
const messagePlay = {
  success: defaultFunction,
  error: defaultFunction
};

const buttonStop: any = defaultFunction;
const coinStop = {
  start: defaultFunction,
  end: defaultFunction
};
const countStop = {
  c1: defaultFunction,
  c2: defaultFunction,
  c3: defaultFunction
};
const crashStop = {
  boost: defaultFunction,
  launch: defaultFunction,
  c4: defaultFunction,
  c3: defaultFunction,
  c2: defaultFunction,
  c1: defaultFunction,
  cend: defaultFunction,
  cashout: defaultFunction,
  multi: defaultFunction
};
const jackpotStop = {
  tick: defaultFunction,
  rollend: defaultFunction,
  win: defaultFunction
};
const towerStop = {
  break: defaultFunction,
  fireworks: defaultFunction,
  lastStar: defaultFunction,
  selectStar: defaultFunction,
  win: defaultFunction
};
const messageStop = {
  success: defaultFunction,
  error: defaultFunction
};

export const SoundContext = createContext({
  buttonPlay,
  coinPlay,
  countPlay,
  crashPlay,
  jackpotPlay,
  towerPlay,
  messagePlay,
  buttonStop,
  coinStop,
  countStop,
  crashStop,
  jackpotStop,
  towerStop,
  messageStop
});

export const SoundProvider: React.FC<React.PropsWithChildren> = ({
  children
}) => {
  const [buttonPlay, { stop: buttonStop }] = useSound(clickSound, {
    interrupt: true
  });

  const [coinStartPlay, { stop: coinStartStop }] = useSound(coinStartSound, {
    interrupt: true
  });

  const [coinEndPlay, { stop: coinEndStop }] = useSound(coinEndSound, {
    interrupt: true
  });

  const [count1Play, { stop: count1Stop }] = useSound(count1Sound, {
    interrupt: true
  });

  const [count2Play, { stop: count2Stop }] = useSound(count2Sound, {
    interrupt: true
  });

  const [count3Play, { stop: count3Stop }] = useSound(count3Sound, {
    interrupt: true
  });

  const [tickPlay, { stop: tickStop }] = useSound(tickSound);

  const [rollendPlay, { stop: rollendStop }] = useSound(rollendSound, {
    volume: 0.8,
    interrupt: true
  });

  const [winPlay, { stop: winStop }] = useSound(winSound, {
    volume: 0.8,
    interrupt: true
  });

  const [successPlay, { stop: successStop }] = useSound(successSound);

  const [errorPlay, { stop: errorStop }] = useSound(errorSound);

  const [towerWinPlay, { stop: towerWinStop }] = useSound(towerWinSound, {
    interrupt: true,
    loop: true
  });

  const [towerBreakPlay, { stop: towerBreakStop }] = useSound(towerBreakSound, {
    interrupt: true,
    volume: 0.7
  });

  const [towerFireworksPlay, { stop: towerFireworksStop }] =
    useSound(towerFireworksSound);

  const [towerLastStarPlay, { stop: towerLastStarStop }] = useSound(
    towerLastStarSound,
    {
      interrupt: true
    }
  );

  const [towerSelectStarPlay, { stop: towerSelectStarStop }] = useSound(
    towerSelectStarSound,
    {
      interrupt: true
    }
  );

  const [crashBoostPlay, { stop: crashBoostStop }] = useSound(crashBoostSound, {
    interrupt: true,
    loop: true
  });

  const [crashLaunchPlay, { stop: crashLaunchStop }] = useSound(
    crashLaunchSound,
    {
      interrupt: true
    }
  );

  const [crashMultiplePlay, { stop: crashMultipleStop }] = useSound(
    crashMultipleSound,
    {
      interrupt: true
    }
  );

  const [crashCashoutPlay, { stop: crashCashoutStop }] = useSound(
    crashCashoutSound,
    {
      interrupt: true
    }
  );

  const [crashCountEndPlay, { stop: crashCountEndStop }] = useSound(
    crashCountEndSound,
    {
      interrupt: true
    }
  );

  const [crashCount4Play, { stop: crashCount4Stop }] = useSound(
    crashCount4Sound,
    {
      interrupt: true
    }
  );

  const [crashCount3Play, { stop: crashCount3Stop }] = useSound(
    crashCount3Sound,
    {
      interrupt: true
    }
  );

  const [crashCount2Play, { stop: crashCount2Stop }] = useSound(
    crashCount2Sound,
    {
      interrupt: true
    }
  );

  const [crashCount1Play, { stop: crashCount1Stop }] = useSound(
    crashCount1Sound,
    {
      interrupt: true
    }
  );

  return (
    <SoundContext.Provider
      value={{
        buttonPlay,
        coinPlay: {
          start: coinStartPlay,
          end: coinEndPlay
        },
        countPlay: {
          c1: count1Play,
          c2: count2Play,
          c3: count3Play
        },
        crashPlay: {
          boost: crashBoostPlay,
          launch: crashLaunchPlay,
          c4: crashCount4Play,
          c3: crashCount3Play,
          c2: crashCount2Play,
          c1: crashCount1Play,
          cend: crashCountEndPlay,
          cashout: crashCashoutPlay,
          multi: crashMultiplePlay
        },
        jackpotPlay: {
          tick: tickPlay,
          rollend: rollendPlay,
          win: winPlay
        },
        towerPlay: {
          break: towerBreakPlay,
          fireworks: towerFireworksPlay,
          lastStar: towerLastStarPlay,
          selectStar: towerSelectStarPlay,
          win: towerWinPlay
        },
        messagePlay: {
          success: successPlay,
          error: errorPlay
        },
        buttonStop,
        coinStop: {
          start: coinStartStop,
          end: coinEndStop
        },
        countStop: {
          c1: count1Stop,
          c2: count2Stop,
          c3: count3Stop
        },
        crashStop: {
          boost: crashBoostStop,
          launch: crashLaunchStop,
          c4: crashCount4Stop,
          c3: crashCount3Stop,
          c2: crashCount2Stop,
          c1: crashCount1Stop,
          cend: crashCountEndStop,
          cashout: crashCashoutStop,
          multi: crashMultipleStop
        },
        jackpotStop: {
          tick: tickStop,
          rollend: rollendStop,
          win: winStop
        },
        towerStop: {
          break: towerBreakStop,
          fireworks: towerFireworksStop,
          lastStar: towerLastStarStop,
          selectStar: towerSelectStarStop,
          win: towerWinStop
        },
        messageStop: {
          success: successStop,
          error: errorStop
        }
      }}
    >
      {children}
    </SoundContext.Provider>
  );
};
