import React, {
  useState,
  useMemo,
  useCallback,
  useEffect,
  useRef
} from 'react';
import { Fireworks } from '@fireworks-js/react';

import explosion0 from 'assets/audio/explosion0.mp3';
import explosion1 from 'assets/audio/explosion1.mp3';
import explosion2 from 'assets/audio/explosion2.mp3';

import { Box, Text, Button, FairnessIcon, Flex, PageSpinner } from 'components';
import {
  StyledTabs,
  StyledTabPanel,
  StyledFlex,
  StyledTowerBox,
  TabContainer
} from './styles';
import { Tab, TabList } from 'react-tabs';
import ManualTab from './Tabs/Manual';
import AutomaticTab from './Tabs/Automatic';
import HistoryTab from './Tabs/History';
import ManualTower from './components/ManualTower';
import { useAppDispatch, useAppSelector } from 'state';
import AutoTower from './components/AutoTower';
import { reset, setAutoStatus, setDifficulty } from 'state/dreamtower/actions';
import { setBalanceType } from 'state/user/actions';
import { fetchDreamtowerRound } from 'utils/fetchData';
import HistoryTower from './components/HistoryTower';
import Logo from './components/Logo';
import { useMatchBreakpoints } from 'hooks';
import { convertChipToBalance } from 'utils/balance';
// import { Firework } from "./animation/firework";

export default function DreamTower() {
  const dispatch = useAppDispatch();
  const { game, auto } = useAppSelector(state => state.dreamtower);
  const meta = useAppSelector(state => state.meta.dreamtower);
  const {
    isHoliday,
    sound: soundEnabled,
    balances
  } = useAppSelector(state => state.user);

  const { isMobile } = useMatchBreakpoints();

  const [fetchSuccess, setFetchSuccess] = useState(false);
  const [tabIndex, setTabIndex] = useState(0);
  const [lastIndex, setLastIndex] = useState(0);
  const [height, setHeight] = useState('calc(100vh - 65px)');
  const [scale, setScale] = useState(1);

  const dreamRef = useRef<any>(null);
  // const santaRef = useRef<any>(null);

  const disabled = useMemo(
    () => game.status === 'playing' || auto.status === 'running',
    [auto.status, game.status]
  );

  useEffect(() => {
    setAutoStatus('');

    const fetchRoundData = async () => {
      setFetchSuccess(false);
      await fetchDreamtowerRound();
      setFetchSuccess(true);
    };

    fetchRoundData();
  }, []);

  useEffect(() => {
    if (balances['coupon'].balance >= convertChipToBalance(0.01))
      dispatch(setBalanceType('coupon'));
    else dispatch(setBalanceType('chip'));
  }, [balances, dispatch]);

  const resize = useCallback(() => {
    const height = window.innerHeight;
    setHeight(isMobile ? 'auto' : `${height - 65}px`);
    setScale(isMobile ? 1 : (height - 65 - 30) / 840);
  }, [isMobile]);

  useEffect(() => {
    resize();
    window.addEventListener('resize', resize);

    return () => {
      window.removeEventListener('resize', resize);
    };
  }, [resize]);

  const onSelect = useCallback(
    (index: number, lastIndex: number) => {
      if (index === lastIndex) return;
      if (index === 2) {
      }
      setTabIndex(index);
      setLastIndex(lastIndex);
      dispatch(reset());
      dispatch(setDifficulty(meta.difficulties[0]));
    },
    [dispatch, meta.difficulties]
  );

  const handleGoBack = useCallback(() => {
    setTabIndex(lastIndex);
    setLastIndex(tabIndex);
  }, [tabIndex, lastIndex]);

  if (!fetchSuccess) return <PageSpinner />;

  return (
    <>
      <StyledFlex
        padding={['30px 12px', '30px 12px', '30px 12px', '30px 25px']}
        position="relative"
      >
        {/* <div style={{ position: 'absolute' }} ref={santaRef}>
          <SantaHorse />
        </div> */}
        {game.status === 'win' && (
          <Fireworks
            style={{ position: 'absolute', width: '100%', height: '100%' }}
            options={{
              rocketsPoint: {
                min: 0,
                max: 100
              },
              hue: {
                min: 0,
                max: 360
              },
              delay: {
                min: 35,
                max: 90
              },
              decay: {
                min: 0.02,
                max: 0.03
              },
              brightness: {
                min: 50,
                max: 80
              },
              lineStyle: 'round',
              acceleration: 1.05,
              friction: 0.97,
              gravity: 1.5,
              particles: 50,
              trace: 3,
              flickering: 55,
              opacity: 0.2,
              explosion: 6,
              intensity: 33,
              traceSpeed: 4,
              sound: {
                enabled: soundEnabled,
                volume: {
                  min: 60,
                  max: 90
                },
                files: [explosion0, explosion1, explosion2]
              }
            }}
          />
        )}
        <TabContainer>
          <Flex flexWrap="wrap" gap={30}>
            <Text fontSize={23} fontWeight={600} color="#768BAD">
              DREAM TOWER
            </Text>

            <Flex gap={30}>
              <Button
                borderLeft="2px solid #4F617B"
                borderRadius="0px"
                background="#070C12"
                fontSize={14}
                fontWeight={500}
                color="#4F617B"
                p="8px 10px"
                nonClickable={true}
              >
                <FairnessIcon size={16} />
                Fair Game
              </Button>
            </Flex>
          </Flex>
          <Box mt={'40px'}>
            <StyledTabs onSelect={onSelect} selectedIndex={tabIndex}>
              <TabList>
                <Flex>
                  <Tab disabled={disabled}>
                    MANUAL <b />
                  </Tab>
                  <Tab disabled={disabled}>
                    AUTOMATIC <b />
                  </Tab>
                  <Tab disabled={disabled}>
                    GAME HISTORY <b />
                  </Tab>
                </Flex>
              </TabList>
              <StyledTabPanel>
                <ManualTab />
              </StyledTabPanel>
              <StyledTabPanel>
                <AutomaticTab />
              </StyledTabPanel>
              <StyledTabPanel>
                <HistoryTab />
              </StyledTabPanel>
            </StyledTabs>
          </Box>
        </TabContainer>
        <Box mx="auto" height={height}>
          <StyledTowerBox ref={dreamRef} isHoliday={isHoliday} scale={scale}>
            <Logo />
            {tabIndex === 0 && <ManualTower dreamRef={dreamRef} />}
            {tabIndex === 1 && <AutoTower dreamRef={dreamRef} />}
            {tabIndex === 2 && (
              <HistoryTower goBackHandler={handleGoBack} dreamRef={dreamRef} />
            )}
          </StyledTowerBox>
        </Box>
      </StyledFlex>
    </>
  );
}
