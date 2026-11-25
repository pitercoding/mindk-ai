import { useEffect, useState } from "react";

import ConfirmDialog from "@/components/common/ConfirmDialog";
import DocumentList from "@/components/documents/DocumentList";
import DocumentUpload from "@/components/documents/DocumentUpload";

import { useToast } from "@/context/toast-context";

import {
    deleteDocument,
    getDocuments,
    searchDocuments,
} from "@/services/documentService";

import type { Document } from "@/types/document";


export default function DocumentsPage() {

    const toast = useToast();

    const [documents, setDocuments] = useState<Document[]>([]);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [deletingId, setDeletingId] = useState<number | null>(null);
    const [documentToDelete, setDocumentToDelete] = useState<Document | null>(null);
    const [search, setSearch] = useState("");


    async function loadDocuments() {
        setIsLoading(true);
        setError(null);

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

            toast.error("Failed to search documents.");

        }
    }


    function requestDelete(id: number) {
        setDocumentToDelete(
            documents.find((document) => document.id === id) ?? null,
        );
    }


    async function handleDelete(target: Document) {

        setDocumentToDelete(null);
        setDeletingId(target.id);

        try {

            await deleteDocument(target.id);

            setDocuments((previous) =>
                previous.filter((document) => document.id !== target.id)
            );

            toast.success("Document deleted.");

        } catch (error) {

            console.error(
                "Failed to delete document:",
                error,
            );

            toast.error("Failed to delete document.");

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
                        onDelete={requestDelete}
                    />
                )
            }

            {documentToDelete && (
                <ConfirmDialog
                    title="Delete document?"
                    message={`"${documentToDelete.name}" will be permanently deleted.`}
                    onConfirm={() => handleDelete(documentToDelete)}
                    onCancel={() => setDocumentToDelete(null)}
                />
            )}

        </section>
    );
}
