import React, { FC, useCallback } from "react";
import styled from "styled-components";
import { PaginationComponentProps } from "react-data-table-component";
import { ReactComponent as PrevIcon } from "assets/imgs/icons/prev.svg";
import { ReactComponent as NextIcon } from "assets/imgs/icons/next.svg";
import { Text } from "components/Text";
import { Button } from "components/Button";

const PaginationWrapper = styled.nav`
  display: flex;
  flex: 1 1 auto;
  justify-content: flex-end;
  align-items: center;
  box-sizing: border-box;
  padding-right: 8px;
  padding-left: 8px;
  padding-bottom: 15px;
  border-bottom: 1px solid #2a3d57;
  width: 100%;
`;

const pad = (d: number) => {
  return d < 10 ? "0" + d.toString() : d.toString();
};

const Pagination: FC<PaginationComponentProps> = ({
  rowsPerPage,
  rowCount,
  onChangePage,
  currentPage,
}) => {
  const totalPage = Math.ceil(rowCount / rowsPerPage);

  const range = `Page ${pad(currentPage)}/${pad(totalPage)}`;

  const handlePrevious = useCallback(() => {
    if (currentPage === 1) return;
    onChangePage(currentPage - 1, rowCount);
  }, [currentPage, onChangePage, rowCount]);

  const handleNext = useCallback(() => {
    if (currentPage === totalPage) return;
    onChangePage(currentPage + 1, rowCount);
  }, [currentPage, onChangePage, rowCount, totalPage]);

  return (
    <PaginationWrapper>
      <Text color={"#BACFEE"} fontSize="14px" mr="30px">
        {range}
      </Text>
      <Button onClick={handlePrevious} background="transparent">
        <PrevIcon />
      </Button>
      <Button onClick={handleNext} background="transparent">
        <NextIcon />
      </Button>
    </PaginationWrapper>
  );
};

export default Pagination;
