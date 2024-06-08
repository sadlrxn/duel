import React from "react";
import { ComponentMeta } from "@storybook/react";

import { Box } from "components/Box";
import TopItem, { UserInfo } from ".";
import { Options } from "react-select";

export default {
  title: "home/TopItem",
  component: TopItem,
  argTypes: {},
} as ComponentMeta<typeof TopItem>;

const TOP_WIN_OPTIONS: Options<{ label: string; value: string }> = [
  { label: "Last 24 hours", value: "last24hours" },
  { label: "Last 7 days", value: "last7days" },
  { label: "Last 30 days", value: "last30days" },
];

const TOP_PLAYERS_OPTIONS: Options<{ label: string; value: string }> = [
  { label: "All Games", value: "all" },
  { label: "Blackjack", value: "blackjack" },
  { label: "Coinflip", value: "coinflip" },
  { label: "Jackpot", value: "jackpot" },
];

const win_data: UserInfo[] = [
  {
    user: {
      avatar: "asdf",
      name: "eroist",
    },
    info: {
      earned: "145,232",
    },
  },
];

const player_data: UserInfo[] = [
  {
    user: {
      avatar: "asdf",
      name: "eroist",
    },
    info: {
      exp: "145,232",
      tier: "Dueler V",
    },
  },
];

export const Default = () => (
  <>
    <Box width={"1000px"} mt={"50px"}>
      <TopItem
        title="Top Win"
        options={TOP_WIN_OPTIONS}
        isWinComp={true}
        data={win_data}
      />
    </Box>
    <Box width={"1000px"} mt={"50px"}>
      <TopItem
        title="Top Players"
        options={TOP_PLAYERS_OPTIONS}
        isWinComp={false}
        data={player_data}
      />
    </Box>
  </>
);
