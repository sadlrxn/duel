import React from "react";
import { ComponentMeta } from "@storybook/react";

import { Box } from "components/Box";
import GameBox from ".";

import CoinflipImg from "assets/imgs/home/coinflip.png";
import BlackjackImg from "assets/imgs/home/blackjack.png";
import JackpotImg from "assets/imgs/home/jackpot.png";

export default {
  title: "home/GameBox",
  component: GameBox,
  argTypes: {},
} as ComponentMeta<typeof GameBox>;

export const Default = () => (
  <Box mt={"50px"} width="500px">
    <GameBox img={CoinflipImg} />

    <GameBox img={BlackjackImg} />
    <GameBox img={JackpotImg} />
  </Box>
);
