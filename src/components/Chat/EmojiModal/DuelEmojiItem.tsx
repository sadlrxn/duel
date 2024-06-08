import React, { useCallback } from "react";
import { LazyLoadImage } from "react-lazy-load-image-component";

import { Span } from "components/Text";
import { DuelEmoji, generateDuelEmojiUrl } from "config";

import { useAppDispatch, useAppSelector } from "state";
import { setChatMsg } from "state/actions";

import { DuelItemContainer } from "../styles";

interface DuelEmojiItemProps {
  emoji: DuelEmoji;
  onDismiss?: () => void;
  inputRef: React.RefObject<HTMLInputElement | HTMLTextAreaElement>;
  showShortName?: boolean;
}

export default function DuelEmojiItem({
  emoji,
  onDismiss,
  inputRef,
  showShortName = false,
}: DuelEmojiItemProps) {
  const dispatch = useAppDispatch();
  const { msg } = useAppSelector((state) => state.chat);

  const onClickEmoji = useCallback(() => {
    dispatch(setChatMsg(`${msg} :{${emoji.id}}`));
    onDismiss && onDismiss();
    inputRef.current && inputRef.current.focus();
  }, [dispatch, onDismiss, msg, emoji.id, inputRef]);

  return (
    <DuelItemContainer alignItems="center" onClick={onClickEmoji}>
      <LazyLoadImage
        width={35}
        height={35}
        src={generateDuelEmojiUrl(emoji)}
        alt={emoji.id}
      />
      {showShortName && <Span>:{emoji.id}</Span>}
    </DuelItemContainer>
  );
}
