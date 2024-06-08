import React from 'react';

export default function ExpandIcon({
  collapse = true
}: {
  collapse?: boolean;
}) {
  return (
    <svg
      width="9"
      height="7"
      viewBox="0 0 9 7"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
      style={{
        transform: collapse ? `rotate(180deg)` : `none`
      }}
    >
      <path
        d="M5.26799 6.0791C4.86821 6.55848 4.13179 6.55848 3.73201 6.0791L0.238901 1.89047C-0.304217 1.2392 0.15888 0.25 1.00689 0.25H7.99311C8.84112 0.25 9.30422 1.2392 8.7611 1.89046L5.26799 6.0791Z"
        fill="currentColor"
      />
    </svg>
  );
}
