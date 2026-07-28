import type { Document } from "@/types/document";

import DocumentCard from "./DocumentCard";


interface DocumentListProps {
    documents: Document[];
}


export default function DocumentList({
    documents,
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
                    />

                ))
            }

        </section>
    );
}
