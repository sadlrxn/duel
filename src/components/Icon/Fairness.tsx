export default function Fairness({
  size = 15,
  color = "#768BAD",
}: {
  size?: number;
  color?: string;
}) {
  return (
    <svg
      viewBox="0 0 14 16"
      width={size}
      height={size}
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
    >
      <path
        d="M6.72551 15.5919C6.65653 15.5919 6.58754 15.581 6.52146 15.5592C3.76794 14.653 2.20819 12.9857 1.38548 11.747C0.363796 10.209 0 8.64488 0 7.69291V3.1647C0 2.90039 0.159025 2.66221 0.403734 2.56055L6.47572 0.0495591C6.63547 -0.0165197 6.81555 -0.0165197 6.9753 0.0495591L13.0473 2.56128C13.2913 2.66221 13.451 2.90039 13.451 3.16543V7.69364C13.451 8.64488 13.0872 10.2097 12.0655 11.7477C11.2428 12.9872 9.68308 14.6544 6.92955 15.5599C6.86348 15.5817 6.79449 15.5926 6.72551 15.5926V15.5919ZM1.30705 3.60184V7.69291C1.30705 8.91283 2.30913 12.6917 6.72551 14.2478C11.1419 12.6917 12.144 8.91283 12.144 7.69291V3.60184L6.72551 1.36024L1.30705 3.60184Z"
        fill={color}
      />
      <path
        d="M5.90932 10.4689C5.90932 10.4689 5.90423 10.4689 5.90133 10.4689C5.72778 10.4668 5.56295 10.3963 5.44168 10.2722L3.48691 8.26874C3.23494 8.01023 3.24003 7.59633 3.49853 7.34436C3.75704 7.09239 4.17094 7.09747 4.42291 7.35598L5.92166 8.89177L9.30547 5.58928C9.56398 5.33731 9.97715 5.34239 10.2299 5.6009C10.4818 5.85941 10.4767 6.2733 10.219 6.52528L6.36751 10.2845C6.24552 10.4036 6.08141 10.4704 5.91077 10.4704L5.90932 10.4689Z"
        fill={color}
      />
    </svg>
  );
}