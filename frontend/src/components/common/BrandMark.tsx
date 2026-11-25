/*
    MindK AI monogram: a "K" built as a small network (a stem branching
    into two nodes). Uses currentColor for the strokes so it follows the
    surrounding color (accent by default).
*/

export default function BrandMark({ size = 28 }: { size?: number }) {
    return (
        <svg
            className="brand-mark"
            viewBox="0 0 32 32"
            width={size}
            height={size}
            fill="none"
            aria-hidden="true"
            focusable="false"
        >
            <rect
                x="1"
                y="1"
                width="30"
                height="30"
                rx="9"
                fill="currentColor"
                fillOpacity="0.12"
                stroke="currentColor"
                strokeOpacity="0.55"
                strokeWidth="1.6"
            />
            <path
                d="M11 8v16M11 16l9.5-7M11 16l9.5 7"
                stroke="currentColor"
                strokeWidth="2.2"
                strokeLinecap="round"
                strokeLinejoin="round"
            />
            <circle cx="21.5" cy="8.2" r="2.2" fill="currentColor" />
            <circle cx="21.5" cy="23.8" r="2.2" fill="currentColor" />
        </svg>
    );
}
