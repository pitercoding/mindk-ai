import type { ReactNode, SVGProps } from "react";

/*
    App icons as inline SVG (24x24 grid, stroke based).
    They inherit the text color through currentColor, so hover and active
    states work without extra rules, and they look the same on every OS.
*/

export interface IconProps {
    size?: number;
}

function Icon({
    children,
    size = 18,
    ...props
}: SVGProps<SVGSVGElement> & IconProps & { children: ReactNode }) {
    return (
        <svg
            viewBox="0 0 24 24"
            width={size}
            height={size}
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

export function HomeIcon(props: IconProps) {
    return (
        <Icon {...props}>
            <path d="M3 11.5 12 4l9 7.5" />
            <path d="M5.5 10v9.5h13V10" />
            <path d="M10 19.5v-5h4v5" />
        </Icon>
    );
}

export function NotesIcon(props: IconProps) {
    return (
        <Icon {...props}>
            <rect x="4" y="3.5" width="16" height="17" rx="2" />
            <path d="M8 8.5h8M8 12.5h8M8 16.5h5" />
        </Icon>
    );
}

export function DocumentsIcon(props: IconProps) {
    return (
        <Icon {...props}>
            <path d="M14 3.5H7A2 2 0 0 0 5 5.5v13a2 2 0 0 0 2 2h10a2 2 0 0 0 2-2v-10z" />
            <path d="M14 3.5v5h5" />
        </Icon>
    );
}

export function LinksIcon(props: IconProps) {
    return (
        <Icon {...props}>
            <path d="M10 14a4 4 0 0 0 5.7 0l3-3a4 4 0 0 0-5.7-5.7l-1 1" />
            <path d="M14 10a4 4 0 0 0-5.7 0l-3 3a4 4 0 0 0 5.7 5.7l1-1" />
        </Icon>
    );
}

export function PdfIcon(props: IconProps) {
    return (
        <Icon {...props}>
            <path d="M14 3.5H7A2 2 0 0 0 5 5.5v13a2 2 0 0 0 2 2h10a2 2 0 0 0 2-2v-10z" />
            <path d="M14 3.5v5h5" />
            <path d="M8.5 16.5v-3h1.3a1.2 1.2 0 0 1 0 2.4H8.5" />
        </Icon>
    );
}

export function SettingsIcon(props: IconProps) {
    return (
        <Icon {...props}>
            <circle cx="12" cy="12" r="3" />
            <path d="M12 3.5v2.2M12 18.3v2.2M3.5 12h2.2M18.3 12h2.2M6 6l1.6 1.6M16.4 16.4 18 18M18 6l-1.6 1.6M7.6 16.4 6 18" />
        </Icon>
    );
}

export function EditIcon(props: IconProps) {
    return (
        <Icon {...props}>
            <path d="M4 20h4L19.5 8.5a2.1 2.1 0 0 0-3-3L5 17z" />
            <path d="m14.5 7.5 3 3" />
        </Icon>
    );
}

export function TrashIcon(props: IconProps) {
    return (
        <Icon {...props}>
            <path d="M4 7h16" />
            <path d="M9.5 7V4.5h5V7" />
            <path d="M6.5 7l.8 12.2a1.5 1.5 0 0 0 1.5 1.3h6.4a1.5 1.5 0 0 0 1.5-1.3L17.5 7" />
            <path d="M10 11v5M14 11v5" />
        </Icon>
    );
}

export function UploadIcon(props: IconProps) {
    return (
        <Icon {...props}>
            <path d="M12 16V4.5" />
            <path d="m7.5 9 4.5-4.5L16.5 9" />
            <path d="M4.5 15.5v3a1.5 1.5 0 0 0 1.5 1.5h12a1.5 1.5 0 0 0 1.5-1.5v-3" />
        </Icon>
    );
}

export function CloseIcon(props: IconProps) {
    return (
        <Icon {...props}>
            <path d="M6 6l12 12M18 6 6 18" />
        </Icon>
    );
}

export function AlertIcon(props: IconProps) {
    return (
        <Icon {...props}>
            <path d="M12 4 2.8 19.5h18.4z" />
            <path d="M12 10v4.5M12 17.4v.1" />
        </Icon>
    );
}

export function CheckIcon(props: IconProps) {
    return (
        <Icon {...props}>
            <circle cx="12" cy="12" r="9" />
            <path d="m8 12.5 2.8 2.8L16.2 9.5" />
        </Icon>
    );
}

export function InfoIcon(props: IconProps) {
    return (
        <Icon {...props}>
            <circle cx="12" cy="12" r="9" />
            <path d="M12 11v5.2M12 7.8v.1" />
        </Icon>
    );
}
