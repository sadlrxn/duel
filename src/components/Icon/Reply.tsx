export default function Reply({
  size = 15,
  color = '#788CA9'
}: {
  size?: number;
  color?: string;
}) {
  return (
    <svg
      width={size}
      height={size}
      viewBox="0 0 15 14"
      fill="none"
      stroke={color}
      xmlns="http://www.w3.org/2000/svg"
    >
      <path
        d="M14 13V8.82964C14 7.40157 12.8423 6.24391 11.4143 6.24391H1M1 6.24391L6.26262 11.4877M1 6.24391L6.26262 1"
        strokeWidth="1.73333"
        strokeMiterlimit="10"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
    </svg>
  );
}
