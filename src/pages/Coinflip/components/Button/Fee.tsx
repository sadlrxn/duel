import { FairnessIcon } from "components/Icon";
import ToolTipButton from "./Tooltip";

interface IProps {
  percentage: number;
  className?: string;
  tooltipPosition?: "top" | "bottom";
}

const GameFee = ({
  percentage,
  className = "",
  tooltipPosition = "top",
}: IProps): JSX.Element => {
  return (
    <ToolTipButton
      text={`${percentage}% Fee`}
      icon={<FairnessIcon />}
      {...{ className, tooltipPosition }}
    >
      {/* <p>
        The platform service fee for this game is {percentage}%. For more
        information on service fees, please visit our{" "}
        <a
          target="_blank"
          rel="noreferrer"
          href="https://app.gitbook.com/o/5VfsngU64V7lyHejBuB2/s/2fN8tRg5IZECQqhTBEyh/games/coin-flip"
          className="game-data-link"
        >
          Knowledge Base.
        </a>
      </p> */}
    </ToolTipButton>
  );
};

export default GameFee;
