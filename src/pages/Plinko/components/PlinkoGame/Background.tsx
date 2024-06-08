import React from 'react';

interface BackgroundProps {
  win?: boolean;
  path: string;
  width?: number;
  height?: number;
}

export default function Background({ win = false, path }: BackgroundProps) {
  return (
    <>
      <g>
        <path
          d={path}
          fill="#1A172C"
          fillOpacity="0.8"
          // filter="url(#container-blur)"
          stroke={win ? 'url(#container-stroke-gradient)' : '#2C274F'}
        />
      </g>
    </>
  );
}
