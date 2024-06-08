import React from "react";
import Chip1Img from "assets/imgs/coins/chip-1.png";
import Chip2Img from "assets/imgs/coins/chip-2.png";
import Chip3Img from "assets/imgs/coins/chip-3.png";
import { Flex } from "components/Box";
export default function PageSpinner() {
  return (
    <Flex justifyContent="center" alignItems={"center"} height="100%">
      <img src={Chip1Img} className="chip-1-spinner" alt="chip-1" />
      <img src={Chip2Img} className="chip-2-spinner" alt="chip-2" />
      <img src={Chip3Img} className="chip-3-spinner" alt="chip-3" />
    </Flex>
  );
}
