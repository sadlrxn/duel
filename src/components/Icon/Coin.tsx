import React from "react";

export default function Coin({size = 19}: {size?: number}) {
  return (
    <svg
      width={size}
      height={size}
      viewBox={`0 0 ${size} ${size}`}
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
    >
      <circle
        cx={size / 2}
        cy={size / 2}
        r={size / 2.3}
        fill="#FFE24B"
        stroke="#FFB31F"
        strokeWidth={size / 6.5}
      />
    </svg>
  );
}

