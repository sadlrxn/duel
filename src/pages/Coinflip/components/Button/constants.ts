export const TRANSITION_CONFIG = {
  top: {
    from: {
      opacity: 0,
      transform: "scale(0.5) translateY(0%) translateX(-100%)",
    },
    to: {
      opacity: 1,
      transform: "scale(1) translateY(-110%) translateX(-50%)",
    },
  },

  bottom: {
    from: {
      opacity: 0,
      transform: "scale(0.5) translateY(0px) translateX(-100%)",
    },

    to: {
      opacity: 1,
      transform: "scale(1) translateY(35px) translateX(-50%)",
    },
  },
};
