import styled from "styled-components";
import { border, layout, space, typography, color } from "styled-system";

import { BaseButtonProps } from "./types";

const BaseButton = styled.button<BaseButtonProps>`
  border: 0;
  ${border}
  ${layout}
  ${space}
  ${typography}
  ${color}
  cursor: pointer;
`;

export default BaseButton;
