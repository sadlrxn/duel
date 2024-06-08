import { BoxProps } from "components/Box";

export type Handler = () => void;

export interface InjectedProps {
  onDismiss?: Handler;
}

export interface ModalProps extends InjectedProps, BoxProps {
  hideCloseButton?: boolean;
  onBack?: Handler;
  bodyPadding?: string;
  headerBackground?: string;
  minWidth?: string;
}
