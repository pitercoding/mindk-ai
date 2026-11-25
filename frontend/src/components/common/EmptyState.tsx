import type { ReactNode } from "react";

interface EmptyStateProps {
    icon: ReactNode;
    title: string;
    description?: string;
    /** Buttons or links rendered under the text. */
    actions?: ReactNode;
    tone?: "default" | "error";
    /** Smaller paddings, for use inside narrow panels. */
    compact?: boolean;
}

export default function EmptyState({
    icon,
    title,
    description,
    actions,
    tone = "default",
    compact = false,
}: EmptyStateProps) {

    const className = [
        "empty-state",
        tone === "error" ? "empty-state-error" : "",
        compact ? "empty-state-compact" : "",
    ].filter(Boolean).join(" ");

    return (
        <div
            className={className}
            role={tone === "error" ? "alert" : undefined}
        >

            <span className="empty-state-icon">
                {icon}
            </span>

            <h3 className="empty-state-title">
                {title}
            </h3>

            {description && (
                <p className="empty-state-description">
                    {description}
                </p>
            )}

            {actions && (
                <div className="empty-state-actions">
                    {actions}
                </div>
            )}

        </div>
    );
}
