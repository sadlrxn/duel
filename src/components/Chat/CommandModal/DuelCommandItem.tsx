import React, { useCallback } from "react";
import styled from "styled-components";

import { DuelCommand } from "config";
import { useAppDispatch } from "state";
import { setChatMsg } from "state/actions";

interface DuelCommandItemProps {
  command: DuelCommand;
  onDismiss?: () => void;
  inputRef: React.RefObject<HTMLInputElement | HTMLTextAreaElement>;
}

export default function DuelCommandItem({
  command,
  onDismiss,
  inputRef,
}: DuelCommandItemProps) {
  const dispatch = useAppDispatch();

  const onClickCommand = useCallback(() => {
    dispatch(
      setChatMsg(
        `/${command.pattern.substring(1, command.pattern.indexOf(" "))} `
      )
    );
    onDismiss && onDismiss();
    inputRef.current && inputRef.current.focus();
  }, [dispatch, onDismiss, command.pattern, inputRef]);
  return (
    <StyledDuelCommandItem onClick={onClickCommand}>
      {command.pattern}
    </StyledDuelCommandItem>
  );
}

const StyledDuelCommandItem = styled.div`
  color: white;
  border-radius: 3px;
  padding: 11px 7px;
  cursor: pointer;
  user-select: none;
  &:hover {
    background: #26374e;
  }
`;
