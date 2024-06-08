/// <reference types="react-scripts" />

declare module '*.mp3' {
  const src: string;
  export default src;
}

declare module '*.wav' {
  const src: string;
  export default src;
}

declare module '*.splinecode' {
  const src: string;
  export default src;
}

declare var balanceDecimal: number = 5;
