import type { Document } from "@/types/document";

import DocumentCard from "./DocumentCard";


interface DocumentListProps {
    documents: Document[];
    deletingId: number | null;
    onDelete: (id: number) => void;
}


export default function DocumentList({
    documents,
    deletingId,
    onDelete,
}: DocumentListProps) {


    if (documents.length === 0) {

        return (
            <p>
                No documents yet.
            </p>
        );
    }


    return (

        <section className="document-list">

            {
                documents.map((document) => (

                    <DocumentCard
                        key={document.id}
                        document={document}
                        isDeleting={deletingId === document.id}
                        onDelete={onDelete}
                    />

                ))
            }

        </section>
    );
}
