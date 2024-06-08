export default function Star1({ size = 30, angle = 0, style, ...props }: any) {
  return (
    <svg
      width={`${size}px`}
      height={`${size}px`}
      transform={`rotate(${angle}deg)`}
      viewBox="0 0 28 30"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
      style={style}
      {...props}
    >
      <path
        d="M14.9969 11.3568L13.8135 3.43447L13.1561 11.4181L0.437012 1.36627L10.8156 14.0822L3.3297 15.3355L10.877 16.0366L1.37495 29.4887L13.3928 18.5166L14.5849 26.4302L15.2336 18.4466L27.9527 28.4984L17.5741 15.7824L25.06 14.5292L17.5039 13.8369L27.0147 0.375977L14.9969 11.3568Z"
        fill="white"
      />
    </svg>
  );
}
