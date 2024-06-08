import { ReactComponent as FairnessIocon } from "assets/imgs/icons/fairness.svg";
import ToolTipButton from "./Tooltip";

interface IProps {
  className?: string;
  tooltipPosition?: "top" | "bottom";
}

const GameFairness = ({
  className = "",
  tooltipPosition = "top",
}: IProps): JSX.Element => {
  return (
    <ToolTipButton
      text="Fairness"
      icon={<FairnessIocon color="#4F617B" />}
      {...{ className, tooltipPosition }}
    >
      {/* <p>
        Duelana uses{" "}
        <a
          href="https://www.random.org/"
          target="_blank"
          rel="noreferrer"
          className="game-data-link"
        >
          Random.org
        </a>{" "}
        to provide provably fair number generation for all of our games. For
        more information on how this process works, please visit our{" "}
        <a
          target="_blank"
          rel="noreferrer"
          href="https://app.gitbook.com/o/5VfsngU64V7lyHejBuB2/s/2fN8tRg5IZECQqhTBEyh/more-info/provably-fair"
          className="game-data-link"
        >
          Knowledge Base.
        </a>
      </p> */}
    </ToolTipButton>
  );
};

export default GameFairness;
