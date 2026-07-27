import type { Document } from "@/types/document";

interface DocumentCardProps {
    document: Document;
}

export default function DocumentCard({
    document,
}: DocumentCardProps) {

    return (
        <article className="document-card">

            <header>
                <h3>
                    {document.name}
                </h3>

                <span>
                    {document.type}
                </span>
            </header>

            <p>
                {document.content.slice(0, 120)}
                {document.content.length > 120 && "..."}
            </p>

        </article>
    );
}