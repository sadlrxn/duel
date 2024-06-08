import { scales, variants } from "./types";

export type ButtonThemeType = typeof dark;

export const light = {
  background: "",
  text: "",
  textSecondary: "",
};

export const dark = {};

export const scaleVariants = {
  [scales.DEFAULT]: {},
  [scales.LG]: {
    height: "52px",
  },
  [scales.MD]: {
    height: "46px",
  },
  [scales.SM]: {
    height: "42px",
  },
  [scales.XS]: {
    height: "38px",
  },
};

export const styleVariants = {
  [variants.PRIMARY]: {
    backgroundColor: "success",
    borderRadius: "32px",
  },
  [variants.SECONDARY]: {
    backgroundColor: "#242F42",
    borderRadius: "5px",
  },
};
