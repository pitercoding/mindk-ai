import { useEffect, useState } from "react";

import DocumentList from "@/components/documents/DocumentList";
import DocumentUpload from "@/components/documents/DocumentUpload";

import {
    deleteDocument,
    getDocuments,
    searchDocuments,
} from "@/services/documentService";

import type { Document } from "@/types/document";


export default function DocumentsPage() {

    const [documents, setDocuments] = useState<Document[]>([]);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [actionError, setActionError] = useState<string | null>(null);
    const [deletingId, setDeletingId] = useState<number | null>(null);
    const [search, setSearch] = useState("");


    async function loadDocuments() {
        setIsLoading(true);
        setError(null);
        setActionError(null);

        try {

            const response = await getDocuments();

            setDocuments(response);

        } catch (error) {

            console.error(
                "Failed to load documents:",
                error,
            );

            setError("Failed to load documents.");

        } finally {

            setIsLoading(false);
        }
    }


    async function handleSearch(
        value: string,
    ) {

        setSearch(value);
        setActionError(null);

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

            setActionError("Failed to search documents.");

        }
    }


    async function handleDelete(id: number) {

        if (!confirm("Delete this document?")) {
            return;
        }

        setActionError(null);
        setDeletingId(id);

        try {

            await deleteDocument(id);

            setDocuments((previous) =>
                previous.filter((document) => document.id !== id)
            );

        } catch (error) {

            console.error(
                "Failed to delete document:",
                error,
            );

            setActionError("Failed to delete document.");

        } finally {
            setDeletingId(null);
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

            {actionError && (
                <p className="documents-error">
                    {actionError}
                </p>
            )}

            {
                isLoading ? (
                    <p>
                        Loading documents...
                    </p>
                ) : error ? (
                    <p className="documents-error">
                        {error}
                    </p>
                ) : (
                    <DocumentList
                        documents={documents}
                        deletingId={deletingId}
                        onDelete={handleDelete}
                    />
                )
            }

        </section>
    );
}
