import { useEffect, useState } from "react";

import type { Document } from "@/types/document";

import { getDocuments } from "@/services/documentService";

import DocumentCard from "./DocumentCard";


export default function DocumentList() {

    const [documents, setDocuments] = useState<Document[]>([]);
    const [isLoading, setIsLoading] = useState(true);

    async function loadDocuments() {

        try {

            const response = await getDocuments();

            setDocuments(response);

        } catch (error) {

            console.error("Failed to load documents:", error,);

        } finally {

            setIsLoading(false);
        }
    }


    useEffect(() => {

        loadDocuments();

    }, []);


    if (isLoading) {

        return (
            <p>Loading documents...</p>
        );
    }


    return (

        <section className="document-list">

            {
                documents.length === 0
                    ? (
                        <p>No documents yet.</p>
                    )
                    
                    : documents.map((document) => (

                        <DocumentCard
                            key={document.id}
                            document={document}
                        />

                    ))
            }

        </section>
    );
}
