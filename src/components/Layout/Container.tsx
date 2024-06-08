import { Box, BoxProps } from "../Box";

const Container: React.FC<React.PropsWithChildren<BoxProps>> = ({
  children,
  ...props
}) => (
  <Box px={["16px", "24px"]} mx="auto" maxWidth="1645px" {...props}>
    {children}
  </Box>
);

export default Container;
