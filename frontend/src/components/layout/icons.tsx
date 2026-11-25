import type { ReactNode, SVGProps } from "react";

/*
    Sidebar icons as inline SVG (24x24 grid, stroke based).
    They inherit the text color through currentColor, so hover and active
    states work without extra rules, and they look the same on every OS.
*/

function Icon({ children, ...props }: SVGProps<SVGSVGElement> & { children: ReactNode }) {
    return (
        <svg
            viewBox="0 0 24 24"
            width="18"
            height="18"
            fill="none"
            stroke="currentColor"
            strokeWidth="1.8"
            strokeLinecap="round"
            strokeLinejoin="round"
            aria-hidden="true"
            focusable="false"
            {...props}
        >
            {children}
        </svg>
    );
}

export function HomeIcon() {
    return (
        <Icon>
            <path d="M3 11.5 12 4l9 7.5" />
            <path d="M5.5 10v9.5h13V10" />
            <path d="M10 19.5v-5h4v5" />
        </Icon>
    );
}

export function NotesIcon() {
    return (
        <Icon>
            <rect x="4" y="3.5" width="16" height="17" rx="2" />
            <path d="M8 8.5h8M8 12.5h8M8 16.5h5" />
        </Icon>
    );
}

export function DocumentsIcon() {
    return (
        <Icon>
            <path d="M14 3.5H7A2 2 0 0 0 5 5.5v13a2 2 0 0 0 2 2h10a2 2 0 0 0 2-2v-10z" />
            <path d="M14 3.5v5h5" />
        </Icon>
    );
}

export function LinksIcon() {
    return (
        <Icon>
            <path d="M10 14a4 4 0 0 0 5.7 0l3-3a4 4 0 0 0-5.7-5.7l-1 1" />
            <path d="M14 10a4 4 0 0 0-5.7 0l-3 3a4 4 0 0 0 5.7 5.7l1-1" />
        </Icon>
    );
}

export function PdfIcon() {
    return (
        <Icon>
            <path d="M14 3.5H7A2 2 0 0 0 5 5.5v13a2 2 0 0 0 2 2h10a2 2 0 0 0 2-2v-10z" />
            <path d="M14 3.5v5h5" />
            <path d="M8.5 16.5v-3h1.3a1.2 1.2 0 0 1 0 2.4H8.5" />
        </Icon>
    );
}

export function SettingsIcon() {
    return (
        <Icon>
            <circle cx="12" cy="12" r="3" />
            <path d="M12 3.5v2.2M12 18.3v2.2M3.5 12h2.2M18.3 12h2.2M6 6l1.6 1.6M16.4 16.4 18 18M18 6l-1.6 1.6M7.6 16.4 6 18" />
        </Icon>
    );
}
