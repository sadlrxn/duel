import React from 'react';
import styled, { keyframes } from 'styled-components';

const WIDTH = 14;
const HEIGHT = 20;
const SPACING = 25;
const SPREAD = 3;
const OPACITY = 0.4;

const flash_1 = keyframes`
  0%, 100% { background: rgba(0,247,165,1);
  box-shadow: 0px ${HEIGHT / 6}px ${
  WIDTH * 2
}px ${SPREAD}px rgba(0,247,165,1);} 
  50% { background: rgba(0,247,165,${OPACITY});
  box-shadow: 0px ${HEIGHT / 6}px ${
  WIDTH * 2
}px ${SPREAD}px rgba(0,247,165,0.2);}
`;

const flash_2 = keyframes`
  0%, 100% { background: rgba(0,255,255,1);
  box-shadow: 0px ${HEIGHT / 6}px ${WIDTH * 2} ${SPREAD}px rgba(0,255,255,1);} 
  50% { background: rgba(0,255,255,${OPACITY});
  box-shadow: 0px ${HEIGHT / 6}px ${WIDTH * 2} ${SPREAD}px rgba(0,255,255,0.2);}
`;

const flash_3 = keyframes`
  0%, 100% { background: rgba(247,0,148,1);
  box-shadow: 0px ${HEIGHT / 6}px ${WIDTH * 2} ${SPREAD}px rgba(247,0,148,1);} 
  50% { background: rgba(247,0,148,${OPACITY});
  box-shadow: 0px ${HEIGHT / 6}px ${WIDTH * 2} ${SPREAD}px rgba(247,0,148,0.2);}
`;

const Ul = styled.ul`
  /* text-align: center; */
  white-space: nowrap;
  position: fixed;
  z-index: 1;
  pointer-events: none;
  left: -15px;
  top: 39px;
  width: 100vw;
  height: 80px;
  overflow: hidden;
  li {
    position: relative;
    animation-fill-mode: both;
    animation-iteration-count: infinite;
    list-style: none;
    margin: 0;
    padding: 0;
    display: block;
    width: ${WIDTH}px;
    height: ${HEIGHT}px;
    border-radius: 50%;
    margin: ${SPACING / 2}px;
    display: inline-block;
    background: rgba(0, 247, 165, 1);
    box-shadow: 0px ${HEIGHT / 6}px ${WIDTH * 2}px ${SPREAD}px
      rgba(0, 247, 165, 1);
    animation-name: ${flash_1};
    animation-duration: 2s;
    &:nth-child(2n + 1) {
      background: rgba(0, 255, 255, 1);
      box-shadow: 0px ${HEIGHT / 6}px ${WIDTH * 2}px ${SPREAD}px
        rgba(0, 255, 255, 0.5);
      animation-name: ${flash_2};
      animation-duration: 0.4s;
    }
    &:nth-child(4n + 2) {
      background: rgba(247, 0, 148, 1);
      box-shadow: 0px ${HEIGHT / 6}px ${WIDTH * 2}px ${SPREAD}
        rgba(247, 0, 148, 1);
      animation-name: ${flash_3};
      animation-duration: 1.1s;
    }
    &:nth-child(odd) {
      animation-duration: 1.8s;
    }
    &:nth-child(3n + 1) {
      animation-duration: 1.4s;
    }
    &:before {
      content: '';
      position: absolute;
      background: #222;
      width: ${WIDTH - 2}px;
      height: ${HEIGHT / 3}px;
      border-radius: 3px;
      top: ${0 - HEIGHT / 6}px;
      left: 1px;
    }
    &:after {
      content: '';
      top: ${-HEIGHT / 2}px;
      left: ${WIDTH - 3}px;
      position: absolute;
      width: ${SPACING + 12}px;
      height: ${(HEIGHT / 3) * 2}px;
      border-bottom: solid #222 2px;
      border-radius: 50%;
    }
    &:last-child:after {
      content: none;
    }
    &:first-child {
      margin-left: -${SPACING}px;
    }
  }
`;

export default function LightRope() {
  return (
    <>
      <Ul>
        {Array(70)
          .fill(1)
          .map((_, index) => {
            return <li key={index} />;
          })}
      </Ul>
    </>
  );
}
