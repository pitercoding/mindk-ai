import EmptyState from "@/components/common/EmptyState";
import { DocumentsIcon, SearchIcon, UploadIcon } from "@/components/common/icons";

import type { Document } from "@/types/document";

import DocumentCard from "./DocumentCard";


interface DocumentListProps {
    documents: Document[];
    deletingId: number | null;
    onDelete: (id: number) => void;
    search?: string;
}


export default function DocumentList({
    documents,
    deletingId,
    onDelete,
    search = "",
}: DocumentListProps) {


    if (documents.length === 0) {

        if (search.trim()) {

            return (
                <EmptyState
                    icon={<SearchIcon size={20} />}
                    title="No documents found"
                    description={`Nothing matches "${search.trim()}". Try a different search term.`}
                />
            );
        }

        return (
            <EmptyState
                icon={<DocumentsIcon size={20} />}
                title="No documents yet"
                description="Upload Markdown or text files and ask questions about them in the chat."
                actions={
                    <button
                        type="button"
                        className="btn btn-primary"
                        onClick={() =>
                            window.document
                                .querySelector<HTMLInputElement>(
                                    ".document-upload input[type=file]",
                                )
                                ?.click()
                        }
                    >
                        <UploadIcon size={16} />
                        Choose a file
                    </button>
                }
            />
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
