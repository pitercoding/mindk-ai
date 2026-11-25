import { TrashIcon } from "@/components/common/icons";

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
                            ? "document-action-button document-action-button-busy"
                            : "document-action-button document-action-button-danger"
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

            <span>
                {document.type}
            </span>

            <p>
                {document.content.slice(0, 120)}
                {document.content.length > 120 && "..."}
            </p>

        </article>
    );
}