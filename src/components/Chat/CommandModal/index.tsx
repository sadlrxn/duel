import React, { useEffect, useState, useMemo, useRef } from "react";
import { AnimatePresence, domAnimation, LazyMotion, m } from "framer-motion";
import { useClickAnyWhere } from "usehooks-ts";
import styled from "styled-components";

import { Flex } from "components/Box";
import { Text } from "components/Text";
import { checkUserHasPermission } from "config";
import { useAppSelector } from "state";
import DuelCommandItem from "./DuelCommandItem";
import {
  animationHandler,
  animationMap,
  animationVariants,
  appearAnimation,
  disappearAnimation,
} from "utils/animationToolkit";

interface CommandModalProps {
  show: boolean;
  inputRef: React.RefObject<HTMLInputElement | HTMLTextAreaElement>;
  onDismiss: () => void;
}

export default function CommandModal({
  show,
  onDismiss,
  inputRef,
}: CommandModalProps) {
  const [rendered, setRendered] = useState(false);
  const { msg, commands } = useAppSelector((state) => state.chat);
  const { role } = useAppSelector((state) => state.user);
  const animationRef = useRef<HTMLDivElement>(null);

  const availableCommands = useMemo(() => {
    return commands
      .filter((c) => checkUserHasPermission(role, c.role))
      .filter(
        (command) =>
          command.pattern.indexOf(
            msg.substring(
              1,
              msg.indexOf(" ") > -1 ? msg.indexOf(" ") : msg.length
            )
          ) > -1
      );
  }, [msg, role, commands]);

  useEffect(() => {
    setRendered(true);
  }, []);

  useEffect(() => {
    if (availableCommands.length === 0) rendered && onDismiss && onDismiss();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [availableCommands]);

  useClickAnyWhere(() => {
    rendered && onDismiss && onDismiss();
  });

  return (
    <LazyMotion features={domAnimation}>
      <AnimatePresence>
        {show && (
          <StyledCommandModal
            ref={animationRef}
            onAnimationStart={() => animationHandler(animationRef.current)}
            {...animationMap}
            variants={animationVariants}
            transition={{ duration: 0.3 }}
          >
            <Text
              textAlign="center"
              color="white"
              textTransform="uppercase"
              mb={11}
              fontWeight={600}
            >
              Available Commands
            </Text>
            <Flex flexDirection="column" gap={4}>
              {availableCommands.map((command) => (
                <DuelCommandItem
                  key={command.pattern}
                  command={command}
                  inputRef={inputRef}
                />
              ))}
            </Flex>
          </StyledCommandModal>
        )}
      </AnimatePresence>
    </LazyMotion>
  );
}

const StyledCommandModal = styled(m.div)`
  position: fixed;
  padding: 11px;
  width: 339px;
  right: 15px;
  bottom: 140px;
  background: linear-gradient(180deg, #0c1725 0%, #18283e 100%);
  border: 1px solid #26374e;
  box-shadow: 0px 4px 22px rgba(0, 0, 0, 0.4);
  border-radius: 6px;
  flex-direction: column;
  gap: 8px;
  &.appear {
    animation: ${appearAnimation} 0.3s ease-in-out forwards;
  }
  &.disappear {
    animation: ${disappearAnimation} 0.3s ease-in-out forwards;
  }
`;
