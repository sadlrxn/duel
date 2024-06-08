import React from "react";

export default function Dot({ color }: { color: string }) {
  return (
    <svg
      width="5"
      height="6"
      viewBox="0 0 5 6"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
    >
      <circle cx="2.5" cy="2.98975" r="2.5" fill={color} />
    </svg>
  );
}
