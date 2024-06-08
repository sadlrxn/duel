import { ReactComponent as HistoryIcon } from "assets/imgs/icons/history.svg";
import ToolTipButton from "./Tooltip";
// import { formatThousandNumber } from "@utils/utils";

interface IStats {
  total: number;
  games_won: number;
  games_lose: number;
  profits: number;
}

interface IProps {
  onClick?: () => void;
  stats?: IStats;
  className?: string;
  tooltipPosition?: "top" | "bottom";
}

const GameHistory = ({
  onClick,
  className = "",
  tooltipPosition = "top",
}: IProps): JSX.Element => {
  return (
    <ToolTipButton
      text="History"
      icon={<HistoryIcon color="#4F617B" />}
      {...{ className, tooltipPosition, onClick }}
    >
      {/* <p>
        <b>Total:</b> {formatThousandNumber(stats.total)}
      </p>
      <p>
        <b>Won:</b> {formatThousandNumber(stats.games_won)}
      </p>
      <p>
        <b>Lost:</b> {formatThousandNumber(stats.games_lose)}
      </p>
      <p>
        <b>Profits:</b> {formatThousandNumber(stats.profits)}
      </p> */}
    </ToolTipButton>
  );
};

export default GameHistory;
