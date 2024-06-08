import styled from "styled-components";
import Container from "./Container";
import { BoxProps } from "../Box";

const StyledPage = styled(Container)`
  min-height: calc(100vh - 65px);
  padding-top: 16px;
  padding-bottom: 16px;
`;

export const PageMeta: React.FC<React.PropsWithChildren> = () => {
  return <></>;
  // return (
  //   <head>
  //     <title>{pageTitle}</title>
  //     <meta property="og:title" content={title} />
  //     <meta property="og:description" content={description} />
  //     <meta property="og:image" content={image} />
  //   </head>
  // );
};

interface PageProps extends BoxProps {}

const Page: React.FC<React.PropsWithChildren<PageProps>> = ({
  children,
  ...props
}) => {
  return (
    <>
      <PageMeta />
      <StyledPage {...props}>{children}</StyledPage>
    </>
  );
};

export default Page;
