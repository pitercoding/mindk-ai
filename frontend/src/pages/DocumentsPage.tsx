import { useEffect, useState } from "react";

import ConfirmDialog from "@/components/common/ConfirmDialog";
import EmptyState from "@/components/common/EmptyState";
import { AlertIcon } from "@/components/common/icons";
import { CardSkeletons } from "@/components/common/Skeleton";
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


    // Does not touch state before the first await, so it can run from an
    // effect (isLoading starts as true). Use reloadDocuments for manual reloads.
    async function loadDocuments() {

        try {

            const response = await getDocuments();

            setDocuments(response);
            setError(null);

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


    function reloadDocuments() {
        setIsLoading(true);
        setError(null);

        return loadDocuments();
    }


    async function handleSearch(
        value: string,
    ) {

        setSearch(value);

        if (!value.trim()) {

            reloadDocuments();

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

        // Deferred to a microtask: the loader sets state once its request
        // resolves, and react-hooks/set-state-in-effect disallows a direct call.
        Promise.resolve().then(loadDocuments);
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
                onUploadSuccess={reloadDocuments}
            />

            {
                isLoading ? (
                    <CardSkeletons label="Loading documents" />
                ) : error ? (
                    <EmptyState
                        tone="error"
                        icon={<AlertIcon size={20} />}
                        title={error}
                        description="Check your connection and try again."
                        actions={
                            <button
                                type="button"
                                className="btn btn-secondary"
                                onClick={reloadDocuments}
                            >
                                Try again
                            </button>
                        }
                    />
                ) : (
                    <DocumentList
                        documents={documents}
                        deletingId={deletingId}
                        onDelete={requestDelete}
                        search={search}
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
