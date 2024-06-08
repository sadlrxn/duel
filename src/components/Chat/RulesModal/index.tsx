import React, { useEffect, useState } from 'react';
import { useClickAnyWhere, useWindowSize } from 'usehooks-ts';
import styled from 'styled-components';

import { Modal, ModalProps } from 'components/Modal';
import { useMatchBreakpoints } from 'hooks';
import { chatWidth } from 'theme';

export default function RulesModal({ onDismiss }: ModalProps) {
  const [rendered, setRendered] = useState(false);
  const { width } = useWindowSize();
  const { isMobile } = useMatchBreakpoints();

  useEffect(() => {
    setRendered(true);
  }, []);

  useClickAnyWhere(() => {
    rendered && onDismiss && onDismiss();
  });

  const rules = [
    'Treat everyone with kindness and respect.',
    'No harassment, sexism, racism, or hate speech will be tolerated.',
    'No spam, shilling, or self-promotion will be tolerated.',
    'No NSFW or obscene content will be tolerated.',
    'No begging for tips or loans.',
    'Respect the Duel team and staff members on duty.',
    'English only.',
    'Type “/” to see all commands.',
    'Failing to adhere to these rules may result in mute, kick, or ban per discretion by the Duel team. Any impersonations of team members or server bots will result in an instant ban.',
    'No calling games unless asked'
  ];

  return (
    <StyledRuleModal
      style={{ left: isMobile ? `10px` : `${width - chatWidth}px` }}
    >
      <h2>Chat Rules</h2>
      <ul>
        {rules.map((rule, i) => (
          <li key={i}>{rule}</li>
        ))}
      </ul>
    </StyledRuleModal>
  );
}

const StyledRuleModal = styled(Modal)`
  position: fixed;
  padding: 25px;
  width: 350px;
  height: max-content;
  bottom: 15px;
  background: linear-gradient(180deg, #0c1725 0%, #18283e 100%);
  border: 1px solid #26374e;
  box-shadow: 0px 4px 22px rgba(0, 0, 0, 0.4);
  border-radius: 6px;
  display: flex;
  flex-direction: column;
  gap: 16px;
  h2 {
    font-family: 'Inter';
    font-style: normal;
    font-weight: 600;
    font-size: 18px;
    line-height: 22px;
    letter-spacing: 0.18em;
    color: #fff;
    margin: 0px;
    text-transform: uppercase;
  }
  ul {
    color: #b2d1ff;
    font-family: 'Inter';
    font-style: normal;
    font-weight: 500;
    font-size: 14px;
    line-height: 20px;
    padding-inline-start: 25px;
    margin: 0px;
  }
  ${({ theme }) => theme.mediaQueries.sm} {
    transform: translateX(-100%);
  }
`;
