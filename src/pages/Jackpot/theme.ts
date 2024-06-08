export type JackpotTheme = typeof dark;

const breakpointMap: { [key: string]: number } = {
  sm: 650,
  md: 852,
  xl: 900, // 1150
};

const eq = {
  sm: `@element #main and (min-width: ${breakpointMap.sm}px)`,
  md: `@element #main and (min-width: ${breakpointMap.md}px)`,
  xl: `@element #main and (min-width: ${breakpointMap.xl}px)`,
};

export const light: JackpotTheme = {
  eq,
  modal: "#202f44",
  input: "#03060999",
  progress: "#182738",
};

export const dark = {
  eq,
  modal: "#202f44",
  input: "#03060999",
  progress: "#182738",
};
