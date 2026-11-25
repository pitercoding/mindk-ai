import type { ReactNode } from "react";

/*
    Loading placeholders. Each group is announced once to screen readers
    ("Loading ...") while the decorative blocks stay hidden from them.
*/

interface SkeletonGroupProps {
    label: string;
    children: ReactNode;
    className?: string;
}

function SkeletonGroup({ label, children, className }: SkeletonGroupProps) {
    return (
        <div
            className={className}
            role="status"
            aria-live="polite"
            aria-busy="true"
        >
            <span className="sr-only">{label}</span>
            {children}
        </div>
    );
}

function Block({ className = "" }: { className?: string }) {
    return <span className={`skeleton ${className}`} aria-hidden="true" />;
}

/** Placeholder cards that mimic note and document cards. */
export function CardSkeletons({
    count = 3,
    label = "Loading",
}: {
    count?: number;
    label?: string;
}) {
    return (
        <SkeletonGroup label={label} className="skeleton-cards">
            {Array.from({ length: count }, (_, index) => (
                <div className="skeleton-card" key={index} aria-hidden="true">
                    <Block className="skeleton-title" />
                    <Block className="skeleton-line" />
                    <Block className="skeleton-line skeleton-line-short" />
                </div>
            ))}
        </SkeletonGroup>
    );
}

/** Placeholder rows that mimic list items (knowledge base, conversations). */
export function RowSkeletons({
    count = 4,
    label = "Loading",
    className = "",
}: {
    count?: number;
    label?: string;
    className?: string;
}) {
    return (
        <SkeletonGroup label={label} className={`skeleton-rows ${className}`}>
            {Array.from({ length: count }, (_, index) => (
                <Block className="skeleton-row" key={index} />
            ))}
        </SkeletonGroup>
    );
}
