import Slider from 'react-slick';
import styled from 'styled-components';
import { Box } from 'components';
import firstImg from 'assets/imgs/home/carousel/slide_1.jpg';
import secondImg from 'assets/imgs/home/carousel/slide_2.jpg';
import thirdImg from 'assets/imgs/home/carousel/slide_3.jpg';
import fourthImg from 'assets/imgs/home/carousel/slide_4.jpg';

import firstMobileImg from 'assets/imgs/home/carousel/slide_mobile_1.jpg';
import secondMobileImg from 'assets/imgs/home/carousel/slide_mobile_2.jpg';
import thirdMobileImg from 'assets/imgs/home/carousel/slide_mobile_3.jpg';
import fourthMobileImg from 'assets/imgs/home/carousel/slide_mobile_4.jpg';

export const StyledCarousel = styled(Slider)`
  width: 100%;

  box-sizing: border-box;
  border: 1px solid #4f617b;
  border-radius: 11px;
  height: 628px;
  .width_700 & {
    height: 372px;
  }

  .slick-list {
    width: 100%;
    overflow: visible;

    .slick-active {
    }
  }

  .slick-arrow:after {
    width: 18px;
    height: 18px;
    display: block;
    content: '';
    background: transparent;
    opacity: 0.2;
    position: absolute;
    top: 0;
    filter: blur(18px);
  }

  .slick-arrow {
    z-index: 10;

    :before {
      background: #96a8c2;
      height: 44px;
      position: absolute;
      top: -11px;
      left: 32px;
      pointer-events: none;
      width: 2px;
    }

    :hover:before {
      background: #4fff8b;
    }

    :hover:after {
      background: #4fff8b;
    }

    :hover svg path {
      stroke: #4fff8b;
    }
  }

  svg {
    width: 18px;
    height: 18px;
  }

  .slick-prev {
    left: 13px;
    transform: rotate(180deg);
  }

  .slick-next {
    right: 13px;
  }

  .slick-prev:before,
  .slick-next:before {
    content: '';
  }

  .slick-dots {
    bottom: -1px;
  }

  .slick-dot-box {
    width: 145px;
    height: 26px;
    margin: auto;
    display: flex;
    align-items: center;

    padding: 0px;

    svg {
      position: absolute;
      width: 145px;
      height: 26px;
      left: 50%;
      transform: translateX(-50%);
    }

    ul {
      display: flex;
      justify-content: space-between;
      padding: 0px 30px;
      margin: 0;
    }
    li {
      display: flex;
      align-items: center;
      width: auto !important;
    }
    li button {
      background: #121c26;
      width: 11px;
      height: 4px;
      padding: 0;
      border-radius: 20.8609px;
    }

    li button:before {
      display: none;
    }

    .slick-active button {
      background: #4fff8b;
      width: 30px !important;
    }
  }

  @keyframes float-icon {
    from {
      transform: translateY(0);
    }

    to {
      transform: translateY(10px);
    }
  }
`;

const Slide = styled(Box)`
  position: relative;
  display: flex;
  flex-direction: column;
  padding: 51.8px 34px;
  border-radius: 11px;
  background-size: 100% 100%;
  background-repeat: no-repeat;
  background-position: center top;
  height: 626px;
  .width_700 & {
    background-size: 100% 100%;
    padding: 65px 75px;
    height: 370px;
  }
`;

export const StyledSlideA = styled(Slide)`
  align-items: center;
  gap: 17.73px;
  background-image: url(${firstMobileImg});
  .width_700 & {
    gap: 20.36px;
    background-image: url(${firstImg});
    align-items: start;
    justify-content: end;
  }
`;

export const StyledSlideB = styled(Slide)`
  flex-direction: column;
  align-items: center;
  gap: 17.73px;
  /* overflow: hidden; */
  background-image: url(${secondMobileImg});
  .width_700 & {
    gap: 20.36px;
    background-image: url(${secondImg});
  }
`;

export const StyledSlideC = styled(Slide)`
  background-image: url(${thirdMobileImg});
  align-items: center;
  justify-content: center;
  gap: 24px;
  overflow: hidden;
  .width_700 & {
    overflow: visible;
    background-image: url(${thirdImg});
  }
`;

export const StyledSlideD = styled(Slide)`
  background-image: url(${fourthMobileImg});
  justify-content: end;
  align-items: center;
  gap: 20px;
  overflow: hidden;
  .width_700 & {
    justify-content: start;
    overflow: visible;
    align-items: start;
    background-image: url(${fourthImg});
  }
`;
