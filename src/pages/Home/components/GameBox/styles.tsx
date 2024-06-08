import { Text } from "components/Text";
import styled from "styled-components";

export const GameStatus = styled(Text)`
  display: flex;
  height: 28px;

  align-items: center;
  background: rgba(79, 255, 139, 0.2);
  padding: 0px 10px;
  border-radius: 22.8571px;

  svg {
    width: 7.18px;
    height: 7.18px;
    margin-right: 5px;
  }

  div {
    font-family: "Inter";
    font-style: normal;
    font-weight: 400;
    font-size: 14.359px;
    line-height: 17px;

    color: #4fff8b;
  }
`;
