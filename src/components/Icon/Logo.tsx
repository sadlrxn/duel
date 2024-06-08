import React from 'react';

export default function Logo({
  width = 93,
  height = 22
}: {
  width?: number;
  height?: number;
}) {
  return (
    <svg
      width={width}
      height={height}
      viewBox="0 0 93 22"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
    >
      <g clipPath="url(#clip0_2821_35135)">
        <path
          fillRule="evenodd"
          clipRule="evenodd"
          d="M32.4787 3.18775L28.7139 21.9893H47.5381L51.9386 0.0131226H46.3896L43.1185 3.1904L40.625 15.6347H38.0851L41.2988 0.0131226H35.7416L32.4787 3.18775Z"
          fill="white"
        />
        <path
          fillRule="evenodd"
          clipRule="evenodd"
          d="M54.5768 3.17722L50.8065 21.9999H69.6308L70.9021 15.6454L60.2677 15.6242L60.7451 13.2412H67.6638L68.5668 8.7401L61.6236 8.72686L62.1229 6.3545H70.8802L72.1543 -6.10352e-05L57.8478 0.0025867L54.5768 3.17722Z"
          fill="white"
        />
        <path
          fillRule="evenodd"
          clipRule="evenodd"
          d="M76.6747 3.17722L72.9044 21.9999H91.7287L93 15.6454L82.3657 15.6242L85.4949 -6.10352e-05H79.954L76.6747 3.17722Z"
          fill="white"
        />
        <path
          fillRule="evenodd"
          clipRule="evenodd"
          d="M4.40051 0.0131226H25.8274L29.093 3.1904L27.8408 9.41522L14.4538 21.9893H0L2.65177 8.75064H10.8362L9.20752 16.8818L20.393 6.37562L3.12919 6.36768L4.40051 0.0131226Z"
          fill="white"
        />
      </g>
      <defs>
        <clipPath id="clip0_2821_35135">
          <rect width="93" height="22" fill="white" />
        </clipPath>
      </defs>
    </svg>
  );
}
