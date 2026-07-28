import { useEffect, useState } from "react";

import DocumentList from "@/components/documents/DocumentList";
import DocumentUpload from "@/components/documents/DocumentUpload";

import {
    getDocuments,
    searchDocuments,
} from "@/services/documentService";

import type { Document } from "@/types/document";


export default function DocumentsPage() {

    const [documents, setDocuments] = useState<Document[]>([]);
    const [isLoading, setIsLoading] = useState(true);
    const [search, setSearch] = useState("");


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


    async function handleSearch(
        value: string,
    ) {

        setSearch(value);

        if (!value.trim()) {

            loadDocuments();

            return;
        }


        try {

            const response =
                await searchDocuments(value);

            setDocuments(response);

        } catch (error) {

            console.error(
                "Failed to search documents:",
                error,
            );

        }
    }


    useEffect(() => {

        loadDocuments();

    }, []);


    return (
        <section className="documents-page">

            <header className="page-header">

                <div>
                    <h1>
                        Documents
                    </h1>

                    <p>
                        Manage your knowledge documents.
                    </p>
                </div>

            </header>


            <div className="documents-search">

                <input
                    type="text"
                    placeholder="Search documents..."
                    value={search}
                    onChange={(event) =>
                        handleSearch(event.target.value)
                    }
                />

            </div>


            <DocumentUpload
                onUploadSuccess={loadDocuments}
            />


            {
                isLoading
                    ? (
                        <p>
                            Loading documents...
                        </p>
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
