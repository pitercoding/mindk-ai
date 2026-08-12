import { useMemo, useState } from "react";

import { useSelectedKnowledge } from "@/context/SelectedKnowledgeContext";

import type { Document } from "@/types/document";
import type { Note } from "@/types/note";

interface KnowledgeBaseProps {
    notes: Note[];
    documents: Document[];
    isLoading: boolean;
    error: string | null;
}

export default function KnowledgeBase({
    notes,
    documents,
    isLoading,
    error,
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

                <button>...</button>
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

                <p className="knowledge-status">
                    Loading knowledge base...
                </p>

            ) : error ? (

                <p className="knowledge-status knowledge-status-error">
                    {error}
                </p>

            ) : isEmpty ? (

                <p className="knowledge-status">
                    No notes or documents yet.
                </p>

            ) : hasNoResults ? (

                <p className="knowledge-status">
                    No results for "{search}".
                </p>

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
                                    icon="📄"
                                    label={note.title}
                                    isActive={
                                        selected?.type === "note" &&
                                        selected.item.id === note.id
                                    }
                                    onSelect={() => selectNote(note)}
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
                                    icon="📑"
                                    label={document.name}
                                    isActive={
                                        selected?.type === "document" &&
                                        selected.item.id === document.id
                                    }
                                    onSelect={() => selectDocument(document)}
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
    icon: string;
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

            <h3>{icon} {label}</h3>

        </article>
    );
}
