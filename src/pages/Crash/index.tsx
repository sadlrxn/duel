import { useState, useCallback, useEffect, useRef, useMemo } from 'react';
import styled, { css } from 'styled-components';
import { Tab, TabList, Tabs } from 'react-tabs';

import { Topbar, Flex, Grid, Box, PageSpinner } from 'components';
import { useAppSelector } from 'state';

import { StyledTabPanel, StyledTabs } from 'components/GameTabs';
import { CrashProvider } from 'contexts/Crash/Provider';
import { useSound, useCrash } from 'hooks';

import { Background, Stars, Sky, CrashGame, History } from './components';
import { Manual, AutoBet } from './Tabs';

const getProperty = (elt: Element, property: string) => {
  return window.getComputedStyle(elt).getPropertyValue(property);
};

function Crash() {
  const isLoaded = useAppSelector(state => state.crash.fetch);
  const status = useAppSelector(state => state.crash.status);
  const globalTime = useAppSelector(state => state.crash.time);
  const { autoBetEnable } = useCrash();

  const { crashPlay, crashStop } = useSound();

  const [tabIndex, setTabIndex] = useState(0);
  const [_, setLastIndex] = useState(0);

  const [tabMaxHeight, setTabMaxHeight] = useState('none');
  const [gameHeight, setGameHeight] = useState(0);
  const [gameWidth, setGameWidth] = useState(0);
  const [bottomGap, setBottomGap] = useState(0);

  const containerRef = useRef<HTMLDivElement>(null);
  const gameRef = useRef<HTMLDivElement>(null);

  const resizeObserver = useMemo(
    () =>
      new ResizeObserver(() => {
        if (!containerRef || !containerRef.current) return;
        if (!gameRef || !gameRef.current) return;

        const parent = containerRef.current.parentElement;
        if (!parent) return;

        // 65 - Navbar, 34 - Topbar, 40 - History, 36 - gap * 2,
        // 94 - 74 + gap(20)

        const width = +getProperty(parent, 'width').slice(0, -2);
        setGameWidth(width);

        const height = +getProperty(parent, 'height').slice(0, -2);
        const bottomGap = width < 700 ? 74 : height > 1000 ? 250 : 150;
        setBottomGap(bottomGap);

        const gameHeight = gameRef.current.offsetTop + 65 + 20;
        setGameHeight(gameHeight);

        if (width > 700) setTabMaxHeight(`calc(100vh - ${gameHeight}px)`);
        else setTabMaxHeight('none');
      }),
    []
  );

  useEffect(() => {
    if (!containerRef || !containerRef.current) return;
    resizeObserver.observe(containerRef.current);
  }, [resizeObserver]);

  const boostSoundRef = useRef(false);

  useEffect(() => {
    let timeElapsed = 0;
    switch (status) {
      case 'bet':
        crashStop.boost();
        boostSoundRef.current = false;
        break;
      case 'ready':
        timeElapsed = Date.now() - globalTime;
        if (timeElapsed < 300) crashPlay.launch();
        break;
      case 'play':
        if (boostSoundRef.current) break;
        // crashPlay.boost();
        boostSoundRef.current = true;
        break;
      case 'explosion':
        boostSoundRef.current = false;
        crashStop.boost();
        break;
      case 'back':
        boostSoundRef.current = false;
        crashStop.boost();
        break;
    }
  }, [crashPlay, crashStop, globalTime, status]);

  useEffect(() => {
    return () => {
      crashStop.boost();
      crashStop.launch();
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const onSelect = useCallback(
    (index: number, lastIndex: number) => {
      if (index === lastIndex) return;
      if (autoBetEnable && tabIndex === 1) return;
      // if (index === 1) return;
      setTabIndex(index);
      setLastIndex(lastIndex);
    },
    [autoBetEnable, tabIndex]
  );

  return (
    <Container ref={containerRef}>
      {isLoaded ? (
        <>
          <Background />
          <Stars />
          <Sky />
          <CustomTopbar title="CRASH" />
          <History />
          <GameContainer ref={gameRef}>
            <GameTabWrapper>
              <StyledGameTabs>
                <Tabs onSelect={onSelect} selectedIndex={tabIndex}>
                  <TabList>
                    <Flex>
                      <Tab>
                        MANUAL
                        <b />
                      </Tab>
                      <Tab>
                        AUTOMATIC
                        <b />
                      </Tab>
                    </Flex>
                  </TabList>
                  <StyledTabPanel>
                    <Box className="container" maxHeight={tabMaxHeight}>
                      <Manual
                        gameHeight={gameHeight}
                        bottomGap={bottomGap}
                        gameWidth={gameWidth}
                      />
                    </Box>
                  </StyledTabPanel>
                  <StyledTabPanel>
                    <Box className="container">
                      <AutoBet
                        gameHeight={gameHeight}
                        bottomGap={bottomGap}
                        gameWidth={gameWidth}
                      />
                    </Box>
                  </StyledTabPanel>
                </Tabs>
              </StyledGameTabs>
            </GameTabWrapper>
            <CrashGame
              gameHeight={`calc(100vh - ${gameHeight + bottomGap}px)`}
            />
          </GameContainer>
        </>
      ) : (
        <PageSpinner />
      )}
    </Container>
  );
}

export default function CrashWrapper() {
  return (
    <CrashProvider>
      <Crash />
    </CrashProvider>
  );
}

const CustomTopbar = styled(Topbar)`
  margin-left: 16px;

  .width_700 & {
    margin-left: 0px;
  }
`;

const StyledGameTabs = styled(StyledTabs)`
  font-size: 16px;
  height: 100%;

  .react-tabs {
    height: 100%;
  }

  .width_700 & {
    font-size: 12px;
  }

  .react-tabs__tab-list {
    display: grid;
    grid-template-columns: repeat(auto-fill, 1fr);
    margin: 0;
    border: none;

    div {
      background: #202e3f;
      border-radius: 0px;

      width: 100%;

      .width_700 & {
        border-radius: 10px 10px 0 0;
      }
    }

    .react-tabs__tab {
      bottom: 0;
      padding: 12px 15px;
      width: 100%;

      border-top-right-radius: 10px;

      font-family: 'Inter';
      font-style: normal;
      font-weight: 600;
      font-size: 14px;

      width: 100%;

      ${({ theme }) => theme.mediaQueries.md} {
        font-size: 14px;
      }

      line-height: 10px;
      text-align: center;
      letter-spacing: 0.18em;
      color: #697e9c;
      text-transform: uppercase;

      :focus:after {
        display: none !important;
      }

      b {
        position: absolute;
        left: 50%;
        transform: translateX(-50%);
        background-color: #4fff8b;
        width: 0%;
        height: 2px;
        border-radius: 2px;
        overflow: hidden;
        transition: all 0.3s ease-in;

        top: calc(100% - 2px);

        .width_700 & {
          top: -2px;
        }
      }

      :hover b {
        width: 60%;
      }
    }

    .react-tabs__tab--selected {
      background: #202e3f;
      border-width: 0px;

      border: 0px solid #4f617b;
      border-radius: 0 0 0 0;
      border-bottom: 0;
      color: #4fff8b;

      .width_700 & {
        background: #101b2c;
        /* background: rgba(19, 32, 49, 0.75); */
        border-width: 2px;
        border-radius: 10px 10px 0 0;
      }

      ::before {
        position: absolute;
        background: #101b2c;
        ${({ tabselectedbackground }) =>
          tabselectedbackground &&
          css`
            background: ${tabselectedbackground};
          `}
        content: '';
        display: block;
        width: 100%;
        left: 0;
        bottom: -2px;
        pointer-events: none;
        z-index: 10;

        height: 0px;

        .width_700 & {
          height: 2px;
        }
      }

      b {
        width: 60%;
      }
    }
  }

  .react-tabs__tab-panel--selected {
    .container {
      position: relative;

      background: #131e2d;
      border-width: 0px;
      border-radius: 0px;

      border: 2px solid #4f617b;
      border-radius: 0px 0px 0px 0px;
      border-width: 0px;

      padding: 10px 14px;
      .width_700 & {
        padding: 18px 20px;
      }

      backdrop-filter: blur(5px);

      .width_700 & {
        border-width: 2px;
        border-radius: 0px 0px 10px 10px;
        background: linear-gradient(
          180deg,
          rgba(19, 32, 49, 0.75) 0%,
          rgba(26, 41, 60, 0.75) 65.67%
        );
      }
    }
  }
`;

const GameTabWrapper = styled(Box)`
  z-index: 1;
  order: 3;
  height: 100%;
  .width_700 & {
    order: 0;
  }
`;

const GameContainer = styled(Grid)`
  gap: 20px;
  width: 100%;
  height: 100%;

  .width_700 & {
    grid-template-columns: 344px 1fr;
    max-height: calc(100% - 110px);
  }
`;

const Container = styled(Flex)`
  flex-direction: column;
  gap: 18px;
  height: 100%;

  max-width: 1000px;
  padding-top: 30px;

  .width_700 & {
    margin-left: auto;
    margin-right: auto;
    margin-bottom: 0px;

    padding: 30px 25px;
  }
`;
