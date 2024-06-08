import React from 'react';
import { Settings } from 'react-slick';
import 'slick-carousel/slick/slick.css';
import 'slick-carousel/slick/slick-theme.css';

import { ReactComponent as ChevronIcon } from 'assets/imgs/icons/chevron.svg';
import { ReactComponent as TrapezeIcon } from 'assets/imgs/icons/trapeze.svg';

import { StyledCarousel } from './styles';
import SlideA from './SlideA';
// import SlideB from './SlideB';
import SlideC from './SlideC';
import SlideD from './SlideD';

function ArrowIcon(props: any) {
  const { className, onClick } = props;
  return (
    <div className={className} onClick={onClick}>
      <ChevronIcon />
    </div>
  );
}

export default function index() {
  const settings: Settings = {
    dots: true,
    infinite: true,
    fade: true,
    speed: 1000,
    slidesToShow: 1,
    slidesToScroll: 1,
    prevArrow: <ArrowIcon />,
    nextArrow: <ArrowIcon />,
    autoplay: false,
    autoplaySpeed: 10000,
    pauseOnHover: false,
    appendDots: dots => (
      <div>
        <div className="slick-dot-box">
          <TrapezeIcon />
          <ul> {dots} </ul>
        </div>
      </div>
    )
  };

  return (
    <StyledCarousel {...settings}>
      <SlideD />
      <SlideA />
      {/* <SlideB /> */}
      <SlideC />
    </StyledCarousel>
  );
}
