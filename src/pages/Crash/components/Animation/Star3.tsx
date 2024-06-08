export default function Star3({ size = 30, angle = 0, ...props }: any) {
  return (
    <svg
      width={`${size}px`}
      height={`${size}px`}
      transform={`rotate(${angle}deg)`}
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
      {...props}
    >
      <path
        d="M15.0197 19.002L22.9177 32.4805L18.0614 17.4334L31.3064 8.6785L16.6677 14.331L8.76096 0.852539L13.6434 15.8997L0.380859 24.6546L15.0197 19.002Z"
        fill="white"
      />
    </svg>
  );
}
