export type CoinflipTheme = typeof dark;

export const light: CoinflipTheme = {
  title: "#4F617B",
  name: "#4E607A",
  greenDark: "#1E4B40",
  yellow: "#FFE24B",
  duel: "#4b88ff",
  ana: "#5d24ff",
  rnd: "#FFBF74",
  duelSecondary: "#85ffe2",
  anaSecondary: "#DFA7FF",
  border: "#1A293D",
  borderSecondary: "#0D141E",
  borderButton: "#374355",
  private: "#8192AA",
  gradients: {
    background:
      "linear-gradient(268.42deg, #1A293D 0%, rgba(26, 41, 61, 0) 100.18%)",
    duel: "linear-gradient(149.33deg, #4FFF8B 7.45%, #4B88FF 98.66%)",
    ana: "linear-gradient(149.33deg, #E392FF 7.45%, #9A1AFF 98.66%)",
    chip: "linear-gradient(90deg, #503B00 0%, #2F2814 100%)",
    side_duel: "linear-gradient(90deg, #25544D 0.62%, #1A293D 51.94%)",
    side_ana: "linear-gradient(90deg, #1A293D 51.8%, #422554 100%)",
    button: "linear-gradient(180deg, #070B10 0%, rgba(7, 11, 16, 0) 162.5%)",
    ok: "linear-gradient(180deg,rgba(7,11,16,0),#1a382a)",
    cancel: "linear-gradient(180deg, #070b10 0%, #4a170b 100%)",
  },
};

export const dark = {
  title: "#4F617B",
  name: "#4E607A",
  greenDark: "#1E4B40",
  yellow: "#FFE24B",
  duel: "#4b88ff",
  ana: "#5d24ff",
  rnd: "#FFBF74",
  duelSecondary: "#85ffe2",
  anaSecondary: "#DFA7FF",
  border: "#1A293D",
  borderSecondary: "#0D141E",
  borderButton: "#374355",
  private: "#8192AA",
  gradients: {
    background:
      "linear-gradient(268.42deg, #1A293D 0%, rgba(26, 41, 61, 0) 100.18%)",
    duel: "linear-gradient(149.33deg, #4FFF8B 7.45%, #4B88FF 98.66%)",
    ana: "linear-gradient(149.33deg, #E392FF 7.45%, #9A1AFF 98.66%)",
    chip: "linear-gradient(90deg, #503B00 0%, #2F2814 100%)",
    side_duel: "linear-gradient(90deg, #25544d 0%, transparent 100%)",
    side_ana: "linear-gradient(90deg, transparent 0%, #422554 100%)",
    button: "linear-gradient(180deg, #070B10 0%, rgba(7, 11, 16, 0) 162.5%)",
    ok: "linear-gradient(180deg,rgba(7,11,16,0),#1a382a)",
    cancel: "linear-gradient(180deg, #070b10 0%, #4a170b 100%)",
  },
};
