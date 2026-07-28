import { useEffect, useState } from "react";

import DocumentList from "@/components/documents/DocumentList";
import DocumentUpload from "@/components/documents/DocumentUpload";

import { getDocuments } from "@/services/documentService";

import type { Document } from "@/types/document";


export default function DocumentsPage() {

    const [documents, setDocuments] = useState<Document[]>([]);
    const [isLoading, setIsLoading] = useState(true);


    async function loadDocuments() {

        try {

            const response = await getDocuments();

            setDocuments(response);

        } catch (error) {

            console.error(
                "Failed to load documents:",
                error,
            );

        } finally {

            setIsLoading(false);
        }
    }


    useEffect(() => {

        loadDocuments();

    }, []);


    return (
        <section className="documents-page">

            <header className="page-header">

                <div>
                    <h1>Documents</h1>

                    <p>
                        Manage your knowledge documents.
                    </p>
                </div>

            </header>


            <DocumentUpload
                onUploadSuccess={loadDocuments}
            />


            {
                isLoading
                    ? (
                        <p>Loading documents...</p>
                    )
                    : (
                        <DocumentList
                            documents={documents}
                        />
                    )
            }

        </section>
    );
}
