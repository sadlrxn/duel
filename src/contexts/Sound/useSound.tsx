import { useContext } from "react";
import { SoundContext } from "./Provider";

const useSound = () => {
  const soundPlay = useContext(SoundContext);

  return soundPlay;
};

export default useSound;
