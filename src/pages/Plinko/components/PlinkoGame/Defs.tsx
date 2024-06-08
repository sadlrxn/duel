import React from 'react';

export default function Defs() {
  return (
    <defs>
      <linearGradient id="multiplier-gradient-0" gradientTransform="rotate(90)">
        <stop offset="0%" stopColor="#00ACC4" />
        <stop offset="100%" stopColor="#0038FF" />
      </linearGradient>
      <linearGradient id="multiplier-gradient-1" gradientTransform="rotate(90)">
        <stop offset="0%" stopColor="#0075FF" />
        <stop offset="100%" stopColor="#1000C4" />
      </linearGradient>
      <linearGradient id="multiplier-gradient-2" gradientTransform="rotate(90)">
        <stop offset="0%" stopColor="#2700C4" />
        <stop offset="95.83%" stopColor="#AD00FF" />
      </linearGradient>
      <linearGradient id="multiplier-gradient-3" gradientTransform="rotate(90)">
        <stop offset="0%" stopColor="#6D00C4" />
        <stop offset="100%" stopColor="#FF008A" />
      </linearGradient>
      <linearGradient id="multiplier-gradient-4" gradientTransform="rotate(90)">
        <stop offset="0%" stopColor="#C400A4" />
        <stop offset="95.83%" stopColor="#FF0000" />
      </linearGradient>
      <linearGradient id="multiplier-gradient-5" gradientTransform="rotate(90)">
        <stop offset="0%" stopColor="#C40C00" />
        <stop offset="95.83%" stopColor="#FF7A00" />
      </linearGradient>
      <linearGradient id="multiplier-gradient-6" gradientTransform="rotate(90)">
        <stop offset="0%" stopColor="#C44600" />
        <stop offset="100%" stopColor="#FFC700" />
      </linearGradient>

      <radialGradient id="ball-gradient-1" cx="0" cy="0" r="1">
        <stop stopColor="#FFDF8E" />
        <stop offset="0.34375" stopColor="#FFD913" />
        <stop offset="0.703125" stopColor="#FF7C32" />
        <stop offset="1" stopColor="#C73D02" />
      </radialGradient>

      <radialGradient id="ball-gradient-2" cx="0" cy="0" r="1">
        <stop stopColor="#F68EFF" />
        <stop offset="0.34375" stopColor="#B413FF" />
        <stop offset="0.703125" stopColor="#7432FF" />
        <stop offset="1" stopColor="#4D02C7" />
      </radialGradient>

      <radialGradient
        id="pin-bounce-gradient"
        cx="0"
        cy="0"
        r="1"
        gradientUnits="userSpaceOnUse"
        gradientTransform="translate(7 7) rotate(90) scale(7)"
      >
        <stop stopColor="#F1CAFF" />
        <stop offset="1" stopColor="#F0C4FF" stopOpacity="0.8" />
      </radialGradient>

      <radialGradient id="container-stroke-gradient" cx="0" cy="0" r="1">
        <stop offset="0.0861746" stopColor="#FF0000" />
        <stop offset="0.285704" stopColor="#FFD600" />
        <stop offset="0.451019" stopColor="#24FF00" />
        <stop offset="0.598033" stopColor="#00FFF3" />
        <stop offset="0.778823" stopColor="#2100FF" />
        <stop offset="0.932357" stopColor="#DB00FF" />
      </radialGradient>

      <filter id="pin-shadow-1" x="-50%" y="-50%" width="200%" height="200%">
        <feDropShadow dx="0" dy="0" stdDeviation="1" floodColor="#D932AF" />
      </filter>

      <filter
        id="pin-shadow-2"
        x="-500%"
        y="-500%"
        width="1000%"
        height="1000%"
      >
        <feDropShadow dx="0" dy="0" stdDeviation="6" floodColor="#FFFFFF" />
        <feDropShadow dx="0" dy="0" stdDeviation="1" floodColor="#D932AF" />
      </filter>

      <filter id="ball-shadow" x="-50%" y="-50%" width="200%" height="200%">
        <feDropShadow dx="0" dy="0" stdDeviation="4" floodColor="#FF4D0180" />
      </filter>

      <filter id="container-blur">
        <feGaussianBlur stdDeviation="5" in="FillPaint" />
      </filter>
    </defs>
  );
}
