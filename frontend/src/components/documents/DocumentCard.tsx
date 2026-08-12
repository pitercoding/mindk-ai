import type { Document } from "@/types/document";

interface DocumentCardProps {
    document: Document;
    onDelete: (id: number) => void;
}

export default function DocumentCard({
    document,
    onDelete,
}: DocumentCardProps) {

    return (
        <article className="document-card">

            <header className="document-card-header">
                <h3>
                    {document.name}
                </h3>

                <button
                    className="document-action-button"
                    type="button"
                    onClick={() => onDelete(document.id)}
                >
                    🗑️
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