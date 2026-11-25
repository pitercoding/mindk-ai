import { TrashIcon } from "@/components/common/icons";

import { formatDate, toPlainText } from "@/utils/text";

import type { Document } from "@/types/document";

interface DocumentCardProps {
    document: Document;
    isDeleting: boolean;
    onDelete: (id: number) => void;
}

export default function DocumentCard({
    document,
    isDeleting,
    onDelete,
}: DocumentCardProps) {

    const preview = toPlainText(document.content);
    const typeLabel = document.type.replace(/^\./, "").toUpperCase();
    const date = formatDate(document.created_at);

    return (
        <article className={
            isDeleting
                ? "document-card document-card-deleting"
                : "document-card"
        }>

            <header className="document-card-header">
                <h3>
                    {document.name}
                </h3>

                <button
                    className={
                        isDeleting
                            ? "btn btn-secondary btn-sm"
                            : "btn btn-secondary btn-icon btn-danger-hover"
                    }
                    type="button"
                    aria-label={`Delete ${document.name}`}
                    title="Delete"
                    disabled={isDeleting}
                    onClick={() => onDelete(document.id)}
                >
                    {isDeleting ? "Deleting..." : <TrashIcon size={16} />}
                </button>
            </header>

            {preview ? (
                <p className="card-preview">{preview}</p>
            ) : (
                <p className="card-preview card-preview-empty">
                    No content yet
                </p>
            )}

            <footer className="card-meta">
                <span className="document-card-type">{typeLabel}</span>

                {date && <time dateTime={document.created_at}>{date}</time>}
            </footer>

        </article>
    );
}