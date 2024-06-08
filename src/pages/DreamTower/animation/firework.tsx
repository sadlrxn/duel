import React, { useRef, useEffect } from "react";
import gsap from "gsap";
import styled, { css } from "styled-components";

const randomNumber = (val: number) => {
  return Math.random() * val;
};

function fireworkAnimation({ ref }: { ref: any }) {
  const tl = gsap.timeline({
    repeat: -1,
    repeatRefresh: true,
  });

  if (!ref) return tl;

  const q = gsap.utils.selector(ref);

  tl.fromTo(
    ref.current,
    {
      opacity: 0,
      top: () => `${70 + randomNumber(30)}%`,
    },
    {
      delay: () => randomNumber(1.7),
      opacity: 1,
      top: () => `${10 + randomNumber(50)}%`,
      duration: 1,
    }
  )
    .fromTo(
      q("div"),
      {
        left: () => `${randomNumber(100)}%`,
      },
      {
        left: () => `+=${randomNumber(30) - 20}%`,
        duration: 1,
        ease: "none",
      },
      "<"
    )
    .fromTo(
      q("div"),
      {
        width: "7px",
        opacity: 1,
      },
      {
        opacity: 0,
        width: () => `${650 + Math.floor(randomNumber(100)) - 200}px`,
        duration: () => 1.2 + randomNumber(0.3) - 0.6,
      },
      ">"
    );

  return tl;
}

const FireworkComponent = styled.div<{ color?: string }>`
  &,
  &::before,
  &::after {
    ${({ color }) => {
      const col = color ?? "#ff08";
      return css`
        content: "";
        position: absolute;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
        width: 0.5vmin;
        aspect-ratio: 1;
        /* random backgrounds */
        background: radial-gradient(circle, ${col} 0.2vmin, #0000 0) 50% 00%,
          radial-gradient(circle, ${col} 5px, #0000 0) 00% 50%,
          radial-gradient(circle, ${col} 7px, #0000 0) 50% 99%,
          radial-gradient(circle, ${col} 4px, #0000 0) 99% 50%,
          radial-gradient(circle, ${col} 5px, #0000 0) 80% 90%,
          radial-gradient(circle, ${col} 7px, #0000 0) 95% 90%,
          radial-gradient(circle, ${col} 7px, #0000 0) 10% 60%,
          radial-gradient(circle, ${col} 4px, #0000 0) 31% 80%,
          radial-gradient(circle, ${col} 5px, #0000 0) 80% 10%,
          radial-gradient(circle, ${col} 4px, #0000 0) 90% 23%,
          radial-gradient(circle, ${col} 5px, #0000 0) 45% 20%,
          radial-gradient(circle, ${col} 7px, #0000 0) 13% 24%;
        background-size: 7px 7px;
        background-repeat: no-repeat;
      `;
    }}
  }

  &::before {
    width: 100%;
    transform: translate(-50%, -50%) rotate(25deg);
  }

  &::after {
    width: 100%;
    transform: translate(-50%, -50%) rotate(37deg);
  }
`;

const Container = styled.div`
  position: absolute;
  width: 100%;
  height: 0;
`;

interface FireworkProps {
  color?: string;
}

export function Firework({ color }: FireworkProps) {
  const fireworkRef = useRef<any>(null);

  useEffect(() => {
    if (!fireworkRef) return;

    //@ts-ignore
    const tl = fireworkAnimation({ ref: fireworkRef });

    return () => {
      tl.kill();
    };
  }, []);

  return (
    <Container ref={fireworkRef}>
      <FireworkComponent color={color} />
    </Container>
  );
}
