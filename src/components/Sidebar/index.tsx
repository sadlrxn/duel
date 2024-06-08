import React, { useCallback, useState, useEffect, useMemo } from 'react';
import { Menu, MenuItem, SubMenu, useProSidebar } from 'react-pro-sidebar';
import { Link, NavLink, useLocation } from 'react-router-dom';
import useHover from '@react-hook/hover';
import { useAppSelector } from 'state';
import { useMatchBreakpoints } from 'hooks';
import { getJackpotProgress } from 'pages/Jackpot/utils';

import { StyledSidebar } from './styles';

import {
  CoinflipIcon,
  JackpotIcon,
  DreamtowerIcon,
  PlinkoIcon,
  GrandJackpotIcon,
  CrashIcon,
  WagerRaceIcon,
  DuelbotsIcon,
  FairIcon,
  SupportIcon,
  DocsIcon
} from 'components/Icon/Sidebar';
import ExpandIcon from './ExpandIcon';

const convertNumberToString = (val: number, maxLength: number = 2) => {
  return val.toString().padStart(maxLength, '0');
};

export default function SidebarFn({
  logout,
  ...props
}: {
  login: () => Promise<void>;
  logout: () => Promise<void>;
}) {
  const { pathname } = useLocation();
  const target = React.useRef(null);
  const isHovering = useHover(target, { enterDelay: 200, leaveDelay: 200 });

  const meta = useAppSelector(state => state.meta.grandJackpot);
  const grandJackpot = useAppSelector(state => state.grandJackpot.game);
  const { collapseSidebar } = useProSidebar();
  const { isHoliday, isAuthenticated } = useAppSelector(state => state.user);
  // const { isAuthenticated, balance } = useAppSelector((state) => state.user);

  const { isMobile } = useMatchBreakpoints();

  const [_countText, setCountText] = useState('00:00');

  const subMenuActive = useMemo(() => {
    if (
      pathname.includes('/duelbots/dashboard') ||
      pathname.includes('/duelbots/myduelbots') ||
      pathname.includes('/duelbots/staking')
    )
      return true;
    else return false;
  }, [pathname]);

  useEffect(() => {
    if (isMobile) return;
    if (isHovering) collapseSidebar(false);
    else collapseSidebar(true);
  }, [isHovering]);

  useEffect(() => {
    const interval = setInterval(() => {
      const { count } = getJackpotProgress({
        countingTime: meta.countingTime,
        updatedTime: grandJackpot.time,
        status: grandJackpot.status,
        rollingTime: meta.rollingTime,
        winnerTime: meta.winnerTime
      });

      setCountText(
        `${convertNumberToString(
          Math.floor(count / 3600)
        )}:${convertNumberToString(Math.floor(count / 60) % 60)}}`
      );
    }, 1000);

    return () => {
      clearInterval(interval);
    };
  }, [grandJackpot.status, grandJackpot.time, meta]);

  return (
    <StyledSidebar
      breakPoint="md"
      width="220px"
      transitionDuration={0}
      collapsedWidth="74px"
      // defaultCollapsed={isMobile ? false : true}
      ref={target}
    >
      <Menu
        renderExpandIcon={({ open }) =>
          open ? <ExpandIcon /> : <ExpandIcon collapse={false} />
        }
      >
        {/* <Box display={collapsed ? 'none' : 'block'}>
          <Text
            color={'#686d7b'}
            fontSize={'14px'}
            p=" 8px 0px 8px 30px"
            fontWeight={600}
          >
            GAMES
          </Text>
        </Box> */}
        <MenuItem
          icon={<CoinflipIcon />}
          component={<NavLink to="/coinflip" />}
          style={{ marginTop: '14px' }}
        >
          Coin Flip
        </MenuItem>
        <MenuItem icon={<JackpotIcon />} component={<NavLink to="/jackpot" />}>
          Jackpot
        </MenuItem>
        <MenuItem
          icon={<DreamtowerIcon />}
          component={<NavLink to="/dream-tower" />}
        >
          Dream Tower
        </MenuItem>
        {/* <MenuItem icon={<PlinkoIcon />} component={<NavLink to="/plinko" />}>
          Plinko
        </MenuItem> */}
        <MenuItem
          icon={<CrashIcon />}
          component={<NavLink to="/crash" />}
          style={{ marginBottom: '23px' }}
        >
          Crash
        </MenuItem>
        {/* <MenuItem
          icon={<GrandJackpotIcon />}
          component={<NavLink to="/grandjackpot" />}
          className="submenu-grand"
          style={{ marginBottom: '23px' }}
        >
          Grand Jackpot
        </MenuItem> */}

        {/* <Box display={collapsed ? 'none' : 'block'}>
          <Text
            color={'#686d7b'}
            fontSize={'14px'}
            p=" 8px 0px 8px 30px"
            fontWeight={600}
          >
            EXTRAS
          </Text>
        </Box> */}

        <MenuItem
          icon={<WagerRaceIcon />}
          component={<NavLink to="/daily-race" />}
        >
          Wager Race
        </MenuItem>

        {isAuthenticated && (
          <SubMenu
            label="DuelBots"
            icon={<DuelbotsIcon />}
            active={subMenuActive}
          >
            {/* <MenuItem component={<NavLink to="/duelbots/dashboard" />}>
              Dashboard
            </MenuItem> */}
            <MenuItem component={<NavLink to="/duelbots/myduelbots" />}>
              My Duelbots
            </MenuItem>
            <MenuItem component={<NavLink to="/duelbots/staking" />}>
              Staking
            </MenuItem>
          </SubMenu>
        )}

        <MenuItem icon={<FairIcon />} component={<NavLink to="/fair" />}>
          Provably Fair
        </MenuItem>
        <MenuItem
          icon={<DocsIcon />}
          component={
            <a
              href="https://docs.duel.win/welcome/overview"
              target="_blank"
              rel="noreferrer"
            />
          }
        >
          Docs
        </MenuItem>
      </Menu>
      <Menu>
        <MenuItem
          icon={<SupportIcon />}
          component={
            <a
              href="https://discord.gg/FrNpcfmmk2"
              target="_blank"
              rel="noreferrer"
            />
          }
          style={{ marginBottom: '14px' }}
        >
          Support
        </MenuItem>
      </Menu>
    </StyledSidebar>
  );
}
