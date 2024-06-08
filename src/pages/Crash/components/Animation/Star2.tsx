export default function Star2({ size = 30, angle = 0, ...props }: any) {
  return (
    <svg
      width={`${size}px`}
      height={`${size}px`}
      transform={`rotate(${angle}deg)`}
      viewBox="0 0 30 30"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
      {...props}
    >
      <path
        d="M16.9775 12.3998L18.4764 5.80959L15.5136 11.8828L8.59742 0.0869141L12.8576 13.2587L6.59009 11.9529L12.3054 14.8098L0.708252 22.478L13.5238 17.5265L12.0336 24.1168L14.9964 18.0436L21.9038 29.8395L17.6525 16.6677L23.9112 17.9735L18.1959 15.1166L29.793 7.44838L16.9775 12.3998Z"
        fill="white"
      />
    </svg>
  );
}
