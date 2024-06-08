// interface CardProps {
//   radius?: number;
//   count?: number;
//   width?: number;
//   height?: number;
// }

const PI = 3.141592;

export default function Card({
  radius = 870,
  count = 50,
  width = 120,
  height = 146,
  ...props
}: any) {
  const x0 = 0.101807;
  const x = (height * PI) / count;
  const y1 = radius * (1 - Math.sqrt(1 - (PI / count) * (PI / count)));
  const y2 = height;
  const path = `M${x0} ${y1}L${x} ${y2}Q${width / 2} ${y2 - y1 * 1.5} ${
    width - x0 - x
  } ${y2}L${width - x0} ${y1}Q${width / 2} ${-y1} ${x0} ${y1} Z`;

  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      width={width}
      height={height}
      viewBox={`0 0 ${width} ${height}`}
      fill="url(#paint0_linear_5316_23169)"
      {...props}
    >
      <path d={path} />
      <defs>
        <linearGradient
          id="paint0_linear_5316_23169"
          x1={width / 2}
          y1="0"
          x2={width / 2}
          y2={height}
          gradientUnits="userSpaceOnUse"
        >
          <stop stopColor="#1A2A3E" />
          <stop offset="1" stopColor="#142131" />
        </linearGradient>
      </defs>
    </svg>
  );
}
