/**
 * Converts markdown into plain, single-line text for card previews,
 * so symbols such as "#", "- [x]" and "|" never leak into the UI.
 */
export function toPlainText(markdown: string): string {

    return markdown
        // fenced code blocks (content is noise in a preview)
        .replace(/```[\s\S]*?```/g, " ")
        // table separator rows: |---|:---:|
        .replace(/^\s*\|?[\s:|-]*-[\s:|-]*\|?\s*$/gm, " ")
        // images and links -> their visible text
        .replace(/!\[([^\]]*)\]\([^)]*\)/g, "$1")
        .replace(/\[([^\]]+)\]\([^)]*\)/g, "$1")
        // headings, blockquotes, task boxes and list markers
        .replace(/^\s{0,3}#{1,6}\s+/gm, "")
        .replace(/^\s*>\s?/gm, "")
        .replace(/^\s*[-*+]\s+\[[ xX]\]\s+/gm, "")
        .replace(/^\s*(?:[-*+]|\d+\.)\s+/gm, "")
        // horizontal rules
        .replace(/^\s*([-*_])(?:\s*\1){2,}\s*$/gm, " ")
        // table pipes
        .replace(/\|/g, " ")
        // emphasis and inline code (underscores are left alone on purpose,
        // so snake_case names are not mangled)
        .replace(/(\*\*|~~)(.+?)\1/g, "$2")
        .replace(/\*(.+?)\*/g, "$1")
        .replace(/`([^`]*)`/g, "$1")
        // collapse whitespace
        .replace(/\s+/g, " ")
        .trim();
}

export function countWords(text: string): number {

    const trimmed = text.trim();

    return trimmed ? trimmed.split(/\s+/).length : 0;
}

export function formatDate(value: string): string {

    const date = new Date(value);

    if (Number.isNaN(date.getTime())) {
        return "";
    }

    return date.toLocaleDateString(undefined, {
        day: "2-digit",
        month: "short",
        year: "numeric",
    });
}

function escapeRegExp(value: string): string {
    return value.replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
}

/** Removes a leading "# Title" heading when it only repeats the note title. */
export function stripTitleHeading(content: string, title: string): string {
    return content.replace(
        new RegExp(`^#\\s+${escapeRegExp(title)}\\s*\\n+`),
        "",
    );
}
