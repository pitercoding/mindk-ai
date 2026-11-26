import { useMemo, useState, type ReactNode } from "react";

import { Link } from "react-router-dom";

import EmptyState from "@/components/common/EmptyState";
import {
    AlertIcon,
    DocumentsIcon,
    NotesIcon,
    SearchIcon,
} from "@/components/common/icons";
import { RowSkeletons } from "@/components/common/Skeleton";

import { useSelectedKnowledge } from "@/context/selected-knowledge-context";

import type { Document } from "@/types/document";
import type { Note } from "@/types/note";

interface KnowledgeBaseProps {
    notes: Note[];
    documents: Document[];
    isLoading: boolean;
    error: string | null;
    onItemSelected?: () => void;
    onRetry?: () => void;
}

export default function KnowledgeBase({
    notes,
    documents,
    isLoading,
    error,
    onItemSelected,
    onRetry,
}: KnowledgeBaseProps) {

    const {
        selected,
        selectNote,
        selectDocument,
    } = useSelectedKnowledge();

    const [search, setSearch] = useState("");

    const query = search.trim().toLowerCase();

    const filteredNotes = useMemo(
        () => notes.filter((note) =>
            note.title.toLowerCase().includes(query)
        ),
        [notes, query],
    );

    const filteredDocuments = useMemo(
        () => documents.filter((document) =>
            document.name.toLowerCase().includes(query)
        ),
        [documents, query],
    );

    const isEmpty = notes.length === 0 && documents.length === 0;

    const hasNoResults =
        !isEmpty &&
        query !== "" &&
        filteredNotes.length === 0 &&
        filteredDocuments.length === 0;

    return (
        <section className="dashboard-panel knowledge-panel">

            <header className="panel-header">
                <div>
                    <h2>KNOWLEDGE BASE</h2>

                    <span>Recent</span>
                </div>
            </header>

            <div className="documents-search knowledge-search">
                <input
                    type="text"
                    placeholder="Search knowledge..."
                    value={search}
                    onChange={(event) => setSearch(event.target.value)}
                />
            </div>

            {isLoading ? (

                <RowSkeletons
                    count={5}
                    label="Loading knowledge base"
                />

            ) : error ? (

                <EmptyState
                    compact
                    tone="error"
                    icon={<AlertIcon size={18} />}
                    title={error}
                    actions={
                        onRetry && (
                            <button
                                type="button"
                                className="btn btn-secondary btn-sm"
                                onClick={onRetry}
                            >
                                Try again
                            </button>
                        )
                    }
                />

            ) : isEmpty ? (

                <EmptyState
                    compact
                    icon={<NotesIcon size={18} />}
                    title="Nothing here yet"
                    description="Add a note or upload a document to start building your knowledge base."
                    actions={
                        <>
                            <Link to="/notes" className="btn btn-primary btn-sm">
                                New note
                            </Link>

                            <Link to="/documents" className="btn btn-secondary btn-sm">
                                Upload
                            </Link>
                        </>
                    }
                />

            ) : hasNoResults ? (

                <EmptyState
                    compact
                    icon={<SearchIcon size={18} />}
                    title="No results"
                    description={`Nothing matches "${search.trim()}".`}
                />

            ) : (

                <div className="knowledge-list">

                    {filteredNotes.length > 0 && (
                        <div className="knowledge-group">
                            <h4 className="knowledge-group-title">
                                Notes
                            </h4>

                            {filteredNotes.map((note) => (
                                <KnowledgeItem
                                    key={`note-${note.id}`}
                                    icon={<NotesIcon size={16} />}
                                    label={note.title}
                                    isActive={
                                        selected?.type === "note" &&
                                        selected.item.id === note.id
                                    }
                                    onSelect={() => {
                                        selectNote(note);
                                        onItemSelected?.();
                                    }}
                                />
                            ))}
                        </div>
                    )}

                    {filteredDocuments.length > 0 && (
                        <div className="knowledge-group">
                            <h4 className="knowledge-group-title">
                                Documents
                            </h4>

                            {filteredDocuments.map((document) => (
                                <KnowledgeItem
                                    key={`document-${document.id}`}
                                    icon={<DocumentsIcon size={16} />}
                                    label={document.name}
                                    isActive={
                                        selected?.type === "document" &&
                                        selected.item.id === document.id
                                    }
                                    onSelect={() => {
                                        selectDocument(document);
                                        onItemSelected?.();
                                    }}
                                />
                            ))}
                        </div>
                    )}

                </div>
            )}

        </section>
    );
}

interface KnowledgeItemProps {
    icon: ReactNode;
    label: string;
    isActive: boolean;
    onSelect: () => void;
}

function KnowledgeItem({
    icon,
    label,
    isActive,
    onSelect,
}: KnowledgeItemProps) {

    return (
        <article
            className={`knowledge-item ${isActive ? "knowledge-item-active" : ""
                }`}
            role="button"
            tabIndex={0}
            onClick={onSelect}

            onKeyDown={(event) => {

                if (
                    event.key === "Enter" ||
                    event.key === " "
                ) {
                    event.preventDefault();
                    onSelect();
                }
            }}
        >

            <h3>
                <span className="knowledge-item-icon">{icon}</span>
                <span className="knowledge-item-label">{label}</span>
            </h3>

        </article>
    );
}
