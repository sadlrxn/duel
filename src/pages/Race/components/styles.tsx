import styled from 'styled-components';
import { Box } from 'components';
import { ReactComponent as CrownIcon } from 'assets/imgs/icons/crown.svg';
import startBurstImg from 'assets/imgs/race/star-burst.png';

export const StyledTable = styled.table`
  width: 100%;
  font-family: 'Inter';
  font-size: 14px;
  color: #b2d1ff;
  border-collapse: separate;
  border-spacing: 0px 10px;

  th {
    font-family: 'Inter';
    font-style: normal;
    font-weight: 500;
    font-size: 14px;
    line-height: 17px;

    letter-spacing: 0.17em;
    text-transform: uppercase;
    color: #b2d1ff;
    padding: 8px 20px;
  }

  tbody tr {
    background-color: #182738;
    cursor: pointer;
  }

  tbody tr.me {
    td {
      border-width: 1px 0px 1px 0px;
      border-style: solid;
      border-color: #49f884;
    }

    td:first-child {
      border-width: 1px 0px 1px 1px;
    }

    td:last-child {
      border-width: 1px 1px 1px 0px;
    }
  }

  td,
  tbody th {
    height: 50px;
    padding: 8px 20px;
  }

  tr td:first-child,
  tr th:first-child {
    border-top-left-radius: 8px;
    border-bottom-left-radius: 8px;
  }
  tr td:last-child,
  tr th:last-child {
    border-top-right-radius: 8px;
    border-bottom-right-radius: 8px;
  }
`;

export const StyledLeaderBox = styled(Box)`
  background-image: url(${startBurstImg});
  background-size: 100% auto;
  background-repeat: no-repeat;
  background-position: center top;
`;

export const StyledCrownIcon = styled(CrownIcon)`
  width: 15px;
  height: 15px;
  ${({ theme }) => theme.mediaQueries.md} {
    width: 23px;
    height: 23px;
  }
`;
