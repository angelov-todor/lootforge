"use client";

interface DieIconProps {
  size?: number;
  color?: string;
}

export function DieIcon({ size = 24, color = "currentColor" }: DieIconProps) {
  return (
    <svg
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
    >
      <path
        d="M12 2L2 8.5V15.5L12 22L22 15.5V8.5L12 2Z"
        stroke={color}
        strokeWidth="1.5"
        strokeLinejoin="round"
      />
      <path
        d="M12 2L12 22"
        stroke={color}
        strokeWidth="1.5"
      />
      <path
        d="M2 8.5L12 12L22 8.5"
        stroke={color}
        strokeWidth="1.5"
      />
      <circle cx="8" cy="14" r="0.8" fill={color} />
      <circle cx="12" cy="7" r="0.8" fill={color} />
      <circle cx="16" cy="14" r="0.8" fill={color} />
    </svg>
  );
}
