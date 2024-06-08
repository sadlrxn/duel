import React from "react";
import { useTheme } from "styled-components";

import { Duel, Ana, RandomIcon } from "components";
import { Circle } from "./Coin.styles";

export default function Coin({
  side = "duel",
  size = 72,
  scale = 0.54,
  ...props
}: any) {
  const theme = useTheme();

  return (
    <Circle
      width={size}
      height={size}
      minWidth={size}
      minHeight={size}
      background={
        side === "duel"
          ? theme.colors.gradients.duel
          : side === "ana"
          ? theme.colors.gradients.ana
          : theme.colors.gradients.rnd
      }
      {...props}
    >
      {side === "duel" ? (
        <Duel size={size * scale} />
      ) : side === "ana" ? (
        <Ana size={size * scale} />
      ) : (
        <RandomIcon size={size * scale} />
      )}
    </Circle>
  );
}
